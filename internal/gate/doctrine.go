package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/co2-lab/anchors/internal/similarity"
)

// --- the VERTICAL AXIS: product doctrine and whoever realizes it ---
//
// Every spec has a TARGET, and co-location ties the triad to one directory. A business
// rule that holds for three screens belongs to none of the three — it lives in
// `product/<name>.doctrine.md`, and the spec points at it with `@realizes`.
//
// The gates here answer distinct questions, and each catches a silence the others cannot
// see:
//
//	does the doctrine the PLAN cites exist?   plan-doctrine-exists    (fails)
//	is the doctrine tied to any spec?         doctrine-realized       (informs)
//	does the doctrine the SPEC cites exist?   spec-doctrine-exists    (fails)

// tbdLineRE — this axis's declared way out: `@TBD: <reason>` (to be developed).
//
// It is `@TBD` and not `@no-...` because the two say DIFFERENT things, and the
// difference decides the finding's fate. `@no-<thing>: <reason>` asserts "it will never
// have one, and here is why" — a permanent waiver, and the gate goes quiet for good.
// `@TBD` asserts "it does not have one YET" — assumed debt, which keeps showing up until
// somebody pays it.
//
// On the product axis almost every case is the second: the cross-cutting rule is DECIDED
// before it is implemented, and that is how product works. Treating it as a permanent
// waiver would erase from the radar exactly the work that remains — and the `Pending`
// verdict marked as debt becomes an issue in `future/`, which is where it belongs.
//
// A BARE marker does not count: the pattern demands the `:` and text after it. Debt with
// no written reason is the silence the gates exist to end.
//
// The BACKTICK guard comes from `triad-complete`, where the lesson was already paid for:
// a revision EXPLAINING the removal of a waiver cites the marker ("the `@TBD: ...` waiver
// is gone"), and without the guard that citation REACTIVATES it. An active marker is
// never inside backticks — it is the declaration, not a mention of one.
var tbdLineRE = regexp.MustCompile("(?i)(^|[^`])@TBD[^\\S\\n]*:[^\\S\\n]*\\S+")

// doctrineRuleRE finds a catalogued rule in the THREE valid forms (heading, table row,
// bold bullet) — the same grammar the specs use, because doctrine catalogues rules the
// same way.
var doctrineRuleRE = regexp.MustCompile(`(?m)(?:^#{1,6}\s+|^\s*\|\s*` + "`?" + `|^\s*-\s+\*\*)([A-Z0-9]{3,6}-[A-Z]{1,2}[0-9]{2})`)

// doctrineSeedCiteRE finds the doctrine a plan cites inside backticks.
var doctrineSeedCiteRE = regexp.MustCompile("`([^`]+\\.doctrine\\.md)`")

