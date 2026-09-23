package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const completeFlag = "| Scenario | When the value | Then |\n" +
	"| --- | --- | --- |\n" +
	"| `CHKUT-G01` | `= \"off\"` | the old checkout answers |\n" +
	"| `CHKUT-G02` | `= \"on\"` | the new checkout answers |\n" +
	"| `CHKUT-G03` | absent | the header's default holds |\n"

func flagNode() mapx.Node {
	return mapx.Node{ID: "flags/checkout.flag.md", Kind: mapx.KindFlag, Code: "CHKUT"}
}

// --- flag-scenarios-complete ---

// The ABSENT case is the one that breaks in production and the one nobody writes.
func TestFlagScenariosComplete(t *testing.T) {
	v, _ := checkFlagScenariosComplete(completeFlag, flagNode(), "", nil, nil)
	if v != Pass {
		t.Errorf("the flag declares `absent` and the verdict was %v", v)
	}

	noAbsent := "| `CHKUT-G01` | `= \"on\"` | turns on |\n"
	v, msg := checkFlagScenariosComplete(noAbsent, flagNode(), "", nil, nil)
	if v != Fail {
		t.Errorf("the flag does NOT declare the absent case and the verdict was %v", v)
	}
	if !strings.Contains(msg, "@no-absent") {
		t.Errorf("the verdict does not name the waiver: %q", msg)
	}
}

// The waiver is `@no-absent` with a reason — a bare marker does not count, by the same
// rule as every other opt-out: the reason is what answers the question six months later.
func TestFlagScenariosComplete_waiver(t *testing.T) {
	noAbsent := "| `CHKUT-G01` | `= \"on\"` | turns on |\n"

	withReason := noAbsent + "\n@no-absent: read from a local constant, never missing\n"
	if v, _ := checkFlagScenariosComplete(withReason, flagNode(), "", nil, nil); v != Pass {
		t.Errorf("the waiver has a written reason and the verdict was %v", v)
	}

	bare := noAbsent + "\n@no-absent\n"
	if v, _ := checkFlagScenariosComplete(bare, flagNode(), "", nil, nil); v != Fail {
		t.Error("the BARE marker was accepted — the reason is mandatory")
	}
}

func TestFlagScenariosComplete_skips(t *testing.T) {
	notFlag := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkFlagScenariosComplete(completeFlag, notFlag, "", nil, nil); v != Skip {
		t.Errorf("not a flag and the verdict was %v", v)
	}
	if v, _ := checkFlagScenariosComplete("# no table\n", flagNode(), "", nil, nil); v != Skip {
		t.Error("a flag with no scenario at all should Skip")
	}
}

// --- flag-scenario-grammar ---

// The grammar is fixed in order to REFUSE. Prose passes any ruler and confronts nothing.
func TestFlagScenarioGrammar(t *testing.T) {
	if v, _ := checkFlagScenarioGrammar(completeFlag, flagNode(), "", nil, nil); v != Pass {
		t.Errorf("every condition is in the grammar and the verdict was %v", v)
	}

	prose := "| `CHKUT-G01` | when the user is a beta tester | turns on |\n"
	v, msg := checkFlagScenarioGrammar(prose, flagNode(), "", nil, nil)
	if v != Fail {
		t.Errorf("the condition is prose and the verdict was %v", v)
	}
	if !strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("the verdict does not say WHICH scenario: %q", msg)
	}
}

// --- flag-scenario-exists ---

func projectWithFlag(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "flags"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, "flags", "checkout.flag.md")
	if err := os.WriteFile(p, []byte(completeFlag), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestFlagScenarioExists(t *testing.T) {
	root := projectWithFlag(t)
	spec := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}

	exists := "### CRED-V01 — validates the limit   @gated-by CHKUT-G02\n"
	if v, _ := checkFlagScenarioExists(exists, spec, root, nil, nil); v != Pass {
		t.Errorf("the scenario exists and the verdict was %v", v)
	}

	missing := "### CRED-V01 — validates the limit   @gated-by CHKUT-G99\n"
	v, msg := checkFlagScenarioExists(missing, spec, root, nil, nil)
	if v != Fail {
		t.Errorf("the scenario does NOT exist and the verdict was %v", v)
	}
	if !strings.Contains(msg, "CHKUT-G99") {
		t.Errorf("the verdict does not carry the code: %q", msg)
	}

	if v, _ := checkFlagScenarioExists("### CRED-V01 — no citation\n", spec, root, nil, nil); v != Skip {
		t.Error("a spec with no citation should Skip")
	}
}

