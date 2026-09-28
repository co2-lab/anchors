package scan

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// ---- Walk ----

func TestWalk_onlyDeclaredLayersEnter(t *testing.T) {
	t.Run("RPSCR-B01: Only files in a declared layer enter the scan", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(filepath.Join(root, "src", "a.ts"), "x"))
	must(t, writeDeep(filepath.Join(root, "docs", "readme.md"), "x"))
	cfg := &config.Config{Layers: map[string]config.Layer{"code": {Pattern: "src/*.ts", Kind: "code"}}}
	files, err := Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Path != "src/a.ts" || files[0].Layer != "code" || files[0].Kind != "code" {
		t.Errorf("want only src/a.ts in layer code, got %+v", files)
	}
}

func TestWalk_honoursTheIgnoreSet(t *testing.T) {
	t.Run("RPSCR-B02: What the ignore set excludes is not scanned", func(t *testing.T) {})
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, ".gitignore"), []byte("probe*.ts\n"), 0o644))
	for _, rel := range []string{"src/a.ts", "src/probe1.ts", "src/a.ts.swp", "node_modules/x.ts"} {
		must(t, writeDeep(filepath.Join(root, rel), "x"))
	}
	cfg := &config.Config{Layers: map[string]config.Layer{"code": {Pattern: "**/*.ts", Kind: "code"}}}
	files, err := Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got := paths(files); len(got) != 1 || got[0] != "src/a.ts" {
		t.Errorf("want only src/a.ts, got %v", got)
	}
}

// A worktree or nested clone inside the tree is another checkout: its files are copies,
// and mapping them duplicated the whole project (1360 nodes from four agent worktrees).
func TestWalkSkipsNestedCheckouts(t *testing.T) {
	t.Run("RPSCR-B03: A nested checkout is not scanned", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(dir, "src"), 0o755))
	must(t, os.WriteFile(filepath.Join(dir, "src", "a.ts"), []byte("export const a = 1\n"), 0o644))
	wt := filepath.Join(dir, "tools", "worktrees", "agent-1")
	must(t, os.MkdirAll(filepath.Join(wt, "src"), 0o755))
	must(t, os.WriteFile(filepath.Join(wt, ".git"), []byte("gitdir: /elsewhere\n"), 0o644))
	must(t, os.WriteFile(filepath.Join(wt, "src", "a.ts"), []byte("export const a = 1\n"), 0o644))
	clone := filepath.Join(dir, "vendor-clone")
	must(t, os.MkdirAll(filepath.Join(clone, ".git"), 0o755))
	must(t, os.WriteFile(filepath.Join(clone, "b.ts"), []byte("export const b = 1\n"), 0o644))

	cfg := &config.Config{Layers: map[string]config.Layer{"code": {Pattern: "**/*.ts", Kind: "code"}}}
	files, err := Walk(dir, cfg)
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(paths(files), " ")
	if !strings.Contains(joined, "src/a.ts") {
		t.Fatalf("the project's own file was not scanned: %v", joined)
	}
	if strings.Contains(joined, "worktrees") || strings.Contains(joined, "vendor-clone") {
		t.Errorf("a nested checkout was scanned as project material: %v", joined)
	}
}

// THE PROGRESS FILE DOES NOT ENTER THE MAP — the whole point of the separation.
//
// A file that exists to CHANGE cannot be confronted by the gates that demand a justification
// for change: the plan-change gate would accuse it at every new `[x]`, and the only way out
// would be to write a revision saying "the work moved on".
//
// The test runs the real `Walk`, with the plan and its companion in the same folder and the
// SAME layer matching both — the real situation. A test that only called `IsProgressFile`
// would prove the function answers right, not that the scan uses it.
func TestWalk_progressStaysOutOfTheMap(t *testing.T) {
	t.Run("RPSCR-B04: A progress companion stays out of the map", func(t *testing.T) {})
	dir := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(dir, "plans"), 0o755))
	plan := "<!-- @anchors\n  code: MTUAO\n  layer: plan\n-->\n# Plan\n\n### MTUAO-W01 — phase\n"
	must(t, os.WriteFile(filepath.Join(dir, "plans", "0017-mutacao.md"), []byte(plan), 0o644))
	prog := "# Progress — MTUAO\n\n## MTUAO-W01\n\n- [x] done\n"
	must(t, os.WriteFile(filepath.Join(dir, "plans", "0017-mutacao-progress.md"), []byte(prog), 0o644))

	// The layer matches BOTH files: that is what makes the test honest. If the glob excluded
	// the companion, the test would pass without the `Walk` guard doing anything.
	cfg := &config.Config{
		Layers: map[string]config.Layer{"plan": {Pattern: "plans/*.md", Kind: "plan"}},
	}

	files, err := Walk(dir, cfg)
	if err != nil {
		t.Fatal(err)
	}

	var foundPlan, foundProgress bool
	for _, f := range files {
		switch f.Path {
		case "plans/0017-mutacao.md":
			foundPlan = true
		case "plans/0017-mutacao-progress.md":
			foundProgress = true
		}
	}
	if !foundPlan {
		t.Errorf("the PLAN has to enter the map — without it the test proves nothing about the "+
			"companion; files seen: %v", paths(files))
	}
	if foundProgress {
		t.Errorf("the progress file ENTERED the map: the gates would demand a justification "+
			"at every `[x]`; files seen: %v", paths(files))
	}
}

func paths(fs []File) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, f.Path)
	}
	return out
}

// A vendored workflow enters the map with no scenario codes of its own: the examples in its
// comments are the Anchors project's vocabulary.
func TestWalk_upstreamWorkflowCarriesNoCodes(t *testing.T) {
	t.Run("RPSCR-B05: An upstream workflow carries no codes", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".github/workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".github/workflows/anchors-board.yml"), []byte(boardWorkflow), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{Layers: map[string]config.Layer{
		"workflow": {Pattern: ".github/workflows/*.yml", Kind: "code"},
	}}
	files, err := Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected the workflow in the scan, got %v", files)
	}
	f := files[0]
	if !f.Upstream || len(f.Codes) != 0 || f.Parent != "" {
		t.Errorf("expected upstream, no codes, no parent; got upstream=%v codes=%v parent=%q", f.Upstream, f.Codes, f.Parent)
	}
}

