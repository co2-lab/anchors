// @anchors
//   code: PCTPA
//   ref: PCSPA

package ops

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// pipelineProject is a project with the branch flow and the stages of a review-then-production
// release, with hotfixes going straight to production.
func pipelineProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	yaml := `version: 8
branch_flow:
  merges:
    review: [develop, "hotfix/*"]
    production: [review, "hotfix/*"]
  protected: [review, production]
pipeline:
  stages:
    review-release:
      branches: [review]
      requires: [flow, { gates: ci }]
      verify: verify-review
    production-release:
      branches: [production]
      requires:
        - flow
        - any:
            - passed: review-release
            - from: "hotfix/*"
      verify: verify-prod
    lint:
      run: run-lint
      verify: verify-lint
`
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// forge stands in for GitHub and the commands: the pull requests of a commit, the statuses
// recorded, the gates' outcome and the commands run.
type forge struct {
	statuses map[string]string // sha|context → state
	gatesErr error
	ran      []string
	failing  map[string]bool
}

func fakeForge(t *testing.T, prs ...pullRequest) *forge {
	t.Helper()
	f := &forge{statuses: map[string]string{}, failing: map[string]bool{}}
	pulls(t, prs...)
	oldS, oldSet, oldG, oldR := statusOf, setStatus, runGates, runCommand
	statusOf = func(_, sha, ctx string) (string, error) { return f.statuses[sha+"|"+ctx], nil }
	setStatus = func(_, sha, ctx, state, _ string) error { f.statuses[sha+"|"+ctx] = state; return nil }
	runGates = func(string, string) error { return f.gatesErr }
	runCommand = func(_, script string) error {
		f.ran = append(f.ran, script)
		if f.failing[script] {
			return errors.New("exit status 1")
		}
		return nil
	}
	t.Cleanup(func() { statusOf, setStatus, runGates, runCommand = oldS, oldSet, oldG, oldR })
	return f
}

func runPipeline(t *testing.T, root string, args ...string) (string, error) {
	t.Helper()
	cmd := newPipelineCmd()
	cmd.SetArgs(append(args, "--root", root))
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	var err error
	out := captureStdout(t, func() { err = cmd.Execute() })
	return out, err
}

func withHead(p pullRequest, sha string) pullRequest {
	p.Head.SHA = sha
	return p
}

func TestPipeline_aStageRunsOnItsBranchWhenItsRequirementsHold(t *testing.T) {
	t.Run("PCSPA-B01: can-run passes a stage on its branch when every requirement holds, naming each, and refuses it otherwise", func(t *testing.T) {})
	root := pipelineProject(t)
	f := fakeForge(t, pr(4, "develop", "review", "r1"))
	out, err := runPipeline(t, root, "can-run", "review-release", "--branch", "review", "--sha", "r1")
	if err != nil || !strings.Contains(out, "✓ flow: the merge of #4 from develop") || !strings.Contains(out, "✓ gates: ci: the ci gates pass") {
		t.Errorf("review-release may run on review: %v\n%s", err, out)
	}
	if _, err := runPipeline(t, root, "can-run", "review-release", "--branch", "production", "--sha", "r1"); err == nil || !strings.Contains(err.Error(), "runs on review") {
		t.Errorf("a stage does not run on another branch: %v", err)
	}
	f.gatesErr = errors.New("exit 1")
	if out, err := runPipeline(t, root, "can-run", "review-release", "--branch", "review", "--sha", "r1"); err == nil || !strings.Contains(out, "✗ gates: ci") {
		t.Errorf("failing gates refuse the stage: %v\n%s", err, out)
	}
	pulls(t, pr(5, "feat/x", "review", "r2"))
	f.gatesErr = nil
	if out, err := runPipeline(t, root, "can-run", "review-release", "--branch", "review", "--sha", "r2"); err == nil || !strings.Contains(out, "✗ flow") {
		t.Errorf("a sha that did not come by the flow refuses the stage: %v\n%s", err, out)
	}
}

func TestPipeline_productionTakesWhatPassedReview(t *testing.T) {
	t.Run("PCSPA-B02: passed: asks whether the head the sha was merged from passed that stage; from: whether it came from that head", func(t *testing.T) {})
	root := pipelineProject(t)
	f := fakeForge(t, withHead(pr(9, "review", "production", "p1"), "r1"))
	if out, err := runPipeline(t, root, "can-run", "production-release", "--branch", "production", "--sha", "p1"); err == nil || !strings.Contains(out, "neither this commit nor the head of #9") {
		t.Errorf("a review head that did not pass review refuses production: %v\n%s", err, out)
	}
	f.statuses["r1|anchors/pipeline/review-release"] = "success"
	if out, err := runPipeline(t, root, "can-run", "production-release", "--branch", "production", "--sha", "p1"); err != nil || !strings.Contains(out, "passed review-release") {
		t.Errorf("the head that passed review lets production run: %v\n%s", err, out)
	}
	pulls(t, withHead(pr(10, "hotfix/login", "production", "p2"), "h1"))
	if out, err := runPipeline(t, root, "can-run", "production-release", "--branch", "production", "--sha", "p2"); err != nil || !strings.Contains(out, "comes from hotfix/login") {
		t.Errorf("a hotfix goes to production without review: %v\n%s", err, out)
	}
}

