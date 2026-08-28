package checklog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEspelhoDuplicaSaidaNoArquivo(t *testing.T) {
	dir := t.TempDir()
	e := Abrir(dir, true, "# cabeçalho\n\n")
	if e == nil {
		t.Fatal("Abrir devolveu nil")
	}
	fmt.Println("linha do relatório")
	e.Fechar()

	b, err := os.ReadFile(filepath.Join(dir, Dir, "check-all.txt"))
	if err != nil {
		t.Fatalf("ler espelho: %v", err)
	}
	got := string(b)
	if !strings.Contains(got, "# cabeçalho") {
		t.Errorf("cabeçalho ausente:\n%s", got)
	}
	if !strings.Contains(got, "linha do relatório") {
		t.Errorf("saída não espelhada:\n%s", got)
	}
}

// O `--all` custa minutos e o `--changed` roda a cada commit. Se os dois
// escrevessem no mesmo arquivo, o pre-commit apagaria a foto completa — que é
// justamente a cara de reproduzir.
func TestEscoposNaoSeSobrescrevem(t *testing.T) {
	dir := t.TempDir()

	e := Abrir(dir, true, "# all\n")
	fmt.Println("foto completa")
	e.Fechar()

	e = Abrir(dir, false, "# changed\n")
	fmt.Println("incremental")
	e.Fechar()

	all, err := os.ReadFile(filepath.Join(dir, Dir, "check-all.txt"))
	if err != nil {
		t.Fatalf("ler check-all: %v", err)
	}
	if !strings.Contains(string(all), "foto completa") {
		t.Errorf("o --changed sobrescreveu a foto do --all:\n%s", all)
	}
	chg, err := os.ReadFile(filepath.Join(dir, Dir, "check-changed.txt"))
	if err != nil {
		t.Fatalf("ler check-changed: %v", err)
	}
	if strings.Contains(string(chg), "foto completa") {
		t.Errorf("escopos misturados no mesmo arquivo:\n%s", chg)
	}
}

// Saída maior que o buffer do pipe (64KB no Linux/macOS): sem alguém drenando
// do outro lado, o comando travaria ao escrever o próprio relatório.
func TestSaidaLongaNaoTrava(t *testing.T) {
	dir := t.TempDir()
	e := Abrir(dir, true, "")
	linha := strings.Repeat("x", 200)
	for i := 0; i < 2000; i++ { // ~400KB
		fmt.Println(linha)
	}
	fim := make(chan struct{})
	go func() { e.Fechar(); close(fim) }()
	select {
	case <-fim:
	case <-time.After(10 * time.Second):
		t.Fatal("Fechar travou — o pipe não estava sendo drenado")
	}

	b, err := os.ReadFile(filepath.Join(dir, Dir, "check-all.txt"))
	if err != nil {
		t.Fatalf("ler espelho: %v", err)
	}
	if n := strings.Count(string(b), linha); n != 2000 {
		t.Errorf("espelho truncado: %d linhas de 2000", n)
	}
}

// O espelho é conveniência. Se o diretório não puder ser criado, o check tem de
// seguir escrevendo na tela — trocar a varredura pelo relatório seria trocar o
// essencial pelo acessório.
func TestFalhaAoAbrirNaoDerruba(t *testing.T) {
	dir := t.TempDir()
	// Um ARQUIVO onde o `.anchors/` deveria ser: o MkdirAll falha.
	if err := os.WriteFile(filepath.Join(dir, Dir), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	e := Abrir(dir, true, "# nada\n")
	if e != nil {
		t.Error("Abrir devolveu espelho onde não podia criar o diretório")
	}
	e.Fechar() // seguro com nil
	if c := e.Caminho(); c != "" {
		t.Errorf("Caminho() = %q, queria vazio", c)
	}
	if os.Stdout == nil {
		t.Error("stdout ficou inválido após falha")
	}
}

func TestCabecalhoRegistraContexto(t *testing.T) {
	quando := time.Date(2026, 8, 18, 14, 32, 0, 0, time.UTC)
	h := Cabecalho("anchors check --all", "abc1234", "fix: algo", 3, quando)

	for _, quer := range []string{
		"anchors check --all",
		"2026-08-18 14:32:00",
		"abc1234 fix: algo",
		"3 arquivos modificados",
	} {
		if !strings.Contains(h, quer) {
			t.Errorf("cabeçalho sem %q:\n%s", quer, h)
		}
	}

	if limpa := Cabecalho("c", "abc", "s", 0, quando); !strings.Contains(limpa, "árvore: limpa") {
		t.Errorf("árvore limpa não registrada:\n%s", limpa)
	}
	// Sem git (repo novo, ou git ausente) o cabeçalho não pode inventar um HEAD.
	if semGit := Cabecalho("c", "", "", 0, quando); strings.Contains(semGit, "HEAD:") {
		t.Errorf("HEAD inventado sem git:\n%s", semGit)
	}
}
