package configs

import (
	"log"
	"os"
	"strings"

	"github.com/ShopOnGO/ShopOnGO/configs"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

type Config struct {
	Mongo        MongoConfig
	Dlq          DlqConfig
	Consumer     Consumer
	SMTP         SMTPConfig
	SMTPreset    SMTPreset
	LogLevel     logger.LogLevel
	FileLogLevel logger.LogLevel
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
	// logger
	logLevelStr := os.Getenv("NOTIFICATION_SERVICE_LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "INFO"
	}
	LogLevel := configs.ParseLogLevel(logLevelStr)
	fileLogLevelStr := os.Getenv("NOTIFICATION_SERVICE_FILE_LOG_LEVEL")
	if fileLogLevelStr == "" {
		fileLogLevelStr = "INFO"
	}
	FileLogLevel := configs.ParseLogLevel(fileLogLevelStr)

	return &Config{
		Mongo: MongoConfig{
			URI:      os.Getenv("MONGO_URI"),
			Database: os.Getenv("MONGO_DB"),
		},
		Dlq: DlqConfig{
			Broker: os.Getenv("KAFKA_BROKERS"),
			Topic:  os.Getenv("KAFKA_NOTIFY_TOPIC"),
		},
		Consumer: Consumer{
			Broker: os.Getenv("KAFKA_BROKERS"),
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
		LogLevel:     LogLevel,
		FileLogLevel: FileLogLevel,
	}
}
