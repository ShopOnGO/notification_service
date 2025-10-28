package configs

import (
	"log"
	"os"
	"strings"
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
	Topics []string
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

	topicsStr := os.Getenv("KAFKA_NOTIFICATION_TOPICS")
	if topicsStr == "" {
		log.Fatal("KAFKA_NOTIFICATION_TOPICS is not set in env")
	}
	topics := strings.Split(topicsStr, ",")

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
			Topics: topics,
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
