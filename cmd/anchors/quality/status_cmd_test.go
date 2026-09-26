package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/mapx"
)

// statusProject is an assembled project (PROJECT.md, anchors.yaml, map) with the given
// extra files, inside a directory that looks like a git repository.
func statusProject(t *testing.T, yaml string, g *mapx.Graph, files map[string]string) string {
	t.Helper()
	dir := repoVazio(t)
	all := map[string]string{"PROJECT.md": "# P\n", "anchors.yaml": yaml}
	for k, v := range files {
		all[k] = v
	}
	for name, body := range all {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if g != nil {
		if err := mapx.Save(g, filepath.Join(dir, mapx.DefaultPath)); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func workGraph() *mapx.Graph {
	return &mapx.Graph{Nodes: []mapx.Node{{ID: "a.go", Kind: mapx.KindCode}}}
}

const statusYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
gates:
  - name: always-green
    on: [code]
    run: "true"
`

// The local queue answers in the cycle's order: what is in progress first, then what
// is waiting, then the task queue, then "nothing pending".
func TestStatusLocalPicksTheFirstPendingStep(t *testing.T) {
	englishOutput(t)
	cases := []struct {
		name  string
		files map[string]string
		want  string
	}{
		{"doing wins", map[string]string{"issues/doing/a.md": "x", "issues/todo/b.md": "x", ".anchors/tasks/pending__1.yaml": "x"}, "work in `issues/doing`"},
		{"then todo", map[string]string{"issues/todo/b.md": "x", ".anchors/tasks/pending__1.yaml": "x"}, "`issues/todo` has work"},
		{"then tasks", map[string]string{".anchors/tasks/pending__1.yaml": "x"}, "`anchors next` pulls the next task"},
		{"nothing", nil, "nothing pending. `anchors doctor`"},
	}
	for _, c := range cases {
		dir := statusProject(t, statusYAML, workGraph(), c.files)
		out, err := runQ(t, newStatusCmd(), "--root", dir)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, c.want) {
			t.Errorf("%s: want %q in:\n%s", c.name, c.want, out)
		}
		for _, other := range cases {
			if other.want != c.want && strings.Contains(out, other.want) {
				t.Errorf("%s: only the first pending step is named, also got %q:\n%s", c.name, other.want, out)
			}
		}
	}
}

func TestStatusReportsConfigMapAndCleanInformativeGates(t *testing.T) {
	englishOutput(t)
	dir := statusProject(t, statusYAML, workGraph(), map[string]string{"issues/todo/b.md": "x", "issues/todo/c.md": "x"})

	out, err := runQ(t, newStatusCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"✓ anchors.yaml (1 layers, 1 gates)",
		"✓ map (1 nodes, 0 edges)",
		// an informative gate that passes everywhere is a candidate for blocking
		"○ 1 informative gate(s) CLEAN: always-green",
		"  issues: 2 in todo, 0 in doing",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestStatusStopsWithoutMap(t *testing.T) {
	englishOutput(t)
	dir := statusProject(t, statusYAML, nil, nil)
	out, err := runQ(t, newStatusCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "○ no map") {
		t.Errorf("without the map the next step is map build:\n%s", out)
	}
	if strings.Contains(out, "queue:") {
		t.Errorf("the queue is not the question before the map exists:\n%s", out)
	}
}

func TestStatusFailsOnABrokenConfig(t *testing.T) {
	englishOutput(t)
	dir := statusProject(t, "version: 2\nnot_a_key: 1\n", nil, nil)
	if _, err := runQ(t, newStatusCmd(), "--root", dir); err == nil || !strings.Contains(err.Error(), "load anchors.yaml") {
		t.Errorf("a config that does not load must fail; got %v", err)
	}
}

const statusGitHubYAML = `version: 2
layers:
  code:
    kind: code
    pattern: "*.go"
workflow:
  mode: github
  repo: acme/app
  labels: [anchors]
  integration_branch: develop
  protected_branches: [develop, main]
`

func TestStatusGitHubPointsAtMissingPipelines(t *testing.T) {
	englishOutput(t)
	dir := statusProject(t, statusGitHubYAML, workGraph(), nil)

	out, err := runQ(t, newStatusCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "queue: GitHub (acme/app, label [anchors])") {
		t.Errorf("github mode names the repository and label:\n%s", out)
	}
	if !strings.Contains(out, "workflow pipeline(s) missing") || strings.Contains(out, "work enters via PR") {
		t.Errorf("without pipelines, status stops at doctor --fix:\n%s", out)
	}
}

func seedPipelines(t *testing.T, dir string) {
	t.Helper()
	if _, _, err := initx.SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{Mode: config.ModeGitHub}}); err != nil {
		t.Fatal(err)
	}
}

// With this agent's card open, status says so before telling it to claim another.
func TestStatusGitHubShowsTheAgentsOwnCard(t *testing.T) {
	englishOutput(t)
	fakeGH(t, "write")
	t.Setenv("ANCHORS_AGENT", "host/session")
	dir := statusProject(t, statusGitHubYAML, workGraph(), nil)
	seedPipelines(t, dir)

	out, err := runQ(t, newStatusCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"✓ workflow pipelines in place",
		"  work enters via PR to `develop` · protected: develop, main",
		"  → YOU ALREADY HAVE WORK:\n    #12 Fix login [in-progress]\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "ask the claim pipeline for work") {
		t.Errorf("with a card of its own, the agent must not be told to claim another:\n%s", out)
	}
}

func TestStatusGitHubWithoutIdentityClaimsNext(t *testing.T) {
	englishOutput(t)
	fakeGH(t, "write")
	t.Setenv("ANCHORS_AGENT", "")
	dir := statusProject(t, statusGitHubYAML, workGraph(), nil)
	seedPipelines(t, dir)

	out, err := runQ(t, newStatusCmd(), "--root", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "ask the claim pipeline for work") || strings.Contains(out, "YOU ALREADY HAVE WORK") {
		t.Errorf("without ANCHORS_AGENT nobody's cards can be listed; the step is claim:\n%s", out)
	}

	// an assembled project with only the seeded guides: the step is the first plan
	guides := statusProject(t, statusGitHubYAML, &mapx.Graph{Nodes: []mapx.Node{{ID: "G.md", Kind: mapx.KindGuide}}}, nil)
	seedPipelines(t, guides)
	out, err = runQ(t, newStatusCmd(), "--root", guides)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "project is assembled and has no work yet") {
		t.Errorf("no real work: the first plan comes before any card:\n%s", out)
	}
}

func TestCountFilesIgnoresDirectoriesAndMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if n := countFiles(dir); n != 1 {
		t.Errorf("countFiles = %d, want 1 (the subdirectory does not count)", n)
	}
	if n := countFiles(filepath.Join(dir, "absent")); n != 0 {
		t.Errorf("a missing directory is 0, got %d", n)
	}
}
