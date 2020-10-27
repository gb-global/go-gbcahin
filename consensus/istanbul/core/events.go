package core

import (
	"gbchain-org/go-gbchain/consensus/istanbul"
)

type backlogEvent struct {
	src istanbul.Validator
	msg *message
}

type timeoutEvent struct{}
