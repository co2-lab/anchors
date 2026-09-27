package gate

import (
	"github.com/co2-lab/anchors/internal/config"
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
	t.Run("FLSCF-B04: A flag that declares the absent case passes completeness", func(t *testing.T) {})
	t.Run("FLSCF-B05: A flag without the absent case fails, naming the waiver", func(t *testing.T) {})
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
	t.Run("FLSCF-B06: The absent waiver needs a written reason, outside backticks", func(t *testing.T) {})
	t.Run("FLSCF-X01: No waiver is accepted without a written reason", func(t *testing.T) {})
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
	t.Run("FLSCF-B01: The flag gates skip what is not a flag, and a flag with no scenario", func(t *testing.T) {})
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
	t.Run("FLSCF-B02: A flag whose every condition is in the grammar passes", func(t *testing.T) {})
	t.Run("FLSCF-B03: A condition in prose fails the grammar, naming the scenario", func(t *testing.T) {})
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
	t.Run("FLSCF-B07: Only a citation of a G code is confronted", func(t *testing.T) {})
	t.Run("FLSCF-B08: A citation of a declared scenario passes", func(t *testing.T) {})
	t.Run("FLSCF-B09: Citations of scenarios that do not exist fail, each named once and sorted", func(t *testing.T) {})
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
	t.Run("FLSCF-B14: Coverage is judged per scenario", func(t *testing.T) {})
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
	t.Run("FLSCF-B14: Coverage is judged per scenario", func(t *testing.T) {})
	n := flagNode()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G01", "CHKUT-G02", "CHKUT-G03"}}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	if v, msg := checkFlagCovered(completeFlag, n, "", g, nil); v != Pass {
		t.Errorf("every scenario has a test and the verdict was %v: %s", v, msg)
	}
}

// Without a map there is no way to know what was proven — and Pass here would be a lie.
func TestFlagCovered_noMapDoesNotClaimPass(t *testing.T) {
	t.Run("FLSCF-I01: A gate that could not measure never answers Pass", func(t *testing.T) {})
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
	t.Run("FLSCF-B15: A scenario no test names fails as having no test", func(t *testing.T) {})
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
	t.Run("FLSCF-B16: A written test not yet proven is told apart: not ingested, or ingested and not green", func(t *testing.T) {})
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
	t.Run("FLSCF-B15: A scenario no test names fails as having no test", func(t *testing.T) {})
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
	t.Run("FLSCF-B11: A scenario no rule cites fails governance", func(t *testing.T) {})
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
	t.Run("FLSCF-B11: A scenario no rule cites fails governance", func(t *testing.T) {})
	g := citing("CHKUT-G01", "CHKUT-G02", "CHKUT-G03")
	if v, msg := checkFlagScenarioGoverns(completeFlag, flagNode(), "", g, nil); v != Pass {
		t.Errorf("all cited and the verdict was %v: %s", v, msg)
	}
}

// The per-scenario waiver: some paths exist and need no rule to name them — the `off` that
// returns to the old behaviour, already governed by the rules that always held.
func TestFlagScenarioGoverns_perScenarioWaiver(t *testing.T) {
	t.Run("FLSCF-B12: A scenario with a reasoned governance waiver is not charged", func(t *testing.T) {})
	t.Run("FLSCF-X01: No waiver is accepted without a written reason", func(t *testing.T) {})
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
	t.Run("FLSCF-I01: A gate that could not measure never answers Pass", func(t *testing.T) {})
	notFlag := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}
	if v, _ := checkFlagScenarioGoverns(completeFlag, notFlag, "", citing(), nil); v != Skip {
		t.Errorf("not a flag and the verdict was %v", v)
	}
	if v, _ := checkFlagScenarioGoverns(completeFlag, flagNode(), "", nil, nil); v == Pass {
		t.Error("with no map the gate claimed Pass — it had no way to know who cites")
	}
}

