package testing

import (
	"fmt"
	"sync"

	"gbchain-org/go-gbchain/log"
	"gbchain-org/go-gbchain/p2p/enode"
)

type TestPeer interface {
	ID() enode.ID
	Drop()
}

// TestPeerPool is an example peerPool to demonstrate registration of peer connections
type TestPeerPool struct {
	lock  sync.Mutex
	peers map[enode.ID]TestPeer
}

func NewTestPeerPool() *TestPeerPool {
	return &TestPeerPool{peers: make(map[enode.ID]TestPeer)}
}

func (p *TestPeerPool) Add(peer TestPeer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	log.Trace(fmt.Sprintf("pp add peer  %v", peer.ID()))
	p.peers[peer.ID()] = peer

}

func (p *TestPeerPool) Remove(peer TestPeer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.peers, peer.ID())
}

func (p *TestPeerPool) Has(id enode.ID) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, ok := p.peers[id]
	return ok
}

func (p *TestPeerPool) Get(id enode.ID) TestPeer {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.peers[id]
}
