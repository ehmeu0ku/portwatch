// Package labelstore provides a thread-safe, JSON-backed registry that maps
// (proto, port) keys to operator-defined labels.
//
// Labels are purely informational annotations; they are attached to events by
// the correlator pipeline so that downstream formatters and notifiers can
// surface human-readable names alongside raw port numbers.
//
// Usage:
//
//	store := labelstore.New("/var/lib/portwatch/labels.json")
//	_ = store.Load()          // ignore error when file does not yet exist
//	store.Set(labelstore.Key{Proto: "tcp", Port: 443}, labelstore.Label{Name: "https"})
//	_ = store.Save()
package labelstore
