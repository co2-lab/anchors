package flow

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/doct"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/co2-lab/anchors/internal/settings"

	"github.com/co2-lab/anchors/internal/config"

	"github.com/co2-lab/anchors/internal/queue"
	"github.com/spf13/cobra"
)

// A fila é o mecanismo que desacopla "algo mudou" de "alguém trabalha nisso". O
// watcher enfileira; a IA puxa. Dois comandos para a IA-conversa e o worker:
//   anchors queue        lista o trabalho vivo (a IA-conversa e o humano inspecionam)
//   anchors next         puxa+reivindica o próximo item (o worker chama)
//   anchors done <id>    fecha a task (o worker chama ao terminar o passo)
//
// Ver o playbook (`anchors guide`) para como a IA usa isto sem se prender.

func newQueueCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "queue",
		Short: "List live tasks (the work the watcher queued)",
		Long: `Shows the pending and claimed tasks — the work the watcher queued
on seeing changes. It is read-only; it claims nothing.

The conversation AI and the human use this to KNOW what there is to do, without
committing. To take work, use 'anchors next' (ideally in a worker/subagent).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			tasks, err := queue.List(absRoot)
			if err != nil {
				return err
			}
			if len(tasks) == 0 {
				fmt.Println("empty queue — no pending work")
				return nil
			}
			fmt.Printf("%d task(s) in the queue:\n\n", len(tasks))
			claimed, triage := 0, 0
			for _, t := range tasks {
				mark := "○"
				if t.State == queue.Claimed {
					mark = "◐"
					claimed++
				}
				if t.SuggestedNext == "triage" {
					triage++
				}
				fmt.Printf("%s [%s] %s\n", mark, t.State, t.ID)
				fmt.Printf("    changed:    %s (%s)\n", t.Changed, t.Kind)
				fmt.Printf("    suggestion: %s — %s%s\n", t.SuggestedNext, t.Reason,
					seedTally(absRoot, t))
				if t.ClaimedBy != "" {
					fmt.Printf("    by:         %s\n", t.ClaimedBy)
				}
			}
			// dicas de higiene da fila (o atrito da fila poluída)
			if claimed > 0 {
				fmt.Printf("\n◐ %d claimed — if a worker died, 'anchors reclaim' returns them to the queue\n", claimed)
			}
			if triage > 0 {
				fmt.Printf("~ %d in 'triage' (unmapped kind) — handle or discard with 'anchors drop <id>'\n", triage)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}

func newNextCmd() *cobra.Command {
	var root, worker string
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Pull and claim the next item from the queue (the worker calls this)",
		Long: `ATOMICALLY claims the next pending task and prints it. Called by the
worker (ideally a background subagent; or, if your AI client has no background,
a dedicated blocking session — never the main conversation).

The claim is atomic: two workers in different terminals NEVER take the same task,
so you can run 'anchors next' in parallel in several sessions.

When you FINISH the step (code written, check passed), close it with 'anchors done <id>'.
If the queue is empty, it prints that and exits with code 0.

In github mode the queue is the board: 'next' asks the claim pipeline and WAITS
(up to 3 minutes) for the card it hands out. It never dispatches a second claim
while one of yours is still pending — re-running it waits on that same claim.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if worker == "" {
				worker = defaultWorkerID()
			}

			// NO MODO GITHUB A FILA É O BOARD, e não `.anchors/tasks/`.
			//
			// O `config.go` declara isso desde o começo — *"mode: github — a fila são as
			// issues/cards do repositório"*, *"um modo OU outro, nunca um com o outro de
			// reserva"*, *"no modo github, `.anchors/tasks/` não deve existir"* — e este
			// comando desobedecia: chamava `queue.Claim` em qualquer modo.
			//
			// Medido no blue-eyes, com `mode: github` no anchors.yaml:
			//
			//	$ anchors next
			//	fila vazia — nada a fazer
			//
			// Havia 84 cards abertos, um por spec, cada um pedindo código, feature, teste
			// e documentação. A resposta estava certa sobre a fila local (vazia, porque
			// nada mudara desde o último `check`) e FALSA sobre o projeto — e quem a leu
			// concluiu que o trabalho tinha acabado. O ciclo parou com 117 specs e zero
			// features.
			cfg, _ := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if cfg.GitHubMode() {
				// A identidade do BOARD, não a da fila local: ali o PID é a resposta
				// certa (um processo, um worker); aqui a sessão atravessa invocações.
				return nextFromBoard(absRoot, cfg, agentID())
			}

			t, err := queue.Claim(absRoot, worker, nowStamp())
			if err != nil {
				return err
			}
			if t == nil {
				// PARTIDA A FRIO. O ciclo é PULL e o watcher enfileira na MUDANÇA — então
				// um plano que está lá, parado, não gera nada. Medido em cinco execuções:
				// `anchors next` respondia "fila vazia" com um plano promovido ao lado
				// listando 11 specs por fazer, e o orquestrador tinha de escolher a
				// primeira unidade à mão, fora do fluxo.
				//
				// Semear aqui, e não num comando novo: quem chega ao `next` está
				// perguntando "o que faço agora", e é essa pergunta que o plano parado
				// responde. Um `anchors plan start` seria mais um passo para lembrar — e o
				// que ninguém lembra de rodar não existe.
				if n, err := seedFromPlans(absRoot); err == nil && n > 0 {
					fmt.Printf("empty queue — seeded with %d plan(s) that still have work\n\n", n)
					if t, err = queue.Claim(absRoot, worker, nowStamp()); err != nil {
						return err
					}
				}
			}
			if t == nil {
				fmt.Println("empty queue — nothing to do")
				return nil
			}
			fmt.Printf("task claimed: %s\n\n", t.ID)
			fmt.Printf("  changed:    %s (%s)\n", t.Changed, t.Kind)
			fmt.Printf("  origin:     %s\n", t.Origin)
			fmt.Printf("  suggestion: %s\n", t.SuggestedNext)
			fmt.Printf("  reason:     %s%s\n\n", t.Reason, seedTally(absRoot, *t))
			fmt.Printf("Run the step (see `anchors guide`). For the fine detail of what to\n")
			fmt.Printf("propagate, run: anchors impact %s\n", t.Changed)
			fmt.Printf("When done:      anchors done %s\n", t.ID)

			// A maturação (QUALITY §7) aparece aqui de forma BARATA: o `next` é chamado
			// pelo worker a cada task e precisa ser rápido, então não roda os gates —
			// só conta quantos estão declarados como informativos. Quem quer saber
			// quais estão limpos roda `anchors status` ou `check`, que já medem.
			rememberMaturationCheap(absRoot)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&worker, "worker", "", "worker identifier (default: pid@host)")
	return cmd
}

