package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// PLCFL-B01: A raw skeleton fails confrontation
//
// TestEsqueletoCruEhReprovado guarda o falso verde mais caro medido: uma spec recém-saída
// do `anchors new`, com `layer: TODO` e todas as regras em TODO, atravessava TODOS os
// gates bloqueantes com "✓ pode promover" — inclusive `header-conforme`, que leu
// `layer: TODO` e aprovou. Os gates de header validam a FORMA, e `TODO` é bem-formado.
//
// E o `anchors work spec` prometia, textualmente, "placeholder não preenchido reprova".
func TestEsqueletoCruEhReprovado(t *testing.T) {
	t.Run("PLCFL-B01: A raw skeleton fails confrontation", func(t *testing.T) {
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
		v, msg := checkPlaceholderFilled(cru, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if v != Fail {
			t.Fatalf("esqueleto cru deve reprovar; veio %v (%s)", v, msg)
		}
		for _, quer := range []string{"layer: TODO", "updated_at: TODO"} {
			if !strings.Contains(msg, quer) {
				t.Errorf("a mensagem precisa nomear %q — é o campo que o gate de header aprovava", quer)
			}
		}
	})
}

// PLCFL-I01: A section written on purpose to list pending work is not accused
//
// O falso positivo que derrubaria o gate: uma seção `## TODOs` é uma lista de pendências
// que o autor escreveu DE PROPÓSITO. Medido em 77 specs de um projeto real. Confundi-la
// com "o autor não escreveu nada" acusaria a maioria — e um gate que acusa a maioria é
// desligado no primeiro dia, levando junto os que funcionam.
func TestSecaoDeTodosLegitimaNaoEhAcusada(t *testing.T) {
	t.Run("PLCFL-I01: A section written on purpose to list pending work is not accused", func(t *testing.T) {
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
		if v, msg := checkPlaceholderFilled(spec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil); v != Pass {
			t.Errorf("seção de pendências do autor é legítima; veio %v (%s)", v, msg)
		}
	})
}

func TestPlaceholderFilled(t *testing.T) {
	TestEsqueletoCruEhReprovado(t)
	TestSecaoDeTodosLegitimaNaoEhAcusada(t)

	t.Run("PLCFL-B02: A header field whose value is the marker fails", func(t *testing.T) {
		headerMarkerSpec := `<!-- @anchors
  code: PHVAL
  layer: TODO
  updated_at: 2026-09-20
-->
# Placeholder Value — legitimate purpose

## Rules
| Rule | Effect |
| --- | --- |
| ` + "`PHVAL-B01`" + ` | legitimate effect |
`
		v, msg := checkPlaceholderFilled(headerMarkerSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if v != Fail {
			t.Fatalf("spec with header placeholder must fail; got %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "layer: TODO") {
			t.Errorf("expected failure message to name 'layer: TODO'; got %s", msg)
		}

	})

	t.Run("PLCFL-B03: A table cell holding only the marker fails", func(t *testing.T) {
		cellMarkerSpec := `<!-- @anchors
  code: PHVAL
  layer: gate
  updated_at: 2026-09-20
-->
# Placeholder Value — legitimate purpose

## Rules
| Rule | Effect |
| --- | --- |
| ` + "`PHVAL-B01`" + ` | TODO |
`
		v, msg := checkPlaceholderFilled(cellMarkerSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if v != Fail {
			t.Fatalf("spec with placeholder table cell must fail; got %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "TODO") || !strings.Contains(msg, "PHVAL-B01") {
			t.Errorf("expected failure message to identify the empty rule cell; got %s", msg)
		}
	})

	t.Run("PLCFL-B04: A title or body line opening with the marker fails", func(t *testing.T) {
		titleMarkerSpec := `<!-- @anchors
  code: PHVAL
  layer: gate
  updated_at: 2026-09-20
-->
# Placeholder Value — TODO purpose in one sentence

## Rules
| Rule | Effect |
| --- | --- |
| ` + "`PHVAL-B01`" + ` | legitimate effect |
`
		vTitle, msgTitle := checkPlaceholderFilled(titleMarkerSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if vTitle != Fail {
			t.Fatalf("spec with placeholder in title must fail; got %v (%s)", vTitle, msgTitle)
		}
		if !strings.Contains(msgTitle, "# Placeholder Value — TODO purpose in one sentence") {
			t.Errorf("expected failure message to name the placeholder title; got %s", msgTitle)
		}

		bodyMarkerSpec := `<!-- @anchors
  code: PHVAL
  layer: gate
  updated_at: 2026-09-20
-->
# Placeholder Value — legitimate purpose

## Overview
TODO: what the unit does and for whom.
`
		vBody, msgBody := checkPlaceholderFilled(bodyMarkerSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if vBody != Fail {
			t.Fatalf("spec with placeholder in body line must fail; got %v (%s)", vBody, msgBody)
		}
		if !strings.Contains(msgBody, "TODO: what the unit does and for whom.") {
			t.Errorf("expected failure message to name the placeholder body line; got %s", msgBody)
		}
	})

	t.Run("PLCFL-B05: The verdict names what was left behind", func(t *testing.T) {
		multiSpec := `<!-- @anchors
  code: PHVAL
  layer: TODO
  updated_at: 2026-09-20
-->
# Placeholder Value — legitimate purpose

## Rules
| Rule | Effect |
| --- | --- |
| ` + "`PHVAL-B01`" + ` | TODO |

TODO: this is a very long placeholder line designed to verify that truncation at sixty characters works properly when reporting leftovers
`
		v, msg := checkPlaceholderFilled(multiSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if v != Fail {
			t.Fatalf("expected Fail; got %v (%s)", v, msg)
		}
		if !strings.Contains(msg, "«layer: TODO»") {
			t.Errorf("verdict must explicitly name «layer: TODO»; got %s", msg)
		}
		if !strings.Contains(msg, "«| `PHVAL-B01` | TODO |»") {
			t.Errorf("verdict must explicitly name table cell; got %s", msg)
		}
		if !strings.Contains(msg, "…»") {
			t.Errorf("verdict must truncate long leftovers with ellipsis; got %s", msg)
		}
	})

	t.Run("PLCFL-B06: An artifact with every marker replaced passes", func(t *testing.T) {
		cleanSpec := `<!-- @anchors
  code: PHVAL
  layer: gate
  updated_at: 2026-09-20
-->
# Placeholder Value — legitimate purpose

## Overview
All sections and descriptions are genuinely written.

## Rules
| Rule | Effect |
| --- | --- |
| ` + "`PHVAL-B01`" + ` | legitimate effect |
`
		vSpec, msgSpec := checkPlaceholderFilled(cleanSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if vSpec != Pass || msgSpec != "" {
			t.Fatalf("clean spec must Pass; got %v (%s)", vSpec, msgSpec)
		}

		cleanFeature := `# language: en
@PHVAL
Feature: CleanFeature — every scenario is completed

  @PHVAL-B01 @unit-level
  Scenario: Clean scenario
    Given a valid input
    When processed
    Then it succeeds
`
		vFeat, msgFeat := checkPlaceholderFilled(cleanFeature, mapx.Node{Kind: mapx.KindFeature}, "", nil, nil)
		if vFeat != Pass || msgFeat != "" {
			t.Fatalf("clean feature must Pass; got %v (%s)", vFeat, msgFeat)
		}
	})

	t.Run("PLCFL-X01: The gate does not judge the quality of replacement text", func(t *testing.T) {
		poorQualitySpec := `<!-- @anchors
  code: JUNKX
  layer: gate
  updated_at: 2026-09-20
-->
# Junk — asdf foo bar

## Overview
qwerty minimal meaningless text that replaced the template

## Rules
| Rule | Effect |
| --- | --- |
| ` + "`JUNKX-B01`" + ` | whatever |
`
		v, msg := checkPlaceholderFilled(poorQualitySpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if v != Pass {
			t.Fatalf("gate must not judge prose quality; got %v (%s)", v, msg)
		}
	})

	// PLCFL-X02: the marker itself, in running prose, is NOT the evidence — the POSITION
	// is. This spec uses the real marker the gate knows, deliberately placed outside every
	// value position: every header field, table cell and rule title is written.
	//
	// The earlier version of this test used invented words (WIP, TBD) and passed for the
	// wrong reason — the gate does not know those words at all, so it proved only that
	// unknown words are ignored, never the distinction the rule actually claims.
	t.Run("PLCFL-X02: a marker in running prose is not accused", func(t *testing.T) {
		proseSpec := `<!-- @anchors
  code: CUSTX
  layer: gate
  updated_at: 2026-09-20
-->
# Custom — a real purpose written out in full

## Overview
The author still has TODO items to weigh here, and says so on purpose.

## TODOs
- decide whether the vocabulary becomes configurable

## Rules
| Rule | Effect |
| --- | --- |
| ` + "`CUSTX-B01`" + ` | a written effect |
`
		v, msg := checkPlaceholderFilled(proseSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil)
		if v != Pass {
			t.Fatalf("a marker outside a value position must not be accused; got %v (%s)", v, msg)
		}
	})

	t.Run("skips non-applicable artifact kind", func(t *testing.T) {
		v, _ := checkPlaceholderFilled("layer: TODO", mapx.Node{Kind: mapx.KindCode}, "", nil, nil)
		if v != Skip {
			t.Fatalf("expected Skip for code artifact; got %v", v)
		}
	})
}

// The vocabulary is the PROJECT's (`placeholder_markers`), defaulting to the only word the
// `anchors new` templates write. Decided in PLCFL-Q01.
func TestPlaceholderVocabulary(t *testing.T) {
	fixmeSpec := `<!-- @anchors
  code: PHVAL
  layer: gate
  updated_at: FIXME
-->
# Placeholder Value — legitimate purpose
`
	t.Run("PLCFL-B07: The marker vocabulary is the project's, and TODO by default", func(t *testing.T) {})
	// FIXME is no word the templates write: by default it is an ordinary value.
	if v, msg := checkPlaceholderFilled(fixmeSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, nil); v != Pass {
		t.Errorf("with the default vocabulary FIXME is a value, got %v (%s)", v, msg)
	}
	// Declared by the project, it is a marker in every position.
	cfg := &config.Config{PlaceholderMarkers: []string{"FIXME"}}
	vFixme, msgFixme := checkPlaceholderFilled(fixmeSpec, mapx.Node{Kind: mapx.KindSpec}, "", nil, cfg)
	if vFixme != Fail || !strings.Contains(msgFixme, "updated_at: FIXME") {
		t.Errorf("expected Fail naming 'updated_at: FIXME'; got %v (%s)", vFixme, msgFixme)
	}
	cell := "| `PHVAL-B01` | FIXME |\n"
	if v, _ := checkPlaceholderFilled(cell, mapx.Node{Kind: mapx.KindSpec}, "", nil, cfg); v != Fail {
		t.Errorf("a declared word in a table cell is a marker too, got %v", v)
	}
	if v, _ := checkPlaceholderFilled("<!-- @anchors\n  layer: TODO\n-->\n", mapx.Node{Kind: mapx.KindSpec}, "", nil, cfg); v != Pass {
		t.Errorf("a declared vocabulary REPLACES the default: TODO is a value here, got %v", v)
	}
}
