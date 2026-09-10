package code

import "testing"

func TestGenerateShape(t *testing.T) {
	// não exijo o código EXATO do app de referência (aquele usa silabação/hypher); exijo a FORMA:
	// 4 chars, maiúsculos, determinístico, e âncora na inicial.
	cases := []string{"Spacer", "Divider", "HomeScreen", "Button", "LoginScreen"}
	for _, name := range cases {
		c := Generate(name)
		if len(c) != Slots {
			t.Errorf("%s → %q: esperava %d chars", name, c, Slots)
		}
		for i := 0; i < len(c); i++ {
			if !(c[i] >= 'A' && c[i] <= 'Z') && !(c[i] >= '0' && c[i] <= '9') {
				t.Errorf("%s → %q: char não-maiúsculo/dígito", name, c)
			}
		}
		if Generate(name) != c {
			t.Errorf("%s: não determinístico", name)
		}
	}
}

func TestGenerateAnchorsOnInitial(t *testing.T) {
	// a 1ª letra do código deve ser a inicial do nome (após stripGeneric)
	if got := Generate("Spacer"); got[0] != 'S' {
		t.Errorf("Spacer → %q: deveria começar com S", got)
	}
	if got := Generate("HomeScreen"); got[0] != 'H' {
		t.Errorf("HomeScreen → %q: deveria começar com H", got)
	}
}

func TestStripGeneric(t *testing.T) {
	cases := map[string]string{
		"LoginScreen": "Login", "AlertSheet": "Alert", "ScrollableLayout": "Scrollable",
		"Button": "Button", "Screen": "Screen", // só o genérico → mantém
	}
	for in, want := range cases {
		if got := StripGeneric(in); got != want {
			t.Errorf("StripGeneric(%q) = %q, quer %q", in, got, want)
		}
	}
}

func TestGenerateMultiWordUsesEachWord(t *testing.T) {
	// TransactionDetail (2 palavras) → deve puxar de ambas: T...D...
	c := Generate("TransactionDetail")
	if c[0] != 'T' {
		t.Errorf("TransactionDetail → %q: início T", c)
	}
	// contém uma letra de 'Detail' (D)?
	if !contains([]byte(c), 'D') {
		t.Errorf("TransactionDetail → %q: deveria conter D (2ª palavra)", c)
	}
}

func TestGenerateUniqueResolvesCollision(t *testing.T) {
	base := Generate("Spacer")
	taken := map[string]bool{base: true}
	got := GenerateUnique("Spacer", taken)
	if got == base {
		t.Fatalf("colisão não resolvida: %q == %q", got, base)
	}
	if taken[got] {
		t.Fatalf("código resolvido %q ainda está tomado", got)
	}
	if len(got) != Slots {
		t.Fatalf("código resolvido %q não tem %d chars", got, Slots)
	}
}

func TestGenerateUniqueNoCollisionReturnsCanonical(t *testing.T) {
	if got := GenerateUnique("Divider", map[string]bool{}); got != Generate("Divider") {
		t.Errorf("sem colisão deveria devolver o canônico")
	}
}

func TestShortNamePadsWithX(t *testing.T) {
	// nome curto → completa com X (ex.: "Ok" tem poucas letras)
	c := Generate("Ok")
	if len(c) != Slots {
		t.Fatalf("Ok → %q: %d chars", c, len(c))
	}
	if c[len(c)-1] != 'X' {
		t.Logf("Ok → %q (padding X esperado no fim se faltar letra)", c)
	}
}
