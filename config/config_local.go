// +build local

package config

// AppConf constants
var AppConf = AppConfig{
	Mongo: MongoConfig{
		TimeoutMS:     360000,
		MinPoolSize:   5,
	