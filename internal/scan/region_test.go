package scan

import "testing"

func TestNestedRegionClosesInOrder(t *testing.T) {
	t.Run("SRRGS-B01: A nested region closes before the outer one", func(t *testing.T) {})
	// The case that motivated the region: MLETX-B05 lives INSIDE MLETX-A03 (it is the catch
	// of its try). Without a declared delimitation, inferring the end by indentation gave
	// A03 an interval that ran to the next marker of smaller indentation — and in JSX that
	// swallowed hundreds of lines (measured: 489).
	src := `func persist() {
  // #region [MLETX-A03]: confirming persists the entry.
  put(item)
    // #region [MLETX-B05]: a failed write warns the user.
    catch(err)
    // #endregion [MLETX-B05]
  toast()
  // #endregion [MLETX-A03]
}`
	regs, errs := Regioes(src)
	if len(errs) != 0 {
		t.Fatalf("there should be no pairing defect, got %+v", errs)
	}
	if len(regs) != 2 {
		t.Fatalf("want 2 regions, got %d", len(regs))
	}
	// the INNER one closes first, so it comes out first
	if regs[0].Code != "MLETX-B05" || regs[0].Start != 4 || regs[0].End != 6 {
		t.Errorf("wrong inner region: %+v", regs[0])
	}
	if regs[1].Code != "MLETX-A03" || regs[1].Start != 2 || regs[1].End != 8 {
		t.Errorf("wrong outer region: %+v", regs[1])
	}
}

func TestOneLineJSXRegionDoesNotSwallowTheFile(t *testing.T) {
	t.Run("SRRGS-B02: A deeply indented one-line region does not swallow the file", func(t *testing.T) {})
	// The regression the indentation heuristic produced: a ONE-line marker in JSX, deeply
	// indented, "swallowed" everything up to the next marker of smaller indentation. With a
	// declared region, the interval is what the author wrote.
	src := `<View>
            {/* #region [MLETX-A01]: toggling the type sets out/in. */}
            <TypeSegment />
            {/* #endregion [MLETX-A01] */}
</View>
` + long(500)
	regs, errs := Regioes(src)
	if len(errs) != 0 || len(regs) != 1 {
		t.Fatalf("want 1 region without defect, got %d regions %+v", len(regs), errs)
	}
	if got := regs[0].End - regs[0].Start + 1; got != 3 {
		t.Errorf("the region should cover 3 lines, covered %d", got)
	}
}

func TestRegionRevIgnoresChangesOutsideIt(t *testing.T) {
	t.Run("SRRGS-B03: The region revision ignores changes outside it", func(t *testing.T) {})
	// The central gain: editing OUTSIDE the region does not expire its stamp. It is what
	// separates "which file changed" from "which requirement changed".
	a := "// #region [CODEX-A01]: x\ndo()\n// #endregion [CODEX-A01]\nsomethingElse()"
	b := "// #region [CODEX-A01]: x\ndo()\n// #endregion [CODEX-A01]\nsomethingElseCHANGED()"
	ra, _ := Regioes(a)
	rb, _ := Regioes(b)
	if ra[0].Rev != rb[0].Rev {
		t.Errorf("a change outside the region should not change its rev: %s != %s", ra[0].Rev, rb[0].Rev)
	}
	c := "// #region [CODEX-A01]: x\ndoDIFFERENT()\n// #endregion [CODEX-A01]\nsomethingElse()"
	rc, _ := Regioes(c)
	if ra[0].Rev == rc[0].Rev {
		t.Error("a change INSIDE the region had to change the rev")
	}
}

func TestPairingDefects(t *testing.T) {
	t.Run("SRRGS-E01: A region never closed is reported", func(t *testing.T) {})
	t.Run("SRRGS-E02: A close with no open region is reported", func(t *testing.T) {})
	t.Run("SRRGS-E03: A close naming another code is reported", func(t *testing.T) {})
	cases := []struct {
		name, src, kind string
	}{
		{"never closed", "// #region [CODEX-A01]: x\ndo()", "sem-fecho"},
		{"orphan close", "do()\n// #endregion [CODEX-A01]", "fecho-orfao"},
		{"swapped close", "// #region [CODEX-A01]: x\ndo()\n// #endregion [CODEX-B02]", "fecho-trocado"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, errs := Regioes(c.src)
			if len(errs) != 1 || errs[0].Kind != c.kind {
				t.Fatalf("want 1 defect %q, got %+v", c.kind, errs)
			}
		})
	}
}

