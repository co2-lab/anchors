package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func stampCfg() *config.Config {
	return &config.Config{Derived: &config.Derived{
		MockDetect:   `(?:jest|vi)\.mock\(['"]([^'"]+)`,
		ExportDetect: `^export\s+(?:async\s+)?(?:function|const)\s+(\w+)`,
	}}
}

const realHooks = `import { useQuery } from 'react-query'

export function useBalance(id) {
  return useQuery(['balance', id], () => fetchBalance(id))
}

export function useLimits(id) {
  return useQuery(['limits', id], () => fetchLimits(id))
}
`

func stampProject(t *testing.T, files map[string]string) (string, *mapx.Graph) {
	t.Helper()
	root := t.TempDir()
	g := &mapx.Graph{}
	for name, body := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		kind := mapx.KindCode
		if strings.Contains(name, ".test.") {
			kind = mapx.KindTest
		}
		g.Nodes = append(g.Nodes, mapx.Node{ID: name, Kind: kind})
	}
	return root, g
}

// THE ROUND TRIP: what the generator writes is what the gate recomputes. After stamping,
// the gate passes; after the stamped member changes, it fails.
func TestGenerateStamps_roundTripWithTheGate(t *testing.T) {
	t.Run("MKSTP-B01: A generated stamp passes the gate, and a change to its snippet fails it", func(t *testing.T) {})
	test := "jest.mock('@/src/hooks/balance', () => ({\n  useBalance: jest.fn(),\n}))\n\nit('shows', () => {})\n"
	root, g := stampProject(t, map[string]string{
		"apps/mobile/src/hooks/balance.ts":      realHooks,
		"apps/mobile/src/screens/Home.test.tsx": test,
	})
	testID := "apps/mobile/src/screens/Home.test.tsx"
	node := mapx.Node{ID: testID, Kind: mapx.KindTest}

	if v, _ := checkMockStamped(test, node, root, g, stampCfg()); v != Fail {
		t.Fatalf("before stamping the double is unstamped and the gate should fail, got %v", v)
	}

	out, written, skipped, err := GenerateStamps(test, testID, root, g, stampCfg())
	if err != nil || len(written) != 1 || len(skipped) != 0 {
		t.Fatalf("expected one stamp and no skip, got written=%v skipped=%v err=%v", written, skipped, err)
	}
	if written[0].Anchor != "export function useBalance(id) {" {
		t.Errorf("the stamp should anchor on the export the factory names, got %q", written[0].Anchor)
	}
	if v, msg := checkMockStamped(out, node, root, g, stampCfg()); v != Pass {
		t.Fatalf("after stamping the gate should pass, got %v: %s\n%s", v, msg, out)
	}

	// The stamped member changes — the double may be stale now, and the gate must say so.
	changed := strings.Replace(realHooks, "['balance', id]", "['balance', id, 'v2']", 1)
	if err := os.WriteFile(filepath.Join(root, "apps/mobile/src/hooks/balance.ts"), []byte(changed), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, _ := checkMockStamped(out, node, root, g, stampCfg()); v != Fail {
		t.Errorf("the stamped member changed and the gate still passed: %v", v)
	}
}

// A change to a member the double does NOT cover does not make the stamp diverge — the
// point of stamping per member.
func TestGenerateStamps_perMemberScope(t *testing.T) {
	t.Run("MKSTP-B03: A stamp per factory key covers only that export", func(t *testing.T) {})
	test := "jest.mock('@/src/hooks/balance', () => ({ useBalance: jest.fn() }))\n"
	root, g := stampProject(t, map[string]string{"src/hooks/balance.ts": realHooks, "src/Home.test.tsx": test})
	out, _, _, err := GenerateStamps(test, "src/Home.test.tsx", root, g, stampCfg())
	if err != nil {
		t.Fatal(err)
	}
	other := strings.Replace(realHooks, "['limits', id]", "['limits', id, 'v2']", 1)
	if err := os.WriteFile(filepath.Join(root, "src/hooks/balance.ts"), []byte(other), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, msg := checkMockStamped(out, mapx.Node{ID: "src/Home.test.tsx", Kind: mapx.KindTest}, root, g, stampCfg()); v != Pass {
		t.Errorf("useLimits changed, the double covers only useBalance, and the gate failed: %v / %s", v, msg)
	}
}

// NEVER REWRITTEN: a stamp that exists is left alone, even when it diverges. Refreshing
// it would let the stamp certify itself.
func TestGenerateStamps_neverRewritesAnExistingStamp(t *testing.T) {
	t.Run("MKSTP-I01: an existing stamp is never rewritten", func(t *testing.T) {})
	t.Run("MKSTP-X01: Generator does not refresh a divergent stamp", func(t *testing.T) {})
	test := "// @contract: src/hooks/balance.ts | export function useBalance(id) { | 3 | deadbeef\n" +
		"jest.mock('@/src/hooks/balance', () => ({ useBalance: jest.fn() }))\n"
	root, g := stampProject(t, map[string]string{"src/hooks/balance.ts": realHooks, "src/Home.test.tsx": test})
	out, written, _, err := GenerateStamps(test, "src/Home.test.tsx", root, g, stampCfg())
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 0 || out != test {
		t.Errorf("an existing (and divergent) stamp was touched: written=%v\n%s", written, out)
	}
}

// With no factory key naming an export, one stamp covers the whole module.
func TestGenerateStamps_wholeModuleWhenNoKeyNamesAnExport(t *testing.T) {
	t.Run("MKSTP-B04: An automock gets one stamp over the whole module", func(t *testing.T) {})
	test := "jest.mock('@/src/hooks/balance')\n"
	root, g := stampProject(t, map[string]string{"src/hooks/balance.ts": realHooks, "src/Home.test.tsx": test})
	out, written, _, err := GenerateStamps(test, "src/Home.test.tsx", root, g, stampCfg())
	if err != nil || len(written) != 1 {
		t.Fatalf("expected one whole-module stamp: %v %v", written, err)
	}
	if written[0].Anchor != "import { useQuery } from 'react-query'" {
		t.Errorf("the whole-module stamp should anchor on the first usable line, got %q", written[0].Anchor)
	}
	if v, msg := checkMockStamped(out, mapx.Node{ID: "src/Home.test.tsx", Kind: mapx.KindTest}, root, g, stampCfg()); v != Pass {
		t.Errorf("the whole-module stamp did not pass the gate: %v / %s", v, msg)
	}
}

// Two workspaces with the same path: the test's own workspace wins; a real tie is skipped.
func TestGenerateStamps_ambiguousModule(t *testing.T) {
	t.Run("MKSTP-B02: An ambiguous specifier resolves to the test's workspace, or is skipped", func(t *testing.T) {})
	test := "jest.mock('@/src/hooks/balance')\n"
	root, g := stampProject(t, map[string]string{
		"apps/mobile/src/hooks/balance.ts":  realHooks,
		"apps/landing/src/hooks/balance.ts": realHooks,
		"apps/mobile/src/Home.test.tsx":     test,
	})
	_, written, skipped, err := GenerateStamps(test, "apps/mobile/src/Home.test.tsx", root, g, stampCfg())
	if err != nil || len(written) != 1 || written[0].File != "apps/mobile/src/hooks/balance.ts" {
		t.Fatalf("the mobile test should stamp the mobile module: written=%v skipped=%v", written, skipped)
	}

	_, written, skipped, _ = GenerateStamps(test, "tests/Home.test.tsx", root, g, stampCfg())
	if len(written) != 0 || len(skipped) != 1 {
		t.Errorf("a test outside both workspaces cannot decide and must skip: written=%v skipped=%v", written, skipped)
	}
}

func TestGenerateStamps_thirdPartyIsNotStamped(t *testing.T) {
	t.Run("MKSTP-B05: A third-party double is not stamped", func(t *testing.T) {})
	test := "jest.mock('react-query')\n"
	root, g := stampProject(t, map[string]string{"src/Home.test.tsx": test})
	_, written, skipped, _ := GenerateStamps(test, "src/Home.test.tsx", root, g, stampCfg())
	if len(written) != 0 || len(skipped) != 0 {
		t.Errorf("a third-party module is not governed and must be left alone: %v %v", written, skipped)
	}
}

// A line that occurs twice cannot anchor: the gate refuses an ambiguous anchor, and a stamp
// on it would fail the moment it was written. The generator moves on to a unique line.
func TestGenerateStamps_ambiguousAnchorLineIsAvoided(t *testing.T) {
	t.Run("MKSTP-I02: A repeated line never anchors a stamp", func(t *testing.T) {})
	module := "'use client'\n'use client'\n" + realHooks
	test := "jest.mock('@/src/hooks/balance')\n"
	root, g := stampProject(t, map[string]string{"src/hooks/balance.ts": module, "src/Home.test.tsx": test})
	out, written, _, err := GenerateStamps(test, "src/Home.test.tsx", root, g, stampCfg())
	if err != nil || len(written) != 1 {
		t.Fatalf("expected one stamp: %v %v", written, err)
	}
	if written[0].Anchor == "'use client'" {
		t.Errorf("the stamp anchored on a line that occurs twice")
	}
	if v, msg := checkMockStamped(out, mapx.Node{ID: "src/Home.test.tsx", Kind: mapx.KindTest}, root, g, stampCfg()); v != Pass {
		t.Errorf("the generated stamp did not pass the gate: %v / %s", v, msg)
	}
}

// The header stays out of the hash: `check --fix` rewrites `updated_at:` on every edit, and
// a stamp covering it failed on any change to the module — even one the double never saw.
func TestGenerateStamps_wholeModuleSkipsTheAnchorsHeader(t *testing.T) {
	t.Run("MKSTP-B06: The header stays out of a whole-module stamp", func(t *testing.T) {})
	header := func(date string) string {
		return "'use client'\n// @anchors\n//   ref: BALNC\n//   updated_at: " + date + "\n//\n// notes about the unit\n"
	}
	test := "jest.mock('@/src/hooks/balance')\n"
	root, g := stampProject(t, map[string]string{"src/hooks/balance.ts": header("2026-09-01") + realHooks, "src/Home.test.tsx": test})
	out, written, _, err := GenerateStamps(test, "src/Home.test.tsx", root, g, stampCfg())
	if err != nil || len(written) != 1 {
		t.Fatalf("expected one whole-module stamp: %v %v", written, err)
	}
	if written[0].Anchor != "import { useQuery } from 'react-query'" {
		t.Errorf("the stamp should anchor on the first line after the header, got %q", written[0].Anchor)
	}
	node := mapx.Node{ID: "src/Home.test.tsx", Kind: mapx.KindTest}

	// Only the header changes: the stamp still holds.
	if err := os.WriteFile(filepath.Join(root, "src/hooks/balance.ts"), []byte(header("2026-09-23")+realHooks), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, msg := checkMockStamped(out, node, root, g, stampCfg()); v != Pass {
		t.Errorf("a header-only edit made the stamp diverge: %v / %s", v, msg)
	}

	// The body changes: it still diverges — skipping the header must not blind the stamp.
	changed := strings.Replace(realHooks, "fetchLimits(id)", "fetchLimits(id, true)", 1)
	if err := os.WriteFile(filepath.Join(root, "src/hooks/balance.ts"), []byte(header("2026-09-23")+changed), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, _ := checkMockStamped(out, node, root, g, stampCfg()); v != Fail {
		t.Errorf("a body edit did not make the whole-module stamp diverge: %v", v)
	}
}

// A key ADDED to an already-stamped double gets its own stamp; the existing stamps are left
// alone. Measured in the reference project: `s3ObjectExists` added to a mock that already
// stamped s3GetObject/s3PutObject came out with no stamp.
func TestGenerateStamps_newKeyOnAStampedDoubleGetsAStamp(t *testing.T) {
	t.Run("MKSTP-B07: A key added to a stamped double gets a stamp", func(t *testing.T) {})
	node := mapx.Node{ID: "src/Home.test.tsx", Kind: mapx.KindTest}
	first := "jest.mock('@/src/hooks/balance', () => ({ useBalance: jest.fn() }))\n"
	root, g := stampProject(t, map[string]string{"src/hooks/balance.ts": realHooks, "src/Home.test.tsx": first})
	stamped, _, _, err := GenerateStamps(first, "src/Home.test.tsx", root, g, stampCfg())
	if err != nil {
		t.Fatal(err)
	}

	grown := strings.Replace(stamped, "useBalance: jest.fn() }", "useBalance: jest.fn(), useLimits: jest.fn() }", 1)
	out, written, _, err := GenerateStamps(grown, "src/Home.test.tsx", root, g, stampCfg())
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 1 || !strings.Contains(written[0].Anchor, "useLimits") {
		t.Fatalf("the new key should get exactly one stamp, on useLimits: %+v", written)
	}
	for _, l := range strings.Split(grown, "\n") {
		if strings.Contains(l, "@contract:") && !strings.Contains(out, l) {
			t.Errorf("an existing stamp was touched: %q", l)
		}
	}
	if v, msg := checkMockStamped(out, node, root, g, stampCfg()); v != Pass {
		t.Errorf("the grown double does not pass the gate: %v / %s", v, msg)
	}

	// Idempotent: nothing more to add.
	if _, again, _, _ := GenerateStamps(out, "src/Home.test.tsx", root, g, stampCfg()); len(again) != 0 {
		t.Errorf("a second run added stamps again: %+v", again)
	}

	// A whole-module stamp already covers every member: a new key adds nothing.
	auto := "jest.mock('@/src/hooks/balance')\n"
	whole, _, _, _ := GenerateStamps(auto, "src/Home.test.tsx", root, g, stampCfg())
	withKey := strings.Replace(whole, "jest.mock('@/src/hooks/balance')", "jest.mock('@/src/hooks/balance', () => ({ useLimits: jest.fn() }))", 1)
	if _, extra, _, _ := GenerateStamps(withKey, "src/Home.test.tsx", root, g, stampCfg()); len(extra) != 0 {
		t.Errorf("a member inside a whole-module stamp got a stamp of its own: %+v", extra)
	}
}

func TestGenerateStamps_errors(t *testing.T) {
	t.Run("MKSTP-E01: A double detector that does not compile stops the generator with an error", func(t *testing.T) {
		for _, pattern := range []string{`[`, `jest\.mock`} {
			cfg := &config.Config{Derived: &config.Derived{MockDetect: pattern}}
			out, written, _, err := GenerateStamps("jest.mock('x')\n", "test.ts", "", &mapx.Graph{}, cfg)
			if err == nil || strings.Contains(err.Error(), "does not declare") || out != "jest.mock('x')\n" || len(written) != 0 {
				t.Fatalf("pattern %q: expected the pattern's own error and the content unchanged, got err=%v written=%v", pattern, err, written)
			}
		}
	})

	t.Run("MKSTP-E02: An undeclared double detector stops the generator with an error", func(t *testing.T) {
		cfg := &config.Config{Derived: &config.Derived{}}
		_, _, _, err := GenerateStamps("content", "test.ts", "", &mapx.Graph{}, cfg)
		if err == nil || !strings.Contains(err.Error(), "mock_detect") {
			t.Fatalf("expected an error naming mock_detect, got %v", err)
		}
	})

	t.Run("MKSTP-E03: A module missing from disk is skipped and the other doubles are still stamped", func(t *testing.T) {
		test := "jest.mock('@/src/hooks/balance')\njest.mock('@/src/hooks/other')\n"
		root, g := stampProject(t, map[string]string{
			"src/hooks/balance.ts": realHooks,
			"src/hooks/other.ts":   "export const other = 1\n",
			"src/Home.test.tsx":    test,
		})
		if err := os.Remove(filepath.Join(root, "src/hooks/balance.ts")); err != nil {
			t.Fatal(err)
		}
		_, written, skipped, err := GenerateStamps(test, "src/Home.test.tsx", root, g, stampCfg())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(skipped) != 1 || !strings.Contains(skipped[0].Reason, "src/hooks/balance.ts") {
			t.Fatalf("expected the missing module skipped naming its file, got %v", skipped)
		}
		if len(written) != 1 || written[0].File != "src/hooks/other.ts" {
			t.Fatalf("expected the present module still stamped, got %v", written)
		}
	})
}
