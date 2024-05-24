package raft

import (
	"fmt"

	"github.com/mkramb/mongodb-nats-connector/internal/mongo"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/nats-io/graft"
	_mongo "go.mongodb.org/mongo-driver/mongo"
)

var changeStream *_mongo.ChangeStream = nil

func (s *Server) stateHandler(stateTo graft.State) {
	switch stateTo {

	case graft.LEADER:
		s.Logger.Info("Becoming leader")

		changeStream = s.MongoClient.Watch()
		s.MongoClient.IterateChangeStream(changeStream, func(json []byte) {
			event, err := mongo.DecodeChangeEvent(json)

			if err != nil {
				s.Logger.Error("Unable to decode received change event", "data", string(json))
			} else {
				var opts nats.PublishOptions

				opts.MsgId = event.ResumeToken.Value
				opts.Subject = fmt.Sprintf("%v.%v.%v", event.Ns.Coll, event.OperationType, event.FullDocument.Id.Value)
				opts.Data = json

				s.NatsClient.Publish(&opts)
			}
		})

	default:
		if changeStream != nil {
			changeStream.Close(s.Context)
			changeStream = nil
		}
	}
}
