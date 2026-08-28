package recode

import "testing"

func TestTestIDPrefix(t *testing.T) {
	if got := TestIDPrefix("TCDTX", "lower"); got != "tcdtx" {
		t.Errorf("lower: %q", got)
	}
	if got := TestIDPrefix("TCDTX", ""); got != "" {
		t.Errorf("vazio deveria dar vazio, veio %q", got)
	}
}

func TestRewriteTestIDs(t *testing.T) {
	in := `testID="tcdt-amount" e o outro <View testID='tcdt-row-1' /> mas 'tcdtsomething' não`
	got, n := RewriteTestIDs(in, "tcdt", "tctx")
	if n != 2 {
		t.Fatalf("esperava 2 trocas, veio %d: %q", n, got)
	}
	if !contains(got, `"tctx-amount"`) || !contains(got, `'tctx-row-1'`) {
		t.Errorf("não trocou os testIDs: %q", got)
	}
	if !contains(got, "tcdtsomething") {
		t.Errorf("tocou um token que NÃO é testID (sem hífen): %q", got)
	}
}

func TestRewriteTestIDs_noopWhenEmpty(t *testing.T) {
	in := `testID="tcdt-x"`
	if _, n := RewriteTestIDs(in, "", "tctx"); n != 0 {
		t.Errorf("prefixo vazio deveria ser no-op")
	}
}

func TestCountTestIDPrefix(t *testing.T) {
	in := `testID="txdt-a" 'txdt-b' e "outro-c"`
	if got := CountTestIDPrefix(in, "txdt"); got != 2 {
		t.Errorf("esperava 2, veio %d", got)
	}
	if got := CountTestIDPrefix(in, "tcdt"); got != 0 {
		t.Errorf("prefixo ausente deveria dar 0, veio %d", got)
	}
}

func TestFileMatchesCode(t *testing.T) {
	pats := []string{"**/{{code}}-*.yaml", "**/{{code}}-suite.yaml", "**/*.{{code}}-*.png"}
	ok := []string{
		"apps/.maestro/screens/x/TCDTX-B01.yaml",
		"apps/.maestro/suites/TCDTX-suite.yaml",
		"apps/screens/Foo.TCDTX-VR-loaded.png",
	}
	for _, p := range ok {
		if !FileMatchesCode(p, "TCDTX", pats) {
			t.Errorf("%q deveria casar", p)
		}
	}
	no := []string{
		"apps/.maestro/screens/x/MNDTX-B01.yaml", // outro código
		"apps/components/ActionLink.tsx",         // sem código no nome
	}
	for _, p := range no {
		if FileMatchesCode(p, "TCDTX", pats) {
			t.Errorf("%q NÃO deveria casar", p)
		}
	}
}

func TestRenameFilePath(t *testing.T) {
	cases := map[string]string{
		"a/b/TCDTX-B01.yaml":        "a/b/TCTXX-B01.yaml",
		"a/b/TCDTX-suite.yaml":      "a/b/TCTXX-suite.yaml",
		"a/Foo.TCDTX-VR-loaded.png": "a/Foo.TCTXX-VR-loaded.png",
		"a/b/ActionLink.tsx":        "a/b/ActionLink.tsx", // sem código → inalterado
		// Código MAIOR não é tocado: `TCDTXX` contém `TCDTX` como prefixo, e um replace
		// ingênuo o mutilaria em `TCTXXX`. A conversão dos exemplos de 4→5 chars colidiu
		// aqui — o caso original era `TCDT-B01` vs `TCDTX-B01`, e ao virar 5 os dois ficaram
		// `TCDTX-B01`, chave duplicada no mapa. O que o caso prova é a FRONTEIRA do código,
		// então o vizinho tem de ser mais longo que o alvo.
		"a/b/TCDTXX-B01.yaml": "a/b/TCDTXX-B01.yaml",
	}
	for in, want := range cases {
		if got := RenameFilePath(in, "TCDTX", "TCTXX"); got != want {
			t.Errorf("RenameFilePath(%q) = %q, quer %q", in, got, want)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
