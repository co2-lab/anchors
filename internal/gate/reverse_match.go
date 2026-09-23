package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- THE RETURN LEG OF EACH TRIAD PAIR ---
//
// `spec-feature-match` asks "does every rule have a scenario?" and `feature-test-match`
// asks "does every scenario have a test?". Both walk the ORIGIN looking for the
// destination, and neither walks the destination asking whether the origin still exists.
//
// The asymmetry has a measured cost. A revert deleted the rule `DTSTD-B10` — spec, code
// and test —, and the `.feature` kept its scenario: at the point of the revert that
// scenario did not exist yet, and a parallel PR reintroduced it with no git conflict.
// Probed in this repository, against the exact case:
//
//	spec→feature   orphan scenario (B10, no rule)   Pass, EMPTY message
//	feature→test   test proving the absent B10      not even mentioned
//
// The scenario stayed in the feature asserting a behaviour the spec no longer decides, and
// the test stayed green proving a rule nobody declares. Every gate green.
//
// It is the same asymmetric pair the documentation stamp and `realizes` had already shown
// — "if the spec has a hash and the doc does not reference it, something broke; and vice
// versa". The return leg is a QUESTION OF ITS OWN, and so it gets a gate of its own instead
// of becoming one more verdict of the forward gate: whoever reads "1 scenario without a
// rule" needs to know the accusation is about the feature, not about the spec.

// withoutVariant strips the VARIANT suffix from a scenario code.
//
// A feature numbers variants of the same requirement — `@DTTBD-B01#01`, `#02` — when a
// rule needs more than one scenario to be exercised. The spec declares the rule ONCE
// (`DTTBD-B01`), and it is the rule that decides; variants are how whoever writes the
// scenario slices it.
//
// Measured in the reference app: without this, the gate reported 69 features, ALL because
// of variants — eight scenarios of a feature whose rule exists and is declared. A gate that
// accuses what is right teaches people to ignore it, and would have buried the real case
// (the reverted `DTSTD-B10`) in the noise.
func withoutVariant(code string) string {
	if i := strings.IndexByte(code, '#'); i >= 0 {
		return code[:i]
	}
	return code
}

// isRevision reports whether a matched code is really the prefix of a REVISION
// (`-R0002`).
//
// `anyCodeRE` reads two digits after the letter; a revision has four. Without this check
// `JDDTJ-R0002` arrives as `JDDTJ-R00`, a code that exists nowhere.
func isRevision(body, matched string) bool {
	for _, i := range indicesOf(body, matched) {
		end := i + len(matched)
		if end < len(body) && body[end] >= '0' && body[end] <= '9' {
			return true
		}
	}
	return false
}

// indicesOf returns every position of `sub` in `s`.
func indicesOf(s, sub string) []int {
	var out []int
	for i := 0; ; {
		j := strings.Index(s[i:], sub)
		if j < 0 {
			return out
		}
		out = append(out, i+j)
		i += j + 1
	}
}

// dataStateDefRE finds where a spec DEFINES a data state: `DS-<name>` at the start of a
// heading, a bullet or a table's first cell — the same positions that define a rule. The
// unit prefix is optional, because specs write the short form inside their own unit.
var dataStateDefRE = regexp.MustCompile("(?m)^\\s*(?:#{2,6}\\s+|[-*]\\s+\\**|\\|\\s*)`?\\*{0,2}(?:[A-Z0-9]{3,6}-)?(DS-[A-Za-z0-9-]+)")

// dataStatesOf returns the data states a spec defines, in short form (`DS-data-present`).
//
// `definedRequirements` only knows `-{letter}{NN}` rules, and a data state is not one: it
// is named, not numbered. Measured in a project that defines states in a table
// (`| DS-data-present | … |`) and cites them in the feature as `@TREX-DS-data-present`:
// 94 of 103 `feature-spec-match` failures were this — states the spec DID define.
func dataStatesOf(content string) map[string]bool {
	out := map[string]bool{}
	for _, m := range dataStateDefRE.FindAllStringSubmatch(content, -1) {
		out[strings.TrimRight(m[1], "-")] = true
	}
	return out
}

// visualRegression reports whether a scenario code is the unit's VISUAL baseline
// (`TREX-VR`, `TREX-VR-<state>`).
//
// It is not a rule the spec defines — it is the screen's picture, and `vr-baseline` is the
// gate that charges it (the flow exists, the baseline image exists). Charging it here too
// would report every screen that adopts visual regression as an orphan.
func visualRegression(code string) bool {
	_, rest, ok := strings.Cut(code, "-")
	return ok && (rest == "VR" || strings.HasPrefix(rest, "VR-"))
}

