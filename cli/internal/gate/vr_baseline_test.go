package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const featVR = `Funcionalidade: Tela

  @estado @TCDTX-S01 @nivel-integration
  Cenário: carrega
    Dado algo

  @estado @TCDTX-VR @nivel-vr @P3
  Cenário: aparência da tela carregada
    Dado algo
`

// TestCenarioVRSemBaselineEhAcusado guarda a superfície de prova que nenhum gate
// alcançava. `@nivel-vr` declara que a tela é provada por CAPTURA, não por asserção — e
// sem imagem de referência não há contra o que comparar. O cenário existe, o gate de
// feature o conta como coberto, e a prova prometida não acontece.
func TestCenarioVRSemBaselineEhAcusado(t *testing.T) {
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, "tela"), 0o755))

	v, msg := checkVRBaseline(featVR, mapx.Node{Kind: mapx.KindFeature, ID: "tela/T.feature"}, root, nil, nil)
	if v != Fail {
		t.Fatalf("cenário VR sem baseline deve reprovar; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "TCDTX-VR") {
		t.Errorf("a mensagem precisa nomear o cenário sem imagem; veio: %s", msg)
	}

	// A convenção aceita variante no nome: `TCDTX-VR-loaded.png` prova `TCDTX-VR`.
	must(t, os.WriteFile(filepath.Join(root, "tela", "T.TCDTX-VR-loaded.png"), []byte("png"), 0o644))
	if v, msg := checkVRBaseline(featVR, mapx.Node{Kind: mapx.KindFeature, ID: "tela/T.feature"}, root, nil, nil); v != Pass {
		t.Errorf("com baseline (mesmo com variante no nome) deve passar; veio %v (%s)", v, msg)
	}
}

// Feature sem cenário visual não é assunto — a maioria das features não tem, e cobrar
// baseline de todas seria inventar dever.
func TestFeatureSemCenarioVisualNaoEhAssunto(t *testing.T) {
	sem := "Funcionalidade: X\n\n  @estado @ABCDX-S01 @nivel-unit\n  Cenário: algo\n"
	if v, _ := checkVRBaseline(sem, mapx.Node{Kind: mapx.KindFeature, ID: "x.feature"}, t.TempDir(), nil, nil); v != Skip {
		t.Errorf("sem cenário visual, Skip; veio %v", v)
	}
}

// A tag que nomeia o regime visual vem do PROJETO (`derived.regimes`), não do framework —
// mesmo princípio de `dialect` e `section_titles`. O de-para é TAG → REGIME, e lê-lo
// invertido fazia o gate não encontrar cenário nenhum, silenciosamente.
func TestTagDoRegimeVisualVemDoProjeto(t *testing.T) {
	if got := tagDeRegimeVisual(nil); got != "nivel-vr" {
		t.Errorf("sem config, o default; veio %q", got)
	}
	cfg := &config.Config{Derived: &config.Derived{Regimes: map[string]string{
		"level-visual": "vr", "level-unit": "unit",
	}}}
	if got := tagDeRegimeVisual(cfg); got != "level-visual" {
		t.Errorf("a tag é a CHAVE do de-para (o que aparece na feature); veio %q", got)
	}
}
