package statestore

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Syncer periodically snapshots scanner output into a Store.
type Syncer struct {
	store    *Store
	scanner  scanner.Scanner
	interval time.Duration
}

// NewSyncer creates a Syncer that writes scanner results to store every interval.
func NewSyncer(store *Store, sc scanner.Scanner, interval time.Duration) *Syncer {
	return &Syncer{store: store, scanner: sc, interval: interval}
}

// Run starts the sync loop. It blocks until ctx is cancelled.
func (sy *Syncer) Run(ctx context.Context) error {
	if err := sy.sync(); err != nil {
		return err
	}
	ticker := time.NewTicker(sy.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_ = sy.sync() // best-effort; errors are non-fatal
		}
	}
}

func (sy *Syncer) sync() error {
	states, err := sy.scanner.Scan()
	if err != nil {
		return err
	}
	return sy.store.Set(states)
}
