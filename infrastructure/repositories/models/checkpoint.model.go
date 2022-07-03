
package models

import (
	"context"
	"time"

	logger "github.com/lenoobz/aws-lambda-logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckPointModel struct
type CheckPointModel struct {
	ID              *primitive.ObjectID   `bson:"_id,omitempty"`
	CreatedAt       int64                 `bson:"createdAt,omitempty"`
	ModifiedAt      int64                 `bson:"modifiedAt,omitempty"`
	Enabled         bool                  `bson:"enabled"`
	Deleted         bool                  `bson:"deleted"`
	Schema          string                `bson:"schema,omitempty"`
	PriceCheckPoint *PriceCheckPointModel `bson:"priceCheckPoint,omitempty"`
}

type PriceCheckPointModel struct {
	PageSize  int64 `bson:"size,omitempty"`
	PrevIndex int64 `bson:"prevIndex"`
}

// NewCheckPointModel create checkpoint model
func NewCheckPointModel(ctx context.Context, log logger.ContextLog, pageSize int64, schemaVersion string) (*CheckPointModel, error) {
	return &CheckPointModel{
		ModifiedAt: time.Now().UTC().Unix(),
		Enabled:    true,