// Every flag-side gate skips a node that is not a flag, and a flag that declares no
// scenario; the citation gate skips a node that is not a spec.
func TestFlagGates_skipWhatTheyDoNotConfront(t *testing.T) {
	t.Run("FLSCF-B01: The flag gates skip what is not a flag, and a flag with no scenario", func(t *testing.T) {})
	notFlag := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}
	g := &mapx.Graph{Nodes: []mapx.Node{flagNode()}}
	gates := map[string]func(string, mapx.Node, string, *mapx.Graph) (Verdict, string){
		"flag-scenario-grammar": func(c string, n mapx.Node, r string, g *mapx.Graph) (Verdict, string) {
			return checkFlagScenarioGrammar(c, n, r, g, nil)
		},
		"flag-scenarios-complete": func(c string, n mapx.Node, r string, g *mapx.Graph) (Verdict, string) {
			return checkFlagScenariosComplete(c, n, r, g, nil)
		},
		"flag-scenario-governs": func(c string, n mapx.Node, r string, g *mapx.Graph) (Verdict, string) {
			return checkFlagScenarioGoverns(c, n, r, g, nil)
		},
		"flag-covered": func(c string, n mapx.Node, r string, g *mapx.Graph) (Verdict, string) {
			return checkFlagCovered(c, n, r, g, nil)
		},
	}
	for name, check := range gates {
		if v, _ := check(completeFlag, notFlag, "", g); v != Skip {
			t.Errorf("%s: a spec node should Skip, got %v", name, v)
		}
		if v, _ := check("# no table\n", flagNode(), "", g); v != Skip {
			t.Errorf("%s: a flag with no scenario should Skip, got %v", name, v)
		}
	}
	cite := "### CRED-V01 — x   @gated-by CHKUT-G02\n"
	if v, _ := checkFlagScenarioExists(cite, flagNode(), projectWithFlag(t), nil, nil); v != Skip {
		t.Errorf("flag-scenario-exists: a flag node should Skip, got %v", v)
	}
}

// A waiver quoted in backticks is prose ABOUT the waiver, not the waiver itself.
func TestFlagScenariosComplete_quotedWaiverDoesNotCount(t *testing.T) {
	t.Run("FLSCF-B06: The absent waiver needs a written reason, outside backticks", func(t *testing.T) {})
	noAbsent := "| `CHKUT-G01` | `= \"on\"` | turns on |\n"
	quoted := noAbsent + "\nwaive with `@no-absent: a reason` when needed\n"
	if v, _ := checkFlagScenariosComplete(quoted, flagNode(), "", nil, nil); v != Fail {
		t.Errorf("a waiver quoted in backticks was accepted: %v", v)
	}
}

// Only a `G` code is a flag citation, and the unknown codes come back once each, sorted.
func TestFlagScenarioExists_citationsAndTheirReport(t *testing.T) {
	t.Run("FLSCF-B07: Only a citation of a G code is confronted", func(t *testing.T) {})
	t.Run("FLSCF-B09: Citations of scenarios that do not exist fail, each named once and sorted", func(t *testing.T) {})
	root := projectWithFlag(t)
	spec := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}

	if v, _ := checkFlagScenarioExists("### CRED-V01 — x   @gated-by CRED-B03\n", spec, root, nil, nil); v != Skip {
		t.Errorf("a citation of a B code is not a flag citation, and the verdict was %v", v)
	}

	twice := "### A-V01 @gated-by CHKUT-G98\n### A-V02 @gated-by CHKUT-G97\n### A-V03 @gated-by `CHKUT-G98`\n"
	v, msg := checkFlagScenarioExists(twice, spec, root, nil, nil)
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if strings.Count(msg, "CHKUT-G98") != 1 || !strings.Contains(msg, "CHKUT-G97, CHKUT-G98") {
		t.Errorf("the unknown codes should appear once each, sorted: %q", msg)
	}
}

// Without a map the per-scenario gates answer Pending: neither Pass nor Fail.
func TestFlagGates_noMapIsPending(t *testing.T) {
	t.Run("FLSCF-B10: Without a map the governance gate is Pending", func(t *testing.T) {})
	t.Run("FLSCF-B13: Without a map the coverage gate is Pending", func(t *testing.T) {})
	if v, _ := checkFlagScenarioGoverns(completeFlag, flagNode(), "", nil, nil); v != Pending {
		t.Errorf("flag-scenario-governs with no map: %v", v)
	}
	if v, _ := checkFlagCovered(completeFlag, flagNode(), "", nil, nil); v != Pending {
		t.Errorf("flag-covered with no map: %v", v)
	}
}

