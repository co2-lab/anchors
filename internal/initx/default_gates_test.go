package initx

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/migra"
)

func gateNames(gates []config.Gate) map[string]bool {
	names := map[string]bool{}
	for _, g := range gates {
		names[g.Name] = true
	}
	return names
}

func TestDefaultGates(t *testing.T) {
	t.Run("DFGTD-B02: A project with spec, feature and test is born with the gates of each", func(t *testing.T) {})
	// a project with spec+feature+test → test gates present
	g := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true}, false)
	names := gateNames(g)
	for _, x := range g {
		if x.IsBlocking() && !blockingByNature[x.Name] {
			t.Errorf("gate %s should be born informative", x.Name)
		}
	}
	for _, want := range []string{"spec-complete", "feature-not-empty", "tests-green", "line-coverage", "scenario-coverage"} {
		if !names[want] {
			t.Errorf("default gate %q missing", want)
		}
	}
}

// blockingByNature are the two gates born blocking whatever the age (DFGTD-B06).
var blockingByNature = map[string]bool{
	"no-secret-leaked": true, "doc-required": true,
}

func TestDefaultGatesNoScenarioWithoutSpec(t *testing.T) {
	t.Run("DFGTD-B03: The gates that cross spec and feature are seeded only when both are chosen", func(t *testing.T) {})
	for _, chosen := range []map[string]bool{
		{"test": true},
		{"spec": true, "test": true},
		{"feature": true, "test": true},
	} {
		names := gateNames(DefaultGates(chosen, false))
		if names["scenario-coverage"] || names["spec-feature-match"] {
			t.Errorf("%v: scenario-coverage/spec-feature-match should not exist without spec and feature", chosen)
		}
	}
	names := gateNames(DefaultGates(map[string]bool{"spec": true, "feature": true}, false))
	if !names["scenario-coverage"] || !names["spec-feature-match"] {
		t.Error("with spec and feature both chosen, scenario-coverage and spec-feature-match are seeded")
	}
	// The way back of each pair needs both ends too.
	for _, c := range []struct {
		chosen map[string]bool
		gate   string
		want   bool
	}{
		{map[string]bool{"feature": true}, "feature-spec-match", false},
		{map[string]bool{"spec": true}, "feature-spec-match", false},
		{map[string]bool{"spec": true, "feature": true}, "feature-spec-match", true},
		{map[string]bool{"test": true}, "test-feature-match", false},
		{map[string]bool{"feature": true}, "test-feature-match", false},
		{map[string]bool{"feature": true, "test": true}, "test-feature-match", true},
	} {
		if got := gateNames(DefaultGates(c.chosen, false))[c.gate]; got != c.want {
			t.Errorf("%v: %s seeded = %v, want %v", c.chosen, c.gate, got, c.want)
		}
	}
}

// `anchors init` offers ArtifactNames; choosing every one of them must seed the whole
// catalog. Before `code` was an option, nine canonical gates — no-secret-leaked among
// them — could never be seeded by any answer.
func TestEveryOfferedArtifactSeedsTheFullCatalog(t *testing.T) {
	t.Run("DFGTD-B13: Choosing every artifact init offers seeds every gate of the catalog", func(t *testing.T) {})
	offered := map[string]bool{}
	for _, n := range ArtifactNames() {
		offered[n] = true
	}
	seeded := gateNames(DefaultGates(offered, true))
	for _, g := range canonicalCatalog() {
		if !seeded[g.Name] {
			t.Errorf("gate %q is in the catalog but no init answer seeds it", g.Name)
		}
	}
	if !seeded["no-secret-leaked"] || !seeded["layer-boundary"] {
		t.Error("the code gates must be seeded when every offered artifact is chosen")
	}
}

