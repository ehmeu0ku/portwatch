package sampler_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sampler"
	"github.com/user/portwatch/internal/scanner"
)

type fakeScanner struct {
	mu     sync.Mutex
	calls  int
	states []scanner.PortState
	err    error
}

func (f *fakeScanner) Scan() ([]scanner.PortState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	return f.states, f.err
}

func TestSamplerCallsConsumerImmediately(t *testing.T) {
	fs := &fakeScanner{states: []scanner.PortState{{Port: 80, Proto: "tcp"}}}
	var received []scanner.PortState
	ctx, cancel := context.WithCancel(context.Background())

	s := sampler.New(fs, 10*time.Second, func(states []scanner.PortState) {
		received = states
		cancel()
	})

	_ = s.Run(ctx)

	if len(received) != 1 || received[0].Port != 80 {
		t.Fatalf("expected port 80, got %v", received)
	}
}

func TestSamplerTicksMultipleTimes(t *testing.T) {
	fs := &fakeScanner{states: []scanner.PortState{}}
	var mu sync.Mutex
	count := 0
	ctx, cancel := context.WithCancel(context.Background())

	s := sampler.New(fs, 20*time.Millisecond, func(states []scanner.PortState) {
		mu.Lock()
		count++
		if count >= 3 {
			cancel()
		}
		mu.Unlock()
	})

	_ = s.Run(ctx)

	mu.Lock()
	defer mu.Unlock()
	if count < 3 {
		t.Fatalf("expected at least 3 samples, got %d", count)
	}
}

func TestSamplerReturnsErrorFromScanner(t *testing.T) {
	expected := errors.New("scan failed")
	fs := &fakeScanner{err: expected}
	s := sampler.New(fs, time.Second, func(_ []scanner.PortState) {})

	err := s.Run(context.Background())
	if !errors.Is(err, expected) {
		t.Fatalf("expected scan error, got %v", err)
	}
}
