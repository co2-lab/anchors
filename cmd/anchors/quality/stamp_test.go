// @anchors
//   ref: CNSTC

package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
	"github.com/co2-lab/anchors/internal/mapx"
)

const stampYAML = `version: 2
layers: {}
derived:
  anchor: spec
  mock_detect: "(?:jest|vi)\\.mock\\(['\"]([^'\"]+)"
  export_detect: "^export\\s+(?:async\\s+)?(?:function|const)\\s+(\\w+)"
`

const stampHooks = `export function useBalance(id) {
  return fetchBalance(id)
}
`

func stampFixture(t *testing.T) (string, string) {
	t.Helper()
	test := "src/Home.test.tsx"
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "src/hooks/balance.ts", Kind: mapx.KindCode},
		{ID: test, Kind: mapx.KindTest},
		// a test with nothing to stamp: silent
		{ID: "src/plain.test.tsx", Kind: mapx.KindTest},
		// a test in the map that is gone from disk: reported, not fatal
		{ID: "src/gone.test.tsx", Kind: mapx.KindTest},
	}}
	dir := qProject(t, stampYAML, map[string]string{
		"src/hooks/balance.ts": stampHooks,
		test:                   "jest.mock('@/src/hooks/balance', () => ({ useBalance: jest.fn() }))\n",
		"src/plain.test.tsx":   "it('adds', () => {})\n",
	}, g)
	return dir, test
}

func TestStampWritesTheMissingStamp(t *testing.T) {
	t.Run("CNSTC-B01: With no test named, every test of the map is considered", func(t *testing.T) {})
	t.Run("CNSTC-B02: The missing stamp is written above the double and totalled", func(t *testing.T) {})
	t.Run("CNSTC-I01: Stamping twice writes nothing the second time", func(t *testing.T) {})
	t.Run("CNSTC-E04: A test of the map that cannot be read is reported and skipped", func(t *testing.T) {})
	dir, test := stampFixture(t)

	out, err := runQ(t, newStampCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"src/gone.test.tsx — could not be read",
		test + "\n  + @/src/hooks/balance → src/hooks/balance.ts | export function useBalance(id) { | 3\n",
		"wrote 1 stamp(s) in 1 file(s); 0 double(s) left unstamped",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "src/plain.test.tsx") {
		t.Errorf("a test with no double is not listed:\n%s", out)
	}
	body := readQ(t, filepath.Join(dir, test))
	if !strings.HasPrefix(body, "// @contract: src/hooks/balance.ts | export function useBalance(id) { | 3 | ") {
		t.Errorf("the stamp should be written above the double:\n%s", body)
	}

	// It only ADDS: a second run finds nothing to write.
	out, err = runQ(t, newStampCmd(), "--root", dir, test)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "wrote 0 stamp(s) in 0 file(s)") {
		t.Errorf("an already stamped double is left alone:\n%s", out)
	}
}

func TestStampDryRunWritesNothing(t *testing.T) {
	t.Run("CNSTC-B03: The dry run says what it would write and writes nothing", func(t *testing.T) {})
	dir, test := stampFixture(t)
	before := readQ(t, filepath.Join(dir, test))

	out, err := runQ(t, newStampCmd(), "--root", dir, "--dry-run", test)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "would write 1 stamp(s) in 1 file(s)") {
		t.Errorf("dry-run says what it would write:\n%s", out)
	}
	if after := readQ(t, filepath.Join(dir, test)); after != before {
		t.Errorf("dry-run must not touch the file:\n%s", after)
	}
}

func TestStampRefreshWithNothingStamped(t *testing.T) {
	t.Run("CNSTC-B05: The refresh of a file no double is stamped against says so", func(t *testing.T) {})
	dir, _ := stampFixture(t)

	out, err := runQ(t, newStampCmd(), "--root", dir, "--refresh", "src/hooks/balance.ts")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "src/hooks/balance.ts — no double is stamped against a previous version of it.") {
		t.Errorf("no stamp points at the file:\n%s", out)
	}
}

func TestStampNeedsConfigAndMap(t *testing.T) {
	t.Run("CNSTC-E01: The stamp command without configuration fails", func(t *testing.T) {})
	t.Run("CNSTC-E02: The stamp command without a map points at the map build", func(t *testing.T) {})
	dir := t.TempDir()
	if _, err := runQ(t, newStampCmd(), "--root", dir); err == nil {
		t.Error("no anchors.yaml: must fail")
	}
	if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte(stampYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runQ(t, newStampCmd(), "--root", dir); err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: must point at map build; got %v", err)
	}
}

func TestBlockDiffListsLeftAndCame(t *testing.T) {
	t.Run("CNSTC-B08: The block diff lists what left and what came, counting repeats", func(t *testing.T) {})
	got := blockDiff("a\nb\nc", "a\nB\nc\nd")
	want := []string{"- b", "+ B", "+ d"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("blockDiff = %v, want %v", got, want)
	}
	if got := blockDiff("x\nx\ny", "x\ny"); strings.Join(got, "|") != "- x" {
		t.Errorf("a line repeated in the old block and kept once is one line that left, got %v", got)
	}
	if got := blockDiff("x\nx\ny", "y\nx\nx"); len(got) != 0 {
		t.Errorf("the same lines, repeats included, are no change, got %v", got)
	}
}

