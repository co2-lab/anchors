package code

import "testing"

func TestGenerateFromPath_normalNoCollision(t *testing.T) {
	// basename normal, nada tomado → algoritmo canônico, inalterado.
	got := GenerateFromPath("packages/backend/models/metadata.spec.md", map[string]bool{})
	want := Generate("metadata")
	if got != want {
		t.Fatalf("sem colisão deveria ser o canônico %q, veio %q", want, got)
	}
}

func TestGenerateFromPath_collisionStaysUnique(t *testing.T) {
	// models/metadata pega o canônico; repositories/metadata COLIDE → recebe uma
	// variação ÚNICA (desempate cego, determinístico dado o mapa). O desempate
	// SIMÉTRICO (ambos prefixados) é uma feature de recode à parte — não aqui.
	taken := map[string]bool{}
	first := GenerateFromPath("packages/backend/models/metadata.spec.md", taken)
	taken[first] = true

	second := GenerateFromPath("packages/backend/repositories/metadata.spec.md", taken)
	if second == first {
		t.Fatalf("colisão não resolvida: ambos %q", second)
	}
	if taken[second] {
		t.Errorf("código de desempate %q não é único", second)
	}
	// determinismo dado o mapa: mesma entrada + mesmo taken → mesmo resultado.
	again := GenerateFromPath("packages/backend/repositories/metadata.spec.md", taken)
	if again != second {
		t.Errorf("não-determinístico dado o mapa: %q depois %q", second, again)
	}
}

func TestGenerateFromPath_genericBasenameUsesParent(t *testing.T) {
	// functions/manage-metadata/handler.ts → a unidade é "manage-metadata".
	got := GenerateFromPath("packages/backend/functions/manage-metadata/handler.ts", map[string]bool{})
	want := GenerateUnique("manage-metadata", map[string]bool{})
	if got != want {
		t.Fatalf("basename genérico deveria usar o dir-pai (%q), veio %q", want, got)
	}
	// e NÃO deveria ser o código de "handler"
	if got == Generate("handler") {
		t.Errorf("código genérico de 'handler' vazou: %q", got)
	}
}

func TestIsGenericBasename(t *testing.T) {
	for _, g := range []string{"handler", "index", "resource", "Handler", "INDEX"} {
		if !IsGenericBasename(g) {
			t.Errorf("%q deveria ser genérico", g)
		}
	}
	for _, d := range []string{"metadata", "Login", "useTransactions"} {
		if IsGenericBasename(d) {
			t.Errorf("%q NÃO deveria ser genérico", d)
		}
	}
}
