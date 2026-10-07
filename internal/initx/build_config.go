// @anchors
//   code: BCIBL
//   ref: BLCNB

package initx

import (
	"path"
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
	// The tests are excluded by the convention the project's own test files follow — a
	// `_test.go` read as code by a fixed `**/*.test.*` was the defect.
	excl := append([]string{"**/*.spec.md", "**/*.feature"}, p.testExcludes()...)
	names := layerNames(p.CodeDirs)
	for _, dir := range p.CodeDirs {
		layerName := names[dir]
		c.Layers[layerName] = config.Layer{
			Pattern: dirGlob(dir, p.CodeDirs) + globExts(p.CodeExts),
			Kind:    "code",
			Tags:    []string{layerName},
			Exclude: excl,
		}
	}

	// derived: co-location, se detectada
	if p.Colocated {
		// The spec is the ANCHOR and never one of its own derivatives (ARCHR-B07): listing
		// it here proposed a config that derived the spec from itself.
		files := map[string]config.Padroes{}
		if p.HasFeature {
			files["feature"] = config.Padroes{"{{dir}}/{{name}}.feature"}
		}
		if p.HasTest {
			t := p.TestTemplate()
			if t == "" {
				t = "{{dir}}/{{name}}.test.{{ext}}" // no convention known: the generic form
			}
			files["test"] = config.Padroes{t}
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

	// The dialect: the family the manifest or the code says, so the gates read this
	// language's tests, exports and comments from the first check.
	if p.Family != "" {
		c.Dialect = &config.Dialect{Family: p.Family}
		// What can fail in this stack, written down from the first day: the failure gates
		// ask each fetch for its failure and its handling, on a project that has nothing
		// to unlearn. An existing project is offered them by the governance tips instead.
		c.Dialect.FalliblePatterns = config.FamilyFalliblePatterns(p.Family)
	}

	// governs fica VAZIO — é a parte semântica, preenchida na P&R (guide↔tag).
	return c
}

// LayerNameFor deriva um nome de layer legível de um diretório (apps/mobile →
// "mobile-code"; packages/backend → "backend-code"; the root → "root-code").
func LayerNameFor(dir string) string {
	if dir == "." || dir == "" {
		return "root-code"
	}
	return path.Base(dir) + "-code"
}

// layerNames names each code folder's layer by its last segment, and by its whole path
// when two folders share that segment (src/handlers and lib/handlers): one name for both
// made the second layer overwrite the first.
func layerNames(dirs []string) map[string]string {
	count := map[string]int{}
	for _, d := range dirs {
		count[LayerNameFor(d)]++
	}
	out := map[string]string{}
	for _, d := range dirs {
		n := LayerNameFor(d)
		if count[n] > 1 {
			n = strings.ReplaceAll(d, "/", "-") + "-code"
		}
		out[d] = n
	}
	return out
}

// dirGlob is the part of a code layer's pattern before the extensions: every file under
// the folder, or only its own files when another code folder sits beneath it, so no file
// falls in two layers.
func dirGlob(dir string, dirs []string) string {
	prefix := dir + "/"
	if dir == "." {
		prefix = ""
	}
	for _, o := range dirs {
		if o != dir && (prefix == "" || strings.HasPrefix(o, prefix)) {
			return prefix + "*."
		}
	}
	return prefix + "**/*."
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

// TestPattern is the pattern of the project's test layer: its conventions, joined; the
// generic `**/*.test.*` only when no convention is known.
func (p *Proposal) TestPattern() string {
	gs := p.TestGlobs()
	switch len(gs) {
	case 0:
		return "**/*.test.*"
	case 1:
		return gs[0]
	}
	return "{" + strings.Join(gs, ",") + "}"
}

// testExcludes are the test globs a code layer leaves out.
func (p *Proposal) testExcludes() []string {
	if gs := p.TestGlobs(); len(gs) > 0 {
		return gs
	}
	return []string{"**/*.test.*"}
}
