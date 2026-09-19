package quality

import (
	"fmt"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// `anchors stale` — agora que o check carimba as arestas, a distinção entre "stale"
// e "validado" tem significado. Uma aresta é stale quando nunca foi confrontada ou
// quando uma ponta avançou de rev desde o último carimbo (PROPAGATION §3). Este
// comando é a leitura desse estado: o que precisa ser reconfrontado.
func newStaleCmd() *cobra.Command {
	var root, mapPath string
	cmd := &cobra.Command{
		Use:   "stale",
		Short: "List stale edges — what changed and has not been reconfronted",
		Long: `Walks the map and lists the STALE edges: never validated, or with an
endpoint that advanced a rev since the check's last stamp. It is the confrontation debt —
run 'anchors check' over the targets to reconcile them and re-stamp.`,
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
			// EVIDÊNCIA primeiro: é a dívida mais acionável das duas. Uma aresta stale diz
			// "reconfronte o texto"; uma evidência vencida diz "o placar de teste que você
			// tem não vale mais" — e esta some silenciosamente, porque o teste continua
			// verde no relatório antigo.
			var vencidas []*mapx.EvidenceStale
			for _, n := range g.Nodes {
				if ev := g.EvidenceStaleFor(n.ID); ev != nil {
					vencidas = append(vencidas, ev)
				}
			}
			if len(vencidas) > 0 {
				fmt.Println(i18n.T("stale.expired_test_evidence", len(vencidas)) + "\n")
				for _, ev := range vencidas {
					motivo := i18n.T("stale.test_file_changed")
					if len(ev.Culprit) > 0 {
						motivo = i18n.T("stale.dependencies_changed", len(ev.Culprit), ev.Culprit[0])
						if ev.Own {
							motivo += i18n.T("stale.and_test_itself")
						}
					}
					fmt.Printf("  %s\n      %s\n", ev.Test, motivo)
				}
				fmt.Printf("\n  %s\n\n", i18n.T("stale.evidence_note"))
			}

			stale := g.StaleEdges()
			total := len(g.Edges)
			if len(stale) == 0 {
				fmt.Println(i18n.T("stale.no_stale_edges", total))
				return nil
			}
			// separa "nunca validada" de "avançou de rev" — dívidas de naturezas
			// diferentes: a 1ª é cobertura que nunca rodou; a 2ª é drift real.
			var never, drifted int
			fmt.Println(i18n.T("stale.edges_header", len(stale), total) + "\n")
			for _, e := range stale {
				reason := i18n.T("stale.rev_drifted")
				if e.Stamp == nil {
					reason = i18n.T("stale.never_validated")
					never++
				} else {
					drifted++
				}
				fmt.Printf("  %s ──%s──▶ %s  (%s)\n", e.From, e.Type, e.To, reason)
			}
			fmt.Printf("\n%s\n", i18n.T("stale.summary", never, drifted))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	return cmd
}
