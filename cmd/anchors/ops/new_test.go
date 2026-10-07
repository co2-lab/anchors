// @anchors
//   code: NWTSN
//   ref: NWARN

package ops

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// runNew runs `anchors new` UNDER a root, as the CLI does (a root command with the
// `progress` subcommand would read the kind as an unknown subcommand).
func runNew(t *testing.T, args ...string) (error, string) {
	t.Helper()
	root := &cobra.Command{Use: "anchors"}
	root.AddCommand(newNewCmd())
	return runCmd(t, root, append([]string{"new"}, args...)...)
}

func TestNewRefusesWhatItCannotPlace(t *testing.T) {
	t.Run("NWARN-B01: An unknown kind is refused", func(t *testing.T) {})
	t.Run("NWARN-B02: The name and the output path are required", func(t *testing.T) {})
	t.Run("NWARN-B06: With and without are validated against the kind's sections", func(t *testing.T) {})
	t.Run("NWARN-I01: A refused new leaves nothing behind", func(t *testing.T) {})
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
	t.Run("NWARN-B03: A spec is born with a code no unit uses", func(t *testing.T) {})
	t.Run("NWARN-B12: The artifact is written where out says and never overwritten", func(t *testing.T) {})
	root := t.TempDir()
	taken := code.Generate("Login")
	writeFile(t, root, "anchors.graph.yaml", "version: 7\nnodes:\n    - id: x/Login.spec.md\n      kind: spec\n      code: "+taken+"\nedges: []\n")
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
	t.Run("NWARN-B04: A feature or test takes its identity from the sibling spec, or warns it is orphaned", func(t *testing.T) {})
	t.Run("NWARN-X01: A feature references the spec's identity instead of owning one", func(t *testing.T) {})
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
	t.Run("NWARN-B08: A spec for a declarative layer is refused", func(t *testing.T) {})
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
	t.Run("NWARN-B05: Code pins the identity", func(t *testing.T) {})
	t.Run("NWARN-B13: A plan is born with its progress companion", func(t *testing.T) {})
	root := t.TempDir()
	err, out := runNew(t, "plan", "Foundation", "--root", root, "--out", "plans/0001-foundation.md", "--code", "FNDTN")
	if err != nil {
		t.Fatalf("new plan: %v", err)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "plans", "0001-foundation.md")); !strings.Contains(string(b), "code: FNDTN") {
		t.Errorf("--code did not pin the identity:\n%s", b)
	}
	if _, serr := os.Stat(filepath.Join(root, "plans", "0001-foundation-progress.md")); serr != nil {
		t.Errorf("the progress companion was not created: %v\n%s", serr, out)
	}
	if strings.Contains(out, "no spec found") {
		t.Errorf("a plan OWNS its identity and must not get the orphan warning:\n%s", out)
	}
}

func TestNewListSectionsShowsTheMenu(t *testing.T) {
	t.Run("NWARN-B14: List-sections prints the menu, with presets only for specs", func(t *testing.T) {})
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
	t.Run("NWARN-B06: With and without are validated against the kind's sections", func(t *testing.T) {})
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
	t.Run("NWARN-B03: A spec is born with a code no unit uses", func(t *testing.T) {})
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
	t.Run("NWARN-B15: The target layer is the one of the unit the artifact describes", func(t *testing.T) {})
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

// The REGIME tag comes from the project's mapping. Hard-coding `@nivel-unit` produced
// scenarios no gate confronts in a project with another vocabulary — and `anchors work`
// already taught that "the tag is the PROJECT's and is not translatable": the template
// contradicted the ruler.
func TestUnitRegimeTagComesFromTheProject(t *testing.T) {
	t.Run("NWARN-B11: The unit regime tag comes from the project, or a visible placeholder", func(t *testing.T) {})
	cases := []struct {
		name string
		cfg  *config.Config
		want string
	}{
		{"Portuguese vocabulary", &config.Config{Derived: &config.Derived{
			Regimes: map[string]string{"nivel-unit": "unit", "nivel-e2e": "e2e"}}}, "@nivel-unit"},
		{"English vocabulary", &config.Config{Derived: &config.Derived{
			Regimes: map[string]string{"level-unit": "unit"}}}, "@level-unit"},
		{"own vocabulary", &config.Config{Derived: &config.Derived{
			Regimes: map[string]string{"fast": "unit", "slow": "e2e"}}}, "@fast"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := unitRegimeTag(c.cfg); got != c.want {
				t.Fatalf("tag = %q, want %q", got, c.want)
			}
		})
	}
}

