package raft

import (
	"github.com/nats-io/graft"
)

func (s *Server) stateHandler(state graft.State) {
	switch state {
	case graft.LEADER:
		s.Logger.Info("***Becoming leader***")
	case graft.FOLLOWER, graft.CANDIDATE:
		s.Logger.Info("***Becoming follower***")
	case graft.CLOSED:
		return
	}
}
