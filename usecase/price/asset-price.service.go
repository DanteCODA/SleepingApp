package price

import (
	"context"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
)

// Service sector
type Service struct {
	assetPriceRepo Repo
	log            logger.ContextLog
}

// NewService create new service
func NewService(assetPriceRepo Repo, log logger.ContextLog) *Service {
	return &Service{
		assetPriceRepo: assetPriceRepo,
		log:            log,
	}
}

// AddAssetPrice creates new asset price
func (s *Service) AddAssetPrice(ctx context.Context, assetPrice *entities.AssetPrice) error {
	s.log.Info(ctx, "adding asset price", "ticker", assetPrice.Ticker)
	return s.assetPriceRepo.InsertAssetPrice(ctx, assetPrice)
}
