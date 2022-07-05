
package repos

import (
	"context"
	"fmt"
	"time"

	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/config"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/consts"
	"github.com/lenoobz/aws-yahoo-asset-price-scraper/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AssetMongo struct
type AssetMongo struct {
	db     *mongo.Database
	client *mongo.Client
	log    logger.ContextLog
	conf   *config.MongoConfig
}

// NewAssetMongo creates new asset mongo repo
func NewAssetMongo(db *mongo.Database, log logger.ContextLog, conf *config.MongoConfig) (*AssetMongo, error) {
	if db != nil {
		return &AssetMongo{
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

	return &AssetMongo{
		db:     client.Database(conf.Dbname),
		client: client,
		log:    log,
		conf:   conf,
	}, nil
}

// Close disconnect from database
func (r *AssetMongo) Close() {
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

// CountAssets count number of assets available
func (r *AssetMongo) CountAssets(ctx context.Context) (int64, error) {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.ASSETS_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return 0, fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{}

	// find options
	countOptions := options.Count()

	return col.CountDocuments(ctx, filter, countOptions)
}

// FindAllAssets find all assets
func (r *AssetMongo) FindAllAssets(ctx context.Context) ([]*entities.Asset, error) {
	// create new context for the query
	ctx, cancel := createContext(ctx, r.conf.TimeoutMS)
	defer cancel()

	// what collection we are going to use
	colname, ok := r.conf.Colnames[consts.ASSETS_COLLECTION]
	if !ok {
		r.log.Error(ctx, "cannot find collection name")
		return nil, fmt.Errorf("cannot find collection name")
	}
	col := r.db.Collection(colname)

	// filter
	filter := bson.D{}

	// find options
	findOptions := options.Find()

	cur, err := col.Find(ctx, filter, findOptions)

	// only run defer function when find success
	if cur != nil {
		defer func() {
			if deferErr := cur.Close(ctx); deferErr != nil {
				err = deferErr
			}
		}()
	}

	// find was not succeed
	if err != nil {
		r.log.Error(ctx, "find query failed", "error", err)
		return nil, err
	}

	var assets []*entities.Asset

	// iterate over the cursor to decode document one at a time
	for cur.Next(ctx) {
		// decode cursor to activity model
		var asset entities.Asset
		if err = cur.Decode(&asset); err != nil {
			r.log.Error(ctx, "decode failed", "error", err)
			return nil, err
		}

		assets = append(assets, &asset)
	}

	if err := cur.Err(); err != nil {
		r.log.Error(ctx, "iterate over cursor failed", "error", err)
		return nil, err
	}

	return assets, nil
}