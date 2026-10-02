// @anchors
//   ref: PRJTS

package gate

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
	"github.com/co2-lab/anchors/internal/testlist"
)

func TestProjectTests_SourceResolution(t *testing.T) {
	t.Run("PRJTS-B01: The project's own source wins over its family's, and neither declares nothing", func(t *testing.T) {})
	ts := testsSource(&config.Config{Dialect: &config.Dialect{Family: "ts"}})
	if ts.Pattern == "" || ts.Script != "" {
		t.Errorf("a ts project reads through its family's pattern, got %+v", ts)
	}
	own := testsSource(&config.Config{Dialect: &config.Dialect{Family: "ts", Tests: &config.TestsSource{Script: "list"}}})
	if own != (testlist.Source{Script: "list"}) {
		t.Errorf("the project's own script wins, got %+v", own)
	}
	for _, cfg := range []*config.Config{nil, {Dialect: &config.Dialect{Family: "python"}}} {
		if testsSource(cfg).Declared() {
			t.Errorf("%+v declares no source", cfg)
		}
	}
	tests, declared, err := projectTests(t.TempDir(), &mapx.Graph{}, nil)
	if tests != nil || declared || err != nil {
		t.Errorf("with no source nothing is read, got %v %v %v", tests, declared, err)
	}
}

func TestProjectTests_OnlyTheMapsTestFiles(t *testing.T) {
	t.Run("PRJTS-B02: A pattern reads only the map's test files", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root := t.TempDir()
	for _, f := range []string{"a_test.go", "b_test.go", "c.go"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("t.Run(\"in "+f+"\", nil)\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a_test.go", Kind: mapx.KindTest}, {ID: "c.go", Kind: mapx.KindCode}}}
	tests, declared, err := projectTests(root, g, &config.Config{Dialect: &config.Dialect{Family: "go"}})
	if err != nil || !declared || len(tests) != 1 || tests[0].Title != "in a_test.go" {
		t.Fatalf("only the map's test file is read, got %v (%v, %v)", tests, declared, err)
	}
}

