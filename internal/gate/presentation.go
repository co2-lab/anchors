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

// THE PRESENTATION VALIDATIONS, confronted.
//
// Each row of `Presentation validations` says: this prop or state, under this condition,
// makes the unit look like this. Four questions follow from that shape, and all four are
// answered from the spec's own text — no language, no rendering:
//
//	presentation-exhaustive        every value the prop can take has an appearance decided
//	presentation-conflict          one prop and one condition do not lead to two appearances
//	presentation-copy-single-source the text shown is a message code, not copy repeated here
//	presentation-observable        what changes is something a test can point at

// presentationRow is one row: rule, what it reads, the condition, the appearance.
type presentationRow struct {
	Rule, Reads, Condition, Appearance string
}

func presentationRows(content string, cfg *config.Config, layer string) []presentationRow {
	inside, _ := splitSections(content, []string{"presentation_validations"}, cfg, layer)
	code := ruleCodeRE()
	var out []presentationRow
	for _, l := range inside {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "|") || tableDividerR.MatchString(t) {
			continue
		}
		cells := cellsOf(t)
		if len(cells) < 4 {
			continue
		}
		rule := code.FindString(cells[0])
		if rule == "" || strings.HasPrefix(strings.ToUpper(strings.TrimSpace(cells[1])), "TODO") {
			continue
		}
		out = append(out, presentationRow{Rule: rule, Reads: strings.TrimSpace(cells[1]),
			Condition: strings.TrimSpace(cells[2]), Appearance: strings.TrimSpace(strings.Join(cells[3:], "|"))})
	}
	return out
}

func normCell(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.NewReplacer("`", "", "'", "", "\"", "").Replace(s))), " ")
}

// presentationSkip is the common skip: not a spec, or no row to confront.
func presentationRowsOrSkip(content string, n mapx.Node, cfg *config.Config) ([]presentationRow, Verdict, string, bool) {
	if n.Kind != mapx.KindSpec {
		return nil, Skip, i18n.T("gate.rule_uses.skip_not_spec"), false
	}
	rows := presentationRows(content, cfg, n.Layer)
	if len(rows) == 0 {
		return nil, Skip, i18n.T("gate.presentation.skip_no_rows"), false
	}
	return rows, Pass, "", true
}

// ── exhaustive ──────────────────────────────────────────────────────────────

// valueSets are the values the spec declares for its props and states: the backticked
// values of the props table's second column, and the first column of the tables under a
// backticked heading (a data state, a variant axis). A set of fewer than two values is not
// an enumeration, and is left out.
func valueSets(content string) map[string][]string {
	out := map[string][]string{}
	add := func(name string, values []string) {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" || strings.HasPrefix(name, "todo") {
			return
		}
		for _, v := range values {
			if v = normCell(v); v != "" && !strings.HasPrefix(v, "todo") && !contains(out[name], v) {
				out[name] = append(out[name], v)
			}
		}
	}
	props, _ := splitSections(content, []string{"props"}, nil, "")
	for _, l := range props {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "|") || tableDividerR.MatchString(t) {
			continue
		}
		cells := cellsOf(t)
		if len(cells) < 2 {
			continue
		}
		var vals []string
		for _, m := range backtickedRE.FindAllStringSubmatch(cells[1], -1) {
			for _, v := range regexp.MustCompile(`\s*\\?\|\s*`).Split(m[1], -1) {
				vals = append(vals, v)
			}
		}
		add(strings.Trim(strings.TrimSpace(cells[0]), "`"), vals)
	}
	current := ""
	for _, l := range strings.Split(content, "\n") {
		t := strings.TrimSpace(l)
		if m := headingRE.FindStringSubmatch(t); m != nil {
			current = ""
			if b := backtickedRE.FindStringSubmatch(m[2]); b != nil {
				current = b[1]
			}
			continue
		}
		if current == "" || !strings.HasPrefix(t, "|") || tableDividerR.MatchString(t) {
			continue
		}
		cells := cellsOf(t)
		if isHeaderRow(cells[0]) {
			continue
		}
		add(current, []string{cells[0]})
	}
	for k, v := range out {
		if len(v) < 2 {
			delete(out, k)
		}
	}
	return out
}

// isHeaderRow says whether a first cell is a table's header in any supported language.
func isHeaderRow(cell string) bool {
	c := normCell(cell)
	for _, h := range i18n.AllTranslations("gate.presentation.value_header") {
		for _, w := range strings.Split(h, "|") {
			if c == strings.TrimSpace(w) {
				return true
			}
		}
	}
	return false
}

