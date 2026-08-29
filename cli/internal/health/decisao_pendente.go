package health

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- decisão em aberto como ponta sistêmica ---
//
// O `doctor` é o comando "me dá o panorama", e a decisão pendente não aparecia nele. A
// hierarquia ficava invertida: o que falta MEDIR (sinal de mutação ausente) era listado
// como ponta de atenção, enquanto uma regra pendente de decisão convivia com "0 pontas".
//
// E é a decisão pendente que mais precisa sobreviver à sessão — ela depende de humano e
// pode levar semanas. A assimetria fica clara na comparação: quando o `layer-boundary`
// reprova, o Anchors abre issue e a RESOLVE sozinho quando o código muda; uma pergunta em
// aberto, que dura muito mais e custa mais caro, não gerava sinal nenhum fora do `check`.
//
// Importa mais em projeto que adota o Anchors DEPOIS de pronto: é onde as perguntas se
// acumulam (o código existe, e ninguém lembra por quê), e era onde o doctor menos ajudava.
func checkDecisoesPendentes(g *mapx.Graph, root string, cfg *config.Config) []Finding {
	if g == nil {
		return nil
	}
	type pendencia struct {
		spec   string
		quanto int
	}
	var achados []pendencia
	total := 0
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindSpec {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, n.ID))
		if err != nil {
			continue // o mapa conhece o nó e o arquivo sumiu: outro achado cobre isso
		}
		// Camada VAZIA: o grafo não guarda a camada do nó (ela só existe durante o scan),
		// então a resolução do título cai em projeto > framework. Um projeto que declara
		// `section_titles` por CAMADA e não no topo não é alcançado aqui — o gate, que
		// roda com a camada em mãos, continua correto.
		if q := gate.DecisõesEmAberto(string(b), cfg, ""); q > 0 {
			achados = append(achados, pendencia{n.ID, q})
			total += q
		}
	}
	if total == 0 {
		return nil
	}
	sort.Slice(achados, func(i, j int) bool { return achados[i].quanto > achados[j].quanto })
	// Uma linha por spec, a mais carregada primeiro: a lista serve para LEVAR a alguém
	// que decide, e quem lê precisa saber onde bater primeiro.
	out := make([]Finding, 0, len(achados))
	for _, a := range achados {
		out = append(out, Finding{"decisao-pendente", Warn, a.spec,
			fmt.Sprintf("%d decisão(ões) que a spec não tomou e o código vai precisar. "+
				"Enquanto durar, implementar é adivinhar — e a adivinhação não é confrontada "+
				"por gate nenhum, porque todas as peças existem. Leve a pergunta a quem "+
				"decide e promova a resposta a regra", a.quanto)})
	}
	return out
}
