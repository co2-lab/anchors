// @anchors
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

// vr-states-covered: every state a visual unit's spec registers is proven by visual
// regression — a VR scenario in the feature, a baseline image per state beside the unit,
// and a test that captures them.
//
// `vr-baseline` asks, from the feature, whether a VR scenario that EXISTS has its image.
// Nothing asked whether it exists: an agent configuring a project with screens wrote
// specs with states, features and unit tests, and no visual regression at all — every
// gate green, and no state of any screen protected against a visual change. A state is
// what a screen looks like under a condition, and a capture is what keeps it so.
//
// It confronts the unit's CODE — the screen or the component —, because that is what the
// project tags as visual (`screen`, `component` on its layers; the gate's `tags:` scope it),
// and reads the states from the spec beside it. The naming is the one `vr-baseline` reads:
// `<Unit>.<CODE>-VR-<state>.<ext>` beside the unit.
func checkVRStatesCovered(_ string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode {
		return Skip, ""
	}
	// The unit's main file only: `Button.tsx` stands for `Button.spec.md`; `Button.styles.ts`
	// is a part of it, and confronting each part would repeat every gap.
	base := strings.TrimSuffix(n.ID, path.Ext(n.ID))
	spec, err := readFile(root, base+".spec.md")
	if err != nil {
		return Skip, i18n.T("gate.vr_states.skip_no_spec")
	}
	content := string(spec)
	m := specCodeRE().FindStringSubmatch(content)
	if m == nil {
		return Skip, i18n.T("gate.vr_states.skip_no_code")
	}
	code := m[1]
	states := specStates(content, code, stateLetter(cfg))
	if len(states) == 0 {
		return Skip, i18n.T("gate.vr_states.skip_no_states")
	}
	vr := code + "-VR"

	var gaps []string
	feature, err := readFile(root, base+".feature")
	if err != nil || !containsString(vrScenarios(string(feature), cfg), vr) {
		gaps = append(gaps, i18n.T("gate.vr_states.no_scenario", vr, visualRegimeTag(cfg)))
	}
	var noImage []string
	for _, s := range states {
		found, _ := doublestar.Glob(os.DirFS(root), base+"."+vr+"-"+s+"*."+baselineExts)
		if len(found) == 0 {
			noImage = append(noImage, code+"-"+s)
		}
	}
	if len(noImage) > 0 {
		gaps = append(gaps, i18n.T("gate.vr_states.no_baseline", len(noImage), strings.Join(noImage, ", "),
			path.Base(base)+"."+vr+"-<state>.png"))
	}
	if !hasVRCapture(root, g, base, vr) {
		gaps = append(gaps, i18n.T("gate.vr_states.no_capture", vr))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
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

// hasVRCapture says whether a test of the project captures the unit's VR: a test node
// whose path names `<CODE>-VR` (a capture flow), or one beside the unit whose content
// does (a screenshot test). Baseline images are evidence, not the capture.
func hasVRCapture(root string, g *mapx.Graph, base, vr string) bool {
	if g == nil {
		return false
	}
	unit := path.Base(base)
	for _, t := range g.Nodes {
		if t.Kind != mapx.KindTest || isImage(t.ID) {
			continue
		}
		if strings.Contains(t.ID, vr) {
			return true
		}
		if strings.HasPrefix(path.Base(t.ID), unit+".") || strings.HasPrefix(path.Base(t.ID), unit+"_") {
			if b, err := readFile(root, t.ID); err == nil && strings.Contains(string(b), vr) {
				return true
			}
		}
	}
	return false
}

func isImage(id string) bool {
	switch strings.ToLower(path.Ext(id)) {
	case ".png", ".jpg", ".jpeg", ".webp", ".gif", ".svg":
		return true
	}
	return false
}

func containsString(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
