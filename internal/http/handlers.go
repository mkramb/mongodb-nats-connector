package http

import (
	"encoding/json"
	"net/http"
)

type health string

const (
	UP   health = "UP"
	DOWN health = "DOWN"
)

func (s *Server) registerRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)

	return mux
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]health)
	statusCode := http.StatusOK

	for _, healthCheck := range s.HeathChecks {
		if err := healthCheck.monitor(); err == nil {
			response[healthCheck.name] = UP
		} else {
			response[healthCheck.name] = DOWN
			statusCode = http.StatusServiceUnavailable
		}
	}

	if statusCode != http.StatusOK {
		s.Logger.Warn("Registered health checks are failing")
	}

	writeJson(w, statusCode, response)
}

func writeJson(writer http.ResponseWriter, code int, response any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)

	_ = json.NewEncoder(writer).Encode(response)
}
