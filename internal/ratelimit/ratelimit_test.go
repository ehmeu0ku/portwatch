package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllowFirstCallPasses(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected first call to Allow to return true")
	}
}

func TestAllowSecondCallWithinCooldownBlocked(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("tcp:8080")
	if l.Allow("tcp:8080") {
		t.Fatal("expected second call within cooldown to return false")
	}
}

func TestAllowAfterCooldownPasses(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("tcp:9090")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("tcp:9090") {
		t.Fatal("expected call after cooldown expiry to return true")
	}
}

func TestAllowDifferentKeysPasses(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("tcp:8080")
	if !l.Allow("tcp:9090") {
		t.Fatal("expected different key to pass independently")
	}
}

func TestResetAllowsImmediateReuse(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("tcp:8080")
	l.Reset("tcp:8080")
	if !l.Allow("tcp:8080") {
		t.Fatal("expected Allow to pass after Reset")
	}
}

func TestResetUnknownKeyIsNoop(t *testing.T) {
	// Resetting a key that was never seen should not panic or cause errors.
	l := ratelimit.New(1 * time.Second)
	l.Reset("tcp:1234")
	if !l.Allow("tcp:1234") {
		t.Fatal("expected Allow to pass for key that was only Reset, never seen")
	}
}

func TestPurgeRemovesExpiredKeys(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("tcp:7070")
	time.Sleep(20 * time.Millisecond)
	l.Purge()
	// After purge the key should be gone; Allow should pass again
	if !l.Allow("tcp:7070") {
		t.Fatal("expected Allow to pass after Purge removed expired key")
	}
}

func TestPurgeKeepsActiveKeys(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("tcp:6060")
	l.Purge()
	// Key is still within cooldown, so second Allow should be blocked
	if l.Allow("tcp:6060") {
		t.Fatal("expected active key to remain after Purge")
	}
}
