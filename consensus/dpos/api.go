// Package dpos implements the delegated-proof-of-stake consensus engine.

package dpos

import (
	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/consensus"
	"gbchain-org/go-gbchain/core/types"
	"gbchain-org/go-gbchain/rpc"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the delegated-proof-of-stake scheme.
type API struct {
	chain consensus.ChainReader
	dpos  *DPoS
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dpos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)

}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dpos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
}

// GetSnapshotAtNumber retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtNumber(number uint64) (*Snapshot, error) {
	header := api.chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dpos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil, nil, defaultLoopCntRecalculateSigners)
}
