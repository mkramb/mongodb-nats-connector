package raft

import (
	"fmt"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/nats-io/graft"
	"github.com/nats-io/nats.go"
)

func StartRaft(cfg *config.ConnectorConfig) {
	var (
		opts = &nats.DefaultOptions
		ci   = graft.ClusterInfo{Name: cfg.NatsConfig.ClusterName, Size: cfg.NatsConfig.ClusterSize}
	)

	opts.Url = cfg.NatsConfig.ServerUrl
	rpc, err := graft.NewNatsRpc(opts)

	if err != nil {
		panic(err)
	}

	var (
		errC         = make(chan error)
		stateChangeC = make(chan graft.StateChange)
		handler      = graft.NewChanHandler(stateChangeC, errC)
	)

	node, err := graft.New(ci, handler, rpc, cfg.NatsConfig.ClusterName)

	if err != nil {
		panic(err)
	}

	defer node.Close()

	stateHandler(node.State())

	for {
		select {
		case change := <-stateChangeC:
			stateHandler(change.To)
		case err := <-errC:
			fmt.Printf("Error: %s\n", err)
		}
	}
}
