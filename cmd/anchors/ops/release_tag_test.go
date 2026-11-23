// @anchors
//   code: RLSTC
//   ref: RLSTA

package ops

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

// fakeGitHub answers the release's API calls: a GET of a known path returns its body, of an
// unknown one a 404; a POST is kept, with its fields.
type fakeGitHub struct {
	have   map[string]string
	posts  []string
	fail   error
	failOn map[string]error // a path, or a POST's path, → its error
}

func stubReleaseAPI(t *testing.T) *fakeGitHub {
	t.Helper()
	f := &fakeGitHub{have: map[string]string{}, failOn: map[string]error{}}
	old := releaseAPI
	releaseAPI = func(_ string, args ...string) ([]byte, error) {
		if f.fail != nil {
			return nil, f.fail
		}
		path := args[0]
		if args[0] == "-X" {
			path = args[2]
		}
		if err := f.failOn[path]; err != nil {
			return nil, err
		}
		if args[0] == "-X" {
			f.posts = append(f.posts, strings.Join(args[2:], " "))
			return []byte("{}"), nil
		}
		if body, ok := f.have[args[0]]; ok {
			return []byte(body), nil
		}
		return nil, errors.New("gh api " + args[0] + ": exit status 1: Not Found (HTTP 404)")
	}
	t.Cleanup(func() { releaseAPI = old })
	return f
}

// releasingRepo is releasedRepo whose production stage releases app, at 1.0.1 on production
// after a fix.
func releasingRepo(t *testing.T) (string, string) {
	t.Helper()
	root := releasedRepo(t)
	gitIn(t, root, "checkout", "-q", "production")
	cfg := strings.Replace(releasedProject, "    production-release:\n      branches: [production]\n", "    production-release:\n      branches: [production]\n      release: app\n", 1)
	commitFiles(t, root, "chore: release app", map[string]string{"anchors.yaml": cfg})
	commitFiles(t, root, "fix: a crash\n\nBug: seen in production", nil)
	sha := commitFiles(t, root, "chore: 1.0.1", appVersion("1.0.1"))
	return root, sha
}

func TestRelease_writesTheTag(t *testing.T) {
	t.Run("RLSTA-B01: A stage that releases writes, when it passes, its product's tag on the sha and the tag's Release with the changelog since the release before", func(t *testing.T) {
		root, sha := releasingRepo(t)
		f := fakeForge(t)
		gh := stubReleaseAPI(t)
		out, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha)
		if err != nil || !strings.Contains(out, "app 1.0.1 released as app-v1.0.1") {
			t.Fatalf("%v\n%s", err, out)
		}
		if len(gh.posts) != 2 || gh.posts[0] != "repos/{owner}/{repo}/git/refs -f ref=refs/tags/app-v1.0.1 -f sha="+sha {
			t.Fatalf("the tag: %q", gh.posts)
		}
		rel := gh.posts[1]
		if !strings.Contains(rel, "tag_name=app-v1.0.1") || !strings.Contains(rel, "## 1.0.1 — ") || !strings.Contains(rel, "a crash") || !strings.Contains(rel, "seen in production") || strings.Contains(rel, "release 1.0.0") {
			t.Errorf("the release's body is its changelog since app-v1.0.0: %s", rel)
		}
		if f.statuses[sha+"|anchors/pipeline/production-release"] != "success" {
			t.Errorf("recorded: %v", f.statuses)
		}
	})
	t.Run("RLSTA-B02: A tag already on the sha is a rerun: only a missing Release is written", func(t *testing.T) {
		root, sha := releasingRepo(t)
		fakeForge(t)
		gh := stubReleaseAPI(t)
		gh.have["repos/{owner}/{repo}/git/ref/tags/app-v1.0.1"] = `{"object": {"sha": "` + sha + `", "type": "commit"}}`
		gh.have["repos/{owner}/{repo}/releases/tags/app-v1.0.1"] = `{}`
		if out, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha); err != nil || len(gh.posts) != 0 || !strings.Contains(out, "is released as app-v1.0.1") {
			t.Errorf("tag and release there: %v %q\n%s", err, gh.posts, out)
		}
		delete(gh.have, "repos/{owner}/{repo}/releases/tags/app-v1.0.1")
		if _, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha); err != nil || len(gh.posts) != 1 || !strings.Contains(gh.posts[0], "releases -f tag_name=app-v1.0.1") {
			t.Errorf("release missing: %v %q", err, gh.posts)
		}
	})
}

