// @anchors
//   code: RDGTR
//   ref: GTENG

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/co2-lab/anchors/internal/gitmeta"
)

// WHERE THE GATES READ FILES FROM.
//
// A check reads the tree. The commit hook's check must read what the commit records — the
// index —, or another session's unstaged edit to a spec the commit does not carry bars it:
// scenario-coverage judged the spec by its tree content, with rules nobody proved yet. The
// hook sets the source to the index (`SetFileSource`); every read of a project file by a
// gate goes through `readFile`.

var (
	fileSourceMu sync.RWMutex
	fileSource   func(rel string) ([]byte, error)
	// changedSource says whether a file has changes not yet committed; nil asks the tree.
	changedSource func(root, rel string) (changed, known bool)
)

// SetChangedSource makes the gates ask whether a file changed through `changed` — the
// staged changes, in the commit hook — and returns what restores the previous one. Nil
// asks the tree.
func SetChangedSource(changed func(root, rel string) (bool, bool)) (restore func()) {
	fileSourceMu.Lock()
	prev := changedSource
	changedSource = changed
	fileSourceMu.Unlock()
	return func() {
		fileSourceMu.Lock()
		changedSource = prev
		fileSourceMu.Unlock()
	}
}

// uncommitted says whether the file has changes the last commit does not have, as the
// current source sees them.
func uncommitted(root, rel string) (changed, known bool) {
	fileSourceMu.RLock()
	src := changedSource
	fileSourceMu.RUnlock()
	if src != nil {
		return src(root, rel)
	}
	if dirty, known := gitState(root).dirtyFiles(); known {
		return dirty[filepath.ToSlash(rel)], true
	}
	return gitmeta.UncommittedChanges(root, rel)
}

// GIT, ONCE PER RUN. Asking git about each file — its status, its last commit — cost two
// processes per file: on MIF's 5,003 files, `updated-at-atual` alone took two minutes of a
// 14-minute check. One `git status` and one `git log` answer for every file. The answers
// are a run's: they expire after gitStateTTL, and ForgetRunState drops them when a run
// changes files (`check --fix`).

const gitStateTTL = 2 * time.Minute

type runGit struct {
	root      string
	at        time.Time
	dirtyOnce sync.Once
	dirty     map[string]bool
	known     bool
	datesOnce sync.Once
	dates     map[string]string
}

var (
	gitStateMu sync.Mutex
	gitStates  = map[string]*runGit{}
)

func gitState(root string) *runGit {
	gitStateMu.Lock()
	defer gitStateMu.Unlock()
	if s, ok := gitStates[root]; ok && time.Since(s.at) < gitStateTTL {
		return s
	}
	s := &runGit{root: root, at: time.Now()}
	gitStates[root] = s
	return s
}

func (s *runGit) dirtyFiles() (map[string]bool, bool) {
	s.dirtyOnce.Do(func() { s.dirty, s.known = gitmeta.DirtyFiles(s.root) })
	return s.dirty, s.known
}

// ForgetRunState drops what the gates remember of git for this run: a run that wrote files
// asks again.
func ForgetRunState() {
	gitStateMu.Lock()
	gitStates = map[string]*runGit{}
	gitStateMu.Unlock()
}

// lastCommitDate is the day of the last commit that touched the file, from one `git log`
// of the whole project; false when no commit did.
func lastCommitDate(root, rel string) (string, bool) {
	s := gitState(root)
	s.datesOnce.Do(func() { s.dates = gitmeta.AllCommitDatesUnder(root) })
	d, ok := s.dates[filepath.ToSlash(rel)]
	return d, ok
}

// SetFileSource makes the gates read project files through `read` (a path relative to the
// root), and returns what restores the previous source. Nil reads the tree.
func SetFileSource(read func(rel string) ([]byte, error)) (restore func()) {
	fileSourceMu.Lock()
	prev := fileSource
	fileSource = read
	fileSourceMu.Unlock()
	return func() {
		fileSourceMu.Lock()
		fileSource = prev
		fileSourceMu.Unlock()
	}
}

// readFile reads a project file, `rel` relative to the root, from the current source.
func readFile(root, rel string) ([]byte, error) {
	fileSourceMu.RLock()
	src := fileSource
	fileSourceMu.RUnlock()
	if src != nil {
		if r, err := filepath.Rel(root, filepath.Join(root, rel)); err == nil {
			return src(r)
		}
	}
	return os.ReadFile(filepath.Join(root, rel))
}

// indexedTargets hands an external tool the content the current source holds — the index,
// under `--index` — for the targets whose file on disk differs from it: each is copied to a
// temporary folder, at its own relative path, and the tool gets the copy. `clean` turns the
// copies' paths back into the project's in what the tool printed; `done` removes them. With
// no source set, or no target that differs, the targets are handed over as they are.
//
// A tool reads files by path, and the path is the tree's: another session's unstaged edit
// to a file in the commit's impact was spell-checked instead of what the commit records.
func indexedTargets(root string, targets []string) (out []string, clean func(string) string, done func()) {
	clean, done = func(s string) string { return s }, func() {}
	fileSourceMu.RLock()
	src := fileSource
	fileSourceMu.RUnlock()
	if src == nil || len(targets) == 0 {
		return targets, clean, done
	}
	var dir string
	out = make([]string, 0, len(targets))
	for _, t := range targets {
		want, err := src(filepath.ToSlash(t))
		if err != nil {
			continue // not in the index: not the commit's, and not judged
		}
		have, herr := os.ReadFile(filepath.Join(root, t))
		if herr == nil && string(have) == string(want) {
			out = append(out, t)
			continue
		}
		if dir == "" {
			if dir, err = os.MkdirTemp("", "anchors-index-"); err != nil {
				out = append(out, t) // @resilient: without a place for the copy, the tree's file is the best there is
				dir = ""
				continue
			}
		}
		p := filepath.Join(dir, filepath.FromSlash(t))
		if os.MkdirAll(filepath.Dir(p), 0o755) != nil || os.WriteFile(p, want, 0o644) != nil {
			out = append(out, t) // @resilient: a copy that cannot be written leaves the tree's file
			continue
		}
		out = append(out, filepath.ToSlash(p)) // with `/`, the tool prints the project's path back once the prefix is stripped
	}
	if dir == "" {
		return out, clean, done
	}
	prefixes := []string{dir}
	if real, err := filepath.EvalSymlinks(dir); err == nil && real != dir {
		prefixes = append(prefixes, real) // a tool may print the resolved path (/private/var on macOS)
	}
	clean = func(s string) string {
		for _, d := range prefixes {
			p := d + string(filepath.Separator)
			s = strings.ReplaceAll(s, p, "")
			s = strings.ReplaceAll(s, filepath.ToSlash(p), "")
		}
		return s
	}
	done = func() { _ = os.RemoveAll(dir) }
	return out, clean, done
}
