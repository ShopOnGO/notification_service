package kafka

import (
	"encoding/json"
	"log"
	"notification/internal/notifications"
)

type Dispatcher struct {
	handlers map[string]func(*notifications.Notification) // map категории к обработчику
}

// NewDispatcher создает новый экземпляр Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]func(*notifications.Notification)),
	}
}

// Register регистрирует обработчик для категории
func (d *Dispatcher) Register(cat string, fn func(*notifications.Notification)) {
	d.handlers[cat] = fn
}

// Dispatch обрабатывает входящее сообщение и вызывает нужный обработчик
func (d *Dispatcher) Dispatch(msg []byte) {
	var n notifications.Notification
	if err := json.Unmarshal(msg, &n); err != nil {
		log.Printf("🚨 Parse error: %v | Raw: %s", err, string(msg))
		return
	}

	// Валидация обязательных полей
	if n.UserID == 0 || n.Category == "" {
		log.Printf("⚠️ Invalid notification: %+v", n)
		return
	}

	// Ищем обработчик по категории
	if handler, ok := d.handlers[n.Category]; ok {
		handler(&n) // Передаем только уведомление в обработчик
	} else {
		log.Printf("⚠️ No handler found for category: %s", n.Category)
	}
}
