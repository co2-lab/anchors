package flow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const deliverLocalYAML = "version: 1\n" +
	"layers:\n" +
	"  logic:\n    pattern: \"src/**/*.ts\"\n    kind: code\n" +
	"  spec:\n    pattern: \"src/**/*.spec.md\"\n    kind: spec\n" +
	"gates:\n" +
	"  - name: always-red\n    id: ALWRD\n    on: [code]\n    blocking: false\n" +
	"    run: \"echo 'the loop drops the last page. Details follow here'; exit 1\"\n"

// deliverMap is a map where src/pricing.ts is the unit PRICX, with no mutation signal.
const deliverMap = "version: 4\nnodes:\n" +
	"  - id: src/pricing.ts\n    kind: code\n    rev: a\n    code: PRICX\n    layer: logic\n" +
	"  - id: src/pricing.spec.md\n    kind: spec\n    rev: b\n    code: PRICX\n    layer: spec\n" +
	"edges: []\n"

func runDeliver(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := newDeliverCmd()
	cmd.SetArgs(args)
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	return out, err
}

// Local mode: the record lands in `changes/`, and the declared files are confronted
// against the disk — the untouched one is named, the one in a new folder is not, the
// informative gate that fails on the unit is shown, and the missing mutation signal too.
func TestDeliverCmd_localRecordIsWrittenAndConfronted(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", deliverLocalYAML)
	writeFile(t, root, "anchors.graph.yaml", deliverMap)
	writeFile(t, root, "src/pricing.ts", "export const p = 1\n")
	writeFile(t, root, "src/quiet.ts", "export const q = 1\n")
	writeFile(t, root, "src/pricing.test.ts", "it('x', () => {})\n")
	gitRepo(t, root)
	writeFile(t, root, "src/pricing.ts", "export const p = 2\n")
	writeFile(t, root, "src/newdir/handler.ts", "export const h = 1\n")

	out, err := runDeliver(t, "--root", root, "--stage", "code", "--unit", "src/pricing.ts",
		"--file", "src/pricing.ts,src/quiet.ts,src/newdir/handler.ts",
		"--intent", "implements PRICX-B01", "--date", "2026-09-26")
	if err != nil {
		t.Fatal(err)
	}

	recs, _ := filepath.Glob(filepath.Join(root, "changes", "*.md"))
	if len(recs) != 1 {
		t.Fatalf("expected one record in changes/, got %v\n%s", recs, out)
	}
	rec, _ := os.ReadFile(recs[0])
	if !strings.Contains(string(rec), "implements PRICX-B01") || !strings.Contains(string(rec), "src/pricing.ts") {
		t.Errorf("the record must carry the unit and the intent:\n%s", rec)
	}
	if !strings.Contains(out, "delivery recorded: changes/") {
		t.Errorf("the output names the record:\n%s", out)
	}

	// the confrontation
	untouched := out[strings.Index(out, "declared files that do NOT appear in the diff"):]
	untouched = untouched[:strings.Index(untouched, "Either you did not")]
	if !strings.Contains(untouched, "src/quiet.ts") {
		t.Errorf("the committed and untouched file must be named:\n%s", out)
	}
	if strings.Contains(untouched, "src/pricing.ts") || strings.Contains(untouched, "src/newdir/handler.ts") {
		t.Errorf("a modified file and a file in a new folder were touched:\n%s", untouched)
	}
	if !strings.Contains(out, "always-red @ src/pricing.ts — the loop drops the last page.") ||
		strings.Contains(out, "Details follow here") {
		t.Errorf("the failing informative gate is shown by its first sentence:\n%s", out)
	}
	if !strings.Contains(out, "NO mutation signal ingested") {
		t.Errorf("a unit with a test and no mutation signal must be warned:\n%s", out)
	}
	if !strings.Contains(out, "the watcher is not running") || !strings.Contains(out, "you declared ZERO free decisions") {
		t.Errorf("the watcher hint and the zero-decisions note are expected:\n%s", out)
	}
	if !strings.Contains(out, "anchors work review --for src/pricing.ts") {
		t.Errorf("the next step is the review of the unit:\n%s", out)
	}
}

