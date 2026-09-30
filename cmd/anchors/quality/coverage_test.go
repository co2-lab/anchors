package quality

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// A spec CITES neighbouring rules in prose, and that is the project's style — not a
// defect. The defect was treating the citation as a DECLARATION.
//
// MEASURED in the reference app: 37 of the 55 nodes with `proven_codes` carried a code of another
// unit. `InfraList.spec.md` cites `QSCOP-B02` ONCE, in prose, no InfraList test mentions
// it — and still the map claimed InfraList proved it. That is worse than a gap: a gap
// shows up in the report, and this one vanishes.
//
// The cause is one of address: `CodesInCase` was written for the NAME of a test case,
// where every code present IS the case's code. Applied to the WHOLE spec file it also
// collects what the prose mentions.
func TestCodeCitedInProseIsNotDeclared(t *testing.T) {
	t.Run("CVCMC-B02: A code another unit owns, cited in prose, is not a declared scenario", func(t *testing.T) {})
	dir := t.TempDir()
	spec := filepath.Join(dir, "InfraList.spec.md")
	content := `# InfraList

## INLSN-B01 — the list sorts by severity

The screen uses the effective one (` + "`QSCOP-B02`" + `), not the requested one.

## INLSN-B02 — the verdict comes from the integration
`
	if err := os.WriteFile(spec, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	codes, err := codesInFileOfUnit(spec, "INLSN")
	if err != nil {
		t.Fatal(err)
	}

	hasB01, hasB02, foreign := false, false, ""
	for _, c := range codes {
		switch c {
		case "INLSN-B01":
			hasB01 = true
		case "INLSN-B02":
			hasB02 = true
		default:
			foreign = c
		}
	}

	if !hasB01 || !hasB02 {
		t.Errorf("the unit's OWN rules vanished: %v — the filter cut too much", codes)
	}
	if foreign != "" {
		t.Errorf("`%s` is cited in PROSE and came in as declared — the map would claim "+
			"this unit proved its neighbour's rule, and none of its tests touches it", foreign)
	}
}

// Without the unit's code there is nothing to filter by, and cutting everything would be
// worse than not filtering: the node would lose its own scenarios. It passes through.
func TestWithoutUnitCodeNothingIsFiltered(t *testing.T) {
	t.Run("CVCMC-B03: A spec with no unit code keeps every code it declares", func(t *testing.T) {})
	dir := t.TempDir()
	spec := filepath.Join(dir, "x.spec.md")
	if err := os.WriteFile(spec, []byte("## ABCDE-B01\n## FGHIJ-B02\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	codes, err := codesInFileOfUnit(spec, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(codes) != 2 {
		t.Errorf("without the unit code the filter should pass through, got %v", codes)
	}
}

// coverageGraph: one spec with a gap, one fully proven, code files with and without
// line/mutation signals.
func coverageGraph() *mapx.Graph {
	return &mapx.Graph{Nodes: []mapx.Node{
		{ID: "b.spec.md", Kind: mapx.KindSpec, Code: "BBBBB", Rev: "now", Signal: &mapx.TestSignal{
			AtRev: "before", ProvenCodes: []string{"BBBBB-B01"},
		}},
		{ID: "c.spec.md", Kind: mapx.KindSpec, Code: "CCCCC", Signal: &mapx.TestSignal{
			ProvenCodes: []string{"CCCCC-B01"},
		}},
		{ID: "low.go", Kind: mapx.KindCode, Rev: "r2", Signal: &mapx.TestSignal{
			AtRev: "r1", TotalLines: 10, LineCoverage: 30, PrevLineCoverage: 50,
			MutantsKilled: 1, MutantsSurvived: 3, MutationScore: 25,
		}},
		{ID: "high.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{
			TotalLines: 10, LineCoverage: 95, PrevLineCoverage: 80,
		}},
		{ID: "bare.go", Kind: mapx.KindCode},
	}}
}

func coverageFiles() map[string]string {
	return map[string]string{
		"b.spec.md": "# B\n\n## BBBBB-B01 — one\n\n## BBBBB-B02 — two\n",
		"c.spec.md": "# C\n\n## CCCCC-B01 — one\n",
		"e.spec.md": "# E\n\nno codes here\n",
	}
}

func TestCoveragePanoramaReportsEachQuestion(t *testing.T) {
	t.Run("CVCMC-B06: The panorama answers by scenario, by line and by mutation", func(t *testing.T) {})
	dir := qProject(t, "version: 2\nlayers: {}\n", coverageFiles(), coverageGraph())

	out, err := runQ(t, newCoverageCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"  b.spec.md — 1/2 unproven: [BBBBB-B02]\n",
		"== line coverage (< 70%) ==\n  low.go — 30% ⚠stale\n",
		"  low.go — 25% (3 survivor(s)) ⚠stale\n",
		"  (measured: 1 of 3 code file(s))",
		"summary: 1 spec(s) with an unproven scenario; 1 file(s) < 70% of lines; 3 surviving mutant(s) in 1 measured file(s); 1 stale signal(s) (re-ingest)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "c.spec.md —") || strings.Contains(out, "high.go —") {
		t.Errorf("a fully proven spec and a covered file are not gaps:\n%s", out)
	}
}

// "nothing measured" and "everything fine" are opposite conclusions; the panorama must
// not print the same thing for both.
func TestCoveragePanoramaSaysWhenNothingWasMeasured(t *testing.T) {
	t.Run("CVCMC-I01: Nothing measured is never reported as nothing below the threshold", func(t *testing.T) {})
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, &mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode}}})

	out, err := runQ(t, newCoverageCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"(none — every declared scenario has a green test, OR nothing was ingested)",
		"⚠ no line coverage ingested",
		"⚠ no mutation signal ingested (1 code file(s))",
		"line coverage NOT measured; mutation NOT measured",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestCoveragePanoramaAllMeasuredAndAbove(t *testing.T) {
	t.Run("CVCMC-B07: Every measured file above the threshold is said with the count, survivors included", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "a.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{TotalLines: 5, LineCoverage: 100, MutantsKilled: 4, MutationScore: 100}},
		{ID: "b.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{TotalLines: 5, LineCoverage: 80, MutantsKilled: 4, MutantsSurvived: 1, MutationScore: 80}},
	}}
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, g)

	out, err := runQ(t, newCoverageCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "✓ none of the 2 measured file(s) below the threshold") {
		t.Errorf("measured and above: say so, with the count:\n%s", out)
	}
	// above the threshold is tolerance, not approval: survivors are still reported
	if !strings.Contains(out, "~ 2 measured file(s) above the threshold, but with 1 surviving mutant(s)") {
		t.Errorf("survivors above the threshold must still be named:\n%s", out)
	}

	g.Nodes[1].Signal.MutantsSurvived = 0
	dir = qProject(t, "version: 2\nlayers: {}\n", nil, g)
	out, _ = runQ(t, newCoverageCmd(), "--root", dir)
	if !strings.Contains(out, "✓ 2 measured file(s), no mutant survived") {
		t.Errorf("no survivor at all:\n%s", out)
	}
}

