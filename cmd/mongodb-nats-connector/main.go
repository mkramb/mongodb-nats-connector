package main

import (
	"sync"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/http"
	"github.com/mkramb/mongodb-nats-connector/internal/raft"
)

func main() {
	cfg := config.NewConfig()
	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		http.StartHttp(cfg)
	}()

	go func() {
		defer wg.Done()
		raft.StartRaft(cfg)
	}()

	wg.Wait()
}
