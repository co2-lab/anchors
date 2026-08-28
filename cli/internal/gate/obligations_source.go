package gate

import (
	"fmt"
	"os"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/pack"
)

// allObligations devolve os deveres em vigor: os declarados inline no anchors.yaml MAIS os
// que vêm dos packs adotados.
//
// A ordem importa e é deliberada: pack primeiro, inline depois. Um projeto que adota
// `privacy/lgpd` e declara uma obrigação de mesmo nome está SOBREPONDO conscientemente — a
// declaração local é a mais específica e vence. O contrário faria o pack silenciosamente
// apagar uma decisão do projeto.
//
// O resultado é memoizado por raiz: `check --all` chama isto milhares de vezes, e reler os
// packs a cada nó tornaria o gate o passo mais lento do pipeline.
func allObligations(root string, cfg *config.Config) []config.Obligation {
	if cfg == nil {
		return nil
	}
	if len(cfg.Packs) == 0 {
		return cfg.Obligations
	}
	packsCacheMu.Lock()
	defer packsCacheMu.Unlock()
	if v, ok := packsCache[root]; ok {
		return append(append([]config.Obligation{}, v...), cfg.Obligations...)
	}

	packs, avisos, err := pack.LoadAll(root, cfg.Packs, cfg.PackValues, cfg.Jurisdictions)
	if err != nil {
		// Erro de pack é erro de CONFIG e precisa aparecer. Silenciar transformaria um
		// conjunto inteiro de deveres em nada, com o relatório verde — o pior desfecho.
		fmt.Fprintf(os.Stderr, "anchors: erro ao carregar packs: %v\n", err)
		packsCache[root] = nil
		return cfg.Obligations
	}
	for _, a := range avisos {
		fmt.Fprintf(os.Stderr, "anchors: %s\n", a)
	}

	var doPack []config.Obligation
	for _, p := range packs {
		for _, ob := range p.Obligations {
			doPack = append(doPack, config.Obligation{
				Name:             ob.Name,
				When:             ob.When,
				MustAppearIn:     ob.MustAppearIn,
				IdentifiedBy:     ob.IdentifiedBy,
				IdentifiedAsForm: ob.IdentifiedAsForm,
				// A norma entra no PORQUÊ. É o que faz a mensagem do gate citar a fonte do
				// dever em vez de só afirmá-lo — e o que permite responder a um auditor
				// "onde estou em relação ao Art. 17".
				Because: comFonte(ob.Because, p.Authority, ob.Article),
			})
		}
	}
	packsCache[root] = doPack
	return append(append([]config.Obligation{}, doPack...), cfg.Obligations...)
}

// comFonte compõe o motivo com a norma que o origina.
func comFonte(because, authority, article string) string {
	fonte := article
	if authority != "" && article != "" {
		fonte = authority + ", " + article
	} else if authority != "" {
		fonte = authority
	}
	switch {
	case fonte == "":
		return because
	case because == "":
		return fonte
	default:
		return because + " (" + fonte + ")"
	}
}

var (
	packsCacheMu sync.Mutex
	packsCache   = map[string][]config.Obligation{}
)
