// @anchors
//   code: WRTSW
//   ref: WRPRW

package flow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/issue"
	"github.com/co2-lab/anchors/internal/mapx"
)

// workCfg is a project with a governed layer (guides, tags, extra steps), a layer that
// waives the feature, a layer with a test override, regimes with a surface, and gates.
func workCfg() *config.Config {
	return &config.Config{
		Layers: map[string]config.Layer{
			"logic": {Pattern: "src/**/*.ts", Kind: "code", Tags: []string{"backend"}, Regime: "comportamental",
				Work: map[string][]string{"code": {"run the migrations before the handler"}}},
			"model":   {Pattern: "models/**/*.ts", Kind: "code", OptionalUnitEdges: []string{"covered-by"}},
			"handler": {Pattern: "lambdas/**/*.ts", Kind: "code"},
			"dao":     {Pattern: "daos/**/*.ts", Kind: "code", Regime: "declarativo"},
			"spec":    {Pattern: "**/*.spec.md", Kind: "spec"},
			"feature": {Pattern: "**/*.feature", Kind: "feature"},
			"test":    {Pattern: "src/**/*.test.ts", Kind: "test"},
		},
		Governs: []config.GovernRule{
			{From: "guides/SPEC_GUIDE.md", Governs: "spec"},
			{From: "guides/BACKEND.md", Governs: "backend"},
			{From: "guides/BACKEND.md", Governs: "spec"},
		},
		Derived: &config.Derived{
			Anchor: "code",
			Files: map[string]config.Padroes{
				"spec":    {"{{dir}}/{{name}}.spec.md"},
				"feature": {"{{dir}}/{{name}}.feature"},
				"test":    {"{{dir}}/{{name}}.test.{{ext}}"},
			},
			Overrides: []config.DerivedOverride{
				{When: "handler", Files: map[string]config.Padroes{"test": {"tests/unit/{{module}}.test.{{ext}}"}}},
			},
			Regimes:           map[string]string{"unit-level": "unit", "integration-level": "integration"},
			Surfaces:          map[string]string{"integration": "api"},
			RuleMarkingPolicy: "required",
		},
		Gates: []config.Gate{
			{Name: "sections", Check: "spec-sections", On: []string{"spec"}, Blocking: config.Bool(true)},
			{Name: "sections-again", Check: "spec-sections", On: []string{"spec"}, Blocking: config.Bool(true)},
			{Name: "ftm", Check: "feature-test-match", On: []string{"feature"}},
			{Name: "custom", Check: "not-described", On: []string{"code"}},
		},
	}
}

