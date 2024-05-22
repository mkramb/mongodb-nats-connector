package nats

import (
	"context"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Options struct {
	Context context.Context
	Logger  logger.Logger
	Config  *config.NatsConfig
}

type Client struct {
	Conn      *nats.Conn
	JetStream jetstream.JetStream
	Options
}

func (o Options) NewClient() *Client {
	opts := nats.Options{
		Url:            o.Config.ServerUrl,
		AllowReconnect: true,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
	}

	conn, err := opts.Connect()

	if err != nil {
		o.Logger.Error("Error connecting to nats", logger.AsError(err))
		panic("Error connecting to nats")
	}

	js, _ := jetstream.New(conn)

	return &Client{
		Conn:      conn,
		JetStream: js,
		Options:   o,
	}
}

func (c *Client) Close() {
	c.Logger.Info("Closing nats client")

	if c.Conn.IsConnected() {
		c.Conn.Close()
	}
}
