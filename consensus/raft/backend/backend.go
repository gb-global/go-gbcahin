package backend

import (
	"crypto/ecdsa"
	"errors"
	"sync"
	"time"

	"gbchain-org/go-gbchain/accounts"
	"gbchain-org/go-gbchain/consensus/raft"
	"gbchain-org/go-gbchain/core"
	"gbchain-org/go-gbchain/eth/downloader"
	"gbchain-org/go-gbchain/ethdb"
	"gbchain-org/go-gbchain/event"
	"gbchain-org/go-gbchain/log"
	"gbchain-org/go-gbchain/miner"
	"gbchain-org/go-gbchain/node"
	"gbchain-org/go-gbchain/p2p"
	"gbchain-org/go-gbchain/p2p/enode"
	"gbchain-org/go-gbchain/rpc"
	"gbchain-org/go-gbchain/sub"
)

type RaftService struct {
	blockchain     *core.BlockChain
	chainDb        ethdb.Database // Block chain database
	txMu           sync.Mutex
	txPool         *core.TxPool
	accountManager *accounts.Manager
	downloader     *downloader.Downloader

	raftProtocolManager *ProtocolManager
	startPeers          []*enode.Node

	// we need an event mux to instantiate the blockchain
	eventMux *event.TypeMux
	minter   *miner.Miner
	nodeKey  *ecdsa.PrivateKey
}

func New(ctx *node.ServiceContext, raftId, raftPort uint16, joinExisting bool, e *sub.Ethereum, startPeers []*enode.Node, datadir string, blockTime time.Duration, useDns bool) (*RaftService, error) {
	service := &RaftService{
		eventMux:       ctx.EventMux,
		chainDb:        e.ChainDb(),
		blockchain:     e.BlockChain(),
		txPool:         e.TxPool(),
		accountManager: e.AccountManager(),
		downloader:     e.Downloader(),
		startPeers:     startPeers,
		nodeKey:        ctx.NodeKey(),
	}

	engine, ok := e.Engine().(*raft.Raft)
	if !ok {
		return nil, errors.New("raft service require raft engine") // never occur
	}
	engine.SetId(raftId)

	minerConfig := &e.Config().Miner
	// Reuse Recommit
	minerConfig.Recommit = blockTime

	// Configure the local mining address
	eb, _ := e.Etherbase()

	service.minter = miner.New(service, minerConfig, e.ChainConfig(), service.eventMux, engine, nil)
	service.minter.SetEtherbase(eb)

	var err error
	if service.raftProtocolManager, err = NewProtocolManager(raftId, raftPort, service.blockchain, service.eventMux, startPeers, joinExisting, datadir, service.minter, service.downloader, useDns); err != nil {
		return nil, err
	}

	return service, nil
}

// Backend interface methods:

func (service *RaftService) AccountManager() *accounts.Manager { return service.accountManager }
func (service *RaftService) BlockChain() *core.BlockChain      { return service.blockchain }
func (service *RaftService) ChainDb() ethdb.Database           { return service.chainDb }
func (service *RaftService) EventMux() *event.TypeMux          { return service.eventMux }
func (service *RaftService) TxPool() *core.TxPool              { return service.txPool }

// node.Service interface methods:

func (service *RaftService) Protocols() []p2p.Protocol { return []p2p.Protocol{} }
func (service *RaftService) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "raft",
			Version:   "1.0",
			Service:   NewPublicRaftAPI(service),
			Public:    true,
		},
	}
}

// Start implements node.Service, starting the background data propagation thread
// of the protocol.
func (service *RaftService) Start(p2pServer *p2p.Server) error {
	service.raftProtocolManager.Start(p2pServer)
	return nil
}

// Stop implements node.Service, stopping the background data propagation thread
// of the protocol.
func (service *RaftService) Stop() error {
	service.blockchain.Stop()
	service.raftProtocolManager.Stop()
	service.minter.Stop()
	service.eventMux.Stop()

	//service.chainDb.Close()

	log.Info("Raft stopped")
	return nil
}