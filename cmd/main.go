package main

import (
	"notification/configs"
	"notification/internal/dlq"
	"notification/internal/kafka"
	"notification/internal/notifications"
	"notification/internal/sse"
	"notification/manager"
	"notification/pkg/mongo"

	"github.com/gin-gonic/gin"
)

func main() {
	conf := configs.LoadConfig()
	db := mongo.NewMongo(conf)
	dlqClient := dlq.NewDLQClient([]string{conf.Dlq.Broker}, conf.Dlq.Topic)

	// клиентs ClientManager
	clientManager := manager.NewClientManager()

	//repository
	sseNotificationRepository := notifications.NewNotificationRepository(db.Database)

	// обработчики различного вида : smtp, sse, websocket
	sseDispatcherService := sse.NewNotificationService(sseNotificationRepository, dlqClient, clientManager)

	// Создаем диспетчер категорий сообщений
	dispatcher := kafka.NewDispatcher()

	dispatcher.Register("MESSAGE", sseDispatcherService.HandleMessageNotification)
	dispatcher.Register("FRIEND_REQUEST", sseDispatcherService.HandleFriendRequestNotification)

	go kafka.ConsumeKafka(dispatcher.Dispatch)

	r := gin.Default()

	// SSE endpoint
	r.GET("/sse/:userID", sse.SSEHandler(clientManager))
	r.GET("/sse/status/:userID", sse.SSEStatusHandler(clientManager))

	r.Run(":8079")
}
