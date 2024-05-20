package nats

import (
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/nats-io/nats.go"
)

type Client struct {
	Conn      *nats.Conn
	JetStream nats.JetStreamContext
	Logger    logger.Logger
}

func InitClient(log logger.Logger, url string) *Client {
	conn, err := nats.Connect(url)

	if err != nil {
		log.Error("Error connecting to nats", logger.AsError(err))
		panic("Error connecting to nats")
	}

	js, _ := conn.JetStream()

	return &Client{
		Conn:      conn,
		JetStream: js,
		Logger:    log,
	}
}

func (c *Client) Close() {
	c.Logger.Info("Closing nats client")

	if c.Conn.IsConnected() {
		c.Conn.Close()
	}
}
