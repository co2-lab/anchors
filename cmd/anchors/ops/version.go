// @anchors
//   code: VRSNH
//   ref: VRSNI

package ops

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/shell"
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
	Range       string          `json:"range,omitempty"`
	RangeRule   string          `json:"range_rule,omitempty"`
	Label       string          `json:"label,omitempty"`
	LabelIn     string          `json:"label_in,omitempty"`

	root    string
	product *config.Product
	scheme  version.Scheme
}

// versionState reads a product's state in the working tree: its sources, its latest release,
// the commits since it, the branch's strategy and the version it computes.
func versionState(root, name, branch string) (*versionStateT, error) {
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return nil, err
	}
	cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", config.DefaultFile, err)
	}
	if branch == "" {
		out, err := exec.Command("git", "-C", absRoot, "rev-parse", "--abbrev-ref", "HEAD").Output()
		branch = strings.TrimSpace(string(out))
		if err != nil || branch == "" || branch == "HEAD" {
			return nil, errors.New("no current branch (a detached HEAD, or no repository): name it with --branch")
		}
	}
	return stateAt(absRoot, cfg, name, "--product", branch, "")
}

// readAt reads a file as the commit rev has it, or from the working tree when rev is empty.
func readAt(root, rev string) func(string) ([]byte, error) {
	if rev == "" {
		return func(f string) ([]byte, error) { return os.ReadFile(filepath.Join(root, filepath.FromSlash(f))) }
	}
	return func(f string) ([]byte, error) {
		b, ok := gitmeta.AtCommit(root, rev, f)
		if !ok {
			return nil, fmt.Errorf("not in %s", short(rev))
		}
		return b, nil
	}
}

// stateAt reads a product's state as the commit rev has it — the working tree and HEAD when
// rev is empty —, for the branch whose strategy decides; `how` says where a product is named.
func stateAt(absRoot string, cfg *config.Config, name, how, branch, rev string) (*versionStateT, error) {
	name, p, err := cfg.Versioning.Pick(name, how)
	if err != nil {
		return nil, err
	}
	s, _ := p.SchemeOf()
	st := &versionStateT{Product: name, Scheme: s.Kind, Branch: branch, root: absRoot, product: p, scheme: s,
		Sources: []versionSource{}, Commits: []versionCommit{}}
	svs, err := p.ReadSources(readAt(absRoot, rev))
	if err != nil {
		return nil, err
	}
	for _, sv := range svs {
		st.Sources = append(st.Sources, versionSource{File: sv.File, Version: sv.Version})
	}
	if st.Current, err = config.Agreed(svs); err != nil {
		return nil, err
	}
	to := or(rev, "HEAD")
	g := changelog.Git{Root: absRoot}
	if st.Released, st.ReleasedTag, err = latestRelease(g, p, s, to); err != nil {
		return nil, err
	}
	commits, err := g.Commits(st.ReleasedTag, to)
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
	st.Range, st.RangeRule = p.RangeFor(branch)
	in := version.Input{Scheme: s, Strategy: strategy, Current: st.Current, Released: st.Released, Commits: commits, Today: time.Now()}
	if strategy.Name == version.Intents {
		intents, err := intentsAt(absRoot, cfg, rev)
		if err != nil {
			return nil, err
		}
		for _, it := range intents {
			if it.Product == name {
				in.Bumps = append(in.Bumps, it.Bump)
			}
		}
	}
	switch strategy.Name {
	case version.Promote:
		in.Decided, in.DecidedBy = promoted(absRoot, p, s, strategy.From)
	case version.Script:
		in.Decided, in.DecidedBy = scripted(absRoot, st, strategy.Run)
	}
	res, err := version.Next(in)
	if err != nil {
		return nil, fmt.Errorf("%s on %s (%s): %w", name, branch, strategy.Name, err)
	}
	st.Next, st.Reason, st.Forced = res.Version, res.Reason, res.Forced
	if l := p.LabelFor(branch); l != nil && st.Next != "" && st.Next != st.Released {
		st.Label, st.LabelIn = labelled(absRoot, p, st.Next, l.Name), l.InOrDefault()
	}
	return st, nil
}

// labelled is the next prerelease of a version: `1.1.0-rc.<n>`, n one above the highest the
// product's tags record.
func labelled(absRoot string, p *config.Product, v, name string) string {
	out, _ := exec.Command("git", "-C", absRoot, "tag", "--list").Output()
	return v + "-" + name + "." + strconv.Itoa(nextLabelNumber(p, strings.Fields(string(out)), v, name))
}

// nextLabelNumber is one above the highest n among the tags recording `<v>-<name>.<n>`.
func nextLabelNumber(p *config.Product, tags []string, v, name string) int {
	n := 0
	for _, t := range tags {
		raw, ok := p.VersionOfTag(t)
		if !ok || !strings.HasPrefix(raw, v+"-"+name+".") {
			continue
		}
		if k, err := strconv.Atoi(strings.TrimPrefix(raw, v+"-"+name+".")); err == nil && k > n {
			n = k
		}
	}
	return n + 1
}

// inRange confronts a version with the branch's range; true when the branch has none.
func (st *versionStateT) inRange(raw string) (bool, string, error) {
	if st.Range == "" {
		return true, "", nil
	}
	v, err := st.scheme.Parse(raw)
	if err != nil {
		return false, "", err
	}
	ok, why, err := config.InRange(st.scheme, st.Range, st.Released, v)
	return ok, why + " (" + st.RangeRule + ")", err
}