// --- flag-covered ---

// Per SCENARIO, not per flag: it is precisely the disabled branch that goes unproven under
// a per-flag ruler, and it is the one that will break.
func TestFlagCovered_perScenarioNotPerFlag(t *testing.T) {
	// A test proves only G01. Under a per-FLAG ruler this would pass.
	n := flagNode()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01"}}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	v, msg := checkFlagCovered(completeFlag, n, "", g, nil)
	if v != Fail {
		t.Fatalf("two scenarios without a test and the verdict was %v", v)
	}
	for _, c := range []string{"CHKUT-G02", "CHKUT-G03"} {
		if !strings.Contains(msg, c) {
			t.Errorf("the verdict does not name %s: %q", c, msg)
		}
	}
	if strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("the verdict accuses the scenario that HAS a test: %q", msg)
	}
}

func TestFlagCovered_allProven(t *testing.T) {
	n := flagNode()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01", "CHKUT-G02", "CHKUT-G03"}}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	if v, msg := checkFlagCovered(completeFlag, n, "", g, nil); v != Pass {
		t.Errorf("every scenario has a test and the verdict was %v: %s", v, msg)
	}
}

// Without a map there is no way to know what was proven — and Pass here would be a lie.
func TestFlagCovered_noMapDoesNotClaimPass(t *testing.T) {
	if v, _ := checkFlagCovered(completeFlag, flagNode(), "", nil, nil); v == Pass {
		t.Error("with no map the gate claimed Pass — it had no way to know")
	}
}

// THE TWO QUESTIONS, and why they are kept apart.
//
// "Is there a written test?" is static and always answerable. "Did the test pass?" needs
// ingested execution. Merging them loses the answer to both: on a project that never
// ingested a report (the case of this repository — zero `proven_codes` in the graph), the
// execution-only ruler accused EVERY scenario, including those with tests written and
// passing.
//
// And the static question alone would be the opposite error, a worse one: a written test
// may never have run.
func TestFlagCovered_noTestAtAllIsReportedEvenWithoutIngestion(t *testing.T) {
	root := t.TempDir()
	// A test node that exists and names NO scenario.
	if err := os.WriteFile(filepath.Join(root, "t_test.go"), []byte("func TestNothing(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t_test.go", Kind: mapx.KindTest}}}

	v, msg := checkFlagCovered(completeFlag, flagNode(), root, g, nil)
	if v != Fail {
		t.Fatalf("no test names the scenarios — expected Fail, got %v", v)
	}
	for _, c := range []string{"CHKUT-G01", "CHKUT-G02", "CHKUT-G03"} {
		if !strings.Contains(msg, c) {
			t.Errorf("the verdict does not name %s: %q", c, msg)
		}
	}
}

