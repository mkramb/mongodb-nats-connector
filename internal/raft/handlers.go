package raft

import (
	"context"

	"github.com/nats-io/graft"
	"go.mongodb.org/mongo-driver/mongo"
)

var changeStream *mongo.ChangeStream = nil

func (s *Server) stateHandler(state graft.State) {
	switch state {

	case graft.LEADER:
		s.Logger.Info("Becoming leader, starting watcher")
		changeStream = s.Mongo.Watch(s.Config.Mongo.WatchDatabase, s.Config.Mongo.WatchCollections, s.Config.Mongo.WatchOperations)

	case graft.FOLLOWER, graft.CANDIDATE:
		if changeStream != nil {
			changeStream.Close(context.TODO())
			changeStream = nil
		}

	case graft.CLOSED:
		return

	}
}
