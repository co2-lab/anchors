package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/health"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	var root, mapPath string
	var corrigir bool
	cmd := &cobra.Command{
		Use: "doctor",

		Short: "Raio-X do ecossistema: pontas sistêmicas, saúde e maturidade",
		Long: `O validador de saúde do ecossistema (QUALITY §5.2) — a visão GLOBAL.
Diferente do check (nó contra gate, incremental), o doctor varre o mapa + a
Estrutura + o disco e caça as pontas SISTÊMICAS que nenhum gate local vê:
integridade do mapa (arestas mortas, nós fantasma), órfãos (código sem spec,
identidade ausente), camadas frouxas, e buracos de cobertura de gates.

APRESENTA e REGISTRA, mas NÃO bloqueia — é diagnóstico, roda sob demanda.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar config: %w", err)
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
			}

			rep := health.Diagnose(g, cfg, absRoot)
			printReport(rep)
			if corrigir {
				return repararAmbiente(absRoot, cfg)
			}
			return nil // doctor NUNCA bloqueia — só reporta
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().BoolVar(&corrigir, "fix", false, "cria o que falta no ambiente do modo `github` (pipelines)")
	return cmd
}

func printReport(r health.Report) {
	fmt.Printf("anchors doctor — %d nós, %d arestas, %d camadas\n\n", r.Nodes, r.Edges, r.Layers)

	warnings := r.Warnings()
	if len(r.Findings) == 0 {
		fmt.Println("✓ nenhuma ponta sistêmica encontrada — ecossistema íntegro")
		return
	}

	// agrupa por verificação
	byCheck := map[string][]health.Finding{}
	var order []string
	for _, f := range r.Findings {
		if _, ok := byCheck[f.Check]; !ok {
			order = append(order, f.Check)
		}
		byCheck[f.Check] = append(byCheck[f.Check], f)
	}

	for _, check := range order {
		fs := byCheck[check]
		// o grupo é Warn se QUALQUER finding for Warn (a mais alta vence) — senão o
		// ícone esconderia warns atrás de um info que veio primeiro na ordenação.
		mark := "ℹ"
		nWarn := 0
		for _, f := range fs {
			if f.Severity == health.Warn {
				nWarn++
			}
		}
		if nWarn > 0 {
			mark = "⚠"
		}
		fmt.Printf("%s %s (%d)\n", mark, check, len(fs))
		for _, f := range fs {
			itemMark := " "
			if f.Severity == health.Warn {
				itemMark = "⚠"
			}
			fmt.Printf("  %s %s — %s\n", itemMark, f.Subject, f.Detail)
		}
		fmt.Println()
	}

	fmt.Printf("resumo: %d ponta(s) de atenção, %d achado(s) no total\n",
		len(warnings), len(r.Findings))
	fmt.Println("(diagnóstico — nada foi bloqueado; decida o que conciliar)")
}

// repararAmbiente cria o que falta no ambiente do modo `github`. É o `--fix` do doctor,
// no precedente do `check --fix`: sem a flag, o doctor segue sendo diagnóstico puro
// ("nada foi bloqueado; decida o que conciliar"), e com ela ele age.
//
// Conserta só o que é DELE para consertar: os pipelines são arquivos do repositório,
// versionados e revisáveis num diff, e criá-los não afeta ninguém até serem commitados.
//
// O BOARD não é criado aqui, e a diferença é de natureza: um GitHub Project é estrutura
// COMPARTILHADA pelo time e vive fora do repositório — um board criado por engano polui a
// organização inteira e não se desfaz com `git checkout`. O doctor diz o que falta; criar
// é decisão de quem opera.
func repararAmbiente(root string, cfg *config.Config) error {
	if !cfg.ModoGitHub() {
		fmt.Println("\n--fix: nada a fazer — o ambiente do GitHub só é exigido no `workflow.mode: github`.")
		return nil
	}
	escritos, err := initx.SemeiaWorkflows(root)
	if err != nil {
		return fmt.Errorf("semear os pipelines: %w", err)
	}
	fmt.Println()
	if len(escritos) == 0 {
		fmt.Println("--fix: os pipelines já existem (nenhum foi sobrescrito).")
	} else {
		fmt.Printf("✓ %d pipeline(s) criados em %s:\n", len(escritos), initx.DirWorkflows)
		for _, e := range escritos {
			fmt.Printf("    %s\n", e)
		}
		fmt.Println("  revise, commite e configure `vars.ANCHORS_PROJECT_NUMBER` no repositório.")
	}
	// As LABELS de estado são o único pré-requisito real do fluxo — e criá-las é seguro:
	// label é do repositório, reversível, e não afeta ninguém fora dele.
	if err := criaLabelsDeEstado(cfg); err != nil {
		fmt.Printf("⚠  não deu para criar as labels de estado: %v\n", err)
	}

	// O BOARD é OPCIONAL: o estado vive na label, e o Project só espelha. Dizê-lo aqui
	// evita que alguém conclua que falta configurar algo.
	fmt.Println("\n  o board é OPCIONAL — o estado do trabalho vive nas labels acima.")
	fmt.Println("  para ter um: crie um GitHub Project com estas colunas e ligue a automação")
	fmt.Println("  nativa do Projects (label adicionada → move para a coluna):")
	fmt.Printf("    %s\n", strings.Join(initx.ColunasDoBoard, " · "))
	fmt.Printf("  o Anchors escreve até `%s`; as seguintes são dos pipelines de entrega.\n",
		initx.EstadoFinalDoAnchors)
	return nil
}

// criaLabelsDeEstado garante as labels que carregam o estado do trabalho. São o único
// pré-requisito do fluxo que não é arquivo — e sem elas o `identify` cria cards que o
// `claim` nunca encontra.
func criaLabelsDeEstado(cfg *config.Config) error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("o `gh` não está no PATH")
	}
	repo := cfg.Workflow.Repo
	// Cores por família: a fila em cinza, o trabalho em curso em azul, o que saiu da
	// alçada do Anchors em verde.
	cor := map[string]string{
		"anchors:to-do": "ededed", "anchors:ready-to-review": "ededed",
		"anchors:in-progress": "1d76db", "anchors:in-review": "1d76db",
		"anchors:ready-to-test": "0e8a16", "anchors:in-test": "0e8a16",
		"anchors:ready-to-release": "0e8a16", "anchors:production": "0e8a16",
	}
	var criadas int
	for _, e := range append([]string{cfg.Workflow.Labels[0]}, initx.EstadosDoTrabalho...) {
		c := cor[e]
		if c == "" {
			c = "5319e7" // a label do próprio Anchors
		}
		out, err := exec.Command("gh", "label", "create", e,
			"--repo", repo, "--color", c, "--force").CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", e, strings.TrimSpace(string(out)))
		}
		criadas++
	}
	fmt.Printf("✓ %d label(s) de estado garantidas em %s\n", criadas, repo)
	return nil
}
