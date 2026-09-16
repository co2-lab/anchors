package main

import (
	"strings"
	"testing"
)

// O HOOK PRECISA FUNCIONAR EM WORKTREE, e `$ROOT/.git` não serve.
//
// Num worktree, `.git` é um ARQUIVO que aponta para `…/.git/worktrees/<nome>` — escrever
// dentro dele falha com `Not a directory`.
//
// MEDIDO: um agente trabalhando em worktree levou
//
//	.git/hooks/pre-commit: line 48: …/be-rev1/.git/anchors-freeze-cache: Not a directory
//
// e commitou com `--no-verify`. O hook INTEIRO deixou de rodar por causa do cache — e um
// hook que falha assim é pior que hook nenhum: ele treina quem o usa a contorná-lo.
//
// `git rev-parse --git-dir` devolve o caminho certo nos dois casos.
func TestHookAchaOGitDirEmWorktree(t *testing.T) {
	if strings.Contains(preCommitScript, `"$ROOT/.git/`) {
		t.Error("o hook monta um caminho dentro de `$ROOT/.git` — num worktree isso é um " +
			"ARQUIVO, e a escrita falha com `Not a directory`, derrubando o hook inteiro")
	}
	if !strings.Contains(preCommitScript, "rev-parse --git-dir") {
		t.Error("o hook não usa `git rev-parse --git-dir` — é o único jeito de achar o " +
			"diretório do git que funciona em clone comum E em worktree")
	}
}