// checkFeatureSpecMatch — the return leg of `spec-feature-match`: does every SCENARIO of
// the feature correspond to a rule the spec still DEFINES?
//
// The confrontation is about the feature, and it is the target: what was left orphaned is
// the scenario.
func checkFeatureSpecMatch(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, i18n.T("gate.feature_spec.skip_not_feature")
	}
	if g == nil {
		return pendingNoMap()
	}
	scenarios := parseFeatureScenarios(content)
	if len(scenarios) == 0 {
		return Skip, i18n.T("gate.feature_spec.skip_no_scenarios")
	}

	// The specs THIS feature covers: the `covered-by` edge leaves the spec and arrives here.
	var specPaths []string
	for _, e := range g.Neighbors(n.ID).In {
		if e.Type == mapx.EdgeCoveredBy {
			specPaths = append(specPaths, e.From)
		}
	}
	if len(specPaths) == 0 {
		// With no linked spec there is nothing to confront, and whoever charges the spec's
		// EXISTENCE is co-location. Pending, so the accusation is not duplicated.
		return Pending, i18n.T("gate.feature_spec.pending_no_spec")
	}

	// The union of what the linked specs define. Union and not intersection: a feature may
	// cover more than one spec, and a scenario matching ANY of them has an owner.
	declared := map[string]bool{}
	states := map[string]bool{}
	for _, sp := range specPaths {
		b, err := os.ReadFile(filepath.Join(root, sp))
		if err != nil {
			continue
		}
		for _, c := range definedRequirements(string(b)) {
			declared[c] = true
		}
		for st := range dataStatesOf(string(b)) {
			states[st] = true
		}
	}
	if len(declared) == 0 && len(states) == 0 {
		return Pending, i18n.T("gate.feature_spec.pending_no_requirements")
	}

	var orphans []string
	seen := map[string]bool{}
	for _, sc := range scenarios {
		for _, c := range sc.Codes {
			// ONLY the unit's own codes. A scenario may cite another unit's rule to say
			// what it runs against, and charging that to the local spec would have the gate
			// ask the impossible — the mistake `scenario-coverage` already measured, with 18
			// scenarios charged to a spec that defined 6.
			rule := withoutVariant(c)
			if n.Code != "" && !strings.HasPrefix(rule, n.Code+"-") {
				continue
			}
			if visualRegression(rule) {
				continue
			}
			// A data state is declared by its NAME, with or without the unit prefix.
			if _, short, ok := strings.Cut(rule, "-"); ok && strings.HasPrefix(short, "DS-") && states[short] {
				continue
			}
			if declared[rule] || seen[rule] {
				continue
			}
			seen[rule] = true
			// The code AS WRITTEN in the scenario, so whoever reads the verdict can find it.
			orphans = append(orphans, c)
		}
	}
	if len(orphans) == 0 {
		return Pass, ""
	}
	sort.Strings(orphans)
	return Fail, fmt.Sprintf(i18n.T("gate.feature_spec.orphan"),
		len(orphans), strings.Join(orphans, ", "))
}

// checkTestFeatureMatch — the return leg of `feature-test-match`: does every CODE the test
// claims to prove correspond to a scenario the feature still declares?
//
// The target is the TEST, and the difference matters: a green test proving a reverted rule
// is worse than a missing test, because it ATTESTS. The suite passes, coverage rises, and
// the number says a behaviour is proven when nobody decides it any more.
func checkTestFeatureMatch(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest {
		return Skip, i18n.T("gate.test_feature.skip_not_test")
	}
	if g == nil {
		return pendingNoMap()
	}

	// The features THIS test exercises: `tested-by` leaves the feature and arrives here.
	var featPaths []string
	for _, e := range g.Neighbors(n.ID).In {
		if e.Type == mapx.EdgeTestedBy {
			featPaths = append(featPaths, e.From)
		}
	}
	if len(featPaths) == 0 {
		return Pending, i18n.T("gate.test_feature.pending_no_feature")
	}

	declared := map[string]bool{}
	units := map[string]bool{}
	for _, fp := range featPaths {
		b, err := os.ReadFile(filepath.Join(root, fp))
		if err != nil {
			continue
		}
		for _, sc := range parseFeatureScenarios(string(b)) {
			for _, c := range sc.Codes {
				declared[withoutVariant(c)] = true
				if u, _, ok := strings.Cut(c, "-"); ok {
					units[u] = true
				}
			}
		}
	}
	if len(declared) == 0 {
		return Pending, i18n.T("gate.test_feature.pending_no_scenarios")
	}

	// COMMENTS OUT, by the same ruler as `feature-test-match`: a code cited in a comment is
	// a REFERENCE, not proof.
	body := stripLineComments(content)

	var orphans []string
	seen := map[string]bool{}
	for _, m := range anyCodeRE.FindAllString(body, -1) {
		// A REVISION IS NOT A RULE, and no scenario is charged for it.
		//
		// `JDDTJ-R0002` is a revision — four digits —, and `anyCodeRE` matches `R00`
		// because `R` is one of the canonical letters (for Rule) and the pattern reads two
		// digits. The rest is left over, and the gate reported a code nobody wrote.
		//
		// Measured in the reference app: two of the three findings were this. A scenario
		// exercises the behaviour the revision DECIDED — the revised rule —, and that is the
		// one the feature declares.
		if isRevision(body, m) {
			continue
		}
		// Only the units these features govern: a test cites other units' codes when it
		// builds fixtures, and charging them here would ask the local feature to declare
		// someone else's scenario.
		rule := withoutVariant(m)
		u, _, ok := strings.Cut(rule, "-")
		if !ok || !units[u] {
			continue
		}
		if declared[rule] || seen[rule] {
			continue
		}
		seen[rule] = true
		orphans = append(orphans, m)
	}
	if len(orphans) == 0 {
		return Pass, ""
	}
	sort.Strings(orphans)
	return Fail, fmt.Sprintf(i18n.T("gate.test_feature.orphan"),
		len(orphans), strings.Join(orphans, ", "))
}
