package portfence

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// rule is the JSON schema for a single fence entry.
type rule struct {
	Proto string   `json:"proto"`
	Ports []uint16 `json:"ports"`
}

// LoadFile reads a JSON file containing an array of fence rules and populates
// a new Fence. The file format is:
//
//	[
//	  {"proto": "tcp", "ports": [22, 80, 443]},
//	  {"proto": "udp", "ports": [53]}
//	]
func LoadFile(path string) (*Fence, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("portfence: read %s: %w", path, err)
	}
	var rules []rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("portfence: parse %s: %w", path, err)
	}
	f := New()
	for _, r := range rules {
		proto := strings.ToLower(strings.TrimSpace(r.Proto))
		if proto != "tcp" && proto != "udp" {
			return nil, fmt.Errorf("portfence: unknown protocol %q in %s", r.Proto, path)
		}
		for _, p := range r.Ports {
			f.Allow(proto, p)
		}
	}
	return f, nil
}
