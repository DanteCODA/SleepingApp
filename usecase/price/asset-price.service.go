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
func