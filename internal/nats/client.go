package nats

import (
	"os"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/nats-io/nats.go"
)

type Client struct {
	Conn      *nats.Conn
	JetStream nats.JetStreamContext
}

func InitClient(log logger.Logger, url string) *Client {
	conn, err := nats.Connect(url)

	if err != nil {
		log.Error("Error connecting to nats", logger.AsError(err))
		os.Exit(1)
	}

	js, _ := conn.JetStream()

	return &Client{
		Conn:      conn,
		JetStream: js,
	}
}
