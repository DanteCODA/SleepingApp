package models

import (
	"context"
	"time"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssetPriceModel struct
type AssetPriceModel struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt  int64               `bson:"createdAt,omitempty"`
	ModifiedAt int64               `bson:"modifiedAt,omitempty"`
	Enabled    bool                `bson:"enabled"`
	Deleted    bool                `bson:"deleted"`
	Schema     string              `bson:"schema,omitempty"`
	Source     string              `bson:"source,omitempty"`
	Ticker     string              `bson:"ticker,omitempty"`
	Currency   string              `bson:"currency,omitempty"`
	Price      float64             `bson:"price,omitempty"`
}

// Ne