// --- does the doctrine the PLAN cites exist? ---
//
// A plan promises the cross-cutting rule before it exists — that is what seeding means.
// Without this gate the promise carries no charge: the specs the plan also seeds are born
// with nothing to realize, and they PASS every gate, because a spec with no `@realizes`
// is legitimate. The silence is total.
func checkPlanDoctrineExists(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindPlan {
		return Skip, i18n.T("gate.plan_doctrine_exists.skip_not_plan")
	}
	var missing, deferred []string
	seen := map[string]bool{}
	for _, line := range strings.Split(content, "\n") {
		// `@TBD` holds for the LINE of the citation: a plan may seed five doctrines and
		// declare debt for only one of them.
		//
		// And debt is NOT a waiver: `@no-*` says "it will never have one", `@TBD` says
		// "it does not have one yet". That is why deferred doctrine does not leave the
		// radar — it becomes a Pending marked as DEBT, which keeps showing up until
		// somebody pays it, instead of a Pass that erases it.
		isDeferred := tbdLineRE.MatchString(line)
		for _, m := range doctrineSeedCiteRE.FindAllStringSubmatch(line, -1) {
			d := m[1]
			// A prose citation is not seeding — the same rule `plan-seeds-valid` applies
			// to specs: only a real path is a promise to create a file.
			if seen[d] || !strings.Contains(d, "/") || strings.HasPrefix(filepath.Base(d), "_TEMPLATE") {
				continue
			}
			seen[d] = true
			if _, err := os.Stat(filepath.Join(root, d)); err != nil {
				if isDeferred {
					deferred = append(deferred, d)
				} else {
					missing = append(missing, d)
				}
			}
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Fail, fmt.Sprintf(i18n.T("gate.plan_doctrine_exists.missing"), len(missing), strings.Join(missing, ", "))
	}
	if len(deferred) > 0 {
		sort.Strings(deferred)
		return Pending, fmt.Sprintf(i18n.T("gate.plan_doctrine_exists.deferred"), len(deferred), strings.Join(deferred, ", "))
	}
	if len(seen) == 0 {
		return Skip, i18n.T("gate.plan_doctrine_exists.skip_none")
	}
	return Pass, ""
}

// --- is the doctrine tied to any spec? ---
//
// INFORMATIVE, and the difference matters: a product rule written TODAY to be realized
// next cycle is legitimate work, not a defect. Failing it would have the gate demand that
// implementation keep pace with the decision in the same commit, which is not how product
// works.
//
// But it cannot be silence either: a rule written and never realized is a decision that
// never reached the code, and the doctrine file passes every other gate because it is
// itself well-formed. Nothing else looks at the far end.
func checkDoctrineRealized(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindProduct {
		return Skip, i18n.T("gate.doctrine_realized.skip_not_doctrine")
	}
	if g == nil {
		return pendingNoMap()
	}

	// The rules THIS doctrine catalogues, minus the ones deferred on their own line.
	rules := map[string]bool{}
	deferred := map[string]bool{}
	for _, line := range strings.Split(content, "\n") {
		isDeferred := tbdLineRE.MatchString(line)
		for _, m := range doctrineRuleRE.FindAllStringSubmatch(line, -1) {
			rules[m[1]] = true
			if isDeferred {
				deferred[m[1]] = true
			}
		}
	}
	if len(rules) == 0 {
		return Skip, i18n.T("gate.doctrine_realized.no_rules")
	}

	// Who realizes: the `realizes` edges ARRIVING here carry, in Method, the code of the
	// rule being realized.
	for _, e := range g.Neighbors(n.ID).In {
		if e.Type == mapx.EdgeRealizes {
			delete(rules, e.Method)
		}
	}
	if len(rules) == 0 {
		return Pass, ""
	}
	// A deferred rule leaves the accusation and enters DEBT: `@TBD` declares it is ahead
	// of the implementation on purpose. It does not become a Pass — it stays on the radar.
	var unrealized, owed []string
	for r := range rules {
		if deferred[r] {
			owed = append(owed, r)
		} else {
			unrealized = append(unrealized, r)
		}
	}
	if len(unrealized) > 0 {
		sort.Strings(unrealized)
		return Fail, fmt.Sprintf(i18n.T("gate.doctrine_realized.unrealized"), len(unrealized), strings.Join(unrealized, ", "))
	}
	sort.Strings(owed)
	return Pending, fmt.Sprintf(i18n.T("gate.doctrine_realized.deferred"), len(owed), strings.Join(owed, ", "))
}

// --- does the doctrine the SPEC cites exist? ---
//
// The vertical axis's `ref-resolves`. A reference that does not resolve is worse than an
// absent one: it looks like traceability, the shape gate goes green, and the map gains no
// edge at all — every gate on this axis then confronts a void, in silence.
//
// The characteristic failure is the rule RENAMED in the doctrine without its realizers
// being updated: the file exists, the rule does not, and nothing accuses.
func checkSpecDoctrineExists(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.spec_doctrine_exists.skip_not_spec")
	}
	if g == nil {
		return pendingNoMap()
	}

	// What the spec DECLARES, minus what is deferred on its own line.
	declared := map[string]bool{}
	for _, line := range strings.Split(content, "\n") {
		if tbdLineRE.MatchString(line) {
			continue
		}
		for _, m := range realizesTagRE.FindAllStringSubmatch(line, -1) {
			declared[m[1]] = true
		}
	}
	if len(declared) == 0 {
		return Skip, i18n.T("gate.spec_doctrine_exists.skip_none")
	}

	// What the map RESOLVED: the build only creates the edge when the doctrine exists, so
	// whatever was declared and did not become an edge is what fails to resolve.
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type == mapx.EdgeRealizes {
			delete(declared, e.Method)
		}
	}
	// One edge exists per doctrine FILE, and it confirms the file was found. Whether the
	// rule inside it exists is a separate question, answered by reading the file.
	if len(declared) > 0 {
		live := map[string]bool{}
		for _, e := range g.Neighbors(n.ID).Out {
			if e.Type != mapx.EdgeRealizes {
				continue
			}
			b, err := os.ReadFile(filepath.Join(root, e.To))
			if err != nil {
				continue
			}
			for _, m := range doctrineRuleRE.FindAllStringSubmatch(string(b), -1) {
				live[m[1]] = true
			}
		}
		for r := range declared {
			if live[r] {
				delete(declared, r)
			}
		}
	}
	if len(declared) == 0 {
		return Pass, ""
	}
	broken := make([]string, 0, len(declared))
	for r := range declared {
		broken = append(broken, r)
	}
	sort.Strings(broken)
	return Fail, fmt.Sprintf(i18n.T("gate.spec_doctrine_exists.unknown"), len(broken), strings.Join(broken, ", "))
}

