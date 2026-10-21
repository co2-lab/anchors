// @anchors
//   code: TLCCH
//   ref: GTENG

package gate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- a project's tool is not run again over the same inputs (DESIGN-tool-result-cache.md) ---
//
// A gate that declares what its tool reads besides its targets (`cache.inputs`) keeps the
// tool's verdict under a key of everything that verdict depends on — the gate as declared,
// each target's content, each input's content, the version of Anchors —, and replays it when
// the key comes again. A pass is kept; a failure only where the gate opts in; in a CI only
// where it opts in. With no inputs, the tool runs every time.

var noToolCache atomic.Bool

// SetNoToolCache makes every tool run, whatever it keeps (`check --no-cache`).
func SetNoToolCache(on bool) { noToolCache.Store(on) }

// toolCacheDir is where the kept verdicts live, out of git.
func toolCacheDir(root string) string { return filepath.Join(root, ".anchors", "cache", "tools") }

// toolCacheKeep is how many keys a gate keeps: a branch switched back and forth finds both.
const toolCacheKeep = 5

type toolCacheEntry struct {
	Key     string  `json:"key"`
	Verdict Verdict `json:"verdict"`
	Detail  string  `json:"detail"`
	At      string  `json:"at"`
}

var toolCacheMu sync.Mutex

// usesToolCache says whether this run keeps and replays the gate's tool verdicts.
func usesToolCache(g config.Gate) bool {
	if g.Cache == nil || len(g.Cache.Inputs) == 0 || noToolCache.Load() {
		return false
	}
	return os.Getenv("CI") == "" || g.Cache.CI
}

// cachedTool is run — the gate's tool over its targets — answered from the kept verdicts when
// everything it depends on is as it was.
func cachedTool(g config.Gate, scope string, targets []string, root string, run func() (Verdict, string, bool)) (Verdict, string) {
	if !usesToolCache(g) {
		v, detail, _ := run()
		return v, detail
	}
	key, ok := toolKey(g, scope, targets, root)
	if !ok {
		v, detail, _ := run()
		return v, detail
	}
	path := filepath.Join(toolCacheDir(root), safeName(g.Name)+".json")
	toolCacheMu.Lock()
	entries := readToolCache(path)
	toolCacheMu.Unlock()
	for _, e := range entries {
		if e.Key == key {
			return e.Verdict, strings.TrimSpace(e.Detail + "\n" + i18n.T("gate.tool_cached", e.At))
		}
	}
	v, detail, ran := run()
	keep := ran && (v == Pass || (v == Fail && g.Cache.Failures))
	if !keep {
		return v, detail
	}
	toolCacheMu.Lock()
	defer toolCacheMu.Unlock()
	entries = append([]toolCacheEntry{{Key: key, Verdict: v, Detail: detail, At: time.Now().Format(time.RFC3339)}}, readToolCache(path)...)
	if len(entries) > toolCacheKeep {
		entries = entries[:toolCacheKeep]
	}
	if b, err := json.MarshalIndent(entries, "", "  "); err == nil {
		if os.MkdirAll(filepath.Dir(path), 0o755) == nil {
			_ = os.WriteFile(path, b, 0o644) // @resilient: a cache that cannot be written only runs the tool again next time
		}
	}
	return v, detail
}

func readToolCache(path string) []toolCacheEntry {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var out []toolCacheEntry
	if json.Unmarshal(b, &out) != nil {
		return nil
	}
	return out
}

// toolKey is the hash of everything the tool's verdict depends on; false when one of them
// cannot be read — a target or an input — and the tool must run.
func toolKey(g config.Gate, scope string, targets []string, root string) (string, bool) {
	h := sha256.New()
	write := func(parts ...string) {
		for _, p := range parts {
			h.Write([]byte(p))
			h.Write([]byte{0})
		}
	}
	write("anchors", mapx.GeneratedBy, "gate", g.Name, g.Run, scope)
	sorted := append([]string(nil), targets...)
	sort.Strings(sorted)
	for _, t := range sorted {
		b, err := readFile(root, t)
		if err != nil {
			return "", false
		}
		sum := sha256.Sum256(b)
		write("target", t, hex.EncodeToString(sum[:]))
	}
	files, known := projectFiles(root)
	if !known {
		return "", false
	}
	for _, f := range files {
		for _, glob := range g.Cache.Inputs {
			if ok, _ := doublestar.Match(glob, f); ok {
				b, err := readFile(root, f)
				if err != nil {
					return "", false
				}
				sum := sha256.Sum256(b)
				write("input", f, hex.EncodeToString(sum[:]))
				break
			}
		}
	}
	return hex.EncodeToString(h.Sum(nil)), true
}

// projectFiles are the project's files as git sees them, listed once per run.
func projectFiles(root string) ([]string, bool) {
	s := gitState(root)
	s.filesOnce.Do(func() {
		s.files, s.filesKnown = gitmeta.ProjectFiles(root)
		sort.Strings(s.files)
	})
	return s.files, s.filesKnown
}

func safeName(name string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
}
