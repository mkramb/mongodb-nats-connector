package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
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
		s.Logger.Error("Error handling JSON marshal", slog.Any("err", err))
	}

	_, _ = w.Write(jsonResp)
}
