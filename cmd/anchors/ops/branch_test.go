// @anchors
//   code: BCTBA
//   ref: BCSBR

package ops

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func flowOf() *config.BranchFlow {
	return &config.BranchFlow{
		Merges: map[string][]string{
			"review":     {"develop", "hotfix/*"},
			"production": {"review", "hotfix/*"},
		},
		Protected: []string{"review", "production"},
	}
}

func TestBranchCheck_aPullRequestByItsPair(t *testing.T) {
	t.Run("BCSBR-B01: A pull request is allowed when the flow allows its head into its base, and refused naming what the base takes", func(t *testing.T) {})
	t.Run("BCSBR-E01: A check with a side missing fails naming what is missing", func(t *testing.T) {})
	if msg, err := checkMerge(flowOf(), "develop", "review"); err != nil || !strings.Contains(msg, "develop → review") {
		t.Errorf("develop → review is allowed: %q %v", msg, err)
	}
	_, err := checkMerge(flowOf(), "develop", "production")
	if err == nil || !strings.Contains(err.Error(), "production takes a merge from review, hotfix/*") {
		t.Errorf("develop → production is refused naming what production takes: %v", err)
	}
	if _, err := checkMerge(flowOf(), "", "review"); err == nil || !strings.Contains(err.Error(), "--head and --base") {
		t.Errorf("a pull request with no head names what is missing: %v", err)
	}
	if _, err := checkPush("", flowOf(), "production", ""); err == nil || !strings.Contains(err.Error(), "--ref and --sha") {
		t.Errorf("a push with no sha names what is missing: %v", err)
	}
}

// pulls answers the pull requests of a commit with the given ones.
func pulls(t *testing.T, prs ...pullRequest) {
	t.Helper()
	old := pullsOf
	pullsOf = func(string, string) ([]pullRequest, error) { return prs, nil }
	t.Cleanup(func() { pullsOf = old })
}

func pr(number int, head, base, merge string) pullRequest {
	p := pullRequest{Number: number, MergedAt: "2026-11-06T10:00:00Z", MergeCommitSHA: merge}
	p.Head.Ref, p.Base.Ref = head, base
	return p
}

func TestBranchCheck_aPushByThePullRequestThatMadeIt(t *testing.T) {
	t.Run("BCSBR-B02: A push to a guarded branch is allowed only as the merge of a pull request the flow allows", func(t *testing.T) {})
	t.Run("BCSBR-E02: A push that cannot be confronted fails", func(t *testing.T) {})
	const sha = "abcdef0123456789"
	pulls(t, pr(7, "review", "production", sha))
	if msg, err := checkPush("", flowOf(), "production", sha); err != nil || !strings.Contains(msg, "#7, review → production") {
		t.Errorf("the merge of review into production is allowed: %q %v", msg, err)
	}
	pulls(t, pr(8, "develop", "production", sha))
	if _, err := checkPush("", flowOf(), "production", sha); err == nil || !strings.Contains(err.Error(), "#8 from develop") {
		t.Errorf("the merge of a pull request the flow refuses is refused: %v", err)
	}
	open := pr(9, "review", "production", sha)
	open.MergedAt = ""
	pulls(t, open, pr(10, "review", "production", "another"), pr(11, "review", "review", sha))
	if _, err := checkPush("", flowOf(), "production", sha); err == nil || !strings.Contains(err.Error(), "merge of no pull request") {
		t.Errorf("a direct push — no merged pull request into this branch whose merge is the sha — is refused: %v", err)
	}
	if msg, err := checkPush("", flowOf(), "feat/x", sha); err != nil || !strings.Contains(msg, "does not guard") {
		t.Errorf("a branch the flow does not guard takes any push: %q %v", msg, err)
	}
	old := pullsOf
	pullsOf = func(string, string) ([]pullRequest, error) { return nil, errors.New("no gh") }
	defer func() { pullsOf = old }()
	if _, err := checkPush("", flowOf(), "production", sha); err == nil {
		t.Error("a push that cannot be confronted is refused, not passed")
	}
}

func TestBranchGuard_noCommitOnAProtectedBranch(t *testing.T) {
	t.Run("BCSBR-B03: A commit on a protected branch is refused; another branch, or a detached HEAD, commits", func(t *testing.T) {})
	if err := guardCommit(flowOf(), "production"); err == nil || !strings.Contains(err.Error(), "production is protected") {
		t.Errorf("a commit on production is refused: %v", err)
	}
	for _, b := range []string{"develop", "hotfix/x", ""} {
		if err := guardCommit(flowOf(), b); err != nil {
			t.Errorf("%q commits: %v", b, err)
		}
	}
	if err := guardCommit(nil, "production"); err != nil {
		t.Errorf("no flow declared protects nothing: %v", err)
	}
}

