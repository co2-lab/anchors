// @anchors
//   code: PPCMP
//   ref: PCSPA

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
	"github.com/co2-lab/anchors/internal/shell"
	"github.com/spf13/cobra"
)

// A stage of the pipeline asks the Anchors before it runs, and tells it when it passed. The
// rules are the project's (`pipeline.stages`): the branch, the branch flow the sha came by, the
// gates of a phase, a stage the sha passed before. What passed is recorded where the next stage
// can read it, on the sha itself, as a commit status `anchors/pipeline/<stage>`. A stage the
// configuration does not declare does not run (reported from MIF: the production release ran
// on commits no review had seen, and nothing it consulted said it should not).

// statusContext is the commit status a stage's pass is recorded under.
func statusContext(stage string) string { return "anchors/pipeline/" + stage }

// statusOf is the state of a commit status on sha: success, failure, pending, or empty.
var statusOf = func(root, sha, context string) (string, error) {
	out, err := ghAPI(root, "repos/{owner}/{repo}/commits/"+sha+"/status")
	if err != nil {
		return "", err
	}
	var combined struct {
		Statuses []struct {
			Context string `json:"context"`
			State   string `json:"state"`
		} `json:"statuses"`
	}
	if err := json.Unmarshal(out, &combined); err != nil {
		return "", fmt.Errorf("read the statuses of %s: %w", sha, err)
	}
	for _, s := range combined.Statuses {
		if s.Context == context {
			return s.State, nil
		}
	}
	return "", nil
}

// setStatus records a commit status on sha.
var setStatus = func(root, sha, context, state, description string) error {
	_, err := ghAPI(root, "-X", "POST", "repos/{owner}/{repo}/statuses/"+sha,
		"-f", "state="+state, "-f", "context="+context, "-f", "description="+description)
	return err
}

// runGates runs the gates of a phase over the whole project, as `anchors verify` does.
var runGates = func(root, phase string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	c := exec.Command(exe, "verify", "--phase", phase, "--all", "--no-record", "--root", root)
	c.Stdout, c.Stderr = os.Stderr, os.Stderr
	return c.Run()
}

// runCommand runs a stage's command in the project's root, its output the CI's.
var runCommand = func(root, script string) error {
	c, err := shell.Command(script)
	if err != nil {
		return err
	}
	c.Dir, c.Stdout, c.Stderr = root, os.Stdout, os.Stderr
	return c.Run()
}

func ghAPI(root string, args ...string) ([]byte, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return nil, errors.New("the gh CLI is not in the PATH: a stage's record lives in GitHub's commit statuses")
	}
	c := exec.Command("gh", append([]string{"api"}, args...)...)
	c.Dir = root
	out, err := c.Output()
	if err != nil {
		return nil, fmt.Errorf("gh api %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

func newPipelineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Ask whether a stage of the pipeline may run, and record that it passed (`pipeline:`)",
		Long: `The stages anchors.yaml declares, asked by the pipeline:

  pipeline:
    stages:
      review-release:
        branches: [review]
        requires: [flow, { gates: ci }]
        verify: node scripts/smoke-review.mjs
      production-release:
        branches: [production]
        requires:
          - flow
          - any: [{ passed: review-release }, { from: "hotfix/*" }]
        verify: node scripts/smoke-prod.mjs

  anchors pipeline can-run <stage>   the first job: refuses a stage that may not run
  anchors pipeline record <stage>    the last job: runs verify:, records the pass on the sha
                                     (alias done; record reads clean to shellcheck, done does not)
  anchors pipeline run <stage>       can-run, run:, verify: and the record, in one step

A stage the configuration does not declare does not run. The branch and the sha come from
--branch and --sha, or $GITHUB_REF_NAME and $GITHUB_SHA. The pass is a commit status
anchors/pipeline/<stage>, written and read with the gh CLI (GH_TOKEN, statuses: write).`,
	}
	cmd.AddCommand(newPipelineCanRunCmd(), newPipelineRecordCmd(), newPipelineRunCmd())
	return cmd
}

// stageTarget is what a pipeline command acts on.
type stageTarget struct {
	root, name, branch, sha string
	cfg                     *config.Config
	stage                   config.Stage
}

