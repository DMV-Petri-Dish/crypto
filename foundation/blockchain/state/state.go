// Package state is the core API for the blockchain and implements
// all the business rules and processing
package state

import (
	"sync"

	"github.com/DMV-Petri-Dish/crypto/foundation/blockchain/database"
	"github.com/DMV-Petri-Dish/crypto/foundation/blockchain/genesis"
	"github.com/DMV-Petri-Dish/crypto/foundation/blockchain/peer"
)

type EventHandler func(v string, args ...interface{})

type Worker interface {
	Shutdown()
	Sync()
	SignalStartMining()
	SignalCancelMining(done func())
	SignalShareTx(blockTx database.BlockTx)
}

// =========================================

// Config represents the configuration required to start
// the blockchain node.
type Config struct {
	BeneficiaryID  database.AccountID
	Host           string
	DBPath         string
	SelectStrategy string
	KnownPeers     *peer.PeerSet
	EvHandler      EventHandler
}

// State manages the blockchain database.
type State struct {
	mu sync.RWMutex

	beneficiaryID database.AccountID
	host          string
	dbPath        string
	evHandler     EventHandler

	allowMining bool
	resyncWG    sync.WaitGroup

	knownPeers *peer.PeerSet
	genesis    genesis.Genesis
	mempool    *mempool.Mempool
	db         *database.Database

	Worker Worker
}

// New constructs a new blockchain for data management.
func New(cfg Config) (*State, error) {

	// Build a safe event handler function for use.
	ev := func(v string, args ...interface{}) {
		if cfg.EvHandler != nil {
			cfg.EvHandler(v, args...)
		}
	}

	// Load the genesis file to get starting balances for
	// founders of the block chain.
	genesis, err := genesis.Load()
	if err != nil {
		return nil, err
	}

	storage, err := disk.New(cfg.DBPath)
	//storage, err := memory.New()
	if err != nil {
		return nil, err
	}

	// Access the storage for the blockchain.
	db, err := database.New(genesis, storage, ev)
	if err != nil {
		return nil, err
	}

	// Construct a mempool with the specified sort strategy.
	mempool, err := mempool.NewWithStrategy(cfg.SelectStrategy)
	if err != nil {
		return nil, err
	}

	// Create the State to provide support for managing the blockchain.
	state := State{
		beneficiaryID: cfg.BeneficiaryID,
		host:          cfg.Host,
		dbPath:        cfg.DBPath,
		evHandler:     ev,
		allowMining:   true,

		knownPeers: cfg.KnownPeers,
		genesis:    genesis,
		mempool:    mempool,
		db:         db,
	}

	// The Worker is not set here. The call to worker.Run will assign itself
	// and start everything up and running for the node.

	return &state, nil
}

// Shutdown cleanly brings the node down.
func (s *State) Shutdown() error {
	s.evHandler("state: shutdown: started")
	defer s.evHandler("state: shutdown: completed")

	// Make sure the database file is properly closed.
	defer func() {
		s.db.Close()
	}()

	// Stop all blockchain writing activity.
	s.Worker.Shutdown()

	// Wait for any resync to finish.
	s.resyncWG.Wait()

	return nil
}
