package main

import (
	"fmt"
	"path/filepath"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

func newMapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "map",
		Short: "Opera o mapa de dependências (anchors.graph.yaml)",
	}
	cmd.AddCommand(newMapBuildCmd())
	cmd.AddCommand(newMapShowCmd())
	return cmd
}

func newMapShowCmd() *cobra.Command {
	var root, mapPath string
	var orphans, stats, worklist, onlyPending bool
	cmd := &cobra.Command{
		Use:   "show [arquivo]",
		Short: "Consulta o mapa: vizinhança de um nó, órfãos, estatísticas",
		Long: `Inspeciona o mapa de dependências:
  anchors map show <arquivo>   — a vizinhança do nó (entra ↑ / sai ↓)
  anchors map show --orphans   — nós sem nenhuma aresta (ilhas)
  anchors map show --stats     — resumo (nós por kind, arestas por tipo)`,
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

			switch {
			case worklist:
				return printWorklist(g, absRoot, onlyPending)
			case stats:
				printStats(g.Statistics())
			case orphans:
				printOrphans(g.Orphans())
			case len(args) == 1:
				target := relTo(absRoot, args[0])
				if !nodeExists(g, target) {
					return fmt.Errorf("arquivo %q não está no mapa", target)
				}
				printNeighborhood(g.Neighbors(target))
			default:
				return fmt.Errorf("informe <arquivo>, ou --orphans, ou --stats, ou --worklist")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto")
	cmd.Flags().StringVar(&mapPath, "map", "", "caminho do mapa")
	cmd.Flags().BoolVar(&orphans, "orphans", false, "lista nós sem arestas (ilhas)")
	cmd.Flags().BoolVar(&stats, "stats", false, "resumo do grafo")
	cmd.Flags().BoolVar(&worklist, "worklist", false, "lista os nós em ordem TOPOLÓGICA (pais/réguas antes dos regidos) — para correção em lote sem retrabalho")
	cmd.Flags().BoolVar(&onlyPending, "pending", false, "com --worklist: só os arquivos que TÊM pendência (gates que reprovam)")
	return cmd
}

func printNeighborhood(nb mapx.Neighborhood) {
	fmt.Printf("%s\n\n", nb.Node)
	fmt.Printf("↑ regido por / valida contra (%d):\n", len(nb.In))
	if len(nb.In) == 0 {
		fmt.Println("   (nada — é topo/raiz)")
	}
	for _, e := range nb.In {
		fmt.Printf("   %s  ←%s←  %s\n", nb.Node, e.Type, e.From)
	}
	fmt.Printf("\n↓ propaga para (%d):\n", len(nb.Out))
	if len(nb.Out) == 0 {
		fmt.Println("   (nada — é folha)")
	}
	for _, e := range nb.Out {
		fmt.Printf("   %s  —%s→  %s\n", nb.Node, e.Type, e.To)
	}
}

func printOrphans(orphs []mapx.Node) {
	fmt.Printf("órfãos — nós sem nenhuma aresta (%d):\n", len(orphs))
	if len(orphs) == 0 {
		fmt.Println("   (nenhum — todo nó está conectado)")
	}
	for _, n := range orphs {
		fmt.Printf("   [%s] %s\n", n.Kind, n.ID)
	}
}

func printStats(s mapx.Stats) {
	fmt.Printf("mapa: %d nós, %d arestas\n\n", s.Nodes, s.Edges)
	fmt.Println("nós por kind:")
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindFeature, mapx.KindTest, mapx.KindCode, mapx.KindDoc, mapx.KindGuide, mapx.KindPlan} {
		if n := s.NodesByKind[k]; n > 0 {
			fmt.Printf("   %-9s %d\n", k, n)
		}
	}
	fmt.Println("\narestas por tipo:")
	for _, ty := range []mapx.EdgeType{mapx.EdgeGoverns, mapx.EdgeSpecifies, mapx.EdgeCoveredBy, mapx.EdgeTestedBy, mapx.EdgeReferences} {
		if n := s.EdgesByType[ty]; n > 0 {
			fmt.Printf("   %-11s %d\n", ty, n)
		}
	}
}

func newMapBuildCmd() *cobra.Command {
	var root, out string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Constrói o mapa de dependências a partir do projeto",
		Long: `Percorre o projeto lendo TEXTO (nunca parseando código) e infere as
arestas do mapa por co-location (nomes de arquivo) e por código de cenário
(a identidade estável que atravessa spec→feature→teste). Grava anchors.graph.yaml.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRaiz(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("carregar %s: %w (rode `anchors init` para criar)", config.DefaultFile, err)
			}
			files, err := scan.Walk(absRoot, cfg)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			// carimbo de alteração (updated_at) de cada nó, do git — uma só chamada
			// batch (não 1 por arquivo). Vazio se não for repo git.
			updatedAt := gitmeta.AllCommitDates(absRoot)
			g := mapx.Build(files, cfg, updatedAt)

			outPath := out
			if outPath == "" {
				outPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			// O rebuild não pode apagar a memória de validação: o `anchors work` manda
			// rodar `map build` antes de todo `check`, então sem isto cada etapa zerava o
			// carimbo da anterior e o `anchors stale` acusava o repositório inteiro como
			// "nunca validado".
			if anterior, err := mapx.Load(outPath); err == nil {
				mapx.PreservarCarimbos(g, anterior)
			}
			if err := mapx.Save(g, outPath); err != nil {
				return fmt.Errorf("save: %w", err)
			}

			fmt.Printf("mapa construído: %d nós, %d arestas\n", len(g.Nodes), len(g.Edges))
			fmt.Printf("  escrito em %s\n", outPath)
			printEdgeSummary(g)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "raiz do projeto a escanear")
	cmd.Flags().StringVar(&out, "out", "", "caminho de saída (default: <root>/anchors.graph.yaml)")
	return cmd
}

// printWorklist lista os nós em ordem TOPOLÓGICA — pais/réguas primeiro, regidos
// depois — para a correção em lote ir de cima pra baixo sem retrabalho (consertar um
// filho e o pai depois invalidá-lo). Com onlyPending, roda os gates e mostra só os que
// têm pendência (o que o agente precisa de fato tocar).
func printWorklist(g *mapx.Graph, root string, onlyPending bool) error {
	ordered := g.TopoOrder()

	pending := map[string]bool{}
	if onlyPending {
		cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
		if err != nil {
			return fmt.Errorf("carregar config: %w", err)
		}
		for _, r := range gate.RunWithConfig(cfg.Gates, g.Nodes, root, g, cfg) {
			if r.Verdict == gate.Fail {
				pending[r.Target] = true
			}
		}
	}

	fmt.Println("worklist — ordem topológica (pais/réguas primeiro; processe de cima pra baixo):")
	fmt.Println("# consertar nesta ordem evita retrabalho — o pai já está estável quando o filho é tocado.")
	n := 0
	for _, node := range ordered {
		if onlyPending && !pending[node.ID] {
			continue
		}
		n++
		fmt.Printf("  %s\n", node.ID)
	}
	if onlyPending {
		fmt.Printf("\n%d arquivo(s) com pendência, em ordem.\n", n)
	} else {
		fmt.Printf("\n%d nó(s), em ordem topológica.\n", n)
	}
	return nil
}

func printEdgeSummary(g *mapx.Graph) {
	byType := map[mapx.EdgeType]int{}
	for _, e := range g.Edges {
		byType[e.Type]++
	}
	if len(byType) == 0 {
		return
	}
	fmt.Println("  arestas por tipo:")
	for _, t := range []mapx.EdgeType{mapx.EdgeGoverns, mapx.EdgeSpecifies, mapx.EdgeCoveredBy, mapx.EdgeTestedBy, mapx.EdgeReferences} {
		if n := byType[t]; n > 0 {
			fmt.Printf("    %-11s %d\n", t, n)
		}
	}
}
