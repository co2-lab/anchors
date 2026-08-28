package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/health"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// `anchors audit <arquivo>` é o DOSSIÊ POR ARQUIVO — tudo que pende sobre um alvo num
// só lugar: os gates de qualidade que se aplicam a ele E os achados sistêmicos do
// doctor sobre ele. Feito para o trabalho paralelo: um agente pega um arquivo, roda
// audit, vê TODAS as suas pendências (não só o header), conserta tudo de uma vez, e
// passa ao próximo. Vários agentes varrem o projeto, um arquivo cada.
//
// Default: só o arquivo. Com --impact: estende ao caminho de impacto (a unidade —
// spec↔code↔feature↔test — e o que o alvo propaga/valida).
func newAuditCmd() *cobra.Command {
	var root, mapPath string
	var withImpact bool
	cmd := &cobra.Command{
		Use:   "audit <arquivo>",
		Short: "Dossiê de pendências de UM arquivo (gates + doctor), para correção em lote",
		Long: `Reúne, para um arquivo, TUDO que pende sobre ele:
  • os gates de qualidade que se aplicam ao seu kind (o que 'check' rodaria)
  • os achados sistêmicos do 'doctor' que o citam (colisão, órfão, identidade…)

Default: só o próprio arquivo. Com --impact, inclui os nós do caminho de impacto
(a trinca da unidade e o que o alvo propaga/valida) — para consertar a unidade toda.

Pensado para varredura PARALELA: um agente por arquivo, 'anchors audit <arquivo>',
conserta todas as pendências de uma vez (header, spec, identidade, cobertura…).`,
		Args: cobra.ExactArgs(1),
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
			target := relTo(absRoot, args[0])
			if !nodeExists(g, target) {
				return fmt.Errorf("arquivo %q não está no mapa", target)
			}

			// o conjunto de nós a auditar: o alvo, + o caminho de impacto se pedido.
			ids := map[string]bool{target: true}
			if withImpact {
				imp := g.AnalyzeImpact(target)
				for _, n := range imp.Propagate {
					ids[n] = true
				}
				for _, n := range imp.Validate {
					ids[n] = true
				}
			}
			var nodes []mapx.Node
			for _, n := range g.Nodes {
				if ids[n.ID] {
					nodes = append(nodes, n)
				}
			}

			// 1) GATES sobre esses nós
			gateResults := gate.RunWithConfig(cfg.Gates, nodes, absRoot, g, cfg)
			// 2) DOCTOR — os achados que citam algum desses nós
			rep := health.Diagnose(g, cfg, absRoot)

			return printAudit(target, withImpact, nodes, gateResults, rep, ids)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().BoolVar(&withImpact, "impact", false, "inclui o caminho de impacto (a unidade toda), não só o arquivo")
	return cmd
}

func printAudit(target string, withImpact bool, nodes []mapx.Node, results []gate.Result, rep health.Report, ids map[string]bool) error {
	scope := "o arquivo"
	if withImpact {
		scope = fmt.Sprintf("a unidade (%d nós no caminho de impacto)", len(nodes))
	}
	fmt.Printf("audit: %s — %s\n\n", target, scope)

	// agrupa por nó, para o dossiê ser "por arquivo"
	byNode := map[string][]string{} // nó → linhas de pendência
	pending := 0

	for _, r := range results {
		if r.Verdict == gate.Pass || r.Verdict == gate.Skip {
			continue
		}
		mark := map[gate.Verdict]string{gate.Fail: "✗", gate.Pending: "~", gate.Judge: "⏳"}[r.Verdict]
		line := fmt.Sprintf("  %s [gate] %s — %s", mark, r.Gate, firstLineAudit(r.Detail))
		byNode[r.Target] = append(byNode[r.Target], line)
		if r.Verdict == gate.Fail {
			pending++
		}
	}
	for _, f := range rep.Findings {
		if !ids[f.Subject] {
			continue // só achados que citam um nó do escopo
		}
		sev := "ℹ"
		if f.Severity == health.Warn {
			sev = "⚠"
		}
		byNode[f.Subject] = append(byNode[f.Subject], fmt.Sprintf("  %s [doctor:%s] %s", sev, f.Check, f.Detail))
		if f.Severity == health.Warn {
			pending++
		}
	}

	if len(byNode) == 0 {
		fmt.Println("✓ nada pendente — o arquivo (e o escopo) está conforme.")
		return nil
	}
	// imprime o alvo primeiro, depois os demais nós do impacto
	order := []string{target}
	for id := range byNode {
		if id != target {
			order = append(order, id)
		}
	}
	for _, id := range order {
		lines := byNode[id]
		if len(lines) == 0 {
			continue
		}
		if id == target {
			fmt.Printf("● %s\n", id)
		} else {
			fmt.Printf("○ %s (impacto)\n", id)
		}
		for _, l := range lines {
			fmt.Println(l)
		}
		fmt.Println()
	}
	fmt.Printf("%d pendência(s) acionável(is) — conserte todas antes de passar ao próximo arquivo.\n", pending)
	return nil
}

func firstLineAudit(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
