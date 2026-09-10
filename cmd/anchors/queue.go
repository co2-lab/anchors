package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
		Short: "Lista as tasks vivas (o trabalho que o watcher enfileirou)",
		Long: `Mostra as tasks pendentes e reivindicadas — o trabalho que o watcher
enfileirou ao ver mudanças. É só leitura; não reivindica nada.

A IA-conversa e o humano usam isto para SABER o que há para fazer, sem se prender.
Para pegar trabalho, use 'anchors next' (idealmente num worker/subagente).`,
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
				fmt.Println("fila vazia — nenhum trabalho pendente")
				return nil
			}
			fmt.Printf("%d task(s) na fila:\n\n", len(tasks))
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
				fmt.Printf("    mudou:    %s (%s)\n", t.Changed, t.Kind)
				fmt.Printf("    sugestão: %s — %s%s\n", t.SuggestedNext, t.Reason,
					seedTally(absRoot, t))
				if t.ClaimedBy != "" {
					fmt.Printf("    por:      %s\n", t.ClaimedBy)
				}
			}
			// dicas de higiene da fila (o atrito da fila poluída)
			if claimed > 0 {
				fmt.Printf("\n◐ %d claimed — se um worker morreu, 'anchors reclaim' as devolve à fila\n", claimed)
			}
			if triage > 0 {
				fmt.Printf("~ %d em 'triage' (kind não mapeado) — trate ou descarte com 'anchors drop <id>'\n", triage)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	return cmd
}

func newNextCmd() *cobra.Command {
	var root, worker string
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Puxa e reivindica o próximo item da fila (o worker chama)",
		Long: `Reivindica ATOMICAMENTE a próxima task pendente e a imprime. Chamado pelo
worker (idealmente um subagente em background; ou, se seu cliente de IA não tem
background, uma sessão bloqueante dedicada — nunca a conversa principal).

O claim é atômico: dois workers em terminais diferentes NUNCA pegam a mesma task,
então você pode rodar 'anchors next' em paralelo em várias sessões.

Ao TERMINAR o passo (código escrito, check passou), feche com 'anchors done <id>'.
Se a fila está vazia, imprime isso e sai com código 0.`,
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
					fmt.Printf("fila vazia — semeada com %d plano(s) que ainda têm trabalho\n\n", n)
					if t, err = queue.Claim(absRoot, worker, nowStamp()); err != nil {
						return err
					}
				}
			}
			if t == nil {
				fmt.Println("fila vazia — nada a fazer")
				return nil
			}
			fmt.Printf("task reivindicada: %s\n\n", t.ID)
			fmt.Printf("  mudou:    %s (%s)\n", t.Changed, t.Kind)
			fmt.Printf("  origem:   %s\n", t.Origin)
			fmt.Printf("  sugestão: %s\n", t.SuggestedNext)
			fmt.Printf("  motivo:   %s%s\n\n", t.Reason, seedTally(absRoot, *t))
			fmt.Printf("Execute o passo (veja `anchors guide`). Para o detalhe fino do que\n")
			fmt.Printf("propagar, rode: anchors impact %s\n", t.Changed)
			fmt.Printf("Ao terminar:    anchors done %s\n", t.ID)

			// A maturação (QUALITY §7) aparece aqui de forma BARATA: o `next` é chamado
			// pelo worker a cada task e precisa ser rápido, então não roda os gates —
			// só conta quantos estão declarados como informativos. Quem quer saber
			// quais estão limpos roda `anchors status` ou `check`, que já medem.
			rememberMaturationCheap(absRoot)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&worker, "worker", "", "identificador do worker (default: pid@host)")
	return cmd
}

