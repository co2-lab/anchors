package changelog

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestClassify_Breaking(t *testing.T) {
	t.Run("CHNGL-B01: A breaking change is never silent", func(t *testing.T) {})
	r := Classify("v1", "", []Commit{
		{Hash: "a1", Subject: "refactor(api)!: drop the v1 routes"},
		{Hash: "b2", Subject: "feat: new flags", Body: "why\n\nBREAKING CHANGE: --x is gone"},
	})
	if len(r.Breaking) != 2 || len(r.Features) != 0 {
		t.Fatalf("both are breaking changes, got %+v", r)
	}
	if r.Breaking[1].Subject != "new flags — --x is gone" || r.Breaking[0].Subject != "drop the v1 routes" {
		t.Errorf("the footer's text follows the subject, got %q / %q", r.Breaking[0].Subject, r.Breaking[1].Subject)
	}
}

func TestClassify_Feature(t *testing.T) {
	t.Run("CHNGL-B02: A feat is a feature", func(t *testing.T) {})
	r := Classify("v1", "", []Commit{{Hash: "c3", Subject: "feat(gate): a new check"}})
	if len(r.Features) != 1 || r.Features[0].Scope != "gate" || r.Features[0].Hash != "c3" || r.Features[0].Subject != "a new check" {
		t.Fatalf("a feature with its scope and hash, got %+v", r.Features)
	}
}

func TestClassify_BugAndFix(t *testing.T) {
	t.Run("CHNGL-B03: A fix with a Bug footer is a bug fixed, without it a fix", func(t *testing.T) {})
	r := Classify("v1", "", []Commit{
		{Hash: "d4", Subject: "fix: the count", Body: "the count was off\n\nBug: v0.1.3, the board"},
		{Hash: "e5", Subject: "fix: my own typo"},
	})
	if len(r.Bugs) != 1 || r.Bugs[0].Bug != "v0.1.3, the board" || r.Bugs[0].Hash != "d4" {
		t.Errorf("a bug fixed carries where it was seen, got %+v", r.Bugs)
	}
	if len(r.Fixes) != 1 || r.Fixes[0].Hash != "e5" {
		t.Errorf("a fix without the footer is a fix, got %+v", r.Fixes)
	}
}

func TestClassify_LeftOut(t *testing.T) {
	t.Run("CHNGL-B04: The internal types and free-form subjects are left out", func(t *testing.T) {})
	r := Classify("v1", "", []Commit{{Subject: "chore: bump"}, {Subject: "refactor: move"}, {Subject: "Merge branch x"}})
	if !r.Empty() {
		t.Fatalf("nothing is listed, got %+v", r)
	}
	if (Release{Fixes: []Entry{{}}}).Empty() {
		t.Error("a release with a fix is not empty")
	}
}

func TestClassify_FootersInLastParagraph(t *testing.T) {
	t.Run("CHNGL-B05: Only the last paragraph holds footers", func(t *testing.T) {})
	r := Classify("v1", "", []Commit{{Subject: "fix: x", Body: "Bug: this is prose\n\nmore prose here"}})
	if len(r.Fixes) != 1 || len(r.Bugs) != 0 {
		t.Fatalf("a Bug line mid-body is not a footer, got %+v", r)
	}
}

// repo builds a repository whose history is the given steps: a subject commits, a
// `tag:NAME` tags the last commit.
func repo(t *testing.T, steps ...string) string {
	t.Helper()
	dir := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "-q", "-b", "main")
	for _, s := range steps {
		if name, ok := strings.CutPrefix(s, "tag:"); ok {
			git("tag", name)
			continue
		}
		git("commit", "-q", "--allow-empty", "-m", s)
	}
	return dir
}

func subjects(es []Entry) []string {
	var out []string
	for _, e := range es {
		out = append(out, e.Subject)
	}
	return out
}

