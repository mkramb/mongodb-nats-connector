package raft

import (
	"github.com/nats-io/graft"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var changeStream *mongo.ChangeStream = nil

func (s *Server) stateHandler(stateTo graft.State) {
	switch stateTo {

	case graft.LEADER:
		s.Logger.Info("Becoming leader")

		changeStream = s.MongoClient.Watch()
		s.MongoClient.IterateChangeStream(changeStream, func(changeEvent bson.M) {
			s.Logger.Info("Received data", "fullDocument", changeEvent["fullDocument"])
		})

	default:
		if changeStream != nil {
			changeStream.Close(s.Context)
			changeStream = nil
		}
	}
}
