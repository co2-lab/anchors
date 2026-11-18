// @anchors
//   code: SYMXT
//   ref: SYMXS

package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const symbolReader = "// @anchors\n//   ref: ITEFH\n\nconst cfg = { readerOptions: { verbosity: 0 } };\n" +
	"// the common reader reads without compiling (#isEvalSupported false), via $cfg.readerOptions\n" +
	"export function read() { return cfg.readerOptions; }\n"

func symbolConfig() *config.Config {
	return &config.Config{Dialect: &config.Dialect{Comments: append(config.FamilyComments("ts"), config.FamilyComments("python")...)}}
}

// symbolCase confronts the file rel holding content, in a project whose code is reader.ts.
func symbolCase(t *testing.T, rel string, kind mapx.Kind, content string) (Verdict, string) {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "reader.ts", symbolReader)
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "reader.ts", Kind: mapx.KindCode}}}
	if rel != "reader.ts" {
		writeFile(t, root, rel, content)
		g.Nodes = append(g.Nodes, mapx.Node{ID: rel, Kind: kind})
	}
	return checkSymbolExists(content, mapx.Node{ID: rel, Kind: kind}, root, g, symbolConfig())
}

func TestSymbolExists_aMarkedNameIsInTheCode(t *testing.T) {
	t.Run("SYMXS-B01: A marked name the code no longer has fails, naming where it is cited", func(t *testing.T) {})
	v, msg := symbolCase(t, "reader.ts", mapx.KindCode, symbolReader)
	if v != Fail || !strings.Contains(msg, "reader.ts:5 `isEvalSupported`") || strings.Contains(msg, "readerOptions") {
		t.Errorf("isEvalSupported is gone, cfg.readerOptions is not: %v %s", v, msg)
	}
	fixed := strings.Replace(symbolReader, "#isEvalSupported", "#verbosity", 1)
	if v, msg := symbolCase(t, "reader.ts", mapx.KindCode, fixed); v != Pass {
		t.Errorf("verbosity is in the code: %v %s", v, msg)
	}
}

func TestSymbolExists_onlyRealMarks(t *testing.T) {
	t.Run("SYMXS-B02: Only real marks in comments and in the spec's prose are confronted", func(t *testing.T) {})
	code := "const s = \"#ghostInString\";\n" +
		"// see [docs](#ghostAnchor), https://x.io/a#ghostFrag, #ghost-id, a#ghostGlued and @#$ghost\n" +
		"// #region ghost section, colored #C9BCFF and #fff\n// #endregion\n"
	if v, msg := symbolCase(t, "other.ts", mapx.KindCode, code); v != Skip {
		t.Errorf("no real mark in the ts file: %v %s", v, msg)
	}
	if v, msg := symbolCase(t, "tool.py", mapx.KindCode, "#ghostOpener is the comment's own\nx = 1\n"); v != Skip {
		t.Errorf("the # that opens a comment is no mark: %v %s", v, msg)
	}
	spec := "# Reader\n\n```sh\necho $ghostShell\n```\n\nThe reader calls @@read.\n"
	if v, msg := symbolCase(t, "reader.spec.md", mapx.KindSpec, spec); v != Pass {
		t.Errorf("the fence is not confronted and @@read is in the code: %v %s", v, msg)
	}
	gone := strings.Replace(spec, "@@read", "@@readAll", 1)
	if v, msg := symbolCase(t, "reader.spec.md", mapx.KindSpec, gone); v != Fail || !strings.Contains(msg, "reader.spec.md:7 `readAll`") {
		t.Errorf("readAll is not in the code: %v %s", v, msg)
	}
}

func TestSymbolExists_aNameCitedOnPurpose(t *testing.T) {
	t.Run("SYMXS-B03: A line that cites a name on purpose is waived", func(t *testing.T) {})
	waived := strings.Replace(symbolReader, "readerOptions\nexport", "readerOptions @no-symbol-check: pdfjs 5 removed the option\nexport", 1)
	if v, msg := symbolCase(t, "reader.ts", mapx.KindCode, waived); v == Fail {
		t.Errorf("the line is waived: %v %s", v, msg)
	}
}

func TestSymbolExists_outsideItsSubject(t *testing.T) {
	t.Run("SYMXS-B04: No mark, or another kind of artifact, is skipped", func(t *testing.T) {})
	if v, msg := symbolCase(t, "plain.ts", mapx.KindCode, "// reads the file\nconst x = 1;\n"); v != Skip {
		t.Errorf("no mark: %v %s", v, msg)
	}
	if v, msg := symbolCase(t, "reader.feature", mapx.KindFeature, "Feature: x\n  # #ghost\n"); v != Skip {
		t.Errorf("a feature: %v %s", v, msg)
	}
}

func TestSymbolExists_errors(t *testing.T) {
	t.Run("SYMXS-E01: A code file whose comments the project does not declare is pending", func(t *testing.T) {})
	t.Run("SYMXS-E02: Without a map the gate is pending", func(t *testing.T) {})
	if v, msg := symbolCase(t, "main.zig", mapx.KindCode, "// #ghost\n"); v != Pending || !strings.Contains(msg, "main.zig") {
		t.Errorf("no comment syntax for zig: %v %s", v, msg)
	}
	if v, _ := checkSymbolExists(symbolReader, mapx.Node{ID: "reader.ts", Kind: mapx.KindCode}, t.TempDir(), nil, symbolConfig()); v != Pending {
		t.Errorf("no map: %v", v)
	}
}
