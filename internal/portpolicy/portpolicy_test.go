package portpolicy_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/portwatch/internal/portpolicy"
)

func TestDefaultDenyBlocksUnknownPort(t *testing.T) {
	pol := portpolicy.New(portpolicy.Deny)
	if got := pol.Evaluate(9999, "tcp"); got != portpolicy.Deny {
		t.Fatalf("expected Deny, got %s", got)
	}
}

func TestDefaultAllowPassesUnknownPort(t *testing.T) {
	pol := portpolicy.New(portpolicy.Allow)
	if got := pol.Evaluate(9999, "tcp"); got != portpolicy.Allow {
		t.Fatalf("expected Allow, got %s", got)
	}
}

func TestAllowRuleMatchesExactPort(t *testing.T) {
	pol := portpolicy.New(portpolicy.Deny)
	pol.Add(portpolicy.Rule{Port: 22, Proto: "tcp", Action: portpolicy.Allow})
	if got := pol.Evaluate(22, "tcp"); got != portpolicy.Allow {
		t.Fatalf("expected Allow, got %s", got)
	}
}

func TestRuleProtoMismatchFallsThrough(t *testing.T) {
	pol := portpolicy.New(portpolicy.Deny)
	pol.Add(portpolicy.Rule{Port: 53, Proto: "tcp", Action: portpolicy.Allow})
	// UDP 53 should fall through to default (Deny)
	if got := pol.Evaluate(53, "udp"); got != portpolicy.Deny {
		t.Fatalf("expected Deny for udp/53, got %s", got)
	}
}

func TestEmptyProtoMatchesBoth(t *testing.T) {
	pol := portpolicy.New(portpolicy.Deny)
	pol.Add(portpolicy.Rule{Port: 80, Proto: "", Action: portpolicy.Allow})
	for _, proto := range []string{"tcp", "udp"} {
		if got := pol.Evaluate(80, proto); got != portpolicy.Allow {
			t.Fatalf("expected Allow for %s/80, got %s", proto, got)
		}
	}
}

func TestFirstRuleWins(t *testing.T) {
	pol := portpolicy.New(portpolicy.Allow)
	pol.Add(portpolicy.Rule{Port: 443, Proto: "tcp", Action: portpolicy.Deny})
	pol.Add(portpolicy.Rule{Port: 443, Proto: "tcp", Action: portpolicy.Allow})
	if got := pol.Evaluate(443, "tcp"); got != portpolicy.Deny {
		t.Fatalf("expected first rule (Deny) to win, got %s", got)
	}
}

func writePolicyFile(t *testing.T, v any) string {
	t.Helper()
	b, _ := json.Marshal(v)
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.json")
	if err := os.WriteFile(p, b, 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadFileValid(t *testing.T) {
	path := writePolicyFile(t, map[string]any{
		"default": "deny",
		"rules": []map[string]any{
			{"port": 22, "proto": "tcp", "action": "allow", "comment": "ssh"},
		},
	})
	pol, err := portpolicy.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := pol.Evaluate(22, "tcp"); got != portpolicy.Allow {
		t.Fatalf("expected Allow for ssh, got %s", got)
	}
}

func TestLoadFileMissing(t *testing.T) {
	_, err := portpolicy.LoadFile("/no/such/file.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFileUnknownAction(t *testing.T) {
	path := writePolicyFile(t, map[string]any{
		"default": "permit", // invalid
		"rules":   []any{},
	})
	_, err := portpolicy.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}
