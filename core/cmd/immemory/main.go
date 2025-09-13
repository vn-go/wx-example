package main

import (
	"context"
	"core/messaging"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func main() {
	topic := "xxxxx"
	ctx := context.Background()
	// client := redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// 	DB:   0,
	// })

	bus := messaging.NewInMemoryBus()

	defer bus.Close()

	// Subscriber
	handler := func(msg *message.Message) error {
		log.Printf("[SUB] Received message: %s\n", string(msg.Payload))
		if ok := msg.Ack(); !ok {
			log.Printf("[SUB] Ack failed")
		}
		return nil
	}

	if err := bus.Subscribe(ctx, topic, handler); err != nil {
		log.Fatal(err)
	}

	// Publisher loop
	go func() {
		for i := 0; ; i++ {
			m := message.NewMessage(
				watermill.NewUUID(),
				[]byte("Hello memory "+time.Now().Format(time.RFC3339)),
			)
			if err := bus.Publish(context.Background(), topic, m); err != nil {
				log.Printf("[PUB] publish error: %v", err)
			} else {
				log.Printf("[PUB] Published message: %s", string(m.Payload))
			}
			time.Sleep(2 * time.Second)
		}
	}()

	select {} // block main thread, chạy liên tục
}
