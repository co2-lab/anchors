package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// ref-resolves: o `ref:` de um artefato aponta para a spec que REALMENTE o descreve.
//
// `header-conforme` verifica que o campo EXISTE. Ninguém verificava que ele aponta para o
// lugar certo — e um `ref:` errado é pior que ausente: ele parece rastreabilidade, o gate
// fica verde, e a unidade inteira é atribuída à spec errada. Todo gate relacional que
// dependa dessa aresta passa a confrontar o par errado, em silêncio.
//
// O modo de falha característico é a REFATORAÇÃO que ninguém propagou. Medido num projeto
// real: 49 arquivos de modelo com `ref: DTAX` — a identidade de quando os modelos viviam
// todos num arquivo só. Depois da desfusão, cada um ganhou spec própria (`AUA1`, `BAB1`,
// `BCB1`…), e nenhum `ref:` foi atualizado. Os 49 continuaram apontando para o schema
// inteiro, e nada acusou.
//
// A régua: se existe uma spec IRMÃ (a que a Estrutura co-loca com este arquivo), o `ref:`
// tem de ser o `code:` dela. Sem spec irmã, o gate se cala — quem cobra a ausência da
// peça é o `trinca-completa`, e dois gates acusando o mesmo defeito viram ruído.
func checkRefResolves(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	switch n.Kind {
	case mapx.KindCode, mapx.KindFeature, mapx.KindTest:
	default:
		return Skip, i18n.T("gate.ref_resolves.skip_not_referencing")
	}

	m := refHeaderRE().FindStringSubmatch(content)
	if m == nil {
		return Skip, i18n.T("gate.ref_resolves.skip_no_ref")
	}
	ref := m[1]

	specPath, specCode := siblingSpecOf(root, n.ID)
	if specCode == "" {
		// With no sibling spec there is nothing to confront the `ref:` with — BUT it can
		// still be required to point at SOME existing identity. An invented `ref:` (the
		// code is no spec's `code:` anywhere in the project) is not "cannot tell": it is
		// wrong, and without this check it passed as undetermined.
		//
		// Measured in MIF: 3174 refs, 1830 confronted and 1344 with NO sibling spec — that
		// is, 42% of the corpus fell into Skip. A hand-made `ref: KYBD`, pointing at a code
		// that exists nowhere, was measured landing on `~1` and not on `✗1`.
		//
		// The distinction from `triad-complete` holds: that one charges the ABSENCE of the
		// sibling spec; this one, the identity the `ref:` INVENTS. An infra file with no
		// spec (legitimate) still passes — what does not pass is citing a phantom code.
		if g != nil && !codeExistsInGraph(g, ref) {
			return Fail, fmt.Sprintf(i18n.T("gate.ref_resolves.unknown_code"), ref, ref)
		}
		return Skip, i18n.T("gate.ref_resolves.skip_no_sibling_spec")
	}
	if ref == specCode {
		return Pass, ""
	}
	return Fail, fmt.Sprintf(i18n.T("gate.ref_resolves.unmatched_ref"), ref, specPath, specCode)
}

// codeExistsInGraph says whether some node of the map DECLARES this code as its identity.
//
// Only `CodeDeclarado` counts: an INFERRED identity (the first code that showed up in a
// fixture's text, for instance) owns nothing, and accepting it here would let through the
// ref that points at an example string in documentation.
func codeExistsInGraph(g *mapx.Graph, code string) bool {
	for i := range g.Nodes {
		if g.Nodes[i].CodeDeclarado && g.Nodes[i].Code == code {
			return true
		}
	}
	return false
}

// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func refHeaderRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*ref:\s*([A-Z0-9]` + config.CodeLengthPattern() + `)\b`)
}

// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func specCodeRE() *regexp.Regexp {
	return regexp.MustCompile(`(?m)^\s*(?://|#|<!--|\*)?\s*code:\s*([A-Z0-9]` + config.CodeLengthPattern() + `)\b`)
}

// siblingSpecOf acha a spec co-localizada com o arquivo e devolve (caminho, código dela).
// Segue a convenção de nome: mesmo tronco, sufixo `.spec.md`, mesmo diretório.
func siblingSpecOf(root, rel string) (string, string) {
	dir := filepath.Dir(rel)
	base := filepath.Base(rel)
	for _, suf := range []string{".feature", ".test.ts", ".test.tsx", ".spec.ts", "_test.go", "_test.py", "_spec.rb"} {
		base = strings.TrimSuffix(base, suf)
	}
	if i := strings.Index(base, "."); i > 0 {
		base = base[:i]
	}
	cand := filepath.Join(dir, base+".spec.md")
	b, err := os.ReadFile(filepath.Join(root, cand))
	if err != nil {
		return "", ""
	}
	if m := specCodeRE().FindStringSubmatch(string(b)); m != nil {
		return cand, m[1]
	}
	return "", ""
}
