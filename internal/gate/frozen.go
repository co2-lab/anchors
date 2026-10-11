// @anchors
//   code: FGCFR
//   ref: FRGTF

package gate

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// --- frozen parts: quiet, and told when touched (DESIGN-frozen-units.md) ---
//
// A frozen part is kept as it is by decision, until someone thaws it. Its gates do not fail
// and do not pend: each answers `frozen`, with the freeze's reason, counted apart. Its rules
// leave every question asked of a rule — a scenario, a test, a realizer —, the way a rule
// marked `@no-scenario` or `@retired` does. Two gates still speak about it: one confronts
// the marks with the records, the other tells a change that touches a frozen part.

// Frozen is the verdict on a frozen part: kept as it is by decision, neither a failure nor
// a pending item.
const Frozen Verdict = "frozen"

// freezeGates are the checks that speak about frozen parts, and so run on them.
var freezeGates = map[string]bool{"freeze-consistent": true, "frozen-touched": true}

// The run's freezes, set by RunWithWaiver from the map: which nodes are frozen whole and
// which rules a freeze covers, each with its value as it reads.
var (
	freezeMu     sync.RWMutex
	frozenNodes  map[string]string
	frozenRuleOf map[string]string
)

// setFreezes reads the run's freezes from the map; with no map, nothing is frozen.
func setFreezes(g *mapx.Graph) {
	nodes, rules := map[string]string{}, map[string]string{}
	if g != nil {
		for _, n := range g.Nodes {
			if n.Frozen != "" {
				nodes[n.ID] = g.DescribeFreeze(n.Frozen)
			}
		}
		for r, v := range g.FrozenRules() {
			rules[r] = g.DescribeFreeze(v)
		}
	}
	freezeMu.Lock()
	frozenNodes, frozenRuleOf = nodes, rules
	freezeMu.Unlock()
}

// frozenNode is the freeze that covers a node whole, as it reads, or empty.
func frozenNode(n mapx.Node) string {
	freezeMu.RLock()
	defer freezeMu.RUnlock()
	if v := frozenNodes[n.ID]; v != "" {
		return v
	}
	return n.Frozen
}

// ruleIsFrozen says whether a freeze covers the rule: it leaves every question asked of a
// rule.
func ruleIsFrozen(code string) bool {
	freezeMu.RLock()
	defer freezeMu.RUnlock()
	_, ok := frozenRuleOf[code]
	return ok
}

// frozenResult is a gate's verdict on a frozen node, when the gate is not one of the
// freeze's own.
func frozenResult(g config.Gate, n mapx.Node) (Result, bool) {
	if freezeGates[g.Check] || freezeGates[g.Name] {
		return Result{}, false
	}
	v := frozenNode(n)
	if v == "" {
		return Result{}, false
	}
	return Result{
		Gate: g.Name, Regra: idDoGate(g), Target: n.ID,
		Verdict: Frozen, Blocking: g.IsBlocking(),
		Detail: i18n.T("gate.frozen", v),
	}, true
}

// withoutFrozen are the targets of an aggregate gate a freeze does not cover: a tool asked
// about the set is not asked about a frozen file.
func withoutFrozen(g config.Gate, nodes []mapx.Node) []mapx.Node {
	if freezeGates[g.Check] || freezeGates[g.Name] {
		return nodes
	}
	out := nodes[:0:0]
	for _, n := range nodes {
		if frozenNode(n) == "" {
			out = append(out, n)
		}
	}
	return out
}

// --- freeze-consistent: a freeze is somewhere, with its reason, or nowhere ---
//
// A mark naming a code no doctrine records points at a reason nobody can read; a record no
// mark carries is a freeze that covers nothing — thawed by hand, its entry forgotten.
// Without a doctrine, the marks carry the reason, and there is nothing to confront.
func checkFreezeConsistent(_ string, _ mapx.Node, _ string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	if g == nil {
		return pendingNoMap()
	}
	records := g.FreezeRecords()
	used := map[string]bool{}
	var unknown []string
	for _, p := range g.FrozenParts("") {
		if !scan.IsFreezeCode(p.Value) {
			continue
		}
		used[p.Value] = true
		if _, ok := records[p.Value]; !ok {
			where := p.Path
			if p.Rule != "" {
				where += " " + p.Rule
			}
			unknown = append(unknown, where+" → "+p.Value)
		}
	}
	var orphan []string
	for c, r := range records {
		if !used[c] {
			orphan = append(orphan, c+" ("+r.Doc+")")
		}
	}
	sort.Strings(unknown)
	sort.Strings(orphan)
	var parts []string
	if len(unknown) > 0 {
		parts = append(parts, fmt.Sprintf(i18n.T("gate.freeze_consistent.unknown"), len(unknown), strings.Join(unknown, "; ")))
	}
	if len(orphan) > 0 {
		parts = append(parts, fmt.Sprintf(i18n.T("gate.freeze_consistent.orphan"), len(orphan), strings.Join(orphan, ", ")))
	}
	if len(parts) > 0 {
		return Fail, strings.Join(parts, " ")
	}
	if len(records) == 0 && len(used) == 0 {
		return Skip, i18n.T("gate.freeze_consistent.none")
	}
	return Pass, ""
}

// --- frozen-touched: a change to a frozen part is told, and goes on ---
//
// The order to change a frozen part comes from a person; the gate says the part is frozen,
// and why, to whoever is about to commit the change. It informs and never bars: the
// verdict is a divergence the project declared.
func checkFrozenTouched(_ string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	whole := frozenNode(n)
	var rules []string
	if whole == "" {
		for r, v := range n.FrozenRules {
			desc := v
			if g != nil {
				desc = g.DescribeFreeze(v)
			}
			rules = append(rules, r+" ("+desc+")")
		}
	}
	if whole == "" && len(rules) == 0 {
		return Skip, i18n.T("gate.frozen_touched.not_frozen")
	}
	changed, known := uncommitted(root, n.ID)
	if !known || !changed {
		return Pass, ""
	}
	if whole != "" {
		return Diverge, i18n.T("gate.frozen_touched.whole", whole) + " " + DeclaredMarker
	}
	sort.Strings(rules)
	return Diverge, i18n.T("gate.frozen_touched.rules", strings.Join(rules, "; ")) + " " + DeclaredMarker
}
