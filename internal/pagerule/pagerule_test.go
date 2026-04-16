package pagerule_test

import (
	"testing"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/pagerule"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/severity"
)

func makeEvent(sev severity.Level, port uint16, tags []string) correlator.Event {
	return correlator.Event{
		State: scanner.PortState{Port: port},
		Severity: sev,
		Tags:     tags,
	}
}

func TestNoRulesAlwaysSuppresses(t *testing.T) {
	ev := New(nil)
	if got := ev.Evaluate(makeEvent(severity.Critical, 22, nil)); got != pagerule.ActionSuppress {
		t.Fatalf("expected suppress, got %s", got)
	}
}

func TestMinSeverityPages(t *testing.T) {
	rules := []pagerule.Rule{
		{MinSeverity: severity.Critical, Action: pagerule.ActionPage},
	}
	ev := pagerule.New(rules)
	if got := ev.Evaluate(makeEvent(severity.Critical, 80, nil)); got != pagerule.ActionPage {
		t.Fatalf("expected page, got %s", got)
	}
}

func TestBelowMinSeveritySuppresses(t *testing.T) {
	rules := []pagerule.Rule{
		{MinSeverity: severity.Critical, Action: pagerule.ActionPage},
	}
	ev := pagerule.New(rules)
	if got := ev.Evaluate(makeEvent(severity.Info, 80, nil)); got != pagerule.ActionSuppress {
		t.Fatalf("expected suppress, got %s", got)
	}
}

func TestPortFilterMatches(t *testing.T) {
	rules := []pagerule.Rule{
		{MinSeverity: severity.Info, Ports: []uint16{443}, Action: pagerule.ActionPage},
	}
	ev := pagerule.New(rules)
	if got := ev.Evaluate(makeEvent(severity.Info, 443, nil)); got != pagerule.ActionPage {
		t.Fatalf("expected page, got %s", got)
	}
}

func TestPortFilterMiss(t *testing.T) {
	rules := []pagerule.Rule{
		{MinSeverity: severity.Info, Ports: []uint16{443}, Action: pagerule.ActionPage},
	}
	ev := pagerule.New(rules)
	if got := ev.Evaluate(makeEvent(severity.Info, 80, nil)); got != pagerule.ActionSuppress {
		t.Fatalf("expected suppress, got %s", got)
	}
}

func TestTagFilterMatches(t *testing.T) {
	rules := []pagerule.Rule{
		{MinSeverity: severity.Info, Tags: []string{"ssh"}, Action: pagerule.ActionPage},
	}
	ev := pagerule.New(rules)
	if got := ev.Evaluate(makeEvent(severity.Info, 22, []string{"ssh"})); got != pagerule.ActionPage {
		t.Fatalf("expected page, got %s", got)
	}
}

func TestFirstRuleWins(t *testing.T) {
	rules := []pagerule.Rule{
		{MinSeverity: severity.Info, Action: pagerule.ActionSuppress},
		{MinSeverity: severity.Info, Action: pagerule.ActionPage},
	}
	ev := pagerule.New(rules)
	if got := ev.Evaluate(makeEvent(severity.Warning, 9000, nil)); got != pagerule.ActionSuppress {
		t.Fatalf("expected suppress (first rule), got %s", got)
	}
}

func New(rules []pagerule.Rule) *pagerule.Evaluator { return pagerule.New(rules) }
