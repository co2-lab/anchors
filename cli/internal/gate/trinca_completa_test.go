package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// grafo mínimo: spec → código (specifies), spec → feature (covered-by),
// feature → teste (tested-by). O teste se alcança em DOIS saltos a partir da spec.
func trincaGraph(withCode, withFeature, withTest bool) *mapx.Graph {
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "x.spec.md", Kind: mapx.KindSpec},
			{ID: "x.ts", Kind: mapx.KindCode},
			{ID: "x.feature", Kind: mapx.KindFeature},
			{ID: "x.test.ts", Kind: mapx.KindTest},
		},
	}
	if withCode {
		g.Edges = append(g.Edges, mapx.Edge{From: "x.spec.md", To: "x.ts", Type: mapx.EdgeSpecifies})
	}
	if withFeature {
		g.Edges = append(g.Edges, mapx.Edge{From: "x.spec.md", To: "x.feature", Type: mapx.EdgeCoveredBy})
	}
	if withTest {
		g.Edges = append(g.Edges, mapx.Edge{From: "x.feature", To: "x.test.ts", Type: mapx.EdgeTestedBy})
	}
	return g
}

func specNode() mapx.Node { return mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec} }

// raizComProva escreve o teste `x.test.ts` mencionando `codigo`, para que a referência
// verificável do `@no-test` RESOLVA. Sem isso, todo teste que usa `@no-test` reprovaria
// pela referência órfã — que é justamente o que o gate passou a cobrar.
func raizComProva(t *testing.T, codigo string) string {
	t.Helper()
	root := t.TempDir()
	conteudo := "it('" + codigo + ": prova o comportamento', () => {})\n"
	if err := os.WriteFile(filepath.Join(root, "x.test.ts"), []byte(conteudo), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestTrincaCompleta_trincaInteiraPassa(t *testing.T) {
	g := trincaGraph(true, true, true)
	if v, msg := checkTrincaCompleta("", specNode(), "", g, &config.Config{}); v != Pass {
		t.Errorf("trinca completa deveria passar: %v (%s)", v, msg)
	}
}

func TestTrincaCompleta_testeAlcancadoEmDoisSaltos(t *testing.T) {
	// A regressão que eu mesmo introduzi: procurar `tested-by` DIRETO na spec acusa
	// falta de teste em todo projeto, porque a aresta nasce na FEATURE.
	g := trincaGraph(true, true, true)
	v, msg := checkTrincaCompleta("", specNode(), "", g, &config.Config{})
	if v != Pass {
		t.Fatalf("teste ligado à feature deveria contar p/ a spec: %v (%s)", v, msg)
	}
}

func TestTrincaCompleta_semTesteReprova(t *testing.T) {
	g := trincaGraph(true, true, false)
	v, msg := checkTrincaCompleta("", specNode(), "", g, &config.Config{})
	if v != Fail {
		t.Fatalf("sem teste deveria reprovar, got %v", v)
	}
	if !strings.Contains(msg, "teste") {
		t.Errorf("mensagem deveria citar o que falta: %q", msg)
	}
}

func TestTrincaCompleta_specSozinhaReprovaCitandoAsTresPecas(t *testing.T) {
	// o caso que motivou o gate: spec sem nada atravessava TODOS os gates.
	g := trincaGraph(false, false, false)
	v, msg := checkTrincaCompleta("", specNode(), "", g, &config.Config{})
	if v != Fail {
		t.Fatalf("spec sozinha deveria reprovar, got %v", v)
	}
	for _, peca := range []string{"código", "feature", "teste"} {
		if !strings.Contains(msg, peca) {
			t.Errorf("mensagem deveria citar %q: %q", peca, msg)
		}
	}
}

func TestTrincaCompleta_camadaReconhecidaPula(t *testing.T) {
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Tags: []string{"dao"}}
	g := trincaGraph(false, false, false)
	if v, _ := checkTrincaCompleta("", n, "", g, &config.Config{}); v != Skip {
		t.Errorf("camada reconhecida não tem trinca a cobrar, esperava Skip, got %v", v)
	}
}

func TestTrincaCompleta_trincaOpcionalDispensaPeca(t *testing.T) {
	// opt-out HONESTO: declarado na Estrutura, não escondido num Skip.
	cfg := &config.Config{Layers: map[string]config.Layer{
		"repository": {TrincaOpcional: []string{"tested-by"}},
	}}
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Tags: []string{"repository"}}
	g := trincaGraph(true, true, false) // sem teste
	if v, msg := checkTrincaCompleta("", n, "", g, cfg); v != Pass {
		t.Errorf("camada que dispensa tested-by deveria passar: %v (%s)", v, msg)
	}
}

