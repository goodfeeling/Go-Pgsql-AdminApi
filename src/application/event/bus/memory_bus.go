package bus

import (
	"context"
	"sync"

	"github.com/gbrayhan/microservices-go/src/application/event/model"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"go.uber.org/zap"
)

// InMemoryEventBus 内存事件总线实现
type InMemoryEventBus struct {
	handlers map[string][]model.EventHandler
	mutex    sync.RWMutex
	logger   *logger.Logger
}

// NewInMemoryEventBus 创建内存事件总线
func NewInMemoryEventBus(logger *logger.Logger) EventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]model.EventHandler),
		logger:   logger,
	}
}

// Subscribe 订阅事件
func (eb *InMemoryEventBus) Subscribe(eventType string, handler model.EventHandler) error {
	eb.logger.Info("Subscribing to event", zap.String("eventType", eventType))

	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)

	eb.logger.Debug("Event handler subscribed successfully",
		zap.String("eventType", eventType),
		zap.Int("handlerCount", len(eb.handlers[eventType])))

	return nil
}

// Unsubscribe 取消订阅事件
func (eb *InMemoryEventBus) Unsubscribe(eventType string, handler model.EventHandler) error {
	eb.logger.Info("Unsubscribing from event", zap.String("eventType", eventType))

	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	handlers, exists := eb.handlers[eventType]
	if !exists {
		eb.logger.Warn("No handlers found for event type", zap.String("eventType", eventType))
		return nil
	}

	for i, h := range handlers {
		if h == handler {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			eb.logger.Debug("Event handler unsubscribed successfully",
				zap.String("eventType", eventType),
				zap.Int("remainingHandlers", len(eb.handlers[eventType])))
			break
		}
	}

	return nil
}

// Publish 发布事件
func (eb *InMemoryEventBus) Publish(ctx context.Context, event model.ApplicationEvent) error {
	eb.logger.Info("Publishing event",
		zap.String("eventType", event.EventType()),
		zap.String("eventID", event.EventID()))

	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	handlers, exists := eb.handlers[event.EventType()]
	if !exists || len(handlers) == 0 {
		eb.logger.Debug("No handlers found for event", zap.String("eventType", event.EventType()))
		return nil
	}

	eb.logger.Debug("Found handlers for event",
		zap.String("eventType", event.EventType()),
		zap.Int("handlerCount", len(handlers)))

	// 异步处理所有事件处理器
	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h model.EventHandler) {
			defer wg.Done()

			eb.logger.Debug("Handling event",
				zap.String("eventType", event.EventType()),
				zap.String("eventID", event.EventID()))

			if err := h.Handle(event); err != nil {
				eb.logger.Error("Error handling event",
					zap.String("eventType", event.EventType()),
					zap.String("eventID", event.EventID()),
					zap.Error(err))
				errChan <- err
			} else {
				eb.logger.Debug("Event handled successfully",
					zap.String("eventType", event.EventType()),
					zap.String("eventID", event.EventID()))
			}
		}(handler)
	}

	wg.Wait()
	close(errChan)

	// 处理错误
	for err := range errChan {
		eb.logger.Error("Error in event handling",
			zap.String("eventType", event.EventType()),
			zap.Error(err))
	}

	eb.logger.Info("Event publishing completed",
		zap.String("eventType", event.EventType()),
		zap.String("eventID", event.EventID()))

	return nil
}
