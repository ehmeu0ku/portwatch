package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/reddysteady/portwatch/internal/circuitbreaker"
)

func TestAllowWhenClosedPasses(t *testing.T) {
	br := circuitbreaker.New(3, time.Second)
	if err := br.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestOpensAfterThreshold(t *testing.T) {
	br := circuitbreaker.New(2, time.Second)
	br.RecordFailure()
	br.RecordFailure()
	if br.State() != circuitbreaker.StateOpen {
		t.Fatal("expected circuit to be open")
	}
	if err := br.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestSuccessResetsFailures(t *testing.T) {
	br := circuitbreaker.New(3, time.Second)
	br.RecordFailure()
	br.RecordFailure()
	br.RecordSuccess()
	br.RecordFailure() // only 1 failure after reset — should not open
	if br.State() != circuitbreaker.StateClosed {
		t.Fatal("expected circuit to remain closed")
	}
}

func TestResetsAfterCooldown(t *testing.T) {
	br := circuitbreaker.New(1, 10*time.Millisecond)
	br.RecordFailure()
	if err := br.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen immediately, got %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := br.Allow(); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
	if br.State() != circuitbreaker.StateClosed {
		t.Fatal("expected state closed after cooldown")
	}
}

func TestForceReset(t *testing.T) {
	br := circuitbreaker.New(1, time.Hour)
	br.RecordFailure()
	if br.State() != circuitbreaker.StateOpen {
		t.Fatal("expected open")
	}
	br.Reset()
	if br.State() != circuitbreaker.StateClosed {
		t.Fatal("expected closed after reset")
	}
	if err := br.Allow(); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}
