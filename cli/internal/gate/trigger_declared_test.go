package gate

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func cfgComObrigacoes() *config.Config {
	return &config.Config{Obligations: []config.Obligation{
		{Name: "lgpd-eliminacao", When: "carries: personal-data"},
		{Name: "lgpd-portabilidade", When: "carries: personal-data"},
	}}
}

// TestGatilhoInexistenteEhAcusado guarda o defeito real: 47 specs mandavam declarar
// `carries: pii` e citavam a obrigação `pii-purgavel` — nenhum dos dois existia. Quem
// obedecesse escrevia um header BEM-FORMADO que não disparava obrigação nenhuma, e nada
// acusava, porque o header estava certo na forma.
func TestGatilhoInexistenteEhAcusado(t *testing.T) {
	spec := "Se carregar, declare `carries: pii` no cabeçalho — a obrigação `pii-purgavel` " +
		"passa a exigir que ela seja apagável."
	v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes())
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v (%s)", v, msg)
	}
	if !contemTudo(msg, "carries: pii", "pii-purgavel", "personal-data") {
		t.Errorf("a mensagem precisa nomear o errado E o certo; veio: %s", msg)
	}
}

func TestGatilhoDeclaradoPassa(t *testing.T) {
	spec := "Declare `carries: personal-data` — as obrigações `lgpd-eliminacao` e " +
		"`lgpd-portabilidade` passam a exigir exclusão e exportação."
	if v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes()); v != Pass {
		t.Errorf("vocabulário correto deve passar; veio %v (%s)", v, msg)
	}
}

// O gate não pode acusar `chave: valor` que não é gatilho — senão pega `layer: dao`,
// `code: ABCDX` e meio repositório junto.
func TestNaoConfundeOutrosCamposComGatilho(t *testing.T) {
	spec := "O header declara `layer: dao` e `code: ABCDX`. Veja `ref: DTAXX`."
	if v, _ := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes()); v != Skip {
		t.Errorf("campos que não são gatilho devem ser ignorados; veio %v", v)
	}
}

// Sem vocabulário declarado, o gate não tem contra o que confrontar — e o terceiro estado
// existe para isso: nomear o que falta em vez de fingir aprovação.
func TestSemPacksFicaPendente(t *testing.T) {
	spec := "declare `carries: personal-data`"
	if v, _ := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, &config.Config{}); v != Pending {
		t.Errorf("sem vocabulário declarado o veredito é Pending; veio %v", v)
	}
}

func contemTudo(s string, partes ...string) bool {
	for _, p := range partes {
		if !contem(s, p) {
			return false
		}
	}
	return true
}

func contem(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
