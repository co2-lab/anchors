// @anchors
//   code: RCTSR
//   ref: RCDEO

package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// DRY-RUN by default: the plan is shown, classified by kind of occurrence, and nothing
// on disk changes.
func TestRecode_dryRunWritesNothing(t *testing.T) {
	t.Run("RCDEO-B01: Codes typed in lower case are renamed in upper case", func(t *testing.T) {})
	t.Run("RCDEO-B02: The plan counts each file's occurrences by kind", func(t *testing.T) {})
	t.Run("RCDEO-B03: A recode without the apply switch writes nothing", func(t *testing.T) {})
	root := fixtureProject(t)
	before, _ := os.ReadFile(filepath.Join(root, "src/login.spec.md"))

	out := runCmd(t, newRecodeCmd(), "login", "signn", "--root", root)
	if !strings.Contains(out, "recode LOGIN → SIGNN: 2 file(s), 3 content substitution(s)") {
		t.Errorf("the codes must be upper-cased and named:\n%s", out)
	}
	if !strings.Contains(out, "  src/login.spec.md\n      header         1\n      scenario-code  1\n") {
		t.Errorf("the spec's occurrences must be counted by kind, header first:\n%s", out)
	}
	if strings.Contains(out, "bare-ref") {
		t.Errorf("a kind with no occurrence must be omitted:\n%s", out)
	}
	if !strings.Contains(out, "(dry-run — nothing was written") {
		t.Errorf("the dry run is not announced:\n%s", out)
	}
	after, _ := os.ReadFile(filepath.Join(root, "src/login.spec.md"))
	if string(before) != string(after) {
		t.Error("the dry-run changed the spec")
	}
}

// --apply rewrites every surface and rebuilds the map from the headers.
func TestRecode_applyRewritesAndRebuildsTheMap(t *testing.T) {
	t.Run("RCDEO-B04: Applying rewrites the header and the scenario codes of the spec and the test", func(t *testing.T) {})
	t.Run("RCDEO-B05: Applying rebuilds the map from the headers", func(t *testing.T) {})
	t.Run("RCDEO-I01: After applying, neither the spec nor the map carries the old code", func(t *testing.T) {})
	t.Run("RCDEO-X01: The map is rebuilt from the files, not edited as text", func(t *testing.T) {})
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
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		t.Fatal(err)
	}
	files, err := scan.Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	revOnDisk := map[string]string{}
	for _, f := range files {
		revOnDisk[f.Path] = f.Rev
	}
	for _, n := range loadMap(t, root).Nodes {
		if n.Code == "LOGIN" {
			t.Errorf("the rebuilt map still has the old code on %s", n.ID)
		}
		if n.ID == "src/login.spec.md" {
			if n.Code != "SIGNN" {
				t.Errorf("the rebuilt map has code %q on the spec", n.Code)
			}
			if want := revOnDisk[n.ID]; n.Rev != want {
				t.Errorf("the spec node carries rev %q, the rewritten file is %q — the map was not rebuilt from the files", n.Rev, want)
			}
		}
	}
}

func TestRecode_refusals(t *testing.T) {
	t.Run("RCDEO-E01: A project with no configuration is refused naming the configuration file", func(t *testing.T) {})
	t.Run("RCDEO-E02: Recoding a code no file carries is refused", func(t *testing.T) {})
	t.Run("RCDEO-E04: A recode with a single code is refused", func(t *testing.T) {})
	root := fixtureProject(t)
	if _, err := runCmdErr(newRecodeCmd(), t, "NOPEX", "SIGNN", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "does not appear in any file") {
		t.Errorf("recoding a code that does not exist must be refused: %v", err)
	}
	if _, err := runCmdErr(newRecodeCmd(), t, "LOGIN", "SIGNN", "--root", t.TempDir()); err == nil ||
		!strings.Contains(err.Error(), "anchors.yaml") {
		t.Errorf("no config: got %v", err)
	}
	if _, err := runCmdErr(newRecodeCmd(), t, "LOGIN", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "accepts 2 arg") {
		t.Errorf("a single code: got %v", err)
	}
}

