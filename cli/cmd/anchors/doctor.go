package main

import (
	"fmt"
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
	// O que o --fix NÃO faz precisa aparecer aqui, junto do que ele fez: um usuário que
	// vê "✓ pronto" e não lê o resto conclui que o ambiente está completo.
	fmt.Println("\n  o BOARD não é criado pelo --fix (é estrutura compartilhada do time).")
	fmt.Printf("  crie um GitHub Project com o campo `Status` contendo, nesta ordem:\n    %s\n",
		strings.Join(initx.ColunasDoBoard, " · "))
	fmt.Printf("  o Anchors move até `%s`; as seguintes são dos pipelines de entrega do projeto.\n",
		initx.ColunaFinalDoAnchors)
	return nil
}
