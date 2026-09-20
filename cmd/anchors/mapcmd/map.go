package mapcmd

import (
	"fmt"
	"github.com/co2-lab/anchors/internal/i18n"
	"path/filepath"
	"sort"
	"strings"

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
		Short: "Operate the dependency map (anchors.graph.yaml)",
	}
	cmd.AddCommand(newMapBuildCmd())
	cmd.AddCommand(newMapShowCmd())
	// O merge driver: o git chama isto em vez de mesclar o mapa como texto (#12).
	cmd.AddCommand(newMapMergeCmd())
	return cmd
}

func newMapShowCmd() *cobra.Command {
	var root, mapPath string
	var orphans, stats, worklist, onlyPending bool
	cmd := &cobra.Command{
		Use:   "show [file]",
		Short: "Query the map: a node's neighbourhood, orphans, statistics",
		Long: `Inspects the dependency map:
  anchors map show <file>      — the node's neighbourhood (in ↑ / out ↓)
  anchors map show --orphans   — nodes with no edge at all (islands)
  anchors map show --stats     — summary (nodes by kind, edges by type)`,
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
					return fmt.Errorf("file %q is not in the map", target)
				}
				printNeighborhood(g.Neighbors(target))
			default:
				return fmt.Errorf("provide <file>, or --orphans, or --stats, or --worklist")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().BoolVar(&orphans, "orphans", false, "lists nodes without edges (islands)")
	cmd.Flags().BoolVar(&stats, "stats", false, "graph summary")
	cmd.Flags().BoolVar(&worklist, "worklist", false, "lists the nodes in TOPOLOGICAL order (parents/rulers before the governed) — for batch fixing without rework")
	cmd.Flags().BoolVar(&onlyPending, "pending", false, "with --worklist: only the files that HAVE pending items (failing gates)")
	return cmd
}

func printNeighborhood(nb mapx.Neighborhood) {
	fmt.Printf("%s\n\n", nb.Node)
	fmt.Println(i18n.T("map.show.governed_by", len(nb.In)))
	if len(nb.In) == 0 {
		fmt.Println(i18n.T("map.show.top_root"))
	}
	for _, e := range nb.In {
		fmt.Printf("   %s  ←%s←  %s\n", nb.Node, e.Type, e.From)
	}
	fmt.Println(i18n.T("map.show.propagates_to", len(nb.Out)))
	if len(nb.Out) == 0 {
		fmt.Println(i18n.T("map.show.leaf"))
	}
	for _, e := range nb.Out {
		fmt.Printf("   %s  —%s→  %s\n", nb.Node, e.Type, e.To)
	}
}

func printOrphans(orphs []mapx.Node) {
	fmt.Println(i18n.T("map.show.orphans_header", len(orphs)))
	if len(orphs) == 0 {
		fmt.Println(i18n.T("map.show.none"))
	}
	for _, n := range orphs {
		fmt.Printf("   [%s] %s\n", n.Kind, n.ID)
	}
}

