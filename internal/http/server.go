package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
)

type Server struct {
	Http    *http.Server
	Config  *config.Config
	Logger  logger.Logger
	Context context.Context
}

func NewServer(ctx context.Context, log logger.Logger, cfg *config.Config) *Server {
	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Http.Port),
		IdleTimeout: time.Minute,
	}

	return &Server{
		Http:    httpServer,
		Config:  cfg,
		Logger:  log,
		Context: ctx,
	}
}

func (s *Server) StartHttp() {
	s.Http.Handler = s.registerRoutes()

	go func() {
		if err := s.Http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("Error starting http server", logger.AsError(err))
		}
	}()

	defer s.Http.Shutdown(s.Context)
	defer s.Logger.Info("Closing http client")

	<-s.Context.Done()
}
