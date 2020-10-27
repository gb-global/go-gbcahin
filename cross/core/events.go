package core

import (
	"math/big"

	"gbchain-org/go-gbchain/common"
)

type ConfirmedMakerEvent struct {
	Txs []*CrossTransaction
}

type NewTakerEvent struct {
	Takers []*ReceptTransaction
}

type ConfirmedTakerEvent struct {
	Txs []*ReceptTransaction
}

type SignedCtxEvent struct { // pool event
	Txs      []*CrossTransactionWithSignatures
	CallBack func([]CommitEvent)
}

type CommitEvent struct {
	Tx              *CrossTransactionWithSignatures
	InvalidSigIndex []int
}

type NewFinishEvent struct {
	Finishes []*CrossTransactionModifier
}

type ConfirmedFinishEvent struct {
	Finishes []*CrossTransactionModifier
}

type NewAnchorEvent struct {
	ChainInfo []*RemoteChainInfo
}

type ModType uint8

const (
	Normal = ModType(iota)
	Remote
	Reorg
)

func (t ModType) String() string {
	switch t {
	case Normal:
		return "normal"
	case Remote:
		return "remote"
	case Reorg:
		return "reorg"
	default:
		return "unknown"
	}
}

type CrossTransactionModifier struct {
	Type          ModType
	ID            common.Hash
	Status        CtxStatus
	AtBlockNumber uint64
}

type CrossBlockEvent struct {
	Number          *big.Int
	ConfirmedMaker  ConfirmedMakerEvent
	NewTaker        NewTakerEvent
	ConfirmedTaker  ConfirmedTakerEvent
	NewFinish       NewFinishEvent
	ConfirmedFinish ConfirmedFinishEvent
	NewAnchor       NewAnchorEvent
	ReorgTaker      NewTakerEvent
	ReorgFinish     NewFinishEvent
}

func (e CrossBlockEvent) IsEmpty() bool {
	return len(e.ConfirmedMaker.Txs)|len(e.ConfirmedTaker.Txs)|
		len(e.ConfirmedFinish.Finishes)|len(e.NewTaker.Takers)|
		len(e.NewFinish.Finishes)|len(e.NewAnchor.ChainInfo)|
		len(e.ReorgTaker.Takers)|len(e.ReorgFinish.Finishes) == 0
}
