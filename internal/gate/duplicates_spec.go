// @anchors
//   code: DSGDP
//   ref: GTDPG

package gate

import (
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// The occurrence readers of the spec catalogue (DPDCD-W02): what each gate's section
// declares, one row each, where it is declared.

// lineAt is the 1-based line of content where body — a piece of content — starts.
func lineAt(content, body string) int {
	i := strings.Index(content, body)
	if i < 0 || body == "" {
		return 0
	}
	return strings.Count(content[:i], "\n") + 1
}

// tableRowsAt are the data rows of the first table in body, each with its line in content
// and its cells keyed by the header's names lower-cased (as sectionRows reads them).
func tableRowsAt(content, body string) []struct {
	line int
	row  map[string]string
} {
	start := lineAt(content, body)
	var out []struct {
		line int
		row  map[string]string
	}
	var header []string
	for i, l := range strings.Split(body, "\n") {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "|") {
			if len(out) > 0 {
				break
			}
			continue
		}
		cells := cellsOf(t)
		if header == nil {
			for _, h := range cells {
				header = append(header, strings.ToLower(strings.TrimSpace(h)))
			}
			continue
		}
		if strings.Trim(strings.Join(cells, ""), ":- ") == "" {
			continue
		}
		row := map[string]string{}
		for j, c := range cells {
			if j < len(header) {
				row[header[j]] = strings.TrimSpace(c)
			}
		}
		out = append(out, struct {
			line int
			row  map[string]string
		}{start + i, row})
	}
	return out
}

// envOccurrences: each environment variable a spec declares, by its row.
func envOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) []Occurrence {
	body, ok := sectionText(content, "section.title.environment")
	if !ok {
		return nil
	}
	var out []Occurrence
	for _, r := range tableRowsAt(content, body) {
		v, _ := col(r.row, variableCols...)
		v = strings.Trim(strings.TrimSpace(v), "`")
		if v == "" || strings.HasPrefix(strings.ToUpper(v), "TODO") {
			continue
		}
		out = append(out, Occurrence{Key: v, Line: r.line})
	}
	return out
}

// depRowLineRE is a row of the Dependencies table: it opens with its `DEPn`.
var depRowLineRE = regexp.MustCompile(`^\s*\|\s*` + "`?" + `(DEP\d+)` + "`?" + `\s*\|`)

// depOccurrences: each `DEPn` a spec's Dependencies table declares. Only a row opening with
// it declares it; a rule's uses cite it in another column.
func depOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) []Occurrence {
	var out []Occurrence
	for i, l := range strings.Split(content, "\n") {
		if m := depRowLineRE.FindStringSubmatch(l); m != nil {
			out = append(out, Occurrence{Key: m[1], Line: i + 1})
		}
	}
	return out
}

// domainOccurrences: each entry the Domain declares, by its row — the first cell, its case
// and backticks aside.
func domainOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) []Occurrence {
	body, ok := sectionDomain(content)
	if !ok {
		return nil
	}
	start := lineAt(content, body)
	var out []Occurrence
	for i, l := range strings.Split(body, "\n") {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "|") || separadorRE.MatchString(t) || cabecalhoDominioRE.MatchString(t) || todoOnlyRE.MatchString(t) {
			continue
		}
		cells := cellsOf(t)
		if len(cells) == 0 {
			continue
		}
		k := strings.ToLower(strings.Trim(strings.TrimSpace(cells[0]), "`* "))
		if k == "" {
			continue
		}
		out = append(out, Occurrence{Key: k, Line: start + i})
	}
	return out
}

// openQuestionOccurrences: each open question by its code — the first cell or the item's
// head —, resolved rows included: a code used again after its question was answered names
// two questions.
func openQuestionOccurrences(content string, n mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) []Occurrence {
	body, ok := openDecisionsSectionCfg(content, cfg, n.Layer)
	if !ok {
		return nil
	}
	start := lineAt(content, body)
	var out []Occurrence
	for i, l := range strings.Split(body, "\n") {
		if !itemRE.MatchString(l) {
			continue
		}
		if c := itemCode(l); c != "" {
			out = append(out, Occurrence{Key: c, Line: start + i})
		}
	}
	return out
}

// revisionOccurrences: each revision a spec or plan records, by its code — `CODE-R0001`
// opening the line.
func revisionOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) []Occurrence {
	re := revisionRE()
	var out []Occurrence
	for i, l := range strings.Split(content, "\n") {
		if m := re.FindStringSubmatch(l); m != nil {
			out = append(out, Occurrence{Key: m[1] + "-R" + m[2], Line: i + 1})
		}
	}
	return out
}

// sectionNoteRE is a title's note in parentheses: `States (Estados da Tela)`.
var sectionNoteRE = regexp.MustCompile(`\s*\([^)]*\)\s*$`)

// sectionOccurrences: each catalog section a spec opens — under its catalog title in any
// language, or the title the project gave it —, keyed by its parent heading and its title.
// Every reader of a section finds the first one; a second under the same parent is text no
// gate reads.
func sectionOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) []Occurrence {
	known := map[string]bool{}
	for _, ts := range projectSectionTitles {
		for _, t := range ts {
			known[strings.ToLower(t)] = true
		}
	}
	type head struct {
		depth int
		title string
	}
	var stack []head
	var out []Occurrence
	for i, l := range strings.Split(content, "\n") {
		m := headingRE.FindStringSubmatch(l)
		if m == nil {
			if strings.HasPrefix(l, "# ") {
				stack = nil
			}
			continue
		}
		depth := len(m[1])
		for len(stack) > 0 && stack[len(stack)-1].depth >= depth {
			stack = stack[:len(stack)-1]
		}
		title := strings.TrimSpace(m[2])
		parent := ""
		if len(stack) > 0 {
			parent = stack[len(stack)-1].title
		}
		stack = append(stack, head{depth, title})
		lower := strings.ToLower(title)
		bare := strings.ToLower(sectionNoteRE.ReplaceAllString(title, ""))
		if k, _ := i18n.SectionKeyFor(bare); k == "" && !known[lower] && !known[bare] {
			continue
		}
		key := lower
		if parent != "" {
			key = strings.ToLower(parent) + " › " + lower
		}
		out = append(out, Occurrence{Key: key, Line: i + 1})
	}
	return out
}

func init() {
	occurrenceReaders["env-declared"] = occurrenceReader{read: envOccurrences, verdict: Fail}
	occurrenceReaders["dependency-honored"] = occurrenceReader{read: depOccurrences, verdict: Fail}
	occurrenceReaders["domain-declared"] = occurrenceReader{read: domainOccurrences, verdict: Fail}
	occurrenceReaders["open-questions-resolved"] = occurrenceReader{read: openQuestionOccurrences, verdict: Fail}
	occurrenceReaders["revision-orphans"] = occurrenceReader{read: revisionOccurrences, verdict: Fail}
	occurrenceReaders["spec-sections"] = occurrenceReader{read: sectionOccurrences, verdict: Fail}
}
