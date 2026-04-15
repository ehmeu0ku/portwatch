package suppress_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

func key(kind string) suppress.Key {
	return suppress.Key{Proto: "tcp", Addr: "0.0.0.0", Port: 8080, Kind: kind}
}

func TestAllowFirstCallPasses(t *testing.T) {
	s := suppress.New(time.Minute)
	if !s.Allow(key("new")) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowSecondCallWithinWindowBlocked(t *testing.T) {
	s := suppress.New(time.Minute)
	s.Allow(key("new"))
	if s.Allow(key("new")) {
		t.Fatal("expected second call within window to be suppressed")
	}
}

func TestAllowAfterWindowPasses(t *testing.T) {
	s := suppress.New(10 * time.Millisecond)
	s.Allow(key("new"))
	time.Sleep(20 * time.Millisecond)
	if !s.Allow(key("new")) {
		t.Fatal("expected call after window expiry to be allowed")
	}
}

func TestAllowDifferentKindsAreIndependent(t *testing.T) {
	s := suppress.New(time.Minute)
	s.Allow(key("new"))
	if !s.Allow(key("gone")) {
		t.Fatal("expected different kind to be allowed independently")
	}
}

func TestResetAllowsImmediateReuse(t *testing.T) {
	s := suppress.New(time.Minute)
	s.Allow(key("new"))
	s.Reset(key("new"))
	if !s.Allow(key("new")) {
		t.Fatal("expected allow after reset")
	}
}

func TestPurgeClearsExpiredEntries(t *testing.T) {
	s := suppress.New(10 * time.Millisecond)
	s.Allow(key("new"))
	s.Allow(suppress.Key{Proto: "udp", Addr: "127.0.0.1", Port: 53, Kind: "new"})

	if s.Len() != 2 {
		t.Fatalf("expected 2 tracked keys, got %d", s.Len())
	}

	time.Sleep(20 * time.Millisecond)
	s.Purge()

	if s.Len() != 0 {
		t.Fatalf("expected 0 keys after purge, got %d", s.Len())
	}
}

func TestPurgeKeepsActiveEntries(t *testing.T) {
	s := suppress.New(time.Minute)
	s.Allow(key("new"))
	s.Purge()
	if s.Len() != 1 {
		t.Fatalf("expected active entry to survive purge, got %d", s.Len())
	}
}
