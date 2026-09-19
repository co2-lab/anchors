package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// `anchors report` gera o PAINEL DE CONFIANÇA em docs/ — o retrato do estado dos
// testes do entregável, consolidado do grafo (não reparseia o JUnit; lê os Signal já
// ingeridos). Mergeia as suites de CAMADAS distintas (unit/integration/e2e) num só
// relatório, rotulado por camada. É um DOC TERMINAL — consumo humano, versionado no
// git para dar história; NÃO é uma âncora (o Anchors não rege o próprio output).
// reportSpec descreve um subcomando de relatório: o nome, o arquivo default em docs/,
// e o renderer (grafo+config → markdown). Cada perspectiva é um recorte das fontes
// que já medimos; nenhum inventa dado.
type reportSpec struct {
	name   string
	short  string
	file   string
	render func(ctx reportCtx) string
}

// reportCtx são as fontes que um renderer pode usar.
type reportCtx struct {
	g    *mapx.Graph
	cfg  *config.Config
	root string
	when string
}

var reportSpecs = []reportSpec{
	{"tests", "test confidence: execution/coverage/scenario (merges layers)", "anchors-test-report.md", renderTests},
	{"quality", "quality gate verdict + debt per gate", "anchors-quality-report.md", renderQuality},
	{"structure", "layers, governance, identity collisions, orphans", "anchors-structure-report.md", renderStructure},
	{"config", "state of anchors.yaml: layers/derived/governs/gates and what is missing", "anchors-config-report.md", renderConfig},
	{"issues", "issues (todo/doing/done) and tasks (the queue) — the tracked work", "anchors-issues-report.md", renderIssues},
	{"inconsistencies", "the list of things to fix: every validator error/warning", "anchors-inconsistencies-report.md", renderInconsistencies},
}

func newReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report <perspective>",
		Short: "Generate reports in docs/ by perspective (cuts of what Anchors measures)",
		Long: `Generates markdown reports in docs/, each one a PERSPECTIVE on the project —
cuts of the sources Anchors already measures (graph, signals, gates, issues). They are
TERMINAL docs (human consumption, versioned in git to give history); they do not become anchors.

  anchors report tests            test confidence (execution/coverage/scenario)
  anchors report quality          gate verdict + debt
  anchors report structure        layers, governance, collisions, orphans
  anchors report config           state of anchors.yaml and what is missing
  anchors report issues           issues (todo/doing/done) + tasks (queue)
  anchors report inconsistencies  the list of things to fix (errors/warnings)
  anchors report all              generates them all + an index in docs/anchors/`,
	}
	for _, sp := range reportSpecs {
		cmd.AddCommand(newReportSubCmd(sp))
	}
	cmd.AddCommand(newReportAllCmd())
	return cmd
}

// newReportSubCmd fabrica um subcomando a partir de um reportSpec (casca comum:
// carrega grafo+config, renderiza, escreve em docs/).
func newReportSubCmd(sp reportSpec) *cobra.Command {
	var root, mapPath, out string
	c := &cobra.Command{
		Use:   sp.name,
		Short: sp.short,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := loadReportCtx(root, mapPath)
			if err != nil {
				return err
			}
			dest := out
			if dest == "" {
				dest = filepath.Join(ctx.root, "docs", sp.file)
			}
			if err := writeReport(dest, sp.render(ctx)); err != nil {
				return err
			}
			rel, _ := filepath.Rel(ctx.root, dest)
			fmt.Printf("report generated: %s\n", rel)
			return nil
		},
	}
	c.Flags().StringVar(&root, "root", ".", "project root")
	c.Flags().StringVar(&mapPath, "map", "", "path to the map")
	c.Flags().StringVar(&out, "out", "", "output file (default docs/"+sp.file+")")
	return c
}

