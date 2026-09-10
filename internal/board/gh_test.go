package board

import (
	"os/exec"
	"strings"
	"testing"
)

// O `exit status 4` do `gh` é NÃO AUTENTICADO, e sozinho ele não diz nada.
//
// Medido com um dev novo: o `anchors next` respondeu `exit status 4` e ele não tinha como
// saber que faltava login. Um erro que não nomeia a causa manda a pessoa procurar no lugar
// errado — e o lugar errado costuma ser a própria máquina.
func TestDicaDeAuth_nomeiaACausaDoCodigo4(t *testing.T) {
	// O `exec.ExitError` de verdade, produzido por um comando que sai com 4.
	err := exec.Command("sh", "-c", "exit 4").Run()
	if err == nil {
		t.Fatal("o comando deveria falhar com 4")
	}

	dica := authHint(err, []string{"issue", "list"})
	for _, esperado := range []string{"NÃO AUTENTICADO", "gh auth login", "interativo"} {
		if !strings.Contains(dica, esperado) {
			t.Errorf("a dica não menciona %q:\n%s", esperado, dica)
		}
	}
}

// Outro código de saída NÃO ganha a dica: dizer "faltou login" para quem tem um repositório
// errado no `workflow.repo` manda a pessoa fazer um login que já está feito.
func TestDicaDeAuth_soOCodigo4(t *testing.T) {
	for _, codigo := range []string{"1", "2", "3", "5", "127"} {
		err := exec.Command("sh", "-c", "exit "+codigo).Run()
		if dica := authHint(err, nil); dica != "" {
			t.Errorf("exit %s ganhou a dica de auth:\n%s", codigo, dica)
		}
	}
	// E um erro que não é de saída (comando inexistente) também não.
	if dica := authHint(exec.Command("comando-que-nao-existe").Run(), nil); dica != "" {
		t.Errorf("erro que não é ExitError ganhou a dica:\n%s", dica)
	}
}

// O TETO existe para o `gh` não pendurar o Anchors.
//
// Medido: um dev novo rodou `anchors status` sem `gh auth login`, o comando ficou pendurado
// e ele precisou matá-lo. Do ponto de vista de quem usa, a ferramenta travou — e um comando
// que não volta é pior que um que falha, porque quem falha diz o que fazer.
func TestGhTimeout_existeEEhGeneroso(t *testing.T) {
	if ghTimeout == 0 {
		t.Fatal("sem teto, um `gh` que espere entrada pendura o Anchors para sempre")
	}
	// Generoso de propósito: um teto apertado transformaria rede lenta em erro, e o custo
	// de errar é assimétrico — falhar cedo faz o agente desistir de trabalho que ia
	// funcionar.
	if ghTimeout.Seconds() < 10 {
		t.Errorf("teto de %s é apertado — rede lenta viraria erro", ghTimeout)
	}
	if ghTimeout.Minutes() > 2 {
		t.Errorf("teto de %s é longo — o ponto é não custar a sessão", ghTimeout)
	}
}