// The gates over specs, features and tests used to sit inside the `plan` block: a project
// with the unit and no plans was born without them.
func TestUnitGatesAreSeededWithoutPlans(t *testing.T) {
	t.Run("DFGTD-B14: The gates that run on specs, features and tests are seeded without plans", func(t *testing.T) {})
	unit := map[string]bool{"spec": true, "feature": true, "test": true}
	gates := DefaultGates(unit, false)
	names := gateNames(gates)
	for _, want := range []string{
		"docs-fresh", "doctrine-realized", "spec-doctrine-exists", "feature-spec-match",
		"test-feature-match", "revision-orphans", "flag-scenario-grammar",
		"flag-scenarios-complete", "flag-scenario-exists", "flag-scenario-governs",
		"flag-covered", "doctrine-not-duplicated", "spec-realizes-doctrine",
		"failure-handled", "failure-logged", "failure-declared", "docs-covered",
		"doc-required", "doc-self-contained", "plan-change-justified",
	} {
		if !names[want] {
			t.Errorf("%q runs on the unit and must be seeded without plans", want)
		}
	}
	for _, planOnly := range []string{"plan-seeds-valid", "plan-source-declared", "plan-doctrine-exists", "plan-revised", "phase-ordered"} {
		if names[planOnly] {
			t.Errorf("%q runs on plans only and must not be seeded without them", planOnly)
		}
	}
	for _, g := range gates {
		if g.Name == "plan-change-justified" && !reflect.DeepEqual(g.On, []string{"spec"}) {
			t.Errorf("without plans, plan-change-justified runs on [spec], got %v", g.On)
		}
	}
}

// What choosing plans adds is exactly what runs on plans — nothing over the other
// artifacts hides behind the plan choice.
func TestChoosingPlansAddsOnlyPlanGates(t *testing.T) {
	t.Run("DFGTD-I04: Choosing plans adds only gates that run on plans", func(t *testing.T) {})
	for _, base := range []map[string]bool{
		{"spec": true, "feature": true, "test": true},
		{"spec": true, "feature": true, "test": true, "code": true, "guide": true},
	} {
		without := gateNames(DefaultGates(base, false))
		withPlan := map[string]bool{"plan": true}
		for k := range base {
			withPlan[k] = true
		}
		for _, g := range DefaultGates(withPlan, false) {
			if without[g.Name] {
				continue
			}
			if !slices.Contains(g.On, "plan") {
				t.Errorf("%v: choosing plans added %q, which runs on %v", base, g.Name, g.On)
			}
		}
	}
}

// The vocabulary check confronts the migration's renames against the registered names;
// with only the unit's gates registered, a rename to a code, plan or guide gate read as
// a rename to nothing.
func TestRegisteredGateNamesAreTheFullCatalog(t *testing.T) {
	t.Run("DFGTD-B15: The gate names registered for the vocabulary check are the full catalog", func(t *testing.T) {})
	var want []string
	for _, g := range canonicalCatalog() {
		want = append(want, g.Name)
	}
	if got := config.DefaultGateNamesForTest(); !reflect.DeepEqual(got, want) {
		t.Errorf("registered names = %v\nwant the full catalog = %v", got, want)
	}
}

func TestDefaultGatesEmpty(t *testing.T) {
	t.Run("DFGTD-B01: A project with no artifact chosen is seeded with no gate", func(t *testing.T) {})
	if len(DefaultGates(map[string]bool{}, false)) != 0 {
		t.Error("a project with no artifacts should not be seeded with gates")
	}
	if len(DefaultGates(map[string]bool{}, true)) != 0 {
		t.Error("a new project with no artifacts should not be seeded with gates either")
	}
}

func TestDefaultGatesGuideOnly(t *testing.T) {
	t.Run("DFGTD-B04: Choosing only guides seeds only the guide checklist gate", func(t *testing.T) {})
	g := DefaultGates(map[string]bool{"guide": true}, false)
	if len(g) != 1 || g[0].Name != "guide-checklist" {
		t.Errorf("guide alone should seed only guide-checklist, got %v", gateNames(g))
	}
}

func TestDefaultGatesParentValidFollowsChosenArtifacts(t *testing.T) {
	t.Run("DFGTD-B05: The parent gate confronts exactly the chosen artifacts among spec, plan and code", func(t *testing.T) {})
	find := func(chosen map[string]bool) *config.Gate {
		for _, g := range DefaultGates(chosen, false) {
			if g.Name == "parent-valid" {
				return &g
			}
		}
		return nil
	}
	if g := find(map[string]bool{"spec": true, "code": true, "feature": true}); g == nil ||
		!reflect.DeepEqual(g.On, []string{"spec", "code"}) {
		t.Errorf("parent-valid should run on [spec code], got %+v", g)
	}
	if g := find(map[string]bool{"feature": true, "test": true}); g != nil {
		t.Errorf("with none of spec, plan or code chosen, parent-valid is not seeded, got %+v", g)
	}
}

