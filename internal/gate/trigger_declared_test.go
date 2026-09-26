package gate

import (
	"strings"
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

func TestTriggerDeclared_NaoSpecPula(t *testing.T) {
	t.Run("TRDCT-B01: Non-spec artifacts skip confrontation", func(t *testing.T) {})
	t.Run("TRDCT-I01: Trigger confrontation applies exclusively to specification artifacts", func(t *testing.T) {})
	t.Run("TRDCT-X02: Code and test artifacts are not checked for trigger citations", func(t *testing.T) {})
	spec := "declare `carries: personal-data`"
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindFeature, mapx.KindTest} {
		v, _ := checkTriggerDeclared(spec, mapx.Node{Kind: k}, t.TempDir(), nil, cfgComObrigacoes())
		if v != Skip {
			t.Fatalf("esperava Skip para kind %v, veio %v", k, v)
		}
	}
}

func TestTriggerDeclared_SemGatilhosPula(t *testing.T) {
	t.Run("TRDCT-B02: Specifications without cited triggers skip confrontation", func(t *testing.T) {})
	t.Run("TRDCT-X01: The gate does not require every specification to cite triggers", func(t *testing.T) {})
	spec := "Uma spec qualquer que descreve regras sem citar nenhum gatilho de compliance."
	v, _ := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes())
	if v != Skip {
		t.Fatalf("spec sem gatilhos deve pular, veio %v", v)
	}
}

func TestTriggerDeclared_SemVocabularioFicaPendente(t *testing.T) {
	t.Run("TRDCT-B03: Cited triggers return Pending when no vocabulary is declared", func(t *testing.T) {})
	t.Run("TRDCT-I02: Missing compliance packs produce Pending rather than Pass", func(t *testing.T) {})
	spec := "declare `carries: personal-data`"
	if v, _ := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, &config.Config{}); v != Pending {
		t.Errorf("sem vocabulário declarado o veredito é Pending; veio %v", v)
	}
	if v, _ := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, nil); v != Pending {
		t.Errorf("com cfg nil o veredito é Pending; veio %v", v)
	}
}

func TestTriggerDeclared_NaoConfundeOutrosCamposComGatilho(t *testing.T) {
	t.Run("TRDCT-B04: Non-trigger key-value citations are ignored", func(t *testing.T) {})
	t.Run("TRDCT-I03: Trigger keys are limited to recognized compliance predicates", func(t *testing.T) {})
	codeABCDX := "ABCD" + "X"
	refDTAXX := "DTA" + "XX"
	spec := "O header declara `layer: dao` e `code: " + codeABCDX + "`. Veja `ref: " + refDTAXX + "`."
	if v, _ := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes()); v != Skip {
		t.Errorf("campos que não são gatilho devem ser ignorados; veio %v", v)
	}
}

func TestTriggerDeclared_ProsaSemCraseIgnorada(t *testing.T) {
	t.Run("TRDCT-B05: Natural language mentions without backticks are ignored", func(t *testing.T) {})
	t.Run("TRDCT-X04: Unquoted prose text is not evaluated as symbol citations", func(t *testing.T) {})
	spec := "Esta tela carrega dado pessoal (carries: personal-data sem crases) e não cita obrigação formatada."
	v, _ := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes())
	if v != Skip {
		t.Fatalf("prosa sem crases não deve ser cobrada como citação: %v", v)
	}
}

func TestTriggerDeclared_VocabularioCorretoPassa(t *testing.T) {
	t.Run("TRDCT-B06: Valid declared trigger values pass confrontation", func(t *testing.T) {})
	t.Run("TRDCT-B07: Valid declared obligation names pass confrontation", func(t *testing.T) {})
	t.Run("TRDCT-X03: Code obligation implementation is not evaluated by this gate", func(t *testing.T) {})
	spec := "Declare `carries: personal-data` — as obrigações `lgpd-eliminacao` e " +
		"`lgpd-portabilidade` passam a exigir exclusão e exportação."
	if v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes()); v != Pass {
		t.Errorf("vocabulário correto deve passar; veio %v (%s)", v, msg)
	}
}

