package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/i18n"
	"gopkg.in/yaml.v3"
)

// load writes the YAML to a temporary anchors.yaml and loads it.
func load(t *testing.T, yaml string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

// ── loading the file ──────────────────────────────────────────────────────────

// The friction that motivated KnownFields, measured when putting Anchors on itself: a
// `governs` block written with the wrong keys was dropped in silence, and `map build`
// answered "222 nodes, 0 edges".
func TestLoad_unknownKeyFailsInsteadOfBeingIgnored(t *testing.T) {
	t.Run("CNFGO-B01: An unknown key in the file is a load error naming the key", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\ngovernz: []\n")
	if err == nil {
		t.Fatal("an unknown key must fail — ignoring it turns the declaration into an invisible no-op")
	}
	if !strings.Contains(err.Error(), "governz") {
		t.Errorf("the error must NAME the wrong key, got: %v", err)
	}
}

// The exact case that bit: the top-level keys are right, the mistake is INSIDE a list item.
func TestLoad_wrongKeyInsideGovernsAlsoFails(t *testing.T) {
	t.Run("CNFGO-B01: An unknown key in the file is a load error naming the key", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\ngoverns:\n  - guide: X.md\n    tags: [y]\n")
	if err == nil {
		t.Fatal("a wrong key inside governs[] must fail")
	}
	if !strings.Contains(err.Error(), "guide") && !strings.Contains(err.Error(), "tags") {
		t.Errorf("the error must name the key, got: %v", err)
	}
}

// A NEW field breaks every project whose binary predates it, and the yaml message speaks
// of a Go type. The hint names the other hypothesis: the binary is older than the file.
func TestLoad_unknownKeyNamesTheVersionHypothesis(t *testing.T) {
	t.Run("CNFGO-B02: An unknown key gets the misspelling and old-binary hypotheses", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\nchaveQueNaoExiste: 1\n")
	if err == nil {
		t.Fatal("an unknown key passed — KnownFields exists to catch it")
	}
	msg := err.Error()
	if !strings.Contains(msg, "chaveQueNaoExiste") {
		t.Errorf("the message does not name the key: %s", msg)
	}
	if !strings.Contains(msg, "binary is OLD") {
		t.Errorf("the message does not offer the version hypothesis:\n%s", msg)
	}
	if !strings.Contains(msg, "go install") {
		t.Errorf("the message states the problem but not the fix: %s", msg)
	}
}

// The hint only appears when it fits: malformed YAML has nothing to do with versions.
func TestLoad_syntaxErrorGetsNoVersionHint(t *testing.T) {
	t.Run("CNFGO-B05: An error that is not an unknown key gets no version hint", func(t *testing.T) {})
	_, err := load(t, "version: 1\n  layers: [\n")
	if err == nil {
		t.Fatal("broken YAML passed")
	}
	if strings.Contains(err.Error(), "binary is OLD") || strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("a syntax error got a version hint: %s", err)
	}
}

// The OUTDATED FILE is the hypothesis that happens most, and it is decidable: the file's
// `version:` says it. Measured in a PR of the reference project: the message told the
// author to UPDATE THE BINARY, the inverse of the fix.
func TestLoad_oldFileWithRenamedKeyAdvisesMigration(t *testing.T) {
	t.Run("CNFGO-B03: A renamed key in an older-format file advises the migration", func(t *testing.T) {})
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = func(k string) bool { return k == "trinca_opcional" }

	_, err := load(t, "version: 1\nlayers:\n    spec:\n        pattern: \"**/*.spec.md\"\n"+
		"        kind: spec\n        trinca_opcional: [tested-by]\n")
	if err == nil {
		t.Fatal("the old key should be refused — it was renamed")
	}
	msg := err.Error()
	if !strings.Contains(msg, "anchors migrate") {
		t.Errorf("the message should advise MIGRATING; got:\n%s", msg)
	}
	if strings.Contains(msg, "go install") {
		t.Errorf("the binary is right — advising to update it is the INVERSE advice:\n%s", msg)
	}
}

// BOTH conditions are needed: an old file with a real typo must not be told to migrate.
func TestLoad_oldFileWithTypoDoesNotAdviseMigration(t *testing.T) {
	t.Run("CNFGO-B04: The migration advice needs both an older format and a renamed key", func(t *testing.T) {})
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = func(k string) bool { return k == "trinca_opcional" }

	_, err := load(t, "version: 1\nlayers: {}\nlayerz: 1\n")
	if err == nil {
		t.Fatal("the typo should be refused")
	}
	msg := err.Error()
	if strings.Contains(msg, "anchors migrate") {
		t.Errorf("migration does not fix a typo:\n%s", msg)
	}
	if !strings.Contains(msg, "is misspelled") {
		t.Errorf("the typo should get the misspelling hypothesis; got:\n%s", msg)
	}
}

// With no registry injected (nil), the message is the general one.
func TestLoad_withoutMigrationRegistryFallsBackToTheGeneralMessage(t *testing.T) {
	t.Run("CNFGO-B04: The migration advice needs both an older format and a renamed key", func(t *testing.T) {})
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = nil

	_, err := load(t, "version: 1\nlayers: {}\nqualquerCoisa: 1\n")
	if err == nil {
		t.Fatal("expected an error")
	}
	if strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("without a registry nothing says the key was renamed:\n%s", err)
	}
}

// A file ALREADY in the current format is never "an old file".
func TestLoad_fileInTheCurrentFormatNeverAdvisesMigration(t *testing.T) {
	t.Run("CNFGO-B04: The migration advice needs both an older format and a renamed key", func(t *testing.T) {})
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = func(string) bool { return true } // would say "renamed" for everything

	// The version comes from the CONSTANT: a literal would tie the test to today's format.
	_, err := load(t, fmt.Sprintf("version: %d\nlayers: {}\ntrinca_opcional: 1\n", FormatoAtualDeConfig))
	if err == nil {
		t.Fatal("expected an error")
	}
	if strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("the file is already in the current format — migrating changes nothing:\n%s", err)
	}
	// A trailing comment does not hide the version.
	_, err = load(t, fmt.Sprintf("version: %d  # current\nlayers: {}\ntrinca_opcional: 1\n", FormatoAtualDeConfig))
	if err == nil || strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("a commented version line is still the current format; got:\n%v", err)
	}
}

// A file WITHOUT `version:` is format 1, not the current format: the first files of the
// product did not declare the field.
func TestLoad_fileWithoutVersionIsFormatOne(t *testing.T) {
	t.Run("CNFGO-B03: A renamed key in an older-format file advises the migration", func(t *testing.T) {})
	original := RenamedKey
	defer func() { RenamedKey = original }()
	RenamedKey = func(k string) bool { return k == "trinca_opcional" }

	_, err := load(t, "layers:\n    spec:\n        pattern: \"**/*.spec.md\"\n"+
		"        kind: spec\n        trinca_opcional: [tested-by]\n")
	if err == nil {
		t.Fatal("the old key should be refused")
	}
	if !strings.Contains(err.Error(), "anchors migrate") {
		t.Errorf("a file without `version:` is format 1 and must MIGRATE; got:\n%s", err)
	}
}

// The ID is what the report prints and what a waiver cites: two gates with one ID make
// every citation point at two places.
func TestLoad_repeatedGateIDFails(t *testing.T) {
	t.Run("CNFGO-B06: Two gates sharing an ID fail the load", func(t *testing.T) {})
	_, err := load(t, "version: 1\ngates:\n  - name: a\n  - name: b\n    id: a\n")
	if err == nil {
		t.Fatal("gate b reuses the ID that gate a takes from its name — the load must fail")
	}
	if !strings.Contains(err.Error(), `gate "b"`) || !strings.Contains(err.Error(), `gate "a"`) {
		t.Errorf("the error must name both gates, got: %v", err)
	}
	if _, err := load(t, "version: 1\ngates:\n  - name: a\n    id: x\n  - name: b\n    id: y\n"); err != nil {
		t.Errorf("distinct IDs must load: %v", err)
	}
}

