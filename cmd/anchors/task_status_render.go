package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// codePrefix tira o `[CODIGO]` do título: o código já aparece na linha do card, e
// repeti-lo rouba a largura de onde está a informação.
var codePrefix = regexp.MustCompile(`^\[[^\]]+\]\s*`)

// renderTaskStatus é o FORMATO — a razão de este comando existir.
//
// A ordem das seções não é arbitrária: ela é a ordem em que a informação decide algo.
//
//  1. ONDE ESTÁ  — o card e seu estado. Quem lê precisa saber se isto está entregue ou
//     parado com o nome de alguém antes de qualquer detalhe.
//  2. O VEREDITO — o PR e os checks. É o que estava faltando no relato que motivou o
//     comando: um PR aberto cujos checks ninguém leu parece trabalho entregue.
//  3. O QUE FALTA — pendências, e o que espera decisão de PESSOA em separado, porque é a
//     única pendência que não se resolve continuando a trabalhar.
//  4. O QUE VEM   — o próximo passo, derivado do estado, não da intenção do agente.
//
// As duas linhas que a máquina não sabe (o que se provou, e o que ficou de fora) aparecem
// como lacunas EXPLÍCITAS. Uma lacuna nomeada é preenchida; uma seção ausente não é notada.
func renderTaskStatus(e taskState) string {
	var b strings.Builder

	// --- 1. onde está ---
	if e.Card != nil {
		t := codePrefix.ReplaceAllString(e.Card.Title, "")
		b.WriteString(fmt.Sprintf("Task  #%d · %s\n", e.Card.Number, t))
		b.WriteString(fmt.Sprintf("      %s", describeState(e.Card.State)))
		if e.Card.Owner != "" {
			b.WriteString(fmt.Sprintf(" · dono: %s", e.Card.Owner))
		}
		b.WriteString("\n")
	} else {
		b.WriteString("Task  (nenhum card encontrado — informe `--card N`)\n")
	}

	// --- 2. o veredito ---
	b.WriteString("\n")
	if e.PR != nil {
		b.WriteString(fmt.Sprintf("PR    #%d %s", e.PR.Number, strings.ToLower(e.PR.Estado)))
		if e.PR.Total > 0 {
			b.WriteString(" · checks: " + describeChecks(e.PR.Checks, e.PR.Total))
		} else {
			// O caso que já custou uma sessão: o PR existe e NENHUM check rodou. Sem esta
			// linha, "PR #368 open" lê-se como trabalho conferido.
			b.WriteString(" · nenhum check rodou")
		}
		b.WriteString("\n")
	} else {
		b.WriteString("PR    nenhum para este branch\n")
	}
	b.WriteString(fmt.Sprintf("Git   %s", e.Branch))
	if !e.Clean {
		b.WriteString(" · mudança não commitada")
	}
	if e.Unpushed > 0 {
		b.WriteString(fmt.Sprintf(" · %d commit(s) não enviado(s)", e.Unpushed))
	}
	if e.Clean && e.Unpushed == 0 {
		b.WriteString(" · em dia com o remoto")
	}
	b.WriteString("\n")

	// --- 3. o que falta ---
	if len(e.Blocked) > 0 {
		b.WriteString("\nEsperando decisão de pessoa\n")
		for _, c := range e.Blocked {
			b.WriteString(fmt.Sprintf("  #%-5d %s\n", c.Number,
				codePrefix.ReplaceAllString(c.Title, "")))
		}
		b.WriteString("  (nenhum agente resolve estas; enquanto não houver resposta, o que\n")
		b.WriteString("   depende delas não anda)\n")
	}

	b.WriteString("\nO que provei\n  · <as regras que a suíte confronta, e o que a mutação matou>\n")
	b.WriteString("\nO que ficou de fora\n  · <nada, ou o que foi deixado e por quê>\n")

	// --- 4. o que vem ---
	b.WriteString("\nPróximo\n")
	for _, p := range nextStep(e) {
		b.WriteString("  · " + p + "\n")
	}
	return b.String()
}

// describeState traduz a label para o que ela SIGNIFICA para quem lê. `anchors:in-progress`
// é vocabulário do board; "em andamento, com dono" é a informação.
func describeState(state string) string {
	switch strings.TrimPrefix(state, "anchors:") {
	case "":
		return "estado desconhecido"
	case "closed":
		return "fechado"
	case "to-do":
		return "disponível, sem dono"
	case "in-progress":
		return "em andamento"
	case "ready-to-review":
		return "checks passaram, esperando revisor"
	case "in-review":
		return "em revisão"
	case "ready-to-test":
		return "aceito"
	case "needs-user":
		return "PARADO esperando decisão de pessoa"
	default:
		return strings.TrimPrefix(state, "anchors:")
	}
}

