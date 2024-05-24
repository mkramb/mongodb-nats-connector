package mongo

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
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

type ChangeStreamCallback func(json []byte)

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

func (c *Client) StartWatch() *mongo.ChangeStream {
	c.Logger.Info("Starting mongo watcher")

	collections := c.Config.WatchCollections
	operations := c.Config.WatchOperations

	opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second).SetFullDocument(options.UpdateLookup)
	stream, err := c.Db.Watch(c.Context, constructPipeline(collections, operations), opts)

	if err != nil {
		c.Logger.Error("Error starting mongo change stream", logger.AsError(err))
		panic("Error starting mongo change stream")
	}

	return stream
}

func (c *Client) IterateChangeStream(changeStream *mongo.ChangeStream, callback ChangeStreamCallback) {
	for changeStream.Next(c.Context) {
		json, err := bson.MarshalExtJSON(changeStream.Current, false, false)

		if err != nil {
			c.Logger.Error("Could not marshal change event from bson", logger.AsError(err))
		} else {
			callback(json)
		}
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
