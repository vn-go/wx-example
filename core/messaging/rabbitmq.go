package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/rabbitmq/amqp091-go"

	//"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	//"github.com/ThreeDotsLabs/watermill/amqp"

	//"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"

	//rabbit "github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

// ////////////////////////////////////////////////////////
// RabbitMQ implementation
// ////////////////////////////////////////////////////////
type RabbitMQBus struct {
	publisher  *amqp.Publisher
	subscriber *amqp.Subscriber
}

func GenerateNameConstant(topic string) string {
	return topic
}

// Custom generators cho fixed names (nâng cao từ hàm của bạn)
func GenerateExchangeNameConstant(name string) func(string) string {
	return func(topic string) string { return name }
}

func GenerateRoutingKeyConstant(name string) func(string) string {
	return func(topic string) string { return name }
}
func newSubcriber(amqpURL, exchangeName, queueName string, logger watermill.LoggerAdapter) (*amqp.Subscriber, error) {
	cfg := amqp.Config{
		Connection: amqp.ConnectionConfig{
			AmqpURI: amqpURL,
		},
		Queue: amqp.QueueConfig{
			GenerateName: amqp.GenerateQueueNameConstant(queueName),
			Durable:      true,
			AutoDelete:   false,
		},
		Exchange: amqp.ExchangeConfig{
			GenerateName: amqp.GenerateQueueNameConstant(exchangeName),
			Type:         "direct",
			Durable:      true,
		},
		Marshaler: amqp.DefaultMarshaler{},
		QueueBind: amqp.QueueBindConfig{
			GenerateRoutingKey: func(topic string) string { return queueName },
		},
		Consume: amqp.ConsumeConfig{
			Qos: amqp.QosConfig{PrefetchCount: 1},
		},
		TopologyBuilder: &amqp.DefaultTopologyBuilder{}, // <--- quan trọng
	}

	sub, err := amqp.NewSubscriber(cfg, logger)

	return sub, err
}
func newPublisher(amqpURL, exchangeName, queueName string, logger watermill.LoggerAdapter) (*amqp.Publisher, error) {
	config := amqp.Config{
		Connection: amqp.ConnectionConfig{
			AmqpURI: amqpURL,
			AmqpConfig: &amqp091.Config{
				Heartbeat: 10 * time.Second,
			},
		},
		Exchange: amqp.ExchangeConfig{
			GenerateName: func(topic string) string { return exchangeName },
			Type:         "direct",
			Durable:      true,
		},
		Queue: amqp.QueueConfig{
			GenerateName: func(topic string) string { return queueName },
			Durable:      true,
			AutoDelete:   false,
		},
		Marshaler:       amqp.DefaultMarshaler{},
		TopologyBuilder: &amqp.DefaultTopologyBuilder{}, // <--- quan trọng
	}

	// Publisher
	config.Publish = amqp.PublishConfig{
		GenerateRoutingKey: func(topic string) string { return queueName },
	}
	pub, err := amqp.NewPublisher(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}
	return pub, err

}
func NewRabbitMQBus(amqpURL, exchangeName, queueName string) (*RabbitMQBus, error) {
	if amqpURL == "" {
		return nil, fmt.Errorf("amqpURL cannot be empty")
	}

	logger := watermill.NewStdLogger(false, false)
	pub, err := newPublisher(amqpURL, exchangeName, queueName, logger)
	if err != nil {
		return nil, err
	}
	// Common configuration

	// // Subscriber
	// config.Consume = amqp.ConsumeConfig{}
	// config.QueueBind = amqp.QueueBindConfig{
	// 	GenerateRoutingKey: func(topic string) string { return queueName },
	// }
	sub, err := newSubcriber(amqpURL, exchangeName, queueName, logger)
	if err != nil {
		pub.Close()
		return nil, fmt.Errorf("failed to create subscriber: %w", err)
	}

	return &RabbitMQBus{
		publisher:  pub,
		subscriber: sub,
	}, nil
}

func (b *RabbitMQBus) Publish(ctx context.Context, topic string, msg *message.Message) error {
	return b.publisher.Publish(topic, msg)
}

func (b *RabbitMQBus) Subscribe(ctx context.Context, topic string, handler func(msg *message.Message) error) error {
	msgs, err := b.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg); err != nil {
				if ok := msg.Nack(); !ok {
					log.Printf("[RABBIT SUB] Nack error")
				}
			} else {
				if ok := msg.Ack(); !ok {
					log.Printf("[RABBIT SUB] Ack error")
				}
			}
		}
	}()
	return nil
}

func (b *RabbitMQBus) Close() error {
	if err := b.publisher.Close(); err != nil {
		return err
	}
	return b.subscriber.Close()
}
