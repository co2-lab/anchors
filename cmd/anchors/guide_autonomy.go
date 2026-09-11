package main

import (
	"fmt"
	"strings"

	"github.com/co2-lab/anchors/internal/settings"
)

// --- o agente que NÃO decide o produto ---
//
// O `settings user-issues` fecha uma porta: o claim não entrega card escalonado a quem não
// declarou que decide. Mas a porta que mais se usa é outra, e ela não tem label — é o
// agente PERGUNTAR ao dev que o está rodando.
//
// A pergunta parece inofensiva e não é. O dev que roda o agente conhece o código, e vai
// responder: "sim, pode assumir que o campo é opcional", "usa o padrão mesmo". A resposta é
// razoável e vira decisão de produto tomada por quem não tinha autoridade — sem passar pelo
// plano, sem revisão, e sem deixar rastro de que foi decidido ali.
//
// O escalonamento existe para isso, e a diferença é o REGISTRO: uma issue fica, tem dono, e
// quem decide a lê quando puder. Uma pergunta no terminal some com a sessão.
//
// Por isso o guia MUDA conforme a declaração local. Não é um aviso a mais no fim do texto —
// quem não decide o produto lê uma instrução diferente, no lugar onde ela importa.

// autonomyGuide devolve a seção do guia sobre o que fazer diante do que não se sabe.
func autonomyGuide(root string) string {
	s, _ := settings.Load(root)

	var b strings.Builder

	// A LENTE do perfil, quando ele tem uma.
	//
	// "Revisar" não é uma coisa só: quem procura vazamento de dado e quem procura consulta
	// em laço leem o mesmo código com perguntas diferentes. Um revisor sem lente declarada
	// tende a fazer a revisão que SABE fazer — não a que falta.
	if lens := s.Role.Lens(); lens != "" {
		fmt.Fprintf(&b, "\n## A lente deste perfil (%s)\n\n%s.\n", s.Role.Title(), lens)
	}

	// O SEGUNDO tipo de parada, e ele não é escalonamento nenhum.
	//
	// Medido: um dev novo pediu ao agente dele para contribuir. O agente leu o
	// CONTRIBUTING, montou o plano de onboarding CORRETO, e parou no passo 4 pedindo
	// autorização — `doctor --fix` "protege o branch" e "cria labels", e isso se lê como
	// alteração de estado compartilhado. Pediu OK e o turno acabou.
	//
	// A cautela estava certa: uma ação irreversível sobre repositório alheio merece
	// pergunta. O que faltava era saber que aquele comando é idempotente, e que rodá-lo
	// num repositório já montado não muda nada.
	//
	// Sem esta seção o agente escolhe entre dois erros: pedir autorização para tudo
	// (e não começar) ou não pedir para nada (e mexer no que não devia).
	b.WriteString("### Preparar o ambiente não pede autorização\n\n")
	b.WriteString("Os comandos de PREPARAÇÃO são idempotentes e conferem antes de agir. " +
		"Num projeto\njá montado eles não mudam nada, e rodá-los é o passo zero do " +
		"trabalho:\n\n")
	b.WriteString("    anchors doctor --fix          pipelines, hooks, labels, proteção do branch\n")
	b.WriteString("    anchors settings role <perfil> --date <AAAA-MM-DD>\n")
	b.WriteString("    anchors map build             o mapa que os gates confrontam\n\n")
	b.WriteString("O `settings role` aceita o perfil e a data como ARGUMENTO — sem eles " +
		"ele pergunta,\ne um agente sem terminal fica esperando resposta que não vem.\n\n")
	b.WriteString("O que PEDE autorização é outra coisa: apagar trabalho de alguém, " +
		"forçar push,\nfechar issue que não é sua, mexer em branch protegido. A régua é a " +
		"REVERSIBILIDADE,\nnão o fato de tocar o remoto.\n\n")

	b.WriteString("\n## Quando você não souber\n\n")

	if s.HandlesUserIssues() {
		fmt.Fprintf(&b, "Seu perfil (%s) decide o rumo deste produto.\n\n", s.Role.Title())
		b.WriteString("Ainda assim, a régua vale: o que muda a DIREÇÃO do projeto se " +
			"escreve, não se\nconversa. Uma decisão tomada no meio de uma sessão não " +
			"deixa rastro de por que\nfoi tomada, e quem a herdar não terá como saber se " +
			"foi escolha ou acidente.\n\n")
		b.WriteString("    anchors escalate \"<o que precisa mudar>\" --about <arquivo> " +
			"--for-user\n\n")
		b.WriteString("O escalonado volta para você — e aí a decisão fica registrada no " +
			"card, com o\nque você sabia na hora.\n")
		return b.String()
	}

	// O caso que motiva esta seção existir.
	//
	// A frase distingue DECLARADO de NÃO DECLARADO, e a distinção importa para quem lê:
	// dizer "ficou declarado" a quem nunca declarou é afirmar um fato que não aconteceu, e
	// o leitor vai procurar a declaração que não existe.
	if s.Decided() {
		fmt.Fprintf(&b, "**Seu perfil (%s) NÃO decide o rumo deste produto** — quem decide "+
			"é o\n`product-owner` ou o `architect`.\n\n", s.Role.Title())
	} else {
		b.WriteString("**Você não declarou um perfil**, e sem perfil o agente não decide " +
			"nada.\nDeclare com `anchors settings role`.\n\n")
	}
	b.WriteString("Então, diante de uma ambiguidade, de uma spec que não decide o " +
		"suficiente, ou de\numa escolha que muda o comportamento do produto:\n\n")
	b.WriteString("    anchors escalate \"<o que precisa ser decidido>\" --about <arquivo> " +
		"--for-user\n\n")
	b.WriteString("**Não pergunte a quem está rodando você.** A pessoa conhece o código e " +
		"vai\nresponder — e a resposta é razoável, e vira decisão de produto tomada por " +
		"quem\nnão tinha autoridade para tomá-la. Sem passar pelo plano, sem revisão, e " +
		"sem\nrastro de que foi decidido ali.\n\n")
	b.WriteString("A diferença entre a issue e a pergunta é o REGISTRO: a issue fica, tem " +
		"dono, e\nquem decide a lê quando puder. A pergunta some com a sessão.\n\n")
	b.WriteString("Depois de escalar, **siga para o próximo card**. O escalonado espera " +
		"quem o\nresolve, e ficar parado nele não o resolve mais rápido.\n\n")
	b.WriteString("### O que NÃO é para escalar\n\n")
	b.WriteString("O que a régua já decide. Se a spec responde, ela responde — procurar " +
		"confirmação\ndo que está escrito transforma o escalonamento em ruído, e o ruído " +
		"faz o próximo\nachado real passar batido.\n")
	return b.String()
}

// printAutonomy imprime a seção — usada pelo `guide work` e pelo `guide review`.
func printAutonomy(root string) {
	fmt.Print(autonomyGuide(root))
}
