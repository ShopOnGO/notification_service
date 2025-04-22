package sse

import (
	"encoding/json"
	"log"
	"notification/internal/dlq"
	"notification/internal/notifications"
	"notification/manager"
)

type NotificationService struct {
	repository    *notifications.NotificationRepository // Репозиторий для работы с базой данных
	dlqClient     *dlq.DLQClient                        // Клиент для работы с Dead Letter Queue
	clientManager *manager.ClientManager
}

// NewNotificationService — конструктор для NotificationService
func NewNotificationService(repository *notifications.NotificationRepository, dlqClient *dlq.DLQClient, clientManager *manager.ClientManager) *NotificationService {
	return &NotificationService{
		repository:    repository,
		dlqClient:     dlqClient,
		clientManager: clientManager,
	}
}

func (s *NotificationService) HandleMessageNotification(n *notifications.Notification, key string) {
	wasInDLQ := n.WasInDLQ
	if err := s.SaveNotificationToDB(n); err != nil {
		log.Printf("⚠️ Could not save to DB: %v", err)
	}

	msg, err := json.Marshal(n)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	if err := s.clientManager.SendToUser(n.UserID, string(msg), key, s.dlqClient, wasInDLQ); err != nil {
		log.Printf("Failed to send notification for User %d: %v", n.UserID, err)
	}
}

func (s *NotificationService) HandleFriendRequestNotification(n *notifications.Notification, key string) {
	wasInDLQ := n.WasInDLQ
	if err := s.SaveNotificationToDB(n); err != nil {
		log.Printf("⚠️ Could not save to DB: %v", err)
	}

	msg, err := json.Marshal(n)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	if err := s.clientManager.SendToUser(n.UserID, string(msg), key, s.dlqClient, wasInDLQ); err != nil {
		log.Printf("Failed to send notification for User %d: %v", n.UserID, err)
	}
}

func (s *NotificationService) SaveNotificationToDB(n *notifications.Notification) error {
	if n.WasInDLQ {
		return nil
	}
	_, err := s.repository.Add(n)
	if err != nil {
		log.Printf("Error saving notification to DB: %v", err)
		return err
	}
	return nil
}
