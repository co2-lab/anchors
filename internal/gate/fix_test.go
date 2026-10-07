// @anchors
//   code: FXTSA
//   ref: FXIXX

package gate

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
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
	t.Run("FXIXX-B01: Only a check with a registered fixer is fixable", func(t *testing.T) {})
	if !Fixable("updated-at-atual") {
		t.Error("updated-at-atual has a registered fixer")
	}
	if Fixable("header-conforme") {
		t.Error("a check without a fixer is not fixable")
	}
}

func TestFix_rewritesAStaleDateToTheCommitDate(t *testing.T) {
	t.Run("FXIXX-B02: A stale date on a committed file is rewritten to its last commit date", func(t *testing.T) {})
	t.Run("FXIXX-I01: Only the date changes, the rest of the file is kept byte for byte", func(t *testing.T) {})
	const rel = "a.spec.md"
	dir := fixRepo(t, rel, "<!-- @anchors\n  updated_at: 2020-01-01\n-->\nbody\n")
	nodes := []mapx.Node{{ID: rel, Kind: mapx.KindSpec}}

	got := Fix([]config.Gate{fixGate()}, nodes, dir, nil)

	if len(got) != 1 || !got[0].Fixed || got[0].Gate != "updated-at" || got[0].Target != rel {
		t.Fatalf("expected one successful fix of %s, got %+v", rel, got)
	}
	b, _ := os.ReadFile(filepath.Join(dir, rel))
	if want := "<!-- @anchors\n  updated_at: 2024-03-05\n-->\nbody\n"; string(b) != want {
		t.Errorf("only the date is replaced, by the commit date:\n got %q\nwant %q", b, want)
	}
}

// The detail of a repair is written in the project's language. It was a Portuguese
// literal ("updated_at corrigido para a data do commit") in every project.
func TestFix_detailIsTranslated(t *testing.T) {
	t.Run("FXIXX-B08: The repair detail is written in the project's language", func(t *testing.T) {})
	const rel = "a.spec.md"
	t.Cleanup(func() { i18n.Set(i18n.Default) })
	for lang, want := range map[string]string{
		"en":    i18n.TIn("en", "gate.fix.updated_at_fixed"),
		"pt-BR": i18n.TIn("pt-BR", "gate.fix.updated_at_fixed"),
	} {
		i18n.Set(lang)
		dir := fixRepo(t, rel, "<!-- @anchors\n  updated_at: 2020-01-01\n-->\n")
		got := Fix([]config.Gate{fixGate()}, []mapx.Node{{ID: rel, Kind: mapx.KindSpec}}, dir, nil)
		if len(got) != 1 || got[0].Detail != want || want == "" || strings.HasPrefix(want, "gate.fix.") {
			t.Errorf("[%s] detail %+v, want %q", lang, got, want)
		}
	}
	if i18n.TIn("en", "gate.fix.updated_at_fixed") == i18n.TIn("pt-BR", "gate.fix.updated_at_fixed") {
		t.Errorf("the two languages must differ, or the case proves nothing")
	}
}

func TestFix_leavesACorrectDateAlone(t *testing.T) {
	t.Run("FXIXX-B04: A date that already matches is left alone and reported as nothing", func(t *testing.T) {})
	const rel = "a.spec.md"
	content := "<!-- @anchors\n  updated_at: 2024-03-05\n-->\n"
	dir := fixRepo(t, rel, content)

	if got := Fix([]config.Gate{fixGate()}, []mapx.Node{{ID: rel, Kind: mapx.KindSpec}}, dir, nil); len(got) != 0 {
		t.Errorf("a date that already matches the commit is not a fix: %+v", got)
	}
}

