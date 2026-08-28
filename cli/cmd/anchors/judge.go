package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	var pending bool
	cmd := &cobra.Command{
		Use:   "judge <alvo>",
		Short: "Registra o veredito de uma IA para um gate de julgamento",
		Long: `Grava o veredito de um gate de JULGAMENTO POR IA sobre um alvo — a mesma
dupla saída de um gate determinístico (carimbo no mapa + issue).

  anchors judge --pending                      lista os alvos aguardando julgamento
  anchors judge <alvo> --gate <g> --verdict pass|fail --reason "..."

O fluxo: o 'anchors check' marca os alvos de um gate 'measures: judgment' como
pendentes e os enfileira. A IA (worker) lê o guide do gate, confronta o alvo, e
reporta aqui.

IMPORTANTE — não desperdice o que você já leu. Para julgar, você leu o guide inteiro
e o alvo inteiro; então o --reason NÃO é uma frase de veredito, é o LAUDO COMPLETO. O
corpo da issue É o seu --reason, verbatim (aceita markdown multi-linha). Num 'fail',
liste CADA não-conformidade com: o que está errado, ONDE (arquivo:linha), POR QUÊ
(qual regra do guide), e COMO corrigir. Assim ninguém precisa reprocessar o alvo
depois só para saber o que arrumar. Ex.:

  anchors judge apps/.../AlertsScreen.tsx --gate atomic-design --verdict fail \
    --reason "$(cat <<'EOF'
  ## Laudo — atomic-design @ AlertsScreen
  ### 1. Item de lista inline deveria ser molécula  (l.151-166)
  - **o quê:** o JSX do item (View + ícone + textos) está montado inline no .map
  - **regra:** SCREEN_GUIDE §Pages — telas injetam componentes, não montam blocos
  - **como:** extrair para components/molecules/AlertRow.tsx e usar <AlertRow item=.../>
  ### 2. ...
  EOF
  )"

'fail' abre a issue com esse laudo; 'pass' resolve a issue se havia uma (e o --reason,
se houver, vira observação registrada).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
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
				return fmt.Errorf("informe o alvo (ex: `anchors judge Foo.tsx --gate atomic --verdict fail --reason ...`) ou use --pending")
			}
			if gateName == "" {
				return fmt.Errorf("--gate é obrigatório (o gate de julgamento que você está avaliando)")
			}
			v := strings.ToLower(verdict)
			if v != "pass" && v != "fail" {
				return fmt.Errorf("--verdict deve ser 'pass' ou 'fail'")
			}
			if v == "fail" && strings.TrimSpace(reason) == "" {
				return fmt.Errorf("--reason é obrigatório num veredito 'fail' (explique a violação)")
			}

			target := relTo(absRoot, args[0])
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}
			// O review acontece sobre trabalho NOVO, e no início da unidade só a spec
			// existe: o `anchors work review --for <alvo>.ts` prescreve o caminho do
			// código, que ainda não foi escrito. Recusar aí deixa o revisor sem onde
			// registrar o achado — medido num E2E, foi o que aconteceu.
			//
			// Cair para a peça que EXISTE preserva a identidade (spec e código são a mesma
			// unidade) e mantém o comando prescrito funcionando desde a primeira etapa.
			if !nodeExists(g, target) {
				if alt := pecaExistenteDaUnidade(g, target); alt != "" {
					fmt.Printf("   (o alvo ainda não existe; registrando em `%s`, a peça desta unidade que já está no mapa)\n", alt)
					target = alt
				}
			}
			if !nodeExists(g, target) {
				return fmt.Errorf("alvo %q não está no mapa", target)
			}
			// valida que o gate existe e é de julgamento
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar config: %w", err)
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
				return fmt.Errorf("gate %q não existe ou não é de julgamento (measures: judgment)", gateName)
			}

			now := time.Now()
			failed := v == "fail"
			verdictStr := "ok"
			if failed {
				verdictStr = "issue"
			}
			// 1) CARIMBO — marca a aresta guide→alvo (o confronto da régua contra o
			//    alvo) com o veredito da IA. Se o gate declara guide, carimba essa
			//    aresta; senão, cai para o carimbo por-nó (todas as arestas do alvo).
			stamped := 0
			if gc.Guide != "" && g.StampEdge(gc.Guide, target, verdictStr, now.Format(time.RFC3339)) {
				stamped = 1
			} else {
				// Carimbo por NÓ: o `judge` julga uma UNIDADE, e o `StampEdges` exige as
				// duas pontas confrontadas na mesma rodada — com um alvo só, nenhuma
				// aresta qualifica e o carimbo saía sempre zero. O veredito abria a issue
				// e não tocava o grafo: um `check` posterior não via o achado.
				// leva o NOME do gate: é ele que permite ao `check` posterior saber que
				// este julgamento foi respondido, em vez de perguntar de novo.
				stamped = g.StampNodeByGate(target, verdictStr, now.Format(time.RFC3339), gateName)
			}
			if err := mapx.Save(g, mapPath); err != nil {
				return fmt.Errorf("salvar carimbo: %w", err)
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
					return fmt.Errorf("ler o patch: %w", rerr)
				}
				sg := suggestion.Suggestion{
					ID:   "judge-" + gc.Name + "-" + strings.ReplaceAll(relSlug(target), "/", "-"),
					Gate: gc.Name, Target: target, Origin: suggestion.FromJudgment,
					Why: reason, Patch: string(diff),
				}
				criada, at, serr := suggestion.Open(absRoot, sg)
				if serr != nil {
					return serr
				}
				if criada {
					fmt.Printf("  ↳ sugestão de correção em %s (`anchors suggest show %s`)\n", at, sg.ID)
				}
			}

			if failed {
				created, at, err := issue.Open(absRoot, iss)
				if err != nil {
					return err
				}
				if created {
					fmt.Printf("✗ julgado FAIL — issue aberta em %s/todo/\n", issue.Dir)
				} else if reaberta, rerr := issue.Reabrir(absRoot, iss); rerr != nil {
					return rerr
				} else if reaberta {
					// Achado NOVO sobre unidade já revisada: reabre e acrescenta o laudo.
					// Antes o `--reason` era descartado em silêncio, e o defeito não
					// ficava em lugar nenhum.
					fmt.Printf("✗ julgado FAIL — achado NOVO acrescentado; issue reaberta em %s/todo/\n", issue.Dir)
				} else {
					fmt.Printf("✗ julgado FAIL — mesmo achado já registrado (%s/), nada a fazer\n", at)
				}
			} else {
				if ok, _ := issue.Resolve(absRoot, iss.Key()); ok {
					fmt.Printf("✓ julgado PASS — issue anterior resolvida (→ %s/done/)\n", issue.Dir)
				} else {
					fmt.Printf("✓ julgado PASS\n")
				}
			}
			fmt.Printf("  carimbado: %d aresta(s) com o veredito de IA (gate '%s')\n", stamped, gc.Name)

			// 3) fecha a task judge correspondente, se houver
			closeJudgeTask(absRoot, gc.Name, target)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().StringVar(&gateName, "gate", "", "o gate de julgamento avaliado")
	cmd.Flags().StringVar(&verdict, "verdict", "", "pass | fail")
	cmd.Flags().StringVar(&reason, "reason", "", "o LAUDO completo (markdown multi-linha) — vira o corpo da issue; obrigatório em fail")
	cmd.Flags().StringVar(&patchFile, "patch", "", "arquivo com o diff que CORRIGE o achado — abre uma sugestão aplicável (`anchors suggest`)")
	cmd.Flags().BoolVar(&pending, "pending", false, "lista os alvos aguardando julgamento")
	return cmd
}

func findJudgmentGate(cfg *config.Config, name string) (config.Gate, bool) {
	for _, g := range cfg.Gates {
		if g.Name == name && g.IsJudgment() {
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
		fmt.Println("nenhum alvo aguardando julgamento (rode `anchors check` para descobrir)")
	} else {
		fmt.Printf("\n%d alvo(s) — julgue com: anchors judge <alvo> --gate <g> --verdict pass|fail --reason ...\n", n)
	}
	return nil
}

// closeJudgeTask fecha a task judge daquele gate+alvo, se estiver na fila.
func closeJudgeTask(root, gateName, target string) {
	id := "judge-" + gateName + "-" + strings.ReplaceAll(relSlug(target), "/", "-")
	_ = queue.MarkDone(root, id)
}

// pecaExistenteDaUnidade acha, para um alvo ausente do mapa, outra peça da MESMA unidade
// que já esteja lá — a spec, tipicamente, quando o código ainda não nasceu.
func pecaExistenteDaUnidade(g *mapx.Graph, target string) string {
	base := strings.TrimSuffix(target, filepath.Ext(target))
	for _, suf := range []string{".spec.md", ".feature", ".test.ts", ".test.tsx", ".ts", ".tsx"} {
		if cand := base + suf; cand != target && nodeExists(g, cand) {
			return cand
		}
	}
	return ""
}
