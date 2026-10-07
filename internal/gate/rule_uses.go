// @anchors
//   code: RUGRL
//   ref: RLUSG

package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// THE RULE AND WHAT IT USES.
//
// A spec has the two ends — the rules (behaviours, states, errors…) and the data (the data
// contract, the domain, the props) — and nothing between them: what each rule READS lived,
// when it lived anywhere, in prose. So a field changed in the contract pointed at no rule,
// and a rule could check a datum the contract never offered.
//
// Three sections tie them, and they are read here:
//
//	Validations               `| V01 | field | condition | behaviour |`
//	Presentation validations  `| P01 | prop/state | condition | appearance |`
//	Rule uses                 `| B01 | field, CODE-S01, DEP1 |`
//
// Each row starts with the rule's code and goes on with what it uses, in that order. The
// row is read by POSITION, not by the columns' names: those follow the project's language
// (`Campo`, `Field`), the order does not. The sections are found by their title in any
// supported language, or by the one the project declares in `section_titles`.

// ruleUseSections are the catalog keys of the three sections (`section.title.<key>`).
var ruleUseSections = []string{"validations", "presentation_validations", "rule_uses"}

// ruleUse is one row: the rule, and what it says it uses.
type ruleUse struct {
	Rule string
	Uses []string
}

var (
	headingRE     = regexp.MustCompile(`^(#{2,4})\s+(.+?)\s*$`)
	backtickedRE  = regexp.MustCompile("`([^`]+)`")
	tableDividerR = regexp.MustCompile(`^\|[\s:|-]+\|?$`)
	depRefRE      = regexp.MustCompile(`^DEP\d+$`)
	depRowRE      = regexp.MustCompile(`(?m)^\|\s*(DEP\d+)\s*\|`)
	noUsesRE      = regexp.MustCompile(`@no-uses[^\S\n]*:[^\S\n]*\S`)
)

// ruleUseTitles are the titles the three sections go by: every supported language, and the
// project's own names for them.
func ruleUseTitles(cfg *config.Config, layer string) map[string]bool {
	return sectionTitles(ruleUseSections, cfg, layer)
}

// sectionTitles are the titles the catalog sections `keys` go by: every supported
// language, and the project's own names for them (`section_titles`).
func sectionTitles(keys []string, cfg *config.Config, layer string) map[string]bool {
	out := map[string]bool{}
	for _, k := range keys {
		for _, t := range i18n.AllTranslations("section.title." + k) {
			out[strings.ToLower(t)] = true
		}
		if cfg != nil {
			if t := cfg.SectionTitle(strings.ReplaceAll(k, "_", "-"), "", layer); t != "" {
				out[strings.ToLower(t)] = true
			}
		}
	}
	return out
}

// splitRuleUseSections separates a spec into the text of the three sections and the rest.
func splitRuleUseSections(content string, cfg *config.Config, layer string) (inside, outside []string) {
	return splitSections(content, ruleUseSections, cfg, layer)
}

// splitSections separates a spec into the text of the catalog sections `keys` and the rest.
func splitSections(content string, keys []string, cfg *config.Config, layer string) (inside, outside []string) {
	titles := sectionTitles(keys, cfg, layer)
	in := false
	for _, l := range strings.Split(content, "\n") {
		if m := headingRE.FindStringSubmatch(l); m != nil {
			in = titles[strings.ToLower(m[2])]
		} else if strings.HasPrefix(l, "# ") {
			in = false
		}
		if in {
			inside = append(inside, l)
		} else {
			outside = append(outside, l)
		}
	}
	return inside, outside
}

// ruleUsesOf reads the rows of the three sections: each row that starts with a rule's code,
// with what its second cell says the rule uses. A row whose uses are still a TODO says
// nothing yet, and is left out.
func ruleUsesOf(content string, cfg *config.Config, layer string) []ruleUse {
	inside, _ := splitRuleUseSections(content, cfg, layer)
	code := ruleCodeRE()
	var out []ruleUse
	for _, l := range inside {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "|") || tableDividerR.MatchString(t) {
			continue
		}
		if retiredLine(t) {
			continue // a retired rule uses nothing
		}
		cells := cellsOf(t)
		if len(cells) < 2 {
			continue
		}
		rule := code.FindString(cells[0])
		used := strings.TrimSpace(cells[1])
		if rule == "" || used == "" || strings.HasPrefix(strings.ToUpper(used), "TODO") {
			continue
		}
		out = append(out, ruleUse{Rule: rule, Uses: usedItems(used)})
	}
	return out
}

