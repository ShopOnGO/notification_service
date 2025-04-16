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

// HandleMessageNotification — обработчик уведомления о сообщении
func (s *NotificationService) HandleMessageNotification(n *notifications.Notification) {
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
	if err := s.clientManager.SendToUser(n.UserID, string(msg), s.dlqClient); err != nil {
		log.Printf("Failed to send notification for User %d: %v", n.UserID, err)
	}
}

// HandleFriendRequestNotification — обработчик уведомления о запросе в друзья
func (s *NotificationService) HandleFriendRequestNotification(n *notifications.Notification) {
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
	if err := s.clientManager.SendToUser(n.UserID, string(msg), s.dlqClient); err != nil {
		log.Printf("Failed to send notification for User %d: %v", n.UserID, err)
	}
}

func (s *NotificationService) SaveNotificationToDB(n *notifications.Notification) error {
	_, err := s.repository.Add(n)
	if err != nil {
		log.Printf("Error saving notification to DB: %v", err)
		return err
	}
	return nil
}
