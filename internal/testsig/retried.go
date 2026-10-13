// @anchors
//   code: RTRCR
//   ref: TFSTM

package testsig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// RetriedFile keeps the tests that passed only on a retry under load, so a test that keeps
// running out of time stays in sight instead of hiding behind the retry
// (DESIGN-retry-on-load.md).
const RetriedFile = ".anchors/retried.json"

// retriedKeep bounds the record: the latest entries are what tells a chronic one.
const retriedKeep = 500

// Retried is one test that passed on a retry: which, when, and under what load.
type Retried struct {
	Name string    `json:"name"`
	File string    `json:"file,omitempty"`
	At   time.Time `json:"at"`
	Load float64   `json:"load"`
	CPUs int       `json:"cpus"`
}

// LoadRetried reads the record; a missing or unreadable one is empty.
func LoadRetried(root string) []Retried {
	var out []Retried
	if b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(RetriedFile))); err == nil {
		_ = json.Unmarshal(b, &out)
	}
	return out
}

// AddRetried appends entries to the record, keeping the latest.
func AddRetried(root string, entries []Retried) error {
	if len(entries) == 0 {
		return nil
	}
	all := append(LoadRetried(root), entries...)
	if len(all) > retriedKeep {
		all = all[len(all)-retriedKeep:]
	}
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(root, filepath.FromSlash(RetriedFile))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
