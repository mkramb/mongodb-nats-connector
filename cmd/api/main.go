package main

import (
	"fmt"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/server"
)

func main() {
	cfg := config.NewConfig()

	server := server.NewServer(cfg)
	err := server.ListenAndServe()

	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
