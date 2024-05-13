package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
)

type ServerHttp struct {
	Config *config.Config
	Logger logger.Logger
}

func (s *ServerHttp) StartHttp() {
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", s.Config.Http.Port),
		Handler:     registerRoutes(),
		IdleTimeout: time.Minute,
	}

	err := server.ListenAndServe()

	if err != nil {
		s.Logger.Error("Error starting http server", err)
		os.Exit(1)
	}
}

func registerRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	return mux
}