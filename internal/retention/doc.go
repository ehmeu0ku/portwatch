// Package retention provides a time-bounded in-memory store that automatically
// evicts entries older than a configured TTL.
//
// It is used by the audit and history subsystems to cap unbounded growth during
// long-running daemon sessions without requiring a separate background goroutine.
package retention
