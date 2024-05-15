package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
)

type Server struct {
	Config *config.Config
	Logger logger.Logger
}

func (s *Server) StartHttp() {
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", s.Config.Http.Port),
		Handler:     s.registerRoutes(),
		IdleTimeout: time.Minute,
	}

	err := server.ListenAndServe()

	if err != nil {
		s.Logger.Error("Error starting http server", logger.AsError(err))
		os.Exit(1)
	}
}
