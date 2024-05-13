package nats

import (
	"os"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	Client natsClient
}

type natsClient struct {
	Conn      *nats.Conn
	JetStream nats.JetStreamContext
}

func InitNats(log logger.Logger, url string) *Nats {
	nc, err := nats.Connect(url)

	if err != nil {
		log.Error("Error connecting to nats", logger.AsError(err))
		os.Exit(1)
	}

	js, err := nc.JetStream()

	if err != nil {
		nc.Close()

		log.Error("Error creating jetstream context", logger.AsError(err))
		os.Exit(1)
	}

	return &Nats{Client: natsClient{
		Conn:      nc,
		JetStream: js,
	}}
}
