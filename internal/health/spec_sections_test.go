package health

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// writeSpec creates a spec on disk and returns its node.
func writeSpec(t *testing.T, root, nome, corpo string) mapx.Node {
	t.Helper()
	caminho := filepath.Join(root, nome)
	if err := os.WriteFile(caminho, []byte(corpo), 0o644); err != nil {
		t.Fatalf("writing %s: %v", nome, err)
	}
	// Layer "spec" on purpose: it is what the map usually carries (the FILE's layer).
	// The UNIT's layer is said by the header — and that is what the check must read.
	return mapx.Node{ID: nome, Kind: mapx.KindSpec, Layer: "spec"}
}

const specTela = `<!-- @anchors
  code: SGINS
  layer: screen
-->
# SignIn

## Visão Geral
Entra no app.

## Regras

### SGINS-B01 — regra
Comportamento.
`

// The distinction the check is for: a DECLARED and blind gate is Warn and asks to edit the
// specs; an UNDECLARED gate is Info and asks to adopt the section AND the gate. The
// corrective action differs, so the severity and the text must too.
func TestCheckSpecSections_DeclaredBlindGateIsWarn(t *testing.T) {
	t.Run("SPSCS-B03: A declared gate that is blind gives a warning", func(t *testing.T) {})
	t.Run("SPSCS-B06: The header layer wins over the map's", func(t *testing.T) {})
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{writeSpec(t, root, "SignIn.spec.md", specTela)}}
	cfg := &config.Config{Gates: []config.Gate{{Name: "dependency-honored"}}}

	findings := checkSpecSections(g, cfg, root)
	var found *Finding
	for i := range findings {
		if findings[i].Subject == "data-contract" {
			found = &findings[i]
		}
	}
	if found == nil {
		t.Fatal("expected a data-contract finding: the screen spec lacks the section")
	}
	if found.Check != "secao-ausente" {
		t.Errorf("check = %q, want secao-ausente (the gate is declared)", found.Check)
	}
	if found.Severity != Warn {
		t.Errorf("severity = %v, want Warn: a declared gate is blind", found.Severity)
	}
}

func TestCheckSpecSections_UndeclaredGateIsInfo(t *testing.T) {
	t.Run("SPSCS-B04: An undeclared gate gives an informational recommendation that asks to declare it", func(t *testing.T) {})
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{writeSpec(t, root, "SignIn.spec.md", specTela)}}
	// No gate declared: `route-declared` does not exist in this project.
	cfg := &config.Config{}

	for _, f := range checkSpecSections(g, cfg, root) {
		if f.Subject != "navigation" {
			continue
		}
		if f.Check != "secao-recomendada" {
			t.Errorf("check = %q, want secao-recomendada (gate not declared)", f.Check)
		}
		if f.Severity != Info {
			t.Errorf("severity = %v, want Info", f.Severity)
		}
		if !strings.Contains(f.Detail, "route-declared") {
			t.Errorf("the recommendation must ask to declare the gate: %s", f.Detail)
		}
		return
	}
	t.Fatal("expected a navigation finding")
}

// A present section cannot be accused — and the title counts in ANY supported language,
// because the spec was written under the `lang:` the project had that day.
func TestCheckSpecSections_TitleInAnotherLanguageCounts(t *testing.T) {
	t.Run("SPSCS-B05: A title in another language counts", func(t *testing.T) {})
	root := t.TempDir()
	withEnglishNavigation := specTela + "\n## Navigation\n| Origem | Gatilho |\n| --- | --- |\n"
	g := &mapx.Graph{Nodes: []mapx.Node{writeSpec(t, root, "SignIn.spec.md", withEnglishNavigation)}}

	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		if f.Subject == "navigation" {
			t.Errorf("navigation reported missing, but the section exists with an English title: %s", f.Detail)
		}
	}
}

// Charging a route or testid to a unit that is not a screen is the false positive of the
// legacy validator, which treated every artefact as a screen.
func TestCheckSpecSections_ScreenSectionsNotDemandedFromOtherLayers(t *testing.T) {
	t.Run("SPSCS-B07: Screen sections are not demanded from other layers", func(t *testing.T) {})
	root := t.TempDir()
	logic := `<!-- @anchors
  code: CALCX
  layer: backend-logic
-->
# Calc

## Visão Geral
Soma.

## Regras

### CALCX-B01 — regra
Comportamento.
`
	g := &mapx.Graph{Nodes: []mapx.Node{writeSpec(t, root, "Calc.spec.md", logic)}}

	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		switch f.Subject {
		case "navigation", "testids", "data-contract", "data-states", "states":
			t.Errorf("a screen section (%s) was demanded from a backend-logic unit", f.Subject)
		}
	}
}

