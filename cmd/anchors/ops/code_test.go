package ops

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// codeProject writes anchors.yaml and a map with the given nodes (YAML list items).
func codeProject(t *testing.T, cfg, nodes string) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, config.DefaultFile, cfg)
	writeFile(t, root, mapx.DefaultPath, "version: 4\nnodes:\n"+nodes+"edges: []\n")
	return root
}

// runCode runs `anchors code` UNDER a root, as the CLI does: a root command with
// subcommands would read the unit name as an unknown subcommand.
func runCode(t *testing.T, args ...string) (error, string) {
	t.Helper()
	root := &cobra.Command{Use: "anchors"}
	root.AddCommand(newCodeCmd())
	return runCmd(t, root, append([]string{"code"}, args...)...)
}

// --check answers whether a proposed code is free, and exits non-zero on a collision.
func TestCodeCheckSaysWhoOwnsTheCode(t *testing.T) {
	root := codeProject(t, "version: 1\n", `    - id: src/auth/Login.spec.md
      kind: spec
      code: LOGNS
`)
	err, out := runCode(t, "--root", root, "--check", "logns")
	if !errors.Is(err, errCollision) {
		t.Errorf("a taken code must return errCollision, got %v", err)
	}
	if !strings.Contains(out, "LOGNS is already used by: src/auth") {
		t.Errorf("the owner is not named:\n%s", out)
	}
	err, out = runCode(t, "--root", root, "--check", "FREEX")
	if err != nil || !strings.Contains(out, "✓ FREEX is free") {
		t.Errorf("a free code: %v\n%s", err, out)
	}
}

// Generating for a name whose canonical code is taken gives ANOTHER code, and says why.
func TestCodeGenerateAvoidsATakenCanonical(t *testing.T) {
	canonical := code.Generate("Spacer")
	root := codeProject(t, "version: 1\n", `    - id: ui/Spacer.spec.md
      kind: spec
      code: `+canonical+`
`)
	err, out := runCode(t, "--root", root, "Spacer")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "free code: "+canonical+"\n") {
		t.Errorf("suggested the taken canonical %s:\n%s", canonical, out)
	}
	if !strings.Contains(out, "the canonical "+canonical+" already belongs to ui") {
		t.Errorf("the adjustment is not explained:\n%s", out)
	}

	err, out = runCode(t, "--root", codeProject(t, "version: 1\n", ""), "Spacer")
	if err != nil || !strings.Contains(out, "free code: "+canonical+"\n") {
		t.Errorf("with nothing taken the canonical is the answer: %v\n%s", err, out)
	}
}

// A path in a layer with `code_prefix` starts with the module prefix; a generic basename
// takes its identity from the parent directory.
func TestCodeGenerateFromAPath(t *testing.T) {
	cfg := "layers:\n  auth:\n    pattern: \"src/auth/**/*.ts\"\n    kind: code\n    code_prefix: AU\n"
	root := codeProject(t, cfg, "")
	writeFile(t, root, "src/auth/Session.ts", "export {}\n")
	err, out := runCode(t, "--root", root, "src/auth/Session.ts")
	if err != nil {
		t.Fatal(err)
	}
	want := code.GenerateWithPrefix("Session", "AU")
	if !strings.Contains(out, "free code: "+want) || !strings.Contains(out, "module prefix 'AU'") {
		t.Errorf("want %s from the layer prefix:\n%s", want, out)
	}

	err, out = runCode(t, "--root", root, "src/billing/index.ts")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "generic basename") || strings.Contains(out, "free code: "+code.Generate("index")+"\n") {
		t.Errorf("a generic basename must take the parent dir's identity:\n%s", out)
	}
}

func TestCodeRequiresANameAndAMap(t *testing.T) {
	root := codeProject(t, "version: 1\n", "")
	if err, _ := runCode(t, "--root", root); err == nil || !strings.Contains(err.Error(), "provide the unit name") {
		t.Errorf("no name: %v", err)
	}
	if err, _ := runCode(t, "--root", t.TempDir(), "X"); err == nil || !strings.Contains(err.Error(), "anchors map build") {
		t.Errorf("no map: %v", err)
	}
}

