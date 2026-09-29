package gate

import (
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// DOES THE TEST REACH THE UNIT IT SAYS IT TESTS?
//
// A test can carry every code, match every title and assert — over a copy. The unit's
// function pasted into the test file and exercised there; a test that imports a neighbour
// instead of its unit; an end-to-end test whose `ref:` names one handler and invokes
// another. The gates that read titles see nothing wrong, and the unit can change freely.
//
// Two gates ask it, on the same routine:
//
//	test-exercises-unit    the unit the derivation pairs with the test (`TestedUnits`)
//	test-ref-matches-unit  the unit the test's `ref:` names
//
// They are separate because the project decides per gate what blocks: a copied function
// is a defect anywhere, a `ref:` reached by a call Anchors cannot see is not.
//
// Reaching is read without knowing the language. The test reaches a unit when:
//   - an import line of the test names the unit's module (its file name, or its directory);
//   - the test names something the unit DEFINES, which the dialect's `definition` says how
//     to find (a function, a type, a class), and does not define it itself;
//   - (for the `ref:`) a call the project declares on the gate entry, `invocations`, names
//     the unit — how a test reaches a lambda, a route, a job.
//
// And the test COPIES the unit when it defines a name the unit defines: the test is then
// exercising its own version.

// minDefinedName is the shortest defined name read as the unit's own: a one- or two-letter
// name (`T`, `ok`) is anyone's.
const minDefinedName = 3

// definedNames are the names the text defines, by the dialect's `definition`: the first
// non-empty capture group of each match.
func definedNames(text string, def *regexp.Regexp) map[string]bool {
	out := map[string]bool{}
	for _, m := range def.FindAllStringSubmatch(text, -1) {
		for _, g := range m[1:] {
			if g != "" {
				if len(g) >= minDefinedName {
					out[g] = true
				}
				break
			}
		}
	}
	return out
}

// reach is how a test reaches a unit, and what it copies of it.
type reach struct {
	Reached bool
	Copied  []string
}

// reachOf reads whether `test` reaches the unit at `unit` (a path relative to root), by
// import, by a name the unit defines, or by one of the declared invocations.
func reachOf(test, unit, root string, cfg *config.Config, def *regexp.Regexp, invocations []*regexp.Regexp) reach {
	var r reach
	if importsModule(test, unit, cfg) || invokes(test, unit, invocations) {
		r.Reached = true
	}
	if def == nil {
		return r
	}
	b, err := readFile(root, unit)
	if err != nil {
		return r // @resilient: an unreadable unit defines nothing this reading can name; the map notices a missing file
	}
	own := definedNames(test, def)
	for name := range definedNames(string(b), def) {
		if own[name] {
			r.Copied = append(r.Copied, name)
			continue
		}
		if !r.Reached && regexp.MustCompile(`\b`+regexp.QuoteMeta(name)+`\b`).MatchString(test) {
			r.Reached = true
		}
	}
	sort.Strings(r.Copied)
	return r
}

// importsModule says whether an import line of the test names the unit's module: its file
// name without extension, or its directory as a path (`internal/gate`), which is how a
// package-per-directory language imports it.
func importsModule(test, unit string, cfg *config.Config) bool {
	if importsUnit(test, unit, cfg) {
		return true
	}
	dir := path.Dir(filepath.ToSlash(unit))
	if dir == "." || dir == "" {
		return false
	}
	for _, l := range importLines(test, cfg) {
		if strings.Contains(l, dir) {
			return true
		}
	}
	return false
}

// invokes says whether a declared invocation in the test names the unit: its first
// non-empty capture is the unit's file name without extension, or one of its directories.
func invokes(test, unit string, invocations []*regexp.Regexp) bool {
	if len(invocations) == 0 {
		return false
	}
	u := filepath.ToSlash(unit)
	names := map[string]bool{strings.TrimSuffix(path.Base(u), path.Ext(u)): true}
	for _, seg := range strings.Split(path.Dir(u), "/") {
		if seg != "." && seg != "" {
			names[seg] = true
		}
	}
	for _, re := range invocations {
		for _, m := range re.FindAllStringSubmatch(test, -1) {
			for _, g := range m[1:] {
				if g != "" {
					if names[g] {
						return true
					}
					break
				}
			}
		}
	}
	return false
}

// definitionOf is the dialect's `definition`, compiled; nil when the project has none.
func definitionOf(cfg *config.Config) *regexp.Regexp {
	if cfg == nil {
		return nil
	}
	d := cfg.DialectFor()
	if strings.TrimSpace(d.Definition) == "" {
		return nil
	}
	re, err := regexp.Compile(d.Definition)
	if err != nil {
		return nil // @resilient: the configuration's load refuses a pattern that does not compile
	}
	return re
}

// ── test-exercises-unit ────────────────────────────────────────────────────

// noUnitImportRE is the waiver of a test that does not reach its unit the way this gate
// reads, with its reason: `@no-unit-import: <why>`.
var noUnitImportRE = regexp.MustCompile(`@no-unit-import[^\S\n]*:[^\S\n]*\S+`)

func checkTestExercisesUnit(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest || n.Support {
		return Skip, i18n.T("gate.test_reach.skip_not_test")
	}
	units := testedUnits(g, cfg)[n.ID]
	if len(units) == 0 {
		return Skip, i18n.T("gate.test_exercises_unit.skip_no_unit")
	}
	if noUnitImportRE.MatchString(content) {
		return Pass, "" // declared: the test reaches its unit another way, and says why
	}
	def := definitionOf(cfg)
	if def == nil {
		return Skip, i18n.T("gate.test_reach.skip_no_definition")
	}
	var problems []string
	for _, u := range units {
		r := reachOf(content, u, root, cfg, def, nil)
		if len(r.Copied) > 0 {
			problems = append(problems, i18n.T("gate.test_exercises_unit.copied", u, strings.Join(r.Copied, ", ")))
		}
		if !r.Reached {
			problems = append(problems, i18n.T("gate.test_exercises_unit.not_reached", u))
		}
	}
	if len(problems) > 0 {
		return Fail, strings.Join(problems, "; ")
	}
	return Pass, ""
}

// ── test-ref-matches-unit ──────────────────────────────────────────────────

func checkTestRefMatchesUnit(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest || n.Support {
		return Skip, i18n.T("gate.test_reach.skip_not_test")
	}
	m := refHeaderRE().FindStringSubmatch(content)
	if m == nil {
		return Skip, i18n.T("gate.test_ref_matches_unit.skip_no_ref")
	}
	units := unitFiles(g, m[1])
	if len(units) == 0 {
		return Skip, i18n.T("gate.test_ref_matches_unit.skip_no_code", m[1])
	}
	var invocations []*regexp.Regexp
	for _, p := range gateEntry(cfg, "test-ref-matches-unit").Invocations {
		if re, err := regexp.Compile(p); err == nil {
			invocations = append(invocations, re)
		}
	}
	def := definitionOf(cfg)
	for _, u := range units {
		if reachOf(content, u, root, cfg, def, invocations).Reached {
			return Pass, ""
		}
	}
	sort.Strings(units)
	return Fail, i18n.T("gate.test_ref_matches_unit.fail", m[1], strings.Join(units, ", "))
}
