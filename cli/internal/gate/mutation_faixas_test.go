package gate

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

func noComScore(score, low, high float64, sobrev int) mapx.Node {
	return mapx.Node{
		ID:   "src/regra.ts",
		Kind: mapx.KindCode,
		Signal: &mapx.TestSignal{
			MutantsKilled:   100,
			MutantsSurvived: sobrev,
			MutationScore:   score,
			MutationLow:     low,
			MutationHigh:    high,
		},
	}
}

// TestAbaixoDoAceitavelReprova — a faixa de baixo é a única que barra.
func TestAbaixoDoAceitavelReprova(t *testing.T) {
	v, detalhe := checkMutationScore("", noComScore(56, 70, 90, 142))
	if v != Fail {
		t.Fatalf("56%% com mínimo 70%% tem de reprovar; veio %v", v)
	}
	if !strings.Contains(detalhe, "70") {
		t.Errorf("o laudo precisa dizer contra qual régua reprovou: %q", detalhe)
	}
}

// TestEntreAceitavelEDesejavelNaoBarra é o conceito novo: passou, e mesmo assim aparece.
// Se isto virasse Fail, seria um limiar de 90 disfarçado — e a distinção entre "não
// pode" e "dá para melhorar" se perderia.
func TestEntreAceitavelEDesejavelNaoBarra(t *testing.T) {
	v, detalhe := checkMutationScore("", noComScore(75, 70, 90, 30))
	if v == Fail {
		t.Fatalf("75%% está acima do aceitável (70%%) — não pode reprovar")
	}
	if v != Pending {
		t.Fatalf("a faixa do meio tem de APARECER (Pending), não sumir; veio %v", v)
	}
	if !strings.Contains(detalhe, "aceitável") || !strings.Contains(detalhe, "desejável") {
		t.Errorf("o laudo precisa nomear as duas faixas: %q", detalhe)
	}
	if !strings.Contains(detalhe, "15") {
		t.Errorf("dizer QUANTO falta é o que torna o aviso acionável: %q", detalhe)
	}
}

// TestAcimaDoDesejavelPassaLimpo — sem ruído para quem já chegou lá.
func TestAcimaDoDesejavelPassaLimpo(t *testing.T) {
	v, detalhe := checkMutationScore("", noComScore(92, 70, 90, 5))
	if v != Pass || detalhe != "" {
		t.Errorf("92%% com desejável 90%% passa limpo; veio %v %q", v, detalhe)
	}
}

// TestSemDesejavelVoltaAUmLimiarSo — nenhum projeto é obrigado a adotar o conceito. Sem
// `high` no relatório, o gate se comporta como antes.
func TestSemDesejavelVoltaAUmLimiarSo(t *testing.T) {
	if v, _ := checkMutationScore("", noComScore(75, 70, 0, 30)); v != Pass {
		t.Errorf("sem desejável declarado, 75%% acima do mínimo passa limpo; veio %v", v)
	}
}

// TestReguaVemDoRelatorioNaoDoEngine — o projeto que declara 60 como aceitável tem 65
// aprovado, mesmo o default do engine sendo 70. É o que impede o framework de decidir
// o que é qualidade suficiente para todo mundo.
func TestReguaVemDoRelatorioNaoDoEngine(t *testing.T) {
	if v, _ := checkMutationScore("", noComScore(65, 60, 0, 40)); v != Pass {
		t.Errorf("com mínimo 60 declarado, 65 passa; veio %v", v)
	}
	if v, _ := checkMutationScore("", noComScore(65, 0, 0, 40)); v != Fail {
		t.Errorf("sem régua no relatório, vale o default 70 e 65 reprova; veio %v", v)
	}
}

// TestDesejavelInvalidoEhIgnorado — `high` menor que `low` é config errada do projeto;
// o gate não pode transformar isso numa faixa impossível que reprova sempre.
func TestDesejavelInvalidoEhIgnorado(t *testing.T) {
	if v, _ := checkMutationScore("", noComScore(75, 70, 50, 30)); v != Pass {
		t.Errorf("desejável abaixo do aceitável é incoerente e deve ser ignorado; veio %v", v)
	}
}

