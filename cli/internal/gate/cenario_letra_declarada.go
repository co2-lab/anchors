package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cenario-letra-declarada: a letra do código de cenário existe no vocabulário.
//
// O `rule_types` declara as letras que o projeto reconhece, e o `rule-types` cobra que
// as SEÇÕES da spec usem as certas. Ninguém cobrava o mesmo do código que o cenário
// carrega — e uma letra inventada passa por todos os outros gates: o código casa
// consigo mesmo entre feature e teste, a trinca está completa, e nada percebe que
// `-SG06` não significa nada.
//
// Medido no projeto que originou o gate: 18 códigos com cinco letras não declaradas
// (`N`, `SG`, `OR`, `RC`, `FP`), sempre acompanhadas de tags igualmente fora do
// vocabulário (`@navegacao`, `@estado-dado`). O padrão é reconhecível: alguém precisou
// de uma natureza que o projeto não tinha e a inventou no lugar de declará-la.
//
// Os dois consertos são legítimos, e a escolha é do projeto: declarar a letra em
// `rule_types` (se a natureza faz falta) ou remapear o cenário para uma letra existente
// (se era só apelido de uma que já havia). O gate não escolhe — mostra o que está fora.
//
// PENDING e não FAIL: descobrir que uma natureza foi usada sem registro é informação,
// e decidir entre adotá-la ou remapeá-la é trabalho de quem conhece o domínio.
func checkCenarioLetraDeclarada(content string, n mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, ""
	}
	if cfg == nil || len(cfg.RuleTypes) == 0 {
		return Skip, "o projeto não declara `rule_types` — sem vocabulário, toda letra é válida"
	}
	validas := map[string]bool{}
	for _, rt := range cfg.RuleTypes {
		validas[strings.ToUpper(strings.TrimSpace(rt.Letter))] = true
	}

	// NÃO usa parseFeatureScenarios: o regex dele é montado com as letras DECLARADAS,
	// então um código de letra inventada é justamente o que ele descarta — o gate
	// ficaria cego para o que existe para achar. Aqui a varredura é sobre a forma
	// (`@ABCDX-XX99`), sem consultar o vocabulário, e o vocabulário entra depois, ao
	// julgar cada letra encontrada.
	todos := codigoDeCenarioLivreRE().FindAllStringSubmatch(content, -1)
	if len(todos) == 0 {
		return Skip, "feature sem cenário com código — nada a confrontar"
	}

	// agrupa por letra: repetir a mesma letra desconhecida em oito cenários vira uma
	// linha, não oito — o que o leitor precisa saber é QUAIS letras estão fora.
	porLetra := map[string][]string{}
	for _, m := range todos {
		l := m[2]
		if validas[l] {
			continue
		}
		raiz := m[1] + "-" + m[2] + m[3]
		if !contemStr(porLetra[l], raiz) {
			porLetra[l] = append(porLetra[l], raiz)
		}
	}
	if len(porLetra) == 0 {
		return Pass, ""
	}

	letras := make([]string, 0, len(porLetra))
	for l := range porLetra {
		letras = append(letras, l)
	}
	sort.Strings(letras)
	var partes []string
	for _, l := range letras {
		cods := porLetra[l]
		sort.Strings(cods)
		partes = append(partes, fmt.Sprintf("`%s` em %s", l, strings.Join(cods, ", ")))
	}
	return Pending, fmt.Sprintf("%d letra(s) de código fora do vocabulário declarado em `rule_types`: %s.\n\n"+
		"Uma letra que o projeto não registrou não significa nada para quem lê depois — e passa por todos os "+
		"outros gates, porque o código casa consigo mesmo entre feature e teste. Duas saídas, e as duas são "+
		"legítimas: declare a letra em `rule_types` (se a natureza faz falta ao vocabulário) ou remapeie o "+
		"cenário para uma letra existente (se era apelido de uma que já havia)",
		len(letras), strings.Join(partes, "; "))
}

func contemStr(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

// codigoDeCenarioLivreRE casa a FORMA de um código de cenário sem consultar o
// vocabulário — é o que permite enxergar a letra que o projeto não declarou.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func codigoDeCenarioLivreRE() *regexp.Regexp {
	return regexp.MustCompile(`@([A-Z0-9]` + config.CodeLengthPattern() + `)-([A-Z]{1,2})(\d{2})(?:#\d{2})?\b`)
}
