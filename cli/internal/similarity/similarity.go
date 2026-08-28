// Package similarity mede o quanto dois textos falam da MESMA coisa, ponderando
// cada palavra pelo quanto ela DISCRIMINA dentro de um conjunto.
//
// Por que não contar palavras em comum. Num conjunto de cenários de um mesmo
// componente, palavras como "componente", "deve" e "props" aparecem em quase
// todos: elas não distinguem nada, e ainda assim inflam qualquer contagem crua.
// As que importam são as raras — "borderRadius", "onChange", "opacidade".
//
// A régua é TF-IDF: um termo presente em todo o conjunto pesa 0 (ruído puro), e
// um termo que aparece uma vez pesa muito. O score é a interseção ponderada
// sobre a união ponderada (Jaccard ponderado), em [0,1] e determinístico.
//
// Portado de uma lib TypeScript em produção (`business-logic/similarity.ts`), onde
// resolve o mesmo problema com descrição de lançamento bancário: "JMIX.com* PAC
// 1223" e "JMIX.com* NETFLIX 4491" parecem iguais pelo prefixo do gateway e são
// séries diferentes — o peso desconta o comum e realça o distintivo.
//
// Por que não cosseno, e por que não uma lib. As libs Go de TF-IDF (go-nlp/tfidf,
// dkgv/go-tf-idf) fazem similaridade de COSSENO: tratam cada texto como vetor e
// medem o ângulo. A diferença em relação ao Jaccard ponderado é a FREQUÊNCIA — o
// cosseno distingue quem repete um termo cinco vezes de quem o cita uma. Em título
// de cenário (5 a 12 palavras, nenhuma repetida) essa distinção não existe, e
// medido nos pares reais deste projeto as duas réguas ORDENAM IGUAL: mesmo assunto
// sempre acima de assunto diferente, mudando só a escala. A `go-nlp` ainda puxaria
// `gorgonia.org/tensor` — álgebra matricial para comparar duas frases.
//
// PURO: sem I/O, sem IA, sem estado. O mesmo par devolve o mesmo número sempre —
// condição para um gate poder usá-lo sem virar oráculo.
package similarity

import (
	"math"
	"regexp"
	"strings"
)

// minTokenLen: token de 1 caractere não carrega identidade — em português são
// artigos e preposições, e no código são variáveis de uma letra.
const minTokenLen = 2

var (
	separadorRE = regexp.MustCompile(`[^\p{L}\p{N}]+`)
	soDigitosRE = regexp.MustCompile(`^\d+$`)
)

// Tokenize normaliza e quebra um texto em tokens comparáveis.
//
// Descarta número puro: `9999` em "borderRadius 9999" é o valor, não o assunto, e
// dois textos sobre coisas diferentes que citem o mesmo número casariam por acaso.
func Tokenize(texto string) []string {
	if texto == "" {
		return nil
	}
	var out []string
	for _, t := range separadorRE.Split(strings.ToUpper(texto), -1) {
		if len([]rune(t)) < minTokenLen || soDigitosRE.MatchString(t) {
			continue
		}
		out = append(out, t)
	}
	return out
}

// Weights calcula o peso IDF de cada token sobre um corpus.
//
// `idf = ln(N / df)`: token em TODAS as descrições pesa 0; token em uma só pesa
// ln(N). É o que separa "componente" (em todo cenário do arquivo) de "onChange"
// (só no que trata do evento).
func Weights(corpus []string) map[string]float64 {
	n := len(corpus)
	w := map[string]float64{}
	if n == 0 {
		return w
	}
	df := map[string]int{}
	for _, texto := range corpus {
		vistos := map[string]bool{}
		for _, tok := range Tokenize(texto) {
			if vistos[tok] {
				continue
			}
			vistos[tok] = true
			df[tok]++
		}
	}
	for tok, freq := range df {
		w[tok] = math.Log(float64(n) / float64(freq))
	}
	return w
}

// Score devolve a similaridade de `a` e `b` em [0,1], ponderada pelos pesos.
//
// Fallback: quando o corpus é homogêneo (todos os pesos 0 — ex.: um único texto,
// sem contraste), não há como discriminar ruído de sinal, e o score cai para
// Jaccard NÃO-ponderado. Sem isso, um corpus de um item devolveria 0 para textos
// idênticos — o pior erro possível para quem lê o resultado.
func Score(a, b string, weights map[string]float64) float64 {
	ta, tb := conjunto(a), conjunto(b)
	if len(ta) == 0 || len(tb) == 0 {
		return 0
	}
	var inter, union float64
	var interN, unionN int
	for tok := range unir(ta, tb) {
		p := weights[tok]
		union += p
		unionN++
		if ta[tok] && tb[tok] {
			inter += p
			interN++
		}
	}
	if union == 0 {
		if unionN == 0 {
			return 0
		}
		return float64(interN) / float64(unionN)
	}
	return inter / union
}