// WRITTEN but NEVER RUN: the gate has to say so, not "no test".
//
// It is the case the request named: a test may have been written and never had its result
// collected. The fix is different — run the suite, not write a test.
func TestFlagCovered_writtenButNotRunSaysWhichOfTheTwo(t *testing.T) {
	root := t.TempDir()
	body := "func TestX(t *testing.T) { /* CHKUT-G01 */ }\nconst c = \"CHKUT-G01\"\n"
	if err := os.WriteFile(filepath.Join(root, "t_test.go"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t_test.go", Kind: mapx.KindTest}}}

	v, msg := checkFlagCovered(completeFlag, flagNode(), root, g, nil)
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	// G01 has a WRITTEN test: it must not show up as "no test names it".
	if !strings.Contains(msg, "ingest") {
		t.Errorf("the verdict does not tell written-but-not-run apart: %q", msg)
	}
	if !strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("the verdict does not name the scenario written and not run: %q", msg)
	}
}

// A CODE IN A COMMENT does not count as a written test — the same ruler as
// `feature-test-match`: a citation is a reference, not an implementation.
func TestFlagCovered_codeInCommentIsNotAWrittenTest(t *testing.T) {
	root := t.TempDir()
	commentOnly := "// CHKUT-G01 is handled elsewhere\nfunc TestX(t *testing.T) {}\n"
	if err := os.WriteFile(filepath.Join(root, "t_test.go"), []byte(commentOnly), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t_test.go", Kind: mapx.KindTest}}}

	_, msg := checkFlagCovered(completeFlag, flagNode(), root, g, nil)
	if strings.Contains(msg, "ingest") {
		t.Errorf("a code only in a COMMENT passed as a written test: %q", msg)
	}
}

// And with ingested execution, the proven scenario leaves the accusation.
func TestFlagCovered_ingestedAndGreenPasses(t *testing.T) {
	root := t.TempDir()
	// The signal lives on the FLAG's node, as `scenario-coverage`'s lives on the spec's:
	// ingestion crosses the proven codes with the ones the node DECLARES.
	n := flagNode()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01", "CHKUT-G02", "CHKUT-G03"}}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	if v, msg := checkFlagCovered(completeFlag, n, root, g, nil); v != Pass {
		t.Errorf("all proven and the verdict was %v: %s", v, msg)
	}
}

// --- flag-scenario-governs: the RETURN LEG of `flag-scenario-exists` ---
//
// That one confronts the spec citing a scenario that does not exist; this one, the
// scenario nobody cites. A scenario no rule invokes is a DECLARED path that governs
// nothing.

// The `gated-by` edge ARRIVES at the flag and carries, in Method, the invoked scenario.
func citing(codes ...string) *mapx.Graph {
	g := &mapx.Graph{Nodes: []mapx.Node{flagNode()}}
	for _, c := range codes {
		g.Edges = append(g.Edges, mapx.Edge{
			From: "a.spec.md", To: "flags/checkout.flag.md",
			Type: mapx.EdgeGatedBy, Method: c,
		})
	}
	return g
}

func TestFlagScenarioGoverns_scenarioWithoutRuleIsReported(t *testing.T) {
	t.Run("CHKUT-G0X: a scenario no rule invokes is reported", func(t *testing.T) {})
	// Only G01 is cited; G02 and G03 are left loose.
	v, msg := checkFlagScenarioGoverns(completeFlag, flagNode(), "", citing("CHKUT-G01"), nil)
	if v != Fail {
		t.Fatalf("two scenarios without a rule and the verdict was %v", v)
	}
	for _, c := range []string{"CHKUT-G02", "CHKUT-G03"} {
		if !strings.Contains(msg, c) {
			t.Errorf("the verdict does not name %s: %q", c, msg)
		}
	}
	if strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("it accused the scenario that IS cited: %q", msg)
	}
}

func TestFlagScenarioGoverns_allCitedPasses(t *testing.T) {
	g := citing("CHKUT-G01", "CHKUT-G02", "CHKUT-G03")
	if v, msg := checkFlagScenarioGoverns(completeFlag, flagNode(), "", g, nil); v != Pass {
		t.Errorf("all cited and the verdict was %v: %s", v, msg)
	}
}

// The per-scenario waiver: some paths exist and need no rule to name them — the `off` that
// returns to the old behaviour, already governed by the rules that always held.
func TestFlagScenarioGoverns_perScenarioWaiver(t *testing.T) {
	withWaiver := "| Scenario | When the value | Then |\n| --- | --- | --- |\n" +
		"| `CHKUT-G01` | `= \"off\"` | the old behaviour holds @no-govern: the rules that always held already govern it |\n"
	if v, msg := checkFlagScenarioGoverns(withWaiver, flagNode(), "", citing(), nil); v != Pass {
		t.Errorf("the waiver has a written reason and the verdict was %v: %s", v, msg)
	}

	bare := strings.Replace(withWaiver, "@no-govern: the rules that always held already govern it", "@no-govern", 1)
	if v, _ := checkFlagScenarioGoverns(bare, flagNode(), "", citing(), nil); v != Fail {
		t.Error("the BARE marker was accepted — the reason is mandatory")
	}
}

func TestFlagScenarioGoverns_skips(t *testing.T) {
	notFlag := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkFlagScenarioGoverns(completeFlag, notFlag, "", citing(), nil); v != Skip {
		t.Errorf("not a flag and the verdict was %v", v)
	}
	if v, _ := checkFlagScenarioGoverns(completeFlag, flagNode(), "", nil, nil); v == Pass {
		t.Error("with no map the gate claimed Pass — it had no way to know who cites")
	}
}
