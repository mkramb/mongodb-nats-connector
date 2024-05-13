package raft

import (
	"github.com/nats-io/graft"
)

func (s *ServerRaft) stateHandler(state graft.State) {
	switch state {
	case graft.LEADER:
		s.Logger.Info("***Becoming leader***")
	case graft.FOLLOWER:
		s.Logger.Info("***Becoming follower***")
	case graft.CANDIDATE:
		s.Logger.Info("***Becoming candidate***")
	case graft.CLOSED:
		return
	}
}