func newDoneCmd() *cobra.Command {
	var root, doFile, doKind string
	var all bool
	cmd := &cobra.Command{
		Use:   "done [id]",
		Short: "Fecha task(s) reivindicada(s) (move para o histórico .anchors/done/)",
		Long: `Marca task(s) como concluída(s): saem da fila viva e vão para .anchors/done/.
O worker chama isto DEPOIS de terminar o passo e o 'anchors check' passar.

Fechar EM LOTE, porque a fila real cresce mais rápido que o trabalho:

  anchors done --file src/pricing.ts    todas as tasks daquele arquivo
  anchors done --kind spec              todas as tasks de um tipo
  anchors done --all                    a fila inteira

O watcher enfileira por MUDANÇA, e uma etapa toca vários arquivos: numa rodada real
a fila chegou a 26+ tasks para 8 entregas, e fechar uma a uma fez o orquestrador
desistir — a fila virou paisagem, que é o oposto do que ela existe para ser.`,
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
				fmt.Printf("task concluída: %s → %s\n", args[0], queue.DoneDir)
				return nil
			}
			if !all && doFile == "" && doKind == "" {
				return fmt.Errorf("informe o <id>, ou um filtro: --file <caminho>, --kind <tipo> ou --all")
			}
			tasks, err := queue.List(absRoot)
			if err != nil {
				return err
			}
			alvo := relTo(absRoot, doFile)
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
				fmt.Println("nenhuma task casou o filtro — veja a fila com `anchors queue`")
				return nil
			}
			fmt.Printf("%d task(s) concluída(s) → %s\n", fechadas, queue.DoneDir)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&doFile, "file", "", "fecha todas as tasks deste arquivo")
	cmd.Flags().StringVar(&doKind, "kind", "", "fecha todas as tasks deste kind (spec|code|feature|test|change…)")
	cmd.Flags().BoolVar(&all, "all", false, "fecha a fila inteira")
	return cmd
}

func newDropCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "drop <id>",
		Short: "Descarta uma task da fila sem concluí-la (remove; não arquiva)",
		Long: `Remove uma task viva (pending ou claimed) da fila — para lixo: tasks
obsoletas, duplicatas, ou um plano que caiu como 'triage' e você não quer tratar.
Diferente de 'done' (que arquiva em done/): drop apaga, não vira histórico.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if err := queue.Drop(absRoot, args[0]); err != nil {
				return err
			}
			fmt.Printf("task descartada: %s\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	return cmd
}

func newReclaimCmd() *cobra.Command {
	var root string
	var force bool
	cmd := &cobra.Command{
		Use:   "reclaim",
		Short: "Devolve à fila as tasks presas em claimed (worker morto)",
		Long: `Move de volta para 'pending' toda task que ficou 'claimed' — tipicamente
