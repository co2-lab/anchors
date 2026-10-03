// @anchors
//   ref: JUIJN

package testsig

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseJUnit(t *testing.T) {
	t.Run("JUIJN-B01: A report with a suites root is read case by case", func(t *testing.T) {})
	t.Run("JUIJN-B06: Only codes of passing cases are proven", func(t *testing.T) {})
	xml := `<?xml version="1.0"?>
<testsuites>
  <testsuite name="Spacer" file="src/Spacer.test.tsx">
    <testcase name="SPCRX-V01: axis"/>
    <testcase name="SPCRX-V02: horizontal"><failure message="x"/></testcase>
    <testcase name="SPCRX-X01: no content"><skipped/></testcase>
  </testsuite>
</testsuites>`
	rep, err := ParseJUnit(write(t, "j.xml", xml))
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Cases) != 3 {
		t.Fatalf("expected 3 cases, got %d", len(rep.Cases))
	}
	proven := rep.PassedCodes()
	if !proven["SPCRX-V01"] {
		t.Error("SPCRX-V01 passed, it should be proven")
	}
	if proven["SPCRX-V02"] {
		t.Error("SPCRX-V02 failed, it must NOT be proven")
	}
	if proven["SPCRX-X01"] {
		t.Error("SPCRX-X01 was skipped, it must NOT be proven")
	}
}

func TestParseJUnitSingleSuite(t *testing.T) {
	t.Run("JUIJN-B02: A report whose root is a single suite is read", func(t *testing.T) {})
	// some reporters emit <testsuite> at the root
	xml := `<testsuite name="X" file="a.test.ts"><testcase name="AB01X-S01: ok"/></testsuite>`
	rep, err := ParseJUnit(write(t, "j.xml", xml))
	if err != nil || len(rep.Cases) != 1 {
		t.Fatalf("expected 1 case, got %d err=%v", len(rep.Cases), err)
	}
}

func TestCodesInCase(t *testing.T) {
	got := CodesInCase("SPCRX-V01: vertical axis extends SPCRX-X01")
	if len(got) != 2 {
		t.Fatalf("expected 2 codes, got %v", got)
	}
}

// The project's vocabulary reaches the JUnit reading: a letter declared after the package
// loaded is recognised in a case name, and one no longer declared is not.
func TestCodesInCaseFollowsTheDeclaredLetters(t *testing.T) {
	t.Run("JUIJN-B08: Every code in a case name is extracted in the declared vocabulary", func(t *testing.T) {})
	if got := CodesInCase("SPCRX-V01: vertical axis extends SPCRX-X01"); len(got) != 2 {
		t.Fatalf("every code of the name must be extracted, got %v", got)
	}
	prev := ruleLetters
	defer SetRuleLetters(prev)

	SetRuleLetters("BZ")
	if got := CodesInCase("ABCDX-Z01: a project letter"); len(got) != 1 || got[0] != "ABCDX-Z01" {
		t.Fatalf("a declared letter must be read from the case name, got %v", got)
	}
	SetRuleLetters("B")
	if got := CodesInCase("ABCDX-Z01: no longer declared"); len(got) != 0 {
		t.Errorf("a letter no longer declared must not be read, got %v", got)
	}
}

// SeenCodes counts every case — passed, failed or skipped: a partial run measured all of
// them, and a code it saw FAIL must lose its proof, not keep it.
func TestSeenCodesCountsEveryOutcome(t *testing.T) {
	t.Run("JUIJN-B07: Every case's codes are seen whatever the outcome", func(t *testing.T) {})
	r := &ExecReport{Cases: []CaseResult{
		{Name: "MBDT-B01 passes"},
		{Name: "MBDT-B02 fails", Failed: true},
		{Name: "MBDT-B03 skipped", Skipped: true},
	}}
	seen := r.SeenCodes()
	for _, c := range []string{"MBDT-B01", "MBDT-B02", "MBDT-B03"} {
		if !seen[c] {
			t.Errorf("%s ran and must be seen: %v", c, seen)
		}
	}
	if r.PassedCodes()["MBDT-B02"] {
		t.Error("control: a failed case is not proven")
	}
}