func prompt(t *testing.T, root, rel, artifact string, cfg *config.Config) string {
	t.Helper()
	out, err := composeWorkPrompt(root, rel, artifact, cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func wantAll(t *testing.T, out string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(out, w) {
			t.Errorf("the prompt lacks %q:\n%s", w, out)
		}
	}
}

func wantNone(t *testing.T, out string, nots ...string) {
	t.Helper()
	for _, w := range nots {
		if strings.Contains(out, w) {
			t.Errorf("the prompt must not carry %q:\n%s", w, out)
		}
	}
}

// The spec stage: the layer and its guides (the artifact's AND the layer's, each once, in
// order), where the pieces are born, what the gates demand (each check once), and a
// verification over the spec — not over a code file that does not exist yet.
func TestComposeWorkPrompt_spec(t *testing.T) {
	t.Run("WRPRW-B05: The heading and the role follow the stage", func(t *testing.T) {})
	t.Run("WRPRW-B06: The target section names the layer, or says it is unclassified", func(t *testing.T) {})
	t.Run("WRPRW-B07: The guides to read are the artifact's and the layer's, each once and sorted", func(t *testing.T) {})
	t.Run("WRPRW-B08: The pieces are listed where the project derives them", func(t *testing.T) {})
	t.Run("WRPRW-B10: Each stage says what is not its scope", func(t *testing.T) {})
	t.Run("WRPRW-B12: The gates' demands are listed once, for the stage", func(t *testing.T) {})
	t.Run("WRPRW-B16: Producing stages record the delivery, reviews close with a verdict", func(t *testing.T) {})
	t.Run("WRPRW-B17: The verification runs over the stage's own piece", func(t *testing.T) {})
	t.Run("WRPRW-B18: The waivers are listed for the stage", func(t *testing.T) {})
	t.Run("WRPRW-B19: When the ruler does not decide, the open-decisions section is named", func(t *testing.T) {})
	out := prompt(t, t.TempDir(), "src/pricing.ts", "spec", workCfg())
	wantAll(t, out,
		"# Work: spec of `src/pricing.ts`",
		"- Layer: **logic** (regime: comportamental)", "- Tags: backend",
		"2. `guides/BACKEND.md`", "3. `guides/SPEC_GUIDE.md`",
		"→ `src/pricing.spec.md` — spec", "`src/pricing.feature` — feature", "`src/pricing.test.ts` — test",
		"**Do not write code, feature or test**",
		"**Every rule must be CATALOGUED**",
		"anchors deliver --stage spec --unit src/pricing.ts",
		"anchors check --changed src/pricing.spec.md --no-record --deterministic",
		"anchors check --changed src/pricing.spec.md --deterministic\n",
		"`@no-mark: <reason>`", "`@no-scenario: <reason>`",
		"Write the question in `## "+openSectionTitle()+"`", "write `"+noneValue()+"`",
	)
	if strings.Count(out, "CATALOGUED") != 1 {
		t.Error("a check declared by two gates is described once")
	}
	if strings.Count(out, "guides/BACKEND.md") != 1 {
		t.Error("a guide that governs two tags is listed once")
	}
	wantNone(t, out, "Steps specific to this layer", "Valid regime tags", "Execution signals")
}

// The code stage: the layer's own steps, the marking the project requires, and the
// code opt-outs.
func TestComposeWorkPrompt_code(t *testing.T) {
	t.Run("WRPRW-B11: The procedure follows the stage and the layer", func(t *testing.T) {})
	t.Run("WRPRW-B20: Rule marking follows the project's policy", func(t *testing.T) {})
	out := prompt(t, t.TempDir(), "src/pricing.ts", "code", workCfg())
	wantAll(t, out,
		"**Steps specific to this layer** (declared in anchors.yaml)", "- run the migrations before the handler",
		"this project REQUIRES the marking",
		"`@no-paginate: <reason>`", "`@allow-boundary: <reason>`",
		"anchors check --changed src/pricing.ts --no-record --deterministic",
		"**Stop and read `## "+openSectionTitle()+"` in the spec.**",
	)
	// the gate with no description is not promised
	wantNone(t, out, "not-described", "What the gates will demand")
}

// The test stage: the override path, the regime tags, the feature-test-match demand
// (declared on the feature, owed by the test) marked informative, and the signals.
func TestComposeWorkPrompt_testWithOverrideAndRegimes(t *testing.T) {
	t.Run("WRPRW-B09: Feature and test stages list the regime tags", func(t *testing.T) {})
	t.Run("WRPRW-B15: Tests and unit reviews explain the execution signals", func(t *testing.T) {})
	out := prompt(t, t.TempDir(), "lambdas/push/send.ts", "test", workCfg())
	wantAll(t, out,
		"→ `tests/unit/push.test.ts` — test  ← layer override (not co-located)",
		"## Valid regime tags (use EXACTLY these)",
		"- `@integration-level` (regime integration) → proved on the surface `api`",
		"- `@unit-level` (regime unit)\n",
		"declares which scenario it proves", "*(informative)*",
		"## Execution signals",
		"anchors check --changed tests/unit/push.test.ts --no-record",
		"**Do not mock the logic under test**",
	)
	if strings.Index(out, "@integration-level") > strings.Index(out, "@unit-level") {
		t.Error("the regime tags are listed in order")
	}
}

// The prompt cites what the project and the tool really call things. It told an English
// project to write `## Open Decisions` … `none` in the procedure and `## Decisões em
// aberto` … `nenhuma` in the gates' demands, and cited `trinca_opcional`,
// `regra-implementada` and `testes-passam`, names no configuration or check has.
func TestComposeWorkPrompt_speaksTheCurrentVocabulary(t *testing.T) {
	t.Run("WRPRW-X02: The prompt cites one open-decisions title and value and the current names", func(t *testing.T) {})
	cfg := workCfg()
	cfg.Gates = append(cfg.Gates,
		config.Gate{Name: "open", Check: "open-questions-resolved", On: []string{"spec"}},
		config.Gate{Name: "rules", Check: "rule-implemented", On: []string{"code"}})
	title, none := openSectionTitle(), noneValue()
	spec := prompt(t, t.TempDir(), "src/pricing.ts", "spec", cfg)
	wantAll(t, spec, "Write the question in `## "+title+"`", "The section **`## "+title+"`** is REQUIRED",
		"write `"+none+"`", "or it carries `"+none+"`")
	if title == "Open Decisions" {
		wantNone(t, spec, "Decisões em aberto", "nenhuma")
	}
	all := spec +
		prompt(t, t.TempDir(), "src/pricing.ts", "code", cfg) +
		prompt(t, t.TempDir(), "src/pricing.ts", "test", cfg) +
		prompt(t, t.TempDir(), "models/user.ts", "feature", cfg) +
		prompt(t, t.TempDir(), "models/user.ts", "spec", cfg)
	wantAll(t, all, "`optional_unit_edges` in anchors.yaml", "the check `rule-implemented` confronts it",
		"the check `tests-pass` reads from here")
	wantNone(t, all, "trinca_opcional", "regra-implementada", "testes-passam")
}

// A piece the layer waives is refused at the top, before any production script.
func TestComposeWorkPrompt_waivedPieceStops(t *testing.T) {
	t.Run("WRPRW-B03: A waived piece stops the prompt", func(t *testing.T) {})
	out := prompt(t, t.TempDir(), "models/user.ts", "feature", workCfg())
	wantAll(t, out, "## STOP", "the layer **model** WAIVES the piece `feature`", "Do not create the feature.")
	wantNone(t, out, "You are going to produce", "## Procedure")

	// and in the unit listing of another stage, the waived piece says so
	spec := prompt(t, t.TempDir(), "models/user.ts", "spec", workCfg())
	wantAll(t, spec, "`models/user.feature` — feature  ← WAIVED by this layer (`optional_unit_edges`): do NOT create")
}

// A recognized (declarative) layer has no unit: the unit stages stop, and the code
// stage gets the declarative procedure.
func TestComposeWorkPrompt_declarativeLayer(t *testing.T) {
	t.Run("WRPRW-B04: A declarative layer has no unit piece", func(t *testing.T) {})
	stop := prompt(t, t.TempDir(), "daos/user.ts", "test", workCfg())
	wantAll(t, stop, "## STOP", "declared as RECOGNIZED", "Do not create the test.")

	code := prompt(t, t.TempDir(), "daos/user.ts", "code", workCfg())
	wantAll(t, code, "**Only the file itself.**", "Translate, transport or declare — **do not decide**")
	wantNone(t, code, "daos/user.spec.md")
}

// The review confronts what exists and the author's records — every record of the unit,
// none of another — and closes with the judge commands instead of a delivery.
func TestComposeWorkPrompt_reviewWithDeliveryRecords(t *testing.T) {
	t.Run("WRPRW-B13: A review confronts the unit's delivery records", func(t *testing.T) {})
	t.Run("WRPRW-B14: The findings already recorded for the unit are listed", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "changes/2026-09-01-code.md", "stage: code\nunit: src/pricing.ts\nintent: the code\n")
	writeFile(t, root, "changes/2026-09-02-test.md", "stage: test\nunit: src/pricing.ts\nintent: the tests\n")
	writeFile(t, root, "changes/2026-09-03-other.md", "stage: code\nunit: src/tax.ts\nintent: other unit\n")
	// real issue names: `<date>--<kind>--[gate--][anchor--vs--]target.md`
	stale := issue.Issue{Kind: issue.Stale, Anchor: "src/pricing.spec.md", Target: "src/pricing.ts", Date: "2026-09-01"}.ID()
	review := issue.Issue{Kind: issue.Violation, Gate: "review", Target: "src/pricing.test.ts", Date: "2026-09-02"}.ID()
	writeFile(t, root, "issues/todo/"+stale, "x\n")
	writeFile(t, root, "issues/doing/"+review, "x\n")
	writeFile(t, root, "issues/todo/"+issue.Issue{Kind: issue.Stale, Target: "src/tax.ts", Date: "2026-09-01"}.ID(), "x\n")
	// a unit whose name only CONTAINS the target's stem is another unit
	writeFile(t, root, "issues/todo/"+issue.Issue{Kind: issue.Violation, Gate: "g", Target: "src/pricing-v2.ts", Date: "2026-09-01"}.ID(), "x\n")
	writeFile(t, root, "issues/todo/"+issue.Issue{Kind: issue.Violation, Gate: "g", Target: "lib/src/pricing.ts", Date: "2026-09-01"}.ID(), "x\n")

	out := prompt(t, root, "src/pricing.ts", "review", workCfg())
	wantAll(t, out,
		"# Review of `src/pricing.ts`", "You are the REVIEWER of this unit.",
		"## What you are going to confront",
		"## Delivery records of this unit (2)", "intent: the code", "intent: the tests",
		"**Do not correct the code**",
		"anchors judge src/pricing.ts --gate review --verdict pass",
		"## Findings ALREADY RECORDED about this unit",
		"`issues/doing/"+review+"`", "`issues/todo/"+stale+"`",
		"## Execution signals",
	)
	wantNone(t, out, "other unit", "src-tax", "src-pricing-v2", "lib-src-pricing", "How to RECORD the delivery")
	if strings.Index(out, "intent: the code") > strings.Index(out, "intent: the tests") {
		t.Error("the records come in order")
	}
}

// Without a record, the review is told why — and in github mode the record lives in the
// issue, so "no record" there would be a false negative.
func TestWriteDeliveryRecord_absentRecordDependsOnTheMode(t *testing.T) {
	var local strings.Builder
	writeDeliveryRecord(&local, t.TempDir(), "src/pricing.ts", workCfg())
	if !strings.Contains(local.String(), "**No delivery record** for this unit") {
		t.Errorf("local mode without a record says so: %s", local.String())
	}
	var gh strings.Builder
	writeDeliveryRecord(&gh, t.TempDir(), "src/pricing.ts", cfgGitHub())
	if !strings.Contains(gh.String(), "is in the COMMENTS of the issue") || strings.Contains(gh.String(), "No delivery record") {
		t.Errorf("github mode points to the issue: %s", gh.String())
	}
}

// The review of the whole and of a draft plan: their own role, frontier and procedure,
// and no per-unit delivery section.
func TestComposeWorkPrompt_planReviews(t *testing.T) {
	whole := prompt(t, t.TempDir(), "plans/0001-f.md", "review-plan", workCfg())
	wantAll(t, whole, "# WHOLE review — `plans/0001-f.md`", "You are the REVIEWER OF THE WHOLE.",
		"**Do not review the isolated piece**", "**Follow the DATA end to end.**",
		"## How to CLOSE the review (REQUIRED)")
	wantNone(t, whole, "How to RECORD the delivery")

	draft := prompt(t, t.TempDir(), "plans/0002-g.draft.md", "review-plan-draft", workCfg())
	wantAll(t, draft, "**Do not review the CODE**", "**Hunt the item whose EXISTENCE depends on a decision not made.**",
		"## How to CLOSE the review (REQUIRED)")
	wantNone(t, draft, "How to RECORD the delivery")
}

// A target outside every layer, with no guides and no `derived:`, still gets a prompt
// that says what is missing instead of inventing it.
func TestComposeWorkPrompt_unclassifiedWithoutDerived(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{"spec": {Pattern: "**/*.spec.md", Kind: "spec"}}}
	out := prompt(t, t.TempDir(), "scripts/tool.py", "feature", cfg)
	wantAll(t, out,
		"- Layer: **unclassified**",
		"> No project guide governs this layer",
		"> The project does not declare `derived:`",
		"**Do not create a new scenario code**",
		"anchors check --changed scripts/tool.py --no-record",
	)
	wantNone(t, out, "Valid regime tags")
}

