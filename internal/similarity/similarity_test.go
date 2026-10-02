// @anchors
//   ref: TXSMT

package similarity

import (
	"math"
	"reflect"
	"testing"
)

// The corpus is the scenario titles of a real component: words like "componente" and
// "props" appear in almost every one and distinguish nothing.
var corpusReal = []string{
	"rounded verdadeiro aplica raio de pílula",
	"Componente não busca dados",
	"Card desabilitado exibe o conteúdo mas bloqueia o toque",
	"Toque no card dispara onPress",
	"Componente apenas reflete as props recebidas",
	"Estado vazio exibe a mensagem de nenhuma conta cadastrada",
}

// The pair that started it all: the scenario in the DOMAIN's words ("raio de pílula"),
// the test in the IMPLEMENTATION's ("borderRadius 9999").
//
// The two rulers DISAGREE here — Jaccard 0.36 and cosine 0.61 — and it is the classic case:
// the test is shorter and more specific, which the Jaccard penalises (every extra word of the
// scenario enters the union) and the cosine does not. Borderline is the honest verdict.
func TestDifferentVocabularyIsBorderline(t *testing.T) {
	t.Run("TXSMT-B06: Rulers that disagree make the pair borderline", func(t *testing.T) {})
	w := Weights(corpusReal)
	a, b := "rounded verdadeiro aplica raio de pílula", "rounded aplica borderRadius 9999"
	// The disagreement must BE a disagreement: if both rulers agreed, the verdict would be firm.
	j, c := Score(a, b, w), Cosine(a, b, w)
	if (j >= limiarSimilar) == (c >= limiarSimilar) {
		t.Fatalf("the rulers agree (jaccard %.2f, cosine %.2f) — the case is no longer borderline", j, c)
	}
	if v, score := Classify(a, b, w); v != Limitrofe {
		t.Errorf("verdict %v (score %.2f), want borderline — the rulers disagree on this pair", v, score)
	}
}

func TestScoreIsTheLargerRuler(t *testing.T) {
	t.Run("TXSMT-B09: The reported score is the larger of the two rulers", func(t *testing.T) {})
	w := Weights(corpusReal)
	a, b := "rounded verdadeiro aplica raio de pílula", "rounded aplica borderRadius 9999"
	j, c := Score(a, b, w), Cosine(a, b, w)
	if _, score := Classify(a, b, w); score != math.Max(j, c) || j == c {
		t.Errorf("score %.4f, want the larger of jaccard %.4f and cosine %.4f", score, j, c)
	}
}

// A really different subject: rewriting does not help, someone must decide which side is stale.
func TestDifferentSubjectsAreDivergent(t *testing.T) {
	t.Run("TXSMT-B08: Different subjects are divergent", func(t *testing.T) {})
	w := Weights(corpusReal)
	v, score := Classify(
		"Estado vazio exibe a mensagem de nenhuma conta cadastrada",
		"tocar copiar abre o sheet de cópia seletiva", w)
	if v != Divergente {
		t.Errorf("verdict %v (score %.2f), want divergent", v, score)
	}
}

// Equality is the ruler; similarity does not even enter the field.
func TestEqualTextIsIdentical(t *testing.T) {
	t.Run("TXSMT-B04: Texts differing only in case and punctuation are identical", func(t *testing.T) {})
	w := Weights(corpusReal)
	if v, s := Classify("Toque no card dispara onPress", "toque no card dispara onPress!", w); v != Identico || s != 1 {
		t.Errorf("verdict %v %.2f, want identical 1 — only case and punctuation differ", v, s)
	}
}

func TestBothRulersAboveTheThresholdAreSimilar(t *testing.T) {
	t.Run("TXSMT-B05: Both rulers above the threshold make the pair similar", func(t *testing.T) {})
	corpus := []string{
		"saving the draft keeps the typed title",
		"the draft keeps the typed title after saving it",
		"a tap opens the menu",
		"an empty list shows a message",
	}
	w := Weights(corpus)
	a, b := corpus[0], corpus[1]
	if j, c := Score(a, b, w), Cosine(a, b, w); j < limiarSimilar || c < limiarSimilar {
		t.Fatalf("both rulers must reach the threshold here (jaccard %.2f, cosine %.2f)", j, c)
	}
	if v, s := Classify(a, b, w); v != Similar {
		t.Errorf("verdict %v (score %.2f), want similar", v, s)
	}
}

// The structural boost: `onCustomUnit` appears only in these two texts, strong evidence of the
// same subject even with the rest of the words differing.
func TestSharedRareTokenPullsToSimilar(t *testing.T) {
	t.Run("TXSMT-B07: A shared rare word pulls a low-scoring pair to similar", func(t *testing.T) {})
	corpus := []string{
		"Editar o campo personalizado dispara onCustomUnit",
		"Toque num chip seleciona a unidade",
		"Componente reflete as props",
		"Estado vazio da lista",
	}
	w := Weights(corpus)
	v, score := Classify(
		"Editar o campo personalizado dispara onCustomUnit",
		"onCustomUnit recebe o texto digitado", w)
	if v != Similar {
		t.Errorf("verdict %v (score %.2f), want similar — they share the rare token onCustomUnit", v, score)
	}
}

