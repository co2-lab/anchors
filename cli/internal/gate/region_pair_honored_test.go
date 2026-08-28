package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func regionNode() mapx.Node { return mapx.Node{ID: "src/Tela.tsx", Kind: mapx.KindCode} }

func TestRegionPairParBemFormadoPassa(t *testing.T) {
	src := "// #region [MLETX-A03]: persiste.\nput()\n// #endregion [MLETX-A03]"
	v, msg := checkRegionPairHonored(src, regionNode(), "", nil, nil)
	if v != Pass {
		t.Fatalf("esperava Pass, veio %v — %s", v, msg)
	}
}

func TestRegionPairSemRegiaoEhSkipNaoFalha(t *testing.T) {
	// A delimitação é OPCIONAL: um projeto que não migrou nada continua válido. Se este
	// gate falhasse por ausência, promovê-lo a bloqueante reprovaria o repositório inteiro
	// no dia em que entrasse — e a região deixaria de ser adotável de forma incremental.
	v, _ := checkRegionPairHonored("// MLETX-A03: à moda antiga\nput()", regionNode(), "", nil, nil)
	if v != Skip {
		t.Fatalf("sem região tem de ser Skip, veio %v", v)
	}
}

func TestRegionPairFechoTrocadoFalhaEDizOsDois(t *testing.T) {
	// O defeito que motiva o código no fecho. A mensagem tem de nomear o esperado E o
	// encontrado: dizer só "aninhamento invertido" manda o leitor procurar sozinho.
	src := "// #region [MLETX-A03]: x\nput()\n// #endregion [MLETX-B05]"
	v, msg := checkRegionPairHonored(src, regionNode(), "", nil, nil)
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if !strings.Contains(msg, "MLETX-B05") || !strings.Contains(msg, "MLETX-A03") {
		t.Errorf("a mensagem tem de citar os DOIS códigos, veio: %s", msg)
	}
	if !strings.Contains(msg, "linha 3") {
		t.Errorf("a mensagem tem de apontar a linha, veio: %s", msg)
	}
}

func TestRegionPairSemFechoFalha(t *testing.T) {
	v, msg := checkRegionPairHonored("// #region [MLETX-A03]: x\nput()", regionNode(), "", nil, nil)
	if v != Fail || !strings.Contains(msg, "nunca fechada") {
		t.Fatalf("esperava Fail com 'nunca fechada', veio %v — %s", v, msg)
	}
}

func TestRegionPairSoOlhaCodigoETeste(t *testing.T) {
	// Uma spec em markdown pode conter o texto `#region` num exemplo; cobrar pareamento
	// dela transformaria documentação em erro.
	spec := mapx.Node{ID: "src/Tela.spec.md", Kind: mapx.KindSpec}
	v, _ := checkRegionPairHonored("// #region [MLETX-A03]: exemplo sem fecho", spec, "", nil, nil)
	if v != Skip {
		t.Fatalf("spec tem de ser Skip, veio %v", v)
	}
}

func TestHeaderConformeIsentaRoteiroExecutavel(t *testing.T) {
	// O `.yaml` do runner e2e não carrega cabeçalho `@anchors`: a identidade dele está no
	// NOME (`LOGIX-A01.yaml`) e nas `tags:`, e o formato é do runner, não nosso. Sem esta
	// isenção, trazer os roteiros para o grafo (camada e2e-flow, para que a execução deixe
	// carimbo) transformaria 717 arquivos preexistentes em defeitos retroativos.
	flow := mapx.Node{ID: "apps/mobile/.maestro/screens/auth/LoginScreen/LOGIX-A01.yaml", Kind: mapx.KindTest}
	v, msg := checkHeaderConforme("appId: com.acme.exemplo\ntags:\n  - LOGIX-A01\n", flow)
	if v != Skip {
		t.Fatalf("roteiro .yaml tem de ser Skip, veio %v — %s", v, msg)
	}
	// Um teste em código NOSSO continua obrigado: a isenção é da forma do arquivo, não do
	// papel de "ser teste".
	tsx := mapx.Node{ID: "src/features/auth/screens/LoginScreen.test.tsx", Kind: mapx.KindTest}
	if v, _ := checkHeaderConforme("describe('x', () => {})\n", tsx); v != Fail {
		t.Errorf(".test.tsx sem header tem de falhar, veio %v", v)
	}
}
