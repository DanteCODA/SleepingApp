package assets

import (
	"context"

	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
)

///////////////////////////////////////////////////////////
// Asset Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
	CountAssets(ctx context.Context) (int64, error)
	FindAllAssets(ctx context.Context) ([]*entities.Asset, error)
	FindAssetsFromCheckpoint(ctx context.Context, checkpoint *entities.Checkpoint) ([]*entities.Asset, error)
}

// Writer interface
type Writer interface {
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
