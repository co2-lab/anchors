// @anchors
//   code: MNCMA
//   ref: MNCMP

package mapcmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

func newMapNavCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "nav [screen]",
		Short: "The app's navigation: which screen leads to which",
		Long: `Prints the navigation between the app's screens — the screens' Out tables and the
code's ` + "`@navigates:`" + ` flags:

  anchors map nav              — every screen and where it leads
  anchors map nav GoalEdit     — one screen: where it comes from and where it leads`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			g, err := mapx.Load(filepath.Join(absRoot, mapx.DefaultPath))
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}
			only := ""
			if len(args) == 1 {
				only = args[0]
			}
			lines, ok := NavLines(gate.NavEdges(absRoot, g), only)
			if !ok {
				return fmt.Errorf("%s is no screen of the app", only)
			}
			for _, l := range lines {
				fmt.Fprintln(cmd.OutOrStdout(), l)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}

// screenLabel is a screen's name: its spec's stem.
func screenLabel(spec string) string { return strings.TrimSuffix(filepath.Base(spec), ".spec.md") }

// NavLines are the navigation's lines: `A → B, C` per screen, sorted; with a screen named,
// `← origins` and `→ destinations` of that one. A name that is no screen is not ok.
func NavLines(edges map[string]map[string]bool, only string) ([]string, bool) {
	sorted := func(set map[string]bool) []string {
		var out []string
		for k := range set {
			out = append(out, screenLabel(k))
		}
		sort.Strings(out)
		return out
	}
	if only == "" {
		var froms []string
		for f := range edges {
			froms = append(froms, f)
		}
		sort.Slice(froms, func(i, j int) bool { return screenLabel(froms[i]) < screenLabel(froms[j]) })
		var out []string
		for _, f := range froms {
			out = append(out, screenLabel(f)+" → "+strings.Join(sorted(edges[f]), ", "))
		}
		return out, true
	}
	match := func(spec string) bool {
		l := screenLabel(spec)
		return strings.EqualFold(l, only) || strings.EqualFold(strings.TrimSuffix(l, "Screen"), only)
	}
	in := map[string]bool{}
	var out map[string]bool
	found := false
	for f, tos := range edges {
		if match(f) {
			out, found = tos, true
		}
		for t := range tos {
			if match(t) {
				in[f], found = true, true
			}
		}
	}
	if !found {
		return nil, false
	}
	return []string{"← " + strings.Join(sorted(in), ", "), "→ " + strings.Join(sorted(out), ", ")}, true
}
