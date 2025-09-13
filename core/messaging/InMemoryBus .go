package messaging

import (
	"context"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
)

type InMemoryBus struct {
	mu          sync.RWMutex
	subscribers map[string][]func(msg *message.Message) error
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		subscribers: make(map[string][]func(msg *message.Message) error),
	}
}

func (b *InMemoryBus) Publish(ctx context.Context, topic string, msg *message.Message) error {
	b.mu.RLock()
	handlers := b.subscribers[topic]
	b.mu.RUnlock()

	// Gửi tới tất cả subscriber của topic đó
	for _, h := range handlers {
		go func(handler func(msg *message.Message) error, m *message.Message) {
			_ = handler(m)
		}(h, msg)
	}
	return nil
}

func (b *InMemoryBus) Subscribe(ctx context.Context, topic string, handler func(msg *message.Message) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], handler)
	return nil
}

func (b *InMemoryBus) Close() error {
	// Không cần làm gì vì in-memory
	return nil
}
