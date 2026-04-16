package retention_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/audit"
	"github.com/yourorg/portwatch/internal/retention"
)

func TestAuditStoreRoundTrip(t *testing.T) {
	s := retention.NewAuditStore(time.Minute)
	e := audit.Entry{Timestamp: time.Now(), Kind: "new", Proto: "tcp", Port: 8080}
	s.AddAudit(e)
	entries := s.AuditEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 8080 {
		t.Errorf("unexpected port: %d", entries[0].Port)
	}
}

func TestAuditStoreEvictsExpired(t *testing.T) {
	s := retention.NewAuditStore(time.Minute)
	old := audit.Entry{Timestamp: time.Now().Add(-2 * time.Minute), Kind: "new", Proto: "tcp", Port: 9090}
	recent := audit.Entry{Timestamp: time.Now(), Kind: "new", Proto: "tcp", Port: 443}
	s.AddAudit(old)
	s.AddAudit(recent)
	entries := s.AuditEntries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 live entry, got %d", len(entries))
	}
	if entries[0].Port != 443 {
		t.Errorf("expected port 443, got %d", entries[0].Port)
	}
}
