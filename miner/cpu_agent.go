package miner

import (
	"sync"
	"sync/atomic"

	"gbchain-org/go-gbchain/consensus"
	"gbchain-org/go-gbchain/core/types"
	"gbchain-org/go-gbchain/log"
)

type CpuAgent struct {
	mu sync.Mutex

	workCh        chan *types.Block
	stop          chan struct{}
	quitCurrentOp chan struct{}
	resultCh      chan<- *types.Block

	chain  consensus.ChainReader
	engine consensus.Engine

	isMining int32 // isMining indicates whether the agent is currently mining
}

func NewCpuAgent(chain consensus.ChainReader, engine consensus.Engine) *CpuAgent {
	miner := &CpuAgent{
		chain:  chain,
		engine: engine,
		stop:   make(chan struct{}, 1),
		workCh: make(chan *types.Block, 1),
	}
	return miner
}

func (self *CpuAgent) DispatchWork(block *types.Block) {
	self.workCh <- block
}

func (self *CpuAgent) SubscribeResult(ch chan<- *types.Block) {
	self.resultCh = ch
}

func (self *CpuAgent) Stop() {
	if !atomic.CompareAndSwapInt32(&self.isMining, 1, 0) {
		return // agent already stopped
	}
	self.stop <- struct{}{}
done:
	// Empty work channel
	for {
		select {
		case <-self.workCh:
		default:
			break done
		}
	}
	log.Info("CPU agent stopped")
}

func (self *CpuAgent) Start() {
	if !atomic.CompareAndSwapInt32(&self.isMining, 0, 1) {
		return // agent already started
	}
	go self.update()
}

func (self *CpuAgent) update() {
out:
	for {
		select {
		case work := <-self.workCh:
			self.mu.Lock()
			if self.quitCurrentOp != nil {
				close(self.quitCurrentOp)
			}
			self.quitCurrentOp = make(chan struct{})
			go self.mine(work, self.quitCurrentOp)
			self.mu.Unlock()
		case <-self.stop:
			self.mu.Lock()
			if self.quitCurrentOp != nil {
				close(self.quitCurrentOp)
				self.quitCurrentOp = nil
			}
			self.mu.Unlock()
			break out
		}
	}
}

func (self *CpuAgent) mine(work *types.Block, stop <-chan struct{}) {
	if err := self.engine.Seal(self.chain, work, self.resultCh, stop); err != nil {
		log.Error("Block sealing failed", "err", err)
	}
}

func (self *CpuAgent) GetHashRate() uint64 {
	if pow, ok := self.engine.(consensus.PoW); ok {
		return uint64(pow.Hashrate())
	}
	return 0
}
