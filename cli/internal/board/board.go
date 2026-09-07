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

// Mine devolve o card que JÁ É deste agente, ou nil.
//
// Só LÊ. A atribuição é do pipeline (ver `Ask`), e este cliente nunca escreve
// `anchors-owner` — é o que mantém a serialização que evita dois agentes no mesmo card.
//
// Qualquer estado vivo conta: `to-do` (a sessão anterior parou antes de começar),
// `in-progress`, `ready-to-review` e `in-review` (o card voltou para correção). O que não
// conta é `needs-user`: ali o trabalho espera decisão de gente, e entregá-lo faria o
// agente decidir sozinho.
func (c Client) Mine(agent string) (*Card, error) {
	cards, err := c.list("")
	if err != nil {
		return nil, err
	}
	for _, card := range cards {
		if card.Owner != agent || waiting(card) {
			continue
		}
		if s := liveState(card); s != "" {
			card.State = s
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

// Ask PEDE trabalho ao pipeline, e não o reivindica direto.
//
// A DIFERENÇA É A CORRIDA, e o `BOOTSTRAP.md` §7.6 a descreve: comentários resolvem *quem
// é o dono* e *em que ordem os claims chegaram*, mas **não resolvem a disputa** — dois
// agentes podem comentar quase ao mesmo tempo, ambos lerem antes do outro escrever, e
// ambos se acharem donos. A API do GitHub não oferece compare-and-swap.
//
// Então os agentes param de disputar: eles pedem, e quem atribui é o pipeline —
// serializado por `concurrency: anchors-claim`, nunca duas instâncias juntas. "Este card
// tem dono?" e "atribua a ele" viram uma operação sem ninguém no meio, e a colisão não
// chega a existir.
//
// Por isso o `anchors-owner` é escrito SÓ pelo pipeline, e este cliente lê.
func (c Client) Ask(agent string) error {
	_, err := c.gh("workflow", "run", ClaimWorkflow, "-f", "agent="+agent)
	if err != nil {
		return fmt.Errorf("não consegui pedir trabalho ao pipeline: %w\n"+
			"  (o claim é serializado por `concurrency` — é o que evita dois agentes\n"+
			"   pegarem o mesmo card, e por isso o CLI pede em vez de reivindicar)", err)
	}
	return nil
}

// ClaimWorkflow é o pipeline que atribui trabalho. O nome vem do `initx`, que o semeia.
const ClaimWorkflow = "anchors-claim.yml"
