package gate

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// --- a SPEC tem de bastar por si ---
//
// A spec é lida por duas plateias, e uma delas não tem o repositório aberto: o `docs/*.md`
// compilado herda o texto dela palavra por palavra. Quem lê a documentação para saber o
// que o sistema faz não tem uso para o nome de um arquivo de plano — e navegar até ele é
// justamente o que o mecanismo de documentação existe para eliminar.
//
// O QUE ESTE GATE NÃO COBRA. A referência COM o texto junto está certa e continua:
//
//	> *"O desenho evita expor rota de ingestão: o alarme publica no SNS nativamente."*
//
// O leitor tem o argumento em mãos. O que sobra é o ANDAIME — a frase que anuncia a
// citação e depois não cita — e é ele que manda a pessoa para fora. Medido num projeto
// real: 48 menções em 37 specs, quase todas seguidas do trecho citado.
//
// SEM PALAVRA NENHUMA. O gate não procura "plano", "ver" nem "conforme": o Anchors governa
// projetos em qualquer idioma, e um gate que casa vocabulário passa em silêncio no projeto
// escrito na outra língua — o que é pior que não existir, porque a spec PARECE protegida.
//
// O que ele casa é ESTRUTURA, e só o que o próprio Anchors define:
//
//   - o CAMINHO de um nó do mapa (`plans/0014-x.md`) escrito no corpo — o mapa é quem
//     sabe quais arquivos existem, e nenhuma língua muda um caminho
//   - a forma `{CODIGO}-R000N`, que é a identidade de revisão da doutrina
//
// INFORMATIVO. Ao contrário do `docs-fresh`, aqui não há comando que conserte: reescrever
// uma frase é trabalho de quem a escreveu, e barrar o commit por uma questão de forma
// pararia o fluxo. O gate marca, e a correção acompanha o card que já toca a spec.

// revisionRefRE casa a identidade de revisão da doutrina — forma, não palavra.
var revisionRefRE = regexp.MustCompile("`?\\b[A-Z0-9]{4,6}-R\\d{4}\\b`?")

// carriesText diz se a referência vem ACOMPANHADA do que ela afirma.
//
// A distinção é o gate inteiro: uma referência com o conteúdo junto entrega o argumento ao
// leitor; sozinha, ela o manda procurar. São duas formas de acompanhar, e ambas contam:
//
//	CITAÇÃO      > *"o desenho evita expor rota de ingestão"*
//	EXPLICAÇÃO   A `PLTFR-R0003` corrigiu o escopo desta spec: o contador de conexões
//	             saiu, porque a fonte não o expõe por réplica.
//
// A segunda foi medida num projeto real e é a mais comum: 14 dos 35 achados da primeira
// versão eram referências que EXPLICAVAM o que mudou, e foram acusadas só por não usar
// aspas. Um gate que exige a forma da citação obriga a escrever pior para passar.
//
// O sinal de que há explicação é a linha CONTINUAR depois da referência com prosa
// substantiva — e "substantiva" é medido em caracteres, não em vocabulário: o Anchors
// governa projetos em qualquer idioma, e uma lista de palavras-chave passaria em silêncio
// no projeto escrito na outra língua, que é pior que não ter o gate.
func carriesText(linhas []string, i int, ref string) bool {
	// Aspas de qualquer tradição escrita — um projeto em francês cita com « », um em
	// alemão com „ “, um em japonês com 「 」. A CRASE fica de fora: em markdown ela
	// delimita código, e a própria referência vem entre crases.
	const aspas = `"'“”„‟«»‹›「」『』`

	// A janela olha à FRENTE porque o conteúdo costuma vir depois dos dois-pontos —
	// separado por uma linha em branco, que o markdown exige antes de um blockquote ou
	// de uma tabela. Três linhas cobrem o branco mais o começo do bloco.
	for j := i; j < len(linhas) && j <= i+3; j++ {
		l := strings.TrimSpace(linhas[j])
		if strings.HasPrefix(l, ">") || strings.ContainsAny(l, aspas) {
			return true
		}
		// A TABELA é conteúdo tanto quanto uma citação: "os rótulos externos entram no
		// mapa (`STBDS-R0002`):" seguido da tabela rótulo→estado entrega ao leitor
		// exatamente o que a referência anuncia. Medido: era o achado nº 1 do projeto de
		// referência, e o conteúdo estava duas linhas abaixo.
		if j > i && strings.HasPrefix(l, "|") {
			return true
		}
		// Item de lista, pelo mesmo motivo: a referência que abre uma enumeração é
		// seguida do que ela enumera.
		if j > i && (strings.HasPrefix(l, "- ") || strings.HasPrefix(l, "* ")) {
			return true
		}
	}

	// A EXPLICAÇÃO NA PRÓPRIA LINHA, e só para a referência a REVISÃO.
	//
	// A assimetria é deliberada. Uma revisão citada com o que ela mudou junto — "a
	// `PLTFR-R0003` corrigiu o escopo: o contador de conexões saiu, porque a fonte não o
	// expõe por réplica" — informa o leitor: o código é uma etiqueta e a frase é o
	// conteúdo. Medido: 14 dos 35 achados da primeira versão eram assim, acusados só por
	// não usarem aspas.
	//
	// Um CAMINHO de arquivo não tem esse escape. "O critério está em `plans/0014.md` e
	// vale aqui" é uma linha longa que não diz qual é o critério: o caminho não é etiqueta
	// de nada, é o lugar aonde a pessoa teria de ir. Ali só o conteúdo de verdade conta —
	// a citação, a tabela, a lista — e o tamanho da prosa em volta não muda isso.
	if !revisionRefRE.MatchString(ref) {
		return false
	}
	resto := strings.ReplaceAll(linhas[i], ref, " ")
	// O código de regra é referência a conteúdo da PRÓPRIA spec ou de uma irmã já
	// nomeada, e não conta como prosa ao medir se a linha explica algo. A regex vem do
	// `rule_types.go` porque ela respeita o `code_lengths` do projeto — uma cópia local
	// congelaria o comprimento padrão e ignoraria a declaração de quem usa outro.
	resto = ruleCodeRE().ReplaceAllString(resto, " ")
	resto = strings.Trim(resto, " \t-*_`.:,;()[]—–")
	return len([]rune(resto)) >= minExplanation
}

