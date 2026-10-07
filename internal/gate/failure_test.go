// @anchors
//   code: FLTSB
//   ref: FLRAI

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func failureProject(t *testing.T, spec, code string) (string, mapx.Node, *mapx.Graph, *config.Config) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}
	g := &mapx.Graph{Edges: []mapx.Edge{{From: "x.spec.md", To: "x.go", Type: mapx.EdgeSpecifies}}}
	return root, n, g, &config.Config{Dialect: &config.Dialect{Family: "go"}}
}

const specWithFailure = "| `CREDT-E01` | insufficient balance | refuses |\n"

// A failure catalogued and not handled is a promise the code does not keep: the spec says
// how the unit fails, and nothing in it deals with that.
func TestFailureHandled_declaredAndUntreatedFails(t *testing.T) {
	t.Run("FLRAI-B07: A declared failure with no handling in the governed code fails, naming it", func(t *testing.T) {})
	root, n, g, cfg := failureProject(t, specWithFailure, "func f() int {\n\treturn 1\n}\n")
	v, msg := checkFailureHandled(specWithFailure, n, root, g, cfg)
	if v != Fail {
		t.Fatalf("expected Fail, got %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "CREDT-E01") {
		t.Errorf("the verdict must name the failure: %s", msg)
	}
}

// HANDLING IS NOT "having a catch". An `if err != nil` handles just as much, and nailing
// one syntax would tie the gate to one family of languages.
func TestFailureHandled_anyHandlingShapeCounts(t *testing.T) {
	t.Run("FLRAI-B08: Any handling path in the governed code passes failure-handled", func(t *testing.T) {})
	t.Run("FLRAI-X01: One handling path answers for every declared failure", func(t *testing.T) {})
	code := "func f() error {\n\tif err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n"
	spec := specWithFailure + "| `CREDT-E02` | partner down | retries |\n"
	root, n, g, cfg := failureProject(t, spec, code)
	if v, msg := checkFailureHandled(spec, n, root, g, cfg); v != Pass {
		t.Errorf("an `if err != nil` handles the failure: %v (%s)", v, msg)
	}
}

// The piece that sustains everything else, and the one nobody charges: a handling path
// that swallows the failure without logging is the perfect silence — it happens, nothing
// knows, and no tool downstream has anything to read.
func TestFailureLogged_handledAndSilentFails(t *testing.T) {
	t.Run("FLRAI-B11: A handling that records nothing fails failure-logged, naming the failure", func(t *testing.T) {})
	t.Run("FLRAI-B12: A handling that records the occurrence passes failure-logged", func(t *testing.T) {})
	code := "func f() error {\n\tif err != nil {\n\t\treturn nil\n\t}\n\treturn nil\n}\n"
	root, n, g, cfg := failureProject(t, specWithFailure, code)
	v, msg := checkFailureLogged(specWithFailure, n, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "CREDT-E01") {
		t.Fatalf("a silent handling must fail, naming the failure: %v (%s)", v, msg)
	}
	withLog := "func f() error {\n\tif err != nil {\n\t\tlog.Error(\"x\")\n\t\treturn err\n\t}\n\treturn nil\n}\n"
	root2, n2, g2, cfg2 := failureProject(t, specWithFailure, withLog)
	if v, msg := checkFailureLogged(specWithFailure, n2, root2, g2, cfg2); v != Pass {
		t.Errorf("a handling that logs must pass: %v (%s)", v, msg)
	}
}

// The inverse, and the most common case: somebody wrote a defence and never declared what
// it prevents. The defence may be right — what is missing is the spec saying which failure
// it answers, so whoever reads it later knows whether it still applies.
func TestFailureDeclared_handlingWithoutDeclarationFails(t *testing.T) {
	t.Run("FLRAI-B13: Handling in the code with no failure declared in the spec fails", func(t *testing.T) {})
	code := "func f() error {\n\tif err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n"
	root, n, g, cfg := failureProject(t, "# spec with no failures\n", code)
	if v, msg := checkFailureDeclared("# spec with no failures\n", n, root, g, cfg); v != Fail {
		t.Errorf("handling with no declaration must fail: %v (%s)", v, msg)
	}
}

// A unit whose handling matches are all normal flow closes its failure section with
// `none — <why>`. The reason is mandatory: a bare `none` does not close it.
func TestFailureDeclared_sectionClosedAsNone(t *testing.T) {
	t.Run("FLRAI-B16: An Errors section closed with none and a reason passes failure-declared", func(t *testing.T) {})
	code := "func f(m map[string]int) {\n\tif m == nil {\n\t\tm = map[string]int{}\n\t}\n}\n"
	for name, c := range map[string]struct {
		spec string
		want Verdict
	}{
		"none with a reason": {"# spec\n\n## Errors\n\nnone — the only nil check is a lazy map init\n", Pass},
		"pt-BR none":         {"# spec\n\n## Erros\n\nnenhuma: o nil check inicializa o mapa\n", Pass},
		"bare none":          {"# spec\n\n## Errors\n\nnone\n", Fail},
		"none outside it":    {"# spec\n\nnone — said in the overview\n", Fail},
		"section with prose": {"# spec\n\n## Errors\n\nThe unit fails when… none — later\n", Fail},
	} {
		root, n, g, cfg := failureProject(t, c.spec, code)
		if v, msg := checkFailureDeclared(c.spec, n, root, g, cfg); v != c.want {
			t.Errorf("%s: got %v, want %v (%s)", name, v, c.want, msg)
		}
	}
}

// `@resilient` is a THIRD assertion, different from the two that already existed:
//
//	@no-<thing>: <reason>   "it will never have one"  — permanent waiver
//	@TBD: <reason>          "it does not have one yet" — debt, keeps showing up
//	@resilient: <reason>    "it happens, I know why, and it is handled"
//
// It is not a waiver (the failure is real and happens) nor debt (nothing is pending): it
// is knowledge acquired, and the written reason is what tells it apart from silencing an
// alert.
func TestFailureHandled_resilientStopsTheCharge(t *testing.T) {
	t.Run("FLRAI-B09: A failure marked resilient with a reason is not charged", func(t *testing.T) {})
	spec := "| `CREDT-E01` | insufficient balance | refuses | @resilient: the partner returns null during migration |\n"
	root, n, g, cfg := failureProject(t, spec, "func f() int {\n\treturn 1\n}\n")
	if v, msg := checkFailureHandled(spec, n, root, g, cfg); v != Pass {
		t.Errorf("a failure marked resilient must not be charged: %v (%s)", v, msg)
	}
	silent := "func f() error {\n\tif err != nil {\n\t\treturn nil\n\t}\n\treturn nil\n}\n"
	root, n, g, cfg = failureProject(t, spec, silent)
	if v, msg := checkFailureLogged(spec, n, root, g, cfg); v != Pass {
		t.Errorf("failure-logged does not charge a resilient failure either: %v (%s)", v, msg)
	}
}

// A bare marker waives nothing — the same rule as every other opt-out. `@resilient` with
// no why would be the silence the gates exist to end.
func TestFailureHandled_bareResilientMarkerDoesNotCount(t *testing.T) {
	t.Run("FLRAI-B10: A bare resilient marker exempts nothing", func(t *testing.T) {})
	spec := "| `CREDT-E01` | insufficient balance | refuses | @resilient |\n"
	root, n, g, cfg := failureProject(t, spec, "func f() int {\n\treturn 1\n}\n")
	if v, _ := checkFailureHandled(spec, n, root, g, cfg); v != Fail {
		t.Errorf("a bare @resilient must not waive: got %v", v)
	}
}

// Without a dialect the verdict is UNDETERMINED, never approval: reading code it cannot
// recognise would stamp what was never checked.
func TestFailureHandled_noDialectIsUndetermined(t *testing.T) {
	root, n, g, _ := failureProject(t, specWithFailure, "func f() int { return 1 }\n")
	if v, _ := checkFailureHandled(specWithFailure, n, root, g, &config.Config{}); v != Pending {
		t.Errorf("with no dialect expected Pending, got %v", v)
	}
}

// The THIRD conclusion of the observation layer, and the one no observability tool
// records: the difference between "nobody investigated" and "we investigated and still do
// not know". The second is KNOWLEDGE — without it the next person starts from zero, ruling
// out what somebody already ruled out.
func TestFailureConclusions_readsTheThreeOutcomes(t *testing.T) {
	t.Run("FLRAI-B17: The conclusions read the resilient and observing reasons of each failure whole", func(t *testing.T) {})
	spec := "| `CREDT-E01` | balance | refuses |\n" +
		"| `CREDT-E02` | partner down | retries | @resilient: the partner restarts at 3am daily; the retry covers it |\n" +
		"| `CREDT-E03` | timeout | refuses | @observing: ruled out partner retry and network latency; only on migrated accounts |\n"
	got := FailureConclusions(spec)

	if c := got["CREDT-E01"]; c.Resilient != "" || c.Observing != "" {
		t.Errorf("a failure with no conclusion must carry none: %+v", c)
	}
	// The WHOLE reason is what matters. The marker pattern stops at the first token —
	// it only proves the reason exists — and reading the reason from it truncated
	// `@resilient: the partner restarts...` down to `the`.
	if c := got["CREDT-E02"]; !strings.Contains(c.Resilient, "retry covers it") {
		t.Errorf("the reason must be read whole, got %q", c.Resilient)
	}
	if c := got["CREDT-E03"]; !strings.Contains(c.Observing, "migrated accounts") {
		t.Errorf("the reason must be read whole, got %q", c.Observing)
	}
}

// A table cell ends at the pipe: the reason is what the author wrote in THIS column, and
// swallowing the next one would attribute to the conclusion a text belonging to another
// field.
func TestFailureConclusions_theReasonStopsAtTheCell(t *testing.T) {
	t.Run("FLRAI-B18: A conclusion reason ends at its table cell", func(t *testing.T) {})
	spec := "| `CREDT-E01` | cond | @resilient: the real reason | another column |\n"
	c := FailureConclusions(spec)["CREDT-E01"]
	if strings.Contains(c.Resilient, "another column") {
		t.Errorf("the reason swallowed the next cell: %q", c.Resilient)
	}
	if !strings.Contains(c.Resilient, "the real reason") {
		t.Errorf("the reason was lost: %q", c.Resilient)
	}
}

// failureGates are the three gates of this file, for the rules they share.
var failureGates = map[string]func(string, mapx.Node, string, *mapx.Graph, *config.Config) (Verdict, string){
	"failure-handled":  checkFailureHandled,
	"failure-logged":   checkFailureLogged,
	"failure-declared": checkFailureDeclared,
}

func TestFailureGates_onlyConfrontSpecs(t *testing.T) {
	t.Run("FLRAI-B01: Every failure gate skips an artifact that is not a spec", func(t *testing.T) {})
	root, _, g, cfg := failureProject(t, specWithFailure, "func f() {}\n")
	for name, check := range failureGates {
		for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindTest, mapx.KindFeature} {
			if v, _ := check(specWithFailure, mapx.Node{ID: "x.spec.md", Kind: k}, root, g, cfg); v != Skip {
				t.Errorf("%s on %s: expected Skip, got %v", name, k, v)
			}
		}
	}
}

