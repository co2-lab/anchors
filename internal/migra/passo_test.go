package migra

import (
	"strings"
	"testing"
)

// A CORRENTE de passos é o que torna a atualização barata: acrescentar um formato é
// acrescentar um arquivo `formato_N.go`, e nada mais precisa saber que ele existe.
//
// O que estes testes protegem é a propriedade que faz isso funcionar — que não há BURACO.
// Um projeto no formato 2 indo para o 4 sem o passo que produz o 3 receberia o número novo
// com o conteúdo velho, e o arquivo passaria a mentir sobre o próprio formato.

func TestStepsFrom_devolveOsPassosEmOrdem(t *testing.T) {
	original := steps
	defer func() { steps = original }()
	steps = nil
	Register(Step{To: 3, Why: "terceiro"})
	Register(Step{To: 2, Why: "segundo"})

	got, err := StepsFrom(1, 3)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(got) != 2 || got[0].To != 2 || got[1].To != 3 {
		t.Fatalf("esperava 1→2→3 nessa ordem; veio %v", got)
	}
}

// Registrar fora de ordem não pode produzir migração fora de ordem: os arquivos
// `formato_N.go` são carregados na ordem que o Go decidir, não na numérica.
func TestRegister_ordenaIndependenteDaOrdemDeRegistro(t *testing.T) {
	original := steps
	defer func() { steps = original }()
	steps = nil
	Register(Step{To: 4})
	Register(Step{To: 2})
	Register(Step{To: 3})

	for i, s := range steps {
		if s.To != i+2 {
			t.Fatalf("os passos deveriam estar ordenados por formato; veio %v", steps)
		}
	}
}

// O BURACO é erro, e não silêncio. Sem o passo que produz o 3, um projeto no 2 não pode ir
// para o 4 fingindo que o 3 não existiu — e a mensagem tem de dizer QUAL passo falta, senão
// quem lê não sabe o que escrever.
func TestStepsFrom_buracoNaCorrenteEhErro(t *testing.T) {
	original := steps
	defer func() { steps = original }()
	steps = nil
	Register(Step{To: 2})
	Register(Step{To: 4}) // falta o 3

	_, err := StepsFrom(1, 4)
	if err == nil {
		t.Fatal("um buraco na corrente tem de ser erro")
	}
	if !strings.Contains(err.Error(), "formato 3") {
		t.Errorf("a mensagem deveria nomear o passo que falta; veio: %v", err)
	}
}

// Nada a fazer é resposta válida, e não erro: um projeto já no formato atual chama o
// migrador a cada comando, e responder erro faria o `check` reprovar por estar em dia.
func TestStepsFrom_jaNoFormatoAtualNaoTemPasso(t *testing.T) {
	got, err := StepsFrom(2, 2)
	if err != nil {
		t.Fatalf("estar em dia não é erro: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("nada a migrar deveria devolver zero passos; veio %d", len(got))
	}
}

// O PASSO 1→2 está registrado de verdade — não só a mecânica funciona, o conteúdo existe.
// Sem esta asserção, remover o `formato_2.go` deixaria a suíte verde e o produto sem
// migração nenhuma.
func TestPassoDoFormato2_estaRegistradoEConverte(t *testing.T) {
	var p *Step
	for i := range steps {
		if steps[i].To == 2 {
			p = &steps[i]
		}
	}
	if p == nil {
		t.Fatal("o passo que produz o formato 2 não está registrado")
	}
	if p.RenameKeys["anchors.graph.yaml"]["julgamentos"] != "judgments" {
		t.Error("o passo 1→2 tem de converter `julgamentos` — é o que preserva os carimbos")
	}
	if p.RenameValues["anchors.graph.yaml"]["gate"]["regra-cumprida"] != "rule-fulfilled" {
		t.Error("o passo 1→2 tem de converter os NOMES DE GATE gravados nos carimbos")
	}
	if p.Why == "" {
		t.Error("todo passo precisa dizer o que mudou — é o que quem revisa o diff lê antes de abri-lo")
	}
}
