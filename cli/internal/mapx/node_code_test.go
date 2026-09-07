package mapx

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// A identidade do nó é a DECLARADA no header. Antes, era inferida do primeiro código de
// cenário que aparecia no texto — e uma spec que CITA outra unidade antes de definir a
// sua entrava no mapa com a identidade errada. Medido: uma spec de modelo que abria
// referenciando `DTAXX-B11` era registrada como dona de `DTAXX`, e todo gate relacional
// passava a confrontar a unidade errada, em silêncio.
func TestNodeCodePreferHeader(t *testing.T) {
	casos := []struct {
		nome     string
		f        scan.File
		esperado string
	}{
		{"header declarado vence a citação que vem antes",
			scan.File{HeaderCode: "MTENX", Codes: []string{"DTAXX-B11", "MTENX-B01"}}, "MTENX"},
		{"sem header, infere do primeiro código (fallback)",
			scan.File{Codes: []string{"ABCDX-B01"}}, "ABCDX"},
		{"sem header e sem código: vazio",
			scan.File{}, ""},
		{"header declarado vale mesmo sem código de cenário no corpo",
			scan.File{HeaderCode: "WXYZX"}, "WXYZX"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if got := nodeCode(c.f, nil); got != c.esperado {
				t.Fatalf("code = %q, queria %q", got, c.esperado)
			}
		})
	}
}

// A ÂNCORA IRMÃ vence a inferência pelo texto, para artefato derivado.
//
// Medido no blue-eyes: `GoLiveChecklist.test.ts` citava `ELKAD-B01` numa string de dado —
// o `estadoAtual` de uma dívida fictícia, "o ELKAD-B01 preparou o caminho" — e o arquivo
// entrou no mapa com `code: ELKAD`.
//
// O efeito: o `scenario-coverage` passou a cobrar 21 cenários de OUTRAS specs, e os três
// invariantes que aquele teste realmente provava apareciam como não provados. Um teste
// verde, com rastreabilidade apontando para a unidade errada.
func TestNodeCode_ancoraIrmaVenceOTextoDoDerivado(t *testing.T) {
	f := scan.File{
		Path:  "packages/infra/GoLiveChecklist.test.ts",
		Kind:  "test",
		Codes: []string{"ELKAD-B01", "GLCGL-I01"},
	}
	ancoras := map[string]string{"packages/infra/GoLiveChecklist.test.ts": "GLCGL"}

	if got := nodeCode(f, ancoras); got != "GLCGL" {
		t.Errorf("code = %q, queria GLCGL — o dado de teste foi lido como declaração", got)
	}
}

// O HEADER continua vencendo tudo: quem declara não é sobrescrito pela âncora.
func TestNodeCode_headerVenceAAncora(t *testing.T) {
	f := scan.File{Path: "x/Y.test.ts", HeaderCode: "DECLR", Codes: []string{"OUTRO-B01"}}
	ancoras := map[string]string{"x/Y.test.ts": "ANCOR"}

	if got := nodeCode(f, ancoras); got != "DECLR" {
		t.Errorf("code = %q, queria DECLR — o header é onde o autor diz de quem é o arquivo", got)
	}
}

// Derivado SEM âncora resolvida cai no fallback: é o caso de um arquivo cuja spec irmã
// não declara identidade, e barrar ali deixaria o nó sem código nenhum.
func TestNodeCode_semAncoraCaiNoFallback(t *testing.T) {
	f := scan.File{Path: "x/Y.test.ts", Kind: "test", Codes: []string{"ABCDX-B01"}}
	if got := nodeCode(f, map[string]string{}); got != "ABCDX" {
		t.Errorf("code = %q, queria ABCDX (fallback)", got)
	}
}

func configComAncoraSpec() *config.Config {
	return &config.Config{Derived: &config.Derived{Anchor: "spec"}}
}

// anchorCodeByDerived liga cada derivado ao código da âncora pelo STEM.
func TestAnchorCodeByDerived_ligaPeloStem(t *testing.T) {
	cfg := configComAncoraSpec()
	files := []scan.File{
		{Path: "packages/infra/GoLive.spec.md", Kind: "spec", HeaderCode: "GLCGL"},
		{Path: "packages/infra/GoLive.ts", Kind: "code"},
		{Path: "packages/infra/GoLive.test.ts", Kind: "test"},
		{Path: "packages/infra/GoLive.feature", Kind: "feature"},
		// Outra unidade no MESMO diretório: não pode receber o código da vizinha.
		{Path: "packages/infra/Outra.test.ts", Kind: "test"},
	}

	got := anchorCodeByDerived(files, cfg)

	for _, p := range []string{"packages/infra/GoLive.ts", "packages/infra/GoLive.test.ts",
		"packages/infra/GoLive.feature"} {
		if got[p] != "GLCGL" {
			t.Errorf("%s → %q, queria GLCGL", p, got[p])
		}
	}
	if c, ok := got["packages/infra/Outra.test.ts"]; ok {
		t.Errorf("Outra.test.ts recebeu %q — o stem não casa, e o código da vizinha não é dele", c)
	}
}
