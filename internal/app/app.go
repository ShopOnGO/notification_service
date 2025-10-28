package app

import (
	"context"
	"fmt"
	"log"
	"notification/configs"
	"notification/internal/dispatcher"
	"notification/internal/dlq"
	"notification/internal/notifications"
	"notification/internal/smtpreset"
	"notification/internal/sse"
	"notification/manager"
	"notification/pkg/email/smtp"
	"notification/pkg/mongo"

	shopKafka "github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/gin-gonic/gin"
)

func App() *gin.Engine {
	conf := configs.LoadConfig()
	ctx := context.Background()

	db := mongo.NewMongo(conf)

	dlqClient := dlq.NewDLQClient([]string{conf.Dlq.Broker}, conf.Dlq.Topic)

	clientManager := manager.NewClientManager()

	smtpSender := smtp.NewSMTPSender(
		conf.SMTP.Name,
		conf.SMTP.From,
		conf.SMTP.Pass,
		conf.SMTP.Host,
		conf.SMTP.Port,
	)

	sseNotificationRepository := notifications.NewNotificationRepository(db.Database)

	sseDispatcherService := sse.NewNotificationService(sseNotificationRepository, dlqClient, clientManager)
	smtpDispatcherService := smtpreset.NewSMTPService(smtpSender)

	// consumer handler
	consumerService := dispatcher.NewConsumerService(sseDispatcherService, smtpDispatcherService, dlqClient)

	for _, topic := range conf.Consumer.Topics {
		if topic == "" {
			continue
		}

		log.Printf("[Kafka] Starting consumer for topic: %s in group: %s", topic, "notifications-group")

		clientID := fmt.Sprintf("notification-service-%s", topic)

		consumer := shopKafka.NewConsumer(
			[]string{conf.Consumer.Broker},
			topic,
			"notifications-group",
			clientID,
		)

		go func(c *shopKafka.KafkaService, t string) {
			log.Printf("[Kafka] Consumer for topic %s started", t)
			// Убедитесь, что ваш .Consume() блокирующий и обрабатывает ошибки внутри
			c.Consume(ctx, consumerService.Handler())
		}(consumer, topic)
	}

	smtpConsumer := shopKafka.NewConsumer(
		[]string{conf.Consumer.Broker},
		conf.SMTPreset.Consumer,
		"auth_reset_group",
		"smtp-reset-service",
	)

	// стартуем консюмер
	go smtpConsumer.Consume(ctx, consumerService.Handler())

	// gin router
	router := gin.Default()

	// SSE endpoints
	router.GET("/sse/:userID", sse.SSEHandler(clientManager))
	router.GET("/sse/status/:userID", sse.SSEStatusHandler(clientManager))

	return router
}
