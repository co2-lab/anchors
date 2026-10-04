// @anchors
//   code: HLDGH
//   ref: HDLYD

package gate

import (
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// header-layer-declared: the layer a header declares is one the Estrutura has.
//
// The header's `layer:` decides which templates derive the unit's siblings and where the
// documentation files it. A layer the Estrutura does not have is read as nothing: the unit
// falls back to the default templates, lands in a page of its own, and shows up — if at
// all — as a layer "outside every container" on the architecture page. In the reference
// app a section declared `layer: landing-component` for months; the Estrutura only had
// `landing-feature`, and it surfaced by chance.
//
// A placeholder (`layer: TODO`) is left to `placeholder-filled`: one defect, one gate.
func checkHeaderLayerDeclared(content string, n mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	layer := scan.HeaderLayerOf(content)
	if layer == "" {
		return Skip, i18n.T("gate.header_layer.skip_no_layer")
	}
	if strings.HasPrefix(strings.ToUpper(layer), "TODO") {
		return Skip, i18n.T("gate.header_layer.skip_placeholder")
	}
	if cfg == nil || len(cfg.Layers) == 0 {
		return Skip, i18n.T("gate.header_layer.skip_no_structure")
	}
	if _, ok := cfg.Layers[layer]; ok {
		return Pass, ""
	}
	declared := make([]string, 0, len(cfg.Layers))
	for name := range cfg.Layers {
		declared = append(declared, name)
	}
	sort.Strings(declared)
	for _, name := range declared {
		if strings.EqualFold(name, layer) {
			return Fail, i18n.T("gate.header_layer.fail_case", layer, name)
		}
	}
	return Fail, i18n.T("gate.header_layer.fail", layer, strings.Join(declared, ", "))
}