// The shape of the report: nesting, the file a case inherits, and what fails a case.
func TestParseJUnitShape(t *testing.T) {
	xml := `<testsuites>
  <testsuite name="outer" file="a.test.ts">
    <testcase name="ABCDX-B01: inherits"/>
    <testcase name="ABCDX-B02: own file" file="b.test.ts"/>
    <testcase name="ABCDX-B03: failure"><failure message="x"/></testcase>
    <testcase name="ABCDX-B04: error"><error message="boom"/></testcase>
    <testcase name="ABCDX-B05: skipped"><skipped/></testcase>
    <testsuite name="inner"><testcase name="ABCDX-B06: nested"/></testsuite>
  </testsuite>
</testsuites>`
	rep, err := ParseJUnit(write(t, "j.xml", xml))
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]CaseResult{}
	for _, c := range rep.Cases {
		byName[c.Name[:9]] = c
	}
	t.Run("JUIJN-B03: Nested suites are flattened", func(t *testing.T) {
		if _, ok := byName["ABCDX-B06"]; !ok || len(rep.Cases) != 6 {
			t.Errorf("the nested case must be read, got %d cases", len(rep.Cases))
		}
	})
	t.Run("JUIJN-B04: A case without a file takes its suite's file", func(t *testing.T) {
		if f := byName["ABCDX-B01"].File; f != "a.test.ts" {
			t.Errorf("inherited file = %q, want a.test.ts", f)
		}
		if f := byName["ABCDX-B02"].File; f != "b.test.ts" {
			t.Errorf("own file = %q, want b.test.ts", f)
		}
	})
	t.Run("JUIJN-B05: A failure or an error marks the case failed and a skip marks it skipped", func(t *testing.T) {
		if !byName["ABCDX-B03"].Failed || !byName["ABCDX-B04"].Failed {
			t.Error("a failure and an error must both fail the case")
		}
		if c := byName["ABCDX-B05"]; !c.Skipped || c.Failed {
			t.Errorf("a skipped case must be skipped and not failed, got %+v", c)
		}
		if c := byName["ABCDX-B01"]; c.Failed || c.Skipped {
			t.Errorf("a plain case passes, got %+v", c)
		}
	})
	t.Run("JUIJN-I01: Every proven code is a seen code", func(t *testing.T) {
		seen := rep.SeenCodes()
		passed := rep.PassedCodes()
		if len(passed) == 0 || len(passed) >= len(seen) {
			t.Fatalf("control: some cases pass and some do not, got %d proven of %d seen", len(passed), len(seen))
		}
		for c := range passed {
			if !seen[c] {
				t.Errorf("%s is proven but not seen", c)
			}
		}
	})
}

func TestParseJUnitGarbageAndMissing(t *testing.T) {
	t.Run("JUIJN-B09: A file that is not a JUnit report yields no case", func(t *testing.T) {
		rep, err := ParseJUnit(write(t, "j.xml", "this is not xml"))
		if err != nil || rep == nil || len(rep.Cases) != 0 {
			t.Errorf("want an empty report and no error, got %v, %v", rep, err)
		}
	})
	t.Run("JUIJN-E01: An unreadable report returns the read error", func(t *testing.T) {
		if _, err := ParseJUnit(filepath.Join(t.TempDir(), "none.xml")); !os.IsNotExist(err) {
			t.Errorf("want the not-exist error, got %v", err)
		}
	})
}

// A code proves a rule only from the case's own name.
func TestPassedCodesReadOnlyTheCaseName(t *testing.T) {
	t.Run("JUIJN-X01: A code in the suite or class name proves nothing", func(t *testing.T) {})
	xml := `<testsuite name="ABCDX-B01 suite"><testcase name="renders" classname="ABCDX-B01 class"/><testcase name="ABCDX-B02: named"/></testsuite>`
	rep, err := ParseJUnit(write(t, "j.xml", xml))
	if err != nil {
		t.Fatal(err)
	}
	proven := rep.PassedCodes()
	if proven["ABCDX-B01"] {
		t.Error("a code in the suite or class name must not be proven")
	}
	if !proven["ABCDX-B02"] {
		t.Error("control: a code in the case name is proven")
	}
}

func TestParseJUnitCaseTime(t *testing.T) {
	t.Run("JUIJN-B10: Each case carries its run time", func(t *testing.T) {})
	xml := `<testsuites><testsuite name="s" file="a.test.ts">
  <testcase name="a" time="0.25"/><testcase name="b" time=" 3 "/><testcase name="c"/>
  <testcase name="d" time="soon"/><testcase name="e" time="-1"/>
</testsuite></testsuites>`
	rep, err := ParseJUnit(write(t, "t.xml", xml))
	if err != nil {
		t.Fatal(err)
	}
	var got []float64
	for _, c := range rep.Cases {
		got = append(got, c.Seconds)
	}
	if want := []float64{0.25, 3, 0, 0, 0}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestExecReport_variantsAreScenarios(t *testing.T) {
	rep := &ExecReport{Cases: []CaseResult{
		{Name: "RDCHX-B02#01: unit", File: "a.test.ts"},
		{Name: "RDCHX-B02#02: integration", File: "a.test.ts", Skipped: true},
	}}
	t.Run("JUIJN-B11: The proven and seen codes carry the variant", func(t *testing.T) {
		passed, seen := rep.PassedCodes(), rep.SeenCodes()
		if !passed["RDCHX-B02#01"] || passed["RDCHX-B02#02"] || !seen["RDCHX-B02#02"] {
			t.Errorf("the variant is kept, the skipped one seen and not proven, got %v %v", passed, seen)
		}
	})
	t.Run("JUIJN-I02: A scenario is proven only by its own passing case", func(t *testing.T) {
		passed := rep.PassedCodes()
		if len(passed) != 1 || passed["RDCHX-B02"] || passed["RDCHX-B02#02"] {
			t.Errorf("only the passing variant is proven, got %v", passed)
		}
	})
}

func TestPassedCodes_aGreenCaptureProvesTheState(t *testing.T) {
	t.Run("JUIJN-B12: A green capture of a state proves the state", func(t *testing.T) {})
	r := &ExecReport{Cases: []CaseResult{
		{Name: "BUTTN-VR-S01 - Enabled"},
		{Name: "BUTTN-VR-S02 - Disabled", Failed: true},
	}}
	passed, seen := r.PassedCodes(), r.SeenCodes()
	if !passed["BUTTN-VR-S01"] || !passed["BUTTN-S01"] {
		t.Errorf("the green capture proves its scenario and its state: %v", passed)
	}
	if passed["BUTTN-S02"] || !seen["BUTTN-S02"] {
		t.Errorf("a red capture proves nothing and is still seen: passed %v seen %v", passed, seen)
	}
}
