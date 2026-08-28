package mapx

import (
	"testing"

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
			if got := nodeCode(c.f); got != c.esperado {
				t.Fatalf("code = %q, queria %q", got, c.esperado)
			}
		})
	}
}
