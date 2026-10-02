// @anchors
//   ref: TLVCD

package gate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// checkTestLevelCodes — does each scenario of the feature reference only codes its test
// level accepts?
//
// A scenario is tagged with its test level and references rules of the spec, so the
// project can say which codes belong to which level (`levels`, on the gate's own entry): a visual
// level that accepts only `-VR` codes, a unit level that refuses them. Without the
// declaration every level accepts every code and the gate has nothing to confront.
//
// It is the general form of a project's naming rule between level and code. Measured in
// the reference app: 9 scenarios tagged with its visual level had no `-VR` code, so the
// gate that asks for a baseline image never asked, and 1 `-VR` code sat under the
// integration level.
func checkTestLevelCodes(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, i18n.T("gate.test_level_codes.skip_not_feature")
	}
	levels := declaredLevels(cfg)
	if len(levels) == 0 {
		return Skip, i18n.T("gate.test_level_codes.skip_not_declared")
	}
	var refused []string
	confronted := false
	for _, s := range parseFeatureScenarios(content) {
		for _, tag := range s.Tags {
			level, ok := levels[tag]
			if !ok {
				continue
			}
			confronted = true
			for _, c := range s.Codes {
				if code := RootCode(c); !level.Accepts(code) {
					refused = append(refused, fmt.Sprintf("%s (@%s)", code, tag))
				}
			}
		}
	}
	if !confronted {
		return Skip, i18n.T("gate.test_level_codes.skip_no_level")
	}
	if len(refused) == 0 {
		return Pass, ""
	}
	sort.Strings(refused)
	return Fail, i18n.T("gate.test_level_codes.refused", len(refused), strings.Join(refused, ", "))
}

// declaredLevels gathers the `levels` of every gate entry that runs this check. The option
// lives on the gate that reads it; two entries declaring the same level add their lists
// together, so neither silently drops the other's patterns.
func declaredLevels(cfg *config.Config) map[string]config.TestLevel {
	if cfg == nil {
		return nil
	}
	out := map[string]config.TestLevel{}
	for _, g := range cfg.Gates {
		if g.Check != "test-level-codes" {
			continue
		}
		for tag, l := range g.Levels {
			cur := out[tag]
			cur.Allow = append(cur.Allow, l.Allow...)
			cur.Exclude = append(cur.Exclude, l.Exclude...)
			out[tag] = cur
		}
	}
	return out
}
