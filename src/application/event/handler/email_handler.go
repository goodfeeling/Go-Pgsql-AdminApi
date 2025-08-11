package handler

import (
	"fmt"
	"log"

	"github.com/gbrayhan/microservices-go/src/application/event/model"
)

// EmailEventHandler 邮件事件处理器
type EmailEventHandler struct {
	smtpServer string
	fromEmail  string
}

// NewEmailEventHandler 创建邮件事件处理器
func NewEmailEventHandler(smtpServer, fromEmail string) *EmailEventHandler {
	return &EmailEventHandler{
		smtpServer: smtpServer,
		fromEmail:  fromEmail,
	}
}

// Handle 处理事件
func (h *EmailEventHandler) Handle(event model.ApplicationEvent) error {
	switch event.EventType() {
	case "UserRegistered":
		return h.handleUserRegistered(event)
	case "OrderCreated":
		return h.handleOrderCreated(event)
	default:
		return nil
	}
}

func (h *EmailEventHandler) handleUserRegistered(event model.ApplicationEvent) error {
	payload := event.Payload().(map[string]interface{})
	email := payload["email"].(string)
	username := payload["username"].(string)

	// 发送欢迎邮件
	message := fmt.Sprintf("Welcome %s! Your account has been created.", username)
	log.Printf("Sending email to %s: %s", email, message)

	// 实际的邮件发送逻辑...
	return nil
}

func (h *EmailEventHandler) handleOrderCreated(event model.ApplicationEvent) error {
	payload := event.Payload().(map[string]interface{})
	email := payload["customerEmail"].(string)
	orderID := payload["orderID"].(string)

	// 发送订单确认邮件
	message := fmt.Sprintf("Your order %s has been created successfully.", orderID)
	log.Printf("Sending order confirmation to %s: %s", email, message)

	// 实际的邮件发送逻辑...
	return nil
}
