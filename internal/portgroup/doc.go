// Package portgroup provides a named-group abstraction over sets of port
// numbers. Groups can be defined at startup from configuration and queried
// at alert time to attach human-readable context such as "database",
// "web", or "internal-services" to a detected port event.
package portgroup