// `code list --check` accuses a DECLARED code outside `code_lengths` and proposes the
// canonical code of the unit's name; a merely CITED one is counted, not accused.
func TestCodeListCheckAccusesOnlyDeclaredCodesOfTheWrongLength(t *testing.T) {
	root := codeProject(t, "version: 1\n", `    - id: app/Wallet.spec.md
      kind: spec
      code: WLTX
      code_declared: true
    - id: app/Budget.spec.md
      kind: spec
      code: BDGTS
      code_declared: true
    - id: app/fixture_test.go
      kind: test
      code: FXTR
`)
	err, out := runCmd(t, newCodeListCmd(), "--root", root, "--check")
	if !errors.Is(err, errCollision) {
		t.Errorf("a divergent code must fail the check, got %v", err)
	}
	want := "WLTX → " + code.Generate("Wallet") + "\tapp"
	if !strings.Contains(out, want) {
		t.Errorf("want %q in:\n%s", want, out)
	}
	if strings.Contains(out, "FXTR →") {
		t.Errorf("a cited code was accused:\n%s", out)
	}
	if !strings.Contains(out, "1 conforming") {
		t.Errorf("the conforming count is wrong:\n%s", out)
	}

	// Filtered to a folder where every declared code conforms, the check passes.
	root = codeProject(t, "version: 1\n", `    - id: ok/Budget.spec.md
      kind: spec
      code: BDGTZ
      code_declared: true
    - id: ok/fixture_test.go
      kind: test
      code: FXTQ
    - id: other/Wallet.spec.md
      kind: spec
      code: WLTQ
      code_declared: true
`)
	err, out = runCmd(t, newCodeListCmd(), "--root", root, "--check", "--in", "ok")
	if err != nil {
		t.Fatalf("the filtered check should pass: %v\n%s", err, out)
	}
	if !strings.Contains(out, "✓ 1 declared code(s) with conforming length (5)") ||
		!strings.Contains(out, "1 code(s) only CITED") {
		t.Errorf("summary:\n%s", out)
	}
}

// --json gives each code its folder, file, kind, title and the work order fields.
func TestCodeListJSONCarriesWhatAConsumerNeeds(t *testing.T) {
	root := codeProject(t, "version: 1\n", `    - id: plans/0002-build.md
      kind: plan
      code: PLNBB
      needs: [PLNAA]
      parent: PRDCT
      revises: [PLNZZ]
`)
	writeFile(t, root, "plans/0002-build.md", "# Plano 0002 — Build\n")
	err, out := runCmd(t, newCodeListCmd(), "--root", root, "--json")
	if err != nil {
		t.Fatal(err)
	}
	var got []struct {
		Code, Arquivo, Kind, Titulo, Parent string
		Onde, Needs, Revises                []string
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	if len(got) != 1 {
		t.Fatalf("want one entry, got %+v", got)
	}
	g := got[0]
	if g.Code != "PLNBB" || g.Arquivo != "plans/0002-build.md" || g.Kind != "plan" || g.Titulo != "Build" ||
		g.Parent != "PRDCT" || strings.Join(g.Needs, ",") != "PLNAA" || strings.Join(g.Revises, ",") != "PLNZZ" ||
		strings.Join(g.Onde, ",") != "plans" {
		t.Errorf("entry = %+v", g)
	}

	err, out = runCmd(t, newCodeListCmd(), "--root", codeProject(t, "version: 1\n", ""), "--json")
	if err != nil || strings.TrimSpace(out) != "[]" {
		t.Errorf("an empty project in JSON is `[]`, got %v %q", err, out)
	}
}

func TestCodeListEmptyAndBrokenInputs(t *testing.T) {
	err, out := runCmd(t, newCodeListCmd(), "--root", codeProject(t, "version: 1\n", ""))
	if err != nil || !strings.Contains(out, "the map has no node with identity") {
		t.Errorf("empty map: %v\n%s", err, out)
	}
	root := t.TempDir()
	if err, _ := runCmd(t, newCodeListCmd(), "--root", root); err == nil || !strings.Contains(err.Error(), "load anchors.yaml") {
		t.Errorf("no config: %v", err)
	}
	writeFile(t, root, config.DefaultFile, "version: 1\n")
	if err, _ := runCmd(t, newCodeListCmd(), "--root", root); err == nil || !strings.Contains(err.Error(), "read the map") {
		t.Errorf("no map: %v", err)
	}
	// The per-file lookups degrade to empty instead of failing the listing.
	missing := root + "/nope.yaml"
	if len(kindByFile(missing))+len(needsByFile(missing))+len(parentByFile(missing))+len(revisesByFile(missing)) != 0 {
		t.Error("a missing map must yield empty lookups")
	}
}

func TestUnitNameStripsTheArtifactSuffixes(t *testing.T) {
	for in, want := range map[string]string{
		"src/auth/Login.spec.md":    "Login",
		"src/auth/Login.feature":    "Login",
		"src/auth/Login.test.ts":    "Login",
		"auth/screens/NewLogin.tsx": "NewLogin",
		"Plain":                     "Plain",
	} {
		if got := unitName(in); got != want {
			t.Errorf("unitName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestJoinLens(t *testing.T) {
	for want, in := range map[string][]int{"4": {4}, "4 or 5": {4, 5}, "not declared": nil} {
		if got := joinLens(in); got != want {
			t.Errorf("joinLens(%v) = %q, want %q", in, got, want)
		}
	}
	if errCollision.Error() != "code already in use" {
		t.Errorf("errCollision = %q", errCollision.Error())
	}
}
