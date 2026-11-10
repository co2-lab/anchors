// @anchors
//   code: VRPRC
//   ref: VRPRS

package gate

import (
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// variants-proven: a rule that names its variants is proven once per variant.
//
// Measured in MIF (Sentry MIF-1S): a rule of the alerts screen spoke of two kinds of item,
// audit and notification; one scenario carried its tag, proved the audit, and every gate
// passed — mutation at 100%, since the code shared by both kinds has no branch per kind. The
// notification path reached the store review broken. A tag counts scenarios; it cannot tell
// that the one scenario proves half of what the rule says.
//
// The rule declares its variants on its line, and each variant asks a scenario of the rule
// that names it — by a `@variant:<value>` tag, or by a row of its Examples:
//
//	| `UAUDU-B02` | a tap marks the item seen @variants: audit, notification |
//	| `UAUDU-E02` | a failed write propagates @variants: B02 |
//
// A rule may take its variants from another rule of the unit (`@variants: B02`): an error
// rule shared by every kind then asks the new kind as soon as the rule it follows gains it.
// A rule short of its own variants fails; a rule short of the variants it follows diverges —
// "proven only on the old ones" —, the warning MIF asked for the siblings.
func checkVariantsProven(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.variants.skip_not_spec")
	}
	unit := specCode(content)
	own, follows := declaredVariants(content, unit)
	if len(own) == 0 && len(follows) == 0 {
		return Skip, i18n.T("gate.variants.skip_none")
	}
	// A rule follows another's variants: the other must be a rule of this unit that names
	// its own, or the declaration stands for nothing.
	var bad []string
	for rule, target := range follows {
		if len(own[target]) == 0 {
			bad = append(bad, i18n.T("gate.variants.follows_nothing", rule, target))
		}
	}
	if len(bad) > 0 {
		sort.Strings(bad)
		return Fail, strings.Join(bad, "; ")
	}
	if g == nil {
		return pendingNoMap()
	}
	var features []string
	for _, e := range g.Edges {
		if e.From == n.ID && e.Type == "covered-by" {
			features = append(features, e.To)
		}
	}
	if len(features) == 0 {
		return Skip, i18n.T("gate.variants.skip_no_feature")
	}
	proven := map[string]map[string]bool{} // rule → variant (lower case) → proven
	mark := func(rule, v string) {
		if proven[rule] == nil {
			proven[rule] = map[string]bool{}
		}
		proven[rule][strings.ToLower(v)] = true
	}
	for _, f := range features {
		b, err := readFile(root, f)
		if err != nil {
			continue
		}
		text := string(b)
		lines := strings.Split(text, "\n")
		for _, sc := range parseFeatureScenarios(text) {
			if sc.Line < 1 || sc.Line > len(lines) {
				continue
			}
			for _, m := range variantTagRE.FindAllStringSubmatch(lines[sc.Line-1], -1) {
				for _, c := range append([]string{sc.Code}, sc.Codes...) {
					mark(RootCode(c), m[1])
				}
			}
		}
		// An outline proves each variant one of its Examples rows names.
		for _, o := range exampleTables(text) {
			for _, row := range o.rows {
				for _, v := range row.values {
					mark(RootCode(o.code), v)
				}
			}
		}
	}

	missing := func(rule string, variants []string) []string {
		var out []string
		for _, v := range variants {
			if !proven[rule][strings.ToLower(v)] {
				out = append(out, v)
			}
		}
		return out
	}
	var failed, diverged []string
	for _, rule := range sortedKeysOf(own) {
		if m := missing(rule, own[rule]); len(m) > 0 {
			failed = append(failed, rule+": "+strings.Join(m, ", "))
		}
	}
	followers := make([]string, 0, len(follows))
	for rule := range follows {
		followers = append(followers, rule)
	}
	sort.Strings(followers)
	for _, rule := range followers {
		vs := own[follows[rule]]
		if m := missing(rule, vs); len(m) > 0 {
			diverged = append(diverged, i18n.T("gate.variants.only_on", rule, strings.Join(m, ", "), follows[rule]))
		}
	}
	switch {
	case len(failed) > 0:
		return Fail, i18n.T("gate.variants.missing", strings.Join(append(failed, diverged...), "; "))
	case len(diverged) > 0:
		return Diverge, strings.Join(diverged, "; ")
	}
	return Pass, ""
}

// variantsMarkRE reads `@variants: a, b` on a rule's line, up to the end of the table cell.
var variantsMarkRE = regexp.MustCompile(`@variants[^\S\n]*:[^\S\n]*([^|\n]+)`)

// variantTagRE reads a scenario's `@variant:<value>` tag.
var variantTagRE = regexp.MustCompile(`@variant:([^\s@]+)`)

// declaredVariants reads the variants each rule of the unit declares: its own list, or the
// rule of the same unit it follows (`@variants: B02` or `@variants: UNIT-B02`).
func declaredVariants(content, unit string) (own map[string][]string, follows map[string]string) {
	own, follows = map[string][]string{}, map[string]string{}
	letter := regexp.MustCompile(`^(?:` + regexp.QuoteMeta(unit) + `-)?([A-Z]\d{2})$`)
	for _, line := range strings.Split(content, "\n") {
		m := variantsMarkRE.FindStringSubmatch(line)
		if m == nil || retiredLine(line) {
			continue
		}
		rule := ruleOnLine(line, unit)
		if rule == "" || unit == "" || !strings.HasPrefix(rule, unit+"-") {
			continue
		}
		var vs []string
		for _, v := range strings.Split(m[1], ",") {
			if v = strings.Trim(strings.TrimSpace(v), "`\"'"); v != "" {
				vs = append(vs, v)
			}
		}
		if len(vs) == 1 {
			if lm := letter.FindStringSubmatch(vs[0]); lm != nil {
				follows[rule] = unit + "-" + lm[1]
				continue
			}
		}
		if len(vs) > 0 {
			own[rule] = vs
		}
	}
	return own, follows
}
