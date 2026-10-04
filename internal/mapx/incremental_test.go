// @anchors
//   code: INTSC
//   ref: GRINC

package mapx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// readerOf serves files from a fixed set, as the scan would read them from the disk.
func readerOf(all []scan.File) (Reader, *[]string) {
	var asked []string
	return func(paths []string) ([]scan.File, error) {
		var out []scan.File
		for _, p := range paths {
			asked = append(asked, p)
			for _, f := range all {
				if f.Path == p {
					out = append(out, f)
				}
			}
		}
		return out, nil
	}, &asked
}

// partialEqualsFull builds the map without the new files, adds them, and compares with the
// full build of every file.
func partialEqualsFull(t *testing.T, all []scan.File, cfg *config.Config, fresh ...string) (*Graph, []string) {
	t.Helper()
	var without []scan.File
	for _, f := range all {
		if !slices.Contains(fresh, f.Path) {
			without = append(without, f)
		}
	}
	full := Build(all, cfg, nil)
	g := Build(without, cfg, nil)
	read, asked := readerOf(all)
	added, err := g.AddFiles(fresh, read, cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g.Nodes, full.Nodes) {
		t.Errorf("nodes differ from the full build:\npartial %+v\nfull    %+v", g.Nodes, full.Nodes)
	}
	if !reflect.DeepEqual(g.Edges, full.Edges) {
		t.Errorf("edges differ from the full build:\npartial %+v\nfull    %+v", g.Edges, full.Edges)
	}
	slices.Sort(added)
	if want := slices.Sorted(slices.Values(fresh)); !reflect.DeepEqual(added, want) {
		t.Errorf("added %v, want %v", added, want)
	}
	return g, *asked
}

func TestAddFiles_newSpecEqualsTheFullBuild(t *testing.T) {
	t.Run("GRINC-B01: A new file enters with the same nodes and edges the full build gives", func(t *testing.T) {})
	_, asked := partialEqualsFull(t, testFiles(), testCfg(), "src/Login.spec.md")
	for _, p := range asked {
		if p == "guides/SPEC_GUIDE.md" || p == "guides/FRONTEND_GUIDE.md" {
			t.Errorf("a file outside the unit was read: %s (asked %v)", p, asked)
		}
	}
}

func TestAddFiles_newFeatureReplacesTheUnitsLinks(t *testing.T) {
	t.Run("GRINC-B02: The derivation links inside the unit are replaced, removals included", func(t *testing.T) {})
	g, _ := partialEqualsFull(t, testFiles(), testCfg(), "src/Login.feature")
	if hasEdge(g, "src/Login.spec.md", "src/Login.test.tsx", EdgeTestedBy) {
		t.Error("with a feature, the spec→test link of a unit without one is gone")
	}
	partialEqualsFull(t, testFiles(), testCfg(), "src/Login.test.tsx")
	partialEqualsFull(t, testFiles(), testCfg(), "src/Login.feature", "src/Login.test.tsx")
}

func TestAddFiles_aNewAnchorGivesItsCodeToTheSiblings(t *testing.T) {
	t.Run("GRINC-B03: A new anchor gives its declared code to its derived siblings", func(t *testing.T) {})
	cfg := testCfg()
	cfg.Derived.Anchor = "spec"
	cfg.Derived.Files = map[string]config.Padroes{
		"code": {"{{dir}}/{{name}}.tsx"}, "feature": {"{{dir}}/{{name}}.feature"}, "test": {"{{dir}}/{{name}}.test.tsx"},
	}
	files := testFiles()
	files[1].HeaderCode = "LOGIX"
	files[2].Codes, files[3].Codes = nil, nil
	g, _ := partialEqualsFull(t, files, cfg, "src/Login.spec.md")
	if n := nodeByID(g, "src/Login.feature"); n == nil || n.Code != "LOGIX" {
		t.Errorf("the feature takes the new spec's code, got %+v", n)
	}
}

// The anchor's HEADER layer chooses an override, and it is only in the anchor's text: a
// stand-in anchor takes the default templates. The anchor that may own the new file — beside
// it, or carrying its name — is read.
func TestAddFiles_theAnchorsHeaderLayerChoosesTheOverride(t *testing.T) {
	t.Run("GRINC-B07: The anchor that may own a new file is read, for the override its header chooses", func(t *testing.T) {})
	cfg := &config.Config{
		Layers: map[string]config.Layer{"spec": {Kind: "spec"}, "screen": {Kind: "code"}, "test": {Kind: "test"}},
		Derived: &config.Derived{
			Anchor: "spec",
			Files:  map[string]config.Padroes{"code": {"{{dir}}/{{name}}.ts"}, "test": {"{{dir}}/{{name}}.test.ts"}},
			Overrides: []config.DerivedOverride{{When: "screen", Files: map[string]config.Padroes{
				"code": {"{{dir}}/{{name}}.tsx"}, "test": {"{{dir}}/__tests__/{{name}}.test.tsx"},
			}}},
		},
	}
	files := []scan.File{
		{Path: "app/Tela.spec.md", Layer: "spec", Kind: "spec", HeaderCode: "TELAX", HeaderLayer: "screen"},
		{Path: "app/Tela.tsx", Layer: "screen", Kind: "code"},
		{Path: "app/__tests__/Tela.test.tsx", Layer: "test", Kind: "test"},
	}
	partialEqualsFull(t, files, cfg, "app/Tela.tsx")
	partialEqualsFull(t, files, cfg, "app/__tests__/Tela.test.tsx")

	// A `{{module}}` template names the file after the directory, not the anchor: only
	// sitting beside it says the anchor may own it.
	cfg.Derived.Overrides[0].Files["test"] = config.Padroes{"{{dir}}/{{module}}.test.ts"}
	byDir := []scan.File{
		{Path: "pkg/auth/index.spec.md", Layer: "spec", Kind: "spec", HeaderCode: "AUTHX", HeaderLayer: "screen"},
		{Path: "pkg/auth/auth.test.ts", Layer: "test", Kind: "test"},
	}
	partialEqualsFull(t, byDir, cfg, "pkg/auth/auth.test.ts")
}

