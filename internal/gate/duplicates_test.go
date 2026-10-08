// @anchors
//   code: DPTSA
//   ref: GTDPG

package gate

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const dupSpec = `# Pay

## Rules

### PAYMT-B01 — charges the amount

| ` + "`PAYMT-B01`" + ` | the amount | charged once |

### Loading

| Rule | Strategy |
| --- | --- |
| ` + "`PAYMT-B01`" + ` | all at once |
| ` + "`PAYMT-B02`" + ` | paged |

## Rule uses

| Rule | Uses |
| --- | --- |
| ` + "`PAYMT-B01`" + ` | ` + "`amount`" + ` |

## Navigation

### Out

| Rule | Destination |
| --- | --- |
| ` + "`PAYMT-B02`" + ` | Receipt |

## State Flow

| From | Trigger | To |
| --- | --- | --- |
| ` + "`PAYMT-B02`" + ` | paid | done |

## Eventos / Callbacks

| Rule | Callback | Trigger |
| --- | --- | --- |
| ` + "`PAYMT-B01`" + ` | ` + "`onPay`" + ` | tap |

## Notes

- ` + "`PAYMT-B02`" + ` registers the paging as it is today
- PAYMT-B09: @retired: r3 — merged into B01

## Open Decisions

| ` + "`PAYMT-B02`" + ` | page size? |
`

func TestDuplicates_ruleTypesCountsDefinitions(t *testing.T) {
	t.Run("GTDPG-B02: rule-types counts the rule codes a file defines, not those it cites", func(t *testing.T) {})
	// The project names its events section "Eventos / Callbacks", in the component layer.
	cfg := &config.Config{Layers: map[string]config.Layer{"component": {SectionTitles: config.SectionTitles{"callbacks": "Eventos / Callbacks"}}}}
	occ := definedRuleOccurrences(dupSpec, mapx.Node{}, "", nil, cfg)
	var got []string
	for _, o := range occ {
		got = append(got, o.Key+"@"+strconv.Itoa(o.Line))
	}
	if strings.Join(got, " ") != "PAYMT-B01@5 PAYMT-B01@13 PAYMT-B02@14" {
		t.Errorf("definitions only, a heading with its own row once: %v", got)
	}
	plain := definedRuleOccurrences(dupSpec, mapx.Node{}, "", nil, nil)
	if len(plain) != 4 {
		t.Errorf("without the project's title, the events table is not known to cite: %d", len(plain))
	}
}

func TestDuplicates_theVerdict(t *testing.T) {
	t.Run("GTDPG-B01: A key declared twice turns the gate's verdict into a failure naming the lines, unless switched off or skipped", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "pay.spec.md", dupSpec)
	n := mapx.Node{ID: "pay.spec.md", Kind: mapx.KindSpec}
	g := config.Gate{Name: "rule-types", Check: "rule-types"}
	if !HasDuplicateReader("rule-types") || HasDuplicateReader("non-empty") {
		t.Fatal("rule-types has a reader; non-empty has none")
	}
	v, msg := confrontDuplicates(g, n, root, nil, nil, Pass, "")
	if v != Fail || !strings.Contains(msg, "PAYMT-B01") || !strings.Contains(msg, "5, 13") || strings.Contains(msg, "PAYMT-B02") {
		t.Errorf("the repeated code fails, naming its lines: %v %s", v, msg)
	}
	if v, msg := confrontDuplicates(g, n, root, nil, nil, Fail, "a letter outside the vocabulary"); v != Fail || !strings.HasPrefix(msg, "a letter outside the vocabulary; ") {
		t.Errorf("the gate's own finding is kept beside: %v %s", v, msg)
	}
	off := g
	off.Duplicates = config.Bool(false)
	if v, msg := confrontDuplicates(off, n, root, nil, nil, Pass, ""); v != Pass || msg != "" {
		t.Errorf("duplicates: false switches it off: %v %s", v, msg)
	}
	if v, _ := confrontDuplicates(g, n, root, nil, nil, Skip, "not a spec"); v != Skip {
		t.Error("a skipped node is not counted")
	}
	occurrenceReaders["test-diverge"] = occurrenceReader{read: definedRuleOccurrences, verdict: Diverge}
	defer delete(occurrenceReaders, "test-diverge")
	dv := config.Gate{Name: "d", Check: "test-diverge"}
	if v, _ := confrontDuplicates(dv, n, root, nil, nil, Pass, ""); v != Diverge {
		t.Errorf("a gate that measures repeats as a divergence diverges: %v", v)
	}
	if v, _ := confrontDuplicates(dv, n, root, nil, nil, Fail, "x"); v != Fail {
		t.Errorf("a failure stays a failure: %v", v)
	}
}