// The command: it validates the artifact and the target, loads the project, and
// redirects a derived piece to its unit, saying so.
func TestWorkCmd(t *testing.T) {
	t.Run("WRPRW-B01: The command refuses what it cannot compose", func(t *testing.T) {})
	t.Run("WRPRW-B02: A derived piece is redirected to its unit", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "anchors.yaml", "version: 1\nlayers:\n"+
		"  logic:\n    pattern: \"src/**/*.ts\"\n    kind: code\n"+
		"  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n")
	writeFile(t, root, "src/pricing.ts", "x\n")
	writeFile(t, root, "src/pricing.spec.md", "# s\n")

	var out string
	var err error
	errOut := stderrOf(t, func() {
		out, err = runFlowCmd(t, newWorkCmd(), "code", "--root", root, "--for", filepath.Join(root, "src/pricing.spec.md"))
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(errOut, "`src/pricing.spec.md` is a derived piece, not the unit. Using `src/pricing.ts`") {
		t.Errorf("the redirect is said: %q", errOut)
	}
	if !strings.Contains(out, "# Work: code of `src/pricing.ts`") {
		t.Errorf("the prompt is about the unit:\n%s", out)
	}

	if _, err := runFlowCmd(t, newWorkCmd(), "deploy", "--root", root, "--for", "src/pricing.ts"); err == nil ||
		!strings.Contains(err.Error(), `unknown artifact "deploy"`) {
		t.Errorf("an unknown artifact is refused, got %v", err)
	}
	if _, err := runFlowCmd(t, newWorkCmd(), "code", "--root", root); err == nil || !strings.Contains(err.Error(), "--for") {
		t.Errorf("a missing target is refused, got %v", err)
	}
	if _, err := runFlowCmd(t, newWorkCmd(), "code", "--root", t.TempDir(), "--for", "x.ts"); err == nil ||
		!strings.Contains(err.Error(), "load config") {
		t.Errorf("a project without config is refused, got %v", err)
	}
}

// Without the required policy, marking is a suggestion, not an obligation.
func TestRuleMarking_optionalWithoutPolicy(t *testing.T) {
	if got := ruleMarking(nil); !strings.Contains(got, "if the project uses that pattern") {
		t.Errorf("got %q", got)
	}
	if got := ruleMarking(workCfg()); strings.Contains(got, "if the project uses") {
		t.Errorf("the required policy leaves no escape: %q", got)
	}
}

func TestWaivedPieces(t *testing.T) {
	cfg := workCfg()
	if got := waivedPieces("model", cfg); !got["feature"] || got["test"] || got["spec"] {
		t.Errorf("model waives only the feature: %v", got)
	}
	if got := waivedPieces("ghost", cfg); len(got) != 0 {
		t.Errorf("an unknown layer waives nothing: %v", got)
	}
	if got := waivedPieces("", cfg); len(got) != 0 {
		t.Errorf("no layer waives nothing: %v", got)
	}
}

// guard against the prompt writing files: composing is read-only.
func TestComposeWorkPrompt_writesNothing(t *testing.T) {
	t.Run("WRPRW-X01: Composing writes nothing", func(t *testing.T) {})
	root := t.TempDir()
	prompt(t, root, "src/pricing.ts", "review", workCfg())
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Errorf("composing a prompt must not touch the project: %v", entries)
	}
}

// "(already exists)" is about the PROJECT's files: `anchors work --root X` run from
// anywhere else must still see what exists under X.
func TestWriteUnitPaths_alreadyExistsIsCheckedUnderTheRoot(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "src/billing/charge.spec.md", "# spec\n")
	out := prompt(t, root, "src/billing/charge.ts", "code", workCfg())
	if !strings.Contains(out, "`src/billing/charge.spec.md` — spec  (already exists)") {
		t.Errorf("the spec exists under the root and must be marked so:\n%s", out)
	}
	if strings.Contains(out, "charge.feature` — feature  (already exists)") {
		t.Errorf("the feature does not exist and must not be marked:\n%s", out)
	}
}

