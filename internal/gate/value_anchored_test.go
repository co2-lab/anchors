package gate

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const anchorPattern = `@code-reference-\[([^\]]+)\]-\[([^\]]+)\]`

func cfgAnchor() *config.Config {
	return &config.Config{Derived: &config.Derived{ValueAnchor: anchorPattern}}
}

// project writes the files and returns the root and a map with them as nodes (`.spec.md`
// as specs, the rest as code). Each call builds a NEW graph, so the per-map index cache
// never leaks between tests.
func project(t *testing.T, files map[string]string) (string, *mapx.Graph) {
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
		if strings.HasSuffix(name, ".spec.md") {
			kind = mapx.KindSpec
		}
		g.Nodes = append(g.Nodes, mapx.Node{ID: name, Kind: kind})
	}
	return root, g
}

func runOn(t *testing.T, root string, g *mapx.Graph, file string) (Verdict, string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, file))
	if err != nil {
		t.Fatal(err)
	}
	return checkValueAnchored(string(b), mapx.Node{ID: file, Kind: mapx.KindCode}, root, g, cfgAnchor())
}

// --- LOCAL: the declaration against the code line below it ---

func TestValueAnchored_declarationMatchingItsLinePasses(t *testing.T) {
	t.Run("VLANV-B01: a declaration whose next code line contains the value passes", func(t *testing.T) {})
	root, g := project(t, map[string]string{"theme.ts": "// @code-reference-[COLOR-OK]-[#1F8A5B]\nsuccess: '#1F8A5B',\n"})
	if v, msg := runOn(t, root, g, "theme.ts"); v != Pass {
		t.Errorf("the line carries the value and the verdict was %v: %s", v, msg)
	}
}

func TestValueAnchored_declarationThatLiesFails(t *testing.T) {
	t.Run("VLANV-B02: a declaration whose next code line does not contain the value fails", func(t *testing.T) {})
	root, g := project(t, map[string]string{"theme.ts": "// @code-reference-[COLOR-OK]-[#1F8A5B]\nsuccess: '#2A9D6B',\n"})
	v, msg := runOn(t, root, g, "theme.ts")
	if v != Fail {
		t.Fatalf("the line says another value and the verdict was %v", v)
	}
	for _, want := range []string{"COLOR-OK", "#1F8A5B", "#2A9D6B"} {
		if !strings.Contains(msg, want) {
			t.Errorf("the verdict does not show %q: %q", want, msg)
		}
	}
}

// Comment lines between the declaration and the code are skipped: declarations may be
// stacked, and an explanation may sit between them.
func TestValueAnchored_commentLinesAreSkipped(t *testing.T) {
	t.Run("VLANV-B03: comment and blank lines between declaration and code are skipped", func(t *testing.T) {})
	body := "// @code-reference-[COLOR-OK]-[#1F8A5B]\n" +
		"// @code-reference-[COLOR-BG]-[#FFFFFF]\n" +
		"// the brand pair, used together\n\n" +
		"pair: ['#1F8A5B', '#FFFFFF'],\n"
	root, g := project(t, map[string]string{"theme.ts": body})
	if v, msg := runOn(t, root, g, "theme.ts"); v != Pass {
		t.Errorf("both declarations annotate the code line below and it holds both: %v / %s", v, msg)
	}
}

func TestValueAnchored_declarationWithNoCodeBelowFails(t *testing.T) {
	t.Run("VLANV-B07: a declaration with no code line below it fails", func(t *testing.T) {})
	root, g := project(t, map[string]string{"theme.ts": "x = 1\n// @code-reference-[COLOR-OK]-[#1F8A5B]\n"})
	v, msg := runOn(t, root, g, "theme.ts")
	if v != Fail || !strings.Contains(msg, "COLOR-OK") {
		t.Errorf("a declaration with nothing below it annotates nothing: %v / %s", v, msg)
	}
}

// --- ACROSS: every declaration of the same key declares the same value ---

// The propagation: the colour changed in one place, declaration included, and the other
// place stayed behind. Each file is locally consistent — only the comparison catches it.
func TestValueAnchored_divergentCopiesAreReported(t *testing.T) {
	t.Run("VLANV-B04: declarations of the same key with different values are reported", func(t *testing.T) {})
	root, g := project(t, map[string]string{
		"theme.ts":  "// @code-reference-[COLOR-OK]-[#2A9D6B]\nsuccess: '#2A9D6B',\n",
		"banner.ts": "// @code-reference-[COLOR-OK]-[#1F8A5B]\nfill: '#1F8A5B',\n",
	})
	v, msg := runOn(t, root, g, "banner.ts")
	if v != Fail {
		t.Fatalf("the copies disagree and the verdict was %v", v)
	}
	for _, want := range []string{"theme.ts", "banner.ts", "#2A9D6B", "#1F8A5B"} {
		if !strings.Contains(msg, want) {
			t.Errorf("the verdict does not list %q: %q", want, msg)
		}
	}
}

