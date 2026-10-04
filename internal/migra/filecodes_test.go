// @anchors
//   code: FLTSD
//   ref: MGFCD

package migra

import (
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
	if !CanCarryCode("a.ts", []byte("export const a = 1\n")) || CanCarryCode("a.json", []byte("{}")) ||
		CanCarryCode("a.ts", []byte("x\x00y")) || !CanCarryCode("a.feature", []byte("Feature: a\n")) {
		t.Error("only text with a comment syntax")
	}
	feat := WithHeaderCode("# language: en\n# @anchors\n#   ref: ARNA\n\nFeature: a\n", "ui/a.feature", "ARFTR")
	if feat != "# language: en\n# @anchors\n#   code: ARFTR\n#   ref: ARNA\n\nFeature: a\n" {
		t.Errorf("feature:\n%s", feat)
	}
}

func TestWithHeaderCode(t *testing.T) {
	t.Run("MGFCD-B05: The line goes below the header's opener, or in a new header", func(t *testing.T) {})
	ts := WithHeaderCode("// @anchors\n//   ref: ARNA\n\nexport const a = 1\n", "ui/a.ts", "ARSCR")
	if !strings.HasPrefix(ts, "// @anchors\n//   code: ARSCR\n//   ref: ARNA\n") {
		t.Errorf("ts:\n%s", ts)
	}
	md := WithHeaderCode("<!-- @anchors\n  layer: guide\n-->\n# Guide\n", "g.md", "GUIDE")
	if !strings.HasPrefix(md, "<!-- @anchors\n  code: GUIDE\n  layer: guide\n") {
		t.Errorf("md:\n%s", md)
	}
	sh := WithHeaderCode("#!/bin/sh\necho hi\n", "run.sh", "RUNSH")
	if !strings.HasPrefix(sh, "#!/bin/sh\n# @anchors\n#   code: RUNSH\n") {
		t.Errorf("sh:\n%s", sh)
	}
}

func TestAsRefWithOwnCode(t *testing.T) {
	t.Run("MGFCD-B06: A header carrying its unit's code turns it into a ref beside a code of its own", func(t *testing.T) {})
	in := "# @anchors\n#   code: LOGIN\n#   updated_at: 2026-01-01\n\nFeature: Login\n"
	want := "# @anchors\n#   code: LGFTR\n#   ref: LOGIN\n#   updated_at: 2026-01-01\n\nFeature: Login\n"
	if got := AsRefWithOwnCode(in, "LOGIN", "LGFTR"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	md := "<!-- @anchors\n  code: LOGIN\n-->\n# x\n"
	if got := AsRefWithOwnCode(md, "LOGIN", "LGFTR"); got != "<!-- @anchors\n  code: LGFTR\n  ref: LOGIN\n-->\n# x\n" {
		t.Errorf("markdown header: %q", got)
	}
	if got := AsRefWithOwnCode("package x\n", "LOGIN", "LGFTR"); got != "package x\n" {
		t.Errorf("no header, nothing changes: %q", got)
	}
}
