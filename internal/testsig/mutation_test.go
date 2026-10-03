// @anchors
//   ref: MTINM

package testsig

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Fixture in the Mutation Testing Elements format — the same JSON that Stryker, PIT,
// Infection and mutmut emit. It covers every status that changes the count.
const fixtureMT = `{
  "schemaVersion": "1.0",
  "files": {
    "src/business-logic/pricing.ts": {
      "language": "typescript",
      "source": "…",
      "mutants": [
        {"id":"1","status":"Killed","mutatorName":"ConditionalExpression","location":{"start":{"line":10,"column":1},"end":{"line":10,"column":9}}},
        {"id":"2","status":"Killed","mutatorName":"ArithmeticOperator","location":{"start":{"line":11,"column":1},"end":{"line":11,"column":9}}},
        {"id":"3","status":"Timeout","mutatorName":"EqualityOperator","location":{"start":{"line":12,"column":1},"end":{"line":12,"column":9}}},
        {"id":"4","status":"Survived","mutatorName":"BooleanLiteral","location":{"start":{"line":42,"column":1},"end":{"line":42,"column":9}}},
        {"id":"5","status":"NoCoverage","mutatorName":"StringLiteral","location":{"start":{"line":88,"column":1},"end":{"line":88,"column":9}}},
        {"id":"6","status":"CompileError","mutatorName":"BlockStatement","location":{"start":{"line":99,"column":1},"end":{"line":99,"column":9}}}
      ]
    }
  }
}`

