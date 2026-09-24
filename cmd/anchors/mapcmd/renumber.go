package mapcmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/spf13/cobra"
)

// `anchors renumber` resolves the collision of revision numbers at rebase.
//
// Two open pull requests revising the same spec each take the next free `-R000N`, and both
// take the same one. The second to merge would land a `R0003` that already means something
// else on the base. The number stays — it says how many times the document changed — and
// the branch that arrives second renumbers ITS revisions, never the base's: the base's are
// history other people have already read and cited. See `gate.PlanRenumber`.
//
// It is the sibling of `recode`: `recode` renames an identity code across the project, this
// renames a revision code across what the branch changed.
func newRenumberCmd() *cobra.Command {
	var root, base string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "renumber [spec files...]",
		Short: "Renumber the revisions this branch added when the base already took their number",
		Long: `Finds the revision blocks (` + "`CODE-R000N`" + `) this branch ADDED whose number the base
already uses — two pull requests that each took the next free number — and moves them to
the next free numbers, rewriting every citation of the old code in what the branch changed:
specs, plans, features, tests and code comments.

Only the branch's revisions move. A revision the base has is history others have already
read and cited; renumbering it would point every one of those citations at the wrong change.
Likewise, a citation is rewritten only on a line the branch added — a line that existed
where the branch forked cites what the base means by that code.

With no file given, every Markdown file the branch changed against the base is examined.
Run it after rebasing onto the base (or before — both work):

    git rebase origin/main
    anchors renumber --dry-run                 what would move, and where it is cited
    anchors renumber                           applies
    anchors renumber specs/Pricing.spec.md     only this spec's revisions
    anchors renumber --base origin/develop     against another base

The base defaults to ` + "`origin/<integration_branch>`" + ` (or the local branch when there is no remote).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load %s: %w", config.DefaultFile, err)
			}
			if base == "" {
				base = defaultBase(absRoot, cfg)
			}
			mb, err := gitOut(absRoot, "merge-base", base, "HEAD")
			if err != nil {
				return fmt.Errorf("no merge base between %s and HEAD: %w", base, err)
			}
			mb = strings.TrimSpace(mb)

			changed, err := branchChanged(absRoot, mb)
			if err != nil {
				return err
			}

			var candidates []string
			if len(args) > 0 {
				for _, a := range args {
					candidates = append(candidates, common.RelTo(absRoot, a))
				}
			} else {
				for _, c := range changed {
					if strings.HasSuffix(c, ".md") {
						candidates = append(candidates, c)
					}
				}
			}

			var files []gate.RevisionFile
			for _, c := range candidates {
				b, err := os.ReadFile(filepath.Join(absRoot, c))
				if err != nil {
					return fmt.Errorf("read %s: %w", c, err)
				}
				files = append(files, gate.RevisionFile{
					Path:      c,
					Working:   string(b),
					MergeBase: showAt(absRoot, mb, c),
					Base:      showAt(absRoot, base, c),
				})
			}
			plan, err := gate.PlanRenumber(files)
			if err != nil {
				return err
			}
			if len(plan) == 0 {
				fmt.Printf("✓ no revision this branch added collides with %s — nothing to renumber.\n", base)
				return nil
			}

			renames := map[string]string{}
			fmt.Printf("renumber against %s:\n\n", base)
			for _, r := range plan {
				renames[r.Old] = r.New
				fmt.Printf("  %s → %s   (%s)\n", r.Old, r.New, r.File)
			}

			// Every file the branch changed is a place the old code may be cited — the spec
			// itself, its feature and test, a code comment, the plan that seeded it.
			type rewrite struct {
				path    string
				content string
				n       int
			}
			var rewrites []rewrite
			for _, c := range changed {
				b, err := os.ReadFile(filepath.Join(absRoot, c))
				if err != nil || bytes.IndexByte(b, 0) >= 0 {
					continue // deleted, unreadable or binary: nothing to cite
				}
				out, n := gate.RewriteRevisionCitations(string(b), showAt(absRoot, mb, c), renames)
				if n > 0 {
					rewrites = append(rewrites, rewrite{c, out, n})
				}
			}
			sort.Slice(rewrites, func(i, j int) bool { return rewrites[i].path < rewrites[j].path })
			fmt.Println("\ncitations rewritten (on the lines this branch added):")
			for _, w := range rewrites {
				fmt.Printf("  %-60s %d\n", w.path, w.n)
			}
			fmt.Println()

			if dryRun {
				fmt.Println("(dry-run — nothing was written.)")
				return nil
			}
			for _, w := range rewrites {
				p := filepath.Join(absRoot, w.path)
				info, err := os.Stat(p)
				if err != nil {
					return err
				}
				if err := os.WriteFile(p, []byte(w.content), info.Mode().Perm()); err != nil {
					return fmt.Errorf("write %s: %w", w.path, err)
				}
			}
			fmt.Printf("✓ %d revision(s) renumbered, %d file(s) rewritten.\n", len(plan), len(rewrites))
			fmt.Println("  review with `git diff`; a citation the branch wrote of the BASE's revision would have moved too.")
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&base, "base", "", "the ref the branch merges into (default: origin/<integration_branch>)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be renumbered, without writing")
	return cmd
}

// defaultBase is where the branch merges: the remote copy of the integration branch when it
// exists, because after `git fetch` the local branch is usually behind it.
func defaultBase(root string, cfg *config.Config) string {
	branch := cfg.Workflow.IntegrationBranchOrDefault()
	if _, err := gitOut(root, "rev-parse", "--verify", "--quiet", "origin/"+branch); err == nil {
		return "origin/" + branch
	}
	return branch
}

// branchChanged lists what the branch changed since the merge base — committed, staged,
// unstaged and untracked — relative to the project root.
func branchChanged(root, mb string) ([]string, error) {
	diff, err := gitOut(root, "diff", "--name-only", "--relative", mb)
	if err != nil {
		return nil, fmt.Errorf("git diff: %w", err)
	}
	untracked, err := gitOut(root, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %w", err)
	}
	seen := map[string]bool{}
	var out []string
	for _, l := range strings.Split(diff+"\n"+untracked, "\n") {
		if l = strings.TrimSpace(l); l != "" && !seen[l] {
			seen[l] = true
			out = append(out, l)
		}
	}
	sort.Strings(out)
	return out, nil
}

// showAt is the file's content at a ref, or "" when it did not exist there.
func showAt(root, ref, rel string) string {
	out, err := gitOut(root, "show", ref+":./"+rel)
	if err != nil {
		return ""
	}
	return out
}

func gitOut(root string, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	var stderr bytes.Buffer
	c.Stderr = &stderr
	out, err := c.Output()
	if err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("%w: %s", err, msg)
		}
		return "", err
	}
	return string(out), nil
}
