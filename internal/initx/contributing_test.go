// @anchors
//   code: CNTSD
//   ref: CNGDC

package initx

import (
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestRenderContributingWhole(t *testing.T) {
	t.Run("CNGDC-B01: The whole guide is a title, a seeding note and the section", func(t *testing.T) {})
	cfg := &config.Config{}
	out := RenderContributing(cfg, "")
	if !strings.HasPrefix(out, "# Contributing\n") || !strings.Contains(out, "seeded by `anchors init`") {
		t.Errorf("the guide should open with the title and the seeding note:\n%s", out)
	}
	if !strings.HasSuffix(out, ContributingSection(cfg, "")) {
		t.Error("the guide should end with the section")
	}
}

func TestContributingOrderOfTheWork(t *testing.T) {
	t.Run("CNGDC-B02: The section states the order of the work", func(t *testing.T) {})
	out := ContributingSection(&config.Config{}, "")
	if !strings.HasPrefix(out, "## Working with Anchors\n") {
		t.Errorf("the section should open with its heading:\n%s", out)
	}
	for _, want := range []string{"The spec is the anchor", "write the rule, then the scenario, then the test"} {
		if !strings.Contains(out, want) {
			t.Errorf("the section should say %q", want)
		}
	}
}

func TestContributingWhereEachPieceLives(t *testing.T) {
	t.Run("CNGDC-B03: Where each piece lives", func(t *testing.T) {})
	cfg := &config.Config{
		Layers: map[string]config.Layer{
			"test":    {Kind: "test", Pattern: "**/*_test.go"},
			"spec":    {Kind: "spec", Pattern: "**/*.spec.md"},
			"feature": {Kind: "feature", Pattern: "**/*.feature"},
		},
		Derived: &config.Derived{Anchor: "spec", Files: map[string]config.Padroes{
			"test":    {"{{dir}}/{{name}}_test.go"},
			"feature": {"{{dir}}/{{name}}.feature"},
		}},
	}
	out := ContributingSection(cfg, "")
	s, f, te := strings.Index(out, "- spec: `**/*.spec.md`"), strings.Index(out, "- feature: `**/*.feature`"), strings.Index(out, "- test: `**/*_test.go`")
	if s < 0 || f < s || te < f {
		t.Errorf("the pieces should be listed spec, feature, test:\n%s", out)
	}
	if !strings.Contains(out, "beside its spec: the feature at `{{dir}}/{{name}}.feature`; the test at `{{dir}}/{{name}}_test.go`") {
		t.Errorf("the colocation templates are missing:\n%s", out)
	}
	if strings.Contains(ContributingSection(&config.Config{}, ""), "Where each piece lives") {
		t.Error("with no artifact layer no piece is listed")
	}
}

func TestContributingCodeLayers(t *testing.T) {
	t.Run("CNGDC-B04: The declared code layers, or that there are none", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{
		"handlers-code": {Kind: "code", Pattern: "src/handlers/**/*.go"},
		"api-code":      {Kind: "code", Pattern: "api/**/*.go"},
	}}
	out := ContributingSection(cfg, "")
	a, h := strings.Index(out, "- `api-code` — `api/**/*.go`"), strings.Index(out, "- `handlers-code` — `src/handlers/**/*.go`")
	if a < 0 || h < a {
		t.Errorf("the code layers should be listed by name:\n%s", out)
	}
	if !strings.Contains(out, "illustration, not a layout") {
		t.Error("the layer kinds should be given as illustration")
	}
	if !strings.Contains(ContributingSection(&config.Config{}, ""), "No code layer is declared yet") {
		t.Error("with no code layer the section should say so")
	}
}

func TestContributingQueue(t *testing.T) {
	t.Run("CNGDC-B05: The queue the daily commands name", func(t *testing.T) {})
	gh := &config.Config{Workflow: &config.Workflow{Mode: config.ModeGitHub, Repo: "acme/app"}}
	if !strings.Contains(ContributingSection(gh, ""), "the issues of `acme/app`") {
		t.Error("github mode should name the repository's issues")
	}
	if !strings.Contains(ContributingSection(&config.Config{}, ""), "the local queue") {
		t.Error("with no workflow the local queue should be named")
	}
}

func TestContributingBlocksAndInforms(t *testing.T) {
	t.Run("CNGDC-B06: What blocks and what informs", func(t *testing.T) {})
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "tests-green", Blocking: config.Bool(true)},
		{Name: "spec-complete"},
	}}
	out := ContributingSection(cfg, "")
	if !strings.Contains(out, "bars the commit: `tests-green`.") || !strings.Contains(out, "These only inform: `spec-complete`.") {
		t.Errorf("blocking and informing gates should be split:\n%s", out)
	}
	if !strings.Contains(ContributingSection(&config.Config{Gates: []config.Gate{{Name: "x"}}}, ""), "every gate only informs") {
		t.Error("with no blocking gate it should say every gate only informs")
	}
}

