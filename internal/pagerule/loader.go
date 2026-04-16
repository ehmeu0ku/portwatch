package pagerule

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/portwatch/internal/severity"
)

type ruleJSON struct {
	MinSeverity string   `json:"min_severity"`
	Ports       []uint16 `json:"ports"`
	Tags        []string `json:"tags"`
	Action      string   `json:"action"`
}

// LoadFile reads a JSON array of rules from path.
func LoadFile(path string) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("pagerule: open %s: %w", path, err)
	}
	defer f.Close()

	var raw []ruleJSON
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, fmt.Errorf("pagerule: decode: %w", err)
	}

	rules := make([]Rule, 0, len(raw))
	for i, r := range raw {
		lvl, err := severity.Parse(r.MinSeverity)
		if err != nil {
			return nil, fmt.Errorf("pagerule: rule %d: %w", i, err)
		}
		action := Action(r.Action)
		if action != ActionPage && action != ActionSuppress {
			return nil, fmt.Errorf("pagerule: rule %d: unknown action %q", i, r.Action)
		}
		rules = append(rules, Rule{
			MinSeverity: lvl,
			Ports:       r.Ports,
			Tags:        r.Tags,
			Action:      action,
		})
	}
	return rules, nil
}
