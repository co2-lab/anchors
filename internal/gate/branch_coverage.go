package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// branch-coverage: the branches no test takes.
//
// Line coverage says a line ran. A line with a condition runs whichever way the condition
// goes, and the other way can stay untested for good: in the reference app a "red" state
// was unreachable (the slider's maximum was below the base it had to cross), a percentage
// was always zero, an empty-state skeleton never rendered — every line covered, every
// branch not. The lcov format already carries branches (`BRDA`), whatever the language
// that wrote it, so the gate reads them from the signal the coverage ingestion keeps.
//
// A branch never taken on a line where the mutation run also found mutants no test ran is
// reported apart, as LIKELY DEAD: the tests do not merely skip it, nothing they do reaches
// it. That always fails, whatever the floor.
//
// A branch the project leaves untested on purpose is declared on its line, or the line
// above, with `@no-branch: <why>`.
func checkBranchCoverage(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode {
		return Skip, i18n.T("gate.branch_coverage.skip_not_code")
	}
	if v, msg, ok := coverageAbsence(n); ok {
		return v, msg
	}
	if n.Signal == nil || n.Signal.TotalLines == 0 {
		return Pending, i18n.T("gate.no_line_coverage")
	}
	if n.SignalStale() {
		return Pending, i18n.T("gate.stale_coverage")
	}
	if n.Signal.BranchTotal == 0 {
		return Skip, i18n.T("gate.branch_coverage.skip_no_branches")
	}
	waived := noBranchLines(content)
	byLine := map[int]int{}
	for _, id := range strings.Fields(n.Signal.BranchMissed) {
		line, _ := strconv.Atoi(strings.SplitN(id, ":", 2)[0])
		if !waived[line] {
			byLine[line]++
		}
	}
	missed := 0
	for _, c := range byLine {
		missed += c
	}
	pct := float64(n.Signal.BranchTotal-missed) / float64(n.Signal.BranchTotal) * 100
	floor := 100.0
	if e := gateEntry(cfg, "branch-coverage"); e.MinPercent != nil {
		floor = *e.MinPercent
	}
	var dead []int
	if !n.MutationStale() {
		for _, l := range n.Signal.NoCoverageLines() {
			if byLine[l] > 0 && !containsInt(dead, l) {
				dead = append(dead, l)
			}
		}
	}
	sort.Ints(dead)
	var problems []string
	if pct < floor {
		lines := make([]int, 0, len(byLine))
		for l := range byLine {
			lines = append(lines, l)
		}
		sort.Ints(lines)
		problems = append(problems, i18n.T("gate.branch_coverage.below", pct, floor, missed, n.Signal.BranchTotal, joinInts(lines)))
	}
	if len(dead) > 0 {
		problems = append(problems, i18n.T("gate.branch_coverage.dead", joinInts(dead)))
	}
	if len(problems) > 0 {
		return Fail, strings.Join(problems, "; ")
	}
	return Pass, ""
}

// noBranchRE is the waiver of a branch, with its reason: `@no-branch: <why>`.
var noBranchRE = regexp.MustCompile(`@no-branch[^\S\n]*:[^\S\n]*\S+`)

// noBranchLines are the lines whose branches are waived: the line carrying the waiver,
// and the one after it (a waiver written on the line above the condition).
func noBranchLines(content string) map[int]bool {
	out := map[int]bool{}
	for i, l := range strings.Split(content, "\n") {
		if noBranchRE.MatchString(l) {
			out[i+1], out[i+2] = true, true
		}
	}
	return out
}

func containsInt(xs []int, x int) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

func joinInts(xs []int) string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = fmt.Sprint(x)
	}
	return strings.Join(out, ", ")
}
