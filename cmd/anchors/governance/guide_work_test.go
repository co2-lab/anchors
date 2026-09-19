package governance

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// E o guia PRECISA cobrir o que o card deixou de dizer, senão a centralização perdeu
// justamente a informação que motivou tudo.
func TestGuiaDeTrabalhoCobreOAchadoDoAgente(t *testing.T) {
	// O que se confronta são NOMES DE COMANDO E FLAG — o vocabulário do próprio sistema,
	// que só muda quando a interface muda (e aí o teste falhar é o comportamento certo).
	//
	// PROSA fica de fora de propósito. A primeira versão exigia "MESMO PR", e quando a
	// frase quebrou de linha eu "consertei" embutindo o `\n` na assertiva — deixando o
	// teste preso à LARGURA DO PARÁGRAFO. Qualquer reflow, reescrita ou tradução do guia
	// o quebraria, sem que nada tivesse piorado de verdade.
	for _, exigido := range []string{
		"anchors escalate",        // como registrar o achado
		"--card",                  // o vínculo, sem o qual ele nasce solto
		"--for-user",              // a saída para quando muda a direção
		"anchors judge --pending", // o que barra o commit
		"@TBD",                    // a peça que ainda não nasceu
		"anchors:under-",          // a label que amarra o achado ao trabalho
	} {
		if !strings.Contains(WorkGuide, exigido) {
			t.Errorf("o guia de trabalho deve cobrir %q — o card não diz mais isso", exigido)
		}
	}
}

// O guia tinha de dizer que EMPURRAR NÃO É ENTREGAR, e não dizia. Um agente de outro dev
// diagnosticou certo por que o CI reprovou (mapa gerado antes do commit que consertava a
// config), reconstruiu, empurrou — e encerrou o turno com "aguardando a nova rodada". O
// card ficou `in-progress` com o nome dele nele.
//
// O guia terminava no `pr-body`: nada dizia como esperar, nem que a espera não é uma
// parada. Este teste cobra as duas coisas que faltavam — o comando que BLOQUEIA até o
// veredito, e a recusa explícita da terceira saída (relatar e parar).
func TestGuiaDeTrabalhoCobraEsperarOVereditoDoCI(t *testing.T) {
	for _, want := range []string{
		"--watch",       // o comando que espera pelo agente
		"work review",   // o que falta no PR que é seu
		"escalate",      // a única saída legítima que não é o verde
		"in-progress",   // o custo de parar no meio: o card fica com seu nome
		"half the work", // o diagnóstico correto é metade do trabalho
	} {
		if !strings.Contains(WorkGuide, want) {
			t.Errorf("guia de trabalho deveria cobrir a espera do CI (esperava %q)", want)
		}
	}
}

// A regra é sobre CONTINUAR, e o teste acima casaria mesmo que o guia dissesse o oposto.
// Este confronta a afirmação: o guia tem de negar que empurrar encerra o turno.
func TestGuiaDeTrabalhoNegaQueEmpurrarEncerraOTurno(t *testing.T) {
	if !strings.Contains(WorkGuide, "is not a stopping point") {
		t.Error("o guia deveria NEGAR explicitamente que empurrar um commit é um ponto de parada")
	}
	if !strings.Contains(WorkGuide, "does not close the card") {
		t.Error("o guia deveria dizer que abrir o PR não fecha o card")
	}
}

// QUEM FECHA O CARD É O MERGE, e o guia não dizia.
//
// Ele dizia "abrir o PR não fecha o card" — verdadeiro e insuficiente: não proibia fechar
// à mão, e fechar à mão parece arrumação (o trabalho está pronto, o PR está aberto, o card
// "já era").
//
// Medido: um projeto acumulou 25 PRs verdes esperando revisão com ZERO cards em
// `ready-to-review`. Os cards tinham sido fechados um minuto ANTES de o PR ser aberto —
// o board dizia que não havia nada para revisar enquanto 25 trabalhos esperavam, e o
// claim, que procura revisão primeiro, entregava trabalho novo.
func TestGuiaDeTrabalhoProibeFecharOCardAMao(t *testing.T) {
	if !strings.Contains(WorkGuide, "YOU DO NOT CLOSE THE CARD") {
		t.Error("o guia precisa dizer explicitamente que o agente não fecha o card")
	}
	// A CONSEQUÊNCIA, e não só a proibição: uma regra sem o porquê se descumpre na
	// primeira vez que parece atrapalhar.
	for _, quer := range []string{"review queue empties", "NEW work", "25 green PRs"} {
		if !strings.Contains(WorkGuide, quer) {
			t.Errorf("o guia deveria explicar o efeito de fechar cedo (esperava %q)", quer)
		}
	}
	// E QUEM fecha: sem isso o agente fica sem saber o que acontece com o card.
	if !strings.Contains(WorkGuide, "Who closes it is the MERGE") {
		t.Error("o guia precisa dizer quem fecha o card, não só quem não fecha")
	}
}

// O NOME DA LABEL no guia tem de ser o REAL.
//
// O guia dizia `anchors:sob-<n>` — o nome ANTERIOR à migração para inglês. O prefixo real é
// `anchors:under-`, e quem seguisse o guia procuraria no board uma label que não existe.
//
// Cinco lugares carregavam o nome velho, incluindo este arquivo de teste (que o cobrava
// errado) e o comentário do `escalate` que descreve o mecanismo.
func TestGuiaDeTrabalhoUsaONomeRealDaLabel(t *testing.T) {
	if strings.Contains(WorkGuide, "anchors:sob-") {
		t.Error("o guia usa o nome ANTIGO da label — quem o seguir procura no board algo " +
			"que não existe")
	}
	if !strings.Contains(WorkGuide, initx.PrefixoLabelSob) {
		t.Errorf("o guia deveria usar %q, o prefixo real", initx.PrefixoLabelSob)
	}
}

