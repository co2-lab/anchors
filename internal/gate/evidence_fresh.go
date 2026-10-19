// @anchors
//   code: EFGVD
//   ref: EVFRV

package gate

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
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
func checkEvidenceFresh(content string, n mapx.Node, _ string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest || g == nil {
		return Skip, ""
	}
	if n.Signal == nil || n.Signal.AtRev == "" {
		// Deliberadamente Skip: sem execução ingerida não existe evidência, e portanto não
		// existe evidência vencida. Quem cobra a AUSÊNCIA de teste verde é o `coverage`.
		return Skip, i18n.T("gate.evidence_fresh.skip_no_execution")
	}
	ev := g.EvidenceStaleFor(n.ID)
	if ev == nil {
		if imps := unboundImports(content, n, g, cfg); len(imps) > 0 {
			return Pending, i18n.T("gate.evidence_fresh.pending_unbound", strings.Join(imps, ", "))
		}
		return Pass, fmt.Sprintf(i18n.T("gate.evidence_fresh.pass_fresh"),
			n.Signal.AtRev, len(n.Signal.ClosureRev))
	}

	var b strings.Builder
	b.WriteString(i18n.T("gate.evidence_fresh.fail_stale") + "\n")
	if ev.Own {
		fmt.Fprintf(&b, i18n.T("gate.evidence_fresh.own_file_changed"), ev.AtRev)
	}
	if len(ev.Culprit) > 0 {
		fmt.Fprintf(&b, i18n.T("gate.evidence_fresh.closure_advanced"), len(ev.Culprit))
		// Lista até 5: quem conserta roda o teste UMA vez, independente de quantas
		// dependências mudaram. Despejar 290 caminhos afoga a única linha que importa.
		for i, c := range ev.Culprit {
			if i == 5 {
				fmt.Fprintf(&b, i18n.T("gate.evidence_fresh.and_others"), len(ev.Culprit)-5)
				break
			}
			fmt.Fprintf(&b, "    %s\n", c)
		}
	}
	return Fail, strings.TrimRight(b.String(), "\n")
}

// unboundImports are the project's code files a test imports when its proof reaches none: the
// score was measured against code no change can stale, and "fresh" would say nothing. Asked
// only of a project that put its tests in the dependency chain (`dep-declared` on `test`) —
// without it a test carries no `@dep`, by the project's decision, and its closure holds no
// code by construction. An import of types alone, or of a file the test replaces with a
// double, binds the proof to nothing it runs, and is no such import.
func unboundImports(content string, n mapx.Node, g *mapx.Graph, cfg *config.Config) []string {
	if !gateOnTests(cfg, "dep-declared") {
		return nil
	}
	for id := range g.EvidenceClosure(n.ID) {
		if c := g.Node(id); c != nil && c.Kind == mapx.KindCode {
			return nil
		}
	}
	doubled := map[string]bool{}
	for _, d := range n.Doubles {
		doubled[d] = true
	}
	seen := map[string]bool{}
	var out []string
	for _, imp := range ImportsOf(content, n.ID, cfg.DialectFor(), g) {
		if imp.Target == "" || imp.TypeOnly || doubled[imp.Target] || seen[imp.Target] {
			continue
		}
		if t := g.Node(imp.Target); (t != nil && t.Kind == mapx.KindCode) || imp.Dir {
			seen[imp.Target] = true
			out = append(out, imp.Target)
		}
	}
	return out
}