// cellsOf splits a markdown table row into its cells, trimmed of the outer pipes. A pipe
// escaped as `\|` is text inside a cell — how a table writes a union of values
// (`'sm' \| 'lg'`) — and does not split it.
func cellsOf(row string) []string {
	row = strings.TrimSpace(row)
	row = strings.TrimSuffix(strings.TrimPrefix(row, "|"), "|")
	var out []string
	var cur strings.Builder
	for i := 0; i < len(row); i++ {
		switch {
		case row[i] == '\\' && i+1 < len(row) && row[i+1] == '|':
			cur.WriteString("\\|")
			i++
		case row[i] == '|':
			out = append(out, cur.String())
			cur.Reset()
		default:
			cur.WriteByte(row[i])
		}
	}
	return append(out, cur.String())
}

// usedItems splits a uses cell: the backticked items when there are any, else the
// comma-separated ones.
func usedItems(cell string) []string {
	var out []string
	if ms := backtickedRE.FindAllStringSubmatch(cell, -1); len(ms) > 0 {
		for _, m := range ms {
			out = append(out, strings.TrimSpace(m[1]))
		}
		return out
	}
	for _, p := range strings.Split(cell, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// gateEntry finds the project's entry for a check, to read its own settings.
func gateEntry(cfg *config.Config, check string) config.Gate {
	if cfg != nil {
		for _, g := range cfg.Gates {
			if g.Check == check {
				return g
			}
		}
	}
	return config.Gate{}
}

// defaultUsesExempt are the letters whose items are not rules that read data: an open
// question, a plan's phase, a feature-flag scenario.
var defaultUsesExempt = map[string]bool{"Q": true, config.PhaseLetter: true, "G": true}

// rule-uses-declared: every rule of the letters the gate asks about says what it uses — a
// row of its own in Validations or Presentation validations, or a row in Rule uses. The
// gate entry's `letters` lists the letters asked about; without it, every rule letter but
// an open question's, a phase's and a flag scenario's. A rule whose line carries
// `@no-uses: <why>` is waived, with the reason where the rule is.
func checkRuleUsesDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.rule_uses.skip_not_spec")
	}
	letters := map[string]bool{}
	for _, l := range gateEntry(cfg, "rule-uses-declared").Letters {
		letters[strings.ToUpper(strings.TrimSpace(l))] = true
	}
	waived := map[string]bool{}
	code := ruleCodeRE()
	for _, l := range strings.Split(content, "\n") {
		if noUsesRE.MatchString(l) {
			for _, c := range code.FindAllString(l, -1) {
				waived[c] = true
			}
		}
	}
	declared := map[string]bool{}
	for _, u := range ruleUsesOf(content, cfg, n.Layer) {
		declared[u.Rule] = true
	}
	var asked, missing []string
	for _, r := range definedRequirements(content) {
		letter := r[strings.LastIndex(r, "-")+1 : len(r)-2]
		if len(letters) > 0 && !letters[letter] || len(letters) == 0 && defaultUsesExempt[letter] {
			continue
		}
		asked = append(asked, r)
		if !declared[r] && !waived[r] {
			missing = append(missing, r)
		}
	}
	if len(asked) == 0 {
		return Skip, i18n.T("gate.rule_uses.skip_no_rules")
	}
	if len(missing) > 0 {
		return Fail, i18n.T("gate.rule_uses_declared.fail_missing", len(missing), listCodes(missing))
	}
	return Pass, ""
}

// rule-uses-resolve: what a rule says it uses exists. A field is a name the spec declares —
// the first cell of a row of any of its other tables (the data contract, the domain, the
// props, the signature, the state shape) or a backticked name in a heading (a data state);
// `summary.balance` resolves by `summary`. A `DEPn` is a row of the dependencies table.
// A code is left to `code-reference-valid`, which already knows every code of the project.
func checkRuleUsesResolve(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.rule_uses.skip_not_spec")
	}
	rows := ruleUsesOf(content, cfg, n.Layer)
	if len(rows) == 0 {
		return Skip, i18n.T("gate.rule_uses_resolve.skip_no_rows")
	}
	_, outside := splitRuleUseSections(content, cfg, n.Layer)
	names := declaredNames(outside)
	deps := map[string]bool{}
	for _, m := range depRowRE.FindAllStringSubmatch(strings.Join(outside, "\n"), -1) {
		deps[m[1]] = true
	}
	code := ruleCodeRE()
	var broken []string
	for _, u := range rows {
		for _, it := range u.Uses {
			switch {
			case code.MatchString(it) && code.FindString(it) == it:
				continue
			case depRefRE.MatchString(it):
				if !deps[it] {
					broken = append(broken, fmt.Sprintf("%s → `%s`", u.Rule, it))
				}
			default:
				if !names[strings.ToLower(it)] && !names[strings.ToLower(rootName(it))] {
					broken = append(broken, fmt.Sprintf("%s → `%s`", u.Rule, it))
				}
			}
		}
	}
	if len(broken) > 0 {
		return Fail, i18n.T("gate.rule_uses_resolve.fail_unresolved", len(broken), strings.Join(broken, "; "))
	}
	return Pass, ""
}

