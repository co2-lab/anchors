package flow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// ghDeGravacao põe no PATH um `gh` que registra o que foi chamado, em ordem.
func ghDeGravacao(t *testing.T) (dir string, chamadas func() []string) {
	t.Helper()
	dir = t.TempDir()
	log := filepath.Join(dir, "chamadas.txt")
	// UMA LINHA POR CHAMADA: o `--body` tem quebras de linha, e `echo "$@"` as
	// preservaria — cada linha do corpo viraria uma chamada no registro.
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" | tr '\\n' ' ' >> " + log + "\n" +
		"printf '\\n' >> " + log + "\n"
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return dir, func() []string {
		b, err := os.ReadFile(log)
		if err != nil {
			return nil
		}
		var out []string
		for _, l := range strings.Split(string(b), "\n") {
			if strings.TrimSpace(l) != "" {
				out = append(out, l)
			}
		}
		return out
	}
}

// A ORDEM importa: label, razão, fechar.
//
// Registrar a razão ANTES de aplicar a label é a pior das inconsistências possíveis — se a
// label falhar, o card fica com um comentário dizendo que foi descartado E continua no
// board, afirmando duas coisas contrárias ao mesmo tempo.
//
// Fechar por último, e só depois de tudo dar certo: um card descartado que continua aberto
// seria servido pelo `claim`.
func TestDescartaAplicaLabelAntesDeRegistrarARazao(t *testing.T) {
	_, chamadas := ghDeGravacao(t)
	if err := descarta("acme/projeto", "42", "o arquivo não existe mais"); err != nil {
		t.Fatal(err)
	}
	c := chamadas()
	if len(c) < 3 {
		t.Fatalf("esperava label, comentário e fechamento; veio: %v", c)
	}
	if !strings.Contains(c[0], "edit") || !strings.Contains(c[0], initx.LabelDiscarded) {
		t.Errorf("a PRIMEIRA chamada devia aplicar a label: %q", c[0])
	}
	if !strings.Contains(c[1], "comment") {
		t.Errorf("a razão vem DEPOIS da label: %q", c[1])
	}
	if !strings.Contains(c[2], "close") {
		t.Errorf("fechar vem por último: %q", c[2])
	}
}

// A RAZÃO vai no card, e não só na tela de quem rodou o comando.
func TestDescartaRegistraARazaoNoCard(t *testing.T) {
	_, chamadas := ghDeGravacao(t)
	if err := descarta("acme/projeto", "42", "spec apagada na limpeza do plano 0004"); err != nil {
		t.Fatal(err)
	}
	var comentario string
	for _, l := range chamadas() {
		if strings.Contains(l, "comment") {
			comentario = l
		}
	}
	if !strings.Contains(comentario, "spec apagada na limpeza") {
		t.Errorf("a razão não chegou ao card: %q", comentario)
	}
	// E o caminho de VOLTA: um descarte sem como desfazer é uma remoção com outro nome.
	if !strings.Contains(comentario, initx.LabelDiscarded) {
		t.Error("o comentário não diz qual label remover para trazer o card de volta")
	}
}