func TestTrincaCompleta_naoSpecPula(t *testing.T) {
	n := mapx.Node{ID: "x.ts", Kind: mapx.KindCode}
	if v, _ := checkTrincaCompleta("", n, "", trincaGraph(false, false, false), &config.Config{}); v != Skip {
		t.Errorf("a trinca é cobrada da spec, esperava Skip p/ código, got %v", v)
	}
}

func TestTrincaCompleta_noTestDispensaPorUnidade(t *testing.T) {
	// A dispensa por CAMADA isenta em bloco; esta é da UNIDADE e fica escrita nela.
	// Serve para o caso real: dentro de `services` convivem o gateway de 9 linhas
	// que só repassa a chamada e o módulo de 150 com regra — isentar os dois junto
	// apagaria a cobrança justamente onde ela vale.
	g := trincaGraph(true, true, false) // sem teste
	root := raizComProva(t, "AAAAX-B01")
	spec := "@no-test: gateway de 1 linha sobre apiCall; provado por `AAAAX-B01`\n"
	if v, msg := checkTrincaCompleta(spec, specNode(), root, g, &config.Config{}); v != Pass {
		t.Errorf("`@no-test` com razão dispensa o teste daquela unidade: %v (%s)", v, msg)
	}
}

func TestTrincaCompleta_dispensaExigeRazao(t *testing.T) {
	// Marcador NU não dispensa — senão `@no-test` viraria um jeito silencioso de
	// calar o gate, que é o oposto do opt-out honesto.
	g := trincaGraph(true, true, false)
	for _, nu := range []string{"@no-test\n", "@no-test:\n", "@no-test:   \n"} {
		if v, _ := checkTrincaCompleta(nu, specNode(), "", g, &config.Config{}); v != Fail {
			t.Errorf("marcador sem razão (%q) NÃO pode dispensar: %v", nu, v)
		}
	}
}

func TestTrincaCompleta_noFeatureArrastaOTeste(t *testing.T) {
	// Sem feature não há cenário a provar: cobrar o teste seria exigir a prova de
	// algo que ninguém especificou.
	g := trincaGraph(true, false, false)
	spec := "@no-feature: wiring de infraestrutura, sem regra observável\n"
	if v, msg := checkTrincaCompleta(spec, specNode(), "", g, &config.Config{}); v != Pass {
		t.Errorf("`@no-feature` dispensa feature E teste: %v (%s)", v, msg)
	}
}

func TestTrincaCompleta_dispensaNaoApagaOCodigo(t *testing.T) {
	// O CÓDIGO nunca é dispensável: uma spec sem o arquivo que ela descreve é a
	// própria situação que este gate existe para pegar.
	g := trincaGraph(false, true, true)
	spec := "@no-test: qualquer razão\n@no-feature: qualquer razão\n"
	if v, _ := checkTrincaCompleta(spec, specNode(), "", g, &config.Config{}); v != Fail {
		t.Errorf("a dispensa não pode apagar a exigência do código: %v", v)
	}
}

