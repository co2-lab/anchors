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
	t.Run("DLVRE-B04: Local mode writes the record under changes", func(t *testing.T) {})
	t.Run("DLVRE-B08: The next step is the review of the unit", func(t *testing.T) {})
	t.Run("DLVRE-B09: The watcher hint appears only when the watcher is not running", func(t *testing.T) {})
	t.Run("DLVRE-B10: The zero-decisions note appears only when nothing was declared", func(t *testing.T) {})
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
	t.Run("DLVRE-B02: The first delivery of a unit is accepted by its spec", func(t *testing.T) {})
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
	if !strings.Contains(out, "`src/pricing.ts` does not exist yet — accepted because `src/pricing.spec.md`, "+
		"a piece of the same unit, exists; the record keeps the unit `src/pricing.ts`") || strings.Contains(out, "recording the unit by") {
		t.Errorf("the existing piece must be named, and the unit kept said to be the one given:\n%s", out)
	}
	if strings.Contains(out, "the watcher is not running") || strings.Contains(out, "ZERO free decisions") {
		t.Errorf("hints that do not apply must not be printed:\n%s", out)
	}
	if !strings.Contains(out, "could not confront the declared files") {
		t.Errorf("outside git the confrontation is said NOT to have happened:\n%s", out)
	}
	// The record keeps the unit given, not the piece that was found.
	recs, _ := filepath.Glob(filepath.Join(root, "changes", "*.md"))
	if len(recs) != 1 {
		t.Fatalf("expected one record, got %v", recs)
	}
	rec, _ := os.ReadFile(recs[0])
	if !strings.Contains(string(rec), "unit: src/pricing.ts\n") {
		t.Errorf("the record's unit is the one given:\n%s", rec)
	}
	// the default file is the unit, relative to the root even when given absolute
	if !strings.Contains(string(rec), "- `src/pricing.ts`") || strings.Contains(string(rec), root) {
		t.Errorf("the record's files are relative to the root:\n%s", rec)
	}
}

func TestDeliverCmd_refusals(t *testing.T) {
	t.Run("DLVRE-B01: A delivery missing its required parts is refused", func(t *testing.T) {})
	t.Run("DLVRE-I01: A refused delivery records nothing", func(t *testing.T) {})
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
	t.Run("DLVRE-B05: Github mode records on the card given", func(t *testing.T) {})
	t.Run("DLVRE-X01: Github mode never falls back to changes", func(t *testing.T) {})
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
	t.Run("DLVRE-B06: Github mode finds the card by the unit's code", func(t *testing.T) {})
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
	t.Run("DLVRE-E01: Github mode without a code for the unit fails with the way out", func(t *testing.T) {})
	t.Run("DLVRE-E02: Github mode with no open card for the code fails with the way out", func(t *testing.T) {})
	t.Run("DLVRE-E03: A closed card given with --card is refused", func(t *testing.T) {})
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
	t.Run("DLVRE-B06: Github mode finds the card by the unit's code", func(t *testing.T) {})
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

// `--decision` and `--uncovered` receive PROSE, and prose has commas.
//
// With `StringSliceVar` pflag splits the value on commas. Measured delivering a spec:
//
//	--decision "eight rules and two invariants, in the B/I letter the neighbours use"
//
// became TWO decisions in the record — the second starting with a space and no subject. Of
// five declared decisions came nine items, four of them fragments. The reviewer confronts
// each decision against the disk, and half a sentence cannot be confronted.
func TestDeliver_decisionWithCommaIsNotSplit(t *testing.T) {
	t.Run("DLVRE-B03: Prose with commas is one decision, and files split on commas", func(t *testing.T) {})
	cmd := newDeliverCmd()
	const prose = "eight rules and two invariants, in the B/I letter the neighbours of F03 use"

	if err := cmd.Flags().Parse([]string{
		"--decision", prose,
		"--uncovered", "I01 is proven by reading the write code, which does not exist yet",
	}); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"decision", "uncovered"} {
		v, err := cmd.Flags().GetStringArray(name)
		if err != nil {
			t.Fatalf("--%s is not a StringArray: %v — prose with commas will be split", name, err)
		}
		if len(v) != 1 {
			t.Errorf("--%s became %d items: %q", name, len(v), v)
		}
	}

	d, _ := cmd.Flags().GetStringArray("decision")
	if d[0] != prose {
		t.Errorf("the decision arrived altered:\n  want: %q\n  got:  %q", prose, d[0])
	}
}

// `--file` STILL splits: a file path has no comma, and there the split is a real
// convenience (`--file a.ts,b.ts`). The fix was not to turn everything into StringArray.
func TestDeliver_fileStillSplitsOnComma(t *testing.T) {
	cmd := newDeliverCmd()
	if err := cmd.Flags().Parse([]string{"--file", "a.ts,b.ts"}); err != nil {
		t.Fatal(err)
	}
	f, err := cmd.Flags().GetStringSlice("file")
	if err != nil {
		t.Fatalf("--file stopped being a StringSlice: %v", err)
	}
	if len(f) != 2 || f[0] != "a.ts" || f[1] != "b.ts" {
		t.Errorf("--file a.ts,b.ts became %q", f)
	}
}

// A vendored pipeline has no local code and no card. In `github` mode the delivery used to
// look the unit's code up in the map and refuse — measured in the reference app on
// `anchors-claim.yml`. It must now succeed without demanding a code, and without falling
// back to `changes/`, which is the local mode's mechanism.
func TestDeliver_upstreamPipelineDemandsNoCode(t *testing.T) {
	t.Run("DLVRE-B07: A vendored pipeline in github mode records nothing", func(t *testing.T) {})
	root := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("anchors.yaml", "version: 2\nworkflow:\n  mode: github\n  repo: acme/app\n  labels: [anchors]\n"+
		"layers:\n  workflow:\n    pattern: \".github/workflows/*.yml\"\n    kind: code\n")
	unit := ".github/workflows/anchors-claim.yml"
	write(unit, "# anchors:template — generated by anchors.\non: push\n")

	cmd := newDeliverCmd()
	cmd.SetArgs([]string{"--root", root, "--stage", "code", "--unit", filepath.Join(root, unit),
		"--intent", "doctor --fix refreshed the pipeline", "--date", "2026-09-23"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatalf("an upstream-owned pipeline must not demand a local code: %v", err)
	}
	if !strings.Contains(out, "nothing to record") || !strings.Contains(out, "owned upstream") {
		t.Errorf("the output must say why nothing was recorded:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(root, "changes")); err == nil {
		t.Error("in github mode the record must not fall back to changes/")
	}
}
