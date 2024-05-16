package mongo

import (
	"context"
	"os"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Conn   *mongo.Client
	Logger logger.Logger
}

func InitClient(log logger.Logger, uri string) *Client {
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Error("Error connecting to mongo", logger.AsError(err))
		os.Exit(1)
	}

	return &Client{
		Conn:   conn,
		Logger: log,
	}
}

func (c *Client) Watch(database string, collections, operations []string) *mongo.ChangeStream {
	db := c.Conn.Database(database)

	opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second)
	changeStream, err := db.Watch(context.TODO(), constructPipeline(collections, operations), opts)

	if err != nil {
		c.Logger.Error("Error starting mongo change stream", logger.AsError(err))
		os.Exit(1)
	}

	return changeStream
}
