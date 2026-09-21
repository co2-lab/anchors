package flagx

import "testing"

// The three house styles say the same thing, and the grammar has to hear one statement.
// This is the finding that shaped the alias table: the tools disagree only on spelling.
func TestParse_asTresGrafiasDizemOMesmo(t *testing.T) {
	grupos := map[Op][]string{
		OpGte:        {">= 50", "gte 50", "NUM_GTE 50", "greaterThanOrEqual 50", "Greater Than Inclusive 50"},
		OpStartsWith: {`startsWith "beta"`, `STR_STARTS_WITH "beta"`, `starts with "beta"`},
		OpMatches:    {`matches "^v2"`, `REGEX "^v2"`, `matches regex "^v2"`},
		OpContains:   {`contains "beta"`, `STR_CONTAINS "beta"`, `includes "beta"`},
		OpEq:         {`= "on"`, `== "on"`, `eq "on"`, `Equal "on"`, `NUM_EQ "on"`},
	}
	for want, escritas := range grupos {
		for _, raw := range escritas {
			c, err := Parse(raw)
			if err != nil {
				t.Errorf("Parse(%q) recusou: %v", raw, err)
				continue
			}
			if c.Op != want {
				t.Errorf("Parse(%q).Op = %q, esperado %q", raw, c.Op, want)
			}
		}
	}
}

// O operando perde as aspas, e só uma camada delas.
func TestParse_operando(t *testing.T) {
	casos := map[string]string{
		`= "on"`:        "on",
		`= 'on'`:        "on",
		">= 50":         "50",
		`contains "be"`: "be",
		`= "a \"b\" c"`: `a \"b\" c`,
	}
	for raw, want := range casos {
		c, err := Parse(raw)
		if err != nil {
			t.Fatalf("Parse(%q): %v", raw, err)
		}
		if c.Operand != want {
			t.Errorf("Parse(%q).Operand = %q, esperado %q", raw, c.Operand, want)
		}
	}
}

// O caso AUSENTE é o que quebra em produção e o que ninguém escreve. Ele não compara
// contra nada, e por isso não passa pela divisão de operando.
func TestParse_ausenteEPresente(t *testing.T) {
	for _, raw := range []string{"absent", "is not set", "not set", "unset", "missing", "ABSENT"} {
		c, err := Parse(raw)
		if err != nil {
			t.Errorf("Parse(%q) recusou o caso ausente: %v", raw, err)
			continue
		}
		if c.Op != OpAbsent {
			t.Errorf("Parse(%q).Op = %q, esperado %q", raw, c.Op, OpAbsent)
		}
		if c.Operand != "" {
			t.Errorf("Parse(%q).Operand = %q, o ausente não compara contra nada", raw, c.Operand)
		}
	}
	for _, raw := range []string{"present", "is set", "any value"} {
		c, err := Parse(raw)
		if err != nil || c.Op != OpPresent {
			t.Errorf("Parse(%q) = (%v, %v), esperado OpPresent", raw, c.Op, err)
		}
	}
}

// `>` não pode engolir o `>=`, e `starts` não pode ser lido como operador desconhecido
// carregando `with "beta"` de operando.
func TestParse_prefixoMaisLongoVence(t *testing.T) {
	c, err := Parse(">= 50")
	if err != nil || c.Op != OpGte || c.Operand != "50" {
		t.Errorf(`Parse(">= 50") = (%q, %q, %v), esperado gte/50`, c.Op, c.Operand, err)
	}
	c, err = Parse(`starts with "beta"`)
	if err != nil || c.Op != OpStartsWith || c.Operand != "beta" {
		t.Errorf(`Parse("starts with ...") = (%q, %q, %v), esperado starts-with/beta`, c.Op, c.Operand, err)
	}
	c, err = Parse(`not contains "beta"`)
	if err != nil || c.Op != OpNotContain {
		t.Errorf(`Parse("not contains ...") = (%q, %v), esperado not-contains`, c.Op, err)
	}
}

// A gramática existe para RECUSAR. Uma que aceita tudo não confronta nada — é a diferença
// entre gramática fixa e prosa, que foi a decisão que originou este pacote.
func TestParse_recusaOQueNaoEntende(t *testing.T) {
	recusar := []string{
		"",
		"   ",
		"quando o usuário for beta", // prosa
		"segmentMatch beta-users",   // depende de serviço externo, fora de propósito
		"modulo 3",                  // idem
		">=",                        // operador sem nada para comparar
		`contains`,                  // idem
		"talvez",                    // palavra solta
	}
	for _, raw := range recusar {
		if c, err := Parse(raw); err == nil {
			t.Errorf("Parse(%q) aceitou como %q/%q — deveria recusar", raw, c.Op, c.Operand)
		}
	}
}

// NeedsOperand separa os dois operadores sem valor do resto.
func TestNeedsOperand(t *testing.T) {
	if OpAbsent.NeedsOperand() || OpPresent.NeedsOperand() {
		t.Error("ausente e presente não comparam contra nada")
	}
	for _, op := range []Op{OpEq, OpGte, OpContains, OpMatches, OpRollout} {
		if !op.NeedsOperand() {
			t.Errorf("%q compara contra um valor", op)
		}
	}
}

// A palavra do caso AUSENTE vem do catálogo de traduções, nunca cravada. A primeira
// versão tinha só o inglês, e um arquivo de flag escrito em português dizendo `ausente`
// saía como `ne` com operando `set` — resposta ERRADA, em silêncio, justamente no cenário
// que mais importa.
func TestParse_ausenteEmCadaIdioma(t *testing.T) {
	for _, raw := range []string{"absent", "ausente", "present", "presente"} {
		c, err := Parse(raw)
		if err != nil {
			t.Errorf("Parse(%q) recusou: %v", raw, err)
			continue
		}
		if c.Op != OpAbsent && c.Op != OpPresent {
			t.Errorf("Parse(%q).Op = %q — caiu na divisão de operando", raw, c.Op)
		}
		if c.Operand != "" {
			t.Errorf("Parse(%q).Operand = %q, deveria ser vazio", raw, c.Operand)
		}
	}
}