func TestSwappedCloseDoesNotCascade(t *testing.T) {
	t.Run("SRRGS-B04: A swapped close does not cascade into the following regions", func(t *testing.T) {})
	t.Run("SRRGS-I01: One pairing defect is reported once", func(t *testing.T) {})
	// A swapped close is ONE defect; if it unbalanced the stack, every following region
	// would look wrong and the report would blame correct code.
	src := `// #region [CODEX-A01]: x
do()
// #endregion [CODEX-B02]
// #region [CODEX-A02]: y
do2()
// #endregion [CODEX-A02]`
	regs, errs := Regioes(src)
	if len(errs) != 1 || errs[0].Kind != "fecho-trocado" {
		t.Fatalf("want exactly 1 fecho-trocado defect, got %+v", errs)
	}
	if len(regs) != 2 {
		t.Fatalf("both regions still have to come out, got %d", len(regs))
	}
}

func TestAnonymousCloseClosesTheTop(t *testing.T) {
	t.Run("SRRGS-B05: A close without a code closes the open region", func(t *testing.T) {})
	regs, errs := Regioes("// #region [CODEX-A01]\nx\n// #endregion")
	if len(errs) != 0 || len(regs) != 1 || regs[0].Code != "CODEX-A01" || regs[0].Start != 1 || regs[0].End != 3 {
		t.Errorf("want one region CODEX-A01 from 1 to 3 and no defect, got %+v %+v", regs, errs)
	}
}

func TestFileWithoutRegionIsNotAnError(t *testing.T) {
	t.Run("SRRGS-B06: A file without regions has no defect", func(t *testing.T) {})
	// The delimitation is OPTIONAL: where there is no region, the file's rev holds (the
	// earlier behaviour). A project need not convert anything to stay valid.
	regs, errs := Regioes("// CODEX-A01: old-style comment\ndo()")
	if len(regs) != 0 || len(errs) != 0 {
		t.Errorf("no region = no defect and no interval, got %d/%d", len(regs), len(errs))
	}
}

func long(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "line\n"
	}
	return s
}

func TestComposeRefsResolvesRelativeToTheScript(t *testing.T) {
	t.Run("SRRGS-B07: Composition resolves relative to the script's directory", func(t *testing.T) {})
	// The paths in a flow are relative to ITS directory (`../../utils/login.yaml`).
	// Resolving against the root would produce non-existent targets, and the edge would be
	// discarded in silence — the worst result, because it looks like "no dependency".
	src := `- runFlow: ../../utils/login.yaml
- runFlow:
    file: ../../utils/dismissOsDialogs.yaml
    env:
      SKIP_SETUP: 'true'
- runScript: ../../scripts/dataLoader.js`
	got := ComposeRefs(src, "apps/mobile/.maestro/suites/smoke/SS-03.yaml")
	want := []string{
		"apps/mobile/.maestro/utils/login.yaml",
		"apps/mobile/.maestro/utils/dismissOsDialogs.yaml",
	}
	if len(got) != len(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] want %q, got %q", i, want[i], got[i])
		}
	}
}

func TestComposeRefsIgnoresRunScriptAndDuplicates(t *testing.T) {
	t.Run("SRRGS-B08: Duplicates and self-references are not dependencies", func(t *testing.T) {})
	// `runScript:` points at a .js that is input DATA, not script composition. And the same
	// util cited twice is one dependency, not two; the script running itself is none.
	src := `- runScript: ../scripts/dataLoader.js
- runFlow: ../utils/login.yaml
- runFlow: ../utils/login.yaml
- runFlow: ./flow.yaml`
	got := ComposeRefs(src, "a/b/flow.yaml")
	if len(got) != 1 || got[0] != "a/utils/login.yaml" {
		t.Fatalf("want only login.yaml once, got %v", got)
	}
}

func TestComposeRefsWithoutComposition(t *testing.T) {
	t.Run("SRRGS-B09: A script with no composition has no dependency", func(t *testing.T) {})
	if got := ComposeRefs("- tapOn:\n    id: ':x'", "a/f.yaml"); got != nil {
		t.Errorf("a flow without runFlow has no composition dependency, got %v", got)
	}
}

func TestComposeRefsOnlyFromDeclaredPaths(t *testing.T) {
	t.Run("SRRGS-X01: Composition comes only from declared paths", func(t *testing.T) {})
	if got := ComposeRefs("- tapOn: \"login.yaml\"\n- assertVisible: utils/login.yaml", "a/f.yaml"); got != nil {
		t.Errorf("a script name outside runFlow/file is not composition, got %v", got)
	}
}
