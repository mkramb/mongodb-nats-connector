package mongo

import (
	"context"
	"errors"
	"net/url"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ErrClientDisconnected = errors.New("could not reach mongodb: connection closed")

type Options struct {
	Context context.Context
	Logger  logger.Logger
	Config  *config.MongoConfig
}

type Client struct {
	Db   *mongo.Database
	Conn *mongo.Client
	Options
}

func (o Options) New() *Client {
	parsedURI, err := url.Parse(o.Config.ServerUri)

	if err != nil {
		o.Logger.Error("Invalid mongo URI", logger.AsError(err))
		panic("Invalid mongo URI")
	}

	if parsedURI.Path == "" || parsedURI.Path == "/" {
		o.Logger.Error("Database not provided in URI")
		panic("Database not provided URI")
	}

	conn, err := mongo.Connect(o.Context, options.Client().ApplyURI(o.Config.ServerUri))

	if err != nil {
		o.Logger.Error("Error connecting to mongo", logger.AsError(err))
		panic("Error connecting to mongo")
	}

	database := parsedURI.Path[1:]
	db := conn.Database(database)

	o.Logger.Info("Connected to Mongo")

	return &Client{
		Db:      db,
		Conn:    conn,
		Options: o,
	}
}

func (c *Client) Monitor() error {
	if err := c.Conn.Ping(c.Context, readpref.Primary()); err != nil {
		return ErrClientDisconnected
	}

	return nil
}

func (c *Client) Close() {
	c.Logger.Info("Closing mongo client")
	err := c.Conn.Ping(c.Context, nil)

	if err == nil {
		c.Conn.Disconnect(c.Context)
	}
}
