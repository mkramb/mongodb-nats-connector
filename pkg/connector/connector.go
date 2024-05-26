package connector

import (
	"context"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/http"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/internal/mongo"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/mkramb/mongodb-nats-connector/internal/watcher"
)

type Options struct {
	Context context.Context
	Logger  logger.Logger
	Config  *config.Config
}

type Connector struct {
	NatsClient  *nats.Client
	MongoClient *mongo.Client
	Options
}

func (o Options) New() *Connector {
	natsClient := nats.Options{
		Context: o.Context,
		Logger:  o.Logger,
		Config:  o.Config.Nats,
	}.New()

	mongoClient := mongo.Options{
		Context: o.Context,
		Logger:  o.Logger,
		Config:  o.Config.Mongo,
	}.New()

	go func() {
		defer natsClient.Close()
		defer mongoClient.Close()

		<-o.Context.Done()
	}()

	return &Connector{
		NatsClient:  natsClient,
		MongoClient: mongoClient,
		Options:     o,
	}
}

func (c *Connector) StartHttp() {
	c.Logger.Info("Starting http server")

	httpServer := http.Options{
		Context: c.Context,
		Logger:  c.Logger,
		Config:  c.Config.Http,
	}.New()

	httpServer.RegisterHealthCheck("nats", c.NatsClient.Monitor)
	httpServer.RegisterHealthCheck("mongo", c.MongoClient.Monitor)

	go httpServer.Start()
}

func (c *Connector) StartWatcher() {
	c.Logger.Info("Starting watcher server")

	watcherServer := watcher.Options{
		Context:     c.Context,
		Logger:      c.Logger,
		Config:      c.Config.Raft,
		NatsClient:  c.NatsClient,
		MongoClient: c.MongoClient,
	}.New()

	if c.Config.Raft.ClusterSize > 1 {
		c.Logger.Info("Starting watcher server using raft")
		go watcherServer.StartRaft()
	} else {
		c.Logger.Info("Starting watcher server")
		go watcherServer.Start()
	}
}
