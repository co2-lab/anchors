package ops

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
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
	writeFile(t, root, mapx.DefaultPath, "version: 5\nnodes:\n"+nodes+"edges: []\n")
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
	t.Run("CDCMC-B05: Check answers whether a code is free, ignoring case, and names the owners", func(t *testing.T) {})
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
	t.Run("CDCMC-B01: A free canonical code is the suggestion", func(t *testing.T) {})
	t.Run("CDCMC-B02: A taken canonical is adjusted to a free code naming its owner", func(t *testing.T) {})
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
	t.Run("CDCMC-B03: A path in a layer with a code prefix gets the module prefix", func(t *testing.T) {})
	t.Run("CDCMC-B04: A generic basename takes its identity from the parent directory", func(t *testing.T) {})
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
	t.Run("CDCMC-E01: Without a name or a map the command fails and says what to do", func(t *testing.T) {})
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
	t.Run("CDCMC-B10: The length check accuses only declared codes and proposes the canonical code", func(t *testing.T) {})
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
	t.Run("CDCMC-B08: The JSON list carries each code's folder, file, kind, title and work order fields", func(t *testing.T) {})
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
	t.Run("CDCMC-B12: An empty map says no node has an identity", func(t *testing.T) {})
	t.Run("CDCMC-E02: The list refuses a project without config or without map", func(t *testing.T) {})
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
	t.Run("CDCMC-B11: The unit name drops the artifact suffixes", func(t *testing.T) {})
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

