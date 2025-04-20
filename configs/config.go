package configs

import (
	"os"
)

type Config struct {
	Mongo     MongoConfig
	Dlq       DlqConfig
	Consumer  Consumer
	SMTP      SMTPConfig
	SMTPreset SMTPreset
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
type SMTPreset struct {
	Consumer string
}
type SMTPConfig struct {
	Name string
	From string
	Pass string
	Host string
	Port int
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
		SMTP: SMTPConfig{
			Name: os.Getenv("SMTP_NAME"),
			From: os.Getenv("SMTP_FROM"),
			Pass: os.Getenv("SMTP_PASS"),
			Host: os.Getenv("SMTP_HOST"),
			Port: 587, // TLS
		},
		SMTPreset: SMTPreset{
			Consumer: os.Getenv("KAFKA_SMTP_CONSUMER")},
	}
}
