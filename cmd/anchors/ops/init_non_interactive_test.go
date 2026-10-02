// @anchors
//   ref: ININT

package ops

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// runNonInteractive runs `anchors init --non-interactive` with the given flags and returns
// the error, the JSON it printed (decoded), and the raw output.
func runNonInteractive(t *testing.T, root string, flags ...string) (error, map[string]any, string) {
	t.Helper()
	cmd := newInitCmd()
	cmd.SetArgs(append([]string{"--root", root, "--non-interactive"}, flags...))
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	var err error
	out := captureStdout(t, func() { err = cmd.Execute() })
	var doc map[string]any
	if strings.TrimSpace(out) != "" {
		if jerr := json.Unmarshal([]byte(out), &doc); jerr != nil {
			t.Fatalf("the output is not one JSON document: %v\n%s", jerr, out)
		}
	}
	return err, doc, out
}

func loadWritten(t *testing.T, root string) *config.Config {
	t.Helper()
	cfg, err := config.Load(filepath.Join(root, config.DefaultFile))
	if err != nil {
		t.Fatalf("anchors.yaml was not written or does not load: %v", err)
	}
	return cfg
}

func layerNames(cfg *config.Config) []string {
	var out []string
	for n := range cfg.Layers {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// Without answers the command ASKS: it returns the questions and writes nothing. Writing
// the defaults of a new project would produce an anchors.yaml that governs nothing.
func TestNonInteractiveWithoutAnswersOnlyAsks(t *testing.T) {
	t.Run("ININT-B01: Without answers the command only asks", func(t *testing.T) {})
	t.Run("ININT-X01: The non-interactive mode never prompts", func(t *testing.T) {})
	root := t.TempDir()
	err, doc, out := runNonInteractive(t, root)
	if err != nil {
		t.Fatalf("asking is not an error: %v", err)
	}
	if doc["escrito"] != false {
		t.Errorf("escrito = %v, want false:\n%s", doc["escrito"], out)
	}
	qs, _ := doc["perguntas"].([]any)
	ids := map[string]bool{}
	for _, q := range qs {
		ids[q.(map[string]any)["id"].(string)] = true
	}
	for _, id := range []string{"artifacts", "colocation", "workflow"} {
		if !ids[id] {
			t.Errorf("the question %q is missing from %v", id, ids)
		}
	}
	if doc["precisa_descobrir"] != true {
		t.Errorf("an empty directory needs the DISCOVER phase, got %v", doc["precisa_descobrir"])
	}
	if _, serr := os.Stat(filepath.Join(root, config.DefaultFile)); serr == nil {
		t.Error("anchors.yaml was written without a single answer")
	}
}

// With answers, each one reaches the file: the artifact layers, co-location, the gates,
// the manual workflow and the header guide.
func TestNonInteractiveAppliesTheAnswers(t *testing.T) {
	t.Run("ININT-B02: The given answers reach the configuration", func(t *testing.T) {})
	t.Run("ININT-B09: The success document names the file and the next step", func(t *testing.T) {})
	root := t.TempDir()
	err, doc, out := runNonInteractive(t, root,
		"--artifacts=spec,feature,test", "--colocation", "--workflow=manual")
	if err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if doc["escrito"] != true {
		t.Fatalf("escrito = %v:\n%s", doc["escrito"], out)
	}
	if doc["arquivo"] != filepath.Join(root, config.DefaultFile) {
		t.Errorf("arquivo = %v", doc["arquivo"])
	}
	// Nothing in the directory: the next step is the DISCOVER phase, not `map build`.
	if next, _ := doc["proximo_passo"].(string); !strings.Contains(next, "anchors guide project") {
		t.Errorf("proximo_passo = %q, want the DISCOVER phase", next)
	}

	cfg := loadWritten(t, root)
	for _, want := range []string{"spec", "feature", "test"} {
		if _, ok := cfg.Layers[want]; !ok {
			t.Errorf("layer %q missing; layers = %v", want, layerNames(cfg))
		}
	}
	if len(cfg.Gates) == 0 {
		t.Error("--gates defaults to true, yet no gate was seeded")
	}
	if cfg.Workflow == nil || cfg.Workflow.Mode != config.ModeManual {
		t.Errorf("workflow = %+v, want manual", cfg.Workflow)
	}
	if _, serr := os.Stat(filepath.Join(root, "guides", "HEADER_GUIDE.md")); serr != nil {
		t.Errorf("the header guide was not seeded: %v", serr)
	}
}

// github mode carries the repository and labels; --gates=false and --header=false are
// deliberate NOs and must be honored, not replaced by the defaults.
func TestNonInteractiveGithubModeAndExplicitNos(t *testing.T) {
	t.Run("ININT-B03: A false flag is a deliberate no", func(t *testing.T) {})
	t.Run("ININT-B04: The github workflow carries the repository and labels", func(t *testing.T) {})
	root := t.TempDir()
	err, doc, out := runNonInteractive(t, root, "--artifacts=spec",
		"--workflow=github", "--repo=acme/app", "--labels=anchors,work",
		"--gates=false", "--header=false")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("init: %v\n%s", err, out)
	}
	cfg := loadWritten(t, root)
	if cfg.Workflow == nil || cfg.Workflow.Mode != config.ModeGitHub ||
		cfg.Workflow.Repo != "acme/app" || strings.Join(cfg.Workflow.Labels, ",") != "anchors,work" {
		t.Errorf("workflow = %+v, want github acme/app [anchors work]", cfg.Workflow)
	}
	if len(cfg.Gates) != 0 {
		t.Errorf("--gates=false still seeded %d gate(s)", len(cfg.Gates))
	}
	if _, serr := os.Stat(filepath.Join(root, "guides", "HEADER_GUIDE.md")); serr == nil {
		t.Error("--header=false still seeded the header guide")
	}
}

// One invalid answer refuses the WHOLE set: writing the valid ones would produce a file
// nobody fully decided.
func TestNonInteractiveRefusesTheSetOnOneInvalidAnswer(t *testing.T) {
	t.Run("ININT-B05: One invalid answer refuses the whole set", func(t *testing.T) {})
	t.Run("ININT-I01: Either the whole set is written or nothing is", func(t *testing.T) {})
	root := t.TempDir()
	err, doc, out := runNonInteractive(t, root, "--artifacts=spec", "--workflow=github")
	if err == nil {
		t.Fatalf("github mode without --repo must be refused:\n%s", out)
	}
	if doc["escrito"] != false {
		t.Errorf("escrito = %v, want false", doc["escrito"])
	}
	var refused []string
	for _, s := range doc["respostas"].([]any) {
		m := s.(map[string]any)
		if m["aceita"] == false {
			refused = append(refused, m["id"].(string))
		}
	}
	if !containsStr(refused, "repo") {
		t.Errorf("the refused answers %v do not name `repo`", refused)
	}
	if _, serr := os.Stat(filepath.Join(root, config.DefaultFile)); serr == nil {
		t.Error("anchors.yaml was written although an answer was invalid")
	}
}

func TestNonInteractiveRejectsAMalformedGovernsRule(t *testing.T) {
	t.Run("ININT-E01: A malformed governs rule is refused with the expected form", func(t *testing.T) {})
	root := t.TempDir()
	err, _, _ := runNonInteractive(t, root, "--governs", "guides/A.md")
	if err == nil || !strings.Contains(err.Error(), "GUIDE=tag1,tag2") {
		t.Fatalf("a --governs without `=` must be refused with the expected form, got %v", err)
	}
}

// --defaults is the explicit "accept everything": only then does a call with no answers
// write the file.
func TestNonInteractiveDefaultsWritesWhenAskedTo(t *testing.T) {
	t.Run("ININT-B06: Defaults writes when asked to", func(t *testing.T) {})
	root := t.TempDir()
	err, doc, out := runNonInteractive(t, root, "--defaults")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("--defaults must write: %v\n%s", err, out)
	}
	loadWritten(t, root)
}

