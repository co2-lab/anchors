package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// DRY-RUN by default: the plan is shown, classified by kind of occurrence, and nothing
// on disk changes.
func TestRecode_dryRunWritesNothing(t *testing.T) {
	root := fixtureProject(t)
	before, _ := os.ReadFile(filepath.Join(root, "src/login.spec.md"))

	out := runCmd(t, newRecodeCmd(), "login", "signn", "--root", root)
	if !strings.Contains(out, "recode LOGIN → SIGNN:") {
		t.Errorf("the codes must be upper-cased and named:\n%s", out)
	}
	for _, want := range []string{"src/login.spec.md", "header", "scenario-code", "(dry-run — nothing was written"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in the plan:\n%s", want, out)
		}
	}
	after, _ := os.ReadFile(filepath.Join(root, "src/login.spec.md"))
	if string(before) != string(after) {
		t.Error("the dry-run changed the spec")
	}
}

// --apply rewrites every surface and rebuilds the map from the headers.
func TestRecode_applyRewritesAndRebuildsTheMap(t *testing.T) {
	root := fixtureProject(t)
	out := runCmd(t, newRecodeCmd(), "LOGIN", "SIGNN", "--apply", "--root", root)
	if !strings.Contains(out, "file(s) rewritten") || !strings.Contains(out, "map rebuilt (5 nodes)") {
		t.Errorf("unexpected output of the apply:\n%s", out)
	}
	spec, _ := os.ReadFile(filepath.Join(root, "src/login.spec.md"))
	if strings.Contains(string(spec), "LOGIN") || !strings.Contains(string(spec), "code: SIGNN") ||
		!strings.Contains(string(spec), "SIGNN-B01") {
		t.Errorf("the spec still carries the old code:\n%s", spec)
	}
	test, _ := os.ReadFile(filepath.Join(root, "src/login.test.ts"))
	if !strings.Contains(string(test), "SIGNN-B01") {
		t.Errorf("the test's scenario code was not renamed:\n%s", test)
	}
	for _, n := range loadMap(t, root).Nodes {
		if n.ID == "src/login.spec.md" && n.Code != "SIGNN" {
			t.Errorf("the rebuilt map still has code %q on the spec", n.Code)
		}
	}
}

func TestRecode_refusals(t *testing.T) {
	root := fixtureProject(t)
	if _, err := runCmdErr(newRecodeCmd(), t, "NOPEX", "SIGNN", "--root", root); err == nil {
		t.Error("recoding a code that does not exist must be refused")
	}
	if _, err := runCmdErr(newRecodeCmd(), t, "LOGIN", "SIGNN", "--root", t.TempDir()); err == nil ||
		!strings.Contains(err.Error(), "anchors.yaml") {
		t.Errorf("no config: got %v", err)
	}
}
