package quality

import (
	"errors"
	"fmt"
	"github.com/co2-lab/anchors/internal/i18n"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/cmd/anchors/mapcmd"
	"github.com/co2-lab/anchors/internal/checklog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/health"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

func newCheckCmd() *cobra.Command {
	var root, mapPath, phase, category string
	var changed []string
	var all, noRecord, fix, deterministic, skipSlow, onlyIssues, showDrift, showTiming bool
	var skipRegras string
	var msgPath string
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Run the quality gates (the pipeline)",
		Long: `Confronts the nodes against the quality gates declared in anchors.yaml
and reports the verdict profile (QUALITY §5-§6).

By default it is INCREMENTAL: it runs over the impact path of a changed file
(--changed <file>). With --all, it sweeps every node (expensive; the full picture).

Each gate that fails generates an issue; a blocking gate that fails bars
promotion. Informative gates enter the profile but do not bar (maturation, §7).

--deterministic SKIPS the AI judgment gates (measures: judgment) and runs only
the computable ones. It is the PRE-COMMIT mode: a judge gate cannot bar a commit
(there is no waiting for the AI) nor should it queue judgment on every commit (repeated
garbage). Without that mode, judge becomes invisible (it neither bars nor records).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			// No modo github o achado de gate vira CARD, não arquivo: o `issues/` é a fila do
			// modo local (mover pasta à mão), e manter os dois faz o board esconder o que os
			// gates encontraram.
			if cfg != nil && cfg.GitHubMode() && len(cfg.Workflow.Labels) > 0 {
				issue.UseGitHub(cfg.Workflow.Repo, cfg.Workflow.Labels[0])
			}
			if len(cfg.Gates) == 0 {
				return fmt.Errorf("no gate declared in anchors.yaml (`gates:` section)")
			}
			// --deterministic: remove os gates de julgamento por IA do conjunto.
			if deterministic {
				kept := cfg.Gates[:0:0]
				for _, gt := range cfg.Gates {
					if !gt.IsJudgment() {
						kept = append(kept, gt)
					}
				}
				cfg.Gates = kept
				if len(cfg.Gates) == 0 {
					fmt.Println(i18n.T("check.no_deterministic_gate"))
					return nil
				}
			}
			// Categorização (fase/custo/natureza): filtra QUAIS gates são cobrados agora.
			perspective := config.PerspectiveChange
			if all {
				perspective = config.PerspectiveAll
			}
			// A dispensa vem da flag OU do ambiente — o hook do git passa por variável,
			// porque `git commit` não repassa flags ao pre-commit.
			bruto := skipRegras
			if bruto == "" {
				bruto = os.Getenv("ANCHORS_SKIP_RULES")
			}
			dispensa, erros := gate.ParseWaiver(bruto)
			// A MENSAGEM DE COMMIT também dispensa: `[skip-trinca-completa@WRKSP: motivo]`.
			//
			// O caminho vem por `--commit-msg`, e não de `.git/COMMIT_EDITMSG`: MEDIDO, o
			// git NÃO grava esse arquivo antes do `pre-commit` — nem com `-m`. Lê-lo ali
			// devolve a mensagem do commit ANTERIOR, e a dispensa passa a valer para o
			// commit errado, em silêncio. É o hook `commit-msg` que recebe o arquivo, e é
			// ele quem passa este caminho.
			if msgPath != "" {
				if msg, err := os.ReadFile(msgPath); err == nil {
					daMsg, errosMsg := gate.WaiverFromMessage(string(msg))
					dispensa = dispensa.Merge(daMsg)
					erros = append(erros, errosMsg...)
				}
			}
			if len(erros) > 0 {
				return fmt.Errorf("invalid --skip-rule:\n  %s\n\nThe reason is mandatory: a "+
					"waiver without a written justification is indistinguishable from someone dodging a "+
					"gate that found a defect", strings.Join(erros, "\n  "))
			}
			cfg.Gates = filterGates(cfg.Gates, phase, category, skipSlow, perspective, dispensa)
			if len(cfg.Gates) == 0 {
				fmt.Println(i18n.T("check.no_gate_for_slice", phase, category))
				return nil
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}

			// BINÁRIO MAIS VELHO QUE O MAPA. Ele não falha — grava o formato que conhece,
			// e desfaz o que a versão nova escreveu. Medido: um build local anterior a um
			// rename de campo revertia 26 linhas a cada `check`, e o mapa ficava oscilando
			// entre dois formatos, com conflito a cada PR.
			//
			// O aviso não barra: o binário velho ainda faz o trabalho da versão dele, e
			// derrubar o CI trocaria "faz menos do que devia" por "não faz nada" — a mesma
			// razão do `doctor --check-pipelines`.
			if staleBinaryWarning(g.GeradoPor, common.Version) != "" {
				fmt.Fprintln(os.Stderr, staleBinaryWarning(g.GeradoPor, common.Version))
			}

			// Os arquivos que de fato MUDARAM, distintos do raio de impacto que o
			// `selectNodes` devolve. Um gate que julga a mudança precisa dos primeiros;
			// os demais, do raio. Ver Config.Alterados.
			cfg.Alterados = normalizeChanged(changed, absRoot)
			nodes, scope, err := selectNodes(g, cfg, all, changed, absRoot)
			if err != nil {
				return err
			}

			// --fix: aplica os reparos automáticos (self-healer) ANTES de confrontar,
			// para que o check seguinte já reflita o conserto.
			if fix {
				fixes := gate.Fix(cfg.Gates, nodes, absRoot)
				n := 0
				for _, fr := range fixes {
					if fr.Fixed {
						fmt.Println(i18n.T("check.fix_applied", fr.Gate, fr.Target, fr.Detail))
						n++
					}
				}
				fmt.Println(i18n.T("check.fixes_count", n))
				fmt.Println()
			}

			// Espelha a saída num arquivo para que ela possa ser RELIDA sem
			// re-executar. Abre AQUI, e não no topo: o que vem antes é erro de
			// invocação (config ausente, mapa ilegível), que não é relatório —
			// gravá-lo sobrescreveria uma foto boa com uma mensagem de erro.
			head, assunto, _ := gitmeta.Head(absRoot)
			espelho := checklog.Open(absRoot, all, checklog.Header(
				"anchors "+strings.Join(os.Args[1:], " "),
				head, assunto, gitmeta.DirtyCount(absRoot), time.Now(),
			))
			defer espelho.Close()

			warnIfMapStale(absRoot, mapPath, cfg, g)
			fmt.Println(i18n.T("check.header", scope, len(nodes), len(cfg.Gates)))
			fmt.Println()

			// `all` chega até os gates: é o que permite ao gate que sabe varrer sozinho
			// (`scope_full`) rodar UMA vez sem receber a lista, em vez de receber o projeto
			// inteiro em lotes.
			results := gate.RunWithWaiver(cfg.Gates, nodes, absRoot, g, cfg, all, dispensa)
			profile := gate.Aggregate(results)
			printProfile(profile, onlyIssues, showDrift)
			if showTiming {
				printTiming(profile)
			}
			warnGatesWithoutTarget(cfg.Gates, profile)

			// O LOOP: check → carimbo → issue. Deixa de "reportar" e passa a
			// "registrar": grava o veredito por aresta no mapa (destrava stale) e
			// abre uma issue de violation por fail bloqueante (sobrevive à sessão).
			// Opt-out honesto: --no-record só reporta, não registra.
			pendentes := 0
			if !noRecord {
				if err := recordCheck(absRoot, mapPath, g, profile); err != nil {
					fmt.Fprintf(os.Stderr, "warning: failed to record (stamp/issue): %v\n", err)
				}
				// gates de JULGAMENTO: enfileira uma task `judge` por alvo pendente,
				// para uma IA confrontar e reportar com `anchors judge`.
				if n := enqueueJudgments(absRoot, cfg, profile, all); n > 0 {
					fmt.Println(i18n.T("check.awaiting_ai_judgment", n))
				}
				// PENDENTES é o que está NA FILA, não o que acabou de ser enfileirado.
				//
				// `enqueueJudgments` conta só as tasks CRIADAS agora — na segunda
				// execução ela devolve 0 porque a task já existia. Barrar por esse número
				// deixaria passar exatamente o caso que importa: rodar o check, não
				// julgar, e commitar. Medido no projeto de referência.
				pendentes = queuedJudgments(absRoot)
			}

			if c := espelho.Path(); c != "" {
				rel, err := filepath.Rel(absRoot, c)
				if err != nil {
					rel = c
				}
				fmt.Println()
				fmt.Println(i18n.T("check.output_mirrored", rel))
			}

			// JULGAMENTO PENDENTE BARRA O COMMIT, e só o commit — não o `--all`.
			//
			// A task de julgamento é trabalho DO COMMIT ATUAL, não da fila do projeto:
			// quem mexeu no arquivo é quem tem o contexto para responder. Deixá-la
			// passar empurra para outro alguém uma pergunta que só quem editou sabe
			// responder — e ela some no meio das outras.
			//
			// Por isso ela não vira card: é resolvida ANTES de o trabalho sair da
			// máquina, e o `.anchors/` fica no `.gitignore`. O PR nasce sem pendência.
			//
			// Barra em `--changed` (a perspectiva do pre-commit) e não em `--all`: o
			// `--all` é a foto do projeto inteiro, e derrubá-lo por um julgamento que
			// ninguém acabou de criar tornaria impossível medir um projeto que já tem
			// pendência acumulada.
			if pendentes > 0 && !all {
				fmt.Println()
				fmt.Println(i18n.T("check.blocked_by_judgment", pendentes))
				espelho.Close()
				os.Exit(1)
			}

			// O MAPA ESTÁ VELHO? — e esta pergunta é diferente de "o arquivo está no mapa".
			//
			// O hook já barrava o arquivo REGIDO fora do mapa (o arquivo novo). O que
			// passava era o arquivo que ESTÁ no mapa com a `rev` de uma versão anterior:
			// o `map build` rodou, e o trabalho continuou depois dele.
			//
			// MEDIDO no projeto de referência: das 12 reprovações do pipeline `gates`,
			// SETE foram isto — a maior causa isolada. Nas sete o agente tinha commitado
			// o mapa; ninguém esqueceu de gerá-lo, todos geraram cedo demais.
			//
			// AVISO e não barreira, e a distinção importa. O mapa velho não invalida o
			// trabalho — invalida a FOTO que os gates confrontam. Barrar o commit aqui
			// impediria alguém de salvar trabalho em andamento, que é justamente quando
			// o mapa fica velho. O que faltava era DIZER, na máquina, o que o CI só diz
			// seis minutos depois.
			if stale := mapcmd.StaleMapNodes(absRoot, g, cfg); len(stale) > 0 && !all {
				fmt.Fprintln(os.Stderr)
				fmt.Fprintln(os.Stderr, i18n.T("check.map_stale", len(stale)))
				for i, v := range stale {
					if i == 3 {
						fmt.Fprintln(os.Stderr, i18n.T("check.map_stale_more", len(stale)-3))
						break
					}
					fmt.Fprintln(os.Stderr, "    "+v)
				}
			}

			// DICAS DE GOVERNANÇA (informativo, QUALITY §5.2).
			// Apresenta oportunidades de evolução do projeto (gates canônicos ausentes ou
			// configurações subótimas) sem quebrar commits ou inflar warnings.
			if opps := health.QuickGovernanceHints(g, cfg); len(opps) > 0 && all {
				fmt.Println()
				for _, op := range opps {
					fmt.Printf("ℹ %s (%s): %s\n", i18n.T("check.governance_tip"), op.Subject, op.Detail)
				}
				fmt.Printf("  %s\n", i18n.T("check.doctor_learn_more"))
			}

			if !profile.Passed {
				// `os.Exit` não roda os `defer`: sem fechar aqui, o espelho perderia
				// o fim do relatório exatamente no caso em que ele mais importa — o
				// da reprovação.
				espelho.Close()
				os.Exit(1) // barra: há fail bloqueante
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().StringSliceVar(&changed, "changed", nil, "changed file(s) — repeatable or comma-separated; the gates run ONCE over the union of the impact paths")
	cmd.Flags().BoolVar(&all, "all", false, "scan every node (the full picture; expensive)")
	cmd.Flags().BoolVar(&noRecord, "no-record", false, "report only: neither stamps the map nor opens issues")
	cmd.Flags().BoolVar(&fix, "fix", false, "self-healer: applies the automatic repairs (e.g. fixes updated_at) before confronting")
	cmd.Flags().BoolVar(&deterministic, "deterministic", false, "runs only the computable gates (skips the AI-judgment ones) — the pre-commit mode")
	cmd.Flags().StringVar(&phase, "phase", "", "enforces only the gates of this phase (pre-commit|pre-push|ci|manual)")
	cmd.Flags().StringVar(&category, "category", "", "enforces only the gates of this nature (types|style|traceability…)")
	cmd.Flags().BoolVar(&skipSlow, "skip-slow", false, "skip gates declared `cost: slow`")
	cmd.Flags().StringVar(&skipRegras, "skip-rule", "",
		"waives rules in this run: `id=reason[,id=reason]`. The reason is mandatory")
	cmd.Flags().StringVar(&msgPath, "commit-msg", "",
		"path to the commit message file, from which to read the `[skip-regra@CODIGO: motivo]` markers")
	cmd.Flags().BoolVar(&onlyIssues, "only-issues", false, "omits from the table the gates that passed everything and left no pending item (the total remains in the footer)")
	cmd.Flags().BoolVar(&showDrift, "show-drift", false, "lists ALL the pending items (⚠) with the address of each one; without the flag, only the table counter")
	cmd.Flags().BoolVar(&showTiming, "timing", false, "measures how long each gate took, and names the slowest targets — to find what makes a scan expensive")
	return cmd
}

// filterGates aplica a CATEGORIZAÇÃO: fase, natureza, custo e perspectiva. Cada eixo é
// independente — `when` diz em que momento o gate é cobrado, `cost` diz se cabe num
// loop apertado, `category` diz o que ele mede, `skip_on` diz sobre QUANTO do projeto a
// resposta dele tem valor.
//
// Todos os filtros são permissivos por omissão: gate sem `when`/`cost`/`category`/
// `skip_on` continua rodando como antes. É o que torna a categorização adotável aos
// poucos — declarar um eixo num gate não muda o comportamento dos outros 30.
func filterGates(gates []config.Gate, phase, category string, skipSlow bool, perspective string, dispensa gate.Waiver) []config.Gate {
	out := gates[:0:0]
	for _, g := range gates {
		if !g.RunsIn(phase) {
			continue
		}
		if category != "" && g.Category != category {
			continue
		}
		if skipSlow && g.IsSlow() {
			continue
		}
		// A perspectiva NÃO é escolha do usuário: vem de COMO o check foi chamado
		// (`--all` ou `--changed`). O gate declara em qual delas ele se abstém.
		if g.SkipsOn(perspective) {
			continue
		}
		// DISPENSA por ID: o usuário declarou, com motivo, que esta regra não deve ser
		// confrontada nesta execução. Diferente de `skip_on`, que é decisão do gate sobre
		// perspectiva, esta é decisão de quem chama sobre uma regra específica — e é o que
		// permite commitar a primeira spec de uma unidade (a feature ainda é um card) sem
		// desligar os gates que verificam outra coisa.
		if motivo, ok := dispensa.Waived(gate.RuleID(gateCmdID(g))); ok {
			fmt.Println(i18n.T("check.waived_rule", gateCmdID(g), motivo))
			continue
		}
		out = append(out, g)
	}
	return out
}

// enqueueJudgments cria uma task `judge` por alvo que um gate de julgamento marcou
// como pendente (verdict Judge). A task carrega o gate e o guide, para o worker (a
// IA) saber o que ler e confrontar. Reusa a fila. Idempotente pelo dedup da fila.
func enqueueJudgments(root string, cfg *config.Config, p gate.Profile, varreduraCompleta bool) int {
	// índice gate → (guide, ask) para enriquecer a task
	guideOf := map[string]string{}
	askOf := map[string]string{}
	knownJudgmentGates = knownJudgmentGates[:0]
	for _, g := range cfg.Gates {
		if g.IsJudgment() {
			guideOf[g.Name] = g.Guide
			askOf[g.Name] = g.Ask
			knownJudgmentGates = append(knownJudgmentGates, g.Name)
		}
	}
	// NO MODO `github` A FILA É O BOARD, e o julgamento não foge disso.
	//
	// É a terceira porta do mesmo defeito: o `next` a tinha (v0.1.55), o `deliver` a
	// tinha (v0.1.61), e o `check` continuava enfileirando julgamento em
	// `.anchors/tasks/`. O efeito é o mesmo — a fila local renasce a cada `check`,
	// sozinha, e o `doctor` a acusa como resíduo de uma migração que nunca terminou.
	//
	// Medido no projeto de referência: 76 tasks de julgamento pendentes, TODAS criadas
	// no mesmo dia pelos `check` de uma sessão, num projeto declarado `mode: github`.
	// Ninguém as puxaria: o `next` naquele modo lê o board.
	//
	// O julgamento que REPROVA já vira issue pelo caminho normal (o `check` a abre).
	// O que fica de fora é o julgamento PENDENTE — e registrá-lo numa fila que ninguém
	// lê é pior que não registrar, porque parece registro.
	if cfg.GitHubMode() {
		if len(p.Judged) > 0 {
			fmt.Println()
			fmt.Println(i18n.T("check.github_judgment_pending", len(p.Judged)))
			fmt.Println(i18n.T("check.github_judgment_mode"))
		}
		return 0
	}
	n := 0
	for _, r := range p.Judged {
		reason := fmt.Sprintf("gate '%s'", r.Gate)
		if gd := guideOf[r.Gate]; gd != "" {
			reason += " — leia " + gd
		}
		if ask := askOf[r.Gate]; ask != "" {
			reason += " — pergunta: " + strings.TrimSpace(ask)
		}
		t := queue.Task{
			ID:            "judge-" + r.Gate + "-" + strings.ReplaceAll(relSlug(r.Target), "/", "-"),
			Changed:       r.Target,
			Kind:          "judgment",
			Origin:        "check",
			SuggestedNext: "judge",
			Reason:        reason,
			CreatedAt:     time.Now().Format(time.RFC3339),
		}
		if created, _ := queue.Enqueue(root, t); created {
			n++
		}
	}
	// A limpeza SÓ vale na varredura completa.
	//
	// `dropStaleJudgments` conclui "não foi enfileirado agora, então é obsoleto". Isso é
	// verdade quando o check olhou TODOS os nós; em `--changed` ele olhou um arquivo, e
	// todos os outros alvos do mesmo gate parecem obsoletos por não terem sido olhados.
	//
	// Medido no blue-eyes: dois `check --changed` seguidos em testes diferentes deixaram
	// UMA task na fila. O primeiro julgamento foi apagado pelo segundo check — e o
	// `judge --pending` respondia "nenhum alvo aguardando" com cinco pendentes no
	// `check --all`. Quem confia na fila para saber o que julgar perde trabalho em
	// silêncio, que é o oposto do que esta função existe para fazer.
	if varreduraCompleta {
		dropStaleJudgments(root, cfg, p)
	}
	return n
}

// dropStaleJudgments tira da fila as tasks de julgamento cujo alvo o gate
// NÃO enfileirou nesta rodada.
//
// A fila é persistente e o conjunto de alvos aplicáveis não é: um gate que passa a
// declarar `requires`, um alvo que perde a marcação que o tornava aplicável, uma spec
// renomeada — em todos os casos a task velha sobrevive e a fila passa a mentir. Foi o
// que aconteceu no app de referência: o `no-test-prova-real` ganhou o filtro, o check passou a
// reportar 16 alvos, e `judge --pending` seguia listando 583 porque as tasks de dias
// antes continuavam ali.
//
// Só mexe em task de julgamento vinda do check (`kind: judgment`, `origin: check`), e
// só do gate que rodou nesta rodada — não toca trabalho de outra origem nem de gate
// que não foi cobrado agora.
func dropStaleJudgments(root string, cfg *config.Config, p gate.Profile) {
	rodou := map[string]bool{}
	for _, g := range cfg.Gates {
		if g.IsJudgment() {
			rodou[g.Name] = true
		}
	}
	vivos := map[string]bool{} // "gate\x00alvo" que o check acabou de enfileirar
	for _, r := range p.Judged {
		vivos[r.Gate+"\x00"+r.Target] = true
	}
	tasks, err := queue.List(root)
	if err != nil {
		return
	}
	for _, t := range tasks {
		if t.Kind != "judgment" || t.Origin != "check" {
			continue
		}
		gateName := gateDaTaskJudge(t.ID)
		if gateName == "" || !rodou[gateName] {
			continue
		}
		if vivos[gateName+"\x00"+t.Changed] {
			continue
		}
		_ = queue.Drop(root, t.ID)
	}
}

// gateDaTaskJudge extrai o nome do gate do ID `judge-<gate>-<slug-do-alvo>`.
//
// O ID é montado com `-` como separador e o nome do gate também os contém
// (`no-test-prova-real`), então não dá para partir pelo separador: a leitura é por
// PREFIXO conhecido, comparando com os gates de julgamento declarados.
func gateDaTaskJudge(id string) string {
	const pref = "judge-"
	if !strings.HasPrefix(id, pref) {
		return ""
	}
	resto := id[len(pref):]
	for _, g := range knownJudgmentGates {
		if strings.HasPrefix(resto, g+"-") {
			return g
		}
	}
	return ""
}

// knownJudgmentGates é preenchido por enqueueJudgments a cada rodada — o
// nome do gate vem da config, não de uma lista fixa.
var knownJudgmentGates []string

// relSlug reduz um caminho a algo usável em ID de task.
func relSlug(p string) string { return strings.TrimSuffix(p, filepath.Ext(p)) }

// recordCheck fecha o loop: carimba as arestas confrontadas no mapa (persistindo o
// veredito, o que destrava a detecção de stale) e abre uma issue de violation por
// fail bloqueante. É a passagem de "reporta" para "registra" (QUALITY §5).
func recordCheck(root, mapPath string, g *mapx.Graph, p gate.Profile) error {
	now := time.Now()
	day := now.Format("2006-01-02")

	// 1) CARIMBO — converte os vereditos por nó (do gate) para o mapa e grava.
	gv := p.NodeVerdicts()
	verdicts := make([]mapx.NodeVerdict, len(gv))
	for i, v := range gv {
		verdicts[i] = mapx.NodeVerdict{ID: v.ID, Failed: v.Failed}
	}
	// A DATA, e não o instante. `changed_at` responde "quando esta relação foi
	// confrontada" — uma pergunta de auditoria — e não entra na regra de staleness, que
	// compara `rev` (ver PROPAGATION.md §3). Com precisão de segundo, cada `anchors
	// check` reescrevia as 26 linhas de carimbo do mapa: conflito em todo PR onde duas
	// pessoas rodaram o check, e um diff que muda sozinho sem dizer nada.
	//
	// Com a data, o mapa só muda quando o VEREDITO ou a `rev` mudam — que é quando há o
	// que registrar.
	stamped := g.StampEdges(verdicts, now.Format(time.DateOnly))
	if err := mapx.Save(g, mapPath); err != nil {
		return fmt.Errorf("save stamped map: %w", err)
	}

	// 2) ISSUES — o loop completo, abre E fecha:
	//   • fail bloqueante → ABRE uma violation (se ainda não existe).
	//   • pass de um (gate, alvo) que tinha issue → RESOLVE (move para done/).
	// A chave estável (gate+alvo+kind) é o que liga o pass de hoje à issue de ontem.
	opened, resolved := 0, 0
	adiadas := 0
	for _, r := range p.Results {
		if r.Verdict == gate.Skip {
			continue // gate não se aplica — não mexe na issue
		}
		iss := issue.Issue{
			Kind: issue.Violation, Target: r.Target, Gate: r.Gate,
			Detail: r.Detail, Date: day,
		}
		// DÍVIDA ASSUMIDA vira issue em `future/`. Enquanto era só `Pending`, ela ficava
		// como uma linha no cabeçalho de um arquivo — visível para quem o abrisse, sem
		// estado, sem como ser paga, sem como vencer. O gate dizia que quem declara
		// afirma três coisas ("conhece o dever, ele vale, e QUANDO será pago"), e o
		// "quando" era prosa livre que nada confrontava.
		//
		// Materializada como issue, ela ganha o ciclo de vida que já existe: mover para
		// `todo/` quando chega a hora, e o próprio confronto a fecha quando o dever é
		// cumprido. `future/` e não `todo/` porque quem lê `todo/` pergunta "o que faço
		// AGORA" — afogar essa lista com o que só vence depois é o caminho mais curto
		// para ninguém mais olhar nenhuma das duas.
		if r.Verdict == gate.Pending {
			// DECISÃO EM ABERTO vira issue em `todo/`, e não em `future/`: ela só é
			// resolvida se alguém a VIR e a levar a quem decide. `future/` é o que vence
			// depois — a pergunta não vence, ela trava quem for implementar.
			//
			// É a issue que o usuário fecha: respondendo direto nela, ou pedindo à IA que
			// liste as perguntas abertas. Quando a resposta virar regra e o item sair da
			// spec, o `Pass` do próximo confronto a resolve sozinho, pelo mesmo caminho
			// que já fecha as violações.
			if r.Decisão {
				issDec := iss
				issDec.Kind = issue.Decision
				// DO USUÁRIO: a resposta não está no código, e o agente não a tem. Deixá-la
				// como dele faria o agente retentá-la para sempre — ou, pior, decidir por
				// conta própria, que é exatamente o que este gate existe para evitar.
				issDec.Dono = issue.DonoUsuário
				created, at, err := issue.Open(root, issDec)
				if err != nil {
					return fmt.Errorf("open decision issue: %w", err)
				}
				if created {
					opened++
				} else if at != issue.Todo {
					fmt.Println(i18n.T("check.record_decision_exists", r.Target, at))
				}
				continue
			}
			if !r.Divida {
				continue // indeterminado, não é dívida de ninguém
			}
			iss.Prazo = r.Prazo
			created, at, err := issue.OpenAt(root, iss, issue.Future)
			if err != nil {
				return fmt.Errorf("open debt issue: %w", err)
			}
			switch {
			case created:
				adiadas++
			case at == issue.Todo || at == issue.Doing:
				// Alguém puxou a dívida para o trabalho de agora. Dizer isso importa: a
				// declaração no header continua dizendo "depois", e o repositório já diz
				// "agora" — quem lê só o header não saberia.
				fmt.Println(i18n.T("check.record_debt_exists", r.Gate, r.Target, at))
			}
			continue
		}
		switch r.Verdict {
		case gate.Fail:
			if r.Blocking { // só fail BLOQUEANTE vira issue (informativo não barra nem registra)
				created, at, err := issue.Open(root, iss)
				if err != nil {
					return fmt.Errorf("open issue: %w", err)
				}
				if created {
					opened++
				} else if at != issue.Todo {
					fmt.Println(i18n.T("check.record_issue_exists", r.Gate, r.Target, at))
				}
			}
		case gate.Pass:
			// passou: se havia issue aberta para este (gate, alvo), fecha-a.
			//
			// O KIND faz parte da chave, e o `iss` acima nasce como `violation`. Um gate
			// que abre issue de outro kind precisa fechá-la pelo kind com que abriu — sem
			// isto a issue de decisão nunca fechava: o gate voltava a ✓ e a pergunta
			// continuava em `todo/` para sempre, que é pior que não ter aberto.
			if d := iss; d.Gate == "open-questions-resolved" {
				d.Kind = issue.Decision
				if ok, err := issue.Resolve(root, d.Key()); err != nil {
					return fmt.Errorf("resolve decision issue: %w", err)
				} else if ok {
					resolved++
					fmt.Println(i18n.T("check.record_decision_resolved", r.Target))
				}
			}
			ok, err := issue.Resolve(root, iss.Key())
			if err != nil {
				return fmt.Errorf("resolve issue: %w", err)
			}
			if ok {
				resolved++
				fmt.Println(i18n.T("check.record_issue_resolved", r.Gate, r.Target, issue.Dir))
			}
		}
	}

	fmt.Println()
	fmt.Println(i18n.T("check.record_summary", stamped, opened, resolved, issue.Dir))
	if adiadas > 0 {
		fmt.Println(i18n.T("check.record_debt_future", adiadas, issue.Dir))
	}
	return nil
}

// ExitNotGoverned é o código de saída para "este caminho não é regido pela Estrutura".
// Não é uma reprovação — é a resposta "não tenho jurisdição sobre isto".
//
// Existe como CÓDIGO, e não como texto a ser grepado, porque o pre-commit precisa
// distinguir isto de uma reprovação real. O hook antes fazia `grep "não está no mapa"`
// na saída: um casamento de string frágil que passava a valer para as DUAS situações
// assim que elas compartilharam uma mensagem — e o furo entrou por aí.
const ExitNotGoverned = 3

// errNotGoverned sinaliza o caminho não-regido. Carrega o alvo para a mensagem, e é
// reconhecida em main() para virar ExitNaoRegido em vez do exit 1 genérico.
type errNotGoverned struct{ target string }

func (e errNotGoverned) Error() string {
	return i18n.T("not_governed", e.target)
}

func (e errNotGoverned) As(target any) bool {
	if nr, ok := target.(*common.ErrNotGoverned); ok {
		*nr = common.ErrNotGoverned{Path: e.target}
		return true
	}
	return false
}

// selectNodes decide o conjunto de nós a confrontar: todos (--all) ou o caminho
// de impacto de um arquivo alterado (--changed, incremental).
func selectNodes(g *mapx.Graph, cfg *config.Config, all bool, changed []string, root string) (nodes []mapx.Node, scope string, err error) {
	if all {
		return g.Nodes, "--all", nil
	}
	if len(changed) == 0 {
		return nil, "", fmt.Errorf("provide --changed <file> (incremental) or --all (everything)")
	}

	// UNIÃO dos caminhos de impacto de todos os arquivos, resolvida numa passagem só.
	// Antes o `--changed` aceitava um arquivo, e quem tinha muitos (o pre-commit) era
	// obrigado a chamar o binário N vezes — recarregando config e mapa a cada volta,
	// ~1,2s por arquivo. Pior que lento: os gates relacionais (feature-test-match,
	// trinca-completa) confrontam a UNIDADE, então rodavam repetidos sobre o mesmo
	// conjunto, uma vez por peça dela.
	vistos := map[string]bool{}
	var ordem []string
	var naoRegidos int
	for _, c := range changed {
		ids, err := impactOf(g, cfg, c, root)
		if err != nil {
			var nr errNotGoverned
			if errors.As(err, &nr) {
				// Não-regido não contamina o lote: o Anchors só não tem jurisdição sobre
				// ele. Com TODOS não-regidos, o erro sobe (o chamador decide o exit 3).
				naoRegidos++
				continue
			}
			return nil, "", err
		}
		for _, id := range ids {
			if !vistos[id] {
				vistos[id] = true
				ordem = append(ordem, id)
			}
		}
	}
	if len(ordem) == 0 && naoRegidos == len(changed) {
		// Todos fora da jurisdição. A mensagem nomeia UM alvo e conta o resto: juntar
		// os caminhos numa string só produzia um erro ilegível de 8 KB num commit de 80
		// arquivos, e ainda dava a impressão de existir um arquivo com aquele nome.
		alvo := changed[0]
		if len(changed) > 1 {
			alvo = fmt.Sprintf("%s (and %d more)", changed[0], len(changed)-1)
		}
		return nil, "", errNotGoverned{target: alvo}
	}
	for _, n := range g.Nodes {
		if vistos[n.ID] {
			nodes = append(nodes, n)
		}
	}
	if len(changed) == 1 {
		return nodes, "--changed " + changed[0], nil
	}
	return nodes, fmt.Sprintf("--changed (%d files)", len(changed)), nil
}

// impactOf resolve UM arquivo alterado nos ids de nó que ele arrasta (o alvo, o que
// propaga a partir dele, o que o valida, e as peças da mesma unidade).
func impactOf(g *mapx.Graph, cfg *config.Config, changed, root string) ([]string, error) {
	target := common.RelTo(root, changed)
	if !common.NodeExists(g, target) {
		if _, statErr := os.Stat(filepath.Join(root, target)); statErr != nil {
			return nil, fmt.Errorf("%q exists neither on disk nor in the map — check the path", target)
		}
		// O arquivo EXISTE e não está no mapa. Duas situações COMPLETAMENTE distintas
		// se escondiam aqui sob uma mensagem só, e o pre-commit tratava ambas como
		// benignas — deixando passar exatamente o que o Anchors existe para barrar:
		// um arquivo REGIDO novo (hook, tela, service) commitado sem spec/feature/teste.
		// A distinção é a Estrutura: se o caminho casa uma camada do `layers:`, ele é
		// regido; senão, é um arquivo qualquer do repo (package.json, lockfile, CI).
		//
		// O IGNORE vem ANTES da classificação, e a ordem importa: `issues/` e `changes/`
		// casam a camada `doc` (`**/*.md`), mas o scanner NUNCA os indexa — são a saída
		// do próprio Anchors, não material a trabalhar. Classificar sem consultar o
		// ignore criava um impasse sem saída: o arquivo era dito "regido", não estava no
		// mapa, e `map build` não o acrescentava nunca — o commit ficava barrado para
		// sempre. Medido ao commitar as issues que o próprio `check` tinha resolvido.
		// `SkipDir` (e não `SkipFile`) porque é lá que mora `registroDoAnchors` — a
		// decisão é sobre o DIRETÓRIO-RAIZ do caminho (`issues/…`, `changes/…`).
		ig := scan.LoadIgnoreFor(root, cfg)
		if raiz, _, achou := strings.Cut(filepath.ToSlash(target), "/"); achou && ig.SkipDir(raiz, raiz) {
			return nil, errNotGoverned{target: target}
		}
		if ig.SkipFile(target) {
			return nil, errNotGoverned{target: target}
		}
		// O `-progress.md` é o MESMO impasse que o comentário acima descreve, por outra
		// porta: ele casa a camada `plan` (`plans/*.md`) e o scanner NUNCA o indexa — de
		// propósito, porque um arquivo que existe para mudar não pode ser confrontado por
		// gates que cobram justificativa de mudança (ver `scan.IsProgressFile`).
		//
		// Sem esta linha o resultado era o descrito ali: o arquivo dito "regido", ausente
		// do mapa, e `map build` não o acrescentando nunca. Medido no blue-eyes ao
		// commitar os 17 progressos que o `anchors new progress` acabara de criar — o
		// commit ficava barrado para sempre, pelo próprio mecanismo que separou os dois.
		if scan.IsProgressFile(target) {
			return nil, errNotGoverned{target: target}
		}
		if layer, _ := scan.Classify(target, cfg); layer == "" {
			return nil, errNotGoverned{target: target}
		}
		return nil, fmt.Errorf("%q is GOVERNED (a `layers:` layer) but is not in the map — "+
			"run `anchors map build` first (it is the step that registers the new file). "+
			"While it stays out of the map, NO gate confronts it: the triad is not "+
			"enforced and the pipeline certifies work that was never verified", target)
	}
	imp := g.AnalyzeImpact(target)
	// o alvo + o que propaga a partir dele (os filhos a refazer). A validação
	// (subir) é confronto contra os pais — entra como alvos também, para os gates
	// que se aplicam a eles rodarem.
	ids := map[string]bool{target: true}
	for _, n := range imp.Propagate {
		ids[n] = true
	}
	for _, n := range imp.Validate {
		ids[n] = true
	}
	// As PEÇAS DA MESMA UNIDADE entram junto, sempre. O impacto é calculado por arestas, e
	// spec/feature/teste/código de uma unidade nem sempre estão todos ligados por elas —
	// então `--changed <alvo>.ts` deixava de fora a `.feature` e o `.test.ts` irmãos.
	//
	// A consequência é um verde que não significa nada: medido, o `--changed` do código
	// não dispara `feature-test-match`, `feature-nao-vazia` nem `teste-nao-vazio`, e um
	// worker que siga o comando prescrito pelo `anchors work` declara "todos os
	// bloqueantes ✗0" sem jamais ter rodado o bloqueante mais importante da sua etapa.
	for _, peca := range unitPieces(target) {
		if common.NodeExists(g, peca) {
			ids[peca] = true
		}
	}
	out := make([]string, 0, len(ids))
	for id := range ids {
		out = append(out, id)
	}
	return out, nil
}

// nameWidth é a coluna do nome do gate: o maior nome presente, com um piso.
//
// Era `%-20s` fixo, e cinco gates passam disso (`handler-ddb-inline-passivo` tem
// 26). O nome mais longo empurrava a coluna do veredito e desalinhava a tabela
// inteira — numa lista de 49 linhas, o desalinhamento é o que faz o olho perder
// a coluna que importa.
func nameWidth(nomes []string) int {
	w := 18
	for _, n := range nomes {
		if len([]rune(n)) > w {
			w = len([]rune(n))
		}
	}
	return w
}

// larguras dos contadores: uma por COLUNA, pelo maior número daquela coluna.
//
// Uma largura única para todas desperdiça espaço onde ele não é preciso: numa
// varredura real o `~` chega a 582 e o `✗` fica em 0 ou 1, e a largura comum
// obrigaria a coluna dos fails a reservar três casas para nada. Cada coluna com
// a sua mantém os números alinhados à direita — que é o que permite compará-los
// a olho — sem esticar a tabela.
type counterWidths struct{ pass, fail, drift, skip, judge int }

func places(n int) int { return len(fmt.Sprint(n)) }

func computeWidths(p gate.Profile) counterWidths {
	// `drift` nasce em 0 — e continua 0 se nenhum gate tiver drift, que é o
	// sinal para a coluna inteira não existir. Os outros têm piso 1: eles sempre
	// aparecem, e `%*d` com largura 0 imprimiria colado no símbolo.
	w := counterWidths{pass: 1, fail: 1, drift: 0, skip: 1, judge: 1}
	max := func(atual, n int) int {
		if c := places(n); c > atual {
			return c
		}
		return atual
	}
	for nome, s := range p.ByGate {
		d := driftCount(p, nome)
		w.pass = max(w.pass, s.Pass)
		w.fail = max(w.fail, s.Fail)
		// `max` não serve para o drift: `casas(0)` é 1, e a coluna ganharia
		// largura mesmo sem drift nenhum — que é justamente o sinal de que ela
		// não deve existir. Só um drift REAL abre a coluna.
		if d > 0 {
			w.drift = max(w.drift, d)
		}
		w.skip = max(w.skip, s.Skip+s.Pending-d)
		w.judge = max(w.judge, s.Judge)
	}
	return w
}

// driftColumn devolve a célula do ⚠. A decisão é da TABELA, não da linha:
//
//   - se algum gate tem drift, a coluna existe em TODAS as linhas — vazia vira
//     branco do mesmo tamanho. Omiti-la só nas linhas sem drift empurraria o `~`
//     para a esquerda nelas, e a mesma coluna passaria a existir em dois lugares.
//   - se NENHUM tem, a coluna não existe em lugar nenhum. Reservá-la aí deixaria
//     um buraco no meio de todas as linhas sem nada que o justificasse — e o
//     caso é o comum, não a exceção: no `--changed` a tabela costuma sair
//     inteira sem drift.
//
// `largura == 0` é o sinal de "a tabela não tem drift nenhum".
func driftColumn(drift, largura int) string {
	if largura == 0 {
		return ""
	}
	if drift == 0 {
		// `len("⚠")` são 3 BYTES, mas o símbolo ocupa 1 coluna no terminal.
		// Usar `len` aqui reservava dois espaços a mais e desalinhava justamente
		// as linhas que o branco existia para alinhar.
		return strings.Repeat(" ", largura+1)
	}
	return fmt.Sprintf("⚠%*d", largura, drift)
}

// separador some junto com a coluna: sem isso, a tabela sem drift ficaria com
// dois espaços a mais entre `✗` e `~`.
func driftSeparator(largura int) string {
	if largura == 0 {
		return ""
	}
	return "  "
}

// cleanGate: passou em tudo que olhou e não deixou nada pendente. É o gate que
// não pede nada de ninguém — o candidato a sumir sob `--only-issues`.
func cleanGate(s gate.GateSummary, drift int) bool {
	return s.Fail == 0 && drift == 0 && s.Skip+s.Pending == 0 && s.Judge == 0
}

// printDrift lista as pendências agrupadas por GATE e, dentro dele, por MOTIVO.
//
// A lista crua tratava 2.430 pendências como 2.430 problemas, e elas não são:
// medido no app de referência, 832 delas têm o motivo IDÊNTICO ("sem sinal de mutação
// ingerido") e são um problema só — ninguém rodou a ingestão. Repetir o mesmo
// parágrafo 832 vezes esconde essa leitura em vez de revelá-la.
//
// Onde o motivo se repete, ele é escrito UMA vez e os alvos vêm compactados numa
// linha; onde cada motivo é único (`feature-test-match`, em que a divergência é
// específica de cada unidade), a lista continua alvo a alvo. O agrupamento não
// pode custar o endereço — ele é a parte acionável.
func printDrift(drifts []gate.Result) {
	if len(drifts) == 0 {
		return
	}
	// Agrupa preservando a ordem de aparição: gate → motivo → alvos.
	type grupo struct {
		motivo string
		alvos  []string
	}
	ordemGates := []string{}
	porGate := map[string][]*grupo{}
	total := map[string]int{}
	for _, r := range drifts {
		if _, visto := porGate[r.Gate]; !visto {
			ordemGates = append(ordemGates, r.Gate)
		}
		total[r.Gate]++
		var g *grupo
		for _, cand := range porGate[r.Gate] {
			if cand.motivo == r.Detail {
				g = cand
				break
			}
		}
		if g == nil {
			g = &grupo{motivo: r.Detail}
			porGate[r.Gate] = append(porGate[r.Gate], g)
		}
		g.alvos = append(g.alvos, r.Target)
	}

	fmt.Println()
	fmt.Println(i18n.T("check.drift_summary", len(drifts), len(ordemGates)))
	for _, nome := range ordemGates {
		fmt.Printf("\n  %s — %d\n", nome, total[nome])
		for _, g := range porGate[nome] {
			if len(g.alvos) == 1 {
				fmt.Printf("    ⚠ %s\n", g.alvos[0])
				if g.motivo != "" {
					fmt.Print(indent(g.motivo, "        "))
				}
				continue
			}
			// Motivo repetido: escrito UMA vez, e os alvos listados um por LINHA.
			//
			// Juntá-los com vírgula economizava altura e custava legibilidade: no
			// projeto que originou isto, 583 alvos viraram uma linha de 36 MIL
			// caracteres — o endereço estava lá e ninguém conseguia lê-lo, nem
			// grepar por um arquivo específico. Uma linha por alvo é o formato que
			// as ferramentas de texto esperam.
			if g.motivo != "" {
				fmt.Print(indent(g.motivo, "    "))
			}
			fmt.Println(i18n.T("check.drift_targets_count", len(g.alvos)))
			for _, alvo := range g.alvos {
				fmt.Printf("        %s\n", alvo)
			}
		}
	}
}

// printLegenda explica os símbolos da tabela — e SÓ os que ela usou. Uma legenda
// fixa listaria `⚠` e `⏳` em varreduras que não os têm, e explicar coluna
// ausente é ruído que ensina a pular a legenda inteira.
//
// A distinção que ela precisa carregar é entre `✗` e `~`: falha é o gate tendo
// confrontado e divergido; indeterminado é ele não ter tido o que confrontar.
// Sem isso, `~582` parece um débito de 582 itens, quando é o contrário — é a
// medida de quanto daquele gate não se aplica ali.
func printLegenda(p gate.Profile, w counterWidths) {
	temJudge := false
	for _, s := range p.ByGate {
		if s.Judge > 0 {
			temJudge = true
			break
		}
	}
	partes := []string{
		"✓  " + i18n.T("check.legend.pass"),
		"✗  " + i18n.T("check.legend.fail"),
	}
	if w.drift > 0 {
		partes = append(partes, "⚠  "+i18n.T("check.legend.warn"))
	}
	partes = append(partes, "~  "+i18n.T("check.legend.skip"))
	if temJudge {
		partes = append(partes, "⏳ "+i18n.T("check.legend.judge"))
	}
	// Uma por linha. Em coluna única os cinco símbolos ficam empilhados e
	// comparáveis; na mesma linha a legenda passa de 100 colunas e quebra sozinha
	// num terminal estreito — o oposto do que ela serve.
	fmt.Println()
	for _, l := range partes {
		fmt.Printf("  %s\n", l)
	}
}

// warnGatesWithoutTarget relata os gates DECLARADOS que não apareceram na tabela.
//
// A tabela é montada de `profile.GateNames()`, que só conhece gate AVALIADO. Um gate
// declarado cujo `on:` não casa nenhum nó do mapa nunca é avaliado, então desaparece por
// completo — e o rodapé ainda diz "✓ pode promover, todos os gates passaram".
//
// Para quem está adotando o framework, isso lê como "está tudo certo". Medido ao pôr o
// Anchors no próprio Anchors: declarei 8 gates, o cabeçalho anunciou "8 gates", a tabela
// mostrou 2 e o rodapé deu ✓. Os 6 ausentes eram os de trinca, e o projeto tem 0 spec e 0
// feature — a informação útil ("declarei regra para artefato que não existe aqui") era
// exatamente a que não aparecia.
//
// Não é falha: é declaração sem alvo, e o conserto é do projeto (criar o artefato ou
// remover o gate). Por isso avisa, não barra.
func warnGatesWithoutTarget(declarados []config.Gate, p gate.Profile) {
	avaliados := map[string]bool{}
	for _, n := range p.GateNames() {
		avaliados[n] = true
	}
	var semAlvo []string
	for _, gt := range declarados {
		if gt.Name != "" && !avaliados[gt.Name] {
			semAlvo = append(semAlvo, gt.Name)
		}
	}
	if len(semAlvo) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(i18n.T("check.unused_gates_header", len(semAlvo)))
	for _, n := range semAlvo {
		fmt.Printf("    %s\n", n)
	}
	fmt.Println(i18n.T("check.unused_gates_note"))
}

// printTiming mostra ONDE a varredura gastou o tempo — por gate, e depois os alvos
// individuais mais caros.
//
// O custo do `check --all` e' invisivel no relatorio normal: ele leva minutos num
// projeto real, e nada diz onde. Sem a medida, otimizar e' adivinhar qual gate cobrar —
// e a suspeita natural (o gate com mais alvos) costuma estar errada, porque um `run:`
// externo que sobe um compilador custa mais em UMA execucao do que um checker interno
// em centenas.
//
// Os ALVOS entram porque o gate caro nem sempre e' caro por igual: um unico arquivo
// patologico (uma spec gigante, um teste que o `run:` recompila) aparece aqui com nome
// e endereco, e e' acionavel de um jeito que a media por gate nao e'.
func printTiming(p gate.Profile) {
	type linha struct {
		nome  string
		total time.Duration
		pior  time.Duration
		alvos int
	}
	var linhas []linha
	var total time.Duration
	for nome, s := range p.ByGate {
		n := s.Pass + s.Fail + s.Skip + s.Pending + s.Judge
		linhas = append(linhas, linha{nome, s.Duracao, s.Pior, n})
		total += s.Duracao
	}
	if len(linhas) == 0 {
		return
	}
	// Empate desempatado pelo NOME, e nao deixado ao acaso do mapa: duas execucoes da
	// mesma varredura tem de imprimir a mesma ordem, senao comparar dois relatorios
	// vira ruido.
	sort.Slice(linhas, func(i, j int) bool {
		if linhas[i].total != linhas[j].total {
			return linhas[i].total > linhas[j].total
		}
		return linhas[i].nome < linhas[j].nome
	})

	largura := 0
	for _, l := range linhas {
		if len(l.nome) > largura {
			largura = len(l.nome)
		}
	}

	fmt.Println()
	fmt.Println(i18n.T("check.timing_header", arredonda(total)))
	for _, l := range linhas {
		fmt.Println(i18n.T("check.timing_item", largura, l.nome, arredonda(l.total), l.alvos, arredonda(l.pior)))
	}

	// Os alvos individuais mais caros, atravessando todos os gates.
	piores := append([]gate.Result(nil), p.Results...)
	sort.Slice(piores, func(i, j int) bool {
		if piores[i].Duracao != piores[j].Duracao {
			return piores[i].Duracao > piores[j].Duracao
		}
		if piores[i].Gate != piores[j].Gate {
			return piores[i].Gate < piores[j].Gate
		}
		return piores[i].Target < piores[j].Target
	})
	if len(piores) > 0 && piores[0].Duracao > 0 {
		fmt.Println()
		fmt.Println(i18n.T("check.timing_targets"))
		for i, r := range piores {
			// Zero nao e' informacao: um alvo que nem registrou tempo nao ajuda ninguem
			// a otimizar, e listar dez deles empurraria os reais para fora da tela.
			if i == 10 || r.Duracao == 0 {
				break
			}
			fmt.Println(i18n.T("check.timing_target_item", arredonda(r.Duracao), r.Gate, r.Target))
		}
	}
	fmt.Println(i18n.T("check.timing_note"))
}

// arredonda corta a precisao ao que o leitor usa para DECIDIR. `1.234567ms` e
// `1.23ms` levam a mesma acao, e o digito a mais so' atrapalha a comparacao entre
// linhas — que e' para o que esta tabela existe.
func arredonda(d time.Duration) string {
	switch {
	case d >= time.Second:
		return d.Round(10 * time.Millisecond).String()
	case d >= time.Millisecond:
		return d.Round(100 * time.Microsecond).String()
	}
	return d.Round(time.Microsecond).String()
}

func printProfile(p gate.Profile, onlyIssues, showDrift bool) {
	nomes := p.GateNames()
	wn := nameWidth(nomes)
	w := computeWidths(p)
	limpos := 0

	// perfil por gate
	for _, name := range nomes {
		s := p.ByGate[name]
		tag := i18n.T("check.tag.informative")
		if s.Blocking {
			tag = i18n.T("check.tag.blocking")
		}
		if s.Judge > 0 {
			tag = i18n.T("check.tag.judgment")
			fmt.Printf("  %-*s  %-11s  "+i18n.T("check.judge_pending_table")+"\n", wn, name, tag, w.judge, s.Judge)
			continue
		}
		// O `~` fundia TRÊS coisas: "não se aplica", "não havia o que confrontar" e
		// DRIFT REAL rebaixado a aviso (ex.: o código do cenário casa mas a descrição
		// do teste divergiu). Somar tudo num número faz o drift parecer benigno — e é
		// o oposto: é a única categoria acionável do balde. Separamos em ⚠.
		drift := driftCount(p, name)
		if onlyIssues && cleanGate(s, drift) {
			limpos++
			continue
		}
		// A coluna do ⚠ é reservada SEMPRE, mesmo vazia. Omiti-la nas linhas sem
		// drift empurrava o `~` para a esquerda só nelas, e a coluna passava a
		// existir em dois lugares na mesma tabela — que é o pior caso: o olho
		// desce a lista comparando números que não estão na mesma vertical.
		fmt.Printf("  %-*s  %-11s  ✓%*d  ✗%*d%s%s  ~%*d\n", wn, name, tag,
			w.pass, s.Pass, w.fail, s.Fail,
			driftSeparator(w.drift), driftColumn(drift, w.drift),
			w.skip, s.Skip+s.Pending-drift)
	}
	// O gate omitido continua tendo rodado, e o número diz isso. Sem esta linha o
	// `--only-issues` pareceria uma varredura menor, e não a mesma varredura com a
	// listagem enxuta.
	if limpos > 0 {
		fmt.Println()
		fmt.Println(i18n.T("check.only_issues_omitted", limpos))
	}
	printLegenda(p, w)

	// DRIFT detalhado: só sob `--show-drift`, e aí SEM TETO.
	//
	// Truncar a lista devolvia o problema que ela existe para resolver: quem pede
	// os endereços quer agir sobre eles, e "… e mais 2405" deixa 2405 sem
	// endereço. Ou a lista é completa, ou o contador da tabela já bastava.
	//
	// Por isso o padrão é o contador: numa varredura grande são milhares de
	// linhas, e despejá-las sem ninguém ter pedido enterraria as issues logo
	// abaixo. A tabela continua anunciando `⚠407`, e a flag é o caminho do número
	// para os endereços.
	if showDrift {
		printDrift(driftResults(p))
	}

	// Motivo dos `~` (skip/pending). O contador sozinho não diz nada — quem lê fica em
	// dúvida se `~1` é problema dele. No modo incremental (poucos nós) o motivo cabe na
	// tela e é acionável ("o código ainda não existe"), então mostramos.
	if reasons := skipReasons(p); len(reasons) > 0 && len(p.Results) <= maxResultsForSkipDetail {
		fmt.Println()
		fmt.Println(i18n.T("check.undetermined_summary", len(reasons)))
		for _, r := range reasons {
			fmt.Printf("  ~ %s @ %s\n", r.Gate, r.Target)
			if r.Detail != "" {
				fmt.Print(indent(r.Detail, "      "))
			}
		}
	}

	if len(p.Failures) > 0 {
		fmt.Println()
		fmt.Println(i18n.T("check.issues_summary", len(p.Failures)))
		for _, r := range p.Failures {
			mark := i18n.T("check.tag.informative")
			if r.Blocking {
				mark = i18n.T("check.tag.blocks_mark")
			}
			fmt.Printf("  ✗ [%s] %s @ %s\n", mark, r.Gate, r.Target)
			if r.Detail != "" {
				fmt.Print(indent(r.Detail, "      "))
			}
		}
	}

	fmt.Println()
	// O ACHADO INFORMATIVO precisa sobreviver ao veredito.
	//
	// A última linha é a que se lê, e "✓ pode promover" apagava tudo que veio antes.
	// Medido em três rodadas de um E2E: um gate informativo diagnosticou um defeito real
	// (uma função que truncava a consulta e quebrava a regra que a consome), o autor leu
	// o verde, e registrou por escrito que o arquivo estava correto — contra o gate que
	// ele mesmo citou pelo nome. O sinal existia e sumiu no meio.
	//
	// Informativo não barra, e não deve barrar: é a maturação (QUALITY §7). Mas "não
	// barra" é diferente de "não aconteceu", e o veredito precisa dizer as duas coisas.
	informativos := 0
	// DRIFT é `Pending` COM motivo: o gate olhou, algo divergiu, e ele não barra. Contá-lo
	// aqui é o que evita a linha "sem achado em aberto" aparecer sobre uma spec que declara
	// 3 decisões que ninguém tomou — o veredito diria o contrário do que o gate viu.
	drift := 0
	// naoConfrontado são os Skip PUROS — o gate rodou e o nó não lhe dizia nada. Ficavam
	// fora de toda contagem, então um projeto onde NADA foi medido lia "sem achado em
	// aberto" no rodapé. Contá-los à parte é o que separa "conforme" de "não medido".
	naoConfrontado := 0
	for _, r := range p.Results {
		if r.Verdict == gate.Fail && !r.Blocking {
			informativos++
		}
		if r.Verdict == gate.Pending && r.Detail != "" {
			drift++
		}
		if r.Verdict == gate.Skip {
			naoConfrontado++
		}
	}
	if p.Passed {
		switch {
		case informativos > 0:
			fmt.Println(i18n.T("check.promote_with_info", informativos))
			if drift > 0 {
				fmt.Println(i18n.T("check.promote_drift_note", drift))
			}
			fmt.Println(i18n.T("check.promote_info_note"))
		case drift > 0:
			fmt.Println(i18n.T("check.promote_with_drift", drift))
			fmt.Println(i18n.T("check.pending_note"))
		case naoConfrontado > 0:
			// `Skip` PURO: o gate rodou, o nó não se aplicava, e nada foi medido. Não é
			// achado nem pendência — é COBERTURA AUSENTE, e o veredito não pode chamar
			// isso de "sem achado em aberto".
			//
			// Medido ao adotar o Anchors no próprio Anchors: 222 nós, dois gates com ~67 e
			// ~85, e o rodapé dizia "todos os gates passaram, sem achado em aberto". 152
			// confrontos não aconteceram, e a última linha — a que se lê — afirmava o
			// oposto. Para quem está adotando, é a diferença entre "meu projeto está
			// conforme" e "eu ainda não declarei o que fazia isto ser medido".
			fmt.Println(i18n.T("check.promote_with_unmatched", naoConfrontado))
			fmt.Println(i18n.T("check.promote_unmatched_note"))
		default:
			fmt.Println(i18n.T("check.can_promote_all_clean"))
		}
	} else {
		msg := i18n.T("check.blocked_by_gates", len(p.Blocked))
		if informativos > 0 {
			msg += i18n.T("check.blocked_extra_info", informativos)
		}
		fmt.Println(msg)
	}

	rememberMaturation(p, onlyIssues)
}

// rememberMaturation avisa sobre gate informativo que já está LIMPO.
//
// A maturação (QUALITY §7) tem uma metade que o Anchors não cobrava: um gate nasce
// informativo porque o projeto ainda não cumpre o limiar, e quando passa a cumprir,
// ninguém volta ao anchors.yaml para promovê-lo. O gate fica medindo sem defender — e o
// projeto acha que está protegido por algo que não barra nada.
//
// A promoção continua sendo decisão humana; o que muda é que ela deixa de depender de
// alguém lembrar sozinho. E o lembrete aparece AQUI, onde a pessoa acabou de ler os
// vereditos — um aviso que exige rodar outro comando é um aviso que ninguém vê.
func rememberMaturation(p gate.Profile, onlyIssues bool) {
	prom := gate.PromotableGates(p)
	if len(prom) == 0 {
		return
	}
	// Com `--only-issues` a pessoa pediu para NÃO ver o que está limpo. Nomear os gates
	// aqui contradiria a flag — mas omiti-los por completo esconderia a maturação
	// pendente, que é justamente o que ela precisa saber. Então: conta, não lista.
	if onlyIssues {
		fmt.Println()
		fmt.Println(i18n.T("check.maturation_reminder", len(prom)))
		return
	}
	fmt.Println()
	fmt.Println(i18n.T("check.maturation_promotable", len(prom)))
	for i, g := range prom {
		if i == 6 {
			fmt.Println(i18n.T("check.maturation_more", len(prom)-6))
			break
		}
		fmt.Println(i18n.T("check.maturation_item", g.Gate, g.Passou))
	}
	fmt.Println(i18n.T("check.maturation_note"))
}

// indent prefixa cada linha de s com pad — para o detalhe do gate (que pode ensinar
// o formato esperado em várias linhas) sair legível sob a issue.
func indent(s, pad string) string {
	var b strings.Builder
	for _, line := range strings.Split(strings.TrimRight(breakOccurrences(s), "\n"), "\n") {
		b.WriteString(pad)
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// breakOccurrences põe cada ocorrência do detalhe na sua própria linha.
//
// Dez gates montam a mensagem com `strings.Join(achados, "; ")`, e cinco
// violações no mesmo arquivo saíam numa linha de 800 caracteres: o leitor não
// distingue onde uma acaba e a outra começa, e a repetição do texto da regra
// (idêntico nas cinco) afoga o único dado que varia, que é o número da linha.
//
// A quebra é no separador que os gates JÁ usam, então nenhum deles precisa
// mudar. Basta UM separador para quebrar: duas ocorrências já se confundem numa
// linha só, e é o caso mais comum.
func breakOccurrences(s string) string {
	s = strings.ReplaceAll(s, "; ", ";\n")

	// Listas longas separadas por vírgula também quebram — 17 gates montam a
	// mensagem com `strings.Join(lista, ", ")` SEM cortar, e uma lista grande vira
	// uma linha que ninguém lê nem grepa.
	//
	// O que separa LISTA de PROSA não é o tamanho da linha: é o que vem entre as
	// vírgulas. Item de lista é um token — um caminho, um código, um símbolo — e
	// não tem espaço interno; oração tem. A primeira versão usava "linha longa com
	// 3+ vírgulas" e picou a própria mensagem deste gate ("a spec existe,\n o
	// código existe,\n os dois se referenciam") — a frase que o comentário citava
	// como o caso a NÃO quebrar.
	var out []string
	for _, linha := range strings.Split(s, "\n") {
		if isTokenList(linha) {
			linha = strings.ReplaceAll(linha, ", ", ",\n")
		}
		out = append(out, linha)
	}
	return strings.Join(out, "\n")
}

// listBreakThreshold — abaixo disto a linha cabe na tela e não vale quebrar, seja
// ela lista ou não. Escolhido para caber num terminal de 120 colunas com a
// indentação do detalhe (6 a 8 espaços).
const listBreakThreshold = 110

// isTokenList: a linha é uma enumeração de itens sem espaço interno?
//
// `a.spec.md, b.spec.md, c.spec.md` é lista; `a spec existe, o código existe, e
// os dois se referenciam` é prosa. A distinção é o espaço DENTRO do item: nome de
// arquivo, código de regra e símbolo não têm; oração tem. Exige maioria dos itens
// sem espaço para tolerar o último ("… e mais 3") e o texto que abre a lista.
func isTokenList(linha string) bool {
	if len([]rune(linha)) <= listBreakThreshold {
		return false
	}
	partes := strings.Split(linha, ", ")
	if len(partes) < 4 {
		return false
	}
	tokens := 0
	for _, p := range partes {
		if !strings.Contains(strings.TrimSpace(p), " ") {
			tokens++
		}
	}
	return tokens*2 > len(partes)
}

// maxResultsForSkipDetail — acima disso (varredura --all) a lista de motivos vira
// ruído; o contador basta. No incremental o motivo é o que o usuário precisa.
const maxResultsForSkipDetail = 40

// driftResults devolve os DRIFTs: `Pending` COM motivo — o gate olhou, algo
// divergiu, e ele não barrou.
//
// Eles saíam do mesmo balde que os `~` e por isso desapareciam no `--all`, junto
// com o corte de ruído: a tabela anunciava `⚠407` e o detalhe abaixo não trazia
// nenhum. Contar uma categoria como acionável e depois escondê-la é pior que não
// a separar — o número vira uma acusação sem endereço.
func driftResults(p gate.Profile) []gate.Result {
	var out []gate.Result
	for _, r := range p.Results {
		if r.Verdict == gate.Pending && r.Detail != "" {
			out = append(out, r)
		}
	}
	return out
}

// skipReasons devolve os resultados indeterminados (skip/pending) QUE TÊM motivo
// escrito. Skip silencioso (gate simplesmente não se aplica àquele kind) não vira
// linha — só polui.
func skipReasons(p gate.Profile) []gate.Result {
	var out []gate.Result
	for _, r := range p.Results {
		// `Pending` COM motivo é DRIFT, e tem bloco próprio (`driftResults`).
		// Deixá-lo aqui o faria herdar o corte de ruído do `~` e sumir no --all.
		if r.Verdict == gate.Skip && r.Detail != "" {
			out = append(out, r)
		}
	}
	return out
}

// warnIfMapStale avisa quando o MAPA é mais velho que os arquivos que ele descreve.
// Sem isso, um arquivo criado/alterado depois do último `map build` é INVISÍVEL ao
// check — e o verde passa a não significar nada (o gate não roda sobre o que não está
// no mapa). É a diferença entre "passou" e "não foi olhado".
//
// A conferência é sobre CONTEÚDO — quais arquivos governados o mapa não conhece —, e
// não sobre mtime.
//
// A primeira versão comparava o mtime do `anchors.graph.yaml` com o dos arquivos-fonte.
// Barato, e errado em CI: o `checkout` escreve TODO o repositório no instante do clone,
// então cada arquivo fica microssegundos mais novo que o mapa e o aviso disparava em toda
// execução. Medido no primeiro run do pipeline de gates — 26 arquivos "mudaram" num
// repositório onde nada mudara, e o log do mesmo job dizia, duas linhas acima, que o mapa
// correspondia ao repositório.
//
// Ruído assim custa mais do que parece: um aviso que aparece em todo PR verde ensina a
// ignorá-lo, e aí ele não serve quando for verdadeiro.
//
// Arquivo governado que não está no mapa é INVISÍVEL ao check, e é isso que o aviso
// existe para dizer — a diferença entre "passou" e "não foi olhado".
func warnIfMapStale(root, mapPath string, cfg *config.Config, g *mapx.Graph) {
	if g == nil {
		return
	}
	files, err := scan.Walk(root, cfg)
	if err != nil {
		return
	}
	noMapa := make(map[string]bool, len(g.Nodes))
	for _, n := range g.Nodes {
		noMapa[filepath.ToSlash(n.ID)] = true
	}
	newer := 0
	var sample string
	for _, f := range files {
		if noMapa[filepath.ToSlash(f.Path)] {
			continue
		}
		newer++
		if sample == "" {
			sample = f.Path
		}
	}
	if newer == 0 {
		return
	}
	extra := ""
	if newer > 1 {
		extra = fmt.Sprintf(" (and %d more)", newer-1)
	}
	// A mensagem diz o que foi MEDIDO: "não está no mapa", e não "mudou". Dizer "mudou"
	// mandaria quem lê procurar uma alteração que pode não existir — o caso comum é
	// arquivo NOVO, que nunca esteve lá.
	fmt.Printf("⚠ STALE map: %d governed file(s) are not in the map — e.g.: %s%s\n"+
		"  The gates only see what is in the map; run `anchors map build` before trusting this result.\n\n",
		newer, sample, extra)
}

// driftCount conta, dentro do balde `~` de um gate, os resultados que são DRIFT REAL
// (indeterminados COM motivo escrito) — em oposição aos que apenas não se aplicam.
// Um Pending com detalhe é o gate dizendo "olhei e algo divergiu, mas não vou barrar";
// somá-lo aos "não se aplica" esconde a única parte acionável do número.
func driftCount(p gate.Profile, gateName string) int {
	n := 0
	for _, r := range p.Results {
		if r.Gate != gateName {
			continue
		}
		// DRIFT é PENDING com motivo: o gate OLHOU e algo divergiu, mas não barra.
		// Skip com motivo é o oposto — o gate não se aplica ("não é uma tela"), e
		// contá-lo como drift inflaria o número com ruído benigno.
		if r.Verdict == gate.Pending && r.Detail != "" {
			n++
		}
	}
	return n
}

// unitPieces devolve os caminhos das outras peças da trinca de um alvo — spec,
// feature, teste e código. É a vizinhança que a IDENTIDADE define, e não a que as arestas
// registram: as duas costumam coincidir, mas a segunda depende de o mapa já ter ligado as
// pontas, e o `check` precisa funcionar antes disso.
func unitPieces(target string) []string {
	base := strings.TrimSuffix(target, filepath.Ext(target))
	for _, suf := range []string{".spec.md", ".feature", ".test", ".spec"} {
		base = strings.TrimSuffix(base, suf)
	}
	var out []string
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx", ".ts", ".tsx"} {
		if cand := base + suf; cand != target {
			out = append(out, cand)
		}
	}
	return out
}

// gateCmdID devolve o ID de um gate (ou o nome, quando o ID falta). Duplica a lógica
// de `internal/gate` de propósito: exportá-la só para isto acoplaria os dois pacotes por
// uma função de três linhas.
func gateCmdID(g config.Gate) string {
	if g.ID != "" {
		return g.ID
	}
	return g.Name
}

// normalizeChanged põe os caminhos na mesma forma que os IDs do mapa (relativos à raiz),
// para que a comparação não dependa de como o chamador escreveu o caminho — o pre-commit
// passa relativo, e quem roda à mão costuma passar absoluto.
func normalizeChanged(changed []string, root string) []string {
	out := make([]string, 0, len(changed))
	for _, c := range changed {
		p := c
		if filepath.IsAbs(p) {
			if rel, err := filepath.Rel(root, p); err == nil {
				p = rel
			}
		}
		out = append(out, filepath.ToSlash(filepath.Clean(p)))
	}
	return out
}

// queuedJudgments conta os julgamentos que AGUARDAM resposta.
func queuedJudgments(root string) int {
	tasks, err := queue.List(root)
	if err != nil {
		return 0
	}
	n := 0
	for _, t := range tasks {
		if t.SuggestedNext == "judge" {
			n++
		}
	}
	return n
}

// staleBinaryWarning compara quem gravou o mapa com quem está rodando.
//
// Compara por IGUALDADE, e não por ordem: "dev" não é ordenável contra "0.1.9", e os dois
// builds locais que produziram o defeito se chamavam "dev" — ordenar não teria pego
// nenhum deles. Diferente já é o suficiente para avisar; QUAL é mais novo, quem lê decide.
func staleBinaryWarning(gravouMapa, rodando string) string {
	if gravouMapa == "" || gravouMapa == rodando {
		return ""
	}
	return i18n.T("map.version_mismatch", gravouMapa, rodando)
}
