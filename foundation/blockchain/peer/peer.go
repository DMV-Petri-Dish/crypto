// Package peer maintains the peer related information such as the set
// of known peers and their status.
package peer

import "sync"

// Peer represents information about a Node in the network.
type Peer struct {
	Host string
}

// New constructs a new info value.
func New(host string) Peer {
	return Peer{
		Host: host,
	}
}

// =======================================

// PeerSet represents the data representation to maintain a set of known peers.
type PeerSet struct {
	mu  sync.RWMutex
	set map[Peer]struct{}
}

// NewPeerSet constructs a new info set to manage node peer information.
func NewPeerSet() *PeerSet {
	return &PeerSet{
		set: make(map[Peer]struct{}),
	}
}

// Add adds a new node to the set.
func (ps *PeerSet) Add(peer Peer) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	_, exists := ps.set[peer]
	if !exists {
		ps.set[peer] = struct{}{}
		return true
	}

	return false
}

// Remove removes a node from the set.
func (ps *PeerSet) Remove(peer Peer) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.set, peer)
}