func TestTriggerDeclared_GatilhoInexistenteComSugestao(t *testing.T) {
	t.Run("TRDCT-B08: An undeclared trigger value fails with a nearest-match suggestion", func(t *testing.T) {})
	t.Run("TRDCT-I04: Undeclared triggers produce a blocking Fail verdict", func(t *testing.T) {})
	spec := "Declare `carries: data` no cabeçalho da spec."
	v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes())
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if !contemTudo(msg, "carries: data", "personal-data") {
		t.Errorf("esperava mensagem com sugestão mais próxima, veio: %s", msg)
	}
}

func TestTriggerDeclared_GatilhoSemMatchListaDeclarados(t *testing.T) {
	t.Run("TRDCT-B09: An undeclared trigger without close match lists declared triggers", func(t *testing.T) {})
	spec := "Declare `carries: xyz` no cabeçalho."
	v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes())
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if !contemTudo(msg, "carries: xyz", "personal-data") {
		t.Errorf("esperava lista de declarados na sugestão, veio: %s", msg)
	}
}

func TestTriggerDeclared_ObrigacaoInexistenteFalha(t *testing.T) {
	t.Run("TRDCT-B10: An undeclared obligation name fails confrontation", func(t *testing.T) {})
	spec := "Declare `carries: personal-data` — a obrigação `inexistente-lgpd` é ativada."
	v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes())
	if v != Fail {
		t.Fatalf("esperava Fail para obrigação inexistente, veio %v", v)
	}
	if !contem(msg, "inexistente-lgpd") {
		t.Errorf("deve citar o nome da obrigação inexistente, veio: %s", msg)
	}
}

func TestTriggerDeclared_DeduplicaCitacoes(t *testing.T) {
	t.Run("TRDCT-B11: Duplicate citations of triggers or obligations are deduplicated", func(t *testing.T) {})
	spec := "Declare `carries: pii` na seção A. Lembre de `carries: pii` na seção B."
	v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes())
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if strings.Contains(msg, "2 erro") || strings.Contains(msg, "2 error") {
		t.Errorf("citações idênticas devem ser deduplicadas: %s", msg)
	}
}

func TestTriggerDeclared_MultiplosErrosOrdenados(t *testing.T) {
	t.Run("TRDCT-B12: Multiple vocabulary errors are sorted deterministically", func(t *testing.T) {})
	spec := "Se carregar, declare `carries: pii` no cabeçalho — a obrigação `pii-purgavel` passa a exigir exclusão."
	v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfgComObrigacoes())
	if v != Fail {
		t.Fatalf("esperava Fail, veio %v", v)
	}
	if !contemTudo(msg, "carries: pii", "pii-purgavel") {
		t.Errorf("mensagem deve conter ambos os erros: %s", msg)
	}
	if !strings.Contains(msg, ";") {
		t.Errorf("múltiplos erros devem ser separados por ponto e vírgula: %s", msg)
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
	return strings.Contains(s, sub)
}

