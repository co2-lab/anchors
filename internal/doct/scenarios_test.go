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

// O CORPO do cenário entra na doc, não só o título.
//
// É a diferença entre esta leitura e a do gate: o gate confronta o cenário contra o teste e
// só precisa do código. A documentação precisa dos passos — o título sozinho é manchete.
func TestScenarios_capturaOCorpo(t *testing.T) {
	root, g := projetoComFeature(t)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	specs, _ := c.fnSpecs("layer=infra")
	cs := c.fnScenarios(specs[0])

	if len(cs) != 2 {
		t.Fatalf("achou %d cenários, queria 2", len(cs))
	}
	if cs[0].Code != "GLCGL-B01" || cs[0].Titulo != "item sem artefato não passa" {
		t.Errorf("primeiro cenário = %+v", cs[0])
	}
	if !strings.Contains(cs[0].Corpo, "Então a régua não libera") {
		t.Errorf("o corpo não veio: %q", cs[0].Corpo)
	}
	if strings.Contains(cs[0].Corpo, "Dado uma dívida") {
		t.Errorf("o corpo invadiu o cenário seguinte: %q", cs[0].Corpo)
	}
}

// O `Contexto:` não é do cenário anterior: seus passos não podem vazar para dentro dele.
func TestScenarios_contextoNaoVazaParaOCenarioAnterior(t *testing.T) {
	feat := `Funcionalidade: X

  @ABCDE-B01 @nivel-unit
  Cenário: primeiro
    Dado A

  Contexto:
    Dado que isto é preparação
`
	cs := parseScenarios(feat, "ABCDE")
	if len(cs) != 1 {
		t.Fatalf("achou %d, queria 1", len(cs))
	}
	if strings.Contains(cs[0].Corpo, "preparação") {
		t.Errorf("o Contexto entrou no cenário: %q", cs[0].Corpo)
	}
}

// A tag de REGIME não é confundida com o código de identidade.
func TestScenarios_separaCodigoDeTag(t *testing.T) {
	cs := parseScenarios(featureDeExemplo, "GLCGL")
	if len(cs[0].Tags) != 1 || cs[0].Tags[0] != "nivel-unit" {
		t.Errorf("tags = %v, queria [nivel-unit]", cs[0].Tags)
	}
}

// Feature em INGLÊS é lida igual. Uma documentação que ignora metade dos cenários por
// causa do idioma do arquivo não é a documentação de ninguém.
func TestScenarios_leFeatureEmIngles(t *testing.T) {
	feat := "Feature: X\n\n  @ABCDE-B01\n  Scenario: it works\n    Given a thing\n"
	cs := parseScenarios(feat, "ABCDE")
	if len(cs) != 1 || cs[0].Titulo != "it works" {
		t.Fatalf("cenários = %+v", cs)
	}
}

// A feature é achada pelo MAPA, não por convenção de nome — que quebraria no primeiro
// projeto que organizasse os arquivos de outro jeito.
func TestScenarios_achaAFeaturePelaAresta(t *testing.T) {
	root, g := projetoComFeature(t)
	// A feature muda de nome; a aresta acompanha.
	os.Rename(filepath.Join(root, "pkg/GoLive.feature"), filepath.Join(root, "pkg/Outro.feature"))
	g.Nodes[1].ID = "pkg/Outro.feature"
	g.Edges[0].To = "pkg/Outro.feature"

	c, _ := New(root, g)
	specs, _ := c.fnSpecs("layer=infra")
	if cs := c.fnScenarios(specs[0]); len(cs) != 2 {
		t.Errorf("achou %d cenários com a feature renomeada, queria 2", len(cs))
	}
}

