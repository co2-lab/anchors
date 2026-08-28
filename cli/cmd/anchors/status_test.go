package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// capturaStatus roda o status numa raiz e devolve o que ele imprimiu.
func capturaStatus(t *testing.T, root string) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	errRun := runStatus(root)
	w.Close()
	os.Stdout = orig
	if errRun != nil {
		t.Fatalf("runStatus: %v", errRun)
	}
	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, rerr := r.Read(buf)
		sb.Write(buf[:n])
		if rerr != nil {
			break
		}
	}
	return sb.String()
}

func repoVazio(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".git", "refs", "heads"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// A razão de o comando existir: quem retoma o trabalho dias depois precisa saber onde
// parou. Num projeto que ainda não começou, a resposta é a fase DESCOBRIR — e ela tem de
// ser NOMEADA, senão o agente que abre a conversa começa adivinhando.
func TestStatusApontaAFaseDescobrirNumProjetoNovo(t *testing.T) {
	saida := capturaStatus(t, repoVazio(t))

	if !strings.Contains(saida, "DESCOBRIR") {
		t.Errorf("projeto sem PROJECT.md e sem config: a fase que falta tem de ser nomeada:\n%s", saida)
	}
	if !strings.Contains(saida, "PRÓXIMO PASSO") {
		t.Errorf("o status existe para dizer o que fazer a seguir:\n%s", saida)
	}
}

// O status para no PRIMEIRO passo que falta. Listar tudo de uma vez faria o leitor
// escolher por onde começar — e a ordem do ciclo é justamente o que ele não deveria ter
// de reconstruir sozinho.
func TestStatusParaNoPrimeiroPassoQueFalta(t *testing.T) {
	dir := repoVazio(t)
	if err := os.WriteFile(filepath.Join(dir, "PROJECT.md"), []byte("# Projeto\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	saida := capturaStatus(t, dir)

	if !strings.Contains(saida, "PROJECT.md existe") {
		t.Errorf("o que já foi feito tem de aparecer como feito:\n%s", saida)
	}
	if !strings.Contains(saida, "anchors init") {
		t.Errorf("com PROJECT.md e sem config, o passo é o init:\n%s", saida)
	}
	// E não pode falar do que vem depois: o mapa ainda não é a pergunta.
	if strings.Contains(saida, "map build") {
		t.Errorf("adiantou um passo que ainda não é o próximo:\n%s", saida)
	}
}

// Sem repositório, o status para antes de tudo: metade do framework depende de git, e
// seguir descrevendo o ciclo daria a impressão de que está tudo bem.
func TestStatusParaSemRepositorio(t *testing.T) {
	dir := t.TempDir()
	if temRepoAcima(dir) {
		t.Skipf("o diretório temporário %s está dentro de um repo git", dir)
	}

	saida := capturaStatus(t, dir)

	if !strings.Contains(saida, "git init") {
		t.Errorf("sem repositório, o passo é criar um:\n%s", saida)
	}
	if strings.Contains(saida, "DESCOBRIR") {
		t.Errorf("não deve adiantar a fase seguinte antes de haver substrato:\n%s", saida)
	}
}

// A fila mora onde o modo declara (WORKFLOW.md §2). No modo local ela é `.anchors/tasks/`
// e `issues/` — mostrar board ali descreveria um lugar onde nada acontece.
func TestStatusLocalMostraAFilaLocal(t *testing.T) {
	dir := repoVazio(t)
	if err := os.WriteFile(filepath.Join(dir, "PROJECT.md"), []byte("# P\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte("version: 1\nlayers: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "anchors.graph.yaml"), []byte("version: 1\nnodes: []\nedges: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	saida := capturaStatus(t, dir)

	if !strings.Contains(saida, "fila: local") {
		t.Errorf("modo local tem de mostrar a fila local:\n%s", saida)
	}
	// Projeto montado e vazio: o passo é o primeiro plano — é ele que semeia as specs.
	if !strings.Contains(saida, "plano") {
		t.Errorf("projeto montado e vazio: o próximo passo é o primeiro plano:\n%s", saida)
	}
}

// temRepoAcima evita falso-positivo quando o TempDir cai dentro de um repo.
func temRepoAcima(root string) bool {
	dir := root
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return true
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			return false
		}
		dir = pai
	}
}