// TestRevIgnoresLineEndings: the `rev` is the identity of the CONTENT — it decides whether an
// ingested signal still holds. With `core.autocrlf=true`, the same commit has different
// bytes on Windows and macOS; if the line ending entered the hash, the rev of every file
// would diverge between the two machines and a `map build` on one side would discard the
// signals accumulated on the other. MEASURED on 24/08 in the reference app: one `map build`
// on Windows erased the 1061 signals of the map generated on macOS, without a warning.
func TestRevIgnoresLineEndings(t *testing.T) {
	t.Run("RPSCR-B06: The revision ignores line endings but not content", func(t *testing.T) {})
	lf := []byte("line one\nline two\n")
	crlf := []byte("line one\r\nline two\r\n")
	if shortHash(lf) != shortHash(crlf) {
		t.Fatalf("the same content with different line endings gave different revs: %s vs %s",
			shortHash(lf), shortHash(crlf))
	}
	// And REALLY different content keeps a different rev — the normalization cannot become a
	// collision.
	if shortHash(lf) == shortHash([]byte("line one\nline three\n")) {
		t.Fatal("distinct contents cannot share a rev")
	}
}

func TestWalk_unwalkableRootIsAnError(t *testing.T) {
	t.Run("RPSCR-E01: A root that cannot be walked returns the error", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{"code": {Pattern: "**/*.ts", Kind: "code"}}}
	if _, err := Walk(filepath.Join(t.TempDir(), "missing"), cfg); err == nil {
		t.Error("walking a root that does not exist must return the error")
	}
}

// ---- Classification ----

// TestLayerTieBreakIsStableAndDeclarable: `cfg.Layers` is a map, whose order Go does not
// define. Without a total tie-break, two patterns of the same length would draw the layer at
// every run — and a classification that changes between two runs poisons every gate.
func TestLayerTieBreakIsStableAndDeclarable(t *testing.T) {
	t.Run("RPSCR-B07: Priority, then pattern length, then layer name decide the layer", func(t *testing.T) {})
	t.Run("RPSCR-I01: The classification is stable across runs", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"broad":    {Pattern: "pkg/**/*.ts", Kind: "code"},
		"specific": {Pattern: "pkg/models/**/*.ts", Kind: "spec"},
	}}
	const target = "pkg/models/a.ts"

	for i := 0; i < 50; i++ {
		if l, _ := classify(target, cfg); l != "specific" {
			t.Fatalf("unstable heuristic: run %d gave %q", i, l)
		}
	}

	// Where the heuristic errs, the project DECLARES — and the declaration wins.
	cfg.Layers["broad"] = config.Layer{Pattern: "pkg/**/*.ts", Kind: "code", Priority: 10}
	for i := 0; i < 50; i++ {
		if l, _ := classify(target, cfg); l != "broad" {
			t.Fatalf("a declared `priority` must beat the heuristic; run %d gave %q", i, l)
		}
	}

	// A total tie falls to the layer name.
	tie := &config.Config{Layers: map[string]config.Layer{
		"b": {Pattern: "src/*.ts", Kind: "code"},
		"a": {Pattern: "src/*.ts", Kind: "test"},
	}}
	for i := 0; i < 50; i++ {
		if l, k := classify("src/x.ts", tie); l != "a" || k != "test" {
			t.Fatalf("a total tie goes to the smaller name; run %d gave (%q, %q)", i, l, k)
		}
	}
}

func TestClassify_exclusionRemovesAPath(t *testing.T) {
	t.Run("RPSCR-B08: An exclusion removes a path from its layer", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"c": {Pattern: "src/**", Kind: "code", Exclude: []string{"src/gen/**"}},
	}}
	if l, _ := classify("src/gen/y.ts", cfg); l != "" {
		t.Errorf("an excluded path has no layer, got %q", l)
	}
	if l, _ := classify("src/z.ts", cfg); l != "c" {
		t.Errorf("a path outside the exclusion stays in c, got %q", l)
	}
}

// TestClassifyNormalizesTheWindowsSeparator: `classify` receives paths from more than ten
// callers, and on Windows `filepath.Rel` returns `\`. Without normalizing, the glob match
// fails for EVERY pattern with a directory, and the effect is not a wrong layer — it is NO
// layer: the file vanishes from the map and no gate confronts it again. MEASURED on 24/08 in
// the reference app: `map build` on Windows gave 1857 nodes and 5416 edges where the same
// commit gave 2799/11302, without an error line.
func TestClassifyNormalizesTheWindowsSeparator(t *testing.T) {
	t.Run("RPSCR-B09: A Windows path is classified like its slash form", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"code": {Pattern: "apps/mobile/src/**/*.ts", Kind: "code"},
	}}
	for _, rel := range []string{
		"apps/mobile/src/business-logic/a.ts",
		`apps\mobile\src\business-logic\a.ts`,
	} {
		if l, k := classify(rel, cfg); l != "code" || k != "code" {
			t.Fatalf("classify(%q) = (%q, %q); both forms of the same path have to match", rel, l, k)
		}
	}
}

