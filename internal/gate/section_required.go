// @anchors
//   code: SCRQG
//   ref: SCRQR

package gate

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// section-required: a spec section the unit's code calls for exists in the spec
// (DESIGN-section-required.md).
//
// The gates that read a section confront it once it exists — `testid-consistent` the Test
// IDs, the data gates the data contract. When the code calls for one and the spec has
// none, they have nothing to confront and say nothing: measured at MIF, 67 component
// specs exposed handles with no Test IDs section, and the spec end of every handle went
// unchecked. This gate makes the absence a finding.
//
// What calls for a section is the project's to declare, per layer of the unit's code
// (`layers.<l>.requires_sections`), with a trigger read from what the core already
// measures: a handle exposed, a dependency on a layer, a flag, a pattern of the
// project's. No framework is known here.
func checkSectionRequired(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.section_required.skip_not_spec")
	}
	if g == nil {
		return pendingNoMap()
	}
	var code []mapx.Node
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		if t := g.Node(e.To); t != nil {
			code = append(code, *t)
		}
	}
	if len(code) == 0 {
		// A spec with no code: the absence is unit-complete's to charge.
		return Skip, i18n.T("gate.section_required.skip_no_code")
	}
	have := map[string]bool{}
	for _, t := range headingTitles(content) {
		have[normalSection(t)] = true
	}

	var missing, pending []string
	asked := false
	seen := map[string]bool{}
	for _, c := range code {
		layer, ok := cfg.Layers[c.Layer]
		if !ok {
			continue
		}
		for _, req := range layer.RequiresSections {
			asked = true
			key := normalSection(req.Section)
			if have[key] || seen[key] {
				continue
			}
			why, measured := sectionTrigger(req.When, c, root, g, cfg)
			if !measured {
				pending = append(pending, why)
				continue
			}
			if why != "" {
				seen[key] = true
				missing = append(missing, i18n.T("gate.section_required.line", req.Section, filepath.Base(c.ID), why))
			}
		}
	}
	if !asked {
		return Skip, i18n.T("gate.section_required.skip_nothing_declared")
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Fail, i18n.T("gate.section_required.fail_missing", filepath.Base(n.ID), len(missing), strings.Join(missing, "\n"))
	}
	if len(pending) > 0 {
		sort.Strings(pending)
		return Pending, strings.Join(pending, "\n")
	}
	return Pass, ""
}

// normalSection is a section title as compared: its case and surrounding spaces left out.
func normalSection(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}

// sectionTrigger says why the code calls for the section — what it found, "" when it calls
// for nothing — and whether the trigger could be measured at all.
func sectionTrigger(w config.SectionTrigger, c mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (string, bool) {
	switch {
	case len(w.DependsOnLayer) > 0:
		layers := map[string]bool{}
		for _, l := range w.DependsOnLayer {
			layers[l] = true
		}
		// A dependency on the layer's values counts; one on its types only — `@dep[type]`,
		// the shape of the data and not the data — counts when the entry asks for it.
		kinds := map[string]bool{"value": true}
		if len(w.Kinds) > 0 {
			kinds = map[string]bool{}
			for _, k := range w.Kinds {
				kinds[k] = true
			}
		}
		for _, e := range g.Neighbors(c.ID).Out {
			if e.Type != mapx.EdgeDependsOn {
				continue
			}
			kind := "value"
			if e.Kind == "type" {
				kind = "type"
			}
			if !kinds[kind] {
				continue
			}
			if t := g.Node(e.To); t != nil && layers[t.Layer] {
				if kind == "type" {
					return i18n.T("gate.section_required.why_depends_type", t.ID, t.Layer), true
				}
				return i18n.T("gate.section_required.why_depends", t.ID, t.Layer), true
			}
		}
		return "", true
	}
	b, err := readFile(root, c.ID)
	if err != nil {
		return "", true // @resilient: a code file that cannot be read calls for nothing this gate can name; unit-complete and the map notice a missing file
	}
	src := string(b)
	switch {
	case w.Exposes != "":
		exposed := exposedPatterns(cfg)
		if len(exposed) == 0 {
			return i18n.T("gate.testid.no_patterns", "dialect.test_handles.exposed"), false
		}
		if hs := handlesIn(src, c.ID, exposed); len(hs) > 0 {
			sort.Strings(hs)
			return i18n.T("gate.section_required.why_exposes", hs[0]), true
		}
		if u := unmeasurableHandles(src, c.ID, exposed); len(u) > 0 {
			return i18n.T("gate.testid.unmeasurable", strings.Join(u, ", ")), false
		}
	case w.Flag != "":
		cs := cfg.CommentsFor(c.ID)
		if cs == nil {
			return i18n.T("gate.comments_undeclared", c.ID), false
		}
		flag := regexp.MustCompile(regexp.QuoteMeta(w.Flag) + `\b`)
		for _, sp := range cs.Spans(src) {
			if flag.MatchString(src[sp[0]:sp[1]]) {
				return i18n.T("gate.section_required.why_flag", w.Flag), true
			}
		}
	case w.CodeMatches != "":
		re := regexp.MustCompile(w.CodeMatches) // validated on load
		if m := re.FindString(cfg.CommentsFor(c.ID).Blank(src)); m != "" {
			return i18n.T("gate.section_required.why_matches", strings.TrimSpace(m)), true
		}
	}
	return "", true
}