func stageCommand(use, short string, act func(t stageTarget) error) *cobra.Command {
	var root, branch, sha string
	cmd := &cobra.Command{
		Use:   use + " <stage>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			t := stageTarget{root: absRoot, name: args[0], cfg: cfg,
				branch: or(branch, os.Getenv("GITHUB_REF_NAME")), sha: or(sha, os.Getenv("GITHUB_SHA"))}
			var s config.Stage
			ok := false
			if cfg.Pipeline != nil {
				s, ok = cfg.Pipeline.Stages[t.name]
			}
			if !ok {
				return fmt.Errorf("✗ %s is no stage of the pipeline (pipeline.stages: %s) — a stage the configuration does not declare does not run",
					t.name, strings.Join(cfg.Pipeline.StageNames(), ", "))
			}
			if t.branch == "" || t.sha == "" {
				return errors.New("name the stage's commit: --branch and --sha, or $GITHUB_REF_NAME and $GITHUB_SHA")
			}
			t.stage = s
			return act(t)
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&branch, "branch", "", "the branch the stage runs on (default $GITHUB_REF_NAME)")
	cmd.Flags().StringVar(&sha, "sha", "", "the commit the stage runs on (default $GITHUB_SHA)")
	return cmd
}

func newPipelineCanRunCmd() *cobra.Command {
	return stageCommand("can-run", "Refuse a stage that may not run on this branch and sha", func(t stageTarget) error {
		lines, err := canRun(t)
		fmt.Print(lines)
		return err
	})
}

// newPipelineRecordCmd is the stage's last job. It was named `done`, and shellcheck — through
// actionlint, on every workflow that calls it — read the word as the end of a loop (SC1010,
// reported from MIF); `done` stays as an alias, so a workflow written with it still runs.
func newPipelineRecordCmd() *cobra.Command {
	cmd := stageCommand("record", "Run the stage's verify: and record on the sha that the stage passed", func(t stageTarget) error {
		if strings.TrimSpace(t.stage.Verify) == "" {
			return fmt.Errorf("✗ %s declares no verify: — a stage the Anchors runs is run with `anchors pipeline run %s`", t.name, t.name)
		}
		return finish(t, runCommand(t.root, t.stage.Verify), "verify")
	})
	cmd.Aliases = []string{"done"}
	return cmd
}

func newPipelineRunCmd() *cobra.Command {
	return stageCommand("run", "Run a stage: can-run, its run:, its verify:, and the record", func(t stageTarget) error {
		lines, err := canRun(t)
		fmt.Print(lines)
		if err != nil {
			return err
		}
		if strings.TrimSpace(t.stage.Run) != "" {
			if err := runCommand(t.root, t.stage.Run); err != nil {
				return finish(t, err, "run")
			}
		}
		var verr error
		if strings.TrimSpace(t.stage.Verify) != "" {
			verr = runCommand(t.root, t.stage.Verify)
		}
		return finish(t, verr, "verify")
	})
}

// finish records the stage's outcome on its sha: success when step passed, failure otherwise.
func finish(t stageTarget, stepErr error, step string) error {
	state, desc := "success", "passed"
	if stepErr != nil {
		state, desc = "failure", step+" failed"
	}
	if err := setStatus(t.root, t.sha, statusContext(t.name), state, desc); err != nil {
		return fmt.Errorf("record %s on %s: %w", t.name, short(t.sha), err)
	}
	if stepErr != nil {
		return fmt.Errorf("✗ %s: %s failed (%v), recorded as failure on %s", t.name, step, stepErr, short(t.sha))
	}
	fmt.Printf("✓ %s passed, recorded on %s (%s)\n", t.name, short(t.sha), statusContext(t.name))
	return nil
}

