package gate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// cenario-tipo-alinhado: a tag de natureza do cenário concorda com a letra do código.
//
// A natureza de uma regra é declarada DUAS vezes: na letra do código (`ABCDX-S01` é
// State) e na tag do cenário (`@estado`). Quando as duas discordam, uma delas mente —
// e nenhum outro gate percebe, porque cada um lê só um dos lados: o `rule-types`
// confronta as seções da SPEC, o `feature-test-match` confronta código e descrição.
//
// Medido no projeto que originou o gate: 46 cenários com tag e letra em desacordo.
// Ao ler, os dois casos apareceram, e distingui-los importa:
//
//   - 36 eram TAG errada. "Intro presente é renderizada em itálico" tem corpo
//     `Quando o componente é renderizado / Então devo ver` — estado puro, sem ação do
//     usuário, sob código `-S` correto e tag `@comportamento`. Conserto: a tag.
//   - 10 eram CÓDIGO emprestado. "Tocar em um chip seleciona a prioridade" é
//     comportamento de verdade, pendurado num `-S` porque a spec não catalogava o
//     evento. Conserto: código próprio na tabela de Eventos.
//
// O gate acusa os dois — decidir qual lado está errado exige ler o cenário, e é
// trabalho de quem conhece a unidade. O que ele garante é que a discordância não passe
// silenciosa.
//
// PENDING e não FAIL: o desacordo é dívida herdada em qualquer base que adote o gate
// depois de escrever features, e a correção pede julgamento caso a caso.
func checkCenarioTipoAlinhado(content string, n mapx.Node, _ string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindFeature {
		return Skip, "" // a tag e o código moram na feature
	}
	if cfg == nil || len(cfg.RuleTypes) == 0 {
		return Skip, "o projeto não declara `rule_types` — sem vocabulário, não há o que confrontar"
	}

	// Sem `tags:` declarada em nenhuma letra, o gate não tem mapa: seria adivinhar que
	// `@estado` significa S. Silenciar aqui é honesto — o projeto não pediu esta régua.
	temMapa := false
	for _, rt := range cfg.RuleTypes {
		if len(rt.Tags) > 0 {
			temMapa = true
			break
		}
	}
	if !temMapa {
		return Skip, "nenhuma letra declara `tags:` — o gate precisa do mapa tag→letra para confrontar"
	}

	cenarios := parseFeatureScenarios(content)
	if len(cenarios) == 0 {
		return Skip, "feature sem cenário com código — nada a confrontar"
	}

	var achados []string
	for _, sc := range cenarios {
		letra := letraDoCodigo(sc.Code)
		if letra == "" {
			continue
		}
		// Um cenário pode provar mais de um requisito (`@IPSBX-V02 @IPSBX-M01`), e aí a
		// tag pode estar descrevendo o SEGUNDO. Reunir as letras de todos os códigos
		// evita acusar quem co-etiqueta corretamente — medido: 4 dos 28 achados eram
		// isso, e "corrigi-los" trocaria uma classificação certa por outra.
		letras := map[string]bool{letra: true}
		for _, c := range sc.Codes {
			if l := letraDoCodigo(c); l != "" {
				letras[l] = true
			}
		}
		for _, tag := range sc.Tags {
			nome := strings.TrimPrefix(tag, "@")
			declaradas, conhecida := cfg.LetrasDaTag(nome)
			// Uma tag pode caber sob mais de uma letra (ver LetrasDaTag): basta que a
			// letra do código esteja entre elas para não haver discordância.
			if !conhecida || alguma(declaradas, letras) {
				continue
			}
			achados = append(achados, fmt.Sprintf("%s é `%s` mas o cenário se declara `@%s` (letra %s): %q",
				sc.Code, letra, nome, strings.Join(declaradas, "/"), corta(sc.Title, 48)))
		}
	}
	if len(achados) == 0 {
		return Pass, ""
	}
	sort.Strings(achados)
	return Pending, fmt.Sprintf("%d cenário(s) cuja tag de natureza discorda da letra do código: %s.\n\n"+
		"Uma das duas mente, e qual delas exige ler o cenário: se o corpo é `Quando … é renderizado / Então devo ver`, "+
		"é ESTADO e a tag está errada; se há ação do usuário (`Quando eu toco`), é COMPORTAMENTO e o código foi "+
		"emprestado de um estado vizinho — nesse caso o conserto é dar código próprio na seção que cataloga eventos, "+
		"não trocar a tag para calar o gate",
		len(achados), strings.Join(achados, "; "))
}

// letraDoCodigo extrai a letra de natureza de `ABCDX-S01` (ou `ABCDX-S01#02`). Devolve
// vazio para códigos que não seguem a forma (`ABCDX-DS-alguma-coisa`, `ABCDX-VR`).
func letraDoCodigo(code string) string {
	raiz := CodeRaiz(code)
	i := strings.LastIndex(raiz, "-")
	if i < 0 || i+1 >= len(raiz) {
		return ""
	}
	sufixo := raiz[i+1:]
	// natureza é LETRA(S) + dois dígitos; `DS-...` e `VR` não qualificam
	corte := 0
	for corte < len(sufixo) && sufixo[corte] >= 'A' && sufixo[corte] <= 'Z' {
		corte++
	}
	if corte == 0 || corte > 2 || len(sufixo)-corte != 2 {
		return ""
	}
	for _, c := range sufixo[corte:] {
		if c < '0' || c > '9' {
			return ""
		}
	}
	return sufixo[:corte]
}

func corta(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	return string([]rune(s)[:n]) + "…"
}

// alguma diz se alguma das letras declaradas para a tag está entre as do cenário.
func alguma(declaradas []string, letras map[string]bool) bool {
	for _, d := range declaradas {
		if letras[d] {
			return true
		}
	}
	return false
}
