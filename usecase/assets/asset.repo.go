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
	FindAllAssets(ct