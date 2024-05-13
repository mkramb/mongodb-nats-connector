package nats

import (
	"log/slog"
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

func InitNats(logger logger.Logger, url string) *Nats {
	nc, err := nats.Connect(url)

	if err != nil {
		logger.Error("Error connecting to nats", slog.Any("err", err))
		os.Exit(1)
	}

	js, err := nc.JetStream()

	if err != nil {
		nc.Close()

		logger.Error("Error creating jetstream context", slog.Any("err", err))
		os.Exit(1)
	}

	return &Nats{Client: natsClient{
		Conn:      nc,
		JetStream: js,
	}}
}
