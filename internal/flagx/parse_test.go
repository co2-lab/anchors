// @anchors
//   ref: FLPRF

package flagx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"
)

const example = "<!-- @anchors\n  code: CHKUT\n-->\n" + `# Flag: new-checkout

## Scenarios

| Scenario | When the value | Then |
| --- | --- | --- |
| ` + "`CHKUT-G01`" + ` | ` + "`= \"off\"`" + ` | the old checkout answers |
| ` + "`CHKUT-G02`" + ` | ` + "`= \"on\"`" + ` | the new checkout answers |
| ` + "`CHKUT-G03`" + ` | ` + "`>= 50`" + ` | the rollout decides per user |
| ` + "`CHKUT-G04`" + ` | absent | the header's default holds |
`

func TestParse_readsTheScenarios(t *testing.T) {
	t.Run("FLPRF-B01: The G rows of the table are the scenarios", func(t *testing.T) {})
	t.Run("FLPRF-B02: The flag's code and name", func(t *testing.T) {})
	f := parse(example, "flags/new-checkout.flag.md")
	if len(f.Scenarios) != 4 {
		t.Fatalf("len(Scenarios) = %d, want 4", len(f.Scenarios))
	}
	if f.Code != "CHKUT" {
		t.Errorf("Code = %q, want CHKUT", f.Code)
	}
	if f.Name != "new-checkout" {
		t.Errorf("Name = %q, want new-checkout", f.Name)
	}
	if got := f.Scenarios[0]; got.Cond.Op != OpEq || got.Cond.Operand != "off" {
		t.Errorf("G01 = %q/%q, want eq/off", got.Cond.Op, got.Cond.Operand)
	}
	if got := f.Scenarios[2]; got.Cond.Op != OpGte || got.Cond.Operand != "50" {
		t.Errorf("G03 = %q/%q, want gte/50", got.Cond.Op, got.Cond.Operand)
	}
	if got := f.Scenarios[1].Then; got != "the new checkout answers" {
		t.Errorf("G02.Then = %q", got)
	}
}

// The line is recorded so the verdict can point at it.
func TestParse_recordsTheLine(t *testing.T) {
	t.Run("FLPRF-B03: Each scenario records its line", func(t *testing.T) {})
	f := parse(example, "flags/new-checkout.flag.md")
	if f.Scenarios[0].Line == 0 {
		t.Error("Line was not recorded — the verdict would have nowhere to point")
	}
	if f.Scenarios[1].Line <= f.Scenarios[0].Line {
		t.Error("the lines are not increasing")
	}
}

// The ABSENT case is the question a flag most often forgets. (Its spelling in other
// languages, `ausente`, has its own test in grammar_test.go.)
func TestAbsent(t *testing.T) {
	t.Run("FLPRF-B05: The flag says whether it declares the absent case", func(t *testing.T) {})
	f := parse(example, "x.flag.md")
	if !f.Absent() {
		t.Error("the example declares `absent` and Absent() said no")
	}

	noAbsent := "| `CHKUT-G01` | `= \"on\"` | turns on |\n"
	if parse(noAbsent, "x.flag.md").Absent() {
		t.Error("there is no absent scenario and Absent() said yes")
	}
}

// A condition the grammar refuses becomes a FINDING, it does not vanish. A line that
// vanishes is exactly the silence the gates exist to end.
func TestParse_invalidConditionBecomesAFindingNotAVanishing(t *testing.T) {
	t.Run("FLPRF-B04: A refused condition becomes a finding", func(t *testing.T) {})
	src := "| `CHKUT-G01` | when the user is a beta tester | turns on |\n"
	f := parse(src, "x.flag.md")
	if len(f.Scenarios) != 1 {
		t.Fatalf("the line vanished: len = %d, want 1", len(f.Scenarios))
	}
	if f.Scenarios[0].Err == nil {
		t.Error("the condition is prose and was accepted — the fixed grammar refused nothing")
	}
}

// Only the letter `G` is a flag scenario. A row with another letter in the same table
// does not count.
func TestParse_onlyTheLetterG(t *testing.T) {
	src := "| `CHKUT-B01` | `= \"on\"` | not a flag scenario |\n" +
		"| `CHKUT-G01` | `= \"on\"` | this one is |\n"
	f := parse(src, "x.flag.md")
	if len(f.Scenarios) != 1 {
		t.Fatalf("len = %d, want 1 — only the letter G", len(f.Scenarios))
	}
	if f.Scenarios[0].Code != "CHKUT-G01" {
		t.Errorf("got %q", f.Scenarios[0].Code)
	}
}

// A project with no `flags/` is not an error: it is a project that has not declared a
// flag yet.
func TestLoad_noFolderIsNotAnError(t *testing.T) {
	t.Run("FLPRF-B06: Flags load in a stable order, and none without a folder", func(t *testing.T) {})
	fs, err := Load(t.TempDir())
	if err != nil || fs != nil {
		t.Errorf("Load with no flags/ = (%v, %v), want (nil, nil)", fs, err)
	}
}

func TestLoad_stableOrder(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"zeta", "alpha", "middle"} {
		src := "| `" + map[string]string{"zeta": "ZETAA", "alpha": "ALPHA", "middle": "MIDDL"}[n] +
			"-G01` | `= \"on\"` | turns on |\n"
		if err := os.WriteFile(filepath.Join(dir, n+Suffix), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	fs, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(fs) != 3 {
		t.Fatalf("len = %d", len(fs))
	}
	if fs[0].Name != "alpha" || fs[1].Name != "middle" || fs[2].Name != "zeta" {
		t.Errorf("unstable order: %q, %q, %q", fs[0].Name, fs[1].Name, fs[2].Name)
	}
}

func TestByCode(t *testing.T) {
	t.Run("FLPRF-B07: Scenarios are indexed by code", func(t *testing.T) {})
	idx := ByCode([]Flag{parse(example, "x.flag.md")})
	if s, ok := idx["CHKUT-G03"]; !ok || s.Cond.Op != OpGte {
		t.Errorf("ByCode did not resolve CHKUT-G03: %v/%v", s.Cond.Op, ok)
	}
	if _, ok := idx["CHKUT-G99"]; ok {
		t.Error("ByCode resolved a code that does not exist")
	}
}

func TestLoad_aReadFailureIsAnError(t *testing.T) {
	t.Run("FLPRF-E02: A flags folder or flag file that cannot be read is an error", func(t *testing.T) {})
	testkit.SkipWithoutPOSIXPermissions(t)
	// flags/ exists but is a file: reading it fails with something other than "does not exist".
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, Dir), []byte("x"), 0o644)
	if flags, err := Load(root); err == nil {
		t.Fatalf("Load over an unreadable flags/ = %+v, nil; want the error", flags)
	}

	// A flag file that cannot be read: before, it was skipped and its scenarios vanished.
	root = t.TempDir()
	os.MkdirAll(filepath.Join(root, Dir), 0o755)
	p := filepath.Join(root, Dir, "locked"+Suffix)
	os.WriteFile(p, []byte(example), 0o000)
	if _, err := os.ReadFile(p); err == nil {
		t.Skip("this user reads a 0000 file (root?): no read failure to observe")
	}
	if flags, err := Load(root); err == nil {
		t.Fatalf("Load with an unreadable flag file = %+v, nil; want the error", flags)
	}
}
