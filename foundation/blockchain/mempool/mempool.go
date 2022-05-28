// Package mempool maintains the mempool for the blockchain
package mempool

import (
	"sync"

	"github.com/DMV-Petri-Dish/crypto/foundation/blockchain/database"
)

// Mempool represents a cache of transactions organized by account:nonce
type Mempool struct {
	mu       sync.RWMutex
	pool     map[string]database.BlockTx
	selectFn selector.Func
}

// New constructs a new mempool using the default sort strategy
func New() (*Mempool, error) {
	return NewWithStrategy(selector.StrategyTip)
}

// NewWithStrategy constructs a new mempool with specified sort strategy
func NewWithStrategy(strategy string) (*Mempool, error) {
	selectFn, err := selector.Retrieve(strategy)
	if err != nil {
		return nil, err
	}

	mp := Mempool{
		pool:     make(map[string]database.BlockTx),
		selectFn: selectFn,
	}

	return &mp, nil
}