func TestProjectTests_OncePerScan(t *testing.T) {
	t.Run("PRJTS-B03: The tests are read once per scan", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root := t.TempDir()
	cfg := &config.Config{Dialect: &config.Dialect{Tests: &config.TestsSource{
		Script: `echo run >> runs.txt; echo '{"version":1,"tests":[]}'`}}}
	g := &mapx.Graph{}
	projectTests(root, g, cfg)
	projectTests(root, g, cfg)
	runs := func() int {
		b, _ := os.ReadFile(filepath.Join(root, "runs.txt"))
		return strings.Count(string(b), "run")
	}
	if n := runs(); n != 1 {
		t.Fatalf("the same scan must reuse the reading, the script ran %d times", n)
	}
	projectTests(root, &mapx.Graph{}, cfg)
	if n := runs(); n != 2 {
		t.Fatalf("a new map reads again, the script ran %d times", n)
	}
}

func TestProjectTests_TestsInFileOrder(t *testing.T) {
	t.Run("PRJTS-B04: The tests of a set of files come in file order then line", func(t *testing.T) {})
	all := []testlist.Test{
		{File: "a", Line: 9, Title: "a9"}, {File: "c", Line: 1, Title: "c1"},
		{File: "b", Line: 5, Title: "b5"}, {File: "a", Line: 2, Title: "a2"}, {File: "b", Line: 1, Title: "b1"},
	}
	var got []string
	for _, x := range testsIn(all, []string{"b", "a", "b"}) {
		got = append(got, x.Title)
	}
	if want := []string{"b1", "b5", "a2", "a9"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestProjectTests_ErrorIsReturned(t *testing.T) {
	t.Run("PRJTS-B05: A source error is returned with the declaration", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	cfg := &config.Config{Dialect: &config.Dialect{Tests: &config.TestsSource{Script: "exit 4"}}}
	root, g := t.TempDir(), &mapx.Graph{}
	_, declared, err := projectTests(root, g, cfg)
	if err == nil || !declared {
		t.Fatalf("the error must come back with the source declared, got %v (%v)", err, declared)
	}
	// The same scan asks again: the reused reading carries its error too.
	if _, _, err := projectTests(root, g, cfg); err == nil {
		t.Fatal("the reused reading must return the same error")
	}
}

func TestProjectTests_SupportIsNotATest(t *testing.T) {
	t.Run("PRJTS-B06: Support files are not read as tests", func(t *testing.T) {})
	resetProjectTestsCache()
	t.Cleanup(resetProjectTestsCache)
	root := t.TempDir()
	for _, f := range []string{"a_test.go", "helpers_test.go"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("t.Run(\"in "+f+"\", nil)\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a_test.go", Kind: mapx.KindTest}, {ID: "helpers_test.go", Kind: mapx.KindTest, Support: true}}}
	tests, _, err := projectTests(root, g, &config.Config{Dialect: &config.Dialect{Family: "go"}})
	if err != nil || len(tests) != 1 || tests[0].File != "a_test.go" {
		t.Fatalf("only the test file is read, got %v (%v)", tests, err)
	}
	if got := realTests(g, []string{"helpers_test.go", "a_test.go"}); len(got) != 1 || got[0] != "a_test.go" {
		t.Fatalf("the support file must leave the paths, got %v", got)
	}
	if got := realTests(&mapx.Graph{}, []string{"x_test.go"}); len(got) != 1 {
		t.Fatalf("with no support file the paths stay, got %v", got)
	}
}

func TestProjectTests_underIndexListsTheCommit(t *testing.T) {
	t.Run("PRJTS-B07: Under --index the tests are listed from the commit", func(t *testing.T) {})
	_, root := indexedRepo(t)
	staged := "t.Run(\"staged\", f)\n"
	writeFile(t, root, "a_test.go", staged)
	if out, err := exec.Command("git", "-C", root, "add", "a_test.go").CombinedOutput(); err != nil {
		t.Skipf("git add: %v %s", err, out)
	}
	writeFile(t, root, "a_test.go", strings.Repeat("\n", 40)+"t.Run(\"tree only\", f)\n")
	read, err := scan.IndexReader(root)
	if err != nil {
		t.Fatal(err)
	}
	defer SetFileSource(read)()
	resetProjectTestsCache()
	defer resetProjectTestsCache()
	cfg := &config.Config{Dialect: &config.Dialect{Family: "go", Tests: &config.TestsSource{Pattern: `t\.Run\(`}}}
	g := &mapx.Graph{Nodes: []mapx.Node{{ID: "a_test.go", Kind: mapx.KindTest}, {ID: "b_test.go", Kind: mapx.KindTest}}}
	tests, _, err := projectTests(root, g, cfg)
	if err != nil || len(tests) != 1 || tests[0].Title != "staged" || tests[0].Line != 1 {
		t.Errorf("a pattern scans the staged content, at its lines, got %+v %v", tests, err)
	}

	// a script reads the tree: what it lists of a file the tree holds otherwise is left out
	writeFile(t, root, "b_test.go", "t.Run(\"same\", f)\n")
	if out, err := exec.Command("git", "-C", root, "add", "b_test.go").CombinedOutput(); err != nil {
		t.Skipf("git add: %v %s", err, out)
	}
	read, err = scan.IndexReader(root)
	if err != nil {
		t.Fatal(err)
	}
	defer SetFileSource(read)()
	resetProjectTestsCache()
	script := `printf '{"version":1,"tests":[{"file":"a_test.go","line":41,"title":"tree only"},{"file":"b_test.go","line":1,"title":"same"}]}'`
	cfg = &config.Config{Dialect: &config.Dialect{Family: "go", Tests: &config.TestsSource{Script: script}}}
	tests, _, err = projectTests(root, g, cfg)
	if err != nil || len(tests) != 1 || tests[0].Title != "same" {
		t.Errorf("only the tests of a file the tree holds as staged, got %+v %v", tests, err)
	}
}

// indexedRepo is a repository with the project at `app`: `app/a.ts` committed, staged with
// another text and edited again in the tree, and `app/new.ts` untracked.
func indexedRepo(t *testing.T) (top, root string) {
	t.Helper()
	top = t.TempDir()
	if real, err := filepath.EvalSymlinks(top); err == nil {
		top = real
	}
	root = filepath.Join(top, "app")
	git := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", append([]string{"-C", top, "-c", "user.email=t@t", "-c", "user.name=t"}, args...)...).CombinedOutput(); err != nil {
			t.Skipf("git %v: %v %s", args, err, out)
		}
	}
	git("init", "-q")
	writeFile(t, root, "a.ts", "committed\n")
	git("add", ".")
	git("commit", "-qm", "base")
	writeFile(t, root, "a.ts", "staged\n")
	git("add", "app/a.ts")
	writeFile(t, root, "a.ts", "half done\n")
	writeFile(t, root, "new.ts", "untracked\n")
	return top, root
}
