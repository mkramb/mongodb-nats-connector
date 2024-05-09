package nats

import (
	"fmt"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/nats-io/graft"
)

func StartRaft(cfg *config.ConnectorConfig, natsConnection *NatsClient) {
	cluster := graft.ClusterInfo{Name: cfg.NatsConfig.ClusterName, Size: cfg.NatsConfig.ClusterSize}
	rpc, err := graft.NewNatsRpcFromConn(natsConnection.Conn)

	if err != nil {
		panic(fmt.Sprintf("Error starting graft: %s", err))
	}

	var (
		errC         = make(chan error)
		stateChangeC = make(chan graft.StateChange)
		handler      = graft.NewChanHandler(stateChangeC, errC)
	)

	node, err := graft.New(cluster, handler, rpc, cfg.NatsConfig.ClusterName)

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
