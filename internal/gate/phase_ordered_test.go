package gate

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func TestPhaseOrdered_B01_PlanPhasesExtraiCodigos(t *testing.T) {
	t.Run("PHORP-B01: The PlanPhases function extracts catalogued phase codes", func(t *testing.T) {})
	plan := "### FNDTN-F01 — a árvore\n\n### FNDTN-F02 — a régua\n\n#### FNDTN-F03 — o teste\n"
	phases := PlanPhases(plan)
	expected := []string{"FNDTN-F01", "FNDTN-F02", "FNDTN-F03"}
	if !reflect.DeepEqual(phases, expected) {
		t.Fatalf("esperava %v, veio %v", expected, phases)
	}
}

func TestPhaseOrdered_B02_NaoPlanoPula(t *testing.T) {
	t.Run("PHORP-B02: Non-plan artifacts skip phase ordering confrontation", func(t *testing.T) {})
	plan := "### FNDTN-F01 — a árvore\n"
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindCode, mapx.KindFeature, mapx.KindTest} {
		v, _ := checkPhaseOrdered(plan, mapx.Node{Kind: k}, "", nil, nil)
		if v != Skip {
			t.Fatalf("esperava Skip para kind %v, veio %v", k, v)
		}
	}
}

func TestPhaseOrdered_B03_PlanoSemFasesPula(t *testing.T) {
	t.Run("PHORP-B03: Plans without phase headings skip confrontation", func(t *testing.T) {})
	t.Run("PHORP-X01: Small plans are not required to catalog phases", func(t *testing.T) {})
	plano := mapx.Node{Kind: mapx.KindPlan, Code: "FNDTN"}
	content := "## Objetivo\n\nPlano pequeno sem divisão em fases.\n"
	v, _ := checkPhaseOrdered(content, plano, "", nil, nil)
	if v != Skip {
		t.Fatalf("plano sem fases deve pular; veio %v", v)
	}
}

