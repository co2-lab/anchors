package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

const featVR = `Funcionalidade: Tela

  @estado @TCDTX-S01 @nivel-integration
  Cenário: carrega
    Dado algo

  @estado @TCDTX-VR @vr-level @P3
  Cenário: aparência da tela carregada
    Dado algo
`

func TestVR_NotFeatureReturnsSkip(t *testing.T) {
	t.Run("VRBSV-B01: Artifacts that are not feature files leave with verdict Skip", func(t *testing.T) {})
	v, _ := checkVRBaseline("some content", mapx.Node{Kind: mapx.KindCode, ID: "tela/T.go"}, t.TempDir(), nil, nil)
	if v != Skip {
		t.Fatalf("expected Skip for non-feature, got %v", v)
	}
}

// Feature sem cenário visual não é assunto — a maioria das features não tem, e cobrar
// baseline de todas seria inventar dever.
func TestFeatureSemCenarioVisualNaoEhAssunto(t *testing.T) {
	t.Run("VRBSV-B02: Feature files declaring no visual regression scenarios leave with verdict Skip", func(t *testing.T) {})
	sem := "Funcionalidade: X\n\n  @estado @ABCDX-S01 @nivel-unit\n  Cenário: algo\n"
	if v, _ := checkVRBaseline(sem, mapx.Node{Kind: mapx.KindFeature, ID: "x.feature"}, t.TempDir(), nil, nil); v != Skip {
		t.Errorf("sem cenário visual, Skip; veio %v", v)
	}
}

func TestVR_MatchingBaselinePasses(t *testing.T) {
	t.Run("VRBSV-B03: Visual regression scenarios having matching baseline images pass", func(t *testing.T) {})
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, "tela"), 0o755))
	must(t, os.WriteFile(filepath.Join(root, "tela", "T.TCDTX-VR.png"), []byte("png"), 0o644))

	v, msg := checkVRBaseline(featVR, mapx.Node{Kind: mapx.KindFeature, ID: "tela/T.feature"}, root, nil, nil)
	if v != Pass {
		t.Fatalf("expected Pass with baseline present, got %v (%s)", v, msg)
	}
}

// TestCenarioVRSemBaselineEhAcusado guarda a superfície de prova que nenhum gate
// alcançava. `@vr-level` declara que a tela é provada por CAPTURA, não por asserção — e
// sem imagem de referência não há contra o que comparar. O cenário existe, o gate de
// feature o conta como coberto, e a prova prometida não acontece.
func TestCenarioVRSemBaselineEhAcusado(t *testing.T) {
	t.Run("VRBSV-B04: Visual regression scenarios lacking baseline images fail naming missing codes", func(t *testing.T) {})
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, "tela"), 0o755))

	v, msg := checkVRBaseline(featVR, mapx.Node{Kind: mapx.KindFeature, ID: "tela/T.feature"}, root, nil, nil)
	if v != Fail {
		t.Fatalf("cenário VR sem baseline deve reprovar; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "TCDTX-VR") {
		t.Errorf("a mensagem precisa nomear o cenário sem imagem; veio: %s", msg)
	}
}

func TestVR_VariantSuffixAccepted(t *testing.T) {
	t.Run("VRBSV-B05: Baseline image matching supports naming variants", func(t *testing.T) {})
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, "tela"), 0o755))
	// A convenção aceita variante no nome: `TCDTX-VR-loaded.png` prova `TCDTX-VR`.
	must(t, os.WriteFile(filepath.Join(root, "tela", "T.TCDTX-VR-loaded.png"), []byte("png"), 0o644))
	if v, msg := checkVRBaseline(featVR, mapx.Node{Kind: mapx.KindFeature, ID: "tela/T.feature"}, root, nil, nil); v != Pass {
		t.Errorf("com baseline (mesmo com variante no nome) deve passar; veio %v (%s)", v, msg)
	}
}

// A tag que nomeia o regime visual vem do PROJETO (`derived.regimes`), não do framework —
// mesmo princípio de `dialect` e `section_titles`. O de-para é TAG → REGIME, e lê-lo
// invertido fazia o gate não encontrar cenário nenhum, silenciosamente.
func TestTagDoRegimeVisualVemDoProjeto(t *testing.T) {
	t.Run("VRBSV-B06: Reads visual regime tag from project configuration", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{Regimes: map[string]string{
		"level-visual": "vr", "level-unit": "unit",
	}}}
	if got := visualRegimeTag(cfg); got != "level-visual" {
		t.Errorf("a tag é a CHAVE do de-para (o que aparece na feature); veio %q", got)
	}
}

func TestVR_DefaultVisualRegimeTagFallback(t *testing.T) {
	t.Run("VRBSV-B07: Falls back to default visual regime tag when configuration is missing", func(t *testing.T) {})
	if got := visualRegimeTag(nil); got != "vr-level" {
		t.Errorf("sem config, esperado default vr-level; veio %q", got)
	}
}

