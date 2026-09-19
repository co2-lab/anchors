package gate

import (
	"os"
	"path/filepath"
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

func TestCodeReferenceSoSpec(t *testing.T) {
	t.Run("CRVCD-B01: Non-specification artifacts skip confrontation", func(t *testing.T) {})
	t.Run("CRVCD-X02: Non-specification artifacts are not inspected by this gate", func(t *testing.T) {})
	codeD := "DT" + "AXX"
	codeZ := "ZZ" + "ZZX-B01"
	spec := "code: " + codeD + "\nAplica `" + codeZ + "`.\n"
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindTest, mapx.KindFeature} {
		if v, _ := checkCodeReferenceValid(spec, mapx.Node{Kind: k}, "", grafoComCodigos(codeD), nil); v != Skip {
			t.Errorf("kind %s deveria ser Skip, foi %s", k, v)
		}
	}
}

func TestCodeReferenceSemMapaEhPendente(t *testing.T) {
	t.Run("CRVCD-B02: Confronting without a map graph returns pending", func(t *testing.T) {})
	t.Run("CRVCD-B03: Confronting with an empty identity universe returns pending", func(t *testing.T) {})
	t.Run("CRVCD-I02: Missing map or empty identity universe returns pending rather than pass", func(t *testing.T) {})
	codeD := "DT" + "AXX"
	codeZ := "ZZ" + "ZZX-B01"
	spec := "<!-- @anchors\n  code: " + codeD + "\n-->\nAplica `" + codeZ + "`.\n"
	if v, _ := rodaRef(t, spec, nil); v != Pending {
		t.Fatalf("sem mapa deveria ser Pending, foi %s", v)
	}
	if v, _ := rodaRef(t, spec, &mapx.Graph{}); v != Pending {
		t.Fatalf("mapa vazio deveria ser Pending (não ✓ sobre um universo que não conhece), foi %s", v)
	}
}

func TestCodeReferenceResolvida(t *testing.T) {
	t.Run("CRVCD-B04: Citations resolving to existing units in the map pass", func(t *testing.T) {})
	t.Run("CRVCD-X01: Implementation correctness of referenced requirements is not evaluated", func(t *testing.T) {})
	codeD := "DT" + "AXX"
	codeM := "MT" + "VRX"
	codeR := "RD" + "MDX"
	refM := codeM + "-B01"
	refR := codeR + "-B03"
	spec := "<!-- @anchors\n  code: " + codeD + "\n-->\nAplica `" + refM + "` e `" + refR + "`.\n"
	if v, d := rodaRef(t, spec, grafoComCodigos(codeD, codeM, codeR)); v != Pass {
		t.Fatalf("citação de códigos existentes deveria passar, foi %s (%s)", v, d)
	}
}

// O caso real: uma spec de schema afirmava "Índices que o schema criou (2026-08-11)" e
// referenciava 4 códigos de specs que não existiam. Passou por TODOS os gates — tinha
// código, header, seções, e o gate de dependências não olha prosa. Um leitor futuro a
// toma como registro do que foi feito.
func TestCodeReferenceOrfa(t *testing.T) {
	t.Run("CRVCD-B05: Citations pointing to non-existent units fail", func(t *testing.T) {})
	t.Run("CRVCD-B07: Multiple orphaned requirement citations are reported sorted", func(t *testing.T) {})
	t.Run("CRVCD-I03: Unresolvable external citations always produce a blocking fail verdict", func(t *testing.T) {})
	codeD := "DT" + "AXX"
	codeT := "TX" + "RC"
	codeR := "RD" + "MDX"
	codeM := "MT" + "VRX"
	refR := codeR + "-B03"
	refM := codeM + "-B01"
	spec := "<!-- @anchors\n  code: " + codeD + "\n-->\n# Schema\n\n### " + codeD + "-B11 — MetadataEntry\n" +
		"O consumidor resolve para o mês (`" + refR + "`), aplicando `" + refM + "`.\n"
	v, d := rodaRef(t, spec, grafoComCodigos(codeD, codeT))
	if v != Fail {
		t.Fatalf("citação de código inexistente deveria reprovar, foi %s (%s)", v, d)
	}
	for _, c := range []string{codeR, codeM} {
		if !strings.Contains(d, c) {
			t.Errorf("não nomeou o código órfão %s: %s", c, d)
		}
	}
	// Verifica ordenação alfabética dos órfãos reportados
	idxM := strings.Index(d, codeM)
	idxR := strings.Index(d, codeR)
	if idxM == -1 || idxR == -1 || idxM > idxR {
		t.Errorf("órfãos deveriam vir ordenados alfabeticamente: %s", d)
	}
}

