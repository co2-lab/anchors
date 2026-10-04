// @anchors
//   code: EDGNV
//   ref: ENVDC

package gate

import (
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

var (
	variableCols   = []string{"variable", "variável", "variavel"}
	deprecatedCols = []string{"deprecated", "depreciada", "obsoleta"}
)

// env-declared: the environment variables a unit's code reads are the ones its spec
// declares — and only those.
//
// An environment variable is a contract with whoever deploys: a variable the code reads
// and no spec names is found missing in production, and a variable the spec names and the
// code no longer reads is configured forever for nothing. Like `contract-status-declared`
// with the status codes, both sides are confronted. A variable declared deprecated may stop
// being read: it is on its way out.
//
// The code is the files the spec `specifies`, comments removed, read with the dialect's
// `env_read` — the family's, or every family's together when the project declares none.
// A read by a computed name (`process.env[name]`) is out of reach of the text.
func checkEnvDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec || g == nil {
		return Skip, ""
	}
	declared := map[string]bool{}
	deprecated := map[string]bool{}
	for _, r := range sectionRows(content, "section.title.environment") {
		v, _ := col(r, variableCols...)
		v = strings.Trim(strings.TrimSpace(v), "`")
		if v == "" || strings.HasPrefix(strings.ToUpper(v), "TODO") {
			continue
		}
		declared[v] = true
		if d, _ := col(r, deprecatedCols...); isYes(d) {
			deprecated[v] = true
		}
	}
	var body strings.Builder
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type != mapx.EdgeSpecifies {
			continue
		}
		if b, err := readFile(root, e.To); err == nil {
			body.WriteString(stripLineComments(string(b)))
			body.WriteString("\n")
		}
	}
	read := map[string]bool{}
	if pat := cfg.DialectFor().EnvReadPattern(); pat != "" {
		if re, err := regexp.Compile(pat); err == nil {
			for _, m := range re.FindAllStringSubmatch(body.String(), -1) {
				for _, name := range m[1:] {
					if name != "" {
						read[name] = true
						break
					}
				}
			}
		}
	}
	if len(declared) == 0 && len(read) == 0 {
		return Skip, i18n.T("gate.env_declared.skip_none")
	}
	var undeclared, unread []string
	for v := range read {
		if !declared[v] {
			undeclared = append(undeclared, v)
		}
	}
	for v := range declared {
		if !read[v] && !deprecated[v] {
			unread = append(unread, v)
		}
	}
	sort.Strings(undeclared)
	sort.Strings(unread)
	var gaps []string
	if len(undeclared) > 0 {
		gaps = append(gaps, i18n.T("gate.env_declared.undeclared", len(undeclared), strings.Join(undeclared, ", ")))
	}
	if len(unread) > 0 {
		gaps = append(gaps, i18n.T("gate.env_declared.unread", len(unread), strings.Join(unread, ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// isYes reads a yes in a cell: `yes`, `sim`, `sí`, `true`, or one followed by a note
// (`yes: use PAYMENTS_KEY`).
func isYes(cell string) bool {
	c := strings.ToLower(strings.Trim(strings.TrimSpace(cell), "`*"))
	for _, y := range []string{"yes", "sim", "sí", "si", "true", "x", "✓"} {
		if c == y || strings.HasPrefix(c, y+":") || strings.HasPrefix(c, y+" ") || strings.HasPrefix(c, y+",") {
			return true
		}
	}
	return false
}