func newDoneCmd() *cobra.Command {
	var root, doFile, doKind string
	var all bool
	cmd := &cobra.Command{
		Use:   "done [id]",
		Short: "Close claimed task(s) (moves to the .anchors/done/ history)",
		Long: `Marks task(s) as done: they leave the live queue and go to .anchors/done/.
The worker calls this AFTER finishing the step and 'anchors check' passing.

Closing in BATCH, because the real queue grows faster than the work:

  anchors done --file src/pricing.ts    every task of that file
  anchors done --kind spec              every task of one kind
  anchors done --all                    the whole queue

The watcher queues by CHANGE, and one stage touches several files: in a real round
the queue reached 26+ tasks for 8 deliveries, and closing them one by one made the
orchestrator give up — the queue became scenery, which is the opposite of what it
exists to be.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if len(args) == 1 {
				if err := queue.MarkDone(absRoot, args[0]); err != nil {
					return err
				}
				fmt.Printf("task done: %s → %s\n", args[0], queue.DoneDir)
				return nil
			}
			if !all && doFile == "" && doKind == "" {
				return fmt.Errorf("provide the <id>, or a filter: --file <path>, --kind <type> or --all")
			}
			tasks, err := queue.List(absRoot)
			if err != nil {
				return err
			}
			alvo := common.RelTo(absRoot, doFile)
			fechadas := 0
			for _, t := range tasks {
				if doFile != "" && t.Changed != alvo {
					continue
				}
				if doKind != "" && t.Kind != doKind {
					continue
				}
				if err := queue.MarkDone(absRoot, t.ID); err != nil {
					fmt.Printf("  ✗ %s: %v\n", t.ID, err)
					continue
				}
				fechadas++
			}
			if fechadas == 0 {
				fmt.Println("no task matched the filter — see the queue with `anchors queue`")
				return nil
			}
			fmt.Printf("%d task(s) done → %s\n", fechadas, queue.DoneDir)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&doFile, "file", "", "close every task of this file")
	cmd.Flags().StringVar(&doKind, "kind", "", "close every task of this kind (spec|code|feature|test|change…)")
	cmd.Flags().BoolVar(&all, "all", false, "close the entire queue")
	return cmd
}

func newDropCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "drop <id>",
		Short: "Discard a task from the queue without completing it (removes; does not archive)",
		Long: `Removes a live task (pending or claimed) from the queue — for junk: obsolete
tasks, duplicates, or a plan that landed as 'triage' and you do not want to handle.
Unlike 'done' (which archives into done/): drop deletes, it does not become history.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if err := queue.Drop(absRoot, args[0]); err != nil {
				return err
			}
			fmt.Printf("task discarded: %s\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}

func newReclaimCmd() *cobra.Command {
	var root string
	var force bool
	cmd := &cobra.Command{
		Use:   "reclaim",
		Short: "Return to the queue the tasks stuck in claimed (dead worker)",
		Long: `Moves back to 'pending' every task left 'claimed' — typically orphaned
from a worker that died without closing with 'done'. Run it after a crash so the
work does not stay stuck. The tasks become pullable by 'anchors next' again.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			n, err := reclaimFn(force)(absRoot)
			if err != nil {
				return err
			}
			fmt.Printf("%d task(s) returned to the queue (claimed → pending)\n", n)
			// O ZERO precisa se explicar.
			//
			// Medido: o comando respondia "0 task(s) devolvida(s)" com uma task
			// visivelmente `claimed` na fila. O número estava certo — ela foi reivindicada
			// há minutos, dentro da janela de trabalho —, e o zero sozinho parece defeito.
			// Custou dois comandos para descartar.
			if !force {
				if r := queue.RecentlyHeld(absRoot); r > 0 {
					fmt.Printf("  (%d task(s) claimed RECENTLY stayed — "+
						"someone may be on them right now.\n"+
						"   If you know the worker stopped: `anchors reclaim --force`)\n", r)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&force, "force", false,
		"ALSO release what a LIVE worker claimed (use only if you know it stopped)")
	return cmd
}

// defaultWorkerID identifica o worker sem depender de time/random (que o resto do
// projeto evita): pid + hostname. Basta para distinguir workers concorrentes.
func defaultWorkerID() string {
	host, _ := os.Hostname()
	if host == "" {
		host = "local"
	}
	return fmt.Sprintf("%d@%s", os.Getpid(), host)
}

// agentID é a identidade do agente no BOARD, e é diferente do worker da fila local.
//
// O PID NÃO SERVE AQUI. Na fila local ele é a identidade certa — um processo é um worker, e
// quando ele morre o `reclaim` devolve a task. No board a sessão do agente atravessa MUITAS
// invocações do CLI: cada `anchors next` é um processo novo.
//
// Medido: com `os.Getpid()`, três chamadas seguidas produziram três donos diferentes
// (`46781@host`, `48256@host`, `52350@host`) e três cards `in-progress` para o mesmo
// agente — cada um achando que era de outra pessoa, e nenhum sendo retomado.
//
// O formato é o que o `BOOTSTRAP.md` §7.6 define: `<maquina>/<sessao>`. A sessão vem do
// ambiente porque só o cliente de IA sabe onde ela começa e termina; sem ela, cai no nome
// do usuário — estável entre invocações, e suficiente para um dev com um agente.
func agentID() string {
	host, _ := os.Hostname()
	if host == "" {
		host = "local"
	}
	sessao := strings.TrimSpace(os.Getenv("ANCHORS_SESSION"))
	if sessao == "" {
		// O usuário do sistema é estável e distingue devs na mesma máquina. Não distingue
		// dois agentes do MESMO dev — e é por isso que `ANCHORS_SESSION` existe: quem roda
		// mais de um agente precisa declará-lo.
		sessao = os.Getenv("USER")
		if sessao == "" {
			sessao = "default"
		}
	}
	return host + "/" + sessao
}

// reclaimFn escolhe entre respeitar o worker vivo (padrão) e ignorá-lo (--force).
// O padrão é respeitar: roubar trabalho por suposição produz dois agentes no mesmo
// arquivo, e o custo disso só aparece depois.
func reclaimFn(force bool) func(string) (int, error) {
	if force {
		return queue.ReclaimForce
	}
	return queue.Reclaim
}

// seedFromPlans enfileira uma task por PLANO que ainda tem spec por nascer.
//
// "Ainda tem trabalho" é medido pelo disco, não por checkbox: um plano cujas specs
// semeadas TODAS existem já foi cumprido, e enfileirá-lo seria o ruído que faz a fila
// perder a confiança de quem a puxa. O que ele semeia sai do mesmo lugar que o gate
// `plan-seeds-valid` lê — os caminhos `.spec.md` citados no texto.
func seedFromPlans(root string) (int, error) {
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		return 0, err
	}
	files, err := scan.Walk(root, cfg)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, f := range files {
		if f.Kind != "plan" || len(f.Seeds) == 0 {
			continue
		}
		faltam := 0
		for _, s := range f.Seeds {
			if seedExists(root, s, files) {
				continue
			}
			faltam++
		}
		if faltam == 0 {
			continue // plano cumprido: tudo que ele semeia já existe
		}
		// UM plano por vez. Medido: um repositório real tinha 104 planos em `plans/` (o
		// projeto nunca moveu os cumpridos para `plans/done/`), e semear todos os que
		// "ainda têm trabalho" encheu a fila com 10 tasks de fases antigas.
		//
		// Uma fila de dez frentes simultâneas é uma fila que ninguém puxa — e a partida a
		// frio existe para responder "o que faço AGORA", que tem uma resposta só. Quem
		// quiser outro plano roda o `next` de novo depois de fechar este.
		if n > 0 {
			break
		}
		next, reason := queue.SuggestNext("plan")
		criada, err := queue.Enqueue(root, queue.Task{
			ID:            taskID(f.Path, next),
			Changed:       f.Path,
			Kind:          "plan",
			Origin:        "seed",
			SuggestedNext: next,
			// A RAZÃO NÃO GUARDA O NÚMERO.
			//
			// Ele é derivável (`faltam` e `len(Seeds)` se recontam do disco a qualquer
			// momento), e guardá-lo fazia o texto envelhecer na fila: a task nasce dizendo
			// "6 de 7", duas specs são entregues, e ela continua dizendo 6.
			//
			// Medido no blue-eyes (co2-lab/anchors#11): a razão velha me fez desconfiar de
			// uma correção que eu tinha acabado de publicar — passei quatro comandos
			// investigando um defeito que não existia, porque a contagem estava certa e o
			// texto era velho.
			//
			// Quem imprime recalcula: `contagemDeSementes` lê o disco na hora. Um campo
			// não pode servir a dois propósitos — o que a task É e qual era o estado
			// quando ela nasceu.
			Reason: reason,
		})
		if err == nil && criada {
			n++
		}
	}
	return n, nil
}

// seedExists diz se a spec que o plano semeia já está no repositório.
//
// O plano cita de duas formas legítimas: pelo caminho, ou só pelo NOME do arquivo — que é
// como se escreve em prosa ("a spec de `SubscriptionScreen.spec.md`"). Medido: 10 das 26
// citações de um repositório real são por nome, e tratá-las como ausentes fazia o plano
// parecer eternamente não-cumprido — a fila semeava fases já concluídas na partida a frio.
//
// Por nome, só quando ele é ÚNICO: dois arquivos homônimos tornam a citação ambígua, e
// escolher um seria decidir pelo autor.
func seedExists(root, seed string, files []scan.File) bool {
	if _, err := os.Stat(filepath.Join(root, seed)); err == nil {
		return true
	}
	base := filepath.Base(seed)
	achou := 0
	for _, f := range files {
		if filepath.Base(f.Path) == base {
			achou++
		}
	}
	return achou == 1
}

// rememberMaturationCheap conta os gates informativos sem RODAR nenhum.
//
// O `next` é chamado pelo worker a cada task, e rodar a suíte de gates ali dobraria o
// custo de puxar trabalho. O que ele pode fazer sem custo é ler a declaração: se há gate
// informativo, existe maturação pendente — e quem quiser saber QUAIS estão limpos roda
// `anchors status`, que já mede.
//
// É um lembrete mais fraco de propósito. Um lembrete caro num comando de laço quente é
// um lembrete que alguém vai querer desligar.
func rememberMaturationCheap(root string) {
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		return
	}
	var informativos int
	for _, g := range cfg.Gates {
		if !g.IsBlocking() {
			informativos++
		}
	}
	if informativos == 0 {
		return
	}
	fmt.Printf("\n○ %d informative gate(s) declared — they measure and do not defend.\n", informativos)
	fmt.Println("  `anchors status` shows which are already clean and can become blocking.")
}

// seedTally recalcula, NA HORA DE IMPRIMIR, quantas specs de um plano faltam.
//
// Vive aqui e não no `Reason` porque o número é derivável do disco, e um número guardado
// envelhece: a task nascia dizendo "6 de 7", duas specs eram entregues, e ela continuava
// dizendo 6 (co2-lab/anchors#11).
//
// Devolve string VAZIA para qualquer task que não seja semeadura de plano — o `reason` de
// julgamento (`gate 'X' — pergunta: …`) não envelhece, porque a pergunta é do gate e não
// do estado.
func seedTally(root string, t queue.Task) string {
	if t.Kind != "plan" || t.Origin != "seed" {
		return ""
	}
	// Carrega a config aqui em vez de recebê-la: os dois chamadores não a têm em escopo, e
	// passá-la obrigaria os dois a carregá-la só para isto. Falha em silêncio — a contagem
	// é um extra na mensagem, e um projeto sem config tem outro problema, que o comando
	// que chama já reporta.
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		return ""
	}
	files, err := scan.Walk(root, cfg)
	if err != nil {
		return ""
	}
	for _, f := range files {
		if f.Path != t.Changed || len(f.Seeds) == 0 {
			continue
		}
		faltam := 0
		for _, sd := range f.Seeds {
			if !seedExists(root, sd, files) {
				faltam++
			}
		}
		// A task ANTIGA já traz a contagem gravada no `reason` — somar a recalculada
		// imprimiria o número duas vezes. Medido ao instalar a correção com a fila cheia.
		//
		// Detectar pela frase e não por versão: o campo não guarda quem o escreveu, e uma
		// task sobrevive a qualquer número de atualizações do binário.
		if strings.Contains(t.Reason, "spec(s) of this plan") ||
			strings.Contains(t.Reason, "spec(s) deste plano") {
			return ""
		}
		return fmt.Sprintf(" — %d of %d spec(s) of this plan do not exist yet",
			faltam, len(f.Seeds))
	}
	return ""
}

// nextFromBoard reivindica o próximo card no modo github.
//
// A ordem é a do `anchors-claim.yml`, que já a implementava do lado do pipeline: retomar o
// próprio antes do board, `ready-to-review` antes de `to-do`, e `needs-user` nunca. A
// razão de cada uma está no `internal/board`.
//
// O QUE ELE IMPRIME é a diferença que importa. O card diz "Implementar spec — <título>", e
// o verbo é preciso: CRIAR a spec é o trabalho do card do plano; IMPLEMENTAR a spec é este
// card — código, feature, teste e documentação.
//
// Imprimir a cadeia não conserta o título: ele já está certo. Serve para quem chega ao card
// sem o contexto do fluxo, e é o que faz o próximo passo ser um comando em vez de uma
// dedução — medido, o agente que leu só o título não seguiu para nenhuma das quatro etapas.
func nextFromBoard(root string, cfg *config.Config, agent string) error {
	if cfg.Workflow.Repo == "" {
		return fmt.Errorf("workflow.repo is empty: in github mode it is REQUIRED — " +
			"inferring it from the remote would make Anchors write to another repository when someone " +
			"works on a fork, and a write in the wrong place is not undone with a revert")
	}
	// A DECISÃO LOCAL, antes de pedir trabalho.
	//
	// Num projeto com vários devs, cada um roda o seu agente — e nem todos podem decidir
	// pelo produto. Um agente que pega um card escalonado e pergunta a quem o está
	// rodando obtém uma resposta, e ela pode não ser a do dono: o escalonamento existe
	// justamente para levar a pergunta a quem decide, e um agente prestativo demais o
	// curto-circuita.
	//
	// Pergunta UMA VEZ, e o `.anchors/settings.yaml` guarda a resposta — fora do git,
	// porque é decisão de quem opera, não do projeto. Quem já declarou não é perguntado
	// de novo: perguntar a cada sessão é como se ensina alguém a responder sem ler.
	local, err := ensureLocalDecision(root)
	if err != nil {
		return err
	}

	cli := board.Client{
		Repo: cfg.Workflow.Repo, Labels: cfg.Workflow.Labels,
		UserIssues: local.HandlesUserIssues(),
	}

	// PRIMEIRO o que já é meu: se este agente tem card, não há nada a pedir.
	//
	// Retomar vence a prioridade do board porque o contexto da sessão anterior vale mais
	// que a ordem da fila — e pedir trabalho novo com um card na mão deixaria o primeiro
	// órfão, com a posse registrada e ninguém trabalhando nele.
	card, err := cli.Mine(agent)
	if err != nil {
		return err
	}
	if card != nil {
		fmt.Printf("resuming your card — FINISH it before taking another\n")
		fmt.Printf("  the context of the previous session is worth more than the queue order.\n\n")
	} else {
		// NO CARD: ask the pipeline — it is serialized, and that is what keeps two agents
		// from receiving the same work (see `board.Ask`) — and WAIT for the answer.
		//
		// This used to stop after the dispatch and tell the agent to run `anchors next`
		// again. Each re-run while the claim was pending dispatched another one, which
		// cancelled the pending run or duplicated the claim (blue-eyes #679). Now this
		// same call waits, bounded, for the run it dispatched — or for the one of this
		// agent that is already pending, which it never dispatches twice.
		fmt.Printf("work requested from the pipeline (serialized claim) — waiting for it, up to %s\n",
			board.DefaultClaimTimeout)
		out, err := cli.AskAndWait(agent, board.ClaimWait{})
		if err != nil {
			return err
		}
		if out.Card == nil {
			return reportClaimWithoutCard(out)
		}
		card = out.Card
		fmt.Println()
	}

	fmt.Printf("card claimed: #%d — %s\n\n", card.Number, card.Title)
	fmt.Printf("  state:    %s\n", strings.TrimPrefix(card.State, "anchors:"))
	fmt.Printf("  owner:    %s\n\n", agent)
	printBoardWork(root, card)
	return nil
}

// reportClaimWithoutCard says why the wait ended without a card, and what to do next.
//
// Every branch names the run, so the agent follows THAT run instead of asking again. A
// claim that failed is an error; one that ran and found nothing is not.
func reportClaimWithoutCard(out board.ClaimOutcome) error {
	runRef := fmt.Sprintf("gh run list --workflow %s --limit 5", board.ClaimWorkflow)
	if out.Run != nil {
		runRef = fmt.Sprintf("gh run view %d --log", out.Run.ID)
	}
	switch {
	case out.TimedOut:
		fmt.Printf("the claim did not finish within %s — it is still queued behind other claims, "+
			"or the pipeline is stuck.\n", board.DefaultClaimTimeout)
		fmt.Printf("  follow it:  %s\n", runRef)
		fmt.Printf("  then run `anchors next` again: it waits for this same claim and does NOT " +
			"dispatch another while it is pending.\n")
		return nil
	case out.Run != nil && out.Run.Conclusion == "success":
		fmt.Printf("the claim ran and handed you no card: no free card on the board (or the " +
			"project is frozen).\n")
		fmt.Printf("  why:  %s\n", runRef)
		return nil
	case out.Run != nil:
		return fmt.Errorf("the claim run #%d ended as %q and handed you no card — see `%s`",
			out.Run.ID, out.Run.Conclusion, runRef)
	}
	return nil
}

// printBoardWork diz O QUE ENTREGAR — os quatro artefatos e a ordem entre eles.
//
// O card nasce quando a spec APARECE no repositório, e pede a implementação dela: o card do
// plano CRIA as specs, este IMPLEMENTA uma. São dois trabalhos com entregáveis diferentes,
// e o título de cada um já diz qual.
//
// O que faltava era o CORPO nomear os quatro artefatos e a dependência entre eles. Medido:
// as 84 issues deste tipo ficaram em `to-do` enquanto o agente concluía que o projeto tinha
// terminado — o `next` não as via (lia a fila local), e nada no card dizia que a
// implementação são quatro entregas, nem que o teste nasce da feature.
func printBoardWork(root string, card *board.Card) {
	// O ESTADO manda, e o título é o assunto.
	//
	// Um card em `in-review` está com a revisão por fazer — e as instruções de
	// implementação ali mandam refazer o que já foi entregue. Medido: o card #321, com a
	// trinca completa, recebeu "ENTREGÁVEL: código + feature + teste + documentação" e o
	// agente foi conferir se tinha esquecido algo. O trabalho que faltava era outro.
	if card.State == board.StateInReview {
		printReviewWork(root, card)
		return
	}

	// O ALVO vem do CÓDIGO, não da pasta.
	//
	// O corpo do card traz as duas coisas — `Unidade: \`packages/infra\`` e
	// `Código: \`GLCGL\`` — e a primeira é a PASTA da camada. Medido: `packages/infra`
	// tem nove specs, e `anchors work code --for packages/infra` não sabe de qual delas
	// se trata.
	//
	// O código é único por unidade, e o mapa sabe o arquivo dele. O `code list --json` já
	// expõe isso, e o comentário de lá prevê exatamente este uso: *"quem consome isto
	// precisa NOMEAR o trabalho — e `onde` é a PASTA da unidade, não o arquivo"*.
	alvo := targetOfCode(root, codeFromBody(card.Body))
	if alvo == "" {
		// Sem código resolvido, a pasta é o que há — e é melhor que nada, porque o
		// `work` ao menos encontra a camada.
		alvo = unitFromBody(card.Body)
	}
	unidade := alvo
	if alvo == "" {
		alvo = "<the unit of the card>"
	}

	switch {
	case strings.Contains(card.Title, "Implementar plan"):
		fmt.Println("DELIVERABLE: every spec the plan seeds.")
		fmt.Println("  A single job — the plan is the input, the specs are the output.")
		fmt.Println()
		fmt.Printf("    anchors work spec --for <target the plan lists>\n\n")
		fmt.Println("  The plan's `-progress.md` is where you mark what you delivered, and the gate")
		fmt.Println("  `progress-honest` confronts each `[x]` against the disk.")

	case strings.Contains(card.Title, "Implementar spec"):
		fmt.Println("DELIVERABLE: code + feature + test + documentation.")
		fmt.Println("  The spec is the INPUT of this card — creating it was the work of the plan's")
		fmt.Println("  card. Implementing it is this one, and there are four deliveries.")
		fmt.Println()
		fmt.Printf("    1. anchors work code    --for %s\n", alvo)
		fmt.Printf("    2. anchors work feature --for %s\n", alvo)
		fmt.Printf("    3. anchors work test    --for %s\n", alvo)
		fmt.Println("    4. the documentation evolves TOGETHER: at every change that should be")
		fmt.Println("       documented, update the doc in the same delivery — leaving it for the")
		fmt.Println("       end is what makes it born incomplete.")
		printDocDuties(root, unidade)
		fmt.Println()
		fmt.Println("  The order 1→2→3 is not taste: the feature describes the behavior in")
		fmt.Println("  scenarios, and it is from it that the tests are born (`anchors work test` starts from the")
		fmt.Println("  feature, not from the spec).")

	default:
		fmt.Println("DELIVERABLE: see the body of the card.")
	}

	// A RÉGUA DA AUTONOMIA no próprio card, e não só no guia.
	//
	// O guia é um comando que alguém precisa rodar; o card é o que o agente lê sempre. Um
	// agente que não decide o produto e não abriu o guia perguntaria ao dev — e é a porta
	// que mais se usa, porque não tem label nem gate.
	if !decidesProduct(root) {
		fmt.Println()
		fmt.Println("  YOU DO NOT DECIDE THE DIRECTION OF THIS PRODUCT (`.anchors/settings.yaml`).")
		fmt.Println("  Faced with ambiguity or a choice that changes the behavior:")
		fmt.Println("      anchors escalate \"<what needs to be decided>\" --about <file> --for-user")
		fmt.Println("  Do NOT ask whoever is running you — the answer is reasonable and becomes a")
		fmt.Println("  product decision by someone who had no authority, without any trace.")
	}

	fmt.Printf("\nBefore starting:   anchors guide work\n")
	if unidade != "" {
		fmt.Printf("What to propagate: anchors impact %s\n", unidade)
	}
	// O CORPO DO PR pelo comando, e não à mão.
	//
	// O `Closes #N` é o que fecha o card no merge, e ele é fácil de esquecer quando o
	// corpo é escrito à mão — medido: o card #319 ficou aberto em `in-progress` depois do
	// merge, e o #321 fechou sozinho, porque um PR tinha a linha e o outro não. O estado
	// do board passou a divergir do repositório sem nada acusar.
	fmt.Printf("When done:         anchors pr-body --cards %d  (brings the `Closes` that closes the card)\n", card.Number)
	fmt.Printf("                   open the PR with that body — the pipeline moves the card, not you\n")
}

// unitFromBody extrai o caminho da unidade do corpo do card.
//
// O `identify.yml` o escreve como "Unidade: `<pasta>`" — e é o que transforma o card num
// alvo concreto para o `anchors work`. Sem ele o agente precisa adivinhar sobre o que o
// card fala a partir do título.
var unitInBodyRE = regexp.MustCompile("(?m)^Unidade:\\s*`([^`]+)`")

func unitFromBody(body string) string {
	m := unitInBodyRE.FindStringSubmatch(body)
	if len(m) < 2 {
		return ""
	}
	// A pasta pode vir como lista ("a, b"): o primeiro alvo é o suficiente para o `work`.
	return strings.TrimSpace(strings.SplitN(m[1], ",", 2)[0])
}

// codeFromBody extrai o código da unidade do corpo do card.
//
// O `identify.yml` o escreve como "Código: `XXXXX`", e ele é a identidade ESTÁVEL do
// trabalho: sobrevive a mover o arquivo de pasta, ao contrário do caminho.
var codeInBodyRE = regexp.MustCompile("(?m)^C[óo]digo:\\s*`([^`]+)`")

func codeFromBody(body string) string {
	m := codeInBodyRE.FindStringSubmatch(body)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

// targetOfCode devolve o ARQUIVO da unidade com aquele código, lendo o mapa.
//
// Falha em silêncio (string vazia) de propósito: um card cujo código não está mais no mapa
// aponta para trabalho que mudou de identidade, e o chamador cai na pasta — que ao menos
// leva o `work` à camada certa. Barrar aqui deixaria o agente sem nada.
func targetOfCode(root, code string) string {
	if code == "" {
		return ""
	}
	g, err := mapx.Load(filepath.Join(root, mapx.DefaultPath))
	if err != nil {
		return ""
	}
	for _, n := range g.Nodes {
		if n.Code == code {
			return n.ID
		}
	}
	return ""
}

// printDocDuties nomeia as documentações que ESTA alteração obriga tocar.
//
// "A documentação evolui junto" não diz QUAL documentação, e um agente que lê só isso
// atualiza o que lhe parece documentação — normalmente um README. As obrigações reais são
// declaradas no `docs:` do `anchors.yaml`, e dependem do tipo de projeto: uma API tem um
// contrato que vive fora do código, e um endpoint que não entra nele é invisível para quem
// consome.
//
// Silencioso quando o projeto não declara nada: cobrar OpenAPI de quem não tem API seria
// ruído, e ruído no card é o que faz o agente parar de ler o card.
func printDocDuties(root, unidade string) {
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil || cfg == nil {
		return
	}
	camada := layerOfPath(root, cfg, unidade)
	// PELA UNIDADE, e não só pela camada: a camada erra sozinha. Medido no projeto de
	// referência: `infra` tem nove unidades e apenas UMA toca esquema de dados — cobrar
	// o esquema das outras oito ensina o agente a ignorar o aviso, que é o pior
	// resultado possível (o gate continua lá e ninguém o lê).
	deveres := cfg.RequiredFor(camada, codeOfUnit(root, unidade))
	if len(deveres) == 0 {
		return
	}
	fmt.Println()
	alvo := camada
	if u := strings.TrimSpace(unidade); u != "" {
		alvo = u
	}
	fmt.Printf("       Changing `%s` REQUIRES touching:\n", alvo)
	for _, d := range deveres {
		fmt.Print(doct.Duty(d))
	}
	fmt.Println("       `anchors docs duties --layer " + camada + "` repeats this at any time.")
}

// layerOfPath devolve a camada de um caminho.
//
// Delega ao `scan`, que é a autoridade: reimplementar a leitura dos padrões aqui manteria
// duas versões da mesma regra, e elas divergiriam na primeira mudança da Estrutura — um
// `overrides` novo passaria a valer para o mapa e não para o card.
func layerOfPath(root string, cfg *config.Config, caminho string) string {
	if caminho == "" {
		return ""
	}
	return scan.LayerOfUnit(root, caminho, cfg)
}

// printReviewWork diz o que fazer com um card que está em REVISÃO.
//
// A revisão é a etapa que mais some, e por um motivo específico: ela não produz artefato
// novo. Quem a pula não vê nada faltando — os quatro arquivos estão lá, os gates estão
// verdes, e o card parece pronto. É exatamente por isso que o texto aqui precisa dizer que
// o trabalho EXISTE, e qual é.
func printReviewWork(root string, card *board.Card) {
	alvo := targetOfCode(root, codeFromBody(card.Body))
	if alvo == "" {
		alvo = unitFromBody(card.Body)
	}

	fmt.Println("DELIVERABLE: the REVIEW — and it produces no new file.")
	fmt.Println("  The implementation of this card was already delivered. What is missing is confronting it")
	fmt.Println("  against what it SAYS it did, and that is why the step disappears: nothing")
	fmt.Println("  seems to be missing. The four files are there and the gates are green.")
	fmt.Println()
	if alvo != "" {
		fmt.Printf("    anchors work review --for %s\n\n", alvo)
	} else {
		// A card whose body names no unit (a plan card, a finding under another card)
		// printed `--for ` with nothing after it — a command that does nothing. The PR is
		// what carries the delivery, so that is what to open.
		fmt.Printf("    this card names no unit: review the PR that references it —\n"+
			"    gh pr list --state open --search \"#%d in:body\"\n\n", card.Number)
	}
	fmt.Println("  Green gates do NOT prove it is right — they confront what is")
	fmt.Println("  DECLARABLE. In three rounds of a real E2E, 7 serious defects (silent")
	fmt.Println("  data loss, a rule with no test that proves it, contradiction between two")
	fmt.Println("  rules of the same spec) passed with everything green. None came from a gate.")
	fmt.Println()
	fmt.Println("  The delivery record is in the comments of this issue: it is the DECLARED")
	fmt.Println("  half of the confrontation — what the author says they did, against the disk.")
	fmt.Println()
	fmt.Println("  If the review finds nothing, say so and move on. If it finds something, the finding becomes a")
	fmt.Println("  correction in this card — not a new card.")

	fmt.Printf("\nBefore starting:   anchors guide review\n")
	// O CORPO DO PR pelo comando, e não à mão.
	//
	// O `Closes #N` é o que fecha o card no merge, e ele é fácil de esquecer quando o
	// corpo é escrito à mão — medido: o card #319 ficou aberto em `in-progress` depois do
	// merge, e o #321 fechou sozinho, porque um PR tinha a linha e o outro não. O estado
	// do board passou a divergir do repositório sem nada acusar.
	fmt.Printf("When done:         anchors pr-body --cards %d  (brings the `Closes` that closes the card)\n", card.Number)
	fmt.Printf("                   open the PR with that body — the pipeline moves the card, not you\n")
}

// ensureLocalDecision pergunta, uma vez, se este agente atua nos cards escalonados.
//
// Só pergunta quando NÃO HÁ decisão registrada. As três situações são distintas, e tratar
// "não declarado" como "não" apagaria a única em que a pergunta cabe:
//
//	nil    nunca perguntei      → pergunta agora
//	false  disse que não        → segue, sem perguntar
//	true   disse que sim        → segue, sem perguntar
//
// Fora do terminal (CI, pipeline) não há a quem perguntar: assume o padrão fechado e
// registra, para o próximo `next` não travar esperando entrada que não vem.
func ensureLocalDecision(root string) (settings.Settings, error) {
	s, err := settings.Load(root)
	if err != nil {
		return s, err
	}
	if s.Decided() {
		return s, nil
	}
	if !interactiveTerminal() {
		fmt.Println("(no terminal to ask — this agent has no declared role,")
		fmt.Println(" and without a role it does not act on escalated cards. Declare it with")
		fmt.Println(" `anchors settings role`)")
		fmt.Println()
		return s, nil
	}
	perfil, err := common.AskRole()
	if err != nil {
		return s, err
	}
	s.Role = perfil
	s.UserIssues = nil
	s.Agent = agentID()
	if err := settings.Save(root, s); err != nil {
		return s, err
	}
	fmt.Println()
	common.PrintRole(perfil)
	fmt.Printf("\n  recorded in %s — I will not ask again\n\n", settings.Path(root))
	return s, nil
}

// interactiveTerminal diz se há alguém do outro lado para responder.
//
// O `ModeCharDevice` sozinho NÃO basta, e foi o defeito: `/dev/null` também é char device,
// então um agente que rode `anchors next < /dev/null` — o caso normal de execução
// não-interativa — passava pela guarda, caía na pergunta e morria com `erro: ler a
// resposta: EOF`.
//
// Medido com um dev novo: o agente dele instalou tudo, declarou o perfil, e não conseguiu
// pedir trabalho. Um comando que morre por EOF não diz o que fazer — ele parece defeito da
// ferramenta, e quem lê vai procurar o problema no ambiente.
//
// A verificação certa é o `ioctl` que o terminal responde e o `/dev/null` não. Em Go, sem
// dependência: um `Stat` no dispositivo e a comparação do inode com o do stdin. Mais
// simples e igualmente confiável: `/dev/null` tem tamanho 0 e um terminal não tem tamanho —
// mas o que os separa de verdade é o `Rdev`.
func interactiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	// `/dev/null` é char device como o terminal. O que os distingue é ser o MESMO
	// dispositivo: se o stdin aponta para `/dev/null`, não há ninguém lá.
	nul, err := os.Stat(os.DevNull)
	if err != nil {
		// Sem poder comparar, assume o lado FECHADO: melhor não perguntar do que morrer
		// por EOF no meio de um pipeline.
		return false
	}
	return !os.SameFile(fi, nul)
}

// decidesProduct diz se este agente declarou que decide o rumo do produto.
//
// Erra para o lado FECHADO quando não consegue ler: um erro de leitura que resultasse em
// "pode perguntar" seria a falha silenciosa mais cara deste mecanismo — o agente perguntaria
// e ninguém saberia que a régua não foi aplicada.
func decidesProduct(root string) bool {
	s, err := settings.Load(root)
	if err != nil {
		return false
	}
	return s.HandlesUserIssues()
}