func TestVR_MissingCodesSortedAlphabetically(t *testing.T) {
	t.Run("VRBSV-B08: Multiple missing baseline scenario codes are sorted alphabetically", func(t *testing.T) {})
	featMulti := `Feature: Multi
  @ZEDTX-VR @vr-level
  Scenario: z
  @ALFTX-VR @vr-level
  Scenario: a
`
	v, msg := checkVRBaseline(featMulti, mapx.Node{Kind: mapx.KindFeature, ID: "m.feature"}, t.TempDir(), nil, nil)
	if v != Fail {
		t.Fatalf("expected Fail, got %v", v)
	}
	idxA := strings.Index(msg, "ALFTX-VR")
	idxZ := strings.Index(msg, "ZEDTX-VR")
	if idxA == -1 || idxZ == -1 || idxA >= idxZ {
		t.Fatalf("expected sorted ALFTX-VR before ZEDTX-VR, got %q", msg)
	}
}

func TestVR_OnlyMatchesVisualSuffix(t *testing.T) {
	t.Run("VRBSV-B09: Only scenario codes containing visual regression suffix are matched", func(t *testing.T) {})
	featCoTagged := `Feature: CoTagged
  @TCDTX-S01 @TCDTX-VR @vr-level
  Scenario: co-tagged
`
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "c.TCDTX-VR.png"), []byte("png"), 0o644))
	v, msg := checkVRBaseline(featCoTagged, mapx.Node{Kind: mapx.KindFeature, ID: "c.feature"}, root, nil, nil)
	if v != Pass {
		t.Fatalf("expected Pass because only -VR is checked for image, got %v (%s)", v, msg)
	}
}

func TestVR_InvariantEveryVisualScenarioDemandsImage(t *testing.T) {
	t.Run("VRBSV-I01: Every visual regression scenario declared must correspond to a baseline image", func(t *testing.T) {})
	featDual := `Feature: Dual
  @TCDTX-VR @vr-level
  Scenario: one
  @TCDTY-VR @vr-level
  Scenario: two
`
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "d.TCDTX-VR.png"), []byte("png"), 0o644))
	// Only one of two has baseline image -> must Fail
	v, msg := checkVRBaseline(featDual, mapx.Node{Kind: mapx.KindFeature, ID: "d.feature"}, root, nil, nil)
	if v != Fail || !strings.Contains(msg, "TCDTY-VR") {
		t.Fatalf("expected Fail naming missing TCDTY-VR, got %v (%s)", v, msg)
	}
}

func TestVR_InvariantTagMappingResolution(t *testing.T) {
	t.Run("VRBSV-I02: Visual regime tag mapping treats map key as tag in feature", func(t *testing.T) {})
	cfg := &config.Config{Derived: &config.Derived{Regimes: map[string]string{
		"@custom-visual": "visual-regression",
	}}}
	tag := visualRegimeTag(cfg)
	if tag != "custom-visual" {
		t.Fatalf("expected custom-visual, got %q", tag)
	}
}

func TestVR_ConstraintDoesNotCheckMtime(t *testing.T) {
	t.Run("VRBSV-X01: Does not evaluate baseline staleness using disk modification timestamps", func(t *testing.T) {})
	root := t.TempDir()
	imgPath := filepath.Join(root, "t.TCDTX-VR.png")
	must(t, os.WriteFile(imgPath, []byte("png"), 0o644))
	// Set mtime far in the past (e.g. 5 years ago)
	oldTime := time.Now().Add(-5 * 365 * 24 * time.Hour)
	_ = os.Chtimes(imgPath, oldTime, oldTime)

	feat := `Feature: T
  @TCDTX-VR @vr-level
  Scenario: t
`
	v, msg := checkVRBaseline(feat, mapx.Node{Kind: mapx.KindFeature, ID: "t.feature"}, root, nil, nil)
	if v != Pass {
		t.Fatalf("expected Pass regardless of old mtime, got %v (%s)", v, msg)
	}
}

func TestVR_ConstraintDoesNotCheckGitCommitDates(t *testing.T) {
	t.Run("VRBSV-X02: Does not fail commits based on git commit dates of baseline images", func(t *testing.T) {})
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "g.TCDTX-VR.png"), []byte("png"), 0o644))
	feat := `Feature: G
  @TCDTX-VR @vr-level
  Scenario: g
`
	v, msg := checkVRBaseline(feat, mapx.Node{Kind: mapx.KindFeature, ID: "g.feature"}, root, nil, nil)
	if v != Pass {
		t.Fatalf("expected Pass without consulting git history, got %v (%s)", v, msg)
	}
}
