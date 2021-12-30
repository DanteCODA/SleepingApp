// +build local

package config

// AppConf constants
var AppConf = AppConfig{
	Mongo: MongoConfig{
		TimeoutMS:     360000,
		MinPoolSize:   5,
		MaxPoolSize:   10,
		MaxIdleTimeMS: 360000,
		Host:          "lenoobdev.l8ckp.mongodb.net",
		Username:      "lenoob_dev",
		Password:      "lenoob_dev",
		Dbname:        "povi",
		SchemaVersion: "1",
		Colnames: map[string]string{
			"assets":            "assets",
			"asset_prices":      "asset_prices",
			"scrape_checkpoint": "scrape_checkpoint",