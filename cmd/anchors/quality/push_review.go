// @anchors
//   code: PRCMD
//   ref: PRCMS

package quality

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/pushreview"
	"github.com/spf13/cobra"
)

// newPushReviewCmd records the review of what a push sends, and — as `verify`, from the
// pre-push hook — refuses a push whose evidence is missing (`push_review:`).
//
// The commit hook runs what a script computes; the judgment gates and the review that reads
// what a change means ran only when someone remembered them. The push is the last moment the
// work is still on the machine, so the evidence is asked there: the agent pushing reviews,
// runs the judgment gates while it does, and records the verdict in a file the repository
// keeps, the way the changelog is kept.
func newPushReviewCmd() *cobra.Command {
	var root, ok, fail, by, rangeArg string
	var notes []string
	var noCommit bool
	cmd := &cobra.Command{
		Use:   "push-review",
		Short: "Record the review of what a push sends; the pre-push hook confronts it",
		Long: `Records the review of the commits a push sends, in the file ` + "`push_review.file`" + ` names
(REVIEWS.md by default), and commits it. With ` + "`push_review.required`" + `, the pre-push hook refuses a
commit no review covers, or whose latest review failed; with ` + "`push_review.require_judgment`" + `, one
that reaches a judgment gate no verdict answers in the map it pushes.

  anchors push-review --ok "<what was validated, and how>" --by <who> [--note "<observation>"]...
  anchors push-review --fail "<what failed>" --by <who> [--note "<observation>"]...
  anchors push-review verify            what the pre-push hook runs

The range is what the branch has beyond its upstream (or beyond every remote branch), unless
--range <base>..<tip> says otherwise. Commits that change only the reviews' file and the map
owe no review. ` + "`anchors guide review`" + ` is the ruler for the review itself.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if (ok == "") == (fail == "") {
				return errors.New(i18n.T("push_review.need_verdict"))
			}
			if strings.TrimSpace(by) == "" {
				return errors.New(i18n.T("push_review.need_by"))
			}
			absRoot, cfg, err := loadPushReviewProject(root)
			if err != nil {
				return err
			}
			verdict, text := pushreview.OK, ok
			if fail != "" {
				verdict, text = pushreview.Fail, fail
			}
			revs, label := defaultPushRange(absRoot)
			if rangeArg != "" {
				revs, label = []string{rangeArg}, rangeArg
			}
			file := cfg.PushReview.FileOrDefault()
			commits, files, err := owedCommits(absRoot, file, revs)
			if err != nil {
				return err
			}
			if len(commits) == 0 {
				return errors.New(i18n.T("push_review.nothing_to_review", label))
			}
			answered, pending := judgmentsOver(absRoot, cfg, nil, files)
			if verdict == pushreview.OK && cfg.PushReview != nil && cfg.PushReview.RequireJudgment && len(pending) > 0 {
				return errors.New(i18n.T("push_review.judgment_pending_record", len(pending), bulletList(pending), strings.Join(files, " ")))
			}
			judgment := append(answered, suffixAll(pending, ": "+i18n.T("push_review.awaiting"))...)
			r := pushreview.Review{At: time.Now(), Verdict: verdict, By: by, Range: label, Commits: commits, Text: text, Notes: notes, Judgment: judgment}
			path := filepath.Join(absRoot, filepath.FromSlash(file))
			old, _ := os.ReadFile(path)
			entries := append([]pushreview.Entry{r.Entry()}, pushreview.Parse(string(old))...)
			entries = pushreview.Keep(entries, cfg.PushReview.RetentionDays(), r.At)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(path, []byte(pushreview.Render(entries)), 0o644); err != nil {
				return err
			}
			fmt.Println(i18n.T("push_review.recorded", verdict, len(commits), label, file))
			if noCommit {
				return nil
			}
			msg := fmt.Sprintf("chore(review): %s %s", verdict, label)
			if out, err := exec.Command("git", "-C", absRoot, "add", "--", file).CombinedOutput(); err != nil {
				return fmt.Errorf("git add %s: %v %s", file, err, out)
			}
			c := exec.Command("git", "-C", absRoot, "commit", "-q", "-m", msg, "--", file)
			c.Stdout, c.Stderr = os.Stdout, os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("%s: %w", i18n.T("push_review.commit_failed", file), err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&ok, "ok", "", "the review passed: what was validated, and how")
	cmd.Flags().StringVar(&fail, "fail", "", "the review failed: what failed")
	cmd.Flags().StringVar(&by, "by", "", "who reviewed, as you name them (human:ana, agent:<vendor>/<model>)")
	cmd.Flags().StringArrayVar(&notes, "note", nil, "an observation (repeatable)")
	cmd.Flags().StringVar(&rangeArg, "range", "", "the commits reviewed, as <base>..<tip> (default: what the branch has beyond its upstream)")
	cmd.Flags().BoolVar(&noCommit, "no-commit", false, "write the review without committing it")
	cmd.AddCommand(newPushReviewVerifyCmd())
	return cmd
}

// newPushReviewVerifyCmd is what the pre-push hook runs: it reads the refs git hands the hook
// on stdin, and refuses the push whose commits miss the evidence `push_review:` asks.
func newPushReviewVerifyCmd() *cobra.Command {
	var root, rangeArg string
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Refuse a push whose commits miss their review or their judgments (the pre-push hook)",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, cfg, err := loadPushReviewProject(root)
			if err != nil {
				return err
			}
			if !cfg.PushReview.Enforced() {
				return nil
			}
			var pushes []pushRange
			if rangeArg != "" {
				pushes = []pushRange{{revs: []string{rangeArg}, tip: tipOf(rangeArg)}}
			} else if pushes, err = readPushRefs(cmd.InOrStdin()); err != nil {
				return err
			} else if pushes == nil {
				revs, _ := defaultPushRange(absRoot)
				pushes = []pushRange{{revs: revs, tip: "HEAD"}}
			}
			var problems []string
			for _, p := range pushes {
				problems = append(problems, confrontPush(absRoot, cfg, p)...)
			}
			if len(problems) == 0 {
				return nil
			}
			return errors.New(i18n.T("push_review.refused", strings.Join(problems, "\n\n")))
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&rangeArg, "range", "", "confront <base>..<tip> instead of the refs on stdin")
	return cmd
}

// pushRange is one ref a push sends: the rev-list arguments of its new commits, and its tip.
type pushRange struct {
	revs []string
	tip  string
}

// readPushRefs reads the lines git gives the pre-push hook — `<local ref> <local sha> <remote
// ref> <remote sha>` — into the commits each sends: after the remote's tip, or, for a new
// branch, those no remote branch has. A deletion sends nothing. Nil when there is no input.
func readPushRefs(in io.Reader) ([]pushRange, error) {
	if f, ok := in.(*os.File); ok {
		if st, err := f.Stat(); err == nil && st.Mode()&os.ModeCharDevice != 0 {
			return nil, nil
		}
	}
	var out []pushRange
	seen := false
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) != 4 {
			continue
		}
		seen = true
		local, remote := f[1], f[3]
		if strings.Trim(local, "0") == "" {
			continue
		}
		if strings.Trim(remote, "0") == "" {
			out = append(out, pushRange{revs: []string{local, "--not", "--remotes"}, tip: local})
			continue
		}
		out = append(out, pushRange{revs: []string{remote + ".." + local}, tip: local})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if !seen {
		return nil, nil
	}
	if out == nil {
		out = []pushRange{}
	}
	return out, nil
}

// confrontPush is what one ref misses: the commits no ok review covers, read from the reviews'
// file as the tip has it — the file the push sends —, and the judgments the map at the tip
// leaves unanswered.
func confrontPush(root string, cfg *config.Config, p pushRange) []string {
	file := cfg.PushReview.FileOrDefault()
	commits, files, err := owedCommits(root, file, p.revs)
	if err != nil {
		return []string{err.Error()}
	}
	if len(commits) == 0 {
		return nil
	}
	var problems []string
	if cfg.PushReview.Required {
		content, _ := gitmeta.AtCommit(root, p.tip, file)
		entries := pushreview.Parse(string(content))
		var missing, failed []string
		reported := map[string]bool{}
		for _, c := range commits {
			e, ok := pushreview.Covering(entries, c)
			switch {
			case !ok:
				missing = append(missing, short(c.Hash)+" "+c.Subject)
			case e.Verdict != pushreview.OK && !reported[e.Raw]:
				reported[e.Raw] = true
				failed = append(failed, strings.TrimSpace(e.Raw))
			}
		}
		if len(missing) > 0 {
			problems = append(problems, i18n.T("push_review.missing", len(missing), bulletList(missing), file))
		}
		for _, f := range failed {
			problems = append(problems, i18n.T("push_review.failed", f))
		}
	}
	if cfg.PushReview.RequireJudgment {
		var g *mapx.Graph
		if b, ok := gitmeta.AtCommit(root, p.tip, mapx.DefaultPath); ok {
			g, _ = mapx.LoadBytes(b, mapx.DefaultPath)
		}
		if _, pending := judgmentsOver(root, cfg, g, files); len(pending) > 0 {
			problems = append(problems, i18n.T("push_review.judgment_pending", len(pending), bulletList(pending), strings.Join(files, " ")))
		}
	}
	return problems
}

// owedCommits are the commits of the range that owe a review — those that change more than
// the reviews' file and the map —, and the files they change, sorted.
func owedCommits(root, file string, revs []string) ([]pushreview.Commit, []string, error) {
	all, err := pushreview.Commits(root, revs...)
	if err != nil {
		return nil, nil, err
	}
	var owed []pushreview.Commit
	seen := map[string]bool{}
	var files []string
	for _, c := range all {
		fs := pushreview.Files(root, c.Hash)
		if pushreview.Exempt(fs, file, mapx.DefaultPath) {
			continue
		}
		owed = append(owed, c)
		for _, f := range fs {
			if !seen[f] && f != file && f != mapx.DefaultPath {
				seen[f] = true
				files = append(files, f)
			}
		}
	}
	sort.Strings(files)
	return owed, files, nil
}

// judgmentsOver runs the judgment gates over what the files reach, in g — the map on disk
// when nil —: those a verdict answers, as `gate @ target: verdict`, and those still awaiting
// one, as `gate @ target`.
func judgmentsOver(root string, cfg *config.Config, g *mapx.Graph, files []string) (answered, pending []string) {
	if g == nil {
		var err error
		if g, err = mapx.Load(filepath.Join(root, mapx.DefaultPath)); err != nil {
			return nil, nil
		}
	}
	var gates []config.Gate
	for _, gc := range cfg.Gates {
		if gc.IsJudgment() {
			gates = append(gates, gc)
		}
	}
	if len(gates) == 0 || len(files) == 0 {
		return nil, nil
	}
	var present []string
	for _, f := range files {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(f))); err == nil {
			present = append(present, f)
		}
	}
	nodes, _, err := selectNodes(g, cfg, false, present, root)
	if err != nil || len(nodes) == 0 {
		return nil, nil
	}
	for _, r := range gate.RunWithWaiver(gates, nodes, root, g, cfg, false, gate.Waiver{}) {
		at := r.Gate + " @ " + r.Target
		switch r.Verdict {
		case gate.Judge:
			pending = append(pending, at)
		case gate.Pass:
			answered = append(answered, at+": ok")
		case gate.Fail:
			answered = append(answered, at+": issue")
		}
	}
	sort.Strings(answered)
	sort.Strings(pending)
	return answered, pending
}

