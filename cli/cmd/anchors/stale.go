package main

import (
	"fmt"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/config"

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
		Short: "Lista as arestas stale — o que mudou e ainda não foi reconfrontado",
		Long: `Percorre o mapa e lista as arestas STALE: nunca validadas, ou com uma
ponta que avançou de rev desde o último carimbo do check. É a dívida de confronto —
rode 'anchors check' sobre os alvos para reconciliá-las e recarimbar.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("carregar mapa: %w (rode `anchors map build`)", err)
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
				fmt.Printf("%d evidência(s) de teste VENCIDA(s) — o teste passou contra código que já mudou:\n\n", len(vencidas))
				for _, ev := range vencidas {
					motivo := "o próprio arquivo de teste mudou"
					if len(ev.Culprit) > 0 {
						motivo = fmt.Sprintf("%d dependência(s) mudaram, ex.: %s", len(ev.Culprit), ev.Culprit[0])
						if ev.Own {
							motivo += " (e o próprio teste)"
						}
					}
					fmt.Printf("  %s\n      %s\n", ev.Test, motivo)
				}
				fmt.Printf("\n  Evidência vencida não é defeito do teste: é o placar que envelheceu. Rode-o de novo.\n\n")
			}

			stale := g.StaleEdges()
			total := len(g.Edges)
			if len(stale) == 0 {
				fmt.Printf("nenhuma aresta stale (%d arestas, todas validadas)\n", total)
				return nil
			}
			// separa "nunca validada" de "avançou de rev" — dívidas de naturezas
			// diferentes: a 1ª é cobertura que nunca rodou; a 2ª é drift real.
			var never, drifted int
			fmt.Printf("%d de %d aresta(s) stale:\n\n", len(stale), total)
			for _, e := range stale {
				reason := "avançou de rev"
				if e.Stamp == nil {
					reason = "nunca validada"
					never++
				} else {
					drifted++
				}
				fmt.Printf("  %s ──%s──▶ %s  (%s)\n", e.From, e.Type, e.To, reason)
			}
			fmt.Printf("\n%d nunca validada(s), %d com drift de rev\n", never, drifted)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	return cmd
}
