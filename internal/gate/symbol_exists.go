// @anchors
//   code: SYMXC
//   ref: SYMXS

package gate

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// symbol-exists: a name a comment or a spec marks as a symbol is in the project's code.
//
// Measured on a project (2026-11-12): a comment of the import handler read "the common reader
// reads without compiling code of the file (isEvalSupported false)". The rule code beside
// it was right at all four ends, so every gate passed — and `isEvalSupported` no longer
// existed: pdfjs-dist 5 removed the option and the reader stopped passing it. Only a review
// caught it. Prose cannot be told from a name by its shape, so the project marks the names
// it cites — `@@name`, `#name` or `$name` — and the gate confronts each against the code of
// the map: a comment's text and a spec's text with its fences removed, against the code and
// test files with their comments removed. A line carrying `@no-symbol-check: <reason>`
// cites a name on purpose (one gone since a version) and is not confronted.
func checkSymbolExists(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	var text string
	var starts map[int]bool
	switch n.Kind {
	case mapx.KindSpec:
		text = withoutFences(content)
	case mapx.KindCode:
		cs := cfg.CommentsFor(n.ID)
		if cs == nil {
			return commentsUndeclared(n.ID)
		}
		text, starts = commentText(content, cs)
	default:
		return Skip, i18n.T("gate.symbol_exists.skip_kind")
	}
	cited := citedSymbols(text, starts)
	if len(cited) == 0 {
		return Skip, i18n.T("gate.symbol_exists.skip_none")
	}
	if g == nil {
		return pendingNoMap()
	}
	known := projectSymbols(g, root, cfg)
	var missing []string
	for _, c := range cited {
		for _, part := range strings.Split(c.name, ".") {
			if !known[part] && !known["$"+part] {
				missing = append(missing, fmt.Sprintf("%s:%d `%s`", n.ID, c.line, c.name))
				break
			}
		}
	}
	if len(missing) == 0 {
		return Pass, ""
	}
	return Fail, i18n.T("gate.symbol_exists.missing", len(missing), strings.Join(missing, ", "))
}

type citedSymbol struct {
	name string
	line int
}

// symbolMarkRE reads a marked name: `@@name`, `#name` or `$name`, with dotted parts
// (`#Reader.open`). The character before the marker is checked apart (markedAt): Go's
// regexp has no look-behind.
var symbolMarkRE = regexp.MustCompile(`(@@|#|\$)([A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)*)`)

// editorFoldRE and hexColorRE are what a `#` says in a comment without marking a name: an
// editor's folding directive (`#region`, measured on a project) and a color (`#FFF`, `#C9BCFF`).
var (
	editorFoldRE = regexp.MustCompile(`^(?:end)?region$`)
	hexColorRE   = regexp.MustCompile(`^(?:[0-9A-Fa-f]{3,4}|[0-9A-Fa-f]{6}|[0-9A-Fa-f]{8})$`)
)

var noSymbolCheckRE = regexp.MustCompile(`@no-symbol-check:\s*\S`)

// citedSymbols are the names the text marks, each with its line. A marker glued to a word,
// a path or a link (`a#b`, `/x#y`, `](#section)`) is no mark, nor is a name followed by a
// dash (`#section-name`, a CSS id), a folding directive or a color; a `#` that opens a
// comment (`starts`) is the comment's delimiter.
func citedSymbols(text string, starts map[int]bool) []citedSymbol {
	var out []citedSymbol
	offset := 0
	for i, line := range strings.Split(text, "\n") {
		lineStart := offset
		offset += len(line) + 1
		if noSymbolCheckRE.MatchString(line) {
			continue
		}
		for _, m := range symbolMarkRE.FindAllStringSubmatchIndex(line, -1) {
			if !markedAt(line, m[0], m[1]) || starts[lineStart+m[0]] {
				continue
			}
			name := line[m[4]:m[5]]
			if line[m[2]:m[3]] == "#" && (editorFoldRE.MatchString(name) || hexColorRE.MatchString(name)) {
				continue
			}
			out = append(out, citedSymbol{name: name, line: i + 1})
		}
	}
	return out
}

// markedAt says whether the match [from, to) of the line is a mark and not part of a word,
// a path, a link or a longer marker.
func markedAt(line string, from, to int) bool {
	if from > 0 {
		switch b := line[from-1]; {
		case isWordByte(b), b == '@', b == '#', b == '$', b == '/', b == '.', b == '&':
			return false
		case b == '(' && from > 1 && line[from-2] == ']':
			return false
		}
	}
	return to >= len(line) || line[to] != '-'
}

