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
	ctx, shutdownServer := context.WithCancel(context.Background())

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	log := logger.NewLogger()
	cfg := config.NewEnvConfig(ctx, log)

	nats := nats.InitClient(ctx, log, cfg.Nats.ServerUrl)
	mongo := mongo.InitClient(ctx, log, cfg.Mongo.ServerUri)

	defer nats.Close()
	defer mongo.Close()

	http := http.NewServer(ctx, log, cfg)
	raft := raft.NewServer(ctx, log, cfg, nats, mongo)

	log.Info("Starting http server")
	log.Info("Starting raft server")

	go http.StartHttp()
	go raft.StartRaft()

	<-shutdownSignal

	log.Info("Received shutdown signal")
	shutdownServer()
}
