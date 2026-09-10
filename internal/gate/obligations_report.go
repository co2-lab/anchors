package gate

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// ObligationsInForce expõe os deveres em vigor (packs + inline) para fora do pacote —
// é o que o `anchors compliance` consome.
func ObligationsInForce(root string, cfg *config.Config) []config.Obligation {
	return allObligations(root, cfg)
}

// ObligationStatus é o estado de UM dever sobre o projeto inteiro.
//
// A unidade do relatório é o DEVER, não o arquivo. "Art. 17 — 42 nós sujeitos, 41
// cumprem" responde a pergunta de quem audita; "arquivo X viola", repetido 42 vezes, não.
type ObligationStatus struct {
	Name      string
	Targets   []string // os arquivos onde o dever cobra presença
	Subject   int      // nós que DISPARAM o dever (carregam o atributo-gatilho)
	Fulfilled int      // dos sujeitos, quantos cumprem
	Debt      int      // quantos assumiram a dívida (`obligation_pending`)
	Waived    int      // quantos dispensaram com motivo (`obligation_waived`)
	Missing   []string // quem não cumpre — para o modo verbose
}

// EvaluateObligations roda cada dever contra todos os nós do mapa.
//
// Reusa o MESMO caminho do gate (`checkObligationHonored`), em vez de reimplementar a
// avaliação: duas implementações da mesma regra divergem com o tempo, e a que estivesse
// errada seria justamente a do relatório — que é onde alguém confia sem conferir.
func EvaluateObligations(root string, cfg *config.Config, g *mapx.Graph, obs []config.Obligation) []ObligationStatus {
	var out []ObligationStatus
	for _, ob := range obs {
		st := ObligationStatus{Name: ob.Name, Targets: ob.MustAppearIn}
		soUma := &config.Config{
			Obligations: []config.Obligation{ob},
			// sem Packs: as obrigações já vêm resolvidas de fora
		}
		for _, n := range g.Nodes {
			b, err := os.ReadFile(filepath.Join(root, n.ID))
			if err != nil {
				continue
			}
			content := string(b)
			if ob.When == "" || !headerHasAttr(content, ob.When) {
				continue // não dispara: não é sujeito
			}
			st.Subject++
			if waiverFor(content, ob.Name) != "" {
				st.Waived++
				st.Fulfilled++ // dispensado COM motivo conta como resolvido
				continue
			}
			v, _ := checkObligationHonored(content, n, root, g, soUma)
			switch v {
			case Pass:
				st.Fulfilled++
			case Pending:
				st.Debt++
			default:
				st.Missing = append(st.Missing, n.ID)
			}
		}
		sort.Strings(st.Missing)
		out = append(out, st)
	}
	return out
}
