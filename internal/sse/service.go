package sse

import (
	"context"
	"encoding/json"
	"log"
	"notification/internal/dlq"
	"notification/internal/model"
	"notification/internal/storage"
	"notification/pkg/mongo"
)

type NotificationService struct {
	repository *mongo.Mongo   // Репозиторий для работы с базой данных
	dlqClient  *dlq.DLQClient // Клиент для работы с Dead Letter Queue
}

// NewNotificationService — конструктор для NotificationService
func NewNotificationService(repository *mongo.Mongo, dlqClient *dlq.DLQClient) *NotificationService {
	return &NotificationService{
		repository: repository,
		dlqClient:  dlqClient,
	}
}

// HandleMessageNotification — обработчик уведомления о сообщении
func (s *NotificationService) HandleMessageNotification(n model.Notification) {
	if err := s.SaveNotificationToDB(n); err != nil {
		log.Printf("⚠️ Could not save to DB: %v", err)
	}

	// Преобразуем уведомление в строку
	msg, err := json.Marshal(n)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	// Пытаемся отправить уведомление пользователю
	if err := storage.SendToUser(n.UserID, string(msg), s.dlqClient); err != nil {
		log.Printf("Failed to send notification for User %d: %v", n.UserID, err)
	}
}

// HandleFriendRequestNotification — обработчик уведомления о запросе в друзья
func (s *NotificationService) HandleFriendRequestNotification(n model.Notification) {
	if err := s.SaveNotificationToDB(n); err != nil {
		log.Printf("⚠️ Could not save to DB: %v", err)
	}

	// Преобразуем уведомление в строку
	msg, err := json.Marshal(n)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	// Пытаемся отправить уведомление пользователю
	if err := storage.SendToUser(n.UserID, string(msg), s.dlqClient); err != nil {
		log.Printf("Failed to send notification for User %d: %v", n.UserID, err)
	}
}

// SaveNotificationToDB — сохраняет уведомление в базе данных
func (s *NotificationService) SaveNotificationToDB(n model.Notification) error {
	_, err := s.repository.Database.Collection("notifications").InsertOne(context.Background(), n)
	if err != nil {
		log.Printf("Error saving notification to DB: %v", err)
		return err
	}
	return nil
}