func TestValueAnchored_agreeingCopiesPass(t *testing.T) {
	root, g := project(t, map[string]string{
		"theme.ts":  "// @code-reference-[COLOR-OK]-[#1F8A5B]\nsuccess: '#1F8A5B',\n",
		"banner.ts": "// @code-reference-[COLOR-OK]-[#1F8A5B]\nfill: '#1F8A5B',\n",
	})
	if v, msg := runOn(t, root, g, "banner.ts"); v != Pass {
		t.Errorf("the copies agree and the verdict was %v: %s", v, msg)
	}
}

// --- SPEC: the key is a rule that declares its value ---

// The rule is the source: the value changed in the spec, and every place still carrying
// the old one is reported — even when all the copies agree with each other.
func TestValueAnchored_ruleValueInTheSpecIsTheSource(t *testing.T) {
	t.Run("VLANV-B05: a declaration disagreeing with the value its rule declares fails", func(t *testing.T) {})
	root, g := project(t, map[string]string{
		"tokens.spec.md": "| Rule | Token | Value |\n| --- | --- | --- |\n| `TKNSX-R01` | success | `#2A9D6B` |\n",
		"theme.ts":       "// @code-reference-[TKNSX-R01]-[#1F8A5B]\nsuccess: '#1F8A5B',\n",
		"banner.ts":      "// @code-reference-[TKNSX-R01]-[#1F8A5B]\nfill: '#1F8A5B',\n",
	})
	v, msg := runOn(t, root, g, "theme.ts")
	if v != Fail {
		t.Fatalf("the rule now says #2A9D6B and the verdict was %v", v)
	}
	if !strings.Contains(msg, "#2A9D6B") || !strings.Contains(msg, "tokens.spec.md") {
		t.Errorf("the verdict does not show the rule's value and where it lives: %q", msg)
	}
}

func TestValueAnchored_ruleValueMatchingPasses(t *testing.T) {
	root, g := project(t, map[string]string{
		"tokens.spec.md": "| Rule | Token | Value |\n| --- | --- | --- |\n| `TKNSX-R01` | success | `#1F8A5B` |\n",
		"theme.ts":       "// @code-reference-[TKNSX-R01]-[#1F8A5B]\nsuccess: '#1F8A5B',\n",
	})
	if v, msg := runOn(t, root, g, "theme.ts"); v != Pass {
		t.Errorf("the declaration matches the rule and the verdict was %v: %s", v, msg)
	}
}

// Only a table cell that is ENTIRELY a backticked token is a declared value. A heading
// or a prose cell carries backticked identifiers all the time, and reading them as values
// would charge every declaration pointing at an ordinary rule.
func TestValueAnchored_proseRuleDeclaresNoValue(t *testing.T) {
	t.Run("VLANV-X01: a rule whose line declares no value is not charged against the spec", func(t *testing.T) {})
	root, g := project(t, map[string]string{
		"credit.spec.md": "### CREDX-B01 — validates the `limit` before submitting\n\n" +
			"| `CREDX-B02` | uses `calcLimit` to decide |\n",
		"a.ts": "// @code-reference-[CREDX-B01]-[500]\nconst max = 500\n// @code-reference-[CREDX-B02]-[600]\nconst min = 600\n",
	})
	if v, msg := runOn(t, root, g, "a.ts"); v != Pass {
		t.Errorf("neither rule declares a value, and the declarations were charged: %v / %s", v, msg)
	}
}

// --- what the gate is NOT ---

// A value nobody declared is not charged: anchoring is for REPLICATED keys, and whoever
// replicates declares.
func TestValueAnchored_undeclaredValuesAreNotCharged(t *testing.T) {
	t.Run("VLANV-X02: a literal with no declaration is not charged", func(t *testing.T) {})
	root, g := project(t, map[string]string{"windows.ts": "export const WINDOWS = [\n  '15m',\n  '1h',\n]\n"})
	if v, _ := runOn(t, root, g, "windows.ts"); v != Skip {
		t.Errorf("no declaration in the file: expected Skip, got %v", v)
	}
}

