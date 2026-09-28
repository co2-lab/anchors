package testsig

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Fixture in the gremlins format — the exact shape `gremlins unleash --output` writes (see
// internal/report/internal/structure.go in gremlins). It covers each status of the enum
// (internal/mutator/mutator.go): KILLED, LIVED, NOT COVERED, NOT VIABLE, TIMED OUT.
const fixtureGremlins = `{
  "go_module": "example-service",
  "test_efficacy": 40.0,
  "mutations_coverage": 75.0,
  "mutants_total": 6,
  "mutants_killed": 2,
  "mutants_lived": 1,
  "mutants_not_viable": 1,
  "mutants_not_covered": 1,
  "elapsed_time": 12.5,
  "files": [
    {
      "file_name": "src/domain/bankaccount.go",
      "mutations": [
        {"type":"CONDITIONALS_NEGATION","status":"KILLED","line":10,"column":4},
        {"type":"ARITHMETIC_BASE","status":"TIMED OUT","line":11,"column":8},
        {"type":"INVERT_NEGATIVES","status":"LIVED","line":42,"column":2},
        {"type":"CONDITIONALS_BOUNDARY","status":"NOT COVERED","line":88,"column":6},
        {"type":"INCREMENT_DECREMENT","status":"NOT VIABLE","line":99,"column":1},
        {"type":"INVERT_LOGICAL","status":"RUNNABLE","line":120,"column":3}
      ]
    }
  ]
}`

func TestParseMutationFormat_Gremlins(t *testing.T) {
	t.Run("GRING-B01: The listed files are read into one result per file", func(t *testing.T) {})
	t.Run("GRING-B02: Killed and timed-out mutants count as killed", func(t *testing.T) {})
	t.Run("GRING-B06: The thresholds stay zero", func(t *testing.T) {})
	t.Run("GRING-X01: The report's own efficacy figure is not used", func(t *testing.T) {})
	dir := t.TempDir()
	p := filepath.Join(dir, "gremlins.json")
	if err := os.WriteFile(p, []byte(fixtureGremlins), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutationFormat(p, "gremlins")
	if err != nil {
		t.Fatal(err)
	}
	fm, ok := rep.Files["src/domain/bankaccount.go"]
	if !ok {
		t.Fatalf("the file did not match; got %v", rep.Files)
	}
	// KILLED + TIMED OUT = 2 killed (a timeout is death by hanging).
	if fm.TimedOut != 1 {
		t.Errorf("TimedOut = %d, want 1 counted apart", fm.TimedOut)
	}
	if fm.Killed != 2 {
		t.Errorf("Killed = %d, want 2", fm.Killed)
	}
	// LIVED is the one survivor; NOT COVERED is counted apart, as MTE does (GRING-B08).
	if fm.Survived != 1 || fm.NoCoverage != 1 {
		t.Errorf("Survived = %d, NoCoverage = %d, want 1 and 1", fm.Survived, fm.NoCoverage)
	}
	// NOT COVERED, NOT VIABLE and RUNNABLE stay out of the denominator: 2/(2+1), not the
	// report's own 40.0.
	if want := float64(2) / 3 * 100; fm.Score != want {
		t.Errorf("Score = %v, want %v", fm.Score, want)
	}
	// The survivors' lines are what the author needs to act.
	if len(fm.SurvivedAt) != 1 || fm.SurvivedAt[0] != 42 {
		t.Errorf("SurvivedAt = %v, want [42]", fm.SurvivedAt)
	}
	// gremlins writes no threshold in the report — Low/High stay absent (zero) and the
	// engine falls back to the default, instead of inheriting a ruler invented here.
	if rep.Low != 0 || rep.High != 0 {
		t.Errorf("thresholds = %v/%v, want 0/0 (gremlins does not emit them)", rep.Low, rep.High)
	}
}

// Swapping one format for the other must FAIL with a message that teaches — it is the
// likely mistake of whoever declares the wrong `format:` in anchors.yaml, and silence here
// would become "0 files" without explanation.
func TestParseMutationFormat_MTEReadAsGremlins(t *testing.T) {
	t.Run("GRING-E01: A canonical-format report under the gremlins format is refused naming format", func(t *testing.T) {})
	mte := filepath.Join(t.TempDir(), "mte.json")
	if err := os.WriteFile(mte, []byte(fixtureMT), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseMutationFormat(mte, "gremlins"); err == nil {
		t.Error("reading MTE as gremlins should fail")
	} else if !strings.Contains(err.Error(), "format") {
		t.Errorf("the message should point at `format:`; got: %v", err)
	}
}

func readGremlins(t *testing.T, body string) (*MutationReport, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "g.json")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return ParseMutationFormat(p, "gremlins")
}

func TestGremlinsStatusesAndPaths(t *testing.T) {
	mustRead := func(t *testing.T, body string) *MutationReport {
		t.Helper()
		rep, err := readGremlins(t, body)
		if err != nil {
			t.Fatal(err)
		}
		return rep
	}
	t.Run("GRING-B03: A mutant that lived counts as survived with its line", func(t *testing.T) {
		fm := mustRead(t, `{"files":[{"file_name":"a.go","mutations":[{"status":"LIVED","line":42},{"status":"KILLED","line":7}]}]}`).Files["a.go"]
		if fm.Survived != 1 || len(fm.SurvivedAt) != 1 || fm.SurvivedAt[0] != 42 {
			t.Errorf("want 1 survived at 42, got %+v", fm)
		}
	})
	t.Run("GRING-B04: Not viable, runnable, skipped and unknown statuses stay out of the score", func(t *testing.T) {
		fm := mustRead(t, `{"files":[{"file_name":"a.go","mutations":[
			{"status":"KILLED","line":1},{"status":"LIVED","line":2},
			{"status":"NOT VIABLE","line":3},{"status":"RUNNABLE","line":4},
			{"status":"SKIPPED","line":5},{"status":"SOMETHING NEW","line":6}]}]}`).Files["a.go"]
		if fm.Score != 50 || fm.Killed != 1 || fm.Survived != 1 {
			t.Errorf("want 50 with 1 killed and 1 survived, got %+v", fm)
		}
	})
	t.Run("GRING-B05: The status is matched ignoring case and spaces", func(t *testing.T) {
		fm := mustRead(t, `{"files":[{"file_name":"a.go","mutations":[
			{"status":"killed","line":1},{"status":" TimedOut ","line":2},{"status":"Timed out","line":4},{"status":"Lived","line":3}]}]}`).Files["a.go"]
		if fm.Killed != 3 || fm.Survived != 1 {
			t.Errorf("want 3 killed and 1 survived, got %+v", fm)
		}
	})
	t.Run("GRING-B07: File paths are normalized as in the canonical reading", func(t *testing.T) {
		rep := mustRead(t, `{"files":[{"file_name":"./cmd/a.go","mutations":[{"status":"KILLED","line":1}]},
			{"file_name":"/home/ci/svc/src/b.go","mutations":[{"status":"KILLED","line":1}]}]}`)
		for _, k := range []string{"cmd/a.go", "src/b.go"} {
			if _, ok := rep.Files[k]; !ok {
				t.Errorf("want %s, got %v", k, keys(rep))
			}
		}
	})
	t.Run("GRING-I01: The score is killed over killed plus survived", func(t *testing.T) {
		rep := mustRead(t, `{"files":[
			{"file_name":"a.go","mutations":[{"status":"KILLED"},{"status":"KILLED"},{"status":"KILLED"},{"status":"LIVED"}]},
			{"file_name":"b.go","mutations":[{"status":"KILLED"},{"status":"LIVED"},{"status":"LIVED"},{"status":"LIVED"}]}]}`)
		for k, want := range map[string]float64{"a.go": 75, "b.go": 25} {
			fm := rep.Files[k]
			if fm.Score != want || fm.Score != float64(fm.Killed)/float64(fm.Killed+fm.Survived)*100 {
				t.Errorf("%s: score %v, want %v (%+v)", k, fm.Score, want, fm)
			}
		}
	})
	t.Run("GRING-E02: A report without files is refused", func(t *testing.T) {
		if _, err := readGremlins(t, `{"go_module":"m","files":[]}`); err == nil {
			t.Error("an empty file list must be refused")
		}
	})
}

// The two readings disagreed on the same facts: a NOT COVERED mutant was a survivor here and
// NoCoverage (out of the score) in the canonical reading, and a file where no mutant ran
// scored 0 here and 100 there. The gate compared scores from different arithmetics.
func TestGremlinsAgreesWithCanonical(t *testing.T) {
	t.Run("GRING-B08: A mutant no test covered is counted apart and does not enter the score", func(t *testing.T) {
		rep, err := readGremlins(t, `{"files":[{"file_name":"a.go","mutations":[
			{"status":"KILLED","line":1},{"status":"NOT COVERED","line":2}]}]}`)
		if err != nil {
			t.Fatal(err)
		}
		fm := rep.Files["a.go"]
		if fm.NoCoverage != 1 || fm.Survived != 0 || len(fm.SurvivedAt) != 0 || fm.Score != 100 || !reflect.DeepEqual(fm.NoCoverageAt, []int{2}) {
			t.Errorf("want 1 no-coverage, no survivor, score 100; got %+v", fm)
		}
	})
	t.Run("GRING-B09: A file where no mutant ran scores 100", func(t *testing.T) {
		rep, err := readGremlins(t, `{"files":[{"file_name":"a.go","mutations":[
			{"status":"NOT VIABLE","line":1},{"status":"SKIPPED","line":2}]},
			{"file_name":"b.go","mutations":[{"status":"NOT COVERED","line":1}]}]}`)
		if err != nil {
			t.Fatal(err)
		}
		if a, b := rep.Files["a.go"], rep.Files["b.go"]; a.Score != 100 || b.Score != 100 || b.NoCoverage != 1 {
			t.Errorf("want both 100 with b's no-coverage kept; got a=%+v b=%+v", a, b)
		}
	})
}