// --layers keeps only the chosen code layers: the others are pruned from the file.
func TestNonInteractiveLayersPrunesTheOthers(t *testing.T) {
	t.Run("ININT-B08: Layers prunes the other code layers", func(t *testing.T) {})
	t.Run("ININT-B09: The success document names the file and the next step", func(t *testing.T) {})
	root := t.TempDir()
	for i := 0; i < 10; i++ {
		writeFile(t, root, filepath.Join("src", "api", "a"+string(rune('a'+i))+".go"), "package api\n")
		writeFile(t, root, filepath.Join("src", "web", "w"+string(rune('a'+i))+".go"), "package web\n")
	}
	_, asked, _ := runNonInteractive(t, root)
	var options []string
	for _, q := range asked["perguntas"].([]any) {
		m := q.(map[string]any)
		if m["id"] == "layers" {
			for _, o := range m["opcoes"].([]any) {
				options = append(options, o.(string))
			}
		}
	}
	if len(options) != 2 {
		t.Fatalf("expected two code layers to choose from, got %v", options)
	}
	keep, drop := options[0], options[1]

	err, doc, out := runNonInteractive(t, root, "--artifacts=spec", "--layers="+keep)
	if err != nil || doc["escrito"] != true {
		t.Fatalf("init: %v\n%s", err, out)
	}
	names := layerNames(loadWritten(t, root))
	if !containsStr(names, keep) || containsStr(names, drop) {
		t.Errorf("--layers=%s: layers = %v, want %s kept and %s pruned", keep, names, keep, drop)
	}
	// A project with code does not need the DISCOVER phase: the next step is the map.
	if next, _ := doc["proximo_passo"].(string); !strings.Contains(next, "anchors map build") {
		t.Errorf("proximo_passo = %q, want `anchors map build`", next)
	}
}