// O GUIA PRECISA ENSINAR O CLAIM — ele é o primeiro comando.
//
// O guia se anuncia como "a régua de quem pegou um card" e nunca ensinava COMO pegar.
// Mencionava `status`, `check`, `escalate`, `judge`, `pr-body` — oito comandos — e não o
// `claim`.
//
// MEDIDO no projeto de referência: 38 de 44 PRs abertos tinham o card ainda em `to-do`, e
// a fila de revisão mostrava UM item enquanto 46 trabalhos esperavam. Os agentes escolhiam
// o card à mão, implementavam e abriam PR — fazendo o que o guia ensinava.
//
// O que se perde não é o registro: é a FILA. O claim serve `ready-to-review` antes de
// `to-do`, e um card que nunca entra na coluna de revisão faz o próximo agente encontrá-la
// vazia.
func TestGuiaEnsinaOClaimComoPrimeiroPasso(t *testing.T) {
	if !strings.Contains(WorkGuide, "anchors next") {
		t.Fatal("o guia não menciona o `anchors next` — quem o seguir escolhe o card à " +
			"mão, e o board deixa de descrever quem está com o quê")
	}

	// ANTES de tudo: o claim precisa vir na primeira seção, não perdido no meio. Um
	// comando ensinado depois de "como escrever a spec" chega tarde demais.
	iClaim := strings.Index(WorkGuide, "anchors next")
	iOrdem := strings.Index(WorkGuide, "## A ordem")
	if iOrdem > 0 && iClaim > iOrdem {
		t.Error("o `next` é ensinado depois da seção de ordem do board — ele é o " +
			"PRIMEIRO comando, e ensiná-lo tarde é o mesmo que não ensinar")
	}

	// E o PORQUÊ: "rode este comando" sem a razão vira passo cerimonial, e é o primeiro
	// que alguém pula quando tem pressa.
	if !strings.Contains(WorkGuide, "find it empty") {
		t.Error("o guia não diz POR QUE o claim importa — sem a razão (a fila de revisão " +
			"que fica vazia) ele vira cerimônia, e cerimônia se pula")
	}
}

// ESCALAR SEM LER A FILA CRIA PERGUNTA REPETIDA.
//
// Nada impede N agentes de escalarem a MESMA decisão: cada um acha o problema trabalhando,
// escreve com o contexto da sua hora, e não tem como saber que outro já escalou.
//
// MEDIDO no projeto de referência: 11 escaladas numa hora, 6 na seguinte. Uma triagem
// reagrupou as 23 abertas e concluiu que eram SETE decisões — o resto era o mesmo assunto
// por outro ângulo. Quem decide não lê 23 relatos para achar 7 perguntas, e enquanto a
// fila cresce assim, a escalada de quem abriu também espera.
//
// NÃO É O AVISO QUE JÁ EXISTE. O `escalate` confere cards abertos sobre o mesmo arquivo
// com a label do FLUXO ("pode haver alguém nisso") — é colisão de trabalho em curso. Aqui
// a pergunta é outra: a fila de DECISÃO, e o que fazer quando o assunto já está lá.
//
// POR QUE O GUIA E NÃO O COMANDO: só quem leu os dois textos sabe se são a mesma pergunta.
// Comparar o `--about` erra (mesmo arquivo não é mesmo assunto — a escala de valores e
// onde a preferência mora são decisões distintas sobre a mesma spec) e casar palavras erra
// pior, porque o Anchors é multi-idioma e dois agentes descrevem o mesmo achado com
// palavras diferentes. Um aviso automático aqui teria a FORMA de régua e o conteúdo de
// chute.
//
// O que se confronta é o VOCABULÁRIO — a label e o comando de consulta —, nunca a prosa:
// a assertiva sobre frase prende o teste à largura do parágrafo e quebra em qualquer
// reflow ou tradução.
func TestGuiaDeTrabalhoMandaLerAFilaAntesDeEscalar(t *testing.T) {
	for _, exigido := range []string{
		"anchors:needs-user", // a fila que já existe
		"gh issue list",      // como olhá-la antes de abrir outra
	} {
		if !strings.Contains(WorkGuide, exigido) {
			t.Errorf("o guia não manda consultar %q antes de escalar — sem isso a fila de "+
				"decisão cresce com a mesma pergunta escrita N vezes", exigido)
		}
	}

	// A CONSULTA tem de vir antes do `--for-user` no texto: um agente lê em ordem, e
	// instrução depois do comando é instrução que chega tarde.
	iConsulta := strings.Index(WorkGuide, "gh issue list --label anchors:needs-user")
	iEscala := strings.Index(WorkGuide, "--for-user")
	if iConsulta < 0 || iEscala < 0 {
		t.Fatal("não achei a consulta e a escalada no guia — o vocabulário mudou")
	}
	if iConsulta > iEscala {
		t.Error("a consulta à fila aparece DEPOIS do `--for-user` no guia — quem lê em " +
			"ordem já abriu o card quando chega na instrução de conferir a fila")
	}
}
