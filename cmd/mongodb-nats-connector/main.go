package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/http"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/mkramb/mongodb-nats-connector/internal/raft"
)

func main() {
	log := logger.NewJSONLogger()
	cfg := config.NewEnvConfig(log)

	nats := nats.InitNats(log, cfg.Nats.ServerUrl)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	http := http.ServerHttp{Config: cfg, Logger: log}
	raft := raft.ServerRaft{
		Nats:   nats,
		Config: cfg,
		Logger: log,
	}

	go raft.StartRaft()
	go http.StartHttp()

	<-shutdown
}
