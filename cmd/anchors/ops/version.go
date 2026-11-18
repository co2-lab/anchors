// @anchors
//   code: VRSNH
//   ref: VRSNI

package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	var root, product, branch string
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show a product's version, its latest release and the next version its branch asks for",
		Long: `Reads the ` + "`versioning:`" + ` block of anchors.yaml and shows, for one product:

  sources    the version each source file holds (they must agree)
  released   the latest release: the highest tag of the product reachable from HEAD
  branch     the strategy the branch resolves to, and the rule that decided it
  next       the version that strategy computes, and why

A ` + "`Version:`" + ` footer on a commit since the release forces the next version. ` + "`--json`" + ` prints the
same as data, with the commits since the release already classified: it is the context a
project's own script receives to decide a version.

` + "`anchors version bump`" + ` writes a version into every source.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := versionState(root, product, branch)
			if err != nil {
				return err
			}
			if asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				enc.SetEscapeHTML(false)
				return enc.Encode(st)
			}
			printVersionState(cmd.OutOrStdout(), st)
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&root, "root", ".", "project root")
	cmd.PersistentFlags().StringVar(&product, "product", "", "the product of `versioning:` (needed when there are several)")
	cmd.PersistentFlags().StringVar(&branch, "branch", "", "the branch whose strategy decides (default: the current branch)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "print the state as JSON")
	cmd.AddCommand(newVersionBumpCmd(&root, &product, &branch))
	return cmd
}

func newVersionBumpCmd(root, product, branch *string) *cobra.Command {
	return &cobra.Command{
		Use:   "bump [part|version]",
		Short: "Write the next version into every source of the product",
		Long: `With no argument, writes the version the branch's strategy computes. A part of the
scheme (` + "`major`" + `, ` + "`minor`" + `, ` + "`patch`" + `, a pattern's group) raises the version the sources hold; an
exact version is written as given. Whatever decides it, the version must be readable, above the
one the sources hold and not below the latest release. Only the version's bytes change in each
file.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := versionState(*root, *product, *branch)
			if err != nil {
				return err
			}
			want := ""
			if len(args) == 1 {
				want = args[0]
			}
			return bumpVersion(cmd.OutOrStdout(), st, want)
		},
	}
}

// versionSource is one source's file and the version it holds.
type versionSource struct {
	File    string `json:"file"`
	Version string `json:"version"`
}

// versionCommit is a commit since the release, as a script reads it.
type versionCommit struct {
	Hash     string            `json:"hash"`
	Subject  string            `json:"subject"`
	Type     string            `json:"type,omitempty"`
	Breaking bool              `json:"breaking,omitempty"`
	Footers  map[string]string `json:"footers,omitempty"`
}

// versionStateT is what `anchors version` shows, and what `--json` prints.
type versionStateT struct {
	Product     string          `json:"product"`
	Scheme      string          `json:"scheme"`
	Branch      string          `json:"branch"`
	Rule        string          `json:"rule"`
	Strategy    string          `json:"strategy"`
	Current     string          `json:"current"`
	Sources     []versionSource `json:"sources"`
	Released    string          `json:"released"`
	ReleasedTag string          `json:"released_tag"`
	Commits     []versionCommit `json:"commits"`
	Next        string          `json:"next"`
	Reason      string          `json:"reason"`
	Forced      bool            `json:"forced"`

	root    string
	product *config.Product
	scheme  version.Scheme
}

// versionState reads a product's state: its sources, its latest release, the commits since
// it, the branch's strategy and the version it computes.
func versionState(root, name, branch string) (*versionStateT, error) {
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return nil, err
	}
	cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", config.DefaultFile, err)
	}
	name, p, err := pickProduct(cfg.Versioning, name)
	if err != nil {
		return nil, err
	}
	if branch == "" {
		out, err := exec.Command("git", "-C", absRoot, "rev-parse", "--abbrev-ref", "HEAD").Output()
		branch = strings.TrimSpace(string(out))
		if err != nil || branch == "" || branch == "HEAD" {
			return nil, errors.New("no current branch (a detached HEAD, or no repository): name it with --branch")
		}
	}
	s, _ := p.SchemeOf()
	st := &versionStateT{Product: name, Scheme: s.Kind, Branch: branch, root: absRoot, product: p, scheme: s,
		Sources: []versionSource{}, Commits: []versionCommit{}}
	svs, err := p.ReadSources(func(f string) ([]byte, error) { return os.ReadFile(filepath.Join(absRoot, filepath.FromSlash(f))) })
	if err != nil {
		return nil, err
	}
	for _, sv := range svs {
		st.Sources = append(st.Sources, versionSource{File: sv.File, Version: sv.Version})
	}
	if st.Current, err = config.Agreed(svs); err != nil {
		return nil, err
	}
	g := changelog.Git{Root: absRoot}
	if st.Released, st.ReleasedTag, err = latestRelease(g, p, s); err != nil {
		return nil, err
	}
	commits, err := g.Commits(st.ReleasedTag, "HEAD")
	if err != nil {
		return nil, err
	}
	for _, c := range commits {
		st.Commits = append(st.Commits, describeCommit(c))
	}
	strategy, rule := p.NextFor(branch)
	st.Strategy, st.Rule = strategy.Name, rule
	if rule == "" {
		st.Rule = "default"
	}
	res, err := version.Next(version.Input{Scheme: s, Strategy: strategy, Current: st.Current, Released: st.Released, Commits: commits, Today: time.Now()})
	if err != nil {
		return nil, fmt.Errorf("%s on %s (%s): %w", name, branch, strategy.Name, err)
	}
	st.Next, st.Reason, st.Forced = res.Version, res.Reason, res.Forced
	return st, nil
}

