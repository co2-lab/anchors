// @anchors
//   code: VAPLV
//   ref: VAPLA

package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/pathref"
	"github.com/spf13/cobra"
)

// A version computed is a version no one wrote yet: it enters the repository the way the
// product's `apply:` says. By default through a pull request the CI opens, carrying the bump
// in every source, the changelog of the release and the intents it consumed: merging it is
// the release, nobody writes the version by hand, and the CI never pushes to a protected
// branch. A project that bumps by hand has the stage confirm it; one versioned by its tags
// alone has the version written into the build's workspace.

// prAPI runs the gh CLI for the release's pull request.
var prAPI = func(root string, args ...string) ([]byte, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return nil, errors.New("the gh CLI is not in the PATH: the release pull request is opened through it")
	}
	c := exec.Command("gh", args...)
	c.Dir = root
	out, err := c.Output()
	if ee, ok := err.(*exec.ExitError); ok {
		return out, fmt.Errorf("gh %s: %v: %s", args[0]+" "+args[1], err, strings.TrimSpace(string(ee.Stderr)))
	}
	return out, err
}

func newVersionApplyCmd(root, product, branch *string) *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Bring the version the branch asks for into the repository, as the product's apply: says",
		Long: `apply: release_pr (the default) opens, or updates, the pull request anchors/release/<product> →
the branch, carrying the version in every source, the release's block in every changelog output
and the intents it consumed. Merging it is the release.

apply: commit confirms the sources already hold the version the branch asks for: a bump made by
hand, with anchors version bump.

apply: build writes the version into the sources of the workspace, commits nothing, and prints
version= and tag= (also into $GITHUB_OUTPUT): the tag alone records it.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := versionState(*root, *product, *branch)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(st.root, config.DefaultFile))
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			target := st.Next
			if st.Label != "" && st.LabelIn != "tag" {
				target = st.Label
			}
			if target == "" {
				return fmt.Errorf("%s: %s", st.Branch, st.Reason)
			}
			switch st.product.ApplyOrDefault() {
			case "commit":
				if st.Current != target {
					return fmt.Errorf("✗ the sources hold %s and %s asks for %s (%s): run anchors version bump", st.Current, st.Branch, target, st.Reason)
				}
				fmt.Fprintf(w, "✓ %s is at %s, as %s asks (%s)\n", st.Product, target, st.Branch, st.Reason)
				return nil
			case "build":
				return applyBuild(w, st, target)
			}
			return applyReleasePR(w, cfg, st, target)
		},
	}
}

// applyBuild writes the version into the workspace's sources and says it.
func applyBuild(w io.Writer, st *versionStateT, target string) error {
	if len(st.Sources) > 0 && st.Current != target {
		if err := bumpVersion(w, st, target); err != nil {
			return err
		}
	}
	tagged := target
	if st.Label != "" && st.LabelIn == "tag" {
		tagged = st.Label
	}
	tag := st.product.TagFor(tagged)
	fmt.Fprintf(w, "version=%s\ntag=%s\n", tagged, tag)
	if out := os.Getenv("GITHUB_OUTPUT"); out != "" {
		f, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := fmt.Fprintf(f, "version=%s\ntag=%s\n", tagged, tag); err != nil {
			return err
		}
	}
	return nil
}

// applyReleasePR prepares the release on a branch of its own and opens, or updates, its
// pull request into the branch. Nothing to release opens nothing.
func applyReleasePR(w io.Writer, cfg *config.Config, st *versionStateT, target string) error {
	if len(st.Sources) == 0 {
		return fmt.Errorf("%s has no source to bump: a product versioned by its tags alone uses apply: build", st.Product)
	}
	if target == st.Released {
		fmt.Fprintf(w, "%s: nothing to release since %s (%s)\n", st.Product, st.ReleasedTag, st.Reason)
		return nil
	}
	git := func(args ...string) (string, error) {
		c := exec.Command("git", append([]string{"-C", st.root}, args...)...)
		out, err := c.CombinedOutput()
		if err != nil {
			return string(out), fmt.Errorf("git %s: %v: %s", args[0], err, strings.TrimSpace(string(out)))
		}
		return string(out), nil
	}
	head := "anchors/release/" + st.Product
	if _, err := git("checkout", "-q", "-B", head); err != nil {
		return err
	}
	defer func() { _, _ = git("checkout", "-q", st.Branch) }()
	if st.Current != target {
		if err := bumpVersion(w, st, target); err != nil {
			return err
		}
	}
	touched := []string{}
	for _, s := range st.Sources {
		touched = append(touched, s.File)
	}
	body, files, err := writeRelease(st, cfg, target)
	if err != nil {
		return err
	}
	touched = append(touched, files...)
	consumed, err := consumeIntents(st, cfg, target)
	if err != nil {
		return err
	}
	touched = append(touched, consumed...)
	var add []string
	for _, f := range touched {
		_, statErr := os.Stat(filepath.Join(st.root, filepath.FromSlash(f)))
		tracked, _ := git("ls-files", "--", f)
		if statErr == nil || strings.TrimSpace(tracked) != "" {
			add = append(add, f)
		}
	}
	if len(add) > 0 {
		if _, err := git(append([]string{"add", "-A", "--"}, add...)...); err != nil {
			return err
		}
	}
	if _, err := git("add", "-u"); err != nil {
		return err
	}
	title := fmt.Sprintf("release %s %s", st.Product, target)
	if _, err := git("diff", "--cached", "--quiet"); err == nil {
		fmt.Fprintf(w, "%s %s: %s already carries the release, nothing to propose\n", st.Product, target, st.Branch)
		return nil
	}
	args := []string{"commit", "-q", "-m", "chore(release): " + st.Product + " " + target}
	if out, _ := git("config", "user.email"); strings.TrimSpace(out) == "" {
		args = append([]string{"-c", "user.name=anchors", "-c", "user.email=anchors@users.noreply.github.com"}, args...)
	}
	if _, err := git(args...); err != nil {
		return err
	}
	if _, err := git("push", "-q", "--force", "origin", head); err != nil {
		return err
	}
	out, err := prAPI(st.root, "pr", "list", "--head", head, "--base", st.Branch, "--state", "open", "--json", "number")
	if err != nil {
		return err
	}
	var open []struct {
		Number int `json:"number"`
	}
	if err := json.Unmarshal(out, &open); err != nil {
		return fmt.Errorf("read the open pull requests: %w", err)
	}
	if len(open) > 0 {
		n := fmt.Sprint(open[0].Number)
		if _, err := prAPI(st.root, "pr", "edit", n, "--title", title, "--body", body); err != nil {
			return err
		}
		fmt.Fprintf(w, "✓ #%s updated: %s\n", n, title)
		return nil
	}
	if _, err := prAPI(st.root, "pr", "create", "--base", st.Branch, "--head", head, "--title", title, "--body", body); err != nil {
		return err
	}
	fmt.Fprintf(w, "✓ opened: %s (%s → %s)\n", title, head, st.Branch)
	return nil
}

// writeRelease writes the release's block into every changelog output, replacing the block a
// previous run wrote for it, and returns the first block — the pull request's body — and the
// files written.
func writeRelease(st *versionStateT, cfg *config.Config, target string) (string, []string, error) {
	outputs := []config.ChangelogOutput{}
	if pc := st.product.Changelog; pc != nil {
		outputs = pc.Outputs
	}
	if len(outputs) == 0 && cfg.Changelog != nil {
		outputs = []config.ChangelogOutput{{Changelog: *cfg.Changelog}}
	}
	commits, err := changelog.Git{Root: st.root}.Commits(st.ReleasedTag, "HEAD")
	if err != nil {
		return "", nil, err
	}
	today := time.Now().UTC().Format("2006-01-02")
	rel := changelog.Classify(target, today, commits)
	var body string
	var files []string
	for i, o := range outputs {
		label := fmt.Sprintf("versioning.%s.changelog.outputs[%d].template", st.Product, i)
		var blocks []changelog.Block
		if o.FromIntents() {
			blocks, err = pendingIntentBlock(st, cfg, &o.Changelog, label, target, today)
		} else {
			blocks, err = renderBlocks(st.root, &o.Changelog, label, []changelog.Release{rel})
		}
		if err != nil {
			return "", nil, err
		}
		if len(blocks) == 0 {
			continue
		}
		if body == "" {
			body = blocks[0].Text
		}
		path := o.PathOrDefault()
		if o.ModeOrDefault() == "per_version" {
			name := target + ".md"
			if o.Format != nil && o.Format.File != "" {
				name = strings.ReplaceAll(o.Format.File, "{version}", target)
			}
			path = path + "/" + name
			if err := os.MkdirAll(filepath.Join(st.root, filepath.FromSlash(o.PathOrDefault())), 0o755); err != nil {
				return "", nil, err
			}
			_ = os.Remove(filepath.Join(st.root, filepath.FromSlash(o.PathOrDefault()), "unreleased.md"))
			files = append(files, o.PathOrDefault())
			if err := os.WriteFile(filepath.Join(st.root, filepath.FromSlash(path)), []byte(strings.TrimRight(blocks[0].Text, "\n")+"\n"), 0o644); err != nil {
				return "", nil, err
			}
		} else {
			abs := filepath.Join(st.root, filepath.FromSlash(path))
			existing, err := os.ReadFile(abs)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return "", nil, err
			}
			if err := os.WriteFile(abs, []byte(changelog.Replace(string(existing), i18n.T("changelog.title"), blocks[0])), 0o644); err != nil {
				return "", nil, err
			}
		}
		files = append(files, path)
	}
	return body, files, nil
}

// pendingIntentBlock renders the release from the product's intents waiting for it, headed
// as the release it becomes.
func pendingIntentBlock(st *versionStateT, cfg *config.Config, cl *config.Changelog, label, target, date string) ([]changelog.Block, error) {
	tmpl := changelog.DefaultIntentsTemplate
	if cl.Template != "" {
		b, err := pathref.Read(st.root, cl.Template)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		tmpl = string(b)
	}
	all, errs := changelog.ReadIntents(filepath.Join(st.root, filepath.FromSlash(intentsDir(cfg))))
	if len(errs) > 0 {
		return nil, fmt.Errorf("intents: %w", errs[0])
	}
	var mine []changelog.Intent
	for _, in := range all {
		if in.Product == st.Product {
			mine = append(mine, in)
		}
	}
	if len(mine) == 0 {
		return nil, nil
	}
	types := (*config.Intents)(nil).TypesOrDefault()
	if cfg.Changelog != nil {
		types = cfg.Changelog.Intents.TypesOrDefault()
	}
	text, err := changelog.RenderIntents(tmpl, changelog.GroupIntents(target, date, mine, types), typeHeading)
	if err != nil {
		return nil, err
	}
	return []changelog.Block{{Version: target, Text: text}}, nil
}

// consumeIntents takes the product's intents out of the waiting ones: deleted, or archived
// under released/<version>/ where the changelog reads them again. It returns the paths
// touched.
func consumeIntents(st *versionStateT, cfg *config.Config, target string) ([]string, error) {
	dir := intentsDir(cfg)
	abs := filepath.Join(st.root, filepath.FromSlash(dir))
	all, errs := changelog.ReadIntents(abs)
	if len(errs) > 0 {
		return nil, fmt.Errorf("intents: %w", errs[0])
	}
	archive := cfg.Changelog != nil && cfg.Changelog.Intents != nil && cfg.Changelog.Intents.Consume == "archive"
	var touched []string
	for _, in := range all {
		if in.Product != st.Product {
			continue
		}
		from := filepath.Join(abs, in.File)
		touched = append(touched, dir+"/"+in.File)
		if !archive {
			if err := os.Remove(from); err != nil {
				return nil, err
			}
			continue
		}
		to := filepath.Join(abs, "released", target, in.File)
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			return nil, err
		}
		if err := os.Rename(from, to); err != nil {
			return nil, err
		}
		touched = append(touched, dir+"/released/"+target+"/"+in.File)
	}
	return touched, nil
}
