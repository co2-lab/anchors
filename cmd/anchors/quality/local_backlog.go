package quality

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/queue"
)

// localBacklog is what is still open in the LOCAL mode after a full check: the issues in
// `issues/todo/` and `issues/doing/`, and the tasks of `.anchors/tasks/` not done.
//
// `check --all` said only what THIS run opened or resolved ("N new issue(s), M
// resolved"). What was already open from earlier runs stayed silent — visible only to
// whoever opened the folder — and a clean-looking check sat on
// top of a backlog nobody was working. The full check is the sweep, so it is where the
// backlog is said; the pre-commit (`--changed`) stays quiet, as a line on every commit
// would be noise.
type localBacklog struct {
	Todo, Doing, ForUser  int
	Pending, Claimed, Old int
}

func (b localBacklog) empty() bool {
	return b.Todo == 0 && b.Doing == 0 && b.Pending == 0 && b.Claimed == 0
}

func readLocalBacklog(root string) localBacklog {
	var b localBacklog
	if l, err := issue.List(root, issue.Todo); err == nil {
		b.Todo = len(l)
	}
	if l, err := issue.List(root, issue.Doing); err == nil {
		b.Doing = len(l)
	}
	if l, err := issue.ListByOwner(root, issue.Todo, issue.DonoUsuário); err == nil {
		b.ForUser = len(l)
	}
	if tasks, err := queue.List(root); err == nil {
		for _, t := range tasks {
			switch t.State {
			case queue.Pending:
				b.Pending++
			case queue.Claimed:
				b.Claimed++
				if queue.ClaimIsOld(t) {
					b.Old++
				}
			}
		}
	}
	return b
}

// printLocalBacklog prints the backlog, or nothing when there is none.
func printLocalBacklog(b localBacklog) {
	if b.empty() {
		return
	}
	var linhas []string
	if b.Todo > 0 || b.Doing > 0 {
		linhas = append(linhas, i18n.T("check.backlog_issues", b.Todo, issue.Dir, b.ForUser, b.Doing, issue.Dir))
	}
	if b.Pending > 0 || b.Claimed > 0 {
		linhas = append(linhas, i18n.T("check.backlog_tasks", b.Pending, b.Claimed, b.Old))
	}
	fmt.Println()
	fmt.Println(i18n.T("check.backlog_head"))
	fmt.Println(strings.Join(linhas, "\n"))
}
