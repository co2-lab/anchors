package main

import (
	"regexp"

	"github.com/spf13/cobra"
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
		if !strings.Contains(workGuide, exigido) {
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
		"--watch",     // o comando que espera pelo agente
		"work review", // o que falta no PR que é seu
		"escalate",    // a única saída legítima que não é o verde
		"in-progress", // o custo de parar no meio: o card fica com seu nome
		"metade",      // o diagnóstico correto é metade do trabalho
	} {
		if !strings.Contains(workGuide, want) {
			t.Errorf("guia de trabalho deveria cobrir a espera do CI (esperava %q)", want)
		}
	}
}

// A regra é sobre CONTINUAR, e o teste acima casaria mesmo que o guia dissesse o oposto.
// Este confronta a afirmação: o guia tem de negar que empurrar encerra o turno.
func TestGuiaDeTrabalhoNegaQueEmpurrarEncerraOTurno(t *testing.T) {
	if !strings.Contains(workGuide, "não é um ponto de parada") {
		t.Error("o guia deveria NEGAR explicitamente que empurrar um commit é um ponto de parada")
	}
	if !strings.Contains(workGuide, "não fecha o card") {
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
	if !strings.Contains(workGuide, "NÃO FECHA O CARD") {
		t.Error("o guia precisa dizer explicitamente que o agente não fecha o card")
	}
	// A CONSEQUÊNCIA, e não só a proibição: uma regra sem o porquê se descumpre na
	// primeira vez que parece atrapalhar.
	for _, quer := range []string{"fila de revisão esvazia", "trabalho NOVO", "25 PRs"} {
		if !strings.Contains(workGuide, quer) {
			t.Errorf("o guia deveria explicar o efeito de fechar cedo (esperava %q)", quer)
		}
	}
	// E QUEM fecha: sem isso o agente fica sem saber o que acontece com o card.
	if !strings.Contains(workGuide, "Quem fecha é o MERGE") {
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
	if strings.Contains(workGuide, "anchors:sob-") {
		t.Error("o guia usa o nome ANTIGO da label — quem o seguir procura no board algo " +
			"que não existe")
	}
	if !strings.Contains(workGuide, initx.PrefixoLabelSob) {
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
	if !strings.Contains(workGuide, "anchors next") {
		t.Fatal("o guia não menciona o `anchors next` — quem o seguir escolhe o card à " +
			"mão, e o board deixa de descrever quem está com o quê")
	}

	// ANTES de tudo: o claim precisa vir na primeira seção, não perdido no meio. Um
	// comando ensinado depois de "como escrever a spec" chega tarde demais.
	iClaim := strings.Index(workGuide, "anchors next")
	iOrdem := strings.Index(workGuide, "## A ordem")
	if iOrdem > 0 && iClaim > iOrdem {
		t.Error("o `next` é ensinado depois da seção de ordem do board — ele é o " +
			"PRIMEIRO comando, e ensiná-lo tarde é o mesmo que não ensinar")
	}

	// E o PORQUÊ: "rode este comando" sem a razão vira passo cerimonial, e é o primeiro
	// que alguém pula quando tem pressa.
	if !strings.Contains(workGuide, "encontrá-la vazia") {
		t.Error("o guia não diz POR QUE o claim importa — sem a razão (a fila de revisão " +
			"que fica vazia) ele vira cerimônia, e cerimônia se pula")
	}
}

// TODO COMANDO QUE O GUIA MANDA RODAR PRECISA EXISTIR.
//
// A v0.1.120 publicou um guia ensinando `anchors claim` — que NÃO EXISTE. O comando é
// `anchors next`; `claim` é o nome do PIPELINE (`anchors-claim.yml`), e eu o confundi com
// o do CLI.
//
// O custo foi imediato: cinco agentes receberam a instrução e o primeiro a rodá-la levou
// `erro: unknown command "claim"`. Um guia que manda rodar o que não existe é pior que um
// guia omisso — o omisso faz procurar, este faz desistir.
//
// A régua confronta o guia com os comandos REGISTRADOS no cobra, que é a fonte de verdade:
// se um comando for renomeado, ela acusa no mesmo commit.
func TestGuiaSoCitaComandoQueExiste(t *testing.T) {
	registrados := map[string]bool{}
	var colhe func(*cobra.Command)
	colhe = func(c *cobra.Command) {
		registrados[c.Name()] = true
		for _, f := range c.Commands() {
			colhe(f)
		}
	}
	colhe(newRootCmd())

	// `anchors <algo>` no texto do guia — o `<algo>` precisa ser um comando de verdade.
	// A captura para no primeiro não-letra: `anchors check --changed` cita `check`.
	re := regexp.MustCompile(`anchors ([a-z][a-z-]*)`)
	vistos := map[string]bool{}
	for _, m := range re.FindAllStringSubmatch(workGuide, -1) {
		nome := m[1]
		if vistos[nome] {
			continue
		}
		vistos[nome] = true
		// `anchors <artefato>` é um placeholder do próprio guia, não um comando.
		if nome == "guide" || registrados[nome] {
			continue
		}
		t.Errorf("o guia manda rodar `anchors %s`, e esse comando não existe — quem o "+
			"seguir recebe `unknown command` e para", nome)
	}
	if len(vistos) == 0 {
		t.Fatal("nenhum `anchors <comando>` no guia — o regex quebrou e o teste passaria vazio")
	}
}
