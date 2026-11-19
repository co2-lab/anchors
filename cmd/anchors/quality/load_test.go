// @anchors
//   code: LDQRT
//   ref: LDQRS

package quality

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/load"
)

func TestAdmitRun(t *testing.T) {
	t.Run("LDQRS-B01: A suite's run holds a lease at the level its files' scope costs, released after; with no load: or with the queue off it holds none", func(t *testing.T) {
		t.Setenv("ANCHORS_LOAD_DIR", t.TempDir())
		t.Setenv("ANCHORS_LOAD", "")
		t.Setenv("CI", "")
		t.Setenv("GITHUB_ACTIONS", "")
		root := t.TempDir()
		suite := config.Suite{Layer: "unit", Workspace: "packages/backend"}
		write := func(body string) {
			if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte("version: 8\n"+body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		write("")
		admitRun(root, "test", suite, nil)()
		if st, _ := load.Read(); len(st.Leases) != 0 {
			t.Errorf("no load: holds nothing: %+v", st)
		}
		write("load:\n  tasks:\n    mutation:\n      scopes:\n        \"packages/backend/**\": heavy\n")
		release := admitRun(root, "mutation", suite, nil)
		st, _ := load.Read()
		if len(st.Leases) != 1 || st.Leases[0].Level != "heavy" || st.Leases[0].Task != "mutation unit" {
			t.Errorf("a whole run is placed by its workspace: %+v", st)
		}
		release()
		release = admitRun(root, "mutation", suite, []string{filepath.Join(root, "apps/web/a.ts")})
		if st, _ := load.Read(); len(st.Leases) != 1 || st.Leases[0].Level != "medium" {
			t.Errorf("files outside every scope take the default: %+v", st)
		}
		release()
		skipLoadQueue()
		admitRun(root, "mutation", suite, nil)()
		if st, _ := load.Read(); len(st.Leases) != 0 {
			t.Errorf("the commit hook's queue is off: %+v", st)
		}
	})
}
