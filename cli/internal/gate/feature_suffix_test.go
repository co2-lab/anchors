package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// Uma regra tem legitimamente mais de um cenário — caminho feliz e alternativos.
// Sem sufixo, os N cenários carregam o mesmo código e nada os distingue: o gate
// compara os N títulos com o mesmo teste e no máximo um casa.
func TestSufixoDaIdentidadeACadaCenario(t *testing.T) {
	feature := `
  @USBPX-B01#01 @USBPX-B02#01 @nivel-unit
  Cenário: busca pontos por userId+month
    Então o repository é consultado

  @USBPX-B01#02 @nivel-unit
  Cenário: sem usuário logado nada é buscado
    Então o repository não é consultado
`
	cenarios := parseFeatureScenarios(feature)
	if len(cenarios) != 2 {
		t.Fatalf("esperava 2 cenários, veio %d", len(cenarios))
	}
	if cenarios[0].Code != "USBPX-B01#01" {
		t.Errorf("código do 1º: %q, queria USBPX-B01#01", cenarios[0].Code)
	}
	if cenarios[1].Code != "USBPX-B01#02" {
		t.Errorf("código do 2º: %q, queria USBPX-B01#02", cenarios[1].Code)
	}
	// O segundo código da tag-line também carrega o sufixo.
	if len(cenarios[0].Codes) != 2 || cenarios[0].Codes[1] != "USBPX-B02#01" {
		t.Errorf("códigos do 1º: %v", cenarios[0].Codes)
	}
}

// Retrocompatível: o sufixo é opcional. Um projeto que nunca o adote continua
// funcionando — e essa é a condição para o gate poder mudar sem quebrar ninguém.
func TestCodigoSemSufixoContinuaValendo(t *testing.T) {
	cenarios := parseFeatureScenarios(`
  @SAUTX-B01 @nivel-unit
  Cenário: Hidratar carrega a sessão
    Então o usuário fica disponível
`)
	if len(cenarios) != 1 || cenarios[0].Code != "SAUTX-B01" {
		t.Fatalf("código sem sufixo quebrou: %+v", cenarios)
	}
}

// Os gates que falam de REGRA precisam da raiz; os que falam de CENÁRIO, do código
// inteiro. Confundir os dois faria a regra `USBPX-B01` parecer três regras.
func TestCodeRaizSeparaRegraDeCenario(t *testing.T) {
	casos := map[string]string{
		"USBPX-B01#02": "USBPX-B01",
		"USBPX-B01":    "USBPX-B01",
		"ATLNX-VR":     "ATLNX-VR",
		"MNMTX-DS-kv":  "MNMTX-DS-kv",
	}
	for entrada, quer := range casos {
		if got := CodeRaiz(entrada); got != quer {
			t.Errorf("CodeRaiz(%q) = %q, queria %q", entrada, got, quer)
		}
	}
}

// Um requisito catalogado sem sufixo é coberto por cenários COM sufixo: `#01`/`#02`
// dão identidade a cada caso, não criam requisito novo. Sem isto, numerar cenários
// faria o spec-feature-match acusar de repente todo requisito com mais de um caso —
// medido no app de referência: 78 specs reprovaram assim que 269 cenários ganharam sufixo.
func TestSpecFeatureMatchAceitaSufixoNoCenario(t *testing.T) {
	root := t.TempDir()
	spec := "business-logic/dedup.spec.md"
	feat := "business-logic/dedup.feature"
	writeFile(t, root, spec, "# Dedup\n\n## Rules\n\n| Regra | Descrição |\n| --- | --- |\n| `DDTDX-B01` | duplicata automática |\n")
	writeFile(t, root, feat, `# language: pt
Funcionalidade: Dedup

  @comportamento @DDTDX-B01#01 @nivel-unit @P1
  Cenário: Duplicata automática quando descrição e valor idênticos
    Dado x
    Então y

  @comportamento @DDTDX-B01#02 @nivel-unit @P1
  Cenário: Duplicata automática quando só a data difere
    Dado x
    Então y
`)
	g := &mapx.Graph{
		Nodes: []mapx.Node{{ID: spec, Kind: mapx.KindSpec}, {ID: feat, Kind: mapx.KindFeature}},
		Edges: []mapx.Edge{{From: spec, To: feat, Type: "covered-by"}},
	}
	n := mapx.Node{ID: spec, Kind: mapx.KindSpec}
	content, err := os.ReadFile(filepath.Join(root, spec))
	if err != nil {
		t.Fatal(err)
	}
	v, detail := checkSpecFeatureMatch(string(content), n, root, g, regimeCfg())
	if v == Fail {
		t.Errorf("sufixo de cenário não deveria reprovar o requisito raiz: %s", detail)
	}
}
