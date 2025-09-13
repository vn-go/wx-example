package messaging

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	nats2 "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/nats.go"
)

// ////////////////////////////////////////////////////////
// NATS implementation
// ////////////////////////////////////////////////////////
type NatsBus struct {
	publisher  *nats2.Publisher
	subscriber *nats2.Subscriber
	natsURL    string
}

func NewNatsBus(natsURL, queueGroup string) (*NatsBus, error) {
	logger := watermill.NewStdLogger(false, false)

	pub, err := nats2.NewPublisher(nats2.PublisherConfig{
		URL:       natsURL,
		Marshaler: &nats2.JSONMarshaler{},
	}, logger)
	if err != nil {
		return nil, err
	}

	sub, err := nats2.NewSubscriber(nats2.SubscriberConfig{
		URL:              natsURL,
		QueueGroupPrefix: queueGroup,
		SubscribersCount: 1,
		AckWaitTimeout:   30 * time.Second,
		Unmarshaler:      &nats2.JSONMarshaler{},
		JetStream: nats2.JetStreamConfig{
			Disabled: false,
		},
	}, logger)
	if err != nil {
		return nil, err
	}

	return &NatsBus{
		publisher:  pub,
		subscriber: sub,
		natsURL:    natsURL,
	}, nil
}

func (b *NatsBus) Publish(ctx context.Context, topic string, msg *message.Message) error {
	return b.publisher.Publish(topic, msg)
}

type initStreamNameFromTopic struct {
	val  string
	once sync.Once
}

var caheStreamNameFromTopic sync.Map

func streamNameFromTopic(topic string) string {
	parts := strings.Split(topic, ".")
	if len(parts) > 0 {
		return strings.ToUpper(parts[0]) + "_STREAM"
	}
	return "DEFAULT_STREAM"
}
func (b *NatsBus) Subscribe(ctx context.Context, topic string, handler func(msg *message.Message) error) error {
	// Đảm bảo stream tồn tại trước khi subscribe
	streamName := streamNameFromTopic(topic)
	subject := topic
	if err := setupStream(b.natsURL, streamName, subject); err != nil {
		return err
	}

	msgs, err := b.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg); err != nil {
				if ok := msg.Nack(); !ok {
					log.Printf("[NATS SUB] Nack error")
				}
			} else {
				if ok := msg.Ack(); !ok {
					log.Printf("[NATS SUB] Ack error")
				}
			}
		}
	}()
	return nil
}
func (b *NatsBus) Close() error {
	if err := b.publisher.Close(); err != nil {
		return err
	}
	return b.subscriber.Close()
}

type initSetupStream struct {
	once sync.Once
	err  error
}

var cachesetupStream sync.Map

func setupStream(natsURL, streamName, subject string) error {
	actually, _ := cachesetupStream.LoadOrStore(fmt.Sprintf("%s/%s@%s", streamName, subject, natsURL), &initSetupStream{})
	init := actually.(*initSetupStream)
	init.once.Do(func() {
		init.err = setupStreamInternal(natsURL, streamName, subject)
	})
	return init.err

}
func setupStreamInternal(natsURL, streamName, subject string) error {
	// Kết nối tới NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return err
	}
	defer nc.Drain()

	// Lấy JetStream context
	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	// Tạo stream nếu chưa tồn tại
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.FileStorage, // hoặc MemoryStorage
		//Durable:  true, //<-- unknown field Durable in struct literal of type "github.com/nats-io/nats.go".StreamConfigcompilerMissingLitField
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return err
	}

	log.Printf("Stream %s created for subject %s\n", streamName, subject)
	return nil
}