// pickProduct is the product asked for, or the only one.
func pickProduct(v config.Versioning, name string) (string, *config.Product, error) {
	names := make([]string, 0, len(v))
	for n := range v {
		names = append(names, n)
	}
	sort.Strings(names)
	switch {
	case len(v) == 0:
		return "", nil, errors.New("anchors.yaml declares no `versioning:` block: nothing is versioned")
	case name != "":
		if p, ok := v[name]; ok {
			return name, p, nil
		}
		return "", nil, fmt.Errorf("versioning has no product %q (%s)", name, strings.Join(names, ", "))
	case len(v) > 1:
		return "", nil, fmt.Errorf("versioning has %d products — name one with --product (%s)", len(v), strings.Join(names, ", "))
	}
	return names[0], v[names[0]], nil
}

// latestRelease is the highest version among the product's tags reachable from HEAD. A tag
// of the product the scheme cannot read is an error: skipped, a release would go unseen.
func latestRelease(g changelog.Git, p *config.Product, s version.Scheme) (string, string, error) {
	tags, err := g.Tags("HEAD")
	if err != nil {
		return "", "", err
	}
	var best *version.V
	bestTag := ""
	for _, t := range tags {
		raw, ok := p.VersionOfTag(t)
		if !ok {
			continue
		}
		v, err := s.Parse(raw)
		if err != nil {
			return "", "", fmt.Errorf("tag %s: %w", t, err)
		}
		if best == nil || version.Compare(v, *best) > 0 {
			best, bestTag = &v, t
		}
	}
	if best == nil {
		return "", "", nil
	}
	return best.Raw, bestTag, nil
}

func describeCommit(c changelog.Commit) versionCommit {
	vc := versionCommit{Hash: c.Hash, Subject: c.Subject}
	vc.Type, vc.Breaking = version.CommitType(c.Subject)
	for _, l := range changelog.FooterLines(c.Body) {
		k, v, _ := strings.Cut(l, ":")
		if vc.Footers == nil {
			vc.Footers = map[string]string{}
		}
		k = strings.TrimSpace(k)
		vc.Footers[k] = strings.TrimSpace(v)
		vc.Breaking = vc.Breaking || k == "BREAKING CHANGE"
	}
	return vc
}

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func printVersionState(w io.Writer, st *versionStateT) {
	fmt.Fprintf(w, "%s (%s)\n", st.Product, st.Scheme)
	if len(st.Sources) == 0 {
		fmt.Fprintln(w, "  sources    — (the version is the tag alone)")
	}
	for i, s := range st.Sources {
		label := "sources "
		if i > 0 {
			label = "        "
		}
		fmt.Fprintf(w, "  %s   %-10s %s\n", label, s.Version, s.File)
	}
	fmt.Fprintf(w, "  released   %-10s %s\n", orDash(st.Released), orDash(st.ReleasedTag))
	fmt.Fprintf(w, "  branch     %s → %s (%s)\n", st.Branch, st.Strategy, st.Rule)
	fmt.Fprintf(w, "  next       %-10s %s\n", orDash(st.Next), st.Reason)
	fmt.Fprintf(w, "  commits    %d since %s\n", len(st.Commits), orDash(st.ReleasedTag))
}

// bumpVersion writes the version into every source: the strategy's, a part raised from the
// sources', or the one given. The version is confronted before a byte is written.
func bumpVersion(w io.Writer, st *versionStateT, want string) error {
	p, s := st.product, st.scheme
	if len(p.Source) == 0 {
		return fmt.Errorf("%s has no source to write: its version is the tag alone", st.Product)
	}
	target, why := st.Next, st.Reason
	switch {
	case want == "":
		if target == "" {
			return fmt.Errorf("%s: %s", st.Branch, st.Reason)
		}
	case s.ValidPart(want) && want != "":
		cur, err := s.Parse(st.Current)
		if err != nil {
			return err
		}
		if target, err = s.Increment(cur, want, time.Now()); err != nil {
			return err
		}
		why = "the " + want + " of " + st.Current + " raised"
	default:
		target, why = want, "given"
	}
	v, err := s.Parse(target)
	if err != nil {
		return err
	}
	cur, err := s.Parse(st.Current)
	if err != nil {
		return err
	}
	switch c := version.Compare(v, cur); {
	case c == 0:
		fmt.Fprintf(w, "%s is already at %s (%s)\n", st.Product, target, why)
		return nil
	case c < 0:
		return fmt.Errorf("%s is below %s, the version the sources hold: a version only moves forward", target, st.Current)
	}
	if st.Released != "" {
		rel, _ := s.Parse(st.Released)
		if version.Compare(v, rel) <= 0 {
			return fmt.Errorf("%s is not above %s, the latest release (%s)", target, st.Released, st.ReleasedTag)
		}
	}
	type write struct {
		path string
		body []byte
	}
	var writes []write
	for _, src := range st.Sources {
		path := filepath.Join(st.root, filepath.FromSlash(src.File))
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out, err := p.Rewrite(src.File, b, target)
		if err != nil {
			return err
		}
		writes = append(writes, write{path, out})
	}
	for _, wr := range writes {
		info, err := os.Stat(wr.path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(wr.path, wr.body, info.Mode().Perm()); err != nil {
			return err
		}
	}
	fmt.Fprintf(w, "%s %s → %s (%s)\n", st.Product, st.Current, target, why)
	for _, src := range st.Sources {
		fmt.Fprintf(w, "  %s\n", src.File)
	}
	return nil
}
