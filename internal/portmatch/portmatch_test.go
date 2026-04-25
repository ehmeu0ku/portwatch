package portmatch_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmatch"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(port uint16, proto string, tags ...string) scanner.PortState {
	return scanner.PortState{Port: port, Proto: proto, Tags: tags}
}

func TestExactPortMatches(t *testing.T) {
	e := portmatch.Expr{Port: 80}
	if !e.Match(makeState(80, "tcp")) {
		t.Fatal("expected match on port 80")
	}
	if e.Match(makeState(81, "tcp")) {
		t.Fatal("expected no match on port 81")
	}
}

func TestRangeMatches(t *testing.T) {
	e := portmatch.Expr{PortMin: 1000, PortMax: 2000}
	if !e.Match(makeState(1500, "tcp")) {
		t.Fatal("expected match inside range")
	}
	if !e.Match(makeState(1000, "tcp")) {
		t.Fatal("expected match at lower bound")
	}
	if !e.Match(makeState(2000, "tcp")) {
		t.Fatal("expected match at upper bound")
	}
	if e.Match(makeState(999, "tcp")) {
		t.Fatal("expected no match below range")
	}
	if e.Match(makeState(2001, "tcp")) {
		t.Fatal("expected no match above range")
	}
}

func TestProtoFilter(t *testing.T) {
	e := portmatch.Expr{Port: 443, Proto: "tcp"}
	if !e.Match(makeState(443, "tcp")) {
		t.Fatal("expected tcp match")
	}
	if e.Match(makeState(443, "udp")) {
		t.Fatal("expected no match on udp")
	}
}

func TestTagFilter(t *testing.T) {
	e := portmatch.Expr{Tag: "web"}
	if !e.Match(makeState(8080, "tcp", "web", "http")) {
		t.Fatal("expected match with web tag")
	}
	if e.Match(makeState(8080, "tcp", "db")) {
		t.Fatal("expected no match without web tag")
	}
}

func TestWildcardMatchesAll(t *testing.T) {
	e := portmatch.Expr{}
	if !e.Match(makeState(12345, "udp", "custom")) {
		t.Fatal("empty expr should match anything")
	}
}

func TestAnyMatch(t *testing.T) {
	m := portmatch.New([]portmatch.Expr{
		{Port: 22},
		{Port: 80},
	})
	if !m.AnyMatch(makeState(22, "tcp")) {
		t.Fatal("expected any-match on port 22")
	}
	if m.AnyMatch(makeState(9999, "tcp")) {
		t.Fatal("expected no any-match on port 9999")
	}
}

func TestAllMatch(t *testing.T) {
	m := portmatch.New([]portmatch.Expr{
		{Proto: "tcp"},
		{Tag: "web"},
	})
	if !m.AllMatch(makeState(80, "tcp", "web")) {
		t.Fatal("expected all-match")
	}
	if m.AllMatch(makeState(80, "udp", "web")) {
		t.Fatal("expected all-match failure on proto")
	}
}

func TestMatcherString(t *testing.T) {
	m := portmatch.New([]portmatch.Expr{{Port: 443, Proto: "tcp"}})
	s := m.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
