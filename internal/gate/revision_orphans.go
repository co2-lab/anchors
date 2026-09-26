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

// --- the revision that changed what a word means, and did not tell whom ---
//
// A rule does not live alone. It shares vocabulary with its siblings, and it is that
// vocabulary a revision changes — not only the text of the rule it rewrites.
//
// MEASURED in the reference app, and it is what produced this gate. `NTCNN-R0002` changed
// the bell badge from a COUNT to a DOT and named the rules it rewrote: `B03`, `B04`, `B07`.
// The invariant `I02` was not named — and it is titled "the badge never COUNTS what the
// list does not show", with a body saying "the NUMBER on the bell matches what appears on
// opening". The invariant governed arithmetic the revision had abolished.
//
// NOTHING ACCUSED IT. `B03` was correct, `I02` well-formed, the triad complete, the suite
// green. The contradiction surfaced MONTHS later, when another agent went to implement and
// could not tell which of the two to follow — and it became a decision that had to
// escalate to the user, with nobody left remembering the context. Seven contradictions of
// this exact shape surfaced in a single batch.
//
// THE RULER IS CO-CITATION, not meaning. Asking whether two rules contradict each other
// requires reading them, and that is judgement — it would make this a judge, not a gate.
// What a machine decides alone is narrower and sufficient: which rules of this unit share
// the vocabulary of the rules the revision touched, and were not mentioned.

// revisesRE and checkedRE match the two fields of a revision.
//
// The keywords come from the TRANSLATION CATALOG, never hardcoded. It is the lesson the
// flow axis already taught: the first version of `fits` carried `(?:Encaixa|Fits)` nailed
// into the pattern, and a project writing in Spanish had no way to declare anything
// without editing the engine.
func revisesRE() *regexp.Regexp { return keywordListRE("revision.keyword.revises") }
func checkedRE() *regexp.Regexp { return keywordListRE("revision.keyword.checked") }

// keywordListRE builds the pattern for a field that lists rule codes:
// `**Revises:** ` + "`B03`, `B07`" + `.
func keywordListRE(key string) *regexp.Regexp {
	return regexp.MustCompile("(?im)^[^\\S\\n]*(?:>|#{1,6})?[^\\S\\n]*(?:\\*\\*)?(?:" +
		strings.Join(escapeKeywords(i18n.AllTranslations(key)), "|") +
		")(?:\\*\\*)?[^\\S\\n]*:[^\\S\\n]*(\\S.*)$")
}

// escapeKeywords prepares translated keywords to enter a regex alternation.
func escapeKeywords(xs []string) []string {
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		if x != "" {
			out = append(out, regexp.QuoteMeta(x))
		}
	}
	return out
}

// ruleRefRE finds the SHORT codes a field lists (`B03`, `I02`).
//
// Short rather than full (`NTCNN-B03`) because that is how a revision is written: it is
// already inside the unit, and repeating the prefix on every item would be noise. The long
// form matches too — whoever writes switches between the two without thinking about it.
var ruleRefRE = regexp.MustCompile(`\b(?:[A-Z0-9]{3,6}-)?([A-Z]\d{2})\b`)