// mapWithCodes writes a minimal map with identity nodes in different workspaces.
func mapWithCodes(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "anchors.graph.yaml")
	const y = `version: 5
nodes:
    - id: apps/mobile/src/features/auth/LoginScreen.spec.md
      kind: spec
      code: LOGI
    - id: apps/mobile/src/features/auth/LoginScreen.tsx
      kind: code
      code: LOGI
    - id: packages/backend/services/security.spec.md
      kind: spec
      code: SGSB
    - id: apps/mobile/src/components/atoms/Button.spec.md
      kind: spec
      code: BTTN
edges: []
`
	if err := os.WriteFile(p, []byte(y), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// runList runs `code list` on a map and returns only what reaches STDOUT.
func runList(t *testing.T, mapPath string, args ...string) string {
	t.Helper()
	cmd := newCodeListCmd()
	var out strings.Builder
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(append([]string{"--map", mapPath}, args...))
	// The command prints to os.Stdout; capture it through a pipe.
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := cmd.Execute()
	w.Close()
	os.Stdout = orig
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	return string(buf[:n])
}

func TestCodeListEnumeratesFromTheMap(t *testing.T) {
	t.Run("CDCMC-B06: The list prints one sorted code per line with its folder, the summary kept off stdout", func(t *testing.T) {})
	t.Run("CDCMC-X01: The codes come from the map's identity field", func(t *testing.T) {})
	// Why the command exists: the alternative is grepping for a code pattern, which matches
	// a mention in prose, a comment and a file name — and depends on the length of the code,
	// which varies per project (`code_lengths`). Here the source is the map's structured field.
	out := runList(t, mapWithCodes(t))
	for _, c := range []string{"LOGI\tapps/mobile/src/features/auth", "SGSB\tpackages/backend/services", "BTTN"} {
		if !strings.Contains(out, c) {
			t.Errorf("the code line %q should appear, got:\n%s", c, out)
		}
	}
	// Sorted: whoever reads looks for a code, and an unsorted list forces looking twice.
	if i, j := strings.Index(out, "BTTN"), strings.Index(out, "LOGI"); i > j {
		t.Errorf("the output must be sorted by code, got:\n%s", out)
	}
	// The summary goes to stderr, so the list pipes clean.
	if strings.Contains(out, "code(s) in use") {
		t.Errorf("the summary reached stdout:\n%s", out)
	}
}

func TestCodeListOneLinePerCodeEvenWithSeveralFiles(t *testing.T) {
	t.Run("CDCMC-B06: The list prints one sorted code per line with its folder, the summary kept off stdout", func(t *testing.T) {})
	// LOGI is in the spec AND in the .tsx — ONE unit, not two. Repeating the line would make
	// the count lie about how many identities exist.
	out := runList(t, mapWithCodes(t))
	if n := strings.Count(out, "LOGI"); n != 1 {
		t.Errorf("LOGI should appear on ONE line (same folder), appeared %d× in:\n%s", n, out)
	}
}

func TestCodeListFiltersByWorkspace(t *testing.T) {
	t.Run("CDCMC-B07: The list filters by path prefix and names a filter that matched nothing", func(t *testing.T) {})
	// The real question in a monorepo: "the codes of THIS workspace".
	out := runList(t, mapWithCodes(t), "--in", "packages/backend")
	if !strings.Contains(out, "SGSB") {
		t.Errorf("SGSB is under packages/backend and should appear, got:\n%s", out)
	}
	for _, outside := range []string{"LOGI", "BTTN"} {
		if strings.Contains(out, outside) {
			t.Errorf("%s is outside the filter and should not appear, got:\n%s", outside, out)
		}
	}
}

func TestCodeListFilterWithNoResultDoesNotClaimAnEmptyProject(t *testing.T) {
	t.Run("CDCMC-B07: The list filters by path prefix and names a filter that matched nothing", func(t *testing.T) {})
	// "no code under X" differs from "the project has no code" — the second message would
	// send the user to run `map build` for nothing.
	out := runList(t, mapWithCodes(t), "--in", "apps/web")
	if !strings.Contains(out, `no code in use under "apps/web"`) {
		t.Errorf("the message must cite the filter that did not match, got:\n%s", out)
	}
}

// The artifact's title NAMES the work in a card, so the part that repeats the kind and the
// number ("Plano 0001 — ") is dropped: the consumer already has the kind and the code, and
// the repetition only makes the text grow.
func TestFileTitleDropsTheRedundantPrefix(t *testing.T) {
	t.Run("CDCMC-B09: The title drops the text before the dash", func(t *testing.T) {})
	dir := t.TempDir()
	cases := map[string]string{
		"# Login\n":                      "Login",
		"# Plano 0001 — Fundação\n":      "Fundação",
		"# Spec 0042 — Recuperação\n":    "Recuperação",
		"no title\n":                     "",
		"<!-- @anchors -->\n\n# After\n": "After",
	}
	for content, want := range cases {
		if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := fileTitle(dir, "a.md"); got != want {
			t.Errorf("%q → %q, want %q", content, got, want)
		}
	}
	// A file that is not markdown has no title to extract.
	if got := fileTitle(dir, "x.go"); got != "" {
		t.Errorf("a non-markdown file returned %q", got)
	}
}

// The file and declared-identity indexes were package-level maps that takenCodes only
// ever added to: a second run in the same process (the MCP server, a test binary) still
// saw the first map's declared WLTX and accused a code the second map only CITES.
func TestCodeListCheckForgetsThePreviousMap(t *testing.T) {
	t.Run("CDCMC-B13: A second check in the same process judges only its own map", func(t *testing.T) {})
	first := codeProject(t, "version: 1\n", `    - id: app/Wallet.spec.md
      kind: spec
      code: WLTX
      code_declared: true
`)
	if err, out := runCmd(t, newCodeListCmd(), "--root", first, "--check"); !errors.Is(err, errCollision) {
		t.Fatalf("the first map has a declared divergent code: %v\n%s", err, out)
	}
	second := codeProject(t, "version: 1\n", `    - id: app/fixture_test.go
      kind: test
      code: WLTX
`)
	err, out := runCmd(t, newCodeListCmd(), "--root", second, "--check")
	if err != nil || strings.Contains(out, "WLTX →") || !strings.Contains(out, "1 code(s) only CITED") {
		t.Errorf("the second map only cites WLTX; got %v:\n%s", err, out)
	}
	// And the JSON's file comes from the map being read, not a previous one.
	err, out = runCmd(t, newCodeListCmd(), "--root", second, "--json")
	if err != nil || !strings.Contains(out, `"arquivo": "app/fixture_test.go"`) {
		t.Errorf("the file must come from the second map; got %v:\n%s", err, out)
	}
}

// `--fix` was declared and never read, while the output told people to run it: the
// command did nothing and said it would. The fix is one `anchors recode` per divergence,
// reviewed in its dry-run — so the hint names that, and the flag is gone.
func TestCodeListCheckPointsToRecodeNotToAMissingFix(t *testing.T) {
	t.Run("CDCMC-B14: The length check points each divergence to anchors recode, and there is no --fix", func(t *testing.T) {})
	if newCodeListCmd().Flags().Lookup("fix") != nil {
		t.Error("--fix is declared but does nothing")
	}
	root := codeProject(t, "version: 1\n", `    - id: app/Wallet.spec.md
      kind: spec
      code: WLTX
      code_declared: true
`)
	_, out := runCmd(t, newCodeListCmd(), "--root", root, "--check")
	if strings.Contains(out, "--fix") {
		t.Errorf("the output still recommends --fix:\n%s", out)
	}
	if want := "anchors recode WLTX " + code.Generate("Wallet"); !strings.Contains(out, want) {
		t.Errorf("want the exact recode command %q in:\n%s", want, out)
	}
}
