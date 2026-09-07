// Package board reivindica trabalho no BOARD do repositório — as issues, no `mode: github`.
//
// POR QUE ISTO EXISTE. O `workflow.mode` é excludente por decisão declarada no
// `config.go`: `local` põe a fila em `.anchors/tasks/`; `github` põe a fila nas issues. E
// o comentário de lá é literal — *"no modo github, `.anchors/tasks/` não deve existir"*.
//
// O `anchors next` não seguia isso: ele chamava `queue.Claim` em qualquer modo. Medido no
// blue-eyes, com `mode: github` declarado no `anchors.yaml`:
//
//	$ anchors next
//	fila vazia — nada a fazer
//
// Havia OITENTA E QUATRO cards abertos no board, um por spec, cada um pedindo código,
// feature, teste e documentação. A resposta estava certa sobre a fila local (vazia, porque
// nada mudara desde o último `check`) e falsa sobre o projeto — e quem a leu concluiu que o
// trabalho tinha acabado. O ciclo parou na metade: 117 specs, zero features.
//
// A REGRA DE PRIORIDADE é a do `anchors-claim.yml`, que já a implementava do lado do
// pipeline. Ela não é arbitrária:
//
//  1. RETOMAR O PRÓPRIO vence a prioridade do board. Um card que este agente já
//     reivindicou carrega o contexto da sessão dele — mandá-lo para outro joga fora o que
//     já foi lido, e deixa dois agentes com metade do entendimento cada.
//
//  2. `ready-to-review` ANTES de `to-do`. Trabalho quase pronto vale mais que trabalho
//     não começado: o primeiro vira entrega com uma revisão, o segundo precisa do ciclo
//     inteiro.
//
//  3. `needs-user` NUNCA. É o card waiting esperando decisão de gente, e entregá-lo a um
//     agente o faz decidir sozinho — que é exatamente o que o `escalate` existe para
//     impedir.
package board

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Estados do ciclo, na ordem em que o claim os oferece.
const (
	StateReadyToReview = "anchors:ready-to-review"
	StateToDo          = "anchors:to-do"
	StateInProgress    = "anchors:in-progress"
	StateInReview      = "anchors:in-review"
	StateNeedsUser     = "anchors:needs-user"
	// A variante em português existe em projetos que a declararam antes de o nome
	// canônico se firmar. Um card com ela é tão waiting quanto o outro.
	StateNeedsUserPt = "anchors:precisa-do-usuario"
)

// OwnerMarker é o prefixo do comentário que declara a posse.
//
// A posse vive num COMENTÁRIO e não no `assignee` porque as perguntas são diferentes: o
// `assignee` responde "qual pessoa é responsável", e o Anchors precisa de "qual AGENTE
// está com o trabalho" — e dois agentes na mesma máquina têm o mesmo usuário do GitHub.
//
// E comentário tem `created_at`: "o último" é uma pergunta com resposta definida, o que a
// lista de `assignees` não oferece (é um conjunto, sem ordem prometida pela API).
const OwnerMarker = "anchors-owner:"

// Client é o repositório e as labels que marcam trabalho do Anchors.
type Client struct {
	Repo   string
	Labels []string
}

// Card é o trabalho que o board entrega.
type Card struct {
	Number int
	Title  string
	Body   string
	Labels []string
	// Owner é o dono declarado no último comentário `anchors-owner:`, ou vazio.
	Owner string
	// State é a label de estado em que ele foi encontrado.
	State string
}

type rawCard struct {
	Number   int                     `json:"number"`
	Title    string                  `json:"title"`
	Body     string                  `json:"body"`
	Labels   []struct{ Name string } `json:"labels"`
	Comments []struct{ Body string } `json:"comments"`
}

