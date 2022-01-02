// +build prod

package config

import "os"

var host = os.Getenv("MONGO_DB_HOST")
var username = os.Getenv("MONGO_DB_USERNAME")
var password = os.Getenv("MONGO_DB_PASSWORD")

// AppConf constants
var AppConf = AppConfig{
	Mongo: MongoConfig{
		TimeoutMS:     360000,
		MinPoolSize:   5,
		MaxPoolSize:   10,
		MaxIdleTimeMS: 360000,
		Host:          host,
		Username:      username,
		Password:      password,
		Dbname:        "povi",
		SchemaVersion: "1",
		Colnames: map[string]string{
	