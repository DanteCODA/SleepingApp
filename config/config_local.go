// +build local

package config

// AppConf constants
var AppConf = AppConfig{
	Mongo: MongoConfig{
		TimeoutMS:     360000,
		MinPoolSize:   5,
		MaxPoolSize:   10,
		MaxIdleTimeMS: 360000,
		H