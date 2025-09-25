package broker

import (
	"fmt"
	"log"
	"strings"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestNats(t *testing.T) {
	urls := []string{
		"nats://127.0.0.1:4222", // nats1
		"nats://127.0.0.1:4223", // nats2
		"nats://127.0.0.1:4224", // nats3
	}
	buLog, er := NewDefaultLogger("./logs/log.txt")
	assert.NoError(t, er)
	bus, err := NewNatsBus(strings.Join(urls, ","), buLog)
	assert.NoError(t, err)
	t.Log(bus)
	// nc, err := nats.Connect(nats.DefaultURL) // "nats://127.0.0.1:4222"
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//defer bus.Close()

	// // Subscribe vào 1 subject
	err = bus.SubscribeRaw("test", func(m *Msg) error {
		fmt.Println(string(m.Data))
		return nil

	})
	if err != nil {
		log.Fatal(err)
	}
	err = bus.PublishRaw(t.Context(), "test", []byte("hello"))
	assert.NoError(t, err)
	// _, err = nc.Subscribe("foo", func(m *nats.Msg) {
	// 	fmt.Printf("Nhận message: %s\n", string(m.Data))
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // Publish 1 message
	// err = nc.Publish("foo", []byte("Xin chào NATS!"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Đảm bảo message đã được gửi đi
	// nc.Flush()

	// if err := nc.LastError(); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Đã gửi message, nhấn Ctrl+C để thoát...")
	select {} // block main goroutine
}
func TestRabbitMq(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	assert.NoError(t, err)
	ch, err := conn.Channel()
	assert.NoError(t, err)
	//failOnError(err, "Failed to open a channel")
	defer ch.Close()
	// 3. Khai báo queue (nếu chưa có thì sẽ tạo mới)
	// 1. Tạo exchange loại fanout
	err = ch.ExchangeDeclare(
		"logs",   // exchange name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	assert.NoError(t, err)
	//failOnError(err, "Failed to declare an exchange")

	// 2. Gửi message
	body := "Hello Broadcast!"
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key (fanout không cần)
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	assert.NoError(t, err)
	fmt.Println(" [x] Sent:", body)
}
