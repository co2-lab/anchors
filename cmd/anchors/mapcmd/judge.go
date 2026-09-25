package mapcmd

import (
	"fmt"
	"github.com/co2-lab/anchors/internal/i18n"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/queue"
	"github.com/co2-lab/anchors/internal/suggestion"
	"github.com/spf13/cobra"
)

// `anchors judge` é o VERBO do medidor de julgamento por IA (QUALITY §5.2). O CLI
// não julga — a IA que opera o Anchors lê o guide, confronta o alvo, e reporta o
// veredito aqui. O CLI faz a MESMA contabilidade de um gate determinístico: carimba
// a aresta no mapa e abre (ou resolve) a issue. Assim o julgamento subjetivo entra
// no mesmo loop anti-drift — o carimbo leva a rev do alvo, então o veredito envelhece
// (fica stale) se o alvo mudar depois.
func newJudgeCmd() *cobra.Command {
	var root, mapPath, gateName, verdict, reason, patchFile string
	var pending, recordIssues bool
	cmd := &cobra.Command{
		Use:   "judge <target>",
		Short: "Record an AI's verdict for a judgment gate",
		Long: `Records the verdict of an AI JUDGMENT gate over a target — the same
double output as a deterministic gate (stamp on the map + issue).

  anchors judge --pending                      lists the targets awaiting judgment
  anchors judge <alvo> --gate <g> --verdict pass|fail|waived --reason "..."

The flow: 'anchors check' marks the targets of a 'measures: judgment' gate as
pending and queues them. The AI (worker) reads the gate's guide, confronts the target, and
reports here.

IMPORTANT — do not waste what you have already read. To judge, you read the whole guide
and the whole target; so --reason is NOT a verdict sentence, it is the FULL REPORT. The
issue body IS your --reason, verbatim (multi-line markdown accepted). On a 'fail',
list EVERY non-conformity with: what is wrong, WHERE (file:line), WHY
(which rule of the guide), and HOW to fix it. That way nobody has to reprocess the target
later just to find out what to fix. E.g.:

  anchors judge apps/.../AlertsScreen.tsx --gate atomic-design --verdict fail \
    --reason "$(cat <<'EOF'
  ## Report — atomic-design @ AlertsScreen
  ### 1. Inline list item should be a molecule  (l.151-166)
  - **what:** the item's JSX (View + icon + texts) is assembled inline in the .map
  - **rule:** SCREEN_GUIDE §Pages — screens inject components, they do not assemble blocks
  - **how:** extract to components/molecules/AlertRow.tsx and use <AlertRow item=.../>
  ### 2. ...
  EOF
  )"

'fail' opens the issue with that report; 'pass' resolves the issue if there was one (and the
--reason, if any, becomes a recorded observation).

'waived' is for when THE TARGET OF THE QUESTION DOES NOT EXIST. A spec that declares
'@TBD: code' states that the code has not been written yet — the normal flow, since the spec
is the anchor and is born first. Then the question "does the excerpt REALIZE what the rule
describes?" has no excerpt, and the other two verdicts lie: 'pass' states that the code realizes
the rule (and the stamp stays on the map looking like a real verification), 'fail' fails work nobody
got wrong. The --reason is mandatory and names the absence: which piece is missing, and where
it is declared.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}

			// --pending: lista os alvos que aguardam julgamento (tasks judge na fila)
			if pending {
				return listPendingJudgments(absRoot)
			}

			if len(args) != 1 {
				return fmt.Errorf("provide the target (e.g.: `anchors judge Foo.tsx --gate atomic --verdict fail --reason ...`) or use --pending")
			}
			if gateName == "" {
				return fmt.Errorf("--gate is mandatory (the judgment gate you are evaluating)")
			}
			v := strings.ToLower(verdict)
			if err := ValidateVerdict(v, reason); err != nil {
				return err
			}

			target := relTo(absRoot, args[0])
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}
			// O review acontece sobre trabalho NOVO, e no início da unidade só a spec
			// existe: o `anchors work review --for <alvo>.ts` prescreve o caminho do
			// código, que ainda não foi escrito. Recusar aí deixa o revisor sem onde
			// registrar o achado — medido num E2E, foi o que aconteceu.
			//
			// Cair para a peça que EXISTE preserva a identidade (spec e código são a mesma
			// unidade) e mantém o comando prescrito funcionando desde a primeira etapa.
			if !nodeExists(g, target) {
				if alt := unitExistingPiece(g, target); alt != "" {
					fmt.Printf("   (the target does not exist yet; recording on `%s`, the piece of this unit that is already in the map)\n", alt)
					target = alt
				}
			}
			if !nodeExists(g, target) {
				return fmt.Errorf("target %q is not in the map", target)
			}
			// valida que o gate existe e é de julgamento
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
			gc, ok := findJudgmentGate(cfg, gateName)
			// `review` é o julgamento do CICLO, não um gate do projeto: ele não roda sobre
			// todos os nós (declará-lo assim marcava 828 alvos como pendentes — uma fila
			// que ninguém puxa), e sim sobre o que acabou de ser entregue, quando o
			// watcher o enfileira depois do `deliver`.
			//
			// Aceitá-lo aqui é o que fecha o ciclo `deliver → review → issue`. Sem isto o
			// achado do revisor não tinha onde ser registrado: ficava no relatório da
			// sessão e morria com ela — inclusive o ACHADO CRUZADO, o defeito que um
			// revisor encontra numa unidade que não é a dele.
			if !ok && gateName == "review" {
				gc, ok = config.Gate{Name: "review", Blocking: config.Bool(false)}, true
			}
			if !ok {
				return fmt.Errorf("gate %q does not exist or is not a judgment gate (measures: judgment)", gateName)
			}

			now := time.Now()
			failed := v == "fail"
			// O carimbo distingue os três: `ok` afirma que foi verificado, `issue` que
			// reprovou, e `dispensado` que NÃO HAVIA o que verificar. Gravar `ok` no
			// terceiro caso perderia justamente a informação que importa — quem lesse o
			// mapa depois veria uma verificação que não aconteceu.
			verdictStr := "ok"
			switch v {
			case "fail":
				verdictStr = "issue"
			case "waived":
				verdictStr = "waived"
			}
			// 1) CARIMBO — marca a aresta guide→alvo (o confronto da régua contra o
			//    alvo) com o veredito da IA. Se o gate declara guide, carimba essa
			//    aresta; senão, cai para o carimbo por-nó (todas as arestas do alvo).
			stamped := 0
			if gc.Guide != "" && g.StampEdge(gc.Guide, target, verdictStr, now.Format(time.DateOnly)) {
				stamped = 1
			} else {
				// Carimbo por NÓ: o `judge` julga uma UNIDADE, e o `StampEdges` exige as
				// duas pontas confrontadas na mesma rodada — com um alvo só, nenhuma
				// aresta qualifica e o carimbo saía sempre zero. O veredito abria a issue
				// e não tocava o grafo: um `check` posterior não via o achado.
				// leva o NOME do gate: é ele que permite ao `check` posterior saber que
				// este julgamento foi respondido, em vez de perguntar de novo.
				stamped = g.StampNodeByGate(target, verdictStr, now.Format(time.DateOnly), gateName)
			}
			if err := mapx.Save(g, mapPath); err != nil {
				return fmt.Errorf("save stamp: %w", err)
			}

			// 2) ISSUE — fail abre violation (com o reason da IA); pass resolve.
			iss := issue.Issue{
				Kind: issue.Violation, Target: target, Gate: gc.Name,
				Detail: reason, Date: now.Format("2006-01-02"),
			}
			// 2b) SUGESTÃO — quem julgou leu o alvo inteiro e sabe qual seria a
			// correção. Sem isto ela só cabia em prosa dentro do laudo, e quem lesse
			// teria de reconstruí-la à mão. Com `--patch`, vira um diff aplicável.
			if failed && patchFile != "" {
				diff, rerr := os.ReadFile(patchFile)
				if rerr != nil {
					return fmt.Errorf("read the patch: %w", rerr)
				}
				sg := suggestion.Suggestion{
					ID:   "judge-" + gc.Name + "-" + strings.ReplaceAll(common.RelSlug(target), "/", "-"),
					Gate: gc.Name, Target: target, Origin: suggestion.FromJudgment,
					Why: reason, Patch: string(diff),
				}
				criada, at, serr := suggestion.Open(absRoot, sg)
				if serr != nil {
					return serr
				}
				if criada {
					fmt.Printf("  ↳ fix suggestion at %s (`anchors suggest show %s`)\n", at, sg.ID)
				}
			}

			// MANUAL MODE writes no issue unless asked: the verdict is in the map (above),
			// and the report is printed here, the one place it stays.
			if !judgeWritesIssue(cfg, recordIssues) {
				switch {
				case failed:
					fmt.Printf("✗ judged FAIL — stamped in the map; no issue written (mode: manual — --record-issues writes it)\n")
					if strings.TrimSpace(reason) != "" {
						fmt.Printf("\n%s\n\n", strings.TrimSpace(reason))
					}
				case v == "waived":
					fmt.Println(i18n.T("judge.waived"))
				default:
					fmt.Println(i18n.T("judge.pass"))
				}
			} else if failed {
				created, at, err := issue.Open(absRoot, iss)
				if err != nil {
					return err
				}
				if created {
					fmt.Printf("✗ judged FAIL — issue opened at %s/todo/\n", issue.Dir)
				} else if reaberta, rerr := issue.Reopen(absRoot, iss); rerr != nil {
					return rerr
				} else if reaberta {
					// Achado NOVO sobre unidade já revisada: reabre e acrescenta o laudo.
					// Antes o `--reason` era descartado em silêncio, e o defeito não
					// ficava em lugar nenhum.
					fmt.Printf("✗ judged FAIL — NEW finding added; issue reopened at %s/todo/\n", issue.Dir)
				} else {
					fmt.Printf("✗ judged FAIL — same finding already recorded (%s/), nothing to do\n", at)
				}
			} else if v == "waived" {
				// A palavra IMPORTA aqui. Anunciar "PASS" desfaria o ponto inteiro do
				// terceiro veredito: quem lê a saída ficaria com a impressão de que o
				// alvo foi verificado e aprovado — a mesma confusão que o `pass`
				// mentiroso produzia, agora vinda do próprio comando.
				if ok, _ := issue.Resolve(absRoot, iss.Key()); ok {
					fmt.Println(i18n.T("judge.waived_resolved", issue.Dir))
				} else {
					fmt.Println(i18n.T("judge.waived"))
				}
			} else {
				if ok, _ := issue.Resolve(absRoot, iss.Key()); ok {
					fmt.Println(i18n.T("judge.pass_resolved", issue.Dir))
				} else {
					fmt.Println(i18n.T("judge.pass"))
				}
			}
			fmt.Printf("  stamped: %d edge(s) with the AI verdict (gate '%s')\n", stamped, gc.Name)

			// 3) fecha a task judge correspondente, se houver
			closeJudgeTask(absRoot, gc.Name, target)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().StringVar(&gateName, "gate", "", "the judgment gate being evaluated")
	cmd.Flags().StringVar(&verdict, "verdict", "", "pass | fail")
	cmd.Flags().StringVar(&reason, "reason", "", "the full REPORT (multi-line markdown) — becomes the issue body; mandatory on fail")
	cmd.Flags().StringVar(&patchFile, "patch", "", "file with the diff that FIXES the finding — opens an applicable suggestion (`anchors suggest`)")
	cmd.Flags().BoolVar(&pending, "pending", false, "lists the targets awaiting judgment")
	cmd.Flags().BoolVar(&recordIssues, "record-issues", false, "mode: manual — write the issue (on fail) and close it (on pass); by default none is written")
	return cmd
}

// findJudgmentGate acha o gate pelo nome, aceitando também o nome ANTIGO.
//
// O `Load` canoniza os nomes ao ler o arquivo: um `anchors.yaml` que declara
// `mock-detect-cobre-o-dialeto` chega aqui com `g.Name == "mock-detect-covers-dialect"`.
// Comparar só contra `g.Name` recusava justamente o nome que a pessoa tem na tela —
// medido no blue-eyes: `--gate mock-detect-cobre-o-dialeto` respondia "gate não existe"
// com o gate declarado, visível, três linhas acima no próprio arquivo.
//
// A mensagem era pior que o erro: ela manda procurar um gate que está ali, e não diz que
// o nome mudou. Canonizar o ARGUMENTO fecha o buraco pela mesma tabela que o `Load` usa,
// então os dois nomes nunca divergem.
func findJudgmentGate(cfg *config.Config, name string) (config.Gate, bool) {
	canonico := name
	for _, g := range cfg.Gates {
		if (g.Name == name || g.Name == canonico) && g.IsJudgment() {
			return g, true
		}
	}
	return config.Gate{}, false
}

func listPendingJudgments(root string) error {
	tasks, err := queue.List(root)
	if err != nil {
		return err
	}
	n := 0
	for _, t := range tasks {
		if t.SuggestedNext == "judge" {
			n++
			fmt.Printf("○ %s\n    %s\n", t.Changed, t.Reason)
		}
	}
	if n == 0 {
		fmt.Println("no target awaiting judgment (run `anchors check` to find out)")
	} else {
		fmt.Printf("\n%d target(s) — judge with: anchors judge <target> --gate <g> --verdict pass|fail|waived --reason ...\n", n)
	}
	return nil
}

// closeJudgeTask fecha a task judge daquele gate+alvo, se estiver na fila.
func closeJudgeTask(root, gateName, target string) {
	id := "judge-" + gateName + "-" + strings.ReplaceAll(common.RelSlug(target), "/", "-")
	_ = queue.MarkDone(root, id)
}

// unitExistingPiece acha, para um alvo ausente do mapa, outra peça da MESMA unidade
// que já esteja lá — a spec, tipicamente, quando o código ainda não nasceu.
func unitExistingPiece(g *mapx.Graph, target string) string {
	base := strings.TrimSuffix(target, filepath.Ext(target))
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx", ".ts", ".tsx"} {
		if cand := base + suf; cand != target && common.NodeExists(g, cand) {
			return cand
		}
	}
	return ""
}

// ValidateVerdict confronta o veredito recebido e o motivo que o acompanha.
//
// Extraída do `RunE` para ser TESTÁVEL: o contrato dos três vereditos é o que impede o
// `pass` mentiroso, e um contrato sem teste é uma intenção.
//
// `dispensado` é o terceiro desfecho, e existe porque os outros dois MENTEM quando o alvo
// da pergunta não existe.
//
// Uma spec que declara `@TBD: code` afirma que o código ainda não foi escrito — o fluxo
// normal, já que a spec é a âncora e nasce primeiro. O gate `regra-cumprida` pergunta "o
// trecho REALIZA o que a regra descreve?", e não há trecho. Com só `pass`/`fail`, quem
// julga escolhe entre afirmar que o código realiza a regra (e o carimbo fica no mapa
// parecendo verificação) ou reprovar trabalho que ninguém errou.
//
// Medido no blue-eyes (#76): a saída usada foi `pass`, com o motivo explicando que não
// havia o que medir. Funcionou uma vez e ensina o hábito errado — carimbar julgamento sem
// olhar é o que corrói o valor de `measures: judgment`.
func ValidateVerdict(v, reason string) error {
	// O VALOR ANTIGO ainda é aceito, e vira o canônico.
	//
	// `dispensado` já está em mapas commitados (`verdict: dispensado`) e em scripts.
	// Recusá-lo faria o `check` reperguntar julgamentos que alguém já respondeu — e o
	// carimbo antigo continuaria no mapa, sem que nada os ligasse.
	if v == "dispensado" {
		v = "waived"
	}
	if v != "pass" && v != "fail" && v != "waived" {
		return fmt.Errorf("--verdict must be 'pass', 'fail' or 'waived' " +
			"(waived: the target of the question does not exist — the spec declares it `@TBD`)")
	}
	// O motivo é obrigatório nos dois desfechos que NÃO são aprovação: um `fail` sem
	// violação escrita não é acionável, e um `dispensado` sem a ausência nomeada é
	// indistinguível de um gate desligado.
	if strings.TrimSpace(reason) != "" {
		return nil
	}
	switch v {
	case "fail":
		return fmt.Errorf("--reason is mandatory on a 'fail' verdict (explain the violation)")
	case "waived":
		return fmt.Errorf("--reason is mandatory on a 'waived' verdict " +
			"(name the absence: which piece is missing, and where it is declared)")
	}
	return nil
}

// judgeWritesIssue says whether `anchors judge` opens and closes the issue of its verdict.
// Always, except in manual mode, where it does only with --record-issues — the verdict is
// stamped in the map either way.
func judgeWritesIssue(cfg *config.Config, recordIssues bool) bool {
	if cfg != nil && cfg.ManualMode() {
		return recordIssues
	}
	return true
}
