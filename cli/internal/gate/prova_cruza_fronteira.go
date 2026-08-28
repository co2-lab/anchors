package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// prova-cruza-fronteira: quando uma REGRA afirma relação com outra unidade — "espelha
// o backend", "fonte única com X", "os mesmos valores que Y" — o código regido tem de
// IMPORTAR aquela unidade. Senão a afirmação é prosa, e a prova é local.
//
// POR QUE ESTE GATE EXISTE. Numa auditoria de 51 divergências spec×código (app de referência,
// 2026-08), os três achados MAIS GRAVES tinham a mesma forma: dois lados definiam a
// mesma coisa, cada lado tinha teste, e cada teste confrontava a PRÓPRIA cópia.
//
//	códigos RFB      `12` = "Terreno" no app, `12` = "Casa" no backend — e o
//	                 documento entregue à Receita saía com a rubrica errada
//	filtro do saldo  `=== 'statement'` no app, `!== 'invoice'` no backend
//	fronteira mínimo `<=` em duas telas, `<` na terceira e na auditoria
//
// Em todos, 53 gates ficavam verdes — corretamente, pelas réguas que tinham. A trinca
// estava completa: regra declarada, cenário escrito, teste com o título casando. O que
// nenhum gate perguntava era se a PROVA alcança o outro lado.
//
// O caso que motivou o desenho está VIVO no repo e ainda não divergiu: `SEATX-B01`
// declara "o preço espelha o backend `orgBilling.ts` — divergir aqui mente sobre a
// cobrança", o cenário diz "os dois valores são os mesmos que o backend cobra", e o
// teste faz `expect(SEAT_PRICE.individual).toBe(15)`. Prova que é 15; não prova que é
// o mesmo que o backend cobra. Mudar o backend para 18 mantém tudo verde.
//
// O QUE O GATE NÃO FAZ, e é deliberado: não procura conceitos duplicados pelo projeto.
// Isso exigiria varrer o produto cartesiano das camadas, ou julgamento. Ele age só
// sobre o que a spec DECLARA — mesma economia do `dependency-honored`, que só cobra
// símbolo prometido entre crases. Duplicação que ninguém nomeou não é alcançada; o
// gate impede que ela VOLTE depois de nomeada.

// marcaFonteUnicaRE é a MARCAÇÃO DECLARADA da relação: `***{fonte-unica}***`.
// Negrito+itálico para não colidir com a ênfase comum (`**email**`, `**puro**`,
// que o repo já usa) — e porque no markdown renderizado ela salta da linha.
var marcaFonteUnicaRE = regexp.MustCompile(`\*\*\*?\{fonte-unica\}\*\*\*?`)

// marcaDonoRE é o CARIMBO DE DONO: `(@fonte-unica)`, que só o arquivo que DETÉM o
// conceito carrega. Os que espelham trazem apenas `{fonte-unica}` + o alvo.
//
// A assimetria é o que dá ao gate uma pergunta que a marcação sozinha não daria:
// quantos donos existem? Zero dono é conceito órfão (ninguém é a fonte, todos
// espelham algo que não está declarado); dois donos é a divergência já instalada,
// com os dois lados se achando a origem.
var marcaDonoRE = regexp.MustCompile(`\(@fonte-unica\)`)

// afirmacaoDeRelacaoRE são as formas em que uma regra afirma dependência de outra
// unidade. Português e inglês, porque a spec segue o idioma do projeto.
var afirmacaoDeRelacaoRE = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bfonte[- ]única\b`),
	regexp.MustCompile(`(?i)\bsingle[- ]source\b`),
	regexp.MustCompile(`(?i)\bespelha(m|r)?\b`),
	regexp.MustCompile(`(?i)\bmirrors?\b`),
	regexp.MustCompile(`(?i)\bo?s? mesmos? (valores?|que|do)\b`),
	regexp.MustCompile(`(?i)\bre-?export(a|ado|ada)?\b`),
	regexp.MustCompile(`(?i)\bmesma (regra|lógica|conta|fronteira|tabela)\b`),
}

// codigoDeRegraAlvoRE captura o ALVO como CÓDIGO DE REGRA (`RDSRX-B04`) — a forma
// preferida. O código é identidade estável: já é usado em toda a aplicação (spec,
// feature, teste, comentário no código), e sobrevive à refatoração que move ou
// renomeia o arquivo.
//
// O caminho de arquivo é frágil por construção, e o repo tem a prova: o comentário
// de `orgPlans.ts` citava `functions/_shared/orgBilling.ts` — caminho que não
// existia mais. Uma regra apontando para o nada não avisa ninguém.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func codigoDeRegraAlvoRE() *regexp.Regexp {
	return regexp.MustCompile("`([A-Z0-9]" + config.CodeLengthPattern() + ")-[A-Z]\\d{2}(?:#\\d{2})?`")
}

