package mapcmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/flowx"
	"github.com/co2-lab/anchors/internal/mapx"
)

func newFlowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flow",
		Short: "Operate the work flows (the puzzle of actions and results)",
		Long: `A FLOW is work driven by shape rather than by memory.

An ACTION (flows/actions/*.action.md) declares what it does and which RESULTS it
offers. A FLOW (flows/*.flow.md) fits actions together and says where each result
goes. What used to be a rule to remember ("NEVER close with a blocking gate red")
becomes the absence of a fitting: the BARRED result has no link to ` + "`done`" + `.`,
	}
	cmd.AddCommand(newFlowBuildCmd(), newFlowShowCmd(), newFlowNextCmd())
	return cmd
}

func newFlowBuildCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Scan flows/ and write the flow graph into the map",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			fg, err := flowx.Build(absRoot)
			if err != nil {
				return err
			}
			if fg == nil {
				fmt.Println("no flow declared — create flows/<name>.flow.md")
				return nil
			}
			path := filepath.Join(absRoot, mapx.DefaultPath)
			g, err := mapx.Load(path)
			if err != nil {
				// The map does not exist yet, and the flow is no reason to create it: a
				// map with no nodes would have the relational gates confront a void.
				return fmt.Errorf("build the map first (`anchors map build`): %w", err)
			}
			g.Flow = fg
			if err := mapx.Save(g, path); err != nil {
				return err
			}
			steps, results := 0, 0
			for _, s := range fg.States {
				if flowx.IsResult(s.Code) {
					results++
				} else {
					steps++
				}
			}
			fmt.Printf("flow built: %d step(s), %d result(s), %d link(s)\n",
				steps, results, len(fg.Transitions))
			// The FINDING the puzzle shape makes possible: an action declares everything
			// it can answer, and a result no flow routes is a hole — whoever gets it
			// improvises, which is the whole problem.
			if un := flowx.Unhandled(fg); len(un) > 0 {
				fmt.Printf("\n⚠ %d result(s) no flow handles:\n", len(un))
				for _, s := range un {
					fmt.Printf("    %-12s %s\n", s.Code, s.Title)
				}
				fmt.Println("  Not a failure: a project may legitimately not cover every branch.")
				fmt.Println("  What it must not do is not know.")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}

