package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// TestTargetObrigatorioQuandoDeclarado — o `run:` que usa {{target}} não pode rodar sem
// alvo. Substituir por vazio faria `npx stryker run --mutate ` mutar o PROJETO INTEIRO:
// horas de rodada, e nada do que se pediu. Falhar aqui custa um segundo.
func TestTargetObrigatorioQuandoDeclarado(t *testing.T) {
	_, err := montaComando("npx stryker run --mutate {{target}}", "")
	if err == nil {
		t.Fatal("sem --target, o comando com {{target}} tem de falhar")
	}
	if !strings.Contains(err.Error(), "--target") {
		t.Errorf("o erro precisa dizer o que fazer; veio %q", err)
	}
}

func TestTargetSubstitui(t *testing.T) {
	got, err := montaComando("stryker run --mutate {{target}}", "business-logic/dedup.ts")
	if err != nil {
		t.Fatal(err)
	}
	if got != "stryker run --mutate business-logic/dedup.ts" {
		t.Errorf("substituição errada: %q", got)
	}
}

// TestComandoSemPlaceholderIgnoraTarget — a maioria das suítes não recebe alvo (`yarn
// test:unit` roda tudo). Exigir --target nelas, ou anexá-lo ao fim, quebraria o comando.
func TestComandoSemPlaceholderIgnoraTarget(t *testing.T) {
	got, err := montaComando("yarn test:unit", "qualquer/coisa.ts")
	if err != nil {
		t.Fatal(err)
	}
	if got != "yarn test:unit" {
		t.Errorf("o comando sem placeholder não deve ser tocado; veio %q", got)
	}
}

// TestCaminhoDoRelatorioEhRelativoARaiz — o anchors.yaml declara caminho de projeto
// ("packages/backend/reports/x.json"), e o comando pode rodar de qualquer diretório.
// Resolver contra a raiz é o que faz `anchors mutation` funcionar de dentro de um
// subdiretório, como todo o resto do CLI já faz.
func TestCaminhoDoRelatorioEhRelativoARaiz(t *testing.T) {
	raiz := t.TempDir()
	got := caminhoAbs(raiz, "packages/backend/reports/mutation.json")
	if !filepath.IsAbs(got) {
		t.Fatalf("devia virar absoluto; veio %q", got)
	}
	if !strings.HasSuffix(filepath.ToSlash(got), "packages/backend/reports/mutation.json") {
		t.Errorf("o sufixo do caminho se perdeu: %q", got)
	}
}

// TestSemRelatorioNaoInventaCaminho — vazio tem de continuar vazio. Se virasse a própria
// raiz, a ingestão tentaria ler um diretório como se fosse relatório.
func TestSemRelatorioNaoInventaCaminho(t *testing.T) {
	if got := caminhoAbs(t.TempDir(), "  "); got != "" {
		t.Errorf("relatório não declarado tem de continuar vazio; veio %q", got)
	}
}

// TestEncadeamentoRecusaComandoDesconhecido — `--then` aceita só o que faz sentido
// depois de ingerir. Aceitar qualquer nome transformaria a flag num executor de shell
// disfarçado, que é justamente o que estes comandos NÃO são.
func TestEncadeamentoRecusaComandoDesconhecido(t *testing.T) {
	err := executaEncadeados("deploy", t.TempDir())
	if err == nil {
		t.Fatal("--then deploy tinha de ser recusado")
	}
	if !strings.Contains(err.Error(), "check") {
		t.Errorf("o erro precisa dizer o que É aceito; veio %q", err)
	}
}

// TestEncadeamentoVazioNaoFazNada — o opt-in é o default. Sem `--then`, roda e para.
func TestEncadeamentoVazioNaoFazNada(t *testing.T) {
	if err := executaEncadeados("", t.TempDir()); err != nil {
		t.Errorf("sem --then não há encadeamento a fazer; veio %v", err)
	}
	if err := executaEncadeados("  ,  ", t.TempDir()); err != nil {
		t.Errorf("separadores vazios não são comando; veio %v", err)
	}
}

// TestSuiteQuePassaSemRelatorioNaoQuebra — declarar `run:` sem `junit:` é legítimo (um
// atalho), e não pode ser tratado como erro. O comando avisa que nada foi ingerido,
// porque o silêncio aqui faria o usuário achar que o gate ia mudar de cor.
func TestSuiteQuePassaSemRelatorioNaoQuebra(t *testing.T) {
	raiz := t.TempDir()
	marca := filepath.Join(raiz, "rodou.txt")
	cs := comandoSuite{nome: "test", secao: "tests"}
	s := []config.Suite{{Layer: "unit", Run: "printf ok > " + filepath.ToSlash(marca)}}

	if err := rodaSuites(cs, s, raiz, "", nil); err != nil {
		t.Fatalf("suíte sem relatório não devia falhar: %v", err)
	}
	if _, err := os.Stat(marca); err != nil {
		t.Errorf("o comando declarado não rodou: %v", err)
	}
}

