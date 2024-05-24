package nats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var ErrClientDisconnected = errors.New("could not reach nats: connection closed")

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

type PublishOptions struct {
	MsgId   string
	Subject string
	Data    []byte
}

func (o Options) New() *Client {
	opts := nats.Options{
		AllowReconnect: true,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
		Url:            o.Config.ServerUrl,
	}

	conn, err := opts.Connect()

	if err != nil {
		o.Logger.Error("Error connecting to nats", logger.AsError(err))
		panic("Error connecting to nats")
	}

	js, _ := jetstream.New(conn)

	o.Logger.Info("Connected to Nats")

	return &Client{
		Conn:      conn,
		JetStream: js,
		Options:   o,
	}
}

func (c *Client) PublishEvent(opts *PublishOptions) {
	c.Logger.Info("Emitting event to nats jetstream", "subject", opts.Subject)

	_, err := c.JetStream.PublishMsg(c.Context, &nats.Msg{
		Subject: fmt.Sprintf("%v.%v", c.Config.StreamName, opts.Subject),
		Data:    opts.Data,
	}, jetstream.WithMsgID(opts.MsgId))

	if err != nil {
		c.Logger.Error("Could not publish nats message",
			"data", opts.Data, "subject", opts.Subject, logger.AsError(err))
	}
}

func (c *Client) Monitor() error {
	if closed := c.Conn.IsClosed(); closed {
		return ErrClientDisconnected
	}

	return nil
}

func (c *Client) Close() {
	c.Logger.Info("Closing nats client")

	if c.Conn.IsConnected() {
		c.Conn.Close()
	}
}
