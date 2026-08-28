package main

import (
	"errors"
	"os/exec"
	"testing"
)

// O código de saída do FILHO precisa atravessar a fronteira do processo.
//
// O `verify` é uma fachada que reexecuta o próprio binário. O `c.Run()` devolve
// um `*exec.ExitError` genérico, e o `main` — que converte `errNaoRegido` em
// `ExitNaoRegido` — não o reconhecia. O `check` saía com 3 ("não tenho
// jurisdição"), o `verify` traduzia para 1, e o pre-commit barrava um commit só
// de configuração: exatamente o caso que o código 3 existe para permitir.
func TestExitNaoRegidoAtravessaOSubprocesso(t *testing.T) {
	// `false` sai com 1; `sh -c 'exit 3'` reproduz o código do filho.
	err := traduzSaidaDoFilho(exec.Command("sh", "-c", "exit 3").Run())

	var nr errNaoRegido
	if !errors.As(err, &nr) {
		t.Fatalf("exit 3 do filho não virou errNaoRegido: %v (%T)", err, err)
	}
}

// Qualquer outro código continua sendo falha: só o 3 tem tratamento próprio.
func TestOutrosCodigosContinuamSendoFalha(t *testing.T) {
	err := traduzSaidaDoFilho(exec.Command("sh", "-c", "exit 1").Run())

	var nr errNaoRegido
	if errors.As(err, &nr) {
		t.Errorf("exit 1 foi tratado como não-regido — reprovação estaria sendo engolida")
	}
	if err == nil {
		t.Error("exit 1 não devolveu erro")
	}
}

func TestSucessoNaoDevolveErro(t *testing.T) {
	if err := traduzSaidaDoFilho(exec.Command("sh", "-c", "exit 0").Run()); err != nil {
		t.Errorf("exit 0 devolveu erro: %v", err)
	}
}
