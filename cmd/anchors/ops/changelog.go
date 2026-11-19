// @anchors
//   code: CHCMC
//   ref: CLGCM

package ops

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/pathref"
	"github.com/co2-lab/anchors/internal/version"
	"github.com/spf13/cobra"
)

func newChangelogCmd() *cobra.Command {
	var root, from, to, product, branch string
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
the project's ` + "`lang`" + `, and ` + "`changelog.template`" + ` replaces the built-in template.

With a ` + "`versioning:`" + ` block, the changelog is a product's (` + "`--product`" + `, needed when there are
several): its releases lie between its own tags, ordered by its scheme, and each is headed by
the version the tag records. A branch that releases (` + "`release: true`" + `, or every branch when no
rule names one) lists its releases; another lists the candidate, what came after the latest
release, under the version its strategy computes. ` + "`--branch`" + ` picks the branch, by default the
first that releases, else the current one. ` + "`--write`" + ` writes every output of
` + "`versioning.<product>.changelog.outputs`" + `, or the ` + "`changelog:`" + ` block when it declares none, and only
on a branch the changelog is written on (` + "`changelog: true`" + `).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := loadChangelogConfig(absRoot)
			if err != nil {
				return err
			}
			var cl *config.Changelog
			if cfg != nil {
				cl = cfg.Changelog
			}
			outputs := []*config.Changelog{cl}
			labels := []string{"changelog.template"}
			var rels []changelog.Release
			if cfg != nil && len(cfg.Versioning) > 0 {
				pr, err := productReleases(absRoot, cfg.Versioning, product, branch, from, to, all, unreleased)
				if err != nil {
					return err
				}
				if write && !pr.product.WritesChangelog(pr.branch) {
					return fmt.Errorf("the changelog of %s is not written on %s: mark the branch changelog: true", pr.name, or(pr.branch, "a detached HEAD"))
				}
				if pc := pr.product.Changelog; pc != nil && len(pc.Outputs) > 0 {
					outputs, labels = outputs[:0], labels[:0]
					for i := range pc.Outputs {
						outputs = append(outputs, &pc.Outputs[i].Changelog)
						labels = append(labels, fmt.Sprintf("versioning.%s.changelog.outputs[%d].template", pr.name, i))
					}
				}
				rels = pr.releases
			} else if rels, err = pickReleases(changelog.Git{Root: absRoot}, from, to, all, unreleased); err != nil {
				return err
			}
			for i, out := range outputs {
				blocks, err := renderBlocks(absRoot, out, labels[i], rels)
				if err != nil {
					return err
				}
				if !write {
					printBlocks(cmd, blocks)
					return nil
				}
				if i > 0 && len(blocks) == 0 {
					continue
				}
				if err := writeChangelog(cmd, absRoot, out, blocks); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&from, "from", "", "list every release after this tag")
	cmd.Flags().StringVar(&to, "to", "HEAD", "the ref the history is read up to")
	cmd.Flags().BoolVar(&all, "all", false, "list every release")
	cmd.Flags().BoolVar(&unreleased, "unreleased", false, "also list what came after the last tag")
	cmd.Flags().BoolVar(&write, "write", false, "write to the file(s) `changelog:` names instead of printing")
	cmd.Flags().StringVar(&product, "product", "", "the product of `versioning:` (needed when there are several)")
	cmd.Flags().StringVar(&branch, "branch", "", "the branch whose changelog is read (default: the first that releases, else the current one)")
	return cmd
}

// renderBlocks renders the releases with an output's template, the built-in one when it
// names none; an empty release renders nothing.
func renderBlocks(absRoot string, cl *config.Changelog, label string, rels []changelog.Release) ([]changelog.Block, error) {
	tmpl := changelog.DefaultTemplate
	if cl != nil && cl.Template != "" {
		b, err := pathref.Read(absRoot, cl.Template)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		tmpl = string(b)
	}
	t := func(k string) string { return i18n.T(k) }
	var blocks []changelog.Block
	for _, r := range rels {
		if r.Empty() {
			continue
		}
		text, err := changelog.Render(tmpl, r, t)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, changelog.Block{Version: r.Version, Text: text})
	}
	return blocks, nil
}

