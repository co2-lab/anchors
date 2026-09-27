package doct

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const featureDeExemplo = `# language: pt
Funcionalidade: A régua de go-live

  Contexto:
    Dado que o projeto tem uma régua

  @GLCGL-B01 @nivel-unit
  Cenário: item sem artefato não passa
    Dado um item que não cita artefato
    Quando confiro a régua
    Então a régua não libera

  @GLCGL-B02 @nivel-unit
  Cenário: dívida aberta bloqueia
    Dado uma dívida em aberto
    Quando confiro a régua
    Então a régua não libera
`

func projetoComFeature(t *testing.T) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "pkg"), 0o755)
	os.WriteFile(filepath.Join(root, "pkg/GoLive.spec.md"), []byte(specDeExemplo), 0o644)
	os.WriteFile(filepath.Join(root, "pkg/GoLive.feature"), []byte(featureDeExemplo), 0o644)
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "pkg/GoLive.spec.md", Kind: mapx.KindSpec, Code: "GLCGL", Layer: "spec"},
			{ID: "pkg/GoLive.feature", Kind: mapx.KindFeature, Code: "GLCGL", Layer: "feature"},
		},
		Edges: []mapx.Edge{
			{From: "pkg/GoLive.spec.md", To: "pkg/GoLive.feature", Type: "covered-by"},
		},
	}
	return root, g
}

// The scenario's BODY goes into the docs, not only the title: a gate only needs the code,
// documentation needs the steps — a title alone is a headline.
func TestScenarios_capturesTheBody(t *testing.T) {
	t.Run("GSRGH-B03: The body is the scenario's non-blank lines until the next scenario", func(t *testing.T) {})
	root, g := projetoComFeature(t)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	specs, _ := c.fnSpecs("layer=infra")
	cs := c.fnScenarios(specs[0])

	if len(cs) != 2 {
		t.Fatalf("found %d scenarios, want 2", len(cs))
	}
	if cs[0].Code != "GLCGL-B01" || cs[0].Titulo != "item sem artefato não passa" {
		t.Errorf("first scenario = %+v", cs[0])
	}
	if !strings.Contains(cs[0].Corpo, "Então a régua não libera") {
		t.Errorf("the body did not come: %q", cs[0].Corpo)
	}
	if strings.Contains(cs[0].Corpo, "Dado uma dívida") {
		t.Errorf("the body ran into the next scenario: %q", cs[0].Corpo)
	}
	blank := parseScenarios("Feature: X\n  @ABCDE-B01\n  Scenario: s\n    Given a\n\n    When b\n", "ABCDE")
	if blank[0].Corpo != "    Given a\n    When b" {
		t.Errorf("a blank line inside the scenario should be dropped: %q", blank[0].Corpo)
	}
}

// A `Contexto:` (Background) is not the previous scenario's: its steps must not leak into it.
func TestScenarios_backgroundDoesNotLeakIntoThePreviousScenario(t *testing.T) {
	t.Run("GSRGH-B04: A Background heading after a scenario does not leak into it", func(t *testing.T) {})
	feat := `Funcionalidade: X

  @ABCDE-B01 @nivel-unit
  Cenário: primeiro
    Dado A

  Contexto:
    Dado que isto é preparação
`
	cs := parseScenarios(feat, "ABCDE")
	if len(cs) != 1 {
		t.Fatalf("found %d, want 1", len(cs))
	}
	if strings.Contains(cs[0].Corpo, "preparação") {
		t.Errorf("the Background entered the scenario: %q", cs[0].Corpo)
	}
	for _, heading := range []string{"Background:", "Rule: r", "Feature: Y"} {
		en := parseScenarios("Feature: X\n  @ABCDE-B01\n  Scenario: s\n    Given a\n  "+heading+"\n    Given leaked\n", "ABCDE")
		if strings.Contains(en[0].Corpo, "leaked") {
			t.Errorf("%q did not close the scenario: %q", heading, en[0].Corpo)
		}
	}
}

// The REGIME tag is not mistaken for the identity code.
func TestScenarios_separatesCodeFromTag(t *testing.T) {
	t.Run("GSRGH-B02: The first code-shaped tag is the code and the others are tags", func(t *testing.T) {})
	cs := parseScenarios(featureDeExemplo, "GLCGL")
	if cs[0].Code != "GLCGL-B01" {
		t.Errorf("code = %q, want GLCGL-B01", cs[0].Code)
	}
	if len(cs[0].Tags) != 1 || cs[0].Tags[0] != "nivel-unit" {
		t.Errorf("tags = %v, want [nivel-unit]", cs[0].Tags)
	}
	// Only the FIRST code-shaped tag is the identity; a later one is just a tag.
	two := parseScenarios("Feature: X\n  @ABCDE-B01 @a @ABCDE-B02\n  Scenario: s\n    Given a\n", "ABCDE")
	if two[0].Code != "ABCDE-B01" || strings.Join(two[0].Tags, ",") != "a,ABCDE-B02" {
		t.Errorf("code = %q, tags = %v; want ABCDE-B01 and [a ABCDE-B02]", two[0].Code, two[0].Tags)
	}
}

// A feature in ENGLISH reads the same. A documentation that ignores half the scenarios
// because of the file's language is nobody's documentation.
func TestScenarios_readsEnglishFeature(t *testing.T) {
	feat := "Feature: X\n\n  @ABCDE-B01\n  Scenario: it works\n    Given a thing\n"
	cs := parseScenarios(feat, "ABCDE")
	if len(cs) != 1 || cs[0].Titulo != "it works" {
		t.Fatalf("scenarios = %+v", cs)
	}
}

