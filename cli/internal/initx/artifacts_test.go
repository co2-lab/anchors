package initx

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestApplyArtifactChoice_addsAndRemoves(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{
		"spec":        {Kind: "spec", Tags: []string{"spec"}}, // já existe
		"mobile-code": {Kind: "code", Tags: []string{"mobile"}},
	}}
	// usuário quer spec + feature (mas NÃO test/guide)
	ApplyArtifactChoice(cfg, map[string]bool{"spec": true, "feature": true}, nil)

	if _, ok := cfg.Layers["feature"]; !ok {
		t.Error("feature escolhida deveria ter sido adicionada")
	}
	if _, ok := cfg.Layers["spec"]; !ok {
		t.Error("spec escolhida deveria permanecer")
	}
	if _, ok := cfg.Layers["test"]; ok {
		t.Error("test não escolhida não deveria existir")
	}
	// camada de código não é tocada por ApplyArtifactChoice
	if _, ok := cfg.Layers["mobile-code"]; !ok {
		t.Error("camada de código não deveria ser afetada")
	}
}

func TestApplyArtifactChoice_emptyProject(t *testing.T) {
	// projeto vazio: nada detectado, mas o usuário escolhe spec+feature+test.
	cfg := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg, map[string]bool{"spec": true, "feature": true, "test": true}, nil)

	for _, name := range []string{"spec", "feature", "test"} {
		if _, ok := cfg.Layers[name]; !ok {
			t.Errorf("num projeto vazio, a layer escolhida %q deveria ser criada", name)
		}
	}
}

func TestApplyArtifactChoice_guideDir(t *testing.T) {
	cfg := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg, map[string]bool{"guide": true}, map[string]string{"guide": "docs/guides"})
	if got := cfg.Layers["guide"].Pattern; got != "docs/guides/*.md" {
		t.Errorf("pattern do guide = %q, want docs/guides/*.md", got)
	}
}

func TestApplyArtifactChoice_planLayer(t *testing.T) {
	// plan escolhido, sem dir detectado → pattern default plans/*.md (específico, para
	// vencer o coringa doc e o watcher sugerir 'specify' em vez de 'triage').
	cfg := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg, map[string]bool{"plan": true}, nil)
	if got := cfg.Layers["plan"].Pattern; got != "plans/*.md" {
		t.Errorf("pattern do plan = %q, want plans/*.md", got)
	}
	if got := cfg.Layers["plan"].Kind; got != "plan" {
		t.Errorf("kind do plan = %q, want plan", got)
	}
	// com dir detectado → respeita o dir
	cfg2 := &config.Config{Layers: map[string]config.Layer{}}
	ApplyArtifactChoice(cfg2, map[string]bool{"plan": true}, map[string]string{"plan": "docs/plans"})
	if got := cfg2.Layers["plan"].Pattern; got != "docs/plans/*.md" {
		t.Errorf("pattern do plan com dir = %q, want docs/plans/*.md", got)
	}
}

func TestApplyColocation(t *testing.T) {
	cfg := &config.Config{}
	// A ÂNCORA é a spec, e ela não aparece entre os derivados: da spec nascem o código,
	// a feature e o teste. Com os quatro escolhidos, sobram TRÊS derivados.
	ApplyColocation(cfg, true, map[string]bool{
		"spec": true, "code": true, "feature": true, "test": true,
	})
	if cfg.Derived == nil || len(cfg.Derived.Files) != 3 {
		t.Fatalf("esperava derived com 3 arquivos, got %+v", cfg.Derived)
	}
	if cfg.Derived.Anchor != "spec" {
		t.Errorf("a âncora é a spec, veio %q", cfg.Derived.Anchor)
	}
	if _, ehDerivada := cfg.Derived.Files["spec"]; ehDerivada {
		t.Error("a spec é a ORIGEM — listá-la como derivada inverte a doutrina")
	}

	// Sem spec escolhida não há âncora, e a co-location não se declara: um projeto sem
	// spec não tem trinca a localizar.
	semSpec := &config.Config{}
	ApplyColocation(semSpec, true, map[string]bool{"code": true, "test": true})
	if semSpec.Derived != nil {
		t.Error("sem spec não há âncora — não deveria haver derived")
	}
	// desligar co-location → derived some
	ApplyColocation(cfg, false, map[string]bool{"spec": true})
	if cfg.Derived != nil {
		t.Error("co-location desligada deveria zerar o derived")
	}
	// co-location ligada mas sem artefatos co-localizáveis (só guide) → sem derived
	ApplyColocation(cfg, true, map[string]bool{"guide": true})
	if cfg.Derived != nil {
		t.Error("sem artefatos co-localizáveis, derived deveria ficar nil")
	}
}
