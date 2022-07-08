
package scraper

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/google/uuid"
	corid "github.com/lenoobz/aws-lambda-corid"
	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/config"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/assets"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/price"
)

// PriceScraper struct
type PriceScraper struct {
	ScrapePriceJob *colly.Collector
	priceService   *price.Service
	assetService   *assets.Service
	log            logger.ContextLog
	errorTickers   []string
}

// NewAssetPriceScraper create new price scraper
func NewAssetPriceScraper(assetService *assets.Service, priceService *price.Service, log logger.ContextLog) *PriceScraper {
	scrapePriceJob := newScraperJob()

	return &PriceScraper{
		ScrapePriceJob: scrapePriceJob,
		assetService:   assetService,
		priceService:   priceService,
		log:            log,
	}
}

// newScraperJob creates a new colly collector with some custom configs
func newScraperJob() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(config.AllowDomain),
		colly.Async(true),
	)

	// Overrides the default timeout (10 seconds) for this collector
	c.SetRequestTimeout(30 * time.Second)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  config.DomainGlob,
		Parallelism: 2,
		RandomDelay: 2 * time.Second,
	})

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	return c
}

// configJobs configs on error handler and on response handler for scaper jobs
func (s *PriceScraper) configJobs() {
	s.ScrapePriceJob.OnError(s.errorHandler)
	s.ScrapePriceJob.OnScraped(s.scrapedHandler)
	s.ScrapePriceJob.OnHTML("div[id=quote-header-info]", s.processPriceResponse)
}

// ScrapeAllAssetPrices scrape all assets price
func (s *PriceScraper) ScrapeAllAssetPrices() {
	ctx := context.Background()

	s.configJobs()

	assets, err := s.assetService.GetAllAssets(ctx)
	if err != nil {
		s.log.Error(ctx, "get assets list failed", "error", err)
	}

	for _, asset := range assets {
		reqContext := colly.NewContext()
		reqContext.Put("ticker", asset.Ticker)
		reqContext.Put("currency", asset.Currency)

		url := config.GetPriceByTickerURL(asset.Ticker)

		s.log.Info(ctx, "scraping asset price", "ticker", asset.Ticker)
		if err := s.ScrapePriceJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scraping asset price", "error", err, "ticker", asset.Ticker)
		}
	}

	s.ScrapePriceJob.Wait()
}

// ScrapeAssetPricesFromCheckpoint scrape all assets price from checkpoint
func (s *PriceScraper) ScrapeAssetPricesFromCheckpoint(pageSize int64) {
	ctx := context.Background()

	s.configJobs()

	assets, err := s.assetService.GetAssetsFromCheckpoint(ctx, pageSize)
	if err != nil {
		s.log.Error(ctx, "get assets list failed", "error", err)
	}

	for _, asset := range assets {
		reqContext := colly.NewContext()
		reqContext.Put("ticker", asset.Ticker)
		reqContext.Put("currency", asset.Currency)

		url := config.GetPriceByTickerURL(asset.Ticker)

		s.log.Info(ctx, "scraping asset price", "ticker", asset.Ticker)
		if err := s.ScrapePriceJob.Request("GET", url, nil, reqContext, nil); err != nil {
			s.log.Error(ctx, "scraping asset price failed", "error", err, "ticker", asset.Ticker)
		}
	}

	s.ScrapePriceJob.Wait()
}

///////////////////////////////////////////////////////////
// Scraper Handler
///////////////////////////////////////////////////////////

// errorHandler generic error handler for all scaper jobs
func (s *PriceScraper) errorHandler(r *colly.Response, err error) {
	ctx := context.Background()
	s.log.Error(ctx, "failed to request url", "url", r.Request.URL, "error", err)
	s.errorTickers = append(s.errorTickers, r.Request.Ctx.Get("ticker"))
}

func (s *PriceScraper) scrapedHandler(r *colly.Response) {
	ctx := context.Background()
	foundPrice := r.Ctx.Get("foundPrice")
	if foundPrice == "" {
		s.log.Error(ctx, "price not found", "ticker", r.Request.Ctx.Get("ticker"))
		s.errorTickers = append(s.errorTickers, r.Request.Ctx.Get("ticker"))
	}
}

func (s *PriceScraper) processPriceResponse(e *colly.HTMLElement) {
	// create correlation if for processing fund list
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)

	ticker := e.Request.Ctx.Get("ticker")
	currency := e.Request.Ctx.Get("currency")
	s.log.Info(ctx, "processPriceResponse", "ticker", ticker)

	foundPrice := false

	assetPrice := entities.AssetPrice{
		Ticker:   ticker,
		Currency: currency,
	}

	e.ForEach("span", func(_ int, span *colly.HTMLElement) {
		txt := span.Attr("data-reactid")
		if strings.EqualFold(txt, "31") {
			p := strings.Replace(span.DOM.Text(), ",", "", -1)

			val, err := strconv.ParseFloat(p, 64)
			if err != nil {
				s.log.Error(ctx, "parse price failed", "error", err, "ticker", ticker, "raw-value", txt)
				return
			}

			assetPrice.Price = val
			foundPrice = true
		}
	})

	if foundPrice {
		e.Response.Ctx.Put("foundPrice", "true")

		if err := s.priceService.AddAssetPrice(ctx, &assetPrice); err != nil {
			s.log.Error(ctx, "add price failed", "error", err, "ticker", ticker)
		}
	}