// A failure rule is a CATALOGUED item: a heading, a table row or a bold bullet. A code
// cited in prose is a mention, not a declaration.
func TestFailureHandled_readsTheThreeCataloguedForms(t *testing.T) {
	t.Run("FLRAI-B02: A failure rule is read in the heading, table row and bullet forms, not in prose", func(t *testing.T) {})
	spec := "### CREDT-E01 — balance\n| `CREDT-E02` | partner down | retries |\n- **CREDT-E03** timeout\n"
	root, n, g, cfg := failureProject(t, spec, "func f() int {\n\treturn 1\n}\n")
	v, msg := checkFailureHandled(spec, n, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "CREDT-E01, CREDT-E02, CREDT-E03") {
		t.Errorf("the three forms are declared failures, sorted: %v (%s)", v, msg)
	}
	prose := "The unit may raise CREDT-E01 when the balance is short.\n"
	if v, msg := checkFailureHandled(prose, n, root, g, cfg); v != Skip {
		t.Errorf("a code in prose declares nothing: %v (%s)", v, msg)
	}
}

func TestFailureHandledAndLogged_skipASpecWithNoFailure(t *testing.T) {
	t.Run("FLRAI-B03: A spec that declares no failure is skipped by failure-handled and failure-logged", func(t *testing.T) {})
	const spec = "# spec\n\n| `CREDT-B01` | a behaviour |\n"
	root, n, g, cfg := failureProject(t, spec, "func f() int {\n\treturn 1\n}\n")
	if v, msg := checkFailureHandled(spec, n, root, g, cfg); v != Skip {
		t.Errorf("failure-handled: %v (%s)", v, msg)
	}
	if v, msg := checkFailureLogged(spec, n, root, g, cfg); v != Skip {
		t.Errorf("failure-logged: %v (%s)", v, msg)
	}
}

