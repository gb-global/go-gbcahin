package scrypt

import (
	"errors"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/common/hexutil"
	"gbchain-org/go-gbchain/core/types"
)

var errScryptStopped = errors.New("scrypt stopped")

// API exposes scrypt related methods for the RPC interface.
type API struct {
	powScrypt *PowScrypt // Make sure the mode of scrypt is normal.
}

// TODO: Scrypt work define
// GetWork returns a work package for external miner.
//
// The work package consists of 3 strings:
//   result[0] - 32 bytes hex encoded current block header pow-hash
//   result[1] - 32 bytes hex encoded boundary condition ("target"), 2^256/difficulty
//   result[2] - hex encoded block number
func (api *API) GetWork() ([2]string, error) {
	if api.powScrypt.config.PowMode != ModeNormal && api.powScrypt.config.PowMode != ModeTest {
		return [2]string{}, errors.New("not supported")
	}

	var (
		workCh = make(chan [2]string, 1)
		errc   = make(chan error, 1)
	)

	select {
	case api.powScrypt.fetchWorkCh <- &sealWork{errc: errc, res: workCh}:
	case <-api.powScrypt.exitCh:
		return [2]string{}, errScryptStopped
	}

	select {
	case work := <-workCh:
		return work, nil
	case err := <-errc:
		return [2]string{}, err
	}
}

// SubmitWork can be used by external miner to submit their POW solution.
// It returns an indication if the work was accepted.
// Note either an invalid solution, a stale work a non-existent work will return false.
func (api *API) SubmitWork(nonce types.BlockNonce, hash common.Hash) bool {
	if api.powScrypt.config.PowMode != ModeNormal && api.powScrypt.config.PowMode != ModeTest {
		return false
	}

	var errc = make(chan error, 1)

	digest, _ := ScryptHash(hash[:], nonce.Uint64())

	select {
	case api.powScrypt.submitWorkCh <- &mineResult{
		nonce:     nonce,
		mixDigest: common.BytesToHash(digest),
		hash:      hash,
		errc:      errc,
	}:
	case <-api.powScrypt.exitCh:
		return false
	}

	err := <-errc
	return err == nil
}

// SubmitHashrate can be used for remote miners to submit their hash rate.
// This enables the node to report the combined hash rate of all miners
// which submit work through this node.
//
// It accepts the miner hash rate and an identifier which must be unique
// between nodes.
func (api *API) SubmitHashRate(rate hexutil.Uint64, id common.Hash) bool {
	if api.powScrypt.config.PowMode != ModeNormal && api.powScrypt.config.PowMode != ModeTest {
		return false
	}

	var done = make(chan struct{}, 1)

	select {
	case api.powScrypt.submitRateCh <- &hashrate{done: done, rate: uint64(rate), id: id}:
	case <-api.powScrypt.exitCh:
		return false
	}

	// Block until hash rate submitted successfully.
	<-done

	return true
}

// GetHashrate returns the current hashrate for local CPU miner and remote miner.
func (api *API) GetHashrate() uint64 {
	return uint64(api.powScrypt.Hashrate())
}