// --- the gate ---
func checkRevisionOrphans(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.revision_orphans.skip_not_spec")
	}

	revs := revisionRE().FindAllStringSubmatch(content, -1)
	if len(revs) == 0 {
		return Skip, i18n.T("gate.revision_orphans.no_revision")
	}

	// The rules THIS spec defines — ALL of them, including those waived from a scenario.
	//
	// `definedRequirements`, from `spec-feature-match`, drops a rule carrying
	// `@no-scenario`, and it is right to there: it asks "does this requirement have a
	// scenario?", and the waiver answers. Here the question is different — "does this rule
	// assert something the revision changed?" — and a scenario waiver says nothing about it.
	//
	// MEASURED on the spec that produced the gate: `NTCNN-B07` carries `@no-scenario` (it
	// is a layout decision, with no rendering to exercise) and is one of the three rules
	// `R0002` rewrote. Reading it the neighbour's way, the gate reported `B07` as an unknown
	// code — a false positive on the very rule that motivated it.
	titleByShort := ruleTitles(content)
	revised := codesIn(revisesRE(), content)
	checked := codesIn(checkedRE(), content)
	// No rules AND nothing revised: nothing to confront. With a `Revises:` naming codes,
	// a spec whose rules were not recognised is exactly the case B04 exists for — every
	// named code is unknown — and skipping here let that revision pass unseen.
	if len(titleByShort) == 0 && len(revised) == 0 {
		return Skip, i18n.T("gate.revision_orphans.no_rules")
	}

	// WITHOUT `Revises:` THE GATE ABSTAINS instead of accusing.
	//
	// The field is new and the revisions already written do not carry it: measured in the
	// reference app, 439 revisions in 141 specs, none with `Revises:` — and the first
	// version of this gate failed 174 specs at once, all with the same message.
	//
	// A hundred and seventy-four identical findings are not a queue, they are noise:
	// whoever opens the report learns to scroll past them, and the REAL finding — the
	// orphaned sibling — gets lost among them. Pending states what is true ("this revision
	// cannot be confronted") without charging whoever wrote before the ruler existed.
	//
	// The charge comes from `plan-change-justified`, which anchors on `--changed`: a NEW
	// revision is born under the new ruler, and that is where the field becomes required.
	if len(revised) == 0 {
		return Pending, fmt.Sprintf(i18n.T("gate.revision_orphans.no_revises"), len(revs))
	}

	// A named code the spec does not define: this axis's `ref-resolves`.
	var unknown []string
	for c := range revised {
		if _, ok := titleByShort[c]; !ok {
			unknown = append(unknown, c)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return Fail, fmt.Sprintf(i18n.T("gate.revision_orphans.unknown_rule"),
			len(unknown), strings.Join(unknown, ", "))
	}

	// The VOCABULARY of the revised rules.
	revisedTerms := map[string]bool{}
	for c := range revised {
		for t := range termsOf(titleByShort[c]) {
			revisedTerms[t] = true
		}
	}
	if len(revisedTerms) == 0 {
		return Pass, ""
	}

	type orphan struct {
		code  string
		terms []string
	}
	var orphans []orphan
	for short, title := range titleByShort {
		// A rule never accuses itself, and what was already checked leaves the list.
		if revised[short] || checked[short] {
			continue
		}
		var shared []string
		for t := range termsOf(title) {
			if revisedTerms[t] {
				shared = append(shared, t)
			}
		}
		// See `discriminates` for why a single cleaned domain word is enough.
		if !discriminates(shared) {
			continue
		}
		sort.Strings(shared)
		orphans = append(orphans, orphan{short, shared})
	}
	if len(orphans) == 0 {
		return Pass, ""
	}
	sort.Slice(orphans, func(i, j int) bool { return orphans[i].code < orphans[j].code })

	var b strings.Builder
	fmt.Fprintf(&b, i18n.T("gate.revision_orphans.orphans"), len(orphans))
	for _, o := range orphans {
		fmt.Fprintf(&b, "\n    %s — %q\n      %s", o.code,
			trimTitle(titleByShort[o.code]), strings.Join(o.terms, ", "))
	}
	b.WriteString("\n" + i18n.T("gate.revision_orphans.how_to_clear"))
	return Fail, b.String()
}

// codesIn gathers the short codes a field lists, across all its occurrences.
func codesIn(re *regexp.Regexp, content string) map[string]bool {
	out := map[string]bool{}
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		for _, r := range ruleRefRE.FindAllStringSubmatch(m[1], -1) {
			out[r[1]] = true
		}
	}
	return out
}

// ruleTitles indexes the TITLE of each defined rule by its short code.
//
// The title and not the body: it is where the rule ASSERTS what it asserts, in one line,
// and it is what a revision contradicts when it contradicts. The body carries
// justification prose, and including it would make almost every rule share vocabulary
// with almost every other.
func ruleTitles(content string) map[string]string {
	out := map[string]string{}
	re := defineRuleCaptureRE()
	for _, line := range strings.Split(content, "\n") {
		m := re.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// A WAIVER IN A COMMENT is not part of what the rule asserts, and letting it in
		// poisons the comparison: measured on the real spec, `B07` carries two `@no-*` with
		// a written reason and came out with 34 terms — against 4 to 6 for its siblings —,
		// sharing vocabulary with almost all of them by accident of prose.
		if i := strings.Index(line, "<!--"); i >= 0 {
			line = line[:i]
		}
		_, short, ok := strings.Cut(m[1], "-")
		if !ok {
			continue
		}
		out[short] = line
	}
	return out
}

