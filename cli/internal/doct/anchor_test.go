package doct

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// A âncora é o que faz o índice funcionar. Errar aqui produz um link clicável que não vai
// a lugar nenhum — o navegador fica onde está, e o índice PARECE funcionar.
func TestGitHubAnchor(t *testing.T) {
	casos := []struct{ heading, esperado string }{
		{"GLCGL-B01 — Cada item tem um artefato que o PROVA",
			"glcgl-b01--cada-item-tem-um-artefato-que-o-prova"},
		{"GLCGL — GoLiveChecklist — a régua que registra o que não foi feito",
			"glcgl--golivechecklist--a-régua-que-registra-o-que-não-foi-feito"},
		// Os ACENTOS FICAM: o GitHub preserva `é` na âncora, e transliterar produziria
		// um link que não resolve. É o erro que uma biblioteca de slug faria por padrão.
		{"A régua e a dívida", "a-régua-e-a-dívida"},
		{"Visão Geral", "visão-geral"},
		// Pontuação some; o espaço que a cercava não vira hífen duplo por acaso — vira,
		// e é o comportamento do GitHub.
		{"O que é, e o que não é", "o-que-é-e-o-que-não-é"},
		{"`código` entre crases", "código-entre-crases"},
		{"  espaços nas pontas  ", "espaços-nas-pontas"},
		{"Camada: infra", "camada-infra"},
	}
	for _, c := range casos {
		if got := GitHubAnchor(c.heading); got != c.esperado {
			t.Errorf("GitHubAnchor(%q)\n  = %q\n  queria %q", c.heading, got, c.esperado)
		}
	}
}

// TODO LINK GERADO TEM DE RESOLVER. Este é o teste que faltava quando 483 de 812 links
// saíram quebrados: o índice montava a âncora da regra sempre, e nas camadas grandes a
// página de destino resume — a regra não é heading lá, e a âncora não existe.
//
// Todos eram clicáveis, e todos paravam no mesmo lugar. Um índice assim é pior que não ter
// índice, porque parece funcionar.
func TestBuild_todoLinkGeradoResolve(t *testing.T) {
	// Uma camada PEQUENA (página completa) e uma GRANDE (página resumida): o defeito só
	// aparece na segunda, e um teste com uma camada só passaria em falso.
	specs := map[string]string{}
	for i := 0; i < DefaultLayout().MaxUnits+2; i++ {
		code := fmt.Sprintf("BIG%02d", i)
		specs[fmt.Sprintf("big/U%02d.spec.md", i)] = fmt.Sprintf(`---
code: %s
layer: grande
---

# %s — a unidade %d

## Visão Geral

Texto.

## Regras

### %s-B01 — a regra da unidade %d

Corpo.
`, code, code, i, code, i)
	}
	specs["peq/Only.spec.md"] = specDeExemplo // layer: infra, com regras

	root, g := projetoDeTeste(t, specs)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := c.InitScaffolds(false); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Build(false); err != nil {
		t.Fatal(err)
	}

	// As duas camadas precisam existir, e em lados opostos do limiar — senão o teste
	// não exercita o caso que quebrou.
	if !c.layerIsBig("grande") {
		t.Fatal("a camada `grande` devia passar do limite — o teste não cobre o defeito")
	}
	if c.layerIsBig("infra") {
		t.Fatal("a camada `infra` devia ficar abaixo do limite")
	}

	headings, links := colheDocs(t, filepath.Join(root, OutDir))
	if len(links) == 0 {
		t.Fatal("nenhum link gerado — o teste não prova nada")
	}
	for _, l := range links {
		anc, ok := headings[l.destino]
		if !ok {
			t.Errorf("%s aponta para `%s`, que não existe", l.origem, l.destino)
			continue
		}
		if !anc[l.ancora] {
			t.Errorf("%s aponta para `%s#%s` — a âncora não existe no destino",
				l.origem, l.destino, l.ancora)
		}
	}
}

type linkGerado struct{ origem, destino, ancora string }

var (
	headingRE = regexp.MustCompile(`(?m)^#{1,6}\s+(.+?)\s*$`)
	linkRE    = regexp.MustCompile(`\]\(([^)#]+)#([^)]+)\)`)
)

