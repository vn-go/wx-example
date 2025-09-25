package core

import (
	"strings"

	"github.com/vn-go/bx"
)

func newBrokerService(cfg *configInfo) (bx.Broker, error) {
	if cfg.Broker.Bus == "rabbit" {
		bus, err := bx.Msg.NewRabbitMQBus(
			cfg.Broker.Rabbit.Url,
			cfg.Broker.Rabbit.Exchange,
			cfg.Broker.Rabbit.Queue,
		)
		if err != nil {
			return nil, err
		}
		return bx.NewBrokerService(cfg.Broker.Topic, bus), nil
	} else if cfg.Broker.Bus == "redis" {
		bus, err := bx.Msg.NewRedisMessage(
			cfg.Broker.Redis.Addr,
			cfg.Broker.Redis.Consumer,
		)
		if err != nil {
			return nil, err
		}
		return bx.NewBrokerService(cfg.Broker.Topic, bus), nil

	} else if cfg.Broker.Bus == "kafka" {
		bus, err := bx.Msg.NewKafkaMessageBus(
			strings.Split(cfg.Broker.Kafka.Brokers, ","),
		)
		if err != nil {
			return nil, err
		}
		return bx.NewBrokerService(cfg.Broker.Topic, bus), nil

	} else if cfg.Broker.Bus == "nats" {
		bus, err := bx.Msg.NewNatsBus(
			cfg.Broker.Nats.Server,
			cfg.Broker.Nats.Group,
		)
		if err != nil {
			return nil, err
		}
		return bx.NewBrokerService(cfg.Broker.Topic, bus), nil

	} else {
		return bx.NewBrokerService(cfg.Broker.Topic, bx.Msg.NewInMemoryBus()), nil
	}

}
