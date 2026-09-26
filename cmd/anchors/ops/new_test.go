package ops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
)

// Product doctrine does NOT inherit the spec's body.
//
// `sectionBody` resolves the body through `section.body.<key>` in the translation
// catalog, and the template's literal only counts when the key is absent there. With bare
// keys (`rules`, `overview`) the doctrine silently inherited the SPEC's text — measured:
// the command emitted `### XXXXX-B01 — TODO rule` and "what the unit does", spec text in
// an artifact that has no unit at all.
func TestNewProduct_doesNotInheritTheSpecBody(t *testing.T) {
	tpl, ok := templates["product"]
	if !ok {
		t.Fatal("the `product` kind is not registered")
	}
	for _, s := range tpl.sections {
		if !strings.HasPrefix(s.Key, "doctrine_") {
			t.Errorf("section %q without the `doctrine_` prefix: it will inherit another kind's body", s.Key)
		}
	}
	body := renderArtifact(tpl, "CreditLimit", "CRLMT", "product/x.doctrine.md", t.TempDir(),
		map[string]bool{"doctrine_title": true, "doctrine_rules": true}, nil, nil)
	if !strings.Contains(body, "CRLMT-R01") {
		t.Errorf("the doctrine should emit `-R01` rules, got:\n%s", body)
	}
	if strings.Contains(body, "-B01 — TODO rule") {
		t.Errorf("the doctrine inherited the spec body:\n%s", body)
	}
	// Doctrine has no TARGET: the layer belongs to the target a spec describes, and a
	// `layer: TODO` here would be a field nobody can fill — the eternal placeholder that
	// `placeholder-filled` exists to accuse.
	if strings.Contains(tpl.headerFn("CRLMT", "product/x.doctrine.md"), "layer:") {
		t.Error("the doctrine header must not declare `layer:`")
	}
}

// runNew runs `anchors new` UNDER a root, as the CLI does (a root command with the
// `progress` subcommand would read the kind as an unknown subcommand).
func runNew(t *testing.T, args ...string) (error, string) {
	t.Helper()
	root := &cobra.Command{Use: "anchors"}
	root.AddCommand(newNewCmd())
	return runCmd(t, root, append([]string{"new"}, args...)...)
}

func TestNewRefusesWhatItCannotPlace(t *testing.T) {
	root := t.TempDir()
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"widget", "X"}, `unknown kind "widget"`},
		{[]string{"spec"}, "provide the unit name"},
		{[]string{"spec", "Login", "--root", root}, "--out"},
		{[]string{"spec", "Login", "--root", root, "--out", "a/Login.spec.md", "--with", "nope"}, `unknown section "nope"`},
		{[]string{"feature", "Login", "--root", root, "--out", "a/Login.feature", "--preset", "store"}, "--preset only applies to `spec`"},
		{[]string{"spec", "Login", "--root", root, "--out", "a/Login.spec.md", "--preset", "nope"}, `unknown preset "nope"`},
	} {
		err, _ := runNew(t, c.args...)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("new %v: got %v, want %q", c.args, err, c.want)
		}
	}
	if entries, _ := os.ReadDir(root); len(entries) != 0 {
		t.Errorf("a refused `new` left files behind: %v", entries)
	}
}

// A spec is born where --out says, with a code no one uses, and it is never overwritten.
func TestNewSpecIsBornWithAFreeCode(t *testing.T) {
	root := t.TempDir()
	taken := code.Generate("Login")
	writeFile(t, root, "anchors.graph.yaml", "version: 4\nnodes:\n    - id: x/Login.spec.md\n      kind: spec\n      code: "+taken+"\nedges: []\n")
	err, out := runNew(t, "spec", "Login", "--root", root, "--out", "src/auth/Login.spec.md", "--with", "errors", "--without", "overview")
	if err != nil {
		t.Fatalf("new spec: %v", err)
	}
	b, rerr := os.ReadFile(filepath.Join(root, "src", "auth", "Login.spec.md"))
	if rerr != nil {
		t.Fatal(rerr)
	}
	got := codeDoHeaderSpec(string(b))
	if got == "" || got == taken {
		t.Errorf("header code = %q; want a code other than the taken %s", got, taken)
	}
	if !strings.Contains(out, "(code: "+got+")") || !strings.Contains(out, "anchors check --changed src/auth/Login.spec.md") {
		t.Errorf("output:\n%s", out)
	}
	if err, _ := runNew(t, "spec", "Login", "--root", root, "--out", "src/auth/Login.spec.md"); err == nil ||
		!strings.Contains(err.Error(), "already exists") {
		t.Errorf("a second `new` onto the same path must refuse, got %v", err)
	}
}

// A feature REFERENCES its sibling spec's code; without a spec it warns that it is born
// orphaned instead of silently minting an identity.
func TestNewFeatureReferencesTheSiblingSpec(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "src/auth/Login.spec.md", "<!-- @anchors\ncode: LGNSP\n-->\n# Login\n")
	err, out := runNew(t, "feature", "Login", "--root", root, "--out", "src/auth/Login.feature")
	if err != nil {
		t.Fatalf("new feature: %v", err)
	}
	if !strings.Contains(out, "identity: `LGNSP` (read from src/auth/Login.spec.md)") {
		t.Errorf("the identity did not come from the spec:\n%s", out)
	}
	b, _ := os.ReadFile(filepath.Join(root, "src", "auth", "Login.feature"))
	if !strings.Contains(string(b), "LGNSP") {
		t.Errorf("the feature does not reference LGNSP:\n%s", b)
	}

	err, out = runNew(t, "feature", "Orphan", "--root", root, "--out", "src/other/Orphan.feature")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no spec found for this target") || !strings.Contains(out, "ORPHANED") {
		t.Errorf("the orphan is not warned:\n%s", out)
	}
}

