// @anchors
//   code: CIGCN
//   ref: CTRIM

package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// contract-impact: a field changed in a spec's data → the rules that use it, and their
// tests, are named — the moment it changes, before the commit.
//
// A spec's file revision says the spec changed, not WHAT changed: editing one field made
// every relation of the spec stale, and nothing told which rules read that field. The rows
// of Validations, Presentation validations and Rule uses say what each rule reads, so a
// field whose row differs from the last commit names its rules: those of this spec that
// use it, and those of the specs that depend on this unit's code and use a field of that
// name. The change is read from git — the spec as it is against the spec at HEAD —, so
// nothing new is kept in the map, and a committed change is no longer "changing".
//
// The same answer feeds the test selection (`ImpactedTests`): a test of an affected rule
// runs even when its own file did not move.

// Impact is one changed field and what reads it.
type Impact struct {
	Spec, Field string
	Rules       []string // the rule codes that use the field
	Tests       []string // the test files whose titles cite one of those rules
}

// fieldRows are the rows a spec declares what its rules read in: every table row outside
// the three rule-use sections, by the name in its first cell, with the row's text. A code
// is a name too: a rule that uses a state (`PAYMT-S01`) is reached when the state's row
// changes, as one that uses a field is. A header row may enter as a "field": it changes only when the
// columns are renamed, and no rule uses a field named after a column.
func fieldRows(content string, cfg *config.Config, layer string) map[string]string {
	_, outside := splitRuleUseSections(content, cfg, layer)
	out := map[string]string{}
	for _, l := range outside {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "|") || tableDividerR.MatchString(t) {
			continue
		}
		cells := cellsOf(t)
		name := strings.TrimSpace(strings.Trim(strings.TrimSpace(cells[0]), "`*"))
		if name == "" || depRefRE.MatchString(name) || strings.HasPrefix(strings.ToUpper(name), "TODO") {
			continue
		}
		out[strings.ToLower(name)] = normCell(t)
	}
	return out
}

