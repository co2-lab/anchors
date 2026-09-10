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

func TestEvidenceFreshCarimboFrescoPassa(t *testing.T) {
	g := grafoEv()
	g.Nodes[0].Signal = &mapx.TestSignal{
		Passed: 1, AtRev: "t1",
		ClosureRev: map[string]string{"utils/login.yaml": "u1"},
	}
	v, msg := checkEvidenceFresh("", g.Nodes[0], "", g, nil)
	if v != Pass {
		t.Fatalf("fecho intacto tem de passar, veio %v — %s", v, msg)
	}
	// A mensagem diz CONTRA O QUE conferiu: é a diferença entre "ninguém olhou" e
	// "olhei e está de pé".
	if !strings.Contains(msg, "1 dependência") {
		t.Errorf("a mensagem devia dizer o tamanho do fecho conferido, veio: %s", msg)
	}
}

func TestEvidenceFreshVenceQuandoDependenciaMuda(t *testing.T) {
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
	if !strings.Contains(msg, "rode-o de novo") {
		t.Errorf("a mensagem tem de dizer o conserto, veio: %s", msg)
	}
}

func TestEvidenceFreshTruncaListaLonga(t *testing.T) {
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
	if !strings.Contains(msg, "e 15 outra(s)") {
		t.Errorf("devia truncar em 5 e somar o resto, veio:\n%s", msg)
	}
}

func TestEvidenceFreshSoOlhaTeste(t *testing.T) {
	g := grafoEv()
	code := mapx.Node{ID: "src/Tela.tsx", Kind: mapx.KindCode, Rev: "x"}
	if v, _ := checkEvidenceFresh("", code, "", g, nil); v != Skip {
		t.Errorf("nó de código não tem placar de execução, veio %v", v)
	}
}
