// @anchors
//   ref: APISP

package gate

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// The coherence of an API spec, asked before anything is compiled or run:
//
//	api-contracts-resolve   every contract the body and the responses cite is a spec with a
//	                        Domain, and no response leaves its contract unsaid
//	api-errors-declared     every error response has its code and message, under a status the
//	                        Responses declare
//	error-codes-honored     every error code the spec declares is one the unit's code emits
//
// The OpenAPI build already fails on a contract that is no spec; asked here, per unit, the
// failure names the unit while it is being written, and the pre-commit stops the commit that
// broke it instead of the docs build that came later.

// apiUnit is an API unit read from its main code file: the spec beside it, when that spec
// has an Endpoint section.
type apiUnit struct {
	base, spec, code string
	content          string
}

func readAPIUnit(n mapx.Node, root string) (*apiUnit, string, bool) {
	if n.Kind != mapx.KindCode {
		return nil, "", false
	}
	base := strings.TrimSuffix(n.ID, path.Ext(n.ID))
	b, err := readFile(root, base+".spec.md")
	if err != nil {
		return nil, i18n.T("gate.vr_states.skip_no_spec"), false
	}
	if !endpointSectionRE.Match(b) {
		return nil, i18n.T("gate.contract_tested.skip_not_api"), false
	}
	u := &apiUnit{base: base, spec: base + ".spec.md", content: string(b)}
	if m := specCodeRE().FindStringSubmatch(u.content); m != nil {
		u.code = m[1]
	}
	return u, "", true
}

// sectionRows reads the first table of the section the catalog key names, under its title
// in any language, as rows keyed by the header's cells lower-cased.
func sectionRows(content, key string) []map[string]string {
	if rest, ok := sectionText(content, key); ok {
		var header []string
		var rows []map[string]string
		for _, line := range strings.Split(rest, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "|") {
				if len(rows) > 0 {
					break
				}
				continue
			}
			cells := cellsOf(line)
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
			for i, c := range cells {
				if i < len(header) {
					row[header[i]] = strings.TrimSpace(c)
				}
			}
			rows = append(rows, row)
		}
		return rows
	}
	return nil
}

// sectionText is the body of the section the catalog key names, under its title in any
// language — the title alone, or followed by a note in parentheses (`## States (Estados da
// Tela)`) —, at any heading level, up to the next heading of the same level or above.
//
// Any level, because specs nest: `### Validações` under `## Rules`, as a project wrote them,
// was a section the gates did not find, and each of its units left without a verdict.
func sectionText(content, key string) (string, bool) {
	for _, title := range i18n.AllTranslations(key) {
		re := regexp.MustCompile(`(?m)^(#{2,6})\s+` + regexp.QuoteMeta(title) + `[ \t]*(?:\([^)\n]*\))?[ \t]*$`)
		m := re.FindStringSubmatchIndex(content)
		if m == nil {
			continue
		}
		level := m[3] - m[2]
		rest := content[m[1]:]
		end := regexp.MustCompile(`(?m)^#{2,` + fmt.Sprint(level) + `}\s`)
		if loc := end.FindStringIndex(rest); loc != nil {
			rest = rest[:loc[0]]
		}
		return rest, true
	}
	return "", false
}

// col is a row's cell under the first of the header names that the row has.
func col(row map[string]string, names ...string) (string, bool) {
	for _, n := range names {
		if v, ok := row[n]; ok {
			return v, true
		}
	}
	return "", false
}

var (
	contractCols  = []string{"contract", "contrato"}
	statusCols    = []string{"status"}
	errorCodeCols = []string{"error code", "código de erro", "codigo de erro", "código de error", "codigo de error"}
	messageCols   = []string{"message", "mensagem", "mensaje"}
	ruleCols      = []string{"rule", "regra", "regla"}
)

var citedCodeRE = regexp.MustCompile("`([A-Z0-9][A-Z0-9-]*)`")

