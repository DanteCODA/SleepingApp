package price

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
	InsertAssetPrice(ctx context.Context, assetPrice *entities.AssetPrice) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
