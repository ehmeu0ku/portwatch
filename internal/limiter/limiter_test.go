package limiter_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/limiter"
)

func TestAllowFirstCallPasses(t *testing.T) {
	l := limiter.New(time.Second, 3)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected first call to pass")
	}
}

func TestAllowUpToBurst(t *testing.T) {
	l := limiter.New(time.Second, 3)
	for i := 0; i < 3; i++ {
		if !l.Allow("tcp:9000") {
			t.Fatalf("expected call %d to pass", i+1)
		}
	}
	if l.Allow("tcp:9000") {
		t.Fatal("expected call beyond burst to be blocked")
	}
}

func TestAllowDifferentKeysAreIndependent(t *testing.T) {
	l := limiter.New(time.Second, 1)
	if !l.Allow("tcp:80") {
		t.Fatal("expected tcp:80 to pass")
	}
	if !l.Allow("tcp:443") {
		t.Fatal("expected tcp:443 to pass independently")
	}
}

func TestAllowAfterWindowPasses(t *testing.T) {
	l := limiter.New(20*time.Millisecond, 1)
	l.Allow("tcp:8080")
	if l.Allow("tcp:8080") {
		t.Fatal("expected second call within window to be blocked")
	}
	time.Sleep(30 * time.Millisecond)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected call after window expiry to pass")
	}
}

func TestResetAllowsImmediateReuse(t *testing.T) {
	l := limiter.New(time.Second, 1)
	l.Allow("tcp:8080")
	if l.Allow("tcp:8080") {
		t.Fatal("expected blocked before reset")
	}
	l.Reset("tcp:8080")
	if !l.Allow("tcp:8080") {
		t.Fatal("expected pass after reset")
	}
}

func TestLenReflectsActiveBuckets(t *testing.T) {
	l := limiter.New(time.Second, 5)
	if l.Len() != 0 {
		t.Fatalf("expected 0 buckets, got %d", l.Len())
	}
	l.Allow("a")
	l.Allow("b")
	if l.Len() != 2 {
		t.Fatalf("expected 2 buckets, got %d", l.Len())
	}
	l.Reset("a")
	if l.Len() != 1 {
		t.Fatalf("expected 1 bucket after reset, got %d", l.Len())
	}
}
