// @anchors
//   code: SNGTS
//   ref: SNGTU

package gate

import (
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// single-test-per-unit: a unit has ONE test file per test layer.
//
// Two files testing one unit in one layer split what the unit is proven by, and each looks
// complete on its own: in the reference app a hook had `useIapNative.test.ts` and
// `useIapNative.test.tsx`, and five models had tests both under `__tests__/unit/models` and
// under `__tests__/unit/lambdas/models-*`. Which one to extend, which one a gate reads, which
// one a mutation run counts — each answer was a guess.
//
// The unit a test tests is the project's own derivation, read backwards (`TestedUnits`):
// nothing here assumes a language or a file layout. A split the project means — a stub and
// the native SDK need different mocks — is declared in one of the files with
// `@split-test: <why>`, and the unit passes.

var splitTestRE = regexp.MustCompile(`@split-test[^\S\n]*:[^\S\n]*\S`)

func checkSingleTestPerUnit(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode {
		return Skip, i18n.T("gate.single_test.skip_not_code")
	}
	if g == nil {
		return pendingNoMap()
	}
	layerOf := map[string]string{}
	for _, t := range g.Nodes {
		if t.Kind == mapx.KindTest {
			layerOf[t.ID] = t.Layer
		}
	}
	byLayer := map[string][]string{}
	for test, units := range testedUnits(g, cfg) {
		for _, u := range units {
			if u == n.ID {
				byLayer[layerOf[test]] = append(byLayer[layerOf[test]], test)
			}
		}
	}
	if len(byLayer) == 0 {
		return Skip, i18n.T("gate.single_test.skip_no_test")
	}
	var split []string
	for layer, tests := range byLayer {
		if len(tests) < 2 || declaresSplit(root, tests) {
			continue
		}
		sort.Strings(tests)
		split = append(split, i18n.T("gate.single_test.split_item", layer, strings.Join(tests, ", ")))
	}
	if len(split) == 0 {
		return Pass, ""
	}
	sort.Strings(split)
	return Fail, i18n.T("gate.single_test.fail_split", strings.Join(split, "; "))
}

// declaresSplit says whether one of the files declares the split, with its reason.
func declaresSplit(root string, tests []string) bool {
	for _, t := range tests {
		if b, err := readFile(root, t); err == nil && splitTestRE.Match(b) {
			return true
		}
	}
	return false
}
