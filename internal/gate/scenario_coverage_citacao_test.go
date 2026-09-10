package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O gate cobra os requisitos que a spec DEFINE, não os que ela CITA.
//
// Medido no blue-eyes: a `GoLiveChecklist` define 6 requisitos e o gate cobrava 18
// cenários — 15 deles de outras unidades (`CRPNC-B03`, `MTTLM-B02`, `DTSTD-B06`…), citadas
// na prosa ao justificar as regras dela.
//
// Nenhum daqueles cenários poderia ser provado por um teste desta unidade: o gate pedia o
// impossível, e a mensagem dizia que a spec estava mal coberta. E o efeito colateral é
// pior que o ruído — um gate que sempre reprova é um gate que se aprende a ignorar.
func TestScenarioCoverage_naoCobraOQueASpecApenasCita(t *testing.T) {
	content := `# GoLiveChecklist

## Regras

### GLCGL-B01 — cada item tem um artefato

É o que o ` + "`PLTFR`" + ` estabelece, e a ` + "`DTSTD-B06`" + ` confirma para
infraestrutura. A ` + "`CRPNC-B03`" + ` usa o mesmo raciocínio na rotação.

### GLCGL-B02 — as dívidas têm estado atual

O ` + "`MTTLM-B02`" + ` nasce desligado, e a ` + "`CRPNC-B06`" + ` o liga.
`
	n := mapx.Node{
		ID:   "packages/infra/GoLiveChecklist.spec.md",
		Code: "GLCGL",
		Signal: &mapx.TestSignal{
			ProvenCodes: []string{"GLCGL-B01", "GLCGL-B02"},
			AtRev:       "abc",
		},
		Rev: "abc",
	}

	v, msg := checkScenarioCoverage(content, n)

	if v != Pass {
		t.Errorf("veredito = %v — %s\n  os dois requisitos DEFINIDOS estão provados; o "+
			"resto é citação", v, msg)
	}
	for _, citado := range []string{"DTSTD-B06", "CRPNC-B03", "MTTLM-B02", "CRPNC-B06"} {
		if strings.Contains(msg, citado) {
			t.Errorf("o gate cobrou %q, que esta spec apenas CITA", citado)
		}
	}
}

// E o requisito DEFINIDO sem cenário provado continua sendo cobrado — a correção não pode
// ter desligado o gate.
func TestScenarioCoverage_aindaCobraORequisitoDefinido(t *testing.T) {
	content := "### ABCDX-B01 — a regra\n\nCorpo.\n\n### ABCDX-B02 — outra\n\nCorpo.\n"
	n := mapx.Node{
		Code:   "ABCDX",
		Rev:    "r1",
		Signal: &mapx.TestSignal{ProvenCodes: []string{"ABCDX-B01"}, AtRev: "r1"},
	}

	v, msg := checkScenarioCoverage(content, n)

	if v != Fail {
		t.Fatalf("veredito = %v — o `ABCDX-B02` não tem cenário provado", v)
	}
	if !strings.Contains(msg, "ABCDX-B02") {
		t.Errorf("a mensagem não nomeia o requisito descoberto: %s", msg)
	}
}
