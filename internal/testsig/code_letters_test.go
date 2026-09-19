package testsig

import "testing"

// O DEFAULT deste pacote tem de casar as letras que o `config` declara como canônicas.
//
// Medido no blue-eyes: o `testsig` tinha `SRVAXBNMD` e o `config` tinha `SRVAXBNMDEIQF`.
// Um projeto sem `rule_types:` usa o default do `config` — que inclui `I` de Invariant —
// e este pacote não reconhecia `GLCGL-I01` no nome do caso JUnit.
//
// O sintoma: 11 casos verdes, três cenários provados, e os três invariantes aparecendo
// "sem teste verde" com teste passando. Rastreabilidade falsa na direção mais perigosa —
// o requisito parece descoberto quando está provado, e quem olhar vai escrever o teste
// de novo.
func TestCodesInCase_reconheceInvariante(t *testing.T) {
	got := CodesInCase("GLCGL-I01: item `feito` com o artefato ausente reprova")

	if len(got) != 1 || got[0] != "GLCGL-I01" {
		t.Errorf("códigos = %v — o `-I01` (Invariant) não foi reconhecido", got)
	}
}

// As letras do `config` que faltavam aqui: E, I, Q, F. Todas as quatro.
func TestCodesInCase_reconheceAsLetrasQueFaltavam(t *testing.T) {
	for _, caso := range []struct{ nome, codigo string }{
		{"ABCDX-E01: exemplo", "ABCDX-E01"},
		{"ABCDX-I01: invariante", "ABCDX-I01"},
		{"ABCDX-Q01: decisão em aberto", "ABCDX-Q01"},
		{"ABCDX-F01: fase", "ABCDX-F01"},
	} {
		got := CodesInCase(caso.nome)
		if len(got) != 1 || got[0] != caso.codigo {
			t.Errorf("%q → %v, queria [%s]", caso.nome, got, caso.codigo)
		}
	}
}

// E as canônicas continuam funcionando: a correção não pode ter trocado a lista.
func TestCodesInCase_mantemAsCanonicas(t *testing.T) {
	for _, codigo := range []string{"ABCDX-B01", "ABCDX-R01", "ABCDX-S01", "ABCDX-V01"} {
		if got := CodesInCase(codigo + ": algo"); len(got) != 1 || got[0] != codigo {
			t.Errorf("%s não foi reconhecido: %v", codigo, got)
		}
	}
}

// E a GUARDA contra a divergência voltar.
//
// Duas cópias da mesma lista divergem na primeira letra nova. O comentário do `config` já
// registra que "é a terceira vez que a lista fica para trás de uma letra nova" — e desta
// vez ficou quatro atrás. Este teste é a régua que impede a quinta.
//
// O pacote não importa `config` de propósito (ele não depende de scan nem de config), então
// a lista é comparada por VALOR — e quem mudar uma sem a outra quebra aqui.
func TestRuleLetters_naoDivergeDoConfig(t *testing.T) {
	// O valor de `config.DefaultRuleLetters`, copiado. Se este teste falhar, a lista de lá
	// mudou: alinhe a de cá e atualize esta constante — nas duas, nunca numa.
	const doConfig = "SRVAXBNMDEIQF"

	if ruleLetters != doConfig {
		t.Errorf("ruleLetters = %q, o config declara %q\n"+
			"  Cada letra que falta aqui é um cenário que passa no teste e aparece\n"+
			"  'sem teste verde' — rastreabilidade falsa na direção que faz alguém\n"+
			"  escrever o teste de novo.", ruleLetters, doConfig)
	}
}
