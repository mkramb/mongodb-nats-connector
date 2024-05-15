package mongo

import (
	"context"
	"os"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Conn *mongo.Client
}

func InitClient(log logger.Logger, uri string) *Client {
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Error("Error connecting to mongo", logger.AsError(err))
		os.Exit(1)
	}

	return &Client{
		Conn: conn,
	}
}
