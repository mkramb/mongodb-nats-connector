package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
)

func StartHttp(cfg *config.ConnectorConfig) {
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.HttpConfig.Port),
		Handler:     registerRoutes(),
		IdleTimeout: time.Minute,
	}

	err := server.ListenAndServe()

	if err != nil {
		panic(fmt.Sprintf("Cannot start server: %s", err))
	}
}

func registerRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	return mux
}
