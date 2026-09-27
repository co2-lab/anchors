package migra

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func stepTo(t *testing.T, n int) *Step {
	t.Helper()
	for i := range steps {
		if steps[i].To == n {
			return &steps[i]
		}
	}
	t.Fatalf("the step that produces format %d is not registered", n)
	return nil
}

// The 1→2 STEP is really registered — not only does the mechanism work, the content exists.
// Without this assertion, removing `format_2.go` would leave the suite green and the product
// with no migration at all.
func TestFormat2Step_isRegisteredAndRenamesTheKeys(t *testing.T) {
	t.Run("MGSTM-B01: Format 2 renames the Portuguese keys in their own files", func(t *testing.T) {})
	p := stepTo(t, 2)
	graph := map[string]string{"gerado_por": "generated_by", "code_declarado": "code_declared", "julgamentos": "judgments"}
	for old, want := range graph {
		if got := p.RenameKeys["anchors.graph.yaml"][old]; got != want {
			t.Errorf("map key %q → %q, want %q", old, got, want)
		}
	}
	if got := p.RenameKeys["anchors.yaml"]["trinca_opcional"]; got != "triad_optional" {
		t.Errorf("configuration key trinca_opcional → %q", got)
	}
}

func TestFormat2Step_renamesTheGateNamesWhereverStored(t *testing.T) {
	t.Run("MGSTM-B02: Format 2 renames the Portuguese gate names wherever they are stored", func(t *testing.T) {})
	dir := t.TempDir()
	graph := escreve(t, dir, "anchors.graph.yaml", `version: 1
edges:
    - from: a
      judgments:
        - gate: regra-cumprida
`)
	cfg := escreve(t, dir, "anchors.yaml", `version: 1
gates:
    - name: regra-cumprida
    - id: trinca-completa
    - check: sem-duplicacao
`)
	for _, p := range []string{graph, cfg} {
		if _, err := MigrateFile(p, 2, false); err != nil {
			t.Fatal(err)
		}
	}
	if g := string(mustRead(t, graph)); !strings.Contains(g, "gate: rule-fulfilled") {
		t.Errorf("the map's judgment gate was not renamed:\n%s", g)
	}
	c := string(mustRead(t, cfg))
	for _, want := range []string{"name: rule-fulfilled", "id: triad-complete", "check: no-duplication"} {
		if !strings.Contains(c, want) {
			t.Errorf("expected %q in:\n%s", want, c)
		}
	}
}

func TestFormat3Step_renamesTheEightGateNames(t *testing.T) {
	t.Run("MGSTM-B03: Format 3 renames the eight remaining gate names in the configuration", func(t *testing.T) {})
	if n := len(stepTo(t, 3).RenameValues["anchors.yaml"]["name"]); n != 8 {
		t.Errorf("format 3 should rename eight gate names, renames %d", n)
	}
	p := escreve(t, t.TempDir(), "anchors.yaml", `version: 2
gates:
    - name: regra-implementada
      check: header-conforme
`)
	if _, err := MigrateFile(p, 3, false); err != nil {
		t.Fatal(err)
	}
	text := string(mustRead(t, p))
	if !strings.Contains(text, "name: rule-implemented") || !strings.Contains(text, "check: header-valid") {
		t.Errorf("format 3 did not rename the gate names:\n%s", text)
	}
}

func TestFormat4Step_renamesTheFourKeys(t *testing.T) {
	t.Run("MGSTM-B04: Format 4 renames the four keys that lied about what they hold", func(t *testing.T) {})
	p := escreve(t, t.TempDir(), "anchors.yaml", `version: 3
auto_judgment: true
rule_marking: required
layers:
    spec:
        triad_optional: [tested-by]
rule_types:
    R:
        requires_code: [rules]
`)
	if _, err := MigrateFile(p, 4, false); err != nil {
		t.Fatal(err)
	}
	text := string(mustRead(t, p))
	for _, want := range []string{"enable_auto_judgment:", "optional_triad_edges:", "sections_require_code:", "rule_marking_policy:"} {
		if !strings.Contains(text, want) {
			t.Errorf("expected %q in:\n%s", want, text)
		}
	}
	for _, old := range []string{"\nauto_judgment:", " triad_optional:", " requires_code:", "\nrule_marking:"} {
		if strings.Contains(text, old) {
			t.Errorf("the old key %q is still there:\n%s", old, text)
		}
	}
}

func TestSteps_theChainToTheCurrentFormatHasNoHole(t *testing.T) {
	t.Run("MGSTM-I01: The chain from format 1 to the current format has no hole", func(t *testing.T) {})
	got, err := StepsFrom(1, mapx.FormatoAtual)
	if err != nil {
		t.Fatalf("the chain to the current format has a hole: %v", err)
	}
	if len(got) != mapx.FormatoAtual-1 {
		t.Fatalf("expected %d steps, got %d", mapx.FormatoAtual-1, len(got))
	}
	for i, s := range got {
		if s.To != i+2 {
			t.Errorf("step %d produces format %d", i, s.To)
		}
		if s.Why == "" {
			t.Errorf("the step producing %d does not say what changed — it is what whoever "+
				"reviews the diff reads before opening it", s.To)
		}
	}
}

func TestSteps_aKeyRenamedTwiceEndsUnderItsLatestName(t *testing.T) {
	t.Run("MGSTM-I02: A key renamed by two formats ends under its latest name", func(t *testing.T) {})
	p := escreve(t, t.TempDir(), "anchors.yaml", `layers:
    spec:
        trinca_opcional: [tested-by]
`)
	if _, err := MigrateFile(p, mapx.FormatoAtual, false); err != nil {
		t.Fatal(err)
	}
	text := string(mustRead(t, p))
	if !strings.Contains(text, "optional_triad_edges: [tested-by]") ||
		strings.Contains(text, "trinca_opcional") || strings.Contains(text, "triad_optional:") {
		t.Errorf("expected the key under its latest name only:\n%s", text)
	}
}

func TestFormat3Step_leavesTheMapGateAndTheIdAlone(t *testing.T) {
	t.Run("MGSTM-X01: Format 3 leaves the map's gate and the configuration's id alone", func(t *testing.T) {})
	dir := t.TempDir()
	graph := escreve(t, dir, "anchors.graph.yaml", "version: 2\nedges:\n    - judgments:\n        - gate: regra-implementada\n")
	cfg := escreve(t, dir, "anchors.yaml", "version: 2\ngates:\n    - id: regra-implementada\n")
	for _, p := range []string{graph, cfg} {
		if _, err := MigrateFile(p, 3, false); err != nil {
			t.Fatal(err)
		}
	}
	if !strings.Contains(string(mustRead(t, graph)), "gate: regra-implementada") {
		t.Error("format 3 renamed a gate in the map")
	}
	if !strings.Contains(string(mustRead(t, cfg)), "id: regra-implementada") {
		t.Error("format 3 renamed an id in the configuration")
	}
}

func TestFormat4Step_renamesNothingInTheMap(t *testing.T) {
	t.Run("MGSTM-X02: Format 4 renames nothing in the map", func(t *testing.T) {})
	p := escreve(t, t.TempDir(), "anchors.graph.yaml", "version: 3\nauto_judgment: true\nrule_marking: required\n")
	r, err := MigrateFile(p, 4, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Replaced) != 0 || string(mustRead(t, p)) != "version: 4\nauto_judgment: true\nrule_marking: required\n" {
		t.Errorf("format 4 touched the map: %+v\n%s", r, mustRead(t, p))
	}
}
