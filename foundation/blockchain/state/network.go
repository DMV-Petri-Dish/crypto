package state

import (
	"fmt"
	"net/http"
)

const baseURL = "http://%s/v1/node"

// NetSendBlockToPeers takes the new mined block and sends it to all know peers.
func (s *State) NetSendBlockToPeers(block database.Block) error {
	s.evHandler("state: NetSendBlockToPeers: started")
	defer s.evHandler("state: NetSendBlockToPeers: completed")

	for _, peer := range s.RetrieveKnownPeers() {
		s.evHandler("state: NetSendBlockToPeers: send: block[%s] to peer[%s]", block.Hash(), peer)

		url := fmt.Sprintf("%s/block/propose", fmt.Sprintf(baseURL, peer.Host))

		var status struct {
			Status string `json:"status"`
		}
		if err := send(http.MethodPost, url, database.NewBlockData(block), &status); err != nil {
			return fmt.Errorf("%s: %s", peer.Host, err)
		}
	}

	return nil
}

// NetSendTxToPeers shares a new block transaction with the known peers.
func (s *State) NetSendTxToPeers(tx database.BlockTx) {
	s.evHandler("state: NetSendTxToPeers: started")
	defer s.evHandler("state: NetSendTxToPeers: completed")

	// CORE NOTE: Bitcoin does not send the full transaction immediately to save
	// on bandwidth. A node will send the transaction's mempool key first so the
	// receiving node can check if they already have the transaction or not. If
	// the receiving node doesn't have it, then it will request the transaction
	// based on the mempool key it received.

	// For now, the Ardan blockchain just sends the full transaction.
	for _, peer := range s.RetrieveKnownPeers() {
		s.evHandler("state: NetSendTxToPeers: send: tx[%s] to peer[%s]", tx, peer)

		url := fmt.Sprintf("%s/tx/submit", fmt.Sprintf(baseURL, peer.Host))

		if err := send(http.MethodPost, url, tx, nil); err != nil {
			s.evHandler("state: NetSendTxToPeers: WARNING: %s", err)
		}
	}
}

// NetSendNodeAvailableToPeers shares this node is available to
// participate in the network with the known peers.
func (s *State) NetSendNodeAvailableToPeers() {
	s.evHandler("state: NetSendNodeAvailableToPeers: started")
	defer s.evHandler("state: NetSendNodeAvailableToPeers: completed")

	host := peer.Peer{Host: s.RetrieveHost()}

	for _, peer := range s.RetrieveKnownPeers() {
		s.evHandler("state: NetSendNodeAvailableToPeers: send: host[%s] to peer[%s]", host, peer)

		url := fmt.Sprintf("%s/peers", fmt.Sprintf(baseURL, peer.Host))

		if err := send(http.MethodPost, url, host, nil); err != nil {
			s.evHandler("state: NetSendNodeAvailableToPeers: WARNING: %s", err)
		}
	}
}

// NetRequestPeerStatus looks for new nodes on the blockchain by asking
// known nodes for their peer list. New nodes are added to the list.
func (s *State) NetRequestPeerStatus(pr peer.Peer) (peer.PeerStatus, error) {
	s.evHandler("state: NetRequestPeerStatus: started: %s", pr)
	defer s.evHandler("state: NetRequestPeerStatus: completed: %s", pr)

	url := fmt.Sprintf("%s/status", fmt.Sprintf(baseURL, pr.Host))

	var ps peer.PeerStatus
	if err := send(http.MethodGet, url, nil, &ps); err != nil {
		return peer.PeerStatus{}, err
	}

	s.evHandler("state: NetRequestPeerStatus: peer-node[%s]: latest-blknum[%d]: peer-list[%s]", pr, ps.LatestBlockNumber, ps.KnownPeers)

	return ps, nil
}
