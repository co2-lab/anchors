// @anchors
//   code: VAPLC
//   ref: VAPLA

package ops

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/changelog"
)

// fakePR answers the pull request calls: the open ones it holds, and what was asked.
type fakePR struct {
	open  string
	calls []string
}

func stubPR(t *testing.T) *fakePR {
	t.Helper()
	f := &fakePR{open: "[]"}
	old := prAPI
	prAPI = func(_ string, args ...string) ([]byte, error) {
		f.calls = append(f.calls, strings.Join(args, " "))
		if args[1] == "list" {
			return []byte(f.open), nil
		}
		return nil, nil
	}
	t.Cleanup(func() { prAPI = old })
	return f
}

// applyRepo is releasedRepo, on review, writing CHANGELOG.md from the commits and NOTES.md
// from the intents, which a release archives; a feat came since app-v1.0.0, with its intent.
func applyRepo(t *testing.T, apply string) string {
	t.Helper()
	root := releasedRepo(t)
	cfg := strings.Replace(releasedProject, "    next: semantic\n", "    next: semantic\n    apply: "+apply+"\n    changelog:\n      outputs:\n        - path: CHANGELOG.md\n        - { path: NOTES.md, source: intents }\n", 1)
	cfg = strings.Replace(cfg, "version: 8\n", "version: 8\nchangelog:\n  intents: { consume: archive }\n", 1)
	commitFiles(t, root, "feat: sharing", map[string]string{
		"anchors.yaml":              cfg,
		".anchors/changes/share.md": "---\nproduct: app\nbump: minor\ntype: feat\n---\nReceipts can be shared.\n",
	})
	return root
}

func applyRun(t *testing.T, root string, args ...string) (string, error) {
	t.Helper()
	return versionRun(t, append([]string{"--root", root, "--branch", "review", "apply"}, args...)...)
}

func currentBranch(t *testing.T, root string) string {
	t.Helper()
	out, _ := exec.Command("git", "-C", root, "rev-parse", "--abbrev-ref", "HEAD").Output()
	return strings.TrimSpace(string(out))
}

func TestApply_releasePR(t *testing.T) {
	t.Run("VAPLA-B01: release_pr pushes anchors/release/<product> with the bump, the release's blocks and the intents consumed, and opens its pull request into the branch", func(t *testing.T) {
		root := applyRepo(t, "release_pr")
		pr := stubPR(t)
		out, err := applyRun(t, root)
		if err != nil || !strings.Contains(out, "✓ opened: release app 1.1.0 (anchors/release/app → review)") {
			t.Fatalf("%v\n%s", err, out)
		}
		if currentBranch(t, root) != "review" {
			t.Errorf("back on review: %s", currentBranch(t, root))
		}
		show := func(f string) string {
			b, err := exec.Command("git", "-C", root, "show", "origin/anchors/release/app:"+f).Output()
			if err != nil {
				gitIn(t, root, "fetch", "-q", "origin")
				b, _ = exec.Command("git", "-C", root, "show", "origin/anchors/release/app:"+f).Output()
			}
			return string(b)
		}
		if got := show("app/version.json"); got != `{"version": "1.1.0"}` {
			t.Errorf("the bump: %q", got)
		}
		if got := show("CHANGELOG.md"); !strings.Contains(got, changelog.Marker("1.1.0")) || !strings.Contains(got, "sharing") {
			t.Errorf("CHANGELOG.md: %s", got)
		}
		if got := show("NOTES.md"); !strings.Contains(got, "## 1.1.0 — ") || !strings.Contains(got, "- Receipts can be shared.") {
			t.Errorf("NOTES.md: %s", got)
		}
		if show(".anchors/changes/share.md") != "" || !strings.Contains(show(".anchors/changes/released/1.1.0/share.md"), "Receipts") {
			t.Errorf("the intent is archived")
		}
		subject, _ := exec.Command("git", "-C", root, "log", "-1", "--format=%s", "origin/anchors/release/app").Output()
		if strings.TrimSpace(string(subject)) != "chore(release): app 1.1.0" {
			t.Errorf("commit: %s", subject)
		}
		if last := pr.calls[len(pr.calls)-1]; !strings.HasPrefix(last, "pr create --base review --head anchors/release/app --title release app 1.1.0 --body ## 1.1.0") {
			t.Errorf("pr: %q", pr.calls)
		}
	})
	t.Run("VAPLA-B02: A rerun updates the open pull request, and once the release is on the branch there is nothing to propose", func(t *testing.T) {
		root := applyRepo(t, "release_pr")
		pr := stubPR(t)
		if _, err := applyRun(t, root); err != nil {
			t.Fatal(err)
		}
		pr.open = `[{"number": 12}]`
		commitFiles(t, root, "fix: a crash", nil)
		if out, err := applyRun(t, root); err != nil || !strings.Contains(out, "✓ #12 updated: release app 1.1.0") || !strings.HasPrefix(pr.calls[len(pr.calls)-1], "pr edit 12 --title release app 1.1.0") {
			t.Errorf("update: %v %q\n%s", err, pr.calls, out)
		}
		gitIn(t, root, "fetch", "-q", "origin")
		gitIn(t, root, "merge", "-q", "--no-edit", "origin/anchors/release/app")
		calls := len(pr.calls)
		if out, err := applyRun(t, root); err != nil || !strings.Contains(out, "already carries the release, nothing to propose") || len(pr.calls) != calls {
			t.Errorf("merged: %v\n%s", err, out)
		}
	})
	t.Run("VAPLA-B03: Nothing to release opens nothing", func(t *testing.T) {
		root := releasedRepo(t)
		stubPR(t)
		if out, err := applyRun(t, root); err != nil || !strings.Contains(out, "nothing to release since app-v1.0.0") {
			t.Errorf("%v\n%s", err, out)
		}
	})
}

