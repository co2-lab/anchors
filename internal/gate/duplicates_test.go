// @anchors
//   code: DPTSA
//   ref: GTDPG

package gate

import (
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
