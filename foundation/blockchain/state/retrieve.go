package state

import "github.com/DMV-Petri-Dish/crypto/foundation/blockchain/peer"

// RetrieveHost returns a copy of host information.
func (s *State) RetrieveHost() string {
	return s.host
}

// RetrieveKnownPeers retrieves a copy of the known peer list.
func (s *State) RetrieveKnownPeers() []peer.Peer {
	return s.knownPeers.Copy(s.host)
}
