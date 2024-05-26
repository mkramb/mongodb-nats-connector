package watcher

import (
	"context"
	"fmt"

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

func (s *Server) Start() {
	if s.Config.ClusterSize > 1 {
		s.StartRaft()
	} else {
		s.Logger.Info("Running watcher as a single instance, disabling raft")

		go func() {
			s.watchForChangeEvents()
		}()

		defer s.Logger.Info("Closing watcher server")
		defer s.MongoClient.StopChangeStream()

		<-s.Context.Done()
	}
}

func (s *Server) watchForChangeEvents() {
	s.MongoClient.StartChangeStream(func(event *mongo.ChangeEvent, json []byte) {
		var opts nats.PublishOptions

		opts.MsgId = event.ResumeToken.Value
		opts.Subject = fmt.Sprintf("%v.%v.%v", event.Ns.Coll, event.OperationType, event.FullDocument.Id.Value)
		opts.Data = json

		s.NatsClient.PublishEvent(&opts)
	})
}
