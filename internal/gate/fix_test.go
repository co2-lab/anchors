package gate

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/mapx"
)

// fixRepo creates a real temporary git repository holding `rel` with `content`,
// committed on 2024-03-05 (author date, which is what LastCommitDate reads).
func fixRepo(t *testing.T, rel, content string) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, rel), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t",
			"-c", "commit.gpgsign=false"}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null",
			"GIT_AUTHOR_DATE=2024-03-05T12:00:00", "GIT_COMMITTER_DATE=2024-03-05T12:00:00")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}
	git("init", "-q")
	git("add", rel)
	git("commit", "-q", "-m", "seed")
	return dir
}

func fixGate() config.Gate {
	return config.Gate{Name: "updated-at", Check: "updated-at-atual", On: []string{"spec"}}
}

func TestFixable(t *testing.T) {
	if !Fixable("updated-at-atual") {
		t.Error("updated-at-atual has a registered fixer")
	}
	if Fixable("header-conforme") {
		t.Error("a check without a fixer is not fixable")
	}
}

func TestFix_rewritesAStaleDateToTheCommitDate(t *testing.T) {
	const rel = "a.spec.md"
	dir := fixRepo(t, rel, "<!-- @anchors\n  updated_at: 2020-01-01\n-->\nbody\n")
	nodes := []mapx.Node{{ID: rel, Kind: mapx.KindSpec}}

	got := Fix([]config.Gate{fixGate()}, nodes, dir)

	if len(got) != 1 || !got[0].Fixed || got[0].Gate != "updated-at" || got[0].Target != rel {
		t.Fatalf("expected one successful fix of %s, got %+v", rel, got)
	}
	b, _ := os.ReadFile(filepath.Join(dir, rel))
	if want := "<!-- @anchors\n  updated_at: 2024-03-05\n-->\nbody\n"; string(b) != want {
		t.Errorf("only the date is replaced, by the commit date:\n got %q\nwant %q", b, want)
	}
}

func TestFix_leavesACorrectDateAlone(t *testing.T) {
	const rel = "a.spec.md"
	content := "<!-- @anchors\n  updated_at: 2024-03-05\n-->\n"
	dir := fixRepo(t, rel, content)

	if got := Fix([]config.Gate{fixGate()}, []mapx.Node{{ID: rel, Kind: mapx.KindSpec}}, dir); len(got) != 0 {
		t.Errorf("a date that already matches the commit is not a fix: %+v", got)
	}
}

func TestFix_skipsGatesWithoutFixerNodesOutOfScopeAndMissingFiles(t *testing.T) {
	const rel = "a.spec.md"
	dir := fixRepo(t, rel, "<!-- @anchors\n  updated_at: 2020-01-01\n-->\n")
	other := config.Gate{Name: "x", Check: "header-conforme", On: []string{"spec"}}
	nodes := []mapx.Node{
		{ID: rel, Kind: mapx.KindCode},            // the gate is not `on` code
		{ID: "gone.spec.md", Kind: mapx.KindSpec}, // not on disk
	}
	if got := Fix([]config.Gate{other, fixGate()}, nodes, dir); len(got) != 0 {
		t.Errorf("nothing applies, nothing is fixed: %+v", got)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, rel)); string(b) != "<!-- @anchors\n  updated_at: 2020-01-01\n-->\n" {
		t.Errorf("an out-of-scope file must not be rewritten: %q", b)
	}
}

func TestFix_reportsAFailedWrite(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root writes through read-only permissions")
	}
	const rel = "a.spec.md"
	dir := fixRepo(t, rel, "<!-- @anchors\n  updated_at: 2020-01-01\n-->\n")
	path := filepath.Join(dir, rel)
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })

	got := Fix([]config.Gate{fixGate()}, []mapx.Node{{ID: rel, Kind: mapx.KindSpec}}, dir)

	if len(got) != 1 || got[0].Fixed || got[0].Detail == "" {
		t.Fatalf("a write that fails is reported as not fixed, with the cause: %+v", got)
	}
}

func TestFixUpdatedAt(t *testing.T) {
	const rel = "a.spec.md"
	n := mapx.Node{ID: rel, Kind: mapx.KindSpec}

	t.Run("without the field there is nothing to correct", func(t *testing.T) {
		content := "<!-- @anchors\n  code: ABCDE\n-->\n"
		dir := fixRepo(t, rel, content)
		if got, changed := fixUpdatedAt(content, n, dir); changed || got != content {
			t.Errorf("the fixer must not create the field: %q %v", got, changed)
		}
	})

	t.Run("a file with a pending edit takes today's date", func(t *testing.T) {
		dir := fixRepo(t, rel, "<!-- @anchors\n  updated_at: 2024-03-05\n-->\n")
		edited := "<!-- @anchors\n  updated_at: 2024-03-05\n-->\nnew line\n"
		if err := os.WriteFile(filepath.Join(dir, rel), []byte(edited), 0o644); err != nil {
			t.Fatal(err)
		}
		got, changed := fixUpdatedAt(edited, n, dir)
		want := "<!-- @anchors\n  updated_at: " + gitmeta.Today() + "\n-->\nnew line\n"
		if !changed || got != want {
			t.Errorf("got %q (%v), want %q", got, changed, want)
		}
	})

	t.Run("outside a repository there is no correct date to write", func(t *testing.T) {
		dir := t.TempDir()
		if gitmeta.Check(dir) == gitmeta.Disponível {
			t.Skip("the temp dir sits inside a git repository")
		}
		content := "updated_at: 2020-01-01\n"
		if got, changed := fixUpdatedAt(content, n, dir); changed || got != content {
			t.Errorf("without git the fixer must not guess: %q %v", got, changed)
		}
	})

	t.Run("an ignored file that was never committed is left alone", func(t *testing.T) {
		dir := fixRepo(t, rel, "x\n")
		content := "updated_at: 2020-01-01\n"
		// never committed, and ignored so status reports nothing either
		if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("new.spec.md\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "new.spec.md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if got, changed := fixUpdatedAt(content, mapx.Node{ID: "new.spec.md"}, dir); changed || got != content {
			t.Errorf("no commit and no edit: nothing to compare against: %q %v", got, changed)
		}
	})
}
