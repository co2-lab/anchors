// @anchors
//   code: BYPRC
//   ref: BYPRS

package gate

import (
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// bypass-proven: a path the code skips while a test runs is proven by a suite that runs it.
//
// Measured in MIF (MIF-1S): a hook returned early under `EXPO_PUBLIC_E2E`, so the visual
// captures stayed clean, and the one suite against the real backend never ran the write the
// hook skipped. The unit tests doubled the write; the end-to-end skipped it; nothing ran it.
//
// The project declares its bypasses (`test_bypass.patterns`) and the suite layers whose proofs
// run the skipped path (`test_bypass.proven_by`). A file with a bypass asks a test of its own
// that passed in one of those layers, or `@no-bypass-proof: <reason>` on the bypass's line or
// the line above it.
func checkBypassProven(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if cfg == nil || cfg.TestBypass == nil {
		return Skip, i18n.T("gate.bypass_proven.skip_none_declared")
	}
	if n.Kind != mapx.KindCode {
		return Skip, i18n.T("gate.bypass_proven.skip_not_code")
	}
	res := cfg.TestBypass.PatternREs()
	lines := strings.Split(content, "\n")
	var open []string
	for i, l := range lines {
		if !matchesAny(res, l) {
			continue
		}
		if noBypassProofRE.MatchString(l) || (i > 0 && noBypassProofRE.MatchString(lines[i-1])) {
			continue
		}
		open = append(open, i18n.T("gate.bypass_proven.line", i+1))
	}
	if len(open) == 0 {
		return Skip, i18n.T("gate.bypass_proven.skip_no_bypass")
	}
	if g == nil {
		return pendingNoMap()
	}
	if provenIn(testsOfUnit(n.ID, g, cfg), g, cfg.TestBypass.ProvenBy) {
		return Pass, ""
	}
	return Fail, i18n.T("gate.bypass_proven.unproven", strings.Join(open, ", "), strings.Join(cfg.TestBypass.ProvenBy, ", "))
}

// noBypassProofRE is a bypass's waiver, with its reason.
var noBypassProofRE = regexp.MustCompile(`@no-bypass-proof[^\S\n]*:[^\S\n]*\S`)

func matchesAny(res []*regexp.Regexp, line string) bool {
	for _, re := range res {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}