// Without the patterns that recognise a handling (or, for failure-logged, a record) the
// gates measured nothing, and a Pass would stamp what was never checked.
func TestFailureGates_undeterminedWithoutTheDialectPatterns(t *testing.T) {
	t.Run("FLRAI-B04: Without the dialect patterns the failure gates are Pending", func(t *testing.T) {})
	code := "func f() error {\n\tif err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n"
	root, n, g, _ := failureProject(t, specWithFailure, code)
	for name, check := range failureGates {
		if v, msg := check(specWithFailure, n, root, g, &config.Config{}); v != Pending {
			t.Errorf("%s with no dialect: expected Pending, got %v (%s)", name, v, msg)
		}
	}
	handleOnly := &config.Config{Dialect: &config.Dialect{HandlePatterns: []string{`if\s+err\s*!=\s*nil`}}}
	if v, msg := checkFailureLogged(specWithFailure, n, root, g, handleOnly); v != Pending {
		t.Errorf("failure-logged with no log patterns: expected Pending, got %v (%s)", v, msg)
	}
}

// No governed code, or none that can be read: there is nothing to confront, and the
// verdict stays open rather than green.
func TestFailureGates_undeterminedWithoutGovernedCode(t *testing.T) {
	t.Run("FLRAI-B05: Without governed code that can be read the failure gates are Pending", func(t *testing.T) {})
	root, n, _, cfg := failureProject(t, specWithFailure, "func f() {}\n")
	unreadable := &mapx.Graph{Edges: []mapx.Edge{{From: n.ID, To: "gone.go", Type: mapx.EdgeSpecifies}}}
	otherType := &mapx.Graph{Edges: []mapx.Edge{{From: n.ID, To: "x.go", Type: mapx.EdgeRealizes}}}
	for name, check := range failureGates {
		for label, g := range map[string]*mapx.Graph{"no map": nil, "unreadable file": unreadable, "no specifies edge": otherType} {
			if v, msg := check(specWithFailure, n, root, g, cfg); v != Pending {
				t.Errorf("%s, %s: expected Pending, got %v (%s)", name, label, v, msg)
			}
		}
	}
}

