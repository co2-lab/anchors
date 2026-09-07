package board

import (
	"strings"
	"testing"
)

// O código é procurado no TÍTULO, entre colchetes — é o que o `identify.yml` escreve, e é
// o que dá identidade ao card.
//
// Buscar no CORPO casaria qualquer card que MENCIONE a unidade — o de outra que depende
// dela, por exemplo — e comentar no card errado é pior que não comentar: o registro de
// entrega apareceria como se fosse de outro trabalho.
func TestFindByCode_casaOTituloENaoOCorpo(t *testing.T) {
	casos := []struct {
		nome, titulo, body, procura string
		casa                        bool
	}{
		{"título com o código", "[GLCGL] Implementar spec — a régua", "", "GLCGL", true},
		{"minúsculo na busca", "[GLCGL] Implementar spec", "", "glcgl", true},
		{"outro código no título", "[ALTPL] Implementar spec", "", "GLCGL", false},
		{"o código só no CORPO", "[ALTPL] Implementar spec",
			"depende de `GLCGL` para a régua", "GLCGL", false},
		{"código como prefixo de outro", "[GLCGLX] Implementar spec", "", "GLCGL", false},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			alvo := "[" + strings.ToUpper(c.procura) + "]"
			casou := strings.Contains(strings.ToUpper(c.titulo), alvo)
			if casou != c.casa {
				t.Errorf("título %q com alvo %s: casou=%v, queria %v",
					c.titulo, alvo, casou, c.casa)
			}
		})
	}
}

// Código vazio é recusado: sem ele a busca casaria "[]" e devolveria o primeiro card
// qualquer da lista.
func TestFindByCode_recusaCodigoVazio(t *testing.T) {
	c := Client{Repo: "x/y", Labels: []string{"anchors"}}
	for _, vazio := range []string{"", "   ", "\t"} {
		if _, err := c.FindByCode(vazio); err == nil {
			t.Errorf("código %q passou — casaria qualquer card", vazio)
		}
	}
}

func TestComment_recusaCorpoVazio(t *testing.T) {
	c := Client{Repo: "x/y"}
	if err := c.Comment(1, "   "); err == nil {
		t.Error("corpo vazio passou — um comentário em branco não registra nada")
	}
}
