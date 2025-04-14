package main

import (
	"notification/configs"
	"notification/internal/dlq"
	"notification/internal/kafka"
	"notification/internal/sse"
	"notification/pkg/mongo"

	"github.com/gin-gonic/gin"
)

func main() {
	conf := configs.LoadConfig()
	db := mongo.NewMongo(conf)
	dlqClient := dlq.NewDLQClient([]string{"kafka:9092"}, "notifications-dlq")

	// обработчики различного вида : smtp, sse, websocket
	sseDispatcherService := sse.NewNotificationService(db, dlqClient)

	// Создаем диспетчер
	dispatcher := kafka.NewDispatcher()

	dispatcher.Register("MESSAGE", sseDispatcherService.HandleMessageNotification)
	dispatcher.Register("FRIEND_REQUEST", sseDispatcherService.HandleFriendRequestNotification)

	go kafka.ConsumeKafka(dispatcher.Dispatch)

	r := gin.Default()

	// SSE endpoint
	r.GET("/sse/:userID", sse.SSEHandler)

	r.Run(":8079")
}
