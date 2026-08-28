package gate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// domain-declared: a seção `## Domínio` declara o que a unidade ACEITA — e quem garante
// que o inválido não chega.
//
// A lacuna medida: num projeto real, 71% das specs tinham seção de regras/efeitos e
// apenas 14% diziam o que a unidade aceita. O framework inteiro é construído sobre
// catalogar EFEITOS, e todo defeito de borda encontrado em três rodadas de review
// adversarial morava no que ninguém tinha declarado.
//
// Por que `constraints` não bastava — e essa é a distinção que dá o gate:
//
//	`## Restrições` diz o que a unidade NÃO FAZ  → empurra o dever para FORA
//	`## Domínio`    diz o que a unidade ACEITA   → nomeia QUEM fica com ele
//
// Escrever mais restrições não fecha nada: CRIA órfãos, porque cada "não é meu" precisa
// de alguém do outro lado. Aconteceu, medido: três specs escreveram "não valido a chave"
// — a da regra, a do repositório, e o modelo nem mencionou. Cada uma correta. O dever
// ficou órfão, e o resultado foi perda silenciosa de dado do usuário, com todos os gates
// verdes.
//
// O QUE ESTE GATE NÃO FAZ: julgar se o domínio declarado está CERTO. Saber que `__proto__`
// é um valor perigoso para uma chave de objeto é conhecimento de domínio e de linguagem —
// não é confrontável por regex, e fingir que é produziria falso-positivo em massa ou uma
// lista de casos especiais que envelhece. O gate cobra a DECLARAÇÃO e o DONO; a qualidade
// do que foi declarado é trabalho do review.
func checkDomainDeclared(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "a fronteira de entrada é declarada na spec"
	}

	corpo, achou := seçãoDominio(content)
	if !achou {
		// Ausência NÃO é falha, pela mesma razão do `open-questions-resolved`: exigir a
		// seção de toda spec transformaria instrumento em ritual, e unidade sem entrada
		// externa (um componente que recebe props tipadas do próprio código) não tem
		// domínio a declarar. O gate cobra quem ABRIU a seção — e o preset a abre em quem
		// recebe entrada de fora.
		return Skip, "a spec não declara `## Domínio` — a seção é para quem recebe entrada " +
			"externa; abra-a quando houver valor que possa chegar errado"
	}

	linhas := linhasDeDominio(corpo)
	if len(linhas) == 0 {
		return Fail, "a seção `## Domínio` está vazia. Ou declare o que a unidade aceita " +
			"(uma linha por entrada), ou remova a seção — uma seção vazia AFIRMA que se olhou " +
			"e não se achou nada a declarar, que é diferente de não ter olhado"
	}

	var semDono []string
	for _, l := range linhas {
		if dono := donoDaEntrada(l); dono == "" {
			semDono = append(semDono, primeiraCelula(l))
		}
	}
	if len(semDono) == 0 {
		return Pass, ""
	}
	return Fail, fmt.Sprintf("%d entrada(s) sem dono na coluna `Quem garante`: %s. "+
		"Sem dono, \"fora do domínio\" é só outra forma de dizer \"não é meu problema\" — e "+
		"o problema não fica com ninguém. Foi assim que três specs declararam, cada uma "+
		"corretamente, que não validavam a mesma entrada: o dever ficou órfão e a entrada "+
		"inválida passou. Nomeie quem barra: esta unidade, o chamador, ou a interface",
		len(semDono), strings.Join(semDono, ", "))
}

// dominioRE reconhece o título da seção, com as variações que aparecem na prática.
var dominioRE = regexp.MustCompile(`(?im)^#{1,4}\s*(dom[ií]nio|domain|entradas?\s+aceitas?|fronteira\s+de\s+entrada)\b[^\n]*\n`)