func TestAmbiguities_reportOnlyHeuristicDecisions(t *testing.T) {
	t.Run("RPSCR-B10: A heuristic decision is reported as an ambiguity, a declared priority is not", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"b": {Pattern: "src/*.ts", Kind: "code"},
		"a": {Pattern: "src/*.ts", Kind: "test"},
		"c": {Pattern: "src/**", Kind: "x"},
	}}
	got := Ambiguities([]File{{Path: "src/x.ts"}, {Path: "src/z.go"}}, cfg)
	if len(got) != 1 || got[0].Arquivo != "src/x.ts" || got[0].Vencedora != "a" ||
		strings.Join(got[0].Perdedoras, ",") != "b,c" {
		t.Errorf("want src/x.ts won by a over b and c, got %+v", got)
	}
	declared := &config.Config{Layers: map[string]config.Layer{
		"b": {Pattern: "src/*.ts", Kind: "code", Priority: 1},
		"a": {Pattern: "src/*.ts", Kind: "test"},
	}}
	if got := Ambiguities([]File{{Path: "src/x.ts"}}, declared); len(got) != 0 {
		t.Errorf("a winner by declared priority is not an ambiguity, got %+v", got)
	}
}

// The UNIT's layer is not the file's, and the spec is the case that separates the two.
//
// A spec matches two patterns: its own (`**/*.spec.md`, layer `spec`) and the one of the
// unit it governs. `ClassifyPath` returns the first — right for the map, where the spec IS a
// node of layer `spec` — and wrong for whoever asks "which layer is this unit".
//
// Measured: `anchors docs duties --unit <x>.spec.md` answered "no mandatory documentation"
// for a lambda that owes the OpenAPI. With the unit's `.ts` the answer was right, and the
// card points at the SPEC — the path the agent uses.
func TestLayerOfUnit(t *testing.T) {
	t.Run("RPSCR-B11: The unit's layer comes from the header before the path", func(t *testing.T) {})
	root := t.TempDir()
	cfg := &config.Config{Layers: map[string]config.Layer{
		"spec":    {Pattern: "**/*.spec.md", Kind: "spec"},
		"lambdas": {Pattern: "pkg/lambdas/**/*.ts", Kind: "code"},
	}}

	write := func(rel, body string) {
		must(t, writeDeep(filepath.Join(root, rel), body))
	}
	write("pkg/lambdas/push/X.spec.md", "<!-- @anchors\n  code: XPTOX\n  layer: lambdas\n-->\n# X\n")
	write("pkg/lambdas/push/X.ts", "export const x = 1;\n")
	write("pkg/lambdas/push/NoHeader.spec.md", "# no header\n")
	write("pkg/lambdas/push/Declared.ts", "// @anchors\n//   layer: handler\n\nexport const y = 1;\n")

	cases := []struct{ name, rel, want string }{
		{"the spec declares the unit's layer in the header", "pkg/lambdas/push/X.spec.md", "lambdas"},
		{"the code already has the unit's layer by its pattern", "pkg/lambdas/push/X.ts", "lambdas"},
		// Without a header it falls to `ClassifyPath` — and there the spec is of layer `spec`.
		// Inventing another answer would be guessing.
		{"a spec without a header falls to the pattern", "pkg/lambdas/push/NoHeader.spec.md", "spec"},
		{"a code file's commented header declares its layer", "pkg/lambdas/push/Declared.ts", "handler"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := LayerOfUnit(root, c.rel, cfg); got != c.want {
				t.Errorf("LayerOfUnit(%s) = %q, want %q", c.rel, got, c.want)
			}
		})
	}

	// And `ClassifyPath` keeps returning the FILE's layer — the map depends on it: the spec
	// is a node of layer `spec`, and changing that would reclassify 84 nodes in a real project.
	if l, _ := ClassifyPath("pkg/lambdas/push/X.spec.md", cfg); l != "spec" {
		t.Errorf("ClassifyPath of the spec = %q, want `spec` — the map depends on it", l)
	}
}

// ---- Codes ----

func TestExtractCodesIgnoresComments(t *testing.T) {
	t.Run("RPSCR-B12: Codes cited in comments are not owned, and each code is listed once", func(t *testing.T) {})
	src := `// see the scenario FOOOX-VR of another screen
const x = 1  // BARRX-S01 here too
export const y = "BAZZX-B01"  // this one is in code (a string), it counts`
	codes := extractCodes([]byte(src))
	has := func(c string) bool {
		for _, x := range codes {
			if x == c {
				return true
			}
		}
		return false
	}
	if has("FOOOX-VR") || has("BARRX-S01") {
		t.Errorf("codes in comments should not count: %v", codes)
	}
	if !has("BAZZX-B01") {
		t.Errorf("a code in code (outside a comment) should count: %v", codes)
	}
	if got := extractCodes([]byte("ABCDX-B01 ABCDX-B01 ABCDX-B02 ABCDX-B01")); strings.Join(got, ",") != "ABCDX-B01,ABCDX-B02" {
		t.Errorf("each code once, in order of first appearance; got %v", got)
	}
}

func TestStripBlockComment(t *testing.T) {
	t.Run("RPSCR-B12: Codes cited in comments are not owned, and each code is listed once", func(t *testing.T) {})
	src := "/* AAAAX-S01 in the block */ const z = \"BBBBX-V01\""
	codes := extractCodes([]byte(src))
	for _, c := range codes {
		if c == "AAAAX-S01" {
			t.Error("a code in a block comment should not count")
		}
	}
}

// TestProjectVocabularyIsSeenByTheScan guards the SILENT failure mode: when the letter class
// was welded into the regex, a project declaring its own (`rule_types`) had the gates seeing
// the code and the scan not. The result was not an error — it was the unit missing from the
// map, with nobody reporting it.
func TestProjectVocabularyIsSeenByTheScan(t *testing.T) {
	t.Run("RPSCR-B13: The project's rule letters are recognised", func(t *testing.T) {})
	defer SetRuleLetters(config.DefaultRuleLetters)

	const text = "scenario KVALX-Z01: the key is immutable"
	if got := extractCodes([]byte(text)); len(got) != 0 {
		t.Fatalf("with the canonical vocabulary, `Z` is not a valid letter — got %v", got)
	}

	SetRuleLetters("SRVAXBNMDEIZ") // the project declared `Z` for Policy
	got := extractCodes([]byte(text))
	if len(got) != 1 || got[0] != "KVALX-Z01" {
		t.Errorf("the scan must see the letter DECLARED by the project; got %v", got)
	}
}

