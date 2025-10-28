package dispatcher

import (
	"encoding/json"
	"fmt"
	"log"

	"notification/internal/dlq"
	"notification/internal/notifications"
	"notification/internal/smtpreset"
	"notification/internal/sse"

	"github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"
	"github.com/segmentio/kafka-go"
)

var allowedCategoriesWithoutUserID = []string{
	"AUTHRESET",
	"RESET_CODE",
}

type ConsumerService struct {
	dispatcher *kafkaService.Dispatcher
}

func NewConsumerService(
	sseHandler *sse.NotificationService,
	smtpHandler *smtpreset.SMTPService,
	dlqClient *dlq.DLQClient,
) *ConsumerService {
	d := kafkaService.NewDispatcher()

	// Регистрируем каждый обработчик под своим префиксом
	d.Register("notification-", wrapHandler(sseHandler.HandleMessageNotification, dlqClient))
	d.Register("reset-", wrapHandler(smtpHandler.HandleResetNotification, dlqClient))

	//d.Register("notification-AddNote", wrapHandler(sseHandler.HandleMessageNotification))
	//d.Register("notification-FriendReq", wrapHandler(sseHandler.HandleFriendRequestNotification))

	return &ConsumerService{dispatcher: d}
}

func (cs *ConsumerService) Handler() func(msg kafka.Message) error {
	return cs.dispatcher.Dispatch
}

