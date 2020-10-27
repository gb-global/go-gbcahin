// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/core/types"
)

// Constants to match up protocol versions and messages
const (
	Eth63 = 63
	Eth64 = 64
)

var (
	EthProtocol = Protocol{
		Name:     "eth",
		Versions: []uint{Eth64, Eth63},
		Lengths:  map[uint]uint64{Eth64: 17, Eth63: 17},
	}
)

// Protocol defines the protocol of the consensus
type Protocol struct {
	// Official short name of the protocol used during capability negotiation.
	Name string
	// Supported versions of the eth protocol (first is primary).
	Versions []uint
	// Number of implemented message corresponding to different protocol versions.
	Lengths map[uint]uint64
}

// Broadcaster defines the interface to enqueue blocks to fetcher and find peer
type Broadcaster interface {
	// Enqueue add a block into fetcher queue
	Enqueue(id string, block *types.Block)
	// FindPeers retrives peers by addresses
	FindPeers(map[common.Address]bool) map[common.Address]Peer
}

// Peer defines the interface to communicate with peer
type Peer interface {
	// Send sends the message to this peer
	Send(msgcode uint64, data interface{}) error
}
