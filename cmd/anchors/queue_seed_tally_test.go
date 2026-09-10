package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/queue"
)

// O NÚMERO NÃO PODE MORAR NO `reason`.
//
// Ele é derivável (as sementes se recontam do disco), e guardá-lo fazia o texto envelhecer
// na fila: a task nascia dizendo "6 de 7", duas specs eram entregues, e ela continuava
// dizendo 6 (co2-lab/anchors#11).
//
// Medido no blue-eyes: a razão velha me fez desconfiar de uma correção que eu tinha
// acabado de publicar — quatro comandos investigando um defeito que não existia, porque a
// contagem estava certa e o texto era velho.
func TestSeedTally(t *testing.T) {
	monta := func(t *testing.T, existentes ...string) string {
		t.Helper()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
			t.Fatal(err)
		}
		yaml := "version: 1\nlayers:\n" +
			"  plan:\n    pattern: \"plans/*.md\"\n    kind: plan\n" +
			"  spec:\n    pattern: \"**/*.spec.md\"\n    kind: spec\n"
		if err := os.WriteFile(filepath.Join(root, "anchors.yaml"), []byte(yaml), 0o644); err != nil {
			t.Fatal(err)
		}
		plano := "- [ ] `src/A.spec.md` — a\n- [ ] `src/B.spec.md` — b\n- [ ] `src/C.spec.md` — c\n"
		if err := os.WriteFile(filepath.Join(root, "plans", "0001-x.md"), []byte(plano), 0o644); err != nil {
			t.Fatal(err)
		}
		for _, f := range existentes {
			p := filepath.Join(root, f)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte("# spec\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return root
	}
	task := queue.Task{
		Changed: "plans/0001-x.md", Kind: "plan", Origin: "seed",
		Reason: "um plano foi semeado — gere as specs que ele lista",
	}

	t.Run("conta o que falta AGORA, não na criação", func(t *testing.T) {
		root := monta(t)
		if got := seedTally(root, task); !strings.Contains(got, "3 de 3") {
			t.Errorf("com zero specs no disco, queria \"3 de 3\"; veio %q", got)
		}
		// a MESMA task, depois de duas entregas
		root = monta(t, "src/A.spec.md", "src/B.spec.md")
		if got := seedTally(root, task); !strings.Contains(got, "1 de 3") {
			t.Errorf("com duas specs entregues, queria \"1 de 3\"; veio %q", got)
		}
	})

	// A task ANTIGA já traz a contagem gravada — somar a recalculada imprimiria o número
	// duas vezes. Detectar pela frase, e não por versão: o campo não guarda quem o
	// escreveu, e uma task sobrevive a qualquer número de atualizações do binário.
	t.Run("task antiga com contagem gravada não duplica", func(t *testing.T) {
		root := monta(t)
		velha := task
		velha.Reason += " — 6 de 7 spec(s) deste plano ainda não existem"
		if got := seedTally(root, velha); got != "" {
			t.Errorf("duplicaria a contagem: %q", got)
		}
	})

	// O `reason` de julgamento não envelhece — a pergunta é do gate, não do estado.
	t.Run("só task de semeadura recebe contagem", func(t *testing.T) {
		root := monta(t)
		for _, outra := range []queue.Task{
			{Changed: "x.spec.md", Kind: "judgment", Origin: "check", Reason: "gate 'y' — pergunta: ..."},
			{Changed: "plans/0001-x.md", Kind: "plan", Origin: "humano", Reason: "alguém pediu"},
		} {
			if got := seedTally(root, outra); got != "" {
				t.Errorf("kind=%s origin=%s recebeu contagem: %q", outra.Kind, outra.Origin, got)
			}
		}
	})
}
