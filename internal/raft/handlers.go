package raft

import (
	"github.com/nats-io/graft"
	"go.mongodb.org/mongo-driver/mongo"
)

var changeStream *mongo.ChangeStream = nil

func (s *Server) stateHandler(stateTo graft.State) {
	switch stateTo {

	case graft.LEADER:
		s.Logger.Info("Becoming leader")

		changeStream = s.Mongo.Watch()
		s.Mongo.IterateChangeStream(changeStream, func(data []byte) {
			s.Logger.Info("Received data", "data", string(data))
		})

	default:
		if changeStream != nil {
			changeStream.Close(s.Context)
			changeStream = nil
		}
	}
}
