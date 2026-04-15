// Package tagger maps port numbers to human-readable service names.
//
// It ships with a built-in table of well-known ports (SSH, HTTP, MySQL, …) and
// accepts caller-supplied custom mappings that take precedence over the
// defaults.  Tags are purely informational — they are attached to alert
// messages and audit log entries to make output easier to read at a glance.
//
// Usage:
//
//	t := tagger.New(map[uint16]string{9200: "elasticsearch"})
//	name := t.Tag(state)   // e.g. "http" for port 80
package tagger
