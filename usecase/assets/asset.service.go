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
func NewService(assetRepo Repo, checkpointService checkpoint.Service, log logger.ContextLog) *Service {
	return &Service{
		assetRepo:         assetRepo,
		checkpointService: checkpointService,
		log:               log,
	}
}

// GetAllAssets gets all assets
func (s *Service) GetAllAssets(ctx context.Context) ([]*entities.Asset, error) {
	s.log.Info(ctx, "getting all assets")
	return s.assetRepo.FindAllAssets(ctx)
}

// GetAssetsFromCheckpoint gets all assets from checkpoint
func (s *Service) GetAssetsFromCheckpoint(ctx context.Context, pageSize int64) ([]*entities.Asset, error) {
	s.log.Info(ctx, "get