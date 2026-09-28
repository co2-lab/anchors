package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// examples-match: the rows of an outline's Examples are the cases its tests run.
//
// A scenario outline says "for each of these rows"; the parameterised test that proves it
// runs its own table. Nothing kept the two tables together, and they drifted apart in the
// reference app in eight features: the Examples said `mesada` and the test `allowance`,
// `simplificado` against `simplified`, "Sem provisão" against "Nenhum", "Crédito" against
// "Cartão de crédito", and a row the code no longer had at all.
//
// The gate reads no test library: every value of a row must appear in the body of a test
// that cites the scenario's code, as a whole token and with its case — "Crédito" inside
// "Cartão de crédito" is not the row. The body is read as `test-has-assertion` reads it
// (`labels` and the tests script's `end` included).
//
// Examples often show what the SCREEN says while the test uses a key. A column whose
// header carries `(label)` — in any supported language — is display text and is not
// looked for; the key goes in a column beside it, and that one is.
func checkExamplesMatch(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, i18n.T("gate.examples_match.skip_not_feature")
	}
	outlines := exampleTables(content)
	if len(outlines) == 0 {
		return Skip, i18n.T("gate.examples_match.skip_no_examples")
	}
	if !testsSource(cfg).Declared() {
		return Skip, i18n.T("gate.test_has_assertion.skip_no_source")
	}
	tests, _, err := projectTests(root, g, cfg)
	if err != nil {
		return Fail, err.Error()
	}
	labels := gateEntry(cfg, "examples-match").Labels
	files := map[string][]string{}
	bodyOf := func(file string) []string {
		if l, ok := files[file]; ok {
			return l
		}
		b, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			files[file] = nil // @resilient: an unreadable test carries no row; the map notices a missing file
			return nil
		}
		files[file] = strings.Split(string(b), "\n")
		return files[file]
	}
	confronted := false
	var gaps []string
	for _, o := range outlines {
		var bodies []string
		for _, t := range tests {
			if strings.Contains(t.Title, o.code) {
				bodies = append(bodies, testBody(bodyOf(t.File), t, labels))
			}
		}
		if len(bodies) == 0 {
			continue // no test cites the scenario: scenario-coverage's to say
		}
		confronted = true
		text := strings.Join(bodies, "\n")
		for _, r := range o.rows {
			var missing []string
			for _, v := range r.values {
				if !tokenIn(text, v) {
					missing = append(missing, "`"+v+"`")
				}
			}
			if len(missing) > 0 {
				gaps = append(gaps, fmt.Sprintf("%s %s %d: %s", o.code, i18n.T("gate.examples_match.line"), r.line, strings.Join(missing, ", ")))
			}
		}
	}
	if !confronted {
		return Skip, i18n.T("gate.examples_match.skip_no_test")
	}
	if len(gaps) > 0 {
		return Fail, i18n.T("gate.examples_match.fail", len(gaps), strings.Join(gaps, "; "))
	}
	return Pass, ""
}

type exampleRow struct {
	line   int
	values []string
}

type exampleOutline struct {
	code string
	rows []exampleRow
}

// examplesRE opens an examples table, in any language of the Gherkin table.
var examplesRE = regexp.MustCompile(`^\s*(?:` + strings.Join(quoteAll(config.GherkinExamplesAlternatives()), "|") + `)\s*:`)

// exampleTables reads each coded scenario's Examples: the rows under each table's header,
// with the values of the columns that are not labels.
func exampleTables(content string) []exampleOutline {
	var out []exampleOutline
	var cur *exampleOutline
	pending := ""
	inTable, header := false, []bool(nil)
	flush := func() {
		if cur != nil && len(cur.rows) > 0 {
			out = append(out, *cur)
		}
		cur = nil
	}
	for i, ln := range strings.Split(content, "\n") {
		t := strings.TrimSpace(ln)
		switch {
		case featScenarioCodeRE.MatchString(ln) && strings.HasPrefix(t, "@"):
			flush()
			m := featScenarioCodeRE.FindStringSubmatch(ln)
			pending, inTable = m[1], false
		case featTitleRE.MatchString(ln):
			flush()
			if pending != "" {
				cur = &exampleOutline{code: pending}
			}
			pending, inTable = "", false
		case examplesRE.MatchString(ln):
			inTable, header = true, nil
		case inTable && strings.HasPrefix(t, "|"):
			cells := cellsOf(t)
			if header == nil {
				header = make([]bool, len(cells))
				for j, h := range cells {
					header[j] = isLabelColumn(h)
				}
				continue
			}
			if cur == nil {
				continue
			}
			var vals []string
			for j, c := range cells {
				if j < len(header) && header[j] {
					continue
				}
				if v := strings.Trim(strings.TrimSpace(c), "\"'`"); v != "" {
					vals = append(vals, v)
				}
			}
			if len(vals) > 0 {
				cur.rows = append(cur.rows, exampleRow{line: i + 1, values: vals})
			}
		case t != "" && !strings.HasPrefix(t, "#"):
			inTable = false
		}
	}
	flush()
	return out
}

// isLabelColumn says whether an Examples header marks its column as display text:
// `(label)`, in any supported language.
func isLabelColumn(header string) bool {
	h := strings.ToLower(header)
	for _, w := range i18n.AllTranslations("gate.examples_match.label_marker") {
		for _, m := range strings.Split(w, "|") {
			if m = strings.TrimSpace(m); m != "" && strings.Contains(h, "("+m+")") {
				return true
			}
		}
	}
	return false
}

// tokenIn says whether v appears in text as a whole token: not glued to a letter, a digit
// or an underscore on either side, and with its case.
func tokenIn(text, v string) bool {
	return regexp.MustCompile(`(?:^|[^\p{L}\p{N}_])` + regexp.QuoteMeta(v) + `(?:$|[^\p{L}\p{N}_])`).MatchString(text)
}
