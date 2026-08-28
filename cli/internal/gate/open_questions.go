package gate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// open-questions-resolved: uma spec com pergunta em aberto NÃO está pronta para implementar.
//
// A classe de defeito é a AMBIGUIDADE NÃO RESOLVIDA, e ela é a mais barata de evitar e a
// mais cara de descobrir tarde. O caminho é sempre o mesmo: a spec não decide algo que o
// código precisa; quem implementa escolhe uma leitura defensável e segue; a escolha nunca
// é confrontada com quem tinha a resposta; o produto sai com a leitura errada. Nenhum gate
// pega, porque todas as peças existem e se referenciam — o defeito é uma decisão que
// ninguém tomou.
//
// O que este gate acrescenta ao "não chute, registre e reporte" que o `anchors work` já
// diz: um LUGAR declarado para o registro. Sem lugar, "registrar" vira comentário no PR
// que morre no merge. Com lugar, a pergunta é um item de trabalho visível, e a spec só
// fica implementável quando a seção esvazia.
//
// O ciclo pretendido:
//  1. quem escreve a spec percebe o que não sabe → escreve em `## Decisões em aberto`;
//  2. o gate acusa enquanto houver item → a spec não passa por pronta;
//  3. a pergunta é levada a quem decide (a seção é o lugar onde ela está visível);
//  4. a resposta VIRA REGRA (com código) e o item some da seção.
//
// A pergunta não é apagada: ela é PROMOVIDA a regra. Item que some sem regra nova é
// pergunta varrida para debaixo do tapete — por isso o gate pede que o item aponte para
// onde a resposta vai morar.
//
// Fechar a seção com "nenhuma" é o opt-out honesto (CONCEPT §5.1): afirma que se olhou e
// não há dúvida, em vez de omitir a seção e deixar a questão sem ter sido feita.
func checkOpenQuestions(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "só a spec decide — a pergunta em aberto é dela"
	}

	corpo, achou := seçãoDecisõesEmAberto(content)
	if !achou {
		// A seção é DECLARADA no catálogo (`Default: true` em todo preset), e o esqueleto
		// que o `anchors new` emite já a traz fechada com `nenhuma`. Então "não tem a
		// seção" quer dizer uma de duas coisas, e elas não se confundem:
		//
		//   • a spec nasceu ANTES da prática — dívida, não defeito. Medido no repositório
		//     que originou o gate: 586 specs sem a seção, e o preset gerando-a hoje.
		//   • alguém a apagou — e aí o gate precisa dizer, porque apagar é exatamente o
		//     movimento que a mensagem de falha chama de "varrer para debaixo do tapete".
		//
		// Nenhum sinal técnico separa as duas, então o veredito é `Pending`: nomeia a
		// dívida sem acusar. `Skip` seria pior — afirma "não tive o que confrontar", e o
		// leitor entende "não se aplica"; um gate BLOQUEANTE com ✓0 ✗0 ~586 parece
		// vigilante e não vigia nada.
		//
		// Não IMPEDE (ver `Result.Impede`): dívida de migração não barra promoção.
		return Pending, "a spec não declara `## Decisões em aberto`, que o preset desta " +
			"camada emite por padrão — não dá para saber se ela decidiu tudo ou se a seção " +
			"foi apagada. Ao tocá-la, feche com `nenhuma` (que AFIRMA que se olhou) ou " +
			"escreva o que ainda não está decidido"
	}

	itens := itensEmAberto(corpo)
	if len(itens) == 0 {
		return Pass, ""
	}

	var nums []string
	for i, it := range itens {
		if i >= 3 {
			nums = append(nums, fmt.Sprintf("… e mais %d", len(itens)-3))
			break
		}
		nums = append(nums, "«"+resumir(it)+"»")
	}
	// PENDING, não FAIL. A distinção é o ponto do gate, e errá-la o inverte:
	//
	// `Fail` é o veredito de quem fez algo ERRADO. Item em aberto não é erro — é a spec
	// dizendo a verdade sobre o que ainda não sabe, que é exatamente o comportamento que
	// este gate existe para produzir. Medido num E2E real: uma spec registrou 3 decisões
	// pendentes legítimas, levou ✗, e o agente leu como punição à honestidade. O caminho
	// mais curto para o verde passava a ser APAGAR a seção — o oposto do que o gate quer,
	// e nada distinguiria a spec que não tem dúvida da que escondeu as suas.
	//
	// `Pending` é o terceiro estado (QUALITY): nomeia o que falta e o risco, sem fingir
	// aprovação nem acusar quem declarou — quem escreveu não é penalizado por ter olhado.
	//
	// Mas pendência não basta para IMPEDIR: `Pending` sozinho nunca barrou, e por um
	// tempo este comentário afirmou que barrava. A consequência apareceu num E2E: três
	// specs fecharam com "✓ pode promover" carregando decisões que ninguém tomou, e a
	// fila enfileirou `code` para elas.
	//
	// O que barra é `Result.Impede`, marcado para este gate: "há decisão por tomar" é
	// diferente de "não tive o que confrontar", e só o primeiro impede a promoção. Culpa
	// e prontidão são eixos distintos — declarar o que não se sabe não é defeito, e ainda
	// assim não deixa a spec pronta.
	return Pending, fmt.Sprintf("%d decisão(ões) que a spec ainda NÃO tomou, e o código vai "+
		"precisar: %s. Registrá-las aqui é o certo — o defeito seria decidir por conta "+
		"própria na hora de implementar. O caminho de saída é UM: leve a pergunta a quem "+
		"decide e PROMOVA a resposta a regra (com código). Apagar o item sem regra nova é "+
		"varrer a pergunta para debaixo do tapete, e o gate não distingue isso de ter "+
		"decidido — quem apaga assume a decisão silenciosamente",
		len(itens), strings.Join(nums, ", "))
}