// In a NEW project the gates are born blocking, because the premise of maturation
// ("imposing the gate as blocking would stop the project") does not apply: there is no
// debt to accommodate. The gate stops nothing; it prevents the FIRST deviation, which is
// when fixing costs least.
func TestNewProjectIsBornWithBlockingGates(t *testing.T) {
	t.Run("DFGTD-B07: A new project is born with its gates blocking", func(t *testing.T) {})
	t.Run("DFGTD-I01: The list of seeded gates does not change with the age of the project", func(t *testing.T) {})
	newGates := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true}, true)
	existing := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true}, false)

	if len(newGates) != len(existing) {
		t.Fatalf("the gate list does not change with the age of the project: %d vs %d", len(newGates), len(existing))
	}
	for i := range newGates {
		if newGates[i].Name != existing[i].Name {
			t.Errorf("gate %d: %q in a new project, %q in an existing one", i, newGates[i].Name, existing[i].Name)
		}
		if !dependOnIngestedSignal[newGates[i].Name] && !newGates[i].IsBlocking() {
			t.Errorf("new project: %q should be born blocking", newGates[i].Name)
		}
	}
}

func TestExistingProjectIsBornInformative(t *testing.T) {
	t.Run("DFGTD-B06: An existing project is born with its gates informative, except the two blocking by nature", func(t *testing.T) {})
	all := map[string]bool{"spec": true, "feature": true, "test": true, "code": true, "guide": true, "plan": true}
	got := map[string]bool{}
	for _, g := range DefaultGates(all, false) {
		if g.IsBlocking() {
			got[g.Name] = true
		}
	}
	if !reflect.DeepEqual(got, blockingByNature) {
		t.Errorf("an existing project should only have the gates blocking by nature blocking; got %v", got)
	}
}

// Gates that depend on an INGESTED signal (test, coverage, mutation) stay informative
// even in a new project: with no report they would block for LACK OF DATA, not for a
// defect — the commit would fail before the suite even existed.
func TestGateWithoutSignalDoesNotBlockEvenInNewProject(t *testing.T) {
	t.Run("DFGTD-B08: A gate that depends on an ingested signal stays informative even in a new project", func(t *testing.T) {})
	g := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true, "code": true}, true)

	checked := 0
	for _, gate := range g {
		if !dependOnIngestedSignal[gate.Name] {
			continue
		}
		checked++
		if gate.IsBlocking() {
			t.Errorf("%q depends on an ingested signal and was born blocking — it would block for lack of data", gate.Name)
		}
	}
	if checked == 0 {
		t.Fatal("no gate depending on an ingested signal was seeded: the test confronted nothing")
	}
	// The gates the spec names as reading an ingested report, confronted by name so that
	// dropping one from the exemption is caught.
	byName := map[string]config.Gate{}
	for _, gate := range g {
		byName[gate.Name] = gate
	}
	for _, name := range []string{
		"tests-green", "line-coverage", "coverage-delta", "mutation-score", "scenario-coverage",
		"sbom-generated", "dependency-vulnerable", "no-duplication", "license-compatible",
		"circular", "deadcode", "spellcheck",
	} {
		gate, ok := byName[name]
		if !ok {
			t.Errorf("%q should be seeded for spec, feature, test and code", name)
			continue
		}
		if gate.IsBlocking() {
			t.Errorf("%q reads an ingested report and was born blocking in a new project", name)
		}
	}
}

