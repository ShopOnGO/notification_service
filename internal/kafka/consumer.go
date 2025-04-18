package kafka

import (
	"context"
	"log"
	"notification/configs"

	"github.com/segmentio/kafka-go"
)

// ConsumeKafka читает сообщения из Kafka и передает их в диспетчер
func ConsumeKafka(conf *configs.Config, dispatch func([]byte)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{conf.Consumer.Broker},
		Topic:    conf.Consumer.Topic,
		GroupID:  "notification-service",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error due reading from Kafka:", err)
			continue
		}
		dispatch(m.Value) // Передаем сообщение в диспетчер
	}
}
