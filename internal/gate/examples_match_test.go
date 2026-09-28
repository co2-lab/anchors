package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const planOutline = `Feature: plans

  @PLANX-B01 @unit-level
  Scenario Outline: The plan names its kind
    Given the plan <plan>
    Then it reads <kind>

    Examples:
      | plan      | kind       |
      | allowance | "monthly"  |
      | savings   | yearly     |
      | goal      | once       |
`

// runExamples writes the feature and the test files and confronts the feature.
func runExamples(t *testing.T, cfg *config.Config, feature string, tests map[string]string) (Verdict, string) {
	t.Helper()
	resetProjectTestsCache()
	root := t.TempDir()
	writeFile(t, root, "plan.feature", feature)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "plan.feature", Kind: mapx.KindFeature}}}
	for f, body := range tests {
		writeFile(t, root, f, body)
		g.Nodes = append(g.Nodes, mapx.Node{ID: f, Kind: mapx.KindTest})
	}
	return checkExamplesMatch(feature, g.Nodes[0], root, g, cfg)
}

func examplesCfg(labels bool, tests *config.TestsSource) *config.Config {
	return &config.Config{
		Dialect: &config.Dialect{Family: "ts", Tests: tests},
		Gates:   []config.Gate{{Name: "e", Check: "examples-match", On: []string{"feature"}, Labels: labels}},
	}
}

func TestExamplesMatch_aRowTheTestsDoNotRun(t *testing.T) {
	t.Run("EXMCH-B01: A row the tests do not run fails", func(t *testing.T) {})
	v, msg := runExamples(t, examplesCfg(false, nil), planOutline, map[string]string{
		"a.test.ts": "it.each([['allowance', 'monthly']])('PLANX-B01: kind', (p, k) => {\n  expect(kind(p)).toBe(k)\n})\n",
		"b.test.ts": "it('PLANX-B01: savings', () => {\n  expect(kind('savings')).toBe('yearly')\n})\nit('other', () => { goal() })\n",
	})
	if v != Fail || !strings.Contains(msg, "1 example row(s)") || !strings.Contains(msg, "PLANX-B01 line 12: `goal`, `once`") {
		t.Errorf("the third row fails by code, line and values, got %v: %s", v, msg)
	}
	all := "it.each([['allowance','monthly'],['savings','yearly'],['goal','once']])('PLANX-B01: kind', () => {\n  run()\n})\n"
	if v, msg := runExamples(t, examplesCfg(false, nil), planOutline, map[string]string{"a.test.ts": all}); v != Pass {
		t.Errorf("every row carried passes, got %v: %s", v, msg)
	}
}

func TestExamplesMatch_wholeTokens(t *testing.T) {
	t.Run("EXMCH-B02: A value is a whole token with its case", func(t *testing.T) {})
	for _, c := range []struct {
		text, v string
		want    bool
	}{
		{"x('Crédito')", "Crédito", true},
		{"x('Cartão de crédito')", "Crédito", false},
		{"x('CreditoCard')", "Credito", false},
		{"x('simplified')", "simplificado", false},
		{"n = 100", "10", false},
		{"n = 10", "10", true},
		{"key_goal", "goal", false},
		{"goal", "goal", true},
	} {
		if got := tokenIn(c.text, c.v); got != c.want {
			t.Errorf("tokenIn(%q, %q) = %v, want %v", c.text, c.v, got, c.want)
		}
	}
	rows := exampleTables("@PLANX-B01\nScenario Outline: x\nExamples:\n| a | b | c |\n| \"x\" | `y` |  |\n| | | |\n")
	if len(rows) != 1 || len(rows[0].rows) != 1 || strings.Join(rows[0].rows[0].values, ",") != "x,y" {
		t.Errorf("quotes dropped, empty cells and rows ask nothing, got %+v", rows)
	}
}

func TestExamplesMatch_labelColumns(t *testing.T) {
	t.Run("EXMCH-B03: A label column is display text", func(t *testing.T) {})
	feature := "@PLANX-B01\nScenario Outline: x\nExamples:\n| plan (label) | key | tipo (rótulo) | kind |\n| Mesada | allowance | Mensal | monthly |\n"
	test := "it('PLANX-B01: x', () => { k('allowance', 'monthly') })\n"
	if v, msg := runExamples(t, examplesCfg(false, nil), feature, map[string]string{"a.test.ts": test}); v != Pass {
		t.Errorf("label columns are not looked for, got %v: %s", v, msg)
	}
	for _, h := range []string{"Plan (LABEL)", "x (etiqueta)", "y (rotulo)"} {
		if !isLabelColumn(h) {
			t.Errorf("%q marks a label", h)
		}
	}
	if isLabelColumn("label") || isLabelColumn("plan") {
		t.Error("a header without the marker is a key column")
	}
}

