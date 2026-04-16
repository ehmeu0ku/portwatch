package correlator_test

import (
	"testing"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/process"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/severity"
	"github.com/user/portwatch/internal/tagger"
)

func makeState(port uint16, proto string) scanner.PortState {
	return scanner.PortState{Port: port, Proto: proto, Inode: 0}
}

func newCorrelator() *correlator.Correlator {
	r := process.NewResolver("/proc")
	e := enricher.New(r)
	t := tagger.New(nil)
	s := severity.New(t)
	return correlator.New(e, t, s)
}

func TestCorrelateNewEvent(t *testing.T) {
	c := newCorrelator()
	ps := makeState(80, "tcp")
	ev := c.Correlate(correlator.KindNew, ps)

	if ev.Kind != correlator.KindNew {
		t.Fatalf("expected KindNew, got %s", ev.Kind)
	}
	if ev.State.Port != ps.Port {
		t.Fatalf("expected port %d, got %d", ps.Port, ev.State.Port)
	}
	if ev.Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestCorrelateGoneEvent(t *testing.T) {
	c := newCorrelator()
	ps := makeState(443, "tcp")
	ev := c.Correlate(correlator.KindGone, ps)

	if ev.Kind != correlator.KindGone {
		t.Fatalf("expected KindGone, got %s", ev.Kind)
	}
}

func TestCorrelatePrivilegedPortIsCritical(t *testing.T) {
	c := newCorrelator()
	ps := makeState(22, "tcp")
	ev := c.Correlate(correlator.KindNew, ps)

	if ev.Severity != severity.Critical {
		t.Fatalf("expected Critical, got %s", ev.Severity)
	}
}

func TestCorrelateTaggedPort(t *testing.T) {
	c := newCorrelator()
	ps := makeState(80, "tcp")
	ev := c.Correlate(correlator.KindNew, ps)

	if ev.Tag == "" {
		t.Fatal("expected non-empty tag for port 80")
	}
}

func TestCorrelateHighPortIsLowerSeverity(t *testing.T) {
	c := newCorrelator()
	ps := makeState(51234, "tcp")
	ev := c.Correlate(correlator.KindNew, ps)

	if ev.Severity == severity.Critical {
		t.Fatalf("high unprivileged port should not be Critical")
	}
}
