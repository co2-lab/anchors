package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/spf13/cobra"
)

// --- por que um comando, e não uma instrução no guia ---
//
// Um agente trabalha e responde. O QUE ele responde depende exclusivamente dele: não havia
// formato para relato de andamento, nem para pendências, nem para próximos passos — e o que
// varia primeiro é justamente o que decide se alguém continua.
//
// Medido nesta sessão: o agente de outro dev diagnosticou corretamente uma falha de CI,
// consertou, empurrou, e encerrou o turno com "aguardando a nova rodada". O relato estava
// certo. A informação que faltava — *o card está `in-progress` com meu nome, e o veredito do
// check que eu disparei não foi lido* — não faltou por descuido: nada a pedia.
//
// Uma instrução no guia dizendo "relate o estado do card" resolveria mal, porque o agente
// teria de DESCOBRIR o estado para escrever, e cada um descobriria de um jeito (ou
// inventaria). Este comando descobre o que a máquina sabe e imprime; o agente acrescenta só
// o que a máquina NÃO sabe — o que ele provou, e o que decidiu deixar de fora.
//
// Não é um gate, e não barra nada — como o `pr-body`, é uma ferramenta que torna o caminho
// certo mais barato que o improviso.

// taskState é o que a máquina consegue afirmar sobre a rodada, sem o agente digitar nada.
type taskState struct {
	Card     *board.Card
	Branch   string
	Clean    bool // a árvore de trabalho não tem mudança pendente
	Unpushed int  // commits à frente do remoto
	PR       *branchPR
	Blocked  []board.Card // os `needs-user`: o que espera decisão de pessoa
}

// branchPR é o PR do branch atual e o veredito dos checks dele.
type branchPR struct {
	Number int
	Estado string // OPEN, MERGED, CLOSED
	Checks map[string]int
	Total  int
}

func newTaskStatusCmd() *cobra.Command {
	var root string
	var cardNum int
	cmd := &cobra.Command{
		Use:   "task-status",
		Short: "O relato da rodada: o que está entregue, o que ficou, e o que vem",
		Long: `Imprime o estado da task em que você está trabalhando.

O que a MÁQUINA sabe, ela descobre: o card e seu estado no board, o branch, o PR e
o veredito dos checks, o que ainda não foi enviado, e as decisões paradas
esperando uma pessoa.

O que a máquina NÃO sabe, ela pede: o que você provou, e o que decidiu deixar de
fora. Essas duas linhas são suas.

Existe porque nada ditava o formato do relato, e o que varia primeiro é o que
decide se alguém continua — o estado do card, e se o veredito do CI foi lido.

    anchors task-status              # descobre o card pelo ANCHORS_AGENT
    anchors task-status --card 303   # o card explícito`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return err
			}
			e := collectTaskState(absRoot, cfg, cardNum)
			fmt.Print(renderTaskStatus(e))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().IntVar(&cardNum, "card", 0, "número do card (padrão: descobre pelo ANCHORS_AGENT)")
	return cmd
}

// collectTaskState reúne o que dá para afirmar. Cada fonte falha em silêncio de propósito:
// um relato incompleto é útil, e um comando que aborta porque o `gh` não respondeu deixa o
// agente sem formato nenhum — de volta ao improviso que este comando existe para evitar.
func collectTaskState(root string, cfg *config.Config, cardNum int) taskState {
	var e taskState

	e.Branch = strings.TrimSpace(outputOf("git", "-C", root, "rev-parse", "--abbrev-ref", "HEAD"))
	e.Clean = strings.TrimSpace(outputOf("git", "-C", root, "status", "--porcelain")) == ""
	if n := strings.TrimSpace(outputOf("git", "-C", root, "rev-list", "--count", "@{u}..HEAD")); n != "" {
		e.Unpushed, _ = strconv.Atoi(n)
	}

	if cfg.GitHubMode() {
		cli := board.Client{Repo: cfg.Workflow.Repo, Labels: cfg.Workflow.Labels}
		if cardNum > 0 {
			e.Card = cardByNumber(cli, cardNum)
		} else if c, err := cli.Mine(agentID()); err == nil {
			e.Card = c
		}
		e.Blocked = escalatedCards(cli)
	}
	e.PR = currentBranchPR(root)
	return e
}

// outputOf roda o comando e devolve a saída, ou string vazia. O erro é DESCARTADO: ver
// acima — este comando prefere um relato parcial a nenhum.
func outputOf(nome string, args ...string) string {
	out, err := exec.Command(nome, args...).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func cardByNumber(cli board.Client, n int) *board.Card {
	out := outputOf("gh", "issue", "view", strconv.Itoa(n), "--json", "number,title,labels,state")
	if out == "" {
		return nil
	}
	var r struct {
		Number int                     `json:"number"`
		Title  string                  `json:"title"`
		State  string                  `json:"state"`
		Labels []struct{ Name string } `json:"labels"`
	}
	if json.Unmarshal([]byte(out), &r) != nil {
		return nil
	}
	c := &board.Card{Number: r.Number, Title: r.Title}
	for _, l := range r.Labels {
		c.Labels = append(c.Labels, l.Name)
		if strings.HasPrefix(l.Name, "anchors:") {
			c.State = l.Name
		}
	}
	// Um card FECHADO não tem label de estado que valha: a label sobra do último trânsito,
	// e dizer `ready-to-review` de um card fechado é pior que não dizer nada.
	if r.State == "CLOSED" {
		c.State = "closed"
	}
	return c
}

// escalatedCards são os cards `needs-user`: trabalho parado esperando uma decisão que nenhum
// agente pode tomar. Entram no relato porque é o único tipo de pendência que NÃO se resolve
// continuando a trabalhar — e um relato que a omite convida o leitor a esperar por algo que
// só ele pode destravar.
func escalatedCards(cli board.Client) []board.Card {
	out := outputOf("gh", "issue", "list", "--label", "anchors:needs-user",
		"--state", "open", "--limit", "50", "--json", "number,title")
	if out == "" {
		return nil
	}
	var rs []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	if json.Unmarshal([]byte(out), &rs) != nil {
		return nil
	}
	var cs []board.Card
	for _, r := range rs {
		cs = append(cs, board.Card{Number: r.Number, Title: r.Title})
	}
	return cs
}

func currentBranchPR(root string) *branchPR {
	out := outputOf("gh", "pr", "view", "--json", "number,state,statusCheckRollup")
	if out == "" {
		return nil
	}
	var r struct {
		Number int    `json:"number"`
		State  string `json:"state"`
		Checks []struct {
			Conclusion string `json:"conclusion"`
			Status     string `json:"status"`
		} `json:"statusCheckRollup"`
	}
	if json.Unmarshal([]byte(out), &r) != nil {
		return nil
	}
	p := &branchPR{Number: r.Number, Estado: r.State, Checks: map[string]int{}}
	for _, c := range r.Checks {
		p.Total++
		// Um check EM CURSO não é um check que passou, e a diferença é a que decide se o
		// turno pode terminar: `SUCCESS` com 3 de 4 é um relato que mente por omissão.
		switch {

		case c.Conclusion == "SUCCESS" || c.Conclusion == "NEUTRAL" || c.Conclusion == "SKIPPED":
			p.Checks["passou"]++
		default:
			p.Checks["reprovou"]++
		}
	}
	return p
}
