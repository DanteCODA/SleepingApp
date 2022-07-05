
package repos

import (
	"context"
	"fmt"
	"time"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/config"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/consts"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/infrastructure/repositories/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CheckpointMongo struct
type CheckpointMongo struct {
	db     *mongo.Database
	client *mongo.Client
	log    logger.ContextLog
	conf   *config.MongoConfig
}

// NewCheckpointMongo creates new checkpoint mongo repo
func NewCheckpointMongo(db *mongo.Database, log logger.ContextLog, conf *config.MongoConfig) (*CheckpointMongo, error) {
	if db != nil {
		return &CheckpointMongo{
			db:   db,
			log:  log,
			conf: conf,
		}, nil
	}

	// set context with timeout from the config
	// create new context for the query
	ctx, cancel := createContext(context.Background(), conf.TimeoutMS)
	defer cancel()

	// set mongo client options
	clientOptions := options.Client()

	// set min pool size
	if conf.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(conf.MinPoolSize)
	}

	// set max pool size
	if conf.MaxPoolSize > 0 {