package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func featNode() mapx.Node { return mapx.Node{Kind: mapx.KindFeature, ID: "x.feature"} }

// Dois cenários com o MESMO código são indistinguíveis: nada liga um deles a um
// teste específico, e os gates relacionais comparam N títulos contra o mesmo teste.
func TestCenarioIdentidadeAcusaCodigoRepetido(t *testing.T) {
	v, msg := checkCenarioIdentidade(`
  @USBPX-B01 @nivel-unit
  Cenário: busca pontos por userId+month
    Então o repository é consultado

  @USBPX-B01 @nivel-unit
  Cenário: sem usuário logado nada é buscado
    Então o repository não é consultado
`, featNode(), "", nil, nil)

	if v != Pending {
		t.Fatalf("veredito %v, queria Pending", v)
	}
	if !strings.Contains(msg, "USBPX-B01") {
		t.Errorf("mensagem sem o código repetido: %s", msg)
	}
	// A mensagem tem de dizer COMO resolver, com o caso real do projeto.
	if !strings.Contains(msg, "#01") {
		t.Errorf("mensagem não ensina o sufixo: %s", msg)
	}
}

// Numerados, os dois cenários passam a ter identidade própria — que é o ponto.
func TestCenarioIdentidadeAceitaSufixo(t *testing.T) {
	v, msg := checkCenarioIdentidade(`
  @USBPX-B01#01 @nivel-unit
  Cenário: busca pontos por userId+month
    Então o repository é consultado

  @USBPX-B01#02 @nivel-unit
  Cenário: sem usuário logado nada é buscado
    Então o repository não é consultado
`, featNode(), "", nil, nil)

	if v != Pass {
		t.Errorf("veredito %v (%s), queria Pass — os sufixos distinguem os cenários", v, msg)
	}
}

// Uma regra com UM cenário é o caso comum: não pode acusar nada.
func TestCenarioIdentidadeNaoAcusaCodigoUnico(t *testing.T) {
	v, _ := checkCenarioIdentidade(`
  @SAUTX-B01 @nivel-unit
  Cenário: Hidratar carrega a sessão
    Então o usuário fica disponível

  @SAUTX-B02 @nivel-unit
  Cenário: Login popula o usuário
    Então o usuário fica disponível
`, featNode(), "", nil, nil)
	if v != Pass {
		t.Errorf("veredito %v, queria Pass", v)
	}
}

// O gate só fala de feature. Um nó de código ou spec não é assunto dele.
func TestCenarioIdentidadeSoOlhaFeature(t *testing.T) {
	v, _ := checkCenarioIdentidade("qualquer coisa", mapx.Node{Kind: mapx.KindCode}, "", nil, nil)
	if v != Skip {
		t.Errorf("veredito %v, queria Skip para nó que não é feature", v)
	}
}
