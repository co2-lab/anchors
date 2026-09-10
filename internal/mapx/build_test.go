package mapx

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

func testCfg() *config.Config {
	return &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec":   {Kind: "spec", Tags: []string{"spec"}},
			"screen": {Kind: "code", Tags: []string{"frontend"}},
			"guide":  {Kind: "guide", Tags: []string{"guide"}},
		},
		Derived: &config.Derived{
			Anchor: "code",
			Files: map[string]config.Padroes{
				"spec":    {"{{dir}}/{{name}}.spec.md"},
				"feature": {"{{dir}}/{{name}}.feature"},
				"test":    {"{{dir}}/{{name}}.test.{{ext}}"},
			},
		},
		Governs: []config.GovernRule{
			{From: "guides/SPEC_GUIDE.md", Governs: "spec"},
			{From: "guides/FRONTEND_GUIDE.md", Governs: "frontend"},
		},
	}
}

func testFiles() []scan.File {
	return []scan.File{
		{Path: "src/Login.tsx", Layer: "screen", Kind: "code", Rev: "a"},
		{Path: "src/Login.spec.md", Layer: "spec", Kind: "spec", Rev: "b", Codes: []string{"LOGIX-A01"}},
		{Path: "src/Login.feature", Layer: "feature", Kind: "feature", Rev: "c", Codes: []string{"LOGIX-A01"}},
		{Path: "src/Login.test.tsx", Layer: "test", Kind: "test", Rev: "d", Codes: []string{"LOGIX-A01"}},
		{Path: "guides/SPEC_GUIDE.md", Layer: "guide", Kind: "guide", Rev: "e"},
		{Path: "guides/FRONTEND_GUIDE.md", Layer: "guide", Kind: "guide", Rev: "f"},
	}
}

func hasEdge(g *Graph, from, to string, typ EdgeType) bool {
	for _, e := range g.Edges {
		if e.From == from && e.To == to && e.Type == typ {
			return true
		}
	}
	return false
}

func TestBuild_colocation(t *testing.T) {
	g := Build(testFiles(), testCfg(), nil)

	// a trinca liga: spec→código (specifies), spec→feature (covered-by), feature→teste
	cases := []struct {
		from, to string
		typ      EdgeType
	}{
		{"src/Login.spec.md", "src/Login.tsx", EdgeSpecifies},
		{"src/Login.spec.md", "src/Login.feature", EdgeCoveredBy},
		{"src/Login.feature", "src/Login.test.tsx", EdgeTestedBy},
	}
	for _, c := range cases {
		if !hasEdge(g, c.from, c.to, c.typ) {
			t.Errorf("faltou aresta de co-location: %s →%s→ %s", c.from, c.typ, c.to)
		}
	}
}

func TestBuild_derivedOverride(t *testing.T) {
	// Layout centralizado (STRUCTURE.md §2.2): a âncora handler mora em
	// functions/<módulo>/handler.ts; spec+feature co-localizadas, mas o TESTE mora numa
	// região central (__tests__/unit/lambdas/<módulo>.test.ts). O override + {{module}}
	// devem achar o teste lá — a co-location pura (irmão) não acharia.
	cfg := &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec":    {Kind: "spec", Tags: []string{"spec"}},
			"handler": {Kind: "code", Tags: []string{"backend", "handler"}},
		},
		Derived: &config.Derived{
			Anchor: "code",
			Files: map[string]config.Padroes{
				"spec":    {"{{dir}}/{{name}}.spec.md"},
				"feature": {"{{dir}}/{{name}}.feature"},
				"test":    {"{{dir}}/{{name}}.test.{{ext}}"},
			},
			Overrides: []config.DerivedOverride{
				{When: "handler", Files: map[string]config.Padroes{
					"test": {"__tests__/unit/lambdas/{{module}}.test.ts"},
				}},
			},
		},
	}
	files := []scan.File{
		{Path: "functions/run-audits/handler.ts", Layer: "handler", Kind: "code", Rev: "a"},
		{Path: "functions/run-audits/handler.spec.md", Layer: "spec", Kind: "spec", Rev: "b"},
		{Path: "functions/run-audits/handler.feature", Layer: "feature", Kind: "feature", Rev: "c"},
		{Path: "__tests__/unit/lambdas/run-audits.test.ts", Layer: "test", Kind: "test", Rev: "d"},
	}
	g := Build(files, cfg, nil)

	// spec→feature ainda liga por co-location (irmãos)
	if !hasEdge(g, "functions/run-audits/handler.spec.md", "functions/run-audits/handler.feature", EdgeCoveredBy) {
		t.Error("spec deveria cobrir a feature co-localizada")
	}
	// feature→teste liga pelo override {{module}} (teste NÃO co-localizado)
	if !hasEdge(g, "functions/run-audits/handler.feature", "__tests__/unit/lambdas/run-audits.test.ts", EdgeTestedBy) {
		t.Error("o override {{module}} deveria ligar a feature ao teste central")
	}
}

