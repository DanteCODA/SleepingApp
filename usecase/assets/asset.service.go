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
	s.log.Info(ctx, "getting assets from checkpoint")
	numAssets, err := s.assetRepo.CountAssets(ctx)
	if err != nil {
		s.log.Error(ctx, "count assets failed", "error", err)
	}

	checkpoint, err := s.checkpointService.UpdateCheckpoint(ctx, pageSize, numAssets)
	if err != nil {
		s.log.Error(ctx, "find and update checkpoint failed", "error", err)
	}

	if checkpoint == nil {
		s.log.Error(ctx, "checkpoint is nil", "checkpoint", checkpoint)
		return nil, nil
	}

	return s.assetRepo.FindAssetsFromCheckpoint(ctx, checkpoint)
}