// A handling that only exists in a comment handles nothing.
func TestFailureHandled_aHandlingInACommentDoesNotCount(t *testing.T) {
	t.Run("FLRAI-B06: A handling written only in a comment line does not count", func(t *testing.T) {})
	code := "func f() int {\n\t// if err != nil { return err }\n\treturn 1\n}\n"
	root, n, g, cfg := failureProject(t, specWithFailure, code)
	if v, msg := checkFailureHandled(specWithFailure, n, root, g, cfg); v != Fail {
		t.Errorf("a commented-out handling must not count: %v (%s)", v, msg)
	}
}

// The two gates never charge the same failure together: with no handling at all
// failure-handled charges it and failure-logged steps aside; with handling that records
// nothing it is the other way round.
func TestFailureHandledAndLogged_neverChargeTheSameDefectTwice(t *testing.T) {
	t.Run("FLRAI-I01: failure-handled and failure-logged never both fail the same spec", func(t *testing.T) {})
	for name, c := range map[string]struct {
		code            string
		handled, logged Verdict
	}{
		"no handling":         {"func f() int {\n\treturn 1\n}\n", Fail, Skip},
		"handling, no record": {"func f() error {\n\tif err != nil {\n\t\treturn nil\n\t}\n\treturn nil\n}\n", Pass, Fail},
	} {
		root, n, g, cfg := failureProject(t, specWithFailure, c.code)
		h, _ := checkFailureHandled(specWithFailure, n, root, g, cfg)
		l, _ := checkFailureLogged(specWithFailure, n, root, g, cfg)
		if h != c.handled || l != c.logged {
			t.Errorf("%s: handled=%v logged=%v, want %v and %v", name, h, l, c.handled, c.logged)
		}
	}
}

