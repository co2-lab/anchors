package gate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// MarcaDecisaoEmAberto distingue, no laudo, o Pending que BARRA (há decisão por tomar) do
// que é dívida de migração (a spec nasceu antes da prática). Os dois são `Pending`, e só
// o primeiro impede a promoção e vira issue.
//
// É um marcador ESTÁVEL, não prosa: o gate que o consulta não pode depender da redação —
// nem do idioma — do texto que ele mesmo escreveu.
const MarcaDecisaoEmAberto = "[decisao-em-aberto]"

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

	// A CAMADA do nó não está no grafo, então a resolução do título cai em projeto >
	// framework. Passar o cfg é o que importa: sem ele, o gate ignoraria o
	// `section_titles` que o projeto declarou e voltaria a afirmar ausência sobre uma
	// seção que existe com outro nome.
	corpo, achou := seçãoDecisõesEmAbertoCfg(content, cfg, "")
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

	// SEM CÓDIGO a pergunta não tem identidade, e sem identidade ela não sobrevive: não
	// vira issue rastreável, não resiste a uma reescrita da spec, e não dá para dizer
	// depois que a regra `X-B03` nasceu da pergunta `X-Q01`. O achado é próprio e vem
	// ANTES da contagem — uma pergunta anônima é um problema diferente de uma pergunta
	// em aberto, e confundir os dois esconde o primeiro.
	var anonimos []string
	for _, it := range itens {
		if códigoDoItem(it) == "" {
			anonimos = append(anonimos, "«"+resumir(it)+"»")
		}
	}
	if len(anonimos) > 0 {
		return Pending, fmt.Sprintf("%d decisão(ões) em aberto SEM CÓDIGO: %s. Uma pergunta "+
			"anônima não vira issue rastreável e não sobrevive a uma reescrita da spec — e "+
			"quando a resposta virar regra, nada liga a regra à pergunta que a originou. "+
			"Dê um código a cada uma (`{CODE}-Q01`), na primeira coluna da tabela",
			len(anonimos), strings.Join(anonimos, ", "))
	}

	var nums []string
	for i, it := range itens {
		if i >= 3 {
			nums = append(nums, fmt.Sprintf("… e mais %d", len(itens)-3))
			break
		}
		// CÓDIGO e assunto juntos: o código identifica (e não muda quando alguém
		// reescreve a frase), o texto diz do que se trata. Só o código faria o relatório
		// exigir abrir a spec para saber o que se perguntou.
		nums = append(nums, códigoDoItem(it)+" «"+resumir(it)+"»")
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
	// O MARCADOR `[decisao-em-aberto]` é o que o `gate.go` consulta para saber que este
	// Pending BARRA e vira issue. Antes ele casava a prosa ("que a spec ainda NÃO tomou"),
	// e isso é uma bomba-relógio: no dia em que o laudo for traduzido — ou só reescrito —
	// o gate para de barrar e de abrir issue, em silêncio, e ninguém liga uma coisa à
	// outra.
	//
	// Marcador e não campo novo porque a assinatura do check é `(Verdict, string)` e é
	// compartilhada por dezenas de gates; mudá-la para um caso obrigaria a tocar todos.
	return Pending, fmt.Sprintf(MarcaDecisaoEmAberto+" %d decisão(ões) que a spec ainda NÃO tomou, e o código vai "+
		"precisar: %s. Registrá-las aqui é o certo — o defeito seria decidir por conta "+
		"própria na hora de implementar. O caminho de saída é UM: leve a pergunta a quem "+
		"decide e PROMOVA a resposta a regra (com código). Apagar o item sem regra nova é "+
		"varrer a pergunta para debaixo do tapete, e o gate não distingue isso de ter "+
		"decidido — quem apaga assume a decisão silenciosamente",
		len(itens), strings.Join(nums, ", "))
}

// seçãoDecisõesEmAberto extrai o corpo da seção.
//
// O TÍTULO VEM DA CONFIG (`section_titles.open`), e não de uma lista de idiomas cravada
// aqui. A lista existia — português e inglês —, e uma spec em qualquer outro idioma caía
// no ramo "a spec não declara a seção" TENDO a seção, com perguntas dentro: o gate
// afirmava ausência onde havia conteúdo, que é a falha mais cara que ele podia ter.
//
// O mecanismo já existia e este gate não o usava: `TituloDaSecao` resolve camada >
// projeto > framework, e foi criado justamente porque título de seção é CONTEÚDO do
// projeto, não mecanismo do Anchors — que não impõe léxico nem idioma (ver `dialect`).
//
// A lista de variantes continua como ÚLTIMO recurso, para o projeto que não declarou
// nada: tirá-la quebraria as specs que já existem, e o padrão do framework é português.
var decisõesRE = regexp.MustCompile(`(?im)^#{1,4}\s*(decis(ões|oes|ão|ao)\s+em\s+aberto|em\s+aberto|quest(ões|oes)\s+em\s+aberto|open\s+questions|pend(ências|encias)\s+de\s+decis(ão|ao))\b[^\n]*\n`)

