package monitor

import "github.com/user/portwatch/internal/scanner"

// key uniquely identifies a port state for comparison purposes.
type key struct {
	Proto string
	Addr  string
	Port  uint16
}

func stateKey(s scanner.PortState) key {
	return key{Proto: s.Proto, Addr: s.Addr, Port: s.Port}
}

// diff compares two slices of PortState and returns:
//   - added: ports present in current but not in previous
//   - removed: ports present in previous but not in current
func diff(previous, current []scanner.PortState) (added, removed []scanner.PortState) {
	prevMap := make(map[key]struct{}, len(previous))
	for _, s := range previous {
		prevMap[stateKey(s)] = struct{}{}
	}

	currMap := make(map[key]struct{}, len(current))
	for _, s := range current {
		currMap[stateKey(s)] = struct{}{}
	}

	for _, s := range current {
		if _, ok := prevMap[stateKey(s)]; !ok {
			added = append(added, s)
		}
	}

	for _, s := range previous {
		if _, ok := currMap[stateKey(s)]; !ok {
			removed = append(removed, s)
		}
	}

	return added, removed
}
