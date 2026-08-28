package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/scan"

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
			absRoot, err := config.AbsRaiz(root)
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
				fmt.Printf("    sugestão: %s — %s\n", t.SuggestedNext, t.Reason)
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
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			if worker == "" {
				worker = defaultWorkerID()
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
				if n, err := semearDosPlanos(absRoot); err == nil && n > 0 {
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
			fmt.Printf("  motivo:   %s\n\n", t.Reason)
			fmt.Printf("Execute o passo (veja `anchors guide`). Para o detalhe fino do que\n")
			fmt.Printf("propagar, rode: anchors impact %s\n", t.Changed)
			fmt.Printf("Ao terminar:    anchors done %s\n", t.ID)

			// A maturação (QUALITY §7) aparece aqui de forma BARATA: o `next` é chamado
			// pelo worker a cada task e precisa ser rápido, então não roda os gates —
			// só conta quantos estão declarados como informativos. Quem quer saber
			// quais estão limpos roda `anchors status` ou `check`, que já medem.
			lembraMaturacaoBarato(absRoot)
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
			absRoot, err := config.AbsRaiz(root)
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
			absRoot, err := config.AbsRaiz(root)
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
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			n, err := reclaimFn(force)(absRoot)
			if err != nil {
				return err
			}
			fmt.Printf("%d task(s) devolvida(s) à fila (claimed → pending)\n", n)
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

// reclaimFn escolhe entre respeitar o worker vivo (padrão) e ignorá-lo (--force).
// O padrão é respeitar: roubar trabalho por suposição produz dois agentes no mesmo
// arquivo, e o custo disso só aparece depois.
func reclaimFn(force bool) func(string) (int, error) {
	if force {
		return queue.ReclaimForce
	}
	return queue.Reclaim
}

// semearDosPlanos enfileira uma task por PLANO que ainda tem spec por nascer.
//
// "Ainda tem trabalho" é medido pelo disco, não por checkbox: um plano cujas specs
// semeadas TODAS existem já foi cumprido, e enfileirá-lo seria o ruído que faz a fila
// perder a confiança de quem a puxa. O que ele semeia sai do mesmo lugar que o gate
// `plan-seeds-valid` lê — os caminhos `.spec.md` citados no texto.
func semearDosPlanos(root string) (int, error) {
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
			if seedExiste(root, s, files) {
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
			Reason: fmt.Sprintf("%s — %d de %d spec(s) deste plano ainda não existem",
				reason, faltam, len(f.Seeds)),
		})
		if err == nil && criada {
			n++
		}
	}
	return n, nil
}

// seedExiste diz se a spec que o plano semeia já está no repositório.
//
// O plano cita de duas formas legítimas: pelo caminho, ou só pelo NOME do arquivo — que é
// como se escreve em prosa ("a spec de `SubscriptionScreen.spec.md`"). Medido: 10 das 26
// citações de um repositório real são por nome, e tratá-las como ausentes fazia o plano
// parecer eternamente não-cumprido — a fila semeava fases já concluídas na partida a frio.
//
// Por nome, só quando ele é ÚNICO: dois arquivos homônimos tornam a citação ambígua, e
// escolher um seria decidir pelo autor.
func seedExiste(root, seed string, files []scan.File) bool {
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

// lembraMaturacaoBarato conta os gates informativos sem RODAR nenhum.
//
// O `next` é chamado pelo worker a cada task, e rodar a suíte de gates ali dobraria o
// custo de puxar trabalho. O que ele pode fazer sem custo é ler a declaração: se há gate
// informativo, existe maturação pendente — e quem quiser saber QUAIS estão limpos roda
// `anchors status`, que já mede.
//
// É um lembrete mais fraco de propósito. Um lembrete caro num comando de laço quente é
// um lembrete que alguém vai querer desligar.
func lembraMaturacaoBarato(root string) {
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
