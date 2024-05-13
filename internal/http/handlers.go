package http

import (
	"encoding/json"
	"net/http"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
)

func (s *ServerHttp) registerRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)

	return mux
}

func (s *ServerHttp) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(map[string]string{
		"message": "It's healthy",
	})

	if err != nil {
		s.Logger.Error("Error handling JSON marshal", logger.AsError(err))
	}

	_, _ = w.Write(jsonResp)
}