func TestParseMutation(t *testing.T) {
	t.Run("MTINM-B02: Killed and timed-out mutants count as killed", func(t *testing.T) {})
	t.Run("MTINM-B03: A survivor counts as survived with its line", func(t *testing.T) {})
	t.Run("MTINM-B06: A mutant that failed to compile stays out of the count and the score", func(t *testing.T) {})
	dir := t.TempDir()
	p := filepath.Join(dir, "mutation.json")
	if err := os.WriteFile(p, []byte(fixtureMT), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm, ok := rep.Files["src/business-logic/pricing.ts"]
	if !ok {
		t.Fatalf("file not found; keys = %v", keys(rep))
	}

	// Timeout counts as KILLED: the test noticed the mutation (it hung because of it).
	if fm.Killed != 3 {
		t.Errorf("killed = %d, want 3 (2 Killed + 1 Timeout)", fm.Killed)
	}
	if fm.TimedOut != 1 {
		t.Errorf("timed out = %d, want 1 counted apart", fm.TimedOut)
	}
	// A SURVIVOR is only what the test EXECUTED and did not notice. NoCoverage used to count
	// as a survivor here, and the decision changed on 08/25: "is there a test that runs this
	// line?" is the COVERAGE gate's question. Adding both, the mutation gate reported 187
	// files (measured in the reference app) it does not own — and among them the REAL zeros,
	// where the test runs and checks nothing, got lost.
	if fm.Survived != 1 {
		t.Errorf("survived = %d, want 1 — only the real Survived", fm.Survived)
	}
	if fm.NoCoverage != 1 {
		t.Errorf("noCoverage = %d, want 1 — counted apart, out of the score", fm.NoCoverage)
	}
	// CompileError stays OUT of the denominator: it says nothing about the test's quality.
	if want := 75.0; fm.Score != want {
		t.Errorf("score = %.1f, want %.1f (3 killed of 4 EXECUTED)", fm.Score, want)
	}
	// The survivors' lines are what make the score actionable — and only lines a test
	// actually reached enter. Line 88, uncovered, is the other gate's business.
	if len(fm.SurvivedAt) != 1 || fm.SurvivedAt[0] != 42 {
		t.Errorf("survived lines = %v, want [42]", fm.SurvivedAt)
	}
}

func TestParseMutationRefusesGarbage(t *testing.T) {
	refused := func(t *testing.T, content string) {
		t.Helper()
		p := filepath.Join(t.TempDir(), "x.json")
		os.WriteFile(p, []byte(content), 0o644)
		if _, err := ParseMutation(p); err == nil {
			t.Fatal("accepted invalid input silently — the worst outcome: a phantom score")
		}
	}
	t.Run("MTINM-E02: A report that is not the canonical format is refused", func(t *testing.T) {
		refused(t, "this is not json")
	})
	t.Run("MTINM-E03: A report with no file and no schema version is refused", func(t *testing.T) {
		refused(t, `{"files":{}}`)
	})
}

func keys(r *MutationReport) []string {
	var out []string
	for k := range r.Files {
		out = append(out, k)
	}
	return out
}

// Fixture of a file whose mutants were ALL ignored by the tool — the case of the constants
// table with `ignoreStatic` on. It is not a missing test: it is the instrumenter saying
// there is no experiment to run there.
const fixtureTudoIgnorado = `{
  "schemaVersion": "1.0",
  "files": {
    "src/constants/categories.ts": {
      "language": "typescript",
      "source": "…",
      "mutants": [
        {"id":"1","status":"Ignored","mutatorName":"StringLiteral","location":{"start":{"line":6,"column":1},"end":{"line":6,"column":9}}},
        {"id":"2","status":"Ignored","mutatorName":"ObjectLiteral","location":{"start":{"line":7,"column":1},"end":{"line":7,"column":9}}}
      ]
    }
  }
}`

// A file where every mutant was ignored has NO score, and keeps the count. It scored 100,
// and the map wrote `mutation_score: 100` for a file with nothing measured (measured: 15
// mutants, all ignored). The counter, not a score, is what tells the gate the tool ran.
func TestParseMutationAllIgnoredHasNoScore(t *testing.T) {
	t.Run("MTINM-I01: A file where nothing ran has no score and keeps the ignored count", func(t *testing.T) {})
	dir := t.TempDir()
	p := filepath.Join(dir, "mutation.json")
	if err := os.WriteFile(p, []byte(fixtureTudoIgnorado), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm := rep.Files["src/constants/categories.ts"]

	if fm.Score != 0 {
		t.Errorf("score = %.1f, want none (0) — nothing was measured", fm.Score)
	}
	if fm.Killed != 0 || fm.Survived != 0 {
		t.Errorf("killed/survived = %d/%d, want 0/0 — ignored is neither killed nor alive",
			fm.Killed, fm.Survived)
	}
	// The counter is what says the tool ran here and found nothing to measure.
	if fm.Ignored != 2 {
		t.Errorf("ignored = %d, want 2", fm.Ignored)
	}
}

// And the ignored mutant does not enter the score of a file that HAS a measure — it stays
// out of the experiment both ways. Without this, a `disable` on a line would inflate the
// unit's number.
func TestParseMutationIgnoredDoesNotInflateScore(t *testing.T) {
	t.Run("MTINM-B05: Ignored mutants are counted apart and stay out of the score", func(t *testing.T) {})
	misto := `{"schemaVersion":"1.0","files":{"a.ts":{"mutants":[
	  {"id":"1","status":"Killed","location":{"start":{"line":1,"column":1},"end":{"line":1,"column":2}}},
	  {"id":"2","status":"Survived","location":{"start":{"line":2,"column":1},"end":{"line":2,"column":2}}},
	  {"id":"3","status":"Ignored","location":{"start":{"line":3,"column":1},"end":{"line":3,"column":2}}},
	  {"id":"4","status":"Ignored","location":{"start":{"line":4,"column":1},"end":{"line":4,"column":2}}}
	]}}}`
	dir := t.TempDir()
	p := filepath.Join(dir, "m.json")
	if err := os.WriteFile(p, []byte(misto), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm := rep.Files["a.ts"]
	if fm.Score != 50 {
		t.Errorf("score = %.1f, want 50 (1 killed of 2 that ran) — the 2 ignored "+
			"cannot enter the denominator NOR the numerator", fm.Score)
	}
	if fm.Ignored != 2 {
		t.Errorf("ignored = %d, want 2", fm.Ignored)
	}
}

// A file NO test executes: every mutant is NoCoverage. The mutation gate does not own
// that — the coverage gate charges "no test runs this line" —, so there is no score and the
// number is recorded apart.
//
// MEASURED in the reference app on 08/25, and it is the common case, not the exception:
// `repositories/transactions.ts` gives 175 NoCoverage and zero survivors, because its proof
// is INTEGRATION, which the mutation config excludes on purpose (a test talking to a real
// service cannot tell "I broke the rule" from "the network wobbled").
func TestParseMutationUncoveredIsNotSurvivor(t *testing.T) {
	t.Run("MTINM-B07: A file where no mutant ran has no score", func(t *testing.T) {})
	uncovered := `{"schemaVersion":"1.0","files":{"repo.ts":{"mutants":[
	  {"id":"1","status":"NoCoverage","location":{"start":{"line":1,"column":1},"end":{"line":1,"column":2}}},
	  {"id":"2","status":"NoCoverage","location":{"start":{"line":2,"column":1},"end":{"line":2,"column":2}}},
	  {"id":"3","status":"NoCoverage","location":{"start":{"line":3,"column":1},"end":{"line":3,"column":2}}}
	]}}}`
	dir := t.TempDir()
	p := filepath.Join(dir, "m.json")
	if err := os.WriteFile(p, []byte(uncovered), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm := rep.Files["repo.ts"]

	if fm.Survived != 0 {
		t.Errorf("survived = %d, want 0 — uncovered is not a survivor", fm.Survived)
	}
	if fm.NoCoverage != 3 {
		t.Errorf("noCoverage = %d, want 3", fm.NoCoverage)
	}
	if fm.Score != 0 {
		t.Errorf("score = %.1f, want none (0) — no mutant executed, nothing was measured",
			fm.Score)
	}
}

// And the contrast that gives the separation its meaning: when part of the file IS
// executed, the mutants that survive THERE are still this gate's finding. The real case was
// `models/holidays.ts`: 32 uncovered and 15 real survivors — it got lost among 187 zeros,
// and the number now says "the test runs and does not check", which is actionable.
func TestParseMutationPartialCountsOnlyExecuted(t *testing.T) {
	t.Run("MTINM-B04: Uncovered mutants are counted apart and stay out of the score", func(t *testing.T) {})
	partial := `{"schemaVersion":"1.0","files":{"h.ts":{"mutants":[
	  {"id":"1","status":"NoCoverage","location":{"start":{"line":1,"column":1},"end":{"line":1,"column":2}}},
	  {"id":"2","status":"NoCoverage","location":{"start":{"line":2,"column":1},"end":{"line":2,"column":2}}},
	  {"id":"3","status":"Survived","location":{"start":{"line":3,"column":1},"end":{"line":3,"column":2}}},
	  {"id":"4","status":"Killed","location":{"start":{"line":4,"column":1},"end":{"line":4,"column":2}}}
	]}}}`
	dir := t.TempDir()
	p := filepath.Join(dir, "m.json")
	if err := os.WriteFile(p, []byte(partial), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil {
		t.Fatal(err)
	}
	fm := rep.Files["h.ts"]
	if !reflect.DeepEqual(fm.NoCoverageAt, []int{1, 2}) {
		t.Errorf("the uncovered mutants' lines are recorded, got %v", fm.NoCoverageAt)
	}

	if fm.Score != 50 {
		t.Errorf("score = %.1f, want 50 (1 killed of 2 EXECUTED) — the 2 uncovered "+
			"cannot enter the denominator", fm.Score)
	}
	if fm.NoCoverage != 2 || fm.Survived != 1 || fm.Killed != 1 {
		t.Errorf("noCov/surv/killed = %d/%d/%d, want 2/1/1",
			fm.NoCoverage, fm.Survived, fm.Killed)
	}
}

// The default (empty) format stays the canonical one: no existing project changes
// behaviour because of the format choice.
func TestParseMutationFormat_EmptyIsMTE(t *testing.T) {
	t.Run("MTINM-B01: The canonical format is read under each of its names", func(t *testing.T) {})
	p := filepath.Join(t.TempDir(), "mte.json")
	if err := os.WriteFile(p, []byte(fixtureMT), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, format := range []string{"", "mutation-testing-elements", " MTE ", "Stryker"} {
		rep, err := ParseMutationFormat(p, format)
		if err != nil {
			t.Fatalf("format %q: %v", format, err)
		}
		if _, ok := rep.Files["src/business-logic/pricing.ts"]; !ok {
			t.Fatalf("format %q should read MTE; got %v", format, rep.Files)
		}
	}
}

func TestParseMutationFormat_Unknown(t *testing.T) {
	t.Run("MTINM-E01: An unknown format is refused naming the accepted ones", func(t *testing.T) {})
	p := filepath.Join(t.TempDir(), "x.json")
	if err := os.WriteFile(p, []byte(fixtureGremlins), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ParseMutationFormat(p, "stryker4s-xml")
	if err == nil {
		t.Fatal("an unknown format should fail")
	}
	// The message must LIST the accepted ones — otherwise the operator guesses.
	if !strings.Contains(err.Error(), "gremlins") {
		t.Errorf("the message should list the accepted formats; got: %v", err)
	}
}

// Reading a gremlins report as the canonical format is refused (the other half of MTINM-E02).
func TestParseMutationFormat_GremlinsReadAsMTE(t *testing.T) {
	p := filepath.Join(t.TempDir(), "grem.json")
	if err := os.WriteFile(p, []byte(fixtureGremlins), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseMutationFormat(p, ""); err == nil {
		t.Error("reading gremlins as MTE should fail")
	}
}

func TestParseMutationThresholdsPathsAndMissing(t *testing.T) {
	read := func(t *testing.T, body string) *MutationReport {
		t.Helper()
		p := filepath.Join(t.TempDir(), "m.json")
		os.WriteFile(p, []byte(body), 0o644)
		rep, err := ParseMutation(p)
		if err != nil {
			t.Fatal(err)
		}
		return rep
	}
	one := `{"status":"Killed","location":{"start":{"line":1}}}`
	t.Run("MTINM-B08: The thresholds are read from the report", func(t *testing.T) {
		rep := read(t, `{"thresholds":{"low":60,"high":80},"files":{"a.ts":{"mutants":[`+one+`]}}}`)
		if rep.Low != 60 || rep.High != 80 {
			t.Errorf("thresholds = %v/%v, want 60/80", rep.Low, rep.High)
		}
	})
	t.Run("MTINM-X01: Thresholds absent from the report stay zero", func(t *testing.T) {
		rep := read(t, `{"files":{"a.ts":{"mutants":[`+one+`]}}}`)
		if rep.Low != 0 || rep.High != 0 {
			t.Errorf("thresholds = %v/%v, want 0/0", rep.Low, rep.High)
		}
	})
	t.Run("MTINM-B09: File paths are normalized to the map's form", func(t *testing.T) {
		rep := read(t, `{"files":{"./lib/a.ts":{"mutants":[`+one+`]},"/home/ci/app/src/b.ts":{"mutants":[`+one+`]}}}`)
		if _, ok := rep.Files["lib/a.ts"]; !ok {
			t.Errorf("want lib/a.ts, got %v", keys(rep))
		}
		if _, ok := rep.Files["src/b.ts"]; !ok {
			t.Errorf("want src/b.ts, got %v", keys(rep))
		}
	})
	t.Run("MTINM-E04: A missing report returns the read error", func(t *testing.T) {
		if _, err := ParseMutation(filepath.Join(t.TempDir(), "none.json")); !os.IsNotExist(err) {
			t.Errorf("want the not-exist error, got %v", err)
		}
	})
}

func TestParseMutation_aReportWithNoFileHadNothingToMutate(t *testing.T) {
	t.Run("MTINM-B10: A report in the format with no file had nothing to mutate", func(t *testing.T) {})
	p := filepath.Join(t.TempDir(), "mutation.json")
	if err := os.WriteFile(p, []byte(`{"schemaVersion":"1.0","thresholds":{"high":80,"low":60},"files":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := ParseMutation(p)
	if err != nil || len(rep.Files) != 0 {
		t.Fatalf("a run with nothing to mutate is an empty report, not an error: %v %v", rep, err)
	}
}