// `anchors stamp --refresh <module>`, end to end in a git repository: the answer names the
// test, the line and the member, and shows how the stamped block changed from HEAD.
func TestStampRefresh_listsTheDoublesAndTheChange(t *testing.T) {
	t.Run("CNSTC-B04: The refresh lists the doubles stamped against the old version with the block change", func(t *testing.T) {})
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

// refreshFixture stamps a double against the current module, outside git, then changes
// the module to newMod: the double is now stamped against the previous version.
func refreshFixture(t *testing.T, newMod string) (root, mod, test string, g *mapx.Graph) {
	t.Helper()
	mod, test = "src/hooks/balance.ts", "src/screens/Home.test.tsx"
	root = qProject(t, "", map[string]string{
		mod:  stampHooks,
		test: "// @contract: " + mod + " | export function useBalance(id) { | 3 | 00000000\njest.mock('@/src/hooks/balance')\n",
	}, nil)
	g = &mapx.Graph{Nodes: []mapx.Node{{ID: mod, Kind: mapx.KindCode}, {ID: test, Kind: mapx.KindTest}}}
	if _, err := gate.RefreshStamps(g, root, mod, true); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, mod), []byte(newMod), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, mod, test, g
}

func TestStampRefreshDryRunOnlyLists(t *testing.T) {
	t.Run("CNSTC-B07: The refresh in dry-run lists the doubles and writes nothing", func(t *testing.T) {})
	root, mod, test, g := refreshFixture(t, "export function useBalance(id) {\n  return fetchBalance(id, 1)\n}\n")
	before := readQ(t, filepath.Join(root, test))
	out := captureStdout(t, func() {
		if err := refreshStamps(root, g, []string{mod}, true); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, mod+" changed — 1 double(s)") || !strings.Contains(out, "would update the stamps above (--dry-run: nothing written).") {
		t.Errorf("dry-run must list the double and say nothing was written:\n%s", out)
	}
	if readQ(t, filepath.Join(root, test)) != before {
		t.Error("dry-run must not rewrite the stamp")
	}
}

func TestStampRefreshWithoutHeadVersion(t *testing.T) {
	t.Run("CNSTC-B06: A refreshed block with no HEAD version says there is nothing to compare", func(t *testing.T) {})
	root, mod, test, g := refreshFixture(t, "export function useBalance(id) {\n  return fetchBalance(id, 1)\n}\n")
	before := readQ(t, filepath.Join(root, test))
	out := captureStdout(t, func() {
		if err := refreshStamps(root, g, []string{mod}, false); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "(no HEAD version of the block to compare)") {
		t.Errorf("outside git there is no HEAD block to diff:\n%s", out)
	}
	if readQ(t, filepath.Join(root, test)) == before {
		t.Error("the stamp should have been updated")
	}
}

func TestStampRefreshAnchorGoneIsNotRefreshed(t *testing.T) {
	t.Run("CNSTC-E03: A stamp whose anchor is gone is reported and not refreshed", func(t *testing.T) {})
	root, mod, test, g := refreshFixture(t, "export function useSaldo(id) {\n  return fetchBalance(id)\n}\n")
	before := readQ(t, filepath.Join(root, test))
	out := captureStdout(t, func() {
		if err := refreshStamps(root, g, []string{mod}, false); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, test+":1  `export function useBalance(id) {`\n      NOT refreshed:") {
		t.Errorf("a stamp whose anchor is gone must be reported as not refreshed:\n%s", out)
	}
	if readQ(t, filepath.Join(root, test)) != before {
		t.Error("a stamp whose anchor is gone must not be rewritten")
	}
}

func TestStampNeverRewritesADivergentStamp(t *testing.T) {
	t.Run("CNSTC-X01: A divergent stamp is left as it is", func(t *testing.T) {})
	dir, test := stampFixture(t)
	divergent := "// @contract: src/hooks/balance.ts | export function useBalance(id) { | 3 | 00000000\n" +
		"jest.mock('@/src/hooks/balance', () => ({ useBalance: jest.fn() }))\n"
	if err := os.WriteFile(filepath.Join(dir, test), []byte(divergent), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runQ(t, newStampCmd(), "--root", dir, test)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "wrote 0 stamp(s) in 0 file(s)") {
		t.Errorf("a stamped double gets no new stamp:\n%s", out)
	}
	if readQ(t, filepath.Join(dir, test)) != divergent {
		t.Error("a divergent stamp must never be rewritten by stamp")
	}
}

func TestStampRefreshHandedATest(t *testing.T) {
	t.Run("CNSTC-B09: The refresh handed a test names the modules to refresh", func(t *testing.T) {})
	root, mod, test, g := refreshFixture(t, "export function useBalance(id) {\n  return 2\n}\n")
	before, _ := os.ReadFile(filepath.Join(root, test))
	out := captureStdout(t, func() {
		if err := refreshStamps(root, g, []string{test}, false); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, test+" holds doubles stamped against other modules") || !strings.Contains(out, "anchors stamp --refresh "+mod) {
		t.Errorf("the refresh names the module and the command:\n%s", out)
	}
	if after, _ := os.ReadFile(filepath.Join(root, test)); string(after) != string(before) {
		t.Error("handed the test, the refresh changes nothing")
	}
}
