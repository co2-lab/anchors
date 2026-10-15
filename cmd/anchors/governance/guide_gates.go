// @anchors
//   code: GTGDG
//   ref: GVGDG

package governance

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/spf13/cobra"
)

// gatesSite is where each gate has its page.
const gatesSite = "https://anchors.co2lab.io/docs/gates/"

// newGuideGatesCmd explains the gates to an agent: what each one measures, whether it bars
// here, and the ways out of a failure. An agent learned a gate only from its failure, after
// the work; the catalog the site's pages are made from is in the binary, and nothing showed
// it on the terminal.
func newGuideGatesCmd() *cobra.Command {
	var root string
	var all bool
	cmd := &cobra.Command{
		Use:   "gates [gate]",
		Short: "the gates: what each one measures, whether it bars here, and the ways out of a failure",
		Long: `With no gate, lists the gates this project declares in anchors.yaml — what each one
measures, and whether it bars a commit or only informs. With --all, every gate of the catalog.

With a gate's name, explains it: what it measures, what it confronts, whether it bars here,
the question an AI answers for a judgment gate, and the ways out of a failure.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, _ := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if len(args) == 1 {
				return explainGate(cmd.OutOrStdout(), cfg, args[0])
			}
			listGates(cmd.OutOrStdout(), cfg, all)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&all, "all", false, "list every gate of the catalog, declared or not")
	return cmd
}

// gateCatalog is every gate of the catalog, by name, as a new project would be given it.
func gateCatalog() map[string]config.Gate {
	chosen := map[string]bool{"spec": true, "feature": true, "test": true, "guide": true, "plan": true, "code": true, "doc": true, "product": true, "flag": true}
	out := map[string]config.Gate{}
	for _, g := range initx.DefaultGates(chosen, true) {
		out[g.Name] = g
	}
	return out
}

// declaredGates are the gates the project declares, by name.
func declaredGates(cfg *config.Config) map[string]config.Gate {
	out := map[string]config.Gate{}
	if cfg == nil {
		return out
	}
	for _, g := range cfg.Gates {
		out[g.Name] = g
	}
	return out
}

// describe is what a gate measures: its own description, or the catalog's.
func describe(g config.Gate, catalog map[string]config.Gate) string {
	if d := g.Describe(); d != "" {
		return d
	}
	return catalog[g.Name].Describe()
}

func barsOrInforms(g config.Gate) string {
	if g.IsBlocking() {
		return "bars"
	}
	return "informs"
}

func listGates(w io.Writer, cfg *config.Config, all bool) {
	fmt.Fprint(w, "# Gates guide\n\n")
	catalog, declared := gateCatalog(), declaredGates(cfg)
	names := make([]string, 0, len(declared))
	for name := range declared {
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		fmt.Fprintln(w, "This project declares no gate (anchors.yaml, `gates:`).")
	} else {
		fmt.Fprintf(w, "The %d gates this project declares (anchors.yaml) — `bars` fails a commit, `informs` reports:\n\n", len(names))
		for _, name := range names {
			g := declared[name]
			fmt.Fprintf(w, "  %-32s %-8s %s\n", name, barsOrInforms(g), oneLine(describe(g, catalog)))
		}
	}
	if all {
		var rest []string
		for name := range catalog {
			if _, ok := declared[name]; !ok {
				rest = append(rest, name)
			}
		}
		sort.Strings(rest)
		if len(rest) > 0 {
			fmt.Fprintf(w, "\nThe %d gates of the catalog this project does not declare:\n\n", len(rest))
			for _, name := range rest {
				fmt.Fprintf(w, "  %-32s %s\n", name, oneLine(catalog[name].Describe()))
			}
		}
	} else if extra := len(catalog) - len(names); extra > 0 {
		fmt.Fprintf(w, "\n`anchors guide gates --all` lists the catalog's other gates too.\n")
	}
	fmt.Fprintln(w, "\n`anchors guide gates <gate>` explains one: what it measures, and the ways out of a failure.")
}

func explainGate(w io.Writer, cfg *config.Config, name string) error {
	catalog, declared := gateCatalog(), declaredGates(cfg)
	g, isDeclared := declared[name]
	c, inCatalog := catalog[name]
	if !isDeclared && !inCatalog {
		return fmt.Errorf("no gate %q in this project nor in the catalog — `anchors guide gates --all` lists them", name)
	}
	if !isDeclared {
		g = c
	}
	fmt.Fprintf(w, "# Gate: %s\n\n", name)
	if d := strings.TrimSpace(describe(g, catalog)); d != "" {
		fmt.Fprintf(w, "%s%s.\n\n", strings.ToUpper(d[:1]), strings.TrimSuffix(d[1:], "."))
	}
	switch {
	case isDeclared && g.IsBlocking():
		fmt.Fprintln(w, "In this project: declared, and it BARS — a failure stops the commit.")
	case isDeclared:
		fmt.Fprintln(w, "In this project: declared, and it INFORMS — a failure is reported, the commit goes on.")
	default:
		fmt.Fprintln(w, "In this project: not declared. It runs once anchors.yaml declares it under `gates:`.")
	}
	switch {
	case g.Measures == config.MeasuresJudgment:
		fmt.Fprintln(w, "Measured by: a judgment — an AI or a person answers the question below, through `anchors judge`.")
	case g.Run != "":
		fmt.Fprintf(w, "Measured by: the command `%s`.\n", g.Run)
	case g.Check != "":
		fmt.Fprintf(w, "Measured by: Anchors' own check `%s`.\n", g.Check)
	}
	if len(g.On) > 0 {
		fmt.Fprintf(w, "Confronts: %s.\n", strings.Join(g.On, ", "))
	}
	if len(g.Presupposes) > 0 {
		fmt.Fprintf(w, "Presupposes: %s — until anchors.yaml declares them, the gate is pending and asks nothing.\n", strings.Join(g.Presupposes, ", "))
	}
	if strings.TrimSpace(g.Ask) != "" {
		fmt.Fprintf(w, "\nThe question:\n\n  %s\n", strings.ReplaceAll(strings.TrimSpace(g.Ask), "\n", "\n  "))
	}
	fmt.Fprintf(w, `
When it fails, the failure's detail names what is missing. The ways out, in order:
  1. Fix what it names, then confront the change again: anchors check --changed <file>
  2. If the failure is right in general but not here, waive it for that unit, with the reason,
     in the commit message: [skip-%s@<CODE>: why]
  3. If the project decides the gate should only inform, that is anchors.yaml's: blocking: false
     — a decision for whoever owns the project, not a way around one failure.

More: %s%s
`, name, gatesSite, name)
	return nil
}

// oneLine is the first line of a text.
func oneLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return s
}
