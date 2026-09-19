package health

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/i18n"
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
func checkPendingDecisions(g *mapx.Graph, root string, cfg *config.Config) []Finding {
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
		if q := gate.OpenDecisions(string(b), cfg, n.Layer); q > 0 {
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
			i18n.T("health.pending_decisions", a.quanto)})
	}
	return out
}
