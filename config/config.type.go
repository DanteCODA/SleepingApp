package config

// MongoConfig struct
type MongoConfig struct {
	TimeoutMS     uint64
	MinPoolSize   uint64
	MaxPoolSize   uint64
	MaxIdleTimeMS uint64
	Schem