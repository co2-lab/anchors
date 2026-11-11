// @anchors
//   code: PRCMT
//   ref: PRCMS

package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

const pushReviewYAML = "version: 2\nlayers: {}\npush_review:\n  required: true\n  retention: 30d\n"

// pushRepo is a git repository with the project committed as its base, and one commit on top
// changing a.go. It returns the root and the base's hash.
func pushRepo(t *testing.T, yaml string, g *mapx.Graph) (string, string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	// The review's commit is made with the user's identity; a CI runner has none.
	for _, k := range []string{"GIT_AUTHOR_NAME", "GIT_COMMITTER_NAME"} {
		t.Setenv(k, "t")
	}
	for _, k := range []string{"GIT_AUTHOR_EMAIL", "GIT_COMMITTER_EMAIL"} {
		t.Setenv(k, "t@t")
	}
	root := qProject(t, yaml, checkFiles(), g)
	gitIn(t, root, "init", "-q")
	gitIn(t, root, "add", ".")
	gitIn(t, root, "commit", "-q", "-m", "chore: base")
	base := pushHead(t, root)
	writeQ(t, root, "a.go", "package a\n\nconst x = 1\n")
	gitIn(t, root, "commit", "-q", "-am", "feat: x")
	return root, base
}

func pushHead(t *testing.T, root string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func writeQ(t *testing.T, root, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func lastSubject(t *testing.T, root string) string {
	t.Helper()
	out, _ := exec.Command("git", "-C", root, "log", "-1", "--format=%s").Output()
	return strings.TrimSpace(string(out))
}

func TestPushReview_records(t *testing.T) {
	t.Run("PRCMS-B01: A review is written to the reviews' file and committed alone", func(t *testing.T) {})
	t.Run("PRCMS-B07: Writing a review drops the reviews older than the retention", func(t *testing.T) {})
	englishOutput(t)
	root, base := pushRepo(t, pushReviewYAML, nil)
	writeQ(t, root, "REVIEWS.md", "# Push reviews\n\n## 2026-01-01 10:00 UTC · ok · old\n<!-- push-review at=\"2026-01-01T10:00:00Z\" verdict=\"ok\" by=\"human:ana\" -->\n\n- `aaaaaaa` old\n")
	writeQ(t, root, "b.go", "package a\n\nconst y = 2\n")
	out, err := runQ(t, newPushReviewCmd(), "--root", root, "--range", base+"..HEAD", "--ok", "the unit tests and the spec", "--by", "agent:claude", "--note", "no UI touched")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if !strings.Contains(out, "review ok of 1 commit(s)") {
		t.Errorf("the record says what it wrote:\n%s", out)
	}
	body := readQ(t, filepath.Join(root, "REVIEWS.md"))
	for _, want := range []string{"· ok · " + base + "..HEAD", `by="agent:claude"`, " feat: x\n", "**Validated**\n\nthe unit tests and the spec", "- no UI touched"} {
		if !strings.Contains(body, want) {
			t.Errorf("the file lacks %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "· old") {
		t.Errorf("a review older than the retention stays:\n%s", body)
	}
	if s := lastSubject(t, root); s != "chore(review): ok "+base+"..HEAD" {
		t.Errorf("the review's commit: %q", s)
	}
	st, _ := exec.Command("git", "-C", root, "status", "--porcelain").Output()
	if !strings.Contains(string(st), "b.go") || strings.Contains(string(st), "REVIEWS.md") {
		t.Errorf("the commit takes the reviews' file alone: %s", st)
	}
}

func TestPushReview_refusesWithoutVerdictOrAuthor(t *testing.T) {
	t.Run("PRCMS-E01: A review with no verdict, two verdicts or no author is refused", func(t *testing.T) {})
	englishOutput(t)
	root, base := pushRepo(t, pushReviewYAML, nil)
	for _, args := range [][]string{
		{"--by", "x"},
		{"--ok", "a", "--fail", "b", "--by", "x"},
		{"--ok", "a"},
	} {
		if _, err := runQ(t, newPushReviewCmd(), append([]string{"--root", root, "--range", base + "..HEAD"}, args...)...); err == nil {
			t.Errorf("%v was accepted", args)
		}
	}
	if _, err := runQ(t, newPushReviewCmd(), "--root", root, "--range", "HEAD..HEAD", "--ok", "a", "--by", "x"); err == nil || !strings.Contains(err.Error(), "nothing to review") {
		t.Errorf("an empty range: %v", err)
	}
}

func TestPushReviewVerify_asksAnOkReview(t *testing.T) {
	t.Run("PRCMS-B02: A push whose commits no review covers is refused, naming them", func(t *testing.T) {})
	t.Run("PRCMS-B03: A push whose latest review failed is refused with that review", func(t *testing.T) {})
	t.Run("PRCMS-B04: A push covered by an ok review, and the review's own commit, go through", func(t *testing.T) {})
	englishOutput(t)
	root, base := pushRepo(t, pushReviewYAML, nil)
	rng := base + "..HEAD"
	_, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng)
	if err == nil || !strings.Contains(err.Error(), "PUSH REFUSED") || !strings.Contains(err.Error(), "1 commit(s) no review covers") || !strings.Contains(err.Error(), "feat: x") {
		t.Fatalf("an unreviewed push: %v", err)
	}
	if _, err := runQ(t, newPushReviewCmd(), "--root", root, "--range", rng, "--fail", "the error path is untested", "--by", "agent:claude"); err != nil {
		t.Fatal(err)
	}
	_, err = runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng)
	if err == nil || !strings.Contains(err.Error(), "FAILED") || !strings.Contains(err.Error(), "the error path is untested") {
		t.Fatalf("a failed review: %v", err)
	}
	if _, err := runQ(t, newPushReviewCmd(), "--root", root, "--range", rng, "--ok", "the error path now has its test", "--by", "agent:claude"); err != nil {
		t.Fatal(err)
	}
	if _, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng); err != nil {
		t.Errorf("a push an ok review covers: %v", err)
	}
}

func TestPushReviewVerify_readsTheHookRefs(t *testing.T) {
	t.Run("PRCMS-B05: The refs git hands the hook say what the push sends; a new branch sends what no remote has, a deletion nothing", func(t *testing.T) {})
	englishOutput(t)
	root, base := pushRepo(t, pushReviewYAML, nil)
	head := pushHead(t, root)
	for name, in := range map[string]string{
		"an update":  "refs/heads/main " + head + " refs/heads/main " + base + "\n",
		"new branch": "refs/heads/main " + head + " refs/heads/main " + strings.Repeat("0", 40) + "\n",
	} {
		cmd := newPushReviewVerifyCmd()
		cmd.SetIn(strings.NewReader(in))
		if _, err := runQ(t, cmd, "--root", root); err == nil || !strings.Contains(err.Error(), "feat: x") {
			t.Errorf("%s: %v", name, err)
		}
	}
	cmd := newPushReviewVerifyCmd()
	cmd.SetIn(strings.NewReader("(delete) " + strings.Repeat("0", 40) + " refs/heads/old " + base + "\n"))
	if _, err := runQ(t, cmd, "--root", root); err != nil {
		t.Errorf("a deletion sends nothing: %v", err)
	}
}

func TestPushReviewVerify_offWithoutTheSection(t *testing.T) {
	t.Run("PRCMS-B06: Without push_review, the push confronts nothing", func(t *testing.T) {})
	root, base := pushRepo(t, "version: 2\nlayers: {}\n", nil)
	if _, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", base+"..HEAD"); err != nil {
		t.Errorf("a project that asks nothing: %v", err)
	}
}