// checkAPIContractsResolve: is every contract the body and the responses cite a spec with a
// Domain, and does every response say its contract?
func checkAPIContractsResolve(_ string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	u, why, ok := readAPIUnit(n, root)
	if !ok {
		return Skip, why
	}
	specByCode := map[string]string{}
	if g != nil {
		for _, x := range g.Nodes {
			if x.Kind == mapx.KindSpec && x.Code != "" {
				specByCode[x.Code] = x.ID
			}
		}
	}
	var unsaid, missing, noDomain []string
	check := func(where, cellv string, present bool) {
		cellv = strings.TrimSpace(cellv)
		switch {
		case !present:
			return
		case cellv == "" || strings.HasPrefix(strings.ToUpper(cellv), "TODO"):
			unsaid = append(unsaid, where)
		case cellv == "—" || cellv == "-":
		default:
			m := citedCodeRE.FindStringSubmatch(cellv)
			if m == nil {
				missing = append(missing, cellv)
				return
			}
			id, ok := specByCode[m[1]]
			if !ok {
				missing = append(missing, m[1])
				return
			}
			b, err := readFile(root, id)
			if err != nil || len(sectionRows(string(b), "section.title.domain")) == 0 {
				noDomain = append(noDomain, m[1])
			}
		}
	}
	for _, r := range sectionRows(u.content, "section.title.request_body") {
		v, ok := col(r, contractCols...)
		check(i18n.T("gate.api_spec.where_body"), v, ok)
	}
	for _, r := range sectionRows(u.content, "section.title.responses") {
		v, ok := col(r, contractCols...)
		st, _ := col(r, statusCols...)
		check(st, v, ok)
	}
	var gaps []string
	if len(unsaid) > 0 {
		gaps = append(gaps, i18n.T("gate.api_spec.contract_unsaid", strings.Join(unsaid, ", ")))
	}
	if len(missing) > 0 {
		gaps = append(gaps, i18n.T("gate.api_spec.contract_missing", strings.Join(dedupe(missing), ", ")))
	}
	if len(noDomain) > 0 {
		gaps = append(gaps, i18n.T("gate.api_spec.contract_no_domain", strings.Join(dedupe(noDomain), ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// checkAPIErrorsDeclared: does every error response carry its code and message, under a
// status the Responses declare (exactly, or in a range like `4xx`)?
func checkAPIErrorsDeclared(_ string, n mapx.Node, root string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
	u, why, ok := readAPIUnit(n, root)
	if !ok {
		return Skip, why
	}
	errs := sectionRows(u.content, "section.title.error_responses")
	if len(errs) == 0 {
		return Skip, i18n.T("gate.api_spec.skip_no_errors")
	}
	declared := map[string]bool{}
	for _, r := range sectionRows(u.content, "section.title.responses") {
		if st, ok := col(r, statusCols...); ok {
			declared[strings.ToLower(strings.Trim(st, "`* "))] = true
		}
	}
	var undeclared, incomplete []string
	for _, r := range errs {
		rule, _ := col(r, ruleCols...)
		rule = strings.Trim(rule, "` ")
		st, _ := col(r, statusCols...)
		st = strings.ToLower(strings.Trim(st, "`* "))
		if len(st) == 3 && !declared[st] && !declared[st[:1]+"xx"] {
			undeclared = append(undeclared, fmt.Sprintf("%s (%s)", rule, st))
		}
		code, _ := col(r, errorCodeCols...)
		msg, _ := col(r, messageCols...)
		if blank(code) || blank(msg) {
			incomplete = append(incomplete, rule)
		}
	}
	var gaps []string
	if len(undeclared) > 0 {
		gaps = append(gaps, i18n.T("gate.api_spec.status_undeclared", strings.Join(undeclared, ", ")))
	}
	if len(incomplete) > 0 {
		gaps = append(gaps, i18n.T("gate.api_spec.error_incomplete", strings.Join(incomplete, ", ")))
	}
	if len(gaps) == 0 {
		return Pass, ""
	}
	return Fail, strings.Join(gaps, "; ")
}

// checkErrorCodesHonored: is every error code the spec declares one the unit's code emits?
//
// The code read is the unit's main file and every file its spec `specifies`, comments out:
// a code written only in a comment is no answer the client receives.
func checkErrorCodesHonored(_ string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	u, why, ok := readAPIUnit(n, root)
	if !ok {
		return Skip, why
	}
	var codes []string
	for _, r := range sectionRows(u.content, "section.title.error_responses") {
		if c, _ := col(r, errorCodeCols...); !blank(c) {
			codes = append(codes, strings.Trim(c, "` "))
		}
	}
	if len(codes) == 0 {
		return Skip, i18n.T("gate.api_spec.skip_no_errors")
	}
	paths := []string{n.ID}
	if g != nil {
		for _, e := range g.Neighbors(u.spec).Out {
			if e.Type == mapx.EdgeSpecifies && e.To != n.ID {
				paths = append(paths, e.To)
			}
		}
	}
	var body strings.Builder
	for _, p := range paths {
		if b, err := readFile(root, p); err == nil {
			body.WriteString(stripLineComments(string(b)))
			body.WriteString("\n")
		}
	}
	var absent []string
	for _, c := range dedupe(codes) {
		if !strings.Contains(body.String(), c) {
			absent = append(absent, c)
		}
	}
	if len(absent) == 0 {
		return Pass, ""
	}
	return Fail, i18n.T("gate.api_spec.code_absent", strings.Join(absent, ", "), strings.Join(paths, ", "))
}

func blank(s string) bool {
	s = strings.Trim(strings.TrimSpace(s), "`\"")
	return s == "" || s == "—" || s == "-" || strings.HasPrefix(strings.ToUpper(s), "TODO")
}

func dedupe(xs []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range xs {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	sort.Strings(out)
	return out
}
