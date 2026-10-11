// @anchors
//   code: GGTGD
//   ref: GDGTG

package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/initx"
	"github.com/co2-lab/anchors/internal/mapx"
)

func guidanceNode(id string) mapx.Node {
	return mapx.Node{ID: id, Kind: mapx.KindDoc, Layer: config.GuidanceLayer, Tags: []string{config.GuidanceLayer}}
}

func TestGuidanceLinked(t *testing.T) {
	t.Run("GDGTG-B01: The guidance points to the agreement", func(t *testing.T) {})
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "CONTRIBUTING.md"), []byte("# Contributing\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{Nodes: []mapx.Node{guidanceNode("CONTRIBUTING.md"), guidanceNode("AGENTS.md"), guidanceNode("README.md")}}
	for _, c := range []struct {
		id, content string
		want        Verdict
		says        string
	}{
		{"CONTRIBUTING.md", "# Contributing\n", Pass, ""},
		{"AGENTS.md", "Read `CONTRIBUTING.md` first.\n", Pass, ""},
		{"README.md", "# The project\n", Fail, "`README.md` does not point to `CONTRIBUTING.md`"},
	} {
		v, msg := checkGuidanceLinked(c.content, guidanceNode(c.id), root, g, nil)
		if v != c.want || !strings.Contains(msg, c.says) {
			t.Errorf("%s: %v %s", c.id, v, msg)
		}
	}
	if v, msg := checkGuidanceLinked("see CONTRIBUTING.md", guidanceNode("CLAUDE.md"), t.TempDir(), &mapx.Graph{}, nil); v != Fail || !strings.Contains(msg, "`CONTRIBUTING.md` does not exist") {
		t.Errorf("no agreement: %v %s", v, msg)
	}
	if v, _ := checkGuidanceLinked("x", mapx.Node{ID: "a.go", Kind: mapx.KindCode}, root, g, nil); v != Skip {
		t.Error("outside the layer: skipped")
	}
}

func TestGuidanceGatesInTheCatalog(t *testing.T) {
	t.Run("GDGTG-B02: The catalog offers the guidance gates", func(t *testing.T) {})
	got := map[string]config.Gate{}
	for _, g := range initx.DefaultGates(map[string]bool{"code": true}, false) {
		got[g.Name] = g
	}
	linked, honored := got["guidance-linked"], got["guidance-honored"]
	if linked.Check != "guidance-linked" || linked.IsBlocking() || linked.Tags[0] != config.GuidanceLayer {
		t.Errorf("guidance-linked: %+v", linked)
	}
	if !honored.IsJudgment() || honored.IsBlocking() || !honored.StaleOnConfig || honored.Tags[0] != config.GuidanceLayer ||
		!strings.Contains(honored.Ask, "anchors guide guidance") || !strings.Contains(honored.Ask, "CONTRADICT") ||
		!strings.Contains(honored.Ask, "NOT judge style") || !strings.Contains(honored.Ask, "PROPOSE the text") {
		t.Errorf("guidance-honored: %+v", honored)
	}
}

func TestHeaderGatesLeaveHeadlessLayers(t *testing.T) {
	t.Run("GDGTG-B03: A headless layer is asked no header", func(t *testing.T) {})
	cfg := &config.Config{Layers: map[string]config.Layer{config.GuidanceLayer: config.CanonicalGuidance(), "screen": {Kind: "code"}}}
	n := guidanceNode("README.md")
	if v, _ := checkersWithGraph["header-valid"]("# The project\n", n, "", nil, cfg); v != Skip {
		t.Errorf("header-valid: %v", v)
	}
	if v, _ := checkHeaderLayerDeclared("<!-- @anchors\n  layer: nope\n-->\n", n, "", nil, cfg); v != Skip {
		t.Errorf("header-layer-declared: %v", v)
	}
	prev := fixConfig
	fixConfig = cfg
	defer func() { fixConfig = prev }()
	if _, changed, _ := fixMissingHeader("# The project\n", n, "", &mapx.Graph{}); changed {
		t.Error("the fix writes no header in a headless layer")
	}
}