func TestPushReview_judgments(t *testing.T) {
	t.Run("PRCMS-B08: A push that reaches a judgment gate no verdict answers in the pushed map is refused; an ok review is not recorded over it", func(t *testing.T) {})
	englishOutput(t)
	yaml := "version: 2\nlayers: {}\ngates:\n  - name: j\n    on: [code]\n    measures: judgment\npush_review:\n  require_judgment: true\n"
	root, base := pushRepo(t, yaml, checkGraph())
	rng := base + "..HEAD"
	_, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng)
	if err == nil || !strings.Contains(err.Error(), "j @ a.go") {
		t.Fatalf("an unanswered judgment: %v", err)
	}
	if _, err := runQ(t, newPushReviewCmd(), "--root", root, "--range", rng, "--ok", "x", "--by", "y", "--no-commit"); err == nil || !strings.Contains(err.Error(), "j @ a.go") {
		t.Errorf("an ok review over an unanswered judgment: %v", err)
	}
	g := checkGraph()
	g.Nodes[0].Judgments = []mapx.Judgment{{Gate: "j", Verdict: "ok", ValidatedFromRev: "r1"}}
	if err := mapx.Save(g, filepath.Join(root, mapx.DefaultPath)); err != nil {
		t.Fatal(err)
	}
	if _, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng); err == nil {
		t.Error("a verdict the push does not send answers nothing")
	}
	gitIn(t, root, "commit", "-q", "-am", "chore: judged")
	if _, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng); err != nil {
		t.Errorf("a judgment the pushed map answers: %v", err)
	}
}

func TestPushReviewVerify_gatesWithoutWaivers(t *testing.T) {
	t.Run("PRCMS-B09: With require_gates, a blocking gate failing over what the push sends refuses it, whatever a commit message waived, read as the pushed tip holds the files", func(t *testing.T) {})
	englishOutput(t)
	yaml := "version: 2\nlayers:\n  code:\n    kind: code\n    pattern: \"*.go\"\ngates:\n  - name: no-forbidden\n    on: [code]\n    blocking: true\n    run: \"! grep -q FORBIDDEN \\\"$1\\\"\"\npush_review:\n  require_gates: true\n"
	root, base := pushRepo(t, yaml, checkGraph())
	writeQ(t, root, "a.go", "package a\n\n// FORBIDDEN\n")
	gitIn(t, root, "commit", "-q", "-am", "feat: y\n\n[skip-no-forbidden@a.go: fixed before the push]")
	waiver := pushHead(t, root)[:8]
	rng := base + "..HEAD"
	_, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng)
	if err == nil || !strings.Contains(err.Error(), "no-forbidden @ a.go") || !strings.Contains(err.Error(), "waived in "+waiver+" feat: y") {
		t.Fatalf("a failure a commit message waived: %v", err)
	}
	writeQ(t, root, "a.go", "package a\n")
	if _, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng); err == nil {
		t.Error("a fix not committed is not what the push sends")
	}
	gitIn(t, root, "commit", "-q", "-am", "fix: y")
	if _, err := runQ(t, newPushReviewVerifyCmd(), "--root", root, "--range", rng); err != nil {
		t.Errorf("the committed fix lets the push through: %v", err)
	}
}
