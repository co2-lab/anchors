// @anchors
//   code: FLTSD
//   ref: MGFCD

package migra

import (
	"github.com/co2-lab/anchors/internal/testkit"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/scan"
)

func TestWidenedCode(t *testing.T) {
	t.Run("MGFCD-B01: A four-character code is widened, keeping it as the prefix", func(t *testing.T) {})
	taken := map[string]bool{"ARNAR": true}
	got := WidenedCode("ARNA", "ArenaScreen", taken)
	if len(got) != 5 || !strings.HasPrefix(got, "ARNA") || got == "ARNAR" {
		t.Errorf("widened = %q", got)
	}
}

func TestFileCodeNameAndCode(t *testing.T) {
	t.Run("MGFCD-B02: The files of one unit get different code names", func(t *testing.T) {})
	t.Run("MGFCD-B03: A file's code is new and of five characters", func(t *testing.T) {})
	cases := map[string]scan.File{
		"ArenaScreen screen":  {Path: "ui/ArenaScreen.tsx", Layer: "screen"},
		"ArenaScreen feature": {Path: "ui/ArenaScreen.feature", Layer: "feature"},
		"ArenaScreen test":    {Path: "ui/ArenaScreen.test.tsx", Layer: "test"},
		"tokens code":         {Path: "theme/tokens.ts", Kind: "code"},
	}
	for want, f := range cases {
		if got := FileCodeName(f); got != want {
			t.Errorf("%s: name = %q, want %q", f.Path, got, want)
		}
	}
	f := scan.File{Path: "ui/ArenaScreen.test.tsx", Layer: "test"}
	first := FileCode(f, map[string]bool{})
	second := FileCode(f, map[string]bool{first: true})
	if len(first) != 5 || len(second) != 5 || first == second {
		t.Errorf("codes = %q, %q", first, second)
	}
}

func TestCanCarryCode(t *testing.T) {
	t.Run("MGFCD-B04: Only a text file with a comment syntax carries the line", func(t *testing.T) {})
	if !CanCarryCode([]byte("export const a = 1\n"), testkit.Comments("a.ts")) || CanCarryCode([]byte("{}"), testkit.Comments("a.json")) ||
		CanCarryCode([]byte("x\x00y"), testkit.Comments("a.ts")) || !CanCarryCode([]byte("Feature: a\n"), testkit.Comments("a.feature")) {
		t.Error("only text with a comment syntax")
	}
	feat := WithHeaderCode("# language: en\n# @anchors\n#   ref: ARNA\n\nFeature: a\n", "ARFTR", testkit.Comments("ui/a.feature"))
	if feat != "# language: en\n# @anchors\n#   code: ARFTR\n#   ref: ARNA\n\nFeature: a\n" {
		t.Errorf("feature:\n%s", feat)
	}
}

func TestWithHeaderCode(t *testing.T) {
	t.Run("MGFCD-B05: The line goes below the header's opener, or in a new header", func(t *testing.T) {})
	ts := WithHeaderCode("// @anchors\n//   ref: ARNA\n\nexport const a = 1\n", "ARSCR", testkit.Comments("ui/a.ts"))
	if !strings.HasPrefix(ts, "// @anchors\n//   code: ARSCR\n//   ref: ARNA\n") {
		t.Errorf("ts:\n%s", ts)
	}
	md := WithHeaderCode("<!-- @anchors\n  layer: guide\n-->\n# Guide\n", "GUIDE", testkit.Comments("g.md"))
	if !strings.HasPrefix(md, "<!-- @anchors\n  code: GUIDE\n  layer: guide\n") {
		t.Errorf("md:\n%s", md)
	}
	sh := WithHeaderCode("#!/bin/sh\necho hi\n", "RUNSH", testkit.Comments("run.sh"))
	if !strings.HasPrefix(sh, "#!/bin/sh\n# @anchors\n#   code: RUNSH\n") {
		t.Errorf("sh:\n%s", sh)
	}
}

func TestAsRefWithOwnCode(t *testing.T) {
	t.Run("MGFCD-B06: A header carrying its unit's code turns it into a ref beside a code of its own", func(t *testing.T) {})
	in := "# @anchors\n#   code: LOGIN\n#   updated_at: 2026-01-01\n\nFeature: Login\n"
	want := "# @anchors\n#   code: LGFTR\n#   ref: LOGIN\n#   updated_at: 2026-01-01\n\nFeature: Login\n"
	if got := AsRefWithOwnCode(in, "LOGIN", "LGFTR", testkit.AnyComment()); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	md := "<!-- @anchors\n  code: LOGIN\n-->\n# x\n"
	if got := AsRefWithOwnCode(md, "LOGIN", "LGFTR", testkit.AnyComment()); got != "<!-- @anchors\n  code: LGFTR\n  ref: LOGIN\n-->\n# x\n" {
		t.Errorf("markdown header: %q", got)
	}
	if got := AsRefWithOwnCode("package x\n", "LOGIN", "LGFTR", testkit.AnyComment()); got != "package x\n" {
		t.Errorf("no header, nothing changes: %q", got)
	}
}

func TestRenames(t *testing.T) {
	t.Run("MGFCD-B07: The renamed codes are read old to current, a code renamed twice to the last one", func(t *testing.T) {})
	root := t.TempDir()
	if got := Renames(root); len(got) != 0 {
		t.Errorf("no file, no rename: %v", got)
	}
	body := "# header\r\n2026-10-04:\r\n  ARNA: ARNAA\r\n  LOGI: LOGIN\r\n2026-11-01:\r\n  ARNAA: ARENA\r\n"
	if err := os.WriteFile(filepath.Join(root, RenamesFile), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	got := Renames(root)
	if got["ARNA"] != "ARENA" || got["ARNAA"] != "ARENA" || got["LOGI"] != "LOGIN" || len(got) != 3 {
		t.Errorf("got %v", got)
	}
}

func TestWithHeaderDate(t *testing.T) {
	t.Run("MGFCD-B08: The header's updated_at is set to the day given, and only in the header", func(t *testing.T) {})
	in := "// @anchors\n//   code: ARSCR\n//   updated_at: 2026-09-01\n\n// updated_at: 2020-01-01 in the body\n"
	want := "// @anchors\n//   code: ARSCR\n//   updated_at: 2026-10-04\n\n// updated_at: 2020-01-01 in the body\n"
	if got := WithHeaderDate(in, "2026-10-04", testkit.AnyComment()); got != want {
		t.Errorf("got %q", got)
	}
	if got := WithHeaderDate("// @anchors\n//   code: ARSCR\n", "2026-10-04", testkit.AnyComment()); got != "// @anchors\n//   code: ARSCR\n" {
		t.Errorf("a header without the date is left as it was: %q", got)
	}
}
