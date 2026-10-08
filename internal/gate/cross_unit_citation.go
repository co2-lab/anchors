// @anchors
//   code: CUCGC
//   ref: CRUCT

package gate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cross-unit-citation: does a rule of another unit the spec cites live in the product?
//
// A spec may REFERENCE another unit — point at its identity, as a screen's Navigation names
// the screen it leads to — and may CITE another unit's rule by its code. What it cites,
// though, is content two units share, and a shared rule lives in the product doctrine, which
// comes first: the cited rule realizes a rule of the doctrine (`@realizes` on its line), so
// neither spec watches the other, both follow the product above them
// (DESIGN-dependencies-out-of-the-spec.md, DOOSD-D02, D04, D05).
//
// Left out: the navigation sections — a reference —, the lines realizing the product, aliases,
// revisions and retired rules. A prose that names another screen without its code is a
// judgment, asked in the review guide, not here.
func checkCrossUnitCitation(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.cross_unit.skip_not_spec")
	}
	if g == nil {
		return pendingNoMap()
	}
	own := headerCode(content)
	units := map[string]string{} // another spec's unit code → its spec
	for _, x := range g.Nodes {
		if x.Kind == mapx.KindSpec && x.Code != "" && x.Code != own && x.ID != n.ID {
			units[x.Code] = x.ID
		}
	}
	references := sectionTitles([]string{"navigation"}, cfg, n.Layer)
	type level struct {
		depth int
		ref   bool
	}
	var stack []level
	inRef := func() bool {
		for _, l := range stack {
			if l.ref {
				return true
			}
		}
		return false
	}
	re := refCodeRE()
	cited := map[string][]int{}
	for i, line := range strings.Split(content, "\n") {
		if m := headingRE.FindStringSubmatch(line); m != nil {
			depth := len(m[1])
			for len(stack) > 0 && stack[len(stack)-1].depth >= depth {
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, level{depth, references[strings.ToLower(strings.TrimSpace(m[2]))]})
			continue
		}
		if strings.HasPrefix(line, "# ") {
			stack = nil
		}
		if inRef() || retiredLine(line) || strings.Contains(line, "@realizes") || strings.Contains(line, "REF[") ||
			revisesRE().MatchString(line) || checkedRE().MatchString(line) || revisionRE().MatchString(line) {
			continue
		}
		for _, m := range re.FindAllStringSubmatch(line, -1) {
			if _, ok := units[m[1]]; ok {
				code := m[0]
				cited[code] = append(cited[code], i+1)
			}
		}
	}
	// A cited rule that realizes the product is shared the right way.
	texts := map[string]string{}
	for code := range cited {
		spec := units[strings.SplitN(code, "-", 2)[0]]
		if _, ok := texts[spec]; !ok {
			b, _ := readFile(root, spec)
			texts[spec] = string(b)
		}
		if realizesProduct(texts[spec], code) {
			delete(cited, code)
		}
	}
	if len(cited) == 0 {
		return Pass, ""
	}
	codes := make([]string, 0, len(cited))
	for c := range cited {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	var parts []string
	for _, c := range codes {
		parts = append(parts, fmt.Sprintf("`%s` (%s)", c, joinInts(cited[c])))
	}
	return Fail, i18n.T("gate.cross_unit.cites", len(codes), strings.Join(parts, ", "))
}

// realizesProduct says whether a rule, in its spec's text, is defined on a line that realizes a
// rule of the product doctrine.
func realizesProduct(spec, code string) bool {
	re := defineRuleCaptureRE()
	for _, line := range strings.Split(spec, "\n") {
		if m := re.FindStringSubmatch(line); m != nil && m[1] == code && strings.Contains(line, "@realizes") {
			return true
		}
	}
	return false
}