func TestReleases_BetweenTags(t *testing.T) {
	t.Run("CHNGL-B06: Each release holds the commits after the previous tag", func(t *testing.T) {})
	g := Git{Root: repo(t, "feat: one", "tag:v0.1.0", "feat: two", "feat: three", "tag:v0.2.0", "feat: after")}
	tags, err := g.Tags("HEAD")
	if err != nil || strings.Join(tags, ",") != "v0.1.0,v0.2.0" {
		t.Fatalf("tags in version order, got %v %v", tags, err)
	}
	rels, err := g.Releases(tags, tags, "HEAD", false)
	if err != nil || len(rels) != 2 {
		t.Fatalf("two releases, got %v %v", rels, err)
	}
	if rels[0].Version != "v0.2.0" || strings.Join(subjects(rels[0].Features), ",") != "two,three" || rels[0].Date == "" {
		t.Errorf("the newest first, holding its own commits oldest first, got %+v", rels[0])
	}
	if rels[1].Version != "v0.1.0" || strings.Join(subjects(rels[1].Features), ",") != "one" {
		t.Errorf("the first release holds what came before it, got %+v", rels[1])
	}
	only, _ := g.Releases(tags, tags[1:], "HEAD", false)
	if len(only) != 1 || len(only[0].Features) != 2 {
		t.Errorf("a sub-list still starts from the previous tag of the whole list, got %+v", only)
	}
}

func TestReleases_Unreleased(t *testing.T) {
	t.Run("CHNGL-B07: What came after the last tag is unreleased", func(t *testing.T) {})
	g := Git{Root: repo(t, "feat: one", "tag:v0.1.0", "fix: after")}
	tags, _ := g.Tags("HEAD")
	rels, err := g.Releases(tags, tags, "HEAD", true)
	if err != nil || len(rels) != 2 || rels[0].Version != "" || len(rels[0].Fixes) != 1 || rels[1].Version != "v0.1.0" {
		t.Fatalf("an unreleased release on top, got %+v %v", rels, err)
	}
	none, err := Git{Root: repo(t, "feat: first")}.Releases(nil, nil, "HEAD", true)
	if err != nil || len(none) != 1 || len(none[0].Features) != 1 {
		t.Errorf("with no tag, the whole history is unreleased, got %+v %v", none, err)
	}
}

func tr(k string) string { return "<" + k + ">" }

func TestRender_Default(t *testing.T) {
	t.Run("CHNGL-B08: The built-in template translates its headings", func(t *testing.T) {})
	r := Release{Version: "v1.2.0", Date: "2026-09-27",
		Features: []Entry{{Scope: "gate", Subject: "a check", Hash: "a1"}},
		Bugs:     []Entry{{Subject: "the count", Hash: "b2", Bug: "v1.1.0, the board"}}}
	got, err := Render(DefaultTemplate, r, tr)
	if err != nil {
		t.Fatal(err)
	}
	want := "## v1.2.0 — 2026-09-27\n\n### <changelog.features>\n\n- **gate**: a check (a1)\n\n### <changelog.bugs>\n\n- the count (b2) — v1.1.0, the board\n"
	if got != want {
		t.Fatalf("want\n%q\ngot\n%q", want, got)
	}
	un, _ := Render(DefaultTemplate, Release{Fixes: []Entry{{Subject: "x", Hash: "c3"}}}, tr)
	if !strings.HasPrefix(un, "## <changelog.unreleased>\n\n### <changelog.fixes>\n") || strings.Contains(un, "breaking") {
		t.Errorf("unreleased heading, only the non-empty sections, got %q", un)
	}
}

func TestWritten(t *testing.T) {
	t.Run("CHNGL-B09: A written release is marked by its version", func(t *testing.T) {})
	out := Prepend("", "Changelog", []Block{{Version: "", Text: "u"}, {Version: "v2", Text: "b"}, {Version: "v1", Text: "a"}})
	got := Written(out)
	if len(got) != 3 || !got["v1"] || !got["v2"] || !got["unreleased"] {
		t.Fatalf("each written version is listed, got %v\n%s", got, out)
	}
	if !strings.Contains(out, Marker("v2")+"\nb\n") {
		t.Errorf("a block starts with its marker, got %q", out)
	}
}

