package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
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
		return Skip, "só quem REFERENCIA carrega `ref:` — a spec é dona (`code:`)"
	}

	m := refHeaderRE().FindStringSubmatch(content)
	if m == nil {
		return Skip, "o artefato não declara `ref:` — a ausência é do gate header-conforme"
	}
	ref := m[1]

	specPath, specCode := specIrmaDe(root, n.ID)
	if specCode == "" {
		return Skip, "sem spec irmã para confrontar — a ausência da peça é do gate trinca-completa"
	}
	if ref == specCode {
		return Pass, ""
	}
	return Fail, fmt.Sprintf("`ref: %s` não é a identidade da spec que descreve este arquivo — "+
		"`%s` declara `code: %s`. Um `ref:` errado PARECE rastreabilidade: o gate de header "+
		"fica verde, e a unidade inteira passa a ser atribuída à spec errada, com todo gate "+
		"relacional confrontando o par errado. Costuma ser refatoração não propagada (a spec "+
		"foi dividida, o `ref:` ficou apontando para a antiga)",
		ref, specPath, specCode)
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

// specIrmaDe acha a spec co-localizada com o arquivo e devolve (caminho, código dela).
// Segue a convenção de nome: mesmo tronco, sufixo `.spec.md`, mesmo diretório.
func specIrmaDe(root, rel string) (string, string) {
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
