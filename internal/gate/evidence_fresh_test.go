package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func grafoEv() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "suites/SS-03.yaml", Kind: mapx.KindTest, Rev: "t1"},
			{ID: "utils/login.yaml", Kind: mapx.KindTest, Rev: "u1"},
		},
		Edges: []mapx.Edge{
			{From: "suites/SS-03.yaml", To: "utils/login.yaml", Type: mapx.EdgeDependsOn},
		},
	}
}

func TestEvidenceFreshSemCarimboEhSkipNaoFail(t *testing.T) {
	t.Run("EVFRV-B03: A test with no execution stamp is skipped, not failed", func(t *testing.T) {})
	t.Run("EVFRV-I01: Absence of proof and expired proof are never the same finding", func(t *testing.T) {})
	// A régua que o Adriel fixou ao promover o gate: enquanto não rodou, NÃO fica vencido
	// — fica esperando a execução. Ausência de prova é outra dívida, com outro conserto
	// (rodar a primeira vez, não revalidar), e misturá-las é o defeito que afoga a lista
	// do `stale` de arestas: 30 "nunca validadas" perdidas em 1419 "avançou de rev".
	g := grafoEv()
	v, msg := checkEvidenceFresh("", g.Nodes[0], "", g, nil)
	if v != Skip {
		t.Fatalf("sem sinal tem de ser Skip, veio %v — %s", v, msg)
	}
}

// O Skip é SILÊNCIO, não carimbo: aprovar afirmaria um frescor que ninguém mediu.
func TestEvidenceFreshSemCarimboTampoucoAprova(t *testing.T) {
	t.Run("EVFRV-I02: A test that never ran is never approved either", func(t *testing.T) {})
	t.Run("EVFRV-X01: The gate does not charge the absence of a green test", func(t *testing.T) {})
	g := grafoEv()
	v, msg := checkEvidenceFresh("", g.Nodes[0], "", g, nil)
	if v == Pass {
		t.Errorf("teste que nunca rodou não pode ser aprovado — %s", msg)
	}
	if v == Fail {
		t.Errorf("cobrar o teste ausente é régua do `coverage`, não desta — %s", msg)
	}
}

func TestEvidenceFreshCarimboFrescoPassa(t *testing.T) {
	t.Run("EVFRV-B04: A test whose closure is intact passes", func(t *testing.T) {})
	t.Run("EVFRV-B05: The passing verdict says what it checked against", func(t *testing.T) {})
	t.Run("EVFRV-X03: The gate does not read the project's configuration", func(t *testing.T) {})
	g := grafoEv()
	g.Nodes[0].Signal = &mapx.TestSignal{
		Passed: 1, AtRev: "t1",
		ClosureRev: map[string]string{"utils/login.yaml": "u1"},
	}
	// A configuração vai NIL de propósito: a verdade confrontada vive no mapa, e uma
	// régua que dependesse de ajuste poderia ser desligada por um padrão que ninguém
	// escolheu.
	v, msg := checkEvidenceFresh("", g.Nodes[0], "", g, nil)
	if v != Pass {
		t.Fatalf("fecho intacto tem de passar, veio %v — %s", v, msg)
	}
	// A mensagem diz CONTRA O QUE conferiu: é a diferença entre "ninguém olhou" e
	// "olhei e está de pé".
	if !strings.Contains(msg, "1 dependência") && !strings.Contains(msg, "1 dependency") &&
		!strings.Contains(msg, "1 dependencia") {
		t.Errorf("a mensagem devia dizer o tamanho do fecho conferido, veio: %s", msg)
	}
}

func TestEvidenceFreshVenceQuandoDependenciaMuda(t *testing.T) {
	t.Run("EVFRV-B06: A test whose dependency advanced a revision fails", func(t *testing.T) {})
	t.Run("EVFRV-B07: The failing verdict names the culprit", func(t *testing.T) {})
	t.Run("EVFRV-B08: The failing verdict states the fix", func(t *testing.T) {})
	t.Run("EVFRV-X02: The gate does not run the test nor judge whether the change broke it", func(t *testing.T) {})
	g := grafoEv()
	g.Nodes[0].Signal = &mapx.TestSignal{
		Passed: 1, AtRev: "t1",
		ClosureRev: map[string]string{"utils/login.yaml": "u1"},
	}
	g.Nodes[1].Rev = "u2" // o util mudou; o roteiro NÃO
	v, msg := checkEvidenceFresh("", g.Nodes[0], "", g, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v — %s", v, msg)
	}
	if !strings.Contains(msg, "utils/login.yaml") {
		t.Errorf("a mensagem tem de nomear o culpado, veio: %s", msg)
	}
	if !strings.Contains(msg, "rode-o de novo") && !strings.Contains(msg, "run it again") &&
		!strings.Contains(msg, "ejecútelo de nuevo") {
		t.Errorf("a mensagem tem de dizer o conserto, veio: %s", msg)
	}
}

