package nats

import (
	"fmt"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/nats-io/graft"
	"github.com/nats-io/nats.go"
)

type StateHandler func(state graft.State)

func NewRaft(cfg *config.NatsConfig, handleState StateHandler) {
	var (
		opts = &nats.DefaultOptions
		ci   = graft.ClusterInfo{Name: cfg.ClusterName, Size: cfg.ClusterSize}
	)

	opts.Url = cfg.ServerUrl
	rpc, err := graft.NewNatsRpc(opts)

	if err != nil {
		panic(err)
	}

	var (
		errC         = make(chan error)
		stateChangeC = make(chan graft.StateChange)
		handler      = graft.NewChanHandler(stateChangeC, errC)
	)

	node, err := graft.New(ci, handler, rpc, cfg.LogPath)

	if err != nil {
		panic(err)
	}

	defer node.Close()

	handleState(node.State())

	for {
		select {
		case change := <-stateChangeC:
			handleState(change.To)
		case err := <-errC:
			fmt.Printf("Error: %s\n", err)
		}
	}
}