const catalogueSpec = "# Pay\n\n> **Code**: `PAYMT`\n\n" +
	"## Domain\n\n| Input | Accepts | Outside | Who guarantees |\n| --- | --- | --- | --- |\n" +
	"| `amount` | > 0 | 0 | this unit |\n| `Amount` | cents | — | this unit |\n| `currency` | ISO | — | the caller |\n\n" +
	"## Environment Variables\n\n| Variable | Type | Required |\n| --- | --- | --- |\n" +
	"| PAY_URL | string | yes |\n| PAY_KEY | string | yes |\n| `PAY_URL` | string | no |\n\n" +
	"## Output Contract\n\n| Status | When |\n| --- | --- |\n| 200 | paid |\n| 4xx | rejected |\n| 422 | invalid |\n| 4XX | other |\n| 422 | too large |\n\n" +
	"## Dependencies\n\n| Dep | File | Uses |\n| --- | --- | --- |\n| DEP1 | `a.ts` | x |\n| DEP2 | `b.ts` | y |\n| `DEP1` | `c.ts` | z |\n\n" +
	"## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `PAYMT-B01` | DEP1, DEP2 |\n\n" +
	"## Revisions\n\nPAYMT-R0001: amount in cents\nPAYMT-R0002: currency added\n**PAYMT-R0001** — renamed\n\n" +
	"## Rules\n\n### Validations\n\n| `PAYMT-V01` | x |\n\n### Validations\n\n| `PAYMT-V02` | y |\n\n" +
	"## Open Decisions\n\n| Code | Question |\n| --- | --- |\n| `PAYMT-Q01` | page size? |\n| ~~`PAYMT-Q02`~~ | answered |\n| `PAYMT-Q02` | reused |\n"

func keysOf(occ []Occurrence) []string {
	var out []string
	for k := range duplicatesOf(occ) {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func TestDuplicates_theSpecCatalogue(t *testing.T) {
	t.Run("GTDPG-B03: The spec catalogue's gates count their declarations: an environment variable, a Domain entry, an open question, a revision, a section under one parent", func(t *testing.T) {})
	n := mapx.Node{ID: "pay.spec.md", Kind: mapx.KindSpec}
	check := func(name string, occ []Occurrence, want string) {
		t.Helper()
		if got := strings.Join(keysOf(occ), ","); got != want {
			t.Errorf("%s: repeated %q, want %q (%+v)", name, got, want, occ)
		}
	}
	check("env", envOccurrences(catalogueSpec, n, "", nil, nil), "PAY_URL")
	check("domain", domainOccurrences(catalogueSpec, n, "", nil, nil), "amount")
	check("open", openQuestionOccurrences(catalogueSpec, n, "", nil, nil), "PAYMT-Q02")
	check("revisions", revisionOccurrences(catalogueSpec, n, "", nil, nil), "PAYMT-R0001")
	check("sections", sectionOccurrences(catalogueSpec, n, "", nil, nil), "rules › validations")
	// The lines name where each is declared.
	for _, o := range duplicatesOf(envOccurrences(catalogueSpec, n, "", nil, nil))["PAY_URL"] {
		if l := strings.Split(catalogueSpec, "\n")[o.Line-1]; !strings.Contains(l, "PAY_URL") {
			t.Errorf("line %d is not the variable's row: %q", o.Line, l)
		}
	}
}

func TestDuplicates_flagsCodeAndFeatures(t *testing.T) {
	t.Run("GTDPG-B04: A flag scenario code, a symbol under two used-by flags, a testID in two rows of the inventory and an identical Examples row are counted; a testID in prose and a used-by on two symbols are not", func(t *testing.T) {})
	n := mapx.Node{ID: "x"}
	keys := func(occ []Occurrence) string { return strings.Join(keysOf(occ), ",") }
	usedBy := "// @used-by: AAAAA\n// @used-by: BBBBB\nexport const PALETTE = {}\n\n// @used-by: AAAAA\nexport const OTHER = 1\n"
	if got := keys(usedByOccurrences(usedBy, n, "", nil, nil)); got != "PALETTE" {
		t.Errorf("two used-by flags on one symbol: %q", got)
	}
	cfg := &config.Config{Derived: &config.Derived{TestHandle: "testID"}}
	ids := "# S\n\n## Test IDs\n\n| testID | Element |\n| --- | --- |\n| `pay-button` | the button |\n| `pay-total` | `pay-button` shown above |\n| `pay-button` | again |\n"
	if got := keys(testIDOccurrences(ids, n, "", nil, cfg)); got != "pay-button" {
		t.Errorf("a testID in two rows, not one cited in another cell: %q", got)
	}
	flag := "| `CHKUT-G01` | `= \"off\"` | old |\n| `CHKUT-G02` | `= \"on\"` | new |\n| `CHKUT-G01` | absent | old |\n"
	if got := keys(flagScenarioOccurrences(flag, mapx.Node{ID: "checkout.flag.md"}, "", nil, nil)); got != "CHKUT-G01" {
		t.Errorf("a flag scenario code twice: %q", got)
	}
	feature := "Feature: x\n\n  @PAYMT-B01\n  Scenario Outline: pays\n    Examples:\n      | amount | result |\n      | 10 | ok |\n      | 20 | ok |\n      | 10 | ok |\n"
	if got := keys(exampleRowOccurrences(feature, n, "", nil, nil)); !strings.HasPrefix(got, "PAYMT-B01 (Examples at line 6) | 10") || strings.Contains(got, "20") {
		t.Errorf("an identical Examples row: %q", got)
	}
	// Two outlines of one rule, each with its own table and columns, share a value: no repeat.
	two := "Feature: x\n\n  @PAYMT-S02#03\n  Scenario Outline: small\n    Examples:\n      | small |\n      | true |\n\n" +
		"  @PAYMT-S02#04\n  Scenario Outline: inverted\n    Examples:\n      | on |\n      | true |\n"
	if got := keys(exampleRowOccurrences(two, n, "", nil, nil)); got != "" {
		t.Errorf("rows of two tables are no repeat: %q", got)
	}
}
