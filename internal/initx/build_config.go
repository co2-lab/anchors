// @anchors
//   code: BCIBL
//   ref: BLCNB

package initx

import (
	"path"
	"sort"
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
	families := p.Families()
	for _, dir := range p.CodeDirs {
		layerName := names[dir]
		exts, dialect := p.CodeExts, ""
		// Roots written in several languages: each layer reads its own language's files,
		// with its family's dialect.
		if len(families) > 1 {
			dialect = p.familyOf(dir)
			exts = p.extsOf(dialect)
		}
		c.Layers[layerName] = config.Layer{
			Pattern: dirGlob(dir, p.CodeDirs) + globExts(exts),
			Kind:    "code",
			Tags:    []string{layerName},
			Exclude: excl,
			Dialect: dialect,
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
	// language's tests, exports and comments from the first check — one per family when the
	// roots are written in several.
	if len(families) > 1 {
		c.Dialects = map[string]*config.Dialect{}
		for _, fam := range families {
			c.Dialects[fam] = p.familyDialect(fam, p.extsOf(fam))
		}
	} else {
		c.Dialect = p.familyDialect(p.Family, p.CodeExts)
	}

	// governs fica VAZIO — é a parte semântica, preenchida na P&R (guide↔tag).
	return c
}

// hasEntry says whether the comments list holds an entry for the glob.
func hasEntry(list []config.Comments, files string) bool {
	for _, c := range list {
		if c.Files == files {
			return true
		}
	}
	return false
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

// Families are the dialect families of the project's code roots, sorted.
func (p *Proposal) Families() []string {
	seen := map[string]bool{}
	var out []string
	for _, dir := range p.CodeDirs {
		if f := p.DirFamily[dir]; f != "" && !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	sort.Strings(out)
	return out
}

// familyOf is a code root's family; a root no known extension tells takes the project's.
func (p *Proposal) familyOf(dir string) string {
	if f := p.DirFamily[dir]; f != "" {
		return f
	}
	if p.Family != "" {
		return p.Family
	}
	return p.Families()[0]
}

// extsOf are the project's code extensions of a family.
func (p *Proposal) extsOf(fam string) []string {
	var out []string
	for _, e := range p.CodeExts {
		if extFamily[e] == fam {
			out = append(out, e)
		}
	}
	if len(out) == 0 {
		return p.CodeExts
	}
	return out
}

// familyDialect is the dialect a family's files start with: the family's example of every
// lexicon field the gates read, and how the files of those extensions comment. Nil when
// neither the family nor the extensions say anything.
func (p *Proposal) familyDialect(fam string, exts []string) *config.Dialect {
	var d *config.Dialect
	if fam != "" {
		d = &config.Dialect{Family: fam}
		// The language's lexicon, from the family's example: the gates read only what the
		// project declares, and a new project starts with it written down.
		if ex := config.FamilyLexicon(fam); ex != nil {
			*d = *ex
			d.Family = fam
		}
		// What can fail in this stack, written down from the first day: the failure gates
		// ask each fetch for its failure and its handling, on a project that has nothing
		// to unlearn. An existing project is offered them by the governance tips instead.
		d.FalliblePatterns = config.FamilyFalliblePatterns(fam)
		// How the handles the inference found are exposed and queried, from the family's
		// example: the handle gates read only these patterns, and a new project starts
		// with them written down, to adapt to its own helpers.
		d.TestHandles = config.FamilyTestHandles(fam, p.TestHandle)
		// How a navigation call is written, from the family's example: nav-annotated reads
		// only this pattern.
		d.NavigationCall = config.FamilyNavigationCall(fam)
		// How its imports and declarations are written, from the family's example: the chain
		// gates read only these patterns.
		if ex := config.FamilyChain(fam); ex != nil {
			d.Imports, d.ImportNames, d.Declarations = ex.Imports, ex.ImportNames, ex.Declarations
			d.TypeImportPattern = ex.TypeImportPattern
			if d.ImportResolve == nil {
				d.ImportResolve = ex.ImportResolve
			}
		}
	}
	// How the files comment, from the examples covering the extensions found and the
	// family's configuration and scripts: a file no entry covers has no header read.
	files := make([]string, 0, len(exts))
	for _, ext := range exts {
		files = append(files, "x"+ext)
	}
	comments, _ := config.CommentsCovering(files)
	for _, fc := range config.FamilyComments(fam) {
		if !hasEntry(comments, fc.Files) {
			comments = append(comments, fc)
		}
	}
	if len(comments) > 0 {
		if d == nil {
			d = &config.Dialect{}
		}
		d.Comments = comments
	}
	return d
}

// ApplyDialects gives, in a project of several dialects, each family its own test layer —
// `test-<family>`, by the conventions its files follow, naming its dialect — in place of the
// one test layer the artifact choice made, and the colocation's test template to each. A
// layer of tests names the dialect its files are written in, and one layer holds one.
func (p *Proposal) ApplyDialects(cfg *config.Config) {
	if len(cfg.Dialects) < 2 {
		return
	}
	test, ok := cfg.Layers["test"]
	if !ok {
		return
	}
	var tmpl config.Padroes
	if cfg.Derived != nil {
		tmpl = cfg.Derived.Files["test"]
	}
	split := false
	for _, fam := range p.Families() {
		var convs []TestConvention
		for _, c := range p.TestConventions {
			if c.Family == fam {
				convs = append(convs, c)
			}
		}
		if len(convs) == 0 {
			if c, ok := familyDefault[fam]; ok {
				convs = []TestConvention{c}
			}
		}
		if len(convs) == 0 {
			continue
		}
		globs := make([]string, 0, len(convs))
		for _, c := range convs {
			globs = append(globs, c.Glob())
		}
		pattern := globs[0]
		if len(globs) > 1 {
			pattern = "{" + strings.Join(globs, ",") + "}"
		}
		name := "test-" + fam
		cfg.Layers[name] = config.Layer{Pattern: pattern, Kind: test.Kind, Tags: test.Tags, Dialect: fam}
		if tmpl != nil {
			cfg.Derived.Files[name] = config.Padroes{convs[0].Template()}
		}
		split = true
	}
	if split {
		delete(cfg.Layers, "test")
		if tmpl != nil {
			delete(cfg.Derived.Files, "test")
		}
		return
	}
	// No family has a convention to split by: the layer is read by the project's main family,
	// and the doctor names the files of another (GVOPG-B16) — a layer naming none would not load.
	test.Dialect = p.familyOf("")
	cfg.Layers["test"] = test
}