func TestContributingGuideFolder(t *testing.T) {
	t.Run("CNGDC-B07: The project's guide folder", func(t *testing.T) {})
	if !strings.Contains(ContributingSection(&config.Config{}, "docs/guides"), "`docs/guides/`") {
		t.Error("the guide folder should be named")
	}
	if strings.Contains(ContributingSection(&config.Config{}, ""), "own guides are in") {
		t.Error("with no guide folder no folder should be named")
	}
}

func TestContributingNamesOnlyWhatIsDeclared(t *testing.T) {
	t.Run("CNGDC-I01: The guide never names what the configuration does not hold", func(t *testing.T) {})
	cfg := &config.Config{
		Layers: map[string]config.Layer{"core-code": {Kind: "code", Pattern: "core/**/*.go"}},
		Gates:  []config.Gate{{Name: "tests-green", Blocking: config.Bool(true)}},
	}
	out := ContributingSection(cfg, "")
	if n := strings.Count(out, "-code`"); n != 1 {
		t.Errorf("one code layer declared, %d listed:\n%s", n, out)
	}
	if strings.Contains(out, "These only inform") {
		t.Errorf("no informing gate is declared:\n%s", out)
	}
}

func TestContributingWritesNothing(t *testing.T) {
	t.Run("CNGDC-X01: Rendering writes nothing to disk", func(t *testing.T) {})
	dir := t.TempDir()
	t.Chdir(dir)
	if RenderContributing(&config.Config{}, "guides") == "" {
		t.Fatal("the guide should come back as text")
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("rendering wrote to disk: %v", entries)
	}
}

func TestContributingRecommendsTheMonitor(t *testing.T) {
	t.Run("CNGDC-B08: The guide recommends watching long commands and agents", func(t *testing.T) {})
	s := ContributingSection(nil, "")
	for _, want := range []string{"### Long commands and agents: watch them", "**Strongly recommended:**", "`anchors monitor run -- <command>`, in the background", "Keep `anchors monitor` running", "`anchors install-hooks --agent`", "show that summary to the user as it comes and in full", "must let those lines through"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in:\n%s", want, s)
		}
	}
}

func TestAgentEntryFilesPointToTheGuide(t *testing.T) {
	t.Run("CNGDC-B09: The agents' entry files send the agent to the contributing guide", func(t *testing.T) {})
	files := AgentEntryFiles()
	if len(files) != 2 || files[0].Name != "AGENTS.md" || files[1].Name != "CLAUDE.md" {
		t.Fatalf("files: %+v", files)
	}
	for _, f := range files {
		if !strings.Contains(f.Content, "Read `CONTRIBUTING.md` before any change") {
			t.Errorf("%s does not send the agent to the guide:\n%s", f.Name, f.Content)
		}
	}
	for _, want := range []string{"anchors guide", "anchors work <stage> --for <target>", "anchors next", "anchors monitor run -- <command>", "show that summary to the user as it comes", "--no-verify"} {
		if !strings.Contains(files[0].Content, want) {
			t.Errorf("AGENTS.md lacks %q", want)
		}
	}
	if !strings.Contains(files[1].Content, "\n@AGENTS.md\n") {
		t.Errorf("CLAUDE.md imports AGENTS.md:\n%s", files[1].Content)
	}
}

func TestContributingNamesWorkAndForbidsNoVerify(t *testing.T) {
	t.Run("CNGDC-B10: The guide names the task prompt and forbids skipping the hooks", func(t *testing.T) {})
	s := ContributingSection(nil, "")
	for _, want := range []string{"`anchors work <stage> --for <target>`", "Never commit with `--no-verify`", "`[skip-<rule>@<CODE>: why]`"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in:\n%s", want, s)
		}
	}
}
