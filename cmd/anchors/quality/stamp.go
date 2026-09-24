package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

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
	var refresh []string
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
    anchors stamp --dry-run                 says what it would write

AFTER CHANGING A FUNCTION, refresh the stamps that point at its file — and read the answer:

    anchors stamp --refresh src/hooks/balance.ts

It lists every double stamped against the previous version (test, line, member, old and
new hash, and how the stamped block changed from HEAD), then updates those hashes. Each
double listed reproduces the OLD contract: adjust it in the same commit. A stamp whose
anchor line is gone (the member was renamed or removed) is not refreshed — re-stamp it by
hand. With --dry-run it only lists.`,
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
			if len(refresh) > 0 {
				return refreshStamps(absRoot, g, refresh, dryRun)
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
	cmd.Flags().StringSliceVar(&refresh, "refresh", nil,
		"after changing this file: list the doubles stamped against its previous version and update their stamps")
	return cmd
}

// refreshStamps is `anchors stamp --refresh <file>`: for each changed file, the doubles
// stamped against its previous version — with how the stamped block changed — and their
// stamps updated. The list is the point: it is the work the change created.
func refreshStamps(root string, g *mapx.Graph, files []string, dryRun bool) error {
	total := 0
	for _, f := range files {
		rel := common.RelTo(root, f)
		found, err := gate.RefreshStamps(g, root, rel, !dryRun)
		if err != nil {
			return err
		}
		if len(found) == 0 {
			fmt.Printf("%s — no double is stamped against a previous version of it.\n", rel)
			continue
		}
		total += len(found)
		antigo, _ := exec.Command("git", "-C", root, "show", "HEAD:"+rel).Output()
		novo, _ := os.ReadFile(filepath.Join(root, rel))
		fmt.Printf("%s changed — %d double(s) stamped against the previous version:\n\n", rel, len(found))
		for _, r := range found {
			if r.Err != "" {
				fmt.Printf("  %s:%d  `%s`\n      NOT refreshed: %s — the member was renamed or removed;\n"+
					"      adjust the double, delete this stamp and run `anchors stamp`.\n\n", r.Test, r.Line, r.Anchor, r.Err)
				continue
			}
			fmt.Printf("  %s:%d  `%s`  (%s → %s)\n", r.Test, r.Line, r.Anchor, r.Old, r.New)
			velho, okV := gate.StampSnippet(string(antigo), r.Anchor, r.Count)
			atual, okA := gate.StampSnippet(string(novo), r.Anchor, r.Count)
			if okV && okA {
				for _, l := range blockDiff(velho, atual) {
					fmt.Printf("      %s\n", l)
				}
			} else {
				fmt.Printf("      (no HEAD version of the block to compare)\n")
			}
			fmt.Println()
		}
	}
	if total == 0 {
		return nil
	}
	if dryRun {
		fmt.Printf("would update the stamps above (--dry-run: nothing written).\n")
		return nil
	}
	fmt.Printf("Stamps updated. Each double above reproduced the OLD contract: adjust it to the new one\n" +
		"in this same commit — the stamp now says the double matches, and only you have checked it.\n")
	return nil
}

// blockDiff lists the lines of a stamped block that left (-) and that came in (+).
func blockDiff(velho, atual string) []string {
	va := strings.Split(velho, "\n")
	aa := strings.Split(atual, "\n")
	emA := map[string]int{}
	for _, l := range aa {
		emA[l]++
	}
	emV := map[string]int{}
	for _, l := range va {
		emV[l]++
	}
	var out []string
	for _, l := range va {
		if emA[l] > 0 {
			emA[l]--
			continue
		}
		out = append(out, "- "+l)
	}
	for _, l := range aa {
		if emV[l] > 0 {
			emV[l]--
			continue
		}
		out = append(out, "+ "+l)
	}
	return out
}
