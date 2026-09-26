package issue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// A DEDUPLICAÇÃO no github vem do marcador no CORPO, e a confirmação tem de ser exata.
//
// A busca do GitHub é por texto e devolve aproximações: sem conferir o marcador, um
// achado sobre `Foo.spec.md` casaria o card de `FooBar.spec.md` — e o Anchors fecharia o
// card errado, que é pior que não fechar nenhum.
func TestMarcadorEhExatoENaoPrefixo(t *testing.T) {
	corpoDeOutro := fmt.Sprintf(KeyMarker, "trinca-completa:packages/Foo.spec.md:violation")
	marcaProcurada := fmt.Sprintf(KeyMarker, "trinca-completa:packages/Foo.spec.md")

	// O marcador do outro card CONTÉM o prefixo do procurado, e mesmo assim não pode
	// casar: são achados de alvos diferentes.
	if strings.Contains(corpoDeOutro, marcaProcurada) {
		t.Error("o marcador precisa fechar com `-->`, senão um achado casa o card de outro " +
			"cujo alvo apenas COMEÇA igual")
	}
}

// O título nomeia o card sem depender do formato da chave — prendê-lo à chave o tornaria
// ilegível para quem lê a lista de issues, que é onde ele aparece.
func TestTituloDizOQueEhSemAChave(t *testing.T) {
	g := GitHub{Repo: "acme/x", Label: "anchors"}
	for kind, esperado := range map[Kind]string{
		Violation: "Violação",
		Decision:  "Decisão",
		Stale:     "Desatualizado",
		Conflict:  "Conflito",
	} {
		got := g.title(Issue{Kind: kind, Gate: "triad-complete", Target: "a/b.spec.md"})
		if !strings.Contains(got, esperado) {
			t.Errorf("título de %s deveria dizer %q, veio %q", kind, esperado, got)
		}
		if !strings.Contains(got, "a/b.spec.md") || !strings.Contains(got, "triad-complete") {
			t.Errorf("o título deve nomear o gate e o alvo; veio %q", got)
		}
	}
}

// O DESTINO PADRÃO é arquivo. Um projeto local, ou qualquer chamador que não configurou
// nada, não pode acabar falando com a rede sem pedir.
func TestDestinoPadraoEhArquivo(t *testing.T) {
	UseFiles()
	if target != nil {
		t.Fatal("sem configurar, o destino tem de ser o arquivo")
	}
	UseGitHub("acme/x", "anchors")
	if target == nil || target.Repo != "acme/x" {
		t.Fatal("UsarGitHub deveria rotear para o repositório declarado")
	}
	UseFiles() // não vaza para os outros testes do pacote
}

// fakeGH puts on the PATH a `gh` that records each call (one line per call) and answers
// `issue list` with the given JSON. `failOn` is a shell glob: a call that matches it
// exits 1 with an error message. Nothing reaches the network.
func fakeGH(t *testing.T, listJSON, failOn string) (calls func() []string) {
	t.Helper()
	dir := t.TempDir()
	log := filepath.Join(dir, "calls.txt")
	script := "#!/bin/sh\n" +
		"printf '%s' \"$*\" | tr '\\n' ' ' >> '" + log + "'\nprintf '\\n' >> '" + log + "'\n" +
		"case \"$*\" in\n"
	if failOn != "" {
		script += failOn + ") echo 'HTTP 502' >&2; exit 1 ;;\n"
	}
	script += "'issue list'*) cat <<'__EOF__'\n" + listJSON + "\n__EOF__\n;;\nesac\nexit 0\n"
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return func() []string {
		b, _ := os.ReadFile(log)
		return strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	}
}

func cardJSON(number int, state, key string) string {
	body, _ := json.Marshal("some text\n\n" + fmt.Sprintf(KeyMarker, key) + "\n")
	return fmt.Sprintf(`{"number":%d,"state":%q,"body":%s}`, number, state, body)
}

var ghFinding = Issue{Kind: Violation, Gate: "triad-complete", Target: "a/Foo.spec.md", Detail: "no feature", Date: "2026-09-26"}

func TestGitHubOpen_createsCardWithLabels(t *testing.T) {
	key := ghFinding.Key()
	// The search returns a near miss (a longer key that CONTAINS ours): it must not match.
	calls := fakeGH(t, "["+cardJSON(7, "OPEN", key+"-other")+"]", "")
	g := GitHub{Repo: "acme/x", Label: "anchors"}

	i := ghFinding
	i.Dono = DonoUsuário
	created, at, err := g.Open(i, Todo)
	if err != nil || !created || at != Todo {
		t.Fatalf("Open = %v, %q, %v; want a new card in todo", created, at, err)
	}
	got := calls()
	if len(got) != 2 {
		t.Fatalf("want a search and a create, got %q", got)
	}
	if !strings.HasPrefix(got[0], "issue list --state all --limit 500 --search "+key) || !strings.HasSuffix(got[0], "--repo acme/x") {
		t.Errorf("search call = %q", got[0])
	}
	for _, frag := range []string{
		"issue create --title [triad-complete] Violação @ a/Foo.spec.md",
		fmt.Sprintf(KeyMarker, key),
		"--label anchors --label " + initx.LabelToDo + " --label " + initx.LabelNeedsUser + " --repo acme/x",
	} {
		if !strings.Contains(got[1], frag) {
			t.Errorf("create call lacks %q:\n%s", frag, got[1])
		}
	}
}

