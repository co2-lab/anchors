package gate

import (
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cenario-identidade: dois cenários da MESMA feature não podem ter o mesmo código.
//
// Uma regra legitimamente tem vários cenários — caminho feliz e alternativos. O que
// não pode é os dois serem indistinguíveis: com o mesmo código, nada liga UM cenário
// a UM teste, e os gates relacionais comparam N títulos contra o mesmo teste (no
// máximo um casa; os outros viram divergência que ninguém consegue resolver).
//
// A saída é o SUFIXO de cenário: `@USBPX-B01#01` e `@USBPX-B01#02` mantêm a regra
// legível no prefixo e dão identidade a cada caso.
//
// Medido no projeto que originou o gate: 204 códigos repetidos, e 65 deles com tags
// de TIPO conflitantes no mesmo código (`@estado` num cenário, `@comportamento` no
// outro) — sinal de código emprestado, não de regra com dois caminhos. Um desses
// (`GLFLX-S02`) provou ser defeito: a spec o define como "Foco", e o segundo cenário
// falava de `onChangeText`, comportamento que não tinha código próprio.
//
// PENDING e não FAIL: numerar cenários é migração, e o gate nasce sobre uma base que
// não conhecia a notação. Quem já migrou fica verde; quem não, vê o que falta.
func checkScenarioIdentity(content string, n mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, "" // só confronta features
	}
	cenarios := parseFeatureScenarios(content)
	if len(cenarios) == 0 {
		return Skip, i18n.T("gate.scenario_identity.skip_no_scenarios_with_code")
	}

	// Agrupa por código COMPLETO (com sufixo, se houver): é ele que identifica o
	// cenário. Dois cenários com `#01` e `#02` são distintos; dois sem sufixo, não.
	porCodigo := map[string][]string{}
	for _, c := range cenarios {
		porCodigo[c.Code] = append(porCodigo[c.Code], c.Title)
	}

	var repetidos []string
	for cod, titulos := range porCodigo {
		if len(titulos) < 2 {
			continue
		}
		repetidos = append(repetidos, i18n.T("gate.scenario_identity.repeated_item",
			cod, len(titulos), strings.Join(summarizeTitles(titulos), " / ")))
	}
	if len(repetidos) == 0 {
		return Pass, ""
	}
	sort.Strings(repetidos)
	return Pending, i18n.T("gate.scenario_identity.pending_repeated_codes",
		len(repetidos), strings.Join(repetidos, "; "),
		firstCode(repetidos), firstCode(repetidos))
}

// summarizeTitles encurta os títulos para a mensagem caber — o endereço é o código,
// o título só ajuda a reconhecer qual cenário é qual.
func summarizeTitles(ts []string) []string {
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		if len([]rune(t)) > 40 {
			t = string([]rune(t)[:40]) + "…"
		}
		out = append(out, t)
	}
	return out
}

// firstCode extrai o código do primeiro achado, para o exemplo da mensagem
// falar do caso REAL do projeto em vez de um genérico.
func firstCode(repetidos []string) string {
	if len(repetidos) == 0 {
		return "XXXXX-B01"
	}
	if i := strings.Index(repetidos[0], " "); i > 0 {
		return repetidos[0][:i]
	}
	return repetidos[0]
}
