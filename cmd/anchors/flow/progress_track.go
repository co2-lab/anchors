// @anchors
//   code: PCCPR
//   ref: LDCML

package flow

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/progress"
	"github.com/spf13/cobra"
)

// `anchors progress` is the only writer of anchors.progress.md (DESIGN-progress-tracking.md):
// each subcommand is one change, made under the file's lock, so agents working in parallel
// never lose each other's, and the file keeps the shape Anchors reads back.

func newProgressTrackCmd() *cobra.Command {
	var root string
	cmd := &cobra.Command{
		Use:   "progress",
		Short: "Track the project's demands, stages and weighted tasks in anchors.progress.md",
		Long: `Keeps anchors.progress.md, the project's progress: its demands (a plan in development, or any
other piece of work), their optional stages, and their tasks, each with a weight — an estimate
of its effort in the unit anchors.yaml names (progress.unit, default points). Progress is the
done weight over the weight not cancelled.

  anchors progress add demand "Checkout" --plan plans/0012-checkout.md
  anchors progress add stage D01 "Payment"
  anchors progress add task D01.S1 "refund rules" --weight 5
  anchors progress start D01.T4         (held by you; --take to take another agent's)
  anchors progress done D01.T4
  anchors progress reopen D01.T4 | cancel D01.T4 --reason "…" | weigh D01.T4 8
  anchors progress rename D01.S2 "Payment"     a new title; the id and the history stay

Work the plan does not list, done for it, is recorded before it is done — where it arises,
marked so it is seen:
  anchors progress add task D01.S2 "fix the date parser" --weight 2 --unplanned
  anchors progress unplanned D01.T31 [--undo]   mark a task already added
  anchors progress move D01.T31 D01.S2         take a task to the stage it arose in
  anchors progress drop D01.S8                 remove a stage left with no task
  anchors progress summary              the demands and their stages, with their progress
  anchors progress show [D01] [--json]  every task

Only these commands write the file: it is never edited by hand, and a hand edit is refused.
Each change takes the file's lock, so agents working at once lose nothing.`,
	}
	cmd.PersistentFlags().StringVar(&root, "root", ".", "project root")
	rootOf := func() (string, string, error) {
		abs, err := config.AbsRoot(root)
		if err != nil {
			return "", "", err
		}
		cfg, err := common.LoadConfig(abs)
		if err != nil {
			return "", "", err
		}
		return abs, cfg.ProgressUnit(), nil
	}
	change := func(c *cobra.Command, fn func(*progress.File) (string, error)) error {
		abs, unit, err := rootOf()
		if err != nil {
			return err
		}
		var said string
		if err := progress.Update(abs, unit, func(f *progress.File) error {
			var e error
			said, e = fn(f)
			return e
		}); err != nil {
			return err
		}
		fmt.Fprintln(c.OutOrStdout(), said)
		return nil
	}

	add := &cobra.Command{Use: "add", Short: "Add a demand, a stage or a task"}
	var plan string
	addDemand := &cobra.Command{
		Use: "demand [title] [--plan <plan>]", Short: "Add a demand — a plan in development, or any piece of work", Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			title := ""
			if len(args) == 1 {
				title = args[0]
			}
			if plan == "" {
				if title == "" {
					return fmt.Errorf("a demand needs a title, or --plan for the plan it develops")
				}
				return change(c, func(f *progress.File) (string, error) {
					id, err := f.AddDemand(title, "")
					return "✓ demand " + id + " — " + title, err
				})
			}
			// A plan becomes a demand titled as the plan, one stage per phase it declares.
			abs, _, err := rootOf()
			if err != nil {
				return err
			}
			rel := common.RelTo(abs, plan)
			content, err := os.ReadFile(filepath.Join(abs, filepath.FromSlash(rel)))
			if err != nil {
				return fmt.Errorf("read the plan: %w", err)
			}
			return change(c, func(f *progress.File) (string, error) {
				id, err := f.PlanDemand(rel, string(content), title)
				return "✓ demand " + id + " — the plan " + rel, err
			})
		},
	}
	addDemand.Flags().StringVar(&plan, "plan", "", "the plan the demand develops: its phases become stages")
	addStage := &cobra.Command{
		Use: "stage <demand> <title>", Short: "Add a stage to a demand", Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				id, err := f.AddStage(args[0], args[1])
				return "✓ stage " + id + " — " + args[1], err
			})
		},
	}
	var weight float64
	var estimate, unplanned bool
	addTask := &cobra.Command{
		Use: "task <demand|stage> <title> --weight <n>", Short: "Add a task, with its weight", Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				id, err := f.AddTask(args[0], args[1], weight, estimate)
				if err == nil && unplanned {
					err = f.MarkUnplanned(id, true)
				}
				return "✓ task " + id + " — " + args[1], err
			})
		},
	}
	addTask.Flags().Float64Var(&weight, "weight", 0, "the task's weight — an estimate of its effort (required)")
	addTask.Flags().BoolVar(&estimate, "estimate", false, "the weight is a placeholder still to be estimated")
	addTask.Flags().BoolVar(&unplanned, "unplanned", false, "work the plan did not list, found on the way: kept where it arose, marked so it is seen")
	add.AddCommand(addDemand, addStage, addTask)

	var by string
	var take bool
	who := func() string {
		if by != "" {
			return by
		}
		return agentID()
	}
	start := &cobra.Command{
		Use: "start <task>", Short: "Mark a task in progress, held by you", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				return "~ " + args[0] + " in progress, held by " + who(), f.Start(args[0], who(), time.Now().Format(progress.StartedAt), take)
			})
		},
	}
	start.Flags().BoolVar(&take, "take", false, "take over a task another agent holds")
	done := &cobra.Command{
		Use: "done <task>", Short: "Mark a task done", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				return "✓ " + args[0] + " done by " + who(), f.Done(args[0], who(), gitmeta.Today())
			})
		},
	}
	for _, c := range []*cobra.Command{start, done} {
		c.Flags().StringVar(&by, "by", "", "who (default: this agent's identity)")
	}
	reopen := &cobra.Command{
		Use: "reopen <task>", Short: "Put a task back to do", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				return "↺ " + args[0] + " reopened", f.Reopen(args[0])
			})
		},
	}
	var reason string
	cancel := &cobra.Command{
		Use: "cancel <task> --reason <why>", Short: "Take a task out of the count, with the reason", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				return "− " + args[0] + " cancelled", f.Cancel(args[0], reason)
			})
		},
	}
	cancel.Flags().StringVar(&reason, "reason", "", "why the task leaves the count (required)")
	weigh := &cobra.Command{
		Use: "weigh <task> <weight>", Short: "Set a task's weight", Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			w, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return fmt.Errorf("the weight %q is not a number", args[1])
			}
			return change(c, func(f *progress.File) (string, error) {
				return "⚖ " + args[0] + " weighs " + args[1], f.Weigh(args[0], w)
			})
		},
	}

	rename := &cobra.Command{
		Use: "rename <demand|stage|task> <title>", Short: "Give a demand, a stage or a task a new title; its id and history stay", Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				return "✎ " + args[0] + " — " + args[1], f.Rename(args[0], args[1])
			})
		},
	}

	move := &cobra.Command{
		Use: "move <task> <demand|stage>", Short: "Take a task to another stage of its demand; its id and history stay", Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				return "→ " + args[0] + " now in " + args[1], f.Relocate(args[0], args[1])
			})
		},
	}
	drop := &cobra.Command{
		Use: "drop <stage>", Short: "Remove a stage left with no task", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				return "✗ stage " + args[0] + " removed", f.DropStage(args[0])
			})
		},
	}

	var undo bool
	markUnplanned := &cobra.Command{
		Use: "unplanned <task> [--undo]", Short: "Mark a task as work the plan did not list (--undo takes the mark off)", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			return change(c, func(f *progress.File) (string, error) {
				said := "◇ " + args[0] + " unplanned"
				if undo {
					said = "◇ " + args[0] + " planned"
				}
				return said, f.MarkUnplanned(args[0], !undo)
			})
		},
	}
	markUnplanned.Flags().BoolVar(&undo, "undo", false, "take the mark off")

	var asJSON bool
	read := func() (*progress.File, error) {
		abs, _, err := rootOf()
		if err != nil {
			return nil, err
		}
		return progress.Load(abs)
	}
	summary := &cobra.Command{
		Use: "summary", Short: "What is in progress: each demand and its stages, with their progress — never the tasks",
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			f, err := read()
			if err != nil {
				return err
			}
			return printSummary(c.OutOrStdout(), f, asJSON)
		},
	}
	show := &cobra.Command{
		Use: "show [demand]", Short: "Every task of the file, or of one demand", Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			c.SilenceUsage = true
			f, err := read()
			if err != nil {
				return err
			}
			if len(args) == 1 {
				d := f.Demand(args[0])
				if d == nil {
					return fmt.Errorf("no demand %s", args[0])
				}
				f = &progress.File{Unit: f.Unit, Demands: []*progress.Demand{d}}
			}
			if asJSON {
				return writeJSON(c.OutOrStdout(), f)
			}
			fmt.Fprint(c.OutOrStdout(), strings.TrimSuffix(progressBody(f), "\n")+"\n")
			return nil
		},
	}
	for _, c := range []*cobra.Command{summary, show} {
		c.Flags().BoolVar(&asJSON, "json", false, "print JSON")
	}
	cmd.AddCommand(add, start, done, reopen, cancel, weigh, rename, move, drop, markUnplanned, summary, show)
	return cmd
}

