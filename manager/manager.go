package manager

import (
	"log"
	"notification/internal/dlq"
	"sync"
	"time"
)

type ClientManager struct {
	mu      sync.RWMutex
	clients map[uint32]chan string
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[uint32]chan string),
	}
}

// AddClient добавляет нового клиента
func (m *ClientManager) AddClient(userID uint32, ch chan string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[userID] = ch
}

// RemoveClient удаляет клиента
func (m *ClientManager) RemoveClient(userID uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, userID)
}

// GetClientChannel возвращает канал клиента, если он существует
func (m *ClientManager) GetClientChannel(userID uint32) (chan string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ch, ok := m.clients[userID]
	return ch, ok
}

// IsConnected проверяет наличие активного подключения
func (m *ClientManager) IsConnected(userID uint32) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.clients[userID]
	return ok
}

func (m *ClientManager) SendToUser(userID uint32, msg string, dlqClient *dlq.DLQClient) error {
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
		ch, ok := m.GetClientChannel(userID)
		if ok {
			if err := m.sendWithRetry(ch, msg, 3, dlqClient); err != nil {
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

func (m *ClientManager) sendWithRetry(ch chan string, msg string, maxRetries int, dlqClient *dlq.DLQClient) error {
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
