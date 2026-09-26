package quality

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

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

// The delta passes when nothing dropped. (A drop ends in os.Exit(1), which an in-process
// test cannot observe; the logic that decides it is the same comparison checked here.)
func TestCoverageDeltaWithoutDrop(t *testing.T) {
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
	dir := qProject(t, "version: 2\nlayers: {}\n", nil, nil)
	if _, err := runQ(t, newCoverageCmd(), "--root", dir); err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: must point at map build; got %v", err)
	}
}