func TestWalk_recordsTheDeclaredIdentityAndFlags(t *testing.T) {
	t.Run("RPSCR-B14: The declared identity and the annotations are recorded", func(t *testing.T) {})
	root := t.TempDir()
	spec := "<!-- @anchors\n  code: OWNRX\n  layer: core\n-->\n# Owner\n\nSee DTAXX-B11 first.\n\n" +
		"@anchors-shared-code\n@noPropagation\n\n### OWNRX-B01 — a rule\n"
	must(t, writeDeep(filepath.Join(root, "a.spec.md"), spec))
	must(t, writeDeep(filepath.Join(root, "b.spec.md"), "# plain\n\n### PLANX-B01 — a rule\n"))
	if HeaderCodeOf(spec) != "OWNRX" || HeaderCodeOf("# plain\n") != "" {
		t.Errorf("HeaderCodeOf reads the declared identity, and nothing when none is declared")
	}
	cfg := &config.Config{Layers: map[string]config.Layer{"spec": {Pattern: "*.spec.md", Kind: "spec"}}}
	files, err := Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	byPath := map[string]File{}
	for _, f := range files {
		byPath[f.Path] = f
	}
	a, b := byPath["a.spec.md"], byPath["b.spec.md"]
	if a.HeaderCode != "OWNRX" || a.Codes[0] != "DTAXX-B11" {
		t.Errorf("the declared identity is OWNRX even though the body cites DTAXX-B11 first; got %q, codes %v", a.HeaderCode, a.Codes)
	}
	if a.HeaderLayer != "core" {
		t.Errorf("the declared unit layer is core, got %q", a.HeaderLayer)
	}
	if !a.SharedCode || !a.NoPropagation {
		t.Errorf("both annotations are recorded: shared=%v noprop=%v", a.SharedCode, a.NoPropagation)
	}
	if b.HeaderCode != "" || b.SharedCode || b.NoPropagation {
		t.Errorf("a file without header or annotations records none: %+v", b)
	}
}

// ---- Header declarations ----

func TestParentDe_onlyInsideTheHeader(t *testing.T) {
	t.Run("RPSCR-B15: The parent is read only inside the header", func(t *testing.T) {})
	if p := parentDe([]byte(boardWorkflow)); p != "" {
		t.Errorf("a `parent:` in a workflow's body is not a header declaration, got %q", p)
	}
	spec := "<!-- @anchors\n  code: INMTN\n  parent: DTBSA-W01\n  layer: lambdas\n-->\n# InstanceMetrics\n\nparent: NOPE-W09\n"
	if p := parentDe([]byte(spec)); p != "DTBSA-W01" {
		t.Errorf("the header's parent should be read, got %q", p)
	}
	code := "// @anchors\n//   ref: PRICX\n//   parent: PRICX-W02\npackage p\n\n// parent: NOPE-W09\n"
	if p := parentDe([]byte(code)); p != "PRICX-W02" {
		t.Errorf("a line-comment header's parent should be read, got %q", p)
	}
	noParent := "<!-- @anchors\n  code: INMTN\n-->\n# InstanceMetrics\n\nparent: NOPE-W09\n"
	if p := parentDe([]byte(noParent)); p != "" {
		t.Errorf("the HTML header ends at `-->`, got %q", p)
	}
	noParentCode := "// @anchors\n//   ref: PRICX\npackage p\n\n// parent: NOPE-W09\n"
	if p := parentDe([]byte(noParentCode)); p != "" {
		t.Errorf("a line-comment header ends at the first non-comment line, got %q", p)
	}
	oneLine := "<!-- @anchors parent: X -->\nparent: NOPE-W09\n"
	if p := parentDe([]byte(oneLine)); p != "" {
		t.Errorf("a one-line header closes on its own line, got %q", p)
	}
}

func TestNeedsFor_pathsForPlansPhasesForSpecs(t *testing.T) {
	t.Run("RPSCR-B16: Needs are plan paths for a plan and phase codes for a spec", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(filepath.Join(root, "plans", "a.md"), "x"))
	plan := []byte("<!-- @anchors\n  code: PLANX\n  needs: `plans/a.md`, plans/b.md\n-->\n")
	if got := needsFor("plan", plan, root, "plans/c.md"); strings.Join(got, ",") != "plans/a.md,plans/b.md" {
		t.Errorf("a plan needs plan paths, got %v", got)
	}
	spec := []byte("<!-- @anchors\n  code: SPECX\n  needs: FNDTN-W02, plans/a.md, FNDTN-B01\n-->\n")
	if got := needsFor("spec", spec, root, "s.spec.md"); strings.Join(got, ",") != "FNDTN-W02" {
		t.Errorf("a spec needs only phase codes, got %v", got)
	}
	if got := needsFor("code", spec, root, "x.go"); got != nil {
		t.Errorf("a code file needs nothing, got %v", got)
	}
}

func TestRevisesDe_onlyPlans(t *testing.T) {
	t.Run("RPSCR-B17: Only a plan revises", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(filepath.Join(root, "plans", "a.md"), "x"))
	hdr := []byte("<!-- @anchors\n  code: PLANX\n  revises: `plans/a.md` -->\n")
	if got := revisesDe("plan", hdr, root, "plans/c.md"); strings.Join(got, ",") != "plans/a.md" {
		t.Errorf("a plan revises plans/a.md, got %v", got)
	}
	if got := revisesDe("spec", hdr, root, "s.spec.md"); got != nil {
		t.Errorf("a spec revises nothing, got %v", got)
	}
}

