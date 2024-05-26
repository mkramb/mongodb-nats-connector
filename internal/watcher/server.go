package watcher

import (
	"context"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/internal/mongo"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/nats-io/graft"
)

type Options struct {
	Context     context.Context
	Logger      logger.Logger
	Config      *config.RaftConfig
	NatsClient  *nats.Client
	MongoClient *mongo.Client
}

type Server struct {
	Cluster *graft.ClusterInfo
	Options
}

func (o Options) New() *Server {
	cluster := &graft.ClusterInfo{Name: o.Config.ClusterName, Size: o.Config.ClusterSize}

	return &Server{
		Cluster: cluster,
		Options: o,
	}
}

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

func (s *Server) Start() {
	go func() {
		s.startWatching()
	}()

	defer s.Logger.Info("Closing watcher server")
	defer s.MongoClient.StopWatcher()

	<-s.Context.Done()
}
