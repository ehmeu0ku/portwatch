package dispatch_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/circuitbreaker"
	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/dispatch"
	"github.com/user/portwatch/internal/formatter"
	"github.com/user/portwatch/internal/pagerule"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/severity"
)

func makeEvent(port uint16, sev severity.Level) correlator.Event {
	return correlator.Event{
		Kind:     correlator.KindNew,
		Severity: sev,
		State:    scanner.PortState{Port: port, Proto: "tcp"},
	}
}

func defaultDispatcher(n *fakeNotifier) *dispatch.Dispatcher {
	rules := pagerule.New([]pagerule.Rule{
		{Action: pagerule.ActionPage, MinSeverity: severity.Info},
	})
	return dispatch.New(dispatch.Config{
		Rules:   rules,
		Format:  formatter.New(formatter.Text),
		Notify:  n,
		Breaker: circuitbreaker.New(circuitbreaker.DefaultConfig()),
	})
}

type fakeNotifier struct {
	calls []string
	err   error
}

func (f *fakeNotifier) Notify(_ context.Context, msg string) error {
	f.calls = append(f.calls, msg)
	return f.err
}
func (f *fakeNotifier) Name() string { return "fake" }

func TestDispatchSendsMatchingEvent(t *testing.T) {
	n := &fakeNotifier{}
	d := defaultDispatcher(n)
	ev := makeEvent(8080, severity.Warning)
	if err := d.Dispatch(context.Background(), ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(n.calls) != 1 {
		t.Fatalf("expected 1 notify call, got %d", len(n.calls))
	}
}

func TestDispatchSuppressesWhenRuleBlocks(t *testing.T) {
	n := &fakeNotifier{}
	rules := pagerule.New([]pagerule.Rule{
		{Action: pagerule.ActionPage, MinSeverity: severity.Critical},
	})
	d := dispatch.New(dispatch.Config{
		Rules:   rules,
		Format:  formatter.New(formatter.Text),
		Notify:  n,
		Breaker: circuitbreaker.New(circuitbreaker.DefaultConfig()),
	})
	if err := d.Dispatch(context.Background(), makeEvent(9090, severity.Info)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(n.calls) != 0 {
		t.Fatalf("expected 0 notify calls, got %d", len(n.calls))
	}
}

func TestDispatchRecordsFailureOnNotifyError(t *testing.T) {
	n := &fakeNotifier{err: errors.New("backend down")}
	d := defaultDispatcher(n)
	err := d.Dispatch(context.Background(), makeEvent(443, severity.Critical))
	if err == nil {
		t.Fatal("expected error from failed notifier")
	}
}