// defaultPushRange is what the branch has beyond its upstream, or — with none — beyond every
// remote branch, as rev-list arguments and as the range a review names.
func defaultPushRange(root string) ([]string, string) {
	head, _ := gitmeta.Commit(root, "HEAD")
	if up, ok := gitmeta.Commit(root, "@{upstream}"); ok {
		return []string{up + "..HEAD"}, short(up) + ".." + short(head)
	}
	return []string{"HEAD", "--not", "--remotes"}, short(head)
}

func loadPushReviewProject(root string) (string, *config.Config, error) {
	absRoot, err := config.AbsRoot(root)
	if err != nil {
		return "", nil, err
	}
	cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
	if err != nil {
		return "", nil, fmt.Errorf("load config: %w", err)
	}
	return absRoot, cfg, nil
}

func tipOf(r string) string {
	if i := strings.LastIndex(r, ".."); i >= 0 {
		return r[i+2:]
	}
	return r
}

func short(h string) string {
	if len(h) > 8 {
		return h[:8]
	}
	return h
}

func bulletList(items []string) string {
	return "  - " + strings.Join(items, "\n  - ")
}

func suffixAll(items []string, suffix string) []string {
	out := make([]string, len(items))
	for i, s := range items {
		out[i] = s + suffix
	}
	return out
}
