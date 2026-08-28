package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// O pior silêncio que o confronto tinha: sem git ele devolvia `nil`, e `nil` é o mesmo
// valor de "conferi e está tudo certo". Quem lia a saída sem aviso nenhum concluía que
// os arquivos declarados conferiam — uma afirmação que ninguém tinha verificado.
func TestConfrontoDizQuandoNaoTeveComoOlhar(t *testing.T) {
	dir := t.TempDir()
	if _, _, ok := repoAcimaDe(dir); ok {
		t.Skipf("o diretório temporário %s está dentro de um repo git", dir)
	}
	declarados := []string{filepath.Join(dir, "a.go")}
	if err := os.WriteFile(declarados[0], []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	avisos, confrontou := arquivosNaoTocados(dir, declarados)

	if confrontou {
		t.Fatal("sem repositório não há como confrontar — dizer que confrontou é mentir")
	}
	if len(avisos) != 0 {
		t.Errorf("sem confronto não há achado a reportar: %v", avisos)
	}
}

// A contrapartida: num repo de verdade o confronto ACONTECE, e um arquivo declarado
// que ninguém tocou continua sendo acusado.
func TestConfrontoAcusaArquivoDeclaradoENaoTocado(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git não instalado")
	}
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "t@t"},
		{"config", "user.name", "t"},
	} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("preparo %v: %s", args, out)
		}
	}
	// Committed e intocado desde então: o confronto tem de reparar nele.
	quieto := filepath.Join(dir, "quieto.go")
	if err := os.WriteFile(quieto, []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-m", "base"}} {
		c := exec.Command("git", args...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("preparo %v: %s", args, out)
		}
	}

	avisos, confrontou := arquivosNaoTocados(dir, []string{quieto})

	if !confrontou {
		t.Fatal("com repositório o confronto tem de acontecer")
	}
	if len(avisos) != 1 {
		t.Fatalf("arquivo declarado e não tocado deveria ser acusado, veio %v", avisos)
	}
}

// repoAcimaDe é o helper local do teste — evita depender do formato de `git status`
// para saber se o TempDir caiu dentro de um repo.
func repoAcimaDe(root string) (string, string, bool) {
	dir := root
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, "", true
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			return "", "", false
		}
		dir = pai
	}
}