// The apply stops at the first file it cannot write. When some were already rewritten,
// the command must say so: the project is half converted.
func TestRecode_partialApplySaysHalfConverted(t *testing.T) {
	t.Run("RCDEO-E03: A write failure after some files changed says the project is half converted", func(t *testing.T) {})
	root := fixtureProject(t)
	locked := filepath.Join(root, "src/login.test.ts")
	if err := os.Chmod(locked, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o644) })
	if f, err := os.OpenFile(locked, os.O_WRONLY, 0); err == nil {
		f.Close()
		t.Skip("the file is writable despite its mode (running as root?)")
	}
	out, err := runCmdErr(newRecodeCmd(), t, "LOGIN", "SIGNN", "--apply", "--root", root)
	if err == nil {
		t.Fatal("a write failure must fail the recode")
	}
	if !strings.Contains(out, "1 file(s) had ALREADY been changed") || !strings.Contains(out, "half converted") {
		t.Errorf("the partial state was not reported:\n%s", out)
	}
}

// The rebuild after --apply must keep what the map knew beyond the headers: the
// judgments and stamps of the edges that survive, and the flow graph that only
// `flow build` fills. It used to build the map from scratch, and a recode silently
// erased every AI judgment and the whole flow.
func TestRecode_applyKeepsJudgmentsAndFlow(t *testing.T) {
	t.Run("RCDEO-B06: Applying keeps the judgments, stamps and flow of the previous map", func(t *testing.T) {})
	root := fixtureProject(t)
	g := loadMap(t, root)
	judged := -1
	for i, e := range g.Edges {
		if e.From == "src/login.spec.md" || e.To == "src/login.spec.md" {
			judged = i
			break
		}
	}
	if judged < 0 {
		t.Fatalf("the fixture has no edge touching the spec: %+v", g.Edges)
	}
	g.Edges[judged].Julgamentos = []mapx.Judgment{{Gate: "atomic", Verdict: "ok"}}
	g.Edges[judged].Stamp = &mapx.Stamp{Verdict: "ok"}
	g.Flow = &mapx.FlowGraph{States: []mapx.FlowState{{Code: "WORKR-T01", Title: "pull", Flow: "flows/w.flow.md"}}}
	// An edge between files that do not exist: the rebuild cannot keep it, and its judgment
	// is lost — which the report must say.
	g.Edges = append(g.Edges, mapx.Edge{From: "src/ghost.ts", To: "src/ghost.spec.md", Type: g.Edges[judged].Type,
		Julgamentos: []mapx.Judgment{{Gate: "review", Verdict: "ok"}}})
	if err := mapx.Save(g, filepath.Join(root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	key := g.Edges[judged]

	out := runCmd(t, newRecodeCmd(), "LOGIN", "SIGNN", "--apply", "--root", root)
	if !strings.Contains(out, "⚠ the map LOST judgment stamp(s):") || !strings.Contains(out, "review") {
		t.Errorf("the judgment lost with the ghost edge was not reported:\n%s", out)
	}

	after := loadMap(t, root)
	if after.Flow == nil || len(after.Flow.States) != 1 {
		t.Errorf("the recode erased the flow graph: %+v", after.Flow)
	}
	found := false
	for _, e := range after.Edges {
		if e.From == key.From && e.To == key.To && e.Type == key.Type {
			found = true
			if len(e.Julgamentos) != 1 || e.Julgamentos[0].Gate != "atomic" {
				t.Errorf("the recode erased the judgment of %s → %s: %+v", e.From, e.To, e.Julgamentos)
			}
			if e.Stamp == nil || e.Stamp.Verdict != "ok" {
				t.Errorf("the recode erased the stamp of %s → %s: %+v", e.From, e.To, e.Stamp)
			}
		}
	}
	if !found {
		t.Errorf("the judged edge %s → %s is gone from the rebuilt map", key.From, key.To)
	}
}

// The project's recode dialect adds testIDs and file renames to the plan; the report
// counts both and lists each rename. Without the dialect none of these lines appears.
func TestRecode_dialectReportsTestIDsAndRenames(t *testing.T) {
	t.Run("RCDEO-B07: The plan reports the testIDs and the file renames of the project's dialect", func(t *testing.T) {})
	root := fixtureProjectWith(t, "recode:\n  test_id: lower\n  file_patterns: [\"**/{{code}}-*.png\"]\n", fixtureSpec)
	writeProjectFile(t, root, "src/login.ts", "export const login = () => <b testID=\"login-button\" />\n")
	writeProjectFile(t, root, "assets/LOGIN-icon.png", "png")

	out := runCmd(t, newRecodeCmd(), "LOGIN", "SIGNN", "--root", root)
	if !strings.Contains(out, ", 1 testID(s), 1 file(s) to rename\n\n") {
		t.Errorf("the testIDs and renames are not counted:\n%s", out)
	}
	if !strings.Contains(out, "\n  rename (git mv):\n      assets/LOGIN-icon.png → assets/SIGNN-icon.png\n") {
		t.Errorf("the rename is not listed:\n%s", out)
	}

	plain := runCmd(t, newRecodeCmd(), "LOGIN", "SIGNN", "--root", fixtureProject(t))
	if !strings.Contains(plain, "content substitution(s)\n\n") || strings.Contains(plain, "rename (git mv)") {
		t.Errorf("without the dialect, no testID or rename line may appear:\n%s", plain)
	}
}

// The dialect's testID prefix is absent but the unit's files carry testIDs with another
// prefix: a previous manual rename. The recode does not guess, and says so.
func TestRecode_dialectWarnsAboutDivergentTestIDs(t *testing.T) {
	t.Run("RCDEO-B08: The plan warns when the files carry testIDs with a prefix other than the dialect's", func(t *testing.T) {})
	root := fixtureProjectWith(t, "recode:\n  test_id: lower\n", fixtureSpec)
	writeProjectFile(t, root, "src/login.test.ts", "test('LOGIN-B01: x', () => render(<b testID=\"old-button\" />))\n")

	out := runCmd(t, newRecodeCmd(), "LOGIN", "SIGNN", "--root", root)
	if !strings.Contains(out, "\n  ⚠ the expected testID prefix \"login\" was not found") {
		t.Errorf("the divergent testID prefix is not reported:\n%s", out)
	}
	if plain := runCmd(t, newRecodeCmd(), "LOGIN", "SIGNN", "--root", fixtureProject(t)); strings.Contains(plain, "⚠") {
		t.Errorf("a plan with nothing to warn about printed a warning:\n%s", plain)
	}
}

// When the very first write fails, nothing changed: the project is not half converted and
// the command must not say it is.
func TestRecode_firstWriteFailureIsNotHalfConverted(t *testing.T) {
	t.Run("RCDEO-B09: A write failure before any file changed fails without saying the project is half converted", func(t *testing.T) {})
	root := fixtureProject(t)
	locked := filepath.Join(root, "src/login.spec.md") // the first file of the plan
	if err := os.Chmod(locked, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o644) })
	if f, err := os.OpenFile(locked, os.O_WRONLY, 0); err == nil {
		f.Close()
		t.Skip("the file is writable despite its mode (running as root?)")
	}
	out, err := runCmdErr(newRecodeCmd(), t, "LOGIN", "SIGNN", "--apply", "--root", root)
	if err == nil {
		t.Fatal("a write failure must fail the recode")
	}
	if strings.Contains(out, "ALREADY") || strings.Contains(out, "half converted") {
		t.Errorf("nothing was changed, yet the project is called half converted:\n%s", out)
	}
}
