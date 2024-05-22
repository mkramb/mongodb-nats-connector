package mongo

import (
	"context"
	"net/url"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

type ChangeStreamCallback func(data []byte)

func (o Options) NewClient() *Client {
	parsedURI, err := url.Parse(o.Config.ServerUri)

	if err != nil {
		o.Logger.Error("Invalid MongoDB URI", logger.AsError(err))
		panic("Invalid MongoDB URI")
	}

	if parsedURI.Path == "" || parsedURI.Path == "/" {
		o.Logger.Error("Database not provided in MongoDB URI")
		panic("Database not provided in MongoDB URI")
	}

	conn, err := mongo.Connect(o.Context, options.Client().ApplyURI(o.Config.ServerUri))

	if err != nil {
		o.Logger.Error("Error connecting to mongo", logger.AsError(err))
		panic("Error connecting to mongo")
	}

	database := parsedURI.Path[1:]
	db := conn.Database(database)

	return &Client{
		Db:      db,
		Conn:    conn,
		Options: o,
	}
}

func (c *Client) Watch() *mongo.ChangeStream {
	c.Logger.Info("Starting mongo watcher")

	collections := c.Config.WatchCollections
	operations := c.Config.WatchOperations

	opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second)
	stream, err := c.Db.Watch(c.Context, constructPipeline(collections, operations), opts)

	if err != nil {
		c.Logger.Error("Error starting mongo change stream", logger.AsError(err))
		panic("Error starting mongo change stream")
	}

	return stream
}

func (c *Client) IterateChangeStream(changeStream *mongo.ChangeStream, callback ChangeStreamCallback) {
	for changeStream.Next(c.Context) {
		data, err := bson.MarshalExtJSON(changeStream.Current, false, false)

		if err != nil {
			c.Logger.Error("Could not decode mongo change event", logger.AsError(err))
		}

		callback(data)
	}
}

func (c *Client) Close() {
	c.Logger.Info("Closing mongo client")
	err := c.Conn.Ping(c.Context, nil)

	if err == nil {
		c.Conn.Disconnect(c.Context)
	}
}