func TestFix_skipsGatesWithoutFixerNodesOutOfScopeAndMissingFiles(t *testing.T) {
	t.Run("FXIXX-B05: Only fixable gates, the nodes they apply to and files on disk are touched", func(t *testing.T) {})
	const rel = "a.spec.md"
	dir := fixRepo(t, rel, "<!-- @anchors\n  updated_at: 2020-01-01\n-->\n")
	other := config.Gate{Name: "x", Check: "header-conforme", On: []string{"spec"}}
	nodes := []mapx.Node{
		{ID: rel, Kind: mapx.KindCode},            // the gate is not `on` code
		{ID: "gone.spec.md", Kind: mapx.KindSpec}, // not on disk
	}
	if got := Fix([]config.Gate{other, fixGate()}, nodes, dir, nil); len(got) != 0 {
		t.Errorf("nothing applies, nothing is fixed: %+v", got)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, rel)); string(b) != "<!-- @anchors\n  updated_at: 2020-01-01\n-->\n" {
		t.Errorf("an out-of-scope file must not be rewritten: %q", b)
	}
}

func TestFix_reportsAFailedWrite(t *testing.T) {
	t.Run("FXIXX-E01: A write that fails is reported as not fixed, with the cause", func(t *testing.T) {})
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

	got := Fix([]config.Gate{fixGate()}, []mapx.Node{{ID: rel, Kind: mapx.KindSpec}}, dir, nil)

	if len(got) != 1 || got[0].Fixed || got[0].Detail == "" {
		t.Fatalf("a write that fails is reported as not fixed, with the cause: %+v", got)
	}
}

func TestFixUpdatedAt(t *testing.T) {
	const rel = "a.spec.md"
	n := mapx.Node{ID: rel, Kind: mapx.KindSpec}

	t.Run("without the field there is nothing to correct", func(t *testing.T) {
		t.Run("FXIXX-X01: A header without the date field is never given one", func(t *testing.T) {})
		content := "<!-- @anchors\n  code: ABCDE\n-->\n"
		dir := fixRepo(t, rel, content)
		if got, changed := fixUpdatedAt(content, n, dir); changed || got != content {
			t.Errorf("the fixer must not create the field: %q %v", got, changed)
		}
	})

	t.Run("a file with a pending edit takes today's date", func(t *testing.T) {
		t.Run("FXIXX-B03: A file with an uncommitted edit takes today's date", func(t *testing.T) {})
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
		t.Run("FXIXX-B06: Outside a git repository the file is left untouched", func(t *testing.T) {})
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
		t.Run("FXIXX-B07: A file never committed and not edited is left untouched", func(t *testing.T) {})
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

func TestFixMissingHeader_fromTheMap(t *testing.T) {
	t.Run("FXIXX-B09: The fix writes the missing header from the map", func(t *testing.T) {})
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "src/pay.spec.md", Kind: mapx.KindSpec, Code: "PAYMT"},
			{ID: "src/pay.feature", Kind: mapx.KindFeature},
			{ID: "src/pay.go", Kind: mapx.KindCode},
			{ID: "src/pay_test.go", Kind: mapx.KindTest},
		},
		Edges: []mapx.Edge{
			{From: "src/pay.spec.md", To: "src/pay.go", Type: mapx.EdgeSpecifies},
			{From: "src/pay.spec.md", To: "src/pay.feature", Type: mapx.EdgeCoveredBy},
			{From: "src/pay.feature", To: "src/pay_test.go", Type: mapx.EdgeTestedBy},
		},
	}
	cases := []struct {
		n      mapx.Node
		in     string
		prefix string
	}{
		{mapx.Node{ID: "src/pay.go", Kind: mapx.KindCode}, "package src\n", "// @anchors\n//   code: CODE\n//   ref: PAYMT\n\npackage src\n"},
		{mapx.Node{ID: "src/pay_test.go", Kind: mapx.KindTest}, "package src\n", "// @anchors\n//   code: CODE\n//   ref: PAYMT\n\n"},
		{mapx.Node{ID: "guides/A.md", Kind: mapx.KindGuide, Layer: "guide"}, "# A\n", "<!-- @anchors\n  code: CODE\n  layer: guide\n-->\n\n# A\n"},
		{mapx.Node{ID: "src/pay.go", Kind: mapx.KindCode}, "// @anchors\n//   updated_at: 2026-01-01\npackage src\n", "// @anchors\n//   code: CODE\n//   ref: PAYMT\n//   updated_at: 2026-01-01\n"},
	}
	for _, c := range cases {
		out, changed, _ := fixMissingHeader(c.in, c.n, "", g)
		out = generatedCodeRE.ReplaceAllString(out, "code: CODE")
		if !changed || !strings.HasPrefix(out, c.prefix) {
			t.Errorf("%s: got %q, want it to start with %q", c.n.ID, out, c.prefix)
		}
	}
	sh := mapx.Node{ID: "src/run.sh", Kind: mapx.KindCode}
	gs := &mapx.Graph{Nodes: []mapx.Node{{ID: "s.spec.md", Kind: mapx.KindSpec, Code: "RUNXX"}, sh}, Edges: []mapx.Edge{{From: "s.spec.md", To: "src/run.sh", Type: mapx.EdgeSpecifies}}}
	if out, _, _ := fixMissingHeader("#!/bin/sh\necho hi\n", sh, "", gs); !strings.HasPrefix(generatedCodeRE.ReplaceAllString(out, "code: CODE"), "#!/bin/sh\n# @anchors\n#   code: CODE\n#   ref: RUNXX\n") {
		t.Errorf("the shebang stays first: %q", out)
	}
	if _, changed, _ := fixMissingHeader("package x\n", mapx.Node{ID: "x/orphan.go", Kind: mapx.KindCode}, "", g); changed {
		t.Error("a file of no unit is left as it is")
	}
	t.Run("FXIXX-B10: The fix gives a file with an identity and no code of its own a code unique in the map", func(t *testing.T) {})
	in := "// @anchors\n//   ref: PAYMT\n\npackage src\n"
	out, changed, _ := fixMissingHeader(in, mapx.Node{ID: "src/pay.go", Kind: mapx.KindCode, Layer: "logic"}, "", g)
	m := regexp.MustCompile(`code: ([A-Z0-9]{5})\n`).FindStringSubmatch(out)
	if !changed || m == nil || m[1] == "PAYMT" || !strings.Contains(out, "ref: PAYMT") {
		t.Errorf("a code of its own, beside the ref: %q", out)
	}
	if again, changed, _ := fixMissingHeader(out, mapx.Node{ID: "src/pay.go", Kind: mapx.KindCode}, "", g); changed || again != out {
		t.Errorf("a header with its own code is left as it is: %q", again)
	}
}

