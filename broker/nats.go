package broker

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsBus struct {
	nc     *nats.Conn
	logger Logger
}

type initNatsBusSubscribeRaw struct {
	val  *nats.Subscription
	err  error
	once sync.Once
}

var initNatsBusSubscribeRawMap sync.Map

func (bus *NatsBus) SubscribeRaw(topic string, fn func(m *Msg) error) error {
	a, _ := initNatsBusSubscribeRawMap.LoadOrStore(topic, &initNatsBusSubscribeRaw{})
	i := a.(*initNatsBusSubscribeRaw)
	i.once.Do(func() {
		i.val, i.err = bus.nc.Subscribe(topic, func(msg *nats.Msg) { // tao go routine de chay cai nay lien tuc
			err := fn(&Msg{
				Data: msg.Data,
			})
			if err != nil {
				bus.logger.Error(err, topic)

			}

		})

	})
	if i.err != nil {
		initNatsBusSubscribeRawMap.Delete(topic)
		return i.err
	}
	return nil

}
func (bus *NatsBus) PublishRaw(ctx context.Context, topic string, data []byte) error {
	return bus.nc.Publish(topic, data)
}
func (bus *NatsBus) Close() error {
	var err error
	initNatsBusSubscribeRawMap.Range(func(key, value any) bool {
		err = value.(*initNatsBusSubscribeRaw).val.Unsubscribe()
		return true

	})
	if err != nil {
		return err
	}
	return bus.nc.Drain()
}
func NewNatsBus(servers string, logger Logger) (Bus, error) {

	nc, err := nats.Connect(
		servers,
		nats.MaxReconnects(-1), // reconnect vô hạn
		nats.ReconnectWait(nats.DefaultReconnectWait), // 2s giữa các lần reconnect)
	)
	if err != nil {
		return nil, err
	}
	return &NatsBus{
		nc:     nc,
		logger: logger,
	}, nil
}
