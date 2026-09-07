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

// `in-review` NÃO é retomado. Ele é trabalho EM CURSO tanto quanto `in-progress` — a
// diferença entre os dois é QUE trabalho acontece, não SE acontece. Um card ali tem um
// revisor, e devolvê-lo ao dono original o faz reabrir o que já entregou.
//
// Medido: o card #321, entregue e em revisão, voltava a cada `anchors next` — a fila
// parecia ter trabalho e o próximo card de verdade nunca era oferecido.
func TestMine_soRetomaOQueEstaPelaMetade(t *testing.T) {
	casos := []struct {
		estado string
		retoma bool
	}{
		{StateInProgress, true},     // pela metade: é meu, e eu continuo
		{StateInReview, false},      // em curso, com OUTRO: não é meu trabalho
		{StateReadyToReview, false}, // parado, mas vem pela FILA e não pela retomada
		{StateToDo, false},          // idem
		{"anchors:ready-to-test", false},
	}
	for _, caso := range casos {
		t.Run(caso.estado, func(t *testing.T) {
			c := Card{Labels: []string{"anchors", caso.estado}}
			got := has(c.Labels, StateInProgress)
			if got != caso.retoma {
				t.Errorf("%s: retoma=%v, queria %v", caso.estado, got, caso.retoma)
			}
		})
	}
}

// `emCurso` cobre os dois estados de trabalho acontecendo — é a régua que separa o que o
// claim pode oferecer do que já tem dono trabalhando.
func TestEmCurso(t *testing.T) {
	for _, e := range []string{StateInProgress, StateInReview} {
		if !emCurso(Card{Labels: []string{e}}) {
			t.Errorf("%s devia contar como trabalho em curso", e)
		}
	}
	for _, e := range []string{StateToDo, StateReadyToReview} {
		if emCurso(Card{Labels: []string{e}}) {
			t.Errorf("%s está PARADO — o claim pode oferecê-lo", e)
		}
	}
}
