// @anchors
//   code: DCTDP
//   ref: DCGDP

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const chainScreen = `import {
  PALETTE,
  withAlpha as alpha,
} from '../theme/tokens'
import useToast from '@/hooks/useToast' // @dep: WRONG
import React from 'react'
import type { Goal } from './types' // @no-dep: types only
const legacy = { case: 1, member: 2 }
export { withAlpha } from '../theme/tokens'

export function Arena() {}
`

const chainTokens = `// @used-by: OLDXX
export const PALETTE = {}

export function withAlpha() {}

// @used-by: ARSCR
export const unused = 1
`

// chainProject is a TS project: a screen importing the tokens (names spanning lines, one
// aliased), a hook through an alias with a wrong flag, a package outside, and a waived import.
func chainProject(t *testing.T) (string, *mapx.Graph, *config.Config) {
	t.Helper()
	root := t.TempDir()
	for rel, body := range map[string]string{
		"src/ui/Arena.tsx":      chainScreen,
		"src/theme/tokens.ts":   chainTokens,
		"src/hooks/useToast.ts": "module.exports = function useToast() {}\n",
		"src/ui/types.ts":       "export type Goal = {}\n",
	} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(root, rel)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, rel), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := &mapx.Graph{Nodes: []mapx.Node{
		{ID: "src/ui/Arena.tsx", Kind: mapx.KindCode, FileCode: "ARSCR"},
		{ID: "src/theme/tokens.ts", Kind: mapx.KindCode, FileCode: "TOKNS"},
		{ID: "src/hooks/useToast.ts", Kind: mapx.KindCode, FileCode: "USTST"},
		{ID: "src/ui/types.ts", Kind: mapx.KindCode, FileCode: "TYPSU"},
	}}
	cfg := &config.Config{Dialect: &config.Dialect{Family: "ts", ImportResolve: &config.ImportResolve{Aliases: map[string]string{"@/": "src/"}}}}
	importersMu.Lock()
	delete(importersCache, root)
	importersMu.Unlock()
	return root, g, cfg
}

func read(t *testing.T, root, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestDepChain_importsAreReadAndResolved(t *testing.T) {
	t.Run("DCGDP-B01: The imports are read by the dialect's pattern and resolved to files of the map", func(t *testing.T) {})
	root, g, cfg := chainProject(t)
	imps := ImportsOf(read(t, root, "src/ui/Arena.tsx"), "src/ui/Arena.tsx", cfg.DialectFor(), g)
	if len(imps) != 5 {
		t.Fatalf("five import statements: %+v", imps)
	}
	if re := imps[4]; re.First != re.Line || strings.Join(re.Symbols, ",") != "withAlpha" {
		t.Errorf("an export-from on its own line is its own statement, with its own names only: %+v", re)
	}
	tok := imps[0]
	if tok.First != 1 || tok.Line != 4 || tok.Target != "src/theme/tokens.ts" || strings.Join(tok.Symbols, ",") != "PALETTE,withAlpha" {
		t.Errorf("an import whose names span lines is one statement, its path on its last line, the alias read as the name: %+v", tok)
	}
	if imps[1].Target != "src/hooks/useToast.ts" || imps[2].Target != "" {
		t.Errorf("an alias resolves into the project, a package stays outside: %+v %+v", imps[1], imps[2])
	}
	if ImportsOf("import x from './y'\n", "a.ts", config.Dialect{}, g) != nil {
		t.Error("with no pattern for the dialect, nothing is read")
	}
}

func TestDepChain_declaredAndHonored(t *testing.T) {
	t.Run("DCGDP-B02: dep-declared names each import of a governed file with no flag, and the code it would carry", func(t *testing.T) {})
	t.Run("DCGDP-B03: dep-honored names each flag whose code is not the one its import resolves to", func(t *testing.T) {})
	root, g, cfg := chainProject(t)
	n := g.Nodes[0]
	content := read(t, root, n.ID)
	v, msg := checkDepDeclared(content, n, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "4 `../theme/tokens` → TOKNS") || strings.Contains(msg, "react") || strings.Contains(msg, "types") {
		t.Errorf("only the unflagged import of a governed file is named: %v %s", v, msg)
	}
	v, msg = checkDepHonored(content, n, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "WRONG") || !strings.Contains(msg, "USTST") {
		t.Errorf("the wrong code is named with the right one: %v %s", v, msg)
	}
	if v, _ := checkDepDeclared("", mapx.Node{ID: "x_test.ts", Kind: mapx.KindTest}, root, g, cfg); v != Skip {
		t.Error("a test is tied by its ref, not by the chain")
	}
}

func TestDepChain_usedByAndTheFixers(t *testing.T) {
	t.Run("DCGDP-B04: used-by-declared names each imported symbol whose flag is missing or wrong, and each flag nobody imports", func(t *testing.T) {})
	t.Run("DCGDP-B05: The fixers write the dependency and used-by flags, correct a wrong code, and remove a stale flag", func(t *testing.T) {})
	root, g, cfg := chainProject(t)
	tokens := g.Nodes[1]
	v, msg := checkUsedByDeclared(read(t, root, tokens.ID), tokens, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "`PALETTE` says `@used-by: OLDXX`") || !strings.Contains(msg, "`withAlpha` is imported by ARSCR") || !strings.Contains(msg, "`unused`") || !strings.Contains(msg, "anchors guide header") {
		t.Errorf("wrong, missing and stale flags are named: %v %s", v, msg)
	}
	gates := []config.Gate{{Name: "dep-declared", Check: "dep-declared", On: []string{"code"}}, {Name: "used-by-declared", Check: "used-by-declared", On: []string{"code"}}}
	res := FixWithConfig(gates, g.Nodes, root, g, cfg)
	movedOf := map[string]bool{}
	for _, r := range res {
		movedOf[r.Target] = movedOf[r.Target] || r.LinesMoved
	}
	if movedOf["src/ui/Arena.tsx"] || !movedOf["src/theme/tokens.ts"] {
		t.Errorf("the flag on the import's line moves none, the used-by line inserted moves: %+v", movedOf)
	}
	arena := read(t, root, "src/ui/Arena.tsx")
	if !strings.Contains(arena, "} from '../theme/tokens' // @dep: TOKNS") || !strings.Contains(arena, "// @dep: USTST") || strings.Contains(arena, "WRONG") {
		t.Errorf("the flags are written and corrected:\n%s", arena)
	}
	if hook := read(t, root, "src/hooks/useToast.ts"); !strings.HasPrefix(hook, "// @used-by: ARSCR\nmodule.exports") {
		t.Errorf("a default import is the module's default, flagged above its declaration:\n%s", hook)
	}
	tok := read(t, root, "src/theme/tokens.ts")
	if !strings.Contains(tok, "// @used-by: ARSCR\nexport const PALETTE") || !strings.Contains(tok, "// @used-by: ARSCR\nexport function withAlpha") || strings.Contains(tok, "OLDXX") || strings.Contains(tok, "@used-by: ARSCR\nexport const unused") {
		t.Errorf("the used-by flags are exactly who imports each symbol:\n%s", tok)
	}
	for i, n := range g.Nodes[:2] {
		for _, check := range []func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string){checkDepDeclared, checkDepHonored, checkUsedByDeclared} {
			if v, msg := check(read(t, root, n.ID), n, root, g, cfg); v == Fail {
				t.Errorf("node %d after the fix: %s", i, msg)
			}
		}
	}
}

