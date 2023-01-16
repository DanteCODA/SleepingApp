package checkpoint

import (
	"context"

	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
)

///////////////////////////////////////////////////////////
// Asset Price Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
}

// Writer interface
type Writer interface {
	UpdateCheckpoint(ctx context.Context, pageSize int64, numAssets int64) (*entities.Checkpoint, error)
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
