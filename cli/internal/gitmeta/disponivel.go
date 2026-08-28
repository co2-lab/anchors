package gitmeta

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Disponibilidade é o diagnóstico de por que uma operação de git não pode acontecer.
// Existe porque "não deu para usar o git" tem causas com CONSERTOS diferentes, e um
// erro genérico (`exit status 128`) manda o usuário investigar a camada errada — o
// mesmo custo que o WORKFLOW.md §2 registra no caso do `login.yaml`.
type Disponibilidade int

const (
	// Disponível — git instalado e a raiz está sob um repositório.
	Disponível Disponibilidade = iota
	// SemBinário — `git` não está no PATH. Conserto: instalar.
	SemBinário
	// SemRepo — git existe, mas a raiz não está sob repositório. Conserto: `git init`.
	SemRepo
)

// Verifica classifica a raiz. Barata (LookPath + stat), pensada para ser chamada no
// ponto de uso, imediatamente antes de tentar a operação — é lá que se sabe QUAL ação
// vai ficar incompleta, e é isso que a mensagem precisa dizer.
func Verifica(root string) Disponibilidade {
	if _, err := exec.LookPath("git"); err != nil {
		return SemBinário
	}
	dir := root
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return Disponível
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			return SemRepo
		}
		dir = pai
	}
}

// Explica devolve a frase que nomeia a causa e o conserto, para ser embutida no erro
// do comando. `acao` é o que o comando ia fazer, na voz do comando ("listar os
// arquivos staged", "instalar o pre-commit") — assim a mensagem liga o sintoma à
// causa numa linha só, em vez de deixar o usuário adivinhar qual das duas faltas é.
//
// Devolve "" quando o git está disponível: aí a falha é outra, e inventar uma causa
// de git seria pior do que repassar o erro cru.
func Explica(d Disponibilidade, acao string) string {
	switch d {
	case SemBinário:
		return "não deu para " + acao + ": o git não está instalado (ou não está no PATH)"
	case SemRepo:
		return "não deu para " + acao + ": este projeto não está sob git — rode `git init`"
	default:
		return ""
	}
}
