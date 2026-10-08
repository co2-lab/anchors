// @anchors
//   code: DMGDP
//   ref: GTDPG

package gate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/flagx"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// The occurrence readers of flags, code and features (DPDCD-W03).

// flagScenarioOccurrences: each scenario a flag file declares, by its code.
func flagScenarioOccurrences(content string, n mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) []Occurrence {
	var out []Occurrence
	for _, s := range flagx.ParseContent(content, n.ID).Scenarios {
		if s.Code != "" {
			out = append(out, Occurrence{Key: s.Code, Line: s.Line})
		}
	}
	return out
}

// usedByOccurrences: each symbol a used-by flag stands above. Two flags on one symbol say two
// things of who imports it, and the gate reads only one of them.
func usedByOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) []Occurrence {
	var out []Occurrence
	for _, u := range scan.UsedByIn([]byte(content)) {
		if u.Symbol != "" {
			out = append(out, Occurrence{Key: u.Symbol, Line: u.Line})
		}
	}
	return out
}

// testIDOccurrences: each testID the spec's inventory declares — the first cell of a row of
// its test identifiers section. A handle the code exposes in two render branches is no repeat:
// only the spec's declarations count.
func testIDOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) []Occurrence {
	attr := testHandleAttr(cfg)
	loc := surfaceSectionRE.FindStringIndex(content)
	if attr == "" || loc == nil {
		return nil
	}
	body := content[loc[1]:]
	if end := regexp.MustCompile(`(?m)^#{1,6}\s`).FindStringIndex(body); end != nil {
		body = body[:end[0]]
	}
	start := strings.Count(content[:loc[1]], "\n") + 1
	var out []Occurrence
	for i, line := range strings.Split(body, "\n") {
		cells := rowCells(line)
		if cells == nil {
			continue
		}
		for _, m := range declaredTestIDRE.FindAllStringSubmatch(cells[0], -1) {
			id := strings.TrimPrefix(m[1], ":")
			if strings.EqualFold(id, attr) {
				continue
			}
			out = append(out, Occurrence{Key: id, Line: start + i})
		}
	}
	return out
}

// exampleRowOccurrences: each row of an outline's Examples, by the outline's code and the
// row's values — the label columns aside, as examples-match reads them.
func exampleRowOccurrences(content string, _ mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) []Occurrence {
	var out []Occurrence
	for _, o := range exampleTables(content) {
		if len(o.rows) == 0 {
			continue
		}
		// Within one table: two outlines of a rule (`#03`, `#04`) may share a row's values,
		// each in its own columns — the table is named by where it starts.
		table := fmt.Sprintf("%s (Examples at line %d)", o.code, o.rows[0].line-1)
		for _, r := range o.rows {
			out = append(out, Occurrence{Key: table + " | " + strings.Join(r.values, " | "), Line: r.line})
		}
	}
	return out
}

func init() {
	occurrenceReaders["flag-scenario-grammar"] = occurrenceReader{read: flagScenarioOccurrences, verdict: Fail}
	occurrenceReaders["used-by-declared"] = occurrenceReader{read: usedByOccurrences, verdict: Fail}
	occurrenceReaders["testid-consistent"] = occurrenceReader{read: testIDOccurrences, verdict: Fail}
	occurrenceReaders["examples-match"] = occurrenceReader{read: exampleRowOccurrences, verdict: Fail}
}