// realizesTagRE — the `@realizes` tag in a spec. Deliberately duplicated from `scan`:
// this gate reads the CONTENT it is handed rather than the already-scanned node, because
// it needs to know which LINE the tag sits on in order to pair it with the waiver.
var realizesTagRE = regexp.MustCompile("@realizes\\s+`?([A-Z0-9]{3,6}-[A-Z]{1,2}[0-9]{2})`?")

// minCorpusForIDF is the floor below which the similarity ruler measures nothing.
//
// Measured against real sentences: a 2-text corpus scores 0.00 for a near-identical pair
// (every shared word weighs zero, and the shared words ARE the evidence); with 4 the same
// pair scores 0.51 and is correctly classified.
const minCorpusForIDF = 4

// minCopyScore is the floor a SIMILAR pair must clear to be called a copy.
//
// It matches the similarity library's own threshold, and it is restated here because the
// verdict alone is not enough: `Classify` promotes a low-scoring pair to "similar" when
// the two share a rare token, which on this axis is the norm rather than evidence — a
// spec that realizes a rule is expected to speak its vocabulary.
const minCopyScore = 0.5

// --- is the doctrine DUPLICATED in the spec? ---
//
// The defect this whole axis exists to eliminate. Before `product/` there were only two
// ways out for a rule spanning three screens: copy it into all three (and watch them
// diverge at the first change), or pick an arbitrary owner. The copy is the common one,
// because it reads well — each spec is complete on its own.
//
// And it is invisible to every other gate: both texts are well-formed, both catalogue
// their rules, both have a complete triad. Nothing compares one against the other.
//
// The ruler is SIMILARITY, not equality: whoever copies almost always adjusts a word.
// Exact comparison would catch only the laziest case and report green on the rest.
func checkDoctrineNotDuplicated(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.doctrine_not_duplicated.skip_not_spec")
	}
	if g == nil {
		return pendingNoMap()
	}

	// The doctrine text of each rule this spec realizes, read from the far end.
	doctrineText := map[string]string{}
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeRealizes {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, e.To))
		if err != nil {
			continue
		}
		for code, text := range ruleTexts(string(b)) {
			doctrineText[code] = text
		}
	}
	if len(doctrineText) == 0 {
		return Skip, i18n.T("gate.doctrine_not_duplicated.skip_none")
	}

	local := ruleTexts(content)

	// The corpus is EVERY rule of both sides — never just the pair being compared.
	//
	// `Weights` is IDF: a token appearing in every document of the corpus weighs zero,
	// because it separates nothing. With a two-text corpus that is exactly what happens
	// to the shared words, which are the evidence of a copy — measured, two sentences
	// differing by one word scored 0.00, and only a byte-identical copy was caught.
	//
	// With the full set of rules the shared vocabulary of the domain ("the", "limit")
	// stays cheap while what only these two say stays expensive, which is the signal the
	// gate is after.
	var corpus []string
	for _, t := range doctrineText {
		corpus = append(corpus, t)
	}
	for _, t := range local {
		corpus = append(corpus, t)
	}
	// FEW RULES: IDF has nothing to weigh, and the gate cannot measure.
	//
	// Measured: with a 2-rule corpus two sentences differing by one word score 0.00; from
	// 4 rules on, the same pair scores 0.51 and is correctly called similar. Below that
	// floor only a byte-identical copy would be caught, and reporting Pass would state
	// something that was never checked — the worst failure of a measuring instrument.
	if len(corpus) < minCorpusForIDF {
		return Pending, fmt.Sprintf(i18n.T("gate.doctrine_not_duplicated.pending_small_corpus"), len(corpus), minCorpusForIDF)
	}
	weights := similarity.Weights(corpus)

	var found []string
	for _, r := range parseRealizesWithLines(content) {
		// `@TBD` on the line is DEBT, not a waiver: the wording is still being worked
		// out, and charging it now would push whoever is writing to paraphrase for the
		// gate instead of for the reader.
		if r.deferred {
			continue
		}
		mine, ok := local[r.from]
		if !ok || mine == "" {
			continue
		}
		theirs, ok := doctrineText[r.to]
		if !ok || theirs == "" {
			continue
		}
		verdict, score := similarity.Classify(mine, theirs, weights)
		// IDENTICO passa direto; SIMILAR ainda precisa do SCORE, e a exigencia extra
		// existe por um caso medido.
		//
		// `Classify` promove a "similar" um par que compartilha um token RARO mesmo com
		// score baixo — evidencia estrutural de que dois textos tratam do mesmo assunto.
		// E' certo na origem da lib e ERRADO aqui: a spec DEVE usar o vocabulario da
		// doutrina que realiza, e' o proposito do `@realizes`. Sem o piso, as tres
		// primeiras arestas reais deste repositorio foram acusadas com 13%, 17% e 8% —
		// textos que nao se parecem em nada, unidos por compartilhar "dispensa" e "gate".
		if verdict == similarity.Identico ||
			(verdict == similarity.Similar && score >= minCopyScore) {
			found = append(found, fmt.Sprintf(i18n.T("gate.doctrine_not_duplicated.item"), r.from, r.to, score*100))
		}
	}
	if len(found) == 0 {
		return Pass, ""
	}
	sort.Strings(found)
	return Fail, fmt.Sprintf(i18n.T("gate.doctrine_not_duplicated.copied"), len(found), strings.Join(found, ", "))
}

