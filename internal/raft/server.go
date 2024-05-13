package raft

import (
	"fmt"

	"github.com/mkramb/mongodb-nats-connector/internal/config"
	"github.com/mkramb/mongodb-nats-connector/internal/logger"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
	"github.com/nats-io/graft"
)

type ServerRaft struct {
	Nats   *nats.Nats
	Config *config.Config
	Logger logger.Logger
}

func (s *ServerRaft) StartRaft() {
	cluster := graft.ClusterInfo{Name: s.Config.Nats.ClusterName, Size: s.Config.Nats.ClusterSize}
	rpc, err := graft.NewNatsRpcFromConn(s.Nats.Client.Conn)

	if err != nil {
		panic(fmt.Sprintf("Error starting graft: %s", err))
	}

	var (
		errC         = make(chan error)
		stateChangeC = make(chan graft.StateChange)
		handler      = graft.NewChanHandler(stateChangeC, errC)
	)

	node, err := graft.New(cluster, handler, rpc, s.Config.Nats.ClusterName)

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
