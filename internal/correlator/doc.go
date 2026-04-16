// Package correlator combines enrichment, tagging, and severity scoring
// into a single pipeline step.
//
// # Overview
//
// When the monitor detects a port-state change it hands the raw
// scanner.PortState to the Correlator.  The Correlator:
//
//  1. Enriches the state with process information via the enricher package.
//  2. Resolves a human-readable service tag via the tagger package.
//  3. Assigns a severity level via the severity package.
//
// The result is a PortEvent that downstream consumers (notifiers, audit
// recorders, the event bus) can use without knowing about any individual
// sub-system.
package correlator