// ruleTexts maps each catalogued rule code to the text that describes it — the rest of
// the line the code opens.
//
// The line is enough because all three catalogued forms put the description right after
// the code: the heading, the table row, the bold bullet. Reading the paragraph below a
// heading would drag in prose belonging to no rule.
func ruleTexts(content string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		m := doctrineRuleRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		rest := line[strings.Index(line, m[1])+len(m[1]):]
		// Strip what is punctuation of the FORM and not of the text: the em dash of a
		// heading, the pipes of a table row, the closing backtick.
		rest = strings.Trim(rest, "`|— -")
		if i := strings.Index(rest, "@realizes"); i >= 0 {
			rest = rest[:i]
		}
		if s := strings.TrimSpace(rest); s != "" {
			out[m[1]] = s
		}
	}
	return out
}

// realizesOnLine is one `@realizes` declaration with the context the gate needs: which
// local rule made it, and whether the line carries debt.
type realizesOnLine struct {
	from     string
	to       string
	deferred bool
}

func parseRealizesWithLines(content string) []realizesOnLine {
	var out []realizesOnLine
	current := ""
	for _, line := range strings.Split(content, "\n") {
		if m := doctrineRuleRE.FindStringSubmatch(line); m != nil {
			current = m[1]
		} else if strings.TrimSpace(line) == "" {
			current = ""
		}
		deferred := tbdLineRE.MatchString(line)
		for _, m := range realizesTagRE.FindAllStringSubmatch(line, -1) {
			out = append(out, realizesOnLine{from: current, to: m[1], deferred: deferred})
		}
	}
	return out
}

