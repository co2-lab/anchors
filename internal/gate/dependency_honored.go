// @anchors
//   code: DHGDP
//   ref: DEPHN

package gate

import (
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// checkDependencyHonored: does the spec still declare its dependencies in a table?
//
// It used to confront the table's methods with the code. A spec precedes the code, though,
// and the files a unit imports and the methods it calls are born with it: they are declared
// where the import is, by `@dep:`, and confronted there (DESIGN-dependencies-out-of-the-spec.md).
// What is left to ask of the spec is that the table is gone — a divergence to migrate, which
// `anchors migrate` removes.
func checkDependencyHonored(content string, n mapx.Node, _ string, _ *mapx.Graph, _ *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, i18n.T("gate.dependency_honored.skip_not_spec")
	}
	rows := 0
	for _, l := range strings.Split(content, "\n") {
		if depRowLineRE.MatchString(l) {
			rows++
		}
	}
	if rows == 0 {
		return Pass, ""
	}
	return Diverge, i18n.T("gate.dependency_honored.table", rows)
}

// depRowLineRE is a row of a Dependencies table: it opens with its `DEPn`. A `DEPn` cited in
// another column — a rule's uses — is no row of it.
var depRowLineRE = regexp.MustCompile(`^\s*\|\s*` + "`?" + `(DEP\d+)` + "`?" + `\s*\|`)
