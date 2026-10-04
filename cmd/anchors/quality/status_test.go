// @anchors
//   code: STTSC
//   ref: PRSTP

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

// captureStatus runs status on a root and returns what it printed.
func captureStatus(t *testing.T, root string) string {
	t.Helper()
	var err error
	out := captureStdout(t, func() { err = runStatus(root) })
	if err != nil {
		t.Fatalf("runStatus: %v", err)
	}
	return out
}

func emptyRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".git", "refs", "heads"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// The reason the command exists: whoever picks the work back up days later needs to know
// where it stopped. In a project not started yet the answer is the DISCOVER phase — and it
// has to be NAMED, or the agent opening the conversation starts by guessing.
func TestStatusNamesTheDiscoverPhaseInANewProject(t *testing.T) {
	t.Run("PRSTP-B03: A project with neither PROJECT.md nor configuration is not started", func(t *testing.T) {})
	englishOutput(t)
	out := captureStatus(t, emptyRepo(t))

	if !strings.Contains(out, "the DISCOVER phase") {
		t.Errorf("no PROJECT.md and no config: the missing phase has to be named:\n%s", out)
	}
	if !strings.Contains(out, "NEXT STEP") {
		t.Errorf("status exists to say what to do next:\n%s", out)
	}
}

// Status stops at the FIRST missing step. Listing everything at once would make the reader
// choose where to start — and the cycle's order is exactly what they should not have to
// rebuild on their own.
func TestStatusStopsAtTheFirstMissingStep(t *testing.T) {
	t.Run("PRSTP-B04: A project with PROJECT.md and no configuration is sent to init", func(t *testing.T) {})
	t.Run("PRSTP-I01: Status never names a step beyond the first one missing", func(t *testing.T) {})
	englishOutput(t)
	dir := emptyRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "PROJECT.md"), []byte("# Project\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStatus(t, dir)

	if !strings.Contains(out, "✓ PROJECT.md exists") {
		t.Errorf("what was already done has to show as done:\n%s", out)
	}
	if !strings.Contains(out, "anchors init") {
		t.Errorf("with PROJECT.md and no config, the step is init:\n%s", out)
	}
	// And it must not talk about what comes after: the map is not the question yet.
	if strings.Contains(out, "map build") {
		t.Errorf("it jumped to a step that is not the next one:\n%s", out)
	}
}

// Without a repository, status stops before anything: half the framework depends on git,
// and going on to describe the cycle would suggest everything is fine.
func TestStatusStopsWithoutRepository(t *testing.T) {
	t.Run("PRSTP-B01: A directory with no git repository stops at git init", func(t *testing.T) {})
	englishOutput(t)
	dir := t.TempDir()
	if temRepoAcima(dir) {
		t.Skipf("the temporary directory %s is inside a git repository", dir)
	}

	out := captureStatus(t, dir)

	if !strings.Contains(out, "git init") {
		t.Errorf("without a repository, the step is to create one:\n%s", out)
	}
	if strings.Contains(out, "DISCOVER") {
		t.Errorf("it must not jump to the next phase before there is a substrate:\n%s", out)
	}
}

func TestStatusWarnsWithoutGitBinaryAndGoesOn(t *testing.T) {
	t.Run("PRSTP-B02: Without the git binary status warns and goes on", func(t *testing.T) {})
	englishOutput(t)
	t.Setenv("PATH", t.TempDir())
	out := captureStatus(t, emptyRepo(t))
	if !strings.Contains(out, "⚠ git not installed") {
		t.Errorf("without git the missing binary is named:\n%s", out)
	}
	if !strings.Contains(out, "the DISCOVER phase") {
		t.Errorf("a missing git binary does not stop the cycle:\n%s", out)
	}
}

// The queue lives where the mode declares. In local mode it is `.anchors/tasks/` and
// `issues/` — showing a board there would describe a place where nothing happens.
func TestStatusLocalShowsTheLocalQueue(t *testing.T) {
	t.Run("PRSTP-B08: An assembled local project with no work is sent to the first plan", func(t *testing.T) {})
	englishOutput(t)
	dir := emptyRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "PROJECT.md"), []byte("# P\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte("version: 2\nlayers: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "anchors.graph.yaml"), []byte("version: 7\nnodes: []\nedges: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStatus(t, dir)

	if !strings.Contains(out, "queue: local") {
		t.Errorf("local mode must show the local queue:\n%s", out)
	}
	// Assembled and empty project: the step is the first plan — it is what seeds the specs.
	if !strings.Contains(out, "project is assembled and has no work yet") {
		t.Errorf("assembled and empty project: the next step is the first plan:\n%s", out)
	}
}

// temRepoAcima avoids a false positive when the TempDir falls inside a repository.
func temRepoAcima(root string) bool {
	dir := root
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}
}

// statusProject is an assembled project (PROJECT.md, anchors.yaml, map) with the given
// extra files, inside a directory that looks like a git repository.
func statusProject(t *testing.T, yaml string, g *mapx.Graph, files map[string]string) string {
	t.Helper()
	dir := emptyRepo(t)
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
	t.Run("PRSTP-B07: The local queue names the first pending step in the cycle's order", func(t *testing.T) {})
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
	t.Run("PRSTP-B06: The configuration and the map are counted, and clean informative gates are named", func(t *testing.T) {})
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
	t.Run("PRSTP-B05: A configured project with no map is sent to the map build", func(t *testing.T) {})
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
	t.Run("PRSTP-E01: A configuration that does not load fails the status", func(t *testing.T) {})
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
	t.Run("PRSTP-B09: The github queue with missing pipelines stops at the doctor fix", func(t *testing.T) {})
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
	t.Run("PRSTP-B10: The github queue states the pull-request flow and the protected branches", func(t *testing.T) {})
	t.Run("PRSTP-B11: The agent's own open cards come before claiming new work", func(t *testing.T) {})
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
	t.Run("PRSTP-B12: Without an agent identity the next step is to claim work", func(t *testing.T) {})
	t.Run("PRSTP-B13: A github project with no real work is sent to the first plan", func(t *testing.T) {})
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

func TestStatusChangesNothing(t *testing.T) {
	t.Run("PRSTP-X01: Status leaves the project as it found it", func(t *testing.T) {})
	englishOutput(t)
	dir := statusProject(t, statusYAML, workGraph(), map[string]string{"issues/todo/b.md": "x", ".anchors/tasks/pending__1.yaml": "x"})
	snapshot := func() map[string]string {
		m := map[string]string{}
		filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				b, _ := os.ReadFile(p)
				m[p] = string(b)
			}
			return nil
		})
		return m
	}
	before := snapshot()
	if _, err := runQ(t, newStatusCmd(), "--root", dir); err != nil {
		t.Fatal(err)
	}
	after := snapshot()
	if len(after) != len(before) {
		t.Errorf("status created or removed files: %d before, %d after", len(before), len(after))
	}
	for p, c := range before {
		if after[p] != c {
			t.Errorf("status changed %s", p)
		}
	}
}
