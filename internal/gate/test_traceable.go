package gate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/scan"
)

// teste-rastreavel: um teste ligado a uma feature precisa DIZER o que prova.
//
// O buraco que ele fecha é sutil e passou duas vezes num projeto real: o teste EXISTE,
// PASSA e cobre o comportamento certo — mas não cita nenhum código de cenário. Para todo
// gate relacional ele é invisível:
//
//   - `feature-test-match` confronta cenário↔teste POR CÓDIGO: sem código no teste, ele
//     reporta os cenários como não implementados, acusando o arquivo errado;
//   - `trinca-completa` vê o arquivo e dá a peça por presente;
//   - `tests-green` vê a execução passar.
//
// O resultado é o pior dos dois mundos: o trabalho foi feito e o pipeline diz que não.
// Quem for consertar escreve um segundo teste do mesmo comportamento, porque não tem como
// saber que o primeiro já o cobria.
//
// Medido: num projeto real, um teste compartilhado por 6 telas usava `PH1A`/`PH2A` onde
// as specs declaravam `PHA1`/`PHA2` — dígito e letra trocados. As três telas tinham teste
// e NENHUMA estava rastreada; a busca por código não as alcançava. No mesmo projeto, o
// teste de uma store cobria 5 comportamentos sem citar um código sequer.
//
// A régua é a mais fraca possível de propósito: basta UM código da unidade aparecer no
// arquivo. Não se cobra um código por caso de teste — isso é do `feature-test-match`, que
// já o faz por cenário. Aqui a pergunta é só "este teste se declara?".
func checkTestTraceable(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindTest {
		return Skip, i18n.T("gate.test_traceable.skip_not_test")
	}
	if g == nil {
		return pendingNoMap()
	}

	// Só se cobra de teste que PROVA uma feature. Um teste sem feature ligada não tem
	// cenário a citar, e exigir código dele seria pedir referência a nada.
	feature, temFeature := featureProvenByTest(n, g)
	if !temFeature {
		return Skip, i18n.T("gate.test_traceable.skip_no_linked_feature")
	}

	codigos := featureCodes(root, feature)
	if len(codigos) == 0 {
		return Skip, i18n.T("gate.test_traceable.skip_no_scenario_codes")
	}

	for _, c := range codigos {
		if strings.Contains(content, c) {
			return Pass, ""
		}
	}
	return Fail, i18n.T("gate.test_traceable.test_untraceable",
		feature, firstOnes(codigos, 3), codigos[0])
}

// featureProvenByTest acha a feature de onde parte a aresta `tested-by` para este teste.
func featureProvenByTest(n mapx.Node, g *mapx.Graph) (string, bool) {
	for _, e := range g.Edges {
		if e.Type == mapx.EdgeTestedBy && e.To == n.ID {
			return e.From, true
		}
	}
	return "", false
}

// featureCodes lê os códigos de cenário que a feature declara.
func featureCodes(root, feature string) []string {
	b, err := os.ReadFile(filepath.Join(root, feature))
	if err != nil {
		return nil
	}
	vistos := map[string]bool{}
	var out []string
	for _, c := range scan.ScenarioCodeRE().FindAllString(string(b), -1) {
		if !vistos[c] {
			vistos[c] = true
			out = append(out, c)
		}
	}
	return out
}

func firstOnes(xs []string, n int) string {
	if len(xs) > n {
		xs = xs[:n]
	}
	return strings.Join(xs, ", ")
}
