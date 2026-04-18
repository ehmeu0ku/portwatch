// Package statestore provides a simple persistent store for scanner.PortState
// values. It serialises the current port snapshot to a JSON file so that the
// monitor can compare the live scan against the previous run and detect
// changes that occurred while portwatch was not running.
//
// Typical usage:
//
//	store, err := statestore.Load("/var/lib/portwatch/state.json")
//	prev := store.Get()
//	// … run scanner …
//	_ = store.Set(current)
package statestore
