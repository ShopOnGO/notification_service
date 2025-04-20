package kafka

import (
	"encoding/json"
	"log"
	"notification/internal/notifications"
)

var allowedCategoriesWithoutUserID = []string{
	"AUTHRESET",
	"RESET_CODE", //нету пока что
	// другие категории, где userID не обязателен
}

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
	log.Printf("[DLQ] 📦 Raw message: %s", string(msg))

	if err := json.Unmarshal(msg, &n); err != nil {
		log.Printf("🚨 Parse error: %v | Raw: %s", err, string(msg))
		return
	}

	if n.Category == "" {
		log.Printf("⚠️ Invalid notification (missing category): %+v", n)
		return
	}

	// Если userID == 0, проверяем, что категория находится в списке допустимых
	if n.UserID == 0 {
		// Проверяем, если категория не в списке разрешённых
		if !contains(allowedCategoriesWithoutUserID, n.Category) {
			log.Printf("🚨 Invalid notification: userID is required for category %s", n.Category)
			return
		}
	}

	// Ищем обработчик по категории
	if handler, ok := d.handlers[n.Category]; ok {
		handler(&n) // Передаем только уведомление в обработчик
	} else {
		log.Printf("⚠️ No handler found for category: %s", n.Category)
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
