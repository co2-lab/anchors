package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const oldFormatMap = "version: 1\ngerado_por: anchors\nnodes: []\nedges: []\n"

// --dry-run says what would change, key by key, and writes nothing.
func TestMigrateDryRunReportsWithoutWriting(t *testing.T) {
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