// minExplanation é o tamanho a partir do qual o resto da linha conta como explicação.
//
// Trinta caracteres. Abaixo disso a linha é um rótulo com a referência colada — "Decidido
// pelo usuário (NTCNN-R0002)" tem 21 e não diz O QUE foi decidido; acima, há uma oração.
// O número é uma régua grosseira e assumidamente imperfeita: erra para o lado de deixar
// passar, porque este é um gate informativo e falso positivo em massa é o que faz alguém
// desligá-lo.
const minExplanation = 30

// codeFenceRE abre e fecha uma cerca. Dentro dela, um caminho é exemplo — não
// referência que manda o leitor a outro arquivo.
var codeFenceRE = regexp.MustCompile("^\\s*```")

func checkDocSelfContained(content string, n mapx.Node, _ string, g *mapx.Graph, _ *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "o corpo que vira documentação é o da SPEC"
	}
	if g == nil {
		return Skip, "sem mapa não há como saber quais caminhos são nós"
	}

	// OS CAMINHOS QUE O MAPA CONHECE. É daqui que sai a lista do que conta como
	// referência — e não de um padrão `plans/*.md` escrito à mão, que só valeria para
	// projetos que chamam a pasta de `plans`. Um projeto em espanhol com `planes/`, ou
	// qualquer outra Estrutura, é coberto pelo mesmo código.
	//
	// Só nós que NÃO são o próprio arquivo: uma spec que cita o caminho de si mesma está
	// se identificando, não mandando ninguém a lugar nenhum.
	caminhos := map[string]bool{}
	for _, no := range g.Nodes {
		if no.ID != n.ID && no.Kind != mapx.KindCode {
			caminhos[no.ID] = true
			caminhos[path.Base(no.ID)] = true
		}
	}

	linhas := strings.Split(content, "\n")
	var achados []string
	dentroDeCodigo := false
	for i, l := range linhas {
		if codeFenceRE.MatchString(l) {
			dentroDeCodigo = !dentroDeCodigo
			continue
		}
		// A própria citação não é acusada: ela É o texto que o leitor precisa.
		if dentroDeCodigo || strings.HasPrefix(strings.TrimSpace(l), ">") {
			continue
		}

		if c := citedNodePath(l, caminhos); c != "" && !carriesText(linhas, i, c) {
			achados = append(achados, fmt.Sprintf("linha %d — aponta `%s`: %s",
				i+1, c, shorten(l)))
			continue
		}
		// A revisão citada no corpo. A seção que as REGISTRA não é acusada porque suas
		// linhas são itens de lista com a explicação junto — e explicação é texto, que o
		// `quotesText` já reconhece pelas aspas. Onde não houver, o achado é legítimo:
		// um código de revisão sem uma palavra sobre o que mudou não informa ninguém.
		if m := revisionRefRE.FindString(l); m != "" && !carriesText(linhas, i, m) {
			achados = append(achados, fmt.Sprintf("linha %d — aponta %s: %s",
				i+1, strings.Trim(m, "`"), shorten(l)))
		}
	}

	if len(achados) == 0 {
		return Pass, ""
	}
	const maxAchados = 6
	if len(achados) > maxAchados {
		achados = append(achados[:maxAchados],
			fmt.Sprintf("… e mais %d", len(achados)-maxAchados))
	}
	return Fail, fmt.Sprintf(
		"%d referência(s) que só APONTAM, sem trazer o texto:\n    %s\n"+
			"  O corpo desta spec vira documentação: quem a lê em `docs/` não tem o "+
			"repositório aberto, e navegar até o arquivo apontado é o que a documentação "+
			"existe para eliminar.\n"+
			"  TRAGA O TEXTO — a referência com o trecho citado junto NÃO é acusada.",
		len(achados), strings.Join(achados, "\n    "))
}

// citedNodePath devolve o caminho de nó que a linha menciona, se houver.
func citedNodePath(linha string, caminhos map[string]bool) string {
	for _, tok := range strings.FieldsFunc(linha, func(r rune) bool {
		return r == ' ' || r == '`' || r == '(' || r == ')' || r == ',' ||
			r == ';' || r == '\t' || r == '"' || r == '\''
	}) {
		tok = strings.TrimRight(tok, ".:—–-")
		if caminhos[tok] {
			return tok
		}
	}
	return ""
}

// resumo encurta a linha para caber na mensagem sem esconder o que ela diz.
func shorten(l string) string {
	l = strings.TrimSpace(l)
	const max = 70
	if len(l) <= max {
		return l
	}
	return l[:max] + "…"
}
