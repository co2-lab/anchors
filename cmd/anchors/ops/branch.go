// @anchors
//   code: BRCMB
//   ref: BCSBR

package ops

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/spf13/cobra"
)

// The branch flow is confronted where the merge happens. A merge into a guarded branch goes
// through a pull request, which names both sides, so the pull request is checked by its pair;
// a push that reaches a guarded branch is checked by the pull request that made it, because a
// repository with no server-side protection takes a direct push (reported from a project: a
// private repository on GitHub's free plan applies no ruleset). On the machine, a guarded
// branch takes no commit of its own: the hooks refuse it.

func newBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Confront merges and commits with the project's branch flow (`branch_flow:`)",
		Long: `Confronts the branch flow anchors.yaml declares:

  branch_flow:
    merges:
      review:     [develop, "hotfix/*"]
      production: [review, "hotfix/*"]
    protected: [review, production]

  anchors branch check --head develop --base review   a pull request's pair
  anchors branch check --push --ref production --sha <sha>
                                                        a push: the merged pull request that made it
  anchors branch guard                                   the hooks: no commit on a protected branch

In a CI, check reads its sides from the environment: GITHUB_EVENT_NAME, GITHUB_HEAD_REF and
GITHUB_BASE_REF for a pull request, GITHUB_REF_NAME and GITHUB_SHA for a push. A push asks
GitHub, with the gh CLI (GH_TOKEN), for the merged pull request whose merge commit is the
sha.`,
	}
	cmd.AddCommand(newBranchCheckCmd(), newBranchGuardCmd())
	return cmd
}

// pullRequest is what a push is confronted with: the merged pull request that made a commit.
type pullRequest struct {
	Number         int    `json:"number"`
	MergedAt       string `json:"merged_at"`
	MergeCommitSHA string `json:"merge_commit_sha"`
	Head           struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

// pullsAPI asks GitHub's API for a list of pull requests of the repository at root.
var pullsAPI = func(root, path string) ([]byte, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return nil, errors.New("the gh CLI is not in the PATH: a push is confronted with the pull request GitHub knows")
	}
	cmd := exec.Command("gh", "api", path)
	cmd.Dir = root
	return cmd.Output()
}

// pullsOf lists the pull requests GitHub associates with a commit. When none of them was merged
// as that commit, the recently closed pull requests are read for the one whose merge commit it
// is: the association is an index GitHub fills apart from the merge, and it can lag or never
// come. Measured on a project: a pull request merged between 502s answered `merged: true` with its
// merge commit, and `commits/{sha}/pulls` returned `[]` past fifteen minutes — the release
// stopped at the guard.
var pullsOf = func(root, sha string) ([]pullRequest, error) {
	prs, err := readPulls(root, "repos/{owner}/{repo}/commits/"+sha+"/pulls", sha)
	if err != nil {
		return nil, err
	}
	for _, pr := range prs {
		if pr.MergeCommitSHA == sha {
			return prs, nil
		}
	}
	recent, err := readPulls(root, "repos/{owner}/{repo}/pulls?state=closed&sort=updated&direction=desc&per_page=100", sha)
	if err != nil {
		return prs, nil
	}
	for _, pr := range recent {
		if pr.MergeCommitSHA == sha {
			prs = append(prs, pr)
		}
	}
	return prs, nil
}

func readPulls(root, path, sha string) ([]pullRequest, error) {
	out, err := pullsAPI(root, path)
	if err != nil {
		return nil, fmt.Errorf("ask GitHub for the pull requests of %s: %w", sha, err)
	}
	var prs []pullRequest
	if err := json.Unmarshal(out, &prs); err != nil {
		return nil, fmt.Errorf("read the pull requests of %s: %w", sha, err)
	}
	return prs, nil
}

