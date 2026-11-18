// @anchors
//   code: RBCRP
//   ref: RPBUG

package flow

import (
	"fmt"
	"io"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/co2-lab/anchors/cmd/anchors/common"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/spf13/cobra"
)

// upstreamRepo is where Anchors' own bugs are reported.
const upstreamRepo = "co2-lab/anchors"

// anchorsBug is a bug in Anchors itself, in the sections of the repository's bug form.
type anchorsBug struct {
	What, Expected, Repro string
}

func (b anchorsBug) title() string { return "[bug] " + firstLineOfReason(b.What) }

// body is the issue's text: the same sections as the repository's bug form, so a report
// from an agent reads like one from a person.
func (b anchorsBug) body() string {
	var s strings.Builder
	s.WriteString("### What happened\n\n" + strings.TrimSpace(b.What) + "\n\n")
	if e := strings.TrimSpace(b.Expected); e != "" {
		s.WriteString("### What should happen\n\n" + e + "\n\n")
	}
	if r := strings.TrimSpace(b.Repro); r != "" {
		s.WriteString("### Minimal case\n\n```shell\n" + r + "\n```\n\n")
	}
	s.WriteString("### Version\n\n" + common.Version + "\n\n### Platform\n\n" + runtime.GOOS + "/" + runtime.GOARCH + "\n\n")
	s.WriteString("---\nReported by `anchors report-bug`.\n")
	return s.String()
}

// projectMarks are the strings that would identify the project in a public report: its
// absolute path, the user's home, and its repository on the platform.
func projectMarks(root string, cfg *config.Config) []string {
	var marks []string
	if abs, err := config.AbsRoot(root); err == nil && abs != "" && abs != string(filepath.Separator) {
		marks = append(marks, abs)
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" && home != string(filepath.Separator) {
		marks = append(marks, home)
	}
	if cfg != nil && cfg.Workflow != nil && strings.TrimSpace(cfg.Workflow.Repo) != "" {
		marks = append(marks, strings.TrimSpace(cfg.Workflow.Repo))
	}
	if out, err := exec.Command("git", "-C", root, "remote", "get-url", "origin").Output(); err == nil {
		if r := repoOfRemote(strings.TrimSpace(string(out))); r != "" && !strings.EqualFold(r, upstreamRepo) {
			marks = append(marks, r)
		}
	}
	return marks
}

// repoOfRemote is the `owner/name` of a git remote URL, in its https or ssh form.
func repoOfRemote(u string) string {
	u = strings.TrimSuffix(u, ".git")
	if i := strings.Index(u, "://"); i >= 0 {
		u = u[i+3:]
		if j := strings.Index(u, "/"); j >= 0 {
			u = u[j+1:]
		}
	} else if i := strings.Index(u, ":"); i >= 0 {
		u = u[i+1:]
	}
	parts := strings.Split(strings.Trim(u, "/"), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2] + "/" + parts[len(parts)-1]
}

// leaks are the project marks a text carries, compared without case.
func leaks(text string, marks []string) []string {
	low := strings.ToLower(text)
	var out []string
	for _, m := range marks {
		if m != "" && strings.Contains(low, strings.ToLower(m)) {
			out = append(out, m)
		}
	}
	return out
}

// sendAnchorsBug reports the bug at Anchors' repository and returns the issue's address.
//
// An OPEN issue with the same title gets a "seen again" comment with the release and the
// platform, instead of a second issue: the same defect met by several projects is one
// issue with one comment per sighting. A refusal prints the prefilled new-issue link a
// person can open, and is returned as the error.
func sendAnchorsBug(b anchorsBug) (string, error) {
	titulo := b.title()
	if out, err := exec.Command("gh", "issue", "list", "--repo", upstreamRepo, "--state", "open",
		"--search", firstLineOfReason(b.What)+" in:title", "--json", "url,title",
		"--jq", `.[] | .title + "\t" + .url`).Output(); err == nil {
		for _, l := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			t, u, ok := strings.Cut(l, "\t")
			if ok && strings.EqualFold(strings.TrimSpace(t), titulo) {
				_ = exec.Command("gh", "issue", "comment", u, "--repo", upstreamRepo,
					"--body", fmt.Sprintf("Seen again — anchors %s · %s/%s.", common.Version, runtime.GOOS, runtime.GOARCH)).Run()
				fmt.Printf("already reported to Anchors: %s (a comment says it was seen again)\n", u)
				return u, nil
			}
		}
	}
	corpo := b.body()
	out, err := exec.Command("gh", "issue", "create", "--repo", upstreamRepo,
		"--title", titulo, "--body", corpo, "--label", "bug").CombinedOutput()
	if err != nil {
		q := neturl.Values{"title": {titulo}, "body": {corpo}}
		fmt.Printf("· warning: could not report to Anchors (%s) — open it by hand:\n  https://github.com/%s/issues/new?%s\n",
			strings.TrimSpace(string(out)), upstreamRepo, q.Encode())
		return "", fmt.Errorf("could not report to Anchors")
	}
	u := strings.TrimSpace(string(out))
	fmt.Printf("reported to Anchors: %s\n", u)
	return u, nil
}

