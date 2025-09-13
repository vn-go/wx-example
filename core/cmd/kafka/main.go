package main

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/vn-go/bx"
)

func main() {
	ctx := context.Background()

	bus, err := bx.Msg.NewKafkaMessageBus([]string{"localhost:9092"})
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// Subscribe trước
	err = bus.Subscribe(ctx, "test-topic", func(msg *message.Message) error {
		log.Printf("📩 Received: %s", string(msg.Payload))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Publish
	err = bus.Publish(ctx, "test-topic", message.NewMessage(watermill.NewUUID(), []byte("Hello from Watermill!")))
	if err != nil {
		log.Fatal(err)
	}

	select {} // giữ chương trình chạy
}
