// @anchors
//   code: PRVLT
//   ref: PRVLS

package pushreview

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var at = time.Date(2026, 11, 11, 12, 40, 0, 0, time.UTC)

func review(verdict string, hashes ...string) Review {
	var cs []Commit
	for _, h := range hashes {
		cs = append(cs, Commit{Hash: h, PatchID: "f" + h[:4] + "0000000", Subject: "fix: " + h[:4]})
	}
	return Review{At: at, Verdict: verdict, By: "agent:claude", Range: "aaaaaaa..bbbbbbb", Commits: cs, Text: "what was seen"}
}

func TestReview_entry(t *testing.T) {
	t.Run("PRVLS-B01: A review is written as an entry", func(t *testing.T) {})
	r := review(OK, "aaaa1111aaaa1111", "bbbb2222bbbb2222")
	r.Judgment = []string{"rule-test-coherent @ a.spec.md: ok"}
	r.Notes = []string{"the flows ran on iOS only"}
	e := r.Entry()
	for _, want := range []string{
		"## 2026-11-11 12:40 UTC · ok · aaaaaaa..bbbbbbb\n",
		`<!-- push-review at="2026-11-11T12:40:00Z" verdict="ok" by="agent:claude" -->`,
		"- `aaaa1111aaaa1111` `p:faaaa0000000` fix: aaaa\n",
		"**Judgment gates**\n\n- rule-test-coherent @ a.spec.md: ok\n",
		"**Validated**\n\nwhat was seen\n",
		"**Notes**\n\n- the flows ran on iOS only\n",
	} {
		if !strings.Contains(e.Raw, want) {
			t.Errorf("the entry lacks %q:\n%s", want, e.Raw)
		}
	}
	if e.Verdict != OK || e.By != "agent:claude" || len(e.Commits) != 2 || !e.At.Equal(at) {
		t.Errorf("what the machine reads: %+v", e)
	}
	if f := review(Fail, "aaaa1111aaaa1111").Entry(); !strings.Contains(f.Raw, "**Failed**\n\nwhat was seen") {
		t.Errorf("a failed review says what failed:\n%s", f.Raw)
	}
}

func TestParse_roundTrip(t *testing.T) {
	t.Run("PRVLS-B02: The file is read back into its entries", func(t *testing.T) {})
	newer, older := review(OK, "bbbb2222bbbb2222").Entry(), review(Fail, "aaaa1111aaaa1111").Entry()
	content := Render([]Entry{newer, older}) + "\n## Notes by hand\n\nnot a review\n"
	got := Parse(content)
	if len(got) != 2 || got[0].Verdict != OK || got[1].Verdict != Fail || got[1].By != "agent:claude" ||
		got[0].Commits[0] != newer.Commits[0] || !got[1].At.Equal(at) {
		t.Fatalf("the entries: %+v", got)
	}
	again := Render(got)
	if !strings.HasPrefix(again, "# Push reviews\n") || len(Parse(again)) != 2 || Render(Parse(again)) != again {
		t.Errorf("written again, the file reads back the same:\n%s", again)
	}
}

func TestKeep_retention(t *testing.T) {
	t.Run("PRVLS-B03: Writing a review drops the entries older than the retention", func(t *testing.T) {})
	es := []Entry{{At: at, Verdict: OK}, {At: at.AddDate(0, 0, -40), Verdict: OK}, {Verdict: Fail}}
	got := Keep(es, 30, at)
	if len(got) != 2 || !got[0].At.Equal(at) || !got[1].At.IsZero() {
		t.Errorf("thirty days: %+v", got)
	}
	if len(Keep(es, 0, at)) != 3 {
		t.Error("no retention keeps every review")
	}
}

