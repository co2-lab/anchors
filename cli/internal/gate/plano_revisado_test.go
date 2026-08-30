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
		Revises: []string{"plans/0001.md"}}
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

	// COM o aviso no topo mas SEM marcar as partes: pendência. O aviso diz que o plano
	// mudou e não diz ONDE — quem lê a fase 3 não sabe se ela é uma das que mudaram, e
	// descobrir custa reler o plano inteiro.
	soTopo := "# Plano 0001\n\n> `@revised-by: plans/0007.md` — a fase 2 mudou de ordem.\n\n## Objetivo\n"
	if v, msg := checkPlanoRevisado(soTopo, antigo, "", g, nil); v != Pending {
		t.Errorf("aviso de topo sem marcar as partes é pendência; veio %v (%s)", v, msg)
	}

	// COM o aviso E a parte marcada, passa.
	completo := soTopo + "\n## Fases\n\n### ANTGO-F02 — a régua\n\n" +
		"> `@amended-by: plans/0007.md` — esta fase passou a vir depois da F03.\n"
	if v, msg := checkPlanoRevisado(completo, antigo, "", g, nil); v != Pass {
		t.Errorf("com aviso e parte marcada deveria passar; veio %v (%s)", v, msg)
	}

	// O aviso PRECISA estar no topo: depois de 40 linhas, quem lê de cima para baixo já
	// seguiu a primeira decisão antes de saber que ela foi revista.
	tarde := "# Plano 0001\n" + strings.Repeat("\ntexto\n", 30) +
		"\n> `@revised-by: plans/0007.md` — tarde demais.\n" +
		"> `@amended-by: plans/0007.md` — a parte está marcada; o problema é o topo.\n"
	if v, _ := checkPlanoRevisado(tarde, antigo, "", g, nil); v != Fail {
		t.Errorf("aviso fora do topo não avisa ninguém; veio %v", v)
	}
}

// `revises` apontando para plano que não existe é uma revisão que não revisa nada — e
// o plano que se pretendia revisar continua sendo lido como vigente.
func TestRevisesParaPlanoInexistente(t *testing.T) {
	orfao := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan,
		Revises: []string{"plans/nao-existe.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{orfao}}
	v, msg := checkPlanoRevisado("# Plano\n", orfao, "", g, nil)
	if v != Fail {
		t.Fatalf("revises órfão deveria reprovar; veio %v", v)
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
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}

	// O plano REVISOR recebe o lembrete — é pendência, não falha: ele não fez nada de
	// errado, só falta um passo que está noutro arquivo.
	v, msg := checkPlanoRevisado("# Plano 0007\n", novo, "", g, nil)
	if v != Pending {
		t.Fatalf("quem revisa deveria receber PENDING com o lembrete; veio %v", v)
	}
	if !strings.Contains(msg, "plans/0001.md") || !strings.Contains(msg, "@revised-by") {
		t.Errorf("o lembrete deveria dizer ONDE escrever e O QUÊ: %s", msg)
	}

	// E o lembrete SOME quando o aviso existe do outro lado — senão ele viraria ruído
	// permanente, e ruído permanente é o que treina a equipe a ignorar o gate.
	//
	// O gate roda por nó: com o aviso posto em 0001, o veredito DELE passa. Aqui
	// verificamos que o revisor não fica pendente para sempre por causa disso.
	comAviso := "# Plano 0001\n\n> `@revised-by: plans/0007.md` — mudou a ordem.\n\n" +
		"### ANTGO-F02\n\n> `@amended-by: plans/0007.md` — esta fase mudou de ordem.\n"
	if v, _ := checkPlanoRevisado(comAviso, antigo, "", g, nil); v != Pass {
		t.Errorf("o plano revisado com aviso e parte marcada deveria passar; veio %v", v)
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
		Revises: []string{"plans/0001.md"}}
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
		[]byte("# Plano 0001\n\n> `@revised-by: plans/0007.md` — mudou.\n\n> `@amended-by: plans/0007.md` — a F01.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if v, msg := checkPlanoRevisado("# Plano 0007\n", novo, root, g, nil); v == Pending {
		t.Errorf("com o aviso escrito, o lembrete não pode continuar; veio %v (%s)", v, msg)
	}
}

// AS DUAS FORMAS de marcar uma parte, e a diferença entre elas é o que o registro
// precisa preservar.
//
// Numa parte JÁ IMPLEMENTADA, o texto fica: ele descreve o que FOI feito, e reescrevê-lo
// faria o registro mentir sobre o passado. Numa parte AINDA NÃO implementada, o texto é
// reescrito com o comportamento novo — mas o original fica ao lado, porque quem leu o
// plano antes precisa saber o que mudou.
func TestMarcacaoDeParteAceitaAsDuasFormas(t *testing.T) {
	antigo := mapx.Node{ID: "plans/0001.md", Kind: mapx.KindPlan}
	novo := mapx.Node{ID: "plans/0007.md", Kind: mapx.KindPlan,
		Revises: []string{"plans/0001.md"}}
	g := &mapx.Graph{Nodes: []mapx.Node{antigo, novo}}
	topo := "# Plano 0001\n\n> `@revised-by: plans/0007.md` — mudou.\n\n"

	// JÁ IMPLEMENTADA: o texto original permanece, com o aviso ABAIXO.
	implementada := topo + "## Fases\n\n### ANTGO-F01 — a árvore\n\nCria os pacotes.\n\n" +
		"> `@amended-by: plans/0007.md` — daqui em diante os pacotes ganham mutação.\n" +
		"> O texto acima descreve o que FOI implementado.\n"
	if v, msg := checkPlanoRevisado(implementada, antigo, "", g, nil); v != Pass {
		t.Errorf("parte implementada com aviso abaixo deveria passar; veio %v (%s)", v, msg)
	}

	// NÃO IMPLEMENTADA: o texto é o novo, e o original fica preservado no aviso.
	futura := topo + "## Fases\n\n### ANTGO-F04 — o CI com mutação\n\n" +
		"> `@amended-by: plans/0007.md`.\n" +
		"> **Era:** o CI roda instalar, lint, build e teste.\n" +
		"> **Por quê:** cobertura não responde se o teste PROVA a linha.\n"
	if v, msg := checkPlanoRevisado(futura, antigo, "", g, nil); v != Pass {
		t.Errorf("parte futura com o original preservado deveria passar; veio %v (%s)", v, msg)
	}
}
