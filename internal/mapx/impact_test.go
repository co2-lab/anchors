package mapx

import (
	"slices"
	"testing"
)

// grafo: guide rege spec do Login E spec do Home (régua compartilhada). Trinca do
// Login completa. Prova: subir do Login pega o guide (valida), mas NÃO desce do
// guide para o Home (subir não re-propaga).
func impactGraph() *Graph {
	return &Graph{
		Version: 1,
		Nodes: []Node{
			{ID: "Login.spec.md", Kind: KindSpec},
			{ID: "Login.tsx", Kind: KindCode},
			{ID: "Login.feature", Kind: KindFeature},
			{ID: "Login.test.tsx", Kind: KindTest},
			{ID: "SPEC_GUIDE.md", Kind: KindGuide},
			{ID: "Home.spec.md", Kind: KindSpec},
			{ID: "Home.tsx", Kind: KindCode},
		},
		Edges: []Edge{
			{From: "Login.spec.md", To: "Login.tsx", Type: EdgeSpecifies},
			{From: "Login.spec.md", To: "Login.feature", Type: EdgeCoveredBy},
			{From: "Login.feature", To: "Login.test.tsx", Type: EdgeTestedBy},
			{From: "SPEC_GUIDE.md", To: "Login.spec.md", Type: EdgeGoverns},
			{From: "SPEC_GUIDE.md", To: "Home.spec.md", Type: EdgeGoverns},
			{From: "Home.spec.md", To: "Home.tsx", Type: EdgeSpecifies},
		},
	}
}

// Mexer na SPEC do Login: propaga (desce) para código+feature+teste; valida (sobe)
// contra o guide. NÃO deve tocar o Home (nem descendo nem subindo).
func TestAnalyzeImpact_specChange(t *testing.T) {
	imp := impactGraph().AnalyzeImpact("Login.spec.md")

	wantProp := []string{"Login.feature", "Login.test.tsx", "Login.tsx"}
	if !slices.Equal(imp.Propagate, wantProp) {
		t.Errorf("Propagate = %v, want %v", imp.Propagate, wantProp)
	}
	wantVal := []string{"SPEC_GUIDE.md"}
	if !slices.Equal(imp.Validate, wantVal) {
		t.Errorf("Validate = %v, want %v", imp.Validate, wantVal)
	}
	// o irmão Home NUNCA aparece — subir até o guide não re-desce para os irmãos
	if slices.Contains(imp.Propagate, "Home.spec.md") || slices.Contains(imp.Validate, "Home.spec.md") ||
		slices.Contains(imp.Propagate, "Home.tsx") || slices.Contains(imp.Validate, "Home.tsx") {
		t.Error("o alvo irmão Home não deveria ser alcançado (subir não re-propaga)")
	}
}

// Mexer no CÓDIGO (folha da trinca): não propaga para ninguém (nada depende do
// código abaixo); valida subindo contra spec e guide.
func TestAnalyzeImpact_codeChange(t *testing.T) {
	imp := impactGraph().AnalyzeImpact("Login.tsx")
	if len(imp.Propagate) != 0 {
		t.Errorf("mexer no código não deveria propagar para baixo; got %v", imp.Propagate)
	}
	wantVal := []string{"Login.spec.md", "SPEC_GUIDE.md"}
	if !slices.Equal(imp.Validate, wantVal) {
		t.Errorf("Validate = %v, want %v", imp.Validate, wantVal)
	}
}

// Mexer no GUIDE (pai de alto grau): propaga (desce) para TODAS as specs que ele
// rege — a onda global. É o comportamento correto de mudar uma régua.
func TestAnalyzeImpact_guideChange(t *testing.T) {
	imp := impactGraph().AnalyzeImpact("SPEC_GUIDE.md")
	// desce para as specs regidas e, delas, para os filhos das specs
	for _, want := range []string{"Login.spec.md", "Home.spec.md", "Login.tsx", "Home.tsx"} {
		if !slices.Contains(imp.Propagate, want) {
			t.Errorf("mudar o guide deveria propagar para %q (onda global); got %v", want, imp.Propagate)
		}
	}
}

// @noPropagation: um filho marcado não deixa a onda descer POR ELE.
func TestAnalyzeImpact_noPropagation(t *testing.T) {
	g := impactGraph()
	// marca Login.feature como @noPropagation → a onda que desce da spec alcança a
	// feature (ela é filha direta), mas NÃO continua da feature para o teste.
	for i := range g.Nodes {
		if g.Nodes[i].ID == "Login.feature" {
			g.Nodes[i].NoPropagation = true
		}
	}
	imp := g.AnalyzeImpact("Login.spec.md")
	if !slices.Contains(imp.Propagate, "Login.feature") {
		t.Error("a feature (filha direta) ainda deve ser alcançada")
	}
	if slices.Contains(imp.Propagate, "Login.test.tsx") {
		t.Error("o teste NÃO deveria ser alcançado — a feature @noPropagation não deixa a onda descer por ela")
	}
}

func TestAnalyzeImpact_isolado(t *testing.T) {
	g := &Graph{Nodes: []Node{{ID: "solto.md"}}}
	imp := g.AnalyzeImpact("solto.md")
	if len(imp.Propagate) != 0 || len(imp.Validate) != 0 {
		t.Errorf("nó isolado não tem impacto; got prop=%v val=%v", imp.Propagate, imp.Validate)
	}
}
