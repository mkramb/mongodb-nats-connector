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

	nats := nats.Options{
		Context: ctx,
		Logger:  log,
		Config:  cfg.Nats,
	}.NewClient()

	mongo := mongo.Options{
		Context: ctx,
		Logger:  log,
		Config:  cfg.Mongo,
	}.NewClient()

	defer nats.Close()
	defer mongo.Close()

	http := http.Options{
		Context: ctx,
		Logger:  log,
		Config:  cfg.Http,
	}.NewServer()

	raft := raft.Options{
		Context: ctx,
		Logger:  log,
		Config:  cfg.Nats,
		Nats:    nats,
		Mongo:   mongo,
	}.NewServer()

	log.Info("Starting http server")
	log.Info("Starting raft server")

	go http.Start()
	go raft.Start()

	<-shutdownSignal

	log.Info("Received shutdown signal")
	shutdownServer()
}
