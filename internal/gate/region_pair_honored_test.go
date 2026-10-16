// @anchors
//   code: RPHTR
//   ref: RPHRG

package gate

import (
	"github.com/co2-lab/anchors/internal/testkit"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func regionNode() mapx.Node { return mapx.Node{ID: "src/Tela.tsx", Kind: mapx.KindCode} }

func TestRegionPairParBemFormadoPassa(t *testing.T) {
	t.Run("RPHRG-B03: Balanced regions closing with matching identity codes pass", func(t *testing.T) {})
	t.Run("RPHRG-X03: Internal code semantics inside regions are not evaluated", func(t *testing.T) {})
	tagA := "MLETX" + "-A03"
	src := "// #region [" + tagA + "]: persiste.\nput()\n// #endregion [" + tagA + "]"
	v, msg := checkRegionPairHonored(src, regionNode(), "", nil, testkit.WithComments(nil))
	if v != Pass {
		t.Fatalf("esperava Pass, veio %v — %s", v, msg)
	}
}

func TestRegionPairSemRegiaoEhSkipNaoFalha(t *testing.T) {
	t.Run("RPHRG-B02: Code or test artifacts without region markers skip", func(t *testing.T) {})
	t.Run("RPHRG-I01: Region absence is never charged as a failure", func(t *testing.T) {})
	t.Run("RPHRG-X01: The gate does not mandate region markers across all source files", func(t *testing.T) {})
	// A delimitação é OPCIONAL: um projeto que não migrou nada continua válido. Se este
	// gate falhasse por ausência, promovê-lo a bloqueante reprovaria o repositório inteiro
	// no dia em que entrasse — e a região deixaria de ser adotável de forma incremental.
	tagA := "MLETX" + "-A03"
	v, _ := checkRegionPairHonored("// "+tagA+": à moda antiga\nput()", regionNode(), "", nil, testkit.WithComments(nil))
	if v != Skip {
		t.Fatalf("sem região tem de ser Skip, veio %v", v)
	}
}

func TestRegionPairFechoTrocadoFalhaEDizOsDois(t *testing.T) {
	t.Run("RPHRG-B06: An end region marker closing with a different code fails", func(t *testing.T) {})
	t.Run("RPHRG-I02: Pairing defects produce a blocking Fail verdict", func(t *testing.T) {})
	t.Run("RPHRG-I03: End markers must explicitly match opening codes", func(t *testing.T) {})
	// O defeito que motiva o código no fecho. A mensagem tem de nomear o esperado E o
	// encontrado: dizer só "aninhamento invertido" manda o leitor procurar sozinho.
	tagA := "MLETX" + "-A03"
	tagB := "MLETX" + "-B05"
	src := "// #region [" + tagA + "]: x\nput()\n// #endregion [" + tagB + "]"
	v, msg := checkRegionPairHonored(src, regionNode(), "", nil, testkit.WithComments(nil))
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if !strings.Contains(msg, tagB) || !strings.Contains(msg, tagA) {
		t.Errorf("a mensagem tem de citar os DOIS códigos, veio: %s", msg)
	}
	if !strings.Contains(msg, "linha 3") && !strings.Contains(msg, "line 3") {
		t.Errorf("a mensagem tem de apontar a linha, veio: %s", msg)
	}
}

func TestRegionPairSemFechoFalha(t *testing.T) {
	t.Run("RPHRG-B04: An opened region that is never closed fails", func(t *testing.T) {})
	tagA := "MLETX" + "-A03"
	v, msg := checkRegionPairHonored("// #region ["+tagA+"]: x\nput()", regionNode(), "", nil, testkit.WithComments(nil))
	if v != Fail || (!strings.Contains(msg, "nunca fechada") && !strings.Contains(msg, "never closed")) {
		t.Fatalf("esperava Fail com 'nunca fechada' / 'never closed', veio %v — %s", v, msg)
	}
}

func TestRegionPairFechoOrfaoFalha(t *testing.T) {
	t.Run("RPHRG-B05: An end region marker without an opening marker fails as an orphan close", func(t *testing.T) {})
	tagA := "MLETX" + "-A03"
	src := "put()\n// #endregion [" + tagA + "]"
	v, msg := checkRegionPairHonored(src, regionNode(), "", nil, testkit.WithComments(nil))
	if v != Fail {
		t.Fatalf("fecho órfão deveria reprovar, veio: %v", v)
	}
	if !strings.Contains(msg, tagA) {
		t.Errorf("a mensagem deve conter o código do fecho órfão, veio: %s", msg)
	}
	if !strings.Contains(msg, "linha 2") && !strings.Contains(msg, "line 2") {
		t.Errorf("a mensagem deve apontar a linha do fecho órfão, veio: %s", msg)
	}
}

func TestRegionPairErrosOrdenadosPorLinha(t *testing.T) {
	t.Run("RPHRG-B07: Multiple pairing errors are ordered sequentially by line number", func(t *testing.T) {})
	tagA := "MLETX" + "-A03"
	tagB := "MLETX" + "-B05"
	tagC := "MLETX" + "-C01"
	src := "// #endregion [" + tagA + "]\n// #region [" + tagB + "]\n// #endregion [" + tagC + "]\n// #region [" + tagA + "]"
	v, msg := checkRegionPairHonored(src, regionNode(), "", nil, testkit.WithComments(nil))
	if v != Fail {
		t.Fatalf("múltiplos erros deveriam reprovar, veio %v", v)
	}
	idxL1 := strings.Index(msg, " 1:")
	idxL3 := strings.Index(msg, " 3:")
	idxL4 := strings.Index(msg, " 4:")
	if idxL1 == -1 || idxL3 == -1 || idxL4 == -1 || !(idxL1 < idxL3 && idxL3 < idxL4) {
		t.Errorf("erros devem estar ordenados por linha; msg: %s", msg)
	}
}

func TestRegionPairSoOlhaCodigoETeste(t *testing.T) {
	t.Run("RPHRG-B01: Confronting an artifact that is neither code nor test skips", func(t *testing.T) {})
	t.Run("RPHRG-X02: Region markers inside specifications or documentation are ignored", func(t *testing.T) {})
	// Uma spec em markdown pode conter o texto `#region` num exemplo; cobrar pareamento
	// dela transformaria documentação em erro.
	spec := mapx.Node{ID: "src/Tela.spec.md", Kind: mapx.KindSpec}
	tagA := "MLETX" + "-A03"
	v, _ := checkRegionPairHonored("// #region ["+tagA+"]: exemplo sem fecho", spec, "", nil, testkit.WithComments(nil))
	if v != Skip {
		t.Fatalf("spec tem de ser Skip, veio %v", v)
	}
}