// seçãoDecisõesEmAberto extrai o corpo da seção. Aceita as variações de escrita que
// aparecem na prática — o gate não deve reprovar por causa de um acento ou de um sinônimo.
var decisõesRE = regexp.MustCompile(`(?im)^#{1,4}\s*(decis(ões|oes|ão|ao)\s+em\s+aberto|em\s+aberto|quest(ões|oes)\s+em\s+aberto|open\s+questions|pend(ências|encias)\s+de\s+decis(ão|ao))\b[^\n]*\n`)

func seçãoDecisõesEmAberto(content string) (string, bool) {
	loc := decisõesRE.FindStringIndex(content)
	if loc == nil {
		return "", false
	}
	resto := content[loc[1]:]
	// a seção vai até o próximo cabeçalho de mesmo nível ou acima
	if fim := regexp.MustCompile(`(?m)^#{1,4}\s`).FindStringIndex(resto); fim != nil {
		resto = resto[:fim[0]]
	}
	return resto, true
}

// fechadaRE reconhece o fechamento honesto: a afirmação de que se olhou e não há dúvida.
var fechadaRE = regexp.MustCompile(`(?i)^\s*[-*]?\s*(nenhuma|nenhum|none|n/?a|sem\s+pend[êe]ncias?|vazio|—|-)\s*\.?\s*$`)

// itemRE: um item é uma linha de lista ou de tabela. Prosa solta explicando a seção não
// conta — senão o texto de abertura viraria uma pendência fantasma.
var itemRE = regexp.MustCompile(`(?m)^\s*(?:[-*+]\s+|\d+[.)]\s+|\|)`)

func itensEmAberto(corpo string) []string {
	var itens []string
	for _, linha := range strings.Split(corpo, "\n") {
		if strings.TrimSpace(linha) == "" || fechadaRE.MatchString(linha) {
			continue
		}
		if !itemRE.MatchString(linha) {
			continue
		}
		// linha separadora de tabela (`|---|---|`) não é item
		if regexp.MustCompile(`^\s*\|[\s:|-]+\|?\s*$`).MatchString(linha) {
			continue
		}
		// cabeçalho de tabela também não
		if regexp.MustCompile(`(?i)^\s*\|\s*(pergunta|d[úu]vida|quest(ão|ao)|item|decis(ão|ao))\b`).MatchString(linha) {
			continue
		}
		// item RESOLVIDO fica no histórico sem bloquear: marcado com [x] ou riscado.
		// A resposta vira regra, mas o rastro de que a pergunta existiu tem valor.
		if regexp.MustCompile(`(?i)^\s*[-*+]\s*\[x\]|~~`).MatchString(linha) {
			continue
		}
		itens = append(itens, strings.TrimSpace(linha))
	}
	return itens
}

// resumir corta o item para caber na mensagem do gate sem perder o assunto.
func resumir(s string) string {
	s = strings.TrimLeft(s, "-*+|0123456789.) \t")
	s = strings.TrimSpace(strings.SplitN(s, "|", 2)[0])
	if len([]rune(s)) > 60 {
		return string([]rune(s)[:57]) + "…"
	}
	return s
}
