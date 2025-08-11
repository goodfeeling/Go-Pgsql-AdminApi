package handler

import (
	"log"

	"github.com/gbrayhan/microservices-go/src/application/event/model"
)

// NotificationEventHandler 通知事件处理器
type NotificationEventHandler struct {
	// notificationService NotificationService
}

// NewNotificationEventHandler 创建通知事件处理器
func NewNotificationEventHandler() *NotificationEventHandler {
	return &NotificationEventHandler{}
}

// Handle 处理事件
func (h *NotificationEventHandler) Handle(event model.ApplicationEvent) error {
	switch event.EventType() {
	case model.UserRegisteredEventType:
		return h.handleUserRegistered(event)
	case model.OrderCreatedEventType:
		return h.handleOrderCreated(event)
	default:
		return nil
	}
}

func (h *NotificationEventHandler) handleUserRegistered(event model.ApplicationEvent) error {
	payload := event.Payload().(map[string]interface{})
	userID := payload["userID"].(string)
	username := payload["username"].(string)

	// 创建欢迎通知
	log.Printf("Creating welcome notification for user %s (ID: %s)", username, userID)

	// 实际的通知创建逻辑...
	return nil
}

func (h *NotificationEventHandler) handleOrderCreated(event model.ApplicationEvent) error {
	payload := event.Payload().(map[string]interface{})
	userID := payload["customerID"].(string)
	orderID := payload["orderID"].(string)

	// 创建订单通知
	log.Printf("Creating order notification for user %s, order %s", userID, orderID)

	// 实际的通知创建逻辑...
	return nil
}
