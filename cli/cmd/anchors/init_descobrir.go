package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-isatty"

	"github.com/co2-lab/anchors/internal/initx"
)

// etapaDescobrir é o passo do `init` que reconhece a fase que ainda não aconteceu.
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
func etapaDescobrir(root string, p *initx.Proposal) bool {
	if !initx.PrecisaDescobrir(root, p) {
		return true
	}
	if initx.DetectaOperador(temTTY(), os.Getenv) == initx.OperadorIA {
		imprimeOrdemDeServico()
		return true
	}
	return instruiPessoa(root)
}

// imprimeOrdemDeServico é o texto para uma IA operando o CLI. Não pergunta nada: quem lê
// isto pode agir, e o que ela precisa é da TAREFA, com os passos na ordem.
func imprimeOrdemDeServico() {
	fmt.Println(`
┌─ A fase DESCOBRIR não aconteceu neste projeto ──────────────────────────────┐

Não há PROJECT.md e não há código: o "init" INFERE a Estrutura do disco, e aqui
não há o que inferir. As perguntas a seguir sairiam sem resposta boa, e a decisão
de stack acabaria tomada por acidente no primeiro arquivo que alguém criar.

VOCÊ (o agente que está operando este CLI) deve, NESTA ORDEM:

  1. rodar  anchors guide project  e seguir aquela régua à risca;
  2. conduzir a entrevista de 5 etapas NESTA CONVERSA, com o usuário —
     uma pergunta por vez, esperando a resposta antes da próxima;
  3. fazer a revisão de inconsistências ao final;
  4. escrever PROJECT.md (as decisões) e INSIGHTS.md (a transcrição) na raiz;
  5. só então rodar  anchors init  de novo — as perguntas dele são respondidas
     pelo que o PROJECT.md decidiu.

A entrevista roda na conversa, nunca num worker de background: é o usuário quem
responde, e delegar a um subagente que não fala com ele não produz resposta.

└─────────────────────────────────────────────────────────────────────────────┘`)
	fmt.Println("\nSeguindo com o init assim mesmo (as respostas ficarão sem base).")
	fmt.Println()
}

// instruiPessoa é o caminho de quem está sozinho no terminal. O Anchors não conduz a
// entrevista (não embute modelo — `guide project` diz isso explicitamente), então ele
// instrui e, quando sabe qual IA está instalada, oferece abri-la com o prompt pronto.
func instruiPessoa(root string) bool {
	fmt.Println(`
⚠ A fase DESCOBRIR ainda não aconteceu neste projeto.

  O ` + "`init`" + ` INFERE a Estrutura do que está no disco, e aqui não há o que inferir:
  sem PROJECT.md e sem código, as perguntas a seguir saem sem resposta boa — e a
  decisão de stack acaba tomada por acidente no primeiro arquivo que alguém criar.

  Antes dele existe uma entrevista de 5 etapas (propósito → linguagem → arquitetura
  → estrutura → ferramental), conduzida por uma IA, que produz PROJECT.md e
  INSIGHTS.md. O Anchors não a conduz: ele não embute modelo, só fornece a régua.`)

	nome := initx.NomeDoAgente(os.Getenv)
	comando := initx.ComandoParaAbrirIA(os.Getenv)

	if len(comando) > 0 {
		fmt.Printf("\n  Detectei %s nesta máquina.\n", nome)
		if askConfirmDefault("Abrir "+nome+" agora com o prompt da entrevista?", true) {
			if erroDePrompt {
				return false
			}
			if err := abreIA(root, comando); err != nil {
				fmt.Printf("  ⚠ não deu para abrir: %v\n", err)
				imprimePassoAPasso()
			}
			return true
		}
		if erroDePrompt {
			return false
		}
	}
	imprimePassoAPasso()
	return true
}

// imprimePassoAPasso é a saída para quem prefere conduzir sozinho — ou para quando o
// Anchors não sabe qual IA abrir. O prompt vai inteiro, pronto para colar.
func imprimePassoAPasso() {
	fmt.Println(`
  Para fazer você mesmo:

    1. abra sua ferramenta de IA neste diretório;
    2. cole o prompt abaixo;
    3. responda as 5 etapas (uma pergunta por vez);
    4. volte aqui e rode ` + "`anchors init`" + ` de novo.

  ── prompt ───────────────────────────────────────────────────────────────────`)
	fmt.Printf("  %s\n", quebraEm(initx.PromptDescobrir, 76, "  "))
	fmt.Println("  ─────────────────────────────────────────────────────────────────────────────")
	fmt.Println("\n  Seguindo com o init assim mesmo (as respostas ficarão sem base).")
	fmt.Println()
}

// abreIA executa a ferramenta detectada com o prompt. Roda com os fluxos herdados: a
// entrevista É uma conversa, e capturar a saída deixaria o usuário diante de um processo
// mudo que não dá para responder.
//
// Sem shell: o argv vai direto para o processo. O prompt tem aspas e parênteses, e passá-lo
// por `sh -c` faria o escape ser a única coisa entre o texto e o interpretador.
func abreIA(root string, argv []string) error {
	c := exec.Command(argv[0], argv[1:]...)
	c.Dir = root
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return c.Run()
}

// quebraEm reflui o texto numa largura, prefixando as linhas seguintes. Um prompt de 400
// caracteres numa linha só é impossível de ler no terminal e feio de copiar.
func quebraEm(s string, largura int, prefixo string) string {
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

// temTTY diz se há terminal interativo na entrada. É a evidência mais fraca de quem
// opera (um pipe qualquer produz o mesmo resultado), por isso `DetectaOperador` só
// recorre a ela depois de procurar as variáveis que um agente declara.
func temTTY() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}