func (c Client) gh(args ...string) ([]byte, error) {
	args = append(args, "--repo", c.Repo)
	out, err := exec.Command("gh", args...).Output()
	if err != nil {
		return out, fmt.Errorf("gh %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

// list traz os cards abertos que têm a label do Anchors E a de estado pedida.
//
// As duas labels, porque `--label` repetido no `gh` é E e não OU: sem a do Anchors o
// comando puxaria issue de produto, que não tem a forma que o ciclo espera.
func (c Client) list(state string) ([]Card, error) {
	if len(c.Labels) == 0 {
		return nil, fmt.Errorf("workflow.labels vazio: no modo github ele é obrigatório — " +
			"sem ele o claim puxaria qualquer issue do repositório")
	}
	args := []string{"issue", "list", "--state", "open", "--limit", "200",
		"--json", "number,title,body,labels,comments"}
	for _, l := range c.Labels {
		args = append(args, "--label", l)
	}
	if state != "" {
		args = append(args, "--label", state)
	}
	out, err := c.gh(args...)
	if err != nil {
		return nil, err
	}
	var raw []rawCard
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("resposta do gh não é o JSON esperado: %w", err)
	}
	var cards []Card
	for _, r := range raw {
		cards = append(cards, Card{
			Number: r.Number, Title: r.Title, Body: r.Body,
			Labels: labelNames(r), Owner: lastOwner(r), State: state,
		})
	}
	return cards, nil
}

func labelNames(r rawCard) []string {
	var out []string
	for _, l := range r.Labels {
		out = append(out, l.Name)
	}
	return out
}

// lastOwner devolve o dono do ÚLTIMO comentário `anchors-owner:`.
//
// O último e não o primeiro: a posse muda de mão (um card rejeitado volta ao autor
// original), e cada reivindicação fica registrada — nada é sobrescrito.
func lastOwner(r rawCard) string {
	owner := ""
	for _, cm := range r.Comments {
		b := strings.TrimSpace(cm.Body)
		if strings.HasPrefix(b, OwnerMarker) {
			owner = strings.TrimSpace(strings.TrimPrefix(b, OwnerMarker))
		}
	}
	return owner
}

func has(labels []string, name string) bool {
	for _, l := range labels {
		if l == name {
			return true
		}
	}
	return false
}

// waiting diz se o card está esperando decisão de gente.
func waiting(c Card) bool {
	return has(c.Labels, StateNeedsUser) || has(c.Labels, StateNeedsUserPt)
}

// Claim reivindica o próximo card para este agente, na ordem de prioridade do
// `anchors-claim.yml`. Devolve nil quando não há trabalho.
func (c Client) Claim(agent string) (*Card, error) {
	// 1. O PRÓPRIO, em qualquer estado vivo. Retomar vence o board.
	mine, err := c.list("")
	if err != nil {
		return nil, err
	}
	for _, card := range mine {
		if card.Owner != agent || waiting(card) {
			continue
		}
		if has(card.Labels, StateToDo) || has(card.Labels, StateInProgress) ||
			has(card.Labels, StateReadyToReview) || has(card.Labels, StateInReview) {
			card.State = liveState(card)
			return &card, nil
		}
	}

	// 2. E DEPOIS o board, na ordem: quase pronto antes de não começado.
	for _, estado := range []string{StateReadyToReview, StateToDo} {
		cards, err := c.list(estado)
		if err != nil {
			return nil, err
		}
		for _, card := range cards {
			if waiting(card) {
				continue
			}
			// Card de outro agente não é livre: ele pode estar no meio de uma sessão.
			if card.Owner != "" && card.Owner != agent {
				continue
			}
			return &card, nil
		}
	}
	return nil, nil
}

func liveState(c Card) string {
	for _, s := range []string{StateInProgress, StateInReview, StateReadyToReview, StateToDo} {
		if has(c.Labels, s) {
			return s
		}
	}
	return ""
}

// Take marca o card como deste agente e o move para `in-progress`.
//
// O comentário vem ANTES da label: se a label entrasse primeiro e o comentário falhasse, o
// card ficaria `in-progress` sem dono — e o próximo claim o veria como livre, entregando o
// mesmo trabalho a dois agentes.
func (c Client) Take(number int, agent string) error {
	if _, err := c.gh("issue", "comment", fmt.Sprint(number),
		"--body", OwnerMarker+" "+agent); err != nil {
		return err
	}
	_, err := c.gh("issue", "edit", fmt.Sprint(number),
		"--add-label", StateInProgress, "--remove-label", StateToDo)
	return err
}
