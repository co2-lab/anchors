package ops

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func renderScreen(t *testing.T, cfg *config.Config) string {
	t.Helper()
	sections, ordem, err := resolveSectionsWithPreset(specTemplate, "screen", nil, nil)
	if err != nil {
		t.Fatalf("resolvendo preset screen: %v", err)
	}
	return renderArtifact(specTemplate, "Login", "LGNOI", "Login.spec.md", t.TempDir(), sections, ordem, cfg)
}

// O TÍTULO SEGUE O `lang:`. Antes disto o catálogo emitia português para todo projeto —
// inclusive `lang: en` e `lang: es` —, contradizendo o default do próprio framework.
func TestTituloDeSecaoSegueOIdiomaDoProjeto(t *testing.T) {
	casos := []struct {
		lang   string
		espera string
		proibe string
	}{
		{"en", "## Overview", "## Visão Geral"},
		{"es", "## Visión General", "## Visão Geral"},
		{"pt-BR", "## Visão Geral", "## Overview"},
	}
	for _, c := range casos {
		out := renderScreen(t, &config.Config{Lang: c.lang})
		if !strings.Contains(out, c.espera) {
			t.Errorf("lang=%s: esperava %q no artefato", c.lang, c.espera)
		}
		if strings.Contains(out, c.proibe) {
			t.Errorf("lang=%s: o artefato trouxe %q, de outro idioma", c.lang, c.proibe)
		}
	}
}

// Sem `lang:` declarado vale o DEFAULT do framework, que é inglês — não a língua de quem
// escreveu o framework.
func TestSemLangDeclaradoUsaODefaultIngles(t *testing.T) {
	out := renderScreen(t, &config.Config{})
	if !strings.Contains(out, "## Overview") {
		t.Error("sem `lang:` o default é inglês (i18n.Default), e o artefato deveria nascer em inglês")
	}
}

// O LÉXICO DO PROJETO vence o idioma: quem declarou "Modelo de Dado" quer isso, mesmo num
// projeto `lang: en`. A precedência é section_titles > rule_types > idioma.
func TestLexicoDoProjetoVenceOIdioma(t *testing.T) {
	cfg := &config.Config{
		Lang:          "en",
		SectionTitles: config.SectionTitles{"overview": "Panorama"},
	}
	out := renderScreen(t, cfg)
	if !strings.Contains(out, "## Panorama") {
		t.Error("`section_titles` deveria vencer a tradução do framework")
	}
	if strings.Contains(out, "## Overview") {
		t.Error("com léxico próprio declarado, o título do framework não deveria aparecer")
	}
}

// Toda seção do catálogo de spec tem de ter tradução: uma chave faltando faz a seção sair
// em português dentro de um artefato inglês, e o defeito passa despercebido porque o resto
// do arquivo está certo. Foi o que aconteceu com `a11y`.
func TestTodaSecaoDeSpecTemTraducao(t *testing.T) {
	out := renderScreen(t, &config.Config{Lang: "en"})
	// O CORPO também é traduzido, não só o título: cabeçalho de tabela e dica de TODO
	// saíam em português dentro de um artefato inglês, e o defeito passava despercebido
	// porque o título ao lado estava certo.
	for _, pt := range []string{"TODO propósito", "**Código**", "Obrigatório", "Quem garante", "nenhuma"} {
		if strings.Contains(out, pt) {
			t.Errorf("o artefato `lang: en` trouxe %q no CORPO — falta a chave section.body.*", pt)
		}
	}
	// Nenhum título do artefato inglês pode estar em português.
	for _, pt := range []string{"Visão Geral", "Regras", "Dependências", "Acessibilidade", "Decisões em aberto"} {
		if strings.Contains(out, "## "+pt) {
			t.Errorf("o artefato `lang: en` trouxe a seção %q em português — falta a chave section.title.*", pt)
		}
	}
}
