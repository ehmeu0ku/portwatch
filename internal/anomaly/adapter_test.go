package anomaly_test

import (
	"testing"

	"github.com/user/portwatch/internal/anomaly"
	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/scorecard"
)

func makeEvent(kind string, port uint16) correlator.Event {
	return correlator.Event{
		Kind:  kind,
		State: scanner.PortState{Proto: "tcp", Port: port},
	}
}

func TestKeyFromEvent(t *testing.T) {
	e := makeEvent("new", 8080)
	got := anomaly.KeyFromEvent(e)
	if got != "tcp:8080" {
		t.Fatalf("unexpected key %q", got)
	}
}

func TestRecorderIncrementsOnNewEvent(t *testing.T) {
	sc := scorecard.New()
	r := anomaly.NewRecorder(sc)
	r.Record(makeEvent("new", 9000))
	snap, ok := sc.Get("tcp:9000")
	if !ok {
		t.Fatal("expected scorecard entry")
	}
	if snap.Alerts != 1 {
		t.Fatalf("expected 1 alert, got %d", snap.Alerts)
	}
}

func TestRecorderDoesNotAlertOnGoneEvent(t *testing.T) {
	sc := scorecard.New()
	r := anomaly.NewRecorder(sc)
	r.Record(makeEvent("gone", 9001))
	snap, ok := sc.Get("tcp:9001")
	if !ok {
		t.Fatal("expected scorecard entry")
	}
	if snap.Alerts != 0 {
		t.Fatalf("expected 0 alerts, got %d", snap.Alerts)
	}
	if snap.Seen != 1 {
		t.Fatalf("expected 1 seen, got %d", snap.Seen)
	}
}
