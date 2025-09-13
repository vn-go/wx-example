package messaging

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

// MessageBus defines the interface for a pub/sub system
type MessageBus interface {
	Publish(ctx context.Context, topic string, msg *message.Message) error
	Subscribe(ctx context.Context, topic string, handler func(msg *message.Message) error) error
	Close() error
}
