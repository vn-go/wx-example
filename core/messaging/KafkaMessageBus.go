package messaging

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type KafkaMessageBus struct {
	publisher  *kafka.Publisher
	subscriber *kafka.Subscriber
	logger     watermill.LoggerAdapter
}

func NewKafkaMessageBus(brokers []string) (*KafkaMessageBus, error) {
	logger := watermill.NewStdLogger(false, false)

	pub, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:                brokers,
			Unmarshaler:            kafka.DefaultMarshaler{},
			ConsumerGroup:          "my-consumer-group",
			InitializeTopicDetails: &sarama.TopicDetail{NumPartitions: 1, ReplicationFactor: 1},
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	return &KafkaMessageBus{
		publisher:  pub,
		subscriber: sub,
		logger:     logger,
	}, nil
}

func (k *KafkaMessageBus) Publish(ctx context.Context, topic string, msg *message.Message) error {
	return k.publisher.Publish(topic, msg)
}

func (k *KafkaMessageBus) Subscribe(ctx context.Context, topic string, handler func(msg *message.Message) error) error {
	messages, err := k.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-messages:
				if !ok {
					return
				}
				if err := handler(msg); err != nil {
					k.logger.Error("handler error", err, nil)
					msg.Nack()
				} else {
					msg.Ack()
				}
			}
		}
	}()
	return nil
}

func (k *KafkaMessageBus) Close() error {
	if err := k.publisher.Close(); err != nil {
		return err
	}
	return k.subscriber.Close()
}
