package messaging

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type RedisMessageBus struct {
	publisher  *redisstream.Publisher
	subscriber *redisstream.Subscriber
	logger     watermill.LoggerAdapter
}

func newRedisSubscriber(redisAddr string, logger watermill.LoggerAdapter) (*redisstream.Subscriber, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})
	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:       client,
		Unmarshaller: redisstream.DefaultMarshallerUnmarshaller{},
	}, logger)
	return sub, err
}
func newRedisPublisher(redisAddr string, logger watermill.LoggerAdapter) (*redisstream.Publisher, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})

	config := redisstream.PublisherConfig{
		Client: client,

		Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
		//ConsumerGroup: consumerGroup,
	}

	pub, err := redisstream.NewPublisher(config, logger)
	return pub, err
}
func NewRedisMessageBus(redisAddr string, consumerGroup string) (*RedisMessageBus, error) {
	logger := watermill.NewStdLogger(false, false)
	pub, err := newRedisPublisher(redisAddr, logger)
	if err != nil {
		return nil, err
	}

	sub, err := newRedisSubscriber(redisAddr, logger)
	if err != nil {
		return nil, err
	}

	return &RedisMessageBus{
		publisher:  pub,
		subscriber: sub,
		logger:     logger,
	}, nil
}

func (r *RedisMessageBus) Publish(ctx context.Context, topic string, msg *message.Message) error {
	return r.publisher.Publish(topic, msg)
}

func (r *RedisMessageBus) Subscribe(ctx context.Context, topic string, handler func(msg *message.Message) error) error {
	messages, err := r.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			if err := handler(msg); err != nil {
				msg.Nack()
			} else {
				msg.Ack()
			}
		}
	}()
	return nil
}

func (r *RedisMessageBus) Close() error {
	if err := r.publisher.Close(); err != nil {
		return err
	}
	if err := r.subscriber.Close(); err != nil {
		return err
	}
	return nil
}