// reportUpstream is `escalate --upstream`: the reason alone, and a warning instead of a
// failure — the finding is already recorded where it was found.
func reportUpstream(motivo string, marks []string) string {
	b := anchorsBug{What: motivo}
	if found := leaks(b.What, marks); len(found) > 0 {
		fmt.Printf("· warning: not reported to Anchors — the reason names this project or this machine (%s). "+
			"Rewrite it in Anchors' terms with `anchors report-bug`.\n", strings.Join(found, ", "))
		return ""
	}
	u, _ := sendAnchorsBug(b)
	return u
}

// newReportBugCmd reports a bug in Anchors itself, from any project and in any mode.
func newReportBugCmd() *cobra.Command {
	var root, expected, repro string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "report-bug <what happened>",
		Short: "Report a bug in Anchors itself at github.com/co2-lab/anchors",
		Long: `When what is wrong is ANCHORS — a command, a gate, the map, a file Anchors seeds —
and not this project's configuration or code, report it where it is fixed for every
project: github.com/co2-lab/anchors. ` + "`anchors guide report-bug`" + ` says how to tell
the two apart and what to do while the fix does not come.

  anchors report-bug "<what happened>" --expected "<what should happen>" [--repro <file>|-]
  anchors report-bug ... --dry-run      print the issue, send nothing

The issue has the sections of the repository's bug form, with the release and the
platform filled in. An open issue with the same title gets a "seen again" comment
instead of a duplicate.

That repository is PUBLIC. Write in Anchors' terms — the command, the gate, what it did,
what it should do — and a made-up minimal case, never this project's code, names or
data. The command refuses a text that carries the project's path, the home folder or the
project's repository.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			what := strings.TrimSpace(strings.Join(args, " "))
			if what == "" {
				return fmt.Errorf("say what happened: the command or gate, and what it did")
			}
			if strings.TrimSpace(expected) == "" {
				return fmt.Errorf("`--expected` is required: what should have happened is half of a bug report")
			}
			caso := ""
			switch repro {
			case "":
			case "-":
				b, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return err
				}
				caso = string(b)
			default:
				b, err := os.ReadFile(repro)
				if err != nil {
					return fmt.Errorf("read the minimal case: %w", err)
				}
				caso = string(b)
			}
			bug := anchorsBug{What: what, Expected: expected, Repro: caso}
			var cfg *config.Config
			if abs, err := config.AbsRoot(root); err == nil {
				if c, err := config.Load(filepath.Join(abs, config.DefaultFile)); err == nil {
					cfg = c
				}
			}
			if found := leaks(bug.title()+"\n"+bug.body(), projectMarks(root, cfg)); len(found) > 0 {
				return fmt.Errorf("the report names this project or this machine (%s), and the repository is public — "+
					"rewrite it in Anchors' terms, with a made-up minimal case", strings.Join(found, ", "))
			}
			if dryRun {
				fmt.Printf("%s\n\n%s", bug.title(), bug.body())
				fmt.Println("\n(dry run — nothing was sent)")
				return nil
			}
			_, err := sendAnchorsBug(bug)
			return err
		},
	}
	cmd.Flags().StringVar(&root, "root", ".", "project root")
	cmd.Flags().StringVar(&expected, "expected", "", "what should have happened (required)")
	cmd.Flags().StringVar(&repro, "repro", "", "a file with the minimal case — made up, not the project's — or - for standard input")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print the issue and send nothing")
	return cmd
}
