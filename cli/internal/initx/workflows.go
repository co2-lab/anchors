package initx

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// workflowsFS carrega os pipelines que o modo `github` do fluxo de trabalho pressupõe.
//
// Viajam no binário e são COPIADOS para `.github/workflows/` — a mesma régua dos packs.
// A diferença importa: no projeto, o pipeline é material, versionado e revisável num
// diff. Um time que precise ajustar o ritmo do `stale`, trocar a label ou mudar a
// permissão edita o arquivo; um pipeline que vivesse só no binário obrigaria a esperar um
// release para mudar o que é decisão do time.
//
//go:embed workflows
var workflowsFS embed.FS

// Workflow é um pipeline do fluxo GitHub, com o que o Anchors precisa saber para
// verificar (doctor) e semear (doctor --fix).
type Workflow struct {
	// Arquivo é o nome em `.github/workflows/`.
	Arquivo string
	// Papel é o que ele faz, para a mensagem do doctor dizer o que fica sem acontecer.
	Papel string
	// ExigeSerial diz se este pipeline PRECISA de `concurrency` sem cancelamento. Para
	// os que atribuem ou criam cards, a serialização não é otimização: é o que impede
	// duas execuções de atribuírem o mesmo card ou criarem o card duas vezes.
	ExigeSerial bool
}

// WorkflowsDoFluxo são os pipelines que o modo `github` pressupõe. A lista é a fonte da
// verdade tanto do que o doctor confere quanto do que o `--fix` semeia — uma lista só,
// para que verificar e consertar nunca discordem sobre o que deveria existir.
var WorkflowsDoFluxo = []Workflow{
	{
		Arquivo:     "anchors-identify.yml",
		Papel:       "cria o card de todo artefato que chegou ao repositório sem um",
		ExigeSerial: true,
	},
	{
		Arquivo:     "anchors-claim.yml",
		Papel:       "atribui trabalho aos agentes (é o que elimina a corrida por um card)",
		ExigeSerial: true,
	},
	{
		Arquivo:     "anchors-stale.yml",
		Papel:       "libera cards cujo dono sumiu, preservando o histórico",
		ExigeSerial: true,
	},
}

// DirWorkflows é onde os pipelines moram no projeto.
const DirWorkflows = ".github/workflows"

// ColunasDoBoard são os valores do campo `Status` do GitHub Project, na ordem do fluxo.
// O estado do trabalho é a COLUNA — não uma label (ver BOOTSTRAP.md §7.13): estado em
// dois lugares dessincroniza, e um board que mente sobre onde o trabalho está vira
// trabalho duplicado, porque é olhando o board que os agentes escolhem o que pegar.
//
// O par `READY TO X` / `IN X` é o que torna a fila legível: um diz "disponível para
// alguém pegar", o outro "alguém está fazendo".
var ColunasDoBoard = []string{
	"TO DO",
	"IN PROGRESS",
	"READY TO REVIEW",
	"IN REVIEW",
	"READY TO TEST",
	"IN TEST",
	"READY TO RELEASE",
	"PRODUCTION",
}

// ColunasDisponiveis são as colunas de onde um agente TIRA trabalho — as `READY TO ...`
// e o `TO DO` inicial, em ORDEM DE PRIORIDADE: da direita para a esquerda do board.
//
// A regra é essa e vale para todo o fluxo: o trabalho mais ADIANTADO vem primeiro.
// Terminar o que está quase pronto antes de começar coisa nova é o que impede o board de
// encher de trabalho pela metade — e trabalho pela metade não entrega nada, ocupa revisor,
// e envelhece até o contexto de quem o escreveu se perder.
//
// Daí a ordem invertida em relação a ColunasDoBoard: `READY TO REVIEW` (mais à direita)
// antes de `TO DO` (mais à esquerda).
var ColunasDisponiveis = []string{
	"READY TO REVIEW",
	"TO DO",
}

// ColunaFinalDoAnchors é a última coluna que o Anchors ESCREVE. Da seguinte em diante
// (`IN TEST`, `READY TO RELEASE`, `PRODUCTION`) quem move são os pipelines de entrega do
// projeto — cada time tem o seu, e o Anchors não tem o que dizer sobre quando um teste de
// aceitação passou ou um deploy aconteceu. Ele continua LENDO essas colunas (é o que o
// `anchors status` mostra), mas não escreve nelas.
const ColunaFinalDoAnchors = "READY TO TEST"

// ColunasQueOAnchorsEscreve são as que os pipelines do Anchors movem. Serve ao doctor:
// uma coluna ausente aqui quebra o fluxo; uma ausente depois é problema do time.
func ColunasQueOAnchorsEscreve() []string {
	for i, c := range ColunasDoBoard {
		if c == ColunaFinalDoAnchors {
			return ColunasDoBoard[:i+1]
		}
	}
	return ColunasDoBoard
}

// FaltaWorkflow diz quais dos pipelines do fluxo não existem no projeto. Só presença —
// a coerência do conteúdo é outra pergunta, respondida por `SemConcurrency`.
func FaltaWorkflow(root string) []Workflow {
	var faltam []Workflow
	for _, w := range WorkflowsDoFluxo {
		if _, err := os.Stat(filepath.Join(root, DirWorkflows, w.Arquivo)); err != nil {
			faltam = append(faltam, w)
		}
	}
	return faltam
}

// SemConcurrency diz quais pipelines existentes estão sem a serialização que o desenho
// pressupõe. Um `anchors-claim.yml` sem `concurrency` roda em paralelo consigo mesmo e
// atribui o mesmo card a dois agentes — o problema exato que o pipeline existe para
// eliminar, de volta e agora invisível, porque o arquivo ESTÁ lá.
//
// Verifica também `cancel-in-progress: false`: cancelar a execução em curso poderia
// matá-la depois de ela já ter criado metade dos cards.
func SemConcurrency(root string) []Workflow {
	var quebrados []Workflow
	for _, w := range WorkflowsDoFluxo {
		if !w.ExigeSerial {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, DirWorkflows, w.Arquivo))
		if err != nil {
			continue // ausente é outro achado (FaltaWorkflow); não contar duas vezes
		}
		texto := string(b)
		if !strings.Contains(texto, "concurrency:") || !strings.Contains(texto, "cancel-in-progress: false") {
			quebrados = append(quebrados, w)
		}
	}
	return quebrados
}

// SemeiaWorkflows copia para `.github/workflows/` os pipelines que faltam. Devolve os
// nomes escritos.
//
// NUNCA sobrescreve um arquivo existente. Um pipeline que o time editou — outro ritmo de
// stale, uma permissão a mais, um passo de build próprio — é trabalho deliberado, e
// reescrevê-lo pelo padrão apagaria a customização sem avisar. É a régua que o
// `install-hooks` já usa com um pre-commit alheio.
func SemeiaWorkflows(root string) ([]string, error) {
	dir := filepath.Join(root, DirWorkflows)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("criar %s: %w", DirWorkflows, err)
	}
	var escritos []string
	for _, w := range WorkflowsDoFluxo {
		dest := filepath.Join(dir, w.Arquivo)
		if _, err := os.Stat(dest); err == nil {
			continue // já existe: é do time, não nosso
		}
		conteudo, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			return escritos, fmt.Errorf("ler o template %s: %w", w.Arquivo, err)
		}
		if err := os.WriteFile(dest, conteudo, 0o644); err != nil {
			return escritos, fmt.Errorf("escrever %s: %w", dest, err)
		}
		escritos = append(escritos, w.Arquivo)
	}
	sort.Strings(escritos)
	return escritos, nil
}
