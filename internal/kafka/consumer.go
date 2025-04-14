package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// ConsumeKafka читает сообщения из Kafka и передает их в диспетчер
func ConsumeKafka(dispatch func([]byte)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "notifications-topic",
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
