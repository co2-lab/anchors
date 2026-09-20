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