func TestPrepend_KeepsWhatItHolds(t *testing.T) {
	t.Run("CHNGL-B10: New releases go on top and the file keeps what it holds", func(t *testing.T) {})
	existing := "# History\n\n" + Marker("v1") + "\n## v1 edited by hand\n\nOld notes, written before anchors.\n"
	got := Prepend(existing, "Changelog", []Block{{Version: "v2", Text: "## v2\n"}, {Version: "v1", Text: "## v1 regenerated\n"}})
	want := "# History\n\n" + Marker("v2") + "\n## v2\n\n" + Marker("v1") + "\n## v1 edited by hand\n\nOld notes, written before anchors.\n"
	if got != want {
		t.Fatalf("want\n%q\ngot\n%q", want, got)
	}
	if again := Prepend(got, "Changelog", []Block{{Version: "v2", Text: "x"}}); again != got {
		t.Errorf("a release the file holds is not written again, got %q", again)
	}
	untitled := Prepend("hand-written\n", "Changelog", []Block{{Version: "v1", Text: "## v1"}})
	if untitled != Marker("v1")+"\n## v1\n\nhand-written\n" {
		t.Errorf("without a title the releases go at the very top, got %q", untitled)
	}
}

func TestPrepend_ReplacesUnreleased(t *testing.T) {
	t.Run("CHNGL-B11: The unreleased block is replaced", func(t *testing.T) {})
	existing := "# C\n\n" + Marker("") + "\n## old unreleased\n\n" + Marker("v1") + "\n## v1\n"
	got := Prepend(existing, "C", []Block{{Text: "## new unreleased\n"}})
	want := "# C\n\n" + Marker("") + "\n## new unreleased\n\n" + Marker("v1") + "\n## v1\n"
	if got != want {
		t.Fatalf("want\n%q\ngot\n%q", want, got)
	}
	last := Prepend("# C\n\n"+Marker("")+"\n## old\n", "C", []Block{{Text: "## new"}})
	if strings.Contains(last, "old") {
		t.Errorf("an unreleased block at the end is replaced up to the end, got %q", last)
	}
	kept := Prepend(existing, "C", []Block{{Version: "v2", Text: "## v2"}})
	if !strings.Contains(kept, "old unreleased") {
		t.Errorf("without a new unreleased block the old one stays, got %q", kept)
	}
}

func TestPrepend_Title(t *testing.T) {
	t.Run("CHNGL-B12: An empty file gets a title", func(t *testing.T) {})
	if got := Prepend("", "Registro", []Block{{Version: "v1", Text: "## v1\n"}}); got != "# Registro\n\n"+Marker("v1")+"\n## v1\n" {
		t.Fatalf("got %q", got)
	}
}

func TestRender_BadTemplate(t *testing.T) {
	t.Run("CHNGL-E01: A template that does not parse is an error", func(t *testing.T) {})
	if _, err := Render("{{.Version", Release{}, tr); err == nil || !strings.Contains(err.Error(), "changelog template") {
		t.Errorf("a parse error names the template, got %v", err)
	}
	if _, err := Render("{{.Nope}}", Release{}, tr); err == nil || !strings.Contains(err.Error(), "changelog template") {
		t.Errorf("an execution error names the template, got %v", err)
	}
}

func TestGit_Error(t *testing.T) {
	t.Run("CHNGL-E02: A git failure names the command", func(t *testing.T) {})
	_, err := Git{Root: t.TempDir()}.Tags("HEAD")
	if err == nil || !strings.Contains(err.Error(), "git tag") {
		t.Fatalf("the error names the git command, got %v", err)
	}
	if _, err := (Git{Root: t.TempDir()}).Commits("", "HEAD"); err == nil {
		t.Error("log in a non-repository fails")
	}
}
