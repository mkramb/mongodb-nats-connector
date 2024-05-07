package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
)

type Server struct{}

func NewServer(cfg *config.ConnectorConfig) *http.Server {
	NewServer := &Server{}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HttpConfig.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	nats.NewRaft(cfg.NatsConfig, NewServer.stateHandler)

	return server
}
