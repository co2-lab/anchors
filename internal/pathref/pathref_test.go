// @anchors
//   code: PTRFT
//   ref: PTRFS

package pathref

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// serve answers fetch from a map of URL → body, a missing URL with a 404 and every URL with
// an unreachable server when down; it clears the per-process memo and names the branch.
func serve(t *testing.T, branch string, pages map[string]string, down *bool) {
	t.Helper()
	oldFetch, oldBranch := fetch, currentBranch
	t.Cleanup(func() { fetch, currentBranch = oldFetch, oldBranch })
	memoMu.Lock()
	memo = map[string]fetched{}
	memoMu.Unlock()
	currentBranch = func(string) string { return branch }
	fetch = func(u string) ([]byte, error) {
		if down != nil && *down {
			return nil, errors.New("connection refused")
		}
		if body, ok := pages[u]; ok {
			return []byte(body), nil
		}
		return nil, errNotFound
	}
}

func forget() {
	memoMu.Lock()
	memo = map[string]fetched{}
	memoMu.Unlock()
}

const guide = "https://raw.example/std/{{branch|develop}}/GUIDE.md"

func TestRead_aLocalPathOrAURL(t *testing.T) {
	t.Run("PTRFS-B01: A path reads the project's file, and a URL the external file", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "GUIDE.md"), []byte("local"), 0o644); err != nil {
		t.Fatal(err)
	}
	serve(t, "", map[string]string{"https://raw.example/std/v1/GUIDE.md": "pinned"}, nil)
	if b, err := Read(root, "GUIDE.md"); err != nil || string(b) != "local" {
		t.Errorf("a local path: %q %v", b, err)
	}
	if b, err := Read(root, "https://raw.example/std/v1/GUIDE.md"); err != nil || string(b) != "pinned" {
		t.Errorf("a URL: %q %v", b, err)
	}
}

func TestRead_theBranchMirrorThenTheDefault(t *testing.T) {
	t.Run("PTRFS-B02: {{branch}} reads the current branch's mirror, and the default when the mirror does not exist", func(t *testing.T) {})
	root := t.TempDir()
	pages := map[string]string{
		"https://raw.example/std/develop/GUIDE.md":   "develop",
		"https://raw.example/std/feature/x/GUIDE.md": "feature",
	}
	serve(t, "feature/x", pages, nil)
	if u, b, err := Resolve(root, guide); err != nil || string(b) != "feature" || u != "https://raw.example/std/feature/x/GUIDE.md" {
		t.Errorf("the mirror: %s %q %v", u, b, err)
	}
	serve(t, "feature/y", pages, nil)
	if u, b, err := Resolve(root, guide); err != nil || string(b) != "develop" || u != "https://raw.example/std/develop/GUIDE.md" {
		t.Errorf("the default, the mirror missing: %s %q %v", u, b, err)
	}
	serve(t, "", pages, nil)
	if b, err := Read(root, guide); err != nil || string(b) != "develop" {
		t.Errorf("the default, no branch checked out: %q %v", b, err)
	}
}

func TestRead_offlineTheKeptCopy(t *testing.T) {
	t.Run("PTRFS-B03: With the server unreachable, the last copy kept answers", func(t *testing.T) {})
	root := t.TempDir()
	down := false
	serve(t, "feature/y", map[string]string{"https://raw.example/std/develop/GUIDE.md": "develop"}, &down)
	if _, err := Read(root, guide); err != nil {
		t.Fatal(err)
	}
	down = true
	forget()
	if b, err := Read(root, guide); err != nil || string(b) != "develop" {
		t.Errorf("offline, the default's copy (the mirror never kept): %q %v", b, err)
	}
}

func TestRead_aURLThatAnswersNothingFails(t *testing.T) {
	t.Run("PTRFS-E01: A URL that answers nothing, with or without the default, or unreachable with no copy, fails naming what was tried", func(t *testing.T) {})
	root := t.TempDir()
	serve(t, "feature/x", nil, nil)
	_, err := Read(root, guide)
	if err == nil || !strings.Contains(err.Error(), "feature/x/GUIDE.md") || !strings.Contains(err.Error(), "develop/GUIDE.md") {
		t.Errorf("neither the mirror nor the default: %v", err)
	}
	serve(t, "", nil, nil)
	if _, err := Read(root, "https://raw.example/std/{{branch}}/GUIDE.md"); err == nil || !strings.Contains(err.Error(), "declare a default") {
		t.Errorf("no branch, no default: %v", err)
	}
	down := true
	serve(t, "main", nil, &down)
	if _, err := Read(root, "https://raw.example/std/v1/GUIDE.md"); err == nil || !strings.Contains(err.Error(), "no copy of it is kept") {
		t.Errorf("unreachable, nothing kept: %v", err)
	}
}

func TestFetch_overHTTP(t *testing.T) {
	t.Run("PTRFS-B01: A path reads the project's file, and a URL the external file", func(t *testing.T) {})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			_, _ = w.Write([]byte("body"))
		case "/boom":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	if b, err := fetch(srv.URL + "/ok"); err != nil || string(b) != "body" {
		t.Errorf("a 200: %q %v", b, err)
	}
	if _, err := fetch(srv.URL + "/missing"); !errors.Is(err, errNotFound) {
		t.Errorf("a 404 is not found: %v", err)
	}
	if _, err := fetch(srv.URL + "/boom"); err == nil || errors.Is(err, errNotFound) {
		t.Errorf("a 500 is an error, not a missing file: %v", err)
	}
	if _, err := fetch("http://127.0.0.1:1/x"); err == nil {
		t.Error("an unreachable server is an error")
	}
	if !IsURL("https://a/b") || IsURL("docs/a.md") {
		t.Error("IsURL")
	}
}

func TestCurrentBranch_fromTheCheckoutOrTheCI(t *testing.T) {
	t.Run("PTRFS-B02: {{branch}} reads the current branch's mirror, and the default when the mirror does not exist", func(t *testing.T) {})
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init", "-q", "-b", "develop").CombinedOutput(); err != nil {
		t.Skipf("git: %v %s", err, out)
	}
	t.Setenv("GITHUB_HEAD_REF", "")
	t.Setenv("GITHUB_REF_TYPE", "")
	if b := currentBranch(root); b != "develop" {
		t.Errorf("the branch checked out, no commit yet: %q", b)
	}
	t.Setenv("GITHUB_HEAD_REF", "feature/pr")
	if b := currentBranch(root); b != "feature/pr" {
		t.Errorf("a pull request's head: %q", b)
	}
	t.Setenv("GITHUB_HEAD_REF", "")
	t.Setenv("GITHUB_REF_TYPE", "branch")
	t.Setenv("GITHUB_REF_NAME", "main")
	if b := currentBranch(t.TempDir()); b != "main" {
		t.Errorf("a CI's ref, outside a checkout: %q", b)
	}
}