func TestBuild_governsByTag(t *testing.T) {
	g := Build(testFiles(), testCfg(), nil)

	// SPEC_GUIDE (tag spec) rege a spec
	if !hasEdge(g, "guides/SPEC_GUIDE.md", "src/Login.spec.md", EdgeGoverns) {
		t.Error("SPEC_GUIDE deveria reger a spec (tag spec)")
	}
	// FRONTEND_GUIDE (tag frontend) rege o código de screen
	if !hasEdge(g, "guides/FRONTEND_GUIDE.md", "src/Login.tsx", EdgeGoverns) {
		t.Error("FRONTEND_GUIDE deveria reger o código (tag frontend)")
	}
	// FRONTEND_GUIDE NÃO rege a spec (tag errada) — sem produto cartesiano
	if hasEdge(g, "guides/FRONTEND_GUIDE.md", "src/Login.spec.md", EdgeGoverns) {
		t.Error("FRONTEND_GUIDE NÃO deveria reger a spec (não tem tag spec)")
	}
}

func TestBuild_noScenarioDupOnColocated(t *testing.T) {
	// LOGIX-A01 está em spec/feature/test co-localizados; não deve gerar aresta
	// por identidade duplicando a co-location (furo #1 corrigido).
	g := Build(testFiles(), testCfg(), nil)
	refs := 0
	for _, e := range g.Edges {
		if e.Origin == OriginInferred {
			refs++
		}
	}
	if refs != 0 {
		t.Errorf("não deveria haver arestas inferred entre co-localizados; got %d", refs)
	}
}

func TestStale(t *testing.T) {
	g := &Graph{
		Nodes: []Node{{ID: "a", Rev: "r2"}, {ID: "b", Rev: "r5"}},
	}
	// aresta sem carimbo → stale
	if !g.Stale(Edge{From: "a", To: "b"}) {
		t.Error("aresta sem carimbo deveria ser stale")
	}
	// carimbo bate com as revs atuais → não stale
	fresh := Edge{From: "a", To: "b", Stamp: &Stamp{ValidatedFromRev: "r2", ValidatedToRev: "r5"}}
	if g.Stale(fresh) {
		t.Error("aresta com carimbo atual não deveria ser stale")
	}
	// uma ponta avançou → stale
	old := Edge{From: "a", To: "b", Stamp: &Stamp{ValidatedFromRev: "r2", ValidatedToRev: "r4"}}
	if !g.Stale(old) {
		t.Error("aresta com o alvo avançado deveria ser stale")
	}
}

func TestDependsOnEdges(t *testing.T) {
	files := []scan.File{
		{Path: "src/Login.spec.md", Layer: "spec", Kind: "spec", Rev: "b",
			Deps: []scan.Dep{
				{Code: "DEP1", File: "src/auth.store.ts", Method: "useAuthStore", Layer: "store"},
				{Code: "DEP2", File: "src/useAuth.ts", Method: "signIn", Layer: "hook"},
				{Code: "DEP3", File: "src/inexistente.ts", Method: "x", Layer: "hook"}, // alvo ausente → sem aresta
			}},
		{Path: "src/auth.store.ts", Layer: "store", Kind: "code", Rev: "c"},
		{Path: "src/useAuth.ts", Layer: "hook", Kind: "code", Rev: "d"},
	}
	edges := dependsOnEdges(files)
	if len(edges) != 2 {
		t.Fatalf("esperava 2 arestas (a 3ª aponta arquivo inexistente), veio %d: %+v", len(edges), edges)
	}
	// a aresta vai da SPEC para o ARQUIVO, com method+dep como metadados, origem declared
	var e1 *Edge
	for i := range edges {
		if edges[i].To == "src/auth.store.ts" {
			e1 = &edges[i]
		}
	}
	if e1 == nil {
		t.Fatal("aresta p/ auth.store.ts não construída")
	}
	if e1.From != "src/Login.spec.md" || e1.Type != EdgeDependsOn || e1.Origin != OriginDeclared {
		t.Errorf("aresta mal formada: %+v", *e1)
	}
	if e1.Method != "useAuthStore" || e1.Dep != "DEP1" {
		t.Errorf("metadados method/dep perdidos: %+v", *e1)
	}
}