func TestPhaseOrdered_B04_FaseEmProsaFicaPendente(t *testing.T) {
	t.Run("PHORP-B04: Plans with phase-like sections lacking codes return Pending", func(t *testing.T) {})
	t.Run("PHORP-I04: Phase detection identifies level-three sections regardless of language", func(t *testing.T) {})
	plano := mapx.Node{Kind: mapx.KindPlan, Code: "FNDTN"}
	prosa := "## Fases\n\n### Fase 1 — a árvore\n\n### Fase 2 — depende da Fase 1\n"
	v, msg := checkPhaseOrdered(prosa, plano, "", nil, nil)
	if v != Pending {
		t.Fatalf("fase em prosa deve retornar Pending; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "FNDTN") {
		t.Errorf("mensagem deve conter o código do plano; veio %s", msg)
	}
}

func TestPhaseOrdered_B05_OrdemValidaPassa(t *testing.T) {
	t.Run("PHORP-B05: Plans declaring valid backward phase dependencies pass", func(t *testing.T) {})
	t.Run("PHORP-X02: Phase duration and calendar timing are not verified", func(t *testing.T) {})
	plano := mapx.Node{Kind: mapx.KindPlan, Code: "FNDTN"}
	ok := "### FNDTN-F01 — a árvore\n\n### FNDTN-F02 — a régua (depende de FNDTN-F01)\n"
	v, msg := checkPhaseOrdered(ok, plano, "", nil, nil)
	if v != Pass {
		t.Fatalf("ordem válida deve passar; veio %v (%s)", v, msg)
	}
}

func TestPhaseOrdered_B06_CodigoDuplicadoFalha(t *testing.T) {
	t.Run("PHORP-B06: Plans with duplicate phase codes fail", func(t *testing.T) {})
	plano := mapx.Node{Kind: mapx.KindPlan, Code: "FNDTN"}
	repetida := "### FNDTN-F01 — a árvore\n\n### FNDTN-F01 — outra fase\n"
	v, msg := checkPhaseOrdered(repetida, plano, "", nil, nil)
	if v != Fail {
		t.Fatalf("código repetido deve falhar; veio %v", v)
	}
	if !strings.Contains(msg, "FNDTN-F01") {
		t.Errorf("mensagem deve conter o código duplicado; veio %s", msg)
	}
}

func TestPhaseOrdered_B07_FaseInexistenteNoPlanoFalha(t *testing.T) {
	t.Run("PHORP-B07: Phases depending on uncatalogued phase codes fail", func(t *testing.T) {})
	plano := mapx.Node{Kind: mapx.KindPlan, Code: "FNDTN"}
	fantasma := "### FNDTN-F01 — a árvore (depende de FNDTN-F09)\n"
	v, msg := checkPhaseOrdered(fantasma, plano, "", nil, nil)
	if v != Fail {
		t.Fatalf("dependência de fase fantasma deve reprovar; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "FNDTN-F09") {
		t.Errorf("mensagem deve conter a fase ausente; veio %s", msg)
	}
}

func TestPhaseOrdered_B08_DependeDeSiMesmoFalha(t *testing.T) {
	t.Run("PHORP-B08: Phases depending on themselves fail", func(t *testing.T) {})
	plano := mapx.Node{Kind: mapx.KindPlan, Code: "FNDTN"}
	auto := "### FNDTN-F01 — a árvore (depende de FNDTN-F01)\n"
	v, msg := checkPhaseOrdered(auto, plano, "", nil, nil)
	if v != Fail {
		t.Fatalf("auto-dependência deve reprovar; veio %v", v)
	}
	if !strings.Contains(msg, "FNDTN-F01") {
		t.Errorf("mensagem deve citar o código da fase com auto-dependência; veio %s", msg)
	}
}

func TestPhaseOrdered_B09_DependeDoFuturoFalha(t *testing.T) {
	t.Run("PHORP-B09: Phases depending on future phases fail", func(t *testing.T) {})
	t.Run("PHORP-I01: Phase dependencies are strictly acyclic and backward-directed", func(t *testing.T) {})
	plano := mapx.Node{Kind: mapx.KindPlan, Code: "FNDTN"}
	invertida := "### FNDTN-F01 — a árvore (depende de FNDTN-F02)\n\n### FNDTN-F02 — a régua\n"
	v, msg := checkPhaseOrdered(invertida, plano, "", nil, nil)
	if v != Fail {
		t.Fatalf("dependência de fase futura deve reprovar; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "DEPOIS") && !strings.Contains(msg, "AFTER") {
		t.Errorf("mensagem deve citar que a fase vem depois; veio %s", msg)
	}
}

func TestPhaseOrdered_B10_PhaseExistsPassa(t *testing.T) {
	t.Run("PHORP-B10: Specifications declaring existing phase dependencies pass", func(t *testing.T) {})
	planoID := "plans/0001.md"
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: planoID, Kind: mapx.KindPlan}}}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	plano := "### FNDTN-F01 — a árvore\n\n### FNDTN-F02 — a régua\n\n### FNDTN-F03 — o teste\n"
	if err := os.WriteFile(filepath.Join(dir, planoID), []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{Kind: mapx.KindSpec, Needs: []string{"FNDTN-F01", "FNDTN-F02"}}
	if v, msg := checkPhaseExists("", spec, dir, g, nil); v != Pass {
		t.Fatalf("fases existentes devem passar; veio %v (%s)", v, msg)
	}
}

func TestPhaseOrdered_B11_PhaseExistsAusenteFalha(t *testing.T) {
	t.Run("PHORP-B11: Specifications declaring missing phase dependencies fail", func(t *testing.T) {})
	t.Run("PHORP-I02: Phase and parent targets must exist in the map", func(t *testing.T) {})
	planoID := "plans/0001.md"
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: planoID, Kind: mapx.KindPlan}}}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	plano := "### FNDTN-F01 — a árvore\n"
	if err := os.WriteFile(filepath.Join(dir, planoID), []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{Kind: mapx.KindSpec, Needs: []string{"FNDTN-F01", "FNDTN-F09"}}
	v, msg := checkPhaseExists("", spec, dir, g, nil)
	if v != Fail {
		t.Fatalf("fase inexistente deve reprovar; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "FNDTN-F09") {
		t.Errorf("mensagem deve conter a fase ausente FNDTN-F09; veio %s", msg)
	}
	if strings.Contains(msg, "FNDTN-F01") {
		t.Errorf("fase existente não deve ser reportada como ausente; veio %s", msg)
	}
}

func TestPhaseOrdered_B12_ParentValidoPassa(t *testing.T) {
	t.Run("PHORP-B12: Artifacts declaring valid parents pass", func(t *testing.T) {})
	t.Run("PHORP-X03: Both artifact codes and phase codes are accepted as parents", func(t *testing.T) {})
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	plano := "### FNDTN-F01 — a árvore\n"
	if err := os.WriteFile(filepath.Join(dir, "plans/0001.md"), []byte(plano), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "plans/0001.md", Kind: mapx.KindPlan, Code: "FNDTN"},
		{ID: "a.spec.md", Kind: mapx.KindSpec, Code: "WRKSP", Parent: "FNDTN-F01"},
		{ID: "b.spec.md", Kind: mapx.KindSpec, Code: "OTHER", Parent: "FNDTN"},
	}}
	// Pai fase catalogada
	if v, msg := checkParentValid("", g.Nodes[1], dir, g, nil); v != Pass {
		t.Fatalf("pai fase catalogada deve passar; veio %v (%s)", v, msg)
	}
	// Pai código de artefato
	if v, msg := checkParentValid("", g.Nodes[2], dir, g, nil); v != Pass {
		t.Fatalf("pai código de artefato deve passar; veio %v (%s)", v, msg)
	}
}

