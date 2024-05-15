package raft

import (
	"os"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/internal/mongo"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/nats-io/graft"
)

type Server struct {
	Nats   *nats.Client
	Mongo  *mongo.Client
	Config *config.Config
	Logger logger.Logger
}

func (s *Server) StartRaft() {
	cluster := graft.ClusterInfo{Name: s.Config.Nats.ClusterName, Size: s.Config.Nats.ClusterSize}
	rpc, err := graft.NewNatsRpcFromConn(s.Nats.Conn)

	if err != nil {
		s.Logger.Error("Error starting raft", logger.AsError(err))
		os.Exit(1)
	}

	var (
		errC         = make(chan error)
		stateChangeC = make(chan graft.StateChange)
		handler      = graft.NewChanHandler(stateChangeC, errC)
	)

	node, err := graft.New(cluster, handler, rpc, s.Config.Nats.ClusterName)

	if err != nil {
		s.Logger.Error("Error starting new raft node", logger.AsError(err))
		os.Exit(1)
	}

	defer node.Close()

	s.stateHandler(node.State())

	for {
		select {
		case change := <-stateChangeC:
			s.stateHandler(change.To)
		case err := <-errC:
			s.Logger.Error("Error processing raft state", logger.AsError(err))
			os.Exit(1)
		}
	}
}