// O ESQUELETO compila. Um scaffold que o próprio compilador não processa entregaria ao
// time um erro de sintaxe como ponto de partida.
func TestInitScaffolds_oEsqueletoCompila(t *testing.T) {
	root, g := projetoComFeature(t)
	c, _ := New(root, g)

	escritos, _, err := c.InitScaffolds(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(escritos) < 4 {
		t.Fatalf("escreveu %d templates: %v", len(escritos), escritos)
	}

	res, err := c.Build(false)
	if err != nil {
		t.Fatalf("o esqueleto que o Anchors gera não compila: %v", err)
	}
	if len(res.Written) != len(escritos) {
		t.Errorf("gerou %d docs de %d templates", len(res.Written), len(escritos))
	}

	// O CONTEÚDO entra na página da CAMADA — é ela que o índice promete.
	b, err := os.ReadFile(filepath.Join(root, OutDir, "camadas", "infra.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "Então a régua não libera") {
		t.Errorf("o corpo do cenário não entrou na página da camada:\n%s", b)
	}

	// E o índice APONTA para lá, com uma âncora que existe. Um link que não resolve é o
	// pior defeito de um índice: ele é clicável, e o navegador fica onde está.
	idx, err := os.ReadFile(filepath.Join(root, OutDir, "comportamento.md"))
	if err != nil {
		t.Fatal(err)
	}
	alvo := "camadas/infra.md#" + GitHubAnchor("GLCGL-B01 — item sem artefato não passa")
	if !strings.Contains(string(idx), alvo) {
		t.Errorf("o índice não aponta para `%s`:\n%s", alvo, idx)
	}
	if !strings.Contains(string(b), "#### GLCGL-B01 — item sem artefato não passa") {
		t.Errorf("a âncora do índice não existe na página de destino:\n%s", b)
	}
}

// Rodar de novo não apaga o que o time editou: o template já editado carrega a moldura
// escrita à mão, que é justamente a parte não gerada.
func TestInitScaffolds_naoSobrescreveSemForce(t *testing.T) {
	root, g := projetoComFeature(t)
	c, _ := New(root, g)
	c.InitScaffolds(false)

	alvo := filepath.Join(root, Dir, "arquitetura.md"+SufixoTemplate)
	editado := "{{/* editado pelo time */}}\n# Nossa arquitetura\n"
	os.WriteFile(alvo, []byte(editado), 0o644)

	_, pulados, err := c.InitScaffolds(false)
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(alvo); string(b) != editado {
		t.Fatalf("o init apagou a moldura escrita à mão:\n%s", b)
	}
	if len(pulados) == 0 {
		t.Error("pulou em silêncio — quem roda de novo procura por que nada mudou")
	}
}

// TEMPLATE EM SUBPASTA compila. O `os.ReadDir` raso pulava `doct/camadas/*.tmpl` em
// silêncio — metade da matriz sumia da documentação com o build verde, que é exatamente o
// modo de falha que este mecanismo existe para fechar.
func TestBuild_desceEmSubpasta(t *testing.T) {
	root, g := projetoComFeature(t)
	sub := filepath.Join(root, Dir, "camadas")
	os.MkdirAll(sub, 0o755)
	os.WriteFile(filepath.Join(sub, "infra.md"+SufixoTemplate),
		[]byte(`{{range specs "layer=infra"}}{{.Code}}{{end}}`), 0o644)

	c, _ := New(root, g)
	res, err := c.Build(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Written) != 1 || res.Written[0] != "camadas/infra.md" {
		t.Fatalf("escritos = %v — o template de subpasta foi pulado", res.Written)
	}
	b, err := os.ReadFile(filepath.Join(root, OutDir, "camadas", "infra.md"))
	if err != nil {
		t.Fatalf("nada em docs/camadas/infra.md: %v", err)
	}
	if !strings.Contains(string(b), "GLCGL") {
		t.Errorf("conteúdo:\n%s", b)
	}
	// E o `Stale` enxerga o mesmo conjunto: um gate cego às subpastas nunca acusaria a
	// defasagem das páginas de camada.
	if s, _ := c.Stale(); len(s) != 0 {
		t.Errorf("acabou de compilar e o Stale acusa %v", s)
	}
}

// Todo scaffold tem NOME e CORPO. Um `Nome` vazio faz o `init` tentar escrever sobre o
// próprio diretório `doct/` — e a mensagem ("is a directory") não diz qual scaffold está
// quebrado, o que transforma um erro de digitação numa caçada.
func TestScaffolds_todosTemNomeECorpo(t *testing.T) {
	todos := append(Scaffolds(), ScaffoldLayer("infra"))
	for i, s := range todos {
		if strings.TrimSpace(s.Nome) == "" {
			t.Errorf("scaffold #%d sem nome (Porque: %q)", i, s.Porque)
		}
		if !strings.HasSuffix(s.Nome, SufixoTemplate) {
			t.Errorf("scaffold %q não termina em %s", s.Nome, SufixoTemplate)
		}
		if strings.TrimSpace(s.Corpo) == "" {
			t.Errorf("scaffold %q sem corpo", s.Nome)
		}
		if strings.TrimSpace(s.Porque) == "" {
			t.Errorf("scaffold %q não diz o que a página responde", s.Nome)
		}
	}
}
