// @anchors
//   ref: DUPLC

package gate

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// checkDuplication — does this file take part in a block copied from somewhere else?
//
// It reads jscpd's REPORT, not its exit code. jscpd exits 0 with any duplication unless
// a `threshold` is configured, so a gate that read the exit code approved what it did
// not measure: in the reference app, 35 clones (1.17% of the lines), exit 0, and
// `check --all` announced the gate "clean, ready to become blocking". The report lists
// every clone with both files and their lines, and that is what this check judges.
//
// One verdict per FILE, not one for the project: the file that holds a copy fails,
// naming the other side, and every other file passes. With `--changed` only the files of
// the change are judged, so a commit answers for the clones it touches and not for old
// debt elsewhere. jscpd still scans the whole project, because a new copy only exists
// relative to its original.
//
// The calibration is the project's `.jscpd.json`, jscpd's own mechanism (`minLines`,
// `ignore`, `jscpd:ignore-start` markers), which jscpd reads by itself. Its `threshold`
// keeps jscpd's meaning: the percentage of duplicated lines the project tolerates. At or
// under it the clones are reported without failing; without it, any clone fails.
func checkDuplication(_ string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	rep, err := duplicationReport(root, g)
	if err != nil {
		return Pending, i18n.T("gate.duplication.no_report", err)
	}
	var mine []string
	for _, c := range rep.Duplicates {
		a, b := reportPath(root, c.FirstFile.Name), reportPath(root, c.SecondFile.Name)
		switch n.ID {
		case a:
			mine = append(mine, cloneLine(c.FirstFile, b, c.SecondFile, c.Lines, a == b))
		case b:
			mine = append(mine, cloneLine(c.SecondFile, a, c.FirstFile, c.Lines, a == b))
		}
	}
	if len(mine) == 0 {
		return Pass, ""
	}
	sort.Strings(mine)
	mine = dedupSorted(mine)
	list := strings.Join(mine, "\n")
	if threshold, ok := jscpdThreshold(root); ok && rep.Statistics.Total.Percentage <= threshold {
		return Diverge, i18n.T("gate.duplication.within_threshold",
			rep.Statistics.Total.Percentage, threshold, list)
	}
	return Fail, i18n.T("gate.duplication.clones", len(mine), list)
}

// jscpdReport is the part of `jscpd-report.json` the check reads.
type jscpdReport struct {
	Duplicates []struct {
		Lines      int        `json:"lines"`
		FirstFile  jscpdPlace `json:"firstFile"`
		SecondFile jscpdPlace `json:"secondFile"`
	} `json:"duplicates"`
	Statistics struct {
		Total struct {
			Percentage float64 `json:"percentage"`
		} `json:"total"`
	} `json:"statistics"`
}

// jscpdPlace is one side of a clone: the file and its first and last line.
type jscpdPlace struct {
	Name  string `json:"name"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// cloneLine describes one clone from the side of the file being judged. A clone inside
// a single file names only the line ranges.
func cloneLine(here jscpdPlace, other string, there jscpdPlace, lines int, sameFile bool) string {
	if sameFile {
		return fmt.Sprintf("  :%d-%d ≡ :%d-%d (%d lines)", here.Start, here.End, there.Start, there.End, lines)
	}
	return fmt.Sprintf("  :%d-%d ≡ %s:%d-%d (%d lines)", here.Start, here.End, other, there.Start, there.End, lines)
}

// reportPath turns a path of the report into a node ID: relative to the root, with
// forward slashes. jscpd writes paths as it walked them, relative to where it ran.
func reportPath(root, p string) string {
	if filepath.IsAbs(p) {
		if rel, err := filepath.Rel(root, p); err == nil {
			p = rel
		}
	}
	return filepath.ToSlash(filepath.Clean(p))
}

// jscpdThreshold reads the `threshold` the project declares in `.jscpd.json`.
func jscpdThreshold(root string) (float64, bool) {
	b, err := readFile(root, ".jscpd.json")
	if err != nil {
		return 0, false // @resilient: no `.jscpd.json` is a project that declares no tolerance
	}
	var c struct {
		Threshold *float64 `json:"threshold"`
	}
	if json.Unmarshal(b, &c) != nil || c.Threshold == nil {
		return 0, false // @resilient: a file jscpd itself cannot read gives no tolerance either
	}
	return *c.Threshold, true
}

// jscpdPackage is the jscpd release the gate runs, pinned. Unpinned, `npx` fetched the
// newest, and jscpd 5.4.0 (2026-09-30) needs a per-platform package `npx` does not install:
// it printed an install hint, wrote no report, and every file of every project read "the
// duplication was not measured". A new release is adopted by changing this line, after
// running it on a real project.
const jscpdPackage = "jscpd@5.3.3"

// duplicationCommand builds the jscpd run that writes its JSON report into outDir. It is
// a variable so the tests can stand in for the tool.
var duplicationCommand = func(root, outDir string) *exec.Cmd {
	cmd := exec.Command("npx", "--yes", jscpdPackage, ".", "--reporters", "json", "--output", outDir, "--silent")
	cmd.Dir = root
	return cmd
}

// ONE jscpd RUN PER SCAN. The answer is the same for every file — the clones of the
// whole project — and the check runs once per target, as docs-fresh does; the key is
// the root and the graph, so two scans of the same session (before and after `--fix`)
// each get their own run.
var (
	dupMu    sync.Mutex
	dupRoot  string
	dupGraph *mapx.Graph
	dupRep   *jscpdReport
	dupErr   error
	dupOK    bool
)

// duplicationReport runs jscpd at most once per scan and returns its report.
func duplicationReport(root string, g *mapx.Graph) (*jscpdReport, error) {
	dupMu.Lock()
	defer dupMu.Unlock()
	if dupOK && dupRoot == root && dupGraph == g {
		return dupRep, dupErr
	}
	rep, err := runJscpd(root)
	dupRoot, dupGraph, dupRep, dupErr, dupOK = root, g, rep, err, true
	return rep, err
}

// runJscpd runs the tool and parses its report. The exit code is not read: with a
// threshold exceeded jscpd exits 1 and still writes the report, and that is a result, not
// a failure. Only a missing or unreadable report is.
func runJscpd(root string) (*jscpdReport, error) {
	dir, err := os.MkdirTemp("", "anchors-jscpd-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	out, _ := duplicationCommand(root, dir).CombinedOutput()
	b, err := os.ReadFile(filepath.Join(dir, "jscpd-report.json"))
	if err != nil {
		return nil, fmt.Errorf("jscpd wrote no report: %s", lastLine(string(out)))
	}
	var rep jscpdReport
	if err := json.Unmarshal(b, &rep); err != nil {
		return nil, fmt.Errorf("jscpd report unreadable: %w", err)
	}
	return &rep, nil
}

// lastLine is the last non-empty line of a tool's output: where it says why it stopped.
func lastLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

// resetDuplicationCache clears the memory between scans, for the tests.
func resetDuplicationCache() {
	dupMu.Lock()
	defer dupMu.Unlock()
	dupOK, dupRoot, dupGraph, dupRep, dupErr = false, "", nil, nil, nil
}
