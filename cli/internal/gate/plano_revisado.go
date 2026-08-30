package gate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- o plano REVISADO avisa quem o lê ---
//
// Planejar erra, e o erro aparece quando o trabalho começa. Editar o plano antigo perderia
// o registro: um plano é ÂNCORA, e editá-lo depois de implementado faz o documento
// descrever algo que não foi o que aconteceu.
//
// Por isso a revisão é um plano NOVO com `supersedes:`. O custo dessa escolha é a leitura
// fora de ordem — quem abre o plano antigo não sabe que existe um mais novo, e segue uma
// decisão que foi revista. É o defeito que este gate impede.
//
// A régua é simples: se algum plano declara `supersedes: X`, então X tem de dizer isso no
// topo. O aviso não é cosmético — é a única coisa que alcança quem chegou pelo caminho
// errado.

// avisoDeRevisaoRE casa o bloco que o plano revisado deve carregar.
var avisoDeRevisaoRE = regexp.MustCompile(`(?im)^\s*>?\s*\*{0,2}(revisado|superado|substituído)\s+por\*{0,2}\s*:?\s*\S`)

// checkPlanoRevisado confronta os dois lados do `supersedes`.
func checkPlanoRevisado(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindPlan {
		return Skip, "só plano revisa plano"
	}
	if g == nil {
		return Pending, "sem mapa não dá para achar quem revisa quem"
	}

	// Quem revisa ESTE plano?
	var revisores []string
	for _, o := range g.Nodes {
		if o.Kind != mapx.KindPlan || o.ID == n.ID {
			continue
		}
		for _, alvo := range o.Supersedes {
			if alvo == n.ID {
				revisores = append(revisores, o.ID)
			}
		}
	}
	if len(revisores) == 0 {
		// Se ESTE declara supersedes, o alvo precisa existir: apontar para um plano que
		// não está no mapa é uma revisão que não revisa nada.
		var ausentes []string
		for _, alvo := range n.Supersedes {
			achou := false
			for _, o := range g.Nodes {
				if o.ID == alvo {
					achou = true
				}
			}
			if !achou {
				ausentes = append(ausentes, alvo)
			}
		}
		if len(ausentes) > 0 {
			return Fail, fmt.Sprintf("declara `supersedes: %s`, e esse plano não está no "+
				"mapa. Uma revisão que aponta para o nada não avisa ninguém — e o plano que "+
				"se pretendia revisar continua sendo lido como vigente",
				strings.Join(ausentes, ", "))
		}
		return Skip, "nenhum plano revisa este"
	}

	// Revisado: o aviso é obrigatório, e no TOPO — quem lê de cima para baixo precisa
	// saber antes de seguir a primeira decisão.
	if !avisoDeRevisaoRE.MatchString(topo(content)) {
		return Fail, fmt.Sprintf("é revisado por %s e não avisa quem o lê. Um plano é "+
			"registro do que se decidiu, e continua valendo como registro — mas quem chega "+
			"nele sem saber da revisão segue uma decisão que foi revista. Escreva no TOPO: "+
			"`> **Revisado por** %s — <o que mudou e por quê>`",
			strings.Join(revisores, ", "), revisores[0])
	}
	return Pass, ""
}

// topo devolve o começo do documento — onde o aviso tem de estar para ser visto.
//
// 40 linhas cobre o header, o título e a primeira seção: um aviso depois disso é
// encontrado por quem já leu o plano inteiro, e aí ele não avisou nada.
func topo(content string) string {
	linhas := strings.Split(content, "\n")
	if len(linhas) > 40 {
		linhas = linhas[:40]
	}
	return strings.Join(linhas, "\n")
}
