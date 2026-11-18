// @anchors
//   code: CMTSA
//   ref: CMMRC

package config

import (
	"strings"
	"testing"
)

func withComments(families ...string) *Config {
	var cs []Comments
	for _, f := range families {
		cs = append(cs, FamilyComments(f)...)
	}
	return &Config{Dialect: &Dialect{Comments: cs}}
}

func TestCommentsFor(t *testing.T) {
	t.Run("CMMRC-B01: A file's syntax is the first entry that covers it, or Anchors' own for its formats", func(t *testing.T) {})
	c := withComments("ts", "python")
	for rel, want := range map[string]string{
		"src/a.tsx": "//", "main.py": "#", "x/a.spec.md": "<!--", "x/a.feature": "#", `win\b.ts`: "//",
	} {
		if got := c.CommentsFor(rel).LineMarker(); got != want {
			t.Errorf("%s: %q, want %q", rel, got, want)
		}
	}
	if c.CommentsFor("data/x.json") != nil {
		t.Error("a JSON file no entry covers has no syntax")
	}
}

func TestComments_blank(t *testing.T) {
	t.Run("CMMRC-B02: Blanking keeps strings and lines, and turns comments into spaces", func(t *testing.T) {})
	ts := withComments("ts").CommentsFor("a.ts")
	src := "const a = 'no // comment' // gone\n/* a\n b */ x()\nconst q = \"say \\\" // still\"\nconst t = `line\n// inside`\n"
	got := ts.Blank(src)
	for _, want := range []string{"const a = 'no // comment'", "\n      x()", "\"say \\\" // still\"", "`line\n// inside`"} {
		if !strings.Contains(got, want) {
			t.Errorf("blanked lacks %q:\n%s", want, got)
		}
	}
	for _, gone := range []string{"gone", "a\n b"} {
		if strings.Contains(got, gone) {
			t.Errorf("blanked keeps %q:\n%s", gone, got)
		}
	}
	var none *Comments
	if none.Blank(src) != src {
		t.Error("no syntax blanks nothing")
	}
}

func TestComments_lines(t *testing.T) {
	t.Run("CMMRC-B03: Lines that are only a comment, the code of a content, a line cut", func(t *testing.T) {})
	ts := withComments("ts").CommentsFor("a.ts")
	lines := []string{"/*", " * text", " */", "call() // why", ""}
	only := ts.CommentOnly(lines)
	if !only[0] || !only[1] || !only[2] || only[3] || only[4] {
		t.Errorf("comment-only lines: %v", only)
	}
	if got := ts.Code(strings.Join(lines, "\n")); got != "call()\n\n" {
		t.Errorf("code: %q", got)
	}
	if got := ts.CutLine("call() // why"); got != "call()" {
		t.Errorf("cut: %q", got)
	}
}

func TestComments_writer(t *testing.T) {
	t.Run("CMMRC-B04: The writer's marker and where a comment cannot go", func(t *testing.T) {})
	c := withComments("ts", "ocaml")
	ts, ml, md := c.CommentsFor("a.tsx"), c.CommentsFor("p.ml"), c.CommentsFor("a.md")
	if ts.LineMarker() != "//" || ts.CloseFor("//") != "" || ts.TakesLineComment("<View>") || !ts.TakesLineComment("call()") {
		t.Errorf("ts: %q, takes after a tag %v", ts.LineMarker(), ts.TakesLineComment("<View>"))
	}
	if ml.LineMarker() != "(*" || ml.CloseFor("(*") != " *)" || md.LineMarker() != "<!--" || md.CloseFor("<!--") != " -->" {
		t.Errorf("ocaml %q %q, markdown %q %q", ml.LineMarker(), ml.CloseFor("(*"), md.LineMarker(), md.CloseFor("<!--"))
	}
}

func TestComments_validation(t *testing.T) {
	t.Run("CMMRC-B05: A malformed entry refuses the config", func(t *testing.T) {})
	for entry, says := range map[string]string{
		"- files: '**/*.x'":                       "declares no marker",
		"- line: ['']":                            "empty marker",
		"- files: '[a'\n      line: ['#']":        "not a valid glob",
		"- line: ['#']\n      no_line_after: '['": "no_line_after` does not compile",
		"- block: [['/*']]":                       "[open, close]",
	} {
		_, err := load(t, "version: 8\nlayers: {}\ndialect:\n  comments:\n    "+entry+"\n")
		if err == nil || !strings.Contains(err.Error(), says) {
			t.Errorf("%s: refused saying %q, got %v", entry, says, err)
		}
	}
}

func TestFamilyComments(t *testing.T) {
	t.Run("CMMRC-B06: The families have an example of their comments, never applied by itself", func(t *testing.T) {})
	c := withComments("ts", "python", "ocaml")
	if got := c.CommentsFor("a.tsx").Blank("<A>{/* x */}</A>"); got != "<A>         </A>" {
		t.Errorf("ts reads JSX's comment: %q", got)
	}
	if got := c.CommentsFor("a.py").Blank(`s = """a # b"""  # c`); got != `s = """a # b"""     ` {
		t.Errorf("python reads its triple-quoted strings: %q", got)
	}
	if c.CommentsFor("src/p.ml").LineMarker() != "(*" || c.CommentsFor("src/dune").LineMarker() != ";" {
		t.Error("ocaml reads its blocks and dune's files")
	}
	if len(AllComments()) < len(FamilyComments("ts")) || FamilyComments("cobol") != nil {
		t.Error("every family's together, and an unknown family has none")
	}
	if (&Config{Dialect: &Dialect{Family: "go"}}).CommentsFor("main.go") != nil {
		t.Error("a project that declared none has no syntax for its code")
	}
}

func TestComments_blankKeepsLines(t *testing.T) {
	t.Run("CMMRC-I01: Blanking never moves a line", func(t *testing.T) {})
	src := "a // b\n/* c\nd */ e\n\"f // g\"\n`h\ni` # j\n(* k\nl *) m\n<!-- n\no -->\n"
	for _, cs := range AllComments() {
		cs := cs
		got := cs.Blank(src)
		if len(got) != len(src) {
			t.Fatalf("%s: length %d, want %d", cs.Files, len(got), len(src))
		}
		for i := range src {
			if (src[i] == '\n') != (got[i] == '\n') {
				t.Fatalf("%s: line break moved at %d", cs.Files, i)
			}
		}
	}
}

func TestComments_noneBuiltIn(t *testing.T) {
	t.Run("CMMRC-X01: No language's comment is built in", func(t *testing.T) {})
	var c *Config
	for _, rel := range []string{"main.go", "a.py", "p.ml"} {
		if c.CommentsFor(rel) != nil {
			t.Errorf("%s has a syntax nobody declared", rel)
		}
	}
}

func TestCommentsCovering(t *testing.T) {
	t.Run("CMMRC-B07: The examples' entries covering a set of files, and the extensions none covers", func(t *testing.T) {})
	got, uncovered := CommentsCovering([]string{"a.go", "./web/b.tsx", "run.sh", "c.zig", "Makefile", "a.spec.md"})
	var globs []string
	for _, c := range got {
		globs = append(globs, c.Files)
	}
	want := []string{"**/*.{ts,tsx,js,jsx,mjs,cjs,mts,cts}", "**/*.{yaml,yml,toml,sh,bash,zsh}", "**/*.go"}
	if strings.Join(globs, " ") != strings.Join(want, " ") {
		t.Errorf("entries in the examples' order, got %v", globs)
	}
	if strings.Join(uncovered, " ") != ".zig Makefile" {
		t.Errorf("the uncovered, sorted, got %v", uncovered)
	}
}
