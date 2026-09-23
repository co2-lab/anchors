package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// value-anchored: a KEY that is replicated across the code is declared in a comment, and
// the code line right below it carries the value.
//
//	// @code-reference-[COLOR-SUCCESS]-[#1F8A5B]
//	success: '#1F8A5B',
//
// It exists for the value that lives in more than one place. Changing a colour means
// changing it everywhere, and nothing says where "everywhere" is: the copies are plain
// literals, and the one that was forgotten keeps compiling. Declaring the key turns the
// literal into a symbol the gate can follow.
//
// Two checks, and the second is what makes the first worth having:
//
//	LOCAL   the next CODE line below the declaration (comment lines skipped) contains the
//	        declared value — the declaration does not lie about its own line.
//	ACROSS  every declaration of the same key declares the same value — so changing it in
//	        one place, declaration included, reports every place that stayed behind.
//
// The local check alone does not propagate: someone updates the colour and its
// declaration in one file, and both files are locally consistent while disagreeing with
// each other.
//
// THE KEY MAY HAVE A SOURCE IN THE SPEC, or be only a replicated reference:
//
//	SPEC RULE  the key is a rule code (`TKNS-R01`) whose defining line in the spec
//	           declares the value in backticks (| `TKNS-R01` | success | `#1F8A5B` |).
//	           Every declaration is confronted with it: change the value in the rule, and
//	           every place still carrying the old one is reported.
//	FREE KEY   no source in any spec — the copies themselves are the truth, and what is
//	           charged is that they agree.
//
// A rule whose defining line declares no value is not charged against the spec: most rules
// are prose, and charging them would fail every declaration that points at one.
//
// WHAT CHANGED, and why. The first version charged an anchor on EVERY value of every
// closed set (`export const X = [...]`). The decision it encoded — every set item is a
// domain decision with an address — was replaced by a narrower one: anchoring is for keys
// that are REPLICATED, and whoever replicates declares. A value nobody declared is not
// charged. With that, the gate no longer needs to find where a set begins or ends: it
// follows declarations.
//
// WHAT IT DOES NOT DO: judge whether a value is right, nor decide which comment shape is a
// declaration — the project declares the pattern (`derived.value_anchor`).
func checkValueAnchored(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindCode {
		return Skip, i18n.T("gate.value_anchored.skip_not_code")
	}
	anchorRE := valueAnchorDe(cfg)
	if anchorRE == nil {
		return Skip, i18n.T("gate.value_anchored.skip_no_value_anchor")
	}
	if g == nil {
		return pendingNoMap()
	}

	decls := declarationsIn(content, n.ID, anchorRE)
	if len(decls) == 0 {
		return Skip, i18n.T("gate.value_anchored.skip_no_declaration")
	}

	var b strings.Builder

	// LOCAL: the declaration against the line it annotates.
	var lying []valueDecl
	for _, d := range decls {
		if !d.matches() {
			lying = append(lying, d)
		}
	}
	if len(lying) > 0 {
		fmt.Fprintf(&b, i18n.T("gate.value_anchored.lying_header"), len(lying), n.ID)
		for _, d := range lying {
			if d.code == "" {
				fmt.Fprintf(&b, i18n.T("gate.value_anchored.no_code_below"), d.line, d.key, d.value)
				continue
			}
			b.WriteString(fmt.Sprintf(i18n.T("gate.value_anchored.lying_item"),
				d.line, d.key, d.value, d.codeLine, strings.TrimSpace(d.code)))
		}
		b.WriteString(i18n.T("gate.value_anchored.lying_footer"))
	}

	idx := anchorIndexFor(root, g, anchorRE)

	// SPEC: a key that is a rule declaring a value in its defining line.
	var stale []valueDecl
	for _, d := range decls {
		if want, ok := idx.specValues[d.key]; ok && !want[d.value] {
			stale = append(stale, d)
		}
	}
	if len(stale) > 0 {
		fmt.Fprintf(&b, i18n.T("gate.value_anchored.spec_header"), len(stale))
		for _, d := range stale {
			fmt.Fprintf(&b, i18n.T("gate.value_anchored.spec_item"),
				d.line, d.key, d.value, idx.specFile[d.key], strings.Join(sortedKeys(idx.specValues[d.key]), "`, `"))
		}
		b.WriteString(i18n.T("gate.value_anchored.spec_footer"))
	}

	// ACROSS: the declarations of the same key elsewhere.
	index := idx.byKey
	reported := map[string]bool{}
	for _, d := range decls {
		if reported[d.key] {
			continue
		}
		sites := index[d.key]
		values := map[string]bool{}
		for _, s := range sites {
			values[s.value] = true
		}
		if len(values) < 2 {
			continue
		}
		reported[d.key] = true
		fmt.Fprintf(&b, i18n.T("gate.value_anchored.divergent_header"), d.key, len(values), len(sites))
		for _, s := range sites {
			fmt.Fprintf(&b, i18n.T("gate.value_anchored.divergent_item"), s.file, s.line, s.value)
		}
	}
	if len(reported) > 0 {
		b.WriteString(i18n.T("gate.value_anchored.divergent_footer"))
	}

	if b.Len() == 0 {
		return Pass, ""
	}
	return Fail, strings.TrimRight(b.String(), "\n")
}

