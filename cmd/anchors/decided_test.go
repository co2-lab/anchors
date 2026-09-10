package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// O `escalate --for-user` gravava um estado que ele não sabia reverter.
//
// Medido no blue-eyes (co2-lab/anchors#10): a decisão saiu, as revisões foram aplicadas
// nos 4 arquivos, a issue de decisão foi fechada — e o card continuou com
// `anchors:needs-user`, que é a label que faz o claim pular. Removi à mão.
//
// A instrução de remover existia, mas dentro do corpo do card que o WORKFLOW abre. Quem
// escala pelo CLI recebia só `· card #4 parado até a decisão`.
func TestDecided_exigeCardEResolucao(t *testing.T) {
	casos := []struct {
		nome  string
		args  []string
		trava string // o que a mensagem de erro tem de dizer
	}{
		{"sem card", []string{"--resolution", "ABCDE-R0001: x"}, "--card"},
		{"sem resolução", []string{"--card", "4"}, "--resolution"},
		{"resolução em branco", []string{"--card", "4", "--resolution", "   "}, "--resolution"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			cmd := newDecidedCmd()
			cmd.SetArgs(c.args)
			cmd.SetOut(os.Stderr)
			cmd.SetErr(os.Stderr)
			err := cmd.Execute()
			if err == nil {
				t.Fatal("aceitou sem o argumento obrigatório")
			}
			if !strings.Contains(err.Error(), c.trava) {
				t.Errorf("o erro não diz o que falta (%q):\n%s", c.trava, err)
			}
		})
	}
}

// A `--resolution` obrigatória não é burocracia: a saída de uma decisão em aberto é UMA —
// a resposta vira REGRA, com código. Liberar o card sem dizer qual revisão nasceu dela
// deixaria a decisão sem rastro, e é a mesma exigência que o `open-questions-resolved`
// faz em outro contexto ("apagar o item sem regra nova é varrer a pergunta para debaixo
// do tapete").
func TestDecided_aMensagemDizPorQueAResolucaoEhObrigatoria(t *testing.T) {
	cmd := newDecidedCmd()
	cmd.SetArgs([]string{"--card", "4"})
	cmd.SetOut(os.Stderr)
	cmd.SetErr(os.Stderr)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("aceitou sem --resolution")
	}
	// não basta dizer que falta: tem de dizer POR QUE, senão a próxima pessoa passa
	// qualquer string para satisfazer o comando
	if !strings.Contains(err.Error(), "vira REGRA") {
		t.Errorf("o erro não explica o motivo da exigência:\n%s", err)
	}
}

// No modo LOCAL a decisão mora em `issues/`, e resolvê-la é mover o arquivo — não há
// label a remover. Responder isso é melhor que falhar no `gh`.
func TestDecided_recusaForaDoModoGithub(t *testing.T) {
	root := t.TempDir()
	yaml := "version: 1\nlayers:\n  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n"
	if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := newDecidedCmd()
	cmd.SetArgs([]string{"--root", root, "--card", "4", "--resolution", "ABCDE-R0001: x"})
	cmd.SetOut(os.Stderr)
	cmd.SetErr(os.Stderr)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("aceitou fora do modo github")
	}
	if !strings.Contains(err.Error(), "issues/") {
		t.Errorf("o erro não aponta onde a decisão mora no modo local:\n%s", err)
	}
}