// A one-item corpus has no contrast to weigh: without the fallback, two identical texts would
// score 0 — the worst possible error.
func TestHomogeneousCorpusFallsBackToPlainJaccard(t *testing.T) {
	t.Run("TXSMT-B03: A single-text corpus falls back to the unweighted count", func(t *testing.T) {})
	w := Weights([]string{"toque dispara onPress"})
	if got := Score("toque dispara onPress", "toque dispara onPress", w); got != 1 {
		t.Errorf("score %.2f for identical texts in a homogeneous corpus, want 1", got)
	}
	if got := Cosine("toque dispara onPress", "toque dispara onPress", w); math.Abs(got-1) > 1e-9 {
		t.Errorf("cosine %.2f for identical texts in a homogeneous corpus, want 1", got)
	}
}

// A pure number is a value, not a subject: two texts about different things that cite the same
// number must not match because of it.
func TestPureNumberIsNotAToken(t *testing.T) {
	t.Run("TXSMT-B01: Words are upper-cased, and one-character and digit-only words are dropped", func(t *testing.T) {})
	if toks := Tokenize("borderRadius 9999"); len(toks) != 1 || toks[0] != "BORDERRADIUS" {
		t.Errorf("tokens = %v, want only BORDERRADIUS", toks)
	}
	if toks := Tokenize("a Bc 12 d-ef"); !reflect.DeepEqual(toks, []string{"BC", "EF"}) {
		t.Errorf("tokens = %v, want [BC EF]", toks)
	}
}

func TestWeights_idfOverTheCorpus(t *testing.T) {
	t.Run("TXSMT-B02: A word in every text weighs nothing and a rarer word weighs the log of its rarity", func(t *testing.T) {})
	w := Weights([]string{"xx yy", "xx zz"})
	if w["XX"] != 0 {
		t.Errorf("XX is in every text and must weigh 0, got %v", w["XX"])
	}
	if math.Abs(w["YY"]-math.Ln2) > 1e-12 {
		t.Errorf("YY is in one of two texts and must weigh ln 2, got %v", w["YY"])
	}
}

// Identity counts EVERY letter-or-digit run, numbers and one-character words included.
// Before, it compared the scoring tokens, which drop numbers: "radius 9999" and "radius 0"
// came out identical with 1.00, and a doctrine rule restated with a different value was
// reported as a verbatim copy; "123" against "123" came out divergent with 0.00.
func TestNumbersCountForIdentity(t *testing.T) {
	t.Run("TXSMT-B10: Texts that differ only by a number are not identical", func(t *testing.T) {})
	w := Weights(corpusReal)
	if v, s := Classify("radius 9999", "radius 0", w); v == Identico {
		t.Errorf("verdict %v %.2f for texts with different numbers, want anything but identical", v, s)
	}
	if v, s := Classify("the card has 2 rows", "the card has 3 rows", w); v == Identico {
		t.Errorf("verdict %v %.2f for texts with different one-digit numbers, want anything but identical", v, s)
	}
	if v, s := Classify("123", "123", w); v != Identico || s != 1 {
		t.Errorf("verdict %v %.2f for the same number, want identical 1", v, s)
	}
}

// Two texts with nothing to tokenize are still the same text: "" against "" came out
// divergent with 0.00, because identity demanded at least one token.
func TestTokenlessEqualTextsAreIdentical(t *testing.T) {
	t.Run("TXSMT-B11: Equal texts without any word are identical", func(t *testing.T) {})
	w := Weights(corpusReal)
	for _, p := range [][2]string{{"", ""}, {"!!", "!!"}, {"!!", "??"}} {
		if v, s := Classify(p[0], p[1], w); v != Identico || s != 1 {
			t.Errorf("Classify(%q, %q) = %v %.2f, want identical 1 — no word differs", p[0], p[1], v, s)
		}
	}
}

// Both rulers fall back to the unweighted count under the SAME condition: no word of the pair
// carries weight. Before, the cosine fell back as soon as ONE side weighed nothing, while the
// Jaccard stayed weighted — the rulers disagreed by construction and the pair came out a false
// borderline (Jaccard 0, cosine 0.82).
func TestRulersFallBackTogether(t *testing.T) {
	t.Run("TXSMT-B12: The rulers fall back together, so one weightless side does not make a borderline", func(t *testing.T) {})
	corpus := []string{"componente props", "componente props onChange", "componente props valor"}
	w := Weights(corpus)
	a, b := corpus[0], corpus[1]
	j, c := Score(a, b, w), Cosine(a, b, w)
	if j != c {
		t.Errorf("jaccard %.2f and cosine %.2f disagree on a pair whose shared words weigh nothing", j, c)
	}
	if v, s := Classify(a, b, w); v != Divergente {
		t.Errorf("verdict %v %.2f, want divergent — only noise words are shared", v, s)
	}
}

func TestVerdictStringIsEnglish(t *testing.T) {
	t.Run("TXSMT-B13: A verdict prints its English name", func(t *testing.T) {})
	got := []string{Identico.String(), Similar.String(), Limitrofe.String(), Divergente.String()}
	if want := []string{"identical", "similar", "borderline", "divergent"}; !reflect.DeepEqual(got, want) {
		t.Errorf("verdict names = %v, want %v", got, want)
	}
}
