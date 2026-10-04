// @anchors
//   code: VSCGV
//   ref: VRSTC

package gate

import (
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// baselineExts are the image formats a visual baseline may be saved in.
const baselineExts = "{png,jpg,jpeg,webp,gif,svg}"

// The visual-regression gates ask four questions of a visual unit — a screen, a component:
//
//	vr-states-covered       does every state of the spec have a VR scenario in the feature?
//	vr-scenarios-tested     does every VR scenario have a VR test, and a baseline image?
//	vr-scenarios-of-states  is every VR scenario of a state the spec registers?
//	vr-tests-of-scenarios   is every VR test of a VR scenario the feature declares?
//
// An agent setting up a project with screens wrote specs with states, features and unit
// tests, and no visual regression at all — every gate green, and no state of any screen
// protected against a visual change. A state is what a screen looks like under a
// condition, and a capture is what keeps it so. The four questions tie each state to its
// capture both ways, so neither a state goes uncaptured nor a capture outlives its state.
//
// They run on the unit's main CODE file, because that is what a project tags as visual
// (`screen`, `component` on its layers), and read the spec, the feature and the tests of
// the unit from it. A VR scenario is one tagged with the project's visual regime and the
// state's code (`@BUTTN-S01` or `@BUTTN-VR-S01`); a VR test names `BUTTN-VR-S01` in its
// path or its text; a baseline is `<Unit>.BUTTN-VR-S01[-variant].<ext>` beside the unit.

// vrUnit is what the four questions read of a visual unit.
type vrUnit struct {
	base, unit, code string
	states           []string            // the states the spec registers: S01, S02…
	messages         []string            // the messages the spec catalogs: M01, M02…
	targets          []string            // what is captured: the states and the messages
	exempt           map[string]bool     // the states the spec exempts with `@no-vr: <reason>`
	unreasoned       []string            // the states marked `@no-vr` with no reason: still asked
	scenarios        map[string]bool     // the states the feature's VR scenarios are of
	stateless        []string            // VR scenarios of no state: the codes on their line
	tests            map[string][]string // the states VR tests name → the tests
}

// readVRUnit reads the unit a code node stands for, or says why there is none to confront.
func readVRUnit(n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (*vrUnit, string, bool) {
	if n.Kind != mapx.KindCode {
		return nil, "", false
	}
	// The unit's main file only: `Button.tsx` stands for `Button.spec.md`; `Button.styles.ts`
	// is a part of it, and confronting each part would repeat every gap.
	base := strings.TrimSuffix(n.ID, path.Ext(n.ID))
	spec, err := readFile(root, base+".spec.md")
	if err != nil {
		return nil, i18n.T("gate.vr_states.skip_no_spec"), false
	}
	m := specCodeRE().FindStringSubmatch(string(spec))
	if m == nil {
		return nil, i18n.T("gate.vr_states.skip_no_code"), false
	}
	u := &vrUnit{base: base, unit: path.Base(base), code: m[1], scenarios: map[string]bool{}, tests: map[string][]string{}}
	letter := stateLetter(cfg)
	u.states = registeredStates(string(spec), u.code, letter)
	// A MESSAGE is captured too: an error or a refusal shows on the screen as the message
	// the spec catalogs, and the state is the same — what changes is the text.
	u.messages = sectionCodes(string(spec), "section.title.messages", u.code)
	u.targets = append(append([]string(nil), u.states...), u.messages...)
	sort.Strings(u.targets)
	u.exempt, u.unreasoned = exemptStates(string(spec), u.code, `[A-Z]`)

	stateRE := regexp.MustCompile(`\b` + regexp.QuoteMeta(u.code) + `-(?:VR-)?([A-Z]\d{2})\b`)
	vrCodeRE := regexp.MustCompile(`@(` + regexp.QuoteMeta(u.code) + `-VR[\w-]*)`)
	if feature, err := readFile(root, base+".feature"); err == nil {
		tag := "@" + visualRegimeTag(cfg)
		for _, line := range strings.Split(string(feature), "\n") {
			if !hasTag(line, tag) {
				continue
			}
			found := false
			for _, sm := range stateRE.FindAllStringSubmatch(line, -1) {
				u.scenarios[sm[1]] = true
				found = true
			}
			if !found {
				name := strings.TrimSpace(line)
				if vm := vrCodeRE.FindStringSubmatch(line); vm != nil {
					name = vm[1]
				}
				u.stateless = append(u.stateless, name)
			}
		}
	}

	refRE := regexp.MustCompile(`\b` + regexp.QuoteMeta(u.code) + `-VR-([A-Z]\d{2})\b`)
	if g != nil {
		for _, t := range g.Nodes {
			if t.Kind != mapx.KindTest || isImage(t.ID) || !testOfUnit(t.ID, u) {
				continue
			}
			text := t.ID
			if b, err := readFile(root, t.ID); err == nil {
				text += "\n" + string(b)
			}
			seen := map[string]bool{}
			for _, sm := range refRE.FindAllStringSubmatch(text, -1) {
				if !seen[sm[1]] {
					seen[sm[1]] = true
					u.tests[sm[1]] = append(u.tests[sm[1]], t.ID)
				}
			}
		}
	}
	return u, "", true
}

// testOfUnit says whether a test may hold the unit's VR tests: its path names the unit's
// code, it sits beside the unit under the unit's name, or a folder of its path is named
// after the unit.
func testOfUnit(id string, u *vrUnit) bool {
	b := path.Base(id)
	return strings.Contains(id, u.code+"-") || strings.HasPrefix(b, u.unit+".") || strings.HasPrefix(b, u.unit+"_") ||
		strings.Contains("/"+id, "/"+u.unit+"/")
}

func hasTag(line, tag string) bool {
	for _, f := range strings.Fields(line) {
		if f == tag {
			return true
		}
	}
	return false
}

// checkVRStatesCovered: does every state of the spec have a VR scenario in the feature?
//
// EVERY state of a screen or a component is asked — only visual units are confronted at all
// —, and the exception is written where the state is: `@no-vr: <reason>` on its heading or
// its row, for a state with no visual value of its own (a transient loading, a state that
// looks like another). An exemption with no reason does not exempt: it cannot be told from
// a capture nobody wanted to make.
func checkVRStatesCovered(_ string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	u, why, ok := readVRUnit(n, root, g, cfg)
	if !ok {
		return Skip, why
	}
	if len(u.targets) == 0 {
		return Skip, i18n.T("gate.vr_states.skip_no_states")
	}
	var missing []string
	for _, s := range u.targets {
		if !u.scenarios[s] && !u.exempt[s] {
			missing = append(missing, u.code+"-"+s)
		}
	}
	var gaps []string
	if len(missing) > 0 {
		gaps = append(gaps, i18n.T("gate.vr_states.no_scenario", len(missing), strings.Join(missing, ", "), visualRegimeTag(cfg)))
	}
	if len(u.unreasoned) > 0 {
		gaps = append(gaps, i18n.T("gate.vr_states.no_reason", len(u.unreasoned), strings.Join(u.unreasoned, ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// checkVRScenariosTested: does every VR scenario have a VR test, and a baseline image?
func checkVRScenariosTested(_ string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	u, why, ok := readVRUnit(n, root, g, cfg)
	if !ok {
		return Skip, why
	}
	if len(u.scenarios) == 0 {
		return Skip, i18n.T("gate.vr_states.skip_no_vr_scenarios")
	}
	var noTest, noImage []string
	for _, s := range sortedKeys(u.scenarios) {
		if len(u.tests[s]) == 0 {
			noTest = append(noTest, u.code+"-VR-"+s)
		}
		if found, _ := doublestar.Glob(os.DirFS(root), u.base+"."+u.code+"-VR-"+s+"*."+baselineExts); len(found) == 0 {
			noImage = append(noImage, u.code+"-VR-"+s)
		}
	}
	var gaps []string
	if len(noTest) > 0 {
		gaps = append(gaps, i18n.T("gate.vr_states.no_test", len(noTest), strings.Join(noTest, ", ")))
	}
	if len(noImage) > 0 {
		gaps = append(gaps, i18n.T("gate.vr_states.no_baseline", len(noImage), strings.Join(noImage, ", "), u.unit+"."+u.code+"-VR-<state>.png"))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// checkVRScenariosOfStates: is every VR scenario of a state the spec registers?
func checkVRScenariosOfStates(_ string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	u, why, ok := readVRUnit(n, root, g, cfg)
	if !ok {
		return Skip, why
	}
	if len(u.scenarios) == 0 && len(u.stateless) == 0 {
		return Skip, i18n.T("gate.vr_states.skip_no_vr_scenarios")
	}
	known := map[string]bool{}
	for _, s := range u.targets {
		known[s] = true
	}
	var orphans, exempted []string
	for _, s := range sortedKeys(u.scenarios) {
		switch {
		case !known[s]:
			orphans = append(orphans, u.code+"-"+s)
		case u.exempt[s]:
			exempted = append(exempted, u.code+"-"+s)
		}
	}
	var gaps []string
	if len(orphans) > 0 {
		gaps = append(gaps, i18n.T("gate.vr_states.scenario_no_state", len(orphans), strings.Join(orphans, ", ")))
	}
	if len(exempted) > 0 {
		gaps = append(gaps, i18n.T("gate.vr_states.scenario_exempt", len(exempted), strings.Join(exempted, ", ")))
	}
	if len(u.stateless) > 0 {
		gaps = append(gaps, i18n.T("gate.vr_states.scenario_stateless", len(u.stateless), strings.Join(u.stateless, ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// checkVRTestsOfScenarios: is every VR test of a VR scenario the feature declares?
func checkVRTestsOfScenarios(_ string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	u, why, ok := readVRUnit(n, root, g, cfg)
	if !ok {
		return Skip, why
	}
	if len(u.tests) == 0 {
		return Skip, i18n.T("gate.vr_states.skip_no_vr_tests")
	}
	var orphans []string
	for _, s := range sortedKeysOf(u.tests) {
		if !u.scenarios[s] {
			orphans = append(orphans, u.code+"-VR-"+s+" ("+strings.Join(u.tests[s], ", ")+")")
		}
	}
	if len(orphans) == 0 {
		return Pass, ""
	}
	return Fail, i18n.T("gate.vr_states.test_no_scenario", len(orphans), strings.Join(orphans, "; "))
}

// stateLetter is the letter of the project's State rule type: the one whose term or
// sections name a state, else the canonical `S`.
func stateLetter(cfg *config.Config) string {
	if cfg != nil {
		for _, rt := range cfg.RuleTypes {
			names := append([]string{rt.Term}, rt.Sections...)
			for _, t := range names {
				t = strings.ToLower(strings.TrimSpace(t))
				if strings.HasPrefix(t, "state") || strings.HasPrefix(t, "estado") {
					if l := strings.ToUpper(strings.TrimSpace(rt.Letter)); len(l) == 1 {
						return l
					}
				}
			}
		}
	}
	return "S"
}

// registeredStates are the states the spec REGISTERS: the codes in its States section when
// it has one — a state code cited elsewhere (a transition to a state that does not exist)
// registers nothing —, else every state code of the spec.
func registeredStates(content, code, letter string) []string {
	if text, ok := sectionText(content, "section.title.states"); ok {
		return specStates(text, code, letter)
	}
	return specStates(content, code, letter)
}

// specStates are the state codes the spec registers — `S01`, `S02` — in order, once each.
func specStates(content, code, letter string) []string {
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(code) + `-(` + regexp.QuoteMeta(letter) + `\d{2})\b`)
	seen := map[string]bool{}
	var out []string
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	sort.Strings(out)
	return out
}

// noVRRE is the exemption `@no-vr`, and the reason written after its colon.
var noVRRE = regexp.MustCompile(`@no-vr\b(?::\s*([^|]*))?`)

// exemptStates are the states whose declaration carries `@no-vr: <reason>` on the same line
// — the state's heading or its row in the states table —, and the ones marked with no
// reason, which stay asked.
func exemptStates(content, code, letterClass string) (map[string]bool, []string) {
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(code) + `-(` + letterClass + `\d{2})\b`)
	exempt := map[string]bool{}
	var unreasoned []string
	for _, line := range strings.Split(content, "\n") {
		m := noVRRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		reasoned := strings.TrimSpace(strings.Trim(strings.TrimSpace(m[1]), "-—>")) != ""
		for _, sm := range re.FindAllStringSubmatch(line, -1) {
			if reasoned {
				exempt[sm[1]] = true
			} else {
				unreasoned = append(unreasoned, code+"-"+sm[1])
			}
		}
	}
	sort.Strings(unreasoned)
	return exempt, unreasoned
}

// sectionCodes are the codes of the unit — their suffix, `M01` — in the first column of the
// table of the section the catalog key names, in any language.
func sectionCodes(content, key, code string) []string {
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(code) + `-([A-Z]\d{2})\b`)
	seen := map[string]bool{}
	var out []string
	for _, r := range sectionRows(content, key) {
		first, _ := col(r, ruleCols...)
		if m := re.FindStringSubmatch(first); m != nil && !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	sort.Strings(out)
	return out
}

func isImage(id string) bool {
	switch strings.ToLower(path.Ext(id)) {
	case ".png", ".jpg", ".jpeg", ".webp", ".gif", ".svg":
		return true
	}
	return false
}

func sortedKeysOf(m map[string][]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
