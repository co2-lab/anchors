package ops

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-isatty"

	"github.com/co2-lab/anchors/internal/initx"
)

// discoverStep é o passo do `init` que reconhece a fase que ainda não aconteceu.
//
// `anchors init` INFERE a Estrutura do disco. Num diretório vazio não há o que inferir —
// e ele pergunta "quais diretórios de código tratar como camadas?", cuja resposta honesta
// é "nenhum ainda". A fase DESCOBRIR (`anchors guide project`) existe para isso: uma
// entrevista de 5 etapas que produz PROJECT.md e INSIGHTS.md, e só depois dela o `init`
// tem o que perguntar.
//
// O que o Anchors DIZ aqui depende de quem está do outro lado, e a diferença não é de
// tom: para uma pessoa a saída é INSTRUÇÃO (a fase existe, e aqui está como começá-la);
// para uma IA é ORDEM DE SERVIÇO (leia o guide e conduza a entrevista nesta conversa).
// Um texto só para os dois falha nos dois — a pessoa recebe instruções que não sabe
// executar, e a IA recebe um convite quando precisava de uma tarefa.
//
// Devolve false quando o `init` deve PARAR: só acontece quando um prompt não pôde rodar.
func discoverStep(root string, p *initx.Proposal) bool {
	if !initx.PrecisaDescobrir(root, p) {
		return true
	}
	if initx.DetectOperator(hasTTY(), os.Getenv) == initx.OperadorIA {
		printWorkOrder()
		return true
	}
	return instructPerson(root)
}

// printWorkOrder é o texto para uma IA operando o CLI. Não pergunta nada: quem lê
// isto pode agir, e o que ela precisa é da TAREFA, com os passos na ordem.
func printWorkOrder() {
	fmt.Println(`
┌─ The DISCOVER phase did not happen in this project ─────────────────────────┐

There is no PROJECT.md and no code: "init" INFERS the Structure from the disk, and here
there is nothing to infer. The questions below would come out without a good answer, and the
stack decision would end up taken by accident in the first file someone creates.

YOU (the agent operating this CLI) must, IN THIS ORDER:

  1. run  anchors guide project  and follow that ruler to the letter;
  2. conduct the 5-stage interview IN THIS CONVERSATION, with the user —
     one question at a time, waiting for the answer before the next;
  3. do the inconsistency review at the end;
  4. write PROJECT.md (the decisions) and INSIGHTS.md (the transcript) at the root;
  5. only then run  anchors init  again — its questions are answered
     by what PROJECT.md decided.

The interview runs in the conversation, never in a background worker: it is the user who
answers, and delegating to a subagent that does not talk to them produces no answer.

└─────────────────────────────────────────────────────────────────────────────┘`)
	fmt.Println("\nProceeding with the init anyway (the answers will have no basis).")
	fmt.Println()
}

// instructPerson é o caminho de quem está sozinho no terminal. O Anchors não conduz a
// entrevista (não embute modelo — `guide project` diz isso explicitamente), então ele
// instrui e, quando sabe qual IA está instalada, oferece abri-la com o prompt pronto.
func instructPerson(root string) bool {
	fmt.Println(`
⚠ The DISCOVER phase has not happened in this project yet.

  ` + "`init`" + ` INFERS the Structure from what is on disk, and here there is nothing to infer:
  without PROJECT.md and without code, the questions below come out without a good answer — and the
  stack decision ends up taken by accident in the first file someone creates.

  Before it there is a 5-stage interview (purpose → language → architecture
  → structure → tooling), conducted by an AI, that produces PROJECT.md and
  INSIGHTS.md. Anchors does not conduct it: it embeds no model, it only provides the ruler.`)

	nome := initx.AgentName(os.Getenv)
	comando := initx.CommandToOpenAI(os.Getenv)

	if len(comando) > 0 {
		fmt.Printf("\n  I detected %s on this machine.\n", nome)
		if askConfirmDefault("Open "+nome+" now with the interview prompt?", true) {
			if erroDePrompt {
				return false
			}
			if err := openAI(root, comando); err != nil {
				fmt.Printf("  ⚠ could not open it: %v\n", err)
				printStepByStep()
			}
			return true
		}
		if erroDePrompt {
			return false
		}
	}
	printStepByStep()
	return true
}

// printStepByStep é a saída para quem prefere conduzir sozinho — ou para quando o
// Anchors não sabe qual IA abrir. O prompt vai inteiro, pronto para colar.
func printStepByStep() {
	fmt.Println(`
  To do it yourself:

    1. open your AI tool in this directory;
    2. paste the prompt below;
    3. answer the 5 stages (one question at a time);
    4. come back here and run ` + "`anchors init`" + ` again.

  ── prompt ───────────────────────────────────────────────────────────────────`)
	fmt.Printf("  %s\n", breakAt(initx.PromptDescobrir, 76, "  "))
	fmt.Println("  ─────────────────────────────────────────────────────────────────────────────")
	fmt.Println("\n  Proceeding with the init anyway (the answers will have no basis).")
	fmt.Println()
}

// openAI executa a ferramenta detectada com o prompt. Roda com os fluxos herdados: a
// entrevista É uma conversa, e capturar a saída deixaria o usuário diante de um processo
// mudo que não dá para responder.
//
// Sem shell: o argv vai direto para o processo. O prompt tem aspas e parênteses, e passá-lo
// por `sh -c` faria o escape ser a única coisa entre o texto e o interpretador.
func openAI(root string, argv []string) error {
	c := exec.Command(argv[0], argv[1:]...)
	c.Dir = root
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return c.Run()
}

// breakAt reflui o texto numa largura, prefixando as linhas seguintes. Um prompt de 400
// caracteres numa linha só é impossível de ler no terminal e feio de copiar.
func breakAt(s string, largura int, prefixo string) string {
	var linhas []string
	var atual string
	for _, palavra := range strings.Fields(s) {
		if atual == "" {
			atual = palavra
			continue
		}
		if len(atual)+1+len(palavra) > largura {
			linhas = append(linhas, atual)
			atual = palavra
			continue
		}
		atual += " " + palavra
	}
	if atual != "" {
		linhas = append(linhas, atual)
	}
	return strings.Join(linhas, "\n"+prefixo)
}

// hasTTY diz se há terminal interativo na entrada. É a evidência mais fraca de quem
// opera (um pipe qualquer produz o mesmo resultado), por isso `DetectaOperador` só
// recorre a ela depois de procurar as variáveis que um agente declara.
func hasTTY() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}
