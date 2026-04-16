package retention

import (
	"time"

	"github.com/yourorg/portwatch/internal/audit"
)

// auditEntry wraps audit.Entry so it satisfies the Entry interface.
type auditEntry struct{ inner audit.Entry }

func (a auditEntry) Timestamp() time.Time { return a.inner.Timestamp }

// AuditStore is a retention.Store pre-configured for audit.Entry values.
type AuditStore struct{ *Store }

// NewAuditStore returns an AuditStore with the given TTL.
func NewAuditStore(ttl time.Duration) *AuditStore {
	return &AuditStore{New(ttl)}
}

// AddAudit wraps e and appends it to the store.
func (a *AuditStore) AddAudit(e audit.Entry) {
	a.Add(auditEntry{inner: e})
}

// AuditEntries returns all live audit entries.
func (a *AuditStore) AuditEntries() []audit.Entry {
	raw := a.Entries()
	out := make([]audit.Entry, len(raw))
	for i, r := range raw {
		out[i] = r.(auditEntry).inner
	}
	return out
}
