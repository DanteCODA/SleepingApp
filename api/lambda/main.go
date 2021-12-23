package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/config"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/consts"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/infrastructure/repositories/repos"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/infrastructure/scraper"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/assets"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/checkpoint"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/price"
)

func main() {
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context) {
	log.Println("lambda handler is called")

	appConf := config.AppConf

	// create new logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer zap.Close()

	// create new repository
	assetPriceRepo, err := repos.NewAssetPriceMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create asset price mongo failed")
	}
	defer assetPriceRepo.Close()

	// create new repository
	assetRepo, err := repos.NewAssetMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create asset mongo failed")
	}
	defer assetRepo.Close()

	// create new repository
	checkpointRepo, err := repos.NewCheckpointMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create checkpoint mongo failed")
	}
	defer checkpointRepo.Close()

	// create new services
	checkpointService := checkpoint.NewService(checkpointRepo, zap)
	assetService := assets.NewService(assetRepo, *checkpointService, zap)
	priceService := price.NewService(assetPriceRepo, zap)

	// create new scraper jobs
	job := scraper.NewAssetPriceScraper(assetService, priceService, zap)
	job.ScrapeAssetPricesFromCheckpoint(consts.PAGE_SIZE)
	defer job.Close()
}
