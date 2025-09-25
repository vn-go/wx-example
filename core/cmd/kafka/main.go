package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vn-go/bx"
)

type TestService struct {
}
type TestType struct {
	Code string
}

func main() {
	ctx := context.Background()

	bus, err := bx.Msg.NewKafkaMessageBus([]string{"localhost:9092"})
	if err != nil {
		log.Fatal(err)
	}
	// tstSvc := &TestService{}
	bk := bx.NewBrokerService("test", bus)
	defer bk.CloseAll()
	bk.Subcribe("test-001", &TestType{}, func(ctx context.Context, msg *bx.MsgItem) error {
		fmt.Println("OK")
		return nil
	})
	// defer bus.Close()

	// // Subscribe trước
	// err = bus.Subscribe(ctx, "test-topic", func(msg *message.Message) error {
	// 	log.Printf("📩 Received: %s", string(msg.Payload))
	// 	return nil
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	for i := 0; i < 10; i++ {
		bk.Publish(ctx, "test-001", &TestType{
			Code: "123356",
		})
	}
	// err = bus.Publish(ctx, "test-topic", message.NewMessage(watermill.NewUUID(), []byte("Hello from Watermill!")))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// }
	// Publish

	select {} // giữ chương trình chạy
}
