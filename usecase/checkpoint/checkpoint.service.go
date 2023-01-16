
package checkpoint

import (
	"context"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
)

// Service sector
type Service struct {
	checkpointRepo Repo
	log            logger.ContextLog
}

// NewService create new service
func NewService(checkpointRepo Repo, log logger.ContextLog) *Service {
	return &Service{
		checkpointRepo: checkpointRepo,
		log:            log,
	}
}

// GetAllAssets gets all assets
func (s *Service) UpdateCheckpoint(ctx context.Context, pageSize int64, numAssets int64) (*entities.Checkpoint, error) {
	s.log.Info(ctx, "updating checkpoint")
	return s.checkpointRepo.UpdateCheckpoint(ctx, pageSize, numAssets)
}