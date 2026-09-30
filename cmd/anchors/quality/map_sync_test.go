package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

type syncRepo struct {
	root string
	cfg  *config.Config
	git  func(args ...string) string
}

func newSyncRepo(t *testing.T, trackMap bool) syncRepo {
	t.Helper()
	root := t.TempDir()
	past := false // the base commits are dated in the past, as a real project's are
	git := func(args ...string) string {
		t.Helper()
		c := exec.Command("git", append([]string{"-C", root, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...)
		if past {
			c.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2026-09-01T12:00:00", "GIT_COMMITTER_DATE=2026-09-01T12:00:00")
		}
		out, err := c.CombinedOutput()
		if err != nil {
			t.Skipf("git %v: %v %s", args, err, out)
		}
		return string(out)
	}
	git("init", "-q")
	yaml := "version: 5\nlayers:\n  spec:\n    pattern: \"src/*.spec.md\"\n    kind: spec\n  code:\n    pattern: \"src/*.ts\"\n    kind: code\nderived:\n  anchor: code\n  files:\n    spec: [\"{{dir}}/{{name}}.spec.md\"]\n"
	touchWrite(t, root, "anchors.yaml", yaml)
	touchWrite(t, root, "src/pay.spec.md", "<!-- @anchors\n  code: PAYMX\n  updated_at: 2026-09-01\n-->\n# Pay\n\n### PAYMX-B01 — charges\n")
	touchWrite(t, root, "src/pay.ts", "export const pay = 1\n")
	touchWrite(t, root, "src/other.ts", "export const other = 1\n")
	git("add", ".")
	past = true
	git("commit", "-qm", "base")
	cfg, err := config.Load(filepath.Join(root, "anchors.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files, _ := scan.Walk(root, cfg)
	g := mapx.Build(files, cfg, gitmeta.AllCommitDates(root))
	for i := range g.Nodes {
		if g.Nodes[i].ID == "src/pay.spec.md" {
			g.Nodes[i].Signal = &mapx.TestSignal{ProvenCodes: []string{"PAYMX-B01"}, AtRev: g.Nodes[i].Rev}
		}
	}
	if err := mapx.Save(g, filepath.Join(root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	if trackMap {
		git("add", mapx.DefaultPath)
		git("commit", "-qm", "map")
	}
	past = false
	return syncRepo{root, cfg, git}
}

// the edit the tests commit: the spec staged, other.ts edited and not staged, new.ts untracked
func (r syncRepo) change(t *testing.T) []touchDecision {
	t.Helper()
	touchWrite(t, r.root, "src/pay.spec.md", "<!-- @anchors\n  code: PAYMX\n  updated_at: 2026-09-01\n-->\n# Pay\n\n### PAYMX-B01 — charges the amount\n")
	r.git("add", "src/pay.spec.md")
	touchWrite(t, r.root, "src/other.ts", "export const other = 2 // not staged\n")
	touchWrite(t, r.root, "src/new.ts", "export const fresh = 1\n")
	bumped, _, err := touchRun(r.root, true, false, gitmeta.Today(), nil, nil)
	if err != nil || len(bumped) != 1 {
		t.Fatalf("the hook dates the staged spec, got %+v %v", bumped, err)
	}
	return bumped
}

// stagedMap is the map the commit will record: the index's.
func (r syncRepo) stagedMap(t *testing.T) *mapx.Graph {
	t.Helper()
	g, err := mapx.LoadBytes([]byte(r.git("show", ":"+mapx.DefaultPath)), "index")
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func nodeOf(g *mapx.Graph, id string) *mapx.Node {
	for i := range g.Nodes {
		if g.Nodes[i].ID == id {
			return &g.Nodes[i]
		}
	}
	return nil
}

func TestSyncMap_theCommittedMapIsTheBuildsMap(t *testing.T) {
	t.Run("MPSYN-B01: The committed map is the one a build of the commit makes", func(t *testing.T) {})
	r := newSyncRepo(t, true)
	bumped := r.change(t)
	if msg, err := syncMapForCommit(r.root, r.cfg, bumped); err != nil || !strings.Contains(msg, "staged") {
		t.Fatalf("the map is synced and staged, got %q %v", msg, err)
	}
	r.git("commit", "-qm", "change")
	committed := r.git("show", "HEAD:"+mapx.DefaultPath)
	// what the CI does: build the commit's files, keeping what the committed map holds
	files, err := scan.WalkStaged(r.root, r.cfg)
	if err != nil {
		t.Fatal(err)
	}
	built := mapx.Build(files, r.cfg, gitmeta.AllCommitDates(r.root))
	tmp := filepath.Join(t.TempDir(), "committed.yaml")
	if err := os.WriteFile(tmp, []byte(committed), 0o644); err != nil {
		t.Fatal(err)
	}
	was, err := mapx.Load(tmp)
	if err != nil {
		t.Fatal(err)
	}
	mapx.PreserveStamps(built, was)
	built.Flow = was.Flow
	if !reflect.DeepEqual(built.Nodes, was.Nodes) || !reflect.DeepEqual(built.Edges, was.Edges) {
		t.Fatalf("a build of the commit differs from the committed map:\nbuilt %+v\ncommitted %+v", built.Nodes, was.Nodes)
	}
}

func TestSyncMap_aDatedFileKeepsItsProofs(t *testing.T) {
	t.Run("MPSYN-B02: A dated file keeps its proofs", func(t *testing.T) {})
	r := newSyncRepo(t, true)
	bumped := r.change(t)
	// the proof was measured at the staged content BEFORE the date: re-measure it there
	g, _ := mapx.Load(filepath.Join(r.root, mapx.DefaultPath))
	if n := nodeOf(g, "src/pay.spec.md"); n != nil {
		n.Rev = scan.ShortHash([]byte(bumped[0].Old))
		n.Signal.AtRev = n.Rev
	}
	if err := mapx.Save(g, filepath.Join(r.root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	if _, err := syncMapForCommit(r.root, r.cfg, bumped); err != nil {
		t.Fatal(err)
	}
	g = r.stagedMap(t)
	n := nodeOf(g, "src/pay.spec.md")
	if n == nil || n.Signal == nil || n.Signal.AtRev != n.Rev || n.Rev != scan.ShortHash([]byte(bumped[0].Content)) {
		t.Fatalf("the proof moves to the dated revision, got %+v", n)
	}
}

func TestSyncMap_theIndexNotTheTree(t *testing.T) {
	t.Run("MPSYN-B03: The index, not the tree", func(t *testing.T) {})
	r := newSyncRepo(t, true)
	bumped := r.change(t)
	if _, err := syncMapForCommit(r.root, r.cfg, bumped); err != nil {
		t.Fatal(err)
	}
	g := r.stagedMap(t)
	if n := nodeOf(g, "src/other.ts"); n == nil || n.Rev != scan.ShortHash([]byte("export const other = 1\n")) {
		t.Errorf("the unstaged edit is not in the map, got %+v", n)
	}
	if nodeOf(g, "src/new.ts") != nil {
		t.Error("the untracked file is not in the map")
	}
}

func TestSyncMap_anUntrackedMapIsLeftAlone(t *testing.T) {
	t.Run("MPSYN-B04: An untracked map is left alone", func(t *testing.T) {})
	r := newSyncRepo(t, false)
	before, _ := os.ReadFile(filepath.Join(r.root, mapx.DefaultPath))
	if msg, err := syncMapForCommit(r.root, r.cfg, r.change(t)); err != nil || msg != "" {
		t.Fatalf("nothing to sync, got %q %v", msg, err)
	}
	if after, _ := os.ReadFile(filepath.Join(r.root, mapx.DefaultPath)); string(after) != string(before) {
		t.Error("an untracked map is not written")
	}
}

func TestSyncMap_aMapThatCannotBeWritten(t *testing.T) {
	t.Run("MPSYN-E01: A map that cannot be written gives the error back", func(t *testing.T) {})
	r := newSyncRepo(t, true)
	bumped := r.change(t)
	mapPath := filepath.Join(r.root, mapx.DefaultPath)
	if err := os.Remove(mapPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(mapPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := syncMapForCommit(r.root, r.cfg, bumped); err == nil {
		t.Fatal("a map that cannot be written gives the error back")
	}
}

func TestSyncMap_untouchedFilesKeepHeadsProofs(t *testing.T) {
	t.Run("MPSYN-B05: A file the commit does not change keeps the proofs HEAD had", func(t *testing.T) {})
	r := newSyncRepo(t, true)
	mapPath := filepath.Join(r.root, mapx.DefaultPath)
	// another session edits the proven spec and does not stage it; a map build drops its proof
	touchWrite(t, r.root, "src/pay.spec.md", "<!-- @anchors\n  code: PAYMX\n  updated_at: 2026-09-01\n-->\n# Pay\n\n### PAYMX-B01 — charges, edited elsewhere\n")
	old, _ := mapx.Load(mapPath)
	files, _ := scan.Walk(r.root, r.cfg)
	g := mapx.Build(files, r.cfg, gitmeta.AllCommitDates(r.root))
	mapx.PreserveStamps(g, old)
	if n := nodeOf(g, "src/pay.spec.md"); n.Signal != nil {
		t.Fatal("the rebuild must have dropped the edited spec's proof")
	}
	if err := mapx.Save(g, mapPath); err != nil {
		t.Fatal(err)
	}
	// this commit only carries other.ts
	touchWrite(t, r.root, "src/other.ts", "export const other = 3\n")
	r.git("add", "src/other.ts")
	if msg, err := syncMapForCommit(r.root, r.cfg, nil); err != nil || !strings.Contains(msg, "staged apart") {
		t.Fatalf("the tree is ahead: the commit's map is staged apart, got %q %v", msg, err)
	}
	if n := nodeOf(r.stagedMap(t), "src/pay.spec.md"); n == nil || n.Signal == nil || n.Signal.ProvenCodes[0] != "PAYMX-B01" {
		t.Errorf("the committed map keeps the spec's proof from HEAD, got %+v", n)
	}
}

func TestSyncMap_theTreeKeepsItsMap(t *testing.T) {
	t.Run("MPSYN-B06: With the tree ahead of the commit, the map on disk stays the tree's", func(t *testing.T) {})
	r := newSyncRepo(t, true)
	mapPath := filepath.Join(r.root, mapx.DefaultPath)
	g, _ := mapx.Load(mapPath)
	nodeOf(g, "src/other.ts").Signal = &mapx.TestSignal{TotalLines: 1, CoveredLines: 1, AtRev: "tree-rev"}
	if err := mapx.Save(g, mapPath); err != nil {
		t.Fatal(err)
	}
	bumped := r.change(t)
	if _, err := syncMapForCommit(r.root, r.cfg, bumped); err != nil {
		t.Fatal(err)
	}
	disk, _ := mapx.Load(mapPath)
	if n := nodeOf(disk, "src/other.ts"); n.Signal == nil || n.Signal.AtRev != "tree-rev" {
		t.Errorf("the map on disk keeps what was measured of the tree, got %+v", n.Signal)
	}
	if n := nodeOf(r.stagedMap(t), "src/pay.spec.md"); n == nil || n.Rev != scan.ShortHash([]byte(bumped[0].Content)) {
		t.Errorf("the staged map is the commit's, got %+v", n)
	}

	clean := newSyncRepo(t, true)
	touchWrite(t, clean.root, "src/pay.spec.md", "<!-- @anchors\n  code: PAYMX\n  updated_at: 2026-09-01\n-->\n# Pay\n\n### PAYMX-B01 — charges clean\n")
	clean.git("add", "src/pay.spec.md")
	if msg, err := syncMapForCommit(clean.root, clean.cfg, nil); err != nil || strings.Contains(msg, "apart") {
		t.Fatalf("a clean tree writes the map as usual, got %q %v", msg, err)
	}
	onDisk, _ := os.ReadFile(filepath.Join(clean.root, mapx.DefaultPath))
	if staged := clean.git("show", ":"+mapx.DefaultPath); staged != string(onDisk) {
		t.Error("with a clean tree the map on disk is the one staged")
	}
}

func TestSyncMap_stagedApartBelowTheTop(t *testing.T) {
	t.Run("MPSYN-B06: With the tree ahead of the commit, the map on disk stays the tree's", func(t *testing.T) {})
	top := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", top, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...).CombinedOutput()
		if err != nil {
			t.Skipf("git %v: %v %s", args, err, out)
		}
		return string(out)
	}
	git("init", "-q")
	git("commit", "-q", "--allow-empty", "-m", "base")
	g := &mapx.Graph{}
	if err := os.MkdirAll(filepath.Join(top, "app"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := stageMapBlob(filepath.Join(top, "app"), g); err != nil {
		t.Fatal(err)
	}
	if staged := git("ls-files", "--cached"); strings.TrimSpace(staged) != "app/"+mapx.DefaultPath {
		t.Errorf("the map is staged in the project's own directory, got %q", staged)
	}
}