func TestExtractHeaderDeps(t *testing.T) {
	t.Run("RPSCR-B18: A non-spec file declares dependencies in its header", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(root+"/apps/mobile/src/theme/tokens.ts", "x"))
	must(t, writeDeep(root+"/apps/mobile/src/utils/cn.ts", "x"))
	// a presentation file (no spec) declaring a dep in the header
	content := []byte(`// @anchors
//   layer: presentation
//   dep: theme/tokens.ts, utils/cn.ts
//   updated_at: 2026-08-09
export const x = 1
`)
	rel := "apps/mobile/src/features/audit/presentation/auditFeed.ts"
	deps := extractHeaderDeps(content, root, rel)
	if len(deps) != 2 {
		t.Fatalf("want 2 header deps, got %d: %+v", len(deps), deps)
	}
	if deps[0].File != "apps/mobile/src/theme/tokens.ts" || deps[1].File != "apps/mobile/src/utils/cn.ts" {
		t.Errorf("header deps badly resolved: %+v", deps)
	}
}

func TestDepsForHeaderVsTable(t *testing.T) {
	t.Run("RPSCR-X02: A spec does not read a dep header line", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(root+"/apps/mobile/src/theme/tokens.ts", "x"))
	hdr := []byte("// @anchors\n//   layer: presentation\n//   dep: theme/tokens.ts\n")
	// code (not a spec) → reads the header dep
	if d := depsFor("code", hdr, root, "apps/mobile/src/f/presentation/p.ts"); len(d) != 1 {
		t.Errorf("code should read the dep from the header, got %+v", d)
	}
	// spec → does NOT read the header dep (it uses the table)
	if d := depsFor("spec", hdr, root, "apps/mobile/src/f/x.spec.md"); d != nil {
		t.Errorf("a spec should not read the dep from the header (only the table), got %+v", d)
	}
}

func TestDepsFor_yamlScriptDependsOnItsComposition(t *testing.T) {
	t.Run("RPSCR-B19: A YAML test script depends on the scripts it composes", func(t *testing.T) {})
	d := depsFor("test", []byte("- runFlow: ../u/l.yaml\n"), t.TempDir(), "t/f/x.yaml")
	if len(d) != 1 || d[0].File != "t/u/l.yaml" || d[0].Method != "runFlow" {
		t.Errorf("want one runFlow dependency on t/u/l.yaml, got %+v", d)
	}
}

// ---- Dependency table ----

func TestExtractDepsParsesTable(t *testing.T) {
	t.Run("RPSCR-B20: A spec's dependency table is read in any catalogue language", func(t *testing.T) {})
	t.Run("RPSCR-B22: The method keeps its backticks", func(t *testing.T) {})
	// create the target files in the tempdir (resolveDepPath stats them)
	root := t.TempDir()
	mk := func(rel string) {
		must(t, writeDeep(root+"/"+rel, "x"))
	}
	mk("apps/mobile/src/stores/auth.store.ts")
	mk("apps/mobile/src/hooks/useAuth.ts")

	spec := `# LoginScreen

> ` + "`" + `code: LOGIX` + "`" + `

## Dependências de Dados

| Cód  | Arquivo                  | Método         | Camada |
| ---- | ------------------------ | -------------- | ------ |
| DEP1 | ` + "`stores/auth.store.ts`" + ` | ` + "`useAuthStore`" + ` | store  |
| DEP2 | ` + "`hooks/useAuth.ts`" + `     | ` + "`signIn`" + `       | hook   |

## Data Contract

| Campo | Origem | Obrigatório |
| ----- | ------ | ----------- |
| ` + "`isLoading`" + ` | DEP1 | ✅ |
`
	specRel := "apps/mobile/src/features/auth/screens/LoginScreen.spec.md"
	deps := extractDeps([]byte(spec), root, specRel)
	if len(deps) != 2 {
		t.Fatalf("want 2 deps, got %d: %+v", len(deps), deps)
	}
	// Method KEEPS the author's backticks (the sign of "this is a SYMBOL", not prose); the
	// other columns come without inline markdown.
	if deps[0].Code != "DEP1" || deps[0].Method != "`useAuthStore`" || deps[0].Layer != "store" {
		t.Errorf("DEP1 badly parsed: %+v", deps[0])
	}
	// resolveDepPath must have resolved to the root-relative path (through the spec's /src/)
	if deps[0].File != "apps/mobile/src/stores/auth.store.ts" {
		t.Errorf("DEP1.File did not resolve to the root: %q", deps[0].File)
	}
	if deps[1].Code != "DEP2" || deps[1].File != "apps/mobile/src/hooks/useAuth.ts" {
		t.Errorf("DEP2 badly parsed: %+v", deps[1])
	}

	// The same table under the English heading, columns in another order.
	en := "## Dependencies\n| File | Code | Method | Layer |\n|--|--|--|--|\n| `nope/x.ts` | DEP1 | `m` | l |\n"
	if d := extractDeps([]byte(en), root, "a/b.spec.md"); len(d) != 1 || d[0].Code != "DEP1" || d[0].File != "nope/x.ts" {
		t.Errorf("the English heading and reordered columns should give one dependency, got %+v", d)
	}
}

func TestExtractDepsNoneWithoutSection(t *testing.T) {
	t.Run("RPSCR-B20: A spec's dependency table is read in any catalogue language", func(t *testing.T) {})
	spec := "# Something\n\n## Data Contract\n\n| Field | Origin |\n|--|--|\n| x | y |\n"
	if d := extractDeps([]byte(spec), t.TempDir(), "a/src/b.spec.md"); d != nil {
		t.Errorf("without a dependencies section it should be nil, got %+v", d)
	}
}

func TestExtractDepsSkipsMalformedCode(t *testing.T) {
	t.Run("RPSCR-B21: A row with a malformed code is not a dependency", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(root+"/src/x.ts", "x"))
	spec := "# S\n## Dependências\n| Cód | Arquivo | Método |\n|--|--|--|\n| notdep | `src/x.ts` | m |\n| DEP1 | `src/x.ts` | n |\n| DEP2 |  | o |\n"
	deps := extractDeps([]byte(spec), root, "src/s.spec.md")
	if len(deps) != 1 || deps[0].Code != "DEP1" {
		t.Errorf("a row with a malformed code or no file should be skipped: %+v", deps)
	}
}

