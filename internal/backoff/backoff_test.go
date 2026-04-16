package backoff_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/backoff"
)

func TestNextIncreasesWithAttempts(t *testing.T) {
	b := backoff.New(backoff.Config{
		Initial:    100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 2.0,
		Jitter:     0,
	})

	d1 := b.Next()
	d2 := b.Next()
	d3 := b.Next()

	if d2 <= d1 {
		t.Errorf("expected d2 > d1, got %v <= %v", d2, d1)
	}
	if d3 <= d2 {
		t.Errorf("expected d3 > d2, got %v <= %v", d3, d2)
	}
}

func TestNextRespectsMax(t *testing.T) {
	max := 500 * time.Millisecond
	b := backoff.New(backoff.Config{
		Initial:    100 * time.Millisecond,
		Max:        max,
		Multiplier: 10.0,
		Jitter:     0,
	})

	for i := 0; i < 10; i++ {
		d := b.Next()
		if d > max {
			t.Errorf("attempt %d: duration %v exceeds max %v", i, d, max)
		}
	}
}

func TestResetClearsAttempt(t *testing.T) {
	b := backoff.New(backoff.DefaultConfig())
	b.Next()
	b.Next()
	if b.Attempt() != 2 {
		t.Fatalf("expected attempt 2, got %d", b.Attempt())
	}
	b.Reset()
	if b.Attempt() != 0 {
		t.Errorf("expected attempt 0 after reset, got %d", b.Attempt())
	}
}

func TestDefaultConfigIsReasonable(t *testing.T) {
	cfg := backoff.DefaultConfig()
	if cfg.Initial <= 0 {
		t.Error("initial must be positive")
	}
	if cfg.Max < cfg.Initial {
		t.Error("max must be >= initial")
	}
	if cfg.Multiplier <= 1.0 {
		t.Error("multiplier must be > 1")
	}
}

func TestJitterProducesVariance(t *testing.T) {
	b := backoff.New(backoff.Config{
		Initial:    200 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 1.0,
		Jitter:     0.5,
	})

	seen := map[time.Duration]bool{}
	for i := 0; i < 20; i++ {
		b.Reset()
		seen[b.Next()] = true
	}
	if len(seen) < 2 {
		t.Error("expected jitter to produce varied durations")
	}
}
