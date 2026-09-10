package initx

import (
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// buildConfig monta um config.Config a partir do que a inferência descobriu.
// É a PROPOSTA — o fluxo interativo do init depois confirma/ajusta cada parte.
func (p *Proposal) buildConfig() *config.Config {
	c := &config.Config{
		Version: 1,
		Layers:  map[string]config.Layer{},
	}

	// NB: as layers de ARTEFATO (spec/feature/test/guide) NÃO são criadas aqui —
	// elas vêm da escolha do usuário no init (ApplyArtifactChoice), que é sempre
	// perguntada (mesmo em projeto vazio). A inferência só pré-marca o detectado
	// (ver DetectedArtifacts). Aqui pré-preenchemos só as layers de CÓDIGO e o
	// derived, que servem de default para as perguntas de granularidade.

	// layers de código — uma por diretório-raiz detectado, com tag pelo nome
	excl := []string{"**/*.spec.md", "**/*.feature", "**/*.test.*"}
	for _, dir := range p.CodeDirs {
		layerName := LayerNameFor(dir)
		c.Layers[layerName] = config.Layer{
			Pattern: dir + "/**/*." + globExts(p.CodeExts),
			Kind:    "code",
			Tags:    []string{layerName},
			Exclude: excl,
		}
	}

	// derived: co-location, se detectada
	if p.Colocated {
		files := map[string]config.Padroes{}
		if p.HasSpecMD {
			files["spec"] = config.Padroes{"{{dir}}/{{name}}.spec.md"}
		}
		if p.HasFeature {
			files["feature"] = config.Padroes{"{{dir}}/{{name}}.feature"}
		}
		if p.HasTest {
			files["test"] = config.Padroes{"{{dir}}/{{name}}.test.{{ext}}"}
		}
		c.Derived = &config.Derived{Anchor: "spec", Files: files}
	}

	// test_handle: como o projeto marca elementos para alcance externo. Só é escrito
	// quando a inferência ACHOU o atributo em uso — um backend, uma lib ou um projeto
	// que não marca nada fica sem a chave, e os gates de inventário de handle pulam.
	// Escrever um default aqui (`testID`) faria esses gates acusarem o repositório
	// inteiro de um projeto que nunca prometeu essa convenção.
	if p.TestHandle != "" {
		if c.Derived == nil {
			c.Derived = &config.Derived{Anchor: "spec"}
		}
		c.Derived.TestHandle = p.TestHandle
	}

	// governs fica VAZIO — é a parte semântica, preenchida na P&R (guide↔tag).
	return c
}

// LayerNameFor deriva um nome de layer legível de um diretório (apps/mobile →
// "mobile-code"; packages/backend → "backend-code").
func LayerNameFor(dir string) string {
	base := filepath.Base(dir)
	return base + "-code"
}

// globExts monta um glob de extensões a partir das extensões de código detectadas.
func globExts(exts []string) string {
	if len(exts) == 1 {
		return strings.TrimPrefix(exts[0], ".")
	}
	return "{" + joinExts(exts) + "}"
}

func joinExts(exts []string) string {
	var bare []string
	for _, e := range exts {
		bare = append(bare, strings.TrimPrefix(e, "."))
	}
	return strings.Join(bare, ",")
}