// canRun confronts the stage with its branches and requirements, and says each one.
func canRun(t stageTarget) (string, error) {
	var b strings.Builder
	if len(t.stage.Branches) > 0 {
		on := false
		for _, p := range t.stage.Branches {
			on = on || config.BranchMatches(p, t.branch)
		}
		if !on {
			return "", fmt.Errorf("✗ %s does not run on %s — it runs on %s", t.name, t.branch, strings.Join(t.stage.Branches, ", "))
		}
	}
	e := &stageEval{t: t}
	var failed []string
	for _, r := range t.stage.Requires {
		ok, why := e.holds(r)
		mark := "✓"
		if !ok {
			mark = "✗"
			failed = append(failed, r.String())
		}
		fmt.Fprintf(&b, "%s %s: %s\n", mark, r, why)
	}
	if len(failed) > 0 {
		return b.String(), fmt.Errorf("✗ %s may not run on %s (%s): %s", t.name, t.branch, short(t.sha), strings.Join(failed, "; "))
	}
	fmt.Fprintf(&b, "✓ %s may run on %s (%s)\n", t.name, t.branch, short(t.sha))
	return b.String(), nil
}

// stageEval answers a stage's requirements, asking GitHub for the sha's pull request once.
type stageEval struct {
	t       stageTarget
	asked   bool
	pr      *pullRequest
	prError error
}

// merged is the merged pull request into the stage's branch whose merge commit is its sha.
func (e *stageEval) merged() (*pullRequest, error) {
	if !e.asked {
		e.asked = true
		prs, err := pullsOf(e.t.root, e.t.sha)
		e.prError = err
		for i, pr := range prs {
			if pr.MergedAt != "" && pr.Base.Ref == e.t.branch && pr.MergeCommitSHA == e.t.sha {
				e.pr = &prs[i]
				break
			}
		}
	}
	return e.pr, e.prError
}

func (e *stageEval) holds(r config.Requirement) (bool, string) {
	t := e.t
	switch {
	case r.Flow:
		if !guards(t.cfg.BranchFlow, t.branch) {
			return true, "the branch flow does not guard " + t.branch
		}
		pr, refused, err := mergingPR(t.root, t.cfg.BranchFlow, t.branch, t.sha)
		if err != nil {
			return false, err.Error()
		}
		if pr == nil {
			if len(refused) > 0 {
				return false, "the merge of " + strings.Join(refused, ", ") + ", which the branch flow refuses"
			}
			return false, "the merge of no pull request"
		}
		return true, fmt.Sprintf("the merge of #%d from %s", pr.Number, pr.Head.Ref)
	case r.Gates != "":
		if err := runGates(t.root, r.Gates); err != nil {
			return false, "the " + r.Gates + " gates fail"
		}
		return true, "the " + r.Gates + " gates pass"
	case r.Passed != "":
		// The sha itself first — a later stage of the same branch reads the stage before it
		// there (reported from MIF: the version pointer flipped only after the production
		// release passed on the same sha, which the head of the review pull request never
		// carries) —, then the head the sha was merged from, the commit a stage of the branch
		// before it passed.
		state, err := statusOf(t.root, t.sha, statusContext(r.Passed))
		if err != nil {
			return false, err.Error()
		}
		if state == "success" {
			return true, fmt.Sprintf("this commit passed %s", r.Passed)
		}
		pr, perr := e.merged()
		if perr != nil || pr == nil || pr.Head.SHA == "" {
			return false, fmt.Sprintf("this commit has not passed %s (%s)", r.Passed, or(state, "no record"))
		}
		of := fmt.Sprintf("the head of #%d (%s)", pr.Number, short(pr.Head.SHA))
		if state, err = statusOf(t.root, pr.Head.SHA, statusContext(r.Passed)); err != nil {
			return false, err.Error()
		}
		if state != "success" {
			return false, fmt.Sprintf("neither this commit nor %s passed %s (%s)", of, r.Passed, or(state, "no record"))
		}
		return true, fmt.Sprintf("%s passed %s", of, r.Passed)
	case r.From != "":
		pr, err := e.merged()
		if err != nil {
			return false, err.Error()
		}
		if pr == nil {
			return false, "the merge of no pull request"
		}
		if !config.BranchMatches(r.From, pr.Head.Ref) {
			return false, fmt.Sprintf("#%d comes from %s", pr.Number, pr.Head.Ref)
		}
		return true, fmt.Sprintf("#%d comes from %s", pr.Number, pr.Head.Ref)
	}
	var whys []string
	for _, a := range r.Any {
		ok, why := e.holds(a)
		if ok {
			return true, a.String() + ": " + why
		}
		whys = append(whys, a.String()+": "+why)
	}
	return false, "none holds — " + strings.Join(whys, "; ")
}
