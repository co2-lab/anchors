package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/mapx"
)

// O plano REVISADO precisa avisar quem o lê.
//
// Editar o plano antigo perderia o registro — um plano é ÂNCORA, e editá-lo depois de
// implementado faz o documento descrever algo que não foi o que aconteceu. Por isso a
// revisão é um plano NOVO. O custo dessa escolha é a leitura fora de ordem: quem abre o
// antigo segue uma decisão que foi revista, e é isso que o gate impede.
func TestPlanoRevisadoAvisa(t *testing.T) {
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan, Code: "ANTGO"}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan, Code: "NOVOO",
		Supersedes: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	// SEM o aviso, reprova — e a mensagem diz o que escrever.
	semAviso := "# Plano 0001\n\n## Objetivo\n\nFazer algo.\n"
	v, msg := checkPlanoRevisado(semAviso, antigo, "", g, nil)
	if v != Fail {
		t.Fatalf("plano revisado sem aviso deveria reprovar; veio %v", v)
	}
	if !strings.Contains(msg, "plans/0007.md") {
		t.Errorf("a mensagem deveria NOMEAR quem revisa: %s", msg)
	}

	// COM o aviso no topo, passa.
	comAviso := "# Plano 0001\n\n> **Revisado por** `plans/0007.md` — a fase 2 mudou de ordem.\n\n## Objetivo\n"
	if v, msg := checkPlanoRevisado(comAviso, antigo, "", g, nil); v != Pass {
		t.Errorf("com aviso no topo deveria passar; veio %v (%s)", v, msg)
	}

	// O aviso PRECISA estar no topo: depois de 40 linhas, quem lê de cima para baixo já
	// seguiu a primeira decisão antes de saber que ela foi revista.
	tarde := "# Plano 0001\n" + strings.Repeat("\ntexto\n", 30) +
		"\n> **Revisado por** `plans/0007.md` — tarde demais.\n"
	if v, _ := checkPlanoRevisado(tarde, antigo, "", g, nil); v != Fail {
		t.Errorf("aviso fora do topo não avisa ninguém; veio %v", v)
	}
}

// `supersedes` apontando para plano que não existe é uma revisão que não revisa nada — e
// o plano que se pretendia revisar continua sendo lido como vigente.
func TestSupersedesParaPlanoInexistente(t *testing.T) {
	orfao := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan,
		Supersedes: []string{"plans/nao-existe.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{orfao}}
	v, msg := checkPlanoRevisado("# Plano\n", orfao, "", g, nil)
	if v != Fail {
		t.Fatalf("supersedes órfão deveria reprovar; veio %v", v)
	}
	if !strings.Contains(msg, "nao-existe.md") {
		t.Errorf("a mensagem deveria nomear o alvo ausente: %s", msg)
	}
}

// Plano que ninguém revisa e que não revisa ninguém: nada a confrontar.
func TestPlanoSemRevisaoPula(t *testing.T) {
	n := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	g := &mapx.Graph{Nodes: []mapx.Node{n}}
	if v, _ := checkPlanoRevisado("# Plano\n", n, "", g, nil); v != Skip {
		t.Errorf("plano sem revisão deveria pular; veio %v", v)
	}
}

// QUEM REVISA é lembrado do passo que falta, e ele é no OUTRO arquivo.
//
// O gate cobra o aviso do lado do plano REVISADO — e quem escreve o revisor só descobriria
// isso quando o outro reprovasse, num arquivo que ele talvez nem tenha aberto. O lembrete
// alcança a pessoa no momento em que ela está fazendo a revisão.
func TestQuemRevisaEhLembradoDoAviso(t *testing.T) {
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan,
		Supersedes: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	// O plano REVISOR recebe o lembrete — é pendência, não falha: ele não fez nada de
	// errado, só falta um passo que está noutro arquivo.
	v, msg := checkPlanoRevisado("# Plano 0007\n", novo, "", g, nil)
	if v != Pending {
		t.Fatalf("quem revisa deveria receber PENDING com o lembrete; veio %v", v)
	}
	if !strings.Contains(msg, "plans/0001.md") || !strings.Contains(msg, "Revisado por") {
		t.Errorf("o lembrete deveria dizer ONDE escrever e O QUÊ: %s", msg)
	}

	// E o lembrete SOME quando o aviso existe do outro lado — senão ele viraria ruído
	// permanente, e ruído permanente é o que treina a equipe a ignorar o gate.
	//
	// O gate roda por nó: com o aviso posto em 0001, o veredito DELE passa. Aqui
	// verificamos que o revisor não fica pendente para sempre por causa disso.
	comAviso := "# Plano 0001\n\n> **Revisado por** `plans/0007.md` — mudou a ordem.\n"
	if v, _ := checkPlanoRevisado(comAviso, antigo, "", g, nil); v != Pass {
		t.Errorf("o plano revisado com aviso deveria passar; veio %v", v)
	}
}

// O LEMBRETE SOME quando o aviso é escrito. Um lembrete que não some vira ruído
// permanente, e ruído permanente é o que treina a equipe a ignorar o gate — o oposto do
// que ele existe para produzir.
func TestLembreteSomeQuandoOAvisoExiste(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "plans"), 0o755); err != nil {
		t.Fatal(err)
	}
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan,
		Supersedes: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	// SEM o aviso no arquivo: o revisor recebe o lembrete.
	if err := os.WriteFile(filepath.Join(root, "plans/0001.md"),
		[]byte("# Plano 0001\n\n## Objetivo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, _ := checkPlanoRevisado("# Plano 0007\n", novo, root, g, nil); v != Pending {
		t.Fatalf("sem o aviso, o revisor deveria receber o lembrete; veio %v", v)
	}

	// COM o aviso: o lembrete some.
	if err := os.WriteFile(filepath.Join(root, "plans/0001.md"),
		[]byte("# Plano 0001\n\n> **Revisado por** `plans/0007.md` — mudou.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, msg := checkPlanoRevisado("# Plano 0007\n", novo, root, g, nil); v == Pending {
		t.Errorf("com o aviso escrito, o lembrete não pode continuar; veio %v (%s)", v, msg)
	}
}