func TestRelease_refused(t *testing.T) {
	t.Run("RLSTA-E01: A tag on another sha, a version the sha cannot give, or GitHub not answering fails the record, recorded as failure", func(t *testing.T) {
		root, sha := releasingRepo(t)
		f := fakeForge(t)
		gh := stubReleaseAPI(t)
		gh.have["repos/{owner}/{repo}/git/ref/tags/app-v1.0.1"] = `{"object": {"sha": "0123456789abcdef", "type": "commit"}}`
		_, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha)
		if err == nil || !strings.Contains(err.Error(), "the tag app-v1.0.1 is already on 0123456") || len(gh.posts) != 0 {
			t.Errorf("released twice: %v %q", err, gh.posts)
		}
		if f.statuses[sha+"|anchors/pipeline/production-release"] != "failure" {
			t.Errorf("recorded as failure: %v", f.statuses)
		}
		delete(gh.have, "repos/{owner}/{repo}/git/ref/tags/app-v1.0.1")
		gh.fail = errors.New("gh api: exit status 1: HTTP 502")
		if _, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha); err == nil || !strings.Contains(err.Error(), "read the tag app-v1.0.1") {
			t.Errorf("no answer: %v", err)
		}
		gh.fail = nil
		broken := commitFiles(t, root, "chore: broken", map[string]string{"app/version.json": `{}`})
		if _, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", broken); err == nil || !strings.Contains(err.Error(), "app/version.json") {
			t.Errorf("no version: %v", err)
		}
	})
}

func TestRelease_GitHubFails(t *testing.T) {
	t.Run("RLSTA-E02: A tag or a Release GitHub does not write, a ref it answers unread, a release it does not say, or no gh, fail the record naming the step", func(t *testing.T) {
		root, sha := releasingRepo(t)
		fakeForge(t)
		gh := stubReleaseAPI(t)
		boom := errors.New("gh api: exit status 1: HTTP 500")
		for _, c := range []struct {
			setup func()
			says  string
		}{
			{func() { gh.failOn["repos/{owner}/{repo}/git/refs"] = boom }, "write the tag app-v1.0.1"},
			{func() { gh.failOn["repos/{owner}/{repo}/releases/tags/app-v1.0.1"] = boom }, "read the release app-v1.0.1"},
			{func() { gh.failOn["repos/{owner}/{repo}/releases"] = boom }, "write the release app-v1.0.1"},
			{func() { gh.have["repos/{owner}/{repo}/git/ref/tags/app-v1.0.1"] = "not json" }, "read the tag app-v1.0.1"},
		} {
			gh.failOn, gh.have, gh.posts = map[string]error{}, map[string]string{}, nil
			c.setup()
			if _, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha); err == nil || !strings.Contains(err.Error(), c.says) {
				t.Errorf("want %q: %v", c.says, err)
			}
		}
		releaseAPI = stubbed
		t.Setenv("PATH", t.TempDir())
		if _, err := releaseAPI(root, "repos/x"); err == nil || !strings.Contains(err.Error(), "gh CLI is not in the PATH") {
			t.Errorf("no gh: %v", err)
		}
	})
	t.Run("RLSTA-E03: A tag of the product the scheme cannot read stops the release body", func(t *testing.T) {
		root, sha := releasingRepo(t)
		gitIn(t, root, "tag", "app-vnext")
		fakeForge(t)
		stubReleaseAPI(t)
		if _, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha); err == nil || !strings.Contains(err.Error(), "app-vnext") {
			t.Errorf("%v", err)
		}
	})
}

// stubbed is the releaseAPI the package was built with.
var stubbed = releaseAPI

