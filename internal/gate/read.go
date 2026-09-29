package gate

import (
	"os"
	"path/filepath"
	"sync"
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
)

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
