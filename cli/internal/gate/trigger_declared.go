package gate

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/pack"
)

// trigger-declared: um gatilho de obrigação CITADO tem de existir.
//
// As obrigações de compliance são disparadas por gatilhos declarados no header
// (`carries: personal-data`). A spec que instrui o autor a declarar um gatilho está
// ensinando o vocabulário — e se ela ensina um nome que nenhum pack declara, o autor
// obedece, escreve o header, e **nenhuma obrigação dispara**.
//
// Isso não falha em lugar nenhum. O header fica bem-formado, os gates de header passam, a
// unidade parece coberta por LGPD, e a cobertura é zero. É a âncora que mente na sua forma
// mais cara: o texto que ENSINA está errado, então todo autor que confiar nele produz o
// mesmo defeito.
//
// Medido num projeto real: 47 specs de modelo mandavam declarar `carries: pii` e citavam a
// obrigação `pii-purgavel`. Nenhum dos dois existia — o vocabulário canônico dos packs era
// `carries: personal-data`, com as obrigações `lgpd-eliminacao` e `lgpd-portabilidade`. As
// 47 vieram do mesmo boilerplate, copiado de spec em spec, e nada acusou.
//
// O gate confronta o que a spec CITA contra o que os packs e o `anchors.yaml` DECLARAM.
// Não exige que a spec cite gatilho algum — cobra só quem cita.
func checkTriggerDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "o vocabulário de gatilho é ensinado pela spec — é ela que o autor lê"
	}
	citados := gatilhosCitados(content)
	if len(citados) == 0 {
		return Skip, "a spec não cita gatilho de obrigação"
	}
	declarados, obrigacoes := vocabularioDeclarado(root, cfg)
	if len(declarados) == 0 {
		// Sem packs nem obrigações, não há vocabulário contra o que confrontar. Calar é o
		// certo: acusar aqui seria cobrar de um projeto que não adotou compliance.
		return Pending, "o projeto não declara `packs:` nem `obligations:` — sem vocabulário " +
			"para confrontar, o gatilho citado não pode ser verificado"
	}

	var erros []string
	for _, c := range citados {
		if declarados[c.valor] {
			continue
		}
		erros = append(erros, fmt.Sprintf("`%s: %s` não é gatilho de nenhum pack%s",
			c.chave, c.valor, sugestao(c.valor, declarados)))
	}
	for _, ob := range obrigacoesCitadas(content) {
		if !obrigacoes[ob] {
			erros = append(erros, fmt.Sprintf("a obrigação `%s` não existe", ob))
		}
	}
	if len(erros) == 0 {
		return Pass, ""
	}
	sort.Strings(erros)
	return Fail, fmt.Sprintf("%d citação(ões) de vocabulário inexistente: %s. "+
		"Quem seguir esta spec escreve um header bem-formado que NÃO dispara obrigação "+
		"nenhuma — e nada acusa, porque o header está correto na forma. Use o vocabulário "+
		"que os packs declaram",
		len(erros), strings.Join(erros, "; "))
}

type gatilhoCitado struct{ chave, valor string }

// gatilhoRE casa a citação de um gatilho no texto: `carries: personal-data`,
// `shared-with: third-party`. A crase é o sinal de que o autor está citando um SÍMBOLO —
// prosa solta ("carrega dado pessoal") não é citação e não é cobrada.
var gatilhoRE = regexp.MustCompile("`([a-z][a-z-]*): ([a-z][a-z0-9-]*)`")

// chavesDeGatilho são os predicados que abrem uma obrigação. Fechada de propósito: sem
// isso o regex pegaria qualquer `chave: valor` entre crases (`layer: dao`, `code: ABCD`)
// e o gate acusaria meio repositório.
var chavesDeGatilho = map[string]bool{
	"carries": true, "processing": true, "renders": true, "shared-with": true,
	"retains": true, "transfers": true,
}

func gatilhosCitados(content string) []gatilhoCitado {
	visto := map[string]bool{}
	var out []gatilhoCitado
	for _, m := range gatilhoRE.FindAllStringSubmatch(content, -1) {
		if !chavesDeGatilho[m[1]] {
			continue
		}
		k := m[1] + ":" + m[2]
		if visto[k] {
			continue
		}
		visto[k] = true
		out = append(out, gatilhoCitado{chave: m[1], valor: m[2]})
	}
	return out
}

// obrigacaoRE casa a citação de uma obrigação pelo nome. Exige a palavra "obrigação" perto
// para não confundir com qualquer identificador entre crases.
var obrigacaoRE = regexp.MustCompile("(?i)obriga[çc][õo]?[eé]?s?[^`\\n]{0,40}`([a-z][a-z0-9-]{2,})`")

func obrigacoesCitadas(content string) []string {
	visto := map[string]bool{}
	var out []string
	for _, m := range obrigacaoRE.FindAllStringSubmatch(content, -1) {
		if visto[m[1]] {
			continue
		}
		visto[m[1]] = true
		out = append(out, m[1])
	}
	return out
}

// vocabularioDeclarado reúne os gatilhos e os nomes de obrigação que o projeto de fato
// tem — dos packs adotados e das obrigações escritas direto no `anchors.yaml`.
func vocabularioDeclarado(root string, cfg *config.Config) (map[string]bool, map[string]bool) {
	gatilhos, obrigacoes := map[string]bool{}, map[string]bool{}
	if cfg == nil {
		return gatilhos, obrigacoes
	}
	registra := func(when, name string) {
		if name != "" {
			obrigacoes[name] = true
		}
		// `when` vem como "carries: personal-data" — guardamos só o VALOR, que é o que a
		// spec cita depois dos dois-pontos.
		if _, v, ok := strings.Cut(when, ":"); ok {
			if v = strings.TrimSpace(v); v != "" {
				gatilhos[v] = true
			}
		}
	}
	for _, ob := range cfg.Obligations {
		registra(ob.When, ob.Name)
	}
	packs, _, err := pack.LoadAll(root, cfg.Packs, cfg.PackValues, cfg.Jurisdictions)
	if err != nil {
		return gatilhos, obrigacoes
	}
	for _, p := range packs {
		for _, ob := range p.Obligations {
			registra(ob.When, ob.Name)
		}
	}
	return gatilhos, obrigacoes
}

// sugestao aponta o gatilho declarado mais parecido — a correção costuma ser um sinônimo
// (`pii` onde o canônico é `personal-data`), e dizer qual economiza a busca.
func sugestao(errado string, declarados map[string]bool) string {
	var melhor string
	for d := range declarados {
		if strings.Contains(d, errado) || strings.Contains(errado, d) {
			if melhor == "" || len(d) < len(melhor) {
				melhor = d
			}
		}
	}
	if melhor == "" {
		var todos []string
		for d := range declarados {
			todos = append(todos, "`"+d+"`")
		}
		sort.Strings(todos)
		if len(todos) > 4 {
			todos = todos[:4]
		}
		return " (declarados: " + strings.Join(todos, ", ") + ")"
	}
	return " — o declarado é `" + melhor + "`"
}
