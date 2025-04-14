package configs

import (
	"os"
)

type Config struct {
	Mongo MongoConfig
}

type MongoConfig struct {
	URI      string
	Database string
}

func LoadConfig() *Config {
	return &Config{
		Mongo: MongoConfig{
			URI:      os.Getenv("MONGO_URI"),
			Database: os.Getenv("MONGO_DB"),
		},
	}
}