func TestTrincaCompleta_noTestComCenarioNaFeatureEhContradicao(t *testing.T) {
	// `@no-test` diz "não há o que provar"; um cenário diz "prova-se assim". As duas
	// afirmações não convivem — e sem esta checagem a contradição fica MUDA: a
	// dispensa satisfaz o trinca-completa e o feature-test-match só cobra quando há
	// teste a confrontar. Resultado: cenário escrito que ninguém prova, tudo verde.
	root := raizComProva(t, "AAAAX-B01")
	if err := os.WriteFile(filepath.Join(root, "x.feature"),
		[]byte("Funcionalidade: X\n\n  Cenário: faz algo\n    Dado a\n    Quando b\n    Então c\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := trincaGraph(true, true, false)
	spec := "@no-test: gateway de uma linha, provado por `AAAAX-B01`\n"

	v, msg := checkTrincaCompleta(spec, specNode(), root, g, &config.Config{})
	if v != Fail {
		t.Fatalf("dispensa + cenário é contradição e deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "@no-test") || !strings.Contains(msg, "cenário") {
		t.Errorf("a mensagem deve nomear a contradição: %s", msg)
	}
}

func TestTrincaCompleta_noTestSemCenarioPassa(t *testing.T) {
	// Feature com cabeçalho e NENHUM cenário (esqueleto) não contradiz a dispensa —
	// não há afirmação de comportamento a provar.
	root := raizComProva(t, "AAAAX-B01")
	if err := os.WriteFile(filepath.Join(root, "x.feature"),
		[]byte("Funcionalidade: X\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := trincaGraph(true, true, false)
	spec := "@no-test: gateway de uma linha, provado por `AAAAX-B01`\n"
	if v, msg := checkTrincaCompleta(spec, specNode(), root, g, &config.Config{}); v != Pass {
		t.Errorf("feature sem cenário não contradiz a dispensa: %v (%s)", v, msg)
	}
}

// A referência do `@no-test` tem de RESOLVER: um código que nenhum teste menciona é
// uma promessa vazia. Pior que a ausência — ela passa a impressão de que a prova foi
// conferida por alguém.
func TestTrincaCompleta_noTestComReferenciaOrfaReprova(t *testing.T) {
	root := raizComProva(t, "AAAAX-B01") // o teste prova B01…
	g := trincaGraph(true, true, false)
	spec := "@no-test: provado por `AAAAX-B99`\n" // …mas a spec alega B99

	v, msg := checkTrincaCompleta(spec, specNode(), root, g, &config.Config{})
	if v != Fail {
		t.Fatalf("referência que não resolve deve reprovar: %v", v)
	}
	if !strings.Contains(msg, "AAAAX-B99") {
		t.Errorf("a mensagem deve nomear o código órfão: %s", msg)
	}
}

// Razão em PROSA, sem código, não basta: é exatamente o "provado na integração" que
// ninguém consegue conferir. É a regra que este gate passou a cobrar.
func TestTrincaCompleta_noTestSemReferenciaReprova(t *testing.T) {
	root := raizComProva(t, "AAAAX-B01")
	g := trincaGraph(true, true, false)
	spec := "@no-test: gateway de uma linha, provado no teste de integração central\n"

	if v, _ := checkTrincaCompleta(spec, specNode(), root, g, &config.Config{}); v != Fail {
		t.Errorf("prosa sem código não é referência verificável: %v", v)
	}
}

// `@no-feature` NÃO precisa de referência: ele afirma que não há comportamento
// observável, e não existe prova a apontar. Exigir o endereço de algo que a spec acabou
// de dizer que não existe seria incoerente — e travaria todo gateway sem regra.
func TestTrincaCompleta_noFeatureNaoExigeReferencia(t *testing.T) {
	g := trincaGraph(true, false, false)
	spec := "@no-feature: wiring de infraestrutura, sem regra observável\n"

	if v, msg := checkTrincaCompleta(spec, specNode(), t.TempDir(), g, &config.Config{}); v != Pass {
		t.Errorf("`@no-feature` dispensa sem exigir referência: %v (%s)", v, msg)
	}
}

// `@TBD` (TO BE DEVELOPED) e `@no-test` dizem coisas DIFERENTES, e a diferença é o tempo.
//
// A spec nasce ANTES do código — é o fluxo normal do Anchors, já que a spec é a âncora.
// Sem o `@TBD`, quem escreve uma spec nova tem duas saídas e ambas são ruins: barrar o
// commit de todo trabalho em andamento, ou declarar `@no-test` mentindo — e aí a cobrança
// some PARA SEMPRE justamente na unidade que mais vai precisar dela.
func TestTBDDispensaSoOQueFoiDeclarado(t *testing.T) {
	casos := []struct {
		marca   string
		quer    []string
		naoQuer []string
	}{
		{"@TBD: code", []string{string(mapx.EdgeSpecifies)},
			[]string{string(mapx.EdgeTestedBy), string(mapx.EdgeCoveredBy)}},
		{"@TBD: code,test", []string{string(mapx.EdgeSpecifies), string(mapx.EdgeTestedBy)},
			[]string{string(mapx.EdgeCoveredBy)}},
		{"@TBD: code, feature, test",
			[]string{string(mapx.EdgeSpecifies), string(mapx.EdgeCoveredBy), string(mapx.EdgeTestedBy)},
			nil},
	}
	for _, c := range casos {
		got := pecasPorDesenvolver("# Spec\n\n> " + c.marca + " — em andamento\n")
		for _, q := range c.quer {
			if !got[q] {
				t.Errorf("%q deveria dispensar %q", c.marca, q)
			}
		}
		for _, n := range c.naoQuer {
			if got[n] {
				t.Errorf("%q NÃO pode dispensar %q — o alvo é obrigatório justamente para "+
					"que a marca não vire interruptor geral do gate", c.marca, n)
			}
		}
	}
}

// Sem `@TBD` nenhum, nada é dispensado — o marcador é opt-in.
func TestSemTBDNadaEhDispensado(t *testing.T) {
	if len(pecasPorDesenvolver("# Spec sem marca nenhuma\n")) != 0 {
		t.Error("sem `@TBD` o gate cobra tudo, como sempre cobrou")
	}
	// E `@TBD` sem alvo não dispensa nada: "está em andamento" sem dizer o quê seria
	// um interruptor geral, que é o oposto do que este marcador é.
	if len(pecasPorDesenvolver("# Spec\n\n> @TBD\n")) != 0 {
		t.Error("`@TBD` sem alvo não pode dispensar nada")
	}
}
