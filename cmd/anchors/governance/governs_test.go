package governance

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/spf13/cobra"
)

// captureStdout collects what fn writes to os.Stdout. The governance commands print
// with fmt.Printf, not through cmd.OutOrStdout, so the pipe is the only way to read them.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	fn()
	w.Close()
	os.Stdout = orig
	return <-done
}

// govProject writes a throwaway project: anchors.yaml, the given files, and the map.
func govProject(t *testing.T, yaml string, files map[string]string, g *mapx.Graph) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
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

// runCmd executes a command built by one of the constructors and returns its stdout.
func runCmd(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	cmd.SetArgs(args)
	cmd.SetErr(io.Discard)
	cmd.SilenceUsage = true
	var err error
	out := captureStdout(t, func() { err = cmd.Execute() })
	return out, err
}

// governsFiles puts the map's nodes on disk: the command resolves its argument against
// the root only when the file exists there.
func governsFiles() map[string]string {
	return map[string]string{"GUIDE.md": "# g\n", "STYLE.md": "# s\n", "a.go": "package a\n", "b.go": "package a\n", "a.spec.md": "# a\n"}
}

func governsGraph() *mapx.Graph {
	return &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "GUIDE.md", Kind: mapx.KindGuide},
			{ID: "STYLE.md", Kind: mapx.KindGuide},
			{ID: "a.go", Kind: mapx.KindCode},
			{ID: "b.go", Kind: mapx.KindCode},
			{ID: "a.spec.md", Kind: mapx.KindSpec},
		},
		Edges: []mapx.Edge{
			{From: "GUIDE.md", To: "a.go", Type: mapx.EdgeGoverns},
			{From: "GUIDE.md", To: "b.go", Type: mapx.EdgeGoverns},
			{From: "GUIDE.md", To: "a.spec.md", Type: mapx.EdgeGoverns},
			{From: "STYLE.md", To: "a.go", Type: mapx.EdgeGoverns},
			// a non-governs edge must not be counted as governance
			{From: "a.spec.md", To: "a.go", Type: mapx.EdgeSpecifies},
		},
	}
}

func TestGovernsBoardRanksGuidesByCount(t *testing.T) {
	dir := govProject(t, "version: 2\nlayers: {}\n", governsFiles(), governsGraph())

	out, err := runCmd(t, newGovernsCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "2 guide(s) with governance") {
		t.Errorf("the board should count the two guides:\n%s", out)
	}
	guide := strings.Index(out, "   3  GUIDE.md")
	style := strings.Index(out, "   1  STYLE.md")
	if guide < 0 || style < 0 || guide > style {
		t.Errorf("GUIDE.md (3) should come before STYLE.md (1), each with its count:\n%s", out)
	}
	if !strings.Contains(out, "total pairs (guide, governed) = 4") {
		t.Errorf("the total must add the direct governs edges only (4):\n%s", out)
	}
}

func TestGovernsDetailGroupsByKind(t *testing.T) {
	dir := govProject(t, "version: 2\nlayers: {}\n", governsFiles(), governsGraph())

	out, err := runCmd(t, newGovernsCmd(), "--root", dir, "GUIDE.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "GUIDE.md governs 3 file(s):") {
		t.Errorf("detail should name the guide and its count:\n%s", out)
	}
	code := strings.Index(out, "[code] 2")
	spec := strings.Index(out, "[spec] 1")
	if code < 0 || spec < 0 || code > spec {
		t.Errorf("kinds should be listed sorted, each with its count:\n%s", out)
	}
	for _, id := range []string{"    a.go", "    b.go", "    a.spec.md"} {
		if !strings.Contains(out, id) {
			t.Errorf("governed file %q missing:\n%s", id, out)
		}
	}
}

func TestGovernsNobody(t *testing.T) {
	dir := govProject(t, "version: 2\nlayers: {}\n", governsFiles(), governsGraph())

	out, err := runCmd(t, newGovernsCmd(), "--root", dir, "a.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "a.go governs nobody") {
		t.Errorf("a code file governs nothing:\n%s", out)
	}

	empty := govProject(t, "version: 2\nlayers: {}\n", nil, &mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode}}})
	out, err = runCmd(t, newGovernsCmd(), "--root", empty)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no guide governs anything") {
		t.Errorf("a map without governs edges has an empty board:\n%s", out)
	}
}

func TestGovernsWithoutMapFails(t *testing.T) {
	dir := govProject(t, "version: 2\nlayers: {}\n", nil, nil)
	_, err := runCmd(t, newGovernsCmd(), "--root", dir)
	if err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("a missing map must fail pointing at `anchors map build`; got %v", err)
	}
}