func TestResolveDepPathAcceptsSrcPrefixed(t *testing.T) {
	t.Run("RPSCR-B23: A declared file resolves through the src fallbacks", func(t *testing.T) {})
	// the author writes "src/hooks/x.ts" (the project's @/src/… convention); it must resolve
	// to the root-relative path without doubling src/ ("apps/mobile/src/src/…").
	root := t.TempDir()
	must(t, writeDeep(root+"/apps/mobile/src/hooks/x.ts", "x"))
	specRel := "apps/mobile/src/features/f/screens/S.spec.md"
	if got := resolveDepPath(root, specRel, "src/hooks/x.ts"); got != "apps/mobile/src/hooks/x.ts" {
		t.Errorf("a decl with src/ should resolve to the root, got %q", got)
	}
	// and the variant WITHOUT src/ keeps working
	if got := resolveDepPath(root, specRel, "hooks/x.ts"); got != "apps/mobile/src/hooks/x.ts" {
		t.Errorf("a decl without src/ should resolve to the root, got %q", got)
	}
	// nothing exists: the path is kept as written, for the map's dead-edge check to report
	if got := resolveDepPath(root, specRel, "nope/x.ts"); got != "nope/x.ts" {
		t.Errorf("an unresolved decl is kept as written, got %q", got)
	}
}

// ---- Rule tags ----

// The `@realizes` tag holds in ALL THREE forms of a catalogued rule — heading, table row
// and bold bullet. A table column would exist only in the middle one, and referencing
// doctrine would force the spec to change format.
func TestExtractRealizes_theThreeRuleForms(t *testing.T) {
	t.Run("RPSCR-B24: A realizes tag pairs with its rule in the three rule forms", func(t *testing.T) {})
	c := "### CREDT-V01 — limit respected    @realizes LIMIT-R03\n" +
		"| `CREDT-V02` | something else | @realizes `LIMIT-R04` |\n" +
		"- **CREDT-B03** third form\n" +
		"  @realizes LIMIT-R05\n"
	got := extractRealizes("spec", c)
	want := []Realizes{
		{From: "CREDT-V01", To: "LIMIT-R03"},
		{From: "CREDT-V02", To: "LIMIT-R04"},
		{From: "CREDT-B03", To: "LIMIT-R05"},
	}
	if len(got) != len(want) {
		t.Fatalf("want %d declarations, got %d: %+v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] got %+v, want %+v", i, got[i], want[i])
		}
	}
}

// A BLANK LINE closes the rule's scope. Without it, an `@realizes` written in prose was
// attributed to the last rule seen — a FALSE edge pointing at the wrong rule, which is
// worse than capturing nothing: the gate would confirm a realization nobody declared.
func TestExtractRealizes_orphanTagDoesNotStealPreviousRule(t *testing.T) {
	t.Run("RPSCR-B25: A tag after a blank line has no owning rule", func(t *testing.T) {})
	c := "### CREDT-B03 — a rule\n\nloose text with @realizes ORFAN-R01\n"
	got := extractRealizes("spec", c)
	if len(got) != 1 {
		t.Fatalf("want 1 declaration, got %+v", got)
	}
	if got[0].From != "" {
		t.Errorf("the orphan tag was attributed to %q — it should have no owner", got[0].From)
	}
}

// Only the SPEC declares realization: code, test and feature catalogue no rules, and
// reading the tag in them would create an edge from something that owns no rule.
func TestExtractRealizes_specOnly(t *testing.T) {
	t.Run("RPSCR-B26: Only a spec declares rule tags", func(t *testing.T) {})
	c := "### CREDT-V01 — x    @realizes LIMIT-R03\n"
	for _, kind := range []string{"code", "test", "feature", "plan", "product"} {
		if got := extractRealizes(kind, c); got != nil {
			t.Errorf("kind %q: want nil, got %+v", kind, got)
		}
	}
}

// The SAME pair repeated says nothing new; one rule realizing SEVERAL, and several
// realizing the same one, are the two legitimate cases that 1-to-many exists to allow.
func TestExtractRealizes_dedupesThePairNotTheCode(t *testing.T) {
	t.Run("RPSCR-B27: A repeated pair is recorded once", func(t *testing.T) {})
	c := "### CREDT-V01 — x    @realizes LIMIT-R03\n" +
		"\n### CREDT-V01 — x again    @realizes LIMIT-R03\n" +
		"\n### CREDT-V02 — y    @realizes LIMIT-R03\n" +
		"\n### CREDT-V03 — z    @realizes LIMIT-R03 @realizes LIMIT-R09\n"
	got := extractRealizes("spec", c)
	if len(got) != 4 {
		t.Fatalf("want 4 (the repeated pair goes, the others stay), got %d: %+v", len(got), got)
	}
}

func TestExtractGatedBy_onlyFlagScenarios(t *testing.T) {
	t.Run("RPSCR-B28: A gated-by tag names only a flag scenario", func(t *testing.T) {})
	c := "### CREDT-B01 — x @gated-by FLAGX-G01 @gated-by FLAGX-B02\n"
	got := extractGatedBy("spec", c)
	if len(got) != 1 || got[0] != (Realizes{From: "CREDT-B01", To: "FLAGX-G01"}) {
		t.Errorf("only the G code is a gate, owned by CREDT-B01; got %+v", got)
	}
	if got := extractGatedBy("code", c); got != nil {
		t.Errorf("a code file declares no gate, got %+v", got)
	}
}

// ---- Seeds ----