func printBlocks(cmd *cobra.Command, blocks []changelog.Block) {
	for i, b := range blocks {
		if i > 0 {
			fmt.Fprintln(cmd.OutOrStdout())
		}
		fmt.Fprint(cmd.OutOrStdout(), strings.TrimRight(b.Text, "\n")+"\n")
	}
	if len(blocks) == 0 {
		fmt.Fprintln(cmd.ErrOrStderr(), "anchors: nothing to list — no feat, fix or breaking change in the range")
	}
}

// productRels is a product's changelog: its releases, and the branch they were read for.
type productRels struct {
	name, branch string
	product      *config.Product
	releases     []changelog.Release
}

// productReleases reads a product's releases between its own tags, in its scheme's order,
// for the branch: the releases where it releases, else the candidate alone.
func productReleases(absRoot string, v config.Versioning, name, branch, from, to string, all, unreleased bool) (*productRels, error) {
	name, p, err := v.Pick(name, "--product")
	if err != nil {
		return nil, err
	}
	s, _ := p.SchemeOf()
	g := changelog.Git{Root: absRoot}
	every, err := g.Tags(to)
	if err != nil {
		return nil, err
	}
	type tagged struct {
		tag string
		v   version.V
	}
	var ts []tagged
	for _, t := range every {
		raw, ok := p.VersionOfTag(t)
		if !ok {
			continue
		}
		pv, err := s.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("tag %s: %w", t, err)
		}
		ts = append(ts, tagged{t, pv})
	}
	sort.SliceStable(ts, func(i, j int) bool { return version.Compare(ts[i].v, ts[j].v) < 0 })
	tags := make([]string, len(ts))
	versionOf := map[string]string{}
	for i, t := range ts {
		tags[i], versionOf[t.tag] = t.tag, t.v.Raw
	}
	if branch == "" {
		if rb := p.ReleaseBranch(); rb != "" && !strings.ContainsAny(rb, "*?[") && !strings.HasPrefix(rb, "re:") {
			branch = rb
		} else if out, err := exec.Command("git", "-C", absRoot, "rev-parse", "--abbrev-ref", "HEAD").Output(); err == nil && strings.TrimSpace(string(out)) != "HEAD" {
			branch = strings.TrimSpace(string(out))
		}
	}
	releasing := p.Releases(branch)
	var want []string
	unrel := true
	if releasing {
		unrel = unreleased || len(tags) == 0
		switch {
		case all:
			want = tags
		case from != "":
			i := indexOf(tags, from)
			if i < 0 {
				i = indexOf(tags, p.TagFor(from))
			}
			if i < 0 {
				return nil, fmt.Errorf("--from %q: no tag of %s reachable from %s", from, name, to)
			}
			want = tags[i+1:]
		case len(tags) > 0:
			want = tags[len(tags)-1:]
		}
	}
	rels, err := g.Releases(tags, want, to, unrel)
	if err != nil {
		return nil, err
	}
	for i := range rels {
		if rels[i].Tag != "" {
			rels[i].Version = versionOf[rels[i].Tag]
			continue
		}
		rev := to
		if rev == "HEAD" {
			rev = ""
		}
		st, err := stateAt(absRoot, v, name, "--product", branch, rev)
		if err != nil {
			if !releasing {
				return nil, err
			}
			continue
		}
		if st.Next != "" && st.Next != st.Released {
			rels[i].Candidate = st.Next
		}
	}
	return &productRels{name: name, branch: branch, product: p, releases: rels}, nil
}

// loadChangelogConfig reads the configuration, and sets the project's language for the
// headings. A project with no anchors.yaml still has a history: it gets the defaults (nil).
func loadChangelogConfig(absRoot string) (*config.Config, error) {
	path := filepath.Join(absRoot, config.DefaultFile)
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	cfg, err := config.Load(path)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", config.DefaultFile, err)
	}
	return cfg, nil
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