órfã de um worker que morreu sem fechar com 'done'. Rode após um crash para o
trabalho não ficar preso. As tasks voltam a ser puxáveis por 'anchors next'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			n, err := reclaimFn(force)(absRoot)
			if err != nil {
				return err
			}
			fmt.Printf("%d task(s) devolvida(s) à fila (claimed → pending)\n", n)
			// O ZERO precisa se explicar.
			//
			// Medido: o comando respondia "0 task(s) devolvida(s)" com uma task
			// visivelmente `claimed` na fila. O número estava certo — ela foi reivindicada
			// há minutos, dentro da janela de trabalho —, e o zero sozinho parece defeito.
			// Custou dois comandos para descartar.
			if !force {
				if r := queue.RecentlyHeld(absRoot); r > 0 {
					fmt.Printf("  (%d task(s) reivindicada(s) RECENTEMENTE ficaram — "+
						"alguém pode estar nelas agora.\n"+
						"   Se você sabe que o worker parou: `anchors reclaim --force`)\n", r)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().BoolVar(&force, "force", false,
		"devolve TAMBÉM o que um worker VIVO reivindicou (use só se souber que ele parou)")
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
	fmt.Printf("\n○ %d gate(s) informativo(s) declarado(s) — medem e não defendem.\n", informativos)
	fmt.Println("  `anchors status` mostra quais já estão limpos e podem virar bloqueantes.")
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
		if strings.Contains(t.Reason, "spec(s) deste plano") {
			return ""
		}
		return fmt.Sprintf(" — %d de %d spec(s) deste plano ainda não existem",
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
		return fmt.Errorf("workflow.repo vazio: no modo github ele é obrigatório — " +
			"inferir do remote faria o Anchors escrever noutro repositório quando alguém " +
			"trabalha num fork, e escrita em lugar errado não se desfaz com revert")
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
		fmt.Printf("retomando o seu card — TERMINE ele antes de pegar outro\n")
		fmt.Printf("  o contexto da sessão anterior vale mais que a ordem da fila.\n\n")
	} else {
		// SEM CARD: pede ao pipeline. Ele é serializado, e é isso que impede dois agentes
		// de receberem o mesmo trabalho (ver `board.Ask`).
		if err := cli.Ask(agent); err != nil {
			return err
		}
		fmt.Printf("trabalho pedido ao pipeline (claim serializado)\n")
		fmt.Printf("  o pipeline escolhe o card, comenta a posse e move para in-progress.\n\n")
		fmt.Printf("Leia o que recebeu:\n")
		fmt.Printf("    anchors next          # em alguns segundos, quando o claim rodar\n")
		fmt.Printf("    gh run list --workflow %s --limit 1   # ver o claim\n", board.ClaimWorkflow)
		return nil
	}

	fmt.Printf("card reivindicado: #%d — %s\n\n", card.Number, card.Title)
	fmt.Printf("  estado:   %s\n", strings.TrimPrefix(card.State, "anchors:"))
	fmt.Printf("  dono:     %s\n\n", agent)
	printBoardWork(root, card)
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
		alvo = "<a unidade do card>"
	}

	switch {
	case strings.Contains(card.Title, "Implementar plan"):
		fmt.Println("ENTREGÁVEL: todas as specs que o plano semeia.")
		fmt.Println("  Um trabalho só — o plano é a entrada, as specs são a saída.")
		fmt.Println()
		fmt.Printf("    anchors work spec --for <alvo que o plano lista>\n\n")
		fmt.Println("  O `-progress.md` do plano é onde você marca o que entregou, e o gate")
		fmt.Println("  `progress-honest` confronta cada `[x]` contra o disco.")

	case strings.Contains(card.Title, "Implementar spec"):
		fmt.Println("ENTREGÁVEL: código + feature + teste + documentação.")
		fmt.Println("  A spec é a ENTRADA deste card — criá-la foi o trabalho do card do")
		fmt.Println("  plano. Implementá-la é este, e são quatro entregas.")
		fmt.Println()
		fmt.Printf("    1. anchors work code    --for %s\n", alvo)
		fmt.Printf("    2. anchors work feature --for %s\n", alvo)
		fmt.Printf("    3. anchors work test    --for %s\n", alvo)
		fmt.Println("    4. a documentação evolui JUNTO: a cada alteração que deveria ser")
		fmt.Println("       documentada, atualize a doc na mesma entrega — deixá-la para o")
		fmt.Println("       fim é o que a faz nascer incompleta.")
		printDocDuties(root, unidade)
		fmt.Println()
		fmt.Println("  A ordem 1→2→3 não é gosto: a feature descreve o comportamento em")
		fmt.Println("  cenários, e é dela que os testes nascem (`anchors work test` parte da")
		fmt.Println("  feature, não da spec).")

	default:
		fmt.Println("ENTREGÁVEL: veja o corpo do card.")
	}

	// A RÉGUA DA AUTONOMIA no próprio card, e não só no guia.
	//
	// O guia é um comando que alguém precisa rodar; o card é o que o agente lê sempre. Um
	// agente que não decide o produto e não abriu o guia perguntaria ao dev — e é a porta
	// que mais se usa, porque não tem label nem gate.
	if !decidesProduct(root) {
		fmt.Println()
		fmt.Println("  VOCÊ NÃO DECIDE O RUMO DESTE PRODUTO (`.anchors/settings.yaml`).")
		fmt.Println("  Diante de ambiguidade ou de escolha que muda o comportamento:")
		fmt.Println("      anchors escalate \"<o que precisa ser decidido>\" --about <arquivo> --for-user")
		fmt.Println("  NÃO pergunte a quem está rodando você — a resposta é razoável e vira")
		fmt.Println("  decisão de produto de quem não tinha autoridade, sem rastro nenhum.")
	}

	fmt.Printf("\nAntes de começar:  anchors guide work\n")
	if unidade != "" {
		fmt.Printf("O que propagar:    anchors impact %s\n", unidade)
	}
	// O CORPO DO PR pelo comando, e não à mão.
	//
	// O `Closes #N` é o que fecha o card no merge, e ele é fácil de esquecer quando o
	// corpo é escrito à mão — medido: o card #319 ficou aberto em `in-progress` depois do
	// merge, e o #321 fechou sozinho, porque um PR tinha a linha e o outro não. O estado
	// do board passou a divergir do repositório sem nada acusar.
	fmt.Printf("Ao terminar:       anchors pr-body --cards %d  (traz o `Closes` que fecha o card)\n", card.Number)
	fmt.Printf("                   abra o PR com esse corpo — o pipeline move o card, você não\n")
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
	fmt.Printf("       Alterar `%s` OBRIGA tocar:\n", alvo)
	for _, d := range deveres {
		fmt.Print(doct.Duty(d))
	}
	fmt.Println("       `anchors docs duties --layer " + camada + "` repete isto a qualquer momento.")
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

	fmt.Println("ENTREGÁVEL: a REVISÃO — e ela não produz arquivo novo.")
	fmt.Println("  A implementação deste card já foi entregue. O que falta é confrontá-la")
	fmt.Println("  contra o que ela DIZ ter feito, e é por isso que a etapa some: nada")
	fmt.Println("  parece faltar. Os quatro arquivos estão lá e os gates estão verdes.")
	fmt.Println()
	fmt.Printf("    anchors work review --for %s\n\n", alvo)
	fmt.Println("  Os gates verdes NÃO provam que está certo — eles confrontam o que é")
	fmt.Println("  DECLARÁVEL. Em três rodadas de um E2E real, 7 defeitos graves (perda")
	fmt.Println("  silenciosa de dado, regra sem teste que a prove, contradição entre duas")
	fmt.Println("  regras da mesma spec) passaram com tudo verde. Nenhum veio de gate.")
	fmt.Println()
	fmt.Println("  O registro de entrega está nos comentários desta issue: é a metade")
	fmt.Println("  DECLARADA do confronto — o que o autor diz ter feito, contra o disco.")
	fmt.Println()
	fmt.Println("  Se a revisão não achar nada, diga isso e siga. Se achar, o achado vira")
	fmt.Println("  correção neste card — não card novo.")

	fmt.Printf("\nAntes de começar:  anchors guide review\n")
	// O CORPO DO PR pelo comando, e não à mão.
	//
	// O `Closes #N` é o que fecha o card no merge, e ele é fácil de esquecer quando o
	// corpo é escrito à mão — medido: o card #319 ficou aberto em `in-progress` depois do
	// merge, e o #321 fechou sozinho, porque um PR tinha a linha e o outro não. O estado
	// do board passou a divergir do repositório sem nada acusar.
	fmt.Printf("Ao terminar:       anchors pr-body --cards %d  (traz o `Closes` que fecha o card)\n", card.Number)
	fmt.Printf("                   abra o PR com esse corpo — o pipeline move o card, você não\n")
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
	if !terminalInterativo() {
		fmt.Println("(sem terminal para perguntar — este agente não tem perfil declarado,")
		fmt.Println(" e sem perfil ele não atua nos cards escalonados. Declare com")
		fmt.Println(" `anchors settings role`)")
		fmt.Println()
		return s, nil
	}
	perfil, err := askRole()
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
	printRole(perfil)
	fmt.Printf("\n  registrado em %s — não perguntarei de novo\n\n", settings.Path(root))
	return s, nil
}

// terminalInterativo diz se há alguém do outro lado para responder.
func terminalInterativo() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
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