// ── escopo velho não decide veredito ───────────────────────────────────────────

func noComEscopos(rev string, iso, full mapx.MutationScope, total float64) mapx.Node {
	return mapx.Node{
		ID: "src/regra.ts", Kind: mapx.KindCode, Rev: rev,
		Signal: &mapx.TestSignal{
			MutantsKilled: 100, MutantsSurvived: 50,
			MutationScore: total, MutationLow: 70,
			AtRev:           rev,
			MutationByScope: map[string]mapx.MutationScope{"isolated": iso, "full": full},
		},
	}
}

// TestEscopoVelhoNaoDecideVeredito — o gate julga pelo ISOLADO. Com um carimbo só para
// o sinal inteiro, reingerir apenas o `full` renovava o carimbo e o isolado antigo
// pegava carona como se fosse atual: o veredito saía de um número medido contra código
// que já tinha mudado.
func TestEscopoVelhoNaoDecideVeredito(t *testing.T) {
	n := noComEscopos("rev2",
		mapx.MutationScope{Score: 30, Survived: 200, AtRev: "rev1"}, // medido antes
		mapx.MutationScope{Score: 95, Survived: 3, AtRev: "rev2"},   // atual
		95)
	v, detalhe := checkMutationScore("", n)
	if v == Fail {
		t.Fatalf("o isolado de rev1 não pode reprovar em rev2; laudo: %q", detalhe)
	}
	if strings.Contains(detalhe, "30") {
		t.Errorf("o número velho não pode aparecer no laudo: %q", detalhe)
	}
}

// TestEscopoAtualAindaDecide — a guarda não pode desligar o par quando os dois estão
// na mesma rev; senão o achado de acoplamento (isolado baixo, completo alto) some.
func TestEscopoAtualAindaDecide(t *testing.T) {
	n := noComEscopos("rev2",
		mapx.MutationScope{Score: 30, Survived: 200, AtRev: "rev2"},
		mapx.MutationScope{Score: 95, Survived: 3, AtRev: "rev2"},
		95)
	if v, _ := checkMutationScore("", n); v != Fail {
		t.Errorf("isolado 30%% na rev atual tem de reprovar; veio %v", v)
	}
}

// TestSinalAntigoSemCarimboDeEscopoContinuaValendo — sinal gravado antes do campo
// existir não tem AtRev por escopo. Tratá-lo como velho faria o gate pedir remedição
// de tudo que já estava no mapa, sem base para afirmar que está desatualizado.
func TestSinalAntigoSemCarimboDeEscopoContinuaValendo(t *testing.T) {
	n := noComEscopos("rev2",
		mapx.MutationScope{Score: 30, Survived: 200},
		mapx.MutationScope{Score: 95, Survived: 3},
		95)
	if v, _ := checkMutationScore("", n); v != Fail {
		t.Errorf("sem carimbo de escopo, o par continua valendo; veio %v", v)
	}
}

// TestLaudoNaoSeContradiz — a frase sobre o delta saía concatenada de forma fixa, então
// um delta BAIXO produzia "os dois escopos concordam … Delta alto significa …" no mesmo
// laudo. Um relatório que se contradiz não é lido: quem o lê para de confiar nele.
func TestLaudoNaoSeContradiz(t *testing.T) {
	n := noComEscopos("r",
		mapx.MutationScope{Score: 58, Survived: 175, AtRev: "r"},
		mapx.MutationScope{Score: 61, Survived: 170, AtRev: "r"},
		61)
	_, detalhe := checkMutationScore("", n)
	if strings.Contains(detalhe, "concordam") && strings.Contains(detalhe, "Delta alto") {
		t.Errorf("laudo contraditório: %q", detalhe)
	}
}
