package health

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func planos(ns ...mapx.Node) *mapx.Graph {
	g := &mapx.Graph{}
	for _, n := range ns {
		n.Kind = mapx.KindPlan
		g.Nodes = append(g.Nodes, n)
	}
	// as arestas que o build criaria: só para alvo existente
	existe := map[string]bool{}
	for _, n := range g.Nodes {
		existe[n.ID] = true
	}
	for _, n := range g.Nodes {
		for _, alvo := range n.Needs {
			if existe[alvo] {
				g.Edges = append(g.Edges, mapx.Edge{From: n.ID, To: alvo, Type: mapx.EdgeNeeds})
			}
		}
	}
	return g
}

// Um `needs` para plano inexistente é INVISÍVEL no grafo — o build descarta a aresta.
// E o efeito é a ordem de trabalho declarada simplesmente não valer: o card nasce antes
// da base existir, sem nada acusar.
func TestNeedsParaPlanoInexistenteEhReportado(t *testing.T) {
	g := planos(
		mapx.Node{ID: "plans/0001-base.md"},
		mapx.Node{ID: "plans/0002-feature.md", Needs: []string{"plans/0099-nao-existe.md"}},
	)

	fs := checkNeedsDosPlanos(g)

	if len(fs) != 1 || fs[0].Check != "needs-quebrado" {
		t.Fatalf("esperava needs-quebrado, veio %+v", fs)
	}
	if !strings.Contains(fs[0].Detail, "0099-nao-existe") {
		t.Errorf("a mensagem tem de nomear o alvo ausente: %s", fs[0].Detail)
	}
}

// Ciclo: nenhum dos planos pode começar, nunca. Sem este check o board fica com cards
// que ninguém pega e ninguém sabe por quê.
func TestCicloDeNeedsEhReportadoComOCaminho(t *testing.T) {
	g := planos(
		mapx.Node{ID: "plans/a.md", Needs: []string{"plans/b.md"}},
		mapx.Node{ID: "plans/b.md", Needs: []string{"plans/a.md"}},
	)

	fs := checkNeedsDosPlanos(g)

	if len(fs) != 1 || fs[0].Check != "needs-ciclo" {
		t.Fatalf("esperava needs-ciclo, veio %+v", fs)
	}
	// O caminho é o que torna o achado acionável — sem ele, "há um ciclo" obriga o
	// leitor a reconstruí-lo à mão.
	if !strings.Contains(fs[0].Detail, "→") {
		t.Errorf("a mensagem tem de mostrar o caminho do ciclo: %s", fs[0].Detail)
	}
}

// A cadeia legítima (A ← B ← C) não é ciclo, e um projeto que declara ordem correta não
// pode receber achado — ruído recorrente treina a equipe a ignorar o doctor.
func TestCadeiaLegitimaNaoEhAchado(t *testing.T) {
	g := planos(
		mapx.Node{ID: "plans/0001-fundacao.md"},
		mapx.Node{ID: "plans/0002-backend.md", Needs: []string{"plans/0001-fundacao.md"}},
		mapx.Node{ID: "plans/0003-tela.md", Needs: []string{"plans/0002-backend.md"}},
	)

	if fs := checkNeedsDosPlanos(g); len(fs) != 0 {
		t.Errorf("cadeia em ordem não é achado: %+v", fs)
	}
}

// Projeto sem plano nenhum: nada a conferir.
func TestSemPlanosNaoReporta(t *testing.T) {
	if fs := checkNeedsDosPlanos(&mapx.Graph{}); len(fs) != 0 {
		t.Errorf("sem planos não há needs a conferir: %+v", fs)
	}
}
