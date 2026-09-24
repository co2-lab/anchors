package quality

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/mapx"
)

// `anchors stamp --refresh <module>`, end to end in a git repository: the answer names the
// test, the line and the member, and shows how the stamped block changed from HEAD.
func TestStampRefresh_listsTheDoublesAndTheChange(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git")
	}
	root := t.TempDir()
	mod := "src/hooks/balance.ts"
	test := "src/screens/Home.test.tsx"
	oldMod := "export function useBalance(id) {\n  return fetchBalance(id)\n}\n"
	write := func(p, body string) {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, p)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, p), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(mod, oldMod)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: mod, Kind: mapx.KindCode}, {ID: test, Kind: mapx.KindTest}}}
	// Stamp it with the hash the gate would compute today.
	write(test, "// @contract: "+mod+" | export function useBalance(id) { | 3 | 00000000\njest.mock('@/src/hooks/balance')\n")
	if _, err := gate.RefreshStamps(g, root, mod, true); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-qm", "base"}} {
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}

	write(mod, "export function useBalance(id) {\n  return fetchBalance(id, { cache: false })\n}\n")
	out := captureStdout(t, func() {
		if err := refreshStamps(root, g, []string{mod}, false); err != nil {
			t.Fatal(err)
		}
	})
	for _, want := range []string{
		mod + " changed — 1 double(s)",
		test + ":1",
		"`export function useBalance(id) {`",
		"-   return fetchBalance(id)",
		"+   return fetchBalance(id, { cache: false })",
		"adjust it to the new one",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the answer should show %q:\n%s", want, out)
		}
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = orig
	b, _ := io.ReadAll(r)
	return string(b)
}
