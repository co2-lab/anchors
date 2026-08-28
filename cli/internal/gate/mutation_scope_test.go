package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func nodeComEscopos(iso, full mapx.MutationScope) mapx.Node {
	return mapx.Node{
		Kind: mapx.KindCode, ID: "x.ts", Rev: "r1",
		Signal: &mapx.TestSignal{
			MutantsKilled: full.Killed, MutantsSurvived: full.Survived,
			MutationScore: full.Score, AtRev: "r1",
			MutationByScope: map[string]mapx.MutationScope{"isolated": iso, "full": full},
		},
	}
}

// O caso real que motivou a mudança: um átomo de UI com 8% isolado e 77% completo.
// Olhando só o completo, ele parecia saudável — e 92% dos mutantes sobrevivem aos
// próprios testes.
func TestEscoposRevelamAcoplamento(t *testing.T) {
	n := nodeComEscopos(
		mapx.MutationScope{Killed: 7, Survived: 81, Score: 8},
		mapx.MutationScope{Killed: 68, Survived: 20, Score: 77})

	v, msg := checkMutationScore("", n)
	if v != Fail {
		t.Fatalf("veredito %v, queria Fail — o isolado está em 8%%", v)
	}
	for _, quer := range []string{"isolado 8%", "completo 77%", "delta 69p", "dependentes"} {
		if !strings.Contains(msg, quer) {
			t.Errorf("mensagem sem %q:\n%s", quer, msg)
		}
	}
}

// O veredito é sobre o ISOLADO: um completo alto não salva uma unidade que não se
// prova sozinha. É a diferença entre "alguém prova" e "o teste desta unidade prova".
func TestVereditoSegueOIsoladoNaoOCompleto(t *testing.T) {
	// completo em 100%, isolado em 30% → ainda reprova
	n := nodeComEscopos(
		mapx.MutationScope{Killed: 3, Survived: 7, Score: 30},
		mapx.MutationScope{Killed: 10, Survived: 0, Score: 100})
	if v, _ := checkMutationScore("", n); v != Fail {
		t.Errorf("veredito %v, queria Fail — completo alto não compensa isolado baixo", v)
	}
}

// Unidade que se prova sozinha passa, ainda que o completo seja maior.
func TestIsoladoAcimaDoLimiarPassa(t *testing.T) {
	n := nodeComEscopos(
		mapx.MutationScope{Killed: 8, Survived: 2, Score: 80},
		mapx.MutationScope{Killed: 9, Survived: 1, Score: 90})
	if v, msg := checkMutationScore("", n); v != Pass {
		t.Errorf("veredito %v (%s), queria Pass", v, msg)
	}
}

// Delta baixo com score baixo é outro diagnóstico: não é acoplamento, é asserção
// faltando — e a mensagem precisa dizer isso, senão o autor procura no lugar errado.
func TestDeltaBaixoApontaAssercaoNaoAcoplamento(t *testing.T) {
	n := nodeComEscopos(
		mapx.MutationScope{Killed: 2, Survived: 8, Score: 20},
		mapx.MutationScope{Killed: 3, Survived: 7, Score: 25})
	_, msg := checkMutationScore("", n)
	if !strings.Contains(msg, "asserção") {
		t.Errorf("mensagem não distingue asserção de acoplamento:\n%s", msg)
	}
}

// Retrocompatível: quem ingere sem escopo continua com a régua antiga sobre o total.
func TestSemEscoposUsaOTotal(t *testing.T) {
	n := mapx.Node{Kind: mapx.KindCode, ID: "x.ts", Rev: "r1",
		Signal: &mapx.TestSignal{MutantsKilled: 8, MutantsSurvived: 2, MutationScore: 80, AtRev: "r1"}}
	if v, msg := checkMutationScore("", n); v != Pass {
		t.Errorf("veredito %v (%s), queria Pass sem escopos", v, msg)
	}
}
