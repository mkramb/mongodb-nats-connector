package raft

import (
	"context"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/internal/mongo"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/nats-io/graft"
)

type Options struct {
	Context context.Context
	Logger  logger.Logger
	Config  *config.NatsConfig
	Nats    *nats.Client
	Mongo   *mongo.Client
}

type Server struct {
	Cluster *graft.ClusterInfo
	Options
}

func (o Options) NewServer() *Server {
	cluster := &graft.ClusterInfo{Name: o.Config.ClusterName, Size: o.Config.ClusterSize}

	return &Server{
		Cluster: cluster,
		Options: o,
	}
}

func (s *Server) Start() {
	natsRpc, err := graft.NewNatsRpcFromConn(s.Nats.Conn)

	if err != nil {
		s.Logger.Error("Error starting RAFT connection", logger.AsError(err))
		return
	}

	var (
		errC         = make(chan error)
		stateChangeC = make(chan graft.StateChange)
		handler      = graft.NewChanHandler(stateChangeC, errC)
	)

	node, err := graft.New(*s.Cluster, handler, natsRpc, s.Config.ClusterName)

	if err != nil {
		s.Logger.Error("Error starting new RAFT node", logger.AsError(err))
	}

	defer node.Close()
	defer natsRpc.Close()
	defer s.Logger.Info("Closing raft connection")

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
