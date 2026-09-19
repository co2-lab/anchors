package governance

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
		fmt.Fprintf(&b, "\n## This role's lens (%s)\n\n%s.\n", s.Role.Title(), lens)
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
	b.WriteString("### Preparing the environment does not ask for authorization\n\n")
	b.WriteString("The PREPARATION commands are idempotent and check before acting. " +
		"In a project\nalready set up they change nothing, and running them is step zero " +
		"of the work:\n\n")
	b.WriteString("    anchors doctor --fix          pipelines, hooks, labels, branch protection\n")
	b.WriteString("    anchors settings role <role> --date <YYYY-MM-DD>\n")
	b.WriteString("    anchors map build             the map the gates confront\n\n")
	b.WriteString("`settings role` accepts the role and the date as ARGUMENTS — without them " +
		"it asks,\nand an agent with no terminal waits for an answer that never comes.\n\n")
	b.WriteString("What DOES ask for authorization is another thing: erasing someone's work, " +
		"force-pushing,\nclosing an issue that is not yours, touching a protected branch. The ruler is " +
		"REVERSIBILITY,\nnot the fact of touching the remote.\n\n")

	b.WriteString("\n## When you do not know\n\n")

	if s.HandlesUserIssues() {
		fmt.Fprintf(&b, "Your role (%s) decides the direction of this product.\n\n", s.Role.Title())
		b.WriteString("Even so, the ruler holds: what changes the project's DIRECTION is " +
			"written, not\ndiscussed. A decision taken mid-session leaves no trace of why " +
			"it was taken,\nand whoever inherits it will have no way to know whether it " +
			"was choice or accident.\n\n")
		b.WriteString("    anchors escalate \"<what needs to change>\" --about <file> " +
			"--for-user\n\n")
		b.WriteString("The escalation comes back to you — and then the decision stays recorded in " +
			"the card, with\nwhat you knew at the time.\n")
		return b.String()
	}

	// O caso que motiva esta seção existir.
	//
	// A frase distingue DECLARADO de NÃO DECLARADO, e a distinção importa para quem lê:
	// dizer "ficou declarado" a quem nunca declarou é afirmar um fato que não aconteceu, e
	// o leitor vai procurar a declaração que não existe.
	if s.Decided() {
		fmt.Fprintf(&b, "**Your role (%s) does NOT decide the direction of this product** — who decides "+
			"is the\n`product-owner` or the `architect`.\n\n", s.Role.Title())
	} else {
		b.WriteString("**You did not declare a role**, and without a role the agent decides " +
			"nothing.\nDeclare one with `anchors settings role`.\n\n")
	}
	b.WriteString("So, faced with an ambiguity, with a spec that does not decide " +
		"enough, or with\na choice that changes the product's behaviour:\n\n")
	b.WriteString("    anchors escalate \"<what needs to be decided>\" --about <file> " +
		"--for-user\n\n")
	b.WriteString("**Do not ask whoever is running you.** That person knows the code and " +
		"will\nanswer — and the answer is reasonable, and becomes a product decision taken by " +
		"someone\nwho had no authority to take it. Without going through the plan, without review, " +
		"and\nwithout a trace that it was decided there.\n\n")
	b.WriteString("The difference between the issue and the question is the RECORD: the issue stays, has " +
		"an owner, and\nwhoever decides reads it when they can. The question vanishes with the session.\n\n")
	b.WriteString("After escalating, **move on to the next card**. The escalation waits for " +
		"whoever\nresolves it, and standing still on it does not resolve it any faster.\n\n")
	b.WriteString("### What is NOT to be escalated\n\n")
	b.WriteString("What the ruler already decides. If the spec answers, it answers — seeking " +
		"confirmation\nof what is written turns escalation into noise, and the noise " +
		"makes the next\nreal finding go unnoticed.\n")
	return b.String()
}

// printAutonomy imprime a seção — usada pelo `guide work` e pelo `guide review`.
func printAutonomy(root string) {
	fmt.Print(autonomyGuide(root))
}
