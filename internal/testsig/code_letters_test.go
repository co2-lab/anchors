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
	const doConfig = "SRVAXBNMDEIQFG"

	if ruleLetters != doConfig {
		t.Errorf("ruleLetters = %q, o config declara %q\n"+
			"  Cada letra que falta aqui é um cenário que passa no teste e aparece\n"+
			"  'sem teste verde' — rastreabilidade falsa na direção que faz alguém\n"+
			"  escrever o teste de novo.", ruleLetters, doConfig)
	}
}

// The project's vocabulary replaces the canonical one in the grammar builder: a declared
// letter is recognised, and an empty value (a config that did not declare it) keeps what
// was there.
//
// These assert on mustCodeRE, the builder the setters feed. CodesInCase reads a regex
// compiled once at package init (junit.go `caseCodeRE`), so today it does NOT see the
// setters -- reported as a bug, and deliberately not asserted here.
func TestSetRuleLetters(t *testing.T) {
	saved := ruleLetters
	t.Cleanup(func() { ruleLetters = saved })

	SetRuleLetters("BZ")
	if got := mustCodeRE().FindAllString("ABCDX-Z01: project letter", -1); len(got) != 1 || got[0] != "ABCDX-Z01" {
		t.Errorf("a declared letter must be recognised, got %v", got)
	}
	if got := mustCodeRE().FindAllString("ABCDX-I01: not declared", -1); len(got) != 0 {
		t.Errorf("a letter outside the project's vocabulary must not match, got %v", got)
	}
	SetRuleLetters("")
	if ruleLetters != "BZ" {
		t.Errorf("an empty value must keep the letters, got %q", ruleLetters)
	}
}

// The code length is the project's too: a 3-character code matches only once declared.
func TestSetCodeLenPattern(t *testing.T) {
	saved := codeLenPattern
	t.Cleanup(func() { codeLenPattern = saved })

	if got := mustCodeRE().FindAllString("ABC-B01: short code", -1); len(got) != 0 {
		t.Fatalf("with the default length a 3-char code must not match, got %v", got)
	}
	SetCodeLenPattern("{3}")
	if got := mustCodeRE().FindAllString("ABC-B01: short code", -1); len(got) != 1 || got[0] != "ABC-B01" {
		t.Errorf("after SetCodeLenPattern({3}) the code must match, got %v", got)
	}
	SetCodeLenPattern("")
	if codeLenPattern != "{3}" {
		t.Errorf("an empty value must keep the pattern, got %q", codeLenPattern)
	}
}

// The project's vocabulary reaches the JUnit reading: a letter declared after the package
// loaded is recognised in a case name, and one no longer declared is not.
func TestCodesInCaseFollowsTheDeclaredLetters(t *testing.T) {
	prev := ruleLetters
	defer SetRuleLetters(prev)

	SetRuleLetters("BZ")
	if got := CodesInCase("ABCDX-Z01: a project letter"); len(got) != 1 || got[0] != "ABCDX-Z01" {
		t.Fatalf("a declared letter must be read from the case name, got %v", got)
	}
	SetRuleLetters("B")
	if got := CodesInCase("ABCDX-Z01: no longer declared"); len(got) != 0 {
		t.Errorf("a letter no longer declared must not be read, got %v", got)
	}
}