// describeChecks nunca resume para uma palavra só.
//
// "passou" com 3 de 4 é exatamente o relato que mente por omissão — e um check EM CURSO
// contado como sucesso é o defeito que este comando existe para impedir.
func describeChecks(m map[string]int, total int) string {
	if len(m) == 1 {
		for k, v := range m {
			return fmt.Sprintf("%d/%d %s", v, total, k)
		}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Ordem fixa, e a pior primeiro: quem lê rápido tem de bater no problema, não no que
	// deu certo. `sort.Strings` daria "em curso, passou, reprovou" — alfabético, e a
	// reprovação no fim.
	weight := map[string]int{"reprovou": 0, "em curso": 1, "passou": 2}
	sort.Slice(keys, func(i, j int) bool { return weight[keys[i]] < weight[keys[j]] })
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%d %s", m[k], k))
	}
	return strings.Join(parts, ", ")
}

// nextStep é DERIVADO do estado, não da intenção do agente.
//
// É a diferença que motivou o comando: "aguardando a nova rodada" é uma intenção, e ela
// encerrou um turno com o card parado. O próximo passo de um PR cujos checks não foram
// lidos é LER os checks — e nenhum agente precisa decidir isso.
func nextStep(e taskState) []string {
	var ps []string

	if !e.Clean {
		ps = append(ps, "há mudança não commitada — `anchors check --changed <arquivos>` antes de commitar")
	}
	if e.Unpushed > 0 {
		ps = append(ps, "há commit não enviado — `git push`")
	}

	switch {
	// Sem PR no branch, o passo depende de ONDE o card está — e não de o branch estar
	// limpo. A primeira versão dizia "abrir o PR" para qualquer card aberto, e num card
	// `in-review` isso é conselho errado: o trabalho já foi entregue, o que falta é
	// julgá-lo. Um próximo passo errado é pior que nenhum: ele parece derivado do estado.
	case e.PR == nil && e.Card != nil && canOpenPR(e):
		ps = append(ps, "abrir o PR — `anchors pr-body --cards <n>` escreve as linhas de fechamento")
	case e.PR == nil && e.Card != nil && inReview(e.Card.State):
		ps = append(ps, fmt.Sprintf("o card #%d está em revisão e este branch não tem PR — "+
			"o PR dele está noutro branch; `gh pr list --search %d` acha", e.Card.Number, e.Card.Number))
	case e.PR != nil && e.PR.Estado == "OPEN" && e.PR.Total == 0:
		ps = append(ps, fmt.Sprintf("nenhum check rodou no PR #%d — não confie no verde que não existe; "+
			"dispare (`gh workflow run`) e espere", e.PR.Number))
	case e.PR != nil && e.PR.Estado == "OPEN" && e.PR.Checks["em curso"] > 0:
		ps = append(ps, fmt.Sprintf("o CI do PR #%d está rodando — `gh pr checks %d --watch` "+
			"BLOQUEIA até o veredito (não encerre o turno aqui)", e.PR.Number, e.PR.Number))
	case e.PR != nil && e.PR.Estado == "OPEN" && e.PR.Checks["reprovou"] > 0:
		ps = append(ps, fmt.Sprintf("%d check(s) reprovaram no PR #%d — é trabalho DESTE card: "+
			"conserte, empurre e espere de novo", e.PR.Checks["reprovou"], e.PR.Number))
	case e.PR != nil && e.PR.Estado == "OPEN":
		ps = append(ps, fmt.Sprintf("os checks do PR #%d passaram — falta a revisão", e.PR.Number))
	}

	if e.Card != nil && e.Card.State == "closed" {
		ps = append(ps, "o card está fechado — `anchors next` pede o próximo ao pipeline de claim")
	}
	if len(ps) == 0 {
		ps = append(ps, "`anchors status` diz onde o projeto está; `anchors next` pede trabalho")
	}
	return ps
}

// canOpenPR: só um card que ESTÁ SENDO implementado por alguém tem PR a abrir. Um card
// `to-do` não foi pego, um `in-review` já foi entregue, e um fechado acabou.
func canOpenPR(e taskState) bool {
	if !e.Clean || e.Unpushed > 0 {
		return false
	}
	return strings.TrimPrefix(e.Card.State, "anchors:") == "in-progress"
}

// inReview cobre os dois estados em que o trabalho saiu das mãos de quem implementou.
func inReview(state string) bool {
	s := strings.TrimPrefix(state, "anchors:")
	return s == "ready-to-review" || s == "in-review"
}