// TestSuiteQueFalhaInterrompe — parar na primeira falha é a regra: as camadas dependem
// umas das outras, e rodar e2e sobre uma unit vermelha só produz ruído sobre uma base
// já quebrada.
func TestSuiteQueFalhaInterrompe(t *testing.T) {
	raiz := t.TempDir()
	depois := filepath.Join(raiz, "nao-devia-existir.txt")
	cs := comandoSuite{nome: "test", secao: "tests"}
	s := []config.Suite{
		{Layer: "unit", Run: "exit 3"},
		{Layer: "e2e", Run: "printf x > " + filepath.ToSlash(depois)},
	}

	err := rodaSuites(cs, s, raiz, "", nil)
	if err == nil {
		t.Fatal("a falha da unit tinha de interromper")
	}
	if !strings.Contains(err.Error(), "unit") {
		t.Errorf("o erro precisa nomear a camada que falhou; veio %q", err)
	}
	if _, serr := os.Stat(depois); serr == nil {
		t.Error("a camada seguinte não podia ter rodado")
	}
}

// ── os dois modos: completo e incremental ──────────────────────────────────────

// TestSemChangedUsaOComandoCompleto — o default é a rodada completa, como no `check`.
func TestSemChangedUsaOComandoCompleto(t *testing.T) {
	s := config.Suite{Run: "jest", RunChanged: "jest --findRelatedTests {{files}}"}
	got, err := escolheComando(s, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "jest" {
		t.Errorf("sem --changed roda o completo; veio %q", got)
	}
}

// TestChangedUsaOComandoIncremental — e os arquivos entram onde o projeto declarou,
// não anexados no fim: a posição do recorte é do comando, não nossa.
func TestChangedUsaOComandoIncremental(t *testing.T) {
	s := config.Suite{Run: "jest", RunChanged: "jest --findRelatedTests {{files}} --ci"}
	got, err := escolheComando(s, []string{"a/x.ts", "b/y.ts"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "jest --findRelatedTests a/x.ts b/y.ts --ci" {
		t.Errorf("substituição de {{files}} errada: %q", got)
	}
}

// TestChangedSemRunChangedRecusa guarda a decisão mais importante dos dois modos: cair
// para a rodada COMPLETA seria caro e, pior, mentiria — o usuário leria "passou"
// achando que rodou o recorte dele.
func TestChangedSemRunChangedRecusa(t *testing.T) {
	s := config.Suite{Run: "yarn test:unit"}
	_, err := escolheComando(s, []string{"a/x.ts"}, "")
	if err == nil {
		t.Fatal("sem run_changed, o modo incremental tem de recusar")
	}
	if !strings.Contains(err.Error(), "run_changed") || !strings.Contains(err.Error(), "{{files}}") {
		t.Errorf("o erro precisa mostrar o que declarar; veio %q", err)
	}
}

// TestMutacaoNaoMutaTeste — mutação altera a REGRA. Mutar o teste inverteria o
// experimento: o teste é o instrumento de medida, não o objeto medido.
func TestMutacaoNaoMutaTeste(t *testing.T) {
	nodes := []mapx.Node{
		{ID: "a/regra.ts", Kind: mapx.KindCode},
		{ID: "a/regra.test.ts", Kind: mapx.KindTest},
		{ID: "a/regra.spec.md", Kind: mapx.KindSpec},
		{ID: "a/regra.feature", Kind: mapx.KindFeature},
	}
	if got := arquivosDoImpacto(nodes, "mutation", ""); len(got) != 1 || got[0] != "a/regra.ts" {
		t.Errorf("mutação recebe só código; veio %v", got)
	}
}

// TestTesteRecebeCodigoETeste — `--findRelatedTests` e equivalentes esperam arquivos
// de fonte; spec e feature não significam nada para um runner.
func TestTesteRecebeCodigoETeste(t *testing.T) {
	nodes := []mapx.Node{
		{ID: "a/regra.ts", Kind: mapx.KindCode},
		{ID: "a/regra.test.ts", Kind: mapx.KindTest},
		{ID: "a/regra.spec.md", Kind: mapx.KindSpec},
	}
	got := arquivosDoImpacto(nodes, "tests", "")
	if len(got) != 2 || got[0] != "a/regra.ts" || got[1] != "a/regra.test.ts" {
		t.Errorf("teste recebe código e teste, sem spec; veio %v", got)
	}
}

// TestCaminhoDeImpactoNaoLevaBarraInvertida — o comando roda dentro de `sh -c`, onde a
// barra invertida é ESCAPE. Medido: com o separador nativo o caminho chega desfeito ao
// runner e a rodada acha 0 teste — sem erro, só um recorte vazio que se lê como "não
// havia o que rodar". É o modo de falha mais caro: silencioso e otimista.
func TestCaminhoDeImpactoNaoLevaBarraInvertida(t *testing.T) {
	raiz := t.TempDir()
	got := arquivosDoImpacto([]mapx.Node{{ID: "packages/backend/x.ts", Kind: mapx.KindCode}}, "tests", raiz)
	if len(got) != 1 {
		t.Fatalf("esperava 1 arquivo; veio %v", got)
	}
	if strings.Contains(got[0], `\`) {
		t.Errorf("caminho passado ao shell não pode levar barra invertida: %q", got[0])
	}
	if !strings.HasSuffix(got[0], "packages/backend/x.ts") {
		t.Errorf("o caminho perdeu o sufixo: %q", got[0])
	}
}
