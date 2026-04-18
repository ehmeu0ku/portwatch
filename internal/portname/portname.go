// Package portname resolves well-known port numbers to human-readable service names.
package portname

import "fmt"

// wellKnown maps port numbers to canonical service names.
var wellKnown = map[uint16]string{
	20:   "ftp-data",
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	9200: "elasticsearch",
	27017: "mongodb",
}

// Resolver resolves port numbers to service names.
type Resolver struct {
	custom map[uint16]string
}

// New returns a Resolver with an optional set of custom mappings that
// override or extend the built-in well-known table.
func New(custom map[uint16]string) *Resolver {
	return &Resolver{custom: custom}
}

// Lookup returns the service name for the given port, or a numeric
// fallback string if the port is not recognised.
func (r *Resolver) Lookup(port uint16) string {
	if r.custom != nil {
		if name, ok := r.custom[port]; ok {
			return name
		}
	}
	if name, ok := wellKnown[port]; ok {
		return name
	}
	return fmt.Sprintf("port-%d", port)
}

// Known returns true when the port has a recognised service name.
func (r *Resolver) Known(port uint16) bool {
	if r.custom != nil {
		if _, ok := r.custom[port]; ok {
			return true
		}
	}
	_, ok := wellKnown[port]
	return ok
}