// adapter pattern
func wrapHandler(fn func(*notifications.Notification, string), dlqClient *dlq.DLQClient) func(msg kafka.Message) error {
	return func(msg kafka.Message) error {
		key := string(msg.Key)
		log.Printf("[DLQ] 📦 Kafka Key: %s | Msg: %s", key, string(msg.Value))

		var n notifications.Notification
		if err := json.Unmarshal(msg.Value, &n); err != nil {
			log.Printf("🚨 [Poison Pill] Parse error: %v | Raw: %s. Sending to DLQ.", err, string(msg.Value))
			// 5. ОТПРАВЛЯЕМ В DLQ
			dlqClient.WriteToDLQ(key, msg.Value, fmt.Sprintf("json_unmarshal_error: %v", err))
			// 6. ВОЗВРАЩАЕМ NIL, ЧТОБЫ ЗАКОММИТИТЬ СООБЩЕНИЕ
			return nil
		}

		if n.Category == "" {
			log.Printf("⚠️ [Poison Pill] Invalid notification (missing category): %+v. Sending to DLQ.", n)
			dlqClient.WriteToDLQ(key, msg.Value, "missing_category")
			return nil // <--- 6. ВОЗВРАЩАЕМ NIL
		}

		// Эта проверка теперь будет работать и для 'PRODUCT_CREATED'
		if n.UserID == 0 && !contains(allowedCategoriesWithoutUserID, n.Category) {
			log.Printf("🚨 [Poison Pill] Invalid notification: userID is required for category %s. Sending to DLQ.", n.Category)
			dlqClient.WriteToDLQ(key, msg.Value, fmt.Sprintf("missing_userid_for_category_%s", n.Category))
			return nil // <--- 6. ВОЗВРАЩАЕМ NIL
		}

		fn(&n, key)
		return nil
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

// package dispatcher

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"notification/internal/notifications"
// 	"notification/internal/smtpreset"
// 	"notification/internal/sse"

// 	shopKafka "github.com/ShopOnGO/ShopOnGO/pkg/kafkaService" // <-- путь к kafkaService
// 	"github.com/segmentio/kafka-go"
// )

// var allowedCategoriesWithoutUserID = []string{
// 	"AUTHRESET",
// 	"RESET_CODE",
// 	// другие категории без userID при необходимости
// }

// type ConsumerService struct {
// 	dispatcher *shopKafka.Dispatcher
// }

// // NewConsumerService создает сервис с зарегистрированными хендлерами
// func NewConsumerService(
// 	sseHandler *sse.NotificationService,
// 	smtpHandler *smtpreset.SMTPService,
// ) *ConsumerService {
// 	d := shopKafka.NewDispatcher()

// 	d.Register("MESSAGE", wrapHandler(sseHandler.HandleMessageNotification))
// 	d.Register("FRIEND_REQUEST", wrapHandler(sseHandler.HandleFriendRequestNotification))
// 	d.Register("AUTHRESET", wrapHandler(smtpHandler.HandleResetNotification))

// 	return &ConsumerService{dispatcher: d}
// }

// // Handler возвращает функцию, которую можно передать в kafkaService.Consume
// func (cs *ConsumerService) Handler() func(msg kafka.Message) error {
// 	return cs.dispatcher.Dispatch
// }

// // Оборачивает хендлер в kafka.Message handler
// func wrapHandler(fn func(*notifications.Notification)) func(msg kafka.Message) error {
// 	return func(msg kafka.Message) error {
// 		log.Printf("[DLQ] 📦 Raw message: %s", string(msg.Value))

// 		var n notifications.Notification
// 		if err := json.Unmarshal(msg.Value, &n); err != nil {
// 			log.Printf("🚨 Parse error: %v | Raw: %s", err, string(msg.Value))
// 			return err
// 		}

// 		if n.Category == "" {
// 			log.Printf("⚠️ Invalid notification (missing category): %+v", n)
// 			return fmt.Errorf("missing category")
// 		}

// 		if n.UserID == 0 && !contains(allowedCategoriesWithoutUserID, n.Category) {
// 			log.Printf("🚨 Invalid notification: userID is required for category %s", n.Category)
// 			return fmt.Errorf("missing userID for category %s", n.Category)
// 		}

// 		fn(&n)
// 		return nil
// 	}
// }

// func contains(slice []string, str string) bool {
// 	for _, s := range slice {
// 		if s == str {
// 			return true
// 		}
// 	}
// 	return false
// }

// import (
// 	"encoding/json"
// 	"log"
// 	"notification/internal/notifications"
// )

// var allowedCategoriesWithoutUserID = []string{
// 	"AUTHRESET",
// 	"RESET_CODE", //нету пока что
// 	// другие категории, где userID не обязателен
// }

// type Dispatcher struct {
// 	handlers map[string]func(*notifications.Notification) // map категории к обработчику
// }

// // NewDispatcher создает новый экземпляр Dispatcher
// func NewDispatcher() *Dispatcher {
// 	return &Dispatcher{
// 		handlers: make(map[string]func(*notifications.Notification)),
// 	}
// }

// // Register регистрирует обработчик для категории
// func (d *Dispatcher) Register(cat string, fn func(*notifications.Notification)) {
// 	d.handlers[cat] = fn
// }

// // Dispatch обрабатывает входящее сообщение и вызывает нужный обработчик
// func (d *Dispatcher) Dispatch(msg []byte) {
// 	var n notifications.Notification
// 	log.Printf("[DLQ] 📦 Raw message: %s", string(msg))

// 	if err := json.Unmarshal(msg, &n); err != nil {
// 		log.Printf("🚨 Parse error: %v | Raw: %s", err, string(msg))
// 		return
// 	}

// 	if n.Category == "" {
// 		log.Printf("⚠️ Invalid notification (missing category): %+v", n)
// 		return
// 	}

// 	// Если userID == 0, проверяем, что категория находится в списке допустимых
// 	if n.UserID == 0 {
// 		// Проверяем, если категория не в списке разрешённых
// 		if !contains(allowedCategoriesWithoutUserID, n.Category) {
// 			log.Printf("🚨 Invalid notification: userID is required for category %s", n.Category)
// 			return
// 		}
// 	}

// 	// Ищем обработчик по категории
// 	if handler, ok := d.handlers[n.Category]; ok {
// 		handler(&n) // Передаем только уведомление в обработчик
// 	} else {
// 		log.Printf("⚠️ No handler found for category: %s", n.Category)
// 	}
// }

// func contains(slice []string, str string) bool {
// 	for _, s := range slice {
// 		if s == str {
// 			return true
// 		}
// 	}
// 	return false
// }
