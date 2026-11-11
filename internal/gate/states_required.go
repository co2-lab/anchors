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
	code := specCode(content)
	stateLine := regexp.MustCompile(`\b` + regexp.QuoteMeta(code) + `-` + regexp.QuoteMeta(stateLetter(cfg)) + `\d{2}\b`)
	var states []string
	for _, l := range strings.Split(content, "\n") {
		if code != "" && stateLine.MatchString(l) {
			states = append(states, strings.ToLower(l))
		}
	}
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

// noStateRE is a spec's waiver of a required state: its name, then a dash and the reason.
var noStateRE = regexp.MustCompile(`@no-state[^\S\n]*:[^\S\n]*([^\s—–-]+)[^\S\n]*[—–-][^\S\n]*\S`)
