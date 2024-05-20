package raft

import (
	"github.com/nats-io/graft"
	"go.mongodb.org/mongo-driver/mongo"
)

var changeStream *mongo.ChangeStream = nil

func (s *Server) stateHandler(stateFrom, stateTo graft.State) {
	switch stateTo {

	case graft.LEADER:
		s.Logger.Info("Becoming leader")
		changeStream = s.Mongo.StartWatch(s.Config.Mongo.WatchCollections, s.Config.Mongo.WatchOperations)

	case graft.FOLLOWER, graft.CANDIDATE:
		if stateFrom == graft.LEADER {
			s.Logger.Info("Becoming follower")
		}

		if changeStream != nil {
			changeStream.Close(s.Context)
			changeStream = nil
		}
	}
}