func TestRelease_label(t *testing.T) {
	t.Run("RLSTA-B03: A stage on a branch that labels its tag releases the next prerelease GitHub does not hold, and a rerun the one already on its sha", func(t *testing.T) {
		root, sha := releasingRepo(t)
		cfg := readFile(t, root, "anchors.yaml")
		cfg = strings.Replace(cfg, "      \"maint/*\":\n        next: manual\n", "      \"maint/*\":\n        next: manual\n      production:\n        label: { name: rc }\n", 1)
		commitFiles(t, root, "chore: label", map[string]string{"anchors.yaml": cfg})
		sha = headOf(t, root)
		fakeForge(t)
		gh := stubReleaseAPI(t)
		gh.have["repos/{owner}/{repo}/git/matching-refs/tags/app-v1.0.1-rc."] = `[{"ref": "refs/tags/app-v1.0.1-rc.1", "object": {"sha": "0000000"}}]`
		if out, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha); err != nil || !strings.Contains(out, "released as app-v1.0.1-rc.2") || !strings.Contains(gh.posts[0], "ref=refs/tags/app-v1.0.1-rc.2") {
			t.Errorf("%v %q\n%s", err, gh.posts, out)
		}
		gh.posts = nil
		gh.have["repos/{owner}/{repo}/git/matching-refs/tags/app-v1.0.1-rc."] = `[{"ref": "refs/tags/app-v1.0.1-rc.1", "object": {"sha": "0000000"}}, {"ref": "refs/tags/app-v1.0.1-rc.2", "object": {"sha": "` + sha + `"}}]`
		gh.have["repos/{owner}/{repo}/git/ref/tags/app-v1.0.1-rc.2"] = `{"object": {"sha": "` + sha + `"}}`
		gh.have["repos/{owner}/{repo}/releases/tags/app-v1.0.1-rc.2"] = `{}`
		if out, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha); err != nil || len(gh.posts) != 0 || !strings.Contains(out, "is released as app-v1.0.1-rc.2") {
			t.Errorf("rerun: %v %q\n%s", err, gh.posts, out)
		}
	})
}

// Reported from a project: its first tagged release followed no tag, the changelog was the
// whole history, and GitHub refused a body past 125000 characters.
func TestRelease_bodyPastGitHubsLimit(t *testing.T) {
	t.Run("RLSTA-B04: A changelog longer than GitHub keeps is cut at the last whole line that fits, saying how many lines were left out", func(t *testing.T) {
		root, _ := releasingRepo(t)
		for i := 0; i < 30; i++ {
			commitFiles(t, root, "fix: crash number "+strconv.Itoa(i), nil)
		}
		sha := commitFiles(t, root, "chore: 1.0.2", appVersion("1.0.2"))
		fakeForge(t)
		gh := stubReleaseAPI(t)
		old := releaseBodyLimit
		releaseBodyLimit = 600
		t.Cleanup(func() { releaseBodyLimit = old })
		if out, err := runPipeline(t, root, "record", "production-release", "--branch", "production", "--sha", sha); err != nil || len(gh.posts) != 2 {
			t.Fatalf("%v %q\n%s", err, gh.posts, out)
		}
		_, body, _ := strings.Cut(gh.posts[1], "-f body=")
		if utf8.RuneCountInString(body) > releaseBodyLimit || !strings.HasPrefix(body, "## 1.0.2") || !strings.Contains(body, "crash number") || !strings.Contains(body, "more lines of the changelog left out") {
			t.Errorf("cut to the limit, saying so (%d characters):\n%s", utf8.RuneCountInString(body), body)
		}
	})
	t.Run("RLSTA-B04: A body within the limit is kept whole", func(t *testing.T) {
		if got := fitRelease("## 1\n- a\n", 100); got != "## 1\n- a\n" {
			t.Errorf("%q", got)
		}
		body := "## 1\n" + strings.Repeat("- a commit\n", 30)
		got := fitRelease(body, 200)
		if utf8.RuneCountInString(got) > 200 || !strings.HasPrefix(got, "## 1\n- a commit\n") {
			t.Fatalf("whole lines, within the limit: %q", got)
		}
		kept := strings.Count(got, "- a commit\n")
		if !strings.Contains(got, "… "+strconv.Itoa(30-kept)+" more lines") {
			t.Errorf("the lines left out are counted (%d kept): %q", kept, got)
		}
	})
}
