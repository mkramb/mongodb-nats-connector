package config

import (
	"context"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/sethvargo/go-envconfig"
)

type Options struct {
	Context context.Context
	Logger  logger.Logger
}

type Config struct {
	Http  *HttpConfig
	Nats  *NatsConfig
	Mongo *MongoConfig
	Raft  *RaftConfig
}

type HttpConfig struct {
	Port int `env:"HTTP_PORT, default=3000"`
}

type NatsConfig struct {
	ServerUrl  string `env:"NATS_SERVER_URL, required"`
	StreamName string `env:"NATS_STREAM_NAME, default=cs"`
}

type RaftConfig struct {
	ClusterSize int    `env:"RAFT_CLUSTER_SIZE, default=3"`
	ClusterName string `env:"RAFT_CLUSTER_NAME, default=connector"`
	LogPath     string `env:"RAFT_LOG_PATH, default=./graft.log"`
}

type MongoConfig struct {
	ServerUri        string   `env:"MONGO_URI, required"`
	WatchCollections []string `env:"MONGO_WATCH_COLLECTIONS, required"`
	WatchOperations  []string `env:"MONGO_WATCH_OPERATIONS, default=insert,update,replace"`
}

func (o Options) New() *Config {
	var c Config

	if err := envconfig.Process(o.Context, &c); err != nil {
		o.Logger.Error("Error parsing environment values", logger.AsError(err))
		panic("Error parsing environment values")
	}

	return &c
}