// tituloDeclaradoRE monta o casador para o título que o PROJETO declarou.
func tituloDeclaradoRE(titulo string) *regexp.Regexp {
	return regexp.MustCompile(`(?im)^#{1,4}\s*` + regexp.QuoteMeta(titulo) + `\b[^\n]*\n`)
}

func seçãoDecisõesEmAberto(content string) (string, bool) {
	return seçãoDecisõesEmAbertoCfg(content, nil, "")
}

func seçãoDecisõesEmAbertoCfg(content string, cfg *config.Config, camada string) (string, bool) {
	// O título declarado VENCE: um projeto que chama a seção de "Pendências de produto"
	// não pode ser cobrado pelo nome que o framework usaria.
	if t := cfg.TituloDaSecao("open", "", camada); t != "" {
		if loc := tituloDeclaradoRE(t).FindStringIndex(content); loc != nil {
			return corpoAPartirDe(content, loc[1]), true
		}
	}
	loc := decisõesRE.FindStringIndex(content)
	if loc == nil {
		return "", false
	}
	return corpoAPartirDe(content, loc[1]), true
}

// corpoAPartirDe devolve o corpo da seção que começa em `ini`, até o próximo cabeçalho de
// mesmo nível ou acima.
func corpoAPartirDe(content string, ini int) string {
	resto := content[ini:]
	// a seção vai até o próximo cabeçalho de mesmo nível ou acima
	if fim := regexp.MustCompile(`(?m)^#{1,4}\s`).FindStringIndex(resto); fim != nil {
		resto = resto[:fim[0]]
	}
	return resto
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
		if regexp.MustCompile(`(?i)^\s*\|\s*(c[óo]digo|id|pergunta|d[úu]vida|quest(ão|ao)|item|decis(ão|ao))\b`).MatchString(linha) {
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

// códigoDoItem extrai o código de identidade da pergunta (`{CODE}-Q01`), ou vazio.
//
// A letra `Q` é a canônica da decisão em aberto, mas o projeto pode declarar outra em
// `rule_types` — daí aceitar qualquer letra do vocabulário: o gate cobra IDENTIDADE, não
// uma letra específica, e recusar `X-D01` num projeto que chama a seção de "Dúvidas"
// seria impor léxico, que a doutrina não faz.
func códigoDoItem(item string) string {
	// SÓ A PRIMEIRA CÉLULA. A linha da tabela carrega outro código na coluna "Vira" — a
	// regra que a resposta vai virar —, e ele NÃO é a identidade da pergunta: são coisas
	// opostas (o que se pergunta e o que a resposta produzirá). Ler a linha inteira daria
	// por identificada toda pergunta que declarou seu destino, que é justamente a que já
	// está mais bem escrita.
	primeira := item
	if strings.HasPrefix(strings.TrimSpace(item), "|") {
		partes := strings.SplitN(strings.TrimSpace(item), "|", 3)
		if len(partes) >= 2 {
			primeira = partes[1]
		}
	}
	return anyCodeRE.FindString(primeira)
}

// resumir corta o item para caber na mensagem do gate sem perder o assunto.
func resumir(s string) string {
	// Numa linha de tabela o ASSUNTO é a célula seguinte ao código — a primeira passou a
	// ser a identidade. Sem isto o resumo devolve o próprio código, e a mensagem fica
	// "PARCX-Q01 «PARCX-Q01»": o leitor precisa abrir a spec para saber o que se
	// perguntou, que é exatamente o que a mensagem existe para evitar.
	if strings.HasPrefix(strings.TrimSpace(s), "|") {
		celulas := strings.Split(strings.Trim(strings.TrimSpace(s), "|"), "|")
		for _, c := range celulas {
			c = strings.TrimSpace(c)
			// Pula a célula que É o código; a primeira com texto de verdade é o assunto.
			if c == "" || anyCodeRE.MatchString(c) {
				continue
			}
			s = c
			break
		}
	}
	s = strings.TrimLeft(s, "-*+|0123456789.) \t")
	s = strings.TrimSpace(strings.SplitN(s, "|", 2)[0])
	if len([]rune(s)) > 60 {
		return string([]rune(s)[:57]) + "…"
	}
	return s
}

// DecisõesEmAberto conta as decisões que a spec ainda não tomou. É o que o `doctor` usa
// para reportar a pendência como ponta sistêmica, no mesmo estatuto de `sinal-ausente`.
//
// Recebe a CONFIG porque o título da seção é do projeto: contar zero numa spec cuja seção
// se chama outra coisa diria "não há decisão pendente" sobre uma spec cheia delas — o
// silêncio que este gate existe para eliminar.
func DecisõesEmAberto(content string, cfg *config.Config, camada string) int {
	corpo, achou := seçãoDecisõesEmAbertoCfg(content, cfg, camada)
	if !achou {
		return 0
	}
	return len(itensEmAberto(corpo))
}
