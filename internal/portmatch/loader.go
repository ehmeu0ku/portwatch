package portmatch

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// exprJSON is the on-disk representation of an Expr.
type exprJSON struct {
	Port    uint16 `json:"port,omitempty"`
	PortMin uint16 `json:"port_min,omitempty"`
	PortMax uint16 `json:"port_max,omitempty"`
	Proto   string `json:"proto,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

// LoadFile reads a JSON array of match expressions from path and returns
// a Matcher. Returns an error if the file cannot be read or parsed.
func LoadFile(path string) (*Matcher, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("portmatch: read %s: %w", path, err)
	}
	var raw []exprJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("portmatch: parse %s: %w", path, err)
	}
	exprs := make([]Expr, 0, len(raw))
	for i, r := range raw {
		proto := strings.ToLower(r.Proto)
		if proto != "" && proto != "tcp" && proto != "udp" {
			return nil, fmt.Errorf("portmatch: entry %d: unknown proto %q", i, r.Proto)
		}
		if r.PortMin > r.PortMax && r.PortMax != 0 {
			return nil, fmt.Errorf("portmatch: entry %d: port_min > port_max", i)
		}
		exprs = append(exprs, Expr{
			Port:    r.Port,
			PortMin: r.PortMin,
			PortMax: r.PortMax,
			Proto:   proto,
			Tag:     r.Tag,
		})
	}
	return New(exprs), nil
}
