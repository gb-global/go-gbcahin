package trigger

import (
	"math/big"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/cross/core"
	"gbchain-org/go-gbchain/event"
)

// Subscriber subscriber block logs, send them to crosschain service
type Subscriber interface {
	SubscribeBlockEvent(ch chan<- core.CrossBlockEvent) event.Subscription
	Stop()
}

// Executor execute transactions on blockchain
type Executor interface {
	SignHash([]byte) ([]byte, error)
	SubmitTransaction([]*core.ReceptTransaction)
	Start()
	Stop()
}

// Validator validate cross transaction on blockchain, check tx signer on contract
type Validator interface {
	VerifyExpire(ctx *core.CrossTransaction) error
	VerifyContract(cws Transaction) error
	//VerifyReorg(ctx Transaction) error
	VerifySigner(ctx *core.CrossTransaction, signChain, storeChainID *big.Int) (common.Address, error)
	UpdateAnchors(info *core.RemoteChainInfo) error
	RequireSignatures() int
	ExpireNumber() int // return -1 if never expired
}

type Transaction interface {
	ID() common.Hash
	ChainId() *big.Int
	DestinationId() *big.Int
	BlockHash() common.Hash
	From() common.Address
}

// ChainRetriever include Validator and provides blockchain retriever
type ChainRetriever interface {
	Validator
	CanAcceptTxs() bool
	ConfirmedDepth() uint64
	CurrentBlockNumber() uint64
	GetTransactionTimeOnChain(Transaction) uint64
	GetTransactionNumberOnChain(Transaction) uint64
	GetConfirmedTransactionNumberOnChain(Transaction) uint64
}
