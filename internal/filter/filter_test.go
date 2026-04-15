package filter_test

import (
	"net"
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(ip string, port uint16) scanner.PortState {
	return scanner.PortState{IP: net.ParseIP(ip), Port: port, Proto: "tcp"}
}

func defaultCfg() *config.Config {
	c := config.DefaultConfig()
	return c
}

func TestApplyPassesAllByDefault(t *testing.T) {
	cfg := defaultCfg()
	f := filter.New(cfg)

	states := []scanner.PortState{
		makeState("192.168.1.1", 8080),
		makeState("10.0.0.1", 443),
	}

	got := f.Apply(states)
	if len(got) != 2 {
		t.Fatalf("expected 2 states, got %d", len(got))
	}
}

func TestApplyFiltersIgnoredPorts(t *testing.T) {
	cfg := defaultCfg()
	cfg.IgnoredPorts = []uint16{22, 80}
	f := filter.New(cfg)

	states := []scanner.PortState{
		makeState("192.168.1.1", 22),
		makeState("192.168.1.1", 80),
		makeState("192.168.1.1", 8080),
	}

	got := f.Apply(states)
	if len(got) != 1 {
		t.Fatalf("expected 1 state, got %d", len(got))
	}
	if got[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", got[0].Port)
	}
}

func TestApplyFiltersLoopbackWhenEnabled(t *testing.T) {
	cfg := defaultCfg()
	cfg.IgnoreLoopback = true
	f := filter.New(cfg)

	states := []scanner.PortState{
		makeState("127.0.0.1", 9000),
		makeState("0.0.0.0", 9001),
	}

	got := f.Apply(states)
	if len(got) != 1 {
		t.Fatalf("expected 1 state, got %d", len(got))
	}
	if got[0].Port != 9001 {
		t.Errorf("expected port 9001, got %d", got[0].Port)
	}
}

func TestApplyEmptyInput(t *testing.T) {
	cfg := defaultCfg()
	f := filter.New(cfg)
	got := f.Apply(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}
