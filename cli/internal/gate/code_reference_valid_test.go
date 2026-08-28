package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func grafoComCodigos(codes ...string) *mapx.Graph {
	g := &mapx.Graph{}
	for _, c := range codes {
		g.Nodes = append(g.Nodes, mapx.Node{ID: c + ".spec.md", Kind: mapx.KindSpec, Code: c})
	}
	return g
}

func rodaRef(t *testing.T, spec string, g *mapx.Graph) (Verdict, string) {
	t.Helper()
	return checkCodeReferenceValid(spec, mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec}, "", g, nil)
}

// O caso real: uma spec de schema afirmava "Índices que o schema criou (2026-08-11)" e
// referenciava 4 códigos de specs que não existiam. Passou por TODOS os gates — tinha
// código, header, seções, e o gate de dependências não olha prosa. Um leitor futuro a
// toma como registro do que foi feito.
func TestCodeReferenceOrfa(t *testing.T) {
	spec := `<!-- @anchors
  code: DTAXX
-->
# Schema

### DTAXX-B11 — MetadataEntry
O consumidor resolve para o mês (` + "`RDMDX-B03`" + `), aplicando ` + "`MTVRX-B01`" + `.
`
	v, d := rodaRef(t, spec, grafoComCodigos("DTAXX", "TXRC"))
	if v != Fail {
		t.Fatalf("citação de código inexistente deveria reprovar, foi %s (%s)", v, d)
	}
	for _, c := range []string{"RDMDX", "MTVRX"} {
		if !strings.Contains(d, c) {
			t.Errorf("não nomeou o código órfão %s: %s", c, d)
		}
	}
}

// Quando as unidades citadas passam a existir, o gate libera — é o ciclo pretendido:
// a citação vira rastreabilidade de verdade.
func TestCodeReferenceResolvida(t *testing.T) {
	spec := "<!-- @anchors\n  code: DTAXX\n-->\nAplica `MTVRX-B01` e `RDMDX-B03`.\n"
	if v, d := rodaRef(t, spec, grafoComCodigos("DTAXX", "MTVRX", "RDMDX")); v != Pass {
		t.Fatalf("citação de códigos existentes deveria passar, foi %s (%s)", v, d)
	}
}

// A spec é dona do próprio código: os requisitos que ela define não são "referência a
// outra unidade" e não podem ser cobrados como órfãos.
func TestCodeReferenceProprioCodigoNaoEhOrfao(t *testing.T) {
	spec := `<!-- @anchors
  code: MTVRX
-->
### MTVRX-B01 — carry-forward
### MTVRX-I02 — editar não muda o passado
Ver ` + "`MTVRX-B01`" + ` acima.
`
	if v, d := rodaRef(t, spec, grafoComCodigos("MTVRX")); v != Pass {
		t.Fatalf("os códigos que a própria spec define não são órfãos, foi %s (%s)", v, d)
	}
}

// O gate é da spec e precisa do universo de identidades: sem mapa, Pendente — nunca ✓.
func TestCodeReferenceSemMapaEhPendente(t *testing.T) {
	spec := "<!-- @anchors\n  code: DTAXX\n-->\nAplica `ZZZZX-B01`.\n"
	if v, _ := rodaRef(t, spec, nil); v != Pending {
		t.Fatalf("sem mapa deveria ser Pending, foi %s", v)
	}
	if v, _ := rodaRef(t, spec, &mapx.Graph{}); v != Pending {
		t.Fatalf("mapa vazio deveria ser Pending (não ✓ sobre um universo que não conhece), foi %s", v)
	}
}

func TestCodeReferenceSoSpec(t *testing.T) {
	spec := "code: DTAXX\nAplica `ZZZZX-B01`.\n"
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindTest, mapx.KindFeature} {
		if v, _ := checkCodeReferenceValid(spec, mapx.Node{Kind: k}, "", grafoComCodigos("DTAXX"), nil); v != Skip {
			t.Errorf("kind %s deveria ser Skip, foi %s", k, v)
		}
	}
}
