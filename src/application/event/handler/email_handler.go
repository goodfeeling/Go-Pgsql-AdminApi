// src/application/event/handler/email_handler.go
package handler

import (
	"fmt"
	"log"

	"github.com/gbrayhan/microservices-go/src/application/event/model"
	"github.com/gbrayhan/microservices-go/src/infrastructure/lib/email"
)

// EmailEventHandler 邮件事件处理器
type EmailEventHandler struct {
	emailService email.EmailService
}

// NewEmailEventHandler 创建邮件事件处理器
func NewEmailEventHandler() *EmailEventHandler {
	return &EmailEventHandler{
		emailService: email.NewSMTPEmailService(),
	}
}

// Handle 处理事件
func (h *EmailEventHandler) Handle(event model.ApplicationEvent) error {
	switch event.EventType() {
	case model.UserRegisteredEventType:
		return h.handleUserRegistered(event)
	case model.OrderCreatedEventType:
		return h.handleOrderCreated(event)
	case model.ForgetPasswordEventType:
		return h.handleForgetPassword(event)
	default:
		return nil
	}
}

func (h *EmailEventHandler) handleForgetPassword(event model.ApplicationEvent) error {
	log.Println("Handling forget password event")
	payload, ok := event.Payload().(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload format for forget password event")
	}

	to, ok := getStringValue(payload, "to")
	if !ok {
		return fmt.Errorf("missing 'to' field in payload")
	}

	subject, _ := getStringValue(payload, "subject")
	if subject == "" {
		subject = "密码重置"
	}

	body, ok := getStringValue(payload, "body")
	if !ok {
		return fmt.Errorf("missing 'body' field in payload")
	}

	log.Printf("Sending forget password email to %s", to)
	res := h.emailService.SendEmail(to, subject, body)
	log.Printf("Email sent: %v", res)
	return res
}

func (h *EmailEventHandler) handleUserRegistered(event model.ApplicationEvent) error {
	payload, ok := event.Payload().(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload format for user registered event")
	}

	emailAddr, ok := getStringValue(payload, "email")
	if !ok {
		return fmt.Errorf("missing 'email' field in payload")
	}

	username, _ := getStringValue(payload, "username")
	if username == "" {
		username = "用户"
	}

	subject := "欢迎注册"
	body := fmt.Sprintf("<h1>欢迎 %s!</h1><p>您的账户已成功创建。</p>", username)

	log.Printf("Sending welcome email to %s", emailAddr)
	res := h.emailService.SendEmail(emailAddr, subject, body)
	log.Printf("Email sent: %v", res)
	return res
}

func (h *EmailEventHandler) handleOrderCreated(event model.ApplicationEvent) error {
	payload, ok := event.Payload().(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload format for order created event")
	}

	emailAddr, ok := getStringValue(payload, "customerEmail")
	if !ok {
		return fmt.Errorf("missing 'customerEmail' field in payload")
	}

	orderID, ok := getStringValue(payload, "orderID")
	if !ok {
		return fmt.Errorf("missing 'orderID' field in payload")
	}

	subject := "订单确认"
	body := fmt.Sprintf("<h1>订单确认</h1><p>您的订单 %s 已成功创建。</p>", orderID)

	log.Printf("Sending order confirmation email to %s", emailAddr)
	return h.emailService.SendHTMLEmail(emailAddr, subject, body)
}

// getStringValue 安全地从map中获取字符串值
func getStringValue(data map[string]interface{}, key string) (string, bool) {
	value, exists := data[key]
	if !exists || value == nil {
		return "", false
	}

	if str, ok := value.(string); ok {
		return str, true
	}

	return fmt.Sprintf("%v", value), true
}
