package storage

import (
	"log"
	"notification/internal/dlq"
	"notification/internal/model"
	"time"
)

func RegisterClient(userID uint32, ch chan string) {
	model.Manager.Mu.Lock()
	defer model.Manager.Mu.Unlock()
	model.Manager.Clients[userID] = ch
}

func UnregisterClient(userID uint32) {
	model.Manager.Mu.Lock()
	defer model.Manager.Mu.Unlock()
	delete(model.Manager.Clients, userID)
}

// SendToUser отправляет сообщение пользователю или пишет в DLQ в случае неудачи
func SendToUser(userID uint32, msg string, dlqClient *dlq.DLQClient) error {
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		model.Manager.Mu.RLock()
		ch, ok := model.Manager.Clients[userID]
		model.Manager.Mu.RUnlock()

		if ok {
			if err := sendWithRetry(ch, msg, 3, dlqClient); err != nil {
				return err
			}
			return nil
		}

		delay := time.Duration(i) * time.Second
		log.Printf("User %d not connected. Retry %d in %v", userID, i+1, delay)
		time.Sleep(delay)
	}

	log.Printf("User %d never connected after %d retries", userID, maxRetries)

	if dlqClient != nil {
		if err := dlqClient.WriteToDLQ([]byte(msg), "user_not_connected"); err != nil {
			log.Printf("DLQ write failed: %v", err)
			return err
		}
	}

	return nil
}

// sendWithRetry пытается отправить сообщение с повторами
func sendWithRetry(ch chan string, msg string, maxRetries int, dlqClient *dlq.DLQClient) error {
	for i := 0; i < maxRetries; i++ {
		select {
		case ch <- msg:
			return nil
		default:
			delay := time.Duration(i) * time.Second
			log.Printf("Retry %d for message '%s' in %v", i+1, msg, delay)
			time.Sleep(delay)
		}
	}

	log.Printf("Failed to send after %d retries: %s", maxRetries, msg)

	if dlqClient != nil {
		if err := dlqClient.WriteToDLQ([]byte(msg), "max_retries"); err != nil {
			log.Printf("DLQ write failed: %v", err)
			return err
		}
	}

	return nil
}
