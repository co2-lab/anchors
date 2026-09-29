package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/spf13/cobra"
)

// KEEP THE EVIDENCE OF A CHANGE THAT PROVES NOTHING NEW.
//
// A file's evidence — the tests that proved it, the stamps of its edges — is tied to the
// revision it was measured at, and any edit ages it. Most of the time that is right. But
// an edit may touch only what Anchors itself reads: a `@realizes` added to a rule line, a
// `#region` around a block, a comment citing a code. Anchors cannot tell that from a
// behaviour change in every language, and rerunning every suite for it is waste. The one
// who made the change can tell, and this is how they say it: the evidence moves to the
// current revision, with the reason recorded on the node.
//
// It also refreshes the `@contract` stamps of the doubles pointing at the file, and lists
// them: a stamp hashes the block it covers as written, comments included, so the same edit
// ages it, and the same declaration covers it.
func newKeepEvidenceCmd() *cobra.Command {
	var root, reason string
	var lines bool
	cmd := &cobra.Command{
		Use:   "keep-evidence <file>...",
		Short: "Declare that a change proves nothing new, keeping the files' evidence fresh",
		Long: `Moves the evidence of each file to its current content — the tests that proved it, the
stamps and judgments of its edges, what other tests recorded of it — because the change
does not touch what was proven. Anchors cannot tell a trace annotation, a region marker or
a comment from a behaviour change in every language; whoever made the change can, and says
why. The reason is required, and recorded on the file's node.

  anchors keep-evidence a.spec.md b.spec.md --reason "only @realizes added to the rule lines"
  anchors keep-evidence src/Screen.tsx --reason "#region markers around the scenarios" 
  anchors keep-evidence src/pay.ts --lines --reason "a comment on a line of its own was reworded"

Coverage and mutation name lines by number, and an edit that added or removed lines moved
them: they are kept only with --lines, when no line moved. When the map was rebuilt after
the change and no longer holds the evidence, it is taken from the map at HEAD.

The ` + "`@contract`" + ` stamps of the doubles pointing at the file are refreshed too, and listed.

Do not use it for a change in behaviour: the evidence is what stops that change from passing
unproven.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reason = strings.TrimSpace(reason)
			if reason == "" {
				return fmt.Errorf("--reason is required: say why the change proves nothing new")
			}
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			mapPath := filepath.Join(absRoot, mapx.DefaultPath)
			if _, err := os.Stat(mapPath); err != nil {
				return fmt.Errorf("no map at %s: run `anchors map build` first", mapx.DefaultPath)
			}
			head := headMap(absRoot)
			today := gitmeta.Today()
			var files []string
			err = mapx.Update(mapPath, func(g *mapx.Graph) error {
				for _, a := range args {
					rel := filepath.ToSlash(common.RelTo(absRoot, a))
					b, err := os.ReadFile(filepath.Join(absRoot, rel))
					if err != nil {
						return fmt.Errorf("%s: %w", rel, err)
					}
					n := g.Node(rel)
					if n == nil {
						return fmt.Errorf("%s is not in the map: run `anchors map build` first", rel)
					}
					if n.Signal == nil && head != nil {
						if h := head.Node(rel); h != nil && h.Signal != nil {
							n.Signal = h.Signal
						}
					}
					carried := g.KeepEvidence(rel, scan.ShortHash(b), reason, today, lines)
					if len(carried) == 0 {
						fmt.Printf("  %s — nothing to keep: its evidence is already at this content, or it has none\n", rel)
						continue
					}
					fmt.Printf("  %s — evidence kept: %s → %s\n", rel, strings.Join(carried, ", "), scan.ShortHash(b))
					files = append(files, rel)
				}
				return nil
			})
			if err != nil {
				return err
			}
			if len(files) > 0 {
				g, err := mapx.Load(mapPath)
				if err != nil {
					return err
				}
				fmt.Println("\n· @contract stamps pointing at these files, refreshed under the same declaration (if the change did alter behaviour, the doubles listed need adjusting):")
				if err := refreshStamps(absRoot, g, files, false); err != nil {
					return err
				}
			}
			fmt.Printf("\n%d file(s) kept their evidence. Reason recorded: %q\n", len(files), reason)
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&reason, "reason", "", "why the change proves nothing new (required)")
	cmd.Flags().BoolVar(&lines, "lines", false, "also keep coverage and mutation: no line moved")
	return cmd
}

// headMap is the map as HEAD has it, or nil when git has none.
func headMap(absRoot string) *mapx.Graph {
	out, err := exec.Command("git", "-C", absRoot, "show", "HEAD:./"+mapx.DefaultPath).Output()
	if err != nil {
		return nil // @resilient: no commit or no committed map — there is no earlier evidence to recover
	}
	g, err := mapx.LoadBytes(out, "HEAD:"+mapx.DefaultPath)
	if err != nil {
		return nil // @resilient: a committed map this binary cannot read has no evidence to offer
	}
	return g
}