// branchProject is a governed project whose anchors.yaml holds the given branch_flow.
func branchProject(t *testing.T, flow string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte("version: 8\nbranch_flow:\n"+flow), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func runBranch(t *testing.T, root string, args ...string) (string, error) {
	t.Helper()
	cmd := newBranchCmd()
	cmd.SetArgs(append(args, "--root", root))
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	var err error
	out := captureStdout(t, func() { err = cmd.Execute() })
	return out, err
}

func TestBranchCmd_readsTheCIsEnvironment(t *testing.T) {
	t.Run("BCSBR-B01: A pull request is allowed when the flow allows its head into its base, and refused naming what the base takes", func(t *testing.T) {})
	root := branchProject(t, "  merges:\n    production: [review]\n  protected: [production]\n")
	t.Setenv("GITHUB_EVENT_NAME", "pull_request")
	t.Setenv("GITHUB_HEAD_REF", "review")
	t.Setenv("GITHUB_BASE_REF", "production")
	if out, err := runBranch(t, root, "check"); err != nil || !strings.Contains(out, "review → production") {
		t.Errorf("the pull request's sides come from the environment: %q %v", out, err)
	}
	if _, err := runBranch(t, root, "check", "--head", "develop"); err == nil {
		t.Error("a flag wins over the environment, and develop → production is refused")
	}
	t.Setenv("GITHUB_EVENT_NAME", "push")
	t.Setenv("GITHUB_REF_NAME", "production")
	t.Setenv("GITHUB_SHA", "abc")
	pulls(t, pr(3, "review", "production", "abc"))
	if out, err := runBranch(t, root, "check"); err != nil || !strings.Contains(out, "#3") {
		t.Errorf("a push event checks the push: %q %v", out, err)
	}
}

func TestBranchCmd_aConfigurationThatDoesNotLoad(t *testing.T) {
	t.Run("BCSBR-E03: A configuration that does not load fails the check and the guard", func(t *testing.T) {})
	root := branchProject(t, "  protected: [\"[\"]\n")
	for _, args := range [][]string{{"check", "--head", "a", "--base", "b"}, {"guard"}} {
		if _, err := runBranch(t, root, args...); err == nil || !strings.Contains(err.Error(), "branch_flow") {
			t.Errorf("%s fails with the load error: %v", args[0], err)
		}
	}
}

func TestBranchCmd_guardReadsTheCurrentBranch(t *testing.T) {
	t.Run("BCSBR-B03: A commit on a protected branch is refused; another branch, or a detached HEAD, commits", func(t *testing.T) {})
	root := branchProject(t, "  protected: [production]\n")
	if out, err := runGit(root, "init", "-q", "-b", "production"); err != nil {
		t.Skipf("git init: %v %s", err, out)
	}
	if _, err := runBranch(t, root, "guard"); err == nil || !strings.Contains(err.Error(), "production is protected") {
		t.Errorf("the guard reads the current branch: %v", err)
	}
	if _, err := runGit(root, "checkout", "-q", "-b", "develop"); err != nil {
		t.Fatal(err)
	}
	if _, err := runBranch(t, root, "guard"); err != nil {
		t.Errorf("develop commits: %v", err)
	}
}

// The association of a commit to its pull request is an index GitHub fills apart from the merge,
// and it can lag or never come: measured on a project, a pull request merged between 502s answered its
// merge commit, and `commits/{sha}/pulls` returned [] past fifteen minutes.
func TestBranchCheck_thePullRequestFoundWhenGitHubDoesNotAssociateIt(t *testing.T) {
	t.Run("BCSBR-B04: A merge GitHub does not associate with its commit is found among the closed pull requests", func(t *testing.T) {})
	merged := pr(7, "review", "production", "3c1b7fcbe17a")
	other := pr(6, "develop", "review", "aaaa")
	asked := []string{}
	old := pullsAPI
	pullsAPI = func(_, path string) ([]byte, error) {
		asked = append(asked, path)
		if strings.Contains(path, "/commits/") {
			return []byte("[]"), nil
		}
		return json.Marshal([]pullRequest{other, merged})
	}
	t.Cleanup(func() { pullsAPI = old })
	msg, err := checkPush("", flowOf(), "production", "3c1b7fcbe17a")
	if err != nil || !strings.Contains(msg, "the merge of #7") {
		t.Errorf("the merged pull request is found by its merge commit, got %q %v", msg, err)
	}
	if len(asked) != 2 || !strings.Contains(asked[1], "pulls?state=closed") {
		t.Errorf("the association first, then the closed pull requests: %v", asked)
	}

	asked = nil
	pullsAPI = func(_, path string) ([]byte, error) {
		asked = append(asked, path)
		return json.Marshal([]pullRequest{merged})
	}
	if _, err := checkPush("", flowOf(), "production", "3c1b7fcbe17a"); err != nil || len(asked) != 1 {
		t.Errorf("an associated merge asks nothing more: %v %v", asked, err)
	}
}
