package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRelToDevolveIdDeNo guarda a portabilidade do mapa. O que relTo devolve é o id de
// um nó, e o id é o mesmo em toda máquina — o mapa é versionado e trafega entre macOS,
// Linux e Windows. Sem a normalização, `filepath.Clean`/`Rel` devolvem
// "packages\backend\x.spec.md" no Windows, a busca no mapa (gravado com "/") não casa, e
// um arquivo que ESTÁ no mapa aparece como "REGIDO mas fora do mapa" — pedindo um
// `map build` que reescreveria o mapa inteiro no dialeto da máquina.
func TestRelToDevolveIdDeNo(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "packages", "backend", "services")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "auth.spec.md"), []byte("# spec\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	const id = "packages/backend/services/auth.spec.md"

	// As três origens legítimas do argumento: relativo à raiz com barra normal (a forma
	// que os prompts do Anchors imprimem), relativo com o separador nativo (o que o
	// shell da máquina completa), e absoluto.
	casos := map[string]string{
		"relativo com barra": id,
		"relativo nativo":    filepath.FromSlash(id),
		"absoluto":           filepath.Join(root, "packages", "backend", "services", "auth.spec.md"),
	}
	for nome, arg := range casos {
		got := relTo(root, arg)
		if got != id {
			t.Errorf("%s: relTo(%q) = %q; queria o id %q", nome, arg, got, id)
		}
		if strings.Contains(got, `\`) {
			t.Errorf("%s: id de nó não pode levar barra invertida: %q", nome, got)
		}
	}
}

// TestRelToNaoInventaCaminhoQuandoFalha — nos caminhos de erro a função devolve o
// argumento como veio; o que ela não pode é devolvê-lo em dialeto de máquina, pelo
// mesmo motivo do caso feliz.
func TestRelToNaoInventaCaminhoQuandoFalha(t *testing.T) {
	root := t.TempDir()
	got := relTo(root, filepath.FromSlash("nao/existe/em/lugar/nenhum.spec.md"))
	if strings.Contains(got, `\`) {
		t.Errorf("nem no caminho de erro o id pode levar barra invertida: %q", got)
	}
}