// Cosseno é a SEGUNDA régua: trata cada texto como vetor (uma dimensão por token,
// valor = peso IDF) e mede o ângulo entre eles.
//
// Ela conta uma coisa diferente do Jaccard. O Jaccard pergunta "que fração do
// vocabulário é comum?" e por isso PUNE o texto mais longo: um teste que detalha
// mais que o cenário perde pontos por cada palavra a mais, ainda que diga o mesmo.
// O cosseno normaliza pelo comprimento — dois textos de tamanhos diferentes sobre
// o mesmo assunto ficam próximos.
//
// É por isso que as duas juntas valem mais que qualquer uma sozinha: onde
// discordam, o par é limítrofe, e dizer isso é mais honesto que fingir um veredito.
func Cosseno(a, b string, weights map[string]float64) float64 {
	ta, tb := conjunto(a), conjunto(b)
	if len(ta) == 0 || len(tb) == 0 {
		return 0
	}
	var num, na, nb float64
	for tok := range ta {
		p := weights[tok]
		na += p * p
		if tb[tok] {
			num += p * p
		}
	}
	for tok := range tb {
		p := weights[tok]
		nb += p * p
	}
	if na == 0 || nb == 0 {
		// Corpus homogêneo: sem peso para projetar, cai na contagem de tokens.
		var comum int
		for tok := range ta {
			if tb[tok] {
				comum++
			}
		}
		den := math.Sqrt(float64(len(ta))) * math.Sqrt(float64(len(tb)))
		if den == 0 {
			return 0
		}
		return float64(comum) / den
	}
	return num / (math.Sqrt(na) * math.Sqrt(nb))
}

func conjunto(s string) map[string]bool {
	out := map[string]bool{}
	for _, t := range Tokenize(s) {
		out[t] = true
	}
	return out
}

func unir(a, b map[string]bool) map[string]bool {
	out := map[string]bool{}
	for k := range a {
		out[k] = true
	}
	for k := range b {
		out[k] = true
	}
	return out
}

// ── Classificação ───────────────────────────────────────────────────────────

// Veredito classifica um par de textos que DEVERIAM ser iguais e não são.
//
// A régua de igualdade é exata (o texto do cenário e o do teste têm de bater);
// a similaridade não decide nada — ela diz ao leitor QUAL é o conserto:
//
//   - Similar   → os dois falam da mesma coisa com palavras diferentes. Reescreva
//     um dos lados para casar com o outro; o conteúdo já está certo.
//   - Divergente→ falam de coisas diferentes. Ou o teste prova outra coisa, ou o
//     cenário descreve o que o código não faz mais. Não é reescrita —
//     é decidir qual dos dois está desatualizado.
type Veredito int

const (
	Identico Veredito = iota
	Similar
	Limitrofe
	Divergente
)

func (v Veredito) String() string {
	switch v {
	case Identico:
		return "idêntico"
	case Similar:
		return "similar"
	case Limitrofe:
		return "limítrofe"
	default:
		return "divergente"
	}
}

// limiarSimilar: acima disto os textos falam do mesmo assunto.
//
// 0,5 e não 0,7 (o corte de duplicata do projeto de origem) porque aqui o custo
// do erro é assimétrico: classificar de "similar" o que é divergente só sugere o
// conserto errado a quem já vai ler os dois textos; o contrário — mandar reescrever
// alguém que precisa DECIDIR qual lado está velho — desperdiça a leitura.
const limiarSimilar = 0.5

// Classifica compara dois textos que deveriam ser iguais.
//
// Compõe três camadas, na ordem em que descartam mais barato (é o desenho da lib
// de origem, onde o pré-filtro por data/valor evita medir texto à toa):
//
//  1. IGUALDADE normalizada — mesma sequência de tokens, ignorando caixa e
//     pontuação. Aqui não há o que medir.
//  2. SIMILARIDADE ponderada — Jaccard sobre IDF do corpus.
//  3. REFORÇO estrutural — pares que compartilham um token RARO (peso alto)
//     sobem de "divergente" para "similar" mesmo com score baixo: um termo que só
//     aparece nesses dois textos é evidência forte de que tratam do mesmo assunto,
//     e ele se dilui num Jaccard cheio de palavras comuns.
func Classifica(a, b string, weights map[string]float64) (Veredito, float64) {
	ta, tb := Tokenize(a), Tokenize(b)
	if len(ta) > 0 && strings.Join(ta, " ") == strings.Join(tb, " ") {
		return Identico, 1
	}
	// DUAS RÉGUAS, e a discordância entre elas é informação.
	//
	// O Jaccard pune o texto mais longo (cada palavra a mais entra na união); o
	// cosseno normaliza pelo comprimento. Quando as duas concordam, o veredito é
	// firme. Quando divergem, o par é LIMÍTROFE — e dizer isso a quem vai
	// consertar vale mais que escolher um número e fingir certeza.
	j := Score(a, b, weights)
	c := Cosseno(a, b, weights)
	jSim, cSim := j >= limiarSimilar, c >= limiarSimilar

	switch {
	case jSim && cSim:
		return Similar, maior(j, c)
	case jSim != cSim:
		return Limitrofe, maior(j, c)
	case compartilhaTokenRaro(ta, tb, weights):
		// Nenhuma das duas alcançou o limiar, mas há um termo RARO em comum — só
		// esses dois textos o usam no arquivo inteiro. É evidência estrutural, de
		// natureza diferente da contagem, e por isso ela desempata.
		return Similar, maior(j, c)
	}
	return Divergente, maior(j, c)
}

func maior(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// pesoRaro: acima disto o token é distintivo o bastante para, sozinho, ligar dois
// textos. `ln(4)` ≈ 1,39 — um termo em no máximo 1/4 do corpus.
const pesoRaro = 1.38

func compartilhaTokenRaro(a, b []string, weights map[string]float64) bool {
	sb := map[string]bool{}
	for _, t := range b {
		sb[t] = true
	}
	for _, t := range a {
		if sb[t] && weights[t] >= pesoRaro {
			return true
		}
	}
	return false
}
