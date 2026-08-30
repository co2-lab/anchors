package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// honraFixture: spec → feature → teste (dois saltos, porque `tested-by` nasce na
// FEATURE). `testSrc` é o conteúdo do teste, o consumidor do handle.
func honraFixture(t *testing.T, testSrc string) (mapx.Node, *mapx.Graph, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.test.tsx"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	// O CÓDIGO e a aresta `specifies` entram no fixture porque o gate único confronta a
	// trinca inteira: sem a ponta que EXPÕE o handle ele responde "spec sem código
	// ligado" e pula, e o teste mediria o Skip em vez do que quer medir. O gate antigo
	// (`testid-honored`) dispensava esta ponta — só perguntava se ALGUÉM consulta.
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}, {ID: "x.feature", Kind: mapx.KindFeature}, {ID: "x.test.tsx", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{
			{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies},
			{From: "x.spec.md", To: "x.feature", Type: mapx.EdgeCoveredBy},
			{From: "x.feature", To: "x.test.tsx", Type: mapx.EdgeTestedBy},
		},
	}
	return spec, g, root
}

func TestTestIDHonored_idConsultadoPassa(t *testing.T) {
	n, g, root := honraFixture(t, `getByTestId(':abcd-screen')`)
	if v, msg := checkTestIDCoerente(secaoOK, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("id consultado pelo teste deveria passar: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_idOrfaoReprova(t *testing.T) {
	// Declarado, exposto e ninguém consulta: custo sem contrapartida. E pior que
	// inútil — parece cobertura, porque a superfície está lá e o inventário completo.
	n, g, root := honraFixture(t, `render(<X />)`)
	v, msg := checkTestIDCoerente(secaoOK, n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("id que ninguém consulta deveria reprovar: %v", v)
	}
	if !strings.Contains(msg, "abcd-screen") {
		t.Errorf("a mensagem deve nomear o id órfão: %s", msg)
	}
}

func TestTestIDHonored_prefixoDinamicoBastaOPrefixo(t *testing.T) {
	// O flow casa `abcd-item-.*` ou constrói `abcd-item-${id}`; o sufixo é dado de
	// runtime. Exigir o id completo tornaria todo template um falso órfão.
	n, g, root := honraFixture(t, "getByTestId(`:abcd-item-${id}`)")
	// O código tem de EXPOR o template — o fixture padrão expõe `:abcd-screen`, que é
	// outro handle. Sem isto o caso mediria "declarado e não exposto", não o que quer.
	if err := os.WriteFile(filepath.Join(root, "x.tsx"),
		[]byte("<View testID={`:abcd-item-${id}`} />"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := "## Superfície de Teste\n\n- `:abcd-item-*`\n"
	if v, msg := checkTestIDCoerente(spec, n, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("prefixo consultado honra o contrato dinâmico: %v (%s)", v, msg)
	}
}

func TestTestIDCoerente_semInventarioAcusaOExposto(t *testing.T) {
	// Este teste MUDOU de expectativa junto com a fusão dos gates, e a mudança é o
	// ponto: no arranjo antigo o `testid-honored` pulava quando a spec não tinha
	// inventário, para não duplicar o achado que o `testid-declared` já daria. Com um
	// gate só não há o que duplicar — e pular deixaria escapar exatamente o caso em que
	// o código expõe handle nenhum declarado, que é a superfície não-contratada.
	n, g, root := honraFixture(t, `render(<X />)`)
	v, msg := checkTestIDCoerente("", n, root, g, cfgHandle("testID"))
	if v != Fail {
		t.Fatalf("spec sem inventário e código expondo handle deve reprovar: %v (%s)", v, msg)
	}
	if !strings.Contains(msg, ":abcd-screen") {
		t.Errorf("o laudo deve nomear o handle exposto e não declarado; veio: %s", msg)
	}
}

func TestTestIDHonored_semHandleDeclaradoPula(t *testing.T) {
	n, g, root := honraFixture(t, `render(<X />)`)
	if v, _ := checkTestIDCoerente(secaoOK, n, root, g, &config.Config{}); v != Skip {
		t.Errorf("sem test_handle o gate deve pular: %v", v)
	}
}

func TestTestIDHonored_flowE2EContaComoConsumidor(t *testing.T) {
	// O flow de ponta a ponta é o consumidor PRINCIPAL do handle e vive FORA do
	// grafo. Se só o teste unitário contasse, todo id usado apenas pelo E2E seria
	// reportado como órfão — acusando exatamente quem cumpre o contrato.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.test.tsx"), []byte(`render(<X />)`), 0o644); err != nil {
		t.Fatal(err)
	}
	flows := filepath.Join(root, "e2e")
	if err := os.MkdirAll(flows, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(flows, "f.yaml"), []byte("- tapOn:\n    id: ':abcd-screen'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// O código e a aresta `specifies` entram porque o gate único confronta a trinca
	// inteira — sem a ponta que EXPÕE, ele pula e o teste mediria o Skip.
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}, {ID: "x.feature", Kind: mapx.KindFeature}, {ID: "x.test.tsx", Kind: mapx.KindTest}},
		Edges: []mapx.Edge{
			{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies},
			{From: "x.spec.md", To: "x.feature", Type: mapx.EdgeCoveredBy},
			{From: "x.feature", To: "x.test.tsx", Type: mapx.EdgeTestedBy},
		},
	}
	// A resolução é em dois passos: `surfaces[e2e]` dá a CHAVE, `files[chave]` dá o
	// caminho. Declarar só `surfaces` (como o app de referência faz hoje) não basta — e o gate tem
	// de tratar isso como "não sei onde procurar", não como "ninguém consome".
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "code", TestHandle: "testID",
		Surfaces: map[string]string{"e2e": "e2e"},
		Files:    map[string]config.Padroes{"e2e": {"e2e/{{name}}.yaml"}},
	}}
	if v, msg := checkTestIDCoerente(secaoOK, spec, root, g, cfg); v != Pass {
		t.Errorf("id usado só pelo flow e2e não é órfão: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_surfaceSemFilesNaoInventaCaminho(t *testing.T) {
	// A armadilha que este teste trava: `surfaces[e2e]` devolve a CHAVE da superfície,
	// não um path. Tratá-la como diretório faz o gate varrer um caminho inexistente,
	// achar zero consumidores e acusar de órfão todo id que só o E2E usa — ou seja,
	// reprovar quem cumpre o contrato. Sem `files[e2e]` não há onde procurar, e a
	// resposta honesta é não afirmar nada sobre a superfície que não se sabe ler.
	n, g, root := honraFixture(t, `getByTestId(':abcd-screen')`)
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "code", TestHandle: "testID",
		Surfaces: map[string]string{"e2e": "e2e"}, // sem files["e2e"]
	}}
	if v, msg := checkTestIDCoerente(secaoOK, n, root, g, cfg); v != Pass {
		t.Errorf("superfície sem caminho declarado não deve virar acusação: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_testeCompartilhadoContaComoConsumidor(t *testing.T) {
	// `PhiScreens.test.tsx` prova 6 telas. A aresta `tested-by` liga a feature ao
	// arquivo, mas nem toda spec do grupo o alcança pelo grafo — e o handle, apesar
	// de consultado, apareceria como órfão. O teste vizinho recupera o caso.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Compartilhado.test.tsx"),
		[]byte(`getByTestId(':abcd-screen')`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	// O código entra (o gate único exige a ponta que EXPÕE), mas a aresta para o TESTE
	// continua ausente — é justamente o que este caso mede: o consumidor alcançado por
	// vizinhança, não pelo grafo.
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies}},
	} // sem aresta p/ teste — é o ponto
	if v, msg := checkTestIDCoerente(secaoOK, spec, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("teste vizinho que consulta o id honra o contrato: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_vizinhoValeMesmoComFlows(t *testing.T) {
	// A guarda `len(consumidores)==0` anulava a rede de segurança: um projeto que
	// declara superfície e2e sempre devolve conteúdo dos flows, então a lista nunca
	// ficava vazia e o teste vizinho nunca era lido. O gate acusava 14 ids que os
	// testes compartilhados consultam.
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Compartilhado.test.tsx"),
		[]byte(`getByTestId(':abcd-screen')`), 0o644); err != nil {
		t.Fatal(err)
	}
	flows := filepath.Join(root, "e2e")
	if err := os.MkdirAll(flows, 0o755); err != nil {
		t.Fatal(err)
	}
	// Flow que existe mas NÃO cita o id — é o que mantinha `consumidores` não-vazio.
	if err := os.WriteFile(filepath.Join(flows, "f.yaml"), []byte("- tapOn:\n    id: ':outro'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	// O código entra (o gate único exige a ponta que EXPÕE), mas a aresta para o TESTE
	// continua ausente — é justamente o que este caso mede: o consumidor alcançado por
	// vizinhança, não pelo grafo.
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies}},
	}
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "code", TestHandle: "testID",
		Surfaces: map[string]string{"e2e": "e2e"},
		Files:    map[string]config.Padroes{"e2e": {"e2e/{{name}}.yaml"}},
	}}
	if v, msg := checkTestIDCoerente(secaoOK, spec, root, g, cfg); v != Pass {
		t.Errorf("teste vizinho deve valer mesmo havendo flows: %v (%s)", v, msg)
	}
}