// The real case: `scope: repo` (the right value is `project`) loaded without complaint and
// three gates had nothing to measure. A wrong enum is a gate turned OFF in silence.
func TestLoad_invalidScopeFails(t *testing.T) {
	t.Run("CNFGO-B07: An unknown scope, full scope, cost, phase or perspective value fails the load", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\ngates:\n  - name: build\n    run: go build ./...\n    scope: repo\n")
	if err == nil {
		t.Fatal("an invalid scope must fail — falling back to per-node silently turns the gate off")
	}
	if !strings.Contains(err.Error(), "repo") || !strings.Contains(err.Error(), "project") {
		t.Errorf("the error must cite the wrong value AND the valid ones, got: %v", err)
	}
}

func TestLoad_invalidCostFails(t *testing.T) {
	t.Run("CNFGO-B07: An unknown scope, full scope, cost, phase or perspective value fails the load", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\ngates:\n  - name: x\n    cost: lento\n")
	if err == nil || !strings.Contains(err.Error(), "lento") {
		t.Fatalf("an invalid cost must fail citing the value, got: %v", err)
	}
}

// `scope_full` and `skip_on` were read but never validated: a typo loaded, and the gate
// silently ran the ordinary scope in the full scan, or ran in the perspective the file
// said it skipped.
func TestLoad_invalidScopeFullAndSkipOnFail(t *testing.T) {
	t.Run("CNFGO-B07: An unknown scope, full scope, cost, phase or perspective value fails the load", func(t *testing.T) {})
	for field, yaml := range map[string]string{
		"projct": "    scope_full: projct\n",
		"node":   "    scope_full: node\n",
		"chnage": "    skip_on: [chnage]\n",
	} {
		_, err := load(t, "version: 1\nlayers: {}\ngates:\n  - name: x\n"+yaml)
		if err == nil || !strings.Contains(err.Error(), field) || !strings.Contains(err.Error(), "gate \"x\"") {
			t.Errorf("%q must fail naming the gate and the value, got: %v", field, err)
		}
	}
}

func TestLoad_invalidPhaseFails(t *testing.T) {
	t.Run("CNFGO-B07: An unknown scope, full scope, cost, phase or perspective value fails the load", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\ngates:\n  - name: x\n    when: [precommit]\n")
	if err == nil || !strings.Contains(err.Error(), "precommit") {
		t.Fatalf("an invalid phase must fail, got: %v", err)
	}
}

// Counter-proof: the legitimate values of the three lists load.
func TestLoad_validEnumsLoad(t *testing.T) {
	t.Run("CNFGO-B07: An unknown scope, full scope, cost, phase or perspective value fails the load", func(t *testing.T) {})
	c, err := load(t, `version: 1
layers: {}
gates:
  - name: a
    scope: project
    cost: slow
    when: [ci]
  - name: b
    scope: batch
    cost: fast
    when: [pre-commit, pre-push, manual]
  - name: c
    scope: node
    scope_full: batch
    skip_on: [change, all]
`)
	if err != nil {
		t.Fatalf("valid enums must load: %v", err)
	}
	if len(c.Gates) != 3 {
		t.Errorf("expected 3 gates, got %d", len(c.Gates))
	}
}

// A project that never declared a workflow keeps working: the queue is local.
func TestLoad_absentWorkflowIsLocal(t *testing.T) {
	t.Run("CNFGO-B08: Local and manual modes refuse the GitHub fields", func(t *testing.T) {})
	c, err := load(t, "version: 1\nlayers: {}\n")
	if err != nil {
		t.Fatalf("a config without a workflow block must load: %v", err)
	}
	if c.GitHubMode() {
		t.Error("without a declaration, the mode is local")
	}
}

// Declaring repo/labels in local mode makes the FILE assert an integration that does not exist.
func TestLoad_localModeRefusesGitHubFields(t *testing.T) {
	t.Run("CNFGO-B08: Local and manual modes refuse the GitHub fields", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\nworkflow:\n  mode: local\n  repo: acme/exemplo\n")
	if err == nil || !strings.Contains(err.Error(), "only apply under `mode: github`") {
		t.Fatalf("expected a misplaced-field error, got %v", err)
	}
}

// `mode: manual` is the local mode without automatic issues: it loads, says so, and like
// local takes no `repo`/`labels`.
func TestLoad_manualWorkflow(t *testing.T) {
	t.Run("CNFGO-B08: Local and manual modes refuse the GitHub fields", func(t *testing.T) {})
	cfg, err := load(t, "version: 1\nlayers: {}\nworkflow:\n  mode: manual\n")
	if err != nil {
		t.Fatalf("mode: manual must load: %v", err)
	}
	if !cfg.ManualMode() || cfg.GitHubMode() {
		t.Errorf("manual mode: ManualMode=%v GitHubMode=%v", cfg.ManualMode(), cfg.GitHubMode())
	}
	if _, err := load(t, "version: 1\nlayers: {}\nworkflow:\n  mode: manual\n  repo: o/r\n"); err == nil {
		t.Error("manual mode with `repo` must be refused, as in local mode")
	}
}

// Inferring from the git remote would send writes to the wrong repository in a fork.
func TestLoad_githubRequiresRepo(t *testing.T) {
	t.Run("CNFGO-B09: GitHub mode requires an owner/name repository and a label", func(t *testing.T) {})
	t.Run("CNFGO-X01: GitHub mode never infers the repository from the git remote", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\nworkflow:\n  mode: github\n  labels: [anchors]\n")
	if err == nil || !strings.Contains(err.Error(), "repo: owner/name") {
		t.Fatalf("expected an error requiring repo, got %v", err)
	}
}

// Without a label, `anchors next` would pull product issues, which lack the cycle's shape.
func TestLoad_githubRequiresLabels(t *testing.T) {
	t.Run("CNFGO-B09: GitHub mode requires an owner/name repository and a label", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\nworkflow:\n  mode: github\n  repo: acme/exemplo\n")
	if err == nil || !strings.Contains(err.Error(), "labels") {
		t.Fatalf("expected an error requiring labels, got %v", err)
	}
}

func TestLoad_githubValid(t *testing.T) {
	t.Run("CNFGO-B09: GitHub mode requires an owner/name repository and a label", func(t *testing.T) {})
	c, err := load(t, "version: 1\nlayers: {}\nworkflow:\n  mode: github\n  repo: acme/exemplo\n  labels: [anchors]\n")
	if err != nil {
		t.Fatalf("a complete config must load: %v", err)
	}
	if !c.GitHubMode() {
		t.Error("mode: github must turn GitHubMode() on")
	}
}

func TestLoad_githubRepoWithoutSlash(t *testing.T) {
	t.Run("CNFGO-B09: GitHub mode requires an owner/name repository and a label", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\nworkflow:\n  mode: github\n  repo: exemplo\n  labels: [anchors]\n")
	if err == nil || !strings.Contains(err.Error(), "owner/name") {
		t.Fatalf("a repo without a slash must fail, got %v", err)
	}
}

// The mode is DECLARED, not guessed: a typo must fail, not silently become local.
func TestLoad_unknownModeHasNoFallback(t *testing.T) {
	t.Run("CNFGO-B10: An unknown workflow mode fails with no fallback", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers: {}\nworkflow:\n  mode: guithub\n")
	if err == nil {
		t.Fatal("an unknown mode must fail, not fall back")
	}
	if !strings.Contains(err.Error(), "no fallback") {
		t.Errorf("the message must say there is no fallback, got: %v", err)
	}
}

// A declared pattern that does not compile is refused at load, naming the field. Before,
// each gate read a failed compile as "not declared" and answered with a false cause.
func TestLoadRefusesAPatternThatDoesNotCompile(t *testing.T) {
	t.Run("CNFGO-B11: A declared pattern that does not compile fails the load naming the field", func(t *testing.T) {})
	for field, yaml := range map[string]string{
		"dialect.cursor":            "dialect:\n  family: go\n  cursor: \"(unclosed\"\n",
		"dialect.guard_patterns[0]": "dialect:\n  family: go\n  guard_patterns: [\"{{param}} (\"]\n",
		"derived.export_detect":     "derived:\n  anchor: spec\n  export_detect: \"^func\\\\s+([A-Z]\"\n",
		"derived.mock_detect":       "derived:\n  anchor: spec\n  mock_detect: \"[\"\n",
	} {
		_, err := load(t, "version: 1\n"+yaml)
		if err == nil || !strings.Contains(err.Error(), field) {
			t.Errorf("%s: a broken pattern must be refused naming the field, got %v", field, err)
		}
	}
	// A valid guard pattern with its placeholder loads.
	if _, err := load(t, "version: 1\ndialect:\n  family: go\n  guard_patterns: [\"if {{param}} == nil\"]\n"); err != nil {
		t.Errorf("a valid pattern must load: %v", err)
	}
}

func TestTestLevelAccepts(t *testing.T) {
	t.Run("CNFGO-B44: A test level's code filter accepts by allow and refuses by exclude", func(t *testing.T) {})
	none := TestLevel{}
	vr := TestLevel{Allow: []string{`-VR$`}}
	unit := TestLevel{Allow: []string{`^SMCS-`}, Exclude: []string{`-VR$`}}
	for _, c := range []struct {
		f    TestLevel
		code string
		want bool
	}{
		{none, "SMCS-B04", true}, {none, "SMCS-VR", true},
		{vr, "SMCS-B04", false}, {vr, "SMCS-VR", true},
		{unit, "SMCS-B04", true}, {unit, "SMCS-VR", false}, {unit, "OTHR-B01", false},
	} {
		if got := c.f.Accepts(c.code); got != c.want {
			t.Errorf("%+v accepts %s = %v, want %v", c.f, c.code, got, c.want)
		}
	}
	_, err := load(t, "version: 1\ngates:\n  - name: test-level-codes\n    check: test-level-codes\n    levels:\n      vr-level:\n        exclude: [\"ok\", \"(\"]\n")
	if err == nil || !strings.Contains(err.Error(), "gates[test-level-codes].levels.vr-level.exclude[1]") {
		t.Errorf("a broken filter pattern must be refused naming level, list and index, got %v", err)
	}
	if _, err := load(t, "version: 1\ngates:\n  - name: test-level-codes\n    check: test-level-codes\n    levels:\n      vr-level:\n        allow: [\"-VR$\"]\n"); err != nil {
		t.Errorf("a valid filter must load: %v", err)
	}
}

func TestLoadChecksTheTestsSource(t *testing.T) {
	t.Run("CNFGO-B45: The tests source is one source with a pattern that compiles", func(t *testing.T) {})
	_, err := load(t, "version: 1\ndialect:\n  family: go\n  tests:\n    pattern: \"x\"\n    script: \"y\"\n")
	if err == nil || !strings.Contains(err.Error(), "both a pattern and a script") {
		t.Errorf("both sources must fail the load naming the conflict, got %v", err)
	}
	_, err = load(t, "version: 1\ndialect:\n  family: go\n  tests:\n    pattern: \"it(\"\n")
	if err == nil || !strings.Contains(err.Error(), "dialect.tests.pattern") {
		t.Errorf("a tests pattern that does not compile must fail naming the field, got %v", err)
	}
	if _, err := load(t, "version: 1\ndialect:\n  family: go\n  tests:\n    script: \"go run ./scripts/tests\"\n"); err != nil {
		t.Errorf("a script alone must load: %v", err)
	}
}

func TestLoadChecksTheSupportGlobs(t *testing.T) {
	t.Run("CNFGO-B46: A layer's support globs must be valid", func(t *testing.T) {})
	_, err := load(t, "version: 1\nlayers:\n  e2e:\n    pattern: \"flows/**/*.yaml\"\n    kind: test\n    support: [\"flows/utils/**\", \"flows/[a\"]\n")
	if err == nil || !strings.Contains(err.Error(), "layers.e2e.support[1]") {
		t.Errorf("a malformed support glob must fail the load naming layer and index, got %v", err)
	}
	if _, err := load(t, "version: 1\nlayers:\n  e2e:\n    pattern: \"flows/**/*.yaml\"\n    kind: test\n    support: [\"flows/utils/**\"]\n"); err != nil {
		t.Errorf("valid support globs must load: %v", err)
	}
}

func TestGateTimeoutCeiling(t *testing.T) {
	t.Run("CNFGO-B47: A gate's timeout ceiling is a share with a default", func(t *testing.T) {})
	if got := (Gate{}).TimeoutCeilingOrDefault(); got != 0.2 {
		t.Errorf("the default ceiling is 0.2, got %v", got)
	}
	if got := (Gate{TimeoutCeiling: 0.5}).TimeoutCeilingOrDefault(); got != 0.5 {
		t.Errorf("a declared ceiling wins, got %v", got)
	}
	for _, v := range []string{"1.5", "-0.1"} {
		_, err := load(t, "version: 1\ngates:\n  - name: mutation-score\n    check: mutation-score\n    timeout_ceiling: "+v+"\n")
		if err == nil || !strings.Contains(err.Error(), "timeout_ceiling") {
			t.Errorf("ceiling %s must fail the load, got %v", v, err)
		}
	}
	if _, err := load(t, "version: 1\ngates:\n  - name: mutation-score\n    check: mutation-score\n    timeout_ceiling: 1\n"); err != nil {
		t.Errorf("a ceiling of 1 is in range: %v", err)
	}
}

func TestGateNoSignal(t *testing.T) {
	t.Run("CNFGO-B48: A gate's no_signal declares targets with their reason", func(t *testing.T) {})
	g := Gate{NoSignal: map[string]string{"cmd/**/guide_*.go": "text constants", "cmd/**/*.go": "anything under cmd"}}
	if r, ok := g.NoSignalFor("cmd/anchors/governance/guide_spec.go"); !ok || r != "anything under cmd" {
		t.Errorf("two matching globs must give the first sorted (`cmd/**/*.go`), got %q %v", r, ok)
	}
	if _, ok := g.NoSignalFor("internal/gate/gate.go"); ok {
		t.Error("a target no glob matches has a signal to measure")
	}
	for body, why := range map[string]string{
		"    no_signal:\n      \"cmd/[x\": \"reason\"\n": "invalid glob",
		"    no_signal:\n      \"cmd/*.go\": \"  \"\n":   "no reason",
	} {
		_, err := load(t, "version: 1\ngates:\n  - name: line-coverage\n    check: line-coverage\n"+body)
		if err == nil || !strings.Contains(err.Error(), why) {
			t.Errorf("%q must fail the load naming %q, got %v", body, why, err)
		}
	}
}

func TestLoad_languageIsSetAtLoad(t *testing.T) {
	t.Run("CNFGO-B12: An unsupported language fails the load", func(t *testing.T) {})
	t.Cleanup(func() { _ = i18n.Set("") })
	_, err := load(t, "version: 1\nlang: xx\n")
	if err == nil || !strings.Contains(err.Error(), "xx") || !strings.Contains(err.Error(), "pt-BR") {
		t.Fatalf("an unsupported language must fail listing the supported ones, got %v", err)
	}
	if _, err := load(t, "version: 1\nlang: es\n"); err != nil {
		t.Fatalf("a supported language must load: %v", err)
	}
	if got := i18n.Current(); got != "es" {
		t.Errorf("after loading lang es the current language is %q", got)
	}
}

func TestLoad_codeLengths(t *testing.T) {
	t.Run("CNFGO-B13: Code lengths outside two to eight fail, and valid ones reach the engine", func(t *testing.T) {})
	prev := append([]int{}, CodeLengths...)
	t.Cleanup(func() { CodeLengths = prev; SetSlotsHook(nil) })
	var hooked []int
	SetSlotsHook(func(ls []int) { hooked = ls })

	for _, bad := range []string{"1", "9"} {
		_, err := load(t, "version: 1\ncode_lengths: ["+bad+"]\n")
		if err == nil || !strings.Contains(err.Error(), "code_lengths: "+bad) {
			t.Errorf("length %s must fail naming it, got %v", bad, err)
		}
	}
	if _, err := load(t, "version: 1\ncode_lengths: [4, 5]\n"); err != nil {
		t.Fatalf("lengths 4 and 5 must load: %v", err)
	}
	if fmt.Sprint(CodeLengths) != "[4 5]" || fmt.Sprint(hooked) != "[4 5]" {
		t.Errorf("engine lengths %v, hook got %v; want [4 5] in both", CodeLengths, hooked)
	}
}

// The lengths are a process global. A file that declares none used to leave them as the
// previous load set them, so a long-running process that loaded a `[4]` project and
// then a project with no declaration read that project's 5-character codes as non-codes.
func TestLoad_noCodeLengthsRestoresTheDefault(t *testing.T) {
	t.Run("CNFGO-B42: A file with no code lengths restores the default, whatever an earlier load set", func(t *testing.T) {})
	prev := append([]int{}, CodeLengths...)
	t.Cleanup(func() { CodeLengths = prev; SetSlotsHook(nil) })
	var hooked []int
	SetSlotsHook(func(ls []int) { hooked = ls })

	if _, err := load(t, "version: 1\ncode_lengths: [4]\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := load(t, "version: 1\n"); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(CodeLengths) != "[5]" || fmt.Sprint(hooked) != "[5]" {
		t.Errorf("after a file with no code_lengths: engine %v, hook %v; want [5] in both", CodeLengths, hooked)
	}
}

func TestLoad_unreadableFile(t *testing.T) {
	t.Run("CNFGO-E01: A file that cannot be read fails the load with the read error", func(t *testing.T) {})
	c, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if c != nil || !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Load of a missing file = %v, %v; want nil and a not-exist error", c, err)
	}
}

// ── canonical gate declarations ───────────────────────────────────────────────

// The merge exists so the wording of a canonical gate lives in ONE place.
func TestLoadMergeCanonical_fillsOmittedFields(t *testing.T) {
	t.Run("CNFGO-B14: A canonical gate inherits every field the project omitted", func(t *testing.T) {})
	SetCanonicalGateResolver(func(name string) (Gate, bool) {
		if name != "gate-canonico" {
			return Gate{}, false
		}
		return Gate{
			Name: "gate-canonico", On: []string{"spec"}, Check: "algo",
			Measures: "o que ele mede", Ask: "a pergunta canônica",
			Scope: ScopeBatch, Cost: "slow", Category: "test",
			When: []string{"ci"},
		}, true
	})
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	// only the NAME: everything else must come from the catalog.
	c, err := load(t, "version: 1\ngates:\n  - name: gate-canonico\n")
	if err != nil {
		t.Fatal(err)
	}
	g := c.Gates[0]
	if g.Check != "algo" || g.Measures != "o que ele mede" || g.Ask != "a pergunta canônica" {
		t.Errorf("text fields were not inherited: %+v", g)
	}
	if g.EffectiveScope() != ScopeBatch || g.Cost != "slow" || g.Category != "test" {
		t.Errorf("scope/cost/category were not inherited: %+v", g)
	}
	if len(g.On) != 1 || g.On[0] != "spec" || len(g.When) != 1 {
		t.Errorf("lists were not inherited: %+v", g)
	}
}

// The project ALWAYS wins where it declares — the merge fills holes, it never overwrites.
func TestLoadMergeCanonical_doesNotOverwriteTheDeclared(t *testing.T) {
	t.Run("CNFGO-B15: A field the project declared wins over the canonical one", func(t *testing.T) {})
	SetCanonicalGateResolver(func(string) (Gate, bool) {
		return Gate{Name: "g", Measures: "canônico", Ask: "pergunta canônica"}, true
	})
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	c, err := load(t, "version: 1\ngates:\n  - name: g\n    measures: \"do projeto\"\n")
	if err != nil {
		t.Fatal(err)
	}
	if c.Gates[0].Measures != "do projeto" {
		t.Errorf("the project's declaration must win: %q", c.Gates[0].Measures)
	}
	if c.Gates[0].Ask != "pergunta canônica" {
		t.Errorf("the omitted field must still be inherited: %q", c.Gates[0].Ask)
	}
}

// A gate the catalog does not know passes intact — the merge must not invent fields.
func TestLoadMergeCanonical_ignoresProjectGate(t *testing.T) {
	t.Run("CNFGO-B17: A gate the catalog does not know loads untouched", func(t *testing.T) {})
	SetCanonicalGateResolver(func(string) (Gate, bool) { return Gate{}, false })
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	c, err := load(t, "version: 1\ngates:\n  - name: meu-gate\n    run: \"bash x.sh\"\n")
	if err != nil {
		t.Fatal(err)
	}
	if c.Gates[0].Check != "" || c.Gates[0].Ask != "" {
		t.Errorf("a project gate must inherit nothing: %+v", c.Gates[0])
	}
}

// The case that motivated `*bool`: the canonical gate BLOCKS and the project declared
// `blocking: false` on purpose; a plain bool would have promoted it to blocking.
func TestLoadMergeCanonical_explicitBlockingFalseWins(t *testing.T) {
	t.Run("CNFGO-B15: A field the project declared wins over the canonical one", func(t *testing.T) {})
	SetCanonicalGateResolver(func(string) (Gate, bool) {
		return Gate{Name: "g", Blocking: Bool(true)}, true
	})
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	c, err := load(t, "version: 1\ngates:\n  - name: g\n    blocking: false\n")
	if err != nil {
		t.Fatal(err)
	}
	if c.Gates[0].IsBlocking() {
		t.Error("an explicit `blocking: false` is a DECISION and cannot be overwritten by the catalog")
	}
}

// The other half: OMITTED is inherited.
func TestLoadMergeCanonical_omittedBlockingIsInherited(t *testing.T) {
	t.Run("CNFGO-B14: A canonical gate inherits every field the project omitted", func(t *testing.T) {})
	SetCanonicalGateResolver(func(string) (Gate, bool) {
		return Gate{Name: "g", Blocking: Bool(true)}, true
	})
	t.Cleanup(func() { SetCanonicalGateResolver(nil) })

	c, err := load(t, "version: 1\ngates:\n  - name: g\n")
	if err != nil {
		t.Fatal(err)
	}
	if !c.Gates[0].IsBlocking() {
		t.Error("an OMITTED field must be inherited from the catalog")
	}
}

// Run and Check are mutually exclusive dispatch mechanisms, and Check takes precedence in
// the runner: inheriting one when the project declared the other would silence the project.
func TestLoadMergeCanonical_runCheckMutuallyExclusive(t *testing.T) {
	t.Run("CNFGO-B16: Declaring run or check inherits neither of the pair", func(t *testing.T) {})
	t.Run("canonical check and project run", func(t *testing.T) {
		SetCanonicalGateResolver(func(string) (Gate, bool) {
			return Gate{Name: "tests-green", Check: "tests-pass"}, true
		})
		defer SetCanonicalGateResolver(nil)

		c, err := load(t, "version: 1\ngates:\n  - name: tests-green\n    run: \"go test ./...\"\n")
		if err != nil {
			t.Fatal(err)
		}
		g := c.Gates[0]
		if g.Run != "go test ./..." {
			t.Errorf("expected run 'go test ./...', got %q", g.Run)
		}
		if g.Check != "" {
			t.Errorf("expected empty check so run is not shadowed, got %q", g.Check)
		}
	})

	t.Run("canonical run and project check", func(t *testing.T) {
		SetCanonicalGateResolver(func(string) (Gate, bool) {
			return Gate{Name: "linter", Run: "golangci-lint run"}, true
		})
		defer SetCanonicalGateResolver(nil)

		c, err := load(t, "version: 1\ngates:\n  - name: linter\n    check: custom-linter\n")
		if err != nil {
			t.Fatal(err)
		}
		g := c.Gates[0]
		if g.Check != "custom-linter" {
			t.Errorf("expected check 'custom-linter', got %q", g.Check)
		}
		if g.Run != "" {
			t.Errorf("expected empty run when check is declared, got %q", g.Run)
		}
	})
}

// ── gate fields ───────────────────────────────────────────────────────────────

// With no declaration anywhere the gate is INFORMATIVE: silence cannot block a commit.
func TestGate_withoutBlockingIsInformative(t *testing.T) {
	t.Run("CNFGO-B18: A gate with no declared severity does not block", func(t *testing.T) {})
	if (Gate{Name: "g"}).IsBlocking() {
		t.Error("a gate with no declared severity must not block")
	}
	if !(Gate{Name: "g", Blocking: Bool(true)}).IsBlocking() {
		t.Error("a gate declared blocking must block")
	}
}

func TestGate_scopeForScan(t *testing.T) {
	t.Run("CNFGO-B19: The scope defaults to one run per target, and the full scan uses scope_full only when it is batch or project", func(t *testing.T) {})
	cases := []struct {
		g           Gate
		incr, whole string
	}{
		{Gate{}, ScopeNode, ScopeNode},
		{Gate{Scope: "repo"}, ScopeNode, ScopeNode},
		{Gate{Scope: ScopeBatch}, ScopeBatch, ScopeBatch},
		{Gate{Scope: ScopeBatch, ScopeFull: ScopeProject}, ScopeBatch, ScopeProject},
		{Gate{Scope: ScopeBatch, ScopeFull: ScopeNode}, ScopeBatch, ScopeBatch},
	}
	for _, c := range cases {
		if got := c.g.ScopeForScan(false); got != c.incr {
			t.Errorf("%+v incremental: %q, want %q", c.g, got, c.incr)
		}
		if got := c.g.ScopeForScan(true); got != c.whole {
			t.Errorf("%+v full: %q, want %q", c.g, got, c.whole)
		}
	}
}

func TestGate_runsIn(t *testing.T) {
	t.Run("CNFGO-B20: A gate with no phases runs in every phase", func(t *testing.T) {})
	all := Gate{}
	for _, p := range []string{PhasePreCommit, PhaseManual, ""} {
		if !all.RunsIn(p) {
			t.Errorf("a gate with no phases must run in %q", p)
		}
	}
	manual := Gate{When: []string{PhaseManual}}
	if manual.RunsIn(PhasePreCommit) {
		t.Error("a manual-only gate must not run in pre-commit")
	}
	if !manual.RunsIn(PhaseManual) || !manual.RunsIn("") {
		t.Error("a manual-only gate runs in manual and in the unnamed phase")
	}
}

func TestGate_skipsOn(t *testing.T) {
	t.Run("CNFGO-B21: A gate participates in every perspective unless skip_on excludes it", func(t *testing.T) {})
	if (Gate{}).SkipsOn(PerspectiveChange) || (Gate{}).SkipsOn(PerspectiveAll) {
		t.Error("with no skip_on the gate takes part in both perspectives")
	}
	g := Gate{SkipOn: []string{PerspectiveChange}}
	if !g.SkipsOn(PerspectiveChange) || g.SkipsOn(PerspectiveAll) {
		t.Errorf("skip_on [change]: change=%v all=%v; want true false", g.SkipsOn(PerspectiveChange), g.SkipsOn(PerspectiveAll))
	}
}

// Silence resolves to the canonical format: a project that already ran cannot change
// behaviour because of the new key.
func TestMutationFormat_defaultIsMTE(t *testing.T) {
	t.Run("CNFGO-B22: The mutation report format comes from the mutation-score gate, normalized, defaulting to the canonical format", func(t *testing.T) {})
	c := &Config{Gates: []Gate{{Name: "mutation-score", Check: "mutation-score"}}}
	if got := c.MutationFormat(); got != FormatMTE {
		t.Errorf("MutationFormat() = %q, want %q", got, FormatMTE)
	}
	if got := (&Config{}).MutationFormat(); got != FormatMTE {
		t.Errorf("no gate: MutationFormat() = %q, want %q", got, FormatMTE)
	}
}

func TestMutationFormat_declared(t *testing.T) {
	t.Run("CNFGO-B22: The mutation report format comes from the mutation-score gate, normalized, defaulting to the canonical format", func(t *testing.T) {})
	c := &Config{Gates: []Gate{
		{Name: "line-coverage", Check: "line-coverage"},
		{Name: "mutation-score", Check: "mutation-score", Format: "gremlins"},
	}}
	if got := c.MutationFormat(); got != FormatGremlins {
		t.Errorf("MutationFormat() = %q, want %q", got, FormatGremlins)
	}
}

// Case and spaces cannot decide the format — the yaml is written by hand.
func TestMutationFormat_normalizes(t *testing.T) {
	t.Run("CNFGO-B22: The mutation report format comes from the mutation-score gate, normalized, defaulting to the canonical format", func(t *testing.T) {})
	c := &Config{Gates: []Gate{{Name: "mutation-score", Format: "  GREMLINS  "}}}
	if got := c.MutationFormat(); got != FormatGremlins {
		t.Errorf("MutationFormat() = %q, want %q", got, FormatGremlins)
	}
}

// ANOTHER gate's `format:` does not apply to mutation.
func TestMutationFormat_ignoresOtherGates(t *testing.T) {
	t.Run("CNFGO-B22: The mutation report format comes from the mutation-score gate, normalized, defaulting to the canonical format", func(t *testing.T) {})
	c := &Config{Gates: []Gate{{Name: "line-coverage", Format: "gremlins"}}}
	if got := c.MutationFormat(); got != FormatMTE {
		t.Errorf("MutationFormat() = %q, want %q (the format belongs to another gate)", got, FormatMTE)
	}
}

func TestGate_checksSectionLanguage(t *testing.T) {
	t.Run("CNFGO-B23: Section language is checked unless the gate turns it off", func(t *testing.T) {})
	if !(Gate{}).ChecksSectionLanguage() || !(Gate{EnforceSectionLanguage: Bool(true)}).ChecksSectionLanguage() {
		t.Error("an undeclared or true setting must check the section language")
	}
	if (Gate{EnforceSectionLanguage: Bool(false)}).ChecksSectionLanguage() {
		t.Error("a gate set false must not check the section language")
	}
}

// ── workflow defaults ─────────────────────────────────────────────────────────

// The default assumes no flow beyond what GitHub gives a new repository.
func TestBranch_defaultAssumesNoFlow(t *testing.T) {
	t.Run("CNFGO-B24: The integration branch defaults to main, and a non-main integration branch protects main too", func(t *testing.T) {})
	var w *Workflow
	if b := w.IntegrationBranchOrDefault(); b != "main" {
		t.Errorf("with no config the integration branch is `main`, got %q", b)
	}
	if p := w.ProtectedBranchesOrDefault(); len(p) != 1 || p[0] != "main" {
		t.Errorf("with no config only main is protected, got %v", p)
	}
}

// The develop→staging→main flow is a project's CONFIGURATION, not a framework rule.
func TestBranch_threeBranchFlowIsConfigurable(t *testing.T) {
	t.Run("CNFGO-B24: The integration branch defaults to main, and a non-main integration branch protects main too", func(t *testing.T) {})
	w := &Workflow{
		IntegrationBranch: "develop",
		ProtectedBranches: []string{"develop", "staging", "main"},
	}
	if b := w.IntegrationBranchOrDefault(); b != "develop" {
		t.Errorf("work arrives in develop, got %q", b)
	}
	if p := w.ProtectedBranchesOrDefault(); len(p) != 3 {
		t.Errorf("the three branches are gates, got %v", p)
	}
}

// Declaring only the integration branch implies protecting main too.
func TestBranch_nonMainIntegrationProtectsMainToo(t *testing.T) {
	t.Run("CNFGO-B24: The integration branch defaults to main, and a non-main integration branch protects main too", func(t *testing.T) {})
	w := &Workflow{IntegrationBranch: "develop"}
	if p := w.ProtectedBranchesOrDefault(); fmt.Sprint(p) != "[develop main]" {
		t.Errorf("expected develop and main protected, got %v", p)
	}
}

func TestWorkflow_requiredApprovals(t *testing.T) {
	t.Run("CNFGO-B25: One approval is required unless the project declares another number, zero included", func(t *testing.T) {})
	var nilW *Workflow
	zero, two := 0, 2
	for _, c := range []struct {
		w    *Workflow
		want int
	}{{nilW, 1}, {&Workflow{}, 1}, {&Workflow{RequiredApprovals: &zero}, 0}, {&Workflow{RequiredApprovals: &two}, 2}} {
		if got := c.w.RequiredApprovalsOrDefault(); got != c.want {
			t.Errorf("RequiredApprovalsOrDefault(%+v) = %d, want %d", c.w, got, c.want)
		}
	}
}

// The DEFAULT is to warn: a stale pipeline still does the old work, and failing it trades
// "does less than it should" for "does nothing". Inverting the default would be silent.
func TestWorkflow_staleAndManualIngestBlockOnlyWhenAsked(t *testing.T) {
	t.Run("CNFGO-B26: A stale pipeline or a manual ingest blocks only when the project asks", func(t *testing.T) {})
	var none *Workflow
	var empty Workflow
	if none.PipelineVelhoBarra() || none.IngestManualBarra() {
		t.Error("without a `workflow:` section the default is to warn")
	}
	if empty.PipelineVelhoBarra() || empty.IngestManualBarra() {
		t.Error("undeclared, the default is to WARN — failing the CI of whoever did not ask is a regression")
	}
	asked := &Workflow{StalePipelineBlocks: true, ManualIngestBlocks: true}
	if !asked.PipelineVelhoBarra() || !asked.IngestManualBarra() {
		t.Error("a project that declared both blocks wants them to block")
	}
}

// ── freeze ────────────────────────────────────────────────────────────────────

// ABSENT means ENABLED — `Enabled` is a pointer so a project that never declared it is
// not born frozen.
func TestFrozen_absentDoesNotFreeze(t *testing.T) {
	t.Run("CNFGO-B27: Only an explicit enabled false freezes the project", func(t *testing.T) {})
	if (&Config{}).Frozen() {
		t.Fatal("a config without `enabled` is frozen — every project that never declared it would be born stopped")
	}
}

func TestFrozen_onlyExplicitFalseFreezes(t *testing.T) {
	t.Run("CNFGO-B27: Only an explicit enabled false freezes the project", func(t *testing.T) {})
	off, on := false, true
	if !(&Config{Enabled: &off}).Frozen() {
		t.Error("`enabled: false` did not freeze")
	}
	if (&Config{Enabled: &on}).Frozen() {
		t.Error("`enabled: true` froze")
	}
}

// Nil-safe: a config that did not load freezes nothing.
func TestFrozen_nilConfigDoesNotFreeze(t *testing.T) {
	t.Run("CNFGO-B27: Only an explicit enabled false freezes the project", func(t *testing.T) {})
	var c *Config
	if c.Frozen() {
		t.Fatal("a nil config answered frozen")
	}
	if c.FreezeReasonText() != "" {
		t.Fatal("a nil config returned a reason")
	}
}

// An ABSENT reason does not become silence: the message asks whoever froze for it.
func TestFrozen_withoutReasonAsksForIt(t *testing.T) {
	t.Run("CNFGO-B28: The freeze reason is shown trimmed, and a missing one asks for freeze_reason", func(t *testing.T) {})
	off := false
	m := (&Config{Enabled: &off}).FreezeReasonText()
	if m == "" {
		t.Fatal("frozen without `freeze_reason` returned an empty reason")
	}
	if !strings.Contains(m, "freeze_reason") {
		t.Errorf("the message does not say where to write the reason: %q", m)
	}
}

func TestFrozen_declaredReasonIsShown(t *testing.T) {
	t.Run("CNFGO-B28: The freeze reason is shown trimmed, and a missing one asks for freeze_reason", func(t *testing.T) {})
	off := false
	c := &Config{Enabled: &off, FreezeReason: "  o plano 0002 aponta para spec inexistente  "}
	// Surrounding spaces are not content: the text goes to a terminal message.
	if got := c.FreezeReasonText(); got != "o plano 0002 aponta para spec inexistente" {
		t.Errorf("reason: %q", got)
	}
}

// A blank reason counts as absent.
func TestFrozen_blankReasonCountsAsAbsent(t *testing.T) {
	t.Run("CNFGO-B28: The freeze reason is shown trimmed, and a missing one asks for freeze_reason", func(t *testing.T) {})
	off := false
	if m := (&Config{Enabled: &off, FreezeReason: "   \n  "}).FreezeReasonText(); !strings.Contains(m, "freeze_reason") {
		t.Errorf("a blank reason did not fall into the notice: %q", m)
	}
}

// ── project vocabulary ────────────────────────────────────────────────────────

func TestSectionTitle_layerThenProjectThenFramework(t *testing.T) {
	t.Run("CNFGO-B29: A section title comes from the layer, then the project, then the framework", func(t *testing.T) {})
	c := &Config{
		SectionTitles: SectionTitles{"rules": "Regras"},
		Layers: map[string]Layer{
			"screen": {SectionTitles: SectionTitles{"rules": "Comportamentos"}},
			"blank":  {SectionTitles: SectionTitles{"rules": "  "}},
		},
	}
	for _, tc := range []struct{ key, def, layer, want string }{
		{"rules", "Rules", "screen", "Comportamentos"},
		{"rules", "Rules", "blank", "Regras"},
		{"rules", "Rules", "other", "Regras"},
		{"states", "States", "screen", "States"},
	} {
		if got := c.SectionTitle(tc.key, tc.def, tc.layer); got != tc.want {
			t.Errorf("SectionTitle(%q, layer %q) = %q, want %q", tc.key, tc.layer, got, tc.want)
		}
	}
}

func TestPlaceholders_defaultToTODO(t *testing.T) {
	t.Run("CNFGO-B30: Placeholder markers default to the templates' marker word", func(t *testing.T) {})
	for _, tc := range []struct {
		markers []string
		want    string
	}{
		{[]string{" FIXME ", " "}, "[FIXME]"},
		{[]string{"", "  "}, "[TODO]"},
		{nil, "[TODO]"},
	} {
		if got := fmt.Sprint((&Config{PlaceholderMarkers: tc.markers}).Placeholders()); got != tc.want {
			t.Errorf("Placeholders(%q) = %s, want %s", tc.markers, got, tc.want)
		}
	}
}

func TestRuleLetters_declaredOrCanonical(t *testing.T) {
	t.Run("CNFGO-B31: Rule letters come from the declared rule types, or the canonical set", func(t *testing.T) {})
	c := &Config{RuleTypes: []RuleType{{Letter: "b"}, {Letter: "B"}, {Letter: "XY"}, {Letter: " s "}, {Letter: ""}}}
	if got := c.RuleLetters(); got != "BS" {
		t.Errorf("RuleLetters() = %q, want BS", got)
	}
	if got := (&Config{RuleTypes: []RuleType{{Letter: "XY"}}}).RuleLetters(); got != DefaultRuleLetters || got != "SRVAXBNMDEIQFG" {
		t.Errorf("with no valid letter RuleLetters() = %q, want the canonical set", got)
	}
}

func TestTagLetters_everyLetterDeclaringTheTag(t *testing.T) {
	t.Run("CNFGO-B32: A scenario tag maps to every letter that declares it", func(t *testing.T) {})
	c := &Config{RuleTypes: []RuleType{
		{Letter: "s", Tags: []string{"@estado-dado"}},
		{Letter: "V", Tags: []string{"@visual", "@estado-dado"}},
		{Letter: "B", Tags: []string{"@comportamento"}},
	}}
	if got, ok := c.TagLetters(" @Estado-Dado "); !ok || fmt.Sprint(got) != "[S V]" {
		t.Errorf("TagLetters = %v, %v; want [S V], true", got, ok)
	}
	if got, ok := c.TagLetters("@smoke"); ok || got != nil {
		t.Errorf("an undeclared tag must be unknown, got %v, %v", got, ok)
	}
}

func TestCodeLengthPattern_singleAndContiguous(t *testing.T) {
	t.Run("CNFGO-B33: The code length pattern matches exactly the declared lengths, contiguous or not", func(t *testing.T) {})
	prev := append([]int{}, CodeLengths...)
	t.Cleanup(func() { CodeLengths = prev })
	for _, tc := range []struct {
		lengths []int
		want    string // which lengths among 2..9 match
	}{
		{[]int{5}, "[5]"},
		{[]int{5, 4}, "[4 5]"},
		// Non-contiguous: the pattern once matched 5 and 7 for [4, 6], and 6 and 8 for
		// [5, 7], because the class before it lengthened each alternative by one.
		{[]int{6, 4}, "[4 6]"},
		{[]int{5, 7}, "[5 7]"},
		{[]int{3, 5, 8}, "[3 5 8]"},
	} {
		CodeLengths = tc.lengths
		re := regexp.MustCompile(`^[A-Z0-9]` + CodeLengthPattern() + `$`)
		var matched []int
		for n := 2; n <= 9; n++ {
			if re.MatchString(strings.Repeat("A", n)) {
				matched = append(matched, n)
			}
		}
		if fmt.Sprint(matched) != tc.want {
			t.Errorf("lengths %v: matched %v, want %s", tc.lengths, matched, tc.want)
		}
	}
}

func TestRequiresCodeIn_matchesTheDeclaredSectionTitle(t *testing.T) {
	t.Run("CNFGO-B41: A rule type catalogues the sections it declares, ignoring case and surrounding spaces", func(t *testing.T) {})
	rt := RuleType{SectionsRequireCode: []string{"Business Rules"}}
	if !rt.RequiresCodeIn(" business rules ") {
		t.Error("the declared section, in another case and with spaces around it, must require a code")
	}
	if rt.RequiresCodeIn("Notes") {
		t.Error("a section the rule type does not declare must not require a code")
	}
	if rt.RequiresCodeIn("Business") {
		t.Error("a part of the declared title is not the declared title")
	}
}

// ── suite selection ───────────────────────────────────────────────────────────

// The real monorepo that motivated the two axes: `unit` exists in TWO workspaces, and
// `integration` in only one.
func suites() []Suite {
	return []Suite{
		{Workspace: "backend", Layer: "unit", Run: "b-unit"},
		{Workspace: "backend", Layer: "integration", Run: "b-it"},
		{Workspace: "mobile", Layer: "unit", Run: "m-unit"},
		{Workspace: "mobile", Layer: "e2e", Run: "m-e2e"},
	}
}

func commands(sel []Suite) string {
	out := make([]string, 0, len(sel))
	for _, s := range sel {
		out = append(out, s.Run)
	}
	return strings.Join(out, ",")
}

// The default of a declarative command is "what was declared", never "nothing".
func TestSuites_noFilterRunsEverything(t *testing.T) {
	t.Run("CNFGO-B34: With no filter every suite is selected", func(t *testing.T) {})
	sel, missing := SelecionaSuites(suites(), nil, nil, nil)
	if len(sel) != 4 || len(missing) != 0 {
		t.Fatalf("with no filter all 4 should run; got %d (missing %v)", len(sel), missing)
	}
}

// The reason the axes are separate: `anchors test unit` means "everyone's unit".
func TestSuites_layerCrossesWorkspaces(t *testing.T) {
	t.Run("CNFGO-B35: The filter axes intersect", func(t *testing.T) {})
	sel, _ := SelecionaSuites(suites(), []string{"unit"}, nil, nil)
	if got := commands(sel); got != "b-unit,m-unit" {
		t.Errorf("unit should take both workspaces; got %q", got)
	}
}

// The opposite axis: "everything in backend", with no layer named.
func TestSuites_workspaceCrossesLayers(t *testing.T) {
	t.Run("CNFGO-B35: The filter axes intersect", func(t *testing.T) {})
	sel, _ := SelecionaSuites(suites(), nil, []string{"backend"}, nil)
	if got := commands(sel); got != "b-unit,b-it" {
		t.Errorf("backend should take both of its layers; got %q", got)
	}
}

// The intersection is the everyday case ("only the backend's unit, which is fast").
func TestSuites_theTwoAxesCombine(t *testing.T) {
	t.Run("CNFGO-B35: The filter axes intersect", func(t *testing.T) {})
	sel, _ := SelecionaSuites(suites(), []string{"unit"}, []string{"backend"}, nil)
	if got := commands(sel); got != "b-unit" {
		t.Errorf("wrong intersection; got %q", got)
	}
}

func TestSuites_severalValuesPerAxis(t *testing.T) {
	t.Run("CNFGO-B35: The filter axes intersect", func(t *testing.T) {})
	sel, _ := SelecionaSuites(suites(), []string{"unit", "e2e"}, []string{"mobile"}, nil)
	if got := commands(sel); got != "m-unit,m-e2e" {
		t.Errorf("wanted both mobile layers; got %q", got)
	}
}

// Whoever wrote anchors.yaml decided unit runs before e2e, which is often a real dependency.
func TestSuites_orderIsTheFilesNotTheUsers(t *testing.T) {
	t.Run("CNFGO-B36: Selected suites keep the order of the file", func(t *testing.T) {})
	sel, _ := SelecionaSuites(suites(), []string{"e2e", "unit"}, nil, nil)
	if got := commands(sel); got != "b-unit,m-unit,m-e2e" {
		t.Errorf("the order is the declaration's; got %q", got)
	}
}

// A wrong name is the common mistake, and the user needs to know on WHICH axis.
func TestSuites_missingNameComesLabelledAndTogether(t *testing.T) {
	t.Run("CNFGO-B38: A name missing from the file is reported with its axis, and an empty combination is not a missing name", func(t *testing.T) {})
	_, missing := SelecionaSuites(suites(), []string{"unit", "smoke"}, []string{"web"}, nil)
	if len(missing) != 2 {
		t.Fatalf("both absences come together; got %v", missing)
	}
	if missing[0] != i18n.T("config.suite.missing_layer", `"smoke"`) {
		t.Errorf("the first should be the smoke layer; got %q", missing[0])
	}
	if !strings.Contains(missing[1], "workspace") || !strings.Contains(missing[1], "web") {
		t.Errorf("the second should be the web workspace; got %q", missing[1])
	}
}

// `integration` and `mobile` exist, but not TOGETHER: reporting a missing layer would send
// the user to fix a name that is right.
func TestSuites_emptyCombinationIsNotAWrongName(t *testing.T) {
	t.Run("CNFGO-B38: A name missing from the file is reported with its axis, and an empty combination is not a missing name", func(t *testing.T) {})
	sel, missing := SelecionaSuites(suites(), []string{"integration"}, []string{"mobile"}, nil)
	if len(missing) != 0 {
		t.Errorf("both names exist; nothing to report as missing. got %v", missing)
	}
	if len(sel) != 0 {
		t.Errorf("the combination does not exist, so nothing runs; got %q", commands(sel))
	}
}

func TestSuites_ignoreCaseAndSpace(t *testing.T) {
	t.Run("CNFGO-B37: Suite names match ignoring case and surrounding spaces", func(t *testing.T) {})
	sel, missing := SelecionaSuites(suites(), []string{" Unit "}, []string{"BackEnd"}, nil)
	if len(sel) != 1 || len(missing) != 0 {
		t.Errorf("case and space should not separate; got %q (missing %v)", commands(sel), missing)
	}
}

// A single-package project has no workspace axis.
func TestSuites_workspaceIsOptional(t *testing.T) {
	t.Run("CNFGO-B35: The filter axes intersect", func(t *testing.T) {})
	t.Run("CNFGO-B39: The declared vocabulary lists each name once, in file order", func(t *testing.T) {})
	simple := []Suite{{Layer: "unit", Run: "u"}, {Layer: "e2e", Run: "e"}}
	sel, missing := SelecionaSuites(simple, []string{"unit"}, nil, nil)
	if len(sel) != 1 || sel[0].Run != "u" || len(missing) != 0 {
		t.Errorf("a suite without a workspace should work; got %q / %v", commands(sel), missing)
	}
	if ws := DeclaredWorkspaces(simple); len(ws) != 0 {
		t.Errorf("with no workspace declared the list is empty; got %v", ws)
	}
}

// No fixed list of layers or workspaces: validating against unit/integration/e2e would be
// the framework deciding its users' vocabulary.
func TestSuites_vocabularyIsTheProjects(t *testing.T) {
	t.Run("CNFGO-X02: Suite names are the project's vocabulary, never a fixed list", func(t *testing.T) {})
	own := []Suite{
		{Workspace: "cobranca", Layer: "contrato", Run: "c"},
		{Workspace: "cobranca", Layer: "carga", Run: "k"},
	}
	sel, missing := SelecionaSuites(own, []string{"carga"}, []string{"cobranca"}, nil)
	if len(sel) != 1 || sel[0].Run != "k" || len(missing) != 0 {
		t.Errorf("the project's own vocabulary should be accepted; got %q / %v", commands(sel), missing)
	}
}

// The lists are what is shown to whoever mistyped; `unit` twice would read like a dump.
func TestSuites_listsDoNotRepeatNames(t *testing.T) {
	t.Run("CNFGO-B39: The declared vocabulary lists each name once, in file order", func(t *testing.T) {})
	if got := strings.Join(DeclaredLayers(suites()), ","); got != "unit,integration,e2e" {
		t.Errorf("layers repeated or out of order: %q", got)
	}
	if got := strings.Join(DeclaredWorkspaces(suites()), ","); got != "backend,mobile" {
		t.Errorf("workspaces repeated or out of order: %q", got)
	}
}

// The scope is the SAME unit measured with different reach: against its own test
// (`isolated`) and against its importers' (`full`).
func mutationSuites() []Suite {
	return []Suite{
		{Workspace: "backend", Layer: "unit", Scope: "full", Run: "b-full"},
		{Workspace: "backend", Layer: "unit", Scope: "isolated", Run: "b-iso"},
		{Workspace: "mobile", Layer: "unit", Scope: "full", Run: "m-full"},
	}
}

// With no scope filter both run: the PAIR is what gives the gate its coupling reading.
func TestSuites_noScopeRunsBoth(t *testing.T) {
	t.Run("CNFGO-B34: With no filter every suite is selected", func(t *testing.T) {})
	sel, _ := SelecionaSuites(mutationSuites(), nil, []string{"backend"}, nil)
	if got := commands(sel); got != "b-full,b-iso" {
		t.Errorf("without --scope both backend scopes run; got %q", got)
	}
}

func TestSuites_scopeFiltersOne(t *testing.T) {
	t.Run("CNFGO-B35: The filter axes intersect", func(t *testing.T) {})
	sel, _ := SelecionaSuites(mutationSuites(), nil, nil, []string{"isolated"})
	if got := commands(sel); got != "b-iso" {
		t.Errorf("only isolated; got %q", got)
	}
}

// The three axes cross, as layer and workspace already did.
func TestSuites_scopeCombinesWithTheOtherAxes(t *testing.T) {
	t.Run("CNFGO-B35: The filter axes intersect", func(t *testing.T) {})
	sel, _ := SelecionaSuites(mutationSuites(), []string{"unit"}, []string{"backend"}, []string{"full"})
	if got := commands(sel); got != "b-full" {
		t.Errorf("intersection of the three axes; got %q", got)
	}
}

// The two backend suites have IDENTICAL (workspace, layer) and differ only in scope:
// without the third axis the isolated one would be unreachable.
func TestSuites_sameLayerAndWorkspaceWithDifferentScopesDoNotCollide(t *testing.T) {
	t.Run("CNFGO-B35: The filter axes intersect", func(t *testing.T) {})
	sel, missing := SelecionaSuites(mutationSuites(), []string{"unit"}, []string{"backend"}, []string{"isolated"})
	if len(missing) != 0 {
		t.Fatalf("nothing missing here; got %v", missing)
	}
	if len(sel) != 1 || sel[0].Run != "b-iso" {
		t.Errorf("the scope must break the tie; got %q", commands(sel))
	}
}

func TestSuites_missingScopeComesLabelled(t *testing.T) {
	t.Run("CNFGO-B38: A name missing from the file is reported with its axis, and an empty combination is not a missing name", func(t *testing.T) {})
	_, missing := SelecionaSuites(mutationSuites(), nil, nil, []string{"parcial"})
	if len(missing) != 1 || missing[0] != i18n.T("config.suite.missing_scope", `"parcial"`) {
		t.Errorf("the absence should be labelled as a scope; got %v", missing)
	}
}

// `tests:` declares no scope, and the axis cannot exclude suites that do not take part in it.
func TestSuites_suiteWithoutScopeStaysWhenNobodyFilters(t *testing.T) {
	t.Run("CNFGO-B34: With no filter every suite is selected", func(t *testing.T) {})
	sel, _ := SelecionaSuites(suites(), []string{"unit"}, nil, nil)
	if len(sel) != 2 {
		t.Errorf("suites with no scope still count; got %q", commands(sel))
	}
}

// ── derived file patterns ─────────────────────────────────────────────────────

func TestDerived_patternsReplaceCodeAndKeepTheRest(t *testing.T) {
	t.Run("CNFGO-B40: Patterns replace the code template and keep the other derived files", func(t *testing.T) {})
	var d Derived
	if err := yaml.Unmarshal([]byte(`
anchor: spec
files:
  code: "{{dir}}/{{name}}.ts"
  feature: "{{dir}}/{{name}}.feature"
  patterns:
    - "tsconfig.base.json"
    - "packages/*/tsconfig.json"
`), &d); err != nil {
		t.Fatal(err)
	}
	got := d.PadroesDe()
	if fmt.Sprint(got["code"]) != "[tsconfig.base.json packages/*/tsconfig.json]" {
		t.Errorf("`patterns` should replace `code`, got %v", got["code"])
	}
	if fmt.Sprint(got["feature"]) != "[{{dir}}/{{name}}.feature]" {
		t.Errorf("`feature` should survive, got %v", got["feature"])
	}
	if _, isLayer := got[PatternKey]; isLayer {
		t.Error("`patterns` is not a derived file and must not remain in the map")
	}
	var plain Derived
	if err := yaml.Unmarshal([]byte("anchor: spec\nfiles:\n  code: \"x.ts\"\n"), &plain); err != nil {
		t.Fatal(err)
	}
	if got := plain.PadroesDe()["code"]; fmt.Sprint(got) != "[x.ts]" {
		t.Errorf("without patterns, `files` rules: %v", got)
	}
}

// ── invariants and I/O failures ───────────────────────────────────────────────

func TestSave_whatIsWrittenLoadsBack(t *testing.T) {
	t.Run("CNFGO-I01: What Save writes, Load reads back", func(t *testing.T) {})
	c, err := load(t, "version: 1\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\ngates:\n  - name: g\n    blocking: true\n")
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := Save(c, p); err != nil {
		t.Fatal(err)
	}
	back, err := Load(p)
	if err != nil {
		t.Fatalf("the saved file must load: %v", err)
	}
	if back.Layers["spec"].Kind != "spec" || len(back.Gates) != 1 || !back.Gates[0].IsBlocking() {
		t.Errorf("reloaded config differs: layers %+v gates %+v", back.Layers, back.Gates)
	}
}

// A `lang: en` project got some load errors in Portuguese: the cost, the three GitHub
// ones and the unknown mode were hardcoded Portuguese, code_lengths was half of each, and
// the gate-ID and workflow checks ran before the language was set, so they came out in
// the language of the previous load.
func TestLoad_refusalsFollowTheProjectLanguage(t *testing.T) {
	t.Run("CNFGO-B43: Every load refusal and the header Save writes are in the project's language", func(t *testing.T) {})
	t.Cleanup(func() { _ = i18n.Set(i18n.Default) })
	refusals := map[string]string{
		"cost":         "gates:\n  - name: x\n    cost: lento\n",
		"no repo":      "workflow:\n  mode: github\n  labels: [a]\n",
		"no label":     "workflow:\n  mode: github\n  repo: o/r\n",
		"bad repo":     "workflow:\n  mode: github\n  repo: r\n  labels: [a]\n",
		"unknown mode": "workflow:\n  mode: guithub\n",
		"local fields": "workflow:\n  mode: local\n  repo: o/r\n",
		"length":       "code_lengths: [9]\n",
		"duplicate id": "gates:\n  - name: a\n  - name: b\n    id: a\n",
	}
	portuguese := regexp.MustCompile(`desconhecido|exige|identificador|não|já é usado|só valem`)
	for name, body := range refusals {
		// A pt-BR load first, so a refusal emitted before the language is set shows up.
		if _, err := load(t, "version: 1\nlang: pt-BR\n"); err != nil {
			t.Fatal(err)
		}
		_, err := load(t, "version: 1\nlang: en\n"+body)
		if err == nil || portuguese.MatchString(err.Error()) {
			t.Errorf("%s, lang en: %v; want an English refusal", name, err)
		}
		_, err = load(t, "version: 1\nlang: pt-BR\n"+body)
		if err == nil || !portuguese.MatchString(err.Error()) {
			t.Errorf("%s, lang pt-BR: %v; want a Portuguese refusal", name, err)
		}
	}

	for lang, want := range map[string]string{
		"":      "# anchors.yaml — Anchors project configuration\n",
		"en":    "# anchors.yaml — Anchors project configuration\n",
		"pt-BR": "# anchors.yaml — configuração do projeto Anchors\n",
	} {
		p := filepath.Join(t.TempDir(), "anchors.yaml")
		if err := Save(&Config{Lang: lang}, p); err != nil {
			t.Fatal(err)
		}
		data, _ := os.ReadFile(p)
		if !strings.HasPrefix(string(data), want) {
			t.Errorf("Save with lang %q wrote %q; want it to start with %q", lang, data, want)
		}
	}
}

func TestSave_unwritablePath(t *testing.T) {
	t.Run("CNFGO-E02: Save to a path that cannot be written returns the write error", func(t *testing.T) {})
	err := Save(&Config{}, filepath.Join(t.TempDir(), "no-such-dir", "anchors.yaml"))
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Save into a missing directory = %v, want a not-exist error", err)
	}
}

// Counter-proof: KnownFields cannot turn a legitimate config into an error. If this breaks
// when a field is added, the field is missing from the struct — do not relax the mode.
func TestLoad_validConfigKeepsLoading(t *testing.T) {
	t.Run("CNFGO-I02: A configuration with every key known keeps loading", func(t *testing.T) {})
	c, err := load(t, `version: 1
layers:
  spec:
    pattern: "**/*.spec.md"
    kind: spec
    tags: [spec]
governs:
  - from: GUIDE.md
    governs: spec
gates:
  - name: spec-completa
    blocking: false
`)
	if err != nil {
		t.Fatalf("a valid config must load: %v", err)
	}
	if len(c.Layers) != 1 || len(c.Governs) != 1 || len(c.Gates) != 1 {
		t.Errorf("config loaded incomplete: %d layers, %d governs, %d gates",
			len(c.Layers), len(c.Governs), len(c.Gates))
	}
}

// `touch.pre_commit: false` is read from the YAML; with no `touch` block the field stays nil.
func TestLoad_touchPreCommitKey(t *testing.T) {
	t.Run("CNFGO-I02: A configuration with every key known keeps loading", func(t *testing.T) {})
	cfg, err := load(t, "version: 1\nlayers: {}\ntouch:\n  pre_commit: false\n  exclude: [anchors.graph.yaml]\n")
	if err != nil {
		t.Fatalf("touch block must load: %v", err)
	}
	if cfg.Touch == nil || cfg.Touch.PreCommit == nil || *cfg.Touch.PreCommit || len(cfg.Touch.Exclude) != 1 {
		t.Errorf("touch = %+v, want pre_commit false and one exclude", cfg.Touch)
	}
	cfg, err = load(t, "version: 1\nlayers: {}\n")
	if err != nil || cfg.Touch != nil {
		t.Errorf("no touch block: Touch = %+v, err %v", cfg.Touch, err)
	}
}

func TestReaders_nilConfigAnswersDefaults(t *testing.T) {
	t.Run("CNFGO-I03: Every reader answers its default on a nil configuration", func(t *testing.T) {})
	var c *Config
	var w *Workflow
	if c.RuleLetters() != DefaultRuleLetters || fmt.Sprint(c.Placeholders()) != "[TODO]" ||
		c.SectionTitle("rules", "Rules", "screen") != "Rules" {
		t.Error("a nil config must answer the vocabulary defaults")
	}
	if tags, ok := c.TagLetters("@estado"); ok || tags != nil {
		t.Error("a nil config knows no tag")
	}
	if c.Frozen() || c.GitHubMode() || c.ManualMode() || c.RouteRegistry() != nil {
		t.Error("a nil config is not frozen, not in GitHub or manual mode, and has no route registry")
	}
	if w.IntegrationBranchOrDefault() != "main" || w.RequiredApprovalsOrDefault() != 1 ||
		w.PipelineVelhoBarra() || w.IngestManualBarra() {
		t.Error("a nil workflow must answer main, one approval, and no blocking switches")
	}
}
