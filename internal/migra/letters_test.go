// @anchors
//   code: LTTSL
//   ref: MGLTR

package migra

import (
	"reflect"
	"strings"
	"testing"
)

var format5Kinds = map[string]string{"PLANA": "plan", "FLOWA": "flow", "ACTNA": "action"}

func TestLetterRenamesBetween(t *testing.T) {
	t.Run("MGLTR-B01: The renames of the steps crossed", func(t *testing.T) {})
	want := []LetterRename{{"plan", "F", "W"}, {"flow", "P", "T"}, {"flow", "R", "O"}, {"action", "R", "O"}}
	if got := LetterRenamesBetween(4, 5); !reflect.DeepEqual(got, want) {
		t.Errorf("from 4 to 5: want %v, got %v", want, got)
	}
	if got := LetterRenamesBetween(1, 4); len(got) != 0 {
		t.Errorf("no step before 5 renames letters, got %v", got)
	}
	if got := LetterRenamesBetween(5, 5); len(got) != 0 {
		t.Errorf("already at 5 there is nothing to cross, got %v", got)
	}
}

func TestRewriteCodeLetters_theThreeKinds(t *testing.T) {
	t.Run("MGLTR-B02: A phase, a step and a result get their new letters", func(t *testing.T) {})
	in := "needs: PLANA-F02\n### FLOWA-P04 — confront\n- `ACTNA-R01` PROMOTABLE → `FLOWA-P06`\n- `FLOWA-R01` a flow fitted as a piece\n"
	got, _ := RewriteCodeLetters(in, format5Kinds, LetterRenamesBetween(4, 5))
	want := "needs: PLANA-W02\n### FLOWA-T04 — confront\n- `ACTNA-O01` PROMOTABLE → `FLOWA-T06`\n- `FLOWA-O01` a flow fitted as a piece\n"
	if got != want {
		t.Fatalf("want\n%q\ngot\n%q", want, got)
	}
}

func TestRewriteCodeLetters_everyOtherCodeStays(t *testing.T) {
	t.Run("MGLTR-B03: Every other code stays", func(t *testing.T) {})
	in := "LOGIN-R01 permission, PLANA-R0004 revision, OTHER-F01 unknown, PLANA-B01 other letter, FLOWA-F01, PLANA-F012"
	got, counts := RewriteCodeLetters(in, format5Kinds, LetterRenamesBetween(4, 5))
	if got != in || len(counts) != 0 {
		t.Fatalf("nothing here is a renamed code, got %q %v", got, counts)
	}
}

func TestRewriteCodeLetters_reportsWhatItReplaced(t *testing.T) {
	t.Run("MGLTR-B04: The rewrite says what it replaced", func(t *testing.T) {})
	_, counts := RewriteCodeLetters("PLANA-F01 and PLANA-F01; FLOWA-P01", format5Kinds, LetterRenamesBetween(4, 5))
	keys := SortedCounts(counts)
	if strings.Join(keys, "|") != "FLOWA-P01 → FLOWA-T01|PLANA-F01 → PLANA-W01" ||
		counts["PLANA-F01 → PLANA-W01"] != 2 || counts["FLOWA-P01 → FLOWA-T01"] != 1 {
		t.Fatalf("two and one, in order, got %v %v", keys, counts)
	}
}

func TestRewriteCodeLetters_isIdempotent(t *testing.T) {
	t.Run("MGLTR-B05: Rewriting twice changes nothing more", func(t *testing.T) {})
	once, _ := RewriteCodeLetters("PLANA-F01 FLOWA-P01 ACTNA-R01", format5Kinds, LetterRenamesBetween(4, 5))
	twice, counts := RewriteCodeLetters(once, format5Kinds, LetterRenamesBetween(4, 5))
	if twice != once || len(counts) != 0 {
		t.Fatalf("a rewritten text stays, got %q %v", twice, counts)
	}
}