func printStats(s mapx.Stats) {
	fmt.Println(i18n.T("map.status.nodes_edges", s.Nodes, s.Edges))
	fmt.Println()
	fmt.Println(i18n.T("map.nodes_by_kind"))
	for _, k := range []mapx.Kind{mapx.KindSpec, mapx.KindFeature, mapx.KindTest, mapx.KindCode, mapx.KindDoc, mapx.KindGuide, mapx.KindPlan} {
		if n := s.NodesByKind[k]; n > 0 {
			fmt.Printf("   %-9s %d\n", k, n)
		}
	}
	fmt.Println("\n" + strings.TrimLeft(i18n.T("map.edges_by_type"), " "))
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
		Short: "Build the dependency map from the project",
		Long: `Walks the project reading TEXT (never parsing code) and infers the
map's edges by co-location (file names) and by scenario code
(the stable identity that crosses spec→feature→test). Writes anchors.graph.yaml.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load %s: %w (run `anchors init` to create it)", config.DefaultFile, err)
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
			var perdidos string
			if anterior, err := mapx.Load(outPath); err == nil {
				mapx.PreserveStamps(g, anterior)
				// PERDA DE CARIMBO é silenciosa, e o mapa continua VÁLIDO.
				//
				// O `PreserveStamps` preserva o que está no arquivo anterior — e num merge
				// o arquivo anterior é o do OUTRO lado. Medido no blue-eyes
				// (co2-lab/anchors#12): resolvi um conflito com `checkout --theirs` +
				// `map build`, e o carimbo de `review` de uma spec desapareceu. Descobri
				// por acaso, contando: 18 onde eu esperava 19.
				//
				// O custo não é o carimbo: é o LAUDO. Ele vive no `--reason` do comando,
				// não no arquivo, e refazer uma revisão adversarial é caro.
				//
				// Remover um nó legitimamente remove os carimbos dele, então isto é AVISO
				// e não reprovação. O que ele faz é transformar perda silenciosa em perda
				// visível, que já é a maior parte do dano.
				perdidos = stampLossWarning(anterior, g)
			}
			if err := mapx.Save(g, outPath); err != nil {
				return fmt.Errorf("save: %w", err)
			}

			fmt.Println(i18n.T("map.built", len(g.Nodes), len(g.Edges)))
			fmt.Println(i18n.T("map.written_to", outPath))
			printEdgeSummary(g)
			printLayerAmbiguities(scan.Ambiguities(files, cfg))
			if perdidos != "" {
				fmt.Print(perdidos)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root to scan")
	cmd.Flags().StringVar(&out, "out", "", "output path (default: <root>/anchors.graph.yaml)")
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
			return fmt.Errorf("load config: %w", err)
		}
		for _, r := range gate.RunWithConfig(cfg.Gates, g.Nodes, root, g, cfg) {
			if r.Verdict == gate.Fail {
				pending[r.Target] = true
			}
		}
	}

	fmt.Println("worklist — topological order (parents/rulers first; process from top to bottom):")
	fmt.Println("# fixing in this order avoids rework — the parent is already stable when the child is touched.")
	n := 0
	for _, node := range ordered {
		if onlyPending && !pending[node.ID] {
			continue
		}
		n++
		fmt.Printf("  %s\n", node.ID)
	}
	if onlyPending {
		fmt.Printf("\n%d file(s) with pending items, in order.\n", n)
	} else {
		fmt.Printf("\n%d node(s), in topological order.\n", n)
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
	fmt.Println(i18n.T("map.edges_by_type"))
	// The ORDER is the triad's (the path the reader walks), and everything else follows
	// alphabetically.
	//
	// It used to be a FIXED list of five types, and whatever was not on it stayed
	// invisible: `depends-on`, `seeds`, `needs` and `realizes` existed in the map and
	// never showed up in the summary. The symptom misled — the `realizes` edges were
	// built, the summary did not list them, and the first conclusion was that the build
	// had not created them.
	known := []mapx.EdgeType{mapx.EdgeGoverns, mapx.EdgeSpecifies, mapx.EdgeCoveredBy, mapx.EdgeTestedBy, mapx.EdgeReferences}
	seen := map[mapx.EdgeType]bool{}
	for _, t := range known {
		seen[t] = true
		if n := byType[t]; n > 0 {
			fmt.Printf("    %-11s %d\n", t, n)
		}
	}
	var rest []string
	for t := range byType {
		if !seen[t] {
			rest = append(rest, string(t))
		}
	}
	sort.Strings(rest)
	for _, t := range rest {
		fmt.Printf("    %-11s %d\n", t, byType[mapx.EdgeType(t)])
	}
}

// printLayerAmbiguities denuncia onde a HEURÍSTICA escolheu a camada.
//
// O comentário do campo `priority` promete: "Declare quando a heurística errar; o `check`
// avisa onde ela decidiu sozinha." Ele não avisava — a `scan.Ambiguities` existia,
// completa e testada, e NINGUÉM a chamava.
//
// O custo do silêncio, medido no blue-eyes (blue-eyes#100): `**/*.test.*` e
// `packages/shared/**/*.ts` casavam o mesmo `AreaStatus.test.ts`, o desempate por
// comprimento escolheu `shared`, e o projeto ficou com ZERO nós `kind: test` tendo 70
// testes verdes. Em cascata, QUATRO gates ficaram cegos — incluindo o `test-traceable`,
// que é bloqueante.
//
// A única pista era a linha "gate declarado sem nada para medir", que aparece por vários
// motivos legítimos (o artefato não existe ainda, a camada não foi declarada). Só se
// descobriu contando os nós do mapa à mão.
//
// Sai no `map build` e não no `check` porque é aqui que a classificação acontece: avisar
// no lugar onde a decisão é tomada é o que liga a causa ao efeito. E o `work` manda rodar
// `map build` antes de todo `check`, então quem segue o fluxo vê.
//
// É AVISO, não reprovação: na maioria das vezes o comprimento acerta, e barrar por
// ambiguidade reprovaria projetos que estão certos. O que não pode é decidir em silêncio.
func printLayerAmbiguities(amb []scan.LayerAmbiguity) {
	if len(amb) == 0 {
		return
	}
	// Agrupa por PAR de camadas: o interessante é "test perde para shared", não a lista
	// de arquivos. Num monorepo o mesmo par produz centenas de linhas idênticas, e
	// despejá-las esconde justamente o padrão que precisa ser visto.
	type par struct{ vencedora, perdedoras string }
	ordem := []par{}
	quantos := map[par]int{}
	exemplo := map[par]string{}
	for _, a := range amb {
		k := par{a.Vencedora, strings.Join(a.Perdedoras, ", ")}
		if quantos[k] == 0 {
			ordem = append(ordem, k)
			exemplo[k] = a.Arquivo
		}
		quantos[k]++
	}
	fmt.Printf("\n⚠ %d file(s) had the layer decided by HEURISTIC — two patterns "+
		"matched and neither declared `priority`:\n", len(amb))
	for _, k := range ordem {
		fmt.Printf("    %s beat %s  (%d file(s), e.g.: %s)\n",
			k.vencedora, k.perdedoras, quantos[k], exemplo[k])
	}
	fmt.Println("  The rule is the pattern LENGTH, which measures verbosity and not precision.")
	fmt.Println("  If the choice is wrong, declare `priority: N` on the layer that should win —")
	fmt.Println("  the wrong layer takes the file out of reach of EVERY gate that measures the right one.")
}

// stampLossWarning compara os carimbos do mapa ANTERIOR com os do novo.
//
// Conta por GATE, e não o total: "o mapa perdeu 1 carimbo" não diz o que refazer, e
// "perdeu 1 de `review`" diz — o laudo de uma revisão adversarial é a coisa mais cara que
// um carimbo representa.
//
// Devolve string vazia quando não houve perda, para o comando não imprimir nada no caso
// normal. Um aviso que sempre aparece deixa de ser lido.
func stampLossWarning(antigo, novo *mapx.Graph) string {
	conta := func(g *mapx.Graph) map[string]int {
		out := map[string]int{}
		for _, e := range g.Edges {
			for _, j := range e.Julgamentos {
				out[j.Gate]++
			}
		}
		return out
	}
	antes, depois := conta(antigo), conta(novo)
	type perda struct {
		gate  string
		antes int
		agora int
	}
	var perdas []perda
	for gate, n := range antes {
		if depois[gate] < n {
			perdas = append(perdas, perda{gate, n, depois[gate]})
		}
	}
	if len(perdas) == 0 {
		return ""
	}
	sort.Slice(perdas, func(i, j int) bool { return perdas[i].gate < perdas[j].gate })

	var b strings.Builder
	b.WriteString("\n⚠ the map LOST judgment stamp(s):\n")
	for _, p := range perdas {
		fmt.Fprintf(&b, "    %-28s %d → %d\n", p.gate, p.antes, p.agora)
	}
	b.WriteString("  Removing a node also removes its stamps — so this may be normal.\n")
	b.WriteString("  What is NOT normal: losing a stamp after resolving a conflict in the map.\n")
	b.WriteString("  The report lives in `anchors judge`'s `--reason`, not in the file: if the stamp\n")
	b.WriteString("  is gone, the judgment goes back to the queue and the evidence has to be redone.\n")
	return b.String()
}
