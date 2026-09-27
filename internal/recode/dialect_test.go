package recode

import "testing"

func TestTestIDPrefix(t *testing.T) {
	t.Run("RCDLR-B01: The lower convention derives the lower-case prefix", func(t *testing.T) {})
	if got := TestIDPrefix("TCDTX", "lower"); got != "tcdtx" {
		t.Errorf("lower: %q", got)
	}
	if got := TestIDPrefix("TCDTX", ""); got != "" {
		t.Errorf("an empty convention should give empty, got %q", got)
	}
}

func TestRewriteTestIDs(t *testing.T) {
	t.Run("RCDLR-B02: Only quoted, hyphenated testID tokens are rewritten", func(t *testing.T) {})
	in := `testID="tcdt-amount" and the other <View testID='tcdt-row-1' /> but 'tcdtsomething' not`
	got, n := RewriteTestIDs(in, "tcdt", "tctx")
	if n != 2 {
		t.Fatalf("want 2 rewrites, got %d: %q", n, got)
	}
	if !contains(got, `"tctx-amount"`) || !contains(got, `'tctx-row-1'`) {
		t.Errorf("the testIDs were not rewritten: %q", got)
	}
	if !contains(got, "tcdtsomething") {
		t.Errorf("touched a token that is NOT a testID (no hyphen): %q", got)
	}
}

func TestRewriteTestIDs_noopWhenEmpty(t *testing.T) {
	t.Run("RCDLR-B03: An empty or unchanged prefix rewrites nothing", func(t *testing.T) {})
	in := `testID="tcdt-x"`
	if _, n := RewriteTestIDs(in, "", "tctx"); n != 0 {
		t.Errorf("an empty prefix should be a no-op")
	}
	if got, n := RewriteTestIDs(in, "tcdt", "tcdt"); n != 0 || got != in {
		t.Errorf("the same prefix should be a no-op, got %q (n=%d)", got, n)
	}
}

func TestCountTestIDPrefix(t *testing.T) {
	t.Run("RCDLR-B04: The occurrences of a prefix are counted", func(t *testing.T) {})
	in := `testID="txdt-a" 'txdt-b' and "outro-c"`
	if got := CountTestIDPrefix(in, "txdt"); got != 2 {
		t.Errorf("want 2, got %d", got)
	}
	if got := CountTestIDPrefix(in, "tcdt"); got != 0 {
		t.Errorf("an absent prefix should give 0, got %d", got)
	}
	if got := CountTestIDPrefix(in+` '-a'`, ""); got != 0 {
		t.Errorf("an empty prefix should give 0, got %d", got)
	}
}

func TestCountAnyTestID(t *testing.T) {
	t.Run("RCDLR-B05: The testID attributes are counted whatever their prefix", func(t *testing.T) {})
	in := `<View testID="zzzz-root" /> <Text testID = 'a-b' /> // the testID "prop" has no value`
	if got := CountAnyTestID(in); got != 2 {
		t.Errorf("want 2 testID attributes, got %d", got)
	}
}

func TestFileMatchesCode(t *testing.T) {
	t.Run("RCDLR-B06: A path matches the file patterns with the exact code", func(t *testing.T) {})
	pats := []string{"**/{{code}}-*.yaml", "**/{{code}}-suite.yaml", "**/*.{{code}}-*.png"}
	ok := []string{
		"apps/.maestro/screens/x/TCDTX-B01.yaml",
		"apps/.maestro/suites/TCDTX-suite.yaml",
		"apps/screens/Foo.TCDTX-VR-loaded.png",
	}
	for _, p := range ok {
		if !FileMatchesCode(p, "TCDTX", pats) {
			t.Errorf("%q should match", p)
		}
	}
	no := []string{
		"apps/.maestro/screens/x/MNDTX-B01.yaml", // another code
		"apps/components/ActionLink.tsx",         // no code in the name
		"apps/.maestro/screens/x/tcdtx-B01.yaml", // the code in another case
	}
	for _, p := range no {
		if FileMatchesCode(p, "TCDTX", pats) {
			t.Errorf("%q should NOT match", p)
		}
	}
}

func TestRenameFilePath(t *testing.T) {
	t.Run("RCDLR-B07: The code is renamed in the file name only", func(t *testing.T) {})
	t.Run("RCDLR-X01: A longer code in a file name is not renamed", func(t *testing.T) {})
	cases := map[string]string{
		"a/b/TCDTX-B01.yaml":        "a/b/TCTXX-B01.yaml",
		"a/b/TCDTX-suite.yaml":      "a/b/TCTXX-suite.yaml",
		"a/Foo.TCDTX-VR-loaded.png": "a/Foo.TCTXX-VR-loaded.png",
		"TCDTX/TCDTX-B01.yaml":      "TCDTX/TCTXX-B01.yaml", // the folder never changes
		"a/b/ActionLink.tsx":        "a/b/ActionLink.tsx",   // no code → unchanged
		// A LONGER code is not touched: `TCDTXX` holds `TCDTX` as a prefix, and a naive replace
		// would mangle it into `TCTXXX`. What the case proves is the code's BOUNDARY, so the
		// neighbour has to be longer than the target.
		"a/b/TCDTXX-B01.yaml": "a/b/TCDTXX-B01.yaml",
	}
	for in, want := range cases {
		if got := RenameFilePath(in, "TCDTX", "TCTXX"); got != want {
			t.Errorf("RenameFilePath(%q) = %q, want %q", in, got, want)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