func contains(xs []string, x string) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

// otherwiseCondition says whether a condition covers every value not named elsewhere, in
// any supported language ("otherwise", "demais valores", "cualquier otro").
func otherwiseCondition(cond string) bool {
	c := normCell(cond)
	for _, t := range i18n.AllTranslations("gate.presentation.otherwise") {
		for _, w := range strings.Split(t, "|") {
			if w = strings.TrimSpace(w); w != "" && strings.Contains(c, w) {
				return true
			}
		}
	}
	return false
}

func checkPresentationExhaustive(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	rows, v, msg, ok := presentationRowsOrSkip(content, n, cfg)
	if !ok {
		return v, msg
	}
	sets := valueSets(content)
	conds := map[string][]string{}
	for _, r := range rows {
		for _, name := range usedItems(r.Reads) {
			conds[strings.ToLower(name)] = append(conds[strings.ToLower(name)], r.Condition)
		}
	}
	var gaps []string
	for name, cs := range conds {
		values := sets[name] // a prop with no declared set has nothing to miss
		covered := false
		for _, c := range cs {
			covered = covered || otherwiseCondition(c)
		}
		if covered {
			continue
		}
		var missing []string
		for _, val := range values {
			hit := false
			for _, c := range cs {
				hit = hit || strings.Contains(normCell(c), val)
			}
			if !hit {
				missing = append(missing, val)
			}
		}
		if len(missing) > 0 {
			gaps = append(gaps, fmt.Sprintf("`%s`: %s", name, strings.Join(missing, ", ")))
		}
	}
	if len(gaps) > 0 {
		sort.Strings(gaps)
		return Fail, i18n.T("gate.presentation_exhaustive.fail", strings.Join(gaps, "; "))
	}
	return Pass, ""
}

// ── conflict ────────────────────────────────────────────────────────────────

func checkPresentationConflict(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	rows, v, msg, ok := presentationRowsOrSkip(content, n, cfg)
	if !ok {
		return v, msg
	}
	first := map[string]presentationRow{}
	var conflicts []string
	for _, r := range rows {
		k := normCell(r.Reads) + "\x00" + normCell(r.Condition)
		if prev, seen := first[k]; seen {
			if normCell(prev.Appearance) != normCell(r.Appearance) {
				conflicts = append(conflicts, fmt.Sprintf("%s × %s", prev.Rule, r.Rule))
			}
			continue
		}
		first[k] = r
	}
	if len(conflicts) > 0 {
		return Fail, i18n.T("gate.presentation_conflict.fail", strings.Join(conflicts, "; "))
	}
	return Pass, ""
}

// ── copy single source ──────────────────────────────────────────────────────

var (
	quotedCopyRE = regexp.MustCompile(`["“«][^"”»]*\p{L}{2,}[^"”»]*["”»]`)
	messageRefRE = regexp.MustCompile(`[A-Z0-9]{3,8}-M\d{2}`)
)

func checkPresentationCopySingleSource(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	rows, v, msg, ok := presentationRowsOrSkip(content, n, cfg)
	if !ok {
		return v, msg
	}
	var copies []string
	for _, r := range rows {
		if m := quotedCopyRE.FindString(r.Appearance); m != "" && !messageRefRE.MatchString(r.Appearance) {
			copies = append(copies, fmt.Sprintf("%s → %s", r.Rule, m))
		}
	}
	if len(copies) > 0 {
		return Fail, i18n.T("gate.presentation_copy.fail", strings.Join(copies, "; "))
	}
	return Pass, ""
}

// ── observable ──────────────────────────────────────────────────────────────

func checkPresentationObservable(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	rows, v, msg, ok := presentationRowsOrSkip(content, n, cfg)
	if !ok {
		return v, msg
	}
	ids, _ := splitSections(content, []string{"testids"}, cfg, n.Layer)
	declared := declaredNames(ids)
	var blind []string
	for _, r := range rows {
		seen := false
		for _, m := range backtickedRE.FindAllStringSubmatch(r.Appearance, -1) {
			seen = seen || declared[strings.ToLower(strings.TrimSpace(m[1]))]
		}
		if !seen {
			blind = append(blind, r.Rule)
		}
	}
	if len(blind) > 0 {
		return Fail, i18n.T("gate.presentation_observable.fail", listCodes(blind))
	}
	return Pass, ""
}
