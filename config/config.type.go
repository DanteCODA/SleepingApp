package config

// MongoConfig struct
type MongoConfig struct {
	TimeoutMS     uint64
	MinPoolSize   uint64
	MaxPoolSize   uint64
	MaxIdleTimeMS uint64
	SchemaVersion string
	Username      string
	Password      string
	Host          string
	Dbname        string
	Colnames      map[string]string
}

// AppConfig struct
type AppConfig struct {
	Mongo MongoConfig
}
