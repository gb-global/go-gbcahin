package retriever

import (
	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/params"

	"gbchain-org/go-gbchain/cross"
	"gbchain-org/go-gbchain/cross/trigger"
	"gbchain-org/go-gbchain/cross/trigger/simpletrigger"
)

type SimpleRetriever struct {
	*ChainInvoke
	*CreditValidator
	pm simpletrigger.ProtocolManager
}

func NewSimpleRetriever(bc simpletrigger.BlockChain, pm simpletrigger.ProtocolManager, contract common.Address,
	config *cross.Config, chainConfig *params.ChainConfig) trigger.ChainRetriever {
	r := new(SimpleRetriever)
	r.pm = pm
	r.ChainInvoke = NewChainInvoke(bc)
	r.CreditValidator = NewCreditleValidator(contract, bc, config, chainConfig)
	r.CreditValidator.SimpleRetriever = r
	return r
}

func (s *SimpleRetriever) CanAcceptTxs() bool {
	return s.pm.CanAcceptTxs()
}

func (s *SimpleRetriever) ConfirmedDepth() uint64 {
	return uint64(simpletrigger.DefaultConfirmDepth)
}
