package initx

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
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

// ProtecaoDeBranch é o que o modo `github` exige da `main`: nada entra sem PR.
//
// É o que torna o ciclo de revisão possível. A §7.9 do BOOTSTRAP diz que o agente sobe o
// código e ABRE O PR, e o card vai para `ready-to-review` — sem PR não há o que revisar,
// e o estado `in-review` fica sem objeto.
//
// E é o que o pipeline de identificação pressupõe: ele dispara na abertura do PR, porque
// o push na main acontece DEPOIS do merge, quando o trabalho já terminou. Push direto na
// main pula o card, pula a revisão, e pula o pipeline.
type ProtecaoDeBranch struct {
	// ExigePR — nada entra na main sem pull request.
	ExigePR bool
	// RevisoesNecessarias — quantas aprovações. Zero é legítimo num time de uma pessoa
	// com agentes: o PR existe para o card ter objeto e para o histórico ficar legível,
	// e exigir aprovação de outra conta travaria o fluxo inteiro.
	RevisoesNecessarias int
}

// ProtecaoExigida é o mínimo que o fluxo pressupõe.
var ProtecaoExigida = ProtecaoDeBranch{ExigePR: true, RevisoesNecessarias: 0}

// DirWorkflows é onde os pipelines moram no projeto.
const DirWorkflows = ".github/workflows"

// EstadosDoTrabalho são as LABELS que carregam o estado de um card, na ordem do fluxo.
//
// O estado é uma LABEL, e não a coluna do Project (ver BOOTSTRAP.md §7.13). A escolha
// anterior foi a coluna, e ela cobrava um preço que só apareceu no uso: escrever num
// Project de organização exige um PAT com escopo `project`, que o `GITHUB_TOKEN` da
// Action não tem — então todo projeto que adotasse o fluxo precisaria criar e manter um
// token pessoal antes de o primeiro card se mover. Atrito de adoção por uma decisão de
// visualização.
//
// Com label, o `GITHUB_TOKEN` basta e nada precisa ser configurado. E o board não some:
// o GitHub Projects tem automação nativa que move o card quando a label muda — a
// sincronia acontece do lado deles, e só UM lado escreve (nós na label, eles no board),
// que era a preocupação original de ter estado em dois lugares.
//
// O par `ready-to-x` / `in-x` é o que torna a fila legível: um diz "disponível para
// alguém pegar", o outro "alguém está fazendo".
var EstadosDoTrabalho = []string{
	"anchors:to-do",
	"anchors:in-progress",
	"anchors:ready-to-review",
	"anchors:in-review",
	"anchors:ready-to-test",
	"anchors:in-test",
	"anchors:ready-to-release",
	"anchors:production",
}

// ColunasDoBoard são os nomes das colunas do Project que ESPELHAM os estados acima.
// O board é opcional: quem o quiser cria as colunas com estes nomes e liga a automação
// nativa do Projects (label adicionada → move para a coluna). Quem não quiser trabalha
// só com issues, e o fluxo funciona igual.
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

// EstadoFinalDoAnchors é o último estado que o Anchors ESCREVE.
const EstadoFinalDoAnchors = "anchors:ready-to-test"

// EstadosRecicláveis são os únicos estados de onde o `stale` tira um card.
//
// A distinção é entre trabalho TRAVADO e trabalho ESPERANDO. Um card em `in-progress`
// sem sinal de vida é um agente que morreu — reciclar devolve o trabalho à fila. Um card
// em `ready-to-review` está parado POR DEFINIÇÃO: ele terminou e aguarda um revisor
// humano, e ninguém comenta enquanto isso. Reciclá-lo joga trabalho pronto de volta em
// `to-do`, como se nunca tivesse sido feito — e o agente seguinte o refaz.
//
// Só entram aqui os estados de trabalho ATIVO, aqueles em que alguém deveria estar
// mexendo agora. Os `ready-to-*` são filas de espera, e esperar não é estar travado.
var EstadosRecicláveis = []string{
	"anchors:in-progress",
	"anchors:in-review",
	"anchors:in-test",
}