func TestDepChain_reExportsAndInlineImports(t *testing.T) {
	t.Run("DCGDP-B06: A re-export declares the names it lists, and an inline import brings the member it reads", func(t *testing.T) {})
	root, g, cfg := chainProject(t)
	for rel, body := range map[string]string{
		"src/theme/barrel.ts": "export { PALETTE } from './tokens' // @dep: TOKNS\nexport type { Goal, Mode } from '../ui/types' // @dep: TYPSU\nexport {\n  alpha,\n  beta as gamma,\n} from './more'\n",
		"src/ui/Shop.tsx":     "import { PALETTE, Mode, gamma } from '../theme/barrel'\n",
		"src/ui/Cart.tsx":     "type R = { item: import('../theme/barrel').Goal }\n",
	} {
		if err := os.WriteFile(filepath.Join(root, rel), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g.Nodes = append(g.Nodes,
		mapx.Node{ID: "src/theme/barrel.ts", Kind: mapx.KindCode, FileCode: "BARRL"},
		mapx.Node{ID: "src/ui/Shop.tsx", Kind: mapx.KindCode, FileCode: "SHOPU"},
		mapx.Node{ID: "src/ui/Cart.tsx", Kind: mapx.KindCode, FileCode: "CARTU"})
	if imps := ImportsOf(read(t, root, "src/ui/Cart.tsx"), "src/ui/Cart.tsx", cfg.DialectFor(), g); len(imps) != 1 || strings.Join(imps[0].Symbols, ",") != "Goal" {
		t.Errorf("the inline import brings the member it reads: %+v", imps)
	}
	gates := []config.Gate{{Name: "dep-declared", Check: "dep-declared", On: []string{"code"}}, {Name: "used-by-declared", Check: "used-by-declared", On: []string{"code"}}}
	FixWithConfig(gates, g.Nodes, root, g, cfg)
	barrel := read(t, root, "src/theme/barrel.ts")
	for _, want := range []string{
		"// @used-by: SHOPU (PALETTE)\nexport { PALETTE }",
		"// @used-by: CARTU (Goal)\n// @used-by: SHOPU (Mode)\nexport type { Goal, Mode }",
		"// @used-by: SHOPU (gamma)\nexport {\n  alpha,",
	} {
		if !strings.Contains(barrel, want) {
			t.Errorf("the flag stands above the list and names its symbol — want %q in:\n%s", want, barrel)
		}
	}
	if cart := read(t, root, "src/ui/Cart.tsx"); !strings.Contains(cart, ".Goal } // @dep: BARRL") {
		t.Errorf("the inline import carries its flag:\n%s", cart)
	}
	n := g.Nodes[len(g.Nodes)-3]
	if v, msg := checkUsedByDeclared(barrel, n, root, g, cfg); v == Fail {
		t.Errorf("after the fix the re-exports answer their importers: %s", msg)
	}
}