func TestTestIDHonored_componenteExercitadoPelaTela(t *testing.T) {
	// `CategoryCard` vive em `trends/components/` e quem o exercita é
	// `TrendsScreen.test.tsx`, em `trends/screens/` — a tela que o renderiza. Sem
	// subir ao módulo, o gate acusa de órfão um handle que o teste consulta.
	root := t.TempDir()
	comp := filepath.Join(root, "trends", "components")
	tela := filepath.Join(root, "trends", "screens")
	for _, d := range []string{comp, tela} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(tela, "TrendsScreen.test.tsx"),
		[]byte(`getByTestId(':abcd-screen')`), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "trends/components/CategoryCard.spec.md", Kind: mapx.KindSpec, Code: "ABCDX"}
	// O código entra (o gate único exige a ponta que EXPÕE) e a aresta parte do ID REAL
	// da spec — que aqui tem caminho. A aresta para o TESTE segue ausente: é o que este
	// caso mede, o consumidor alcançado por vizinhança e não pelo grafo.
	cod := filepath.Join("trends", "components", "CategoryCard.tsx")
	if err := os.WriteFile(filepath.Join(root, cod), []byte(`<View testID=":abcd-screen" />`), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: cod, Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: spec.ID, To: cod, Type: mapx.EdgeSpecifies}},
	}
	if v, msg := checkTestIDCoerente(secaoOK, spec, root, g, cfgHandle("testID")); v != Pass {
		t.Errorf("teste da tela irmã exercita o componente: %v (%s)", v, msg)
	}
}
