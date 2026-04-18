// Package tunables provides a thread-safe container for runtime-adjustable
// operational parameters used across portwatch subsystems.
//
// Parameters include the port scan interval, alert cooldown window, and the
// maximum number of history entries to retain in memory. All getters and
// setters are safe for concurrent use from multiple goroutines.
//
// Typical usage:
//
//	tn := tunables.Defaults()
//	tn.SetScanInterval(10 * time.Second)
//
//	// later, in the scan loop:
//	time.Sleep(tn.ScanInterval())
package tunables