// revision-orphans blocks (RVORP-Q01, decided by the user): it is of the blocking class,
// born blocking in a new project and matured in an existing one.
func TestRevisionOrphansIsOfTheBlockingClass(t *testing.T) {
	t.Run("DFGTD-B12: The gate is of the blocking class", func(t *testing.T) {})
	find := func(gates []config.Gate) *config.Gate {
		for i := range gates {
			if gates[i].Name == "revision-orphans" {
				return &gates[i]
			}
		}
		return nil
	}
	all := map[string]bool{"spec": true, "feature": true, "test": true, "code": true, "guide": true, "plan": true}
	novo, existente := find(DefaultGates(all, true)), find(DefaultGates(all, false))
	if novo == nil || existente == nil {
		t.Fatal("revision-orphans must be seeded in both")
	}
	if !novo.IsBlocking() {
		t.Error("a new project is born with revision-orphans blocking")
	}
	if existente.IsBlocking() {
		t.Error("an existing project takes it through maturation: informative until promoted")
	}
	if dependOnIngestedSignal["revision-orphans"] {
		t.Error("revision-orphans reads the spec alone; it must not wait for an ingested signal")
	}
}

// EVERY judgment gate that asks about a PIECE (code or test) carries the `@TBD`
// instruction.
//
// This is a GUARD AGAINST THE FUTURE: the defect of #76 was not writing the wrong `ask:`,
// it was that nothing demanded the instruction. A new judgment gate would be born without
// it, and the defect would return — under another name, in the gate nobody remembered to
// review.
//
// The exception is declared and justified below: a gate whose target is the WAIVER
// (`@no-test`) does not suffer the problem, because whoever declares a permanent waiver
// states the proof lives elsewhere — not that it is missing.
func TestJudgmentGateOnAPieceCarriesTheTBDInstruction(t *testing.T) {
	t.Run("DFGTD-B09: Every judgment gate that asks about code or a test carries the @TBD instruction", func(t *testing.T) {})
	// Judgment gates whose target is NOT a piece that can be under `@TBD`.
	//
	// `no-test-proof-real` asks about the proof `@no-test` points at — a PERMANENT waiver,
	// meaning "the test lives elsewhere". `@TBD` says the opposite ("still to be written"),
	// and the two do not overlap: the gate is already filtered by `Requires: "@no-test"`.
	exempt := map[string]string{
		"no-test-proof-real": "asks about the proof of a permanent waiver (@no-test), " +
			"not about a missing piece",
	}

	// EVERY artifact: gates are conditional on the choice, and a subset would leave out
	// exactly the gate nobody reviewed.
	gates := DefaultGates(allArtifacts(), false)

	var judgments int
	for _, g := range gates {
		if !g.IsJudgment() {
			continue
		}
		judgments++
		if reason, isExempt := exempt[g.Name]; isExempt {
			if strings.Contains(g.Ask, "@TBD") {
				t.Errorf("gate %q is in the exempt list (%s) and MENTIONS @TBD — "+
					"the list or the instruction is wrong", g.Name, reason)
			}
			continue
		}
		if !strings.Contains(g.Ask, "@TBD") {
			t.Errorf("the judgment gate %q does not instruct about `@TBD`.\n\n"+
				"Without it, facing a spec that declares the piece still to be written the "+
				"question has no subject, and the easy way out is to `pass` to unblock — "+
				"a stamp that stays in the map looking like real verification.\n\n"+
				"Append `tbdInstruction(\"the code\")` (or \"the test\") to the end of "+
				"`Ask:`, or declare the gate in this test's `exempt` list with the reason.\n\n"+
				"current ask: %s", g.Name, g.Ask)
		}
	}
	// If a refactor removed every judgment gate, the loop above would pass without
	// checking anything — and the test would report green over nothing.
	if judgments == 0 {
		t.Fatal("no judgment gate in the defaults: the test confronted nothing")
	}
}

func TestCanonicalGateLooksUpTheFullCatalog(t *testing.T) {
	t.Run("DFGTD-B10: The canonical declaration of a gate is found by name in the catalog of every artifact", func(t *testing.T) {})
	// guide-checklist is seeded only when guides are chosen, and plan-seeds-valid only
	// with plans: both are canonical because the catalog turns every artifact on.
	for _, name := range []string{"guide-checklist", "plan-seeds-valid", "tests-green"} {
		g, ok := CanonicalGate(name)
		if !ok || g.Name != name {
			t.Errorf("CanonicalGate(%q) should find the gate, got %+v %v", name, g, ok)
		}
	}
	if g, ok := CanonicalGate("no-such-gate"); ok || g.Name != "" {
		t.Errorf("an unknown name should not be found, got %+v %v", g, ok)
	}
}

