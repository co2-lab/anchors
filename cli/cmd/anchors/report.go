package main

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
	{"tests", "confiança dos testes: execução/cobertura/cenário (mergeia camadas)", "anchors-test-report.md", renderTests},
	{"quality", "veredito dos gates de qualidade + débito por gate", "anchors-quality-report.md", renderQuality},
	{"structure", "camadas, governança, colisões de identidade, órfãos", "anchors-structure-report.md", renderStructure},
	{"config", "estado do anchors.yaml: layers/derived/governs/gates e o que falta", "anchors-config-report.md", renderConfig},
	{"issues", "issues (todo/doing/done) e tasks (a fila) — o trabalho rastreado", "anchors-issues-report.md", renderIssues},
	{"inconsistencies", "a lista de coisas a arrumar: todo erro/warning dos validators", "anchors-inconsistencies-report.md", renderInconsistencies},
}

func newReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report <perspectiva>",
		Short: "Gera relatórios em docs/ por perspectiva (recortes do que o Anchors mede)",
		Long: `Gera relatórios markdown em docs/, cada um uma PERSPECTIVA do projeto —
recortes das fontes que o Anchors já mede (grafo, sinais, gates, issues). São docs
TERMINAIS (consumo humano, versionados no git para dar história); não viram âncoras.

  anchors report tests            confiança dos testes (execução/cobertura/cenário)
  anchors report quality          veredito dos gates + débito
  anchors report structure        camadas, governança, colisões, órfãos
  anchors report config           estado do anchors.yaml e o que falta
  anchors report issues           issues (todo/doing/done) + tasks (fila)
  anchors report inconsistencies  a lista de coisas a arrumar (erros/warnings)
  anchors report all              gera todos + um índice em docs/anchors/`,
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
			fmt.Printf("relatório gerado: %s\n", rel)
			return nil
		},
	}
	c.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	c.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	c.Flags().StringVar(&out, "out", "", "arquivo de saída (default docs/"+sp.file+")")
	return c
}

func newReportAllCmd() *cobra.Command {
	var root, mapPath string
	c := &cobra.Command{
		Use:   "all",
		Short: "gera todas as perspectivas + um índice em docs/anchors/",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := loadReportCtx(root, mapPath)
			if err != nil {
				return err
			}
			dir := filepath.Join(ctx.root, "docs", "anchors")
			var idx strings.Builder
			idx.WriteString("# Relatórios do Anchors\n\n")
			fmt.Fprintf(&idx, "> Gerados em %s por `anchors report all`.\n\n", ctx.when)
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
			fmt.Printf("%d relatório(s) + índice gerados em docs/anchors/\n", len(reportSpecs))
			return nil
		},
	}
	c.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	c.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	return c
}

func loadReportCtx(root, mapPath string) (reportCtx, error) {
	absRoot, err := config.AbsRaiz(root)
	if err != nil {
		return reportCtx{}, err
	}
	if mapPath == "" {
		mapPath = filepath.Join(absRoot, mapx.DefaultPath)
	}
	g, err := mapx.Load(mapPath)
	if err != nil {
		return reportCtx{}, fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
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
	return fmt.Sprintf("# %s\n\n> Gerado em %s por `anchors report`, a partir do que o Anchors mede.\n> Consumo humano — regenerado a cada `anchors report` (não editar à mão).\n\n", title, when)
}

func renderTests(ctx reportCtx) string {
	g, root, when := ctx.g, ctx.root, ctx.when
	var b strings.Builder
	b.WriteString(reportHeader("Relatório de testes", when))

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
	b.WriteString("## Execução (por camada)\n\n")
	if len(layerTot) == 0 {
		b.WriteString("_Nenhum resultado de execução ingerido._ Rode a suíte e `anchors ingest --junit <r> --layer <camada>`.\n\n")
	} else {
		b.WriteString("| camada | ✓ passou | ✗ falhou | ⊘ pulou |\n|---|---:|---:|---:|\n")
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
			fmt.Fprintf(&b, "⚠ **%d teste(s) falhando** — o entregável não está verde.\n\n", gf)
		}
	}
	if stale > 0 {
		fmt.Fprintf(&b, "⚠ %d arquivo(s) de teste com sinal STALE (mudaram desde a ingestão — reingira).\n\n", stale)
	}

	// --- 2. Cobertura por cenário (o diferencial) ---
	b.WriteString("## Cobertura por cenário (requisitos provados)\n\n")
	// Distingue "medido" (a spec teve sinal ingerido) de "não medido" (nenhuma
	// ingestão tocou). Só specs MEDIDAS entram na conta de provado/lacuna — senão o
	// relatório fingiria que specs nunca testadas "falharam".
	var specGaps, measuredSpecs, unmeasuredSpecs, provTotal, declMeasured int
	var gapLines []string
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindSpec {
			continue
		}
		declared, _ := codesInFile(filepath.Join(root, n.ID))
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
				gapLines = append(gapLines, fmt.Sprintf("- `%s` — %d/%d sem prova: %s", n.ID, len(missing), len(declared), strings.Join(missing, ", ")))
			}
		}
	}
	if measuredSpecs == 0 {
		fmt.Fprintf(&b, "_Nenhuma spec medida ainda_ — %d spec(s) com cenários aguardam ingestão de teste.\n\n", unmeasuredSpecs)
	} else {
		fmt.Fprintf(&b, "**%d/%d requisitos provados** por teste verde, em %d spec(s) medida(s). %d spec(s) com lacuna:\n\n", provTotal, declMeasured, measuredSpecs, specGaps)
		if len(gapLines) == 0 {
			b.WriteString("- (nenhuma — todos os cenários medidos têm teste verde)\n")
		} else {
			b.WriteString(strings.Join(gapLines, "\n") + "\n")
			if specGaps > len(gapLines) {
				fmt.Fprintf(&b, "- … e mais %d spec(s)\n", specGaps-len(gapLines))
			}
		}
		if unmeasuredSpecs > 0 {
			fmt.Fprintf(&b, "\n_%d spec(s) ainda NÃO medidas_ (nenhuma ingestão as tocou) — não contam acima.\n", unmeasuredSpecs)
		}
		b.WriteString("\n")
	}

	// --- 3. Cobertura de linha + regressões ---
	b.WriteString("## Cobertura de linha\n\n")
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
	fmt.Fprintf(&b, "%d arquivo(s) abaixo de %.0f%%:\n\n", below, threshold)
	if len(belowLines) == 0 {
		b.WriteString("- (nenhum, ou cobertura de linha não ingerida)\n")
	} else {
		b.WriteString(strings.Join(belowLines, "\n") + "\n")
		if below > len(belowLines) {
			fmt.Fprintf(&b, "- … e mais %d\n", below-len(belowLines))
		}
	}
	b.WriteString("\n### Regressões de cobertura\n\n")
	if regressed == 0 {
		b.WriteString("- (nenhuma queda desde a ingestão anterior)\n")
	} else {
		b.WriteString(strings.Join(regLines, "\n") + "\n")
	}

	b.WriteString("\n---\n_Anchors — painel de confiança do entregável._\n")
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
