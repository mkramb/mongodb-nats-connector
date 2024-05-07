package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nats-io/graft"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)

	return mux
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(map[string]string{
		"message": "It's healthy",
	})

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) stateHandler(state graft.State) {
	switch state {
	case graft.LEADER:
		fmt.Println("***Becoming leader***")
	case graft.FOLLOWER:
		fmt.Println("***Becoming follower***")
	case graft.CANDIDATE:
		fmt.Println("***Becoming candidate***")
	case graft.CLOSED:
		return
	default:
		panic(fmt.Sprintf("Unknown state: %s", state))
	}
}
