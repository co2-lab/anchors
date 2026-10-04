// @anchors
//   code: MGTSM
//   ref: MGCMM

package ops

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/migra"
	"github.com/co2-lab/anchors/internal/scan"
)

const oldFormatMap = "version: 1\ngerado_por: anchors\nnodes: []\nedges: []\n"

// --dry-run says what would change, key by key, and writes nothing.
func TestMigrateDryRunReportsWithoutWriting(t *testing.T) {
	t.Run("MGCMM-B05: A dry run reports without writing", func(t *testing.T) {})
	root := t.TempDir()
	mapPath := writeFile(t, root, mapx.DefaultPath, oldFormatMap)
	err, out := runCmd(t, newMigrateCmd(), "--root", root, "--dry-run")
	if err != nil {
		t.Fatalf("migrate --dry-run: %v", err)
	}
	want := fmt.Sprintf("anchors.graph.yaml would be migrated: format 1 → %d", mapx.FormatoAtual)
	for _, w := range []string{want, "gerado_por  (1 occurrence(s))", "nothing was written"} {
		if !strings.Contains(out, w) {
			t.Errorf("the dry run does not say %q:\n%s", w, out)
		}
	}
	if b, _ := os.ReadFile(mapPath); string(b) != oldFormatMap {
		t.Errorf("--dry-run wrote the map:\n%s", b)
	}
}

// The migration rewrites the old keys and raises the version; the missing anchors.yaml
// is reported and does not stop the map from migrating. A second run changes nothing.
func TestMigrateRewritesTheMapAndIsIdempotent(t *testing.T) {
	t.Run("MGCMM-B02: A missing file is reported and the other still migrates", func(t *testing.T) {})
	t.Run("MGCMM-B03: A file already current is reported as such", func(t *testing.T) {})
	t.Run("MGCMM-B06: A real migration asks for the commit", func(t *testing.T) {})
	t.Run("MGCMM-I01: A second run changes nothing", func(t *testing.T) {})
	t.Run("MGCMM-X01: The keys renamed are the migration package's", func(t *testing.T) {})
	t.Run("MGCMM-X02: The command reminds of the commit and does not make it", func(t *testing.T) {})
	root := t.TempDir()
	mapPath := writeFile(t, root, mapx.DefaultPath, oldFormatMap)
	err, out := runCmd(t, newMigrateCmd(), "--root", root)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	b, _ := os.ReadFile(mapPath)
	if !strings.Contains(string(b), fmt.Sprintf("version: %d", mapx.FormatoAtual)) ||
		!strings.Contains(string(b), "generated_by: anchors") || strings.Contains(string(b), "gerado_por") {
		t.Errorf("the map was not migrated:\n%s", b)
	}
	if !strings.Contains(out, "· anchors.yaml:") {
		t.Errorf("the missing anchors.yaml is not reported:\n%s", out)
	}
	if !strings.Contains(out, "commit the migration") {
		t.Errorf("the output does not ask for the commit:\n%s", out)
	}

	err, out = runCmd(t, newMigrateCmd(), "--root", root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, fmt.Sprintf("anchors.graph.yaml is already in format %d", mapx.FormatoAtual)) ||
		strings.Contains(out, "commit the migration") {
		t.Errorf("the second run is not a no-op:\n%s", out)
	}
	if b2, _ := os.ReadFile(filepath.Join(root, mapx.DefaultPath)); string(b2) != string(b) {
		t.Error("the second run changed the map")
	}
}

// Both files move to the current format in one run, and each file's renamed keys are listed
// in alphabetical order with their counts, whatever order they appear in the file.
func TestMigrateBringsBothFilesAndListsTheKeysSorted(t *testing.T) {
	t.Run("MGCMM-B01: Both the map and the config reach the current format", func(t *testing.T) {})
	t.Run("MGCMM-B04: The renamed keys are listed in alphabetical order with their counts", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, mapx.DefaultPath,
		"version: 1\njulgamentos: []\ngerado_por: anchors\ncode_declarado: X\nnodes: []\nedges: []\n")
	cfgPath := writeFile(t, root, "anchors.yaml", "version: 1\ntrinca_opcional: true\n")
	err, out := runCmd(t, newMigrateCmd(), "--root", root)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	for _, w := range []string{
		fmt.Sprintf("anchors.graph.yaml migrated: format 1 → %d", mapx.FormatoAtual),
		fmt.Sprintf("anchors.yaml migrated: format 1 → %d", mapx.FormatoAtual),
		"trinca_opcional  (1 occurrence(s))",
	} {
		if !strings.Contains(out, w) {
			t.Errorf("the output does not say %q:\n%s", w, out)
		}
	}
	code, gerado, julg := strings.Index(out, "code_declarado  ("), strings.Index(out, "gerado_por  ("),
		strings.Index(out, "julgamentos  (")
	if code < 0 || gerado < 0 || julg < 0 || code > gerado || gerado > julg {
		t.Errorf("the map's keys are not listed in alphabetical order:\n%s", out)
	}
	if b, _ := os.ReadFile(cfgPath); !strings.Contains(string(b), "optional_unit_edges: true") ||
		!strings.Contains(string(b), fmt.Sprintf("version: %d", mapx.FormatoAtual)) {
		t.Errorf("anchors.yaml was not migrated:\n%s", b)
	}
}