// declaredNames are the names a spec declares outside the three sections: the first cell of
// every table row, and the backticked names of headings.
func declaredNames(lines []string) map[string]bool {
	out := map[string]bool{}
	add := func(s string) {
		s = strings.TrimSpace(strings.Trim(strings.TrimSpace(s), "`*"))
		if s != "" && !strings.HasPrefix(strings.ToUpper(s), "TODO") {
			out[strings.ToLower(s)] = true
			out[strings.ToLower(rootName(s))] = true
		}
	}
	for _, l := range lines {
		t := strings.TrimSpace(l)
		switch {
		case strings.HasPrefix(t, "|") && !tableDividerR.MatchString(t):
			cells := cellsOf(t)
			if ms := backtickedRE.FindAllStringSubmatch(cells[0], -1); len(ms) > 0 {
				for _, m := range ms {
					add(m[1])
				}
			} else {
				add(cells[0])
			}
		case headingRE.MatchString(t):
			for _, m := range backtickedRE.FindAllStringSubmatch(t, -1) {
				add(m[1])
			}
		}
	}
	return out
}

// rootName is the first segment of a dotted or indexed name: `summary` of
// `summary.balance`, `items` of `items[0].id`.
func rootName(s string) string {
	if i := strings.IndexAny(s, ".[("); i > 0 {
		return s[:i]
	}
	return s
}

func listCodes(codes []string) string {
	sort.Strings(codes)
	if len(codes) > 12 {
		return strings.Join(codes[:12], ", ") + fmt.Sprintf(" … (+%d)", len(codes)-12)
	}
	return strings.Join(codes, ", ")
}

// rule-uses-implemented: what a rule says it reads, the code the spec governs reads too. A
// field no code file of the unit mentions leaves the rule orphan of code: declared, maybe
// tested, and checking nothing the implementation touches. The search is plain text — the
// field's name as a whole word, or its first or last segment for a dotted name
// (`summary.balance`) —, so no language is assumed. Codes and dependency rows are not
// fields, and are left to the other gates.
func checkRuleUsesImplemented(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.rule_uses.skip_not_spec")
	}
	rows := ruleUsesOf(content, cfg, n.Layer)
	if len(rows) == 0 {
		return Skip, i18n.T("gate.rule_uses_resolve.skip_no_rows")
	}
	code, found := governedCode(n, root, g)
	if !found {
		return Skip, i18n.T("gate.rule_uses_implemented.skip_no_code")
	}
	codeRE := ruleCodeRE()
	var orphan []string
	for _, u := range rows {
		for _, it := range u.Uses {
			if codeRE.FindString(it) == it || depRefRE.MatchString(it) {
				continue
			}
			if !mentions(code, it) {
				orphan = append(orphan, fmt.Sprintf("%s → `%s`", u.Rule, it))
			}
		}
	}
	if len(orphan) > 0 {
		return Fail, i18n.T("gate.rule_uses_implemented.fail", len(orphan), strings.Join(orphan, "; "))
	}
	return Pass, ""
}

// mentions says whether the code names the field as a whole word — the name itself, or the
// first or last segment of a dotted or indexed one.
func mentions(code, field string) bool {
	parts := strings.FieldsFunc(field, func(r rune) bool { return r == '.' || r == '[' || r == ']' || r == '(' || r == ')' })
	candidates := []string{field}
	if len(parts) > 1 {
		candidates = append(candidates, parts[0], parts[len(parts)-1])
	}
	for _, c := range candidates {
		if c == "" {
			continue
		}
		if regexp.MustCompile(`(^|[^\p{L}\p{N}_])` + regexp.QuoteMeta(c) + `($|[^\p{L}\p{N}_])`).MatchString(code) {
			return true
		}
	}
	return false
}
