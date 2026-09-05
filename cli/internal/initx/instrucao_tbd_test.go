package initx

import (
	"strings"
	"testing"
)

// TODO gate de julgamento que interroga uma PEÇA (código ou teste) carrega a instrução
// do `@TBD`.
//
// Este é o teste que importa, e é uma GUARDA CONTRA O FUTURO: o defeito do #76 não foi
// escrever o `ask:` errado, foi não haver nada que exigisse a instrução. Um gate de
// julgamento novo nasceria sem ela, e o defeito voltaria — com outro nome, no gate que
// ninguém lembrou de revisar.
//
// A exceção é declarada e justificada abaixo: um gate cujo alvo é a DISPENSA (`@no-test`)
// não sofre o problema, porque quem declara dispensa permanente afirma que a prova existe
// em outro lugar — não que ela falta.
func TestJudgment_todoGateSobrePecaCarregaAInstrucaoTBD(t *testing.T) {
	// Gates de julgamento cujo alvo NÃO é uma peça que pode estar sob `@TBD`.
	//
	// `no-test-prova-real` interroga a prova apontada por `@no-test` — uma dispensa
	// PERMANENTE, cujo sentido é "o teste vive em outro lugar". `@TBD` diz o oposto
	// ("falta escrever"), e os dois não se sobrepõem: o gate já é filtrado por
	// `Requires: "@no-test"`.
	isentos := map[string]string{
		"no-test-prova-real": "interroga a prova de uma dispensa permanente (@no-test), " +
			"não uma peça que falta",
	}

	// TODOS os artefatos: os gates são condicionais por escolha, e um subconjunto
	// deixaria de fora justamente o gate que ninguém revisou.
	gates := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true}, false)

	var judgments int
	for _, g := range gates {
		if !g.IsJudgment() {
			continue
		}
		judgments++
		if motivo, isento := isentos[g.Name]; isento {
			if strings.Contains(g.Ask, "@TBD") {
				t.Errorf("gate %q está na lista de isentos (%s) e MENCIONA @TBD — "+
					"a lista ou a instrução está errada", g.Name, motivo)
			}
			continue
		}
		if !strings.Contains(g.Ask, "@TBD") {
			t.Errorf("o gate de julgamento %q não instrui sobre `@TBD`.\n\n"+
				"Sem isso, diante de uma spec que declara a peça por desenvolver a "+
				"pergunta não tem sujeito, e a saída fácil é dar `pass` para destravar — "+
				"um carimbo que fica no mapa parecendo verificação real.\n\n"+
				"Acrescente `instrucaoTBD(\"o código\")` (ou \"o teste\") ao fim do "+
				"`Ask:`, ou declare o gate na lista `isentos` deste teste com o motivo.\n\n"+
				"ask atual: %s", g.Name, g.Ask)
		}
	}
	// Se um refactor removesse todos os gates de julgamento, o loop acima passaria sem
	// conferir nada — e o teste reportaria verde sobre o vazio.
	if judgments == 0 {
		t.Fatal("nenhum gate de julgamento no default: o teste não confrontou nada")
	}
}

// A instrução DIZ O QUE RESPONDER, e diz para NÃO dar pass.
//
// Um texto que só mencionasse `@TBD` passaria no teste acima sem resolver o problema: o
// ponto inteiro é que a saída não seja `pass`.
func TestInstrucaoTBD_proibeOPassEMandaNomearAAusencia(t *testing.T) {
	got := instrucaoTBD("o código")
	for _, exigido := range []string{"@TBD", "DISPENSADO", "pass"} {
		if !strings.Contains(got, exigido) {
			t.Errorf("a instrução não menciona %q:\n%s", exigido, got)
		}
	}
	// A peça entra no texto: sem isso a instrução falaria de "o código" num gate que
	// interroga teste.
	if !strings.Contains(instrucaoTBD("o teste"), "o teste") {
		t.Error("a instrução não usa a peça que recebeu")
	}
	// O `@TBD` desatualizado é a outra metade: uma peça que passou a existir com o
	// marcador ainda no arquivo faz todo gate que o lê dispensar o que devia cobrar.
	if !strings.Contains(got, "desatualizado") {
		t.Errorf("a instrução não cobre o `@TBD` que deixou de ser verdade:\n%s", got)
	}
}
