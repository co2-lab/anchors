// @anchors
//   code: RDWRC
//   ref: RDWRS

package gate

import (
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// reads-written: every data a file's header says it reads is written by some file.
//
// The `reads:` and `writes:` of the headers tie a reader to the writer of the same data
// (`reads-from`), so the reader's evidence stales when the writer changes — measured in MIF:
// seven defects where the writer and the reader of one persisted record disagreed on a field,
// a sign or a vocabulary, each side green against its own fixture. A name no file writes
// makes no edge, and the reader would look tied while it is tied to nothing: a typo, or a
// writer that never declared itself.
func checkReadsWritten(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	var cs *config.Comments
	if cfg != nil {
		cs = cfg.CommentsFor(n.ID)
	}
	_, reads := scan.HeaderData([]byte(content), cs)
	if len(reads) == 0 {
		return Skip, i18n.T("gate.reads_written.skip_none")
	}
	if g == nil {
		return pendingNoMap()
	}
	tied := map[string]bool{}
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type == mapx.EdgeReadsFrom {
			tied[e.Method] = true
		}
	}
	var loose []string
	for _, d := range reads {
		if !tied[d] {
			loose = append(loose, d)
		}
	}
	if len(loose) == 0 {
		return Pass, ""
	}
	sort.Strings(loose)
	return Fail, i18n.T("gate.reads_written.unwritten", strings.Join(loose, ", "))
}