func TestFailureDeclared_noHandlingIsSkipped(t *testing.T) {
	t.Run("FLRAI-B14: Code with no handling is skipped by failure-declared", func(t *testing.T) {})
	root, n, g, cfg := failureProject(t, "# spec\n", "func f() int {\n\treturn 1\n}\n")
	if v, msg := checkFailureDeclared("# spec\n", n, root, g, cfg); v != Skip {
		t.Errorf("no handling, nothing to declare: %v (%s)", v, msg)
	}
}

func TestFailureDeclared_aDeclaredFailurePasses(t *testing.T) {
	t.Run("FLRAI-B15: Handling in the code and a failure declared in the spec passes", func(t *testing.T) {})
	code := "func f() error {\n\tif err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n"
	root, n, g, cfg := failureProject(t, specWithFailure, code)
	if v, msg := checkFailureDeclared(specWithFailure, n, root, g, cfg); v != Pass {
		t.Errorf("a declared failure answers the handling: %v (%s)", v, msg)
	}
}

// The failure regex was fixed at 3..6 characters: in a project that declares a length
// of 7, a spec's `-E` rules were not read and both failure gates skipped it.
func TestFailure_readsEveryDeclaredCodeLength(t *testing.T) {
	t.Run("FLRAI-B19: A failure rule is read at the code length the project declares", func(t *testing.T) {})
	codeLengthsForTest(t, 7)
	spec := "| `CREDITS-E01` | partner down | @resilient: retried by the queue |\n- **CREDITS-E02** timeout\n"
	all, resilient := declaredFailures(spec)
	if strings.Join(all, " ") != "CREDITS-E01 CREDITS-E02" || strings.Join(resilient, " ") != "CREDITS-E01" {
		t.Errorf("declaredFailures = %v, %v; want both failures, E01 resilient", all, resilient)
	}
	if c := FailureConclusions(spec)["CREDITS-E01"]; c.Resilient == "" {
		t.Errorf("FailureConclusions lost the resilient reason of CREDITS-E01: %+v", c)
	}
}

