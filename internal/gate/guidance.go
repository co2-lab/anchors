// @anchors
//   code: GLCGD
//   ref: GDGTG

package gate

import (
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- the guidance layer: the working agreement is one file, and the others point to it ---
//
// README.md, CONTRIBUTING.md, AGENTS.md and CLAUDE.md tell whoever arrives how the project is
// worked on (DESIGN-guidance-layer.md). The agreement lives in CONTRIBUTING.md; the others send
// the reader there, so it is said once. What can be measured is measured here; whether what
// they say contradicts Anchors and the configuration is the judgment `guidance-honored` asks.

// guidanceAgreement is the file that holds the working agreement.
const guidanceAgreement = "CONTRIBUTING.md"

// checkGuidanceLinked: the agreement exists, and every other guidance file points to it.
func checkGuidanceLinked(content string, n mapx.Node, root string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	if !slices.Contains(n.Tags, config.GuidanceLayer) {
		return Skip, i18n.T("gate.guidance_linked.skip_not_guidance")
	}
	if path.Base(n.ID) == guidanceAgreement {
		return Pass, ""
	}
	agreement := path.Join(path.Dir(n.ID), guidanceAgreement)
	if (g == nil || g.Node(agreement) == nil) && !exists(root, agreement) {
		return Fail, i18n.T("gate.guidance_linked.no_agreement", agreement)
	}
	if !strings.Contains(content, guidanceAgreement) {
		return Fail, i18n.T("gate.guidance_linked.no_pointer", n.ID, guidanceAgreement)
	}
	return Pass, ""
}

func exists(root, rel string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil
}