func TestDependsOnEdgesIntegratedInBuild(t *testing.T) {
	files := append(testFiles(),
		scan.File{Path: "src/useLogin.ts", Layer: "hook", Kind: "code", Rev: "g"},
	)
	// injeta uma dep na spec de Login apontando o hook
	for i := range files {
		if files[i].Path == "src/Login.spec.md" {
			files[i].Deps = []scan.Dep{{Code: "DEP1", File: "src/useLogin.ts", Method: "run", Layer: "hook"}}
		}
	}
	g := Build(files, testCfg(), nil)
	if !hasEdge(g, "src/Login.spec.md", "src/useLogin.ts", EdgeDependsOn) {
		t.Error("Build deveria conter a aresta depends-on spec→hook")
	}
}

// TestSeedResolvidaPorNome: um plano cita a spec de duas formas legítimas — pelo caminho,
// ou só pelo NOME, que é como se escreve em prosa ("a spec de `SubscriptionScreen.spec.md`").
// Medido num repositório real: 10 das 26 citações são por nome.
//
// Tratar "citado por nome" como "não existe" fazia o plano parecer eternamente
// não-cumprido, e a partida a frio da fila semeava fases antigas já concluídas.
func TestSeedResolvidaPorNome(t *testing.T) {
	files := []scan.File{
		{Path: "plans/p.md", Kind: "plan", Seeds: []string{"Tela.spec.md", "apps/y/Outra.spec.md", "Ambigua.spec.md"}},
		{Path: "apps/x/Tela.spec.md", Kind: "spec"},
		{Path: "apps/y/Outra.spec.md", Kind: "spec"},
		// duas com o mesmo nome: a citação vira ambígua e NÃO deve gerar aresta —
		// escolher uma seria inventar uma ligação que o autor não declarou.
		{Path: "apps/a/Ambigua.spec.md", Kind: "spec"},
		{Path: "apps/b/Ambigua.spec.md", Kind: "spec"},
	}
	destinos := map[string]bool{}
	for _, e := range seedEdges(files) {
		destinos[e.To] = true
	}
	if !destinos["apps/x/Tela.spec.md"] {
		t.Error("citação por NOME devia resolver para o caminho único")
	}
	if !destinos["apps/y/Outra.spec.md"] {
		t.Error("citação por CAMINHO devia continuar funcionando")
	}
	if destinos["apps/a/Ambigua.spec.md"] || destinos["apps/b/Ambigua.spec.md"] {
		t.Error("nome duplicado é ambíguo — não pode gerar aresta")
	}
}

func TestBuild_ancoraSpec(t *testing.T) {
	// A spec ancorando, que é a forma canônica: os arquivos derivam dela. A extensão do
	// código NÃO está no nome da spec (`Login.spec.md`), então o template do derivado
	// precisa declará-la literalmente — {{ext}} aqui seria vazio.
	cfg := &config.Config{
		Version: 1,
		Layers: map[string]config.Layer{
			"spec":   {Kind: "spec", Tags: []string{"spec"}},
			"screen": {Kind: "code", Tags: []string{"frontend"}},
			"test":   {Kind: "test"},
		},
		Derived: &config.Derived{
			Anchor: "spec",
			Files: map[string]config.Padroes{
				"code": {"{{dir}}/{{name}}.ts"},
				"test": {"{{dir}}/{{name}}.test.ts"},
			},
		},
	}
	files := []scan.File{
		{Path: "src/Login.spec.md", Layer: "spec", Kind: "spec", Rev: "a"},
		{Path: "src/Login.ts", Layer: "screen", Kind: "code", Rev: "b"},
		{Path: "src/Login.test.ts", Layer: "test", Kind: "test", Rev: "c"},
	}
	g := Build(files, cfg, nil)

	if !hasEdge(g, "src/Login.spec.md", "src/Login.ts", EdgeSpecifies) {
		t.Error("a spec âncora deveria especificar o código derivado")
	}
	// Sem feature declarada, o teste deriva direto da spec — do contrário um projeto que
	// não usa feature perderia a ligação spec→teste inteira.
	if !hasEdge(g, "src/Login.spec.md", "src/Login.test.ts", EdgeTestedBy) {
		t.Error("sem feature, a spec deveria ligar direto ao teste")
	}
	// E a direção NÃO se inverte: o código não especifica a spec.
	if hasEdge(g, "src/Login.ts", "src/Login.spec.md", EdgeSpecifies) {
		t.Error("o código não pode especificar a spec — a spec é a âncora")
	}
}
