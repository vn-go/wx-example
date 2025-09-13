package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Kết nối tới NATS server
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Cannot connect to NATS: %v", err)
	}
	defer nc.Drain() // đảm bảo connection được đóng gọn gàng

	// Khởi tạo JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Cannot get JetStream context: %v", err)
	}

	// Tạo stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "MYSTREAM",
		Subjects: []string{"mytopic.*"},
	})
	if err != nil {
		log.Fatalf("Cannot add stream: %v", err)
	}

	log.Println("Stream MYSTREAM created successfully")
}
