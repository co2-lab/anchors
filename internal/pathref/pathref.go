// @anchors
//   code: PTRFC
//   ref: PTRFS

// Package pathref reads the file a path of anchors.yaml names: a file of the project, or an
// external one by URL — a guide kept in another repository, a changelog template, a pack. It
// is the one place that fetches: no command opens a URL on its own.
//
// A URL may carry `{{branch}}`, replaced by the project's current branch, so a project on
// `develop` reads the `develop` of the external repository; `{{branch|develop}}` names the
// branch read when the mirror does not exist. A URL that answers nothing, with or without
// the default, is an error wherever it is read.
package pathref

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// IsURL says whether a path of anchors.yaml names an external file.
func IsURL(p string) bool {
	return strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "http://")
}

var branchRE = regexp.MustCompile(`\{\{\s*branch\s*(?:\|\s*([^}\s]+)\s*)?\}\}`)

// errNotFound is a URL the server answered it does not have.
var errNotFound = errors.New("not found")

// fetch reads a URL: its body, errNotFound on a 404, another error when the server could not
// be reached or answered otherwise. A variable so the tests answer without a network.
var fetch = func(u string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	// A private repository's raw file asks for the token; it goes only to GitHub's hosts.
	if host := req.URL.Hostname(); host == "raw.githubusercontent.com" || host == "api.github.com" {
		if tok := firstEnv("GH_TOKEN", "GITHUB_TOKEN"); tok != "" {
			req.Header.Set("Authorization", "token "+tok)
		}
	}
	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return nil, errNotFound
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("answered %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// currentBranch is the project's branch: the pull request's head in a CI, else the branch
// checked out, else the ref a CI ran on. Empty when none tells it (a detached checkout).
var currentBranch = func(root string) string {
	if b := os.Getenv("GITHUB_HEAD_REF"); b != "" {
		return b
	}
	// symbolic-ref, not rev-parse: it names the branch of a repository with no commit yet,
	// and fails on a detached checkout instead of answering "HEAD".
	if out, err := exec.Command("git", "-C", root, "symbolic-ref", "--short", "-q", "HEAD").Output(); err == nil {
		if b := strings.TrimSpace(string(out)); b != "" {
			return b
		}
	}
	if os.Getenv("GITHUB_REF_TYPE") == "branch" {
		return os.Getenv("GITHUB_REF_NAME")
	}
	return ""
}

// Candidates are the URLs a path reads, in order: with `{{branch}}`, the current branch's,
// then the default's when one is declared; without it, the URL itself. An error when the
// URL asks for the branch, no branch is known and no default is declared.
func Candidates(root, u string) ([]string, error) {
	m := branchRE.FindStringSubmatch(u)
	if m == nil {
		return []string{u}, nil
	}
	var out []string
	if b := currentBranch(root); b != "" {
		out = append(out, branchRE.ReplaceAllLiteralString(u, b))
	}
	if def := m[1]; def != "" {
		if d := branchRE.ReplaceAllLiteralString(u, def); len(out) == 0 || out[0] != d {
			out = append(out, d)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s asks for the branch, and none is checked out: declare a default, `{{branch|develop}}`", u)
	}
	return out, nil
}

type fetched struct {
	url  string
	body []byte
	err  error
}

var (
	memoMu sync.Mutex
	memo   = map[string]fetched{}
)

// Read reads the file p names: a path relative to root, or a URL. A URL is read once per
// process and kept in `.anchors/cache/`; when the server cannot be reached the last copy kept
// answers, and when the server answers it has none — the branch's nor the default's — Read
// fails, naming each URL tried.
func Read(root, p string) ([]byte, error) {
	if !IsURL(p) {
		return os.ReadFile(filepath.Join(root, filepath.FromSlash(p)))
	}
	_, body, err := Resolve(root, p)
	return body, err
}

// Resolve reads a URL as Read does and says which URL answered — the branch's or the
// default's —, the one a reader is pointed at.
func Resolve(root, u string) (string, []byte, error) {
	memoMu.Lock()
	defer memoMu.Unlock()
	key := root + "\x00" + u
	if f, ok := memo[key]; ok {
		return f.url, f.body, f.err
	}
	f := resolve(root, u)
	memo[key] = f
	return f.url, f.body, f.err
}

func resolve(root, u string) fetched {
	candidates, err := Candidates(root, u)
	if err != nil {
		return fetched{err: err}
	}
	// unreachable is the first URL the server could not be asked about and of which no copy
	// is kept: the next candidate — the default — may still answer, from the server or from
	// its copy, and only when none does is it the error.
	var unreachable error
	for _, c := range candidates {
		if _, err := url.Parse(c); err != nil {
			return fetched{err: fmt.Errorf("%s is not a URL: %w", c, err)}
		}
		body, err := fetch(c)
		switch {
		case err == nil:
			keep(root, c, body)
			return fetched{url: c, body: body}
		case errors.Is(err, errNotFound):
			continue
		}
		if kept, ok := kept(root, c); ok {
			return fetched{url: c, body: kept}
		}
		if unreachable == nil {
			unreachable = fmt.Errorf("%s could not be read (%v), and no copy of it is kept", c, err)
		}
	}
	if unreachable != nil {
		return fetched{err: unreachable}
	}
	return fetched{err: fmt.Errorf("%s does not exist: tried %s", u, strings.Join(candidates, ", "))}
}

func cachePath(root, u string) string {
	sum := sha256.Sum256([]byte(u))
	return filepath.Join(root, ".anchors", "cache", "urls", hex.EncodeToString(sum[:])[:24])
}

// keep writes the copy best-effort: a cache that cannot be written only costs the next
// offline read.
func keep(root, u string, body []byte) {
	p := cachePath(root, u)
	if os.MkdirAll(filepath.Dir(p), 0o755) == nil {
		_ = os.WriteFile(p, body, 0o644)
	}
}

func kept(root, u string) ([]byte, bool) {
	b, err := os.ReadFile(cachePath(root, u))
	return b, err == nil
}

func firstEnv(names ...string) string {
	for _, n := range names {
		if v := os.Getenv(n); v != "" {
			return v
		}
	}
	return ""
}