// EstadosDisponiveis são os estados de onde um agente TIRA trabalho, em ORDEM DE
// PRIORIDADE: da direita para a esquerda do fluxo. O trabalho mais ADIANTADO vem
// primeiro — terminar o que está quase pronto antes de começar coisa nova é o que impede
// o board de encher de trabalho pela metade.
var EstadosDisponiveis = []string{
	"anchors:ready-to-review",
	"anchors:to-do",
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

// MarcadorDeTemplate identifica um pipeline que ainda é o do Anchors — não editado pelo
// time. É o que separa "customizado" de "desatualizado".
//
// Sem essa distinção, `--fix` só sabia que o arquivo EXISTIA, e tratava as duas
// situações igual: preservava as duas. O efeito era que uma correção no template nunca
// alcançava quem já tinha instalado — o projeto ficava com o pipeline defeituoso para
// sempre, e nada avisava. Foi assim que um `stale` que reciclava trabalho pronto
// continuou rodando depois de corrigido na fonte.
const MarcadorDeTemplate = "# anchors:template"

// ÉTemplateIntacto diz se o pipeline instalado ainda é o do Anchors.
//
// A prova é o marcador que o template carrega. Editar o arquivo pede que se remova a
// linha — e quem edita sem removê-la perde a customização na próxima atualização. É uma
// troca deliberada: o marcador é explícito no arquivo e diz exatamente isso, e a
// alternativa (hash do conteúdo) marcaria como "customizado" qualquer arquivo que o
// próprio Anchors gerou com um branch de integração diferente.
func ÉTemplateIntacto(root, arquivo string) bool {
	b, err := os.ReadFile(filepath.Join(root, DirWorkflows, arquivo))
	if err != nil {
		return false
	}
	return strings.Contains(string(b), MarcadorDeTemplate)
}

// WorkflowsDesatualizados diz quais pipelines instalados são do Anchors (marcador
// intacto) e diferem do template atual — os que `--fix` pode e deve atualizar.
//
// Um pipeline SEM o marcador não entra aqui mesmo que difira: é do time, e a diferença é
// a customização dele.
func WorkflowsDesatualizados(root string, cfg *config.Config) []Workflow {
	var velhos []Workflow
	for _, w := range WorkflowsDoFluxo {
		if !ÉTemplateIntacto(root, w.Arquivo) {
			continue
		}
		atual, err := os.ReadFile(filepath.Join(root, DirWorkflows, w.Arquivo))
		if err != nil {
			continue
		}
		esperado, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			continue
		}
		esperado = aplicaBranchDeIntegracao(esperado, cfg.Workflow.BranchDeIntegracao())
		if !bytes.Equal(atual, esperado) {
			velhos = append(velhos, w)
		}
	}
	return velhos
}

// SemeiaWorkflows escreve em `.github/workflows/` os pipelines que faltam, e ATUALIZA os
// que ainda são do Anchors e ficaram para trás. Devolve os nomes escritos.
//
// Nunca sobrescreve um pipeline que o time editou — outro ritmo de stale, uma permissão a
// mais, um passo de build próprio é trabalho deliberado, e reescrevê-lo apagaria a
// customização sem avisar. É a régua que o `install-hooks` já usa com um pre-commit
// alheio. A diferença é que agora "editado pelo time" é uma pergunta que se responde
// (o marcador), e não uma suposição a partir da mera existência do arquivo.
// Recebe o CONFIG, e não o branch já extraído: as regras de branch moram no anchors.yaml,
// e quem precisa delas as lê de lá. Passar o valor pronto espalharia a decisão por cada
// chamador, e bastaria um deles ler de outro lugar para o projeto ter dois fluxos.
func SemeiaWorkflows(root string, cfg *config.Config) ([]string, error) {
	dir := filepath.Join(root, DirWorkflows)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("criar %s: %w", DirWorkflows, err)
	}
	var escritos []string
	for _, w := range WorkflowsDoFluxo {
		dest := filepath.Join(dir, w.Arquivo)
		if _, err := os.Stat(dest); err == nil && !ÉTemplateIntacto(root, w.Arquivo) {
			continue // existe e o marcador saiu: é do time, não nosso
		}
		conteudo, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			return escritos, fmt.Errorf("ler o template %s: %w", w.Arquivo, err)
		}
		conteudo = aplicaBranchDeIntegracao(conteudo, cfg.Workflow.BranchDeIntegracao())
		if err := os.WriteFile(dest, conteudo, 0o644); err != nil {
			return escritos, fmt.Errorf("escrever %s: %w", dest, err)
		}
		escritos = append(escritos, w.Arquivo)
	}
	sort.Strings(escritos)
	return escritos, nil
}

// aplicaBranchDeIntegracao troca o branch cravado no template pelo que o projeto
// declarou. A linha alvo é marcada com `# anchors:integration-branch` — um marcador, e
// não uma busca por "main", porque "main" aparece em comentário e em outros contextos, e
// substituir a ocorrência errada quebraria o pipeline de um jeito difícil de ver.
func aplicaBranchDeIntegracao(conteudo []byte, branch string) []byte {
	if branch == "" || branch == "main" {
		return conteudo
	}
	linhas := strings.Split(string(conteudo), "\n")
	for i, l := range linhas {
		if strings.Contains(l, "# anchors:integration-branch") {
			indent := l[:len(l)-len(strings.TrimLeft(l, " "))]
			linhas[i] = indent + "branches: [" + branch + "] # anchors:integration-branch"
		}
	}
	return []byte(strings.Join(linhas, "\n"))
}