// termsOf reduces a title to its vocabulary, as a SET.
//
// It reuses `significantTerms` from `feature-test-match` — the same question ("which
// words of this text say what it is about?"), and two stopword lists would diverge on the
// first word someone added to only one of them.
//
// The CODE is removed first: `NTCNN-B03` would bring `NTCNN` into every title of the
// unit, and a term everyone shares discriminates nothing.
func termsOf(title string) map[string]bool {
	out := map[string]bool{}
	for _, w := range significantTerms(anyCodeRE.ReplaceAllString(title, " ")) {
		if !extraStopwords[w] {
			out[w] = true
		}
	}
	return out
}

// extraStopwords are the words `descStopwords` does not need to drop and this ruler does.
//
// NEGATION costs the most: half the rules of a well-written spec state what the unit does
// NOT do, and "não" linked every rule to every other. Measured on the spec that produced
// the gate — of the six orphans the first version reported, four shared nothing but "não".
//
// It does not go into `descStopwords` because the question there is different:
// `feature-test-match` compares a scenario's description with the test body, and there
// negation DOES discriminate (a test that asserts and one that denies prove different
// things).
//
// The Portuguese and Spanish words are DATA, not prose: they are the negations that appear
// in the specs of projects written in those languages, and the ruler must recognise them.
var extraStopwords = map[string]bool{
	"não": true, "nao": true, "nunca": true, "nenhum": true, "nenhuma": true,
	"not": true, "never": true, "none": true, "no": true,
	"ni": true,
	// The spec's STRUCTURAL vocabulary, not its domain's.
	"regra": true, "rule": true, "spec": true, "unidade": true, "unit": true,
}

// discriminates decides whether the shared vocabulary points at the SAME subject.
//
// ONE DOMAIN WORD IS ENOUGH, and the threshold was measured in both directions.
//
// The first version required TWO shared words, and the case that produced the gate did not
// pass: on the real spec, `I02` shares exactly one word with what `R0002` rewrote —
// `badge` —, and that word carries the whole contradiction.
//
// What makes one word sufficient is CLEANING the title, not counting. Two measurements, on
// the same spec:
//
//	not cleaned   3 reported — `B01` and `B05` entered on "não" alone
//	cleaned       1 reported — `I02`, which is the target
//
// The two sources of noise were NEGATION (half the rules of a well-written spec state what
// the unit does not do) and the WAIVER IN A COMMENT (`B07` came out with 34 terms against
// 4 for its siblings). With both removed, what remains is domain vocabulary — and there a
// single coincidence is already a signal.
//
// A ubiquity filter was also tried (dropping terms that appear in more than two thirds of
// the rules). It dropped `lista` and `badge` — precisely the subject — and made the real
// case accuse nothing. In a well-written spec the domain vocabulary repeats on purpose:
// dropping what repeats is dropping the subject.
func discriminates(shared []string) bool {
	return len(shared) >= 1
}

// trimTitle makes a title readable in the verdict: no `###`, no code, no dash.
func trimTitle(line string) string {
	s := strings.TrimSpace(strings.TrimLeft(line, "#> *|`"))
	if m := anyCodeRE.FindStringIndex(s); m != nil {
		s = s[m[1]:]
	}
	// `|` too: in a table row the code's cell ends in one, and the title came out as
	// "| text |" in the orphan message.
	s = strings.TrimSpace(strings.TrimLeft(s, "—–-: `|"))
	// A waiver in an HTML comment is not part of what the rule asserts.
	if i := strings.Index(s, "<!--"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	return strings.TrimSpace(strings.TrimRight(s, "| "))
}
