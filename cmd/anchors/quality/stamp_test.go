package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	got := blockDiff("a\nb\nc", "a\nB\nc\nd")
	want := []string{"- b", "+ B", "+ d"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("blockDiff = %v, want %v", got, want)
	}
}