func TestCoverageForOneSpec(t *testing.T) {
	t.Run("CVCMC-B01: The scenarios of one spec are listed as proven or not", func(t *testing.T) {})
	t.Run("CVCMC-B04: A spec changed since ingestion flags its signal as stale", func(t *testing.T) {})
	t.Run("CVCMC-B05: A spec with no scenario code has nothing to cover", func(t *testing.T) {})
	t.Run("CVCMC-E02: A spec outside the map is refused", func(t *testing.T) {})
	dir := qProject(t, "version: 2\nlayers: {}\n", coverageFiles(), coverageGraph())

	out, err := runQ(t, newCoverageCmd(), "--root", dir, "b.spec.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"coverage by scenario — b.spec.md:",
		"  ✓ BBBBB-B01 — proven by a green test\n",
		"  ✗ BBBBB-B02 — no passing test\n",
		"1/2 scenario(s) proven  ⚠ STALE SIGNAL",
		"green tests missing for: [BBBBB-B02]",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}

	g := coverageGraph()
	g.Nodes = append(g.Nodes, mapx.Node{ID: "e.spec.md", Kind: mapx.KindSpec})
	dir = qProject(t, "version: 2\nlayers: {}\n", coverageFiles(), g)
	out, err = runQ(t, newCoverageCmd(), "--root", dir, "e.spec.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "e.spec.md declares no scenario codes") {
		t.Errorf("a spec without codes has nothing to cover:\n%s", out)
	}

	if _, err := runQ(t, newCoverageCmd(), "--root", dir, "ghost.spec.md"); err == nil || !strings.Contains(err.Error(), "is not in the map") {
		t.Errorf("a spec outside the map must be refused; got %v", err)
	}
}