// latestRelease is the highest version among the product's tags reachable from `to`. A tag
// of the product the scheme cannot read is an error: skipped, a release would go unseen.
func latestRelease(g changelog.Git, p *config.Product, s version.Scheme, to string) (string, string, error) {
	tags, err := g.Tags(to)
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
		if len(v.Pre) > 0 {
			continue // a prerelease is not the release the next version moves from
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
	if st.Range != "" {
		fmt.Fprintf(w, "  range      %-10s %s\n", st.Range, st.RangeRule)
	}
	if st.Label != "" {
		fmt.Fprintf(w, "  label      %-10s in %s\n", st.Label, st.LabelIn)
	}
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
		if st.Label != "" && st.LabelIn != "tag" {
			target, why = st.Label, st.Reason+", labelled"
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
	if ok, why, err := st.inRange(target); err != nil || !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: the branch %s may not produce it", why, st.Branch)
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
	if p.AfterBump != "" {
		c, err := shell.Command(p.AfterBump)
		if err != nil {
			return err
		}
		c.Dir, c.Stdout, c.Stderr = st.root, w, w
		c.Env = append(os.Environ(), "ANCHORS_PRODUCT="+st.Product, "ANCHORS_VERSION="+target)
		if err := c.Run(); err != nil {
			return fmt.Errorf("after_bump `%s` failed (%v): the sources already hold %s", p.AfterBump, err, target)
		}
	}
	return nil
}

// intentsDir is the directory of the project's change intents.
func intentsDir(cfg *config.Config) string {
	if cfg.Changelog == nil {
		return (*config.Intents)(nil).PathOrDefault()
	}
	return cfg.Changelog.Intents.PathOrDefault()
}

// intentsAt reads the change intents waiting for a release — the files of the intents'
// directory, not those archived under it — as the commit rev has them, or the working tree
// when rev is empty. An intent that does not parse is an error: read as absent, its change
// would ship without moving the version.
func intentsAt(absRoot string, cfg *config.Config, rev string) ([]changelog.Intent, error) {
	dir := intentsDir(cfg)
	if rev == "" {
		in, errs := changelog.ReadIntents(filepath.Join(absRoot, filepath.FromSlash(dir)))
		if len(errs) > 0 {
			return nil, fmt.Errorf("intents: %w", errs[0])
		}
		return in, nil
	}
	out, err := exec.Command("git", "-C", absRoot, "ls-tree", "--name-only", rev, dir+"/").Output()
	if err != nil {
		return nil, fmt.Errorf("intents at %s: %w", short(rev), err)
	}
	var intents []changelog.Intent
	for _, f := range strings.Fields(string(out)) {
		if !strings.HasSuffix(f, ".md") {
			continue
		}
		b, ok := gitmeta.AtCommit(absRoot, rev, f)
		if !ok {
			continue
		}
		in, err := changelog.ParseIntent(path.Base(f), b)
		if err != nil {
			return nil, fmt.Errorf("intents: %w", err)
		}
		intents = append(intents, in)
	}
	return intents, nil
}

// promoted is the version the branch `from` carries at its head — on origin when the
// checkout has it, else the local branch —: what its sources hold, or for a product with no
// source its latest release there. Empty with the reason when it cannot be read.
func promoted(absRoot string, p *config.Product, s version.Scheme, from string) (string, string) {
	ref := ""
	for _, r := range []string{"refs/remotes/origin/" + from, "refs/heads/" + from} {
		if _, ok := gitmeta.Commit(absRoot, r); ok {
			ref = r
			break
		}
	}
	if ref == "" {
		return "", "no branch " + from + " here, nor on origin"
	}
	head, _ := gitmeta.Commit(absRoot, ref)
	if len(p.Source) == 0 {
		v, tag, err := latestRelease(changelog.Git{Root: absRoot}, p, s, ref)
		if err != nil || v == "" {
			return "", fmt.Sprintf("%s (%s) has no release to promote", from, short(head))
		}
		return v, fmt.Sprintf("the release of %s, %s", from, tag)
	}
	v, err := versionAt(absRoot, p, head)
	if err != nil {
		return "", fmt.Sprintf("the version of %s (%s): %v", from, short(head), err)
	}
	return v, fmt.Sprintf("the version %s carries (%s)", from, short(head))
}

// scripted is the version the project's command prints, given the product's state as
// `anchors version --json` prints it, on its stdin. A command that fails or prints nothing
// decides nothing: falling back to a default would publish a version no one chose.
func scripted(absRoot string, st *versionStateT, run string) (string, string) {
	ctx, err := json.Marshal(st)
	if err != nil {
		return "", err.Error()
	}
	c, err := shell.Command(run)
	if err != nil {
		return "", fmt.Sprintf("`%s`: %v", run, err)
	}
	var out, errb bytes.Buffer
	c.Dir, c.Stdin, c.Stdout, c.Stderr = absRoot, bytes.NewReader(ctx), &out, &errb
	if err := c.Run(); err != nil {
		return "", fmt.Sprintf("`%s` failed (%v): %s", run, err, strings.TrimSpace(errb.String()))
	}
	v := strings.TrimSpace(strings.SplitN(strings.TrimSpace(out.String()), "\n", 2)[0])
	if v == "" {
		return "", fmt.Sprintf("`%s` printed no version", run)
	}
	return v, "`" + run + "`"
}