// A unit that consumes a fallible source declares how it fails — the source is in its code,
// or in a dependency on a layer the project marks fallible.
func TestFailureDeclared_fallibleSourceAsksForTheFailure(t *testing.T) {
	t.Run("FLRAI-B20: A unit with a fallible source and no declared failure fails, naming the source", func(t *testing.T) {})
	code := "package x\n\n// useQuery in a comment is no call\nfunc load() {\n\tdata := useQuery(\"budget\") // the fetch\n\t_ = data\n}\n"
	root, n, g, cfg := failureProject(t, "# X\n\n| `CREDT-B01` | shows the budget |\n", code)
	cfg.Dialect.FalliblePatterns = []config.FalliblePattern{{Call: `\buseQuery\(`, Handled: `isError`}}
	v, msg := checkFailureDeclared("# X\n\n| `CREDT-B01` | shows the budget |\n", n, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "x.go:5") || strings.Contains(msg, "x.go:3") {
		t.Errorf("the call names its line, and a comment is no call: %v %s", v, msg)
	}
	if v, msg := checkFailureDeclared("# X\n\n"+specWithFailure, n, root, g, cfg); v != Fail || !strings.Contains(msg, "x.go:5") {
		t.Errorf("a failure that does not name the source does not answer it — \"not found\" is no load failure: %v %s", v, msg)
	}
	for _, spec := range []string{
		"# X\n\n| `CREDT-E02` | `useQuery` fails (no network) | shows the load error and a retry |\n",
		"# X\n\n| `CREDT-B01` | shows the budget |\n\n@no-failure: the value comes from a constant\n",
		"# X\n\n## Errors\n\nnone — the source is a local constant\n",
	} {
		if v, msg := checkFailureDeclared(spec, n, root, g, cfg); v == Fail {
			t.Errorf("a declared failure, a waiver or a closed section answers the source: %v %s\n%s", v, msg, spec)
		}
	}
	t.Run("FLRAI-B21: A dependency on a file of a fallible layer is a fallible source", func(t *testing.T) {})
	cfg.Dialect.FalliblePatterns = nil
	cfg.Layers = map[string]config.Layer{"hook": {Pattern: "hooks/*.ts", Kind: "code", Fallible: true}, "util": {Pattern: "utils/*.ts", Kind: "code"}}
	g.Nodes = []mapx.Node{{ID: "hooks/useBudget.ts", Kind: mapx.KindCode, Layer: "hook"}, {ID: "utils/fmt.ts", Kind: mapx.KindCode, Layer: "util"}}
	g.Edges = append(g.Edges,
		mapx.Edge{From: "x.spec.md", To: "hooks/useBudget.ts", Type: mapx.EdgeDependsOn, Dep: "DEP1"},
		mapx.Edge{From: "x.spec.md", To: "utils/fmt.ts", Type: mapx.EdgeDependsOn, Dep: "DEP2"})
	v, msg = checkFailureDeclared("# X\n\n| `CREDT-B01` | shows the budget |\n", n, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "DEP1 hooks/useBudget.ts") || strings.Contains(msg, "fmt.ts") {
		t.Errorf("the fallible layer's dependency is the source, the other is not: %v %s", v, msg)
	}
	uses := "# X\n\n| `CREDT-E01` | the budget does not load | shows the load error |\n\n## Rule uses\n\n| Rule | Uses |\n| --- | --- |\n| `CREDT-E01` | `DEP1` |\n"
	if v, msg := checkFailureDeclared(uses, n, root, g, cfg); v == Fail {
		t.Errorf("a failure whose uses name the DEPn answers it: %v %s", v, msg)
	}
}

// Each fallible call has its handling beside it — above, in the destructuring that receives
// it, or below, within its window —, or is waived on its line.
func TestFailureHandled_eachFallibleCallHasItsHandling(t *testing.T) {
	t.Run("FLRAI-B22: A fallible call with no handling in its window fails, named by its line", func(t *testing.T) {})
	code := "package x\n\nfunc a() {\n\tconst {\n\t\tdata,\n\t\tisError,\n\t} = useQuery(\"a\")\n}\n\n" + // line 7: handled above
		"func b() {\n\tr := useQuery(\"b\")\n\tif r.isError { return }\n}\n\n" + // line 11: handled below
		"func c() {\n\titems := useQuery(\"c\").data ?? []\n\t_ = items\n}\n\n" + // line 16: not handled
		"func d() {\n\t// @no-handle: the cache answers offline\n\tv := useQuery(\"d\")\n\t_ = v\n}\n" // line 22: waived
	root, n, g, cfg := failureProject(t, specWithFailure, code)
	cfg.Dialect.FalliblePatterns = []config.FalliblePattern{{Call: `\buseQuery\(`, Handled: `\bisError\b`, Window: 3}}
	v, msg := checkFailureHandled(specWithFailure, n, root, g, cfg)
	if v != Fail || !strings.Contains(msg, "x.go:16") {
		t.Fatalf("the unhandled call is named: %v %s", v, msg)
	}
	for _, line := range []string{"x.go:7", "x.go:11", "x.go:22"} {
		if strings.Contains(msg, line) {
			t.Errorf("%s is handled or waived, and is not named: %s", line, msg)
		}
	}
}
