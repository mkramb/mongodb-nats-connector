package watcher

import (
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/nats-io/graft"
)

func (s *Server) StartRaft() {
	rpc, err := graft.NewNatsRpcFromConn(s.NatsClient.Conn)

	if err != nil {
		s.Logger.Error("Error starting raft connection", logger.AsError(err))
		return
	}

	var (
		errC         = make(chan error)
		stateChangeC = make(chan graft.StateChange)
		handler      = graft.NewChanHandler(stateChangeC, errC)
	)

	node, err := graft.New(*s.Cluster, handler, rpc, s.Config.LogPath)

	if err != nil {
		s.Logger.Error("Error starting new raft node", logger.AsError(err))
	}

	defer s.Logger.Info("Closing watcher server")
	defer node.Close()
	defer rpc.Close()

	s.stateHandler(node.State())

	for {
		select {

		case change := <-stateChangeC:
			if change.To == graft.CLOSED {
				s.Logger.Info("RAFT connection is closed")
				return
			} else {
				s.stateHandler(change.To)
			}

		case err := <-errC:
			s.Logger.Error("Error processing raft state", logger.AsError(err))
			return

		case <-s.Context.Done():
			return
		}
	}
}

func (s *Server) stateHandler(stateTo graft.State) {
	switch stateTo {

	case graft.LEADER:
		s.Logger.Info("Becoming leader")
		s.watchForChangeEvents()

	default:
		s.Logger.Info("Becoming follower")
		s.MongoClient.StopChangeStream()
	}
}
