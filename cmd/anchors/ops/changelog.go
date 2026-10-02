// @anchors
//   ref: CLGCM

package ops

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/spf13/cobra"
)

func newChangelogCmd() *cobra.Command {
	var root, from, to string
	var all, unreleased, write bool
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Build the technical changelog from the commits, release by release",
		Long: `This is the TECHNICAL changelog: every feat, fix and breaking change, in the commits' own
words, for whoever works on the code. It is not the product's release notes — for those,
have an agent synthesize a product changelog from this one (` + "`anchors guide changelog`" + `).

A release is what the history holds between two tags. Its commits, in Conventional
Commits, say what goes where:

  breaking changes   a ` + "`!`" + ` after the type, or a ` + "`BREAKING CHANGE:`" + ` footer
  features           feat
  bugs fixed         fix with a ` + "`Bug:`" + ` footer — a defect that shipped
  fixes              fix without it — a correction of work that never reached anyone

The other types (refactor, test, chore, …) are history, not changelog, and are left out.

By default the latest release is printed. ` + "`--from`" + ` takes every release after that tag,
` + "`--all`" + ` every one, ` + "`--unreleased`" + ` adds what came after the last tag.

` + "`--write`" + ` writes to the file the ` + "`changelog:`" + ` block of anchors.yaml names: in ` + "`incremental`" + `
mode (the default) the releases the file does not hold yet go on its top, and what it holds
is kept as it is; in ` + "`per_version`" + ` mode each release gets its own file. The headings follow
the project's ` + "`lang`" + `, and ` + "`changelog.template`" + ` replaces the built-in template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cl, err := loadChangelogConfig(absRoot)
			if err != nil {
				return err
			}
			tmpl := changelog.DefaultTemplate
			if cl != nil && cl.Template != "" {
				b, err := os.ReadFile(filepath.Join(absRoot, cl.Template))
				if err != nil {
					return fmt.Errorf("changelog.template: %w", err)
				}
				tmpl = string(b)
			}
			rels, err := pickReleases(changelog.Git{Root: absRoot}, from, to, all, unreleased)
			if err != nil {
				return err
			}
			t := func(k string) string { return i18n.T(k) }
			var blocks []changelog.Block
			for _, r := range rels {
				if r.Empty() {
					continue
				}
				text, err := changelog.Render(tmpl, r, t)
				if err != nil {
					return err
				}
				blocks = append(blocks, changelog.Block{Version: r.Version, Text: text})
			}
			if !write {
				for i, b := range blocks {
					if i > 0 {
						fmt.Fprintln(cmd.OutOrStdout())
					}
					fmt.Fprint(cmd.OutOrStdout(), strings.TrimRight(b.Text, "\n")+"\n")
				}
				if len(blocks) == 0 {
					fmt.Fprintln(cmd.ErrOrStderr(), "anchors: nothing to list — no feat, fix or breaking change in the range")
				}
				return nil
			}
			return writeChangelog(cmd, absRoot, cl, blocks)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&from, "from", "", "list every release after this tag")
	cmd.Flags().StringVar(&to, "to", "HEAD", "the ref the history is read up to")
	cmd.Flags().BoolVar(&all, "all", false, "list every release")
	cmd.Flags().BoolVar(&unreleased, "unreleased", false, "also list what came after the last tag")
	cmd.Flags().BoolVar(&write, "write", false, "write to the file(s) `changelog:` names instead of printing")
	return cmd
}

// loadChangelogConfig reads the `changelog:` block, and sets the project's language for
// the headings. A project with no anchors.yaml still has a history: it gets the defaults.
func loadChangelogConfig(absRoot string) (*config.Changelog, error) {
	path := filepath.Join(absRoot, config.DefaultFile)
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	cfg, err := config.Load(path)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", config.DefaultFile, err)
	}
	return cfg.Changelog, nil
}

// pickReleases reads the releases the flags ask for, newest first. With no tag at all,
// everything is unreleased.
func pickReleases(g changelog.Git, from, to string, all, unreleased bool) ([]changelog.Release, error) {
	tags, err := g.Tags(to)
	if err != nil {
		return nil, err
	}
	var want []string
	switch {
	case all:
		want = tags
	case from != "":
		i := indexOf(tags, from)
		if i < 0 {
			return nil, fmt.Errorf("--from %q: no such tag reachable from %s", from, to)
		}
		want = tags[i+1:]
	case len(tags) > 0:
		want = tags[len(tags)-1:]
	}
	return g.Releases(tags, want, to, unreleased || len(tags) == 0)
}

func indexOf(xs []string, x string) int {
	for i, y := range xs {
		if y == x {
			return i
		}
	}
	return -1
}

// writeChangelog writes the blocks where the project's mode says.
func writeChangelog(cmd *cobra.Command, absRoot string, cl *config.Changelog, blocks []changelog.Block) error {
	path := filepath.Join(absRoot, cl.PathOrDefault())
	if cl.ModeOrDefault() == "per_version" {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return err
		}
		wrote := 0
		for _, b := range blocks {
			name := b.Version
			if name == "" {
				name = "unreleased"
			}
			file := filepath.Join(path, name+".md")
			// A release's file, once written, is the project's: it may have been edited.
			// Only what is not released yet is rewritten, since it grows until the tag.
			if _, err := os.Stat(file); err == nil && b.Version != "" {
				continue
			}
			if err := os.WriteFile(file, []byte(strings.TrimRight(b.Text, "\n")+"\n"), 0o644); err != nil {
				return err
			}
			wrote++
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", relTo(absRoot, file))
		}
		if wrote == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "every release already has its file")
		}
		return nil
	}
	existing, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	out := changelog.Prepend(string(existing), i18n.T("changelog.title"), blocks)
	if out == string(existing) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s already holds every release\n", relTo(absRoot, path))
		return nil
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", relTo(absRoot, path))
	return nil
}

func relTo(root, p string) string {
	if r, err := filepath.Rel(root, p); err == nil {
		return filepath.ToSlash(r)
	}
	return p
}
