package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/mapx"
)

// newStampCmd writes the `@contract` stamps `mock-stamped` checks, on the doubles that
// have none.
//
// It only ADDS. An existing stamp is never rewritten, even when it diverges: a divergence
// is the gate saying the double may have drifted, and refreshing it would let the stamp
// certify itself. See `gate.GenerateStamps`.
func newStampCmd() *cobra.Command {
	var root, mapPath string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "stamp [test files...]",
		Short: "Write the missing `@contract` stamps on test doubles (never rewrites an existing one)",
		Long: `Writes, above every double of a governed module that has no stamp, the
` + "`@contract`" + ` stamp that the mock-stamped gate recomputes against the real module.

Each key of the double's factory that names an export of the module gets its own stamp,
covering that export's block. With no such key, one stamp covers the whole module.

It only ADDS stamps. An existing stamp is never rewritten, even when it diverges — the
divergence is the gate telling you the double may be stale, and that calls for a person
looking at the double, not a tool refreshing the hash.

With no file given, every test of the map is stamped.

    anchors stamp                           every test
    anchors stamp src/Home.test.tsx         one test
    anchors stamp --dry-run                 says what it would write`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
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

			var tests []string
			if len(args) > 0 {
				for _, a := range args {
					tests = append(tests, common.RelTo(absRoot, a))
				}
			} else {
				for _, n := range g.Nodes {
					if n.Kind == mapx.KindTest {
						tests = append(tests, n.ID)
					}
				}
			}
			sort.Strings(tests)

			var files, stamps, skips int
			for _, t := range tests {
				path := filepath.Join(absRoot, t)
				b, err := os.ReadFile(path)
				if err != nil {
					fmt.Printf("  %s — could not be read: %v\n", t, err)
					continue
				}
				out, written, skipped, err := gate.GenerateStamps(string(b), t, absRoot, g, cfg)
				if err != nil {
					return err
				}
				if len(written) == 0 && len(skipped) == 0 {
					continue
				}
				fmt.Printf("%s\n", t)
				for _, w := range written {
					fmt.Printf("  + %s → %s | %s | %d\n", w.Module, w.File, w.Anchor, w.Count)
				}
				for _, s := range skipped {
					fmt.Printf("  ~ %s — not stamped: %s\n", s.Module, s.Reason)
				}
				stamps += len(written)
				skips += len(skipped)
				if len(written) > 0 {
					files++
					if !dryRun {
						if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
							return err
						}
					}
				}
			}

			verb := "wrote"
			if dryRun {
				verb = "would write"
			}
			fmt.Printf("\n%s %d stamp(s) in %d file(s); %d double(s) left unstamped (see `~` above).\n",
				verb, stamps, files, skips)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "say what would be written, without writing")
	return cmd
}
