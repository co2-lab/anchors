// @anchors
//   code: RUTRL
//   ref: RLUSG

package gate

import (
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

var usesSpecNode = mapx.Node{Kind: mapx.KindSpec, Layer: "screen"}

func TestRuleUsesOf_anyLanguageByPosition(t *testing.T) {
	t.Run("RLUSG-B01: The sections are read in any language and by position", func(t *testing.T) {})
	spec := "## Validações\n\n| Regra | Campo | Condição | Comportamento |\n| --- | --- | --- | --- |\n" +
		"| `PAYMT-V01` | `amount` | > 0 | 400 |\n\n" +
		"## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `PAYMT-B02` | `balance`, `DEP1` |\n| `PAYMT-B03` | TODO: fill |\n\n" +
		"## Cosmetics\n\n| `PAYMT-P01` | `status` | pending | dim |\n\n" +
		"## Rules\n\n| `PAYMT-B09` | not a use row |\n"
	cfg := &config.Config{SectionTitles: map[string]string{"presentation-validations": "Cosmetics"}}
	got := ruleUsesOf(spec, cfg, "screen")
	want := []ruleUse{{"PAYMT-V01", []string{"amount"}}, {"PAYMT-B02", []string{"balance", "DEP1"}}, {"PAYMT-P01", []string{"status"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestUsedItems(t *testing.T) {
	t.Run("RLUSG-B02: A uses cell lists backticked or comma-separated items", func(t *testing.T) {})
	if got := usedItems("reads `summary.balance` and `PAYMT-S01`"); !reflect.DeepEqual(got, []string{"summary.balance", "PAYMT-S01"}) {
		t.Errorf("backticked items, got %v", got)
	}
	if got := usedItems("amount, currency ,DEP2"); !reflect.DeepEqual(got, []string{"amount", "currency", "DEP2"}) {
		t.Errorf("comma-separated items, got %v", got)
	}
}

const declaredSpec = "## Rules\n\n### PAYMT-B01 — pays\n\n### PAYMT-B02 — refunds  @no-uses: pure orchestration\n\n" +
	"### PAYMT-E01 — declined\n\n### PAYMT-Q01 — which gateway?\n\n" +
	"## Validations\n\n| `PAYMT-V01` | `amount` | > 0 | 400 |\n"

func TestRuleUsesDeclared_failsTheSilentRules(t *testing.T) {
	t.Run("RLUSG-B03: A rule that does not say what it uses fails", func(t *testing.T) {})
	v, msg := checkRuleUsesDeclared(declaredSpec, usesSpecNode, "", nil, &config.Config{})
	if v != Fail || !strings.Contains(msg, "PAYMT-B01, PAYMT-E01") || strings.Contains(msg, "B02") || strings.Contains(msg, "Q01") || strings.Contains(msg, "V01") {
		t.Fatalf("the behaviour and the error, not the waived, the question nor the validation, got %v %s", v, msg)
	}
	cfg := &config.Config{Gates: []config.Gate{{Name: "rule-uses-declared", Check: "rule-uses-declared", Letters: []string{"e"}}}}
	if v, msg := checkRuleUsesDeclared(declaredSpec, usesSpecNode, "", nil, cfg); v != Fail || strings.Contains(msg, "B01") || !strings.Contains(msg, "PAYMT-E01") {
		t.Errorf("with letters [E] only the error is asked about, got %v %s", v, msg)
	}
	ok := declaredSpec + "\n## Rule uses\n\n| `PAYMT-B01` | `amount` |\n| `PAYMT-E01` | `DEP1` |\n"
	if v, msg := checkRuleUsesDeclared(ok, usesSpecNode, "", nil, &config.Config{}); v != Pass {
		t.Errorf("every rule declared passes, got %v %s", v, msg)
	}
}

func TestRuleUsesDeclared_skips(t *testing.T) {
	t.Run("RLUSG-B04: Nothing to ask about is a skip", func(t *testing.T) {})
	if v, _ := checkRuleUsesDeclared(declaredSpec, mapx.Node{Kind: mapx.KindFeature}, "", nil, nil); v != Skip {
		t.Errorf("a feature is skipped, got %v", v)
	}
	if v, _ := checkRuleUsesDeclared("## Open\n\n### PAYMT-Q01 — which?\n", usesSpecNode, "", nil, nil); v != Skip {
		t.Errorf("only a question to ask about is skipped, got %v", v)
	}
}

func TestRuleUsesResolve(t *testing.T) {
	t.Run("RLUSG-B05: What a rule uses must exist in the spec", func(t *testing.T) {})
	spec := "## Data contract\n\n| Field | Origin |\n| --- | --- |\n| `amount` | DEP1 |\n| summary | DEP1 |\n\n" +
		"## Data states\n\n### `status`\n\n## Dependencies\n\n| Code | File |\n| --- | --- |\n| DEP1 | `api.ts` |\n\n" +
		"## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n" +
		"| `PAYMT-B01` | `amount`, `summary.balance`, `status`, `DEP1`, `OTHER-S01` |\n" +
		"| `PAYMT-B02` | `ghost`, `DEP9` |\n"
	v, msg := checkRuleUsesResolve(spec, usesSpecNode, "", nil, nil)
	if v != Fail || !strings.Contains(msg, "PAYMT-B02 → `ghost`") || !strings.Contains(msg, "PAYMT-B02 → `DEP9`") || strings.Contains(msg, "B01") {
		t.Fatalf("only the undeclared field and dependency, with their rule, got %v %s", v, msg)
	}
	fixed := strings.Replace(spec, "| `PAYMT-B02` | `ghost`, `DEP9` |\n", "", 1)
	if v, msg := checkRuleUsesResolve(fixed, usesSpecNode, "", nil, nil); v != Pass {
		t.Errorf("everything declared passes, got %v %s", v, msg)
	}
}

func TestRuleUsesResolve_skips(t *testing.T) {
	t.Run("RLUSG-B06: A spec whose rules say nothing yet is a skip", func(t *testing.T) {})
	if v, _ := checkRuleUsesResolve(declaredSpec, mapx.Node{Kind: mapx.KindFeature}, "", nil, nil); v != Skip {
		t.Errorf("a feature is skipped, got %v", v)
	}
	if v, _ := checkRuleUsesResolve("## Rules\n\n### PAYMT-B01 — pays\n", usesSpecNode, "", nil, nil); v != Skip {
		t.Errorf("no rule uses is skipped, got %v", v)
	}
}

func TestRuleUsesImplemented(t *testing.T) {
	t.Run("RLUSG-B07: The fields a rule uses appear in the code the spec governs", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "src/pay.ts", "// ghost is only in a comment\nexport const pay = (amount: number, s) => s.balance > amount\n")
	g := &mapx.Graph{Edges: []mapx.Edge{{From: "src/pay.spec.md", To: "src/pay.ts", Type: mapx.EdgeSpecifies}}}
	spec := "## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `PAYMT-B01` | `amount`, `summary.balance`, `DEP1`, `PAYMT-S01` |\n| `PAYMT-B02` | `ghost` |\n| `PAYMT-B03` | `amounts` |\n| `PAYMT-B04` | `bal` |\n"
	n := mapx.Node{ID: "src/pay.spec.md", Kind: mapx.KindSpec}
	v, msg := checkRuleUsesImplemented(spec, n, root, g, nil)
	if v != Fail || !strings.Contains(msg, "PAYMT-B02 → `ghost`") || !strings.Contains(msg, "PAYMT-B03 → `amounts`") || !strings.Contains(msg, "PAYMT-B04 → `bal`") || strings.Contains(msg, "B01") {
		t.Fatalf("only the fields no code mentions as a word, got %v: %s", v, msg)
	}
	if v, _ := checkRuleUsesImplemented(spec, n, root, &mapx.Graph{}, nil); v != Skip {
		t.Errorf("a spec governing no code is skipped, got %v", v)
	}
	if v, _ := checkRuleUsesImplemented("## Rules\n", n, root, g, nil); v != Skip {
		t.Errorf("a spec with no rule uses is skipped, got %v", v)
	}
}
