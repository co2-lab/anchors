// @anchors
//   code: STTSG
//   ref: MSCMG

package migra

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The CHAIN of steps is what makes upgrading cheap: adding a format is adding a
// `format_N.go` file, and nothing else needs to know it exists.
//
// What these tests protect is the property that makes it work — that there is no HOLE. A
// project on format 2 going to 4 without the step that produces 3 would get the new number
// with the old content, and the file would start lying about its own format.

func TestStepsFrom_returnsTheStepsInOrder(t *testing.T) {
	t.Run("MSCMG-B02: The steps between two formats come in ascending order", func(t *testing.T) {})
	withSteps(t, Step{To: 3, Why: "third"}, Step{To: 2, Why: "second"})

	got, err := StepsFrom(1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].To != 2 || got[1].To != 3 {
		t.Fatalf("expected 1→2→3 in that order; got %v", got)
	}
}

// Registering out of order cannot produce a migration out of order: the `format_N.go` files
// load in the order Go decides, not the numeric one.
func TestRegister_sortsWhateverTheRegistrationOrder(t *testing.T) {
	t.Run("MSCMG-B01: Steps registered out of order are kept sorted", func(t *testing.T) {})
	withSteps(t, Step{To: 4}, Step{To: 2}, Step{To: 3})

	for i, s := range steps {
		if s.To != i+2 {
			t.Fatalf("the steps should be sorted by format; got %v", steps)
		}
	}
}

// The HOLE is an error, not silence. Without the step that produces 3, a project on 2 cannot
// go to 4 pretending 3 never existed — and the message must say WHICH step is missing, or
// whoever reads does not know what to write.
func TestStepsFrom_aHoleInTheChainIsAnError(t *testing.T) {
	t.Run("MSCMG-E01: A hole in the chain is an error naming the missing format", func(t *testing.T) {})
	withSteps(t, Step{To: 2}, Step{To: 4}) // 3 is missing

	_, err := StepsFrom(1, 4)
	if err == nil {
		t.Fatal("a hole in the chain must be an error")
	}
	if !strings.Contains(err.Error(), "format 3") {
		t.Errorf("the message should name the missing step; got: %v", err)
	}
}

func TestStepsFrom_aTargetBeyondTheLastStepIsAnError(t *testing.T) {
	t.Run("MSCMG-E02: A target beyond the last step is an error naming the first missing format", func(t *testing.T) {})
	withSteps(t, Step{To: 2}, Step{To: 3})

	_, err := StepsFrom(1, 5)
	if err == nil || !strings.Contains(err.Error(), "format 4") {
		t.Errorf("expected an error naming format 4, got %v", err)
	}
}

// Nothing to do is a valid answer, not an error: a project already on the current format
// calls the migrator on every command, and answering an error would make `check` fail for
// being up to date.
func TestStepsFrom_alreadyAtTheTargetHasNoStep(t *testing.T) {
	t.Run("MSCMG-B03: A file already at the target needs no step", func(t *testing.T) {})
	got, err := StepsFrom(2, 2)
	if err != nil {
		t.Fatalf("being up to date is not an error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("nothing to migrate should return zero steps; got %d", len(got))
	}
}

func TestStepsFrom_onlyTheStepsInsideTheInterval(t *testing.T) {
	t.Run("MSCMG-X01: Only the steps inside the interval are returned", func(t *testing.T) {})
	withSteps(t, Step{To: 2}, Step{To: 3}, Step{To: 4})

	got, err := StepsFrom(2, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].To != 3 {
		t.Errorf("expected only the step producing 3; got %v", got)
	}
}

func TestStep_letterRenamesAreNotAppliedToTheYAML(t *testing.T) {
	t.Run("MSCMG-B05: A step may rename code letters by kind", func(t *testing.T) {})
	s := Step{To: 5, RenameLetters: []LetterRename{{Kind: "plan", From: "F", To: "W"}}}
	if s.RenameLetters[0] != (LetterRename{"plan", "F", "W"}) {
		t.Fatalf("the step carries the rename, got %+v", s.RenameLetters)
	}
	p := filepath.Join(t.TempDir(), "anchors.graph.yaml")
	if err := os.WriteFile(p, []byte("version: 4\nnodes:\n    - id: plans/p.md\n      code: PLANA\n      codes: [PLANA-F01]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateFile(p, 5, false); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(p); !strings.Contains(string(b), "version: 5") || !strings.Contains(string(b), "PLANA-F01") {
		t.Errorf("the YAML step only raises the version; the letters are the command's, got:\n%s", b)
	}
}

func TestRenamedKey_answersForAnyStepAndFile(t *testing.T) {
	t.Run("MSCMG-B04: A key some step renames is reported as renamed", func(t *testing.T) {})
	for key, want := range map[string]bool{
		"julgamentos":  true,  // format 2, the map
		"rule_marking": true,  // format 4, the configuration
		"no_such_key":  false, // a typo
		"judgments":    false, // a NEW name is not an old key
	} {
		if got := RenamedKey(key); got != want {
			t.Errorf("RenamedKey(%q) = %v, want %v", key, got, want)
		}
	}
}

func TestAllSteps_returnsACopy(t *testing.T) {
	t.Run("MSCMG-I01: The listed steps are a copy of the registry", func(t *testing.T) {})
	listed := AllSteps()
	if len(listed) == 0 {
		t.Fatal("the real registry should have steps")
	}
	why := listed[0].Why
	listed[0].Why = "changed by a caller"
	if AllSteps()[0].Why != why {
		t.Error("changing the listed steps changed the registry")
	}
}
