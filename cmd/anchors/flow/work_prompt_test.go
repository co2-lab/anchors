package flow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// workCfg is a project with a governed layer (guides, tags, extra steps), a layer that
// waives the feature, a layer with a test override, regimes with a surface, and gates.
func workCfg() *config.Config {
	return &config.Config{
		Layers: map[string]config.Layer{
			"logic": {Pattern: "src/**/*.ts", Kind: "code", Tags: []string{"backend"}, Regime: "comportamental",
				Work: map[string][]string{"code": {"run the migrations before the handler"}}},
			"model":   {Pattern: "models/**/*.ts", Kind: "code", OptionalTriadEdges: []string{"covered-by"}},
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
		"`@no-mark: <reason>`",
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
	out := prompt(t, t.TempDir(), "src/pricing.ts", "code", workCfg())
	wantAll(t, out,
		"**Steps specific to this layer** (declared in anchors.yaml)", "- run the migrations before the handler",
		"this project REQUIRES the marking",
		"`@no-paginate: <reason>`",
		"anchors check --changed src/pricing.ts --no-record --deterministic",
		"**Stop and read `## "+openSectionTitle()+"` in the spec.**",
	)
	// the gate with no description is not promised
	wantNone(t, out, "not-described", "What the gates will demand")
}

// The test stage: the override path, the regime tags, the feature-test-match demand
// (declared on the feature, owed by the test) marked informative, and the signals.
func TestComposeWorkPrompt_testWithOverrideAndRegimes(t *testing.T) {
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

// A piece the layer waives is refused at the top, before any production script.
func TestComposeWorkPrompt_waivedPieceStops(t *testing.T) {
	out := prompt(t, t.TempDir(), "models/user.ts", "feature", workCfg())
	wantAll(t, out, "## STOP", "the layer **model** WAIVES the piece `feature`", "Do not create the feature.")
	wantNone(t, out, "You are going to produce", "## Procedure")

	// and in the triad listing of another stage, the waived piece says so
	spec := prompt(t, t.TempDir(), "models/user.ts", "spec", workCfg())
	wantAll(t, spec, "`models/user.feature` — feature  ← WAIVED by this layer (`trinca_opcional`): do NOT create")
}

// A recognized (declarative) layer has no triad: the triad stages stop, and the code
// stage gets the declarative procedure.
func TestComposeWorkPrompt_declarativeLayer(t *testing.T) {
	stop := prompt(t, t.TempDir(), "daos/user.ts", "test", workCfg())
	wantAll(t, stop, "## STOP", "declared as RECOGNIZED", "Do not create the test.")

	code := prompt(t, t.TempDir(), "daos/user.ts", "code", workCfg())
	wantAll(t, code, "**Only the file itself.**", "Translate, transport or declare — **do not decide**")
	wantNone(t, code, "daos/user.spec.md")
}

// The review confronts what exists and the author's records — every record of the unit,
// none of another — and closes with the judge commands instead of a delivery.
func TestComposeWorkPrompt_reviewWithDeliveryRecords(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "changes/2026-09-01-code.md", "stage: code\nunit: src/pricing.ts\nintent: the code\n")
	writeFile(t, root, "changes/2026-09-02-test.md", "stage: test\nunit: src/pricing.ts\nintent: the tests\n")
	writeFile(t, root, "changes/2026-09-03-other.md", "stage: code\nunit: src/tax.ts\nintent: other unit\n")
	writeFile(t, root, "issues/todo/stale-src-pricing-edge.md", "x\n")
	writeFile(t, root, "issues/doing/review-src-pricing-b04.md", "x\n")
	writeFile(t, root, "issues/todo/stale-src-tax.md", "x\n")

	out := prompt(t, root, "src/pricing.ts", "review", workCfg())
	wantAll(t, out,
		"# Review of `src/pricing.ts`", "You are the REVIEWER of this unit.",
		"## What you are going to confront",
		"## Delivery records of this unit (2)", "intent: the code", "intent: the tests",
		"**Do not correct the code**",
		"anchors judge src/pricing.ts --gate review --verdict pass",
		"## Findings ALREADY RECORDED about this unit",
		"`issues/doing/review-src-pricing-b04.md`", "`issues/todo/stale-src-pricing-edge.md`",
		"## Execution signals",
	)
	wantNone(t, out, "other unit", "stale-src-tax", "How to RECORD the delivery")
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
	root := t.TempDir()
	prompt(t, root, "src/pricing.ts", "review", workCfg())
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Errorf("composing a prompt must not touch the project: %v", entries)
	}
}

// "(already exists)" is about the PROJECT's files: `anchors work --root X` run from
// anywhere else must still see what exists under X.
func TestWriteTriadPaths_alreadyExistsIsCheckedUnderTheRoot(t *testing.T) {
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
