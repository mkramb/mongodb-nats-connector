package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/pkg/connector"
)

func main() {
	ctx, shutdownServer := context.WithCancel(context.Background())

	log := logger.New()
	cfg := config.Options{
		Context: ctx,
		Logger:  log,
	}.New()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	connector := connector.Options{
		Context: ctx,
		Logger:  log,
		Config:  cfg,
	}.New()

	connector.StartHttp()
	connector.StartWatcher()

	<-shutdownSignal

	log.Info("Received shutdown signal")

	shutdownServer()
}
