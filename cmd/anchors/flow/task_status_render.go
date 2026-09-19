package flow

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
			b.WriteString(fmt.Sprintf(" · owner: %s", e.Card.Owner))
		}
		b.WriteString("\n")
	} else {
		b.WriteString("Task  (no card found — provide `--card N`)\n")
	}

	// --- 1.5. o que foi DESFEITO ---
	//
	// Antes do veredito do PR, e antes do que falta: uma reversão muda o que o agente
	// pensa que fez. Medido — ele fechou o card à mão, a trava desfez no mesmo minuto, e
	// ele encerrou o turno escrevendo "issue closed e resolvida, nada mais a fazer".
	//
	// O comentário da reversão estava no card e estava correto. Quem já saiu da conversa
	// não o lê — e este relato é o último lugar onde a informação ainda muda o desfecho.
	if len(e.Reverted) > 0 {
		b.WriteString("\n⚠ WHAT YOU DID WAS UNDONE\n")
		for _, r := range e.Reverted {
			b.WriteString("  · " + r + "\n")
		}
		b.WriteString("  The card moves by FACT: the claim delivers, the checks move it to\n")
		b.WriteString("  review, the merge closes it. If the movement was deliberate, put the label\n")
		b.WriteString("  `anchors:manual` on the card and redo it.\n")
	}

	// --- 2. o veredito ---
	b.WriteString("\n")
	if e.PR != nil {
		b.WriteString(fmt.Sprintf("PR    #%d %s", e.PR.Number, strings.ToLower(e.PR.State)))
		if e.PR.Total > 0 {
			b.WriteString(" · checks: " + describeChecks(e.PR.Checks, e.PR.Total))
		} else {
			// O caso que já custou uma sessão: o PR existe e NENHUM check rodou. Sem esta
			// linha, "PR #368 open" lê-se como trabalho conferido.
			b.WriteString(" · no check ran")
		}
		b.WriteString("\n")
	} else {
		b.WriteString("PR    none for this branch\n")
	}
	b.WriteString(fmt.Sprintf("Git   %s", e.Branch))
	if !e.Clean {
		b.WriteString(" · uncommitted change")
	}
	if e.Unpushed > 0 {
		b.WriteString(fmt.Sprintf(" · %d unpushed commit(s)", e.Unpushed))
	}
	if e.Clean && e.Unpushed == 0 {
		b.WriteString(" · up to date with the remote")
	}
	b.WriteString("\n")

	// --- 3. o que falta ---
	if len(e.Blocked) > 0 {
		b.WriteString("\nWaiting on a person's decision\n")
		for _, c := range e.Blocked {
			b.WriteString(fmt.Sprintf("  #%-5d %s\n", c.Number,
				codePrefix.ReplaceAllString(c.Title, "")))
		}
		b.WriteString("  (no agent resolves these; while there is no answer, what\n")
		b.WriteString("   depends on them does not move)\n")
	}

	b.WriteString("\nWhat I proved\n  · <the rules the suite confronts, and what the mutation killed>\n")
	b.WriteString("\nWhat was left out\n  · <nothing, or what was left and why>\n")

	// --- 4. o que vem ---
	b.WriteString("\nNext\n")
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
		return "unknown state"
	case "closed":
		return "closed"
	case "to-do":
		return "available, no owner"
	case "in-progress":
		return "in progress"
	case "ready-to-review":
		return "checks passed, waiting for a reviewer"
	case "in-review":
		return "under review"
	case "ready-to-test":
		return "accepted"
	case "needs-user":
		return "STOPPED waiting on a person's decision"
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
		ps = append(ps, "there is an uncommitted change — `anchors check --changed <files>` before committing")
	}
	if e.Unpushed > 0 {
		ps = append(ps, "there is an unpushed commit — `git push`")
	}

	switch {
	// Sem PR no branch, o passo depende de ONDE o card está — e não de o branch estar
	// limpo. A primeira versão dizia "abrir o PR" para qualquer card aberto, e num card
	// `in-review` isso é conselho errado: o trabalho já foi entregue, o que falta é
	// julgá-lo. Um próximo passo errado é pior que nenhum: ele parece derivado do estado.
	case e.PR == nil && e.Card != nil && canOpenPR(e):
		ps = append(ps, "open the PR — `anchors pr-body --cards <n>` writes the closing lines")
	case e.PR == nil && e.Card != nil && inReview(e.Card.State):
		ps = append(ps, fmt.Sprintf("card #%d is under review and this branch has no PR — "+
			"its PR is on another branch; `gh pr list --search %d` finds it", e.Card.Number, e.Card.Number))
	case e.PR != nil && e.PR.State == "OPEN" && e.PR.Total == 0:
		ps = append(ps, fmt.Sprintf("no check ran on PR #%d — do not trust the green that does not exist; "+
			"trigger it (`gh workflow run`) and wait", e.PR.Number))
	case e.PR != nil && e.PR.State == "OPEN" && e.PR.Checks["em curso"] > 0:
		ps = append(ps, fmt.Sprintf("the CI of PR #%d is running — `gh pr checks %d --watch` "+
			"BLOCKS until the verdict (do not end the turn here)", e.PR.Number, e.PR.Number))
	case e.PR != nil && e.PR.State == "OPEN" && e.PR.Checks["reprovou"] > 0:
		ps = append(ps, fmt.Sprintf("%d check(s) failed on PR #%d — it is work of THIS card: "+
			"fix it, push and wait again", e.PR.Checks["reprovou"], e.PR.Number))
	case e.PR != nil && e.PR.State == "OPEN":
		ps = append(ps, fmt.Sprintf("the checks of PR #%d passed — the review is missing", e.PR.Number))
	}

	if e.Card != nil && e.Card.State == "closed" {
		ps = append(ps, "the card is closed — `anchors next` asks the claim pipeline for the next one")
	}
	if len(ps) == 0 {
		ps = append(ps, "`anchors status` says where the project is; `anchors next` asks for work")
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