func TestApply_commitAndBuild(t *testing.T) {
	t.Run("VAPLA-B04: commit confirms the sources hold the version the branch asks for, and refuses otherwise", func(t *testing.T) {
		root := applyRepo(t, "commit")
		if _, err := applyRun(t, root); err == nil || !strings.Contains(err.Error(), "the sources hold 1.0.0 and review asks for 1.1.0") {
			t.Errorf("not bumped: %v", err)
		}
		writeIn(t, root, "app/version.json", `{"version": "1.1.0"}`)
		if out, err := applyRun(t, root); err != nil || !strings.Contains(out, "✓ app is at 1.1.0, as review asks") {
			t.Errorf("bumped: %v\n%s", err, out)
		}
	})
	t.Run("VAPLA-B05: build writes the version into the workspace's sources, commits nothing, and prints version and tag, into $GITHUB_OUTPUT too", func(t *testing.T) {
		root := applyRepo(t, "build")
		ghOut := filepath.Join(t.TempDir(), "out")
		t.Setenv("GITHUB_OUTPUT", ghOut)
		before := headOf(t, root)
		out, err := applyRun(t, root)
		if err != nil || !strings.Contains(out, "version=1.1.0\ntag=app-v1.1.0\n") || readFile(t, root, "app/version.json") != `{"version": "1.1.0"}` || headOf(t, root) != before {
			t.Errorf("%v\n%s", err, out)
		}
		if b, _ := os.ReadFile(ghOut); string(b) != "version=1.1.0\ntag=app-v1.1.0\n" {
			t.Errorf("GITHUB_OUTPUT: %q", b)
		}
	})
}

func TestApply_refused(t *testing.T) {
	t.Run("VAPLA-E01: release_pr of a product with no source, a branch that never versions, or no gh fails naming why", func(t *testing.T) {
		root := applyRepo(t, "release_pr")
		cfg := readFile(t, root, "anchors.yaml")
		if _, err := versionRun(t, "--root", root, "--branch", "feat/x", "apply"); err == nil || !strings.Contains(err.Error(), "never versions") {
			t.Errorf("none: %v", err)
		}
		noSource := strings.Replace(cfg, "    source: [{ file: app/version.json, path: version }]\n", "", 1)
		writeIn(t, root, "anchors.yaml", strings.Replace(noSource, "      \"maint/*\":\n        next: manual\n", "", 1))
		if _, err := applyRun(t, root); err == nil || !strings.Contains(err.Error(), "uses apply: build") {
			t.Errorf("no source: %v", err)
		}
		writeIn(t, root, "anchors.yaml", cfg)
		old := prAPI
		prAPI = func(string, ...string) ([]byte, error) { return nil, errors.New("the gh CLI is not in the PATH") }
		t.Cleanup(func() { prAPI = old })
		if _, err := applyRun(t, root); err == nil || !strings.Contains(err.Error(), "gh CLI") || currentBranch(t, root) != "review" {
			t.Errorf("no gh: %v, on %s", err, currentBranch(t, root))
		}
		prAPI = stubbedPR
		t.Setenv("PATH", t.TempDir())
		if _, err := prAPI(root, "pr", "list"); err == nil || !strings.Contains(err.Error(), "gh CLI is not in the PATH") {
			t.Errorf("real prAPI without gh: %v", err)
		}
	})
}

// stubbedPR is the prAPI the package was built with.
var stubbedPR = prAPI

func TestApply_perVersionAndDelete(t *testing.T) {
	t.Run("VAPLA-B06: A per_version output gets the release's file, named by format.file, in place of the unreleased one, and intents consumed by delete are gone", func(t *testing.T) {
		root := applyRepo(t, "release_pr")
		cfg := readFile(t, root, "anchors.yaml")
		cfg = strings.Replace(cfg, "        - path: CHANGELOG.md\n", "        - { path: notes, mode: per_version, format: { file: \"v{version}.md\" } }\n", 1)
		cfg = strings.Replace(cfg, "consume: archive", "consume: delete", 1)
		commitFiles(t, root, "chore: per version", map[string]string{"anchors.yaml": cfg, "notes/unreleased.md": "next\n"})
		stubPR(t)
		if _, err := applyRun(t, root); err != nil {
			t.Fatal(err)
		}
		gitIn(t, root, "fetch", "-q", "origin")
		files, _ := exec.Command("git", "-C", root, "ls-tree", "-r", "--name-only", "origin/anchors/release/app").Output()
		list := string(files)
		if !strings.Contains(list, "notes/v1.1.0.md\n") || strings.Contains(list, "notes/unreleased.md") || strings.Contains(list, ".anchors/changes/") {
			t.Errorf("%s", list)
		}
	})
}