// The check reports the PATTERN, not the isolated case: one spec without the section among
// several that have it is its author's decision, and a warning would train the team to
// ignore the output.
func TestCheckSpecSections_MinorityMissingIsSilent(t *testing.T) {
	t.Run("SPSCS-B02: Exactly half or a minority missing is silent", func(t *testing.T) {})
	root := t.TempDir()
	withDomain := `<!-- @anchors
  code: AAAAA
  layer: backend-logic
-->
# A

## Visão Geral
x

## Domínio
| Entrada | Aceita |
| --- | --- |
`
	withoutDomain := `<!-- @anchors
  code: BBBBB
  layer: backend-logic
-->
# B

## Visão Geral
y
`
	g := &mapx.Graph{Nodes: []mapx.Node{
		writeSpec(t, root, "A.spec.md", withDomain),
		writeSpec(t, root, "B.spec.md", withDomain),
		writeSpec(t, root, "C.spec.md", withDomain),
		writeSpec(t, root, "D.spec.md", withoutDomain),
	}}

	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		if f.Subject == "domain" {
			t.Errorf("domain reported with only 1 of 4 missing — it should be silent: %s", f.Detail)
		}
	}

	// Exactly half is still not the pattern; three of four is.
	writeSpec(t, root, "C.spec.md", withoutDomain)
	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		if f.Subject == "domain" {
			t.Errorf("domain reported with 2 of 4 missing — half is not a majority: %s", f.Detail)
		}
	}
	writeSpec(t, root, "B.spec.md", withoutDomain)
	var reported bool
	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		reported = reported || f.Subject == "domain"
	}
	if !reported {
		t.Error("3 of 4 missing is the pattern and must be reported")
	}
}

func TestCheckSpecSections_SixRecommendedSections(t *testing.T) {
	t.Run("SPSCS-B01: A screen spec without any recommended section is told about all six", func(t *testing.T) {})
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{writeSpec(t, root, "SignIn.spec.md", specTela)}}
	got := map[string]bool{}
	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		got[f.Subject] = true
	}
	for _, key := range []string{"navigation", "data-contract", "data-states", "testids", "domain", "states"} {
		if !got[key] {
			t.Errorf("the screen spec lacks %s and it was not reported: %v", key, got)
		}
	}
	if len(got) != 6 {
		t.Errorf("exactly six sections are recommended, got %v", got)
	}
}

func TestCheckSpecSections_OneFindingPerSection(t *testing.T) {
	t.Run("SPSCS-B08: Three specs lacking a section give one finding with the numbers", func(t *testing.T) {})
	t.Run("SPSCS-X01: No finding names a spec file", func(t *testing.T) {})
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{
		writeSpec(t, root, "A.spec.md", specTela),
		writeSpec(t, root, "B.spec.md", specTela),
		writeSpec(t, root, "C.spec.md", specTela),
	}}
	var navigation int
	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		if strings.HasSuffix(f.Subject, ".spec.md") {
			t.Errorf("a finding names a spec file: %+v", f)
		}
		if f.Subject == "navigation" {
			navigation++
			if !strings.Contains(f.Detail, "3 of 3") && !strings.Contains(f.Detail, "3 de 3") {
				t.Errorf("the finding must carry the numbers: %s", f.Detail)
			}
		}
	}
	if navigation != 1 {
		t.Errorf("three specs lacking navigation give one finding, got %d", navigation)
	}
}

func TestCheckSpecSections_NilInputs(t *testing.T) {
	t.Run("SPSCS-B09: A nil map or configuration gives nothing", func(t *testing.T) {})
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{writeSpec(t, root, "SignIn.spec.md", specTela)}}
	if fs := checkSpecSections(nil, &config.Config{}, root); fs != nil {
		t.Errorf("a nil map gives nothing: %+v", fs)
	}
	if fs := checkSpecSections(g, nil, root); fs != nil {
		t.Errorf("a nil configuration gives nothing: %+v", fs)
	}
}

func TestCheckSpecSections_AddingTheSectionClearsIt(t *testing.T) {
	t.Run("SPSCS-I01: Adding the reported section removes the finding", func(t *testing.T) {})
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{writeSpec(t, root, "SignIn.spec.md", specTela)}}
	reported := func() bool {
		for _, f := range checkSpecSections(g, &config.Config{}, root) {
			if f.Subject == "navigation" {
				return true
			}
		}
		return false
	}
	if !reported() {
		t.Fatal("setup: the spec lacks navigation and must be reported")
	}
	writeSpec(t, root, "SignIn.spec.md", specTela+"\n## Navigation\n")
	if reported() {
		t.Fatal("the section was added and navigation is still reported")
	}
}

func TestCheckSpecSections_UnreadableSpecsAreLeftOut(t *testing.T) {
	t.Run("SPSCS-E01: Specs missing on disk are left out of the counts", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "Gone.spec.md", Kind: mapx.KindSpec, Layer: "screen"},
		{ID: "AlsoGone.spec.md", Kind: mapx.KindSpec, Layer: "screen"},
	}}
	if fs := checkSpecSections(g, &config.Config{}, t.TempDir()); len(fs) != 0 {
		t.Fatalf("a spec that cannot be read is not counted: %+v", fs)
	}
}