// unidadeCitadaRE captura o ARQUIVO citado na regra: um caminho com extensão, entre
// crases ou solto. Forma ACEITA (não preferida) — ver codigoDeRegraAlvoRE.
// A 2ª alternativa cobre o arquivo citado SEM caminho (`orgBilling.ts`), que é
// como a maioria das regras o nomeia. A lista de extensões inclui `json`/`yaml`
// porque contrato de provedor e schema entram aqui — e é justamente o caso que a
// dispensa `@no-cross` atende.
var unidadeCitadaRE = regexp.MustCompile("`?([\\w./@-]+/[\\w.-]+\\.\\w{2,4})`?|`([\\w.-]+\\.(?:ts|tsx|js|jsx|mjs|go|py|rb|kt|swift|json|ya?ml))`")

// linhaDeRegraRE isola a linha de uma regra catalogada (`| \x60ABCD-B01\x60 | … |`).
// A afirmação e o arquivo citado têm de estar na MESMA linha: é o que amarra a
// exigência à regra, e não a um parágrafo de prosa em volta.
// Compilado por CHAMADA e não em `var`: o comprimento do código vem da config do
// projeto (`code_lengths`), carregada DEPOIS dos globais. Um `var` congelaria o
// default e a declaração do projeto não teria efeito.
func linhaDeRegraRE() *regexp.Regexp {
	return regexp.MustCompile("(?m)^\\s*\\|\\s*`([A-Z0-9]" + config.CodeLengthPattern() + "-[A-Z]\\d{2}(?:#\\d{2})?)`\\s*\\|(.+)$")
}

// dispensaFronteiraRE é a saída declarada: a relação existe mas não é importável
// (contrato de rede, arquivo gerado, valor que vive num provedor externo).
var dispensaFronteiraRE = regexp.MustCompile(`@no-cross(?:-boundary)?\s*:\s*\S`)

