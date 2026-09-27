package mapcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// A change to the spec propagates DOWN to the code and the test it specifies, and is
// governed by nobody; a change to the code propagates nowhere and is confronted UP with
// the spec and the guide.
func TestImpact_bothDirections(t *testing.T) {
	t.Run("MPCTI-B01: A change to a spec propagates down to the code and the test it specifies", func(t *testing.T) {})
	t.Run("MPCTI-B02: A change to the code is validated up against its spec and its guide", func(t *testing.T) {})
	t.Run("MPCTI-B03: An empty direction is said explicitly instead of printed as an empty list", func(t *testing.T) {})
	root := fixtureProject(t)

	spec := runCmd(t, newImpactCmd(), "src/login.spec.md", "--root", root)
	if !strings.Contains(spec, "impact of: src/login.spec.md") {
		t.Errorf("the origin is not named:\n%s", spec)
	}
	down := spec[strings.Index(spec, "↓"):strings.Index(spec, "↑")]
	for _, id := range []string{"src/login.ts", "src/login.test.ts"} {
		if !strings.Contains(down, id) {
			t.Errorf("%s should be redone when the spec changes:\n%s", id, spec)
		}
	}
	if !strings.Contains(spec, "(nothing — it is not governed by anyone)") {
		t.Errorf("the spec is governed by nobody here:\n%s", spec)
	}

	code := runCmd(t, newImpactCmd(), "src/login.ts", "--root", root)
	if !strings.Contains(code, "(nothing — no child depends on it)") {
		t.Errorf("nothing depends on the code:\n%s", code)
	}
	up := code[strings.Index(code, "↑"):]
	for _, id := range []string{"src/login.spec.md", "guides/CODE.md"} {
		if !strings.Contains(up, id) {
			t.Errorf("the code must be confronted with %s:\n%s", id, code)
		}
	}
}

func TestImpact_errors(t *testing.T) {
	t.Run("MPCTI-E01: A file that is not a node of the map is refused", func(t *testing.T) {})
	t.Run("MPCTI-E02: Without a map the impact query fails and asks for the map build", func(t *testing.T) {})
	root := fixtureProject(t)
	if _, err := runCmdErr(newImpactCmd(), t, "src/ghost.ts", "--root", root); err == nil ||
		!strings.Contains(err.Error(), "is not in the map") {
		t.Errorf("a file outside the map: got %v", err)
	}
	if _, err := runCmdErr(newImpactCmd(), t, "src/login.ts", "--root", t.TempDir()); err == nil ||
		!strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: got %v", err)
	}
}

// The query reads the map and never writes it: running it leaves the map byte for byte
// as it was.
func TestImpact_writesNothing(t *testing.T) {
	t.Run("MPCTI-X01: The impact query leaves the map unchanged", func(t *testing.T) {})
	root := fixtureProject(t)
	path := filepath.Join(root, mapx.DefaultPath)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	runCmd(t, newImpactCmd(), "src/login.spec.md", "--root", root)
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("the impact query changed the map")
	}
}

// relTo returns the id of a node, and the id is the same on every machine — the map is
// versioned and travels between macOS, Linux and Windows. Without the normalisation,
// `filepath.Clean`/`Rel` return "packages\backend\x.spec.md" on Windows, the lookup in
// the map (written with "/") does not match, and a file that IS in the map shows up as
// "GOVERNED but outside the map" — asking for a `map build` that would rewrite the whole
// map in the machine's dialect.
func TestRelTo_returnsTheNodeID(t *testing.T) {
	t.Run("MPCTI-B04: A root-relative, native-separator or absolute argument resolves to the same node id", func(t *testing.T) {})
	root := t.TempDir()
	dir := filepath.Join(root, "packages", "backend", "services")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "auth.spec.md"), []byte("# spec\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	const id = "packages/backend/services/auth.spec.md"

	// The three legitimate origins of the argument: root-relative with a forward slash
	// (the form Anchors' prompts print), relative with the native separator (what the
	// machine's shell completes), and absolute.
	cases := map[string]string{
		"relative with slash": id,
		"relative native":     filepath.FromSlash(id),
		"absolute":            filepath.Join(root, "packages", "backend", "services", "auth.spec.md"),
	}
	for name, arg := range cases {
		got := relTo(root, arg)
		if got != id {
			t.Errorf("%s: relTo(%q) = %q; want the id %q", name, arg, got, id)
		}
		if strings.Contains(got, `\`) {
			t.Errorf("%s: a node id cannot carry a backslash: %q", name, got)
		}
	}
}

// A relative argument that does not exist under the root is read from the directory the
// command was called in, as every shell does.
func TestRelTo_resolvesFromTheWorkingDirectory(t *testing.T) {
	t.Run("MPCTI-B05: A relative argument that does not exist under the root is resolved from the working directory", func(t *testing.T) {})
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "packages", "backend"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "packages", "backend", "x.spec.md"), []byte("# spec\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(filepath.Join(root, "packages"))
	if got := relTo(root, "backend/x.spec.md"); got != "packages/backend/x.spec.md" {
		t.Errorf("relTo from a subdirectory = %q, want packages/backend/x.spec.md", got)
	}
}

// Even on the path where the argument names nothing on disk, the returned id is in the
// map's dialect, for the same reason as the happy case.
func TestRelTo_neverReturnsABackslash(t *testing.T) {
	t.Run("MPCTI-I01: The resolved node id always uses forward slashes", func(t *testing.T) {})
	root := t.TempDir()
	got := relTo(root, filepath.FromSlash("does/not/exist/anywhere.spec.md"))
	if strings.Contains(got, `\`) {
		t.Errorf("not even on the error path can the id carry a backslash: %q", got)
	}
}