// A plan writes `*.spec.md` to talk about a SET ("every `*.spec.md` needs a header"), not
// to cite a file that will be born. Taking it as a seed made the plan look forever
// unfulfilled — the "file" would never exist — and the queue's cold start seeded phases
// already done.
func TestSeedIgnoresGlobInProse(t *testing.T) {
	t.Run("RPSCR-B29: A plan seeds only concrete spec and doctrine paths", func(t *testing.T) {})
	plan := "Every `*.spec.md` needs a header.\n" +
		"The spec of `apps/x/Tela.spec.md` is born in this phase.\n" +
		"Do not copy the `_TEMPLATE_SCREEN.spec.md`.\n" +
		"The rule lives in `a/b.doctrine.md`, and again `apps/x/Tela.spec.md`.\n"

	got := extractSeeds("plan", plan)
	if strings.Join(got, ",") != "apps/x/Tela.spec.md,a/b.doctrine.md" {
		t.Errorf("only the concrete paths are seeds, once each; got %v", got)
	}
	if got := extractSeeds("spec", plan); got != nil {
		t.Errorf("only a plan seeds, got %v", got)
	}
}

// A NAME WITHOUT A DIRECTORY is not a path — it is the file cited in prose.
//
// A revision writes "the `MutualTls.spec.md` moved to PLTFR-W03" when explaining what
// changed, and that is not the promise to create a file: the plan seeds
// `packages/infra/MutualTls.spec.md`, with the whole path.
//
// Measured in the reference app: three such mentions (all in revisions) made `anchors next` say
// "1 of 12 spec(s) of this plan do not exist yet" in a plan with all 9 delivered — and
// seed work to create files whose directory-less name points nowhere. The plan looked
// forever unfulfilled, and the queue handed out impossible work.
func TestSeedIgnoresNameWithoutDirectory(t *testing.T) {
	t.Run("RPSCR-B29: A plan seeds only concrete spec and doctrine paths", func(t *testing.T) {})
	plan := "> **PLTFR-R0004:** the `MutualTls.spec.md` moved to `PLTFR-W03`, and the\n" +
		"> `CertificatePinning.spec.md` stays in 0016.\n\n" +
		"- [ ] `packages/infra/MutualTls.spec.md` — the mTLS channel\n"

	got := extractSeeds("plan", plan)
	if len(got) != 1 || got[0] != "packages/infra/MutualTls.spec.md" {
		t.Errorf("only the path counts as a seed; got %v", got)
	}
}

// ---- Constraints ----

func TestExtractDepsIgnoresDocumentalTable(t *testing.T) {
	t.Run("RPSCR-X01: A guide's dependencies table is not a dependency", func(t *testing.T) {})
	// a documental "Dependencies" table (File|Description, WITHOUT a code column) — like a
	// guide's. It is not a reuse dependency table and must not yield deps.
	root := t.TempDir()
	must(t, writeDeep(root+"/src/x.tsx", "x"))
	doc := "# Guide\n## Dependências\n| Arquivo | Descrição |\n|--|--|\n| `src/x.tsx` | register routes |\n"
	if d := extractDeps([]byte(doc), root, "guides/G.md"); d != nil {
		t.Errorf("a documental table (no code column) should not become deps: %+v", d)
	}
}

func TestDepsForOnlySpec(t *testing.T) {
	t.Run("RPSCR-X01: A guide's dependencies table is not a dependency", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(root+"/src/x.ts", "x"))
	tbl := "# T\n## Dependências\n| Cód | Arquivo | Método |\n|--|--|--|\n| DEP1 | `src/x.ts` | m |\n"
	if d := depsFor("guide", []byte(tbl), root, "guides/G.md"); d != nil {
		t.Errorf("kind guide should not extract deps: %+v", d)
	}
	if d := depsFor("spec", []byte(tbl), root, "a/src/s.spec.md"); len(d) != 1 {
		t.Errorf("kind spec should extract 1 dep, got %+v", d)
	}
}