func TestGitHubOpen_futureDebtHasNoFlowLabel(t *testing.T) {
	calls := fakeGH(t, "[]", "")
	created, at, err := GitHub{Repo: "acme/x", Label: "anchors"}.Open(ghFinding, Future)
	if err != nil || !created || at != Future {
		t.Fatalf("Open = %v, %q, %v; want a new card in future", created, at, err)
	}
	create := calls()[1]
	if strings.Contains(create, initx.LabelToDo) || strings.Contains(create, initx.LabelNeedsUser) {
		t.Errorf("a future debt must carry only the workflow label:\n%s", create)
	}
}

func TestGitHubOpen_existingOpenCardIsLeftAlone(t *testing.T) {
	calls := fakeGH(t, "["+cardJSON(12, "OPEN", ghFinding.Key())+"]", "")
	created, at, err := GitHub{Repo: "acme/x", Label: "anchors"}.Open(ghFinding, Todo)
	if err != nil || created || at != Todo {
		t.Fatalf("Open = %v, %q, %v; want no new card", created, at, err)
	}
	if got := calls(); len(got) != 1 {
		t.Fatalf("only the search may run, got %q", got)
	}
}

func TestGitHubOpen_closedCardIsReopenedWithTheNewReport(t *testing.T) {
	calls := fakeGH(t, "["+cardJSON(12, "CLOSED", ghFinding.Key())+"]", "")
	created, at, err := GitHub{Repo: "acme/x", Label: "anchors"}.Open(ghFinding, Todo)
	if err != nil || !created || at != Todo {
		t.Fatalf("Open = %v, %q, %v; want the card reopened", created, at, err)
	}
	got := calls()
	if len(got) != 3 || got[1] != "issue reopen 12 --repo acme/x" ||
		!strings.HasPrefix(got[2], "issue comment 12 --body ⟲ **O achado voltou.**") || !strings.Contains(got[2], "no feature") {
		t.Fatalf("want reopen + comment with the detail, got %q", got)
	}
}

func TestGitHubOpen_propagatesGhFailures(t *testing.T) {
	for name, tc := range map[string]struct{ list, failOn string }{
		"search fails":  {"[]", "'issue list'*"},
		"bad search":    {"not json", ""},
		"create fails":  {"[]", "'issue create'*"},
		"reopen fails":  {"[" + cardJSON(3, "CLOSED", ghFinding.Key()) + "]", "'issue reopen'*"},
		"comment fails": {"[" + cardJSON(3, "CLOSED", ghFinding.Key()) + "]", "'issue comment'*"},
	} {
		t.Run(name, func(t *testing.T) {
			fakeGH(t, tc.list, tc.failOn)
			created, _, err := GitHub{Repo: "acme/x", Label: "anchors"}.Open(ghFinding, Todo)
			if err == nil || created {
				t.Fatalf("Open = %v, %v; want the failure reported", created, err)
			}
			if tc.failOn != "" && !strings.Contains(err.Error(), "HTTP 502") {
				t.Errorf("the error must carry gh's output, got %v", err)
			}
		})
	}
}

func TestGitHubResolve(t *testing.T) {
	g := GitHub{Repo: "acme/x", Label: "anchors"}
	key := ghFinding.Key()

	calls := fakeGH(t, "["+cardJSON(5, "OPEN", key)+"]", "")
	if ok, err := g.Resolve(key); !ok || err != nil {
		t.Fatalf("Resolve of an open card = %v, %v; want closed", ok, err)
	}
	if got := calls(); len(got) != 2 || !strings.HasPrefix(got[1], "issue close 5 --comment ✓ Resolvido") {
		t.Fatalf("want the card closed with a comment, got %q", got)
	}

	for name, list := range map[string]string{
		"already closed": "[" + cardJSON(5, "CLOSED", key) + "]",
		"not found":      "[]",
	} {
		calls = fakeGH(t, list, "")
		if ok, err := g.Resolve(key); ok || err != nil {
			t.Errorf("%s: Resolve = %v, %v; want nothing done", name, ok, err)
		}
		if got := calls(); len(got) != 1 {
			t.Errorf("%s: only the search may run, got %q", name, got)
		}
	}

	fakeGH(t, "["+cardJSON(5, "OPEN", key)+"]", "'issue close'*")
	if ok, err := g.Resolve(key); ok || err == nil {
		t.Fatalf("a failed close = %v, %v; want the error", ok, err)
	}
}