func TestExamplesMatch_tables(t *testing.T) {
	t.Run("EXMCH-B04: The rows are the coded outline's tables", func(t *testing.T) {})
	feature := `# language: pt
Funcionalidade: planos

  @PLANX-B02
  Esquema do Cenário: dois quadros
    Dado <a>

    Exemplos:
      | a |
      # comment
      | um |

    Exemplos:
      | a |
      | dois |
    Então fim
      | tres |

  Cenário: sem código
    Exemplos:
      | a |
      | quatro |
`
	got := exampleTables(feature)
	if len(got) != 1 || got[0].code != "PLANX-B02" || len(got[0].rows) != 2 ||
		got[0].rows[0].values[0] != "um" || got[0].rows[0].line != 11 || got[0].rows[1].values[0] != "dois" {
		t.Errorf("the coded outline has its two tables' rows, got %+v", got)
	}
	wide := exampleTables("@PLANX-B04\nScenario Outline: x\nExamples:\n| a (label) |\n| shown | extra |\n")
	if len(wide) != 1 || strings.Join(wide[0].rows[0].values, ",") != "extra" {
		t.Errorf("a cell past the header is looked for, got %+v", wide)
	}
	if got := exampleTables("@PLANX-B03\nScenario: plain\n  Given x\n"); len(got) != 0 {
		t.Errorf("a scenario with no rows is left out, got %+v", got)
	}
}

func TestExamplesMatch_labelledTests(t *testing.T) {
	t.Run("EXMCH-B05: The test's body is read as the assertion gate reads it", func(t *testing.T) {})
	feature := "@PLANX-B01\nScenario Outline: x\nExamples:\n| plan |\n| allowance |\n"
	test := "describe('plans', () => {\n  it('PLANX-B01: kinds', () => {})\n  expect(kind('allowance')).toBe(1)\n})\nit('far', () => { kind('savings') })\n"
	if v, msg := runExamples(t, examplesCfg(true, nil), feature, map[string]string{"a.test.ts": test}); v != Pass {
		t.Errorf("with labels the block around carries the row, got %v: %s", v, msg)
	}
	if v, _ := runExamples(t, examplesCfg(false, nil), feature, map[string]string{"a.test.ts": test}); v != Fail {
		t.Errorf("without labels the empty test carries nothing, got %v", v)
	}
}

func TestExamplesMatch_skips(t *testing.T) {
	t.Run("EXMCH-B06: Nothing to confront is skipped", func(t *testing.T) {})
	cfg := examplesCfg(false, nil)
	if v, _ := checkExamplesMatch(planOutline, mapx.Node{Kind: mapx.KindSpec}, "", nil, cfg); v != Skip {
		t.Errorf("a spec is skipped, got %v", v)
	}
	if v, _ := runExamples(t, cfg, "@PLANX-B01\nScenario: x\n  Given y\n", nil); v != Skip {
		t.Errorf("no examples is skipped, got %v", v)
	}
	if v, _ := runExamples(t, &config.Config{Dialect: &config.Dialect{Family: "rust"}}, planOutline, nil); v != Skip {
		t.Errorf("no tests source is skipped, got %v", v)
	}
	if v, _ := runExamples(t, cfg, planOutline, map[string]string{"a.test.ts": "it('other', () => {})\n"}); v != Skip {
		t.Errorf("no test citing the code is skipped, got %v", v)
	}
}

func TestExamplesMatch_unlistableTests(t *testing.T) {
	t.Run("EXMCH-E01: Tests that cannot be listed fail the gate with the reason", func(t *testing.T) {})
	cfg := examplesCfg(false, &config.TestsSource{Script: "exit 4"})
	if v, msg := runExamples(t, cfg, planOutline, nil); v != Fail || !strings.Contains(msg, "exit 4") {
		t.Errorf("a failing script fails naming why, got %v: %s", v, msg)
	}
}