func seçãoDominio(content string) (string, bool) {
	loc := dominioRE.FindStringIndex(content)
	if loc == nil {
		return "", false
	}
	resto := content[loc[1]:]
	if fim := regexp.MustCompile(`(?m)^#{1,4}\s`).FindStringIndex(resto); fim != nil {
		resto = resto[:fim[0]]
	}
	return resto, true
}

// linhasDeDominio extrai as linhas de DADOS da tabela — nem cabeçalho, nem separador,
// nem a prosa explicativa que costuma acompanhar a seção.
func linhasDeDominio(corpo string) []string {
	var out []string
	for _, l := range strings.Split(corpo, "\n") {
		t := strings.TrimSpace(l)
		if !strings.HasPrefix(t, "|") {
			continue
		}
		if separadorRE.MatchString(t) || cabecalhoDominioRE.MatchString(t) {
			continue
		}
		// linha ainda no molde (só TODOs) não conta como declaração
		if todoOnlyRE.MatchString(t) {
			continue
		}
		out = append(out, t)
	}
	return out
}

var (
	separadorRE        = regexp.MustCompile(`^\|[\s:|-]+\|?$`)
	cabecalhoDominioRE = regexp.MustCompile(`(?i)^\|\s*(entrada|input|par[âa]metro|campo|valor)\b`)
	todoOnlyRE         = regexp.MustCompile(`(?i)^(\|\s*TODO[^|]*)+\|?\s*$`)
)

// donoDaEntrada devolve a última célula (a coluna `Quem garante`), vazia se ela não
// nomeia ninguém. "não é meu" e variações NÃO contam como dono — é justamente a resposta
// que cria o órfão.
func donoDaEntrada(linha string) string {
	cels := strings.Split(strings.Trim(linha, "|"), "|")
	if len(cels) < 4 {
		return "" // a tabela precisa das 4 colunas para ter dono
	}
	dono := strings.TrimSpace(cels[len(cels)-1])
	if dono == "" || strings.HasPrefix(dono, "TODO") {
		return ""
	}
	if naoEhDonoRE.MatchString(dono) {
		return ""
	}
	return dono
}

// naoEhDonoRE reconhece a NÃO-RESPOSTA — a frase que parece preencher a coluna e não
// nomeia ninguém.
//
// Não pode ancorar no fim da linha (`$`): o autor escreve "não valido (MTVRX-X04)", com a
// referência à restrição ao lado, e a âncora deixaria passar. Testado contra o caso real
// que motivou o gate — foi exatamente assim que o dever ficou órfão em três specs, cada
// uma "declarando" que não validava. A não-resposta com citação continua sendo
// não-resposta: citar a própria restrição é dizer "não é meu" com fonte.
// Duas famílias, porque a borda é diferente em cada uma:
//   - TRAÇO puro: a célula inteira é `-` ou `—`. Ancorada no fim (`$`), porque `\b` não
//     casa depois de um traço (não é caractere de palavra) e deixaria passar.
//   - FRASE: "não valido", "ninguém", "fora de escopo" — SEM âncora no fim, porque o
//     autor escreve "não valido (MTVRX-X04)", com a referência à restrição ao lado. Foi
//     exatamente assim que o dever ficou órfão em três specs reais; citar a própria
//     restrição é dizer "não é meu" com fonte, e continua sendo não-resposta.
var naoEhDonoRE = regexp.MustCompile(`(?i)(^\s*[-—]+\s*$)|(^\s*(n/?a|ningu[ée]m|nobody|none|nenhum|` +
	`n[ãa]o\s+(é|eh|e)\s+(meu|daqui|desta)|n[ãa]o\s+valid\w*|n[ãa]o\s+se\s+aplica|` +
	`fora\s+de\s+escopo|delegado|outra\s+camada)\b)`)

func primeiraCelula(linha string) string {
	cels := strings.Split(strings.Trim(linha, "|"), "|")
	if len(cels) == 0 {
		return "?"
	}
	c := strings.TrimSpace(cels[0])
	if len(c) > 40 {
		c = c[:37] + "…"
	}
	return "`" + c + "`"
}
