package smtpreset

import (
	"fmt"
	"notification/internal/notifications"
	"notification/pkg/email"
	"notification/pkg/logger"
	"time"
)

type SMTPService struct {
	sender email.Sender
}

func NewSMTPService(sender email.Sender) *SMTPService {
	return &SMTPService{
		sender: sender,
	}
}

func (s *SMTPService) HandleResetNotification(n *notifications.Notification, key string) {
	if n.Category != "AUTHRESET" {
		logger.Error("❌ Неверная категория уведомления: " + n.Category)
		return
	}

	if n.Subtype != "SEND_RESET_CODE" && n.Subtype != "RESET_CODE" {
		logger.Error("❌ Неверный подтип уведомления: " + n.Subtype)
		return
	}

	code, okCode := n.Payload["code"].(string)
	subject, okSubject := n.Payload["subject"].(string)
	expiresAtFloat, okExpiresAt := n.Payload["expiresAt"].(float64)
	payloadEmail, okEmail := n.Payload["email"].(string)

	if !okCode || !okSubject || !okExpiresAt || !okEmail {
		logger.Error("❌ Недостаточно данных для обработки уведомления AUTHRESET")
		return
	}

	expiresAt := int64(expiresAtFloat)
	ttl := expiresAt - time.Now().Unix()

	body := fmt.Sprintf("Ваш код для сброса пароля: %s\n\nОн истекает через %d секунд.", code, ttl)

	emailInput := email.SendEmailInput{
		To:      payloadEmail,
		Subject: subject,
		Body:    body,
	}

	if err := s.sender.Send(emailInput); err != nil {
		logger.Error(fmt.Sprintf("❌ Ошибка при отправке письма: %s", err.Error()))
		return
	}

	logger.Info("📨 Email с кодом сброса пароля отправлен: " + payloadEmail)
}
