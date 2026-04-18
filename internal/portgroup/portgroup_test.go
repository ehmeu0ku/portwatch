package portgroup

import (
	"sort"
	"testing"
)

func TestDefineAndContains(t *testing.T) {
	s := New()
	s.Define("web", []uint16{80, 443, 8080})
	if !s.Contains("web", 80) {
		t.Fatal("expected port 80 in group web")
	}
	if s.Contains("web", 22) {
		t.Fatal("port 22 should not be in group web")
	}
}

func TestContainsMissingGroup(t *testing.T) {
	s := New()
	if s.Contains("nonexistent", 80) {
		t.Fatal("missing group should return false")
	}
}

func TestLookupReturnsAllMatchingGroups(t *testing.T) {
	s := New()
	s.Define("web", []uint16{80, 443})
	s.Define("proxy", []uint16{80, 3128})
	got := s.Lookup(80)
	sort.Strings(got)
	if len(got) != 2 || got[0] != "proxy" || got[1] != "web" {
		t.Fatalf("unexpected groups: %v", got)
	}
}

func TestLookupUnknownPortReturnsEmpty(t *testing.T) {
	s := New()
	s.Define("db", []uint16{5432, 3306})
	if got := s.Lookup(9999); len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestDefineReplacesExistingGroup(t *testing.T) {
	s := New()
	s.Define("web", []uint16{80})
	s.Define("web", []uint16{443})
	if s.Contains("web", 80) {
		t.Fatal("old port should be gone after redefinition")
	}
	if !s.Contains("web", 443) {
		t.Fatal("new port should be present after redefinition")
	}
}

func TestNamesReturnsAllGroups(t *testing.T) {
	s := New()
	s.Define("web", []uint16{80})
	s.Define("db", []uint16{5432})
	names := s.Names()
	sort.Strings(names)
	if len(names) != 2 || names[0] != "db" || names[1] != "web" {
		t.Fatalf("unexpected names: %v", names)
	}
}