// The feature is found through the MAP, not by name convention — which would break in the
// first project that organised its files another way.
func TestScenarios_findsTheFeatureByTheEdge(t *testing.T) {
	t.Run("GSRGH-B05: The feature is found by the map edge, in either direction, not by name", func(t *testing.T) {})
	root, g := projetoComFeature(t)
	// The feature is renamed; the edge follows.
	os.Rename(filepath.Join(root, "pkg/GoLive.feature"), filepath.Join(root, "pkg/Outro.feature"))
	g.Nodes[1].ID = "pkg/Outro.feature"
	g.Edges[0].To = "pkg/Outro.feature"

	c, _ := New(root, g)
	specs, _ := c.fnSpecs("layer=infra")
	if cs := c.fnScenarios(specs[0]); len(cs) != 2 {
		t.Errorf("found %d scenarios with the feature renamed, want 2", len(cs))
	}

	// The edge from the feature to the spec counts the same.
	root2, g2 := projetoComFeature(t)
	g2.Edges[0] = mapx.Edge{From: "pkg/GoLive.feature", To: "pkg/GoLive.spec.md", Type: "covers"}
	c2, _ := New(root2, g2)
	specs2, _ := c2.fnSpecs("layer=infra")
	if cs := c2.fnScenarios(specs2[0]); len(cs) != 2 {
		t.Errorf("found %d scenarios through a feature→spec edge, want 2", len(cs))
	}
}

// A SCENARIO OUTLINE is a scenario. The title regex required the line to START with
// Cenário/Scenario, and `Esquema do Cenário:` never matched: 7 outlines were missing from
// blue-eyes' comportamento.md. The `Exemplos:` table must not pass for a scenario.
func TestScenarios_outlineIsAScenario(t *testing.T) {
	t.Run("GSRGH-B01: Scenarios and outlines open in any dialect, and an examples table does not", func(t *testing.T) {})
	feat := "# language: pt\nFuncionalidade: X\n\n" +
		"  @ABCDE-B01\n  Esquema do Cenário: por <caso>\n    Dado <caso>\n\n" +
		"    Exemplos:\n      | caso |\n      | a    |\n\n" +
		"  @ABCDE-B02\n  Cenário: comum\n    Dado algo\n"
	cs := parseScenarios(feat, "ABCDE")
	if len(cs) != 2 {
		t.Fatalf("want the outline and the scenario, got %+v", cs)
	}
	if cs[0].Titulo != "por <caso>" || cs[1].Titulo != "comum" {
		t.Errorf("titles = %q, %q", cs[0].Titulo, cs[1].Titulo)
	}
	if !strings.Contains(cs[0].Corpo, "| a    |") {
		t.Errorf("the Examples table did not stay in the outline's body: %q", cs[0].Corpo)
	}
	en := parseScenarios("Feature: X\n\n  @ABCDE-B01\n  Scenario Outline: by <c>\n    Given <c>\n", "ABCDE")
	if len(en) != 1 || en[0].Titulo != "by <c>" {
		t.Errorf("English outline = %+v", en)
	}
}

func TestAllScenarios_followsTheSpecSelection(t *testing.T) {
	t.Run("GSRGH-B06: The scenarios of a selection follow the spec selection and its errors", func(t *testing.T) {
		root, g := projetoComFeature(t)
		c, _ := New(root, g)
		cs, err := c.fnAllScenarios("layer=infra")
		if err != nil || len(cs) != 2 {
			t.Fatalf("layer=infra gave %d scenarios, %v; want 2", len(cs), err)
		}
		_, want := c.fnSpecs("layer=nothere")
		if _, err := c.fnAllScenarios("layer=nothere"); err == nil || err.Error() != want.Error() {
			t.Errorf("error = %v, want the spec selection's %v", err, want)
		}
	})
}

func TestScenarios_carryTheSpecCode(t *testing.T) {
	t.Run("GSRGH-I01: Every scenario read carries its spec's code", func(t *testing.T) {
		root, g := projetoComFeature(t)
		c, _ := New(root, g)
		specs, _ := c.fnSpecs("layer=infra")
		for _, s := range c.fnScenarios(specs[0]) {
			if s.Spec != "GLCGL" {
				t.Errorf("scenario %q names spec %q, want GLCGL", s.Titulo, s.Spec)
			}
		}
	})
	t.Run("GSRGH-X01: The body is kept verbatim", func(t *testing.T) {
		cs := parseScenarios("Feature: X\n  @ABCDE-B01\n  Scenario: s\n    Given   a  thing\n", "ABCDE")
		if cs[0].Corpo != "    Given   a  thing" {
			t.Errorf("body = %q, want the line exactly as written", cs[0].Corpo)
		}
	})
}

func TestScenarios_missingFeatureIsSkipped(t *testing.T) {
	t.Run("GSRGH-E01: A linked feature missing on disk contributes no scenario and no error", func(t *testing.T) {
		root, g := projetoComFeature(t)
		g.Nodes = append(g.Nodes, mapx.Node{ID: "pkg/Gone.feature", Kind: mapx.KindFeature})
		g.Edges = append(g.Edges, mapx.Edge{From: "pkg/GoLive.spec.md", To: "pkg/Gone.feature", Type: "covered-by"})
		c, err := New(root, g)
		if err != nil {
			t.Fatal(err)
		}
		specs, _ := c.fnSpecs("layer=infra")
		if cs := c.fnScenarios(specs[0]); len(cs) != 2 {
			t.Errorf("found %d scenarios, want the readable feature's 2", len(cs))
		}
	})
}
