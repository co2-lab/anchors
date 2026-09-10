package main

import (
	"bytes"
	"strings"
	"testing"
)

// O guide de projeto é a régua de uma fase que o agente executa SOZINHO, na conversa
// — não há gate que a confronte depois. Se o texto perder uma das peças, ninguém
// descobre por falha: descobre por um PROJECT.md pela metade meses depois. Daí o
// teste cobrar as peças que a fase não sobrevive sem.
func TestProjectGuideCobreAsPecasDaFase(t *testing.T) {
	casos := map[string]string{
		"os dois arquivos":                        "PROJECT.md",
		"a transcrição":                           "INSIGHTS.md",
		"o descarte (o que dá valor ao INSIGHTS)": "Descartado",
		"a revisão de inconsistências":            "inconsistências",
		"a entrevista roda na conversa":           "CONVERSA",
		"uma pergunta por vez":                    "UMA pergunta por vez",
		"ser opinativo":                           "OPINATIVO",
		"a ponte para o init":                     "anchors init",
	}
	for nome, want := range casos {
		if !strings.Contains(projectGuide, want) {
			t.Errorf("guide de projeto deveria cobrir %s (esperava %q)", nome, want)
		}
	}
}

// As cinco etapas precisam estar TODAS lá e o recorte é técnico: o que o usuário
// pediu (tecnologias, paradigma, estrutura, extensões, editores, indentação) é
// exatamente o que o PROJECT.md existe para responder.
func TestProjectGuideTemAsCincoEtapasTecnicas(t *testing.T) {
	for _, etapa := range []string{
		"Etapa 1 — Propósito e forma",
		"Etapa 2 — Linguagem e runtime",
		"Etapa 3 — Arquitetura e paradigma",
		"Etapa 4 — Estrutura macro e convenções de arquivo",
		"Etapa 5 — Ferramental e formatação",
	} {
		if !strings.Contains(projectGuide, etapa) {
			t.Errorf("falta a %q", etapa)
		}
	}
	for _, assunto := range []string{"indentação", "extensões", "editores", "paradigma", "co-location"} {
		if !strings.Contains(projectGuide, assunto) {
			t.Errorf("o recorte técnico deveria cobrir %q", assunto)
		}
	}
}

// A fase só serve se o agente souber que ela existe: quem lê o playbook (`anchors
// guide`, sem argumento) tem de ser mandado para cá ANTES de planejar. Um guide que
// só é achado por quem já sabe que ele existe é o problema que o USING.md registra.
func TestPlaybookApontaParaAFaseDescobrir(t *testing.T) {
	for _, want := range []string{"anchors guide project", "PROJECT.md", "INSIGHTS.md"} {
		if !strings.Contains(agentGuide, want) {
			t.Errorf("o playbook deveria citar %q", want)
		}
	}
	// e antes do PLANEJAR: um projeto descoberto DEPOIS do plano é um plano escrito
	// sem saber a linguagem.
	if strings.Index(agentGuide, "DESCOBRIR") > strings.Index(agentGuide, "### 1. PLANEJAR") {
		t.Error("DESCOBRIR deveria vir antes de PLANEJAR no playbook")
	}
}

// O subcomando tem de existir e imprimir o guide — é o único caminho até o texto.
func TestGuideProjectSubcomandoImprime(t *testing.T) {
	cmd := newGuideCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"project"})

	sub, _, err := cmd.Find([]string{"project"})
	if err != nil || sub.Name() != "project" {
		t.Fatalf("`anchors guide project` deveria existir (err=%v)", err)
	}
	if sub.Short == "" {
		t.Error("o subcomando deveria ter descrição curta (aparece no --help)")
	}
}
