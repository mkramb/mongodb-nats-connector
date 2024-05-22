package mongo

import (
	"context"
	"net/url"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Db      *mongo.Database
	Conn    *mongo.Client
	Logger  logger.Logger
	Context context.Context
}

type ChangeStreamCallback func(data []byte)

func InitClient(ctx context.Context, log logger.Logger, uri string) *Client {
	parsedURI, err := url.Parse(uri)

	if err != nil {
		log.Error("Invalid MongoDB URI", logger.AsError(err))
		panic("Invalid MongoDB URI")
	}

	if parsedURI.Path == "" || parsedURI.Path == "/" {
		log.Error("Database not provided in MongoDB URI")
		panic("Database not provided in MongoDB URI")
	}

	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Error("Error connecting to mongo", logger.AsError(err))
		panic("Error connecting to mongo")
	}

	database := parsedURI.Path[1:]
	db := conn.Database(database)

	return &Client{
		Db:      db,
		Conn:    conn,
		Logger:  log,
		Context: ctx,
	}
}

func (c *Client) Watch(collections, operations []string) *mongo.ChangeStream {
	c.Logger.Info("Starting mongo watcher")

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