func TestConfigLoadCompletesACanonicalGate(t *testing.T) {
	t.Run("DFGTD-B11: Loading a configuration completes a canonical gate declared by name alone", func(t *testing.T) {})
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 1\ngates:\n  - name: guide-checklist\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := config.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := CanonicalGate("guide-checklist")
	if got := c.Gates[0]; got.Check != want.Check || got.Measures != want.Measures || len(got.On) == 0 {
		t.Errorf("the gate declared by name should inherit the canonical declaration, got %+v", got)
	}
}

func TestSpecsSeedHeaderValid(t *testing.T) {
	t.Run("DFGTD-B16: Choosing specs seeds header-valid on specs and features", func(t *testing.T) {})
	var hv *config.Gate
	for _, g := range DefaultGates(map[string]bool{"spec": true}, false) {
		if g.Name == "header-valid" {
			g := g
			hv = &g
		}
	}
	if hv == nil {
		t.Fatal("choosing specs must seed header-valid")
	}
	if hv.Check != "header-valid" || strings.Join(hv.On, ",") != "spec,feature" || hv.Blocking == nil || *hv.Blocking {
		t.Errorf("header-valid must run its check on spec and feature, informative, got %+v", *hv)
	}
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 1\ngates:\n  - name: header-valid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := config.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(c.Gates[0].On, ","); got != "spec,feature" {
		t.Errorf("a bare header-valid must inherit on: spec,feature, got %q", got)
	}
}

func TestNoDuplicationIsTheNativeCheck(t *testing.T) {
	t.Run("DFGTD-B17: no-duplication is the native duplication check on code files", func(t *testing.T) {})
	g, ok := CanonicalGate("no-duplication")
	if !ok {
		t.Fatal("no-duplication must be in the catalog")
	}
	if g.Check != "duplication" || g.Run != "" || g.Scope != "" || g.ScopeFull != "" ||
		strings.Join(g.On, ",") != "code" || g.NeedsTool != "npx" {
		t.Errorf("no-duplication must be the per-file duplication check on code needing npx, got %+v", g)
	}
}

func TestDefaultGateNamesAreUnique(t *testing.T) {
	t.Run("DFGTD-I02: Every gate of the full catalog has a unique name that is also its id", func(t *testing.T) {})
	seen := map[string]bool{}
	for _, g := range DefaultGates(allArtifacts(), false) {
		if seen[g.Name] {
			t.Errorf("gate name %q is seeded twice", g.Name)
		}
		seen[g.Name] = true
		if g.ID != g.Name {
			t.Errorf("gate %q has id %q", g.Name, g.ID)
		}
	}
}

// The migration's legacy-to-canonical table must point at gates that EXIST.
//
// Each step's table converts a legacy name to the one of its format, and a later step may
// rename it again (format 2 made `trinca-completa` into `triad-complete`, format 6 made that
// `unit-complete`). So a target is followed through the later steps to where a project
// migrated today lands, and THAT must be a gate. A name that matches no gate would convert
// the project to a gate that does not exist — and it would vanish from `check` with nothing
// reporting it, which is worse than the old name.
func TestLegacyVocabularyPointsAtAnExistingGate(t *testing.T) {
	t.Run("DFGTD-I03: Every canonical name the migration renames a legacy gate to is a default gate", func(t *testing.T) {})
	existing := gateNames(DefaultGates(allArtifacts(), false))
	steps := migra.AllSteps()
	// landing follows a renamed value through the steps after the one that wrote it.
	landing := func(from int, file, key, name string) string {
		for _, later := range steps[from+1:] {
			if next, ok := later.RenameValues[file][key][name]; ok {
				name = next
			}
		}
		return name
	}
	for i, step := range steps {
		for file, keys := range step.RenameValues {
			for key, table := range keys {
				if key != "id" && key != "gate" {
					continue // `check:` points at an internal checker, not at a gate
				}
				for old, target := range table {
					canonical := landing(i, file, key, target)
					if !existing[canonical] {
						t.Errorf("%s/%s: %q → %q, and gate %q does not exist in the defaults",
							file, key, old, canonical, canonical)
					}
				}
			}
		}
	}
}

