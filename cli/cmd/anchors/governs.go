package main

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"

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
		Short: "Mostra quem cada guide rege (e quantos) — dimensiona a auditoria por guide",
		Long: `Lê o mapa e mostra a governança:

  anchors governs             o quadro: cada guide e quantos nós ele rege (direto)
  anchors governs <guide>     os arquivos que aquele guide rege

Rege "direto" = as arestas governs que saem do guide (não a onda transitiva; para
o impacto completo use 'anchors impact <guide>'). Guides que regem o MESMO conjunto
sinalizam redundância (candidatos a afinar por tag).`,
		Args: cobra.MaximumNArgs(1),
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

			// modo detalhe: um guide específico
			if len(args) == 1 {
				guide := relTo(absRoot, args[0])
				regidos := g.Governs(guide)
				if len(regidos) == 0 {
					fmt.Printf("%s não rege ninguém (sem regra governs, ou não é guide)\n", guide)
					return nil
				}
				fmt.Printf("%s rege %d arquivo(s):\n\n", guide, len(regidos))
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
				fmt.Println("nenhum guide rege nada (sem regras governs declaradas)")
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
			fmt.Printf("%d guide(s) com governança (rege direto):\n\n", len(rows))
			total := 0
			for _, r := range rows {
				fmt.Printf("  %4d  %s\n", r.n, r.guide)
				total += r.n
			}
			fmt.Printf("\ntotal de pares (guide, regido) = %d — o tamanho de uma auditoria completa por guide\n", total)
			fmt.Println("dica: guides com a MESMA contagem sobre o mesmo escopo são redundantes (afine por tag).")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
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