// With no mapping declared: a VISIBLE TODO, not a guess. A wrong tag goes unnoticed (no
// gate confronts it); a TODO is fixed on the spot.
func TestUnitRegimeTagWithoutMappingDoesNotGuess(t *testing.T) {
	t.Run("NWARN-B11: The unit regime tag comes from the project, or a visible placeholder", func(t *testing.T) {})
	for name, cfg := range map[string]*config.Config{
		"nil config":     nil,
		"no derived":     {},
		"empty derived":  {Derived: &config.Derived{}},
		"no unit regime": {Derived: &config.Derived{Regimes: map[string]string{"slow": "e2e"}}},
	} {
		t.Run(name, func(t *testing.T) {
			got := unitRegimeTag(cfg)
			if !strings.Contains(got, "TODO") {
				t.Fatalf("guessed a tag instead of asking for the declaration: %q", got)
			}
			if strings.Contains(got, "nivel-unit") {
				t.Fatalf("guessed a specific project's vocabulary: %q", got)
			}
		})
	}
}

// The ORDER of the sections is the spec's reading thread, and it belongs to the preset: in
// a store spec, "State Shape" must come before "Invariants" — an invariant cannot be stated
// about a state the reader does not know yet. Before, the order always came from the
// catalog and came out inverted.
func TestSectionOrderComesFromThePreset(t *testing.T) {
	t.Run("NWARN-B07: A preset fixes the sections and their order, and extra sections follow", func(t *testing.T) {})
	chosen, order, err := resolveSectionsWithPreset(specTemplate, "store", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	var keys []string
	for _, s := range sortSections(specTemplate, chosen, order) {
		keys = append(keys, s.Key)
	}
	pos := func(k string) int {
		for i, c := range keys {
			if c == k {
				return i
			}
		}
		t.Fatalf("section %q was not emitted: %v", k, keys)
		return -1
	}
	if pos("state-shape") > pos("invariants") {
		t.Errorf("State Shape after Invariants — the reader meets the invariant before "+
			"knowing what the state is: %v", keys)
	}
	if pos("actions") < pos("state-shape") {
		t.Errorf("Actions before the Shape: %v", keys)
	}
	if keys[len(keys)-1] != "open" {
		t.Errorf("Open Decisions should close the spec: %v", keys)
	}
}

// A section added by --with that the preset did not foresee comes at the end, in catalog
// order — the only criterion available for it.
func TestSectionOrderWithAnExtraSection(t *testing.T) {
	t.Run("NWARN-B07: A preset fixes the sections and their order, and extra sections follow", func(t *testing.T) {})
	chosen, order, err := resolveSectionsWithPreset(specTemplate, "validation", []string{"auth"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := sortSections(specTemplate, chosen, order)
	if len(got) == 0 || got[len(got)-1].Key != "auth" {
		t.Fatalf("the section asked with --with was not emitted after the preset's: %v", got)
	}
}

// The identity of a feature/test is READ from the sibling spec — they reference (`ref:`),
// they do not own. Generating a new code produces an orphan BY CONSTRUCTION: the `ref:`
// points at a spec that does not exist. It happened: `anchors new feature
// metadataVersioning` minted `MTVA` next to a spec that declared `MTVRX`.
func TestRefComesFromTheSiblingSpec(t *testing.T) {
	t.Run("NWARN-B04: A feature or test takes its identity from the sibling spec, or warns it is orphaned", func(t *testing.T) {})
	dir := t.TempDir()
	spec := filepath.Join(dir, "metadataVersioning.spec.md")
	if err := os.WriteFile(spec, []byte("<!-- @anchors\n  code: MTVRX\n-->\n# x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, file := range map[string]string{
		"feature beside":     "metadataVersioning.feature",
		"test beside":        "metadataVersioning.test.ts",
		"go test beside":     "metadataVersioning_test.go",
		"python test beside": "metadataVersioning_test.py",
	} {
		t.Run(name, func(t *testing.T) {
			got, origin := refDaSpecIrma(dir, filepath.Join(dir, file), "metadataVersioning")
			if got != "MTVRX" {
				t.Fatalf("ref = %q, want MTVRX (the sibling spec is beside it); origin=%q", got, origin)
			}
		})
	}
}

// Without a sibling spec, no identity is invented silently: the caller generates one with
// a warning. The likeliest case is not "the spec comes later" — it is a wrong --out.
func TestRefWithoutASiblingSpecInventsNothing(t *testing.T) {
	t.Run("NWARN-B04: A feature or test takes its identity from the sibling spec, or warns it is orphaned", func(t *testing.T) {})
	dir := t.TempDir()
	if got, _ := refDaSpecIrma(dir, filepath.Join(dir, "noSpec.feature"), "noSpec"); got != "" {
		t.Fatalf("with no sibling spec it should return empty, got %q", got)
	}
	// ANOTHER unit's spec in the same directory is not a sibling.
	os.WriteFile(filepath.Join(dir, "otherThing.spec.md"), []byte("<!-- @anchors\n  code: XXXXX\n-->\n"), 0o644)
	if got, _ := refDaSpecIrma(dir, filepath.Join(dir, "noSpec.feature"), "noSpec"); got != "" {
		t.Fatalf("another unit's spec is not a sibling, got %q", got)
	}
}

func renderScreen(t *testing.T, cfg *config.Config) string {
	t.Helper()
	sections, order, err := resolveSectionsWithPreset(specTemplate, "screen", nil, nil)
	if err != nil {
		t.Fatalf("resolving the screen preset: %v", err)
	}
	return renderArtifact(specTemplate, "Login", "LGNOI", "Login.spec.md", t.TempDir(), sections, order, cfg)
}

// THE TITLE FOLLOWS `lang:`. Before, the catalog emitted Portuguese for every project —
// including `lang: en` and `lang: es` —, contradicting the framework's own default.
func TestSectionTitleFollowsTheProjectLanguage(t *testing.T) {
	t.Run("NWARN-B09: Section titles and bodies follow the project lexicon, then its language", func(t *testing.T) {})
	for _, c := range []struct{ lang, want, forbidden string }{
		{"en", "## Overview", "## Visão Geral"},
		{"es", "## Visión General", "## Visão Geral"},
		{"pt-BR", "## Visão Geral", "## Overview"},
	} {
		out := renderScreen(t, &config.Config{Lang: c.lang})
		if !strings.Contains(out, c.want) {
			t.Errorf("lang=%s: expected %q in the artifact", c.lang, c.want)
		}
		if strings.Contains(out, c.forbidden) {
			t.Errorf("lang=%s: the artifact carries %q, from another language", c.lang, c.forbidden)
		}
	}
}

// The plan was written in Portuguese in every project — body, instructions and the
// `(depende de …)` of its phases —, and an English plan that followed its own language
// wrote `(depends on …)`, which the order gate did not read. Its title section also took
// the SPEC's title body. The whole plan follows `lang:` now, and the dependency phrase it
// writes is one the gate reads (PHORP-B14).
func TestPlanFollowsTheProjectLanguage(t *testing.T) {
	t.Run("NWARN-B09: Section titles and bodies follow the project lexicon, then its language", func(t *testing.T) {})
	all := map[string]bool{}
	var order []string
	for _, s := range planTemplate.sections {
		all[s.Key] = true
		order = append(order, s.Key)
	}
	for _, c := range []struct {
		lang      string
		want      []string
		forbidden []string
	}{
		{"en", []string{"## Objective", "## Rationale", "(depends on PLANX-W01)", "## Out of Scope", "## Definition of Done",
			"## What This Plan Revises", "One step is missing", "> **Code**: `PLANX`"},
			[]string{"depende de", "Objetivo", "Motivo", "Fora de escopo", "Falta um passo", "Código", "purpose in one sentence"}},
		{"es", []string{"## Objetivo", "(depende de PLANX-W01)", "## Fuera de alcance", "Falta un paso, y está en el OTRO archivo"},
			[]string{"depends on", "Fora de escopo", "Falta um passo", "purpose in one sentence"}},
		{"pt-BR", []string{"## Objetivo", "(depende de PLANX-W01)", "## Fora de escopo", "Falta um passo, e ele é no OUTRO arquivo"},
			[]string{"depends on", "Out of Scope", "purpose in one sentence"}},
	} {
		out := renderArtifact(planTemplate, "Foundation", "PLANX", "plans/0001.md", t.TempDir(), all, order, &config.Config{Lang: c.lang})
		for _, w := range c.want {
			if !strings.Contains(out, w) {
				t.Errorf("lang=%s: want %q in:\n%s", c.lang, w, out)
			}
		}
		for _, f := range c.forbidden {
			if strings.Contains(out, f) {
				t.Errorf("lang=%s: %q comes from another language or another artifact", c.lang, f)
			}
		}
	}
}

// With no `lang:` declared the framework DEFAULT holds, which is English — not the
// language of whoever wrote the framework.
func TestWithoutADeclaredLangTheDefaultIsEnglish(t *testing.T) {
	t.Run("NWARN-B09: Section titles and bodies follow the project lexicon, then its language", func(t *testing.T) {})
	if out := renderScreen(t, &config.Config{}); !strings.Contains(out, "## Overview") {
		t.Error("without `lang:` the default is English (i18n.Default), and the artifact should be born in English")
	}
}

// THE PROJECT LEXICON beats the language: whoever declared a title wants it, even in a
// `lang: en` project. The precedence is the target layer's titles > the project's titles >
// rule_types > the language.
func TestProjectLexiconBeatsTheLanguage(t *testing.T) {
	t.Run("NWARN-B09: Section titles and bodies follow the project lexicon, then its language", func(t *testing.T) {})
	cfg := &config.Config{
		Lang:          "en",
		SectionTitles: config.SectionTitles{"overview": "Panorama"},
	}
	out := renderScreen(t, cfg)
	if !strings.Contains(out, "## Panorama") {
		t.Error("`section_titles` should beat the framework's translation")
	}
	if strings.Contains(out, "## Overview") {
		t.Error("with its own lexicon declared, the framework's title should not appear")
	}

	// A rule type names the section of its letter when no title is declared for it.
	cfg = &config.Config{Lang: "en", RuleTypes: []config.RuleType{{Letter: "B", Term: "Behaviour", Sections: []string{"Comportamentos"}}}}
	sections, order, _ := resolveSectionsWithPreset(specTemplate, "hook", nil, nil)
	if out := renderArtifact(specTemplate, "X", "ABCDE", "X.spec.md", t.TempDir(), sections, order, cfg); !strings.Contains(out, "## Comportamentos") {
		t.Errorf("the rule type's section title was not used:\n%s", out)
	}

	// The TARGET layer's own title beats the project's.
	root := t.TempDir()
	cfg = &config.Config{Lang: "en", SectionTitles: config.SectionTitles{"overview": "Panorama"},
		Layers: map[string]config.Layer{"api": {Pattern: "api/**/*.go", Kind: "code",
			SectionTitles: config.SectionTitles{"overview": "Resumo"}}}}
	if out := renderArtifact(specTemplate, "X", "ABCDE", filepath.Join(root, "api", "X.spec.md"), root,
		map[string]bool{"overview": true}, nil, cfg); !strings.Contains(out, "## Resumo") {
		t.Errorf("the target layer's title was not used:\n%s", out)
	}
}

// Every section of the spec catalog has a translation: a missing key makes the section come
// out in Portuguese inside an English artifact, unnoticed because the rest of the file is
// right. It happened with `a11y`.
func TestEverySpecSectionHasATranslation(t *testing.T) {
	t.Run("NWARN-B09: Section titles and bodies follow the project lexicon, then its language", func(t *testing.T) {})
	out := renderScreen(t, &config.Config{Lang: "en"})
	// The BODY is translated too, not only the title: a table header and a TODO hint came
	// out in Portuguese inside an English artifact, unnoticed because the title beside them
	// was right.
	for _, pt := range []string{"TODO propósito", "**Código**", "Obrigatório", "Quem garante", "nenhuma"} {
		if strings.Contains(out, pt) {
			t.Errorf("the `lang: en` artifact carries %q in the BODY — a section.body.* key is missing", pt)
		}
	}
	// No title of the English artifact may be in Portuguese.
	for _, pt := range []string{"Visão Geral", "Regras", "Dependências", "Acessibilidade", "Decisões em aberto"} {
		if strings.Contains(out, "## "+pt) {
			t.Errorf("the `lang: en` artifact carries the section %q in Portuguese — a section.title.* key is missing", pt)
		}
	}
}

// A feature speaks the project's Gherkin dialect, and a test is born in its stack's syntax:
// both come from the project's declared dialect.
func TestNewFollowsTheProjectDialect(t *testing.T) {
	t.Run("NWARN-B10: Features and tests follow the project's declared dialect", func(t *testing.T) {})
	cfg := &config.Config{Dialect: &config.Dialect{GherkinLanguage: "pt", Family: "python"}}
	ch, _ := resolveSections(featureTemplate, nil, nil)
	feat := renderArtifact(featureTemplate, "Login", "LGNSP", "/x/Login.feature", t.TempDir(), ch, nil, cfg)
	if !strings.HasPrefix(feat, "# language: pt\n") || !strings.Contains(feat, "Funcionalidade:") ||
		strings.Contains(feat, "Feature:") {
		t.Errorf("the feature does not follow the pt Gherkin dialect:\n%s", feat)
	}
	ch, _ = resolveSections(testTemplate, nil, nil)
	test := renderArtifact(testTemplate, "calcTotal", "LGNSP", "/x/calc_test.py", t.TempDir(), ch, nil, cfg)
	if !strings.Contains(test, "def test_calc_total(") || strings.Contains(test, "describe(") {
		t.Errorf("the test does not follow the python family:\n%s", test)
	}
}

// The help said `--out` defaulted to the root (it is mandatory, and the root is refused)
// and named 3 of the 7 kinds the catalog gives birth to; the unknown-kind refusal did too.
func TestNewHelpTellsTheTruth(t *testing.T) {
	t.Run("NWARN-B16: The help and the unknown-kind refusal name every kind, and --out is mandatory", func(t *testing.T) {})
	cmd := newNewCmd()
	usage := cmd.Flags().Lookup("out").Usage
	if strings.Contains(usage, "default") || !strings.Contains(usage, "mandatory") {
		t.Errorf("--out usage = %q; it is mandatory, with no default", usage)
	}
	err, _ := runNew(t, "widget", "X")
	for kind := range templates {
		if !strings.Contains(cmd.Long, kind) {
			t.Errorf("the help does not name the kind %q:\n%s", kind, cmd.Long)
		}
		if err == nil || !strings.Contains(err.Error(), kind) {
			t.Errorf("the unknown-kind refusal does not name %q: %v", kind, err)
		}
	}
}

func TestNewArtifactEntersTheMap(t *testing.T) {
	t.Run("NWARN-B17: The new artifact enters the map at once", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", "version: 6\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n")
	writeFile(t, root, "anchors.graph.yaml", "version: 7\nnodes: []\nedges: []\n")
	err, out := runNew(t, "spec", "Pay", "--root", root, "--out", "src/Pay.spec.md")
	if err != nil {
		t.Fatalf("new spec: %v", err)
	}
	g, lerr := mapx.Load(filepath.Join(root, "anchors.graph.yaml"))
	if lerr != nil {
		t.Fatal(lerr)
	}
	if len(g.Nodes) != 1 || g.Nodes[0].ID != "src/Pay.spec.md" || !strings.Contains(out, "added to the map: src/Pay.spec.md") {
		t.Fatalf("the spec enters the map at once, got %+v\n%s", g.Nodes, out)
	}
}

func TestNewRefArtifactIsBornWithItsOwnCode(t *testing.T) {
	t.Run("NWARN-B18: A new artifact that refs its unit is born with a code of its own", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", "version: 7\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n  feature:\n    pattern: \"**/*.feature\"\n    kind: feature\n")
	writeFile(t, root, "anchors.graph.yaml", "version: 7\nnodes:\n  - id: src/Pay.spec.md\n    kind: spec\n    code: PAYMT\n    file_code: PAYMT\nedges: []\n")
	writeFile(t, root, "src/Pay.spec.md", "<!-- @anchors\n  code: PAYMT\n-->\n# Pay\n")
	if err, out := runNew(t, "feature", "Pay", "--root", root, "--out", "src/Pay.feature"); err != nil {
		t.Fatalf("new feature: %v\n%s", err, out)
	}
	b, _ := os.ReadFile(filepath.Join(root, "src/Pay.feature"))
	m := regexp.MustCompile(`(?m)^#   code: ([A-Z0-9]{5})$`).FindStringSubmatch(string(b))
	if m == nil || m[1] == "PAYMT" || !strings.Contains(string(b), "#   ref: PAYMT") {
		t.Errorf("the feature refs its unit and carries a code of its own:\n%s", b)
	}
}

func TestOwnSectionTitleSkipsASectionEveryLetterLists(t *testing.T) {
	t.Run("NWARN-B19: A rule section takes the title its letter lists alone, not one every letter shares", func(t *testing.T) {})
	types := []config.RuleType{
		{Letter: "E", Sections: []string{"Uso das regras", "Efeitos"}},
		{Letter: "X", Sections: []string{"Uso das regras", "Restrições"}},
		{Letter: "B", Sections: []string{"Comportamentos", "Uso das regras"}},
		{Letter: "U", Sections: []string{"Uso das regras"}},
	}
	for letter, want := range map[string]string{"E": "Efeitos", "X": "Restrições", "B": "Comportamentos", "U": "Uso das regras", "Z": ""} {
		if got := ownSectionTitle(types, letter); got != want {
			t.Errorf("letter %s: got %q, want %q", letter, got, want)
		}
	}
}

func TestScreenPresetWritesTheLoadStatesAndTheirFailure(t *testing.T) {
	t.Run("NWARN-B20: The screen preset writes the four states of a unit that loads data and the failure of its load, under the titles of States and Errors", func(t *testing.T) {})
	chosen, order, err := resolveSectionsWithPreset(specTemplate, "screen", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	out := renderArtifact(specTemplate, "Probe", "PRBOE", "Probe.spec.md", t.TempDir(), chosen, order, &config.Config{Lang: "en"})
	for _, want := range []string{"## States", "PRBOE-S01 — Loading", "PRBOE-S02 — Empty", "PRBOE-S03 — Load error", "PRBOE-S04 — Loaded", "## Errors / Failures", "`PRBOE-E01`"} {
		if !strings.Contains(out, want) {
			t.Errorf("the screen spec lacks %q", want)
		}
	}
	renamed := renderArtifact(specTemplate, "Probe", "PRBOE", "Probe.spec.md", t.TempDir(), chosen, order,
		&config.Config{Lang: "en", SectionTitles: config.SectionTitles{"states": "Screen States"}})
	if !strings.Contains(renamed, "## Screen States") {
		t.Error("the load states take the title the project gave its States section")
	}
}

func TestScreenPresetKeepsTheLetterInOneHome(t *testing.T) {
	t.Run("NWARN-B21: Two sections of one letter never define the same code, and a section of the rules' letter with no title of its own goes inside the rules section the layer names", func(t *testing.T) {})
	chosen, order, err := resolveSectionsWithPreset(specTemplate, "screen", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// jokenpo's shape: the letter B's first own title is "Efeitos", and the screen layer
	// names its rules section "Rules (Regras de Negócio)".
	cfg := &config.Config{Lang: "pt-BR",
		RuleTypes:     []config.RuleType{{Letter: "B", Sections: []string{"Uso das regras", "Efeitos", "Comportamentos"}}, {Letter: "U", Sections: []string{"Uso das regras"}}},
		SectionTitles: config.SectionTitles{"rules": "Rules (Regras de Negócio)"}}
	out := renderArtifact(specTemplate, "Probe", "PRBOE", "Probe.spec.md", t.TempDir(), chosen, order, cfg)
	if strings.Contains(out, "## Efeitos") || strings.Count(out, "## Rules (Regras de Negócio)") != 1 {
		t.Errorf("one home for the letter's rules, not a second section:\n%s", out)
	}
	rules := out[strings.Index(out, "## Rules (Regras de Negócio)"):]
	rules = rules[:strings.Index(rules[3:], "\n## ")+3]
	if !strings.Contains(rules, "### PRBOE-B01") || !strings.Contains(rules, "### Carregamento") || !strings.Contains(rules, "`PRBOE-B02`") {
		t.Errorf("the loading table goes inside the rules section with the next code:\n%s", rules)
	}
	if n := strings.Count(out, "### PRBOE-B01") + strings.Count(out, "| `PRBOE-B01` | TODO: tudo"); n != 1 {
		t.Errorf("PRBOE-B01 is defined once, got %d:\n%s", n, out)
	}
	plain := renderArtifact(specTemplate, "Probe", "PRBOE", "Probe.spec.md", t.TempDir(), chosen, order, &config.Config{Lang: "en"})
	if !strings.Contains(plain, "`PRBOE-B02` | TODO") || strings.Count(plain, "PRBOE-B01 —") != 1 {
		t.Errorf("with no title declared the sections stay apart, their codes still distinct:\n%s", plain)
	}
}
