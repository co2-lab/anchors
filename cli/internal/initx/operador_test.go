package initx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// envFalso simula variáveis de ambiente sem tocar as reais — a detecção precisa ser
// testável nos dois lados, e a máquina de quem roda a suíte tem as suas.
func envFalso(pares map[string]string) func(string) string {
	return func(k string) string { return pares[k] }
}

// A distinção que estrutura o bootstrap: o que o Anchors DIZ muda conforme quem lê. Uma
// variável de agente é declaração explícita, e vale mais que a ausência de TTY — que um
// pipe qualquer também produz.
func TestVarDeAgenteVenceAAusenciaDeTTY(t *testing.T) {
	comAgente := envFalso(map[string]string{"CLAUDE_CODE_ENTRYPOINT": "cli"})

	// Mesmo COM tty, a variável decide: uma IA pode rodar num PTY.
	if o := DetectaOperador(true, comAgente); o != OperadorIA {
		t.Error("com variável de agente declarada, é IA — mesmo havendo TTY")
	}
	// E sem nada, um terminal é uma pessoa.
	if o := DetectaOperador(true, envFalso(nil)); o != OperadorHumano {
		t.Error("TTY e nenhuma variável de agente: é uma pessoa")
	}
	// Sem TTY sobra pipe, CI ou agente — nenhum é alguém digitando.
	if o := DetectaOperador(false, envFalso(nil)); o != OperadorIA {
		t.Error("sem TTY não há pessoa lendo a saída e respondendo")
	}
}

// Nomear a ferramenta é o que permite dizer "abro o Claude Code?" em vez de "abra sua
// IA" — e é a diferença entre uma oferta acionável e um conselho.
func TestNomeiaAFerramentaDetectada(t *testing.T) {
	if n := NomeDoAgente(envFalso(map[string]string{"CURSOR_TRACE_ID": "x"})); n != "Cursor" {
		t.Errorf("esperava Cursor, veio %q", n)
	}
	if n := NomeDoAgente(envFalso(nil)); n != "" {
		t.Errorf("sem variável conhecida não há nome a dizer, veio %q", n)
	}
}

// O comando é ARGV, não linha de shell: o prompt tem aspas, parênteses e setas, e passá-lo
// por `sh -c` faria o escape ser a única coisa entre o texto e o interpretador.
func TestComandoDeAberturaEhArgvNaoShell(t *testing.T) {
	argv := ComandoParaAbrirIA(envFalso(map[string]string{"CLAUDE_CODE_ENTRYPOINT": "cli"}))

	if len(argv) < 2 {
		t.Fatalf("esperava argv com programa e prompt, veio %v", argv)
	}
	if argv[0] != "claude" {
		t.Errorf("programa errado: %q", argv[0])
	}
	// O prompt inteiro num argumento só — sem aspas de shell coladas nele.
	if strings.HasPrefix(argv[1], "'") || strings.HasSuffix(argv[1], "'") {
		t.Errorf("o prompt veio com aspas de shell: %q", argv[1])
	}
	// Ferramenta desconhecida: nada a oferecer. Um comando que não existe é pior que
	// nenhuma oferta.
	if argv := ComandoParaAbrirIA(envFalso(nil)); argv != nil {
		t.Errorf("sem ferramenta conhecida não há comando a oferecer: %v", argv)
	}
}

// O prompt manda LER o guide, não o resume: `anchors guide project` é a fonte da verdade
// da entrevista, e um resumo desatualizaria em silêncio na primeira vez que ela mudasse.
func TestPromptMandaLerAReguaEmVezDeResumila(t *testing.T) {
	if !strings.Contains(PromptDescobrir, "anchors guide project") {
		t.Error("o prompt tem de mandar rodar o guide — ele é a régua")
	}
	if !strings.Contains(PromptDescobrir, "anchors init") {
		t.Error("o prompt tem de fechar o ciclo apontando de volta para o init")
	}
}

// A fase é necessária só quando não há NADA de onde inferir. Um PROJECT.md é a prova de
// que ela rodou; código no disco dispensa (o init infere dali).
func TestPrecisaDescobrirSoNoVazio(t *testing.T) {
	dir := t.TempDir()

	if !PrecisaDescobrir(dir, &Proposal{}) {
		t.Error("sem PROJECT.md e sem código, a fase não aconteceu")
	}
	if !PrecisaDescobrir(dir, nil) {
		t.Error("proposta nula não deve mascarar o vazio")
	}
	// Código no disco: há o que inferir.
	if PrecisaDescobrir(dir, &Proposal{CodeDirs: []string{"src"}}) {
		t.Error("com código, o init infere do disco — a fase não é necessária")
	}
	// PROJECT.md presente: a fase já rodou, mesmo sem código ainda (que é exatamente o
	// estado em que ela DEVE deixar o projeto).
	if err := os.WriteFile(filepath.Join(dir, "PROJECT.md"), []byte("# P\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if PrecisaDescobrir(dir, &Proposal{}) {
		t.Error("PROJECT.md é a prova de que a entrevista aconteceu")
	}
}
