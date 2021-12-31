// +build prod

package config

import "os"

var host = os.Getenv("MONGO_DB_HOST")
var username = os.Getenv("MONGO_DB_USERNAME")
var password = os.Getenv("MONGO_DB_PASSWORD")

// AppConf constants
var AppConf =