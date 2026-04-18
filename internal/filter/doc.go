// Package filter implements pre-alert suppression logic for portwatch.
//
// A Filter wraps a [config.Config] and exposes a single Apply method that
// strips unwanted [scanner.PortState] entries from a slice before the
// remaining states are forwarded to the alerter or baseline comparison.
//
// Suppression rules (evaluated in order):
//
//  1. Loopback addresses — dropped when Config.IgnoreLoopback is true.
//  2. Explicitly ignored ports — dropped when the port appears in
//     Config.IgnoredPorts.
//
// The filter is stateless and safe for concurrent use; multiple goroutines
// may call Apply simultaneously on the same Filter instance.
//
// Usage:
//
//	f := filter.New(cfg)
//	visible := f.Apply(rawStates)
package filter