func TestAsListAcceptsBothListShapes(t *testing.T) {
	if got := asList([]string{"a", "b"}); strings.Join(got, ",") != "a,b" {
		t.Errorf("[]string: %v", got)
	}
	// A decoded JSON list is []any; a non-string element is dropped, not stringified.
	if got := asList([]any{"a", 3, "c"}); strings.Join(got, ",") != "a,c" {
		t.Errorf("[]any: %v", got)
	}
	if got := asList("a"); got != nil {
		t.Errorf("a scalar is not a list: %v", got)
	}
}

func containsStr(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// A contributor arrives and does not know the spec comes first: the seeded guide says it
// from the configuration just written. A guide the project already has is the project's.
func TestNonInteractiveSeedsContributingOnlyWhenAbsent(t *testing.T) {
	t.Run("ININT-B11: CONTRIBUTING.md is seeded when absent, and an existing one is left as it is", func(t *testing.T) {})
	root := t.TempDir()
	err, doc, out := runNonInteractive(t, root, "--artifacts=spec")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("init: %v\n%s", err, out)
	}
	b, rerr := os.ReadFile(filepath.Join(root, "CONTRIBUTING.md"))
	if rerr != nil || !strings.Contains(string(b), "The spec is the anchor") {
		t.Errorf("CONTRIBUTING.md was not seeded from the configuration: %v\n%s", rerr, b)
	}
	if c, _ := doc["contributing"].(map[string]any); c["escrito"] != true {
		t.Errorf("the success JSON does not say CONTRIBUTING.md was written: %v", doc["contributing"])
	}

	ours := t.TempDir()
	writeFile(t, ours, "CONTRIBUTING.md", "ours\n")
	err, doc, out = runNonInteractive(t, ours, "--artifacts=spec")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if b, _ := os.ReadFile(filepath.Join(ours, "CONTRIBUTING.md")); string(b) != "ours\n" {
		t.Errorf("the project's CONTRIBUTING.md was touched:\n%s", b)
	}
	c, _ := doc["contributing"].(map[string]any)
	if trecho, _ := c["trecho"].(string); c["escrito"] != false || !strings.Contains(trecho, "## Working with Anchors") {
		t.Errorf("the success JSON should carry the section that would be added: %v", doc["contributing"])
	}

	none := t.TempDir()
	if err, _, out := runNonInteractive(t, none, "--artifacts=spec", "--contributing=false"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if _, serr := os.Stat(filepath.Join(none, "CONTRIBUTING.md")); serr == nil {
		t.Error("--contributing=false wrote CONTRIBUTING.md")
	}
}

// --governs was validated and echoed as accepted, and then never written: the file
// governed nothing while the JSON said it did.
func TestNonInteractiveGovernsReachesTheFile(t *testing.T) {
	t.Run("ININT-B10: The governs rules reach the configuration, one rule per tag", func(t *testing.T) {})
	root := t.TempDir()
	err, doc, out := runNonInteractive(t, root, "--artifacts=spec",
		"--governs", "guides/STYLE.md=backend,web", "--governs", "guides/API.md=backend")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("init: %v\n%s", err, out)
	}
	var got []string
	for _, g := range loadWritten(t, root).Governs {
		got = append(got, g.From+"="+g.Governs)
	}
	want := []string{"guides/API.md=backend", "guides/STYLE.md=backend", "guides/STYLE.md=web"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("governs written = %v, want %v", got, want)
	}
}

// A project copying a shell gate script inherits the platform's traps; the answer that
// seeds the gates says to write a gate step in the project's language.
func TestNonInteractiveSeedingGatesAdvisesTheProjectLanguage(t *testing.T) {
	t.Run("ININT-B12: Seeding the gates tells to write gate steps in the project's language", func(t *testing.T) {})
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module x\n")
	err, doc, out := runNonInteractive(t, root, "--artifacts=spec")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if note, _ := doc["gate_scripts"].(string); !strings.Contains(note, "go run ./tools/gates <gate>") || !strings.Contains(note, "shell") {
		t.Errorf("gate_scripts = %q", doc["gate_scripts"])
	}
	other := t.TempDir()
	if err, doc, out := runNonInteractive(t, other, "--artifacts=spec", "--gates=false"); err != nil || doc["gate_scripts"] != nil {
		t.Errorf("no gates seeded, no gate_scripts: %v %v\n%s", err, doc["gate_scripts"], out)
	}
}