func TestPipeline_anUndeclaredStageDoesNotRun(t *testing.T) {
	t.Run("PCSPA-E01: A stage the configuration does not declare does not run", func(t *testing.T) {})
	t.Run("PCSPA-E02: A stage with no branch or sha, by flag or environment, does not run", func(t *testing.T) {})
	root := pipelineProject(t)
	fakeForge(t)
	if _, err := runPipeline(t, root, "can-run", "deploy", "--branch", "review", "--sha", "x"); err == nil || !strings.Contains(err.Error(), "no stage of the pipeline") {
		t.Errorf("an undeclared stage does not run: %v", err)
	}
	t.Setenv("GITHUB_REF_NAME", "")
	t.Setenv("GITHUB_SHA", "")
	if _, err := runPipeline(t, root, "can-run", "lint"); err == nil || !strings.Contains(err.Error(), "--branch and --sha") {
		t.Errorf("a stage with no commit named does not run: %v", err)
	}
}

func TestPipeline_doneRecordsThePass(t *testing.T) {
	t.Run("PCSPA-B03: record runs verify: and records the stage's outcome on the sha as a commit status", func(t *testing.T) {})
	root := pipelineProject(t)
	f := fakeForge(t)
	t.Setenv("GITHUB_REF_NAME", "review")
	t.Setenv("GITHUB_SHA", "r1")
	if out, err := runPipeline(t, root, "record", "review-release"); err != nil || f.statuses["r1|anchors/pipeline/review-release"] != "success" || !strings.Contains(out, "recorded on r1") {
		t.Errorf("a verify that passes records success: %v %v\n%s", err, f.statuses, out)
	}
	f.failing["verify-review"] = true
	if _, err := runPipeline(t, root, "done", "review-release"); err == nil || f.statuses["r1|anchors/pipeline/review-release"] != "failure" {
		t.Errorf("a verify that fails records failure and fails: %v %v", err, f.statuses)
	}
}

func TestPipeline_runRunsTheStage(t *testing.T) {
	t.Run("PCSPA-B04: run confronts the stage, runs its run: and its verify:, and records the outcome", func(t *testing.T) {})
	root := pipelineProject(t)
	f := fakeForge(t)
	if _, err := runPipeline(t, root, "run", "lint", "--branch", "feat/x", "--sha", "l1"); err != nil || strings.Join(f.ran, ",") != "run-lint,verify-lint" || f.statuses["l1|anchors/pipeline/lint"] != "success" {
		t.Errorf("run runs run: then verify:, and records success: %v %v %v", err, f.ran, f.statuses)
	}
	f.ran, f.failing["run-lint"] = nil, true
	if _, err := runPipeline(t, root, "run", "lint", "--branch", "feat/x", "--sha", "l2"); err == nil || strings.Join(f.ran, ",") != "run-lint" || f.statuses["l2|anchors/pipeline/lint"] != "failure" {
		t.Errorf("a run: that fails stops before verify: and records failure: %v %v %v", err, f.ran, f.statuses)
	}
	if _, err := runPipeline(t, root, "run", "review-release", "--branch", "production", "--sha", "x"); err == nil {
		t.Error("run refuses what can-run refuses")
	}
}

func TestPipeline_aStageItCannotConfrontDoesNotRun(t *testing.T) {
	t.Run("PCSPA-E03: A requirement GitHub cannot answer refuses the stage; a stage with no verify: is not recorded", func(t *testing.T) {})
	root := pipelineProject(t)
	fakeForge(t)
	old := pullsOf
	pullsOf = func(string, string) ([]pullRequest, error) { return nil, errors.New("no gh") }
	defer func() { pullsOf = old }()
	if out, err := runPipeline(t, root, "can-run", "production-release", "--branch", "production", "--sha", "p"); err == nil || !strings.Contains(out, "no gh") {
		t.Errorf("a requirement GitHub cannot answer refuses the stage: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte("version: 8\npipeline:\n  stages:\n    only-run:\n      run: x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runPipeline(t, root, "record", "only-run", "--branch", "b", "--sha", "s"); err == nil || !strings.Contains(err.Error(), "anchors pipeline run only-run") {
		t.Errorf("a stage with no verify: is not done: %v", err)
	}
}

func TestPipeline_aLaterStageOfTheSameShaReadsItsRecord(t *testing.T) {
	t.Run("PCSPA-B05: passed: holds when the sha itself carries the stage's success, before asking the head it was merged from", func(t *testing.T) {})
	root := t.TempDir()
	yaml := "version: 8\npipeline:\n  stages:\n    production-release:\n      branches: [production]\n      verify: v\n" +
		"    publish-app-config:\n      branches: [production]\n      requires: [{ passed: production-release }]\n      verify: v\n"
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	f := fakeForge(t, withHead(pr(9, "review", "production", "p1"), "r1"))
	if _, err := runPipeline(t, root, "can-run", "publish-app-config", "--branch", "production", "--sha", "p1"); err == nil {
		t.Error("the pointer does not flip before production passed on the sha")
	}
	f.statuses["p1|anchors/pipeline/production-release"] = "success"
	if out, err := runPipeline(t, root, "can-run", "publish-app-config", "--branch", "production", "--sha", "p1"); err != nil || !strings.Contains(out, "this commit passed production-release") {
		t.Errorf("the production release recorded on the same sha lets the pointer flip: %v\n%s", err, out)
	}
}