// With the watcher running and decisions declared, neither hint is printed; and on the
// first delivery of a unit only its spec exists, which is enough.
func TestDeliverCmd_specStageAcceptsTheSpecAsTheUnitsPiece(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", deliverLocalYAML)
	writeFile(t, root, "src/pricing.spec.md", "# spec\n")
	writeFile(t, root, ".anchors/watch.meta", "started=x\n")

	// The unit is given absolute: a relative path that does not exist under the root is
	// resolved against the working directory (common.RelTo), as a user at the root expects.
	out, err := runDeliver(t, "--root", root, "--stage", "spec", "--unit", filepath.Join(root, "src/pricing.ts"),
		"--intent", "the spec", "--date", "2026-09-26", "--decision", "a, b", "--uncovered", "c")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "recording the unit by `src/pricing.spec.md`") {
		t.Errorf("the existing piece must be named:\n%s", out)
	}
	if strings.Contains(out, "the watcher is not running") || strings.Contains(out, "ZERO free decisions") {
		t.Errorf("hints that do not apply must not be printed:\n%s", out)
	}
	if !strings.Contains(out, "could not confront the declared files") {
		t.Errorf("outside git the confrontation is said NOT to have happened:\n%s", out)
	}
}

func TestDeliverCmd_refusals(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", deliverLocalYAML)
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--unit", "src/a.ts"}, "provide --stage"},
		{[]string{"--stage", "deploy", "--unit", "src/a.ts"}, `unknown stage "deploy"`},
		{[]string{"--stage", "code", "--unit", "src/a.ts", "--intent", " "}, "provide --intent"},
		{[]string{"--stage", "code", "--unit", "src/a.ts", "--intent", "x"}, "provide --date"},
		{[]string{"--stage", "code", "--unit", "src/a.ts", "--intent", "x", "--date", "2026-09-26"}, "does not exist on disk"},
	} {
		_, err := runDeliver(t, append([]string{"--root", root}, c.args...)...)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v: want %q, got %v", c.args, c.want, err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "changes")); err == nil {
		t.Error("a refused delivery must record nothing")
	}
}