// A producing stage records its delivery and a review closes with a verdict — never both.
func TestComposeWorkPrompt_recordOrCloseNeverBoth(t *testing.T) {
	t.Run("WRPRW-I01: No prompt both records a delivery and closes a review", func(t *testing.T) {})
	for _, artifact := range []string{"spec", "code", "feature", "test", "review", "review-plan", "review-plan-draft"} {
		out := prompt(t, t.TempDir(), "src/pricing.ts", artifact, workCfg())
		record := strings.Contains(out, "How to RECORD the delivery")
		closes := strings.Contains(out, "How to CLOSE the review")
		review := strings.HasPrefix(artifact, "review")
		if record == closes || record == review {
			t.Errorf("%s: record=%v close=%v", artifact, record, closes)
		}
	}
}

func cfgUnit() *config.Config {
	return &config.Config{Layers: map[string]config.Layer{
		"logic":   {Pattern: "src/**/*.ts", Kind: "code"},
		"spec":    {Pattern: "**/*.spec.md", Kind: "spec"},
		"feature": {Pattern: "**/*.feature", Kind: "feature"},
		"test":    {Pattern: "**/*.test.ts", Kind: "test"},
	}}
}

// The target of `work` is the UNIT, not a derived piece. Pointing at the spec is the
// predictable mistake (it is the artifact that already exists), and the result was silently
// absurd: paths like `x.spec.spec.md` and `x.spec.feature`, with no warning.
func TestWork_redirectsADerivedPiece(t *testing.T) {
	t.Run("WRPRW-B02: A derived piece is redirected to its unit", func(t *testing.T) {})
	dir := t.TempDir()
	unit := "src/metadataVersioning.ts"
	writeFile(t, dir, unit, "export const x = 1\n")

	for _, piece := range []string{
		"src/metadataVersioning.spec.md",
		"src/metadataVersioning.feature",
		"src/metadataVersioning.test.ts",
	} {
		t.Run(piece, func(t *testing.T) {
			got, found := derivedPieceUnit(dir, piece, cfgUnit(), nil)
			if !found {
				t.Fatalf("%q was not recognised as a derived piece", piece)
			}
			if got != unit {
				t.Fatalf("redirected to %q, want %q", got, unit)
			}
		})
	}
}

