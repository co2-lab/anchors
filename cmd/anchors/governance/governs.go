package governance

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// `anchors governs` responde "quem cada guide rege, e quantos" — a partir do mapa.
// Sem argumento: o quadro de todos os guides (dimensiona a auditoria por guide, e
// expõe redundância — guides que regem o mesmo conjunto). Com um guide: os arquivos
// que ele rege diretamente. É a base para fatiar uma auditoria de julgamento.
func newGovernsCmd() *cobra.Command {
	var root, mapPath string
	cmd := &cobra.Command{
		Use:   "governs [<guide>]",
		Short: "Show who each guide governs (and how many) — sizes the audit per guide",
		Long: `Reads the map and shows the governance:

  anchors governs             the board: each guide and how many nodes it governs (direct)
  anchors governs <guide>     the files that guide governs

Governs "direct" = the governs edges leaving the guide (not the transitive wave; for
the full impact use 'anchors impact <guide>'). Guides that govern the SAME set
signal redundancy (candidates to narrow by tag).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}

			// modo detalhe: um guide específico
			if len(args) == 1 {
				guide := common.RelTo(absRoot, args[0])
				regidos := g.Governs(guide)
				if len(regidos) == 0 {
					fmt.Printf("%s governs nobody (no governs rule, or it is not a guide)\n", guide)
					return nil
				}
				fmt.Printf("%s governs %d file(s):\n\n", guide, len(regidos))
				byKind := map[mapx.Kind][]string{}
				kindOf := nodeKindIndex(g)
				for _, id := range regidos {
					byKind[kindOf[id]] = append(byKind[kindOf[id]], id)
				}
				for _, k := range sortedKinds(byKind) {
					fmt.Printf("  [%s] %d\n", k, len(byKind[k]))
					for _, id := range byKind[k] {
						fmt.Printf("    %s\n", id)
					}
				}
				return nil
			}

			// modo quadro: todos os guides
			summary := g.GovernanceSummary()
			if len(summary) == 0 {
				fmt.Println("no guide governs anything (no governs rules declared)")
				return nil
			}
			type row struct {
				guide string
				n     int
			}
			var rows []row
			for guide, n := range summary {
				rows = append(rows, row{guide, n})
			}
			sort.Slice(rows, func(i, j int) bool {
				if rows[i].n != rows[j].n {
					return rows[i].n > rows[j].n
				}
				return rows[i].guide < rows[j].guide
			})
			fmt.Printf("%d guide(s) with governance (direct):\n\n", len(rows))
			total := 0
			for _, r := range rows {
				fmt.Printf("  %4d  %s\n", r.n, r.guide)
				total += r.n
			}
			fmt.Printf("\ntotal pairs (guide, governed) = %d — the size of a full audit per guide\n", total)
			fmt.Println("tip: guides with the SAME count over the same scope are redundant (narrow by tag).")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	return cmd
}

func nodeKindIndex(g *mapx.Graph) map[string]mapx.Kind {
	m := map[string]mapx.Kind{}
	for _, n := range g.Nodes {
		m[n.ID] = n.Kind
	}
	return m
}

func sortedKinds(m map[mapx.Kind][]string) []mapx.Kind {
	var ks []mapx.Kind
	for k := range m {
		ks = append(ks, k)
	}
	slices.Sort(ks)
	return ks
}
