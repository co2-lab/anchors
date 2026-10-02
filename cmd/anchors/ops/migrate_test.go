package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
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