// The delta passes when nothing dropped. (A drop ends in os.Exit(1); it is observed in a
// child process by TestCoverageDeltaWithDropExitsOne.)
func TestCoverageDeltaWithoutDrop(t *testing.T) {
	t.Run("CVCMC-B12: The delta with no drop says so and counts the improvements", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "high.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{TotalLines: 10, LineCoverage: 95, PrevLineCoverage: 80}},
		{ID: "same.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{TotalLines: 10, LineCoverage: 50, PrevLineCoverage: 50}},
		// no baseline: not compared
		{ID: "new.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{TotalLines: 10, LineCoverage: 10}},
	}}
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, g)

	out, err := runQ(t, newCoverageCmd(), "--root", dir, "--delta")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "✓ no file lost coverage") || !strings.Contains(out, "(1 file(s) improved)") {
		t.Errorf("one improved, one unchanged, one without baseline:\n%s", out)
	}
	if strings.Contains(out, "⚠") {
		t.Errorf("nothing dropped, nothing to warn:\n%s", out)
	}
}

const coverageDiffText = `diff --git a/pkg/a.go b/pkg/a.go
--- a/pkg/a.go
+++ b/pkg/a.go
@@ -1,2 +1,4 @@
 package pkg
+func A() {}
+func B() {}
+// just a comment
diff --git a/README.md b/README.md
--- a/README.md
+++ b/README.md
@@ -1 +1,2 @@
 # r
+more
`

func TestCoverageDiffCrossesTheDiffWithTheLCOV(t *testing.T) {
	t.Run("CVCMC-B09: The changed lines are crossed with the coverage report", func(t *testing.T) {})
	t.Run("CVCMC-X01: The diff coverage needs no map", func(t *testing.T) {})
	// the lcov carries an absolute prefix: the match is by path suffix.
	// line 2 covered, line 3 not, line 4 (the comment) not instrumented.
	lcov := "SF:/build/src/pkg/a.go\nDA:2,1\nDA:3,0\nend_of_record\n"
	dir := qProject(t, "", map[string]string{"change.diff": coverageDiffText, "cov.info": lcov}, nil)

	out, err := runQ(t, newCoverageCmd(), "--root", dir, "--diff-file", dir+"/change.diff", "--lcov", dir+"/cov.info", "--threshold", "50")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"diff coverage:\n",
		"  pkg/a.go — 2 instrumented changed line(s), 50% covered — NO test: [3]\n",
		"diff coverage: 50% (1/2 changed lines covered)",
		"✓ what you changed is covered",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "README.md") {
		t.Errorf("a changed file without coverage is not code — it is skipped:\n%s", out)
	}
}

func TestCoverageDiffGuards(t *testing.T) {
	t.Run("CVCMC-B10: A diff with no instrumented line has nothing to cover", func(t *testing.T) {})
	t.Run("CVCMC-E03: The diff coverage without a coverage report is refused", func(t *testing.T) {})
	t.Run("CVCMC-E06: A diff file that cannot be read fails the diff coverage", func(t *testing.T) {})
	dir := qProject(t, "", map[string]string{"change.diff": coverageDiffText, "cov.info": "SF:other.go\nDA:1,1\nend_of_record\n"}, nil)

	if _, err := runQ(t, newCoverageCmd(), "--root", dir, "--diff-file", dir+"/change.diff"); err == nil || !strings.Contains(err.Error(), "--lcov <file> is mandatory") {
		t.Errorf("--diff without --lcov must be refused; got %v", err)
	}

	out, err := runQ(t, newCoverageCmd(), "--root", dir, "--diff-file", dir+"/change.diff", "--lcov", dir+"/cov.info")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "(no instrumented code line in the diff — nothing to cover)") {
		t.Errorf("no changed file has coverage: nothing to judge:\n%s", out)
	}

	if _, err := runQ(t, newCoverageCmd(), "--root", dir, "--diff-file", dir+"/missing.diff", "--lcov", dir+"/cov.info"); err == nil {
		t.Error("an unreadable diff must fail")
	}
}

func TestCoverageWithoutMapFails(t *testing.T) {
	t.Run("CVCMC-E01: The coverage command without a map points at the map build", func(t *testing.T) {})
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, nil)
	if _, err := runQ(t, newCoverageCmd(), "--root", dir); err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: must point at map build; got %v", err)
	}
}

