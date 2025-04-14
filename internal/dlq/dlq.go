package dlq

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type DLQClient struct {
	writer *kafka.Writer
}

// consumer would be later
func NewDLQClient(brokers []string, topic string) *DLQClient {
	return &DLQClient{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      brokers,
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: int(kafka.RequireOne),
			BatchTimeout: 10 * time.Millisecond,
		},
		)}
} // the copy of main_app kafkaService/kafka.go

func (c *DLQClient) WriteToDLQ(msg []byte, reason string) error {
	fmt.Println("dlq")
	return c.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: msg,
			Headers: []kafka.Header{
				{Key: "X-Error-Reason", Value: []byte(reason)},
			},
		},
	)
}
