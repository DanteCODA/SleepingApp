package assets

import (
	"context"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/usecase/checkpoint"
)

// Service sector
type Service struct {
	assetRepo         Repo
	checkpointService checkpoint.Service
	log               logger.ContextLog
}

// NewService create new service
func NewService(assetRepo Repo, checkpointSer