// @anchors
//   code: IMCMM
//   ref: MPCTI

package mapcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

func newImpactCmd() *cobra.Command {
	var root, mapPath, propagate, format, layer string
	var stopAt []string
	cmd := &cobra.Command{
		Use:   "impact <file>...",
		Short: "Impact analysis: who depends on these files, and the tests to run again",
		Long: `Queries the map and shows, for one or more files, who depends on them — always the side that
depends: whoever imports it or declares @dep of it, reads what it writes, captures it,
realizes it; for a spec, the code, feature and tests it rules.

  • Direct — what depends on the file itself (the default).
  • --propagate indirect — one level further: what depends on the direct ones.
  • --propagate cascade  — the whole chain, up to @noPropagation nodes.

Each item past the first level carries the path that ties it to the file, and an item many
others depend on says how many it pulls: judge whether that branch belongs (--stop-at
<glob> lists it and climbs no further from it). The tests to run again are those reached
and those whose evidence closure holds a reached file — what evidence-fresh judges —, by the
layer of the suite that runs them.

  • Validates against — what the file is confronted with (its spec, feature, guides).

  --format text    grouped by level and layer (default)
  --format files   the tests to run again, one per line (with --layer, one layer's)
  --format json    levels, paths, tests and rulers

Query only — it opens no issue and changes nothing.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q, err := readTiesQuery(root, mapPath, propagate, format, stopAt, args)
			if err != nil {
				return err
			}
			// The suites say the layer of each test; a project with no configuration still
			// gets the dependents, with the tests unassigned.
			cfg, _ := config.Load(filepath.Join(q.absRoot, config.DefaultFile))
			rep := impactReport(q.g, cfg, q.targets, q.maxLevel, q.stop)
			switch format {
			case "files":
				for _, t := range rep.Tests {
					if layer == "" || slices.Contains(t.Layers, layer) {
						fmt.Println(t.File)
					}
				}
			case "json":
				b, err := json.MarshalIndent(rep, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(b))
			default:
				printImpact(rep, layer, propagate)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map (default: <root>/anchors.graph.yaml)")
	tiesFlags(cmd, &propagate, &format, &stopAt)
	cmd.Flags().StringVar(&layer, "layer", "", "only the tests of this layer (the `layer:` of a suite in `tests:`)")
	return common.TakesFiles(cmd)
}

// tiesFlags are the flags `impact` and `dependencies` share.
func tiesFlags(cmd *cobra.Command, propagate, format *string, stopAt *[]string) {
	cmd.Flags().StringVar(propagate, "propagate", "", "past the direct level: indirect (one level more) or cascade (the whole chain)")
	cmd.Flags().StringSliceVar(stopAt, "stop-at", nil, "glob of nodes listed but not walked past (repeatable or comma-separated)")
	cmd.Flags().StringVar(format, "format", "text", "text | files | json")
}

// tiesQuery is what `impact` and `dependencies` read before walking: the map, the files asked
// about as node ids, how far to walk, and where to stop.
type tiesQuery struct {
	absRoot  string
	g        *mapx.Graph
	targets  []string
	maxLevel int
	stop     func(string) bool
}

func readTiesQuery(root, mapPath, propagate, format string, stopAt, args []string) (tiesQuery, error) {
	q := tiesQuery{maxLevel: 1}
	switch propagate {
	case "":
	case "indirect":
		q.maxLevel = 2
	case "cascade":
		q.maxLevel = 0
	default:
		return q, fmt.Errorf("--propagate %q: indirect or cascade (without it, the direct level)", propagate)
	}
	switch format {
	case "text", "files", "json":
	default:
		return q, fmt.Errorf("--format %q: text, files or json", format)
	}
	for _, p := range stopAt {
		if !doublestar.ValidatePattern(p) {
			return q, fmt.Errorf("--stop-at %q is not a valid glob", p)
		}
	}
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return q, err
	}
	if mapPath == "" {
		mapPath = filepath.Join(absRoot, mapx.DefaultPath)
	}
	g, err := mapx.Load(mapPath)
	if err != nil {
		return q, fmt.Errorf("load map: %w (run `anchors map build`)", err)
	}
	q.absRoot, q.g = absRoot, g
	for _, a := range args {
		target := relTo(absRoot, a)
		if !nodeExists(g, target) {
			return q, fmt.Errorf("file %q is not in the map", target)
		}
		if !slices.Contains(q.targets, target) {
			q.targets = append(q.targets, target)
		}
	}
	q.stop = func(id string) bool {
		for _, p := range stopAt {
			if ok, _ := doublestar.Match(p, id); ok {
				return true
			}
		}
		return false
	}
	return q, nil
}

// tieItems turns a walk into the report's items: each with the relation to the node it was
// reached from, its path back to the file asked about, how many it pulls past it, and its
// freeze.
func tieItems(g *mapx.Graph, targets []string, ties []mapx.Dependent, relation map[mapx.EdgeType]string, pulls func(string) int) []impactItem {
	byID := map[string]mapx.Dependent{}
	for _, d := range ties {
		byID[d.ID] = d
	}
	items := []impactItem{}
	for _, d := range ties {
		var path []string
		cur := d.From
		for !slices.Contains(targets, cur) {
			path = append(path, cur)
			cur = byID[cur].From
		}
		path = append(path, cur)
		it := impactItem{ID: d.ID, Level: d.Level, Relation: relation[d.Type], Path: path}
		if d.TypesOnly {
			it.Relation += " (types only)"
		}
		if n := pulls(d.ID); n >= hubPull {
			it.Pulls = n
		}
		if f := g.FrozenAt(d.ID); f != "" {
			it.Frozen = g.DescribeFreeze(f)
		}
		items = append(items, it)
	}
	return items
}

// hubPull is how many direct dependents make an item worth flagging: a branch through it pulls
// that many more, as a navigator that imports every screen.
const hubPull = 10

// impactItem is one dependent as the report shows it.
type impactItem struct {
	ID       string   `json:"id"`
	Level    int      `json:"level"`
	Relation string   `json:"relation"`
	Path     []string `json:"path"`
	Pulls    int      `json:"pulls,omitempty"`
	Frozen   string   `json:"frozen,omitempty"`
}

// impactTest is a test to run again, with the layers of the suites that run it.
type impactTest struct {
	File   string   `json:"file"`
	Layers []string `json:"layers"`
}

// impactOut is the whole answer for the files asked about.
type impactOut struct {
	Origins    []string     `json:"origins"`
	Dependents []impactItem `json:"dependents"`
	Tests      []impactTest `json:"tests"`
	Validates  []string     `json:"validates"`
}

// relations say how an item depends on the one it was reached from, by edge type.
var relations = map[mapx.EdgeType]string{
	mapx.EdgeDependsOn: "depends on it", mapx.EdgeReadsFrom: "reads what it writes", mapx.EdgeCaptures: "captures it",
	mapx.EdgeRealizes: "realizes it", mapx.EdgeGatedBy: "is gated by it", mapx.EdgeNeeds: "needs it done",
	mapx.EdgeGoverns: "is governed by it", mapx.EdgeSpecifies: "is specified by it", mapx.EdgeCoveredBy: "covers it",
	mapx.EdgeTestedBy: "tests it", mapx.EdgeSeeds: "is seeded by it",
}

func impactReport(g *mapx.Graph, cfg *config.Config, targets []string, maxLevel int, stop func(string) bool) impactOut {
	deps := g.Dependents(targets, maxLevel, stop)
	reached := map[string]bool{}
	for _, t := range targets {
		reached[t] = true
	}
	for _, d := range deps {
		// A tie of types alone carries no change at run time: its tests are not run again.
		if !d.TypesOnly {
			reached[d.ID] = true
		}
	}
	rep := impactOut{Origins: targets, Tests: []impactTest{}, Validates: []string{}}
	rep.Dependents = tieItems(g, targets, deps, relations, g.DependentCount)
	seen := map[string]bool{}
	for _, t := range targets {
		for _, r := range g.Rulers(t) {
			if !seen[r] && !slices.Contains(targets, r) {
				seen[r] = true
				rep.Validates = append(rep.Validates, r)
			}
		}
	}
	sort.Strings(rep.Validates)
	for _, t := range g.TestsHolding(reached) {
		if slices.Contains(targets, t) || g.FrozenAt(t) != "" {
			continue
		}
		rep.Tests = append(rep.Tests, impactTest{File: t, Layers: testLayers(g.Node(t), cfg)})
	}
	return rep
}

// testLayers are the layers of the suites that run a test: those its last result names, or,
// before any, the suites whose paths cover it.
func testLayers(n *mapx.Node, cfg *config.Config) []string {
	seen := map[string]bool{}
	var out []string
	if n.Signal != nil {
		for l := range n.Signal.ByLayer {
			seen[l] = true
			out = append(out, l)
		}
	}
	if len(out) == 0 && cfg != nil {
		for _, s := range cfg.Tests {
			if s.Covers(n.ID) && !seen[s.Layer] {
				seen[s.Layer] = true
				out = append(out, s.Layer)
			}
		}
	}
	sort.Strings(out)
	if len(out) == 0 {
		out = []string{"unassigned"}
	}
	return out
}

// relTo converte o arg (que pode ser absoluto ou relativo ao cwd) para o caminho
// relativo à raiz do projeto — a forma como os nós são identificados no mapa.
// relTo converte o argumento do usuário num caminho relativo à RAIZ do projeto.
//
// O argumento pode vir de duas origens legítimas, e elas se resolvem de formas opostas:
// relativo ao diretório onde o comando foi chamado (como todo shell faz) ou relativo à
// raiz do projeto (que é a forma que os próprios prompts do Anchors imprimem, e a que os
// agentes copiam). Enquanto a raiz era sempre o CWD, as duas coincidiam. Desde que os
// comandos passaram a SUBIR até a raiz, deixaram de coincidir — e um `--for` copiado do
// prompt, executado de um subdiretório, virava `packages/backend/packages/backend/…`.
//
// A régua: se o caminho relativo à raiz EXISTE, é ele. É o desempate certo porque a origem
// do argumento não é declarada — só o disco sabe qual leitura corresponde a um arquivo.
//
// A saída é SEMPRE com barra normal, porque o que se devolve é o id de um nó, e o id é o
// mesmo em toda máquina: o mapa é versionado e trafega entre macOS, Linux e Windows. Sem
// isto, `filepath.Clean`/`Rel` devolvem "packages\backend\x.spec.md" no Windows, a busca
// no mapa (gravado com "/") não casa, e o arquivo aparece como "REGIDO mas fora do mapa" —
// pedindo um `map build` que reescreveria o mapa inteiro no dialeto da máquina. Converter
// só na fronteira do disco é o que mantém o mapa intercambiável.
func relTo(root, arg string) string {
	if !filepath.IsAbs(arg) {
		if _, err := os.Stat(filepath.Join(root, arg)); err == nil {
			return filepath.ToSlash(filepath.Clean(arg))
		}
	}
	abs, err := filepath.Abs(arg)
	if err != nil {
		return filepath.ToSlash(arg)
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return filepath.ToSlash(arg)
	}
	return filepath.ToSlash(rel)
}

func nodeExists(g *mapx.Graph, id string) bool {
	for _, n := range g.Nodes {
		if n.ID == id {
			return true
		}
	}
	return false
}

// printLevels prints the items level by level — the direct ones, then each further level with
// the path of each item —, and returns the frozen ones, which a freeze keeps apart (FRUNF-D09).
func printLevels(items []impactItem, what, none, propagate, verb string) []impactItem {
	levels := map[int][]impactItem{}
	var frozen []impactItem
	maxLevel := 0
	for _, it := range items {
		if it.Frozen != "" {
			frozen = append(frozen, it)
			continue
		}
		levels[it.Level] = append(levels[it.Level], it)
		if it.Level > maxLevel {
			maxLevel = it.Level
		}
	}
	fmt.Printf("↓ Direct (%d) — %s:\n", len(levels[1]), what)
	if len(levels[1]) == 0 {
		fmt.Printf("   %s\n", none)
	}
	for level := 1; level <= maxLevel; level++ {
		if level > 1 {
			name := "Indirect"
			if level > 2 {
				name = fmt.Sprintf("Cascade, level %d", level)
			}
			fmt.Printf("\n↓ %s (%d):\n", name, len(levels[level]))
		}
		for _, it := range levels[level] {
			line := fmt.Sprintf("   %s — %s", it.ID, it.Relation)
			if level > 1 {
				line += " · via " + strings.Join(it.Path, " ← ")
			}
			if it.Pulls > 0 {
				line += fmt.Sprintf(" (pulls %d more)", it.Pulls)
			}
			fmt.Println(line)
		}
	}
	if propagate == "" {
		fmt.Printf("   (--propagate indirect | cascade %s further, with the path of each item)\n", verb)
	} else if maxLevel < 2 {
		fmt.Printf("\n↓ Indirect (0):\n   (nothing further — every branch ends at the direct level, or at --stop-at)\n")
	}
	return frozen
}

func printImpact(rep impactOut, layer, propagate string) {
	fmt.Printf("impact of: %s\n\n", strings.Join(rep.Origins, ", "))

	frozen := printLevels(rep.Dependents, "they depend on it", "(nothing — no child depends on it)", propagate, "climbs")

	byLayer := map[string][]string{}
	var layers []string
	for _, t := range rep.Tests {
		for _, l := range t.Layers {
			if layer != "" && l != layer {
				continue
			}
			if _, ok := byLayer[l]; !ok {
				layers = append(layers, l)
			}
			byLayer[l] = append(byLayer[l], t.File)
		}
	}
	sort.Strings(layers)
	fmt.Printf("\n⟳ Tests to run again (%d) — reached, or holding a reached file in their evidence:\n", len(rep.Tests))
	if len(layers) == 0 {
		fmt.Println("   (none)")
	}
	for _, l := range layers {
		fmt.Printf("   %s (%d):\n", l, len(byLayer[l]))
		for _, f := range byLayer[l] {
			fmt.Printf("      %s\n", f)
		}
	}

	fmt.Printf("\n↑ Validates against (%d) — confront; a divergence becomes an issue:\n", len(rep.Validates))
	if len(rep.Validates) == 0 {
		fmt.Println("   (nothing — it is not governed by anyone)")
	}
	for _, n := range rep.Validates {
		fmt.Printf("   %s\n", n)
	}

	if len(frozen) > 0 {
		fmt.Printf("\n❄ Frozen (%d) — reached, and kept as they are by decision; not redone without an order:\n", len(frozen))
		for _, it := range frozen {
			fmt.Printf("   %s (%s)\n", it.ID, it.Frozen)
		}
	}
}
