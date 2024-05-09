package nats

import (
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	Conn      *nats.Conn
	JetStream nats.JetStreamContext
}

func NewNatsClient(url string) (*NatsClient, error) {
	nc, err := nats.Connect(url)

	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()

	if err != nil {
		nc.Close()
		return nil, err
	}

	return &NatsClient{
		Conn:      nc,
		JetStream: js,
	}, nil
}
