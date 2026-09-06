package main

import (
	"slices"
	"testing"
)

// O HOOK NÃO PODE DESPEJAR O RELATÓRIO INTEIRO — muito menos duas vezes.
//
// Medido no blue-eyes: o `pre-commit` e o `commit-msg` chamam o MESMO
// `anchors verify --phase pre-commit --staged` (de propósito: um reporta cedo,
// o outro barra, e só o segundo tem a mensagem em mãos). Sem `--only-issues` a
// tabela de 42 gates saía nas duas, ~88 linhas por commit, e o `remote:` do
// push — a única linha que a pessoa esperava — ficava fora da tela.
//
// O `--only-issues` já existia e ninguém o ligava: é o quarto caso desta
// família no projeto (mecanismo construído, testado, e não conectado).
func TestFaseAutomaticaNaoDespejaORelatorioInteiro(t *testing.T) {
	for _, fase := range []string{"pre-commit", "pre-push", "ci"} {
		args := checkArgs(".", fase, "", "", []string{"a.spec.md"}, false, false, true)
		if !slices.Contains(args, "--only-issues") {
			t.Errorf("fase %q chama o check SEM --only-issues: o hook despeja a tabela de gates inteira\n  args: %v", fase, args)
		}
	}
}

// Na fase MANUAL o relatório é o produto: ali a pessoa digitou o comando para
// ver a foto, e omitir os gates limpos esconderia o que ela pediu.
func TestFaseManualMostraORelatorioCompleto(t *testing.T) {
	args := checkArgs(".", "manual", "", "", []string{"a.spec.md"}, false, false, true)
	if slices.Contains(args, "--only-issues") {
		t.Errorf("fase manual omitiu os gates limpos — a pessoa pediu a foto e recebeu o recorte\n  args: %v", args)
	}
	if slices.Contains(args, "--deterministic") {
		t.Errorf("fase manual virou determinística — ali o julgamento por IA é justamente o ponto")
	}
}

// O que já funcionava continua: a extração de `checkArgs` do RunE não pode ter
// perdido nenhuma flag pelo caminho.
func TestCheckArgsPreservaAsFlagsDaFachada(t *testing.T) {
	args := checkArgs("/r", "pre-commit", "/tmp/MSG", "types", []string{"a.md", "b.md"}, false, true, true)

	for _, esperado := range []string{
		"check", "--root", "/r",
		"--changed", "a.md", "--changed", "b.md",
		"--phase", "pre-commit",
		"--commit-msg", "/tmp/MSG",
		"--category", "types",
		"--skip-slow", "--no-record",
	} {
		if !slices.Contains(args, esperado) {
			t.Errorf("flag %q sumiu na extração\n  args: %v", esperado, args)
		}
	}

	// `--all` e `--changed` são exclusivos: pedir os dois faria o check varrer o
	// projeto inteiro num hook que devia ser incremental.
	if slices.Contains(args, "--all") {
		t.Errorf("--changed veio junto com --all\n  args: %v", args)
	}
}
