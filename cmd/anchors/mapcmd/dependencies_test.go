// @anchors
//   code: DPNDT
//   ref: DPNDS

package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func TestDependencies_levels(t *testing.T) {
	t.Run("DPNDS-B01: The direct dependencies show, and --propagate goes further with the path of each item", func(t *testing.T) {})
	root := tiesProject(t)
	direct := runCmd(t, newDependenciesCmd(), "src/app.tsx", "--root", root)
	if !strings.Contains(direct, "dependencies of: src/app.tsx") || !strings.Contains(direct, "↓ Direct (1)") || !strings.Contains(direct, "src/nav.tsx — is depended on") {
		t.Errorf("the app depends on the navigator:\n%s", direct)
	}
	if strings.Contains(direct, "src/screen.tsx") {
		t.Errorf("the screen is a level down:\n%s", direct)
	}
	cascade := runCmd(t, newDependenciesCmd(), "src/app.tsx", "--root", root, "--propagate", "cascade")
	for _, want := range []string{"src/screen.tsx — is depended on · via src/nav.tsx ← src/app.tsx", "src/util.ts — is depended on · via src/screen.tsx ← src/nav.tsx ← src/app.tsx", "src/screen.spec.md — specifies it"} {
		if !strings.Contains(cascade, want) {
			t.Errorf("want %q in\n%s", want, cascade)
		}
	}
	if stopped := runCmd(t, newDependenciesCmd(), "src/app.tsx", "--root", root, "--propagate", "cascade", "--stop-at", "src/nav.tsx"); strings.Contains(stopped, "src/screen.tsx") {
		t.Errorf("the stop goes no further:\n%s", stopped)
	}
	if both := runCmd(t, newDependenciesCmd(), "src/app.tsx", "src/screen.tsx", "--root", root); !strings.Contains(both, "src/util.ts") || !strings.Contains(both, "src/nav.tsx") {
		t.Errorf("two files at once:\n%s", both)
	}
}

func TestDependencies_formats(t *testing.T) {
	t.Run("DPNDS-B02: The files and json formats print the dependencies", func(t *testing.T) {})
	root := tiesProject(t)
	if got := runCmd(t, newDependenciesCmd(), "src/screen.tsx", "--root", root, "--format", "files"); got != "src/screen.spec.md\nsrc/util.ts\n" {
		t.Errorf("one per line: %q", got)
	}
	js := runCmd(t, newDependenciesCmd(), "src/app.tsx", "--root", root, "--format", "json", "--propagate", "indirect")
	for _, want := range []string{`"origins": [`, `"level": 2`, `"path": [`} {
		if !strings.Contains(js, want) {
			t.Errorf("json: want %s in\n%s", want, js)
		}
	}
}

func TestDependencies_writesNothing(t *testing.T) {
	t.Run("DPNDS-X01: The dependencies query leaves the map unchanged", func(t *testing.T) {})
	root := tiesProject(t)
	path := filepath.Join(root, mapx.DefaultPath)
	before, _ := os.ReadFile(path)
	runCmd(t, newDependenciesCmd(), "src/screen.tsx", "--root", root)
	if after, _ := os.ReadFile(path); string(before) != string(after) {
		t.Error("the dependencies query changed the map")
	}
}
