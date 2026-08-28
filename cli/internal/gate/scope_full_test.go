package gate

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// gravaContagem monta um gate que registra QUANTOS alvos recebeu em cada execução, uma
// linha por invocação. É assim que se distingue "rodou uma vez sobre o projeto" de
// "rodou N vezes com o projeto picado".
func gravaContagem(t *testing.T, root, saida string) string {
	t.Helper()
	return "printf '%s\\n' \"$#\" >> " + saida
}

func linhas(t *testing.T, p string) []string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	var out []string
	for _, l := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

func nós(n int) []mapx.Node {
	out := make([]mapx.Node, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, mapx.Node{ID: "src/arquivo" + string(rune('a'+i%26)) + string(rune('a'+i/26)) + ".ts", Kind: mapx.KindCode})
	}
	return out
}

// TestScopeFullRodaUmaVezSemAlvos é o ponto do scope_full: no `--all` o recorte é o
// projeto inteiro, e passar a lista deixa de fazer sentido — a ferramenta que sabe
// varrer sozinha deve ser chamada UMA vez, sem alvo. Antes disso, o mesmo gate rodava
// em lotes: N invocações onde uma bastava (e, no Windows, N estouros de linha de comando).
func TestScopeFullRodaUmaVezSemAlvos(t *testing.T) {
	root := t.TempDir()
	saida := filepath.ToSlash(filepath.Join(root, "invocacoes.txt"))

	g := config.Gate{
		Name:      "lint",
		On:        []string{string(mapx.KindCode)},
		Scope:     config.ScopeBatch,
		ScopeFull: config.ScopeProject,
		Run:       gravaContagem(t, root, saida),
	}

	RunCompleto([]config.Gate{g}, nós(300), root, nil, &config.Config{}, true)

	got := linhas(t, saida)
	if len(got) != 1 {
		t.Fatalf("no full o gate devia rodar UMA vez; rodou %d", len(got))
	}
	if got[0] != "0" {
		t.Errorf("no full o gate não recebe alvos; recebeu %s", got[0])
	}
}

// TestSemScopeFullContinuaEmLotes — o opt-in não pode mudar quem não pediu. Um gate
// batch sem `scope_full` continua recebendo os alvos, porque o script dele pode muito
// bem sair 0 quando não recebe nada: promovê-lo a project sozinho o deixaria verde sem
// olhar arquivo nenhum.
func TestSemScopeFullContinuaEmLotes(t *testing.T) {
	root := t.TempDir()
	saida := filepath.ToSlash(filepath.Join(root, "invocacoes.txt"))

	g := config.Gate{
		Name:  "lint",
		On:    []string{string(mapx.KindCode)},
		Scope: config.ScopeBatch,
		Run:   gravaContagem(t, root, saida),
	}

	RunCompleto([]config.Gate{g}, nós(300), root, nil, &config.Config{}, true)

	got := linhas(t, saida)
	if len(got) == 0 {
		t.Fatal("o gate batch tem de rodar")
	}
	total := 0
	for _, l := range got {
		if l == "0" {
			t.Fatal("gate batch sem scope_full não pode ser chamado sem alvos")
		}
		n, err := strconv.Atoi(l)
		if err != nil {
			t.Fatalf("linha inesperada no registro de invocações: %q", l)
		}
		total += n
	}
	if total != 300 {
		t.Errorf("os lotes cobrem %d alvos; deviam cobrir os 300", total)
	}
}

// TestScopeFullNaoValeNoIncremental fixa a outra metade da regra: fora do `--all` o
// recorte é pequeno e específico, e é exatamente ele que o gate precisa receber. Aplicar
// scope_full aqui faria um commit de um arquivo varrer o projeto inteiro.
func TestScopeFullNaoValeNoIncremental(t *testing.T) {
	root := t.TempDir()
	saida := filepath.ToSlash(filepath.Join(root, "invocacoes.txt"))

	g := config.Gate{
		Name:      "lint",
		On:        []string{string(mapx.KindCode)},
		Scope:     config.ScopeBatch,
		ScopeFull: config.ScopeProject,
		Run:       gravaContagem(t, root, saida),
	}

	RunCompleto([]config.Gate{g}, nós(3), root, nil, &config.Config{}, false)

	got := linhas(t, saida)
	if len(got) != 1 || got[0] != "3" {
		t.Errorf("no incremental o gate recebe os 3 alvos do recorte; veio %v", got)
	}
}