func newBranchCheckCmd() *cobra.Command {
	var root, head, base, ref, sha, headSHA string
	var push bool
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Refuse a pull request, or a push, the branch flow does not allow",
		Long: `On a pull request, refuses a head the flow does not take into the base, and a head that
does not meet what the stages of the base ask of its version and of the branches it contains
(the version: and contains: requirements of pipeline:). The head commit comes from --head-sha,
or the event's pull_request.head.sha.

On a push, refuses a push to a guarded branch that is not the merge of an allowed pull request.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			if !push && head == "" && base == "" {
				push = os.Getenv("GITHUB_EVENT_NAME") == "push"
			}
			var msg string
			if push {
				msg, err = checkPush(absRoot, cfg.BranchFlow, or(ref, os.Getenv("GITHUB_REF_NAME")), or(sha, os.Getenv("GITHUB_SHA")))
			} else {
				base = or(base, os.Getenv("GITHUB_BASE_REF"))
				if msg, err = checkMerge(cfg.BranchFlow, or(head, os.Getenv("GITHUB_HEAD_REF")), base); err == nil {
					var stages string
					stages, err = checkStages(absRoot, cfg, base, or(headSHA, eventHeadSHA()))
					msg += stages
				}
			}
			if msg != "" {
				fmt.Println(msg)
			}
			return err
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&head, "head", "", "the pull request's head branch (default $GITHUB_HEAD_REF)")
	cmd.Flags().StringVar(&base, "base", "", "the pull request's base branch (default $GITHUB_BASE_REF)")
	cmd.Flags().BoolVar(&push, "push", false, "confront a push instead of a pull request (default when $GITHUB_EVENT_NAME is push)")
	cmd.Flags().StringVar(&ref, "ref", "", "the pushed branch (default $GITHUB_REF_NAME)")
	cmd.Flags().StringVar(&sha, "sha", "", "the pushed commit (default $GITHUB_SHA)")
	cmd.Flags().StringVar(&headSHA, "head-sha", "", "the pull request's head commit, which the stages of the base confront (default: the event's pull_request.head.sha)")
	return cmd
}

// eventHeadSHA is the head commit of the pull request the CI's event names. $GITHUB_SHA is
// not it: on a pull request it is the merge GitHub tries, which already has the base in it.
func eventHeadSHA() string {
	b, err := os.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		return ""
	}
	var ev struct {
		PullRequest struct {
			Head struct {
				SHA string `json:"sha"`
			} `json:"head"`
		} `json:"pull_request"`
	}
	if json.Unmarshal(b, &ev) != nil {
		return ""
	}
	return ev.PullRequest.Head.SHA
}

// checkStages confronts a pull request's head with what the stages of its base ask of the
// version and of the branches it contains, so a refusal shows on the pull request and not
// later in the stage. same_as_passed is left to the stage: nothing has passed on a pull
// request yet. The other requirements — the flow, the gates, a stage passed, the head it
// comes from — are about the merge, which has not happened.
func checkStages(root string, cfg *config.Config, base, sha string) (string, error) {
	var b strings.Builder
	var failed []string
	for _, name := range cfg.Pipeline.StageNames() {
		st := cfg.Pipeline.Stages[name]
		on := len(st.Branches) == 0
		for _, p := range st.Branches {
			on = on || config.BranchMatches(p, base)
		}
		if !on {
			continue
		}
		for _, r := range st.Requires {
			var ok bool
			var why string
			switch {
			case r.Contains != "":
				if sha == "" {
					return b.String(), errHeadSHA
				}
				ok, why = containsHolds(root, r.Contains, sha)
			case r.Version != nil:
				vr := *r.Version
				vr.SameAsPassed = ""
				if !vr.Fresh && vr.Listed == nil && !vr.Expected {
					continue
				}
				if sha == "" {
					return b.String(), errHeadSHA
				}
				ok, why = versionHolds(root, cfg.Versioning, &vr, base, sha, nil)
			default:
				continue
			}
			mark := "✓"
			if !ok {
				mark = "✗"
				failed = append(failed, name+": "+r.String())
			}
			fmt.Fprintf(&b, "\n%s %s: %s: %s", mark, name, r, why)
		}
	}
	if len(failed) > 0 {
		return b.String(), fmt.Errorf("✗ the head (%s) does not meet what the stages of %s ask: %s", short(sha), base, strings.Join(failed, "; "))
	}
	return b.String(), nil
}

var errHeadSHA = errors.New("name the pull request's head commit with --head-sha: the stages of the base confront it, and $GITHUB_SHA is the trial merge")

func or(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}

// checkMerge says whether a pull request of head into base is one the flow allows.
func checkMerge(f *config.BranchFlow, head, base string) (string, error) {
	if head == "" || base == "" {
		return "", errors.New("name the pull request's sides: --head and --base, or $GITHUB_HEAD_REF and $GITHUB_BASE_REF")
	}
	ok, heads := f.MergeAllowed(head, base)
	if !ok {
		return "", fmt.Errorf("✗ %s → %s: the branch flow does not allow it — %s takes a merge from %s",
			head, base, base, takes(heads))
	}
	return fmt.Sprintf("✓ %s → %s: allowed by the branch flow", head, base), nil
}

// checkPush says whether a push to ref is the merge of a pull request the flow allows. A branch
// the flow does not guard takes any push.
func checkPush(root string, f *config.BranchFlow, ref, sha string) (string, error) {
	if ref == "" || sha == "" {
		return "", errors.New("name the push: --ref and --sha, or $GITHUB_REF_NAME and $GITHUB_SHA")
	}
	if !guards(f, ref) {
		return fmt.Sprintf("✓ %s: the branch flow does not guard it", ref), nil
	}
	pr, refused, err := mergingPR(root, f, ref, sha)
	if err != nil {
		return "", err
	}
	if pr != nil {
		return fmt.Sprintf("✓ %s: %s is the merge of #%d, %s → %s, allowed by the branch flow",
			ref, short(sha), pr.Number, pr.Head.Ref, ref), nil
	}
	_, heads := f.MergeAllowed("", ref)
	if len(refused) > 0 {
		return "", fmt.Errorf("✗ %s: %s is the merge of %s, and %s takes a merge from %s",
			ref, short(sha), strings.Join(refused, ", "), ref, takes(heads))
	}
	return "", fmt.Errorf("✗ %s: %s is the merge of no pull request — a guarded branch takes only a merged pull request, from %s",
		ref, short(sha), takes(heads))
}

// mergingPR is the merged pull request into ref whose merge commit is sha and whose head the
// flow allows, or nil with the merged pull requests the flow refuses.
func mergingPR(root string, f *config.BranchFlow, ref, sha string) (*pullRequest, []string, error) {
	prs, err := pullsOf(root, sha)
	if err != nil {
		return nil, nil, err
	}
	var refused []string
	for i, pr := range prs {
		if pr.MergedAt == "" || pr.Base.Ref != ref || pr.MergeCommitSHA != sha {
			continue
		}
		if ok, _ := f.MergeAllowed(pr.Head.Ref, ref); ok {
			return &prs[i], nil, nil
		}
		refused = append(refused, fmt.Sprintf("#%d from %s", pr.Number, pr.Head.Ref))
	}
	return nil, refused, nil
}

// guards says whether the flow rules a branch: an entry names it, or it is protected.
func guards(f *config.BranchFlow, branch string) bool {
	if f == nil {
		return false
	}
	if f.IsProtected(branch) {
		return true
	}
	for base := range f.Merges {
		if config.BranchMatches(base, branch) {
			return true
		}
	}
	return false
}

func takes(heads []string) string {
	if len(heads) == 0 {
		return "no branch"
	}
	return strings.Join(heads, ", ")
}

func short(sha string) string {
	if len(sha) > 8 {
		return sha[:8]
	}
	return sha
}

func newBranchGuardCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "guard",
		Short: "Refuse a commit on a branch the flow protects (the pre-commit and pre-merge-commit hooks)",
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			out, _ := gitOutput(absRoot, "symbolic-ref", "--quiet", "--short", "HEAD")
			return guardCommit(cfg.BranchFlow, strings.TrimSpace(out))
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	return cmd
}

// guardCommit refuses a commit on a protected branch. A detached HEAD — a rebase, a bisect —
// is on no branch, and commits freely.
func guardCommit(f *config.BranchFlow, branch string) error {
	if !f.IsProtected(branch) {
		return nil
	}
	_, heads := f.MergeAllowed("", branch)
	return fmt.Errorf("✗ %s is protected by the branch flow: it takes no commit of its own, only a merged pull request from %s",
		branch, takes(heads))
}
