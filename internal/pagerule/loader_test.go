package pagerule_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/pagerule"
	"github.com/user/portwatch/internal/severity"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "rules.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadFileValid(t *testing.T) {
	path := writeRulesFile(t, `[
		{"min_severity":"critical","action":"page"},
		{"min_severity":"info","ports":[80,443],"tags":["http"],"action":"suppress"}
	]`)
	rules, err := pagerule.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].MinSeverity != severity.Critical {
		t.Errorf("rule 0 severity mismatch")
	}
	if rules[1].Action != pagerule.ActionSuppress {
		t.Errorf("rule 1 action mismatch")
	}
}

func TestLoadFileMissing(t *testing.T) {
	_, err := pagerule.LoadFile("/nonexistent/rules.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFileUnknownAction(t *testing.T) {
	path := writeRulesFile(t, `[{"min_severity":"info","action":"explode"}]`)
	_, err := pagerule.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestLoadFileBadSeverity(t *testing.T) {
	path := writeRulesFile(t, `[{"min_severity":"ultra","action":"page"}]`)
	_, err := pagerule.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for bad severity")
	}
}
