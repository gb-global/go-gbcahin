package cross

import (
	"math/big"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/rpc"

	"gbchain-org/go-gbchain/cross/trigger"
)

type ProtocolChain interface {
	ChainID() *big.Int
	GenesisHash() common.Hash
	RegisterAPIs([]rpc.API) //TODO: 改成由backend自己注册API
}

type ServiceContext struct {
	Config        *Config
	ProtocolChain ProtocolChain
	Subscriber    trigger.Subscriber
	Retriever     trigger.ChainRetriever
	Executor      trigger.Executor
}
