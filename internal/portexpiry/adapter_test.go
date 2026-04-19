package portexpiry

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeState(proto string, port uint16) scanner.PortState {
	return scanner.PortState{Proto: proto, Port: port}
}

func TestKeyFromState(t *testing.T) {
	s := makeState("tcp", 8080)
	if got := KeyFromState(s); got != "tcp:8080" {
		t.Fatalf("unexpected key %q", got)
	}
}

func TestObserveGoneReturnsFalseFirst(t *testing.T) {
	obs := NewObserver(New(time.Minute))
	if obs.ObserveGone(makeState("tcp", 80)) {
		t.Fatal("expected false on first call")
	}
}

func TestObservePresentResetsAbsence(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.now = func() time.Time { return now }
	obs := NewObserver(tr)
	s := makeState("tcp", 443)
	obs.ObserveGone(s)
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	obs.ObservePresent(s)
	// after reset, first absence observation should return false
	if obs.ObserveGone(s) {
		t.Fatal("expected false after presence reset")
	}
}

func TestObserveGoneExpires(t *testing.T) {
	now := time.Now()
	tr := New(30 * time.Second)
	tr.now = func() time.Time { return now }
	obs := NewObserver(tr)
	s := makeState("udp", 53)
	obs.ObserveGone(s)
	tr.now = func() time.Time { return now.Add(31 * time.Second) }
	if !obs.ObserveGone(s) {
		t.Fatal("expected true after TTL elapsed")
	}
}
