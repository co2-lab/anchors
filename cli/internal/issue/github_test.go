package issue

import (
	"fmt"
	"strings"
	"testing"
)

// A DEDUPLICAÇÃO no github vem do marcador no CORPO, e a confirmação tem de ser exata.
//
// A busca do GitHub é por texto e devolve aproximações: sem conferir o marcador, um
// achado sobre `Foo.spec.md` casaria o card de `FooBar.spec.md` — e o Anchors fecharia o
// card errado, que é pior que não fechar nenhum.
func TestMarcadorEhExatoENaoPrefixo(t *testing.T) {
	corpoDeOutro := fmt.Sprintf(MarcadorChave, "trinca-completa:packages/Foo.spec.md:violation")
	marcaProcurada := fmt.Sprintf(MarcadorChave, "trinca-completa:packages/Foo.spec.md")

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
		got := g.titulo(Issue{Kind: kind, Gate: "trinca-completa", Target: "a/b.spec.md"})
		if !strings.Contains(got, esperado) {
			t.Errorf("título de %s deveria dizer %q, veio %q", kind, esperado, got)
		}
		if !strings.Contains(got, "a/b.spec.md") || !strings.Contains(got, "trinca-completa") {
			t.Errorf("o título deve nomear o gate e o alvo; veio %q", got)
		}
	}
}

// O DESTINO PADRÃO é arquivo. Um projeto local, ou qualquer chamador que não configurou
// nada, não pode acabar falando com a rede sem pedir.
func TestDestinoPadraoEhArquivo(t *testing.T) {
	UsarArquivos()
	if destino != nil {
		t.Fatal("sem configurar, o destino tem de ser o arquivo")
	}
	UsarGitHub("acme/x", "anchors")
	if destino == nil || destino.Repo != "acme/x" {
		t.Fatal("UsarGitHub deveria rotear para o repositório declarado")
	}
	UsarArquivos() // não vaza para os outros testes do pacote
}
