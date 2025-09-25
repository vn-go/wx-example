package broker

import "context"

type Msg struct {
	Data []byte
}
type Bus interface {
	Close() error
	SubscribeRaw(topic string, fn func(m *Msg) error) error
	PublishRaw(ctx context.Context, topic string, data []byte) error
}