// The code file itself is NOT redirected — it already is the unit.
func TestWork_doesNotRedirectTheUnit(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "src/x.ts", "x\n")
	if _, found := derivedPieceUnit(dir, "src/x.ts", cfgUnit(), nil); found {
		t.Fatal("the code is the unit — it should not be redirected")
	}
}

// With a map, the `specifies` edge is the precise source: it says exactly which code the
// spec describes, even when the naming convention would not be enough — and a feature or a
// test reaches it through the spec it is linked to.
func TestWork_usesTheMapWhenItExists(t *testing.T) {
	g := &mapx.Graph{Edges: []mapx.Edge{
		{From: "a/different-name.spec.md", To: "b/otherName.ts", Type: "specifies"},
		{From: "a/different-name.spec.md", To: "c/scenarios.feature", Type: "covered-by"},
	}}
	got, found := derivedPieceUnit(t.TempDir(), "a/different-name.spec.md", cfgUnit(), g)
	if !found || got != "b/otherName.ts" {
		t.Fatalf("the map should resolve the target: got=%q found=%v", got, found)
	}
	got, found = derivedPieceUnit(t.TempDir(), "c/scenarios.feature", cfgUnit(), g)
	if !found || got != "b/otherName.ts" {
		t.Fatalf("a feature covered by the spec resolves through it: got=%q found=%v", got, found)
	}
}

