package core

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/consensus/istanbul"
	"gbchain-org/go-gbchain/core/types"
	elog "gbchain-org/go-gbchain/log"
)

func makeBlock(number int64) *types.Block {
	header := &types.Header{
		Difficulty: big.NewInt(0),
		Number:     big.NewInt(number),
		GasLimit:   0,
		GasUsed:    0,
		Time:       0,
	}
	block := &types.Block{}
	return block.WithSeal(header)
}

func newTestProposal() istanbul.Proposal {
	return makeBlock(1)
}

func TestNewRequest(t *testing.T) {
	testLogger.SetHandler(elog.StdoutHandler)

	N := uint64(4)
	F := uint64(1)

	sys := NewTestSystemWithBackend(N, F)

	close := sys.Run(true)
	defer close()

	request1 := makeBlock(1)
	sys.backends[0].NewRequest(request1)

	<-time.After(1 * time.Second)

	request2 := makeBlock(2)
	sys.backends[0].NewRequest(request2)

	<-time.After(1 * time.Second)

	for _, backend := range sys.backends {
		if len(backend.committedMsgs) != 2 {
			t.Errorf("the number of executed requests mismatch: have %v, want 2", len(backend.committedMsgs))
		}
		if !reflect.DeepEqual(request1.Number(), backend.committedMsgs[0].commitProposal.Number()) {
			t.Errorf("the number of requests mismatch: have %v, want %v", request1.Number(), backend.committedMsgs[0].commitProposal.Number())
		}
		if !reflect.DeepEqual(request2.Number(), backend.committedMsgs[1].commitProposal.Number()) {
			t.Errorf("the number of requests mismatch: have %v, want %v", request2.Number(), backend.committedMsgs[1].commitProposal.Number())
		}
	}
}

func TestQuorumSize(t *testing.T) {
	N := uint64(4)
	F := uint64(1)

	sys := NewTestSystemWithBackend(N, F)
	backend := sys.backends[0]
	c := backend.engine.(*core)

	valSet := c.valSet
	for i := 1; i <= 1000; i++ {
		valSet.AddValidator(common.BytesToAddress([]byte(string(i))))
		if 2*c.Confirmations() <= (valSet.Size()+valSet.F()) || 2*c.Confirmations() > (valSet.Size()+valSet.F()+2) {
			t.Errorf("quorumSize constraint failed, expected value (2*Confirmations > Size+F && 2*Confirmations <= Size+F+2) to be:%v, got: %v, for size: %v", true, false, valSet.Size())
		}
	}
}
