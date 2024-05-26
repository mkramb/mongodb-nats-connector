package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
)

type HealthMonitor func() error

type HealthCheck struct {
	name    string
	monitor HealthMonitor
}

type Options struct {
	Context context.Context
	Logger  logger.Logger
	Config  *config.HttpConfig
}

type Server struct {
	Http        *http.Server
	HeathChecks []HealthCheck
	Options
}

func (o Options) New() *Server {
	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", o.Config.Port),
		IdleTimeout: time.Minute,
	}

	return &Server{
		Http:        httpServer,
		HeathChecks: make([]HealthCheck, 0),
		Options:     o,
	}
}

func (s *Server) RegisterHealthCheck(name string, monitor HealthMonitor) {
	s.HeathChecks = append(s.HeathChecks, HealthCheck{
		name:    name,
		monitor: monitor,
	})
}

func (s *Server) Start() {
	s.Http.Handler = s.registerRoutes()

	go func() {
		if err := s.Http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("Error starting http server", logger.AsError(err))
		}
	}()

	defer s.Logger.Info("Closing http client")
	defer s.Http.Shutdown(s.Context)

	<-s.Context.Done()
}
