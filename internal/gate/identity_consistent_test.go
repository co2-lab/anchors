package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cenário mínimo: uma spec com `code`, apontando (specifies) para a unidade cujo
// conteúdo o teste controla. `outros` são códigos de OUTRAS unidades no mapa — é o
// conjunto que distingue reuso deliberado de identidade órfã.
func identidadeFixture(t *testing.T, code, unitSrc string, outros ...string) (mapx.Node, *mapx.Graph, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "x.tsx"), []byte(unitSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	spec := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: code}
	g := &mapx.Graph{
		Nodes: []mapx.Node{spec, {ID: "x.tsx", Kind: mapx.KindCode}},
		Edges: []mapx.Edge{{From: "x.spec.md", To: "x.tsx", Type: mapx.EdgeSpecifies}},
	}
	for i, c := range outros {
		g.Nodes = append(g.Nodes, mapx.Node{
			ID: filepath.Join("outro", string(rune('a'+i))+".spec.md"), Kind: mapx.KindSpec, Code: c,
		})
	}
	return spec, g, root
}

func TestIdentityConsistent_naoSpecPula(t *testing.T) {
	t.Run("IDCND-B01: Confronting an artifact that is not a spec skips", func(t *testing.T) {})
	n := mapx.Node{ID: "x.tsx", Kind: mapx.KindCode, Code: "BGET"}
	if v, _ := checkIdentityConsistent("", n, "", &mapx.Graph{}, &config.Config{}); v != Skip {
		t.Errorf("não-spec deve pular: %v", v)
	}
}

func TestIdentityConsistent_semGrafoNaoAfirma(t *testing.T) {
	t.Run("IDCND-B02: Confronting without a graph returns Pending", func(t *testing.T) {})
	t.Run("IDCND-I01: Without a map graph the relational gate never approves", func(t *testing.T) {})
	n := mapx.Node{ID: "x.spec.md", Kind: mapx.KindSpec, Code: "BGET"}
	if v, _ := checkIdentityConsistent("", n, "", nil, &config.Config{}); v != Pending {
		t.Errorf("sem mapa o gate relacional não pode afirmar nada: %v", v)
	}
}

func TestIdentityConsistent_semCodigoNaoDuplicaAchado(t *testing.T) {
	t.Run("IDCND-B03: Confronting a spec without a declared code skips", func(t *testing.T) {})
	t.Run("IDCND-I03: Absence of a spec code is never double-charged", func(t *testing.T) {})
	t.Run("IDCND-X02: The gate does not enforce code presence on specs", func(t *testing.T) {})
	// A ausência de código é ofensa de `spec-tem-codigo`. Cobrá-la aqui de novo
	// transformaria um achado em dois, o segundo sem nada a acrescentar.
	n, g, root := identidadeFixture(t, "", `<View testID=":qualquer-coisa" />`)
	if v, _ := checkIdentityConsistent("", n, root, g, &config.Config{}); v != Skip {
		t.Errorf("spec sem código deveria pular (é ofensa de outro gate): %v", v)
	}
}

func TestIdentityConsistent_testIDBateComOCodigo(t *testing.T) {
	t.Run("IDCND-B04: A testID prefix matching the spec identity code passes", func(t *testing.T) {})
	t.Run("IDCND-B10: When no orphan testID prefixes or baseline discrepancies exist the gate passes", func(t *testing.T) {})
	n, g, root := identidadeFixture(t, "BGET", `<View testID=":bget-screen" />`)
	// Adiciona aresta specifies para arquivo inexistente
	g.Edges = append(g.Edges, mapx.Edge{From: "x.spec.md", To: "nonexistent.tsx", Type: mapx.EdgeSpecifies})
	if v, msg := checkIdentityConsistent("", n, root, g, &config.Config{}); v != Pass {
		t.Errorf("testID igual ao código deveria passar ignorando arquivo inexistente: %v (%s)", v, msg)
	}
}

func TestIdentityConsistent_siglaOrfaReprova(t *testing.T) {
	t.Run("IDCND-B05: An orphan testID prefix with code shape fails", func(t *testing.T) {})
	t.Run("IDCND-B09: Inconsistency failures cite conflicting acronyms and origin files", func(t *testing.T) {})
	// O caso que motivou o gate: spec diz BGET, testID diz bdge, e `bdge` não é
	// código de unidade alguma — logo não está no dicionário gerado do mapa.
	n, g, root := identidadeFixture(t, "BGET", `<View testID=":bdge-screen" />`)
	v, msg := checkIdentityConsistent("", n, root, g, &config.Config{})
	if v != Fail {
		t.Fatalf("sigla órfã deveria reprovar: %v", v)
	}
	if !strings.Contains(msg, "BDGE") || !strings.Contains(msg, "BGET") {
		t.Errorf("a mensagem deve nomear as duas identidades em conflito: %s", msg)
	}
	if !strings.Contains(msg, "x.tsx") {
		t.Errorf("a mensagem deve nomear o arquivo de origem; msg: %s", msg)
	}
}

