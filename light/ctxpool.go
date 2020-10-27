package light

import (
	"sync"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/core/types"
	cc "gbchain-org/go-gbchain/cross/core"
	"gbchain-org/go-gbchain/params"
)

type CtxPool struct {
	config  *params.ChainConfig
	signer  types.Signer
	mu      sync.RWMutex
	chain   *LightChain
	pending map[common.Hash]*cc.CrossTransactionWithSignatures
}

// NewCtxPool creates a new light cross transaction pool
func NewCtxPool(config *params.ChainConfig, chain *LightChain) *CtxPool {
	pool := &CtxPool{
		config:  config,
		signer:  types.NewEIP155Signer(config.ChainID),
		chain:   chain,
		pending: make(map[common.Hash]*cc.CrossTransactionWithSignatures),
	}
	return pool
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *CtxPool) addTx(cws *cc.CrossTransactionWithSignatures, local bool) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if _, ok := pool.pending[cws.Hash()]; !ok {
		pool.pending[cws.Hash()] = cws
	}
	return nil
}

func (pool *CtxPool) Stats() int {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	return len(pool.pending)
}

func (pool *CtxPool) Pending() (map[common.Hash]*cc.CrossTransactionWithSignatures, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	return pool.pending, nil
}
