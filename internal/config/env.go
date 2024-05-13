package config

import (
	"context"
	"os"

	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Http *HttpConfig
	Nats *NatsConfig
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

func NewEnvConfig(log logger.Logger) *Config {
	ctx := context.Background()

	var c Config

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Error("Error parsing environment values", logger.AsError(err))
		os.Exit(1)
	}

	return &c
}
