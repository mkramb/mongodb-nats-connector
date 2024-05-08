package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/http"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
)

func main() {
	cfg := config.NewConfig()

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	go http.StartHttp(cfg)
	go nats.StartRaft(cfg)

	<-gracefulShutdown
}
