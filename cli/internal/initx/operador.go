package initx

import (
	"os"
	"path/filepath"
	"sort"
)

// Operador é QUEM está do outro lado do `anchors init`. A distinção existe porque a
// fase DESCOBRIR (ver `anchors guide project`) precisa de uma IA conduzindo a
// entrevista, e o que o Anchors deve DIZER muda conforme quem lê a saída:
//
//   - para uma PESSOA, a saída é instrução: aqui está o que falta, e aqui está como
//     começar (inclusive abrindo a IA com o prompt pronto);
//   - para uma IA, a saída é a própria ordem de serviço: leia o guide e conduza a
//     entrevista com o usuário, nesta conversa.
//
// Escrever uma coisa só para os dois falha nos dois: a pessoa recebe instruções que
// não sabe executar, e a IA recebe um convite quando precisava de uma tarefa.
type Operador int

const (
	// OperadorHumano — uma pessoa num terminal.
	OperadorHumano Operador = iota
	// OperadorIA — uma IA operando o CLI (o usuário pediu a ela para iniciar o projeto).
	OperadorIA
)

// agentesConhecidos mapeia variável de ambiente → nome da ferramenta, para o caso em
// que o Anchors precisa NOMEAR quem está operando. A lista é de reconhecimento, não de
// suporte: uma IA fora dela ainda pode ser detectada por outros sinais, e nenhuma
// funcionalidade depende de estar aqui.
var agentesConhecidos = map[string]string{
	"CLAUDE_CODE_ENTRYPOINT": "Claude Code",
	"CLAUDECODE":             "Claude Code",
	"CURSOR_TRACE_ID":        "Cursor",
	"AIDER_MODEL":            "Aider",
	"GEMINI_CLI":             "Gemini CLI",
	"CODEX_SANDBOX":          "Codex",
}

// DetectaOperador diz quem está operando. `temTTY` é injetado (não consultado aqui)
// para manter a função pura e testável: a suíte precisa cobrir os dois lados sem
// depender de como o teste foi invocado.
//
// A ordem das evidências importa. Uma variável de agente é declaração EXPLÍCITA de
// quem está rodando, e vale mais que a ausência de TTY — que é só um indício, e um
// indício que um pipe qualquer também produz.
func DetectaOperador(temTTY bool, env func(string) string) Operador {
	if env == nil {
		env = os.Getenv
	}
	if nomeDoAgente(env) != "" {
		return OperadorIA
	}
	// `AI_AGENT` é genérica o bastante para valer como sinal sem estar na lista de
	// nomes: quem a define está declarando que não é uma pessoa digitando.
	if env("AI_AGENT") != "" {
		return OperadorIA
	}
	// Sem TTY sobra pipe, CI ou agente. Nenhum deles é uma pessoa lendo a saída e
	// digitando a resposta, e o texto para IA é o mais útil dos dois nesse caso: ele
	// descreve o trabalho, e é acionável por quem quer que esteja lendo o log.
	if !temTTY {
		return OperadorIA
	}
	return OperadorHumano
}

// nomeDoAgente devolve o nome da ferramenta de IA detectada, ou "". Percorre em ordem
// estável para que a mensagem não mude entre execuções idênticas.
func nomeDoAgente(env func(string) string) string {
	if env == nil {
		env = os.Getenv
	}
	chaves := make([]string, 0, len(agentesConhecidos))
	for k := range agentesConhecidos {
		chaves = append(chaves, k)
	}
	sort.Strings(chaves)
	for _, k := range chaves {
		if env(k) != "" {
			return agentesConhecidos[k]
		}
	}
	return ""
}

// NomeDoAgente é a versão exportada, para a mensagem poder dizer "abrir o Claude Code"
// em vez de "abrir sua IA".
func NomeDoAgente(env func(string) string) string { return nomeDoAgente(env) }

// PrecisaDescobrir diz se a fase DESCOBRIR ainda não aconteceu neste projeto: não há
// PROJECT.md, e não há código de onde o `init` pudesse inferir a Estrutura.
//
// As duas condições juntas, nunca uma só. Um projeto com código e sem PROJECT.md não
// precisa da fase (o `init` infere do disco); e um PROJECT.md presente é a prova de que
// ela já rodou, mesmo num diretório ainda vazio de código — que é exatamente o estado
// em que ela DEVE deixar o projeto.
func PrecisaDescobrir(root string, p *Proposal) bool {
	if p != nil && (len(p.CodeDirs) > 0 || p.HasSpecMD || p.HasFeature || p.HasTest) {
		return false
	}
	return !TemProjectMD(root)
}

// TemProjectMD diz se o PROJECT.md já existe na raiz. Aceita as duas grafias que
// aparecem na prática — o guide escreve `PROJECT.md`, mas um projeto que já usava
// `project.md` não deve ser mandado refazer a entrevista.
func TemProjectMD(root string) bool {
	for _, nome := range []string{"PROJECT.md", "project.md", "Project.md"} {
		if _, err := os.Stat(filepath.Join(root, nome)); err == nil {
			return true
		}
	}
	return false
}

// PromptDescobrir é o texto que inicia a fase DESCOBRIR numa IA. É o que o Anchors
// oferece copiar/executar para o usuário que está sozinho no terminal, e é o mesmo
// trabalho que ele DESCREVE quando quem opera já é uma IA.
//
// Curto de propósito: manda ler a régua em vez de reproduzi-la. O `anchors guide
// project` é a fonte da verdade da entrevista, e um prompt que a resumisse
// desatualizaria em silêncio na primeira vez que o guide mudasse.
const PromptDescobrir = `Rode "anchors guide project" e siga essa régua à risca: ` +
	`conduza comigo, aqui nesta conversa, a entrevista de 5 etapas da fase DESCOBRIR ` +
	`(propósito e forma → linguagem → arquitetura e paradigma → estrutura macro e ` +
	`convenções → ferramental e formatação). Uma pergunta por vez, esperando minha ` +
	`resposta antes da próxima. No fim, faça a revisão de inconsistências e escreva ` +
	`PROJECT.md e INSIGHTS.md na raiz. Depois disso rodamos "anchors init".`

// ComandoParaAbrirIA devolve o argv que abre a IA detectada já com o prompt, ou nil
// quando não há como saber qual abrir. Só ferramentas cuja invocação por linha de comando
// é estável entram aqui: oferecer um comando que não existe é pior do que não oferecer
// nada.
//
// Devolve ARGV, não uma linha de shell. O prompt tem aspas, parênteses e setas — passá-lo
// por `sh -c` exigiria escapá-lo certo, e um escape errado ou vira comando torto ou
// executa o que não devia. Com argv, o texto é um argumento e ponto.
func ComandoParaAbrirIA(env func(string) string) []string {
	switch nomeDoAgente(env) {
	case "Claude Code":
		return []string{"claude", PromptDescobrir}
	case "Gemini CLI":
		return []string{"gemini", PromptDescobrir}
	case "Aider":
		return []string{"aider", "--message", PromptDescobrir}
	default:
		return nil
	}
}