func newFlowShowCmd() *cobra.Command {
	var root string
	var asMermaid bool
	var direction string
	cmd := &cobra.Command{
		Use:   "show [flow]",
		Short: "Draw the flow — assembled at RUNTIME from the graph",
		Long: `Draws the flow from the graph, never from a stored diagram.

A versioned diagram ages against the flow it describes, and the repository already
measured what that costs (see the docs-fresh gate). Here the drawing is derived,
so it cannot diverge.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fg, err := loadFlowGraph(root)
			if err != nil {
				return err
			}
			names := flowx.Flows(fg)
			var wanted []string
			if len(args) > 0 {
				for _, n := range names {
					if strings.Contains(n, args[0]) {
						wanted = append(wanted, n)
					}
				}
				if len(wanted) == 0 {
					return fmt.Errorf("no flow matching %q — available: %s", args[0], strings.Join(names, ", "))
				}
			} else {
				wanted = names
			}
			for i, n := range wanted {
				if i > 0 {
					fmt.Println()
				}
				if asMermaid {
					drawMermaid(fg, n, direction)
					continue
				}
				drawFlow(fg, n)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&asMermaid, "mermaid", false,
		"emit the diagram as mermaid (boxes and arrows, for rendering)")
	cmd.Flags().StringVar(&direction, "direction", "TD",
		"diagram direction: TD (top-down, tree-like) or LR (left to right)")
	return cmd
}

func newFlowNextCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "next <step>",
		Short: "The VALID exits of a step — the command that drives",
		Long: `Answers "from here, where can I go", instead of listing everything and
asking someone to choose. It is what separates having a flow from having one more
document.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fg, err := loadFlowGraph(root)
			if err != nil {
				return err
			}
			code := strings.ToUpper(args[0])
			st, ok := flowx.StateByCode(fg, code)
			if !ok {
				return fmt.Errorf("step %q is not in the flow graph", code)
			}
			fmt.Printf("%s — %s\n", st.Code, st.Title)
			if st.Fits != "" {
				if nome, ok := flowx.ActionTitle(fg, st.Fits); ok {
					fmt.Printf("  fits: %s (anchors %s)\n", st.Fits, nome)
				} else {
					// A piece declared and never written: the step points at an action that
					// does not exist, and whoever arrives has nothing to run.
					fmt.Printf("  fits: %s  ⚠ no action file declares it\n", st.Fits)
				}
			}
			exits := flowx.Next(fg, code)
			if len(exits) == 0 {
				if st.Terminal {
					fmt.Println("\n▣ terminal — the work ends here")
					return nil
				}
				// A step with no exit and no `@terminal` is the defect the flow exists to
				// make visible: whoever arrives has nowhere to go, and improvises.
				fmt.Println("\n⚠ no exit, and it does not declare @terminal — whoever arrives here is stuck")
				return nil
			}
			fmt.Printf("\nvalid exits (%d):\n", len(exits))
			for _, t := range exits {
				label := t.On
				if label == "" {
					label = "—"
				}
				dest := t.To
				if d, ok := flowx.StateByCode(fg, t.To); ok {
					dest = fmt.Sprintf("%s (%s)", d.Code, d.Title)
				}
				fmt.Printf("  %-12s → %s\n", label, dest)
				if t.When != "" {
					fmt.Printf("  %-12s   %s\n", "", t.When)
				}
				// A SUGESTÃO vem do RESULTADO, e não da transição: ela é o que a ação
				// recomenda a quem recebe aquela resposta, independentemente de para onde
				// este fluxo em particular a encaminha.
				//
				// Sem ela, quem chega a um resultado sabe para onde ir e não o que fazer —
				// e a recomendação continuaria só na prosa da mensagem do gate, que é de
				// onde ela precisava sair.
				if r, ok := flowx.StateByCode(fg, t.On); ok && r.Suggests != "" {
					fmt.Printf("  %-12s   ↳ %s\n", "", r.Suggests)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}

func loadFlowGraph(root string) (*mapx.FlowGraph, error) {
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return nil, err
	}
	g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
	if err != nil {
		return nil, err
	}
	if g.Flow == nil || len(g.Flow.States) == 0 {
		return nil, fmt.Errorf("no flow in the map — run `anchors flow build`")
	}
	return g.Flow, nil
}

// drawFlow draws ONE flow: the steps in file order, each with the piece it fits and the
// exits it offers.
//
// The order is the file's and not topological: whoever wrote the flow wrote it in the
// sequence in which it happens, and reordering by graph would scramble the thread of
// reading — the drawing serves whoever reads, and a graph ordered by traversal is another
// artifact.
func drawFlow(fg *mapx.FlowGraph, flow string) {
	steps := flowx.StatesOf(fg, flow)
	fmt.Printf("%s\n", flow)
	fmt.Println(strings.Repeat("─", len(flow)))

	width := 0
	for _, s := range steps {
		if flowx.IsResult(s.Code) {
			continue
		}
		if n := len(s.Code); n > width {
			width = n
		}
	}

	for _, s := range steps {
		if flowx.IsResult(s.Code) {
			continue
		}
		mark := "●"
		if s.Terminal {
			mark = "▣"
		}
		row := fmt.Sprintf("%s %-*s  %s", mark, width, s.Code, s.Title)
		if s.Fits != "" {
			row += fmt.Sprintf("   [%s]", s.Fits)
		}
		fmt.Println(row)

		exits := flowx.Next(fg, s.Code)
		sort.SliceStable(exits, func(i, j int) bool { return exits[i].On < exits[j].On })
		for i, t := range exits {
			branch := "├─"
			if i == len(exits)-1 {
				branch = "└─"
			}
			label := t.On
			if label == "" {
				label = "sempre"
			}
			fmt.Printf("  %s %-12s → %s", branch, label, t.To)
			if t.When != "" {
				fmt.Printf("  %s", t.When)
			}
			fmt.Println()
		}
	}
}

// drawMermaid emits the flow as a mermaid diagram — boxes and arrows, for rendering.
//
// DERIVED at runtime and never stored: a versioned diagram ages against the flow it
// describes, and this repository already measured what that costs (the `docs-fresh` gate
// exists because of it). Here the drawing cannot diverge, because it is recomputed.
func drawMermaid(fg *mapx.FlowGraph, flow, dir string) {
	steps := flowx.StatesOf(fg, flow)
	dir = strings.ToUpper(strings.TrimSpace(dir))
	if dir != "LR" && dir != "RL" && dir != "BT" {
		dir = "TD"
	}
	fmt.Printf("flowchart %s\n", dir)

	for _, s := range steps {
		if flowx.IsResult(s.Code) {
			continue
		}
		id := mermaidID(s.Code)
		label := escapeMermaid(s.Title)
		if s.Fits != "" {
			// The PIECE in the label: whoever reads the diagram needs to know which command
			// runs there, and hunting for it in another document would break what the
			// drawing is good for.
			//
			// The break is `\n` inside the quotes — mermaid's NATIVE form. `<br/>` looks
			// like the obvious path and only works with `htmlLabels` on; without it the
			// renderer prints the whole label on one line, and the box came out as
			// "pull the next taskACNXT" — the tag swallowed and both phrases glued.
			label += "\\n" + escapeMermaid(s.Fits)
		}
		switch {
		case s.Terminal:
			// A terminal gets its OWN shape (stadium): the end is recognised by the
			// outline before the text is read — which is what the eye does first in a
			// flowchart.
			fmt.Printf("  %s([\"%s\"])\n", id, label)
		default:
			fmt.Printf("  %s[\"%s\"]\n", id, label)
		}
	}

	for _, s := range steps {
		if flowx.IsResult(s.Code) {
			continue
		}
		exits := flowx.Next(fg, s.Code)
		sort.SliceStable(exits, func(i, j int) bool { return exits[i].On < exits[j].On })
		for _, t := range exits {
			label := exitLabel(fg, t)
			fmt.Printf("  %s -->|%s| %s\n", mermaidID(s.Code), escapeMermaid(label), mermaidID(t.To))
		}
	}

	// The ENTRY gets its own highlight, like the dark "START HERE" of a flowchart: whoever
	// opens the drawing needs to find the beginning before reading any box. Without it, in
	// a cyclic flow (the worker has three loops) there is no visual clue where one enters.
	if e, ok := flowx.Entry(fg, flow); ok {
		fmt.Println("  classDef inicio fill:#b5503c,stroke:#8f3f2f,color:#fff,font-weight:bold")
		fmt.Printf("  class %s inicio\n", mermaidID(e.Code))
	}

	// The TERMINALS get their own class — in the reference drawing they are the two
	// coloured boxes at the end, and that is how a flowchart is read: the destination
	// stands out before the path.
	var terms []string
	for _, s := range steps {
		if !flowx.IsResult(s.Code) && s.Terminal {
			terms = append(terms, mermaidID(s.Code))
		}
	}
	if len(terms) > 0 {
		fmt.Println("  classDef fim fill:#3b82f6,stroke:#1d4ed8,color:#fff,font-weight:bold")
		fmt.Printf("  class %s fim\n", strings.Join(terms, ","))
	}
}

// rotuloDaSaida é o que vai NA SETA: o nome curto do resultado, que é o que o leitor
// segue ("BARRADO", "FILA VAZIA") — e não o código, que não diz nada a quem lê o desenho.
func exitLabel(fg *mapx.FlowGraph, t mapx.FlowTransition) string {
	if t.On == "" {
		if t.When != "" {
			return t.When
		}
		return "sempre"
	}
	if r, ok := flowx.StateByCode(fg, t.On); ok {
		// The result's title reads "BARRED: a blocking gate failed" — the arrow fits the
		// first part, and the rest is already in the destination box.
		if nome, _, achou := strings.Cut(r.Title, ":"); achou {
			return strings.TrimSpace(nome)
		}
		return r.Title
	}
	return t.On
}

// mermaidID replaces what mermaid does not accept in an identifier.
func mermaidID(code string) string {
	return strings.ReplaceAll(code, "-", "_")
}

// escapeMermaid protects the quotes in a label.
func escapeMermaid(s string) string {
	return strings.ReplaceAll(s, `"`, "'")
}