// --- does this layer DEMAND doctrine? ---
//
// The rule a layer declares with `requires_doctrine: true`: every catalogued rule of its
// specs must say which product decision it concretises.
//
// Optional by default, and the reason is measured. In this repository, of 841 catalogued
// rules the overwhelming majority is LOCAL to its unit — `SBGRD-B01` ("an artifact that
// is not code leaves without a verdict") belongs to no product; it is a gate's mechanics.
// Demanding it everywhere would force inventing umbrella doctrine just to silence the
// gate, which is the vice `placeholder-filled` exists to catch.
//
// In a PRODUCT application the proportion inverts: almost every screen rule serves a
// product decision, and the one that does not is suspect. Anchors does not know which
// case it is looking at — the Structure does, layer by layer.
func checkSpecRealizesDoctrine(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.spec_realizes_doctrine.skip_not_spec")
	}
	if cfg == nil {
		return Skip, i18n.T("gate.spec_realizes_doctrine.skip_no_config")
	}
	if !layerRequiresDoctrine(n, root, g, cfg) {
		return Skip, i18n.T("gate.spec_realizes_doctrine.skip_layer_does_not_require")
	}

	declared := map[string]bool{}
	deferred := map[string]bool{}
	for _, r := range parseRealizesWithLines(content) {
		if r.from == "" {
			continue // an orphan tag declares nothing for any rule
		}
		declared[r.from] = true
	}
	var naked []string
	for _, line := range strings.Split(content, "\n") {
		m := doctrineRuleRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if tbdLineRE.MatchString(line) {
			deferred[m[1]] = true
			continue
		}
		if !declared[m[1]] && !strings.Contains(line, "@realizes") {
			naked = append(naked, m[1])
		}
	}
	if len(naked) > 0 {
		sort.Strings(naked)
		return Fail, fmt.Sprintf(i18n.T("gate.spec_realizes_doctrine.missing"), len(naked), strings.Join(naked, ", "))
	}
	if len(deferred) > 0 {
		return Pending, fmt.Sprintf(i18n.T("gate.spec_realizes_doctrine.deferred"), len(deferred))
	}
	return Pass, ""
}

// layerRequiresDoctrine answers whether the unit's layer demands `@realizes`.
//
// It resolves the layer by the TARGET the spec describes, never by the spec's own file:
// a `.spec.md` matches the `spec` layer, and the demand is declared on the layer of the
// thing being specified (`screen`, `handler`).
//
// Two routes, and the second is not redundancy: the `specifies` edge only exists AFTER
// the code is born, and at the `spec` step it is not. Without resolving by PATH the gate
// would be blind exactly in the window where the author is writing the spec — which is
// when the demand matters most. The lesson is the same one `triad-complete` already paid
// for, in the comment above its own second route.
func layerRequiresDoctrine(n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) bool {
	tags := append([]string{}, n.Tags...)
	if g != nil {
		for _, e := range g.Neighbors(n.ID).Out {
			if e.Type != mapx.EdgeSpecifies {
				continue
			}
			for _, target := range g.Nodes {
				if target.ID == e.To {
					tags = append(tags, target.Tags...)
				}
			}
		}
	}
	for _, t := range tags {
		if l, ok := cfg.Layers[t]; ok && l.RequiresDoctrine {
			return true
		}
	}
	if base := strings.TrimSuffix(n.ID, ".spec.md"); base != n.ID {
		for _, ext := range []string{".ts", ".tsx", ".go", ".py", ".js"} {
			layer, _ := scan.Classify(base+ext, cfg)
			if layer == "" {
				continue
			}
			if cfg.Layers[layer].RequiresDoctrine {
				return true
			}
		}
	}
	return false
}
