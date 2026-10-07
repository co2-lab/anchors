// @anchors
//   code: MDCMA
//   ref: MDCMP

package mapcmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

func newMapDepsCmd() *cobra.Command {
	var root string
	var up bool
	var depth int
	cmd := &cobra.Command{
		Use:   "deps <CODE|file>",
		Short: "The dependency tree of a file: what it uses, or who uses it (--up)",
		Long: `Prints the dependency tree of a file, by the depends-on edges of the map — the ones the
code's ` + "`@dep:`" + ` flags declare and the ones the specs' Dependencies tables declare:

  anchors map deps TOKNS              — what the file of code TOKNS uses, down the tree
  anchors map deps src/theme/tokens.ts --up   — who uses it, up the tree
  anchors map deps ARSCR --depth 2    — two levels`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}
			start := resolveDepsStart(g, args[0])
			if start == "" {
				return fmt.Errorf("%s is no code nor file of the map", args[0])
			}
			for _, l := range DepsTree(g, start, up, depth) {
				fmt.Fprintln(cmd.OutOrStdout(), l)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().BoolVar(&up, "up", false, "who uses the file, instead of what it uses")
	cmd.Flags().IntVar(&depth, "depth", 0, "levels to print (0 = all)")
	return cmd
}

// resolveDepsStart is the file a code or a path names: a file's own code first, then a
// unit's code (its first file), then the path as given.
func resolveDepsStart(g *mapx.Graph, arg string) string {
	for _, n := range g.Nodes {
		if n.FileCode == arg {
			return n.ID
		}
	}
	for _, n := range g.Nodes {
		if n.Code == arg && n.Kind == mapx.KindCode {
			return n.ID
		}
	}
	if g.Node(filepath.ToSlash(arg)) != nil {
		return filepath.ToSlash(arg)
	}
	return ""
}

// DepsTree is the dependency tree of a file as indented lines — `CODE path` —, down what it
// uses or, with up, who uses it; a file already on the branch is marked and not walked
// again, so a cycle ends.
func DepsTree(g *mapx.Graph, start string, up bool, depth int) []string {
	next := map[string][]string{}
	for _, e := range g.Edges {
		if e.Type != mapx.EdgeDependsOn {
			continue
		}
		if up {
			next[e.To] = append(next[e.To], e.From)
		} else {
			next[e.From] = append(next[e.From], e.To)
		}
	}
	label := func(id string) string {
		if n := g.Node(id); n != nil && n.FileCode != "" {
			return n.FileCode + " " + id
		}
		return id
	}
	var out []string
	var walk func(id string, level int, branch map[string]bool)
	walk = func(id string, level int, branch map[string]bool) {
		prefix := strings.Repeat("  ", level)
		if branch[id] {
			out = append(out, prefix+label(id)+" ↺")
			return
		}
		out = append(out, prefix+label(id))
		if depth > 0 && level >= depth {
			return
		}
		branch[id] = true
		kids := append([]string(nil), next[id]...)
		sort.Strings(kids)
		for _, k := range kids {
			walk(k, level+1, branch)
		}
		delete(branch, id)
	}
	walk(start, 0, map[string]bool{})
	return out
}
