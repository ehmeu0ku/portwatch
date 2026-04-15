// Package audit implements a persistent, append-only audit log for portwatch.
//
// Each port event (new listener detected or listener gone) is serialised as a
// JSON Lines record and appended to a file on disk.  The log can be replayed
// with ReadAll for offline analysis or reporting.
//
// Usage:
//
//	log, err := audit.Open("/var/log/portwatch/audit.jsonl")
//	if err != nil { ... }
//	defer log.Close()
//
//	log.Record(audit.Entry{
//	    Kind:  audit.EventNew,
//	    Proto: "tcp",
//	    Port:  8080,
//	    Addr:  "0.0.0.0",
//	})
package audit
