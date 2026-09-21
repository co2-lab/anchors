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

const specWithFailure = "| `CRED-E01` | insufficient balance | refuses |\n"

// A failure catalogued and not handled is a promise the code does not keep: the spec says
// how the unit fails, and nothing in it deals with that.
func TestFailureHandled_declaredAndUntreatedFails(t *testing.T) {
	root, n, g, cfg := failureProject(t, specWithFailure, "func f() int {\n\treturn 1\n}\n")
	v, msg := checkFailureHandled(specWithFailure, n, root, g, cfg)
	if v != Fail {
		t.Fatalf("expected Fail, got %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "CRED-E01") {
		t.Errorf("the verdict must name the failure: %s", msg)
	}
}

// TREATING IS NOT "having a catch". An `if err != nil` handles just as much, and cravings
// for one syntax would nail the gate to one family of languages.
func TestFailureHandled_anyHandlingShapeCounts(t *testing.T) {
	code := "func f() error {\n\tif err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n"
	root, n, g, cfg := failureProject(t, specWithFailure, code)
	if v, msg := checkFailureHandled(specWithFailure, n, root, g, cfg); v != Pass {
		t.Errorf("an `if err != nil` handles the failure: %v (%s)", v, msg)
	}
}

// The piece that sustains everything else, and the one nobody charges: a handling path
// that swallows the failure without logging is the perfect silence — it happens, nothing
// knows, and no tool downstream has anything to read.
func TestFailureLogged_handledAndSilentFails(t *testing.T) {
	code := "func f() error {\n\tif err != nil {\n\t\treturn nil\n\t}\n\treturn nil\n}\n"
	root, n, g, cfg := failureProject(t, specWithFailure, code)
	v, msg := checkFailureLogged(specWithFailure, n, root, g, cfg)
	if v != Fail {
		t.Fatalf("a silent handling must fail: %v (%s)", v, msg)
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
	code := "func f() error {\n\tif err != nil {\n\t\treturn err\n\t}\n\treturn nil\n}\n"
	root, n, g, cfg := failureProject(t, "# spec with no failures\n", code)
	if v, msg := checkFailureDeclared("# spec with no failures\n", n, root, g, cfg); v != Fail {
		t.Errorf("handling with no declaration must fail: %v (%s)", v, msg)
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
	spec := "| `CRED-E01` | insufficient balance | refuses | @resilient: the partner returns null during migration |\n"
	root, n, g, cfg := failureProject(t, spec, "func f() int {\n\treturn 1\n}\n")
	if v, msg := checkFailureHandled(spec, n, root, g, cfg); v == Fail {
		t.Errorf("a failure marked resilient must not be charged: %s", msg)
	}
}

// A bare marker waives nothing — the same rule as every other opt-out. `@resilient` with
// no why would be the silence the gates exist to end.
func TestFailureHandled_bareResilientMarkerDoesNotCount(t *testing.T) {
	spec := "| `CRED-E01` | insufficient balance | refuses | @resilient |\n"
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