// valueAnchorDe returns the pattern the project declared, or nil if it declared none.
// It requires TWO groups: without the second the declaration asserts no value, and the
// confrontation the gate exists for cannot happen.
func valueAnchorDe(cfg *config.Config) *regexp.Regexp {
	if cfg == nil || cfg.Derived == nil || strings.TrimSpace(cfg.Derived.ValueAnchor) == "" {
		return nil
	}
	re, err := regexp.Compile(cfg.Derived.ValueAnchor)
	if err != nil || re.NumSubexp() < 2 {
		return nil
	}
	return re
}

// valueDecl is one declaration of a replicated key, and the code line it annotates.
type valueDecl struct {
	file, key, value string
	line             int    // the declaration's line (1-based)
	code             string // the first CODE line below it; empty when there is none
	codeLine         int
}

// matches reports whether the annotated code line contains the declared value.
func (d valueDecl) matches() bool {
	return d.code != "" && strings.Contains(d.code, d.value)
}

// declarationsIn finds every declaration in a file and pairs it with the first CODE line
// below it.
//
// Comment lines between the declaration and the code are skipped — several declarations
// may be stacked above the same line, and a longer explanation may sit between them.
// Blank lines are skipped too. Which lines are comments is read from their start (`//`,
// `#`, `/*`, `*`, `--`, `<!--`), the same reading the other gates use.
func declarationsIn(content, file string, anchorRE *regexp.Regexp) []valueDecl {
	lines := strings.Split(content, "\n")
	var out []valueDecl
	for i, l := range lines {
		m := anchorRE.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		d := valueDecl{file: file, key: m[1], value: m[2], line: i + 1}
		for j := i + 1; j < len(lines); j++ {
			t := strings.TrimSpace(lines[j])
			if t == "" || isCommentLine(t) || anchorRE.MatchString(lines[j]) {
				continue
			}
			d.code, d.codeLine = lines[j], j+1
			break
		}
		out = append(out, d)
	}
	return out
}

// isCommentLine reports whether a trimmed line is a whole-line comment.
func isCommentLine(t string) bool {
	for _, p := range []string{"//", "#", "/*", "*", "--", "<!--"} {
		if strings.HasPrefix(t, p) {
			return true
		}
	}
	return false
}

// anchorIndexCache keeps one index per map for the duration of a run: the ACROSS check
// needs every declaration of the project, and building it once per code file would read
// the whole repository once per file — the shape of cost that once made `docs-fresh` 97%
// of a check.
var (
	anchorIndexMu    sync.Mutex
	anchorIndexGraph *mapx.Graph
	anchorIndexVal   *anchorIndex
)

// anchorIndex is what the gate needs from the whole project, built once per map.
type anchorIndex struct {
	byKey      map[string][]valueDecl     // every declaration, by key
	specValues map[string]map[string]bool // rule code → the values its defining line declares
	specFile   map[string]string          // rule code → the spec that defines it
}

// wholeBacktickRE is a table cell that is ENTIRELY one backticked token.
var wholeBacktickRE = regexp.MustCompile("^`([^`]+)`$")

// declaredValues reads the values a rule's defining line declares: the table cells that
// are ENTIRELY one backticked token, other than the rule's own code.
//
// Only whole cells, and only table rows. A heading or a prose cell carries backticked
// identifiers all the time (`### CRED-B01 — validates the \`limit\“), and reading them
// as values would charge every declaration pointing at an ordinary rule.
func declaredValues(line, code string) []string {
	t := strings.TrimSpace(line)
	if !strings.HasPrefix(t, "|") {
		return nil
	}
	var out []string
	for _, cell := range strings.Split(strings.Trim(t, "|"), "|") {
		m := wholeBacktickRE.FindStringSubmatch(strings.TrimSpace(cell))
		if m == nil || m[1] == code {
			continue
		}
		out = append(out, m[1])
	}
	return out
}

// anchorIndexFor returns every declaration of every code file, and the values every spec
// rule declares, in a stable order.
func anchorIndexFor(root string, g *mapx.Graph, anchorRE *regexp.Regexp) *anchorIndex {
	anchorIndexMu.Lock()
	defer anchorIndexMu.Unlock()
	if anchorIndexGraph == g && anchorIndexVal != nil {
		return anchorIndexVal
	}
	idx := &anchorIndex{specValues: map[string]map[string]bool{}, specFile: map[string]string{}}
	defineRE := defineRuleCaptureRE()
	for _, node := range g.Nodes {
		if node.Kind != mapx.KindSpec {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, node.ID))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			m := defineRE.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			for _, v := range declaredValues(line, m[1]) {
				if idx.specValues[m[1]] == nil {
					idx.specValues[m[1]] = map[string]bool{}
					idx.specFile[m[1]] = node.ID
				}
				idx.specValues[m[1]][v] = true
			}
		}
	}
	index := map[string][]valueDecl{}
	for _, node := range g.Nodes {
		if node.Kind != mapx.KindCode {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, node.ID))
		if err != nil {
			continue
		}
		for _, d := range declarationsIn(string(b), node.ID, anchorRE) {
			index[d.key] = append(index[d.key], d)
		}
	}
	for k := range index {
		sort.Slice(index[k], func(i, j int) bool {
			if index[k][i].file != index[k][j].file {
				return index[k][i].file < index[k][j].file
			}
			return index[k][i].line < index[k][j].line
		})
	}
	idx.byKey = index
	anchorIndexGraph, anchorIndexVal = g, idx
	return idx
}

// sortedKeys lists a set in a stable order, for the verdict.
func sortedKeys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