// writeDeep creates a file, creating the directories of its path.
func writeDeep(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// Header keys read over the whole file took prose and code as declarations: a body line
// `code: BOGUS` became the unit's identity, `layer: bogus` its layer, and a `// dep:` comment
// in the code a dependency edge — the same hole `parent:` had (RPSCR-B15).
func TestHeaderKeysOnlyInsideHeader(t *testing.T) {
	t.Run("RPSCR-B30: Header keys are read only inside the header", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(filepath.Join(root, "plans", "a.md"), "x"))
	must(t, writeDeep(filepath.Join(root, "a.ts"), "x"))
	body := "<!-- @anchors\n  updated_at: 2026-09-26\n-->\n# T\n\ncode: BOGUS\nlayer: bogus\nneeds: FNDTN-W02, plans/a.md\nrevises: plans/a.md\n"
	if got := extractHeaderCode(body); got != "" {
		t.Errorf("a body `code:` line is not the identity, got %q", got)
	}
	if got := extractHeaderLayer(body); got != "" {
		t.Errorf("a body `layer:` line is not the unit layer, got %q", got)
	}
	if got := needsFor("spec", []byte(body), root, "s.spec.md"); got != nil {
		t.Errorf("a body `needs:` line is not a spec's needs, got %v", got)
	}
	if got := needsFor("plan", []byte(body), root, "plans/c.md"); got != nil {
		t.Errorf("a body `needs:` line is not a plan's needs, got %v", got)
	}
	if got := revisesDe("plan", []byte(body), root, "plans/c.md"); got != nil {
		t.Errorf("a body `revises:` line is not a revision, got %v", got)
	}
	code := []byte("// @anchors\n//   updated_at: 2026-09-26\npackage p\n\n// dep: a.ts\n")
	if got := extractHeaderDeps(code, root, "p.go"); got != nil {
		t.Errorf("a `dep:` comment in the body is not a dependency, got %+v", got)
	}
	must(t, writeDeep(filepath.Join(root, "x", "B.spec.md"), body))
	cfg := &config.Config{Layers: map[string]config.Layer{"spec": {Pattern: "**/*.spec.md", Kind: "spec"}}}
	if got := LayerOfUnit(root, "x/B.spec.md", cfg); got != "spec" {
		t.Errorf("LayerOfUnit reads the header only, got %q", got)
	}
}

// A file that matches a layer and cannot be read was dropped with no error: the map lost a
// unit and no gate ever saw it — the silence this unit is built against.
func TestWalk_unreadableFileIsAnError(t *testing.T) {
	t.Run("RPSCR-E02: A layer file that cannot be read fails the walk", func(t *testing.T) {})
	if os.Geteuid() == 0 {
		t.Skip("root reads a 0o000 file")
	}
	root := t.TempDir()
	p := filepath.Join(root, "src", "a.ts")
	must(t, writeDeep(p, "x"))
	must(t, os.Chmod(p, 0))
	t.Cleanup(func() { _ = os.Chmod(p, 0o644) })
	cfg := &config.Config{Layers: map[string]config.Layer{"code": {Pattern: "src/*.ts", Kind: "code"}}}
	files, err := Walk(root, cfg)
	if err == nil || !strings.Contains(err.Error(), "src/a.ts") {
		t.Errorf("want an error naming src/a.ts, got %v (files %v)", err, paths(files))
	}
}

// The tags follow the declared code length: a project of 7-character codes gets its
// @realizes and @gated-by edges like any other.
func TestExtractRuleTags_followTheDeclaredCodeLength(t *testing.T) {
	t.Run("RPSCR-B31: Rule tags follow the declared code length", func(t *testing.T) {})
	prev := config.CodeLengths
	config.SetCodeLengths([]int{7})
	defer config.SetCodeLengths(prev)
	c := "### CREDITS-B01 — a rule    @realizes LIMITED-R03 @gated-by FLAGSET-G01\n"
	if got := extractRealizes("spec", c); len(got) != 1 || got[0] != (Realizes{From: "CREDITS-B01", To: "LIMITED-R03"}) {
		t.Errorf("a 7-character @realizes must be read: %+v", got)
	}
	if got := extractGatedBy("spec", c); len(got) != 1 || got[0] != (Realizes{From: "CREDITS-B01", To: "FLAGSET-G01"}) {
		t.Errorf("a 7-character @gated-by must be read: %+v", got)
	}
}

func TestWalk_marksSupportFiles(t *testing.T) {
	t.Run("RPSCR-B32: A file in its layer's support list is marked as support", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(filepath.Join(root, "flows", "screens", "login.yaml"), "x"))
	must(t, writeDeep(filepath.Join(root, "flows", "utils", "loginClean.yaml"), "x"))
	cfg := &config.Config{Layers: map[string]config.Layer{
		"e2e":   {Pattern: "flows/**/*.yaml", Kind: "test", Support: []string{"flows/utils/**"}},
		"other": {Pattern: "docs/**", Kind: "doc", Support: []string{"**"}},
	}}
	files, err := Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, f := range files {
		got[f.Path] = f.Support
	}
	if !got["flows/utils/loginClean.yaml"] || got["flows/screens/login.yaml"] || len(got) != 2 {
		t.Errorf("only the utils file is support, got %v", got)
	}
}

func TestScanPaths_readsOnlyTheGivenFilesAsTheWalk(t *testing.T) {
	t.Run("RPSCR-B33: Only the given files are read, as the walk reads them", func(t *testing.T) {})
	root := t.TempDir()
	must(t, writeDeep(filepath.Join(root, "src", "a.spec.md"), "<!-- @anchors\n  code: AAAAX\n-->\n# A\n\n### AAAAX-B01 — a rule\n"))
	must(t, writeDeep(filepath.Join(root, "src", "b.spec.md"), "# B\n"))
	must(t, writeDeep(filepath.Join(root, "README.md"), "# readme\n"))
	must(t, writeDeep(filepath.Join(root, "node_modules", "x", "c.spec.md"), "# vendored\n"))
	cfg := &config.Config{Layers: map[string]config.Layer{"spec": {Pattern: "**/*.spec.md", Kind: "spec"}}}
	walked, err := Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ScanPaths(root, cfg, []string{"src/a.spec.md", "README.md", "node_modules/x/c.spec.md", "src/gone.spec.md"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !reflect.DeepEqual(got[0], walked[0]) {
		t.Fatalf("only the spec, as the walk reads it: got %+v, walk %+v", got, walked)
	}
	paths, err := GovernedPaths(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, f := range walked {
		want = append(want, f.Path)
	}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("the listing names the walk's paths: %v vs %v", paths, want)
	}
}

func TestWalkStaged_readsTheIndex(t *testing.T) {
	t.Run("RPSCR-B34: The staged walk reads the index, not the tree", func(t *testing.T) {})
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", root, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...).CombinedOutput(); err != nil {
			t.Skipf("git %v: %v %s", args, err, out)
		}
	}
	git("init", "-q")
	must(t, writeDeep(filepath.Join(root, "a.spec.md"), "# A v1\n"))
	git("add", ".")
	git("commit", "-qm", "base")
	must(t, writeDeep(filepath.Join(root, "a.spec.md"), "# A v2, not staged\n"))
	must(t, writeDeep(filepath.Join(root, "b.spec.md"), "# B\n"))
	git("add", "b.spec.md")
	must(t, writeDeep(filepath.Join(root, "c.spec.md"), "# C, untracked\n"))
	cfg := &config.Config{Layers: map[string]config.Layer{"spec": {Pattern: "*.spec.md", Kind: "spec"}}}
	files, err := WalkStaged(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	revs := map[string]string{}
	for _, f := range files {
		revs[f.Path] = f.Rev
	}
	if revs["a.spec.md"] != ShortHash([]byte("# A v1\n")) || revs["b.spec.md"] == "" || revs["c.spec.md"] != "" || len(revs) != 2 {
		t.Fatalf("the index: a as committed, b staged, no c; got %v", revs)
	}
}