// runCoverageInChild runs the coverage command in a child copy of the test binary: a
// failing verdict ends in os.Exit(1), which only a separate process can observe.
func runCoverageInChild(t *testing.T, testName string, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$")
	cmd.Env = append(os.Environ(), "ANCHORS_COVERAGE_CHILD_ARGS="+strings.Join(args, "\x1f"))
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode(), string(out)
	}
	if err != nil {
		t.Fatal(err)
	}
	return 0, string(out)
}

// childCoverage is the child side: it runs the command with the parent's args and returns
// true, or returns false in the parent.
func childCoverage() bool {
	raw := os.Getenv("ANCHORS_COVERAGE_CHILD_ARGS")
	if raw == "" {
		return false
	}
	cmd := newCoverageCmd()
	cmd.SetArgs(strings.Split(raw, "\x1f"))
	_ = cmd.Execute()
	return true
}

func TestCoverageDiffBelowThresholdExitsOne(t *testing.T) {
	if childCoverage() {
		return
	}
	t.Run("CVCMC-B11: Changed lines covered below the threshold fail the diff coverage", func(t *testing.T) {})
	lcov := "SF:pkg/a.go\nDA:2,1\nDA:3,0\nend_of_record\n"
	dir := qProject(t, "", map[string]string{"change.diff": coverageDiffText, "cov.info": lcov}, nil)
	code, out := runCoverageInChild(t, "TestCoverageDiffBelowThresholdExitsOne",
		"--root", dir, "--diff-file", dir+"/change.diff", "--lcov", dir+"/cov.info", "--threshold", "70")
	if code != 1 || !strings.Contains(out, "✗ below the threshold (70%)") {
		t.Errorf("50%% of the changed lines against 70%%: exit 1 with the verdict; got exit %d:\n%s", code, out)
	}
}

func TestCoverageDeltaWithDropExitsOne(t *testing.T) {
	if childCoverage() {
		return
	}
	t.Run("CVCMC-B13: A file that lost line coverage fails the delta", func(t *testing.T) {})
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "low.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{TotalLines: 10, LineCoverage: 30, PrevLineCoverage: 50}},
		{ID: "high.go", Kind: mapx.KindCode, Signal: &mapx.TestSignal{TotalLines: 10, LineCoverage: 95, PrevLineCoverage: 80}},
	}}
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, g)
	code, out := runCoverageInChild(t, "TestCoverageDeltaWithDropExitsOne", "--root", dir, "--delta")
	if code != 1 {
		t.Errorf("a drop must exit 1, got %d:\n%s", code, out)
	}
	for _, want := range []string{"⚠ low.go — 50% → 30% (-20)", "✗ 1 file(s) lost coverage (worst: -20 points)"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "high.go —") {
		t.Errorf("a file that improved is not listed as a drop:\n%s", out)
	}
}

func TestCoveragePanoramaCapsEachSection(t *testing.T) {
	t.Run("CVCMC-B08: Each panorama section shows at most fifteen entries and counts the rest", func(t *testing.T) {})
	g := &mapx.Graph{}
	for i := 0; i < 17; i++ {
		g.Nodes = append(g.Nodes, mapx.Node{ID: fmt.Sprintf("f%02d.go", i), Kind: mapx.KindCode,
			Signal: &mapx.TestSignal{TotalLines: 10, LineCoverage: 10, MutantsKilled: 1, MutantsSurvived: 1, MutationScore: 50}})
	}
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, g)
	out, err := runQ(t, newCoverageCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(out, "% ("); n != 15 {
		t.Errorf("the mutation section shows %d entries, want 15:\n%s", n, out)
	}
	if n := strings.Count(out, "… and 2 more file(s) below the threshold"); n != 2 {
		t.Errorf("both the line and the mutation sections must count the 2 not shown (found %d):\n%s", n, out)
	}
	if !strings.Contains(out, "17 file(s) < 70% of lines") {
		t.Errorf("the summary counts every file, not only the shown ones:\n%s", out)
	}
}

func TestCoverageDiffUnreadableLCOVFails(t *testing.T) {
	t.Run("CVCMC-E04: A coverage report that cannot be read fails the diff coverage", func(t *testing.T) {})
	dir := qProject(t, "", map[string]string{"change.diff": coverageDiffText}, nil)
	_, err := runQ(t, newCoverageCmd(), "--root", dir, "--diff-file", dir+"/change.diff", "--lcov", dir+"/missing.info")
	if err == nil || !strings.Contains(err.Error(), "parse lcov") {
		t.Errorf("an unreadable lcov must fail naming it; got %v", err)
	}
}

