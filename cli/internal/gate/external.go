package gate

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/co2-lab/anchors/internal/mapx"
)

// runExternal invoca um comando externo (jest, eslint, tsc…). O CLI NÃO reimplementa
// a ferramenta — roda e lê o exit code (D5: reimplementamos parsing, não ferramentas
// de teste/lint). O comando do gate vem do anchors.yaml (config do projeto) e pode
// ser composto (ex.: "cd apps/mobile && jest \"$1\""), por isso roda via `sh -c`.
//
// SEGURANÇA: o caminho do alvo NUNCA é interpolado na string do shell — ele é
// passado como ARGUMENTO POSICIONAL (`$1`), que o shell não reinterpreta. Assim um
// nome de arquivo com metacaracteres (foo;rm.tsx) é dado puro, não injeção. O
// comando deve referenciar o alvo como "$1" (com aspas). `{{file}}` é aceito por
// conveniência e reescrito para "$1" antes de rodar.
//
// exit 0 = pass; qualquer outro = fail (com stderr/stdout no detalhe).
func runExternal(command string, n mapx.Node, root string) (Verdict, string) {
	return RunExternalArgs(command, []string{n.ID}, root)
}

// limiteArgv é quanto os alvos podem ocupar, em bytes, numa ÚNICA chamada do shell.
//
// O Windows monta a linha de comando como uma string só e o CreateProcess a corta em
// 32767 chars: um `scope: batch` sobre um projeto inteiro (1484 alvos ≈ 85 KB) falhava
// ANTES de rodar, com "o nome do arquivo ou a extensão é muito grande". Como o exec
// nem chegava a executar, a saída voltava vazia e o gate reprovava MUDO — indistinguível
// de violação real, e verde no macOS, onde o ARG_MAX é de centenas de KB.
//
// O teto do Windows é MUITO menor que os 32767 porque o orçamento não é gasto só aqui:
// o gate costuma empilhar wrappers (sh → yarn → node → binário), e cada camada reexpande
// os caminhos — o sh do MSYS2 converte cada relativo em absoluto ao chamar um .exe
// nativo. Medido no app de referência com o gate eslint: 8379 bytes de alvos ainda passam, 11389 já
// estouram DENTRO do script, longe da vista do Anchors. 6000 deixa a margem.
//
// O teto do Unix fica alto de propósito: onde o lote já cabia, continua sendo UMA
// execução, e o comportamento não muda.
//
// ANCHORS_ARGV_MAX ajusta o teto para quem empilha mais (ou menos) wrapper que isso —
// o número certo é propriedade dos gates do projeto, não do Anchors.
func limiteArgv() int {
	if v := os.Getenv("ANCHORS_ARGV_MAX"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	if runtime.GOOS == "windows" {
		return 6000
	}
	return 100000
}

// RunExternalArgs é o motor: roda o comando com N alvos como argumentos posicionais
// ($1, $2, … e "$@"). Um alvo é o caso por-nó; vários, o `scope: batch`; nenhum, o
// `scope: project` (a ferramenta olha o projeto inteiro e não recebe alvo).
//
// A garantia de segurança é a mesma qualquer que seja a quantidade: os caminhos vão
// em argv, nunca interpolados no script. `{{files}}` é reescrito para "$@" pela mesma
// razão que `{{file}}` vira "$1" — conveniência, sem abrir caminho para injeção.
//
// Quando os alvos não cabem numa linha de comando, o gate roda EM LOTES e os vereditos
// se somam: reprova se qualquer lote reprovar, e o laudo junta a saída de todos. Isso
// pressupõe um gate que julga arquivo a arquivo — o que todo gate de `scope: batch` é,
// por construção (recebe um recorte arbitrário do projeto, nunca o conjunto fechado).
// Um gate que precise ver todos os alvos DE UMA VEZ (achar duplicata entre arquivos,
// por exemplo) deve declarar `scope: project` e varrer sozinho, não depender do lote.
func RunExternalArgs(command string, targets []string, root string) (Verdict, string) {
	script := strings.ReplaceAll(command, "{{file}}", `"$1"`)
	script = strings.ReplaceAll(script, "{{files}}", `"$@"`)

	var falhas []string
	for _, lote := range fatiarAlvos(targets, limiteArgv()-len(script)) {
		out, err := execShell(script, lote, root)
		if err == nil {
			continue
		}
		if detalhe := strings.TrimSpace(out); detalhe != "" {
			falhas = append(falhas, detalhe)
			continue
		}
		// Reprovou sem dizer nada. Acontece quando o próprio exec falha (não há `sh`
		// no PATH, a linha de comando estourou) — sem esta linha o laudo fica vazio e
		// o operador lê "violação de código" onde há problema de ambiente.
		falhas = append(falhas, "gate não produziu saída: "+err.Error())
	}
	if len(falhas) == 0 {
		return Pass, ""
	}

	detail := strings.Join(falhas, "\n")
	// Ferramenta de projeto (tsc/eslint) reporta MUITOS erros de uma vez. O corte em
	// 500 do caso por-nó deixaria o laudo inútil justamente onde ele mais importa —
	// some a lista de arquivo:linha, que é o conteúdo que resolve o problema.
	limite := 500
	if len(targets) != 1 {
		limite = 4000
	}
	if len(detail) > limite {
		detail = detail[:limite] + "\n… (saída truncada)"
	}
	return Fail, detail
}

// execShell roda UM lote. `sh -c '<script>' sh <alvo...>` → os alvos viram $1..$N
// (e "$@") dentro do script, sem passar pela tokenização do shell (não é injeção).
func execShell(script string, targets []string, root string) (string, error) {
	args := append([]string{"-c", script, "sh"}, targets...)
	cmd := exec.Command("sh", args...) //nolint:gosec // comando é da config do projeto; os alvos vão como argv, não interpolados
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// fatiarAlvos divide os alvos em lotes que cabem numa linha de comando.
//
// Sem alvos devolve UM lote vazio, não lote nenhum: `scope: project` é justamente a
// execução sem alvo, e devolver nada faria o gate não rodar (passando por omissão).
// Um alvo sozinho maior que o teto vai só no seu lote — não há como parti-lo, e é
// melhor deixar o SO recusar com a mensagem dele do que silenciar o alvo.
func fatiarAlvos(targets []string, teto int) [][]string {
	if len(targets) == 0 {
		return [][]string{nil}
	}
	if teto < 1 {
		teto = 1
	}

	var lotes [][]string
	var atual []string
	tam := 0
	for _, t := range targets {
		custo := len(t) + 1 // o separador que o SO conta entre um argumento e o próximo
		if len(atual) > 0 && tam+custo > teto {
			lotes = append(lotes, atual)
			atual, tam = nil, 0
		}
		atual = append(atual, t)
		tam += custo
	}
	return append(lotes, atual)
}