func checkProvaCruzaFronteira(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "não é uma spec — a afirmação de relação vive na regra catalogada"
	}
	if g == nil {
		return Pending, "sem mapa carregado — o gate relacional precisa do grafo"
	}

	type exigencia struct {
		regra   string
		arquivo string
	}
	var exigencias []exigencia // marcação declarada: o gate COBRA
	var suspeitas []string     // só prosa: o gate AVISA
	dispensadas := 0
	donos := 0

	for _, m := range linhaDeRegraRE().FindAllStringSubmatch(content, -1) {
		regra, texto := m[1], m[2]
		if dispensaFronteiraRE.MatchString(texto) {
			dispensadas++
			continue // a relação foi declarada como não-importável, com razão
		}
		marcado := marcaFonteUnicaRE.MatchString(texto)

		// O DONO não espelha ninguém: ele É a fonte. Nada a importar.
		if marcado && marcaDonoRE.MatchString(texto) {
			donos++
			continue
		}

		if marcado {
			// Marcação declarada é a RÉGUA. O alvo tem de estar na linha — a
			// marcação diz "espelho algo", e sem o algo não há o que confrontar.
			achouAlvo := false

			// PREFERIDO: o alvo por CÓDIGO DE REGRA. Resolve pelo mapa (o nó que
			// declara `code: RDSR`), então a regra não carrega caminho de arquivo —
			// que muda — e sim a identidade da unidade, que não muda.
			for _, mm := range codigoDeRegraAlvoRE().FindAllStringSubmatch(texto, -1) {
				unidadeAlvo := mm[1]
				if unidadeAlvo == regraUnidade(regra) {
					continue // auto-referência: a própria unidade não é o outro lado
				}
				arquivos := arquivosDaUnidade(g, unidadeAlvo)
				if len(arquivos) == 0 {
					suspeitas = append(suspeitas, regra+" → `"+mm[0]+
						"` (o código da regra-alvo não resolve para nenhuma unidade no mapa — "+
						"confira o código, ou rode `anchors map build`)")
					achouAlvo = true
					continue
				}
				achouAlvo = true
				for _, a := range arquivos {
					exigencias = append(exigencias, exigencia{regra: regra, arquivo: a})
				}
			}

			for _, mm := range unidadeCitadaRE.FindAllStringSubmatch(texto, -1) {
				alvo := mm[1]
				if alvo == "" {
					alvo = mm[2]
				}
				if alvo == "" {
					continue
				}
				achouAlvo = true
				exigencias = append(exigencias, exigencia{regra: regra, arquivo: alvo})
			}
			if !achouAlvo {
				suspeitas = append(suspeitas, regra+" (marcou `{fonte-unica}` sem dizer de QUAL arquivo, "+
					"e sem o carimbo `(@fonte-unica)` que declararia esta unidade como a dona)")
			}
			continue
		}

		// Sem marcação: o VOCABULÁRIO só avisa. É o que faz a convenção ser adotada
		// em vez de esquecida — quem escreveu "espelha o X.ts" em prosa provavelmente
		// não sabia do marcador, e o aviso ensina sem barrar.
		for _, re := range afirmacaoDeRelacaoRE {
			if !re.MatchString(texto) {
				continue
			}
			for _, mm := range unidadeCitadaRE.FindAllStringSubmatch(texto, -1) {
				alvo := mm[1]
				if alvo == "" {
					alvo = mm[2]
				}
				if alvo != "" {
					suspeitas = append(suspeitas, regra+" → `"+alvo+"`")
				}
			}
			break
		}
	}
	// Uma regra que fala de relação em prosa é uma relação que nenhum gate cobra.
	// Sai como Fail — e o Anchors não tem veredito intermediário: quem decide se
	// isso barra é o `blocking` do gate na config. Nasce informativo (não barra,
	// mas aparece), e ao promover para bloqueante a prosa passa a ser cobrada.
	avisar := func() (Verdict, string) {
		sort.Strings(suspeitas)
		return Fail, "regra(s) que falam de relação com outra unidade SEM a marcação " +
			"`***{fonte-unica}***`: " + strings.Join(suspeitas, "; ") + ". Enquanto a " +
			"relação é só prosa, nenhum gate a confronta — e é assim que os dois lados " +
			"evoluem soltos, cada teste provando a própria cópia. Marque a regra (e " +
			"carimbe `(@fonte-unica)` na unidade que DETÉM o conceito)"
	}

	if len(exigencias) == 0 {
		if len(suspeitas) > 0 {
			return avisar()
		}
		// A mensagem distingue os silêncios: "não havia o que cobrar" é diferente de
		// "havia, e foi dispensado com razão", e de "esta unidade é a DONA".
		if donos > 0 {
			return Pass, fmt.Sprintf(
				"%d conceito(s) com o carimbo `(@fonte-unica)` — esta unidade É a fonte, "+
					"não espelha ninguém", donos)
		}
		if dispensadas > 0 {
			return Skip, fmt.Sprintf(
				"%d regra(s) afirmam relação, todas com `@no-cross` declarado — a relação "+
					"existe e não é importável (contrato de rede, arquivo gerado)", dispensadas)
		}
		return Skip, "nenhuma regra declara `***{fonte-unica}***` nem fala de relação com " +
			"outra unidade — nada a confrontar"
	}

	// O código regido: é ele que tem de importar a unidade citada. Exigir o import
	// no TESTE seria falso-positivo — o `BLBLX-B06` do app de referência prova a regra de verdade
	// (exercita recibo/manual/legado) e alcança a função pelo código, não por import
	// próprio. O que importa é a unidade não ter uma CÓPIA.
	var codePaths []string
	for _, e := range g.Neighbors(n.ID).Out {
		if e.Type == mapx.EdgeSpecifies {
			codePaths = append(codePaths, e.To)
		}
	}
	if len(codePaths) == 0 {
		return Pending, "spec sem código ligado (specifies) — nada a confrontar ainda"
	}

	var corpo strings.Builder
	for _, cp := range codePaths {
		b, err := os.ReadFile(filepath.Join(root, cp))
		if err != nil {
			continue
		}
		// SEM comentários, e é o ponto do gate: `orgPlans.ts` cita `orgBilling.ts`
		// duas vezes em comentário e não importa uma linha. Contar o comentário
		// aprovaria justamente o caso que motivou o gate.
		corpo.WriteString(stripLineComments(string(b)))
		corpo.WriteString("\n")
	}
	code := corpo.String()
	if strings.TrimSpace(code) == "" {
		return Pending, "código ligado não pôde ser lido — nada a confrontar"
	}

	var quebradas []string
	for _, ex := range exigencias {
		if importaUnidade(code, ex.arquivo) {
			continue
		}
		quebradas = append(quebradas, fmt.Sprintf("%s → `%s`", ex.regra, ex.arquivo))
	}
	if len(quebradas) > 0 {
		sort.Strings(quebradas)
		return Fail, "regra(s) que afirmam relação com outra unidade sem o código IMPORTAR " +
			"aquela unidade: " + strings.Join(quebradas, "; ") + ". Uma regra que diz " +
			"\"espelha\"/\"fonte única\" e não importa descreve uma CÓPIA: os dois lados " +
			"passam a evoluir soltos, cada teste prova a própria cópia, e todos os gates " +
			"ficam verdes até alguém mudar um lado. Re-exporte da unidade citada, ou " +
			"declare `@no-cross: <razão>` na linha da regra se a relação não for " +
			"importável (contrato de rede, arquivo gerado)"
	}

	// Cobrança em ordem, mas a suspeita não se perde: um arquivo pode ter uma regra
	// marcada (que passou) e outra em prosa (que ninguém confronta).
	if len(suspeitas) > 0 {
		return avisar()
	}

	return Pass, ""
}

