package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// openCardsAbout devolve os cards ABERTOS que já tratam este alvo.
//
// POR QUE ISTO EXISTE. Dois agentes entregaram o MESMO trabalho no mesmo dia. O primeiro
// pegou o card do gate para `MetricCard.spec.md` às 11:38. O segundo, trabalhando noutro
// card, encontrou o mesmo problema às 12:02 e abriu um card NOVO para ele — e os dois PRs
// acrescentaram a mesma seção ao mesmo documento.
//
// O `claim` protege contra dois agentes pegarem o mesmo card. Não protegia contra um agente
// CRIAR um card para trabalho que já está em andamento noutro — e a fila não tinha como
// mostrar isso, porque o card novo nasce legítimo.
//
// O `internal/issue` já deduplica pelo marcador `anchors-issue-key`, mas só para quem cria
// por ele: o `escalate` monta o corpo e chama `gh issue create` direto, e passa ao largo.
// Reconstruir a chave aqui não serve — ela inclui o GATE, e um achado humano não tem gate.
// O alvo é o que as duas formas têm em comum.
//
// AVISO e não recusa: escalar duas vezes o mesmo arquivo é legítimo (dois problemas
// distintos na mesma spec), e recusar transformaria um caso comum em trabalho parado. O que
// faltava era DIZER que já há alguém ali.
func openCardsAbout(target, label string) []string {
	if target == "" || label == "" {
		return nil
	}
	// A BUSCA é por texto e devolve aproximações — o GitHub não expõe outra forma de
	// procurar no corpo. É o bastante para um AVISO: um falso positivo custa uma linha
	// a mais na tela, e um falso negativo devolve o comportamento de antes.
	out, err := exec.Command("gh", "issue", "list",
		"--state", "open", "--label", label, "--limit", "60",
		"--search", target, "--json", "number,title,body").Output()
	if err != nil {
		// Sem a consulta, segue sem o aviso: esta é uma conferência auxiliar, e
		// impedir o `escalate` por causa dela seria pior que a duplicata que ela evita.
		return nil
	}
	var cards []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if json.Unmarshal(out, &cards) != nil {
		return nil
	}

	var achados []string
	for _, c := range cards {
		// CONFIRMAÇÃO pelo caminho exato, como o `internal/issue` faz com o marcador:
		// a busca aproximada casaria `MetricCard.spec.md` com `MetricCardList.spec.md`.
		if !strings.Contains(c.Title, target) && !strings.Contains(c.Body, target) {
			continue
		}
		achados = append(achados, fmt.Sprintf("#%d %s", c.Number, c.Title))
	}
	return achados
}