// A derived piece whose unit does not exist: no target is invented (the caller goes on with
// the original and the rest of the prompt explains what is missing).
func TestWork_withoutAUnitInventsNothing(t *testing.T) {
	if _, found := derivedPieceUnit(t.TempDir(), "src/ghost.spec.md", cfgUnit(), nil); found {
		t.Fatal("with no code on disk and no map, there is no target to deduce")
	}
}

func cfgWithDeclarative() *config.Config {
	return &config.Config{
		Layers: map[string]config.Layer{
			"dao":   {Pattern: "backend/models/**/*.ts", Kind: "code", Regime: "declarativo"},
			"logic": {Pattern: "backend/business-logic/**/*.ts", Kind: "code"},
			"spec":  {Pattern: "**/*.spec.md", Kind: "spec"},
		},
		// `derived:` resolves the unit's paths; without it the prompt has nothing to
		// derive (and says so). The fixture needs it to exercise the common path.
		Derived: &config.Derived{
			Anchor: "code",
			Files:  map[string]config.Padroes{"spec": {"{{dir}}/{{name}}.spec.md"}},
		},
	}
}

// The prompt contradicted itself: it declared "Layer: dao (regime: declarativo)" and three
// lines below listed `metadata.spec.md` among "the pieces and where they are born" — the spec
// the layer FORBIDS — and the procedure said "read the whole spec; it is the ruler".
func TestWork_declarativeLayerAsksForNoSpec(t *testing.T) {
	t.Run("WRPRW-B11: The procedure follows the stage and the layer", func(t *testing.T) {})
	out, err := composeWorkPrompt(t.TempDir(), "backend/models/metadata.ts", "code", cfgWithDeclarative(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "metadata.spec.md") {
		t.Errorf("it listed the spec the declarative layer forbids:\n%s", out)
	}
	if strings.Contains(out, "Read the whole spec") {
		t.Errorf("it told the worker to read a spec that does not exist:\n%s", out)
	}
	// and it says what holds instead
	for _, want := range []string{"has no spec", "does not originate rules", "do not decide"} {
		if !strings.Contains(out, want) {
			t.Errorf("the prompt does not explain the declarative layer (missing %q)", want)
		}
	}
}

// The GOVERNED layer keeps the unit and the normal procedure — the fix must not have
// emptied the common path.
func TestWork_governedLayerKeepsTheUnit(t *testing.T) {
	out, err := composeWorkPrompt(t.TempDir(), "backend/business-logic/recurrence.ts", "code", cfgWithDeclarative(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "recurrence.spec.md") {
		t.Errorf("a governed layer should list the unit's spec:\n%s", out)
	}
	if !strings.Contains(out, "Read the whole spec") {
		t.Errorf("a governed layer should keep the normal procedure")
	}
}

// THE UNIT'S NAME is not the basename minus the last extension. `filepath.Ext("X.spec.md")`
// is `.md`, so `{{name}}` came out `X.spec` and the prompt pointed to
// `GoLiveChecklist.spec.feature` — a file no gate finds, while the map (which cuts the
// COMPOUND suffix) looked for `GoLiveChecklist.feature`.
func TestDerivedPaths_keepsNoSpecInTheName(t *testing.T) {
	t.Run("WRPRW-B08: The pieces are listed where the project derives them", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "spec",
		Files: map[string]config.Padroes{
			"code":    {"{{dir}}/{{name}}.ts"},
			"feature": {"{{dir}}/{{name}}.feature"},
			"test":    {"{{dir}}/{{name}}.test.ts"},
		},
	}}

	files, _ := derivedPaths("packages/infra/GoLiveChecklist.spec.md", "spec", cfg)

	want := map[string]string{
		"code":    "packages/infra/GoLiveChecklist.ts",
		"feature": "packages/infra/GoLiveChecklist.feature",
		"test":    "packages/infra/GoLiveChecklist.test.ts",
	}
	for k, w := range want {
		if got := files[k]; got != w {
			t.Errorf("%s: %q — want %q (the `.spec` in the middle points to a file no gate finds)", k, got, w)
		}
	}
}

// And the FEATURE is a compound-suffix anchor too: `X.feature` → `X`.
func TestDerivedPaths_cutsTheFeatureSuffix(t *testing.T) {
	cfg := &config.Config{Derived: &config.Derived{
		Anchor: "spec",
		Files:  map[string]config.Padroes{"test": {"{{dir}}/{{name}}.test.ts"}},
	}}
	files, _ := derivedPaths("apps/mobile/src/Login.feature", "feature", cfg)
	if got := files["test"]; got != "apps/mobile/src/Login.test.ts" {
		t.Errorf("test: %q — want apps/mobile/src/Login.test.ts", got)
	}
}