func TestCoverageDiffOutsideGitExplains(t *testing.T) {
	t.Run("CVCMC-E05: A git diff outside a repository explains why and offers the diff file", func(t *testing.T) {})
	dir := qProject(t, "", map[string]string{"cov.info": "SF:a.go\nDA:1,1\nend_of_record\n"}, nil)
	_, err := runQ(t, newCoverageCmd(), "--root", dir, "--diff", "main", "--lcov", dir+"/cov.info")
	if err == nil || !strings.Contains(err.Error(), "--diff-file") {
		t.Errorf("outside git the diff must be refused pointing at --diff-file; got %v", err)
	}
}

// The diff coverage prints the same report on every run: the changed files in path order,
// and a changed file whose path suffix matches two coverage entries always crossed with
// the same one (the first in path order). Both came from ranging over maps.
func TestCoverageDiffIsDeterministic(t *testing.T) {
	t.Run("CVCMC-B14: The diff coverage lists the changed files in path order, and resolves a path matching several coverage entries always to the same one", func(t *testing.T) {})
	var diff, lcov strings.Builder
	names := []string{"pkg/f", "pkg/b", "pkg/e", "pkg/a", "pkg/d", "pkg/c"}
	for _, n := range names {
		fmt.Fprintf(&diff, "diff --git a/%[1]s.go b/%[1]s.go\n--- a/%[1]s.go\n+++ b/%[1]s.go\n@@ -1 +1,2 @@\n package pkg\n+func X() {}\n", n)
		fmt.Fprintf(&lcov, "SF:/build/%s.go\nDA:2,1\nend_of_record\n", n)
	}
	// one more changed file, whose suffix eight coverage entries share: only the first in
	// path order covers the line
	diff.WriteString("diff --git a/z.go b/z.go\n--- a/z.go\n+++ b/z.go\n@@ -1 +1,2 @@\n package z\n+func Z() {}\n")
	lcov.WriteString("SF:/a/z.go\nDA:2,1\nend_of_record\n")
	for _, d := range []string{"b", "c", "d", "e", "f", "g", "h"} {
		fmt.Fprintf(&lcov, "SF:/%s/z.go\nDA:2,0\nend_of_record\n", d)
	}
	dir := qProject(t, "", map[string]string{"change.diff": diff.String(), "cov.info": lcov.String()}, nil)

	first := ""
	for i := 0; i < 20; i++ {
		out, err := runQ(t, newCoverageCmd(), "--root", dir, "--diff-file", dir+"/change.diff", "--lcov", dir+"/cov.info", "--threshold", "0")
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			first = out
			last := -1
			for _, n := range []string{"pkg/a", "pkg/b", "pkg/c", "pkg/d", "pkg/e", "pkg/f", "z"} {
				at := strings.Index(out, "  "+n+".go — ")
				if at <= last {
					t.Fatalf("the changed files are not in path order (%s):\n%s", n, out)
				}
				last = at
			}
			if !strings.Contains(out, "  z.go — 1 instrumented changed line(s), 100% covered") {
				t.Errorf("z.go is crossed with the first matching entry in path order (/a/z.go):\n%s", out)
			}
			continue
		}
		if out != first {
			t.Fatalf("run %d printed a different report:\n%s\n--- first:\n%s", i, out, first)
		}
	}
}

func TestCoverage_countsARuleByAllItsVariants(t *testing.T) {
	t.Run("CVCMC-B15: The report counts a rule proven by all its variants", func(t *testing.T) {})
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "v.spec.md", Kind: mapx.KindSpec, Code: "VVVVX", Signal: &mapx.TestSignal{ProvenCodes: []string{"VVVVX-B01#01"}}},
			{ID: "v.feature", Kind: mapx.KindFeature},
		},
		Edges: []mapx.Edge{{From: "v.spec.md", To: "v.feature", Type: mapx.EdgeCoveredBy}},
	}
	files := map[string]string{
		"v.spec.md": "<!-- @anchors\n  code: VVVVX\n-->\n# V\n\n### VVVVX-B01 — does it\n",
		"v.feature": "  @VVVVX-B01#01\n  Scenario: a\n  @VVVVX-B01#02\n  Scenario: b\n",
	}
	dir := qProject(t, "version: 2\nlayers: {}\n", files, g)
	out, err := runQ(t, newCoverageCmd(), "--root", dir, "v.spec.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "✗ VVVVX-B01 — no passing test") {
		t.Errorf("the rule with an unproven variant has no passing test:\n%s", out)
	}
}
