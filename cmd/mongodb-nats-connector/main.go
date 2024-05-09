package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/http"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
)

func main() {
	cfg := config.NewConfig()
	natsClient, err := nats.NewNatsClient(cfg.NatsConfig.ServerUrl)

	if err != nil {
		panic(fmt.Sprintf("Error connecting to nats: %s", err))
	}

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	go nats.StartRaft(cfg, natsClient)
	go http.StartHttp(cfg)

	<-gracefulShutdown
}
