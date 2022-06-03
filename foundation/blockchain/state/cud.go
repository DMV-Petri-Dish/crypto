package state

import "github.com/DMV-Petri-Dish/crypto/foundation/blockchain/peer"

// AddKnownPeer provides the ability to add a new peer to the known peer list
func (s *State) AddKnownPeer(peer peer.Peer) bool {
	return s.knownPeers.Add(peer)
}

// RemoveKnownPeer provides the ability to remove a peer from
// the known peer list.
func (s *State) RemoveKnownPeer(peer peer.Peer) {
	s.knownPeers.Remove(peer)
}
