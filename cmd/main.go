package main

import (
	"notification/configs"
	"notification/internal/dlq"
	"notification/internal/kafka"
	"notification/internal/notifications"
	passwordreset "notification/internal/smtp_passwordreset"
	"notification/internal/sse"
	"notification/manager"
	"notification/pkg/email/smtp"
	"notification/pkg/mongo"

	"github.com/gin-gonic/gin"
)

func main() {
	conf := configs.LoadConfig()
	db := mongo.NewMongo(conf)

	dlqClient := dlq.NewDLQClient([]string{conf.Dlq.Broker}, conf.Dlq.Topic)

	// клиентs ClientManager
	clientManager := manager.NewClientManager()

	//sender
	smtpSender := smtp.NewSMTPSender(
		conf.SMTP.Name,
		conf.SMTP.From,
		conf.SMTP.Pass,
		conf.SMTP.Host,
		conf.SMTP.Port,
	)

	//repository
	sseNotificationRepository := notifications.NewNotificationRepository(db.Database)

	// обработчики различного вида : smtp, sse, websocket
	sseDispatcherService := sse.NewNotificationService(sseNotificationRepository, dlqClient, clientManager)
	smtpDispatcherService := passwordreset.NewSMTPService(smtpSender)

	// Создаем диспетчер категорий сообщений
	dispatcher := kafka.NewDispatcher()

	dispatcher.Register("MESSAGE", sseDispatcherService.HandleMessageNotification)
	dispatcher.Register("FRIEND_REQUEST", sseDispatcherService.HandleFriendRequestNotification)
	dispatcher.Register("AUTHRESET", smtpDispatcherService.HandleResetNotification)

	go kafka.ConsumeKafka(conf, conf.Consumer.Topic, "notifications-group", dispatcher.Dispatch) //когда пользователь не подключен походу нет проверки на was in dlq
	go kafka.ConsumeKafka(conf, conf.SMTPreset.Consumer, "auth_reset_group", dispatcher.Dispatch)

	r := gin.Default()

	// SSE endpoint
	r.GET("/sse/:userID", sse.SSEHandler(clientManager))
	r.GET("/sse/status/:userID", sse.SSEStatusHandler(clientManager))

	r.Run(":8079")
}
