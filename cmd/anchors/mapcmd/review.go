// @anchors
//   code: RVCMA
//   ref: RVCMR

package mapcmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// newReviewCmd records a review of a target for a gate that declares `review:`, or lists
// what is to review.
//
// A review is a second look at what an agent decided, and its product is FINDINGS: each
// one becomes a failing test and then a fix, or is dismissed with a reason. So a review
// records who looked, at which revision, and opens an issue with the findings when there
// are any — it does not stamp a verdict, and it has no `waived`: a review that did not
// happen is a review still due.
func newReviewCmd() *cobra.Command {
	var root, mapPath, gateName, by, findings string
	var pending, recordIssues bool
	cmd := &cobra.Command{
		Use:   "review [target]",
		Short: "Record a review of a target, or list what is to review",
		Long: `A gate that declares ` + "`review:`" + ` marks its targets TO REVIEW, apart from how it
measures: a target is to review until a review is recorded at its current revision, and a
change to it makes the review due again. Reviews inform; they never block.

  anchors review --pending [--gate <g>]                 what is to review, and the question
  anchors review <target> --gate <g> --by <who>         a review that found nothing
  anchors review <target> --gate <g> --by <who> --findings "<report>"

--by says who reviewed, as you name them (human:ana, agent:<vendor>/<model>); the record
keeps it. --findings is the whole report — each finding with what, WHERE (file:line) and
why — and it opens an issue, the way a failing judgment does: each finding is then a
failing test and a fix, or dismissed with its reason in the issue.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := config.AbsRoot(root)
			if err != nil {
				return err
			}
			if mapPath == "" {
				mapPath = filepath.Join(absRoot, mapx.DefaultPath)
			}
			cfg, err := config.Load(filepath.Join(absRoot, config.DefaultFile))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			g, err := mapx.Load(mapPath)
			if err != nil {
				return fmt.Errorf("load map: %w (run `anchors map build`)", err)
			}
			if pending {
				return listDueReviews(cfg, g, absRoot, gateName)
			}
			if len(args) != 1 {
				return fmt.Errorf("%s", i18n.T("review.need_target"))
			}
			gc, ok := reviewGate(cfg, gateName)
			if !ok {
				return fmt.Errorf("%s", i18n.T("review.unknown_gate", gateName))
			}
			if strings.TrimSpace(by) == "" {
				return fmt.Errorf("%s", i18n.T("review.need_by"))
			}
			target := relTo(absRoot, args[0])
			if !nodeExists(g, target) {
				return fmt.Errorf("target %q is not in the map", target)
			}
			if cfg.GitHubMode() && len(cfg.Workflow.Labels) > 0 {
				issue.UseGitHub(cfg.Workflow.Repo, cfg.Workflow.Labels[0])
			}

			now := time.Now()
			found := strings.TrimSpace(findings) != ""
			rec := mapx.Review{Gate: gc.Name, By: strings.TrimSpace(by), At: now.Format(time.DateOnly), Findings: found}
			if err := mapx.Update(mapPath, func(g *mapx.Graph) error {
				if !g.RecordReview(target, rec) {
					return fmt.Errorf("target %q is not in the map", target)
				}
				return nil
			}); err != nil {
				return fmt.Errorf("save review: %w", err)
			}

			if !found {
				fmt.Println(i18n.T("review.recorded_clean", target, gc.Name, rec.By))
				return nil
			}
			if !judgeWritesIssue(cfg, recordIssues) {
				fmt.Println(i18n.T("review.recorded_findings_no_issue", target, gc.Name, rec.By))
				fmt.Printf("\n%s\n\n", strings.TrimSpace(findings))
				return nil
			}
			iss := issue.Issue{Kind: issue.Violation, Target: target, Gate: gc.Name, Detail: findings, Date: now.Format("2006-01-02")}
			created, at, err := issue.Open(absRoot, iss)
			if err != nil {
				return err
			}
			if !created {
				if _, err := issue.Reopen(absRoot, iss); err != nil {
					return err
				}
			}
			fmt.Println(i18n.T("review.recorded_findings", target, gc.Name, rec.By, issue.Dir+"/"+at+"/"))
			return nil
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&mapPath, "map", "", "path to the map")
	cmd.Flags().StringVar(&gateName, "gate", "", "the gate whose review this is (a gate that declares `review:`)")
	cmd.Flags().StringVar(&by, "by", "", "who reviewed, as you name them (human:ana, agent:<vendor>/<model>)")
	cmd.Flags().StringVar(&findings, "findings", "", "the findings, as the whole report; each one opens the issue's work")
	cmd.Flags().BoolVar(&pending, "pending", false, "list the targets to review")
	cmd.Flags().BoolVar(&recordIssues, "record-issues", false, "in manual mode, also write the issue of the findings")
	return cmd
}

// reviewGate is the project's gate of that name, when it declares `review:`.
func reviewGate(cfg *config.Config, name string) (config.Gate, bool) {
	for _, g := range cfg.Gates {
		if g.Name == name && g.Review != nil {
			return g, true
		}
	}
	return config.Gate{}, false
}

// listDueReviews prints what is to review, grouped by gate with its question.
func listDueReviews(cfg *config.Config, g *mapx.Graph, root, only string) error {
	due := gate.ReviewsDue(cfg, g, root, only)
	if len(due) == 0 {
		fmt.Println(i18n.T("review.none_due"))
		return nil
	}
	last := ""
	for _, d := range due {
		if d.Gate != last {
			fmt.Printf("\n%s — %s\n", d.Gate, d.Ask)
			last = d.Gate
		}
		fmt.Printf("  ○ %s\n", d.Target)
	}
	fmt.Println("\n" + i18n.T("review.due_footer", len(due)))
	return nil
}
