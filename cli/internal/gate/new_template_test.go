package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// Este teste é o META-GATE dos templates de `anchors new`: garante que o esqueleto
// que o comando emite NASCE CONFORME — passa nas MESMAS funções de gate que o
// `check` roda (checkHeaderConforme, checkSpecSections). Sem ele, o template poderia
// divergir da régua e ninguém perceberia (o template não está no grafo do projeto).
//
// As strings abaixo são a saída canônica de `new` (default) por kind — mantê-las em
// sincronia com cmd/anchors/new_templates.go é o contrato. Se o gate mudar de régua,
// este teste quebra e força atualizar o template junto.

func headerNode(kind mapx.Kind) mapx.Node {
	// nó regido genérico (não é camada reconhecida) → o header exige code|ref.
	return mapx.Node{ID: "x/Login." + string(kind), Kind: kind}
}

func TestNewTemplate_specIsBornConforming(t *testing.T) {
	// espelha specTemplate (default: title+overview+rules) renderizado p/ "Login"/"LGNOX".
	spec := "<!-- @anchors\n  code: LGNOX\n  updated_at: TODO\n  layer: TODO\n-->\n" +
		"# Login — TODO propósito em uma frase\n\n> **Código**: `LGNOX`\n\n" +
		"## Visão Geral\nTODO: o que a unidade faz e para quem.\n\n" +
		"## Regras\n\n### LGNOX-B01 — TODO regra\nDescreva o comportamento (não a implementação).\n\n"

	if v, msg := checkHeaderConforme(spec, mapx.Node{ID: "x/Login.spec.md", Kind: "spec"}); v != Pass {
		t.Fatalf("spec do `new` reprova header-conforme: %s", msg)
	}
	if v, msg := checkSpecSections(spec, mapx.Node{ID: "x/Login.spec.md"}); v != Pass {
		t.Fatalf("spec do `new` reprova spec-completa: %s", msg)
	}
}

func TestNewTemplate_featureIsBornConforming(t *testing.T) {
	feat := "# language: pt\n# @anchors\n#   ref: LGNOX\n#   updated_at: TODO\n#   layer: feature\n" +
		"\n@LGNOX\nFuncionalidade: Login\n\n" +
		"  @LGNOX-B01 @nivel-unit @P2\n  Cenário: TODO\n    Dado TODO\n    Quando TODO\n    Então o efeito LGNOX-B01 se verifica\n\n"

	if v, msg := checkHeaderConforme(feat, headerNode(mapx.KindFeature)); v != Pass {
		t.Fatalf("feature do `new` reprova header-conforme: %s", msg)
	}
	// non-empty: a feature tem conteúdo além do header.
	if strings.TrimSpace(strings.SplitN(feat, "feature\n", 2)[1]) == "" {
		t.Fatal("feature do `new` está vazia após o header")
	}
}

func TestNewTemplate_testIsBornConforming(t *testing.T) {
	test := "// @anchors\n//   ref: LGNOX\n//   updated_at: TODO\n//   layer: test\n" +
		"\ndescribe('Login', () => {\n  it('[LGNOX-B01] TODO', () => {\n    // TODO\n  })\n})\n"

	if v, msg := checkHeaderConforme(test, headerNode(mapx.KindTest)); v != Pass {
		t.Fatalf("test do `new` reprova header-conforme: %s", msg)
	}
}