// github mode with --card: the record is a comment on THAT card, sent through stdin.
func TestDeliverCmd_githubCardReceivesTheRecord(t *testing.T) {
	root := githubProject(t)
	writeFile(t, root, "src/pricing.ts", "export const p = 1\n")
	body := filepath.Join(t.TempDir(), "body.md")
	calls := scriptedGH(t,
		ghRule{match: "issue view 77 *", out: `{"number":77,"title":"[plano] F02","body":"","state":"OPEN","labels":[{"name":"anchors"}]}`},
		ghRule{match: "issue comment 77 *", stdinTo: body},
	)
	out, err := runDeliver(t, "--root", root, "--stage", "plan", "--unit", "src/pricing.ts",
		"--intent", "delivered phase 2", "--date", "2026-09-26", "--card", "77")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(body)
	if !strings.Contains(string(b), "delivered phase 2") {
		t.Errorf("the comment must carry the record: %q (calls %v)", b, calls())
	}
	if !strings.Contains(out, "delivery recorded on issue #77 — [plano] F02") {
		t.Errorf("the output names the card:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(root, "changes")); err == nil {
		t.Error("github mode must not fall back to changes/")
	}
}

// github mode without --card: the card is found by the unit's code in the map.
func TestDeliverCmd_githubFindsTheCardByTheUnitsCode(t *testing.T) {
	root := githubProject(t)
	writeFile(t, root, "anchors.graph.yaml", deliverMap)
	writeFile(t, root, "src/pricing.ts", "export const p = 1\n")
	body := filepath.Join(t.TempDir(), "body.md")
	scriptedGH(t,
		ghRule{match: "issue list *", out: `[{"number":5,"title":"[OTHER] x"},{"number":6,"title":"[PRICX] Implementar spec — pricing"}]`},
		ghRule{match: "issue comment 6 *", stdinTo: body},
	)
	out, err := runDeliver(t, "--root", root, "--stage", "code", "--unit", "src/pricing.ts",
		"--intent", "implements PRICX-B01", "--date", "2026-09-26")
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(body); !strings.Contains(string(b), "implements PRICX-B01") {
		t.Errorf("the record must go to the PRICX card: %q", b)
	}
	if !strings.Contains(out, "delivery recorded on issue #6") {
		t.Errorf("the output names the card:\n%s", out)
	}
}

func TestDeliverCmd_githubFailuresNameTheWayOut(t *testing.T) {
	root := githubProject(t)
	writeFile(t, root, "src/pricing.ts", "export const p = 1\n")
	scriptedGH(t, ghRule{match: "issue list *", out: `[]`})

	// no map: no code
	_, err := runDeliver(t, "--root", root, "--stage", "code", "--unit", "src/pricing.ts",
		"--intent", "x", "--date", "2026-09-26")
	if err == nil || !strings.Contains(err.Error(), "could not find the CODE") || !strings.Contains(err.Error(), "--card <n>") {
		t.Errorf("without a code the command must say how to name the card, got %v", err)
	}

	// a code, and no open card with it
	writeFile(t, root, "anchors.graph.yaml", deliverMap)
	_, err = runDeliver(t, "--root", root, "--stage", "code", "--unit", "src/pricing.ts",
		"--intent", "x", "--date", "2026-09-26")
	if err == nil || !strings.Contains(err.Error(), "[PRICX]") || !strings.Contains(err.Error(), "reopen it") {
		t.Errorf("a missing card must be named with the way out, got %v", err)
	}

	// --card of a closed issue
	scriptedGH(t, ghRule{match: "issue view 77 *", out: `{"number":77,"title":"t","state":"CLOSED","labels":[{"name":"anchors"}]}`})
	_, err = runDeliver(t, "--root", root, "--stage", "code", "--unit", "src/pricing.ts",
		"--intent", "x", "--date", "2026-09-26", "--card", "77")
	if err == nil || !strings.Contains(err.Error(), "closed") {
		t.Errorf("a closed card must be refused, got %v", err)
	}
}

// The code of a unit is found by the exact file first, then by any piece of the unit.
func TestCodeOfUnit(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.graph.yaml", deliverMap)
	if got := codeOfUnit(root, "src/pricing.ts"); got != "PRICX" {
		t.Errorf("exact file: got %q", got)
	}
	if got := codeOfUnit(root, "src/pricing.feature"); got != "PRICX" {
		t.Errorf("a piece of the unit: got %q", got)
	}
	if got := codeOfUnit(root, "src/other.ts"); got != "" {
		t.Errorf("an unknown unit has no code: got %q", got)
	}
	if got := codeOfUnit(t.TempDir(), "src/pricing.ts"); got != "" {
		t.Errorf("no map, no code: got %q", got)
	}
}

// A mutation signal already ingested silences the warning; no test, no warning either.
func TestMutationNotMeasured(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "src/pricing.ts", "x\n")
	if got := mutationNotMeasured(root, "src/pricing.ts"); got != "" {
		t.Errorf("without a test the subject is another: %q", got)
	}
	writeFile(t, root, "src/pricing.test.ts", "x\n")
	writeFile(t, root, "anchors.graph.yaml", "version: 4\nnodes:\n"+
		"  - id: src/pricing.ts\n    kind: code\n    rev: a\n    signal:\n      mutants_killed: 3\nedges: []\n")
	if got := mutationNotMeasured(root, "src/pricing.ts"); got != "" {
		t.Errorf("a measured unit must not be warned: %q", got)
	}
	if got := mutationNotMeasured(root, ""); got != "" {
		t.Errorf("no unit, no warning: %q", got)
	}
}