func newReportAllCmd() *cobra.Command {
	var root, mapPath string
	c := &cobra.Command{
		Use:   "all",
		Short: "generate every perspective + an index in docs/anchors/",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := loadReportCtx(root, mapPath)
			if err != nil {
				return err
			}
			dir := filepath.Join(ctx.root, "docs", "anchors")
			var idx strings.Builder
			idx.WriteString("# Anchors reports\n\n")
			fmt.Fprintf(&idx, "> Generated on %s by `anchors report all`.\n\n", ctx.when)
			for _, sp := range reportSpecs {
				dest := filepath.Join(dir, sp.file)
				if err := writeReport(dest, sp.render(ctx)); err != nil {
					return err
				}
				fmt.Fprintf(&idx, "- [%s](%s) — %s\n", sp.name, sp.file, sp.short)
			}
			if err := writeReport(filepath.Join(dir, "index.md"), idx.String()); err != nil {
				return err
			}
			fmt.Printf("%d report(s) + index generated in docs/anchors/\n", len(reportSpecs))
			return nil
		},
	}
	c.Flags().StringVar(&root, "root", ".", "project root")
	c.Flags().StringVar(&mapPath, "map", "", "path to the map")
	return c
}

func loadReportCtx(root, mapPath string) (reportCtx, error) {
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return reportCtx{}, err
	}
	if mapPath == "" {
		mapPath = filepath.Join(absRoot, mapx.DefaultPath)
	}
	g, err := mapx.Load(mapPath)
	if err != nil {
		return reportCtx{}, fmt.Errorf("load map: %w (run `anchors map build`)", err)
	}
	cfg, _ := config.Load(filepath.Join(absRoot, config.DefaultFile)) // pode faltar
	return reportCtx{g: g, cfg: cfg, root: absRoot, when: time.Now().Format("2006-01-02 15:04")}, nil
}

func writeReport(dest, body string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, []byte(body), 0o644)
}

func reportHeader(title, when string) string {
	return fmt.Sprintf("# %s\n\n> Generated on %s by `anchors report`, from what Anchors measures.\n> Human consumption — regenerated on every `anchors report` (do not edit by hand).\n\n", title, when)
}