// progressBody is the file as Anchors renders it, without its first and last lines — the
// header that says who writes it, and the sum.
func progressBody(f *progress.File) string {
	lines := strings.Split(strings.TrimRight(progress.Render(f), "\n"), "\n")
	return strings.Join(lines[1:len(lines)-1], "\n") + "\n"
}

// printSummary prints one line per demand and per stage, with the progress and what is
// open; the tasks are left out.
func printSummary(w io.Writer, f *progress.File, asJSON bool) error {
	rows := progress.Summary(f)
	if asJSON {
		if rows == nil {
			rows = []progress.Row{}
		}
		return writeJSON(w, rows)
	}
	unit := f.Unit
	done, total := progress.Totals(f.CountedTasks())
	if len(rows) == 0 {
		fmt.Fprintln(w, "nothing in progress — `anchors progress add demand \"<title>\"` starts one")
		return nil
	}
	fmt.Fprintf(w, "Progress — %s (%s/%s %s)\n", progress.Percent(done, total), num(done), num(total), unit)
	for _, r := range rows {
		indent, plan := "  ", ""
		if r.Stage {
			indent = "      "
		}
		if r.Plan != "" {
			plan = " (" + r.Plan + ")"
		}
		extra := fmt.Sprintf("%d open", r.Open)
		if r.Doing > 0 {
			extra += fmt.Sprintf(", %d in progress", r.Doing)
		}
		if r.Estimate > 0 {
			extra += fmt.Sprintf(", %d weight(s) still estimates", r.Estimate)
		}
		if r.Unplanned > 0 {
			extra += fmt.Sprintf(" · %s of %s %s unplanned", num(r.Unplanned), num(r.Total), unit)
		}
		if r.Frozen != "" && !r.Stage {
			extra += " · ❄ frozen (" + r.Frozen + "), out of the total"
		}
		fmt.Fprintf(w, "%s%s %s%s — %s (%s/%s) · %s\n", indent, r.ID, r.Title, plan, r.Percent, num(r.Done), num(r.Total), extra)
	}
	return nil
}

func num(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

func writeJSON(w io.Writer, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(w, string(b))
	return nil
}
