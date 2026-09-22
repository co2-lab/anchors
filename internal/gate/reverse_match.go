package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- A VOLTA DE CADA PAR DA TRINCA ---
//
// O `spec-feature-match` pergunta "toda regra tem cenário?" e o `feature-test-match`
// pergunta "todo cenário tem teste?". Os dois percorrem a ORIGEM procurando o destino, e
// nenhum percorre o destino perguntando se a origem ainda existe.
//
// A assimetria tem um custo medido. Um revert apagou a regra `DTSTD-B10` — spec, código e
// teste —, e a `.feature` manteve o cenário dela: no ponto do revert aquele cenário ainda
// não existia, e um PR paralelo o reintroduziu sem conflito de git. Sondado neste
// repositório, contra o caso exato:
//
//	spec→feature   cenário órfão (B10 sem regra)   Pass, mensagem VAZIA
//	feature→teste  teste provando B10 ausente      nem mencionado
//
// O cenário ficou na feature afirmando um comportamento que a spec não decide mais, e o
// teste ficou verde provando uma regra que ninguém declara. Todos os gates verdes.
//
// É o mesmo par assimétrico que o carimbo da documentação e o `realizes` já mostraram —
// "se a spec tem hash e a doc não a referencia, algo quebrou; e vice-versa". A volta é
// uma PERGUNTA PRÓPRIA, e por isso ganha gate próprio em vez de virar mais um veredito do
// gate de ida: quem lê "1 cenário sem regra" precisa saber que a acusação é sobre a
// feature, não sobre a spec.

// checkFeatureSpecMatch — a volta do `spec-feature-match`: todo CENÁRIO da feature
// corresponde a uma regra que a spec ainda DEFINE?
//
// O confronto é sobre a feature, e o alvo é ela: quem sobrou órfão foi o cenário.
func checkFeatureSpecMatch(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, i18n.T("gate.feature_spec.skip_not_feature")
	}
	if g == nil {
		return pendingNoMap()
	}
	scenarios := parseFeatureScenarios(content)
	if len(scenarios) == 0 {
		return Skip, i18n.T("gate.feature_spec.skip_no_scenarios")
	}

	// As specs que ESTA feature cobre: a aresta `covered-by` sai da spec e chega aqui.
	var specPaths []string
	for _, e := range g.Neighbors(n.ID).In {
		if e.Type == mapx.EdgeCoveredBy {
			specPaths = append(specPaths, e.From)
		}
	}
	if len(specPaths) == 0 {
		// Sem spec ligada não há o que confrontar, e quem cobra a EXISTÊNCIA da spec é a
		// co-location. Pending para não duplicar a acusação.
		return Pending, i18n.T("gate.feature_spec.pending_no_spec")
	}

	// A união do que as specs ligadas definem. União e não interseção: uma feature pode
	// cobrir mais de uma spec, e o cenário que casa QUALQUER uma delas tem dono.
	declared := map[string]bool{}
	for _, sp := range specPaths {
		b, err := os.ReadFile(filepath.Join(root, sp))
		if err != nil {
			continue
		}
		for _, c := range definedRequirements(string(b)) {
			declared[c] = true
		}
	}
	if len(declared) == 0 {
		return Pending, i18n.T("gate.feature_spec.pending_no_requirements")
	}

	var orphans []string
	seen := map[string]bool{}
	for _, sc := range scenarios {
		for _, c := range sc.Codes {
			// SÓ os códigos da própria unidade. Um cenário pode citar a regra de outra
			// unidade para dizer contra o que ele roda, e cobrar isso da spec local faria
			// o gate pedir o impossível — é o erro que o `scenario-coverage` já mediu, com
			// 18 cenários cobrados de uma spec que definia 6.
			if n.Code != "" && !strings.HasPrefix(c, n.Code+"-") {
				continue
			}
			if declared[c] || seen[c] {
				continue
			}
			seen[c] = true
			orphans = append(orphans, c)
		}
	}
	if len(orphans) == 0 {
		return Pass, ""
	}
	sort.Strings(orphans)
	return Fail, fmt.Sprintf(i18n.T("gate.feature_spec.orphan"),
		len(orphans), strings.Join(orphans, ", "))
}

// checkTestFeatureMatch — a volta do `feature-test-match`: todo CÓDIGO que o teste diz
// provar corresponde a um cenário que a feature ainda declara?
//
// O alvo é o TESTE, e a diferença importa: um teste verde provando regra revertida é pior
// que um teste ausente, porque ele ATESTA. A suíte passa, a cobertura sobe, e o número diz
// que um comportamento está provado quando ninguém mais o decide.
func checkTestFeatureMatch(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest {
		return Skip, i18n.T("gate.test_feature.skip_not_test")
	}
	if g == nil {
		return pendingNoMap()
	}

	// As features que ESTE teste exercita: `tested-by` sai da feature e chega aqui.
	var featPaths []string
	for _, e := range g.Neighbors(n.ID).In {
		if e.Type == mapx.EdgeTestedBy {
			featPaths = append(featPaths, e.From)
		}
	}
	if len(featPaths) == 0 {
		return Pending, i18n.T("gate.test_feature.pending_no_feature")
	}

	declared := map[string]bool{}
	unidades := map[string]bool{}
	for _, fp := range featPaths {
		b, err := os.ReadFile(filepath.Join(root, fp))
		if err != nil {
			continue
		}
		for _, sc := range parseFeatureScenarios(string(b)) {
			for _, c := range sc.Codes {
				declared[c] = true
				if u, _, ok := strings.Cut(c, "-"); ok {
					unidades[u] = true
				}
			}
		}
	}
	if len(declared) == 0 {
		return Pending, i18n.T("gate.test_feature.pending_no_scenarios")
	}

	// COMENTÁRIO FORA, pela mesma régua do `feature-test-match`: um código citado em
	// comentário é REFERÊNCIA, não prova.
	corpo := stripLineComments(content)

	var orphans []string
	seen := map[string]bool{}
	for _, m := range anyCodeRE.FindAllString(corpo, -1) {
		// Só as unidades que estas features governam: um teste cita código de outras
		// unidades ao montar fixture, e cobrá-los aqui seria pedir que a feature local
		// declarasse cenário alheio.
		u, _, ok := strings.Cut(m, "-")
		if !ok || !unidades[u] {
			continue
		}
		if declared[m] || seen[m] {
			continue
		}
		seen[m] = true
		orphans = append(orphans, m)
	}
	if len(orphans) == 0 {
		return Pass, ""
	}
	sort.Strings(orphans)
	return Fail, fmt.Sprintf(i18n.T("gate.test_feature.orphan"),
		len(orphans), strings.Join(orphans, ", "))
}
