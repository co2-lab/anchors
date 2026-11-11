// @anchors
//   code: EXPRC
//   ref: EXPRS

package gate

import (
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// external-proven: a code file that calls a real system the project declares
// (`external_systems:`) is proven by a test of its own that passed in a suite reaching that
// system.
//
// Measured in MIF: fourteen defects crossed a boundary — an index that did not exist, a
// condition outside the key schema, an owner field the API refused —, each proven green by
// unit tests against a double that accepted the wrong thing. The unit tests stay right about
// the unit; what no gate asked is whether any proof of the file ever met the real system.
//
// A file waives it on a line of its own with `@no-external-proof: <reason>` — a call kept
// behind a flag, a system only the CI may reach.
func checkExternalProven(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if cfg == nil || len(cfg.ExternalSystems) == 0 {
		return Skip, i18n.T("gate.external_proven.skip_none_declared")
	}
	if n.Kind != mapx.KindCode {
		return Skip, i18n.T("gate.external_proven.skip_not_code")
	}
	var called []config.ExternalSystem
	for _, s := range cfg.ExternalSystems {
		if re := s.CallsRE(); re != nil && re.MatchString(content) {
			called = append(called, s)
		}
	}
	if len(called) == 0 {
		return Skip, i18n.T("gate.external_proven.skip_no_call")
	}
	if noExternalProofRE.MatchString(content) {
		return Skip, i18n.T("gate.external_proven.skip_waived")
	}
	if g == nil {
		return pendingNoMap()
	}
	tests := testsOfUnit(n.ID, g, cfg)
	var missing []string
	for _, s := range called {
		if !provenIn(tests, g, s.ProvenBy) {
			missing = append(missing, i18n.T("gate.external_proven.unproven", s.Name, strings.Join(s.ProvenBy, ", ")))
		}
	}
	if len(missing) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(missing, "; ") + i18n.T("gate.external_proven.tests_seen", len(tests))
}

// noExternalProofRE is the file's waiver, with its reason.
var noExternalProofRE = regexp.MustCompile(`@no-external-proof[^\S\n]*:[^\S\n]*\S`)

// testsOfUnit are the tests of a code file: those the project's derivation places beside it,
// those that reach it by an edge of their own (an import), and those of the feature that
// covers its spec.
func testsOfUnit(id string, g *mapx.Graph, cfg *config.Config) []string {
	seen := map[string]bool{}
	add := func(t string) {
		if n := g.Node(t); n != nil && n.Kind == mapx.KindTest {
			seen[t] = true
		}
	}
	for t, units := range testedUnits(g, cfg) {
		for _, u := range units {
			if u == id {
				add(t)
			}
		}
	}
	for _, e := range g.Neighbors(id).In {
		add(e.From)
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		for _, cov := range g.Neighbors(e.From).Out {
			if cov.Type != mapx.EdgeCoveredBy {
				continue
			}
			for _, tb := range g.Neighbors(cov.To).Out {
				if tb.Type == mapx.EdgeTestedBy {
					add(tb.To)
				}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// provenIn says whether one of the tests passed in a suite of one of the layers.
func provenIn(tests []string, g *mapx.Graph, layers []string) bool {
	for _, t := range tests {
		n := g.Node(t)
		if n == nil || n.Signal == nil {
			continue
		}
		for _, l := range layers {
			if ex, ok := n.Signal.ByLayer[l]; ok && ex.Passed > 0 {
				return true
			}
		}
	}
	return false
}
