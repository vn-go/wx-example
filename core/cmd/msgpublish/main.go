package main

import (
	"context"
	"log"
	"time"

	// thay bằng đường dẫn module thật

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/vn-go/bx"
)

func main() {
	natsURL := "nats://localhost:4222"
	topic := "abc"
	bus, err := bx.Msg.NewNatsBus(natsURL, "test")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// Subscriber
	handler := func(msg *message.Message) error {
		log.Printf("[SUB] Received message: %s\n", string(msg.Payload))
		if ok := msg.Ack(); !ok {
			log.Printf("[SUB] Ack failed")
		}
		return nil
	}

	ctx := context.Background()
	if err := bus.Subscribe(ctx, topic, handler); err != nil {
		log.Fatal(err)
	}

	// Publisher loop
	go func() {
		for i := 0; ; i++ {
			m := message.NewMessage(
				watermill.NewUUID(),
				[]byte("Hello RabbitMQ "+time.Now().Format(time.RFC3339)),
			)
			if err := bus.Publish(ctx, topic, m); err != nil {
				log.Printf("[PUB] publish error: %v", err)
			} else {
				log.Printf("[PUB] Published message: %s", string(m.Payload))
			}
			time.Sleep(2 * time.Second)
		}
	}()

	select {} // block main thread, chạy liên tục
}