func TestIdentityConsistent_prefixoDeOutraUnidadeEhReusoLegitimo(t *testing.T) {
	t.Run("IDCND-B06: A testID prefix matching another declared unit is accepted as legitimate reuse", func(t *testing.T) {})
	t.Run("IDCND-X03: Components referencing parent screen codes in testIDs are not forbidden", func(t *testing.T) {})
	// SpendingMonthCard (SMCDX) expõe `:homex-spending-card` porque vive dentro da
	// HomeScreen (HOMEX). O testID nomeia ONDE o elemento aparece — é assim que o
	// flow o alcança, partindo da tela. Cobrar unificação aqui quebraria a navegação.
	//
	// O prefixo do testID tem de DERIVAR do código da outra unidade: é isso que o torna
	// reconhecível como reúso em vez de divergência. A conversão dos exemplos de 4→5 chars
	// mudou o código para `HOMEX` e deixou o testID em `:home-`, quebrando o par — e aí o
	// gate reprovava com razão, porque `home` deixou de ser código de unidade alguma.
	n, g, root := identidadeFixture(t, "SMCDX", `<View testID=":homex-spending-card" />`, "HOMEX")
	if v, msg := checkIdentityConsistent("", n, root, g, &config.Config{}); v != Pass {
		t.Errorf("prefixo que é código de outra unidade é reuso, não divergência: %v (%s)", v, msg)
	}
}

func TestIdentityConsistent_baselineDivergenteReprova(t *testing.T) {
	t.Run("IDCND-B07: A visual regression baseline with a divergent code fails", func(t *testing.T) {})
	t.Run("IDCND-I02: Cross-unit reuse is never permitted for visual regression baselines", func(t *testing.T) {})
	n, g, root := identidadeFixture(t, "BGET", `<View testID=":bget-screen" />`, "BDEDX")
	// O baseline é a prova DESTA unidade — aqui não vale a dispensa do
	// código-de-outra-unidade, porque ele não aponta para onde ela aparece.
	if err := os.WriteFile(filepath.Join(root, "x.BDEDX-VR-loaded.png"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	v, msg := checkIdentityConsistent("", n, root, g, &config.Config{})
	if v != Fail {
		t.Fatalf("baseline com sigla divergente deveria reprovar: %v", v)
	}
	if !strings.Contains(msg, "BDED") {
		t.Errorf("a mensagem deve nomear a sigla do baseline: %s", msg)
	}
	if !strings.Contains(msg, "x.BDEDX-VR-loaded.png") {
		t.Errorf("a mensagem deve nomear o arquivo de baseline; msg: %s", msg)
	}
}

func TestIdentityConsistent_palavraCurtaNaoEhSigla(t *testing.T) {
	t.Run("IDCND-B08: Short testID prefixes of three letters or fewer pass", func(t *testing.T) {})
	t.Run("IDCND-X01: Common shorthand prefixes of three letters or fewer are not scrutinized", func(t *testing.T) {})
	// O regex exige 4-5 letras: `tab-`, `btn-` (3) não têm forma de código e não
	// devem entrar no escrutínio — senão o gate vira ruído em todo componente.
	n, g, root := identidadeFixture(t, "BGET", `<View testID=":btn-save" /><View testID=":tab-x" />`)
	if v, msg := checkIdentityConsistent("", n, root, g, &config.Config{}); v != Pass {
		t.Errorf("prefixo curto não tem forma de código: %v (%s)", v, msg)
	}
}

func TestIdentityConsistent_Errors(t *testing.T) {
	t.Run("IDCND-E01: A governed file gone from disk is left out and the others are still confronted", func(t *testing.T) {
		n, g, root := identidadeFixture(t, "BGET", `<View testID=":bdge-screen" />`)
		g.Edges = append([]mapx.Edge{{From: "x.spec.md", To: "gone.tsx", Type: mapx.EdgeSpecifies}}, g.Edges...)
		v, msg := checkIdentityConsistent("", n, root, g, &config.Config{})
		if v != Fail || !strings.Contains(msg, "BDGE") || strings.Contains(msg, "gone.tsx") {
			t.Fatalf("the missing file must be left out and the orphan in the present one still charged; got %v: %s", v, msg)
		}
	})
}