func TestValueAnchored_skips(t *testing.T) {
	t.Run("VLANV-B06: without a declared pattern the gate skips and names the setting", func(t *testing.T) {})
	root, g := project(t, map[string]string{"a.ts": "// @code-reference-[K]-[v]\nv\n"})
	v, msg := checkValueAnchored("// @code-reference-[K]-[v]\nv\n", mapx.Node{ID: "a.ts", Kind: mapx.KindCode}, root, g, &config.Config{})
	if v != Skip || !strings.Contains(msg, "derived.value_anchor") {
		t.Errorf("no pattern: expected Skip naming the setting, got %v: %s", v, msg)
	}
	if v, _ := checkValueAnchored("x", mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}, root, g, cfgAnchor()); v != Skip {
		t.Errorf("not a code file and the verdict was %v", v)
	}
	t.Run("VLANV-I01: a pattern with fewer than two capture groups is treated as not declared", func(t *testing.T) {
		oneGroup := &config.Config{Derived: &config.Derived{ValueAnchor: `@code-reference-\[([^\]]+)\]`}}
		if v, _ := checkValueAnchored("// @code-reference-[K]\nv\n", mapx.Node{ID: "a.ts", Kind: mapx.KindCode}, root, g, oneGroup); v != Skip {
			t.Errorf("a one-group pattern asserts no value and must not enable the gate: %v", v)
		}
	})

	t.Run("VLANV-I03: with no built map the verdict is never approval", func(t *testing.T) {
		if v, _ := checkValueAnchored("// @code-reference-[K]-[v]\nv\n", mapx.Node{ID: "a.ts", Kind: mapx.KindCode}, root, nil, cfgAnchor()); v == Pass {
			t.Error("with no map the gate approved — it could not look across")
		}
	})
}

// Every declaration of the project is indexed once per graph instance, not once per file.
func TestValueAnchored_indexIsCachedPerGraph(t *testing.T) {
	t.Run("VLANV-I02: every declaration of the project is indexed once per map", func(t *testing.T) {})
	root, g := project(t, map[string]string{"theme.ts": "// @code-reference-[COLOR-OK]-[#1F8A5B]\nsuccess: '#1F8A5B',\n"})
	re := regexp.MustCompile(anchorPattern)
	first := anchorIndexFor(root, g, re)
	second := anchorIndexFor(root, g, re)
	if first != second {
		t.Fatalf("expected identical index pointer for the same graph instance, got %p and %p", first, second)
	}
}

// A declaration is a comment of its own. The comment EXPLAINING the syntax carries the
// anchor too, and reading it as a declaration charged the explanation — measured in the
// project that adopted the gate: key "key", failed.
func TestValueAnchored_anchorInsideProseIsAMention(t *testing.T) {
	t.Run("VLANV-B08: only an anchor standing alone on a comment line is a declaration", func(t *testing.T) {})
	// Text AFTER the anchor, and text BEFORE it: each is caught by its own half of the rule.
	prose := "// @code-reference-[key]-[value] is the syntax each entry carries\n" +
		"// the syntax each entry carries is @code-reference-[key]-[value]\n" +
		"/* e.g. `@code-reference-[key]-[value]` */\n" +
		"const doc = '@code-reference-[key]-[value]'\n" +
		"export const C = {}\n"
	root, g := project(t, map[string]string{"seal.ts": prose})
	if v, msg := runOn(t, root, g, "seal.ts"); v != Skip {
		t.Errorf("mentions of the syntax were read as declarations: %v / %s", v, msg)
	}

	// The standalone forms still declare — in every comment style, closing marker included.
	decl := "// @code-reference-[SEAL]-[#0B1F3A]\nnavy: '#0B1F3A',\n" +
		"/* @code-reference-[SEAL]-[#0B1F3A] */\nfill: '#0B1F3A',\n" +
		"# @code-reference-[SEAL]-[#0B1F3A]\nNAVY = '#0B1F3A'\n"
	root, g = project(t, map[string]string{"seal.ts": decl})
	if v, msg := runOn(t, root, g, "seal.ts"); v != Pass {
		t.Errorf("standalone declarations stopped being read: %v / %s", v, msg)
	}
	decl = "// @code-reference-[SEAL]-[#0B1F3A]\nnavy: '#FFFFFF',\n"
	root, g = project(t, map[string]string{"seal.ts": decl})
	if v, _ := runOn(t, root, g, "seal.ts"); v != Fail {
		t.Errorf("a standalone declaration that lies must still fail: %v", v)
	}
}