func TestAddFiles_aNewPlanLinksWhatItDeclares(t *testing.T) {
	t.Run("GRINC-B04: What the new file declares reaches the existing files", func(t *testing.T) {})
	files := append(testFiles(), scan.File{Path: "plans/0001.md", Layer: "plan", Kind: "plan", Rev: "p",
		Seeds: []string{"src/Login.spec.md"}})
	cfg := testCfg()
	cfg.Layers["plan"] = config.Layer{Kind: "plan"}
	partialEqualsFull(t, files, cfg, "plans/0001.md")
}

func TestAddFiles_nothingToAdd(t *testing.T) {
	t.Run("GRINC-B05: A known, ignored or unclassified file adds nothing", func(t *testing.T) {})
	g := Build(testFiles(), testCfg(), nil)
	before := len(g.Nodes)
	read, asked := readerOf(testFiles())
	if added, err := g.AddFiles([]string{"src/Login.tsx"}, read, testCfg(), nil); err != nil || added != nil || len(*asked) != 0 {
		t.Errorf("a known file is not read nor added, got %v %v %v", added, err, *asked)
	}
	if added, _ := g.AddFiles([]string{"README.md"}, read, testCfg(), nil); added != nil || len(g.Nodes) != before {
		t.Errorf("a file the reader does not give adds nothing, got %v", added)
	}
}

// What an EXISTING file declares toward the new one lives in its text, and the map does
// not keep it: the plan's seed of a spec created later waits for the next `map build`.
func TestAddFiles_whatAnExistingFileDeclaresWaits(t *testing.T) {
	t.Run("GRINC-X01: A relation an existing file declares toward the new one waits for the full build", func(t *testing.T) {})
	cfg := testCfg()
	cfg.Layers["plan"] = config.Layer{Kind: "plan"}
	all := append(testFiles(), scan.File{Path: "plans/0001.md", Layer: "plan", Kind: "plan", Seeds: []string{"src/Login.spec.md"}})
	var without []scan.File
	for _, f := range all {
		if f.Path != "src/Login.spec.md" {
			without = append(without, f)
		}
	}
	g := Build(without, cfg, nil)
	read, _ := readerOf(all)
	if _, err := g.AddFiles([]string{"src/Login.spec.md"}, read, cfg, nil); err != nil {
		t.Fatal(err)
	}
	if hasEdge(g, "plans/0001.md", "src/Login.spec.md", EdgeSeeds) {
		t.Error("the existing plan's seed was not read, so it cannot be here")
	}
	if !hasEdge(Build(all, cfg, nil), "plans/0001.md", "src/Login.spec.md", EdgeSeeds) {
		t.Error("the full build has it — that is the difference the next map build closes")
	}
}

func TestAddFilesAt_onDisk(t *testing.T) {
	t.Run("GRINC-B06: On disk the addition runs under the lock, and not without a map", func(t *testing.T) {})
	root := t.TempDir()
	cfg := &config.Config{Layers: map[string]config.Layer{"spec": {Pattern: "**/*.spec.md", Kind: "spec"}}}
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.MkdirAll(filepath.Join(root, "src"), 0o755))
	must(os.WriteFile(filepath.Join(root, "src", "pay.spec.md"), []byte("# Pay\n"), 0o644))
	mapPath := filepath.Join(root, DefaultPath)
	if added, err := AddMissingAt(root, mapPath, cfg); err != nil || added != nil {
		t.Fatalf("with no map nothing is written, got %v %v", added, err)
	}
	if _, err := os.Stat(mapPath); err == nil {
		t.Fatal("no map was created")
	}
	must(Save(&Graph{}, mapPath))
	added, err := AddMissingAt(root, mapPath, cfg)
	if err != nil || len(added) != 1 || added[0] != "src/pay.spec.md" {
		t.Fatalf("the missing spec is added, got %v %v", added, err)
	}
	if g, _ := Load(mapPath); len(g.Nodes) != 1 {
		t.Errorf("the map on disk has it, got %+v", g.Nodes)
	}
	if _, err := os.Stat(LockPath(mapPath)); err == nil {
		t.Error("the lock is released")
	}
}

func TestAddFiles_aReaderThatFails(t *testing.T) {
	t.Run("GRINC-E01: A reader that fails changes nothing", func(t *testing.T) {})
	g := Build(testFiles()[:1], testCfg(), nil)
	before := fmt.Sprint(g.Nodes, g.Edges)
	boom := errors.New("boom")
	_, err := g.AddFiles([]string{"src/Login.spec.md"}, func([]string) ([]scan.File, error) { return nil, boom }, testCfg(), nil)
	if !errors.Is(err, boom) || fmt.Sprint(g.Nodes, g.Edges) != before {
		t.Fatalf("the error comes back and the map is unchanged, got %v", err)
	}
}
