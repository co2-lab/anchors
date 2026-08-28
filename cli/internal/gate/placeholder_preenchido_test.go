package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// TestEsqueletoCruEhReprovado guarda o falso verde mais caro medido: uma spec recém-saída
// do `anchors new`, com `layer: TODO` e todas as regras em TODO, atravessava TODOS os
// gates bloqueantes com "✓ pode promover" — inclusive `header-conforme`, que leu
// `layer: TODO` e aprovou. Os gates de header validam a FORMA, e `TODO` é bem-formado.
//
// E o `anchors work spec` prometia, textualmente, "placeholder não preenchido reprova".
func TestEsqueletoCruEhReprovado(t *testing.T) {
	cru := `<!-- @anchors
  code: CRUXX
  layer: TODO
  updated_at: TODO
-->
# Cru — TODO propósito em uma frase

## Visão Geral
TODO: o que a unidade faz e para quem.

## Regras
| Regra | Efeito |
| --- | --- |
| ` + "`CRUXX-B01`" + ` | TODO |
`
	v, msg := checkPlaceholderPreenchido(cru, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
	if v != Fail {
		t.Fatalf("esqueleto cru deve reprovar; veio %v (%s)", v, msg)
	}
	for _, quer := range []string{"layer: TODO", "updated_at: TODO"} {
		if !strings.Contains(msg, quer) {
			t.Errorf("a mensagem precisa nomear %q — é o campo que o gate de header aprovava", quer)
		}
	}
}

// O falso positivo que derrubaria o gate: uma seção `## TODOs` é uma lista de pendências
// que o autor escreveu DE PROPÓSITO. Medido em 77 specs de um projeto real. Confundi-la
// com "o autor não escreveu nada" acusaria a maioria — e um gate que acusa a maioria é
// desligado no primeiro dia, levando junto os que funcionam.
func TestSecaoDeTodosLegitimaNaoEhAcusada(t *testing.T) {
	spec := `<!-- @anchors
  code: REALX
  layer: business-logic
  updated_at: 2026-08-13
-->
# Real — calcula o saldo

## Regras
| Regra | Efeito |
| --- | --- |
| ` + "`REALX-B01`" + ` | soma as entradas do mês |

## TODOs
- migrar para o novo formato de data
- TODO: avaliar cache
`
	if v, msg := checkPlaceholderPreenchido(spec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil); v != Pass {
		t.Errorf("seção de pendências do autor é legítima; veio %v (%s)", v, msg)
	}
}
