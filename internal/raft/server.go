package raft

import (
	"context"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/internal/mongo"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/nats-io/graft"
)

type Server struct {
	Nats    *nats.Client
	Cluster *graft.ClusterInfo
	Mongo   *mongo.Client
	Config  *config.Config
	Logger  logger.Logger
	Context context.Context
}

func NewServer(ctx context.Context, log logger.Logger, cfg *config.Config, nats *nats.Client, mongo *mongo.Client) *Server {
	cluster := &graft.ClusterInfo{Name: cfg.Nats.ClusterName, Size: cfg.Nats.ClusterSize}

	return &Server{
		Nats:    nats,
		Cluster: cluster,
		Mongo:   mongo,
		Config:  cfg,
		Logger:  log,
		Context: ctx,
	}
}

func (s *Server) StartRaft() {
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

	node, err := graft.New(*s.Cluster, handler, natsRpc, s.Config.Nats.ClusterName)

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
			}

			s.stateHandler(change.To)

		case err := <-errC:
			s.Logger.Error("Error processing raft state", logger.AsError(err))
			return

		case <-s.Context.Done():
			return
		}
	}
}
