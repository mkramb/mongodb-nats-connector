package config

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type ConnectorConfig struct {
	HttpConfig *HttpConfig
	NatsConfig *NatsConfig
}

type HttpConfig struct {
	Port int `env:"HTTP_PORT, default=3000"`
}

type NatsConfig struct {
	ServerUrl   string `env:"NATS_SERVER_URL, required"`
	ClusterSize int    `env:"NATS_CLUSTER_SIZE, default=2"`
	ClusterName string `env:"NATS_CLUSTER_NAME, default=connector"`
	LogPath     string `env:"NATS_LOG_PATH, default=/tmp/graft.log"`
}

func NewConfig() *ConnectorConfig {
	ctx := context.Background()

	var c ConnectorConfig

	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	return &c
}