// Format 5 renamed the letters of plans, flows and actions. They live in the project's own
// files, so the migration rewrites them there — only for the units of those kinds.
func TestMigrateRewritesTheCodeLettersOfPlansFlowsAndActions(t *testing.T) {
	t.Run("MGCMM-B07: Crossing format 5 rewrites the letters of plans, flows and actions", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, "version: 4\nlayers:\n  plan:\n    pattern: \"plans/*.md\"\n    kind: plan\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n")
	writeFile(t, root, mapx.DefaultPath, "version: 4\nnodes: []\nedges: []\n")
	plan := "<!-- @anchors\n  code: PLANA\n-->\n# Plan\n\n### PLANA-F01 — tree\n\n### PLANA-F02 — rules (depende de PLANA-F01)\n"
	flow := "<!-- @anchors\n  code: FLOWA\n-->\n# Flow\n\n### FLOWA-P01 — confront\n\nFits: `ACTNA`\n\n- `ACTNA-R01` PROMOTABLE → `FLOWA-P02`\n\n### FLOWA-P02 — end\n\n> @terminal\n"
	action := "<!-- @anchors\n  code: ACTNA\n-->\n# Action\n\n### ACTNA-R01 — PROMOTABLE: ok\n"
	spec := "<!-- @anchors\n  code: LOGIN\n  needs: PLANA-F02\n-->\n# Login\n\n### LOGIN-R01 — only anonymous\n\nRevision LOGIN-R0001.\n"
	writeFile(t, root, "plans/0001-foundation.md", plan)
	writeFile(t, root, "flows/work.flow.md", flow)
	writeFile(t, root, "flows/actions/check.action.md", action)
	writeFile(t, root, "src/login.spec.md", spec)
	read := func(rel string) string { b, _ := os.ReadFile(filepath.Join(root, rel)); return string(b) }

	err, out := runCmd(t, newMigrateCmd(), "--root", root, "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "src/login.spec.md would be rewritten") || !strings.Contains(out, "PLANA-F02 → PLANA-W02  (1)") ||
		read("plans/0001-foundation.md") != plan {
		t.Fatalf("the dry run lists the rewrites and writes nothing:\n%s", out)
	}

	if err, out = runCmd(t, newMigrateCmd(), "--root", root); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	for rel, want := range map[string]string{
		"plans/0001-foundation.md":      "### PLANA-W02 — rules (depende de PLANA-W01)",
		"flows/work.flow.md":            "- `ACTNA-O01` PROMOTABLE → `FLOWA-T02`",
		"flows/actions/check.action.md": "### ACTNA-O01 — PROMOTABLE",
		"src/login.spec.md":             "needs: PLANA-W02",
	} {
		if !strings.Contains(read(rel), want) {
			t.Errorf("%s should read %q:\n%s", rel, want, read(rel))
		}
	}
	if !strings.Contains(out, "anchors doctor --fix") {
		t.Errorf("the migration points to the doctor for the installed pipelines:\n%s", out)
	}
	if s := read("src/login.spec.md"); !strings.Contains(s, "### LOGIN-R01 — only anonymous") || !strings.Contains(s, "LOGIN-R0001") {
		t.Errorf("the spec's own permission and revision stay:\n%s", s)
	}

	err, out = runCmd(t, newMigrateCmd(), "--root", root)
	if err != nil || strings.Contains(out, "rewritten") {
		t.Errorf("a second run rewrites nothing: %v\n%s", err, out)
	}
}

// Crossing format 7, a four-character code is widened to five wherever it is written — in
// the files, in their names, in the config and the map —, every governed file without a
// code of its own gets one, the map carries each file's measurement to its new revision,
// and the renames are recorded. A second run changes nothing.
func TestMigrateCrossingSevenGivesEveryFileACode(t *testing.T) {
	t.Run("MGCMM-B08: Crossing format 7 widens the four-character codes and renames the files named by them", func(t *testing.T) {})
	t.Run("MGCMM-B09: Crossing format 7 gives every governed file a code of its own", func(t *testing.T) {})
	t.Run("MGCMM-B10: Crossing format 7 carries what each file was measured at to its new revision", func(t *testing.T) {})
	t.Run("MGCMM-B12: Crossing format 7 turns a file carrying its spec's code into a ref with a code of its own", func(t *testing.T) {})
	t.Run("MGCMM-B11: Crossing format 7 records each renamed code in anchors.renames.yaml", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, "version: 6\ncode_lengths: [4, 5]\nlayers:\n"+
		"  spec:\n    pattern: \"src/*.spec.md\"\n    kind: spec\n"+
		"  logic:\n    pattern: \"src/*.ts\"\n    kind: code\n"+
		"  test:\n    pattern: \"src/*.test.ts\"\n    kind: test\n"+
		"  feature:\n    pattern: \"src/*.feature\"\n    kind: feature\n"+
		"derived:\n  anchor: logic\n  files:\n    spec: [\"{{dir}}/{{name}}.spec.md\"]\n")
	writeFile(t, root, "src/Login.spec.md", "<!-- @anchors\n  code: LOGI\n-->\n# Login\n\n### LOGI-B01 — only anonymous\n")
	t.Run("MGCMM-B13: Crossing format 7 refreshes the stamps a widened code broke", func(t *testing.T) {})
	t.Run("MGCMM-B14: Crossing format 7 dates every file it rewrote", func(t *testing.T) {})
	writeFile(t, root, "src/Login.ts", "// @anchors\n//   ref: LOGI\n//   updated_at: 2026-01-01\n\nexport const login = 1 // LOGI-B01\n")
	writeFile(t, root, "src/Login.test.ts", "// @anchors\n//   ref: LOGI\n\n"+
		"// @contract: src/Login.ts | export const login = 1 // LOGI-B01 | 1 | "+hash8("export const login = 1 // LOGI-B01")+"\n"+
		"// @contract: src/Login.ts | export const login = 1 // LOGI-B01 | 1 | deadbeef\n"+
		"test('LOGI-B01: only anonymous', () => {})\n")
	writeFile(t, root, "baselines/LOGI-B01.txt", "a capture\n")
	writeFile(t, root, "scripts/run.sh", "# runs LOGI-B01 and LOGI-VR-S01; LOGI alone and LOGI_URL stay\n")
	writeFile(t, root, "src/Login.feature", "# @anchors\n#   code: LOGI\n\nFeature: Login\n")
	read := func(rel string) string { b, _ := os.ReadFile(filepath.Join(root, rel)); return string(b) }

	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		t.Fatal(err)
	}
	prev := config.CodeLengths
	config.SetCodeLengths([]int{4, 5})
	files, err := scan.Walk(root, cfg)
	config.SetCodeLengths(prev)
	if err != nil {
		t.Fatal(err)
	}
	var nodes strings.Builder
	for _, f := range files {
		fmt.Fprintf(&nodes, "  - id: %s\n    kind: %s\n    rev: %s\n    mutation_score: 0.9\n", f.Path, f.Kind, f.Rev)
	}
	writeFile(t, root, mapx.DefaultPath, "version: 6\nnodes:\n"+nodes.String()+"edges: []\n")

	err, out := runCmd(t, newMigrateCmd(), "--root", root)
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	t.Log(out)
	spec := read("src/Login.spec.md")
	m := regexp.MustCompile(`code: (LOGI[A-Z0-9])\n`).FindStringSubmatch(spec)
	if m == nil {
		t.Fatalf("the spec's code is widened keeping the old one as prefix:\n%s", spec)
	}
	newCode := m[1]
	if !strings.Contains(spec, "### "+newCode+"-B01") || !strings.Contains(read("src/Login.test.ts"), "'"+newCode+"-B01:") {
		t.Errorf("the rule and the test that cite the code speak the new one:\n%s\n%s", spec, read("src/Login.test.ts"))
	}
	if _, err := os.Stat(filepath.Join(root, "baselines", newCode+"-B01.txt")); err != nil {
		t.Errorf("the file named by the code is renamed with it: %v", err)
	}
	t.Run("MGCMM-B15: Crossing format 7 rewrites the rule codes cited in files the project does not govern", func(t *testing.T) {})
	if got, want := read("scripts/run.sh"), "# runs "+newCode+"-B01 and "+newCode+"-VR-S01; LOGI alone and LOGI_URL stay\n"; got != want {
		t.Errorf("an ungoverned file has its rule codes rewritten and nothing else:\n%s", got)
	}
	if !strings.Contains(read(config.DefaultFile), "code_lengths: [5]") {
		t.Errorf("the config reads five-character codes only:\n%s", read(config.DefaultFile))
	}
	for _, rel := range []string{"src/Login.ts", "src/Login.test.ts"} {
		s := read(rel)
		c := regexp.MustCompile(`(?m)^//   code: ([A-Z0-9]{5})$`).FindStringSubmatch(s)
		if c == nil || c[1] == newCode || !strings.Contains(s, "ref: "+newCode) {
			t.Errorf("%s gets a code of its own and keeps its unit's ref:\n%s", rel, s)
		}
	}
	if f := read("src/Login.feature"); !strings.Contains(f, "#   ref: "+newCode+"\n") ||
		!regexp.MustCompile(`#   code: [A-Z0-9]{5}\n`).MatchString(f) || strings.Contains(f, "code: "+newCode) {
		t.Errorf("the feature that carried the spec's code refs it and gets its own:\n%s", f)
	}
	if tst := read("src/Login.test.ts"); !strings.Contains(tst, "| 1 | "+hash8("export const login = 1 // "+newCode+"-B01")+"\n") ||
		!strings.Contains(tst, "| 1 | deadbeef") {
		t.Errorf("the stamp that held follows the widened code, the stale one stays:\n%s", tst)
	}
	if !strings.Contains(read("src/Login.ts"), "updated_at: "+time.Now().Format("2006-01-02")) {
		t.Errorf("a rewritten file carries the day of the migration:\n%s", read("src/Login.ts"))
	}
	if !strings.Contains(read(migra.RenamesFile), "LOGI: "+newCode) {
		t.Errorf("the rename is recorded:\n%s", read(migra.RenamesFile))
	}
	g, err := mapx.Load(filepath.Join(root, mapx.DefaultPath))
	if err != nil {
		t.Fatal(err)
	}
	config.SetCodeLengths([]int{5})
	defer config.SetCodeLengths(prev)
	now, err := scan.Walk(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range now {
		if n := g.Node(f.Path); n != nil && n.Rev != f.Rev {
			t.Errorf("%s: the map still has the revision before the migration (%s, now %s)", f.Path, n.Rev, f.Rev)
		}
	}

	before := map[string]string{}
	for _, rel := range []string{"src/Login.spec.md", "src/Login.ts", "src/Login.test.ts", "src/Login.feature", migra.RenamesFile} {
		before[rel] = read(rel)
	}
	if err, out = runCmd(t, newMigrateCmd(), "--root", root); err != nil {
		t.Fatal(err)
	}
	for rel, s := range before {
		if read(rel) != s {
			t.Errorf("the second run changed %s", rel)
		}
	}
}

// A file the migration cannot write stops it with the write error; run again after the
// cause is fixed, it codes only what is left.
func TestMigrateCrossingSevenFailsOnAFileItCannotWrite(t *testing.T) {
	t.Run("MGCMM-E02: Crossing format 7, a file that cannot be written fails the command, and a second run finishes it", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, "version: 6\nlayers:\n  logic:\n    pattern: \"src/*.ts\"\n    kind: code\n")
	writeFile(t, root, mapx.DefaultPath, "version: 6\nnodes: []\nedges: []\n")
	locked := writeFile(t, root, "src/a.ts", "// @anchors\n//   layer: logic\n\nexport const a = 1\n")
	if err := os.Chmod(locked, 0o444); err != nil {
		t.Fatal(err)
	}
	if f, err := os.OpenFile(locked, os.O_WRONLY, 0); err == nil {
		f.Close()
		t.Skip("this user writes read-only files")
	}
	if err, out := runCmd(t, newMigrateCmd(), "--root", root); err == nil {
		t.Fatalf("a file that cannot be written must fail the migration:\n%s", out)
	}
	if err := os.Chmod(locked, 0o644); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(root, mapx.DefaultPath)); !strings.Contains(string(b), "version: 6") {
		t.Errorf("the map goes back to its format, so the next run crosses the step again:\n%s", b)
	}
	if err, out := runCmd(t, newMigrateCmd(), "--root", root); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if b, _ := os.ReadFile(locked); !regexp.MustCompile(`//   code: [A-Z0-9]{5}\n`).Match(b) {
		t.Errorf("the second run gives the file its code:\n%s", b)
	}
}

// hash8 is a stamp's hash: the first eight hex digits of the snippet's SHA-256.
func hash8(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])[:8]
}