func isWordByte(b byte) bool {
	return b == '_' || b >= '0' && b <= '9' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

// commentText is the content with everything but its comments blanked, lines kept, and the
// offsets where each comment starts — its delimiter, which a `#name` there is not.
func commentText(content string, cs *config.Comments) (string, map[int]bool) {
	b := []byte(content)
	keep := make([]bool, len(b))
	starts := map[int]bool{}
	for _, s := range cs.Spans(content) {
		starts[s[0]] = true
		for j := s[0]; j < s[1] && j < len(b); j++ {
			keep[j] = true
		}
	}
	for j := range b {
		if !keep[j] && b[j] != '\n' {
			b[j] = ' '
		}
	}
	return string(b), starts
}

// withoutFences is a spec's text with its fenced blocks blanked, lines kept: a fence holds
// an example in another language, where `$` and `#` mean what that language says.
func withoutFences(content string) string {
	lines := strings.Split(content, "\n")
	in := false
	for i, l := range lines {
		fence := strings.HasPrefix(strings.TrimSpace(l), "```") || strings.HasPrefix(strings.TrimSpace(l), "~~~")
		if fence || in {
			lines[i] = ""
		}
		if fence {
			in = !in
		}
	}
	return strings.Join(lines, "\n")
}

// projectSymbolsCache holds the names of each map's code, read once per map: every target
// of a run confronts the same project.
var projectSymbolsCache sync.Map // *mapx.Graph → map[string]bool

var symbolIdentRE = regexp.MustCompile(`[A-Za-z_$][A-Za-z0-9_$]*`)

// projectSymbols are the identifiers the map's code and test files hold outside their
// comments — a name that only a comment holds is the citation itself, no symbol.
//
// Each file's identifiers are kept between runs, by the file's revision in the map
// (symbolStore): reading and scanning every code file again took 15 s of a project's commit
// of 8 files, charged to the first target. A file with changes not committed is read again,
// whatever its revision says.
func projectSymbols(g *mapx.Graph, root string, cfg *config.Config) map[string]bool {
	if v, ok := projectSymbolsCache.Load(g); ok {
		return v.(map[string]bool)
	}
	store := loadSymbolStore(root, cfg)
	kept := map[string]symbolEntry{}
	changed := false
	known := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Kind != mapx.KindCode && n.Kind != mapx.KindTest {
			continue
		}
		e, ok := store.Files[n.ID]
		if dirty, _ := uncommitted(root, n.ID); !ok || n.Rev == "" || e.Rev != n.Rev || dirty {
			b, err := readFile(root, n.ID)
			if err != nil {
				continue
			}
			src := string(b)
			if cs := cfg.CommentsFor(n.ID); cs != nil {
				src = cs.Blank(src)
			}
			e = symbolEntry{Rev: n.Rev, Idents: uniqueIdents(src)}
			changed = true
		}
		kept[n.ID] = e
		for _, id := range e.Idents {
			known[id] = true
		}
	}
	if changed || len(kept) != len(store.Files) {
		store.Files = kept
		saveSymbolStore(root, store)
	}
	projectSymbolsCache.Store(g, known)
	return known
}

func uniqueIdents(src string) []string {
	seen := map[string]bool{}
	var out []string
	for _, id := range symbolIdentRE.FindAllString(src, -1) {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

// symbolStore is what symbol-exists keeps between runs: each file's identifiers by its
// revision, for one configuration — a change to anchors.yaml (its comment syntaxes) reads
// every file again.
type symbolStore struct {
	Version int
	Config  string
	Files   map[string]symbolEntry
}

type symbolEntry struct {
	Rev    string
	Idents []string
}

const symbolStoreVersion = 1

// symbolConfigKey says which configuration the names were read with: the comment syntaxes
// decide what is blanked, so another dialect reads every file again.
func symbolConfigKey(cfg *config.Config) string {
	b, _ := json.Marshal(struct {
		Rev      string
		Dialect  *config.Dialect
		Dialects map[string]*config.Dialect
	}{cfg.Rev(), cfg.Dialect, cfg.Dialects})
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:8])
}

// symbolStorePath is in the repository's common git folder, shared by its worktrees; empty
// outside git.
func symbolStorePath(root string) string {
	common, ok := gitmeta.CommonDir(root)
	if !ok {
		return ""
	}
	return filepath.Join(common, "anchors", "symbols.gob")
}

func loadSymbolStore(root string, cfg *config.Config) symbolStore {
	empty := symbolStore{Version: symbolStoreVersion, Config: symbolConfigKey(cfg), Files: map[string]symbolEntry{}}
	p := symbolStorePath(root)
	if p == "" {
		return empty
	}
	f, err := os.Open(p)
	if err != nil {
		return empty
	}
	defer f.Close()
	var s symbolStore
	if gob.NewDecoder(f).Decode(&s) != nil || s.Version != symbolStoreVersion || s.Config != symbolConfigKey(cfg) || s.Files == nil {
		return empty
	}
	return s
}

// saveSymbolStore writes the store through a temporary file and a rename, so a worktree
// reading it at the same time never reads half of one.
func saveSymbolStore(root string, s symbolStore) {
	p := symbolStorePath(root)
	if p == "" || os.MkdirAll(filepath.Dir(p), 0o755) != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), "symbols-*.tmp")
	if err != nil {
		return
	}
	if gob.NewEncoder(tmp).Encode(s) != nil || tmp.Close() != nil {
		_ = os.Remove(tmp.Name())
		return
	}
	if os.Rename(tmp.Name(), p) != nil {
		_ = os.Remove(tmp.Name())
	}
}
