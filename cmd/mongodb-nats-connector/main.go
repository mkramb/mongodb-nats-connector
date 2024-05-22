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

	log := logger.NewLogger()
	cfg := config.Options{
		Context: ctx,
		Logger:  log,
	}.NewEnvConfig()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	natsClient := nats.Options{
		Context: ctx,
		Logger:  log,
		Config:  cfg.Nats,
	}.NewClient()

	mongoClient := mongo.Options{
		Context: ctx,
		Logger:  log,
		Config:  cfg.Mongo,
	}.NewClient()

	defer natsClient.Close()
	defer mongoClient.Close()

	httpServer := http.Options{
		Context: ctx,
		Logger:  log,
		Config:  cfg.Http,
	}.NewServer()

	raftServer := raft.Options{
		Context:     ctx,
		Logger:      log,
		Config:      cfg.Nats,
		NatsClient:  natsClient,
		MongoClient: mongoClient,
	}.NewServer()

	log.Info("Starting http server")
	log.Info("Starting raft server")

	go httpServer.Start()
	go raftServer.Start()

	<-shutdownSignal

	log.Info("Received shutdown signal")
	shutdownServer()
}