func TestNoDefaultGateHasAPortugueseName(t *testing.T) {
	t.Run("DFGTD-X01: No default gate carries a legacy name", func(t *testing.T) {})
	// The legacy names are exactly the list of what may no longer appear, and it now
	// lives in the migration step.
	forbidden := map[string]string{}
	for _, step := range migra.AllSteps() {
		for _, keys := range step.RenameValues {
			for _, table := range keys {
				for old, canonical := range table {
					forbidden[old] = canonical
				}
			}
		}
	}
	if len(forbidden) == 0 {
		t.Fatal("no legacy name found in the migration steps: the test confronted nothing")
	}
	for _, g := range DefaultGates(allArtifacts(), false) {
		if canonical, isLegacy := forbidden[g.Name]; isLegacy {
			t.Errorf("default gate %q still uses the legacy name — it should be %q", g.Name, canonical)
		}
	}
}

// allArtifacts turns on EVERY artifact choice.
//
// The legacy table must be confronted against the WHOLE set of default gates, not a
// subset: a target that only exists when the project chose "plan" would slip through a
// test that does not ask for plans — and the gate would become a nonexistent name exactly
// in the project that uses it.
func allArtifacts() map[string]bool {
	return map[string]bool{
		"spec": true, "feature": true, "test": true,
		"code": true, "plan": true, "guide": true, "doc": true,
	}
}

// What a default gate carries is written into every new project's anchors.yaml. The
// judgment questions, the @TBD instruction and an install hint were Portuguese.
func TestDefaultGateTextsAreEnglish(t *testing.T) {
	t.Run("DFGTD-X02: No default gate writes Portuguese into the project", func(t *testing.T) {})
	for _, g := range DefaultGates(allArtifacts(), false) {
		for field, text := range map[string]string{"ask": g.Ask, "measures": g.Measures, "install_hint": g.InstallHint} {
			words := " " + strings.ToLower(text) + " "
			portuguese := strings.ContainsAny(text, "ãõçáéíóúâêôÃÕÇÁÉÍÓÚ")
			for _, w := range []string{" que ", " de ", " da ", " não ", " um ", " uma ", " para ", "instale ", " acompanha"} {
				portuguese = portuguese || strings.Contains(words, w)
			}
			if portuguese {
				t.Errorf("gate %q: %s is not English: %s", g.Name, field, text)
			}
		}
	}
}

func TestMockDialectJudgmentPresupposesThePattern(t *testing.T) {
	t.Run("DFGTD-B18: The mock dialect judgment presupposes the pattern it asks about", func(t *testing.T) {})
	want := map[string]string{"mock-detect-covers-dialect": "derived.mock_detect", "mock-stamped": "derived.mock_detect", "mock-typed": "derived.mock_contract"}
	for _, g := range DefaultGates(allArtifacts(), false) {
		if field, ok := want[g.Name]; ok {
			if len(g.Presupposes) != 1 || g.Presupposes[0] != field {
				t.Errorf("%s presupposes %v, want [%s]", g.Name, g.Presupposes, field)
			}
			delete(want, g.Name)
		}
	}
	if len(want) > 0 {
		t.Errorf("not among the default gates: %v", want)
	}
}

func TestRuleFulfilledIsJudgedAndReviewed(t *testing.T) {
	t.Run("DFGTD-B19: rule-fulfilled is judged and marked to review", func(t *testing.T) {})
	for _, g := range DefaultGates(allArtifacts(), false) {
		if g.Name == "rule-fulfilled" {
			if !g.IsJudgment() || g.Review == nil || g.Review.Ask == "" {
				t.Errorf("rule-fulfilled should be judged and reviewed: %+v", g)
			}
			return
		}
	}
	t.Fatal("rule-fulfilled is not among the default gates")
}
