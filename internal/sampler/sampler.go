// Package sampler provides periodic sampling of port states,
// collecting snapshots at a configurable interval and forwarding
// them to a consumer function.
package sampler

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Scanner is the interface used to collect port states.
type Scanner interface {
	Scan() ([]scanner.PortState, error)
}

// ConsumerFunc is called with each collected sample.
type ConsumerFunc func([]scanner.PortState)

// Sampler periodically scans and delivers results to a consumer.
type Sampler struct {
	scanner  Scanner
	interval time.Duration
	consumer ConsumerFunc
}

// New creates a Sampler that calls consumer every interval.
func New(s Scanner, interval time.Duration, consumer ConsumerFunc) *Sampler {
	return &Sampler{
		scanner:  s,
		interval: interval,
		consumer: consumer,
	}
}

// Run starts the sampling loop and blocks until ctx is cancelled.
func (s *Sampler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Perform an immediate scan before waiting for the first tick.
	if err := s.sample(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.sample(); err != nil {
				return err
			}
		}
	}
}

func (s *Sampler) sample() error {
	states, err := s.scanner.Scan()
	if err != nil {
		return err
	}
	s.consumer(states)
	return nil
}
