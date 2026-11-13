// @anchors
//   code: STRQC
//   ref: STRQS

package gate

import (
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// states-required: every spec of a layer registers the states the project requires of it.
//
// Measured in MIF: nine released defects were a screen state nobody wrote — an error shown as
// an empty list, a not-found that spun forever, an offline screen with no way back. Each spec
// registered the states its author thought of, and every gate confronted only those. The
// project names the states each layer owes (`required_states: {screen: [loading, error,
// empty]}`), and each spec of that layer registers a state by each name — a state line, its
// heading or its row, carrying the state's code and the name —, or says why not on a line of
// its own: `@no-state: offline — the screen reads only local data`.
//
// A state is registered where it is DEFINED, by its name: the heading `### CODE-S01: Error`, the
// row whose first cell is the code and whose next cell names it, the list item that opens with
// the code. A line that only cites a state — measured in MIF: a message row "Empty state —
// ARSCA-S02", a prose line that mentions an error on its way — registered whatever word it
// carried, and a screen with no error state of its own passed.
func checkStatesRequired(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.states_required.skip_not_spec")
	}
	if cfg == nil || len(cfg.RequiredStates) == 0 {
		return Skip, i18n.T("gate.states_required.skip_none_declared")
	}
	layer := layerOf(n, content)
	required := cfg.RequiredStates[layer]
	if len(required) == 0 {
		return Skip, i18n.T("gate.states_required.skip_layer", layer)
	}
	states := stateNames(content, specCode(content), stateLetter(cfg))
	waived := map[string]bool{}
	for _, m := range noStateRE.FindAllStringSubmatch(content, -1) {
		waived[strings.ToLower(m[1])] = true
	}
	var missing []string
	for _, name := range required {
		if waived[strings.ToLower(name)] {
			continue
		}
		word := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(name) + `\b`)
		found := false
		for _, s := range states {
			if word.MatchString(s) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, name)
		}
	}
	if len(missing) == 0 {
		return Pass, ""
	}
	return Fail, i18n.T("gate.states_required.missing", layer, strings.Join(missing, ", "))
}

// stateNames are the names of the states the spec defines, lowercased: a heading's text, the
// cell after the code in a row, a list item's text up to the first dash or period that ends
// its name.
func stateNames(content, code, letter string) []string {
	if code == "" {
		return nil
	}
	def := regexp.MustCompile("^\\s*(#{2,6}\\s+|[-*]\\s+\\**|\\|\\s*)`?\\*{0,2}" +
		regexp.QuoteMeta(code) + "-" + regexp.QuoteMeta(letter) + `\d{2}\b(.*)$`)
	var out []string
	for _, l := range strings.Split(content, "\n") {
		m := def.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		name := m[2]
		switch strings.TrimSpace(m[1])[0] {
		case '|':
			cells := strings.Split(name, "|")
			if len(cells) > 1 {
				cells = cells[:2]
			}
			name = strings.Join(cells, " ")
		case '-', '*':
			if i := listNameEnd.FindStringIndex(name); i != nil {
				name = name[:i[0]]
			}
		}
		out = append(out, strings.ToLower(name))
	}
	return out
}

// listNameEnd is where a list item's state name ends: a spaced dash or a period.
var listNameEnd = regexp.MustCompile(`\s[—–-]\s|\.(\s|$)`)

// noStateRE is a spec's waiver of a required state: its name, then a dash and the reason.
var noStateRE = regexp.MustCompile(`@no-state[^\S\n]*:[^\S\n]*([^\s—–-]+)[^\S\n]*[—–-][^\S\n]*\S`)