// Execution was ingested and the written test is not among the proven codes: the gate
// says the test did not pass, not that it never ran.
func TestFlagCovered_ingestedButNotGreenSaysSo(t *testing.T) {
	t.Run("FLSCF-B16: A written test not yet proven is told apart: not ingested, or ingested and not green", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "t_test.go"), []byte("const c = \"CHKUT-G01\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	n := flagNode()
	n.Signal = &mapx.TestSignal{ProvenCodes: []string{"CHKUT-G02", "CHKUT-G03"}}
	g := &mapx.Graph{Nodes: []mapx.Node{n, {ID: "t_test.go", Kind: mapx.KindTest}}}
	v, msg := checkFlagCovered(completeFlag, n, root, g, nil)
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	if !strings.Contains(msg, "did NOT pass") || strings.Contains(msg, "no execution ingested") {
		t.Errorf("the verdict should say the written test did not pass: %q", msg)
	}
	if !strings.Contains(msg, "CHKUT-G01") {
		t.Errorf("the verdict does not name the scenario: %q", msg)
	}
}

// The @gated-by regex was fixed at 3..6 characters: with a declared length of 7, a
// citation of a scenario that does not exist was read as no citation and skipped.
func TestFlagScenarioExists_readsEveryDeclaredCodeLength(t *testing.T) {
	t.Run("FLSCF-B17: A gated-by citation is read at the code length the project declares", func(t *testing.T) {})
	codeLengthsForTest(t, 7)
	spec := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}
	v, msg := checkFlagScenarioExists("### CREDITS-V01 — x   @gated-by CHKUTXY-G99\n", spec, projectWithFlag(t), nil, nil)
	if v != Fail || !strings.Contains(msg, "CHKUTXY-G99") {
		t.Errorf("a citation of a missing 7-character scenario = %v (%s); want Fail naming it", v, msg)
	}
}

// An unreadable flags folder is Pending, and the verdict names that cause. The branch
// answered with the "no map loaded" message, which sent the reader to build a map this
// gate never reads.
func TestFlagScenarioExists_unreadableFlagsIsPendingWithTheCause(t *testing.T) {
	t.Run("FLSCF-E02: Flags that cannot be read leave the citation Pending, naming the flags folder", func(t *testing.T) {})
	if os.Geteuid() == 0 {
		t.Skip("root reads a folder without permission")
	}
	root := projectWithFlag(t)
	dir := filepath.Join(root, "flags")
	if err := os.Chmod(dir, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	spec := mapx.Node{ID: "a.spec.md", Kind: mapx.KindSpec}
	v, msg := checkFlagScenarioExists("### CRED-V01 — x   @gated-by CHKUT-G02\n", spec, root, nil, nil)
	_, noMap := pendingNoMap()
	if v != Pending || msg == noMap || !strings.Contains(msg, "flags") {
		t.Errorf("unreadable flags: want Pending naming the flags folder, got %v (%s)", v, msg)
	}
}

// flagWithTest confronts the complete flag with one test file under a configuration.
func flagWithTest(t *testing.T, body string, cfg *config.Config) (Verdict, string) {
	t.Helper()
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "t_test.go"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "t_test.go", Kind: mapx.KindTest}}}
	return checkFlagCovered(completeFlag, flagNode(), root, g, cfg)
}

func TestFlagCovered_TitlesWhenDeclared(t *testing.T) {
	t.Run("FLSCF-B18: With a tests source a flag scenario is written only when a title cites it", func(t *testing.T) {})
	body := "var fixture = \"CHKUT-G01\"\nfunc TestX(t *testing.T) { t.Run(\"unrelated\", nil) }\n"
	goFamily := &config.Config{Dialect: &config.Dialect{Family: "go"}}
	if _, msg := flagWithTest(t, body, goFamily); strings.Contains(msg, "WRITTEN") {
		t.Fatalf("with titles declared, a code in a fixture writes no test, got %s", msg)
	}
	if _, msg := flagWithTest(t, body, nil); !strings.Contains(msg, "WRITTEN") {
		t.Fatalf("without a declaration the code in the file counts as written, got %s", msg)
	}
	titled := "func TestX(t *testing.T) { t.Run(\"CHKUT-G01: on\", nil) }\n"
	if _, msg := flagWithTest(t, titled, goFamily); !strings.Contains(msg, "WRITTEN") {
		t.Fatalf("a title citing the code writes its test, got %s", msg)
	}
}

func TestFlagCovered_FailingSource(t *testing.T) {
	t.Run("FLSCF-E03: A failing tests source fails flag-covered naming the error", func(t *testing.T) {})
	cfg := &config.Config{Dialect: &config.Dialect{Tests: &config.TestsSource{Script: "echo 'lister crashed' >&2; exit 1"}}}
	if v, msg := flagWithTest(t, "func TestX(t *testing.T) {}\n", cfg); v != Fail || !strings.Contains(msg, "lister crashed") {
		t.Fatalf("a failing source must fail naming its error, got %v: %s", v, msg)
	}
}

func TestFlagCovered_FileTheSourceDoesNotDescribe(t *testing.T) {
	t.Run("FLSCF-B18: With a tests source a flag scenario is written only when a title cites it", func(t *testing.T) {})
	ts := &config.Config{Dialect: &config.Dialect{Family: "ts"}}
	if _, msg := flagWithTest(t, "name: 'CHKUT-G01 on'\n", ts); !strings.Contains(msg, "WRITTEN") {
		t.Fatalf("in a file the source lists no test in, the code in it counts as written, got %s", msg)
	}
}