// A spec é dona do próprio código: os requisitos que ela define não são "referência a
// outra unidade" e não podem ser cobrados como órfãos.
func TestCodeReferenceProprioCodigoNaoEhOrfao(t *testing.T) {
	t.Run("CRVCD-B06: Citations matching the specification's own identity code pass", func(t *testing.T) {})
	t.Run("CRVCD-I01: Self-references to a specification's own requirements never fail", func(t *testing.T) {})
	codeM := "MT" + "VRX"
	ref1 := codeM + "-B01"
	ref2 := codeM + "-I02"
	spec := "<!-- @anchors\n  code: " + codeM + "\n-->\n" +
		"### " + ref1 + " — carry-forward\n" +
		"### " + ref2 + " — editar não muda o passado\n" +
		"Ver `" + ref1 + "` acima.\n"
	if v, d := rodaRef(t, spec, grafoComCodigos(codeM)); v != Pass {
		t.Fatalf("os códigos que a própria spec define não são órfãos, foi %s (%s)", v, d)
	}
}

func TestCodeReferenceHeaderOnDisk(t *testing.T) {
	t.Run("CRVCD-B08: Identity ownership is resolved from graph nodes and header metadata", func(t *testing.T) {})
	root := t.TempDir()
	codeA := "EX" + "TNL"
	codeB := "MY" + "SPC"
	refA := codeA + "-B01"
	// Cria arquivo da spec EXTNL no disco sem código explícito no Node
	specAContent := "<!-- @anchors\n  code: " + codeA + "\n-->\n# External Unit\n"
	if err := os.WriteFile(filepath.Join(root, "extnl.spec.md"), []byte(specAContent), 0o644); err != nil {
		t.Fatal(err)
	}
	g := &mapx.Graph{
		Nodes: []mapx.Node{
			{ID: "extnl.spec.md", Kind: mapx.KindSpec}, // Code vazio, deve ler do disco
			{ID: "myspc.spec.md", Kind: mapx.KindSpec, Code: codeB},
		},
	}
	mySpec := "<!-- @anchors\n  code: " + codeB + "\n-->\nCita `" + refA + "` que existe no disco.\n"
	v, d := checkCodeReferenceValid(mySpec, mapx.Node{ID: "myspc.spec.md", Kind: mapx.KindSpec}, root, g, nil)
	if v != Pass {
		t.Fatalf("deveria ler código do disco e passar, foi %s (%s)", v, d)
	}
}

func TestCodeReferenceNonRequirementTokensIgnored(t *testing.T) {
	t.Run("CRVCD-B09: Non-requirement tokens are ignored", func(t *testing.T) {})
	codeM := "MY" + "SPC"
	spec := "<!-- @anchors\n  code: " + codeM + "\n-->\n" +
		"# Title\n\nWords like HTTP, JSON, UUID, UTF8, internal/gate/code.go, and random-dashed-words.\n"
	if v, d := rodaRef(t, spec, grafoComCodigos(codeM)); v != Pass {
		t.Fatalf("tokens arbitrários não devem ser cobrados como referências: %s (%s)", v, d)
	}
}

func TestCodeReferenceNoExternalCitationsPasses(t *testing.T) {
	t.Run("CRVCD-X03: External requirement citations are not mandatory", func(t *testing.T) {})
	codeM := "MY" + "SPC"
	spec := "<!-- @anchors\n  code: " + codeM + "\n-->\n" +
		"# Title\n\nPure self-contained prose without any external requirement codes.\n"
	if v, d := rodaRef(t, spec, grafoComCodigos(codeM)); v != Pass {
		t.Fatalf("spec sem citações externas deve passar: %s (%s)", v, d)
	}
}
