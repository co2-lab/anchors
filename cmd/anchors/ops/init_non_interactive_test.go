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
	root := t.TempDir()
	err, _, _ := runNonInteractive(t, root, "--governs", "guides/A.md")
	if err == nil || !strings.Contains(err.Error(), "GUIDE=tag1,tag2") {
		t.Fatalf("a --governs without `=` must be refused with the expected form, got %v", err)
	}
}

// --defaults is the explicit "accept everything": only then does a call with no answers
// write the file.
func TestNonInteractiveDefaultsWritesWhenAskedTo(t *testing.T) {
	root := t.TempDir()
	err, doc, out := runNonInteractive(t, root, "--defaults")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("--defaults must write: %v\n%s", err, out)
	}
	loadWritten(t, root)
}

// A stack preset fills the code layers from the stack's structure.
//
// The project has a module DIRECTORY but no code yet. The case with code on disk is
// TestNonInteractivePresetWithCodeKeepsThePresetLayers.
func TestNonInteractivePresetFillsTheCodeLayers(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "src", "modules", "users"), 0o755); err != nil {
		t.Fatal(err)
	}
	err, doc, out := runNonInteractive(t, root, "--preset=node-ts", "--artifacts=spec",
		"--governs", "guides/STYLE.md=backend,web")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("init: %v\n%s", err, out)
	}
	names := layerNames(loadWritten(t, root))
	for _, want := range []string{"modules", "core", "common"} {
		if !containsStr(names, want) {
			t.Errorf("the preset layer %q is missing: %v", want, names)
		}
	}
}

// --layers keeps only the chosen code layers: the others are pruned from the file.
func TestNonInteractiveLayersPrunesTheOthers(t *testing.T) {
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

// A stack preset in a project that ALREADY has code keeps the preset's code layers when
// `--layers` is not passed. The default of the `layers` question is inferred BEFORE the
// preset; pruning with it dropped `core` and `common` and left `modules` renamed — while
// the TUI, which asks after the preset, keeps them all.
func TestNonInteractivePresetWithCodeKeepsThePresetLayers(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < 10; i++ {
		n := string(rune('a' + i))
		writeFile(t, root, filepath.Join("src", "modules", "users", n+".service.ts"), "export const x = 1\n")
		writeFile(t, root, filepath.Join("src", "core", n+".ts"), "export const y = 1\n")
		writeFile(t, root, filepath.Join("src", "common", n+".ts"), "export const z = 1\n")
	}
	err, doc, out := runNonInteractive(t, root, "--preset=node-ts", "--artifacts=spec")
	if err != nil || doc["escrito"] != true {
		t.Fatalf("init: %v\n%s", err, out)
	}
	names := layerNames(loadWritten(t, root))
	for _, want := range []string{"modules", "core", "common"} {
		if !containsStr(names, want) {
			t.Errorf("the preset layer %q is missing: %v", want, names)
		}
	}
}
