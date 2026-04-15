package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestWatchdogFiresOnTimeout(t *testing.T) {
	t.Parallel()

	var fired atomic.Bool
	wd := watchdog.New(50*time.Millisecond, func() {
		fired.Store(true)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go wd.Start(ctx)

	// Do NOT send a beat — expect timeout to fire.
	time.Sleep(150 * time.Millisecond)

	if !fired.Load() {
		t.Fatal("expected watchdog timeout callback to be called")
	}
}

func TestWatchdogDoesNotFireWhenBeating(t *testing.T) {
	t.Parallel()

	var fired atomic.Bool
	wd := watchdog.New(80*time.Millisecond, func() {
		fired.Store(true)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go wd.Start(ctx)

	// Beat faster than the timeout to keep the watchdog satisfied.
	for i := 0; i < 5; i++ {
		time.Sleep(20 * time.Millisecond)
		wd.Beat()
	}

	if fired.Load() {
		t.Fatal("watchdog should not have fired while beats were sent")
	}
}

func TestWatchdogStopsOnContextCancel(t *testing.T) {
	t.Parallel()

	wd := watchdog.New(200*time.Millisecond, func() {})

	ctx, cancel := context.WithCancel(context.Background())
	go wd.Start(ctx)

	time.Sleep(20 * time.Millisecond)
	if !wd.Running() {
		t.Fatal("expected watchdog to be running")
	}

	cancel()
	time.Sleep(30 * time.Millisecond)

	if wd.Running() {
		t.Fatal("expected watchdog to have stopped after context cancel")
	}
}

func TestWatchdogBeatResetsTimer(t *testing.T) {
	t.Parallel()

	var count atomic.Int32
	wd := watchdog.New(60*time.Millisecond, func() {
		count.Add(1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go wd.Start(ctx)

	// Beat once just before the first timeout would fire.
	time.Sleep(40 * time.Millisecond)
	wd.Beat()

	// Wait long enough that without the reset a second fire would occur.
	time.Sleep(40 * time.Millisecond)

	if count.Load() != 0 {
		t.Fatalf("expected 0 timeouts after beat reset, got %d", count.Load())
	}
}
