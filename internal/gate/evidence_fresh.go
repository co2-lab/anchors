package gate

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// evidence-fresh: o placar deste teste vale contra o código de HOJE.
//
// Os outros gates confrontam texto contra texto. Este confronta uma EXECUÇÃO contra as
// revs que ela mediu: o teste passou verde, e o `Signal` guarda a rev de cada nó do fecho
// naquele momento. Se qualquer uma avançou, o placar continua escrito e deixou de valer.
//
// É a falha que não faz ruído. Um teste que quebra grita no runner; um placar que
// envelhece continua verde no relatório de ontem, e some para dentro da conclusão de quem
// o lê. Medido no app de referência: `utils/login.yaml` é composto por 290 roteiros — tocá-lo invalidava
// 290 medições, e o sinal de nenhuma delas mudava.
//
// A régua de quando NÃO julgar é tão importante quanto a de quando reprovar:
//
//   - SEM CARIMBO → Skip, nunca Fail. Um teste que nunca rodou não tem placar para
//     expirar: isso é ausência de prova, que é outra dívida, com outro nome e outro
//     conserto (rodar a primeira vez, não revalidar). Reportá-las juntas é o defeito que
//     o `stale` de arestas tem — 30 "nunca validadas" afogadas em 1419 "avançou de rev",
//     e a lista deixa de ser lida.
//   - CARIMBO FRESCO → Pass, e vale dizer contra o que: o número de dependências
//     conferidas é a diferença entre "ninguém olhou" e "olhei e está de pé".
//
// Bloqueante porque o conserto é conhecido, local e barato: rode o teste. Não há
// julgamento de domínio a fazer — diferente de uma letra não declarada, onde as duas
// saídas (adotar ou remapear) são legítimas e a escolha é de quem conhece o produto.
func checkEvidenceFresh(_ string, n mapx.Node, _ string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest || g == nil {
		return Skip, ""
	}
	if n.Signal == nil || n.Signal.AtRev == "" {
		// Deliberadamente Skip: sem execução ingerida não existe evidência, e portanto não
		// existe evidência vencida. Quem cobra a AUSÊNCIA de teste verde é o `coverage`.
		return Skip, "sem execução ingerida — nada a envelhecer (ausência de prova é outra dívida)"
	}
	ev := g.EvidenceStaleFor(n.ID)
	if ev == nil {
		return Pass, fmt.Sprintf("placar de %s confere com o código atual (%d dependência(s) no fecho)",
			n.Signal.AtRev, len(n.Signal.ClosureRev))
	}

	var b strings.Builder
	b.WriteString("o placar deste teste mediu um código que já mudou — rode-o de novo.\n")
	if ev.Own {
		fmt.Fprintf(&b, "  o próprio arquivo de teste mudou (medido em %s)\n", ev.AtRev)
	}
	if len(ev.Culprit) > 0 {
		fmt.Fprintf(&b, "  %d dependência(s) do fecho avançaram de rev:\n", len(ev.Culprit))
		// Lista até 5: quem conserta roda o teste UMA vez, independente de quantas
		// dependências mudaram. Despejar 290 caminhos afoga a única linha que importa.
		for i, c := range ev.Culprit {
			if i == 5 {
				fmt.Fprintf(&b, "    … e %d outra(s)\n", len(ev.Culprit)-5)
				break
			}
			fmt.Fprintf(&b, "    %s\n", c)
		}
	}
	return Fail, strings.TrimRight(b.String(), "\n")
}
