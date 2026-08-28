package similarity

import "testing"

// O corpus são os títulos de cenário de um componente real: palavras como
// "componente" e "props" aparecem em quase todos e não distinguem nada.
var corpusReal = []string{
	"rounded verdadeiro aplica raio de pílula",
	"Componente não busca dados",
	"Card desabilitado exibe o conteúdo mas bloqueia o toque",
	"Toque no card dispara onPress",
	"Componente apenas reflete as props recebidas",
	"Estado vazio exibe a mensagem de nenhuma conta cadastrada",
}

// O par que motivou tudo: cenário no vocabulário do DOMÍNIO ("raio de pílula"),
// teste no da IMPLEMENTAÇÃO ("borderRadius 9999").
//
// As duas réguas DISCORDAM aqui — Jaccard 0,36 e cosseno 0,61 —, e é o caso
// clássico: o teste é mais curto e específico, o que o Jaccard pune (cada palavra
// a mais do cenário entra na união) e o cosseno não. Limítrofe é o veredito
// honesto: os dois provavelmente falam do mesmo, e vale a leitura humana.
func TestVocabularioDiferenteEhLimitrofe(t *testing.T) {
	w := Weights(corpusReal)
	v, score := Classifica(
		"rounded verdadeiro aplica raio de pílula",
		"rounded aplica borderRadius 9999", w)
	if v != Limitrofe {
		t.Errorf("veredito %v (score %.2f), queria limítrofe — as réguas discordam neste par", v, score)
	}
}

// E a discordância tem de ser DE FATO uma discordância: se as duas réguas
// concordassem, o veredito seria firme. Este teste guarda a mecânica.
func TestLimitrofeEhDiscordanciaEntreAsReguas(t *testing.T) {
	w := Weights(corpusReal)
	a, b := "rounded verdadeiro aplica raio de pílula", "rounded aplica borderRadius 9999"
	j, c := Score(a, b, w), Cosseno(a, b, w)
	if (j >= limiarSimilar) == (c >= limiarSimilar) {
		t.Fatalf("as réguas concordam (jaccard %.2f, cosseno %.2f) — o caso deixou de ser limítrofe", j, c)
	}
}

// Assunto realmente diferente: aqui não adianta reescrever, alguém precisa decidir
// qual lado está velho.
func TestAssuntosDiferentesSaoDivergentes(t *testing.T) {
	w := Weights(corpusReal)
	v, score := Classifica(
		"Estado vazio exibe a mensagem de nenhuma conta cadastrada",
		"tocar copiar abre o sheet de cópia seletiva", w)
	if v != Divergente {
		t.Errorf("veredito %v (score %.2f), queria divergente", v, score)
	}
}

// Igualdade é a régua; a similaridade nem entra em campo.
func TestTextoIgualEhIdentico(t *testing.T) {
	w := Weights(corpusReal)
	if v, _ := Classifica("Toque no card dispara onPress", "toque no card dispara onPress!", w); v != Identico {
		t.Errorf("veredito %v, queria idêntico — só caixa e pontuação diferem", v)
	}
}

// O reforço estrutural: `onCustomUnit` só aparece nesses dois textos, e é evidência
// forte de mesmo assunto mesmo com o resto das palavras diferindo.
func TestTokenRaroCompartilhadoPuxaParaSimilar(t *testing.T) {
	corpus := []string{
		"Editar o campo personalizado dispara onCustomUnit",
		"Toque num chip seleciona a unidade",
		"Componente reflete as props",
		"Estado vazio da lista",
	}
	w := Weights(corpus)
	v, score := Classifica(
		"Editar o campo personalizado dispara onCustomUnit",
		"onCustomUnit recebe o texto digitado", w)
	if v != Similar {
		t.Errorf("veredito %v (score %.2f), queria similar — compartilham o token raro onCustomUnit", v, score)
	}
}

// Corpus de um item só não tem contraste para ponderar: sem o fallback, dois
// textos idênticos dariam score 0 — o pior erro possível.
func TestCorpusHomogeneoCaiParaJaccardSimples(t *testing.T) {
	w := Weights([]string{"toque dispara onPress"})
	if got := Score("toque dispara onPress", "toque dispara onPress", w); got != 1 {
		t.Errorf("score %.2f para textos idênticos em corpus homogêneo, queria 1", got)
	}
}

// Número puro é valor, não assunto: dois textos sobre coisas diferentes que citem
// o mesmo número não podem casar por isso.
func TestNumeroPuroNaoEhToken(t *testing.T) {
	if toks := Tokenize("borderRadius 9999"); len(toks) != 1 || toks[0] != "BORDERRADIUS" {
		t.Errorf("tokens = %v, queria só BORDERRADIUS", toks)
	}
}
