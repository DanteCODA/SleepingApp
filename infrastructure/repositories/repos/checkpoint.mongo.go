
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
		clientOptions.SetMaxPoolSize(conf.MaxPoolSize)
	}

	// set max idle time ms
	if conf.MaxIdleTimeMS > 0 {
		clientOptions.SetMaxConnIdleTime(time.Duration(conf.MaxIdleTimeMS) * time.Millisecond)
	}

	// construct a connection string from mongo config object
	cxnString := fmt.Sprintf("mongodb+srv://%s:%s@%s", conf.Username, conf.Password, conf.Host)

	// create mongo client by making new connection
	client, err := mongo.Connect(ctx, clientOptions.ApplyURI(cxnString))
	if err != nil {
		return nil, err
	}

	return &CheckpointMongo{
		db:     client.Database(conf.Dbname),
		client: client,
		log:    log,
		conf:   conf,
	}, nil
}

// Close disconnect from database
func (r *CheckpointMongo) Close() {
	ctx := context.Background()
	r.log.Info(ctx, "close mongo client")

	if r.client == nil {
		return
	}

	if err := r.client.Disconnect(ctx); err != nil {
		r.log.Error(ctx, "disconnect mongo failed", "error", err)
	}
}

///////////////////////////////////////////////////////////////////////////////
// Implement interface
///////////////////////////////////////////////////////////////////////////////

// UpdateCheckpoint updates a checkpoint given page size and number of asset
func (r *CheckpointMongo) UpdateCheckpoint(ctx context.Context, pageSize int64, numAssets int64) (*entities.Checkpoint, error) {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.SCRAPE_CHECKPOINT_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return nil, fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{}

	// find options
	findOptions := options.FindOne()

	cur := col.FindOne(ctx, filter, findOptions)

	// only run defer function when find success
	err := cur.Err()

	if err == mongo.ErrNoDocuments {
		// decode cursor to activity model
		cp, err := models.NewCheckPointModel(ctx, r.log, pageSize, r.conf.SchemaVersion)
		if err != nil {
			r.log.Error(ctx, "create model failed", "error", err)
			return nil, err
		}

		return r.updateCheckPoint(ctx, col, cp)
	}

	// find was not succeed
	if err != nil {
		r.log.Error(ctx, "find query failed", "error", err)
		return nil, err
	}

	var checkpoint models.CheckPointModel
	if err = cur.Decode(&checkpoint); err != nil {
		r.log.Error(ctx, "decode failed", "error", err)
		return nil, err
	}

	if checkpoint.PriceCheckPoint == nil {
		checkpoint.PriceCheckPoint = &models.PriceCheckPointModel{
			PageSize:  pageSize,
			PrevIndex: 0,
		}
		return r.updateCheckPoint(ctx, col, &checkpoint)
	}

	currNumAssets := checkpoint.PriceCheckPoint.PrevIndex*checkpoint.PriceCheckPoint.PageSize + checkpoint.PriceCheckPoint.PageSize
	if currNumAssets >= numAssets {
		checkpoint.PriceCheckPoint.PrevIndex = 0
	} else {
		checkpoint.PriceCheckPoint.PrevIndex = checkpoint.PriceCheckPoint.PrevIndex + 1
	}

	checkpoint.PriceCheckPoint.PageSize = pageSize

	return r.updateCheckPoint(ctx, col, &checkpoint)
}

// updateCheckPoint update checkpoint
func (r *CheckpointMongo) updateCheckPoint(ctx context.Context, col *mongo.Collection, checkpoint *models.CheckPointModel) (*entities.Checkpoint, error) {
	// filter
	filter := bson.D{}
	if checkpoint.ID != nil {
		filter = bson.D{{Key: "_id", Value: checkpoint.ID}}
	} else {
		checkpoint.CreatedAt = time.Now().UTC().Unix()
	}

	// update
	update := bson.D{
		{
			Key:   "$set",
			Value: checkpoint,
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		r.log.Error(ctx, "update one failed", "error", err)
		return nil, err
	}

	return &entities.Checkpoint{
		PageSize:  checkpoint.PriceCheckPoint.PageSize,
		PageIndex: checkpoint.PriceCheckPoint.PrevIndex,
	}, nil
}