var generatedCodeRE = regexp.MustCompile(`code: [A-Z0-9]{5}`)

func TestFixWithConfig_theDialectAndTheMovedLines(t *testing.T) {
	t.Run("FXIXX-B11: FixWithConfig gives the fixers the project's config for that run", func(t *testing.T) {})
	t.Run("FXIXX-B12: Each repair says whether it moved lines", func(t *testing.T) {})
	root, g, cfg := chainProject(t)
	gates := []config.Gate{{Name: "dep-declared", Check: "dep-declared", On: []string{"code"}}, {Name: "used-by-declared", Check: "used-by-declared", On: []string{"code"}}}
	moved := map[string]bool{}
	fixed := map[string]bool{}
	for _, r := range FixWithConfig(gates, g.Nodes, root, g, cfg) {
		fixed[r.Target] = fixed[r.Target] || r.Fixed
		moved[r.Target] = moved[r.Target] || r.LinesMoved
	}
	if !fixed["src/ui/Arena.tsx"] {
		t.Error("with the project's dialect, the fixer reads the screen's imports and flags them")
	}
	if fixConfig != nil {
		t.Error("the config is the fixers' only for that run")
	}
	if moved["src/ui/Arena.tsx"] || !moved["src/theme/tokens.ts"] {
		t.Errorf("a flag on the import's line moves no line, a used-by line inserted does: %+v", moved)
	}
	if !linesMoved("a\nb\nc", "x\na\nc") || linesMoved("a // @dep: X\nb", "a // @dep: Y\nb") {
		t.Error("an insertion with a removal moves lines; a flag rewritten in place does not")
	}
}
