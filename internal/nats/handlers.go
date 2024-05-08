package nats

import (
	"fmt"

	"github.com/nats-io/graft"
)

func stateHandler(state graft.State) {
	switch state {
	case graft.LEADER:
		fmt.Println("***Becoming leader***")
	case graft.FOLLOWER:
		fmt.Println("***Becoming follower***")
	case graft.CANDIDATE:
		fmt.Println("***Becoming candidate***")
	case graft.CLOSED:
		return
	default:
		panic(fmt.Sprintf("Unknown state: %s", state))
	}
}
