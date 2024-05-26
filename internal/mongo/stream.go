package mongo

import (
	"errors"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var changeStream *mongo.ChangeStream = nil

type ChangeStreamCallback func(event *ChangeEvent, json []byte)

type resumeToken struct {
	Value string `bson:"value"`
}

func (c *Client) StartChangeStream(callback ChangeStreamCallback) {
	c.Logger.Info("Starting mongo change stream")

	resumeTokens := c.Db.Collection(c.Config.ResumeTokensCollection)

	for {
		lastResumeToken := &resumeToken{}
		findOneOpts := options.FindOne().
			SetSort(bson.D{{Key: "_id", Value: -1}})

		err := resumeTokens.FindOne(c.Context, bson.D{}, findOneOpts).Decode(lastResumeToken)

		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			c.Logger.Error("Could not fetch or decode resume token", logger.AsError(err))
			panic("Could not fetch or decode resume token")
		}

		changeStreamOpts := options.ChangeStream().
			SetMaxAwaitTime(2 * time.Second).
			SetFullDocument(options.UpdateLookup)

		if lastResumeToken.Value != "" {
			c.Logger.Info("Resuming after token", "token", lastResumeToken.Value)
			changeStreamOpts.SetResumeAfter(bson.D{{Key: "_data", Value: lastResumeToken.Value}})
		}

		stream, err := c.Db.Watch(c.Context, constructPipeline(c.Config.WatchCollections, c.Config.WatchOperations), changeStreamOpts)

		if err != nil {
			c.Logger.Error("Not able to start mongo change stream", logger.AsError(err))
			panic("Not able to start mongo change stream")
		}

		changeStream = stream

		for changeStream.Next(c.Context) {
			json, err := bson.MarshalExtJSON(changeStream.Current, false, false)

			if err != nil {
				c.Logger.Error("Could not marshal change event to json", logger.AsError(err))
				break
			}

			event, err := DecodeChangeEvent(json)

			if err != nil {
				c.Logger.Error("Unable to decode received change event", "data", string(json))
				break
			}

			callback(event, json)

			if _, err = resumeTokens.InsertOne(c.Context, &resumeToken{Value: event.ResumeToken.Value}); err != nil {
				c.Logger.Error("Could not insert resume token", logger.AsError(err))
				break
			}
		}

		c.Logger.Info("Stopped watching mongodb collections")
	}
}

func (c *Client) StopChangeStream() {
	if changeStream != nil {
		c.Logger.Info("Closing mongo change stream")

		changeStream.Close(c.Context)
		changeStream = nil
	}
}