// O TETO DE QUATRO na lista de gatilhos declarados.
//
// Quando nenhum gatilho declarado é próximo o bastante para ser sugerido, o veredito
// lista os que existem — e para de listar no quarto. Um projeto que declara doze
// gatilhos produziria uma mensagem onde o conselho se perde na enumeração, que é o
// oposto do que a sugestão existe para fazer.
//
// Escrito depois de uma mutação SOBREVIVER: trocar o teto por um número inalcançável
// deixava a suíte inteira verde. O comportamento existia, estava certo, e nada o pinava.
func TestTriggerDeclared_ListaDeclaradosTemTetoDeQuatro(t *testing.T) {
	t.Run("TRDCT-B13: The list of declared triggers is capped at four", func(t *testing.T) {})
	cfg := &config.Config{Obligations: []config.Obligation{
		{Name: "ob-a", When: "carries: aaa"}, {Name: "ob-b", When: "carries: bbb"},
		{Name: "ob-c", When: "carries: ccc"}, {Name: "ob-d", When: "carries: ddd"},
		{Name: "ob-e", When: "carries: eee"}, {Name: "ob-f", When: "carries: fff"},
	}}
	// Um gatilho escrito que não se parece com nenhum dos declarados: força o ramo que
	// lista os declarados em vez de sugerir o mais próximo.
	spec := "Declare `carries: zzzzzzzz` nesta unidade."
	v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfg)
	if v != Fail {
		t.Fatalf("gatilho não declarado tem de reprovar, veio %v: %s", v, msg)
	}
	// O teto: no máximo QUATRO nomes de gatilho na mensagem.
	citados := 0
	for _, g := range []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff"} {
		if contem(msg, g) {
			citados++
		}
	}
	if citados > 4 {
		t.Errorf("a lista de declarados passou do teto de quatro (%d citados): %s", citados, msg)
	}
	if citados == 0 {
		t.Errorf("a mensagem não lista gatilho declarado nenhum — o conselho sumiu: %s", msg)
	}
}

// A pack that does not load is never read as "no vocabulary" nor as an unknown trigger.
func TestTriggerDeclared_PackThatDoesNotLoad(t *testing.T) {
	t.Run("TRDCT-E02: A declared pack that does not load leaves the verdict pending", func(t *testing.T) {})
	cfg := cfgComObrigacoes()
	cfg.Packs = []string{"no-such-pack-anywhere"}
	spec := "Declare `carries: personal-data`."
	v, msg := checkTriggerDeclared(spec, mapx.Node{Kind: mapx.KindSpec}, t.TempDir(), nil, cfg)
	if v != Pending || !strings.Contains(msg, "no-such-pack-anywhere") {
		t.Fatalf("a pack that does not load must leave the verdict Pending with its error, got %v (%s)", v, msg)
	}
}

func TestTriggerDeclared_reportDetails(t *testing.T) {
	cfg := &config.Config{Obligations: []config.Obligation{
		{Name: "lgpd-eliminacao", When: "carries: personal-data"},
		{Name: "lgpd-sensivel", When: "carries: sensitive-personal-data"},
	}}
	spec := mapx.Node{Kind: mapx.KindSpec}

	t.Run("TRDCT-B08: An undeclared trigger fails, naming the nearest declared one", func(t *testing.T) {
		v, msg := checkTriggerDeclared("declare `carries: personal`", spec, t.TempDir(), nil, cfg)
		if v != Fail || !strings.Contains(msg, "declared is `personal-data`") {
			t.Errorf("the shortest declared trigger containing the citation is suggested: %v (%s)", v, msg)
		}
	})

	t.Run("TRDCT-B12: Every error is aggregated in the verdict", func(t *testing.T) {
		content := "declare `carries: pii` and `renders: face`\n" +
			"a obrigação `lgpd-inexistente` e a obrigação `lgpd-outra`\n"
		v, msg := checkTriggerDeclared(content, spec, t.TempDir(), nil, cfg)
		if v != Fail {
			t.Fatalf("got %v (%s)", v, msg)
		}
		for _, want := range []string{"`carries: pii`", "`renders: face`", "`lgpd-inexistente`", "`lgpd-outra`"} {
			if !strings.Contains(msg, want) {
				t.Errorf("missing %s in: %s", want, msg)
			}
		}
		if !strings.HasPrefix(msg, "4 citation(s)") {
			t.Errorf("four errors are counted: %s", msg)
		}
	})
}

// Equally close candidates resolve the same way on every run.
func TestSuggestionIsDeterministic(t *testing.T) {
	decl := map[string]bool{"data-pii": true, "data-pie": true, "unrelated-x": true}
	first := suggestion("data", decl)
	for i := 0; i < 50; i++ {
		if got := suggestion("data", decl); got != first {
			t.Fatalf("the suggestion changed between runs: %q then %q", first, got)
		}
	}
	if !strings.Contains(first, "data-pie") {
		t.Errorf("a tie goes to the alphabetically first candidate: %q", first)
	}
}
