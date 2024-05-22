package nats

import (
	"context"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	Conn      *nats.Conn
	JetStream jetstream.JetStream
	Logger    logger.Logger
	Context   context.Context
}

func InitClient(ctx context.Context, log logger.Logger, url string) *Client {
	opts := nats.Options{
		Url:            url,
		AllowReconnect: true,
		MaxReconnect:   -1,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
	}

	conn, err := opts.Connect()

	if err != nil {
		log.Error("Error connecting to nats", logger.AsError(err))
		panic("Error connecting to nats")
	}

	js, _ := jetstream.New(conn)

	return &Client{
		Conn:      conn,
		JetStream: js,
		Logger:    log,
		Context:   ctx,
	}
}

func (c *Client) Close() {
	c.Logger.Info("Closing nats client")

	if c.Conn.IsConnected() {
		c.Conn.Close()
	}
}
