// +build prod

package config

import "os"

var host = os.Getenv("MONGO_DB_HOST")
var username = os.Getenv("MO