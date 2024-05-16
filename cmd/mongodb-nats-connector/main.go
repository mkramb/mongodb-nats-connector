package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/http"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/internal/mongo"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/mkramb/mongodb-nats-connector/internal/raft"
)

func main() {
	log := logger.NewJSONLogger()
	cfg := config.NewEnvConfig(log)

	nats := nats.InitClient(log, cfg.Nats.ServerUrl)
	mongo := mongo.InitClient(log, cfg.Mongo.ServerUri)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	http := http.Server{Config: cfg, Logger: log}
	raft := raft.Server{
		Nats:   nats,
		Mongo:  mongo,
		Config: cfg,
		Logger: log,
	}

	go raft.StartRaft()
	go http.StartHttp()

	defer nats.Conn.Close()
	defer mongo.Conn.Disconnect(context.TODO())

	<-shutdown
}