// A spec for a RECOGNIZED (declarative) layer is refused before the file exists.
func TestNewSpecRefusesADeclarativeLayer(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", "layers:\n  dao:\n    pattern: \"src/dao/**/*.ts\"\n    kind: code\n    regime: declarativo\n")
	err, _ := runNew(t, "spec", "UserDao", "--root", root, "--out", "src/dao/UserDao.spec.md", "--code", "USRDA")
	if err == nil || !strings.Contains(err.Error(), "`dao` is a RECOGNIZED layer") {
		t.Fatalf("want the declarative-layer refusal, got %v", err)
	}
	if _, serr := os.Stat(filepath.Join(root, "src", "dao", "UserDao.spec.md")); serr == nil {
		t.Error("the refused spec was written anyway")
	}
	// The guard is about specs: a feature there is not its business.
	if refuseIfRecognizedLayer(root, filepath.Join(root, "src/dao/UserDao.feature"), nil, "feature") != nil {
		t.Error("the guard must only act on specs")
	}
}

// A plan is born with its progress companion — the state lives there, not in the plan.
func TestNewPlanIsBornWithItsProgress(t *testing.T) {
	root := t.TempDir()
	err, out := runNew(t, "plan", "Foundation", "--root", root, "--out", "plans/0001-foundation.md", "--code", "FNDTN")
	if err != nil {
		t.Fatalf("new plan: %v", err)
	}
	if _, serr := os.Stat(filepath.Join(root, "plans", "0001-foundation-progress.md")); serr != nil {
		t.Errorf("the progress companion was not created: %v\n%s", serr, out)
	}
	if strings.Contains(out, "no spec found") {
		t.Errorf("a plan OWNS its identity and must not get the orphan warning:\n%s", out)
	}
}

func TestNewListSectionsShowsTheMenu(t *testing.T) {
	err, out := runNew(t, "spec", "--list-sections")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Sections of `spec`", "[default ] title", "[optional] errors",
		"Presets (ready-made sets", "store"} {
		if !strings.Contains(out, want) {
			t.Errorf("the menu misses %q:\n%s", want, out)
		}
	}
	_, out = runNew(t, "feature", "--list-sections")
	if strings.Contains(out, "Presets") {
		t.Errorf("presets are a spec thing; the feature menu lists them:\n%s", out)
	}
}

// resolveSections starts from the defaults, then applies --with and --without.
func TestResolveSectionsAppliesWithAndWithout(t *testing.T) {
	chosen, err := resolveSections(specTemplate, []string{"errors"}, []string{"overview"})
	if err != nil {
		t.Fatal(err)
	}
	if !chosen["errors"] || chosen["overview"] || !chosen["title"] || chosen["route"] {
		t.Errorf("chosen = %v", chosen)
	}
}

// The code of a path in a layer with `code_prefix` starts with that prefix; a plain name
// gets the canonical code, avoiding the taken ones.
func TestResolveNewCode(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", "layers:\n  auth:\n    pattern: \"src/auth/**/*.ts\"\n    kind: code\n    code_prefix: AU\n")
	writeFile(t, root, "src/auth/Session.ts", "")
	got, err := resolveNewCode(root, "src/auth/Session.ts")
	if err != nil || !strings.HasPrefix(got, "AU") {
		t.Errorf("resolveNewCode(path) = %q, %v; want the AU prefix", got, err)
	}
	if got, _ := resolveNewCode(root, "src/misc/index.ts"); got == code.Generate("index") {
		t.Errorf("a generic basename outside a prefixed layer must use the parent dir, got %q", got)
	}
	if got, _ := resolveNewCode(root, "Wallet"); got != code.Generate("Wallet") {
		t.Errorf("resolveNewCode(name) = %q, want %q", got, code.Generate("Wallet"))
	}
}

// The target layer is the one of the unit the spec describes: the file that EXISTS wins;
// before the code exists, the most specific extension wins.
func TestTargetLayerPrefersTheExistingTarget(t *testing.T) {
	root := t.TempDir()
	cfg := &config.Config{Layers: map[string]config.Layer{
		"screen":   {Pattern: "app/**/*.tsx", Kind: "code"},
		"catchall": {Pattern: "app/**/*.ts", Kind: "code"},
	}}
	if got := targetLayer(root, filepath.Join(root, "app/P.spec.md"), cfg); got != "screen" {
		t.Errorf("no target on disk: got %q, want screen (.tsx before .ts)", got)
	}
	writeFile(t, root, "app/Q.ts", "")
	if got := targetLayer(root, filepath.Join(root, "app/Q.spec.md"), cfg); got != "catchall" {
		t.Errorf("Q.ts exists: got %q, want catchall", got)
	}
	if got := targetLayer(root, filepath.Join(root, "app/x.ts"), cfg); got != "catchall" {
		t.Errorf("a non-derived path classifies itself: got %q", got)
	}
	if targetLayer(root, "x", nil) != "" {
		t.Error("without config there is no layer")
	}
}