// O `SignalStale` isolado — o próprio arquivo de teste mudou — continua sendo reportado,
// e SEPARADO do fecho: são dois achados com a mesma consequência e origens diferentes.
func TestEvidenceFreshArquivoProprioMudado(t *testing.T) {
	t.Run("EVFRV-B09: A test whose own file changed is reported separately from its closure", func(t *testing.T) {})
	g := grafoEv()
	g.Nodes[0].Signal = &mapx.TestSignal{
		Passed: 1, AtRev: "t1",
		ClosureRev: map[string]string{"utils/login.yaml": "u1"},
	}
	g.Nodes[0].Rev = "t2" // o próprio roteiro mudou; o fecho NÃO
	v, msg := checkEvidenceFresh("", g.Nodes[0], "", g, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v — %s", v, msg)
	}
	if !strings.Contains(msg, "t1") {
		t.Errorf("o laudo tem de dizer em que rev o placar foi medido, veio: %s", msg)
	}
	if strings.Contains(msg, "utils/login.yaml") {
		t.Errorf("o fecho está intacto — o achado do arquivo próprio não pode misturá-lo: %s", msg)
	}
}

func TestEvidenceFreshTruncaListaLonga(t *testing.T) {
	t.Run("EVFRV-B10: The culprit list is truncated at five and the remainder counted", func(t *testing.T) {})
	t.Run("EVFRV-I03: Truncation never hides the size of the problem", func(t *testing.T) {})
	// Quem conserta roda o teste UMA vez, independente de quantas dependências mudaram.
	// Despejar 290 caminhos afoga a única linha que importa.
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t.yaml", Kind: mapx.KindTest, Rev: "r"}}}
	closure := map[string]string{}
	for i := 0; i < 20; i++ {
		id := string(rune('a'+i)) + ".yaml"
		g.Nodes = append(g.Nodes, mapx.Node{ID: id, Kind: mapx.KindTest, Rev: "novo"})
		g.Edges = append(g.Edges, mapx.Edge{From: "t.yaml", To: id, Type: mapx.EdgeDependsOn})
		closure[id] = "velho"
	}
	g.Nodes[0].Signal = &mapx.TestSignal{Passed: 1, AtRev: "r", ClosureRev: closure}
	v, msg := checkEvidenceFresh("", g.Nodes[0], "", g, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if !strings.Contains(msg, "e 15 outra(s)") && !strings.Contains(msg, "and 15 other(s)") &&
		!strings.Contains(msg, "y 15 otra(s)") {
		t.Errorf("devia truncar em 5 e somar o resto, veio:\n%s", msg)
	}
}

func TestEvidenceFreshSoOlhaTeste(t *testing.T) {
	t.Run("EVFRV-B01: An artifact that is not a test leaves without a verdict", func(t *testing.T) {})
	g := grafoEv()
	code := mapx.Node{ID: "src/Tela.tsx", Kind: mapx.KindCode, Rev: "x"}
	if v, _ := checkEvidenceFresh("", code, "", g, nil); v != Skip {
		t.Errorf("nó de código não tem placar de execução, veio %v", v)
	}
}

// Sem mapa não há fecho a percorrer — e o gate não aprova o que não pôde olhar.
func TestEvidenceFreshSemMapaNaoJulga(t *testing.T) {
	t.Run("EVFRV-B02: Without a built map the gate stays quiet", func(t *testing.T) {})
	n := mapx.Node{ID: "suites/SS-03.yaml", Kind: mapx.KindTest, Rev: "t1",
		Signal: &mapx.TestSignal{Passed: 1, AtRev: "t1"}}
	if v, msg := checkEvidenceFresh("", n, "", nil, nil); v != Skip {
		t.Errorf("sem grafo o gate não pode julgar, veio %v — %s", v, msg)
	}
}