func renderTests(ctx reportCtx) string {
	g, root, when := ctx.g, ctx.root, ctx.when
	var b strings.Builder
	b.WriteString(reportHeader("Test report", when))

	// --- 1. Execução, mergeada por CAMADA ---
	layerTot := map[string]*mapx.LayerExec{}
	var testFiles, stale int
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindTest || n.Signal == nil {
			continue
		}
		if n.SignalStale() {
			stale++
		}
		if len(n.Signal.ByLayer) > 0 {
			testFiles++
			for layer, le := range n.Signal.ByLayer {
				if layerTot[layer] == nil {
					layerTot[layer] = &mapx.LayerExec{}
				}
				layerTot[layer].Passed += le.Passed
				layerTot[layer].Failed += le.Failed
				layerTot[layer].Skipped += le.Skipped
			}
		}
	}
	b.WriteString("## Execution (by layer)\n\n")
	if len(layerTot) == 0 {
		b.WriteString("_No execution result ingested._ Run the suite and `anchors ingest --junit <r> --layer <layer>`.\n\n")
	} else {
		b.WriteString("| layer | ✓ passed | ✗ failed | ⊘ skipped |\n|---|---:|---:|---:|\n")
		var gp, gf, gs int
		for _, layer := range sortedKeys(layerTot) {
			t := layerTot[layer]
			fmt.Fprintf(&b, "| %s | %d | %d | %d |\n", layer, t.Passed, t.Failed, t.Skipped)
			gp += t.Passed
			gf += t.Failed
			gs += t.Skipped
		}
		fmt.Fprintf(&b, "| **total** | **%d** | **%d** | **%d** |\n\n", gp, gf, gs)
		if gf > 0 {
			fmt.Fprintf(&b, "⚠ **%d failing test(s)** — the deliverable is not green.\n\n", gf)
		}
	}
	if stale > 0 {
		fmt.Fprintf(&b, "⚠ %d test file(s) with a STALE signal (they changed since ingestion — re-ingest).\n\n", stale)
	}

	// --- 2. Cobertura por cenário (o diferencial) ---
	b.WriteString("## Coverage by scenario (proven requirements)\n\n")
	// Distingue "medido" (a spec teve sinal ingerido) de "não medido" (nenhuma
	// ingestão tocou). Só specs MEDIDAS entram na conta de provado/lacuna — senão o
	// relatório fingiria que specs nunca testadas "falharam".
	var specGaps, measuredSpecs, unmeasuredSpecs, provTotal, declMeasured int
	var gapLines []string
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindSpec {
			continue
		}
		declared, _ := codesInFileOfUnit(filepath.Join(root, n.ID), n.Code)
		if len(declared) == 0 {
			continue
		}
		measured := n.Signal != nil && len(n.Signal.ProvenCodes) > 0
		if !measured {
			unmeasuredSpecs++
			continue
		}
		measuredSpecs++
		declMeasured += len(declared)
		proven := map[string]bool{}
		for _, c := range n.Signal.ProvenCodes {
			proven[c] = true
		}
		var missing []string
		for _, c := range declared {
			if proven[c] {
				provTotal++
			} else {
				missing = append(missing, c)
			}
		}
		if len(missing) > 0 {
			specGaps++
			sort.Strings(missing)
			if len(gapLines) < 20 {
				gapLines = append(gapLines, fmt.Sprintf("- `%s` — %d/%d unproven: %s", n.ID, len(missing), len(declared), strings.Join(missing, ", ")))
			}
		}
	}
	if measuredSpecs == 0 {
		fmt.Fprintf(&b, "_No spec measured yet_ — %d spec(s) with scenarios await test ingestion.\n\n", unmeasuredSpecs)
	} else {
		fmt.Fprintf(&b, "**%d/%d requirements proven** by a green test, across %d measured spec(s). %d spec(s) with a gap:\n\n", provTotal, declMeasured, measuredSpecs, specGaps)
		if len(gapLines) == 0 {
			b.WriteString("- (none — every measured scenario has a green test)\n")
		} else {
			b.WriteString(strings.Join(gapLines, "\n") + "\n")
			if specGaps > len(gapLines) {
				fmt.Fprintf(&b, "- … and %d more spec(s)\n", specGaps-len(gapLines))
			}
		}
		if unmeasuredSpecs > 0 {
			fmt.Fprintf(&b, "\n_%d spec(s) still NOT measured_ (no ingestion touched them) — they do not count above.\n", unmeasuredSpecs)
		}
		b.WriteString("\n")
	}

	// --- 3. Cobertura de linha + regressões ---
	b.WriteString("## Line coverage\n\n")
	const threshold = 70.0
	var below, regressed int
	var belowLines, regLines []string
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindCode || n.Signal == nil || n.Signal.TotalLines == 0 {
			continue
		}
		if n.Signal.LineCoverage < threshold {
			below++
			if len(belowLines) < 15 {
				belowLines = append(belowLines, fmt.Sprintf("- `%s` — %.0f%%", n.ID, n.Signal.LineCoverage))
			}
		}
		if n.Signal.PrevLineCoverage > 0 && n.Signal.LineCoverage < n.Signal.PrevLineCoverage-0.01 {
			regressed++
			regLines = append(regLines, fmt.Sprintf("- `%s` — %.0f%% → %.0f%% (%.0f)", n.ID, n.Signal.PrevLineCoverage, n.Signal.LineCoverage, n.Signal.LineCoverage-n.Signal.PrevLineCoverage))
		}
	}
	fmt.Fprintf(&b, "%d file(s) below %.0f%%:\n\n", below, threshold)
	if len(belowLines) == 0 {
		b.WriteString("- (none, or line coverage not ingested)\n")
	} else {
		b.WriteString(strings.Join(belowLines, "\n") + "\n")
		if below > len(belowLines) {
			fmt.Fprintf(&b, "- … and %d more\n", below-len(belowLines))
		}
	}
	b.WriteString("\n### Coverage regressions\n\n")
	if regressed == 0 {
		b.WriteString("- (no drop since the previous ingestion)\n")
	} else {
		b.WriteString(strings.Join(regLines, "\n") + "\n")
	}

	b.WriteString("\n---\n_Anchors — the deliverable's confidence panel._\n")
	return b.String()
}

func sortedKeys[V any](m map[string]V) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