func TestCovering(t *testing.T) {
	t.Run("PRVLS-B04: A commit is covered by the newest entry that names it", func(t *testing.T) {})
	es := []Entry{review(OK, "aaaa1111aaaa1111").Entry(), review(Fail, "aaaa1111aaaa1111").Entry()}
	if e, ok := Covering(es, Commit{Hash: "aaaa1111aaaa1111"}); !ok || e.Verdict != OK {
		t.Errorf("the newest review covers: %+v %v", e, ok)
	}
	if _, ok := Covering(es, Commit{Hash: "aaaa111"}); !ok {
		t.Error("a short hash covers")
	}
	if _, ok := Covering(es, Commit{Hash: "cccc3333cccc3333", PatchID: "faaaa0000000"}); !ok {
		t.Error("the patch id covers a rebased commit")
	}
	if _, ok := Covering(es, Commit{Hash: "dddd4444dddd4444", PatchID: "fdddd0000000"}); ok {
		t.Error("a commit no entry names is not covered")
	}
}

func TestExempt(t *testing.T) {
	t.Run("PRVLS-B05: A commit that changes only the evidence owes no review", func(t *testing.T) {})
	if !Exempt([]string{"REVIEWS.md", "anchors.graph.yaml"}, "REVIEWS.md", "anchors.graph.yaml") {
		t.Error("the evidence alone owes no review")
	}
	if Exempt([]string{"REVIEWS.md", "a.go"}, "REVIEWS.md", "anchors.graph.yaml") || Exempt(nil, "REVIEWS.md", "anchors.graph.yaml") {
		t.Error("a code file, or nothing, owes a review")
	}
}

func TestCommits_range(t *testing.T) {
	t.Run("PRVLS-B06: The commits of a range are read with their patch id and files", func(t *testing.T) {})
	d := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", d}, args...)...).CombinedOutput()
		if err != nil {
			t.Skipf("git: %v %s", err, out)
		}
		return strings.TrimSpace(string(out))
	}
	git("init", "-q")
	git("config", "user.email", "t@t.co")
	git("config", "user.name", "t")
	for i, f := range []string{"a.txt", "b.txt", "c/d.txt"} {
		os.MkdirAll(filepath.Dir(filepath.Join(d, f)), 0o755)
		os.WriteFile(filepath.Join(d, f), []byte(f), 0o644)
		git("add", f)
		git("commit", "-q", "-m", "feat: "+string(rune('a'+i)))
	}
	first := git("rev-list", "--max-parents=0", "HEAD")
	before, err := Commits(d, first+"..HEAD")
	if err != nil || len(before) != 2 || before[0].Subject != "feat: b" || before[1].Subject != "feat: c" || before[1].PatchID == "" {
		t.Fatalf("the range: %+v %v", before, err)
	}
	git("commit", "-q", "--amend", "-m", "feat: c, reworded")
	after, _ := Commits(d, first+"..HEAD")
	if after[1].Hash == before[1].Hash || after[1].PatchID != before[1].PatchID {
		t.Errorf("a reworded commit keeps its patch id: %+v %+v", before[1], after[1])
	}
	if fs := Files(d, after[1].Hash); len(fs) != 1 || fs[0] != "c/d.txt" {
		t.Errorf("the commit's files: %v", fs)
	}
}

func TestCommits_unreadable(t *testing.T) {
	t.Run("PRVLS-E01: A range git cannot list fails naming it", func(t *testing.T) {})
	t.Run("PRVLS-E02: A commit git gives no files or patch id is covered by hash and owes a review", func(t *testing.T) {})
	d := t.TempDir()
	for _, args := range [][]string{{"init", "-q"}, {"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "--allow-empty", "-m", "chore: empty"}} {
		if out, err := exec.Command("git", append([]string{"-C", d}, args...)...).CombinedOutput(); err != nil {
			t.Skipf("git: %v %s", err, out)
		}
	}
	if _, err := Commits(d, "nope..HEAD"); err == nil || !strings.Contains(err.Error(), "nope..HEAD") {
		t.Errorf("an unknown revision: %v", err)
	}
	cs, err := Commits(d, "HEAD")
	if err != nil || len(cs) != 1 || cs[0].PatchID != "" || len(Files(d, cs[0].Hash)) != 0 {
		t.Fatalf("an empty commit: %+v %v", cs, err)
	}
	if Files(d, "deadbeefdeadbeef") != nil || Exempt(Files(d, cs[0].Hash), "REVIEWS.md", "anchors.graph.yaml") {
		t.Error("a commit with no files owes a review")
	}
}