// changedFields are the names whose row is not what HEAD has — changed or removed. A field
// added is new, and nothing read it before.
func changedFields(now, before map[string]string) []string {
	var out []string
	for name, row := range before {
		if now[name] != row {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// ContractImpacts compares a spec with its last commit and returns, per changed field, the
// rules that use it — in this spec and in the specs that depend on its code — and the
// tests that cite those rules. No repository, a new spec or no change is no impact.
func ContractImpacts(root, specID, content string, g *mapx.Graph, cfg *config.Config) []Impact {
	head, ok := gitmeta.AtHead(root, specID)
	if !ok {
		return nil
	}
	layer := ""
	if g != nil {
		for _, n := range g.Nodes {
			if n.ID == specID {
				layer = n.Layer
			}
		}
	}
	changed := changedFields(fieldRows(content, cfg, layer), fieldRows(head, cfg, layer))
	if len(changed) == 0 {
		return nil
	}
	// What this change's own revision says it revised or checked is answered: the rules a
	// revision added since HEAD names in `Revises:` or `Checked:`. The impact exists only
	// while the change is uncommitted, so with no way to answer it a project whose
	// divergence blocks could never commit a field change (reported from jokenpo).
	//
	// A rule is answered where it lives: by a revision added to its own spec — the changed
	// one or a spec that reads its data —, short code or full; or by one added to the changed
	// spec naming it by its full code. A short code answers only its own spec's rule: `B12`
	// in this spec's revision said nothing of another unit's `B12`. A spec not yet in git
	// has every revision added (reported from jokenpo).
	readers := []specText{{specID, content, layer}}
	readers = append(readers, dependents(root, specID, g)...)
	ownAcked := map[string]map[string]bool{specID: acknowledgedRules(content, head)}
	fullAcked := acknowledgedFullCodes(content, head)
	ackedIn := func(r specText) map[string]bool {
		if a, ok := ownAcked[r.id]; ok {
			return a
		}
		h, _ := gitmeta.AtHead(root, r.id)
		ownAcked[r.id] = acknowledgedRules(r.content, h)
		return ownAcked[r.id]
	}
	var out []Impact
	for _, f := range changed {
		imp := Impact{Spec: specID, Field: f}
		answered := map[string]bool{}
		for _, r := range readers {
			for _, u := range ruleUsesOf(r.content, cfg, r.layer) {
				for _, it := range u.Uses {
					if strings.EqualFold(it, f) || strings.EqualFold(rootName(it), f) {
						imp.Rules = appendOnce(imp.Rules, u.Rule)
						if ackedIn(r)[shortRule(u.Rule)] || fullAcked[u.Rule] {
							answered[u.Rule] = true
						}
					}
				}
			}
		}
		if len(imp.Rules) == 0 {
			continue
		}
		open := false
		for _, r := range imp.Rules {
			if !answered[r] {
				open = true
				break
			}
		}
		if !open {
			continue
		}
		sort.Strings(imp.Rules)
		imp.Tests = testsCiting(root, g, cfg, imp.Rules)
		out = append(out, imp)
	}
	return out
}

type specText struct{ id, content, layer string }

// dependents are the specs that depend on the code this spec governs: they read its data.
func dependents(root, specID string, g *mapx.Graph) []specText {
	if g == nil {
		return nil
	}
	governed := map[string]bool{}
	for _, e := range g.Edges {
		if e.From == specID && e.Type == mapx.EdgeSpecifies {
			governed[e.To] = true
		}
	}
	kindOf := map[string]mapx.Node{}
	for _, n := range g.Nodes {
		kindOf[n.ID] = n
	}
	var out []specText
	seen := map[string]bool{specID: true}
	for _, e := range g.Edges {
		if e.Type != mapx.EdgeDependsOn || !governed[e.To] || seen[e.From] || kindOf[e.From].Kind != mapx.KindSpec {
			continue
		}
		seen[e.From] = true
		if b, err := readFile(root, e.From); err == nil {
			out = append(out, specText{e.From, string(b), kindOf[e.From].Layer})
		}
	}
	return out
}

// testsCiting are the test files whose test titles cite one of the rules.
func testsCiting(root string, g *mapx.Graph, cfg *config.Config, rules []string) []string {
	all, _, err := projectTests(root, g, cfg)
	if err != nil {
		return nil
	}
	var out []string
	for _, t := range all {
		for _, r := range rules {
			if regexp.MustCompile(`(^|[^A-Z0-9-])` + regexp.QuoteMeta(r) + `($|[^0-9])`).MatchString(t.Title) {
				out = appendOnce(out, t.File)
			}
		}
	}
	sort.Strings(out)
	return out
}

func appendOnce(xs []string, x string) []string {
	for _, y := range xs {
		if y == x {
			return xs
		}
	}
	return append(xs, x)
}

// ImpactedTests are the test files of every rule a field changed since the last commit
// reaches, across the specs with uncommitted changes. It is what the test selection adds to
// a run: those tests' own files did not move, and their rule's input did.
func ImpactedTests(root string, g *mapx.Graph, cfg *config.Config) []string {
	if g == nil {
		return nil
	}
	var out []string
	for _, n := range g.Nodes {
		// A clean spec is its HEAD: the check is only what spares a `git show` per spec.
		if n.Kind != mapx.KindSpec || !gitmeta.HasUncommittedChanges(root, n.ID) {
			continue
		}
		b, err := readFile(root, n.ID)
		if err != nil {
			continue
		}
		for _, imp := range ContractImpacts(root, n.ID, string(b), g, cfg) {
			for _, t := range imp.Tests {
				out = appendOnce(out, t)
			}
		}
	}
	sort.Strings(out)
	return out
}

func checkContractImpact(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.rule_uses.skip_not_spec")
	}
	imps := ContractImpacts(root, n.ID, content, g, cfg)
	if len(imps) == 0 {
		return Pass, ""
	}
	var items []string
	for _, imp := range imps {
		tests := "—"
		if len(imp.Tests) > 0 {
			tests = strings.Join(imp.Tests, ", ")
		}
		items = append(items, fmt.Sprintf(i18n.T("gate.contract_impact.item"), imp.Field, strings.Join(imp.Rules, ", "), tests))
	}
	return Diverge, i18n.T("gate.contract_impact.pending", strings.Join(items, "; "))
}

// acknowledgedRules are the rules — by their short code, `R03` — that the revision lines
// added since HEAD name in `Revises:` or `Checked:`.
func acknowledgedRules(content, head string) map[string]bool {
	text := addedLines(content, head)
	out, _ := revisedCodes(text)
	for c := range codesIn(checkedRE(), text) {
		out[c] = true
	}
	return out
}

// addedLines are the lines of content that head does not have, as text.
func addedLines(content, head string) string {
	before := map[string]int{}
	for _, l := range strings.Split(head, "\n") {
		before[l]++
	}
	var added strings.Builder
	for _, l := range strings.Split(content, "\n") {
		if before[l] > 0 {
			before[l]--
			continue
		}
		added.WriteString(l)
		added.WriteString("\n")
	}
	return added.String()
}

// acknowledgedFullCodes are the full rule codes (`ARNAA-A04`) the revisions added since HEAD
// name in `Revises:` or `Checked:` — how the changed spec answers another spec's rule.
func acknowledgedFullCodes(content, head string) map[string]bool {
	text := addedLines(content, head)
	out := map[string]bool{}
	for _, re := range []*regexp.Regexp{revisesRE(), checkedRE()} {
		for _, m := range re.FindAllStringSubmatch(text, -1) {
			for _, c := range fullRuleRefRE.FindAllString(m[1], -1) {
				out[c] = true
			}
		}
	}
	return out
}

// fullRuleRefRE is a rule code with its unit.
var fullRuleRefRE = regexp.MustCompile(`\b[A-Z0-9]{3,6}-[A-Z]\d{2}\b`)

// shortRule is a rule's code without its unit: `DATAR-R03` → `R03`.
func shortRule(code string) string {
	if i := strings.LastIndex(code, "-"); i >= 0 {
		return code[i+1:]
	}
	return code
}