// colheDocs devolve as âncoras de cada página e todos os links internos gerados.
func colheDocs(t *testing.T, dir string) (map[string]map[string]bool, []linkGerado) {
	t.Helper()
	headings := map[string]map[string]bool{}
	var links []linkGerado
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		rel, _ := filepath.Rel(dir, p)
		rel = filepath.ToSlash(rel)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		set := map[string]bool{}
		for _, m := range headingRE.FindAllStringSubmatch(string(b), -1) {
			set[GitHubAnchor(m[1])] = true
		}
		headings[rel] = set
		base := path.Dir(rel)
		for _, m := range linkRE.FindAllStringSubmatch(string(b), -1) {
			links = append(links, linkGerado{
				origem:  rel,
				destino: path.Clean(path.Join(base, m[1])),
				ancora:  m[2],
			})
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return headings, links
}

// As SETAS do diagrama vêm do mapa, agregadas por camada.
func TestLayerDeps_agregaPorCamada(t *testing.T) {
	specs := map[string]string{
		"a/S1.spec.md": "---\ncode: SCRNA\nlayer: screen\n---\n\n# S1 — tela um\n",
		"a/S2.spec.md": "---\ncode: SCRNB\nlayer: screen\n---\n\n# S2 — tela dois\n",
		"b/H1.spec.md": "---\ncode: SHRDA\nlayer: shared\n---\n\n# H1 — util\n",
	}
	root, g := projetoDeTeste(t, specs)
	// Duas telas dependem do mesmo util: uma seta, peso 2.
	g.Edges = []mapx.Edge{
		{From: "a/S1.spec.md", To: "b/H1.spec.md", Type: "needs"},
		{From: "a/S2.spec.md", To: "b/H1.spec.md", Type: "needs"},
		// Aresta que sai das specs não vira seta de camada: um plano não é camada, e
		// desenhá-lo produziria uma seta de um nó que a página não mostra.
		{From: "plans/0001.md", To: "a/S1.spec.md", Type: "needs"},
		// Aresta de outro tipo também não: `tested-by` é a trinca, não dependência.
		{From: "a/S1.spec.md", To: "b/H1.spec.md", Type: "tested-by"},
	}
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}

	deps := c.fnLayerDeps()
	if len(deps) != 1 {
		t.Fatalf("deps = %+v, queria 1 (screen→shared)", deps)
	}
	if deps[0].From != "screen" || deps[0].To != "shared" {
		t.Errorf("seta = %s→%s", deps[0].From, deps[0].To)
	}
	// O PESO distingue a dependência única da sistêmica — sem ele as duas pareceriam
	// a mesma coisa no diagrama.
	if deps[0].Count != 2 {
		t.Errorf("peso = %d, queria 2 (duas telas sustentam a seta)", deps[0].Count)
	}
}

// Dependência DENTRO da mesma camada não vira seta: um laço de um nó para si mesmo não
// informa nada sobre a arquitetura.
func TestLayerDeps_ignoraDependenciaInterna(t *testing.T) {
	specs := map[string]string{
		"a/S1.spec.md": "---\ncode: SCRNA\nlayer: screen\n---\n\n# S1 — um\n",
		"a/S2.spec.md": "---\ncode: SCRNB\nlayer: screen\n---\n\n# S2 — dois\n",
	}
	root, g := projetoDeTeste(t, specs)
	g.Edges = []mapx.Edge{{From: "a/S1.spec.md", To: "a/S2.spec.md", Type: "needs"}}
	c, _ := New(root, g)
	if deps := c.fnLayerDeps(); len(deps) != 0 {
		t.Errorf("deps = %+v — a dependência é interna à camada", deps)
	}
}

// O ID do Mermaid não pode ter hífen nem acento: o parser quebra, e o BLOCO INTEIRO deixa
// de renderizar — aparece como texto cru. Uma camada `feature-hook` bastaria.
func TestMermaidID(t *testing.T) {
	casos := map[string]string{
		"feature-hook": "n_feature_hook",
		"app-infra":    "n_app_infra",
		"lambdas":      "n_lambdas",
		"padrão":       "n_padr_o",
	}
	for in, esperado := range casos {
		if got := MermaidID(in); got != esperado {
			t.Errorf("MermaidID(%q) = %q, queria %q", in, got, esperado)
		}
	}
}