func TestPhaseOrdered_B13_ParentInvalidoFalha(t *testing.T) {
	t.Run("PHORP-B13: Artifacts declaring invalid parents, self-parenting, or parent cycles fail", func(t *testing.T) {})
	t.Run("PHORP-I03: Parent chains are cycle-free and bounded", func(t *testing.T) {})
	dir := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.md", Code: "AAAAA", Parent: "BBBBB"},
		{ID: "b.md", Code: "BBBBB", Parent: "AAAAA"},
	}}
	// Inexistente
	v1, msg1 := checkParentValid("", mapx.Node{Code: "X", Parent: "FNDTN-F09"}, dir, g, nil)
	if v1 != Fail {
		t.Fatalf("pai inexistente deve reprovar; veio %v (%s)", v1, msg1)
	}
	// Auto-pai
	v2, msg2 := checkParentValid("", mapx.Node{Code: "AAAAA", Parent: "AAAAA"}, dir, g, nil)
	if v2 != Fail {
		t.Fatalf("auto-pai deve reprovar; veio %v (%s)", v2, msg2)
	}
	// Ciclo
	v3, msg3 := checkParentValid("", g.Nodes[0], dir, g, nil)
	if v3 != Fail {
		t.Fatalf("ciclo de parent deve reprovar; veio %v (%s)", v3, msg3)
	}
}

func TestLetraDaFaseEhCanonica(t *testing.T) {
	if !strings.Contains(config.DefaultRuleLetters, "F") {
		t.Errorf("`F` (fase) não está em %s — toda spec com `needs: CODE-F01` reprovaria no rule-types", config.DefaultRuleLetters)
	}
	if !phaseRE().MatchString("### ABCDX-F01 — a fase") {
		t.Error("o regex de fase não casa o formato que o próprio gate documenta")
	}
}

func TestPhaseOrdered_Errors(t *testing.T) {
	t.Run("PHORP-E01: With no map a phase or parent reference is not judged", func(t *testing.T) {
		spec := mapx.Node{Kind: mapx.KindSpec, Needs: []string{"FNDTN-F01"}}
		if v, msg := checkPhaseExists("", spec, t.TempDir(), nil, nil); v != Pending {
			t.Fatalf("needs with no map must be Pending, got %v (%s)", v, msg)
		}
		child := mapx.Node{Code: "X", Parent: "FNDTN-F01"}
		if v, msg := checkParentValid("", child, t.TempDir(), nil, nil); v != Pending {
			t.Fatalf("parent with no map must be Pending, got %v (%s)", v, msg)
		}
	})

	t.Run("PHORP-E02: A plan missing from disk does not unresolve the phases of the other plans", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, "plans"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "plans/0002.md"), []byte("### FNDTN-F01 — the tree\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		g := &mapx.Graph{Nodes: []mapx.Node{
			{ID: "plans/0001.md", Kind: mapx.KindPlan},
			{ID: "plans/0002.md", Kind: mapx.KindPlan},
		}}
		spec := mapx.Node{Kind: mapx.KindSpec, Needs: []string{"FNDTN-F01"}}
		if v, msg := checkPhaseExists("", spec, dir, g, nil); v != Pass {
			t.Fatalf("the phase of the plan on disk must resolve needs, got %v (%s)", v, msg)
		}
		child := mapx.Node{Code: "X", Parent: "FNDTN-F01"}
		if v, msg := checkParentValid("", child, dir, g, nil); v != Pass {
			t.Fatalf("the phase of the plan on disk must resolve parent, got %v (%s)", v, msg)
		}
	})
}
