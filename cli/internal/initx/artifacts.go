package initx

import "github.com/co2-lab/anchors/internal/config"

// ArtifactLayer descreve uma camada de artefato oferecível no init: o nome, o
// pattern default e o kind. As perguntas do init são SEMPRE feitas; a inferência
// só decide quais vêm pré-marcadas. Assim um projeto vazio também é configurável
// (o usuário marca o que PRETENDE usar).
type ArtifactLayer struct {
	Name    string
	Kind    string
	Pattern string
}

// KnownArtifactLayers são os tipos de âncora que o init oferece por padrão.
// (Presets de stack — patterns específicos por linguagem — virão depois.)
var KnownArtifactLayers = []ArtifactLayer{
	{Name: "spec", Kind: "spec", Pattern: "**/*.spec.md"},
	{Name: "feature", Kind: "feature", Pattern: "**/*.feature"},
	{Name: "test", Kind: "test", Pattern: "**/*.test.*"},
	{Name: "guide", Kind: "guide", Pattern: "guides/*.md"},
	// plan: a origem do movimento (fase PLANEJAR). Pattern ESPECÍFICO (plans/) para
	// vencer um coringa `**/*.md` da layer 'doc' — senão o plano é classificado como
	// doc e o watcher sugere 'triage' em vez de 'specify'. Ver SuggestNext.
	{Name: "plan", Kind: "plan", Pattern: "plans/*.md"},
}

// DetectedArtifacts devolve os nomes dos artefatos que a inferência achou — usado
// para PRÉ-MARCAR as opções (não para gatear a pergunta).
func (p *Proposal) DetectedArtifacts() map[string]bool {
	m := map[string]bool{}
	if p.HasSpecMD {
		m["spec"] = true
	}
	if p.HasFeature {
		m["feature"] = true
	}
	if p.HasTest {
		m["test"] = true
	}
	if p.GuideDir != "" {
		m["guide"] = true
	}
	if p.PlanDir != "" {
		m["plan"] = true
	}
	return m
}

// ApplyArtifactChoice (re)constrói as camadas de artefato do cfg conforme a escolha
// do usuário — adiciona as escolhidas (mesmo que a inferência não as tenha achado)
// e remove as não escolhidas. Preserva pattern inferido quando havia. `dirs` traz o
// diretório detectado para as layers baseadas em pasta (guide, plan) — quando vazio,
// usa o pattern default. Puro.
func ApplyArtifactChoice(cfg *config.Config, chosen map[string]bool, dirs map[string]string) {
	if cfg.Layers == nil {
		cfg.Layers = map[string]config.Layer{}
	}
	for _, a := range KnownArtifactLayers {
		if chosen[a.Name] {
			if _, exists := cfg.Layers[a.Name]; !exists {
				pattern := a.Pattern
				// guide/plan são baseados em pasta: respeita o dir detectado.
				if dir := dirs[a.Name]; dir != "" && (a.Name == "guide" || a.Name == "plan") {
					pattern = dir + "/*.md"
				}
				cfg.Layers[a.Name] = config.Layer{Pattern: pattern, Kind: a.Kind, Tags: []string{a.Name}}
			}
		} else {
			delete(cfg.Layers, a.Name)
		}
	}
}

// ArtifactNames devolve os nomes conhecidos, em ordem estável (para as opções).
func ArtifactNames() []string {
	out := make([]string, len(KnownArtifactLayers))
	for i, a := range KnownArtifactLayers {
		out[i] = a.Name
	}
	return out
}

// ApplyColocation monta ou remove o `derived` do cfg conforme a escolha do usuário.
// Se use=true, cria os templates de co-location a partir da ÂNCORA, que é a spec.
//
// A spec não aparece entre os derivados — ela é a origem, e da origem nascem o código, a
// feature e o teste. Sem `spec` escolhido não há âncora, e a co-location não se declara:
// um projeto sem spec não tem trinca a localizar.
//
// A extensão do código é declarada, não inferida: `{{ext}}` valia quando o código ancorava
// (a extensão vinha dele), e da spec não há de onde tirá-la. `ext` é o que o projeto
// decidiu no PROJECT.md.
func ApplyColocation(cfg *config.Config, use bool, chosenArtifacts map[string]bool) {
	if !use || !chosenArtifacts["spec"] {
		cfg.Derived = nil
		return
	}
	files := map[string]config.Padroes{}
	if chosenArtifacts["code"] {
		files["code"] = config.Padroes{"{{dir}}/{{name}}.{{ext}}"}
	}
	if chosenArtifacts["feature"] {
		files["feature"] = config.Padroes{"{{dir}}/{{name}}.feature"}
	}
	if chosenArtifacts["test"] {
		files["test"] = config.Padroes{"{{dir}}/{{name}}.test.{{ext}}"}
	}
	if len(files) == 0 {
		cfg.Derived = nil
		return
	}
	cfg.Derived = &config.Derived{Anchor: "spec", Files: files}
}
