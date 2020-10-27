package utils

import (
	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/cross"
	crossBackend "gbchain-org/go-gbchain/cross/backend"
	crossdb "gbchain-org/go-gbchain/cross/database"
	"gbchain-org/go-gbchain/cross/trigger/simpletrigger"
	"gbchain-org/go-gbchain/cross/trigger/simpletrigger/executor"
	"gbchain-org/go-gbchain/cross/trigger/simpletrigger/retriever"
	"gbchain-org/go-gbchain/cross/trigger/simpletrigger/subscriber"
	"gbchain-org/go-gbchain/eth"
	"gbchain-org/go-gbchain/node"
	"gbchain-org/go-gbchain/sub"
)

func RegisterCrossChainService(stack *node.Node, cfg cross.Config, mainCh chan *eth.Ethereum, subCh chan *sub.Ethereum) {
	err := stack.Register(func(sc *node.ServiceContext) (node.Service, error) {
		mainNode := <-mainCh
		subNode := <-subCh
		defer close(mainCh)
		defer close(subCh)
		mainCtx, err := newCreditChainContext(sc, mainNode, cfg, cfg.MainContract, "mainChain_unconfirmed.rlp", "mainChain_queue")
		if err != nil {
			return nil, err
		}
		subCtx, err := newCreditChainContext(sc, subNode, cfg, cfg.SubContract, "subChain_unconfirmed.rlp", "subChain_queue")
		if err != nil {
			return nil, err
		}
		return crossBackend.NewCrossService(sc, mainCtx, subCtx, cfg)
	})
	if err != nil {
		Fatalf("Failed to register the CrossChain service: %v", err)
	}
}

func newCreditChainContext(node *node.ServiceContext, chain simpletrigger.GBChain, config cross.Config,
	contract common.Address, journal string, queue string) (ctx *cross.ServiceContext, err error) {
	edb, err := crossdb.OpenEtherDB(node, queue)
	if err != nil {
		return nil, err
	}
	qdb, err := crossdb.NewQueueDB(edb)
	if err != nil {
		return nil, err
	}

	ctx = &cross.ServiceContext{ProtocolChain: simpletrigger.NewSimpleProtocolChain(chain), Config: &config}
	ctx.Executor, err = executor.NewSimpleExecutor(chain, config.Signer, contract, qdb)
	if err != nil {
		return nil, err
	}
	ctx.Retriever = retriever.NewSimpleRetriever(chain.BlockChain(), chain.ProtocolManager(), contract, ctx.Config, chain.ChainConfig())
	ctx.Subscriber = subscriber.NewSimpleSubscriber(contract, chain.BlockChain(), node.ResolvePath(journal))
	return ctx, nil
}
