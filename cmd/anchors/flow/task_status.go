package flow

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/co2-lab/anchors/internal/telemetry"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
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
	// Reverted: as mudanças que a trava de estado DESFEZ neste card.
	//
	// Medido: um agente fechou o card à mão, a trava reverteu no mesmo minuto, e ele
	// escreveu "issue closed e resolvida, nada mais a fazer" — sem saber. O comentário da
	// reversão estava lá e estava correto; quem já saiu da conversa não o lê.
	//
	// Aqui ele aparece no relato que o guia manda rodar ao FIM do turno — o último
	// momento em que a informação ainda muda o desfecho.
	Reverted []string
}

// branchPR é o PR do branch atual e o veredito dos checks dele.
type branchPR struct {
	Number int
	State  string // OPEN, MERGED, CLOSED
	Checks map[string]int
	Total  int
}

func newTaskStatusCmd() *cobra.Command {
	var root string
	var cardNum int
	cmd := &cobra.Command{
		Use:   "task-status",
		Short: "The round's report: what is delivered, what remains, and what comes next",
		Long: `Prints the state of the task you are working on.

What the MACHINE knows, it discovers: the card and its state on the board, the
branch, the PR and the verdict of the checks, what has not been pushed yet, and the
decisions stopped waiting for a person.

What the machine does NOT know, it asks for: what you proved, and what you decided to
leave out. Those two lines are yours.

It exists because nothing dictated the format of the report, and what varies first is
what decides whether someone continues — the state of the card, and whether the CI
verdict was read.

    anchors task-status              # discovers the card through ANCHORS_AGENT
    anchors task-status --card 303   # the explicit card`,
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
			// O EVENTO QUE MAIS IMPORTA: em que estado o turno terminou.
			//
			// Medido: um agente diagnosticou uma falha de CI, consertou, empurrou e
			// encerrou com "aguardando a nova rodada" — o card ficou `in-progress` com o
			// nome dele e ninguém soube por horas. Nada falhou; nenhum comando retornou
			// erro. Só a SEQUÊNCIA revela o problema, e é ela que este evento registra.
			emitTurnEnded(e)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().IntVar(&cardNum, "card", 0, "card number (default: discovered from ANCHORS_AGENT)")
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
	if e.Card != nil {
		e.Reverted = revertedOn(e.Card.Number)
	}
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
	p := &branchPR{Number: r.Number, State: r.State, Checks: map[string]int{}}
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

// emitTurnEnded registra o estado em que o trabalho parou.
//
// Só NÚMEROS e VOCABULÁRIO: o estado do card, se há PR, quantos checks reprovaram. Nem o
// título do card, nem o nome da branch, nem o repositório — quem investiga um caso
// específico pede o relatório local, não o painel.
func emitTurnEnded(e taskState) {
	attrs := map[string]any{
		"tem_card":     e.Card != nil,
		"tem_pr":       e.PR != nil,
		"arvore_limpa": e.Clean,
		"nao_enviado":  e.Unpushed,
	}
	if e.Card != nil {
		// O ESTADO, sem o prefixo: `in-progress`, não `anchors:in-progress`. É vocabulário
		// do produto, o mesmo em todo projeto.
		attrs["estado"] = strings.TrimPrefix(e.Card.State, "anchors:")
	}
	if e.PR != nil {
		attrs["pr_estado"] = strings.ToLower(e.PR.State)
		attrs["checks_total"] = e.PR.Total
		attrs["checks_reprovaram"] = e.PR.Checks["reprovou"]
		attrs["checks_rodando"] = e.PR.Checks["em curso"]
	}
	if common.Emitter != nil {
		common.Emitter.Emit(telemetry.New(telemetry.TurnEnded, attrs, time.Now))
	}
}

// ehReversao diz se um comentário é uma reversão da trava de estado.
//
// DUAS condições, e as duas importam: o marcador E o autor ser o bot. Um comentário de
// pessoa que por acaso comece com o mesmo símbolo não é uma reversão, e tratá-lo como uma
// faria o relato acusar algo que não aconteceu — um alarme falso gasta a atenção que o
// alarme verdadeiro vai precisar.
//
// Função própria e não um `if` inline: assim o teste chama O QUE O CÓDIGO USA, em vez de
// reescrever a mesma condição ao lado. A primeira versão do teste fazia isso, e sobreviveu
// à mutação que removia a checagem do autor — o teste media a si mesmo.
func ehReversao(corpo, autor string) bool {
	return strings.HasPrefix(corpo, initx.MarcadorDeReversao) &&
		strings.HasPrefix(autor, "github-actions")
}

// revertedOn lista as reversões que a trava de estado fez neste card.
func revertedOn(card int) []string {
	out := outputOf("gh", "issue", "view", strconv.Itoa(card), "--json", "comments")
	if out == "" {
		return nil
	}
	var r struct {
		Comments []struct {
			Body   string `json:"body"`
			Author struct {
				Login string `json:"login"`
			} `json:"author"`
		} `json:"comments"`
	}
	if json.Unmarshal([]byte(out), &r) != nil {
		return nil
	}
	var rev []string
	for _, c := range r.Comments {
		if !ehReversao(c.Body, c.Author.Login) {
			continue
		}
		// A PRIMEIRA LINHA basta: ela diz o que foi revertido e por quem. O comentário
		// inteiro tem seis parágrafos explicando a regra, e despejá-los no terminal faria
		// o relato virar um muro — quem precisa do detalhe abre o card.
		linha := c.Body
		if i := strings.IndexByte(linha, '\n'); i > 0 {
			linha = linha[:i]
		}
		linha = strings.ReplaceAll(linha, "**", "")
		linha = strings.ReplaceAll(linha, "`", "")
		rev = append(rev, strings.TrimSpace(strings.TrimPrefix(linha, initx.MarcadorDeReversao)))
	}
	return rev
}
