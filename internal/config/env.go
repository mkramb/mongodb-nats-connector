package config

import (
	"context"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Http  *HttpConfig
	Nats  *NatsConfig
	Mongo *MongoConfig
}

type HttpConfig struct {
	Port int `env:"HTTP_PORT, default=3000"`
}

type NatsConfig struct {
	ServerUrl   string `env:"NATS_SERVER_URL, required"`
	ClusterSize int    `env:"NATS_CLUSTER_SIZE, default=3"`
	ClusterName string `env:"NATS_CLUSTER_NAME, default=connector"`
	LogPath     string `env:"NATS_LOG_PATH, default=/tmp/graft.log"`
}

type MongoConfig struct {
	ServerUri        string   `env:"MONGO_URI, required"`
	WatchCollections []string `env:"MONGO_WATCH_COLLECTIONS, required"`
	WatchOperations  []string `env:"MONGO_WATCH_OPERATIONS, default=insert,update,replace"`
}

func NewEnvConfig(ctx context.Context, log logger.Logger) *Config {
	var c Config

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Error("Error parsing environment values", logger.AsError(err))
		panic("Error parsing environment values")
	}

	return &c
}
