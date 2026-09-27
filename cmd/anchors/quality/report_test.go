package quality

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// qProject writes a throwaway project: anchors.yaml (when yaml is not empty), the given
// files, and the map (when g is not nil). Returns the root.
func qProject(t *testing.T, yaml string, files map[string]string, g *mapx.Graph) string {
	t.Helper()
	dir := t.TempDir()
	if yaml != "" {
		if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for name, body := range files {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if g != nil {
		if err := mapx.Save(g, filepath.Join(dir, mapx.DefaultPath)); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// runQ executes a command built by a constructor and returns what it printed to stdout.
func runQ(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	var err error
	out := captureStdout(t, func() { err = cmd.Execute() })
	return out, err
}

func readQ(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
