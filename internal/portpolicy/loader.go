package portpolicy

import (
	"encoding/json"
	"fmt"
	"os"
)

type ruleJSON struct {
	Port    uint16 `json:"port"`
	Proto   string `json:"proto"`
	Action  string `json:"action"`
	Comment string `json:"comment"`
}

type policyJSON struct {
	Default string     `json:"default"`
	Rules   []ruleJSON `json:"rules"`
}

// LoadFile reads a JSON policy file and returns a populated Policy.
// The file format is:
//
//	{"default":"deny","rules":[{"port":22,"proto":"tcp","action":"allow"}]}
func LoadFile(path string) (*Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("portpolicy: open %s: %w", path, err)
	}
	defer f.Close()

	var pj policyJSON
	if err := json.NewDecoder(f).Decode(&pj); err != nil {
		return nil, fmt.Errorf("portpolicy: decode %s: %w", path, err)
	}

	defaultAction, err := parseAction(pj.Default)
	if err != nil {
		return nil, fmt.Errorf("portpolicy: default action: %w", err)
	}

	pol := New(defaultAction)
	for i, rj := range pj.Rules {
		act, err := parseAction(rj.Action)
		if err != nil {
			return nil, fmt.Errorf("portpolicy: rule %d: %w", i, err)
		}
		pol.Add(Rule{
			Port:    rj.Port,
			Proto:   rj.Proto,
			Action:  act,
			Comment: rj.Comment,
		})
	}
	return pol, nil
}

func parseAction(s string) (Action, error) {
	switch s {
	case "allow":
		return Allow, nil
	case "deny", "":
		return Deny, nil
	default:
		return Deny, fmt.Errorf("unknown action %q", s)
	}
}
