package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// THE TWO QUESTIONS of `scenario-coverage` — the same ruler as `flag-covered`.
//
// The earlier version asked only "did it pass?", and on a project that never ingested a
// report it answered Pending for everything: "nobody measured", hiding the scenarios
// nobody tested. And the static question alone would be the opposite error — a written
// test may never have run.
//
// Each state has a different fix, which is why the verdict must keep them apart: "no
// test" asks someone to write one; "written and not run" asks someone to run it.

const specWithTwoRequirements = "### CREDX-B01 — validates the limit\n\n### CREDX-B02 — refuses the balance\n"

func specNodeCoverage() mapx.Node {
	return mapx.Node{ID: "credx.spec.md", Kind: mapx.KindSpec, Code: "CREDX"}
}

func rootWithTest(t *testing.T, body string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "credx_test.go"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, &mapx.Graph{Nodes: []mapx.Node{{ID: "credx_test.go", Kind: mapx.KindTest}}}
}

// NO TEST AT ALL: it is reported even without ingested execution. It is the case the
// earlier version hid behind a Pending.
func TestScenarioCoverage_noTestIsReportedEvenWithoutIngestion(t *testing.T) {
	root, g := rootWithTest(t, "func TestNothing(t *testing.T) {}\n")

	v, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil)
	if v != Fail {
		t.Fatalf("no test names the requirements — expected Fail, got %v", v)
	}
	for _, c := range []string{"CREDX-B01", "CREDX-B02"} {
		if !strings.Contains(msg, c) {
			t.Errorf("the verdict does not name %s: %q", c, msg)
		}
	}
}

// WRITTEN AND NEVER RUN: the verdict has to say which of the two problems it is.
func TestScenarioCoverage_writtenButNotRunSaysWhichOfTheTwo(t *testing.T) {
	root, g := rootWithTest(t, "const c = \"CREDX-B01\"\nfunc TestX(t *testing.T) {}\n")

	v, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil)
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(msg, "ingest") {
		t.Errorf("the verdict does not tell written-but-not-run apart: %q", msg)
	}
	// B02 still has no test at all, and both states appear in the same verdict.
	if !strings.Contains(msg, "CREDX-B02") {
		t.Errorf("the verdict lost the requirement with no test: %q", msg)
	}
}

// A CODE IN A COMMENT does not count — a citation is a reference, not an implementation.
func TestScenarioCoverage_commentDoesNotCountAsATest(t *testing.T) {
	root, g := rootWithTest(t, "// CREDX-B01 is covered elsewhere\nfunc TestX(t *testing.T) {}\n")

	_, msg := checkScenarioCoverage(specWithTwoRequirements, specNodeCoverage(), root, g, nil)
	if strings.Contains(msg, "ingest") {
		t.Errorf("a code only in a COMMENT passed as a written test: %q", msg)
	}
}

// And what EXECUTION proved leaves the accusation, which is the original behaviour.
func TestScenarioCoverage_provenPasses(t *testing.T) {
	root, _ := rootWithTest(t, "func TestNothing(t *testing.T) {}\n")
	n := specNodeCoverage()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CREDX-B01", "CREDX-B02"}}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "credx_test.go", Kind: mapx.KindTest}}}

	if v, msg := checkScenarioCoverage(specWithTwoRequirements, n, root, g, nil); v != Pass {
		t.Errorf("both requirements proven and the verdict was %v: %s", v, msg)
	}
}
