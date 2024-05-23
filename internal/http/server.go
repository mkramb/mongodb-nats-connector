package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
)

type Options struct {
	Context context.Context
	Logger  logger.Logger
	Config  *config.HttpConfig
}

type Server struct {
	Http *http.Server
	Options
}

func (o Options) New() *Server {
	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", o.Config.Port),
		IdleTimeout: time.Minute,
	}

	return &Server{
		Http:    httpServer,
		Options: o,
	}
}

func (s *Server) Start() {
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
