package configs

import (
	"os"
)

type Config struct {
	Mongo    MongoConfig
	Dlq      DlqConfig
	Consumer Consumer
}
type DlqConfig struct {
	Broker string
	Topic  string
}
type Consumer struct {
	Broker string
	Topic  string
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
		Dlq: DlqConfig{
			Broker: os.Getenv("KAFKA_BROKER"),
			Topic:  os.Getenv("KAFKA_TOPIC"),
		},
		Consumer: Consumer{
			Broker: os.Getenv("KAFKA_BROKER"),
			Topic:  os.Getenv("KAFKA_CONSUMER"),
		},
	}
}
