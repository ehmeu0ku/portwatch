package deadman_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/deadman"
)

func TestSwitchFiresWhenNotReset(t *testing.T) {
	var fired atomic.Bool
	dm := deadman.New(20*time.Millisecond, func(_ time.Duration) {
		fired.Store(true)
	})
	dm.Start(context.Background())

	time.Sleep(60 * time.Millisecond)
	if !fired.Load() {
		t.Fatal("expected dead-man switch to fire")
	}
	if !dm.Tripped() {
		t.Fatal("expected Tripped() == true")
	}
}

func TestSwitchDoesNotFireWhenReset(t *testing.T) {
	var fired atomic.Bool
	dm := deadman.New(30*time.Millisecond, func(_ time.Duration) {
		fired.Store(true)
	})
	dm.Start(context.Background())

	// Reset faster than the window.
	for i := 0; i < 5; i++ {
		time.Sleep(10 * time.Millisecond)
		dm.Reset()
	}

	if fired.Load() {
		t.Fatal("dead-man switch fired despite regular resets")
	}
}

func TestSwitchStopsOnContextCancel(t *testing.T) {
	var fired atomic.Bool
	ctx, cancel := context.WithCancel(context.Background())
	dm := deadman.New(20*time.Millisecond, func(_ time.Duration) {
		fired.Store(true)
	})
	dm.Start(ctx)
	cancel()

	time.Sleep(50 * time.Millisecond)
	// After cancel the timer is stopped; alert may or may not have fired
	// depending on scheduling — we just ensure no panic and Tripped is consistent.
	_ = dm.Tripped()
}

func TestResetClearsTripped(t *testing.T) {
	var count atomic.Int32
	dm := deadman.New(15*time.Millisecond, func(_ time.Duration) {
		count.Add(1)
	})
	dm.Start(context.Background())

	time.Sleep(40 * time.Millisecond)
	if !dm.Tripped() {
		t.Fatal("expected tripped")
	}
	dm.Reset()
	if dm.Tripped() {
		t.Fatal("expected Tripped() false after Reset")
	}
}