// importaUnidade decide se `code` traz a unidade `alvo` por import. Compara pelo
// NOME BASE sem extensão, não pelo caminho: a mesma unidade é citada na spec como
// `packages/backend/business-logic/orgBilling.ts` e importada no código como
// `@backend/business-logic/orgBilling` — alias e extensão variam por projeto, o
// nome do módulo não.
func importaUnidade(code, alvo string) bool {
	base := strings.TrimSuffix(filepath.Base(alvo), filepath.Ext(alvo))
	if base == "" {
		return false
	}
	// O nome tem de aparecer numa LINHA de import/require/from — não em qualquer
	// lugar do arquivo. Um `orgBilling` mencionado numa string não é dependência.
	for _, linha := range strings.Split(code, "\n") {
		l := strings.TrimSpace(linha)
		if !strings.Contains(l, "import") && !strings.Contains(l, "require") &&
			!strings.Contains(l, "from ") {
			continue
		}
		if strings.Contains(l, base) {
			return true
		}
	}
	return false
}

// regraUnidade extrai o código da UNIDADE de um código de regra (`SEATX-B01` → `SEAT`).
func regraUnidade(regra string) string {
	if i := strings.IndexByte(regra, '-'); i > 0 {
		return regra[:i]
	}
	return regra
}

// arquivosDaUnidade resolve um código de unidade (`PPAO`) para os arquivos de CÓDIGO
// daquela trinca. É o que permite o alvo ser declarado por código de regra em vez de
// caminho de arquivo — o código é identidade estável, o caminho não.
//
// O caminho é via SPEC: só a spec carrega `code:` no header (o arquivo de código
// declara `ref:`, e o mapa não o indexa por código). Da spec, a aresta `specifies`
// leva ao código regido — que é o que precisa ser importado. Importar a spec ou o
// teste não faria sentido.
func arquivosDaUnidade(g *mapx.Graph, unidade string) []string {
	var out []string
	for _, n := range g.Nodes {
		if n.Code != unidade || n.Kind != mapx.KindSpec {
			continue
		}
		for _, e := range g.Neighbors(n.ID).Out {
			if e.Type == mapx.EdgeSpecifies {
				out = append(out, e.To)
			}
		}
	}
	// Fallback: projeto que indexe o código diretamente no nó de código.
	if len(out) == 0 {
		for _, n := range g.Nodes {
			if n.Code == unidade && n.Kind == mapx.KindCode {
				out = append(out, n.ID)
			}
		}
	}
	return out
}
