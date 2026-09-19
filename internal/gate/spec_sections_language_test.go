package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// specEmIngles tem a regra catalogada (o que o gate já cobrava) e os títulos de seção
// em INGLÊS — o caso que o `enforce_section_language` existe para pegar num projeto pt-BR.
const specEmIngles = `<!-- @anchors
  code: LGNOI
-->
# Login

## Overview
Entra no app.

## Rules

### LGNOI-B01 — regra
Comportamento.
`

const specEmPortugues = `<!-- @anchors
  code: LGNOI
-->
# Login

## Visão Geral
Entra no app.

## Regras

### LGNOI-B01 — regra
Comportamento.
`

func gateSecoes(enforce *bool) *config.Config {
	return &config.Config{
		Lang:  "pt-BR",
		Gates: []config.Gate{{Name: "spec-complete", Check: "spec-sections", EnforceSectionLanguage: enforce}},
	}
}

// O PADRÃO é cobrar: omitir a chave não pode significar "não verifica", senão o acervo
// misto nasce em silêncio — que é o defeito que a opção existe para tornar visível.
func TestSpecSections_IdiomaErradoReprovaPorPadrao(t *testing.T) {
	v, msg := checkSpecSections(specEmIngles, mapx.Node{}, "", nil, gateSecoes(nil))
	if v != Fail {
		t.Fatalf("veredito = %v, queria Fail: a spec está em inglês num projeto pt-BR", v)
	}
	if !strings.Contains(msg, "Visão Geral") {
		t.Errorf("a mensagem deveria dizer o título ESPERADO, para o conserto ser óbvio; veio: %s", msg)
	}
}

func TestSpecSections_IdiomaCertoPassa(t *testing.T) {
	if v, msg := checkSpecSections(specEmPortugues, mapx.Node{}, "", nil, gateSecoes(nil)); v != Pass {
		t.Fatalf("veredito = %v (%s), queria Pass", v, msg)
	}
}

// Desligar é decisão legítima e declarada: migração em curso, ou projeto bilíngue.
func TestSpecSections_EnforceFalseNaoCobraIdioma(t *testing.T) {
	nao := false
	if v, msg := checkSpecSections(specEmIngles, mapx.Node{}, "", nil, gateSecoes(&nao)); v != Pass {
		t.Fatalf("veredito = %v (%s), queria Pass com enforce_section_language: false", v, msg)
	}
}

// O léxico PRÓPRIO do projeto não é idioma errado — é seção que o framework não nomeia.
// Acusá-la seria impor o vocabulário do engine ao projeto.
func TestSpecSections_TituloForaDoCatalogoNaoEhIdiomaErrado(t *testing.T) {
	spec := specEmPortugues + "\n## Fora de escopo\nNada.\n\n## Decisões em aberto\nnenhuma\n"
	if v, msg := checkSpecSections(spec, mapx.Node{}, "", nil, gateSecoes(nil)); v != Pass {
		t.Fatalf("veredito = %v (%s), queria Pass: títulos próprios do projeto não são idioma errado", v, msg)
	}
}

// Sem config o gate não pode inventar idioma: cobra só o que sempre cobrou.
func TestSpecSections_SemConfigNaoCobraIdioma(t *testing.T) {
	if v, msg := checkSpecSections(specEmIngles, mapx.Node{}, "", nil, nil); v != Pass {
		t.Fatalf("veredito = %v (%s), queria Pass sem config", v, msg)
	}
}
