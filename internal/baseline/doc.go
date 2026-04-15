// Package baseline provides persistence for the approved set of port listeners
// in portwatch.
//
// A Baseline represents the "known-good" snapshot of ports that are expected to
// be open on the host. When the monitor detects a new listener it can be
// compared against the baseline to decide whether an alert should fire.
//
// Usage:
//
//	// Load existing baseline from disk (or start empty if the file is absent).
//	b, err := baseline.Load("/var/lib/portwatch/baseline.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Check whether a newly observed port is already approved.
//	if !b.Contains(portState) {
//		alerter.Notify(portState)
//	}
//
//	// Approve the port and persist the updated baseline.
//	b.Add(portState)
//	if err := b.Save(); err != nil {
//		log.Printf("warning: could not save baseline: %v", err)
//	}
//
// The JSON file is written with mode 0600 to prevent other users from
// tampering with the approved list.
package baseline
