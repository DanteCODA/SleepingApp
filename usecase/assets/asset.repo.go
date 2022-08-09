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
	CountAssets(ctx c