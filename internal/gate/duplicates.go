// @anchors
//   code: DPGTA
//   ref: GTDPG

package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// Duplicate declarations (DESIGN-duplicate-declarations.md): a gate that controls a kind of
// declaration says what its occurrences are, and the engine confronts a key declared twice in
// one file. Only declarations count — a citation repeats by nature —, and each kind has one
// owner, so a repeat is reported once.

// Occurrence is one declaration a gate controls: its key, the 1-based line it is declared on
// and, when it lives in another file than the node judged, that file.
type Occurrence struct {
	Key  string
	Line int
	File string
}

// occurrenceReader reads the occurrences a gate controls in a node.
type occurrenceReader struct {
	read func(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) []Occurrence
	// verdict is what a repeat makes the gate say: Fail, unless the gate measures repeats
	// as a divergence of its own.
	verdict Verdict
	// hint is the i18n key of what the gate tells about fixing a repeat of its kind — the
	// `#NN` suffix of a scenario —, beside the generic finding; empty for none.
	hint string
}

// occurrenceReaders are the gates with a reader, by check.
var occurrenceReaders = map[string]occurrenceReader{}

// HasDuplicateReader says whether a check confronts the repeats of what it declares.
func HasDuplicateReader(check string) bool {
	_, ok := occurrenceReaders[check]
	return ok
}

// duplicatesOf are the keys declared more than once, each with its occurrences in order.
func duplicatesOf(occ []Occurrence) map[string][]Occurrence {
	by := map[string][]Occurrence{}
	for _, o := range occ {
		if o.Key != "" {
			by[o.Key] = append(by[o.Key], o)
		}
	}
	for k, os := range by {
		if len(os) < 2 {
			delete(by, k)
		}
	}
	return by
}

// confrontDuplicates runs the gate's occurrence reader on a node it measured and, when a key
// is declared twice, turns the verdict into the reader's and adds the finding beside the
// gate's own.
func confrontDuplicates(g config.Gate, n mapx.Node, root string, graph *mapx.Graph, cfg *config.Config, v Verdict, detail string) (Verdict, string) {
	rd, ok := occurrenceReaders[g.Check]
	if !ok || !g.DuplicatesOn() || v == Skip {
		return v, detail
	}
	content, err := readFile(root, n.ID)
	if err != nil {
		return v, detail
	}
	dups := duplicatesOf(rd.read(string(content), n, root, graph, cfg))
	if len(dups) == 0 {
		return v, detail
	}
	keys := make([]string, 0, len(dups))
	for k := range dups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		var at []string
		for _, o := range dups[k] {
			where := fmt.Sprint(o.Line)
			if o.File != "" && o.File != n.ID {
				where = o.File + ":" + where
			}
			at = append(at, where)
		}
		parts = append(parts, i18n.T("gate.duplicates.key", k, len(dups[k]), strings.Join(at, ", ")))
	}
	msg := i18n.T("gate.duplicates.found", strings.Join(parts, "; "))
	if rd.hint != "" {
		msg += " " + i18n.T(rd.hint)
	}
	if detail != "" && v != Pass {
		msg = detail + "; " + msg
	}
	verdict := rd.verdict
	if v == Fail {
		verdict = Fail
	}
	return verdict, msg
}

// citingSections are the catalog sections whose rows open with a rule code they CITE, not
// define: what a rule uses, the open questions, the navigation rows a rule triggers, the
// transitions between states, the events a unit emits — a row declares the callback, and its
// first column names the rule it realizes —, a change history and what a plan revises.
var citingSections = []string{"rule_uses", "open", "navigation", "state_flow", "callbacks", "history", "plan_revision"}

// citingTitles are the titles of the citing sections: the catalog's, in every language, and
// the ones the project gives them in any layer — the file judged need not be in the layer
// that declares the title (a component's spec, judged as a spec).
func citingTitles(cfg *config.Config) map[string]bool {
	out := sectionTitles(citingSections, nil, "")
	for _, k := range citingSections {
		for _, t := range cfg.SectionTitlesFor(strings.ReplaceAll(k, "_", "-")) {
			out[strings.ToLower(t)] = true
		}
		for _, t := range cfg.SectionTitlesFor(k) {
			out[strings.ToLower(t)] = true
		}
	}
	return out
}

// definedRuleOccurrences are the rule codes a file defines — a heading, the first cell of a
// table row, a bold bullet —, each where it is defined, outside the sections that cite codes
// (and their subsections), alias and retired lines. A heading and the rows under it that open
// with its own code are one definition: the occurrence is counted once per heading block.
func definedRuleOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) []Occurrence {
	citing := citingTitles(cfg)
	re := defineRuleCaptureRE()
	type level struct {
		depth    int
		excluded bool
	}
	var stack []level
	excluded := func() bool {
		for _, l := range stack {
			if l.excluded {
				return true
			}
		}
		return false
	}
	var out []Occurrence
	block := map[string]bool{}
	for i, line := range strings.Split(content, "\n") {
		if m := headingRE.FindStringSubmatch(line); m != nil {
			depth := len(m[1])
			for len(stack) > 0 && stack[len(stack)-1].depth >= depth {
				stack = stack[:len(stack)-1]
			}
			title := strings.ToLower(strings.TrimSpace(m[2]))
			stack = append(stack, level{depth, citing[title] || decisionsRE.MatchString(line+"\n")})
			block = map[string]bool{}
		} else if strings.HasPrefix(line, "# ") {
			stack = nil
			block = map[string]bool{}
		}
		if excluded() || retiredLine(line) {
			continue
		}
		if _, ok := ruleAliasTarget(line); ok {
			continue
		}
		m := re.FindStringSubmatch(line)
		if m == nil || block[m[1]] || citingBullet(line) {
			continue
		}
		block[m[1]] = true
		out = append(out, Occurrence{Key: m[1], Line: i + 1})
	}
	return out
}

// citingBulletRE is a list item opening with a code in backticks and no bold: prose that
// names the rule (`- ` + "`CODE-B03` registers…" + `), where a definition is bold or bare.
var citingBulletRE = regexp.MustCompile("^\\s*[-*+]\\s+`")

func citingBullet(line string) bool { return citingBulletRE.MatchString(line) }

func init() {
	occurrenceReaders["rule-types"] = occurrenceReader{read: definedRuleOccurrences, verdict: Fail}
}
