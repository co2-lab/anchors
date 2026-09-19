package health

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// escreveSpec cria uma spec no disco e devolve o nó correspondente.
func escreveSpec(t *testing.T, root, nome, corpo string) mapx.Node {
	t.Helper()
	caminho := filepath.Join(root, nome)
	if err := os.WriteFile(caminho, []byte(corpo), 0o644); err != nil {
		t.Fatalf("escrevendo %s: %v", nome, err)
	}
	// Layer "spec" de propósito: é o que o mapa costuma trazer (a camada do ARQUIVO).
	// Quem diz a camada da UNIDADE é o header — e é isso que o check tem de ler.
	return mapx.Node{ID: nome, Kind: mapx.KindSpec, Layer: "spec"}
}

const specTela = `<!-- @anchors
  code: SGINS
  layer: screen
-->
# SignIn

## Visão Geral
Entra no app.

## Regras

### SGINS-B01 — regra
Comportamento.
`

// A distinção que dá o check: gate DECLARADO e cego é Warn e pede edição das specs;
// gate NÃO declarado é Info e pede adotar a seção E o gate. A ação corretiva é diferente,
// então a severidade e o texto também têm de ser.
func TestCheckSpecSections_GateDeclaradoECegoEhWarn(t *testing.T) {
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{escreveSpec(t, root, "SignIn.spec.md", specTela)}}
	cfg := &config.Config{Gates: []config.Gate{{Name: "dependency-honored"}}}

	achados := checkSpecSections(g, cfg, root)
	var achado *Finding
	for i := range achados {
		if achados[i].Subject == "data-contract" {
			achado = &achados[i]
		}
	}
	if achado == nil {
		t.Fatal("esperava achado para data-contract: a spec de tela não tem a seção")
	}
	if achado.Check != "secao-ausente" {
		t.Errorf("check = %q, queria secao-ausente (o gate está declarado)", achado.Check)
	}
	if achado.Severity != Warn {
		t.Errorf("severidade = %v, queria Warn: há gate declarado e cego", achado.Severity)
	}
}

func TestCheckSpecSections_GateNaoDeclaradoEhInfo(t *testing.T) {
	root := t.TempDir()
	g := &mapx.Graph{Nodes: []mapx.Node{escreveSpec(t, root, "SignIn.spec.md", specTela)}}
	// Nenhum gate declarado: `route-declared` não existe neste projeto.
	cfg := &config.Config{}

	for _, f := range checkSpecSections(g, cfg, root) {
		if f.Subject != "navigation" {
			continue
		}
		if f.Check != "secao-recomendada" {
			t.Errorf("check = %q, queria secao-recomendada (gate não declarado)", f.Check)
		}
		if f.Severity != Info {
			t.Errorf("severidade = %v, queria Info", f.Severity)
		}
		return
	}
	t.Fatal("esperava achado para navigation")
}

// A seção presente não pode ser acusada — e o título vale em QUALQUER idioma suportado,
// porque a spec foi escrita sob o `lang:` que o projeto tinha naquele dia.
func TestCheckSpecSections_TituloEmOutroIdiomaConta(t *testing.T) {
	root := t.TempDir()
	comNavegacaoEmIngles := specTela + "\n## Navigation\n| Origem | Gatilho |\n| --- | --- |\n"
	g := &mapx.Graph{Nodes: []mapx.Node{escreveSpec(t, root, "SignIn.spec.md", comNavegacaoEmIngles)}}

	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		if f.Subject == "navigation" {
			t.Errorf("acusou navigation ausente, mas a seção existe com o título em inglês: %s", f.Detail)
		}
	}
}

// Cobrar rota/testid de uma unidade que não é tela é o falso-positivo do validador
// legado, que tratava todo artefato como tela.
func TestCheckSpecSections_NaoCobraSecaoDeTelaDeOutraCamada(t *testing.T) {
	root := t.TempDir()
	logica := `<!-- @anchors
  code: CALCX
  layer: backend-logic
-->
# Calc

## Visão Geral
Soma.

## Regras

### CALCX-B01 — regra
Comportamento.
`
	g := &mapx.Graph{Nodes: []mapx.Node{escreveSpec(t, root, "Calc.spec.md", logica)}}

	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		switch f.Subject {
		case "navigation", "testids", "data-contract", "data-states", "states":
			t.Errorf("cobrou seção de tela (%s) de uma unidade backend-logic", f.Subject)
		}
	}
}

// O check reporta o PADRÃO, não o caso isolado: uma spec sem a seção entre várias que a
// têm é decisão de quem escreveu, e virar aviso treinaria o time a ignorar a saída.
func TestCheckSpecSections_MinoriaAusenteNaoReporta(t *testing.T) {
	root := t.TempDir()
	comDominio := `<!-- @anchors
  code: AAAAA
  layer: backend-logic
-->
# A

## Visão Geral
x

## Domínio
| Entrada | Aceita |
| --- | --- |
`
	semDominio := `<!-- @anchors
  code: BBBBB
  layer: backend-logic
-->
# B

## Visão Geral
y
`
	g := &mapx.Graph{Nodes: []mapx.Node{
		escreveSpec(t, root, "A.spec.md", comDominio),
		escreveSpec(t, root, "B.spec.md", comDominio),
		escreveSpec(t, root, "C.spec.md", comDominio),
		escreveSpec(t, root, "D.spec.md", semDominio),
	}}

	for _, f := range checkSpecSections(g, &config.Config{}, root) {
		if f.Subject == "domain" {
			t.Errorf("reportou domain com só 1 de 4 ausente — deveria calar: %s", f.Detail)
		}
	}
}
