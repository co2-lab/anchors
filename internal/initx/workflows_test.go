// @anchors
//   ref: FLWRF

package initx

import (
	"embed"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/testkit"
	"gopkg.in/yaml.v3"

	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/config"
)

// A serialização não é detalhe de performance: é o mecanismo que impede duas execuções
// de atribuírem o mesmo card ou de criarem o card duas vezes. Um template que a perdesse
// devolveria a corrida — agora invisível, porque o arquivo ESTÁ lá.
func TestTemplatesSeriaisTrazemConcurrency(t *testing.T) {
	for _, w := range WorkflowsDoFluxo {
		if !w.ExigeSerial {
			continue
		}
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatalf("%s: %v", w.Arquivo, err)
		}
		texto := string(b)
		if !strings.Contains(texto, "concurrency:") {
			t.Errorf("%s exige serialização e não declara `concurrency:`", w.Arquivo)
		}
		// Cancelar a execução em curso pode matá-la DEPOIS de ela ter criado metade dos
		// cards, e o push seguinte não saberia o que ficou pela metade.
		if !strings.Contains(texto, "cancel-in-progress: false") {
			t.Errorf("%s cancela a execução em curso — pode interromper no meio do trabalho", w.Arquivo)
		}
	}
}

// O pipeline de claim é o que faz a concorrência deixar de existir. Se ele passasse a
// disparar em `push`, voltaria a rodar concorrente com outras causas.
func TestClaimSoRodaSobDemanda(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, "workflow_dispatch:") {
		t.Error("o claim é PEDIDO pelo agente — precisa de workflow_dispatch")
	}
	if strings.Contains(texto, "\n  push:") {
		t.Error("claim disparado por push atribuiria trabalho sem ninguém ter pedido")
	}
	// A identidade do agente é máquina+sessão, não o usuário do GitHub: dois agentes na
	// mesma máquina têm o mesmo login.
	if !strings.Contains(texto, "anchors-owner:") {
		t.Error("o dono-agente é registrado por comentário `anchors-owner:`, não pelo assignee")
	}
}

// O pipeline de identificação usa `--json` de propósito: no TSV as pastas vêm juntas por
// ", ", que se confunde com uma vírgula dentro de um caminho.
func TestIdentifyUsaSaidaEstruturada(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-identify.yml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "anchors code list --json") {
		t.Error("o identify deve consumir `anchors code list --json`, não o TSV de leitura humana")
	}
	// Lê issues abertas E fechadas: um artefato cujo trabalho terminou não pode ganhar
	// card novo a cada push.
	if !strings.Contains(string(b), "--state all") {
		t.Error("consultar só issues abertas recriaria o card de todo trabalho já concluído")
	}
}

// O título do card vem do TÍTULO do artefato, não do caminho. `anchors code list`
// devolve a PASTA da unidade — com 16 planos em `plans/`, todo card sairia como
// "Implementar plans", indistinguível dos outros. E o caminho já está no corpo.
func TestIdentifyUsaOTituloDoArtefato(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-identify.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, `--title "[$code] $titulo"`) {
		t.Error("o título do card tem de nomear o trabalho")
	}
	// O verbo importa: o card é uma TAREFA, não um rótulo do arquivo.
	if !strings.Contains(texto, `titulo="Implementar ${kind:-artefato} — $titulo"`) {
		t.Error("o título tem de dizer o que fazer, com que tipo de artefato, e sobre o quê")
	}
	// E precisa de reserva: um arquivo de código não tem `# título`.
	if !strings.Contains(texto, `titulo="Implementar ${kind:-artefato} em $onde"`) {
		t.Error("sem título no artefato, o card precisa de um nome mesmo assim")
	}
}

// O pipeline é SERIALIZADO, então cada segundo dele é fila para todos. Reconstruir o
// mapa SEMPRE custaria 12,7s num projeto de 3.587 nós (medido), e o mapa já chega pronto
// no push — é versionado, e o pre-commit o mantém em dia.
//
// A regra: verifica o barato (só os arquivos do change estão no mapa?) e reconstrói SÓ
// se estiver defasado. Confiar cegamente no mapa commitado seria pior que o custo: o
// pipeline decidiria sobre uma foto velha e deixaria de criar os cards do que acabou de
// chegar, em silêncio.
func TestIdentifyVerificaAntesDeReconstruirOMapa(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-identify.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	// O binário vem da RELEASE, não de `go install`: o módulo declara o caminho da raiz
	// mas vive em `cli/`, então nenhum caminho de install resolve — e baixar o binário é
	// mais rápido num pipeline serializado, onde cada segundo é fila.
	if strings.Contains(texto, "go install github.com/co2-lab/anchors") {
		t.Error("`go install` não resolve neste repositório (go.mod em cli/ declara a raiz)")
	}
	if !strings.Contains(texto, "gh release download") {
		t.Error("o binário tem de vir da release")
	}
	// A comparação é entre o mapa COMMITADO e o que o `map build` produziria — e não
	// entre os arquivos do diff e o mapa. A primeira versão fazia isso e acusava todo
	// arquivo NÃO REGIDO (.gitignore, PROJECT.md, o próprio grafo) como "mapa
	// desatualizado": 4 falsos positivos no primeiro push real, com rebuild à toa.
	if !strings.Contains(texto, "diff -q anchors.graph.yaml") {
		t.Error("a verificação tem de comparar o mapa commitado com o reconstruído")
	}
	if strings.Contains(texto, `grep -qF "id: $f"`) {
		t.Error("perguntar se o arquivo do diff está no mapa confunde `não regido` com `mapa velho`")
	}
	// E quando está defasado, conserta em vez de mandar o dev consertar: o pipeline tem o
	// repositório na mão, e falhar aqui pararia a fila por algo que ele resolve sozinho.
	if !strings.Contains(texto, "anchors map build") {
		t.Error("mapa defasado tem de ser reconstruído pelo próprio pipeline")
	}
}

// O recorte por diff é o que mantém o custo constante conforme o projeto cresce: um
// artefato só fica órfão quando APARECE, e quando aparece ele está no diff.
func TestIdentifyOlhaSoOQueMudou(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-identify.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, "git diff --name-only") {
		t.Error("sem recorte por diff, cada push refaz o trabalho de todos os códigos do projeto")
	}
	// E a varredura completa continua alcançável para reconciliar o passado.
	if !strings.Contains(texto, "tudo:") {
		t.Error("falta a saída de reconciliação (varrer o projeto inteiro sob demanda)")
	}
}

// Todo pipeline precisa poder ser rodado À MÃO. Quando algo dá errado — o cron atrasou,
// um push não disparou, um card ficou preso — a saída é executar o workflow na hora, e
// sem `workflow_dispatch` o botão não existe: só resta um commit vazio ou esperar.
func TestTodoPipelineAceitaExecucaoManual(t *testing.T) {
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatalf("%s: %v", w.Arquivo, err)
		}
		if !strings.Contains(string(b), "workflow_dispatch:") {
			t.Errorf("%s não pode ser rodado à mão — sem saída quando algo trava", w.Arquivo)
		}
	}
}

// O estado do trabalho é uma LABEL, e nenhum pipeline toca o Project.
//
// A decisão anterior era a coluna do board, e ela cobrava um preço que só apareceu no
// uso: escrever num Project de organização exige um PAT com escopo `project` — que o
// GITHUB_TOKEN da Action não tem —, então todo projeto que adotasse o fluxo precisaria
// criar e manter um token pessoal antes de o primeiro card se mover.
//
// Este teste é o que impede a volta: um `gh project` num pipeline reintroduz o atrito.
func TestEstadoVivenaLabelSemTocarOBoard(t *testing.T) {
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatalf("%s: %v", w.Arquivo, err)
		}
		texto := string(b)
		if strings.Contains(texto, "gh project ") {
			t.Errorf("%s escreve no Project — isso exige PAT e vira atrito de adoção", w.Arquivo)
		}
		if strings.Contains(texto, "ANCHORS_PROJECT_TOKEN") {
			t.Errorf("%s ainda pede token de Project", w.Arquivo)
		}
		if strings.Contains(texto, "repository-projects: write") {
			t.Errorf("%s pede permissão de Project que não usa mais", w.Arquivo)
		}
	}

	// O card nasce COM o estado inicial: sem label, ele não aparece para o claim.
	b, _ := fs.ReadFile(workflowsFS, "workflows/anchors-identify.yml")
	if !strings.Contains(string(b), `--label "anchors:to-do"`) {
		t.Error("a issue criada precisa nascer com o estado inicial")
	}
	// E o claim escolhe pelas labels de estado.
	b, _ = fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if !strings.Contains(string(b), "anchors:ready-to-review") {
		t.Error("o claim tem de tirar candidatos das labels de estado")
	}
}

// Toda coluna citada nos pipelines tem de existir na lista oficial. Um nome digitado
// diferente (`TODO` em vez de `TO DO`) não falha em lugar nenhum: o `item-edit` não acha
// a opção, o card não se move, e o trabalho some do fluxo em silêncio.
func TestPipelinesSoUsamColunasDeclaradas(t *testing.T) {
	valida := map[string]bool{}
	for _, c := range WorkStates {
		valida[c] = true
	}
	// A label de ESCALAÇÃO é válida sem ser estado: ela marca QUEM destrava o card, e o
	// card continua na coluna onde o trabalho parou. Tratá-la como estado a faria sair
	// dessa coluna, e o board deixaria de mostrar onde o fluxo travou.
	valida[LabelNeedsUser] = true
	// A label ANTIGA continua válida enquanto durar a migração: os pipelines a aceitam
	// para não abandonar as issues que já a carregam.
	valida[LabelNeedsUserLegacy] = true
	// A SEGUNDA FILA DO USUÁRIO, pela mesma razão do `needs-user`: ela diz que o card
	// espera uma pessoa, e o card continua na coluna onde o trabalho parou. O que a
	// distingue é o que ela PEDE — "confira se isto é seu" em vez de "decida entre A e
	// B" —, e essa diferença é de peso, não de estado.
	valida[LabelNeedsFraming] = true
	// O VÍNCULO DE DESBLOQUEIO, pela mesma razão do `needs-user` acima: ele diz que a
	// entrega DESTE card destrava outro, e o card continua na coluna onde o trabalho está.
	// O prefixo é conferido sem o número, que é por card e não se pode enumerar.
	valida[PrefixoLabelDesbloqueia] = true
	// O VÍNCULO DO ACHADO, irmão do de desbloqueio: `under-<n>` diz que este card nasceu
	// SOB outro, e o card continua na coluna onde o trabalho está. São os dois caminhos
	// para um card ficar preso a outro — o `escalate` produz este, o `unblock` produz
	// aquele —, e o board precisa dos dois para saber se a espera ainda é real.
	valida[PrefixoLabelSob] = true
	// O DESCARTE, pela mesma razão das outras: ele diz que o card saiu do BOARD, não que
	// ele está numa coluna. Tratá-lo como estado o faria aparecer como uma — e o ponto do
	// descarte é exatamente não aparecer.
	valida[LabelDiscarded] = true
	// O OPT-OUT da trava de estado, pela mesma razão das duas acima: ele autoriza mover o
	// card à mão, e o card continua onde o trabalho está. Tratá-lo como estado o faria
	// sair da coluna — e o board deixaria de mostrar o que ele autoriza.
	valida[LabelManual] = true
	// O BLOQUEIO, terceira direção do mesmo vínculo: `blocked-by-<n>` diz que ESTE card
	// espera o #n. Não é estado pela mesma razão do `needs-user` — o card continua na
	// coluna onde o trabalho parou, e é isso que faz o board mostrar ONDE ele travou.
	// Fosse estado, um card bloqueado sairia de `in-progress` e o board diria que o
	// trabalho nunca começou.
	valida[PrefixoLabelBlockedBy] = true
	// A PROCEDÊNCIA PELO PR, irmã do `under-<n>`: ela diz onde o achado foi VISTO, e o
	// card continua na coluna onde o trabalho está. Os dois coexistem — um responde por
	// onde o achado se entrega, o outro por onde se rastreia a revisão.
	valida[PrefixoLabelDePR] = true
	// THE BUG, like `needs-user`: it says what the card waits for (a fix where no agent in
	// the queue edits), not which column it is in.
	valida[LabelBug] = true
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatalf("%s: %v", w.Arquivo, err)
		}
		for _, m := range regexp.MustCompile(`"(anchors:[a-z-]+)"`).FindAllStringSubmatch(string(b), -1) {
			if !valida[m[1]] {
				t.Errorf("%s usa o estado %q, que não está em EstadosDoTrabalho", w.Arquivo, m[1])
			}
		}
	}
}

// O Anchors escreve até `READY TO TEST` e para: as três últimas colunas são dos pipelines
// de entrega do projeto. Um pipeline do Anchors que movesse para lá passaria por cima de
// uma decisão que não é dele (quando um teste passou, quando um deploy aconteceu).
func TestAnchorsNaoEscreveAlemDeReadyToTest(t *testing.T) {
	naoNossas := map[string]bool{
		"anchors:in-test": true, "anchors:ready-to-release": true, "anchors:production": true,
	}
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatal(err)
		}
		for estado := range naoNossas {
			if strings.Contains(string(b), `--add-label "`+estado+`"`) {
				t.Errorf("%s move para %q — além da alçada do Anchors", w.Arquivo, estado)
			}
		}
	}
	if AnchorsFinalState != "anchors:ready-to-test" {
		t.Errorf("o último estado que o Anchors escreve mudou: %q", AnchorsFinalState)
	}
}

// O board é COMPARTILHADO: carrega issues de produto, de infra, do que o time quiser.
// Todo pipeline que ESCREVE em cards tem de filtrar pela label do Anchors — sem isso, um
// agente comentaria `anchors-owner` numa issue de produto e a moveria para `IN PROGRESS`,
// sequestrando trabalho que não é dele. E o dono real dessa issue não tem como saber:
// ninguém procura por um campo do Anchors numa issue que não é do Anchors.
func TestPipelinesSoTocamCardsDoAnchors(t *testing.T) {
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatalf("%s: %v", w.Arquivo, err)
		}
		texto := string(b)
		// A regra vale para quem MEXE em issue. Um pipeline que só lê o repositório e
		// reporta (o `gates`) não tem card a filtrar — exigir a label dele obrigaria a
		// escrever uma consulta que ele não faz, só para satisfazer o teste.
		if !strings.Contains(texto, "gh issue") {
			continue
		}
		// `-f label="$LABEL"` é o mesmo filtro na consulta GraphQL paginada.
		if !strings.Contains(texto, `--label "$LABEL"`) && !strings.Contains(texto, `-f label="$LABEL"`) {
			t.Errorf("%s não filtra pela label do Anchors — pode tocar card de outro fluxo", w.Arquivo)
		}
	}

	// O claim seleciona por DUAS labels ao mesmo tempo: a do Anchors (o quintal) e a do
	// estado (a fila). Consultar só o estado pegaria uma issue de produto que alguém
	// tivesse rotulado igual.
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `--label "$LABEL" --label "$estado"`) {
		t.Error("o claim precisa cruzar a label do Anchors com a do estado")
	}
}

// A liberação de card órfão é um comentário NOVO, nunca uma edição do anterior: o
// histórico de quem passou pelo card é o que permite saber que ele foi liberado por
// inatividade, e não por decisão de alguém.
func TestStalePreservaOHistorico(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-stale.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, "gh issue comment") {
		t.Error("a liberação tem de ser um comentário novo (preserva o histórico)")
	}
	if strings.Contains(texto, "--edit-last") {
		t.Error("editar o comentário anterior apagaria quem tinha o card")
	}
}

// O BOARD NÃO PODE SUBSTITUIR O SITE. `actions/deploy-pages` publica o artefato como o
// site INTEIRO — num projeto que já tem landing page ou documentação no Pages, isso
// trocaria o site pelo board. A perda é silenciosa: só se descobre quando alguém abre o
// endereço e o site sumiu.
//
// Este teste é o que impede a volta ao deploy direto, que é o caminho óbvio e errado.
func TestBoardNaoSubstituiOSiteExistente(t *testing.T) {
	b, err := workflowsFS.ReadFile("workflows/anchors-board.yml")
	if err != nil {
		t.Fatal(err)
	}
	// Só as linhas EXECUTÁVEIS: o template explica em comentário por que não usa
	// `deploy-pages`, e casar o texto inteiro reprovaria a própria explicação.
	var exec []string
	for _, l := range strings.Split(string(b), "\n") {
		if t := strings.TrimSpace(l); t != "" && !strings.HasPrefix(t, "#") {
			exec = append(exec, l)
		}
	}
	texto := strings.Join(exec, "\n")
	for _, proibido := range []string{"deploy-pages", "upload-pages-artifact"} {
		if strings.Contains(texto, proibido) {
			t.Errorf("o board usa `%s`, que publica o artefato como o site INTEIRO — "+
				"um projeto com site existente o perderia", proibido)
		}
	}
	// E a escrita tem de ser numa SUBPASTA, não na raiz do branch.
	if !strings.Contains(texto, "SUBPASTA") {
		t.Error("o board deveria escrever numa subpasta, para não tocar no resto do site")
	}
}

// O CARD ESCALADO sai da fila. Sem isso a escalação não teria efeito: o pipeline
// entregaria a um agente o card que já se sabe que não converge, e a décima primeira
// revisão produziria a décima segunda.
func TestClaimPulaCardEscalado(t *testing.T) {
	b, err := workflowsFS.ReadFile("workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	// O filtro aceita as DUAS labels durante a migração, então o teste confere que a
	// NOVA está presente — a antiga é tolerância, não requisito.
	if !strings.Contains(texto, `index("`+LabelNeedsUser+`")`) {
		t.Error("o claim precisa EXCLUIR o card escalado da lista de disponíveis — " +
			"senão a escalação vira só um rótulo")
	}
	// E precisa escalar em algum limite: contar sem agir deixaria o ciclo rodando.
	if !strings.Contains(texto, "--add-label \""+LabelNeedsUser+"\"") {
		t.Error("o claim precisa APLICAR a label ao atingir o limite")
	}
}

// A label de escalação tem de ser CRIADA no repositório. `gh issue edit` com label
// inexistente não é erro fatal — ele falha em silêncio, e o card ficaria travado sem o
// sinalizador que diz por quê.
func TestLabelDeEscalacaoNaoEhEstado(t *testing.T) {
	for _, e := range WorkStates {
		if e == LabelNeedsUser {
			t.Fatal("a escalação NÃO é estado: o card continua na coluna onde o trabalho " +
				"parou, e o que muda é quem pode destravá-lo")
		}
	}
}

// O SCRIPT de cada pipeline tem de ser bash VÁLIDO, e isso não é óbvio: o YAML valida a
// estrutura do arquivo e não olha o conteúdo de `run:`. Um erro de sintaxe passa por
// todos os testes e só aparece quando o pipeline roda — em produção, com um card já
// atribuído pela metade.
//
// Medido: um heredoc com o delimitador de fecho INDENTADO (o YAML exige a indentação; o
// bash exige a coluna zero) quebrou o claim com "here-document delimited by end-of-file",
// depois de já ter comentado no card.
func TestScriptsDosPipelinesSaoBashValido(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("sem bash no PATH")
	}
	for _, w := range WorkflowsDoFluxo {
		b, err := workflowsFS.ReadFile("workflows/" + w.Arquivo)
		if err != nil {
			t.Fatal(err)
		}
		var doc struct {
			Jobs map[string]struct {
				Steps []struct {
					Name string `yaml:"name"`
					Run  string `yaml:"run"`
				} `yaml:"steps"`
			} `yaml:"jobs"`
		}
		if err := yaml.Unmarshal(b, &doc); err != nil {
			t.Fatalf("%s: %v", w.Arquivo, err)
		}
		for _, job := range doc.Jobs {
			for _, s := range job.Steps {
				if strings.TrimSpace(s.Run) == "" {
					continue
				}
				cmd := pipelineBash(t, "-n")
				cmd.Stdin = strings.NewReader(s.Run)
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Errorf("%s / %q: script inválido:\n%s", w.Arquivo, s.Name, out)
				}
			}
		}
	}
}

// UM PIPELINE DESATUALIZADO falha em silêncio: ele roda, faz o que a versão dele sabia
// fazer, e o que foi corrigido depois simplesmente não acontece. Medido — um passo novo do
// `identify` não rodou por três execuções sem que nada acusasse.
//
// Os pipelines que TÊM o Anchors instalado devem se autoverificar. Os que não têm ficam de
// fora de propósito: instalar o binário só para a conferência os tornaria mais lentos a
// cada execução, e o `claim` é chamado a cada pedido de trabalho.
func TestPipelinesComAnchorsSeAutoverificam(t *testing.T) {
	for _, w := range WorkflowsDoFluxo {
		b, err := workflowsFS.ReadFile("workflows/" + w.Arquivo)
		if err != nil {
			t.Fatal(err)
		}
		texto := string(b)
		if !strings.Contains(texto, "Instalar o Anchors") {
			continue // sem o binário, não há como conferir
		}
		if !strings.Contains(texto, "--check-pipelines") {
			t.Errorf("%s instala o Anchors e não confere se está atualizado — "+
				"um pipeline velho roda e o que foi corrigido depois não acontece", w.Arquivo)
		}
		// E quem decide se isto BARRA é o projeto, pelo `stale_pipeline_blocks` — não o
		// YAML. As duas formas de decidir aqui dentro estão proibidas, e por motivos
		// opostos: `continue-on-error: true` engoliria a escolha de quem prefere barrar,
		// e um `exit 1` escrito à mão barraria todo mundo, inclusive quem quer só o aviso.
		i := strings.Index(texto, "--check-pipelines")
		passo := texto[max(0, i-400):min(len(texto), i+400)]
		if regexp.MustCompile(`(?m)^\s*continue-on-error\s*:`).MatchString(passo) {
			t.Errorf("%s: `continue-on-error` anula quem declarou `stale_pipeline_blocks: "+
				"true` — a decisão é do anchors.yaml, não do YAML do pipeline", w.Arquivo)
		}
		if strings.Contains(passo, "exit 1") {
			t.Errorf("%s: `exit 1` à mão barra todo mundo — o comando já sai com o código "+
				"certo conforme o projeto declarou", w.Arquivo)
		}
	}
}

// UM `select` NO FIM DE UM PIPE não filtra o CAMPO — filtra o objeto inteiro: quando ele
// reprova, o jq não produz valor nenhum, e o item some do array.
//
// Medido no app de referência: o board publicava 3 de 7 cards. Os 4 ausentes eram exatamente os
// de dono LIBERADO — que é o estado normal de quem terminou o trabalho. O board apagava o
// trabalho concluído, e como ele ainda mostrava ALGUNS cards, parecia estar funcionando.
func TestBoardNaoDerrubaCardPeloFiltroDeDono(t *testing.T) {
	b, err := workflowsFS.ReadFile("workflows/anchors-board.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	// `"owner:"` cru também casa dentro de `anchors-owner:` — e casava PRIMEIRO,
	// numa linha de comentário, fazendo a janela de 400 caracteres terminar antes
	// do campo que esta régua existe para medir. A âncora precisa ser o CAMPO.
	i := strings.Index(texto, "\n                    owner:")
	if i < 0 {
		t.Fatal("o board deixou de extrair o dono — se isso é intencional, remova este teste")
	}
	expr := texto[i:min(len(texto), i+400)]
	if strings.Contains(expr, "select(test(") {
		t.Error("`select` na expressão do `dono` derruba o CARD inteiro quando o dono está " +
			"liberado, e liberado é o estado normal de quem terminou. Use `if/then/else`, " +
			"que devolve \"\" e preserva o item")
	}
	if !strings.Contains(expr, "if test(") {
		t.Error("o dono liberado deve virar campo VAZIO (if/then/else), não sumir")
	}
}

// O BOARD INTEIRO cai para zero cards quando UM campo devolve `null`.
//
// `jq` não isola: `test()` sobre `null` aborta o filtro com código 5 e o board devolve
// `{"erro": ...}` em vez de itens. Foi o que aconteceu ao cortar o dono na primeira
// linha — `"" | split("\n")` devolve lista VAZIA, `first` disso é `null`, e todo card
// sem comentário de dono (o caso comum) matava a página.
//
// A régua exige que todo índice sobre `split` tenha um fallback, porque a falha não é
// local: ela apaga o board.
func TestIndiceSobreSplitNoBoardTemFallback(t *testing.T) {
	b, err := workflowsFS.ReadFile("workflows/anchors-board.yml")
	if err != nil {
		t.Fatal(err)
	}
	for _, linha := range strings.Split(string(b), "\n") {
		corte := strings.TrimSpace(linha)
		if strings.HasPrefix(corte, "#") || !strings.Contains(corte, "split(") {
			continue
		}
		// `split(...)[0]` e `split(...)|first` devolvem `null` em lista vazia. Só
		// vale se houver um `//` protegendo o resultado na mesma linha.
		if !semFallbackRE.MatchString(corte) {
			continue
		}
		if !strings.Contains(corte, "//") {
			t.Errorf("índice sobre `split` sem fallback `// \"\"` derruba o board inteiro "+
				"quando a lista vem vazia (`null cannot be matched`):\n  %s", corte)
		}
	}
}

var semFallbackRE = regexp.MustCompile(`split\([^)]*\)\s*(\[[0-9]+\]|\|\s*first)`)

// A FRONTEIRA REAL: alguém tem de confrontar os gates no PR.
//
// O pre-commit roda na máquina de quem commita, e `git commit --no-verify` o contorna
// com uma flag. Medido no projeto de referência: os quatro checks do PR estavam VERDES e
// nenhum deles confrontava gate algum — a frase "nada sobe se não passar" (QUALITY.md §8)
// não era verdade ali.
func TestAlgumPipelineConfrontaOsGatesNoPR(t *testing.T) {
	achou := false
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatal(err)
		}
		texto := string(b)
		// `anchors check` no CORPO DE UM CARD é instrução para quem lê, não execução —
		// e foi assim que este teste passou a acusar o `identify` de não usar `--all`.
		// O que interessa é o pipeline RODAR o comando, num bloco `run:`.
		if !regexp.MustCompile(`(?m)^\s+(anchors|.*&&\s*anchors) check`).MatchString(texto) {
			continue
		}
		achou = true
		if !strings.Contains(texto, "pull_request") {
			t.Errorf("%s roda os gates mas não no PR — a fronteira é ali", w.Arquivo)
		}
		// `--all` e não `--changed`: um PR pode QUEBRAR arquivo que não tocou (renomear
		// símbolo, mover spec), e o raio de impacto do `--changed` é calculado sobre o
		// mapa — que pode estar velho justamente por causa deste PR.
		if !strings.Contains(texto, "check --all") {
			t.Errorf("%s deveria confrontar `--all`: o `--changed` não alcança o que o PR "+
				"quebrou sem tocar", w.Arquivo)
		}
		// NÃO registra: abrir card a cada push de PR encheria o board de trabalho que o
		// autor ainda vai corrigir na revisão seguinte.
		if !strings.Contains(texto, "--no-record") {
			t.Errorf("%s registra a partir do PR — o board viraria ruído; use `--no-record`",
				w.Arquivo)
		}
	}
	if !achou {
		t.Error("nenhum pipeline confronta os gates: o pre-commit é contornável com " +
			"`--no-verify`, e sem isto 'nada sobe se não passar' é só uma frase")
	}
}

// O CARD APONTA O GUIA; QUEM ENSINA É O BINÁRIO.
//
// O gate abre issue do que ELE detecta. O que o agente descobre sozinho — uma config que
// contradiz a doutrina, um caminho que ninguém documentou — não tem quem registre, e o
// caminho barato é consertar na hora: o conserto some do histórico.
//
// A instrução PODERIA ser copiada no corpo de cada card, e a primeira versão foi assim.
// O custo não é só repetição: texto copiado CONGELA. Um card criado hoje carrega a
// instrução de hoje, e quando ela muda, os cards antigos passam a ensinar o errado sem
// que nada acuse. No binário, o guia acompanha a versão.
func TestCardApontaOGuiaDeTrabalho(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-identify.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, "anchors guide work") {
		t.Error("o card deve apontar o guia de trabalho — é onde a instrução vive")
	}
	// E NÃO pode carregar a instrução inteira: seria a versão congelada no dia em que o
	// card nasceu.
	if strings.Contains(texto, "anchors escalate") {
		t.Error("a instrução do `escalate` não pode ser copiada no card: copiada, ela " +
			"congela — use `anchors guide work`, que acompanha a versão do binário")
	}
}

// O VÍNCULO CARD↔PR É DO ANCHORS; A PALAVRA-CHAVE É DA PLATAFORMA.
//
// A primeira versão deste passo casava `closes|fixes|resolves` no corpo — uma exigência
// da PLATAFORMA vazando para dentro da doutrina. Estava errada pelo mesmo motivo que
// tiramos match de prosa dos gates: o Anchors é multi-idioma, e obrigar o texto do PR a
// estar em inglês não é régua do Anchors.
//
// A régua é: os cards que este trabalho fecha estão declarados? Quem sabe QUAIS é o
// Anchors (`anchors-owner:` e `anchors:under-<n>`); quem sabe a SINTAXE é o `pr-body`.
func TestPipelineConfrontaOVinculoENaoAPalavra(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-gates.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, "anchors pr-body") {
		t.Error("o pipeline deve confrontar o corpo com o que o `anchors pr-body` geraria — " +
			"é o que mantém a sintaxe da plataforma fora da doutrina")
	}
	// LER a sintaxe é diferente de EXIGI-LA. O pipeline extrai do corpo os cards que o
	// PR declara — e para isso precisa reconhecer a forma que a plataforma usa. O que não
	// pode é ele MANDAR escrever daquele jeito: essa é a exigência que vazaria para a
	// doutrina, e quem a resolve é o `pr-body`, gerando a linha pronta.
	//
	// A distinção está na mensagem: o aviso fala do que FALTA declarar, não da palavra.
	if strings.Contains(texto, "Use \\`Closes #N\\`") {
		t.Error("o pipeline não pode MANDAR escrever a palavra em inglês: num projeto " +
			"noutro idioma isso é a plataforma vazando para a doutrina. Mande rodar o " +
			"`anchors pr-body`, que gera a linha pronta")
	}
	if !strings.Contains(texto, "--so-sob") {
		t.Error("o pipeline deve conferir o que FALTA (os achados sob os cards declarados), " +
			"e não repetir o que o corpo já diz")
	}
}

// O CLAIM SERIALIZA QUEM PEGA O CARD, e isso não basta.
//
// O card é solto ao fim da implementação (`anchors-owner: (liberado)`) e volta à fila para
// REVISÃO — mas a branch e o PR do primeiro agente continuam de pé. O claim seguinte o
// entrega como se fosse trabalho novo, e o agente que o recebe reimplementa do zero.
//
// Medido no app de referência: TRÊS agentes resolveram a issue #375 em paralelo, sete minutos entre
// o primeiro PR e o terceiro. Os três chegaram à mesma solução correta — o desperdício foi
// de coordenação, não de qualidade. Dois PRs completos foram fechados como duplicata.
func TestClaimPulaCardQueJaTemPRAberto(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	if !strings.Contains(s, "gh pr list --state open") {
		t.Error("o claim precisa consultar os PRs ABERTOS antes de entregar um card")
	}
	// A busca é pelo NÚMERO no corpo do PR — onde o `anchors pr-body` escreve a linha de
	// fechamento. Por branch não serve: o nome dela é livre, e os três PRs do caso medido
	// tinham nomes diferentes.
	if !strings.Contains(s, "in:body") {
		t.Error("a busca deve ser pelo número do card no CORPO do PR, não pelo nome da branch")
	}
	if !strings.Contains(s, "revisar o PR, não reimplementar") {
		t.Error("a mensagem deve dizer o que fazer: um PR aberto pede revisão, não trabalho novo")
	}
	// `continue` e não `break`: o card com PR é PULADO, e o claim segue procurando outro.
	// Um `break` entregaria a fila vazia a quem tem trabalho disponível logo abaixo.
	i := strings.Index(s, "já tem o PR")
	if i < 0 {
		t.Fatal("a linha que anuncia o PR aberto sumiu")
	}
	depois := s[i:min(i+200, len(s))]
	if !strings.Contains(depois, "continue") {
		t.Error("o card com PR deve ser PULADO (`continue`), não encerrar a busca")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// O CARD BLOQUEADO por um card de desbloqueio não volta à fila antes da entrega.
//
// `needs-user` diz que um card espera uma pessoa. Quando a decisão dessa pessoa GERA
// TRABALHO, o trabalho vira um card com `anchors:desbloqueia-<n>` — e o bloqueado espera
// por ELE, não mais pela pessoa.
//
// Sem esta guarda o ciclo se repete: alguém remove o `needs-user` achando que decidir
// bastava, o claim devolve o card, e o agente que o pega encontra o mesmo impasse e escala
// de novo. Medido: um card do projeto de referência esperava uma decisão sobre qual branch
// aplicar uma correção, e o `anchors next` continuava re-servindo o mesmo card.
func TestClaimRespeitaOCardDeDesbloqueio(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	if !strings.Contains(s, "anchors:desbloqueia-$n") {
		t.Error("o claim precisa procurar o card que destrava o candidato")
	}
	if !strings.Contains(s, "--state open") {
		t.Error("só o card de desbloqueio ABERTO bloqueia — fechado significa entregue")
	}
	// A mensagem tem de dizer POR QUEM ele espera: "pulando" sozinho manda quem lê o log
	// procurar a razão no lugar errado.
	if !strings.Contains(s, "(desbloqueio)") {
		t.Error("a linha do log deveria distinguir este motivo dos outros dois que pulam card")
	}
	// `continue`, não `break`: o card bloqueado é pulado e a busca segue.
	i := strings.Index(s, "espera a entrega do #$bloqueio")
	if i < 0 {
		t.Fatal("a linha que anuncia o bloqueio sumiu")
	}
	if !strings.Contains(s[i:min(i+200, len(s))], "continue") {
		t.Error("o card bloqueado deve ser PULADO, não encerrar a busca")
	}
}

// O PIPELINE cuida de UM caso: o desbloqueio entregue.
//
// Há dois jeitos de um card em `needs-user` voltar à fila:
//
//	· decisão simples → `anchors decided --card N --resolution "..."`, que exige a REGRA
//	  que nasceu da decisão. É comando, e não pipeline, porque a resolução não se inventa;
//	· decisão que GEROU TRABALHO → `anchors unblock` cria o card de desbloqueio, e quando
//	  ele fecha NINGUÉM está olhando: não há evento no card bloqueado, e quem decidiu já
//	  seguiu para outra coisa.
//
// Medido: o card #311 do projeto de referência teve o defeito consertado e continuou
// parado. A condição escrita na instrução ("decida, e depois remova a label") tinha sido
// cumprida, e a label ficou.
func TestPipelineDecidedCuidaDoDesbloqueioEntregue(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-decided.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	if !strings.Contains(s, "--remove-label") {
		t.Error("o pipeline precisa REMOVER a label quando os desbloqueios fecham")
	}
	// O CARD SEM desbloqueio nenhum não é caso deste pipeline: ele espera uma PESSOA, e
	// liberá-lo aqui puliaria a resolução que o `decided` exige.
	// A DISTINÇÃO entre "sem desbloqueio nenhum" e "desbloqueio entregue" precisa estar no
	// FLUXO, não só numa busca que ninguém usa.
	//
	// A primeira versão deste teste checava só a presença de `--state all`, e sobreviveu à
	// mutação que remove o `if` inteiro: a busca continuava no arquivo, e o pipeline
	// passava a liberar TODO card escalado — inclusive os que esperam uma pessoa, pulando
	// a resolução que o `anchors decided` exige.
	if !strings.Contains(s, `if [ "$todos" = "0" ]; then`) {
		t.Error("o card SEM desbloqueio nenhum tem de ser pulado: ele espera uma PESSOA, e " +
			"liberá-lo aqui pularia a resolução que o `anchors decided` exige")
	}
	if !strings.Contains(s, `--state all`) {
		t.Error("a contagem de `todos` precisa incluir os fechados — é ela que distingue " +
			"`sem desbloqueio` de `desbloqueio entregue`")
	}
	// E ALGUM ainda aberto mantém a label.
	i := strings.Index(s, "label MANTIDA")
	if i < 0 {
		t.Fatal("a linha que anuncia o bloqueio remanescente sumiu")
	}
	if !strings.Contains(s[i:min(i+120, len(s))], "continue") {
		t.Error("com desbloqueio aberto, o card é PULADO — a label não sai")
	}
}

// A TAG `#solution` NÃO pode voltar: ela criaria um segundo caminho de destrave, mais
// barato que o `anchors decided` e sem exigir a resolução. Dois jeitos de liberar o mesmo
// card, um deles deixando a decisão sem rastro.
func TestPipelineDecidedNaoInterpretaComentario(t *testing.T) {
	b, _ := fs.ReadFile(workflowsFS, "workflows/anchors-decided.yml")
	s := string(b)
	if strings.Contains(s, "#solution") {
		t.Error("o pipeline não deve reconhecer tag em comentário — quem registra a decisão " +
			"é o `anchors decided`, que exige a revisão que nasceu dela")
	}
	if strings.Contains(s, "issue_comment") {
		t.Error("o gatilho é o FECHAMENTO do card de desbloqueio, não um comentário")
	}
}

// O GATILHO precisa cobrir os dois caminhos, e o segundo é o que se esquece.
//
// O comentário é o caminho natural (a pessoa responde, o card anda em segundos). Mas um
// card cujo DESBLOQUEIO fechou não gera comentário nenhum nele — e sem o cron ele ficaria
// parado esperando um evento que não vem.
func TestPipelineDecidedTemOsDoisGatilhos(t *testing.T) {
	b, _ := fs.ReadFile(workflowsFS, "workflows/anchors-decided.yml")
	s := string(b)
	// O `schedule` precisa do CRON junto: a palavra sozinha aparece na prosa que explica
	// por que ele existe, e a mutação que remove as duas linhas do agendamento passava.
	for _, gatilho := range []string{"issues:", "workflow_dispatch:"} {
		if !strings.Contains(s, gatilho) {
			t.Errorf("o pipeline precisa do gatilho %q", gatilho)
		}
	}
	if !strings.Contains(s, "- cron:") {
		t.Error("o `schedule` precisa do cron: uma issue fechada pela API fora do fluxo não " +
			"dispara o evento, e sem o agendamento o card espera para sempre")
	}
}

// A TRAVA DE ESTADO: o card se move pelo FATO, e nada além do pipeline o move.
//
// Medido no projeto de referência: um agente fechou o card #432 às 15:47 e abriu o PR às
// 15:48 — o mesmo padrão em oito PRs seguidos. O efeito não aparece no card, aparece na
// FILA: o card sai de `ready-to-review` antes de alguém revisar, o claim (que procura
// revisão primeiro) não acha nada e entrega trabalho NOVO.
//
// Resultado: 25 PRs verdes esperando revisão com ZERO cards em `ready-to-review`.
func TestTravaDeEstadoRevertOQueNaoVeioDoFluxo(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	// O PRÓPRIO PIPELINE não se vigia: ele move card o tempo todo, e reagir aos próprios
	// eventos seria um laço. É a única distinção que o `actor` permite — os agentes usam a
	// MESMA conta que a pessoa.
	if !strings.Contains(s, "github.event.sender.login != 'github-actions[bot]'") {
		t.Error("o pipeline precisa ignorar os próprios eventos, senão vira laço")
	}

	// OS QUATRO EVENTOS, e cada um tem o seu inverso.
	for _, ev := range []string{"closed", "reopened", "labeled", "unlabeled"} {
		if !strings.Contains(s, ev) {
			t.Errorf("o pipeline precisa reagir a %q", ev)
		}
	}
	for _, inverso := range []string{"gh issue reopen", "gh issue close", "--remove-label", "--add-label"} {
		if !strings.Contains(s, inverso) {
			t.Errorf("falta o inverso %q — reverter é desfazer, não avisar", inverso)
		}
	}

	// O OPT-OUT precisa estar na GUARDA, não só explicado em comentário.
	//
	// A primeira versão procurava a label em qualquer lugar do arquivo, e sobreviveu à
	// mutação que remove o `if` inteiro: `anchors:manual` continuava aparecendo na prosa
	// que explica por que ele existe, e a trava passava a não ter escape nenhum.
	//
	// Terceira vez hoje que a mesma armadilha aparece — procurar a string em vez do que a
	// usa. A asserção agora casa a guarda completa.
	if !strings.Contains(s, `any(. == "`+LabelManual+`")`) {
		t.Errorf("o opt-out %q precisa ser CONSULTADO, não só mencionado — uma trava sem "+
			"escape transforma percalço em trabalho parado", LabelManual)
	}
	if !strings.Contains(s, "movimento manual autorizado") {
		t.Error("o pipeline precisa SAIR quando o opt-out está presente")
	}
}

// SÓ AS LABELS DE ESTADO são revertidas.
//
// O `escalate` põe `needs-user` e `sob-<n>`; o `unblock` põe `desbloqueia-<n>`. Os dois
// fazem isso pela CLI, que autentica como a PESSOA — reverter essas labels desfaria o
// trabalho dos próprios comandos do Anchors.
func TestTravaDeEstadoNaoTocaLabelQueNaoEhEstado(t *testing.T) {
	b, _ := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	s := string(b)

	i := strings.Index(s, `case "$LABEL_MEXIDA" in`)
	if i < 0 {
		t.Fatal("o filtro por label de estado sumiu — o pipeline reverteria `needs-user` e " +
			"`sob-<n>`, desfazendo o que o `escalate` acabou de fazer")
	}
	filtro := s[i:min(i+400, len(s))]
	// A ALÇADA DO ANCHORS para em `ready-to-test` — o que vem depois é dos pipelines de
	// entrega do PROJETO, e reverter ali desfaria trabalho de outro fluxo.
	daAlcada := []string{
		LabelToDo, "anchors:in-progress", "anchors:ready-to-review",
		"anchors:in-review", "anchors:ready-to-test",
	}
	for _, estado := range daAlcada {
		if !strings.Contains(filtro, estado) {
			t.Errorf("o estado %q precisa estar no filtro da reversão", estado)
		}
	}
	// E os de ENTREGA não podem estar: o Anchors escreve até `ready-to-test` e larga.
	for _, depois := range []string{"anchors:in-test", "anchors:ready-to-release", "anchors:production"} {
		if strings.Contains(filtro, depois) {
			t.Errorf("%q é dos pipelines de ENTREGA do projeto — o Anchors não o reverte", depois)
		}
	}
	for _, fora := range []string{LabelNeedsUser, PrefixoLabelDesbloqueia, PrefixoLabelSob} {
		if strings.Contains(filtro, fora) {
			t.Errorf("%q não é label de estado — revertê-la desfaria o `escalate`/`unblock`", fora)
		}
	}
}

// O COMENTÁRIO é o que transforma a reversão em ensinamento. Sem ele, quem mexeu vê o card
// voltar sozinho e conclui que o pipeline está quebrado — e a próxima reação é desligá-lo.
func TestTravaDeEstadoExplicaOQueFezEComoAutorizar(t *testing.T) {
	b, _ := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	s := string(b)
	for _, quer := range []string{"Revertido", "se move pelo FATO", "25 PRs verdes", LabelManual} {
		if !strings.Contains(s, quer) {
			t.Errorf("o comentário da reversão deveria conter %q", quer)
		}
	}
}

// A LISTA COM VÍRGULA precisa ser acusada: ela fecha UM card e parece fechar todos.
//
// O caso real, no projeto de referência. O PR #475 trazia no corpo:
//
//	Closes #419, #473
//
// A plataforma honra só o primeiro número. O PR mergeou, o #419 fechou, e o #473 ficou
// aberto com o trabalho já entregue — invisível, porque o corpo do PR dizia o contrário.
//
// O ESTRAGO não parou aí. Um agente notou e fechou o #473 à mão; a trava de estado reverteu,
// e estava certa pela regra que ela conhece (um card saindo da fila sem merge é exatamente
// o que ela existe para impedir). O card voltou para `in-progress` e ficou três horas parado
// esperando um trabalho que já existia.
//
// AVISO e não reprovação: o vínculo do PR com o card está declarado, e barrar um merge por
// sintaxe seria pior que o defeito. O que faltava era alguém DIZER que os números depois da
// vírgula não valem.
func TestPipelineAcusaListaDeCardsComVirgula(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-pr-checks.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	if !strings.Contains(texto, "extras=") {
		t.Fatal("o pipeline não extrai os números que vêm depois da vírgula — um PR com " +
			"`Closes #419, #473` mergearia deixando o segundo card aberto, em silêncio")
	}

	// A MENSAGEM não pode ditar a sintaxe da plataforma: é a mesma régua do
	// TestPipelineConfrontaOVinculoENaoAPalavra, e o aviso novo é um lugar fácil de
	// quebrá-la — a tentação é mostrar o texto certo em vez de mandar gerá-lo.
	if !strings.Contains(texto, "anchors pr-body --cards $") {
		t.Error("o aviso deve mandar RODAR o `anchors pr-body` com os cards do PR — " +
			"mostrar a linha pronta faria a sintaxe da plataforma vazar para a doutrina")
	}

	// A lista de números precisa chegar ao comando sugerido: um `pr-body` sem os cards
	// certos não resolve o problema de quem o roda.
	if !strings.Contains(texto, `todos="$card`) {
		t.Error("o comando sugerido deve juntar o card já declarado aos que a vírgula " +
			"engoliu — senão ele regenera só o que já estava certo")
	}
}

// O CARD NÃO FECHA NO MERGE — a trava reverte isso, e a razão é a esteira.
//
// Esta régua guardava o oposto: havia uma exceção que MANTINHA o fechamento quando um PR
// mergeado citava o card. A premissa era que o merge encerra o trabalho.
//
// Ele não encerra. A alçada do Anchors acaba em `ready-to-test`; depois vêm `in-test`,
// `ready-to-release` e `production` — do CD do projeto, que o Anchors não rastreia. Fechar
// no merge declara entregue um trabalho com três estados à frente, e quem ia testar perde
// de vista o que testar.
//
// MEDIDO no projeto de referência: 71 cards fechados em `ready-to-test` contra 7 abertos,
// 35 num só dia. A coluna que devia ACUMULAR o que espera teste mostrava só o resíduo — os
// que nunca receberam a palavra de fechamento.
//
// O CASO QUE A EXCEÇÃO PROTEGIA desapareceu com a causa. Ele era: o corpo trazia
// `Closes #419, #473`, a plataforma reconhecia só o primeiro, o #473 ficava aberto com o
// trabalho entregue, alguém fechava à mão e a trava revertia (medido: 3h parado). Com
// `Refs #N` (v0.1.130) nada fecha automaticamente, então não há mais card fechado pela
// vírgula nem fechamento manual a preservar.
//
// QUEM FECHA agora é quem termina a esteira, declarando `anchors:manual` (entregue) ou
// `anchors:discarded` (saiu do roadmap) — e a trava já isenta essas duas antes de chegar
// ao fechamento.
func TestOCardNaoFechaNoMerge(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// A CONSULTA DO TIMELINE não pode voltar: era ela que autorizava manter o fechamento.
	if strings.Contains(texto, "cross-referenced") {
		t.Error("a trava voltou a consultar `cross-referenced` no timeline — era assim que " +
			"ela isentava o fechamento de um card com PR mergeado, e o card precisa ficar " +
			"ABERTO em `ready-to-test` porque a esteira segue no CD")
	}
	if strings.Contains(texto, "Fechamento mantido") {
		t.Error("a trava voltou a ter a mensagem de fechamento mantido — o merge não " +
			"encerra o card, ele o move para `ready-to-test` e ali o card fica aberto")
	}

	// E O FECHAMENTO DECLARADO tem de continuar isento: `manual` é entregue,
	// `discarded` saiu do roadmap. Sem isso não haveria como encerrar um card nunca.
	for _, decl := range []string{"anchors:manual", "anchors:discarded"} {
		if !strings.Contains(texto, decl) {
			t.Errorf("a trava não conhece %q — sem uma forma DECLARADA de encerrar, o card "+
				"fica aberto para sempre e a trava reverte quem tentar fechá-lo", decl)
		}
	}
}

// QUEM LÊ DADOS DE PR PRECISA DECLARAR A PERMISSÃO.
//
// O `GITHUB_TOKEN` só concede o que o pipeline pede. Um workflow que lê `merged_at` sem
// `pull-requests: read` não recebe erro: o campo vem NULO — e o código que depende dele toma
// a decisão errada achando que decidiu certo.
//
// Continua valendo para todo pipeline que lê dados de PR — o `pr-checks` decide por
// `merged` se avança o card, e um campo nulo o faria tratar PR abandonado como entregue.
func TestPipelineQueLeDadosDePRDeclaraAPermissao(t *testing.T) {
	entradas, err := fs.ReadDir(workflowsFS, "workflows")
	if err != nil {
		t.Fatal(err)
	}
	if len(entradas) == 0 {
		t.Fatal("nenhum workflow embutido — o glob quebrou e o teste passaria vazio")
	}
	for _, e := range entradas {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+e.Name())
		if err != nil {
			t.Fatal(err)
		}
		texto := string(b)
		// SÓ A CONSULTA, não o payload do evento.
		//
		// `github.event.pull_request.base.sha` chega no payload que dispara o workflow: o
		// GitHub já o entregou, e nenhuma permissão o altera. Foi o falso positivo desta
		// régua na primeira versão — ela acusou o `anchors-identify.yml`, que está certo.
		//
		// O que precisa da permissão é PERGUNTAR à API: `gh api .../timeline`, `gh pr view`,
		// `gh pr list`. Aí o token decide o que devolve, e sem `pull-requests: read` os
		// campos de PR vêm nulos em vez de dar erro.
		consultaAPI := strings.Contains(texto, "gh api") ||
			strings.Contains(texto, "gh pr ")
		leDadosDePR := strings.Contains(texto, "merged_at") ||
			strings.Contains(texto, ".source.issue.pull_request")
		if consultaAPI && leDadosDePR && !strings.Contains(texto, "pull-requests: read") {
			t.Errorf("%s lê dados de PR e não declara `pull-requests: read` — o campo vem "+
				"NULO em vez de dar erro, e a decisão sai errada em silêncio", e.Name())
		}
	}
}

// TODAS as linhas `Closes` precisam ser conferidas, não só a primeira.
//
// A conferência que já existia olhava o PRIMEIRO número declarado — o que basta enquanto um
// PR fecha um card. Um PR que consolida vários traz uma linha para cada, e as demais
// passavam sem ninguém olhar.
//
// MEDIDO: um PR consolidando 19 trabalhos trouxe 19 linhas `Closes`, e todos os 19 números
// eram de PULL REQUESTS, não dos cards. No GitHub um PR também é uma issue — `gh issue view`
// aceita o número e responde normalmente —, então nada acusou. O merge fechou os PRs, e
// dezoito cards ficaram abertos com o trabalho já entregue.
//
// O que salva é a label: um PR não carrega a label do Anchors. A conferência teria pego o
// erro na primeira linha; o que faltava era ela olhar as outras dezoito.
func TestPipelineConfereTodasAsLinhasDeFechamento(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-pr-checks.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	if !strings.Contains(texto, "outros=") {
		t.Fatal("o pipeline confere só o primeiro `Closes` — num PR que fecha vários cards, " +
			"os demais números entram sem ninguém olhar se são cards de verdade")
	}
	if !strings.Contains(texto, "naoSaoCards") {
		t.Error("o pipeline deve ACUSAR os números que não são cards, e não só ignorá-los: " +
			"um número de PR ali não fecha card nenhum, e o silêncio é o defeito")
	}

	// O MESMO FILTRO da primeira linha. Duas formas de extrair a mesma coisa divergem na
	// primeira vez que uma delas mudar — e aqui divergir significa conferir a linha 1 com
	// uma régua e as linhas 2..N com outra.
	if !strings.Contains(texto, "emBloco = !emBloco") {
		t.Error("a extração das demais linhas deve pular bloco de código como a primeira faz")
	}
	// `tail -n +2` é o que distingue "as demais" de "todas": a primeira já tem a sua
	// própria conferência, com mensagem própria.
	if !strings.Contains(texto, "tail -n +2") {
		t.Error("as demais linhas são da segunda em diante — a primeira já é conferida acima")
	}
}

// O PIPELINE DE LIBERAÇÃO precisa conhecer OS DOIS prefixos de vínculo.
//
// Há dois comandos que prendem um card a outro, e cada um usa o seu:
//
//	anchors unblock   →  anchors:desbloqueia-<n>   a decisão gerou trabalho
//	anchors escalate  →  anchors:under-<n>         o achado virou card próprio
//
// São a mesma relação — há trabalho aberto sob este card —, e o pipeline olhava só o
// primeiro. MEDIDO: o card #311 do projeto de referência teve os dois achados que o
// travavam (#441, #442) fechados e continuou com `needs-user`, aparecendo no board como
// pendência do usuário em `ready-to-test`, com o trabalho já entregue pelo PR #402.
//
// Ninguém era responsável por tirar a label: o comando que libera é o `anchors decided`, e
// quem decidiu já seguiu para outra coisa. O pipeline existe exatamente para esse buraco — e
// não o cobria para metade dos cards.
func TestPipelineDecidedConheceOsDoisPrefixosDeVinculo(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-decided.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// O laço sobre os dois prefixos, e não uma consulta por prefixo fixo.
	if !strings.Contains(texto, "for pref in desbloqueia under") {
		t.Error("o pipeline consulta um prefixo só — metade dos cards escalados (os do " +
			"`escalate`, com `under-<n>`) nunca seria liberada")
	}

	// AMBOS os lados: contar quantos existem E quantos estão abertos. Cobrir só o
	// primeiro faria o pipeline achar que há vínculo e nunca ver que ele fechou.
	if n := strings.Count(texto, "for pref in desbloqueia under"); n < 2 {
		t.Errorf("o laço aparece %d vez(es) — as DUAS consultas precisam dele: a que conta "+
			"os vínculos existentes e a que vê quais ainda estão abertos", n)
	}

	// A construção da lista de abertos não pode perder o que veio do primeiro prefixo.
	if !strings.Contains(texto, `abertos="${abertos:+$abertos, #}$q"`) {
		t.Error("a lista de abertos sobrescreve em vez de acumular — um card com vínculo " +
			"aberto de cada tipo reportaria só um deles")
	}
}

// RESPONDIDO exige as DUAS coisas: já houve trabalho sob o card, E ele acabou.
//
// A primeira versão perguntava só "há card aberto sob este?" — e marcava como respondido
// todo card com zero abertos. Mas um card que NUNCA teve trabalho sob ele também tem zero,
// e ele é o caso mais comum: a decisão que ninguém tomou ainda não gerou trabalho nenhum.
//
// MEDIDO: o #400 do projeto de referência ("nenhum workflow roda `pnpm typecheck`") nunca
// foi respondido e apareceu no board como já-respondido. Um sinal invertido no caso mais
// frequente é pior que não ter sinal: ensina a ignorar a cor.
func TestBoardSoMarcaRespondidoQuandoHouveTrabalhoQueAcabou(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-board.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// As duas contagens: o que EXISTE (qualquer estado) e o que está ABERTO. Uma só não
	// distingue "acabou" de "nunca começou".
	if !strings.Contains(texto, "existem=") {
		t.Fatal("o pipeline conta só os cards ABERTOS sob o card — com essa pergunta, " +
			"um card que nunca teve trabalho sob ele vira `respondido`, e ele é o caso " +
			"mais comum de espera real")
	}
	if !strings.Contains(texto, `[ "$existem" -gt 0 ] && [ "$abertos" = "0" ]`) {
		t.Error("a condição de RESPONDIDO precisa exigir as duas coisas — já houve " +
			"trabalho, e ele acabou")
	}

	// Os dois prefixos, pela mesma razão do pipeline de liberação: `under-` vem do
	// `escalate`, `desbloqueia-` do `unblock`.
	if !strings.Contains(texto, "for pref in under desbloqueia") {
		t.Error("o board precisa contar os DOIS prefixos de vínculo — olhar só um faria " +
			"metade dos cards escalados parecer sem trabalho sob eles")
	}
}

// O CONFLITO QUE SÓ O MAPA CAUSA precisa ser diagnosticado — o GitHub não sabe resolvê-lo.
//
// O `anchors.graph.yaml` é GERADO e muda em todo PR. O Anchors traz um merge driver que sabe
// uni-lo (o `anchors map merge`), mas ele vive no clone de quem trabalha: o merge do servidor
// não o tem, e cai no merge textual.
//
// MEDIDO no projeto de referência: ONZE PRs de trabalho independente — cada um tocando só os
// arquivos da sua unidade — ficaram `CONFLICTING` de uma vez, porque um deles mergeou. O
// único arquivo em comum era o mapa.
//
// Quem abre o PR vê "conflitos" e presume que mexeu no que outro mexeu. Não mexeu — e a saída
// é um `git rebase` local, onde o driver está. Sem alguém DIZER isso, o PR espera.
func TestPipelineDiagnosticaOConflitoDoMapa(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-pr-checks.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	if !strings.Contains(texto, "conflito-do-mapa") {
		t.Fatal("o pipeline não diagnostica o conflito do mapa — onze PRs de trabalho " +
			"independente podem ficar parados sem ninguém saber que a saída é um rebase")
	}

	// A INTERSEÇÃO, e não a presença do mapa no diff. O mapa está em TODO PR; o que
	// distingue o caso mecânico é ele ser o ÚNICO arquivo que os dois lados tocaram.
	if !strings.Contains(texto, "comm -12") {
		t.Error("o diagnóstico precisa cruzar os arquivos dos DOIS lados — só olhar se o " +
			"mapa está no PR acusaria todo PR, e um aviso que sempre sai não é lido")
	}
	if !strings.Contains(texto, `[ "$comuns" = "anchors.graph.yaml" ]`) {
		t.Error("a condição deve exigir que o mapa seja o ÚNICO comum: com mais arquivos " +
			"em comum o conflito é de trabalho, e mandar rebasear esconderia isso")
	}

	// O CAMINHO DE SAÍDA, não só o diagnóstico. Dizer "é o mapa" sem dizer o que fazer
	// deixa o leitor exatamente onde estava.
	if !strings.Contains(texto, "git rebase origin/$base") {
		t.Error("o aviso deve dar o comando que resolve — o diagnóstico sozinho não destrava")
	}
	if !strings.Contains(texto, "anchors install-hooks") {
		t.Error("o aviso deve cobrir o caso de o driver NÃO estar instalado no clone, que " +
			"é quando o rebase pede resolução manual e a pessoa desiste")
	}
}

// O `board.json` É CONTRATO, e contrato do Anchors é em inglês.
//
// Os 18 campos nasceram em português — `numero`, `titulo`, `estado` — enquanto o resto do
// produto é inglês: as flags (`--about`, `--card`), as labels (`anchors:needs-user`), os
// identificadores Go (a régua `TestNenhumIdentificadorEmPortugues` os cobra).
//
// O JSON é o que um projeto de terceiro lê para construir a própria visão. Um contrato meio
// em cada idioma obriga quem o consome a saber os dois — e a decidir, campo a campo, qual
// esperar.
//
// A janela foi agora porque o Anchors ainda roda num projeto só. Depois, cada board
// publicado seria uma migração.
func TestBoardJSONEmIngles(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-board.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// Os nomes que existiam antes. Se um voltar, volta como campo do contrato.
	for _, velho := range []string{
		"numero", "titulo", "atualizado", "criado", "fechado", "autor", "corpo",
		"codigo", "alvo", "estado", "escalado", "vinculo", "destrava", "dono",
		"posse", "comentarios", "esperando", "fases", "itens", "quando",
	} {
		// Só a CHAVE do objeto — a palavra pode aparecer em prosa de comentário, e
		// proibir isso tornaria impossível explicar a migração.
		if regexp.MustCompile(`(?m)^\s{20}` + velho + `:`).MatchString(texto) {
			t.Errorf("o campo %q voltou ao `board.json` — o contrato é em inglês", velho)
		}
	}

	// E os nomes novos precisam estar lá: sem isto, apagar o campo passaria no teste.
	for _, novo := range []string{"number", "title", "state", "created", "closed", "code"} {
		if !regexp.MustCompile(`(?m)^\s{20}` + novo + `:`).MatchString(texto) {
			t.Errorf("o campo %q sumiu do `board.json` — quem lê o contrato o espera", novo)
		}
	}
}

// TODA `$variavel` USADA NUM PIPELINE PRECISA EXISTIR.
//
// O board parou de publicar com `waiting: unbound variable`. Uma renomeação trocou o USO de
// três variáveis sem trocar a DECLARAÇÃO — e o YAML continua válido, a sintaxe do shell
// continua válida, o JS continua válido. Nada acusou até o runner executar.
//
// As três eram da mesma natureza, e a do meio é a que derrubava tudo:
//
//	esperando=true          declarava      $waiting        usava
//	--arg t "$title"        usava          titulo=         declarava
//	--slurpfile itens       declarava      $items          o jq lia
//
// `set -u` pega isso — mas só no runner, depois do merge. Aqui é de graça.
func TestPipelinesNaoUsamVariavelInexistente(t *testing.T) {
	entradas, err := fs.ReadDir(workflowsFS, "workflows")
	if err != nil {
		t.Fatal(err)
	}
	if len(entradas) == 0 {
		t.Fatal("nenhum workflow embutido — o glob quebrou e o teste passaria vazio")
	}

	// O que o AMBIENTE dá: do GitHub Actions, do `env:` do passo, e os laços do shell.
	ambiente := regexp.MustCompile(`(?m)^\s*([A-Z_][A-Z_0-9]*):`)
	// A ATRIBUIÇÃO não vive só no começo da linha: `a=1; b=2` e `if ! x=$(cmd)` são as
	// duas formas que apareceram no fluxo, e exigir início de linha acusava as duas.
	atribui := regexp.MustCompile(`(?:^|[;&|(]|\bif\s+!?\s*|\bthen\s+|\bdo\s+|\s)([a-zA-Z_][\w]*)=`)
	laco := regexp.MustCompile(`\b(?:for|read(?:\s+-r)?)\s+(?:IFS=[^ ]*\s+read\s+-r\s+)?([a-zA-Z_][\w]*)`)
	// o jq declara variáveis por `--arg`, `--argjson` e `--slurpfile`
	jqVar := regexp.MustCompile(`--(?:arg|argjson|slurpfile|rawfile)\s+([a-zA-Z_][\w]*)`)
	// O jq também declara DENTRO do programa: `... as $m`, `... as [$a, $b]`. É a mesma
	// natureza de declaração, e ignorá-la acusava expressões corretas.
	jqAs := regexp.MustCompile(`\bas\s+\$([a-zA-Z_][\w]*)`)
	usa := regexp.MustCompile(`\$\{?([a-z_][a-zA-Z_0-9]*)\}?`)

	for _, e := range entradas {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+e.Name())
		if err != nil {
			t.Fatal(err)
		}
		texto := string(b)

		declaradas := map[string]bool{}
		for _, re := range []*regexp.Regexp{ambiente, atribui, laco, jqVar, jqAs} {
			for _, m := range re.FindAllStringSubmatch(texto, -1) {
				declaradas[m[1]] = true
			}
		}
		// A consulta GraphQL declara as SUAS variáveis no cabeçalho: `query($endCursor:
		// String)`. Elas vivem dentro de aspas simples — o shell não as expande —, e o
		// `--paginate` do gh preenche `$endCursor`. Mesma natureza do `as $x` do jq.
		for _, m := range regexp.MustCompile(`\$([a-zA-Z_][\w]*)\s*:\s*\[?[A-Z]`).
			FindAllStringSubmatch(texto, -1) {
			declaradas[m[1]] = true
		}
		// `read -r a b c d` declara TODAS — o laço captura só a primeira, e um `read` de
		// quatro campos é comum quando a linha vem de um `\t`-separado.
		for _, m := range regexp.MustCompile(`read\s+-r\s+([\w]+(?:\s+[\w]+)*)`).
			FindAllStringSubmatch(texto, -1) {
			for _, nome := range strings.Fields(m[1]) {
				declaradas[nome] = true
			}
		}
		// O `--arg` cuja variável vem na linha SEGUINTE, por continuação com `\`.
		for _, m := range regexp.MustCompile(`(?s)--(?:arg|argjson|slurpfile|rawfile)\s+\\?\s*\n?\s*([a-zA-Z_][\w]*)`).
			FindAllStringSubmatch(texto, -1) {
			declaradas[m[1]] = true
		}

		for _, m := range usa.FindAllStringSubmatch(texto, -1) {
			nome := m[1]
			if declaradas[nome] {
				continue
			}
			// `$takenAt` casa `taken` no prefixo minúsculo — confere o nome inteiro.
			if declaradas[strings.TrimSuffix(nome, "At")] {
				continue
			}
			t.Errorf("%s usa $%s e nada a declara — `set -u` derruba o passo no runner",
				e.Name(), nome)
		}
	}
}

// A AFILIAÇÃO RESOLVE NUM PASSE, e denuncia o que não fecha.
//
// A tentação é repetir passes até estabilizar: numa árvore de N níveis isso custa N
// varreduras da lista inteira. Pior — esconde a inconsistência, porque um vínculo que aponta
// para um card INEXISTENTE se parece com um que ainda não foi visitado.
//
// O desenho certo é registrar o que cada achado DECLARA e resolver depois, seguindo as
// referências com memória. O que não fecha é dado errado na issue, e merece nome.
func TestBoardAfiliaNumPasseEDenunciaOQueNaoFecha(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-board.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// NÃO pode haver laço de repetição sobre a lista inteira.
	if regexp.MustCompile(`for\s+_\w*\s+in\s+range\(`).MatchString(texto) {
		t.Error("a afiliação repete passes — numa árvore de N níveis isso custa N " +
			"varreduras, e um vínculo quebrado fica indistinguível de um ainda não visitado")
	}

	if !strings.Contains(texto, "declara = {}") {
		t.Fatal("a afiliação não registra o que cada achado DECLARA — sem isso ela depende " +
			"da ordem da lista, que é a ordem das issues no GitHub")
	}
	// MEMÓRIA: sem ela, uma cadeia compartilhada é percorrida uma vez por dependente.
	if !strings.Contains(texto, "resolvido = {}") {
		t.Error("a resolução não guarda o resultado — cadeias compartilhadas seriam " +
			"percorridas de novo a cada dependente")
	}

	// A DENÚNCIA, com os três motivos distinguidos. Chamar tudo de ciclo mandaria procurar
	// o que não há: uma cadeia que termina sem unidade é legítima.
	for _, m := range []string{"não existe", "ciclo de vínculos", "não é uma unidade"} {
		if !strings.Contains(texto, m) {
			t.Errorf("o diagnóstico não distingue %q — os três motivos pedem ações "+
				"diferentes de quem for consertar", m)
		}
	}
}

// O CARD DESCARTADO sai do board — e do DADO que o board publica.
//
// Fechar não basta: o board mostra os fechados porque o roadmap precisa deles para desenhar
// o passado. Um card de teste, ou um achado sobre arquivo que não existe mais, fica lá para
// sempre — fechado, sem pai possível, ocupando uma linha na raiz da árvore.
//
// MEDIDO no projeto de referência: das 27 raízes do roadmap, DOZE eram ruído permanente.
//
// O FILTRO É NA ORIGEM, e não na página: filtrar ao desenhar deixaria o card no JSON, e quem
// consome o JSON (um painel próprio, um relatório) o veria. O contrato é o que o board
// publica — se saiu do board, saiu do dado.
func TestBoardNaoPublicaCardDescartado(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-board.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, `any(. == "anchors:discarded") | not`) {
		t.Error("o board publica o card descartado — ele ficaria na raiz do roadmap para " +
			"sempre, sem pai possível e sem trabalho a fazer")
	}
	// Na MESMA consulta que lê as issues. Um filtro depois, sobre o JSON já montado,
	// dependeria de ninguém esquecer de aplicá-lo no caminho novo.
	i := strings.Index(texto, "gh api graphql --paginate")
	j := strings.Index(texto[i:], "> _board/board.json")
	if i < 0 || j < 0 || !strings.Contains(texto[i:i+j], "anchors:discarded") {
		t.Error("o filtro do descartado não está na consulta que lê as issues")
	}
}

// A TRAVA DE ESTADO não pode reabrir o card descartado.
//
// O `anchors discard` fecha o card depois de marcá-lo — e fechar à mão é exatamente o que a
// trava reverte. Sem esta saída, o descartado voltaria a abrir em dez segundos e o `claim` o
// serviria como trabalho novo: o mesmo ciclo que o #410 sofreu, por outra porta.
func TestTravaRespeitaOCardDescartado(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)
	if !strings.Contains(texto, `any(. == "anchors:discarded")`) {
		t.Fatal("a trava reabriria o card descartado — e o `claim` o serviria de novo")
	}
	// ANTES da reversão: conferir depois de reabrir não adianta.
	iDesc := strings.Index(texto, `any(. == "anchors:discarded")`)
	iRev := strings.Index(texto, "gh issue reopen")
	if iDesc > iRev {
		t.Error("a conferência do descartado vem DEPOIS da reversão — o card já teria " +
			"sido reaberto quando ela roda")
	}
}

// A FASE PRECISA SE DISTINGUIR DO CARD no roadmap.
//
// Um `## FNDTN-W01 — o CI` é uma SEÇÃO dentro do arquivo do plano: não tem issue, não tem
// unidade, e nenhum agente a pega porque não há o que pegar. Medido no projeto de referência:
// 43 fases, TODAS sem card.
//
// O plano, ao contrário, É trabalho — `plans/0014-alertas-incidentes.md` é um arquivo regido,
// com card próprio e unidade a cumprir; 18 deles existem, 2 já fechados por agentes.
//
// Desenhar os dois igual sugere que a fase espera alguém, e ninguém virá.
func TestBoardDistingueFaseDeCard(t *testing.T) {
	b, err := fs.ReadFile(boardFS, "board/anchors-board.html")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// A MARCA vem de não ter card (`!l`), e não de uma lista de códigos: uma fase que
	// ganhasse card deixaria de ser estrutura, e a regra tem de acompanhar.
	if !strings.Contains(texto, "const ehFase = !l;") {
		t.Error("a fase é reconhecida por outra coisa que não a ausência de card — se " +
			"uma fase ganhar issue, ela passa a ser trabalho e o desenho tem de seguir")
	}

	// TEXTURA, e não só cor: a distinção precisa sobreviver a quem não distingue cores.
	if !strings.Contains(texto, "repeating-linear-gradient(45deg") {
		t.Error("a barra da fase não tem hachura — cor sozinha não distingue para quem " +
			"não a enxerga, e a textura é o que separa à distância")
	}
	// A COR fora da paleta de ESTADO: a fase não está em coluna nenhuma, e usar cinza,
	// azul ou verde a poria numa.
	if !strings.Contains(texto, "--fase:") {
		t.Error("a fase usa uma cor de estado — ela não tem estado, e isso a poria numa " +
			"coluna onde ela não está")
	}
	// Nos TRÊS blocos de tema: uma cor definida só no claro some no escuro.
	if n := strings.Count(texto, "--fase:"); n < 3 {
		t.Errorf("`--fase` declarada %d vez(es) — faltam os blocos de tema escuro", n)
	}
}

// O STALE LIBERA O DONO, e não move o card.
//
// A versão anterior decidia um destino pela label anterior: `in-review` voltava para
// `ready-to-review`, e todo o resto para `to-do`. O raciocínio do primeiro caso estava certo
// e escrito no próprio código — "mandar para `to-do` faria o próximo agente REIMPLEMENTAR o
// que já existe".
//
// O SEGUNDO CASO PRESUMIA o que não conferia. A justificativa dizia "aqui o trabalho NÃO
// terminou: não há PR, ou os checks não passaram" — mas nada olhava se havia PR. Quando
// havia, e verde, a presunção estava errada e o efeito era exatamente o que o comentário
// condenava três linhas acima: o card voltava para `to-do` e o trabalho pronto se perdia.
//
// Liberar o dono BASTA: o `claim` serve por estado E por dono liberado.
func TestStaleLiberaODonoSemMoverOCard(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-stale.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	if strings.Contains(texto, "--add-label") {
		t.Error("o stale move o card — mover é uma decisão tomada sem os dados que a " +
			"justificariam, e um card com PR verde voltaria para `to-do`")
	}
	// A LIBERAÇÃO continua: sem ela o card fica preso a um dono que sumiu.
	if !strings.Contains(texto, "anchors-owner: (liberado)") {
		t.Fatal("o stale deixou de liberar o dono — o card ficaria preso a quem sumiu, e " +
			"é para isso que este pipeline existe")
	}
	// E o REGISTRO diz onde o card ficou: sem isso, quem lê precisa conferir as labels
	// para descobrir o que aconteceu.
	if !strings.Contains(texto, "estado_atual") {
		t.Error("o registro da liberação não diz em que coluna o card ficou")
	}
}

// O PR CUJO CARD ESTÁ EM `to-do` REPROVA — o claim foi pulado.
//
// Um PR aberto significa que alguém está trabalhando, e o card precisa dizer isso. Em
// `to-do`, ninguém registrou posse: o board mostra como disponível um trabalho que já tem PR.
//
// MEDIDO no projeto de referência: 38 dos 44 PRs abertos tinham o card em `to-do`. O efeito
// não é cosmético — o `claim` prioriza `ready-to-review` sobre `to-do`, e encontrava a fila
// de revisão VAZIA enquanto 46 trabalhos esperavam. Cada agente pegava trabalho novo em vez
// de revisar o que estava pronto.
//
// REPROVA, e não avisa. O `pr-checks` já avisava ("card não está em trabalho ativo — nada a
// fazer"), e o aviso não impediu as 38 ocorrências. Um gate que acusa e deixa passar ensina
// que o passo é opcional.
func TestGateReprovaCardSemClaim(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-gates.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// O ACÚMULO: sem ele o gate confere o último card e esquece os anteriores — um PR que
	// fecha três cards passaria com dois deles em `to-do`.
	if !strings.Contains(texto, `semClaim="$semClaim #$c"`) {
		t.Error("o gate não ACUMULA os cards sem claim — num PR que fecha vários, só o " +
			"último seria conferido")
	}
	// E a JURISDIÇÃO: sem confirmar a label-raiz, o gate reprovaria uma issue comum do
	// projeto que alguém citou num `Closes`.
	if !strings.Contains(texto, `index("anchors")`) {
		t.Error("o gate lê o card sem confirmar que ele é do Anchors — uma issue de outro " +
			"fluxo citada num `Closes` seria reprovada por não ter label de estado")
	}
	if !strings.Contains(texto, "semClaim=") {
		t.Fatal("o gate não confere a coluna do card — um PR com o card em `to-do` passa, " +
			"e o board segue mostrando como disponível um trabalho que já tem PR")
	}
	// REPROVA de verdade: `::error::` sem `exit 1` é aviso com cara de erro.
	i := strings.Index(texto, `if [ -n "$semClaim" ]`)
	if i < 0 || !strings.Contains(texto[i:min(len(texto), i+2200)], "exit 1") {
		t.Error("o gate acusa e deixa passar — foi o que já acontecia no `pr-checks`, e " +
			"não impediu 38 ocorrências")
	}

	// OS ESTADOS QUE PASSAM, e cada um por uma razão diferente:
	//
	//   · trabalho em curso — o card diz que alguém está nele
	//   · de `ready-to-test` em diante — o Anchors saiu da alçada, e um PR que corrige
	//     algo já entregue é legítimo
	//   · sem label de estado — não é card do fluxo, e a régua não tem jurisdição
	for _, e := range []string{"in-progress", "ready-to-review", "in-review"} {
		if !strings.Contains(texto[i-1500:i], e) {
			t.Errorf("o estado %q não está entre os que passam — um card legitimamente "+
				"em curso seria reprovado", e)
		}
	}
	if !strings.Contains(texto[i-1500:i], "ready-to-test") {
		t.Error("`ready-to-test` em diante precisa passar: o Anchors saiu da alçada, e " +
			"cobrar claim ali seria cobrar duas vezes pelo mesmo trabalho")
	}

	// A SAÍDA na mensagem: um gate que reprova sem dizer o que fazer transfere o problema.
	//
	// Esta régua nasceu ERRADA: ela casava o texto `anchors claim`, e com isso EXIGIA
	// que a mensagem ensinasse um comando inexistente. Uma régua que fixa o texto da
	// instrução prova que a mensagem não mudou — não que ela funciona. Agora ela cobra
	// a posse pelo comando real, e o `TestPipelineSoEnsinaComandoQueExiste` confronta
	// todo `anchors <algo>` dos pipelines contra os comandos registrados no cobra.
	if !strings.Contains(texto, "anchors next") {
		t.Error("a mensagem não diz como destravar — quem a lê fica sabendo que errou e " +
			"não o que fazer")
	}
}

// O card PARADO em espera e o card PARADO em andamento nao sao o mesmo problema,
// e cobra-los com o mesmo relogio erra nos dois lados.
//
// Em `to-do` e `ready-to-*` o card esta parado POR DEFINICAO -- ele espera alguem
// pegar. Ter dono ali e' residuo: alguem reivindicou e nao trabalhou, e enquanto o
// nome estiver la o `next` pode pular o card. Duas horas bastam, e o custo de errar
// e' zero: o card ja estava disponivel.
//
// Em `in-progress` e `in-review` ha trabalho ACONTECENDO. Liberar cedo demais rouba
// o card de quem esta trabalhando devagar -- e o novo dono comeca do zero sobre algo
// meio feito. Vinte e quatro horas, porque o sinal de vida (`updatedAt`) so aparece
// quando o agente comenta, move label ou referencia commit: um dev humano pode passar
// um dia inteiro num card sem tocar no issue.
func TestStaleSeparaEsperaDeAndamento(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-stale.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// Os dois relogios existem e sao distintos.
	if !strings.Contains(texto, "HORAS_ESPERA") {
		t.Error("falta `HORAS_ESPERA` — sem ele o card com dono em `to-do`/`ready-to-*` " +
			"espera o mesmo tempo do trabalho em andamento, e o `next` pula um card livre")
	}
	if !strings.Contains(texto, "HORAS_ANDAMENTO") {
		t.Error("falta `HORAS_ANDAMENTO` — sem ele o trabalho em curso e a espera " +
			"compartilham relogio, e um dos dois fica errado")
	}
	if strings.Contains(texto, "HORAS_ATE_STALE") {
		t.Error("`HORAS_ATE_STALE` ainda existe — o relogio unico e' exatamente o que " +
			"esta separacao desfaz")
	}

	// O de espera precisa ser MENOR: e' o caso barato de errar.
	espera := valorDeEnv(texto, "HORAS_ESPERA")
	andamento := valorDeEnv(texto, "HORAS_ANDAMENTO")
	if espera == "" || andamento == "" {
		t.Fatalf("nao consegui ler os dois relogios (espera=%q andamento=%q)", espera, andamento)
	}
	if espera >= andamento {
		t.Errorf("espera=%s e andamento=%s — a espera tem que ser MENOR: liberar cedo "+
			"um card que ja estava disponivel nao custa nada, e liberar cedo trabalho "+
			"em curso rouba o card de quem o faz", espera, andamento)
	}

	// `ready-to-*` entra na varredura de espera. Antes ele era EXCLUIDO de tudo, e a
	// razao valia contra MOVER -- nao contra liberar o dono.
	if !strings.Contains(texto, "ready-to-") {
		t.Error("`ready-to-*` nao aparece — o card pronto com dono residual fica preso, " +
			"e era justamente o que a separacao vem resolver")
	}

	// O que a correcao anterior conquistou nao pode voltar: liberar o dono, nunca mover.
	if strings.Contains(texto, "add-label") {
		t.Error("`add-label` de volta no stale — o pipeline voltou a MOVER card, e mover " +
			"apaga o trabalho ja feito")
	}
}

// valorDeEnv le `NOME: "valor"` do bloco `env:` do YAML.
func valorDeEnv(texto, nome string) string {
	for _, linha := range strings.Split(texto, "\n") {
		l := strings.TrimSpace(linha)
		if !strings.HasPrefix(l, nome+":") {
			continue
		}
		v := strings.TrimSpace(strings.TrimPrefix(l, nome+":"))
		return strings.Trim(v, `"'`)
	}
	return ""
}

// TODO ISSUE NASCE COM ESTADO, e o pipeline garante isso ao ve-la.
//
// MEDIDO no app de referência: 59 de 95 cards abertos estavam FORA do fluxo -- 53 sem label de
// estado nenhuma, 6 so com `anchors:under-N` (que e' bloqueio, nao estado). O `claim` so
// varre `to-do` e `ready-to-review`, entao esse trabalho era invisivel: tres agentes
// pediram card, ouviram "nenhum card livre", e pararam.
//
// A causa: os achados de gate (`[domain-declared] Violacao @ ...`) sao criados por
// caminhos diferentes -- `anchors judge`, `anchors escalate`, o reporter dos gates -- e
// nenhum aplicava estado.
//
// POR QUE NO PIPELINE, e nao em cada criador: consertar um por um deixa os outros, e o
// proximo caminho de criacao nasce quebrado de novo. Nao ha regua que cubra "todo lugar
// que cria issue". O pipeline ve a issue nascer, venha de onde vier.
//
// E' o mesmo raciocinio do `stale` e do proprio guard: o pipeline corrige o FATO, em vez
// de depender de cada um lembrar.
func TestGuardDaEstadoAQuemNasceSemEle(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	texto := string(b)

	// Ele precisa VER a issue nascer.
	//
	// A regua nasceu FRACA aqui: `strings.Contains(texto, "opened")` casava dentro de
	// `reopened`, e passava sem a correcao existir. E' a mesma armadilha que os agentes
	// vem achando nos testes deste projeto -- asserção por substring nao distingue o que
	// mede. O gatilho precisa estar na LISTA de types.
	if !regexp.MustCompile(`types:\s*\[[^]]*\bopened\b`).MatchString(texto) {
		t.Error("o guard nao escuta `opened` na lista de types — uma issue criada sem " +
			"estado passa sem ninguem ver, e some do `claim`")
	}

	// E precisa APLICAR o estado inicial.
	if !strings.Contains(texto, "anchors:to-do") {
		t.Error("o guard nao aplica `anchors:to-do` — ver a issue nascer sem agir nao " +
			"resolve nada")
	}

	// O QUE JA TEM ESTADO nao pode ser tocado: um card em `in-review` que recebesse
	// `to-do` voltaria para a fila e seria reimplementado do zero.
	if !strings.Contains(texto, "anchors:ready-to-review") &&
		!strings.Contains(texto, "startswith(\\\"anchors:\\\")") {
		t.Error("o guard nao confere se JA HA estado antes de aplicar — poria `to-do` " +
			"num card que ja esta em revisao, e o trabalho seria refeito")
	}
}

// REABRIR SEM DECLARAÇÃO É CORREÇÃO, NÃO VIOLAÇÃO.
//
// A trava era simétrica — fechou à mão, reabre; reabriu à mão, fecha —, e a simetria valia
// enquanto fechar significasse "trabalho encerrado". Não significa mais: o card fica aberto
// em `ready-to-test` até alguém terminar a esteira (ver TestOCardNaoFechaNoMerge).
//
// O passivo torna isto obrigatório, não opcional: MEDIDO no projeto de referência, 46 cards
// fecharam sem `anchors:manual` nem `anchors:discarded` — por efeito do `Closes`, não por
// decisão de ninguém. Com a trava simétrica, reabrir esses 46 seria desfeito em segundos, e
// a trava estaria defendendo uma regra revogada exatamente quando alguém tenta consertar o
// que ela ajudou a criar.
//
// O QUE NÃO PODE SUMIR é o outro lado: reabrir um card DECLARADO encerrado continua sendo a
// mexida que a trava reverte.
func TestReaberturaSemDeclaracaoEhMantida(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	i := strings.Index(s, "reopened)")
	if i < 0 {
		t.Fatal("o pipeline não trata `reopened` — a trava perdeu um dos quatro eventos")
	}
	// Só o ramo do `reopened`, e não o arquivo inteiro: o `gh issue close` aparece noutros
	// lugares, e uma busca global passaria mesmo se este ramo fechasse sempre.
	fim := strings.Index(s[i:], "\n            labeled)")
	if fim < 0 {
		t.Fatal("não achei o fim do ramo `reopened` — o `case` mudou de forma")
	}
	ramo := s[i : i+fim]

	// A DECISÃO tem de ser pela declaração, dentro do ramo.
	if !strings.Contains(ramo, "anchors:manual") || !strings.Contains(ramo, "anchors:discarded") {
		t.Error("o ramo `reopened` fecha de volta sem conferir se o encerramento foi " +
			"DECLARADO — assim ele desfaz a reabertura de um card que o `Closes` fechou")
	}
	// E o caminho sem declaração tem de SAIR sem fechar.
	if !strings.Contains(ramo, "exit 0") {
		t.Error("o ramo `reopened` não tem saída sem fechar — a reabertura legítima " +
			"precisa ser mantida, e não só comentada")
	}
	// O `break` não sai de um `case` em sh — sai de laço. Usá-lo aqui deixaria o ramo
	// escorrendo para o código comum de reversão.
	if strings.Contains(ramo, "\n                break\n") {
		t.Error("`break` não encerra um ramo de `case` em sh: o fluxo escorre para a " +
			"reversão comum e o card é fechado mesmo no caminho que devia mantê-lo aberto")
	}
}

// A PROCEDÊNCIA É CONFERIDA NO PIPELINE, porque é o único ponto por onde TODO card passa.
//
// Não há hook a pôr no `gh`: ele é um cliente HTTP e não tem `pre-issue-create`. E se
// tivesse não bastaria — `gh api -X POST`, `curl` e a interface web criam a issue pelo mesmo
// endpoint sem passar por ele. O `anchors escalate` valida e protege só quem o usa.
//
// O QUE ESTE PASSO FAZ é acusar, não impedir: o GitHub não tem "recusar issue". O card
// nasce, e em segundos ganha um comentário dizendo o que falta — antes de alguém pegá-lo.
//
// MEDIDO no projeto de referência: 18 de 27 decisões abertas sem amarração alguma. Todas
// traziam o arquivo (`**Onde:**`) e nenhuma o card de origem; uma escreveu "o achado veio do
// #690" em PROSA, que não se consulta.
func TestGuardCobraProcedenciaDoCardEscalado(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	// SÓ O CARD ESCALADO: um card de trabalho vem do plano, não de outro card, e cobrar
	// procedência dele seria exigir um vínculo que não existe.
	if !strings.Contains(s, "O card escalado nasce com procedência") {
		t.Fatal("o guard não confere procedência de card escalado — o achado nasce solto e " +
			"ninguém consegue perguntar o que aquele trabalho gerou")
	}

	// A GRAFIA É UMA SÓ. O projeto está em beta fechado — não há board de terceiro com a
	// grafia anterior a respeitar, e aceitar duas formas da mesma coisa é superfície de
	// divergência: a primeira vez que uma delas mudasse, a outra ficaria para trás.
	if !strings.Contains(s, PrefixoLabelSob) {
		t.Errorf("a validação de procedência não conhece %q", PrefixoLabelSob)
	}

	// O REMÉDIO junto da acusação. Um aviso que diz o que está errado e não como consertar
	// transfere o trabalho de descobrir para quem já errou.
	if !strings.Contains(s, "--add-label anchors:under-") {
		t.Error("o aviso não diz COMO pôr a procedência — quem o lê tem de ir descobrir o " +
			"nome da label, e a maioria não vai")
	}

	// E A SAÍDA LEGÍTIMA declarada: escalar algo visto de fora de um trabalho é válido, e
	// um aviso sem essa ressalva ensina que o card está errado quando não está.
	if !strings.Contains(s, "nasceu fora de um trabalho") {
		t.Error("o aviso não reconhece o achado que nasce fora de um card — sem isso ele " +
			"acusa de defeito o que é uso legítimo")
	}
}

// O `gates` NÃO ENGOLE O CÓDIGO DE SAÍDA do `doctor --check-pipelines`.
//
// Havia ali um `if ! anchors doctor --check-pipelines; then echo ::warning::; fi`, e o `if`
// ANULAVA a decisão do projeto: ele consome o código de saída, então o passo passava verde
// mesmo com `stale_pipeline_blocks: true` declarado no `anchors.yaml`.
//
// O comentário logo acima dizia "quem decide se isto barra é o PROJETO" — e o código abaixo
// tirava a decisão dele. O `doctor` já implementava a parte certa: sai com 1 quando o
// projeto declarou, e imprime o aviso ele mesmo quando não declarou.
//
// MEDIDO (#798): mutação aplicada ao próprio pipeline — 6 dos 10 passos de verificação
// somem sem que nenhum sinal acuse. Os quatro deste arquivo eram o único caso detectado, e
// por este comando; mas o aviso saía de DENTRO do arquivo verificado, então apagar o passo
// apagava o aviso junto.
//
// O BOARD E O IDENTIFY FICAM COM O AVISO, e a diferença é o que o passo faz: eles publicam
// e movem card, e reprová-los trocaria "o board atualiza com o desenho antigo" por "o board
// não atualiza". Quem decide se o trabalho é aceito barra; quem informa avisa.
func TestGatesNaoEngoleOVereditoDoCheckPipelines(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-gates.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	if !strings.Contains(s, "anchors doctor --check-pipelines") {
		t.Fatal("o `gates` não confere os pipelines — o verificador deixa de se verificar")
	}
	// O `if !` é o que anula: ele consome o código de saída e o passo passa verde.
	if strings.Contains(s, "if ! anchors doctor --check-pipelines") {
		t.Error("o `gates` voltou a envolver o `check-pipelines` num `if !` — o `if` " +
			"consome o código de saída, e o `stale_pipeline_blocks: true` do projeto " +
			"deixa de ter efeito")
	}
	// E `continue-on-error` faria o mesmo por outro caminho.
	iCheck := strings.Index(s, "anchors doctor --check-pipelines")
	trecho := s[maxInt(0, iCheck-600):iCheck]
	if strings.Contains(trecho, "continue-on-error: true") {
		t.Error("o passo do `check-pipelines` ganhou `continue-on-error` — mesmo efeito do " +
			"`if !`: o veredito do projeto não barra nada")
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// O `concurrency` QUE ESCREVE NO CARD É POR PR, não global.
//
// A razão de serializar é a escrita: dois eventos do MESMO PR (o `check_suite` e o
// `synchronize`) podem chegar juntos e mover o mesmo card duas vezes. O grupo por PR
// resolve isso — e PRs diferentes mexem em cards diferentes.
//
// O GRUPO GLOBAL CUSTAVA O WORKFLOW INTEIRO. O GitHub mantém 1 rodando + 1 pendente por
// grupo, e um terceiro run enfileirado CANCELA o pendente. Com fila grande, quase todo PR
// perde a corrida.
//
// E o cancelamento é PIOR que falha: um job que falha aparece vermelho; um job cancelado
// antes de criar o check run não aparece de forma alguma — o rollup fica idêntico ao de um
// PR onde o workflow nunca existiu.
//
// MEDIDO no projeto de referência (#821): com 19 PRs abertos, 10 de 10 conferidos estavam
// SEM o check do `pr-checks`. E ele é quem move o card para `ready-to-test` no merge.
func TestConcurrencyDoPRChecksEhPorPR(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-pr-checks.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)

	i := strings.Index(s, "concurrency:")
	if i < 0 {
		t.Fatal("o `pr-checks` não declara `concurrency` — dois eventos do mesmo PR " +
			"escreveriam estado em cima de estado no card")
	}
	fim := strings.Index(s[i:], "\npermissions:")
	if fim < 0 {
		fim = len(s) - i
	}
	bloco := s[i : i+fim]

	// O GRUPO tem de variar por PR. Um literal fixo serializa o repositório inteiro.
	if !strings.Contains(bloco, "github.event.pull_request.number") {
		t.Error("o `group` do `pr-checks` não varia por PR — com fila grande o GitHub " +
			"cancela os pendentes, e o workflow SOME do rollup em vez de falhar")
	}

	// E O FALLBACK importa: este workflow também roda por `check_suite` e por
	// `workflow_dispatch`, onde `pull_request.number` é vazio. Sem fallback, todos esses
	// eventos cairiam no MESMO grupo (o de sufixo vazio) — o global de volta, por outro
	// caminho.
	if !strings.Contains(bloco, "||") {
		t.Error("o `group` não tem fallback para os eventos sem `pull_request.number` " +
			"(`check_suite`, `workflow_dispatch`) — eles voltariam a compartilhar um grupo")
	}
}

// The `to-do` AT BIRTH IS NOT TAMPERING — and the lock removed it.
//
// Whoever creates the card already with `anchors:to-do` (by hand, through `anchors
// escalate`, through the web) fires `labeled` as a human. The `reverter`'s `labeled`
// branch removed the label twenty seconds later, and `estado-inicial` did not fix it:
// when it looked, the card ALREADY had a state. The two jobs ran together and cancelled
// each other out.
//
// MEASURED in the reference project: 21 open cards with no state, all with the same trail
// — born with `to-do`, `github-actions[bot]` removes it in ~20s. Invisible to `claim`.
//
// The ruler looks ONLY at the `labeled` branch: `anchors:to-do` appears all over the file
// (`estado-inicial` applies it), and a global search would pass with no exception there.
func TestLockDoesNotRemoveTheToDoOfANewCard(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	i := strings.Index(s, "\n            labeled)")
	if i < 0 {
		t.Fatal("the `reverter` does not handle `labeled` — the lock lost one of the events")
	}
	end := strings.Index(s[i:], "\n            unlabeled)")
	if end < 0 {
		t.Fatal("could not find the end of the `labeled` branch — the `case` changed shape")
	}
	branch := s[i : i+end]

	removal := strings.Index(branch, "--remove-label")
	if removal < 0 {
		t.Fatal("the `labeled` branch no longer removes anything — the lock stopped reverting")
	}
	before := branch[:removal]

	// The exception must come BEFORE the removal, and exit without removing.
	if !strings.Contains(before, `"$LABEL_MEXIDA" = "anchors:to-do"`) {
		t.Error("the `labeled` branch removes the `to-do` without asking whether the card is " +
			"being BORN — a card created already in the queue loses its state and vanishes from `claim`")
	}
	if !strings.Contains(before, "exit 0") {
		t.Error("the birth exception does not exit before the removal — the `to-do` is removed " +
			"even when recognised as an entrance")
	}

	// And it only holds for a card with NO OTHER STATE: moving an `in-progress` back to
	// the queue by hand is still tampering.
	for _, state := range []string{"anchors:in-progress", "anchors:ready-to-review",
		"anchors:in-review", "anchors:ready-to-test"} {
		if !strings.Contains(before, state) {
			t.Errorf("the exception does not check %q — a card in that state receiving `to-do` "+
				"by hand would go back to the queue and be reimplemented", state)
		}
	}
}

// The `anchors` LABEL ARRIVES AFTER CREATION, and `estado-inicial` has to see it.
//
// MEASURED in the reference project: #912 was born at 18:39:23 and got `anchors` at
// 18:40:06. When the job ran on `opened`, the card did not belong to Anchors yet, and it
// left without acting. The card stayed stateless, invisible to `claim`, with no revert in
// the trail to explain why.
//
// The ruler reads the JOB's `if`, parsed: a search for `labeled` in the file would match
// the workflow trigger or the `reverter` branch, and pass with no fix in place.
func TestInitialStateSeesTheLabelThatArrivesLater(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Jobs map[string]struct {
			If string `yaml:"if"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	job, ok := doc.Jobs["estado-inicial"]
	if !ok {
		t.Fatal("the `estado-inicial` job vanished from the guard")
	}
	cond := strings.Join(strings.Fields(job.If), " ")
	if !strings.Contains(cond, "github.event.action == 'opened'") {
		t.Error("`estado-inicial` no longer sees the issue being BORN")
	}
	if !strings.Contains(cond, "github.event.action == 'labeled'") ||
		!strings.Contains(cond, "github.event.label.name == 'anchors'") {
		t.Errorf("`estado-inicial` does not react to the `anchors` label arriving after "+
			"creation — the card stays stateless and vanishes from `claim`. if: %q", cond)
	}
}

// ONE PROVENANCE WARNING PER CARD. `estado-inicial` runs on `opened` and on the late
// `anchors` label; a card created with its labels fires both, and the warning was posted
// twice. The step must look for its own earlier comment BEFORE commenting.
func TestGuardProvenanceWarningIsPostedOnce(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-guard.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	i := strings.Index(s, "- name: O card escalado nasce com procedência")
	if i < 0 {
		t.Fatal("the provenance step vanished from the guard")
	}
	step := s[i:]
	if j := strings.Index(step[1:], "\n      - name:"); j >= 0 {
		step = step[:j+1]
	}
	comment := strings.Index(step, "gh issue comment")
	if comment < 0 {
		t.Fatal("the provenance step no longer comments")
	}
	before := step[:comment]
	marker := "⚠ **Este card não diz de qual trabalho nasceu.**"
	if !strings.Contains(before, "select(startswith(\""+marker+"\"))") {
		t.Error("the step comments without checking for its own earlier warning — a card " +
			"created with its labels gets the warning twice")
	}
	if !strings.Contains(step[comment:], marker) {
		t.Error("the check looks for a marker the comment no longer starts with")
	}
	if !strings.Contains(before, "exit 0") {
		t.Error("finding the earlier warning does not stop the step")
	}
}

// --- the claim script, RUN against a fake `gh` ---
//
// Greps over the YAML cannot tell a regex that matches `4/4` from one that does not, nor
// an error swallowed by `|| true` from one that stops the card. These tests run the
// claim's real `run:` script with a `gh` that answers from environment variables, and
// look at what it did: whether it commented `anchors-owner:` on the card.

// fakeGH answers the claim script's `gh` calls. One to-do card ($FAKE_CARD) is on the
// board; $FAKE_FAIL names the one call that fails like an API error; $FAKE_PRS is the
// JSON of the open PRs, filtered by the script's own `--jq` through the real `jq`.
const fakeGH = `#!/usr/bin/env bash
echo "gh $*" >> "$FAKE_LOG"
call="$1 $2"
json=""; jqx=""; labels=""
while [ $# -gt 0 ]; do
  case "$1" in
    --json) json="$2"; shift ;;
    --jq) jqx="$2"; shift ;;
    --label) labels="$labels $2"; shift ;;
  esac
  shift
done
fail() { if [ "${FAKE_FAIL:-}" = "$1" ]; then echo "gh: HTTP 502" >&2; exit 1; fi; }
case "$call" in
  "api graphql")
    if [ -n "${FAKE_GRAPHQL:-}" ]; then printf '%s' "$FAKE_GRAPHQL" | jq -r "$jqx"
    else printf '%s' "${FAKE_MINE:-}"; fi ;;
  "api "*) exit 1 ;;
  "issue list")
    case "$labels" in
      *anchors:desbloqueia-*) fail unblock ;;
      *" anchors:${FAKE_CARD_STATE:-to-do}"*) echo "$FAKE_CARD" ;;
    esac ;;
  "issue view")
    case "$json" in
      labels) case "$jqx" in
          *blocked-by*) fail labels; [ -z "${FAKE_BLOCKED_BY:-}" ] || echo "$FAKE_BLOCKED_BY" ;;
          *) echo "anchors:${FAKE_CARD_STATE:-to-do}" ;;
        esac ;;
      state) fail state; echo "${FAKE_STATE:-CLOSED}" ;;
      title) fail title ;;
      comments) case "$jqx" in
          *"## Revisão"*) echo 0 ;;
          *"implementação concluída"*) echo "${FAKE_AUTHOR:-}" ;;
          *'select(. != "(liberado)")'*) echo "${FAKE_PREV_OWNER:-}" ;;
          *) echo "${FAKE_OWNER:-}" ;;
        esac ;;
    esac ;;
  "pr list") fail pr; printf '%s' "${FAKE_PRS:-[]}" | jq -r "$jqx" ;;
esac
exit 0
`

// runClaim runs the claim script with the fake `gh` and the given environment, and
// returns its output and whether it handed card #4 out.
func runClaim(t *testing.T, env ...string) (string, bool) {
	t.Helper()
	for _, bin := range []string{"bash", "jq"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("no %s on PATH", bin)
		}
	}
	b, err := workflowsFS.ReadFile("workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Jobs map[string]struct {
			Steps []struct {
				Run string `yaml:"run"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	script := ""
	for _, j := range doc.Jobs {
		for _, s := range j.Steps {
			if strings.Contains(s.Run, "anchors-owner: $AGENT") {
				script = s.Run
			}
		}
	}
	if script == "" {
		t.Fatal("the claim step (the one that comments `anchors-owner: $AGENT`) vanished")
	}
	dir := t.TempDir()
	// The script keeps its scratch files in /tmp; each run gets its own.
	script = strings.ReplaceAll(script, "/tmp/", dir+"/")
	bin := filepath.Join(dir, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	testkit.FakeBin(t, bin, "gh", fakeGH)
	logFile := filepath.Join(dir, "gh.log")
	cmd := pipelineBash(t, "-c", script)
	cmd.Env = append(os.Environ(),
		"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"FAKE_LOG="+logFile, "FAKE_CARD=4",
		"GH_REPO=o/r", "AGENT=machine/session", "LABEL=anchors")
	cmd.Env = append(cmd.Env, env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("the claim script failed: %v\n%s", err, out)
	}
	calls, _ := os.ReadFile(logFile)
	return string(out), strings.Contains(string(calls), "gh issue comment 4 --body anchors-owner: machine/session")
}

// Only a REAL reference to the card hides it: a linking keyword before `#N`, or an issue
// URL. Measured in the reference app: a PR body reporting a mutation score of `4/4` matched
// the old `(#|/)N\b`, and card #4 was never handed out again.
func TestClaimOnlyCountsARealReferenceToTheCard(t *testing.T) {
	pr := func(body string) string {
		j, _ := json.Marshal([]map[string]any{{"number": 77, "body": body}})
		return "FAKE_PRS=" + string(j)
	}
	for _, body := range []string{
		"Closes #4", "closes #4", "Refs #4", "Fixes: #4", "RESOLVES #4", "fixed #4",
		"see https://github.com/o/r/issues/4 for context",
	} {
		out, handed := runClaim(t, pr(body))
		if handed {
			t.Errorf("PR body %q references card #4, and the card was still handed out:\n%s", body, out)
		}
		if !strings.Contains(out, "já tem o PR #77") {
			t.Errorf("PR body %q: the skip does not name the PR:\n%s", body, out)
		}
	}
	for _, body := range []string{
		"mutation score 4/4", "Closes #40", "Refs #14", "step #4 of the plan",
		"https://github.com/o/r/pull/4", "https://github.com/o/r/issues/44",
	} {
		if out, handed := runClaim(t, pr(body)); !handed {
			t.Errorf("PR body %q does not reference card #4, and it hid the card:\n%s", body, out)
		}
	}
}

// The guards FAIL CLOSED: when `gh` errors, the card is skipped and the log says so. The
// old `2>/dev/null || true` read an API error as "no blocker" and handed the card out.
func TestClaimGuardsFailClosed(t *testing.T) {
	// The control: with every call answering, the card IS handed out — otherwise the
	// cases below would pass for a reason that has nothing to do with the failure.
	if out, handed := runClaim(t); !handed {
		t.Fatalf("with no failure the card should be handed out:\n%s", out)
	}
	for _, c := range []struct{ name, fail, extra string }{
		{"unblock card lookup", "unblock", ""},
		{"blocked-by labels", "labels", ""},
		{"blocker state", "state", "FAKE_BLOCKED_BY=anchors:blocked-by-9"},
		{"open PR lookup", "pr", ""},
		{"card title (dependencies)", "title", ""},
	} {
		env := []string{"FAKE_FAIL=" + c.fail}
		if c.extra != "" {
			env = append(env, c.extra)
		}
		out, handed := runClaim(t, env...)
		if handed {
			t.Errorf("%s: gh failed and the card was handed out anyway (fail open):\n%s", c.name, out)
		}
		if !strings.Contains(out, "#4") || !(strings.Contains(out, "fail closed") ||
			strings.Contains(out, "gh failed")) {
			t.Errorf("%s: the skip is not logged — nobody reading the run learns why:\n%s", c.name, out)
		}
	}
}

// The claim run carries the agent in its title, exactly as `board.ClaimRunTitle` writes it:
// that title is how `anchors next` finds the run it dispatched and waits for it instead of
// dispatching another.
func TestClaimRunNameIsWhatNextLooksFor(t *testing.T) {
	b, err := workflowsFS.ReadFile("workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		RunName string `yaml:"run-name"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	if want := board.ClaimRunTitle("${{ inputs.agent }}"); doc.RunName != want {
		t.Errorf("run-name = %q, want %q — `anchors next` would not recognise its own claim",
			doc.RunName, want)
	}
}

// workflowStep is one parsed step of an embedded template, as the tests below read it.
type workflowStep struct {
	Name string            `yaml:"name"`
	ID   string            `yaml:"id"`
	If   string            `yaml:"if"`
	Uses string            `yaml:"uses"`
	With map[string]string `yaml:"with"`
	Env  map[string]string `yaml:"env"`
	Run  string            `yaml:"run"`
}

type workflowJob struct {
	Env   map[string]string `yaml:"env"`
	Steps []workflowStep    `yaml:"steps"`
}

func parseWorkflow(t *testing.T, name string) map[string]workflowJob {
	t.Helper()
	b, err := fs.ReadFile(workflowsFS, "workflows/"+name)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Jobs map[string]workflowJob `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return doc.Jobs
}

// writeStubs puts executables named after the map keys in a fresh directory, and returns
// a PATH that finds them before the real ones.
func writeStubs(t *testing.T, stubs map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range stubs {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte("#!/usr/bin/env bash\n"+body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return dir + string(os.PathListSeparator) + os.Getenv("PATH")
}

// EVERY TEMPLATE THAT INSTALLS ANCHORS PICKS THE VERSION THE SAME WAY.
//
// The resolver read `min_version` from anchors.yaml; `gates`, `board` and `identify` used
// `ANCHORS_VERSION: latest`. A project pinned to a version had three of its pipelines jump
// to a newer binary on their own while the fourth stayed put — the same repository judged
// by two different Anchors depending on which workflow ran.
//
// The ruler compares the install SCRIPTS byte for byte: one copy that drifts (a new
// template written from an old one, a fix applied to three of four) fails here.
func TestEveryTemplateResolvesTheAnchorsVersionTheSameWay(t *testing.T) {
	entries, err := fs.ReadDir(workflowsFS, "workflows")
	if err != nil {
		t.Fatal(err)
	}
	var reference, referenceFile string
	installers := 0
	for _, e := range entries {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+e.Name())
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(b), "ANCHORS_VERSION") {
			t.Errorf("%s still declares ANCHORS_VERSION — the version comes from `min_version`", e.Name())
		}
		if !strings.Contains(string(b), "gh release download") {
			continue
		}
		installers++
		var script string
		for _, job := range parseWorkflow(t, e.Name()) {
			for _, s := range job.Steps {
				if strings.Contains(s.Run, "gh release download") {
					if s.Name != "Instalar o Anchors" {
						t.Errorf("%s downloads Anchors in step %q, outside the shared install step", e.Name(), s.Name)
					}
					script = s.Run
				}
			}
		}
		if !strings.Contains(script, "min_version") {
			t.Errorf("%s installs Anchors without reading `min_version`", e.Name())
		}
		if reference == "" {
			reference, referenceFile = script, e.Name()
			continue
		}
		if script != reference {
			t.Errorf("%s installs Anchors with a script different from %s's — the version "+
				"rule has to be the same snippet everywhere", e.Name(), referenceFile)
		}
	}
	if installers < 4 {
		t.Errorf("only %d templates install Anchors; expected gates, board, identify and "+
			"resolve-queue — the glob or the marker changed and this ruler would pass empty", installers)
	}
}

// THE SNIPPET, RUN: `min_version` wins in every spelling the config accepts, and `latest`
// only when the field is absent.
//
// `gh`, `tar`, `sudo` and `anchors` are stubs; what is measured is the TAG handed to
// `gh release download`.
func TestInstallSnippetDownloadsTheMinVersion(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("no bash on PATH")
	}
	var script string
	for _, job := range parseWorkflow(t, "anchors-gates.yml") {
		for _, s := range job.Steps {
			if s.Name == "Instalar o Anchors" {
				script = s.Run
			}
		}
	}
	if script == "" {
		t.Fatal("anchors-gates.yml has no `Instalar o Anchors` step")
	}
	cases := []struct {
		name, yaml string // yaml == "-" means no anchors.yaml at all
		want       string
	}{
		{"plain", "version: 4\nmin_version: 0.1.84\n", "v0.1.84"},
		{"double-quoted", "min_version: \"0.1.84\"\n", "v0.1.84"},
		{"single-quoted", "min_version: '0.1.84'\n", "v0.1.84"},
		{"with a v", "min_version: v0.1.84\n", "v0.1.84"},
		{"trailing comment", "min_version: 0.1.84 # raised for 0.1.80's fix\n", "v0.1.84"},
		{"absent", "version: 4\n", "v9.9.9"},
		{"commented out", "# min_version: 0.1.84\n", "v9.9.9"},
		{"no anchors.yaml", "-", "v9.9.9"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			got := filepath.Join(dir, "downloaded")
			if c.yaml != "-" {
				if err := os.WriteFile(filepath.Join(dir, "anchors.yaml"), []byte(c.yaml), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			path := writeStubs(t, map[string]string{
				"gh": `case "$1 $2" in
  "release view") echo v9.9.9 ;;
  "release download") echo "$3" > "$GOT" ;;
  *) exit 1 ;;
esac`,
				"tar":     "exit 0",
				"sudo":    "exit 0",
				"anchors": "exit 0",
			})
			cmd := pipelineBash(t, "-c", script)
			cmd.Dir = dir
			cmd.Env = append(os.Environ(), "PATH="+path, "GOT="+got)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("the install step failed: %v\n%s", err, out)
			}
			b, err := os.ReadFile(got)
			if err != nil {
				t.Fatalf("nothing was downloaded: %v", err)
			}
			if tag := strings.TrimSpace(string(b)); tag != c.want {
				t.Errorf("downloaded %q, want %q", tag, c.want)
			}
		})
	}
}

// THE RESOLVER PUSHES ONLY WITH AN APP TOKEN.
//
// A push with `GITHUB_TOKEN` leaves every check of the PR at `action_required`: the PR
// looks green with no Anchors check run (reference app). The token has to be minted from
// the App's secrets, and a step's `if:` cannot read `secrets` — so the presence test lives
// in the job's `env`, and the step reads `env`.
func TestResolveQueuePushesOnlyWithAnAppToken(t *testing.T) {
	job, ok := parseWorkflow(t, "anchors-resolve-queue.yml")["varrer"]
	if !ok {
		t.Fatal("the `varrer` job vanished from the resolver")
	}
	has := job.Env["ANCHORS_HAS_APP"]
	if !strings.Contains(has, "secrets.ANCHORS_APP_ID") || !strings.Contains(has, "secrets.ANCHORS_APP_PRIVATE_KEY") {
		t.Errorf("the job env does not test both App secrets: ANCHORS_HAS_APP=%q", has)
	}
	var app, checkout, resolve *workflowStep
	for i := range job.Steps {
		s := &job.Steps[i]
		switch {
		case strings.HasPrefix(s.Uses, "actions/create-github-app-token@"):
			app = s
		case strings.HasPrefix(s.Uses, "actions/checkout@"):
			checkout = s
		case strings.Contains(s.Run, "git push"):
			resolve = s
		}
	}
	if app == nil || checkout == nil || resolve == nil {
		t.Fatalf("missing a step (app token: %v, checkout: %v, push: %v)", app != nil, checkout != nil, resolve != nil)
	}
	if strings.Contains(app.If, "secrets.") || !strings.Contains(app.If, "env.ANCHORS_HAS_APP") {
		t.Errorf("the App token step must branch on env.ANCHORS_HAS_APP (a step `if:` cannot read secrets): if=%q", app.If)
	}
	if app.With["app-id"] == "" || app.With["private-key"] == "" {
		t.Error("the App token step does not pass the App id and private key")
	}
	if tok := checkout.With["token"]; !strings.Contains(tok, "steps."+app.ID+".outputs.token") {
		t.Errorf("the checkout does not use the App token, so `git push` would push as GITHUB_TOKEN: token=%q", tok)
	}
	if cp := resolve.Env["CAN_PUSH"]; !strings.Contains(cp, "steps."+app.ID+".outputs.token") {
		t.Errorf("CAN_PUSH does not come from the minted token: %q", cp)
	}
}

// resolveQueueRepo builds an origin whose `feature` branch conflicts with `develop` only in
// the generated `anchors.graph.yaml`, and a clone of it where the resolver runs.
func resolveQueueRepo(t *testing.T) (origin, work string) {
	t.Helper()
	root := t.TempDir()
	origin = filepath.Join(root, "origin.git")
	seed := filepath.Join(root, "seed")
	work = filepath.Join(root, "work")
	git := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t",
			"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t", "GIT_CONFIG_GLOBAL=/dev/null")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	write := func(content string) {
		if err := os.WriteFile(filepath.Join(seed, "anchors.graph.yaml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(root, "init", "-q", "--bare", origin)
	git(root, "init", "-q", "-b", "develop", seed)
	write("nodes: [a]\n")
	git(seed, "add", ".")
	git(seed, "commit", "-q", "-m", "base")
	git(seed, "push", "-q", origin, "develop")
	git(seed, "checkout", "-q", "-b", "feature")
	write("nodes: [a, feature]\n")
	git(seed, "commit", "-q", "-am", "feature")
	git(seed, "push", "-q", origin, "feature")
	git(seed, "checkout", "-q", "develop")
	write("nodes: [a, other]\n")
	git(seed, "commit", "-q", "-am", "other")
	git(seed, "push", "-q", origin, "develop")
	git(root, "clone", "-q", origin, work)
	return origin, work
}

func revParse(t *testing.T, dir, ref string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", dir, "rev-parse", ref).Output()
	if err != nil {
		t.Fatalf("rev-parse %s: %v", ref, err)
	}
	return strings.TrimSpace(string(out))
}

// runResolveStep runs the resolver's sweep against `work`, with `gh`, `anchors` and `sleep`
// stubbed. The `gh` stub serves PR #7 on branch `feature`, keeps the comments posted in
// `state`, and reports them back to the "already commented?" query.
func runResolveStep(t *testing.T, work, state, canPush string) string {
	t.Helper()
	var script string
	for _, s := range parseWorkflow(t, "anchors-resolve-queue.yml")["varrer"].Steps {
		if strings.Contains(s.Run, "git push") {
			script = s.Run
		}
	}
	path := writeStubs(t, map[string]string{
		"gh": `case "$1 $2" in
  "pr list") echo 7 ;;
  "pr view")
    case "$*" in
      *headRefName*) echo feature ;;
      *isCrossRepository*) echo false ;;
      *comments*) ls "$STATE" | grep -c '^comment-' || true ;;
    esac ;;
  "pr comment")
    n=$(ls "$STATE" | grep -c '^comment-' || true)
    while [ $# -gt 0 ]; do [ "$1" = --body-file ] && cp "$2" "$STATE/comment-$n"; shift; done ;;
  *) exit 1 ;;
esac`,
		"anchors": `case "$1" in
  generated-paths) echo 'anchors\.graph\.yaml' ;;
esac
exit 0`,
		"sleep": "exit 0",
	})
	cmd := pipelineBash(t, "-c", script)
	cmd.Dir = work
	cmd.Env = append(os.Environ(), "PATH="+path, "STATE="+state, "CAN_PUSH="+canPush,
		"GH_REPO=o/r", "BASE=develop", "GITHUB_STEP_SUMMARY="+filepath.Join(state, "summary"),
		"GIT_CONFIG_GLOBAL=/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("the resolve step failed: %v\n%s", err, out)
	}
	return string(out)
}

func postedComments(t *testing.T, state string) []string {
	t.Helper()
	entries, err := os.ReadDir(state)
	if err != nil {
		t.Fatal(err)
	}
	var bodies []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "comment-") {
			b, err := os.ReadFile(filepath.Join(state, e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			bodies = append(bodies, string(b))
		}
	}
	return bodies
}

// WITHOUT AN APP TOKEN, THE RESOLVER COMMENTS — ONCE — AND PUSHES NOTHING.
//
// Run against a real git repository: the PR branch must stay where it was, the comment
// must carry the exact local commands, and a second sweep (the next merge into the base)
// must not repeat it.
func TestResolveQueueWithoutAppTokenCommentsOnceAndDoesNotPush(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("no bash on PATH")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git on PATH")
	}
	origin, work := resolveQueueRepo(t)
	state := t.TempDir()
	before := revParse(t, origin, "refs/heads/feature")

	runResolveStep(t, work, state, "false")

	// The run summary counts the PR as commented, not as resolved.
	summary, err := os.ReadFile(filepath.Join(state, "summary"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(summary), "- **0** PR(s) tiveram") ||
		!strings.Contains(string(summary), "- **1** PR(s) conflict only in generated files") {
		t.Errorf("the summary does not count the PR as commented-not-resolved:\n%s", summary)
	}
	if after := revParse(t, origin, "refs/heads/feature"); after != before {
		t.Error("the resolver PUSHED without an App token — the PR would look green with no Anchors check run")
	}
	comments := postedComments(t, state)
	if len(comments) != 1 {
		t.Fatalf("want exactly one comment on the PR, got %d", len(comments))
	}
	body := comments[0]
	for _, want := range []string{
		"<!-- anchors:resolve-queue:resolve-locally -->",
		"- `anchors.graph.yaml`",
		"git checkout feature\n",
		"git merge --no-commit origin/develop\n",
		"anchors map build\nanchors docs build\n", // map BEFORE docs: the docs read the map
		"git commit --no-edit\n",
		"ANCHORS_APP_ID",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("the comment lacks %q:\n%s", want, body)
		}
	}

	runResolveStep(t, work, state, "false")
	if n := len(postedComments(t, state)); n != 1 {
		t.Errorf("the second sweep commented again: %d comments — every merge into the base "+
			"would add one more to the same PR", n)
	}
}

// WITH AN APP TOKEN, THE RESOLVER KEEPS RESOLVING: it pushes the merge, and says nothing.
func TestResolveQueueWithAppTokenPushes(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("no bash on PATH")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git on PATH")
	}
	origin, work := resolveQueueRepo(t)
	state := t.TempDir()
	before := revParse(t, origin, "refs/heads/feature")

	runResolveStep(t, work, state, "true")

	if after := revParse(t, origin, "refs/heads/feature"); after == before {
		t.Error("with an App token the resolver did not push the resolution")
	}
	if n := len(postedComments(t, state)); n != 0 {
		t.Errorf("with an App token the resolver also commented (%d) — the comment is the fallback", n)
	}
}

// EVERY template the flow seeds carries the marker. Without it, `doctor --fix` treats the
// installed file as the team's own and never updates it — `anchors-resolve-queue.yml` was
// born without it, so its fixes (the GitHub App push, the `min_version` install) reached
// no project that already had it.
func TestEveryFlowTemplateCarriesTheMarker(t *testing.T) {
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), MarcadorDeTemplate) {
			t.Errorf("%s has no %q — once installed, `doctor --fix` never updates it", w.Arquivo, MarcadorDeTemplate)
		}
	}
}

// A card the stale RELEASED in `in-progress` goes back to the queue: the stale keeps the
// state (moving would presume where the work is), so the claim has to offer it. In
// the reference app 15 released cards sat in `in-progress` for days, offered to nobody.
// An `in-progress` card with NO declared owner is not taken: someone moved it by hand.
func TestClaimOffersAReleasedInProgressCard(t *testing.T) {
	if out, handed := runClaim(t, "FAKE_CARD_STATE=in-progress",
		"FAKE_OWNER=anchors-owner: (liberado) — sem progresso há mais de 24h"); !handed {
		t.Errorf("a released in-progress card was not handed out:\n%s", out)
	}
	if out, handed := runClaim(t, "FAKE_CARD_STATE=in-progress"); handed {
		t.Errorf("an in-progress card with no declared owner was taken:\n%s", out)
	}
	if out, handed := runClaim(t, "FAKE_CARD_STATE=in-progress", "FAKE_OWNER=anchors-owner: other/agent"); handed {
		t.Errorf("an in-progress card another agent owns was taken:\n%s", out)
	}
	if out, handed := runClaim(t); !handed {
		t.Errorf("control: a to-do card with no owner must still be handed out:\n%s", out)
	}
}

// A review the stale released in `in-review` is offered too, with the same guard: only an
// explicit `(liberado)`. In the reference app 4 such reviews sat with their PRs open since 19/09.
func TestClaimOffersAReleasedInReviewCard(t *testing.T) {
	if out, handed := runClaim(t, "FAKE_CARD_STATE=in-review",
		"FAKE_OWNER=anchors-owner: (liberado) — sem progresso há mais de 24h"); !handed {
		t.Errorf("a released in-review card was not handed out:\n%s", out)
	}
	if out, handed := runClaim(t, "FAKE_CARD_STATE=in-review"); handed {
		t.Errorf("an in-review card with no declared owner was taken:\n%s", out)
	}
}

// The claim reads "is there a card that is already mine?" in pages of 25: the old
// `gh issue list --json ...,comments` over 200 cards got 502 on every run in the reference app.
// A card already owned is resumed before anything new is handed out.
func TestClaimReadsItsOwnCardInPages(t *testing.T) {
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	i := strings.Index(s, "meu=$(")
	if i < 0 {
		t.Fatal("the own-card lookup vanished from the claim")
	}
	head := s[i : i+200]
	if strings.Contains(head, "gh issue list") || !strings.Contains(head, "gh api graphql --paginate") {
		t.Errorf("the own-card lookup is not the paginated GraphQL query:\n%s", head)
	}
	if !strings.Contains(s, "issues(first: 25") {
		t.Error("the page size is not 25")
	}
	// Wired: when the lookup finds the agent's card, the claim resumes it.
	out, _ := runClaim(t, "FAKE_MINE=9")
	if !strings.Contains(out, "9") {
		t.Errorf("the agent's own card #9 was not resumed:\n%s", out)
	}
}

// A card released by hand in flight is never handed back to WHOEVER RELEASED IT: that was
// the loop of the reference app. Any other agent takes it — skipping it for everyone left 29
// cards in the working columns with no owner (reference app, 2026-09-24). An approved review
// released by the verdict waits for the merge, and is taken by no one.
func TestClaimDoesNotReofferAHandReleasedCardInFlight(t *testing.T) {
	for _, st := range []string{"in-review", "in-progress"} {
		if out, handed := runClaim(t, "FAKE_CARD_STATE="+st,
			"FAKE_OWNER=anchors-owner: (liberado)", "FAKE_PREV_OWNER=machine/session"); handed {
			t.Errorf("%s: a card released by hand was handed back to whoever released it:\n%s", st, out)
		}
		if out, handed := runClaim(t, "FAKE_CARD_STATE="+st,
			"FAKE_OWNER=anchors-owner: (liberado)", "FAKE_PREV_OWNER=other/agent"); !handed {
			t.Errorf("%s: a card another agent released by hand was offered to no one:\n%s", st, out)
		}
		if out, handed := runClaim(t, "FAKE_CARD_STATE="+st,
			"FAKE_OWNER=anchors-owner: (liberado) — revisão concluída: approved por other/agent",
			"FAKE_PREV_OWNER=other/agent"); handed {
			t.Errorf("%s: an approved review waiting for the merge was handed out:\n%s", st, out)
		}
		if out, handed := runClaim(t, "FAKE_CARD_STATE="+st,
			"FAKE_OWNER=anchors-owner: (liberado) — sem progresso há mais de 24h"); !handed {
			t.Errorf("%s: a card the stale released was not handed out:\n%s", st, out)
		}
	}
	// A to-do card released by hand is still free: nothing is in flight there.
	if out, handed := runClaim(t, "FAKE_OWNER=anchors-owner: (liberado)"); !handed {
		t.Errorf("a to-do card released by hand must stay claimable:\n%s", out)
	}
}

// A card WAITING ON A PERSON is not the agent's open work: the own-card lookup skips
// `needs-user` too. the reference app, parked after a rejected review, went back to the same
// agent on every `anchors next`. Runs the workflow's real jq over GraphQL-shaped data.
func TestClaimOwnCardLookupSkipsEscalated(t *testing.T) {
	page := func(extra string) string {
		return `{"data":{"repository":{"issues":{"nodes":[{"number":9,` +
			`"labels":{"nodes":[{"name":"anchors"},{"name":"anchors:in-review"}` + extra + `]},` +
			`"comments":{"nodes":[{"body":"anchors-owner: machine/session"}]}}]}}}}`
	}
	if out, _ := runClaim(t, "FAKE_GRAPHQL="+page("")); !strings.Contains(out, "#9") {
		t.Errorf("control: the agent's own card #9 was not resumed:\n%s", out)
	}
	if out, _ := runClaim(t, "FAKE_GRAPHQL="+page(`,{"name":"anchors:needs-user"}`)); strings.Contains(out, "#9") {
		t.Errorf("a card waiting on a person was resumed as the agent's own work:\n%s", out)
	}
}

// WHOEVER WROTE IT DOES NOT REVIEW IT. A card in review released by the same agent that
// asks is its own work: it stays for another agent (reference app: six self-reviews in one day
// with two agents working). A to-do card handed back to its implementer is still theirs.
func TestClaimDoesNotHandAnAgentItsOwnReview(t *testing.T) {
	released := "FAKE_OWNER=anchors-owner: (liberado) — implementação concluída"
	if out, handed := runClaim(t, "FAKE_CARD_STATE=ready-to-review", released,
		"FAKE_PREV_OWNER=machine/session"); handed {
		t.Errorf("the agent was handed the review of its own work:\n%s", out)
	}
	if out, handed := runClaim(t, "FAKE_CARD_STATE=ready-to-review", released,
		"FAKE_PREV_OWNER=other/agent"); !handed {
		t.Errorf("the review of ANOTHER agent's work was not handed out:\n%s", out)
	}
	if out, handed := runClaim(t, released, "FAKE_PREV_OWNER=machine/session"); !handed {
		t.Errorf("a to-do card released by the same agent must stay claimable:\n%s", out)
	}
}

// The AUTHOR is the owner before the delivery ("implementação concluída"), not whoever
// released last: reviewers release too, and the previous reviewer must not pass for the
// author while the real author is handed its own re-review.
func TestClaimAuthorIsTheOwnerBeforeTheDelivery(t *testing.T) {
	released := "FAKE_OWNER=anchors-owner: (liberado) — revisão concluída"
	// The agent reviewed last time; another agent wrote it → the agent may review again.
	if out, handed := runClaim(t, "FAKE_CARD_STATE=ready-to-review", released,
		"FAKE_PREV_OWNER=machine/session", "FAKE_AUTHOR=other/agent"); !handed {
		t.Errorf("the previous reviewer was taken for the author:\n%s", out)
	}
	// The agent wrote it; someone else released last → still the agent's own work.
	if out, handed := runClaim(t, "FAKE_CARD_STATE=ready-to-review", released,
		"FAKE_PREV_OWNER=other/agent", "FAKE_AUTHOR=machine/session"); handed {
		t.Errorf("the author was handed the review of its own work:\n%s", out)
	}
}

// The claim never hands out a bug card: no agent in the queue can fix the pipeline or the
// tool. Runs the candidate filter of the claim through the real `jq`.
func TestClaimNeverHandsOutABug(t *testing.T) {
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("no jq on PATH")
	}
	b, err := workflowsFS.ReadFile("workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`--jq '(\.\[\] \| select\(\[\.labels\[\]\.name\] \| \(index\("anchors:needs-user"\)[^']*)'`).FindSubmatch(b)
	if m == nil {
		t.Fatal("the claim's candidate filter was not found")
	}
	in := `[{"number":1,"labels":[{"name":"anchors"},{"name":"anchors:to-do"},{"name":"anchors:bug"}]},
	        {"number":2,"labels":[{"name":"anchors"},{"name":"anchors:to-do"}]},
	        {"number":3,"labels":[{"name":"anchors"},{"name":"anchors:needs-user"}]}]`
	cmd := exec.Command("jq", "-r", string(m[1]))
	cmd.Stdin = strings.NewReader(in)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("jq: %v\n%s", err, out)
	}
	if got := strings.TrimSpace(string(out)); got != "2" {
		t.Errorf("candidates = %q, want only #2 (a bug and a decision are never handed out)", got)
	}
}

// the reference app: a `blocked-by-<n>` naming a PULL REQUEST could not be read without
// `pull-requests: read`, and every card waiting on a PR was skipped as still blocked.
func TestClaimCanReadABlockingPullRequest(t *testing.T) {
	b, err := workflowsFS.ReadFile("workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Permissions map[string]string `yaml:"permissions"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	if p := doc.Permissions["pull-requests"]; p != "read" && p != "write" {
		t.Errorf("the claim needs `pull-requests: read` to read the state of a blocking PR, has %q", p)
	}
}

// A REVIEW CARD WITH AN OPEN PR IS THE REVIEW. The open-PR skip exists so nobody
// reimplements a delivered card; applied to `ready-to-review` it emptied the review queue
// (reference app, 2026-09-24: 15 of 20 review cards skipped, no reviewer handed a review).
func TestClaimHandsOutAReviewCardThatHasItsPR(t *testing.T) {
	j, _ := json.Marshal([]map[string]any{{"number": 77, "body": "Closes #4"}})
	prs := "FAKE_PRS=" + string(j)
	if out, handed := runClaim(t, "FAKE_CARD_STATE=ready-to-review", prs); !handed {
		t.Errorf("a review card with its open PR was skipped — there is nothing else to review:\n%s", out)
	}
	if out, handed := runClaim(t, prs); handed {
		t.Errorf("control: a to-do card with an open PR must still be skipped:\n%s", out)
	}
}

// --- The Go side of the flow's pipelines (FLWRF): the list, the doctor's findings, the
// seeding and the board page. The tests above this block prove the templates themselves.

// Every pipeline declared in the list must EXIST embedded. The list is what the doctor
// checks and what `--fix` seeds — a name with no file would make the doctor demand
// something Anchors cannot create, and `--fix` would fail halfway.
func TestFlowEveryDeclaredWorkflowHasATemplate(t *testing.T) {
	t.Run("FLWRF-B01: Every declared pipeline has a carried template, a role and its serialization need", func(t *testing.T) {})
	notSerial := map[string]bool{"anchors-gates.yml": true, "anchors-board.yml": true}
	for _, w := range WorkflowsDoFluxo {
		if _, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo); err != nil {
			t.Errorf("%s is in the list and has no embedded template: %v", w.Arquivo, err)
		}
		if w.Papel == "" {
			t.Errorf("%s has no role — the doctor's message would have nothing to say about what stops happening", w.Arquivo)
		}
		if w.ExigeSerial == notSerial[w.Arquivo] {
			t.Errorf("%s: requires serialization = %v, want %v", w.Arquivo, w.ExigeSerial, !notSerial[w.Arquivo])
		}
	}
}

func flowCfg(branch string) *config.Config {
	return &config.Config{Workflow: &config.Workflow{IntegrationBranch: branch}}
}

func writePipeline(t *testing.T, dir, name, content string) string {
	t.Helper()
	wf := filepath.Join(dir, DirWorkflows)
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(wf, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestFlowMissingWorkflowSeesWhatIsAbsent(t *testing.T) {
	t.Run("FLWRF-B02: A pipeline is missing only when its file is absent", func(t *testing.T) {})
	t.Run("FLWRF-B04: Seeding writes the missing pipelines and returns them sorted", func(t *testing.T) {})
	dir := t.TempDir()

	if missing := MissingWorkflow(dir); len(missing) != len(WorkflowsDoFluxo) {
		t.Fatalf("empty project: expected %d missing, got %d", len(WorkflowsDoFluxo), len(missing))
	}
	// Presence only: any content counts as present.
	writePipeline(t, dir, "anchors-claim.yml", "not even yaml")
	for _, w := range MissingWorkflow(dir) {
		if w.Arquivo == "anchors-claim.yml" {
			t.Error("a present file, whatever its content, is not missing")
		}
	}
	if err := os.Remove(filepath.Join(dir, DirWorkflows, "anchors-claim.yml")); err != nil {
		t.Fatal(err)
	}

	written, _, err := SemeiaWorkflows(dir, flowCfg(""))
	if err != nil {
		t.Fatalf("seeding: %v", err)
	}
	if len(written) != len(WorkflowsDoFluxo) {
		t.Errorf("expected %d written, got %d", len(WorkflowsDoFluxo), len(written))
	}
	if !sort.StringsAreSorted(written) {
		t.Errorf("the written list should be sorted: %v", written)
	}
	if missing := MissingWorkflow(dir); len(missing) != 0 {
		t.Errorf("after seeding nothing should be missing: %v", missing)
	}
	if broken := SemConcurrency(dir); len(broken) != 0 {
		t.Errorf("the seeded templates carry concurrency: %v", broken)
	}
}

// A pipeline the team edited — another stale rhythm, one more permission, a step of its
// own — is deliberate work. Rewriting it with the default would erase the customization
// without warning. Same ruler as `install-hooks` with someone else's pre-commit.
func TestFlowSeedingDoesNotOverwriteWhatTheTeamEdited(t *testing.T) {
	t.Run("FLWRF-B05: Seeding leaves a pipeline without the marker untouched", func(t *testing.T) {})
	t.Run("FLWRF-X01: A file without the marker is never taken over", func(t *testing.T) {})
	dir := t.TempDir()
	mine := "# the team's pipeline, edited by hand\nname: mine\n"
	target := writePipeline(t, dir, WorkflowsDoFluxo[0].Arquivo, mine)

	written, _, err := SemeiaWorkflows(dir, flowCfg(""))
	if err != nil {
		t.Fatalf("seeding: %v", err)
	}

	b, _ := os.ReadFile(target)
	if string(b) != mine {
		t.Error("the pipeline the team edited was overwritten")
	}
	for _, e := range written {
		if e == WorkflowsDoFluxo[0].Arquivo {
			t.Error("the existing file should not be listed as written")
		}
	}
	if outdated := OutdatedWorkflows(dir, flowCfg("")); len(outdated) != 0 {
		t.Errorf("a team-owned file is never outdated: %v", outdated)
	}
}

// A pipeline that still carries the marker is Anchors' own: seeding brings it up to date,
// whatever it holds.
func TestFlowSeedingUpdatesAnIntactTemplate(t *testing.T) {
	t.Run("FLWRF-B04: Seeding writes the missing pipelines and returns them sorted", func(t *testing.T) {})
	dir := t.TempDir()
	name := WorkflowsDoFluxo[0].Arquivo
	target := writePipeline(t, dir, name, "# "+MarcadorDeTemplate+"\nname: old\n")

	written, _, err := SemeiaWorkflows(dir, flowCfg(""))
	if err != nil {
		t.Fatalf("seeding: %v", err)
	}
	want, _ := fs.ReadFile(workflowsFS, "workflows/"+name)
	if b, _ := os.ReadFile(target); string(b) != string(want) {
		t.Error("an intact template should be rewritten with the current one")
	}
	found := false
	for _, e := range written {
		found = found || e == name
	}
	if !found {
		t.Errorf("%s was rewritten and should be listed as written: %v", name, written)
	}
}

// A pipeline PRESENT but without serialization is the worst case: it looks configured and
// brings the race back in silence. It must be a finding of its own, distinct from "missing".
func TestFlowNoConcurrencyCatchesAPipelineThatLooksOK(t *testing.T) {
	t.Run("FLWRF-B03: A serial pipeline without serialization is flagged", func(t *testing.T) {})
	dir := t.TempDir()
	noSerial := "name: claim\non:\n  workflow_dispatch:\njobs:\n  x:\n    runs-on: ubuntu-latest\n"
	writePipeline(t, dir, "anchors-claim.yml", noSerial)
	// Serialized, but cancelling the run in progress: still broken.
	writePipeline(t, dir, "anchors-stale.yml", "concurrency:\n  group: x\n  cancel-in-progress: true\n")
	// Not serial by design: never flagged.
	writePipeline(t, dir, "anchors-gates.yml", noSerial)

	flagged := map[string]bool{}
	for _, w := range SemConcurrency(dir) {
		flagged[w.Arquivo] = true
	}
	if !flagged["anchors-claim.yml"] {
		t.Error("claim without `concurrency` would hand the same card to two agents — it must be a finding")
	}
	if !flagged["anchors-stale.yml"] {
		t.Error("a pipeline that cancels the run in progress must be a finding")
	}
	if flagged["anchors-gates.yml"] {
		t.Error("gates does not require serialization")
	}
	// An absent file is another finding (missing), never counted twice.
	if len(flagged) != 2 {
		t.Errorf("only the two broken files should be flagged, got %v", flagged)
	}
	for _, w := range MissingWorkflow(dir) {
		if w.Arquivo == "anchors-claim.yml" {
			t.Error("the file exists — counting it as missing would report the same problem twice")
		}
	}
}

// The PAGE goes with the pipeline that publishes it: seeding one without the other leaves
// the flow halfway — the workflow runs and fails copying a file that does not exist.
func TestFlowSeedingWritesTheBoardPage(t *testing.T) {
	t.Run("FLWRF-B06: The board page is seeded outside the pipelines folder, with the marker", func(t *testing.T) {})
	dir := t.TempDir()
	if _, _, err := SemeiaWorkflows(dir, flowCfg("")); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, BoardFile))
	if err != nil {
		t.Fatalf("the board page was not seeded: %v", err)
	}
	if !strings.Contains(string(b), MarcadorDeTemplate) {
		t.Error("the page should carry the marker — without it `--fix` never updates it")
	}
	// The HTML stays OUTSIDE `.github/workflows/`: GitHub runs everything there, and an HTML
	// in that folder becomes an invalid workflow — a permanent syntax error in the adopter's
	// repository.
	if strings.Contains(BoardFile, DirWorkflows) {
		t.Errorf("the page cannot live in %s: GitHub would try to run it", DirWorkflows)
	}
}

// `doctor --fix` printed "the pipelines already exist and are up to date" while it
// REWROTE the board page in silence. The change showed up in the `git status` of whoever
// ran the command with nothing having announced it.
//
// The board seeding now returns WHAT IT DID. These tests demand the answers, and the one
// that avoids the noise is: identical is NOT an update.
func TestFlowBoardSeedingSaysWhatItDid(t *testing.T) {
	t.Run("FLWRF-B07: The board page seeding reports created, updated or unchanged", func(t *testing.T) {})
	t.Run("FLWRF-X01: A file without the marker is never taken over", func(t *testing.T) {})
	cfg := flowCfg("")

	t.Run("absent: created", func(t *testing.T) {
		dir := t.TempDir()
		_, board, err := SemeiaWorkflows(dir, cfg)
		if err != nil {
			t.Fatalf("seeding: %v", err)
		}
		if board != BoardCreated {
			t.Errorf("a new page should be BoardCreated, got %v", board)
		}
	})

	t.Run("already identical: untouched, and it does not say it updated", func(t *testing.T) {
		// Saying "updated" when nothing changed trains the reader to ignore the notice —
		// and then it stops serving when the change is real.
		dir := t.TempDir()
		if _, _, err := SemeiaWorkflows(dir, cfg); err != nil {
			t.Fatalf("first seeding: %v", err)
		}
		_, board, err := SemeiaWorkflows(dir, cfg)
		if err != nil {
			t.Fatalf("second seeding: %v", err)
		}
		if board != BoardUnchanged {
			t.Errorf("an identical page should be BoardUnchanged, got %v", board)
		}
	})

	t.Run("Anchors' own and behind: updated", func(t *testing.T) {
		dir := t.TempDir()
		if _, _, err := SemeiaWorkflows(dir, cfg); err != nil {
			t.Fatalf("first seeding: %v", err)
		}
		// An OLD Anchors page: the marker is there, and the content differs.
		target := filepath.Join(dir, BoardFile)
		b, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		old := strings.Replace(string(b), "<meta charset=\"utf-8\">",
			"<meta charset=\"utf-8\">\n<!-- old version -->", 1)
		if old == string(b) {
			t.Fatal("the test could not age the page: the replace target changed")
		}
		if err := os.WriteFile(target, []byte(old), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}

		_, board, err := SemeiaWorkflows(dir, cfg)
		if err != nil {
			t.Fatalf("second seeding: %v", err)
		}
		if board != BoardUpdated {
			t.Errorf("an Anchors page left behind should be BoardUpdated, got %v", board)
		}
	})

	t.Run("edited by the team: untouched, and the content stays", func(t *testing.T) {
		dir := t.TempDir()
		if _, _, err := SemeiaWorkflows(dir, cfg); err != nil {
			t.Fatalf("first seeding: %v", err)
		}
		target := filepath.Join(dir, BoardFile)
		mine := "<!doctype html><p>the page belongs to the team now</p>"
		if err := os.WriteFile(target, []byte(mine), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}

		_, board, err := SemeiaWorkflows(dir, cfg)
		if err != nil {
			t.Fatalf("second seeding: %v", err)
		}
		if board != BoardUnchanged {
			t.Errorf("the team's page should be BoardUnchanged, got %v", board)
		}
		b, _ := os.ReadFile(target)
		if string(b) != mine {
			t.Error("Anchors overwrote a page the team took over")
		}
	})
}

// Outdated means: still Anchors' own (marker intact) and different from what seeding
// would write now.
func TestFlowOutdatedWorkflows(t *testing.T) {
	t.Run("FLWRF-B08: An intact pipeline that differs from its template is outdated", func(t *testing.T) {})
	dir := t.TempDir()
	cfg := flowCfg("")
	if _, _, err := SemeiaWorkflows(dir, cfg); err != nil {
		t.Fatal(err)
	}
	// One intact file falls behind; one becomes the team's; one disappears.
	aged := filepath.Join(dir, DirWorkflows, "anchors-claim.yml")
	b, _ := os.ReadFile(aged)
	if err := os.WriteFile(aged, append(b, []byte("# an older line\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	writePipeline(t, dir, "anchors-stale.yml", "name: ours now\n")
	if err := os.Remove(filepath.Join(dir, DirWorkflows, "anchors-guard.yml")); err != nil {
		t.Fatal(err)
	}

	var got []string
	for _, w := range OutdatedWorkflows(dir, cfg) {
		got = append(got, w.Arquivo)
	}
	if len(got) != 1 || got[0] != "anchors-claim.yml" {
		t.Errorf("outdated = %v, want only anchors-claim.yml", got)
	}
}

// The integration branch is where the work arrives. A project on `develop` gets pipelines
// that trigger on `develop`, and doctor does not call them outdated for it.
func TestFlowIntegrationBranchIsAppliedAndNotOutdated(t *testing.T) {
	t.Run("FLWRF-B09: The integration branch, main when none is declared, replaces every marked branch line", func(t *testing.T) {})
	t.Run("FLWRF-I01: What seeding writes is never outdated for the same configuration", func(t *testing.T) {})
	// `main` and the undeclared branch too: the resolve-queue template carried `develop`,
	// and the early return for `main` left it there — on a `main` project it never ran.
	for _, branch := range []string{"", "main", "develop"} {
		dir := t.TempDir()
		cfg := flowCfg(branch)
		if _, _, err := SemeiaWorkflows(dir, cfg); err != nil {
			t.Fatal(err)
		}
		if outdated := OutdatedWorkflows(dir, cfg); len(outdated) != 0 {
			t.Errorf("branch %q: freshly seeded pipelines reported outdated: %v", branch, outdated)
		}
		expected := branch
		if expected == "" {
			expected = "main"
		}
		marked := 0
		for _, w := range WorkflowsDoFluxo {
			b, err := os.ReadFile(filepath.Join(dir, DirWorkflows, w.Arquivo))
			if err != nil {
				t.Fatal(err)
			}
			tmpl, _ := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
			tl := strings.Split(string(tmpl), "\n")
			for i, l := range strings.Split(string(b), "\n") {
				if !strings.Contains(tl[i], "# anchors:integration-branch") {
					continue
				}
				marked++
				indent := tl[i][:len(tl[i])-len(strings.TrimLeft(tl[i], " "))]
				if want := indent + "branches: [" + expected + "] # anchors:integration-branch"; l != want {
					t.Errorf("branch %q: %s:%d = %q, want %q", branch, w.Arquivo, i+1, l, want)
				}
			}
		}
		if marked == 0 {
			t.Fatal("no template carries a marked branch line — the test proves nothing")
		}
	}
}

// Anchors writes only up to READY TO TEST; the columns after it belong to people.
func TestFlowColumnsAnchorsWrites(t *testing.T) {
	t.Run("FLWRF-B10: Anchors writes the board columns only up to READY TO TEST", func(t *testing.T) {})
	want := []string{"TO DO", "IN PROGRESS", "READY TO REVIEW", "IN REVIEW", "READY TO TEST"}
	if got := ColumnsAnchorsWrites(); !reflect.DeepEqual(got, want) {
		t.Errorf("ColumnsAnchorsWrites = %v, want %v", got, want)
	}
}

func TestFlowPerCardLabels(t *testing.T) {
	t.Run("FLWRF-B11: A per-card label is its prefix followed by the card", func(t *testing.T) {})
	for got, want := range map[string]string{
		LabelSob("44"):         "anchors:under-44",
		LabelDesbloqueia("44"): "anchors:desbloqueia-44",
		LabelDePR("556"):       "anchors:from-pr-556",
		LabelBlockedBy("900"):  "anchors:blocked-by-900",
	} {
		if got != want {
			t.Errorf("label = %q, want %q", got, want)
		}
	}
}

func TestFlowSeedingFailures(t *testing.T) {
	t.Run("FLWRF-E01: A pipelines folder that cannot be created fails the seeding", func(t *testing.T) {})
	t.Run("FLWRF-E02: A pipeline that cannot be written fails the seeding", func(t *testing.T) {})
	t.Run("FLWRF-E03: A template the binary does not carry fails the seeding", func(t *testing.T) {})

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".github"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := SemeiaWorkflows(dir, flowCfg("")); err == nil || !strings.Contains(err.Error(), "creating") {
		t.Errorf("a folder that cannot be created must fail, got %v", err)
	}

	testkit.SkipWithoutPOSIXPermissions(t) // what follows needs a folder that refuses writes
	dir = t.TempDir()
	wf := filepath.Join(dir, DirWorkflows)
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(wf, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(wf, 0o755) })
	if written, _, err := SemeiaWorkflows(dir, flowCfg("")); err == nil || !strings.Contains(err.Error(), "writing") || len(written) != 0 {
		t.Errorf("a pipeline that cannot be written must fail, got %v (written %v)", err, written)
	}

	saved := workflowsFS
	t.Cleanup(func() { workflowsFS = saved })
	workflowsFS = embed.FS{}
	if _, _, err := SemeiaWorkflows(t.TempDir(), flowCfg("")); err == nil || !strings.Contains(err.Error(), "reading the template") {
		t.Errorf("a template the binary does not carry must fail, got %v", err)
	}
}

// ── from the pipelines, run: the stale release ──

// staleFakeGH answers the stale script: the issue list runs the script's own `--jq`
// through the real `jq` over $FAKE_ISSUES; the comments of card N come from
// $FAKE_COMMENTS_N (a JSON array), and $FAKE_API_FAIL names a card whose read fails.
const staleFakeGH = `#!/usr/bin/env bash
echo "gh $*" >> "$FAKE_LOG"
call="$1 $2"; path="$2"
jqx=""
while [ $# -gt 0 ]; do
  case "$1" in --jq) jqx="$2"; shift ;; esac
  shift
done
case "$call" in
  "issue list") printf '%s' "$FAKE_ISSUES" | jq -r "$jqx" ;;
  "api "*)
    n=$(echo "$path" | sed -E 's#.*/issues/([0-9]+)/comments#\1#')
    if [ "${FAKE_API_FAIL:-}" = "$n" ]; then echo "gh: HTTP 502" >&2; exit 1; fi
    var="FAKE_COMMENTS_$n"
    printf '%s' "${!var:-[]}" | jq -r "$jqx" ;;
esac
exit 0
`

func runStale(t *testing.T, env ...string) (out, calls string) {
	t.Helper()
	for _, bin := range []string{"bash", "jq"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("no %s on PATH", bin)
		}
	}
	b, err := workflowsFS.ReadFile("workflows/anchors-stale.yml")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Jobs map[string]struct {
			Steps []struct {
				Run string `yaml:"run"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	script := ""
	for _, j := range doc.Jobs {
		for _, s := range j.Steps {
			if strings.Contains(s.Run, "anchors-owner: (liberado)") {
				script = s.Run
			}
		}
	}
	if script == "" {
		t.Fatal("the stale step vanished")
	}
	// Fixed cut-offs: macOS `date` has no `-d`, and the test needs a known clock.
	script = strings.Replace(script, `limite_espera=$(date -u -d "-${HORAS_ESPERA} hours" +%Y-%m-%dT%H:%M:%SZ)`, `limite_espera=2026-09-24T08:00:00Z`, 1)
	script = strings.Replace(script, `limite_andamento=$(date -u -d "-${HORAS_ANDAMENTO} hours" +%Y-%m-%dT%H:%M:%SZ)`, `limite_andamento=2026-09-23T10:00:00Z`, 1)
	script = strings.Replace(script, `limite_retrabalho=$(date -u -d "-${HORAS_RETRABALHO} hours" +%Y-%m-%dT%H:%M:%SZ)`, `limite_retrabalho=2026-09-24T09:00:00Z`, 1)
	if strings.Contains(script, "date -u -d") {
		t.Fatal("the cut-off lines changed shape; update the test's replacement")
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	testkit.FakeBin(t, bin, "gh", staleFakeGH)
	logFile := filepath.Join(dir, "gh.log")
	cmd := pipelineBash(t, "-c", script)
	cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"FAKE_LOG="+logFile, "GH_REPO=o/r", "LABEL=anchors", "HORAS_ESPERA=2", "HORAS_ANDAMENTO=24", "HORAS_RETRABALHO=1")
	cmd.Env = append(cmd.Env, env...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("the stale script failed: %v\n%s", err, o)
	}
	c, _ := os.ReadFile(logFile)
	return string(o), string(c)
}

// THE STALE JOB, RUN. Measured in the reference app: the list asked for the comments of 200 cards
// and GitHub's GraphQL answered 502/504 for 60 runs in a row. The list now carries no
// comments; the owner comes per card from REST; a card already released is not
// released again (the marker is the PREFIX); and an API failure skips the card.
func TestStaleReleasesByRESTOwnerAndOnlyOnce(t *testing.T) {
	t.Run("FLWRF-B12: The stale pipeline releases an idle owned card once", func(t *testing.T) {})
	old := "2026-09-20T10:00:00Z"
	issues := `[
	 {"number":5,"updatedAt":"` + old + `","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]},
	 {"number":6,"updatedAt":"` + old + `","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]},
	 {"number":7,"updatedAt":"` + old + `","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]},
	 {"number":8,"updatedAt":"2026-09-24T09:00:00Z","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]}
	]`
	out, calls := runStale(t, "FAKE_ISSUES="+issues,
		`FAKE_COMMENTS_5=[{"body":"anchors-owner: a/one"},{"body":"work log"}]`,
		`FAKE_COMMENTS_6=[{"body":"anchors-owner: a/one"},{"body":"anchors-owner: (liberado) — sem progresso há mais de 24h"}]`,
		`FAKE_COMMENTS_7=[{"body":"anchors-owner: a/one"}]`, "FAKE_API_FAIL=7",
		`FAKE_COMMENTS_8=[{"body":"anchors-owner: a/two"}]`)

	for _, line := range strings.Split(calls, "\n") {
		if strings.HasPrefix(line, "gh issue list") && strings.Contains(line, "comments") {
			t.Errorf("the issue list asks for comments again — the query GitHub timed out on: %s", line)
		}
	}
	if !strings.Contains(calls, "gh issue comment 5 --body anchors-owner: (liberado)") ||
		!strings.Contains(calls, "dono anterior: a/one") {
		t.Errorf("the stale owned card #5 was not released:\n%s\n%s", calls, out)
	}
	if strings.Contains(calls, "gh issue comment 6 ") {
		t.Error("card #6 was already released, and it was released again")
	}
	if strings.Contains(calls, "gh issue comment 7 ") || !strings.Contains(out, "#7") {
		t.Errorf("an API failure must skip card #7 and say so:\n%s", out)
	}
	if strings.Contains(calls, "gh issue comment 8 ") {
		t.Error("card #8 had progress inside the window and was released")
	}
}

// A REJECTED card goes back to `to-do` owned by its author, and the author's turn is short:
// 1h, not the 2h of an ordinary waiting card. The clock here: now 10:00, waiting cut-off
// 08:00, rework cut-off 09:00. A card last touched at 08:30 is released only if its owner
// came from a rejection (the "↩ PR #N rejected" line after the owner line).
func TestStaleReleasesARejectedCardAfterOneHour(t *testing.T) {
	t.Run("FLWRF-B13: A rejected card is released after the shorter rework window", func(t *testing.T) {})
	at := "2026-09-24T08:30:00Z"
	issues := `[
	 {"number":5,"updatedAt":"` + at + `","labels":[{"name":"anchors"},{"name":"anchors:to-do"}]},
	 {"number":6,"updatedAt":"` + at + `","labels":[{"name":"anchors"},{"name":"anchors:to-do"}]},
	 {"number":7,"updatedAt":"2026-09-24T07:00:00Z","labels":[{"name":"anchors"},{"name":"anchors:to-do"}]}
	]`
	out, calls := runStale(t, "FAKE_ISSUES="+issues,
		`FAKE_COMMENTS_5=[{"body":"anchors-owner: a/author"},{"body":"↩ PR #9 rejected by a/rev — the card is back in to-do"}]`,
		`FAKE_COMMENTS_6=[{"body":"anchors-owner: a/other"}]`,
		`FAKE_COMMENTS_7=[{"body":"anchors-owner: a/other"}]`)
	if !strings.Contains(calls, "gh issue comment 5 --body anchors-owner: (liberado)") {
		t.Errorf("a rejected card idle for 1h30 must go to any agent:\n%s\n%s", calls, out)
	}
	if strings.Contains(calls, "gh issue comment 6 ") {
		t.Errorf("an ordinary waiting card idle for 1h30 keeps its 2h:\n%s", calls)
	}
	if !strings.Contains(calls, "gh issue comment 7 --body anchors-owner: (liberado)") {
		t.Errorf("an ordinary waiting card idle for 3h must still be released:\n%s", calls)
	}
}

// ── from the pipelines, run: the review verdict ──

// The review outcome as the `anchors/review` commit status (reference app), and the merge
// that happens without it (reference app).
//
// These tests RUN the scripts of `anchors-pr-checks.yml` against a fake `gh` that serves
// JSON fixtures through the real `jq` — the same filters the pipeline uses in CI. Reading
// the YAML for strings would prove the words are there; running it proves which status a
// given card history produces.

const prChecksFile = "workflows/anchors-pr-checks.yml"

type prChecksDoc struct {
	On   map[string]any `yaml:"on"`
	Perm map[string]any `yaml:"permissions"`
	Jobs map[string]struct {
		If      string         `yaml:"if"`
		Needs   any            `yaml:"needs"`
		Outputs map[string]any `yaml:"outputs"`
		Steps   []struct {
			Name string            `yaml:"name"`
			Env  map[string]string `yaml:"env"`
			Run  string            `yaml:"run"`
		} `yaml:"steps"`
	} `yaml:"jobs"`
}

func loadPRChecks(t *testing.T) prChecksDoc {
	t.Helper()
	b, err := fs.ReadFile(workflowsFS, prChecksFile)
	if err != nil {
		t.Fatal(err)
	}
	var doc prChecksDoc
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

// reviewFakeGH is a `gh` that answers reads from JSON fixtures (filtered by the real `jq`, as
// `gh --jq` does) and records every write in calls.log. The fixture of a read is named
// after its positional arguments: `gh pr view 7` → pr_view_7.json, `gh api
// repos/o/r/issues/12/events` → api_repos_o_r_issues_12_events.json.
const reviewFakeGH = `#!/usr/bin/env bash
dir="$(cd "$(dirname "$0")" && pwd)"
jq_expr=""; pos=(); write=""; body=""; fields=()
while [ $# -gt 0 ]; do
  case "$1" in
    --jq) jq_expr="$2"; shift 2 ;;
    --add-label) fields+=("add-label=$2"); shift 2 ;;
    --json|--label|--limit|--state|--repo|--remove-label) shift 2 ;;
    --method) [ "$2" = "POST" ] && write=1; shift 2 ;;
    --body) body="$2"; shift 2 ;;
    -f) fields+=("$2"); shift 2 ;;
    --paginate) shift ;;
    *) pos+=("$1"); shift ;;
  esac
done
case "${pos[0]} ${pos[1]}" in
  "issue comment"|"issue edit") write=1 ;;
esac
if [ -n "$write" ]; then
  { printf 'CALL %s\n' "${pos[*]}"; for f in "${fields[@]}"; do printf 'FIELD %s\n' "$f"; done
    [ -n "$body" ] && printf 'BODY %s\n' "$body"; } >> "$dir/calls.log"
  exit 0
fi
key=$(IFS=_; echo "${pos[*]}" | tr '/' '_')
f="$dir/fixtures/$key.json"
[ -f "$f" ] || { echo "fake gh: no fixture $key" >&2; exit 1; }
if [ -n "$jq_expr" ]; then jq -r "$jq_expr" "$f"; else cat "$f"; fi
`

type ghWorld struct {
	t   *testing.T
	dir string
}

func newGHWorld(t *testing.T) *ghWorld {
	t.Helper()
	for _, bin := range []string{"bash", "jq", "base64"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s not on the PATH", bin)
		}
	}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "fixtures"), 0o755); err != nil {
		t.Fatal(err)
	}
	testkit.FakeBin(t, dir, "gh", reviewFakeGH)
	return &ghWorld{t: t, dir: dir}
}

func (w *ghWorld) fixture(key string, v any) {
	w.t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		w.t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.dir, "fixtures", key+".json"), b, 0o644); err != nil {
		w.t.Fatal(err)
	}
}

// run executes a step script with the fake `gh` first on the PATH. It returns the script
// output, the recorded writes and the GITHUB_OUTPUT content.
func (w *ghWorld) run(script string, env map[string]string) (out, calls, outputs string) {
	w.t.Helper()
	ghOut := filepath.Join(w.dir, "github_output")
	_ = os.WriteFile(ghOut, nil, 0o644)
	_ = os.Remove(filepath.Join(w.dir, "calls.log"))
	cmd := pipelineBash(w.t, "-c", script)
	cmd.Env = append(os.Environ(),
		"PATH="+w.dir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"GITHUB_OUTPUT="+ghOut,
		"GITHUB_STEP_SUMMARY="+filepath.Join(w.dir, "summary"),
		"GH_REPO=o/r", "LABEL=anchors",
		"RUN_URL=https://example.test/run/1",
	)
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		w.t.Fatalf("script failed: %v\n%s", err, b)
	}
	c, _ := os.ReadFile(filepath.Join(w.dir, "calls.log"))
	o, _ := os.ReadFile(ghOut)
	return string(b), string(c), string(o)
}

// Card #12, linked by the PR body, assigned to `agent-b` for review at 10:00.
type reviewWorld struct {
	labels     []string
	owners     []map[string]string // card comments
	assignedAt string              // "" = never moved to in-review
	prComments []map[string]any
	prBody     string
	perms      map[string]string // login → repository permission, as the API reports it
}

func defaultReviewWorld() reviewWorld {
	return reviewWorld{
		labels: []string{"anchors", "anchors:in-review"},
		owners: []map[string]string{
			{"createdAt": "2026-09-20T08:00:00Z", "body": "anchors-owner: agent-a"},
			{"createdAt": "2026-09-20T09:00:00Z", "body": "anchors-owner: (liberado) — implementação concluída"},
			{"createdAt": "2026-09-20T10:00:00Z", "body": "anchors-owner: agent-b"},
		},
		assignedAt: "2026-09-20T10:00:00Z",
		prBody:     "Implements the thing.\n\nRefs #12\n",
	}
}

func prComment(at, who, body string) map[string]any {
	return prCommentBy(at, who, "someone", body)
}

// prCommentBy names the author: the permission fallback looks the LOGIN up.
func prCommentBy(at, who, login, body string) map[string]any {
	return map[string]any{"created_at": at, "author_association": who, "user": map[string]string{"login": login}, "body": body}
}

func (w *ghWorld) seedReview(rw reviewWorld) {
	w.fixture("pr_view_7", map[string]any{
		"body": rw.prBody, "title": "some work", "headRefOid": "abc123",
	})
	var labels []map[string]string
	for _, l := range rw.labels {
		labels = append(labels, map[string]string{"name": l})
	}
	w.fixture("issue_view_12", map[string]any{"labels": labels, "comments": rw.owners})
	w.fixture("issue_list", []any{})
	events := []map[string]any{
		{"event": "labeled", "label": map[string]string{"name": "anchors:in-progress"}, "created_at": "2026-09-20T08:00:00Z"},
	}
	if rw.assignedAt != "" {
		events = append(events, map[string]any{
			"event": "labeled", "label": map[string]string{"name": "anchors:in-review"}, "created_at": rw.assignedAt,
		})
	}
	w.fixture("api_repos_o_r_issues_12_events", events)
	if rw.prComments == nil {
		rw.prComments = []map[string]any{}
	}
	w.fixture("api_repos_o_r_issues_7_comments", rw.prComments)
	for login, perm := range rw.perms {
		w.fixture("api_repos_o_r_collaborators_"+login+"_permission", map[string]string{"permission": perm})
	}
}

func reviewScript(t *testing.T) string {
	t.Helper()
	job, ok := loadPRChecks(t).Jobs["review"]
	if !ok || len(job.Steps) == 0 {
		t.Fatal("`anchors-pr-checks.yml` has no `review` job — nothing publishes `anchors/review`")
	}
	return job.Steps[0].Run
}

var statusState = regexp.MustCompile(`FIELD state=(\w+)`)

func publishedStatus(calls string) (state, desc string) {
	if !strings.Contains(calls, "CALL api repos/o/r/statuses/abc123") {
		return "", ""
	}
	if m := statusState.FindStringSubmatch(calls); m != nil {
		state = m[1]
	}
	for _, l := range strings.Split(calls, "\n") {
		if strings.HasPrefix(l, "FIELD description=") {
			desc = strings.TrimPrefix(l, "FIELD description=")
		}
	}
	return state, desc
}

func TestPRChecksPublishesTheReviewStatus(t *testing.T) {
	t.Run("FLWRF-B14: The review status follows the assigned reviewer's verdict line", func(t *testing.T) {})
	script := reviewScript(t)
	env := map[string]string{"PR": "7", "HEAD_SHA": "abc123"}

	cases := []struct {
		name      string
		edit      func(*reviewWorld)
		wantState string // "" = no status published
		wantDesc  string
	}{
		{
			name: "success only on the assigned reviewer's approved line",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "## Revisão\n\nAll rules hold.\n\nanchors-review: approved by agent-b\n"),
				}
			},
			wantState: "success", wantDesc: "approved by agent-b",
		},
		{
			name: "failure on the assigned reviewer's rejected line",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: rejected by agent-b\r\n"),
				}
			},
			wantState: "failure", wantDesc: "rejected by agent-b",
		},
		{
			name:      "pending while the assigned reviewer has not posted a verdict",
			edit:      func(rw *reviewWorld) {},
			wantState: "pending", wantDesc: "awaiting the verdict of agent-b",
		},
		{
			name: "another agent's approval does not count",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-a"),
				}
			},
			wantState: "pending", wantDesc: "agent-b",
		},
		{
			name: "an approval from before the assignment belongs to an earlier round",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T09:30:00Z", "OWNER", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "pending", wantDesc: "agent-b",
		},
		{
			name: "a verdict line inside a code block is an example",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "Write this when done:\n\n```\nanchors-review: approved by agent-b\n```\n"),
				}
			},
			wantState: "pending",
		},
		{
			name: "a line quoted mid-sentence is not a verdict",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "I will post anchors-review: approved by agent-b later"),
				}
			},
			wantState: "pending",
		},
		{
			name: "a commenter without write access cannot sign the review",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "NONE", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "pending",
		},
		{
			// the reference app: an org with no PUBLIC members. The job's token sees the member's
			// comment as CONTRIBUTOR, and every approval was dropped.
			name: "a private org member's line counts by repository permission",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prCommentBy("2026-09-20T11:00:00Z", "CONTRIBUTOR", "dev", "anchors-review: approved by agent-b"),
				}
				rw.perms = map[string]string{"dev": "write"}
			},
			wantState: "success",
		},
		{
			name: "a contributor with only read access still cannot sign the review",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prCommentBy("2026-09-20T11:00:00Z", "CONTRIBUTOR", "outsider", "anchors-review: approved by agent-b"),
				}
				rw.perms = map[string]string{"outsider": "read"}
			},
			wantState: "pending",
		},
		{
			name: "the last verdict wins",
			edit: func(rw *reviewWorld) {
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
					prComment("2026-09-20T12:00:00Z", "MEMBER", "anchors-review: rejected by agent-b"),
				}
			},
			wantState: "failure",
		},
		{
			name: "the reviewer is the owner AT the assignment, not a later one",
			edit: func(rw *reviewWorld) {
				rw.owners = append(rw.owners, map[string]string{
					"createdAt": "2026-09-20T13:00:00Z", "body": "anchors-owner: agent-c",
				})
				rw.prComments = []map[string]any{
					prComment("2026-09-20T14:00:00Z", "OWNER", "anchors-review: approved by agent-c"),
				}
			},
			wantState: "pending", wantDesc: "agent-b",
		},
		{
			// Timestamps have one-second resolution: a release stamped in the same second
			// as the assignment sorts before it, and `(liberado)` is not a reviewer's name.
			name: "a released owner is never taken for the reviewer",
			edit: func(rw *reviewWorld) {
				rw.owners = append(rw.owners, map[string]string{
					"createdAt": "2026-09-20T10:00:00Z", "body": "anchors-owner: (liberado) — race",
				})
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "success", wantDesc: "approved by agent-b",
		},
		{
			name: "a card back in ready-to-review is pending, whatever the earlier round said",
			edit: func(rw *reviewWorld) {
				rw.labels = []string{"anchors", "anchors:ready-to-review"}
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "pending", wantDesc: "awaiting a reviewer",
		},
		{
			name: "a card never assigned for review gets no status",
			edit: func(rw *reviewWorld) {
				rw.labels = []string{"anchors", "anchors:in-progress"}
				rw.assignedAt = ""
			},
			wantState: "",
		},
		{
			name: "a PR that declares no card gets no status",
			edit: func(rw *reviewWorld) {
				rw.prBody = "Pipeline fix, no card.\n"
				rw.prComments = []map[string]any{
					prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
				}
			},
			wantState: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := newGHWorld(t)
			rw := defaultReviewWorld()
			c.edit(&rw)
			w.seedReview(rw)
			out, calls, outputs := w.run(script, env)
			state, desc := publishedStatus(calls)
			if state != c.wantState {
				t.Fatalf("anchors/review = %q, want %q\noutput:\n%s\ncalls:\n%s", state, c.wantState, out, calls)
			}
			if c.wantState != "" {
				if !strings.Contains(calls, "FIELD context=anchors/review") {
					t.Errorf("the status must be published under the context `anchors/review`:\n%s", calls)
				}
				if !strings.Contains(outputs, "state="+c.wantState) {
					t.Errorf("the job output must carry the state for `mover`, got:\n%s", outputs)
				}
			}
			if c.wantDesc != "" && !strings.Contains(desc, c.wantDesc) {
				t.Errorf("description %q should name %q", desc, c.wantDesc)
			}
		})
	}
}

// The job must actually be wired: a script nobody triggers publishes nothing.
func TestPRChecksReviewJobIsWired(t *testing.T) {
	t.Run("FLWRF-B15: The review job is wired to verdict comments and feeds the mover", func(t *testing.T) {})
	doc := loadPRChecks(t)
	if _, ok := doc.On["issue_comment"]; !ok {
		t.Error("`anchors-pr-checks.yml` does not listen to `issue_comment` — the verdict line " +
			"is a PR comment, and without the trigger `anchors/review` never leaves pending")
	}
	if doc.Perm["statuses"] != "write" {
		t.Error("`anchors-pr-checks.yml` needs `statuses: write` to publish `anchors/review`")
	}
	review := doc.Jobs["review"]
	if !strings.Contains(review.If, "issue_comment") || !strings.Contains(review.If, "anchors-review:") {
		t.Errorf("the `review` job should run on PR comments carrying a verdict line, if: %q", review.If)
	}
	mover := doc.Jobs["mover"]
	if needs, _ := mover.Needs.(string); needs != "review" {
		t.Errorf("`mover` must run after `review` (needs: review) to know the outcome at the merge, got %v", mover.Needs)
	}
	// `always()`: a review job that failed or was skipped must not hold the card.
	if !strings.Contains(mover.If, "always()") {
		t.Errorf("`mover` must run even when `review` failed or was skipped, if: %q", mover.If)
	}
	// A comment is a verdict, and moves no card.
	if !strings.Contains(mover.If, "!= 'issue_comment'") {
		t.Errorf("`mover` must not run on a PR comment, if: %q", mover.If)
	}
	env := mover.Steps[0].Env
	if !strings.Contains(env["REVIEW_STATE"], "needs.review.outputs.state") ||
		!strings.Contains(env["REVIEWER"], "needs.review.outputs.reviewer") {
		t.Errorf("`mover` must receive the review outcome from the `review` job, env: %v", env)
	}
}

// The verdict line the pipeline parses is the one the review guide teaches. Two copies of
// one format drift the first time one of them changes, and here drifting means every
// reviewed PR stays pending.
func TestPRChecksVerdictLineMatchesTheReviewGuide(t *testing.T) {
	t.Run("FLWRF-I02: The parsed verdict line is the one the review guide teaches", func(t *testing.T) {})
	guide, err := os.ReadFile(filepath.Join("..", "..", "cmd", "anchors", "governance", "guide_review.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range []string{"anchors-review: approved by <you>", "anchors-review: rejected by <you>"} {
		if !strings.Contains(string(guide), line) {
			t.Errorf("the review guide no longer teaches %q — the line `anchors-pr-checks.yml` parses", line)
		}
	}
	if !strings.Contains(reviewScript(t), `^anchors-review: (approved|rejected) by `) {
		t.Error("the `review` job no longer parses `anchors-review: approved|rejected by <reviewer>`")
	}
}

// moverScript returns the `mover` script with the `${{ … }}` expressions it inlines replaced
// by the values of one event.
func moverScript(t *testing.T, expr map[string]string) string {
	t.Helper()
	job := loadPRChecks(t).Jobs["mover"]
	script := job.Steps[0].Run
	return regexp.MustCompile(`\$\{\{\s*([^}]*?)\s*\}\}`).ReplaceAllStringFunc(script, func(m string) string {
		k := strings.TrimSpace(m[3 : len(m)-2])
		return expr[k]
	})
}

func (w *ghWorld) seedMerge(labels ...string) {
	w.fixture("pr_view_7", map[string]any{
		"body": "Refs #12\n", "title": "some work", "headRefOid": "abc123",
	})
	var ls []map[string]string
	for _, l := range append([]string{"anchors"}, labels...) {
		ls = append(ls, map[string]string{"name": l})
	}
	w.fixture("issue_view_12", map[string]any{"labels": ls, "comments": []any{}})
	w.fixture("pr_checks_7", []map[string]string{{"bucket": "pass"}})
}

var mergedEvent = map[string]string{
	"inputs.pr": "7", "github.event.action": "closed", "github.event.pull_request.merged": "true",
}

func TestPRChecksMergeWithoutReviewOutcomeIsSaidOnTheCard(t *testing.T) {
	t.Run("FLWRF-B16: A merge without the review outcome is said on the card", func(t *testing.T) {})
	script := moverScript(t, mergedEvent)

	t.Run("merged with anchors/review success: the usual move, no warning", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:in-review")
		_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "success", "REVIEWER": "agent-b"})
		if !strings.Contains(calls, "FIELD add-label=anchors:ready-to-test") {
			t.Fatalf("the merge should move the card to ready-to-test:\n%s", calls)
		}
		if strings.Contains(calls, "WITHOUT the review outcome") {
			t.Errorf("a reviewed merge must not be flagged:\n%s", calls)
		}
	})

	t.Run("merged while pending: the card moves AND says who was reviewing", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:in-review")
		_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "pending", "REVIEWER": "agent-b"})
		if !strings.Contains(calls, "FIELD add-label=anchors:ready-to-test") {
			t.Fatalf("the card must still move to ready-to-test — the work landed:\n%s", calls)
		}
		if !strings.Contains(calls, "PR #7 merged WITHOUT the review outcome") {
			t.Fatalf("the card must say the PR merged without the review outcome:\n%s", calls)
		}
		if !strings.Contains(calls, "reviewer: **agent-b**") || !strings.Contains(calls, "`pending`") {
			t.Errorf("the comment must name the reviewer and the status it had:\n%s", calls)
		}
	})

	t.Run("merged while rejected: flagged too", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:in-review")
		_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "failure", "REVIEWER": "agent-b"})
		if !strings.Contains(calls, "WITHOUT the review outcome") || !strings.Contains(calls, "`failure`") {
			t.Errorf("a merge over a rejection must be flagged on the card:\n%s", calls)
		}
	})

	t.Run("merged before any review was assigned: flagged, saying no reviewer", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:ready-to-review")
		_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "", "REVIEWER": ""})
		if !strings.Contains(calls, "WITHOUT the review outcome") {
			t.Fatalf("a merge with no review status must be flagged:\n%s", calls)
		}
		if !strings.Contains(calls, "no reviewer had been assigned") || !strings.Contains(calls, "not reported") {
			t.Errorf("the comment must say no reviewer was assigned and no status reported:\n%s", calls)
		}
	})
}

// The move to `ready-to-review` is when the review becomes OWED: `mover` publishes
// `anchors/review` pending right there, because the `review` job ran before the move.
func TestPRChecksGreenPRPublishesReviewPending(t *testing.T) {
	t.Run("FLWRF-B17: A green PR publishes the review as pending", func(t *testing.T) {})
	script := moverScript(t, map[string]string{"inputs.pr": "7", "github.event.action": "synchronize"})
	w := newGHWorld(t)
	w.seedMerge("anchors:in-progress")
	_, calls, _ := w.run(script, map[string]string{"REVIEW_STATE": "none", "REVIEWER": ""})
	if !strings.Contains(calls, "FIELD add-label=anchors:ready-to-review") {
		t.Fatalf("a green PR should move the card to ready-to-review:\n%s", calls)
	}
	state, desc := publishedStatus(calls)
	if state != "pending" || !strings.Contains(calls, "FIELD context=anchors/review") {
		t.Errorf("the move to ready-to-review must publish anchors/review pending, got %q:\n%s", state, calls)
	}
	if !strings.Contains(desc, "awaiting a reviewer") {
		t.Errorf("the pending status should say it awaits a reviewer, got %q", desc)
	}
}

// The claim hands the reviewer the instructions, and they must teach the SAME line the
// `review` job parses. The claim used to say "aceito → move the card to ready-to-test",
// which bypasses the status and, done before the merge, hides the #782 comment.
func TestClaimTeachesTheReviewVerdictLine(t *testing.T) {
	t.Run("FLWRF-B18: The claim teaches the reviewer the verdict line", func(t *testing.T) {})
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{"anchors-review: approved by <você>", "anchors-review: rejected by <você>"} {
		if !strings.Contains(s, want) {
			t.Errorf("the reviewer instructions do not teach %q", want)
		}
	}
	if strings.Contains(s, "**aceito** → mova o card") {
		t.Error("the reviewer is still told to move the card on approval")
	}
}

// THE VERDICT FREES THE REVIEWER: after an approved/rejected line, the card no longer names
// the reviewer as owner, so `anchors next` stops resuming it. Pending changes nothing, and
// a card already released is not released twice.
func TestPRChecksVerdictReleasesTheReviewer(t *testing.T) {
	t.Run("FLWRF-B19: A verdict releases the reviewer once", func(t *testing.T) {})
	script := reviewScript(t)
	env := map[string]string{"PR": "7", "HEAD_SHA": "abc123"}
	release := "BODY anchors-owner: (liberado)"

	for _, verdict := range []string{"approved", "rejected"} {
		w := newGHWorld(t)
		rw := defaultReviewWorld()
		rw.prComments = []map[string]any{prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: "+verdict+" by agent-b")}
		w.seedReview(rw)
		_, calls, _ := w.run(script, env)
		if !strings.Contains(calls, "CALL issue comment 12") || !strings.Contains(calls, release) {
			t.Errorf("%s verdict did not release the reviewer:\n%s", verdict, calls)
		}
	}

	w := newGHWorld(t)
	w.seedReview(defaultReviewWorld()) // no verdict yet
	if _, calls, _ := w.run(script, env); strings.Contains(calls, release) {
		t.Errorf("a pending review released the reviewer:\n%s", calls)
	}

	w = newGHWorld(t)
	rw := defaultReviewWorld()
	rw.owners = append(rw.owners, map[string]string{"createdAt": "2026-09-20T11:30:00Z", "body": "anchors-owner: (liberado) — revisão concluída: approved por agent-b"})
	rw.prComments = []map[string]any{prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b")}
	w.seedReview(rw)
	if _, calls, _ := w.run(script, env); strings.Contains(calls, release) {
		t.Errorf("an already released card was released again:\n%s", calls)
	}
}

// the reference app: the body opens with `Refs #13` (a related card, itself in the
// queue) and ends with `Closes #12` (the card under review). The closing line wins in the
// review job AND in the mover, whatever the order in the body.
func TestPRChecksClosingLineWinsOverRefs(t *testing.T) {
	t.Run("FLWRF-B20: The closing line wins over a reference", func(t *testing.T) {})
	const body = "Refs #13 · the pattern is there\n\nImplements the thing.\n\nCloses #12\n"
	related := map[string]any{
		"labels":   []map[string]string{{"name": "anchors"}, {"name": "anchors:ready-to-review"}},
		"comments": []any{},
	}

	t.Run("review: the verdict counts for the closed card", func(t *testing.T) {
		w := newGHWorld(t)
		rw := defaultReviewWorld()
		rw.prBody = body
		rw.prComments = []map[string]any{
			prComment("2026-09-20T11:00:00Z", "OWNER", "anchors-review: approved by agent-b"),
		}
		w.seedReview(rw)
		w.fixture("issue_view_13", related)
		out, calls, _ := w.run(reviewScript(t), map[string]string{"PR": "7", "HEAD_SHA": "abc123"})
		if state, desc := publishedStatus(calls); state != "success" || !strings.Contains(desc, "card #12") {
			t.Fatalf("anchors/review = %q (%q), want success on card #12\noutput:\n%s", state, desc, out)
		}
	})

	t.Run("mover: the merge moves the closed card, not the referenced one", func(t *testing.T) {
		w := newGHWorld(t)
		w.seedMerge("anchors:in-review")
		w.fixture("pr_view_7", map[string]any{"body": body, "title": "some work", "headRefOid": "abc123"})
		w.fixture("issue_view_13", related)
		_, calls, _ := w.run(moverScript(t, mergedEvent), map[string]string{"REVIEW_STATE": "success", "REVIEWER": "agent-b"})
		if !strings.Contains(calls, "CALL issue edit 12") || strings.Contains(calls, "CALL issue edit 13") {
			t.Fatalf("the merge must move #12 and leave #13 alone:\n%s", calls)
		}
	})
}

// REJECTED sends the work back to its author: the card leaves `in-review` for `to-do`,
// owned by the author again (the owner before the reviewer). Their `anchors next` resumes
// it; if they do not ask for work, the stale releases it after the waiting window and any
// agent takes it. Before, it stayed in `in-review` with no owner and was offered to no one
// (reference app, 2026-09-24). APPROVED leaves the card for the merge.
func TestPRChecksRejectedGoesBackToTheAuthor(t *testing.T) {
	t.Run("FLWRF-B21: A rejection sends the card back to its author", func(t *testing.T) {})
	run := func(verdict string) string {
		w := newGHWorld(t)
		rw := defaultReviewWorld()
		rw.prComments = []map[string]any{
			prComment("2026-09-20T11:00:00Z", "OWNER", "## Revisão\n\nFails B02.\n\nanchors-review: "+verdict+" by agent-b\n"),
		}
		w.seedReview(rw)
		_, calls, _ := w.run(reviewScript(t), map[string]string{"PR": "7", "HEAD_SHA": "abc123"})
		return calls
	}

	calls := run("rejected")
	for _, want := range []string{
		"BODY anchors-owner: (liberado) — revisão concluída: rejected por agent-b",
		"FIELD add-label=anchors:to-do",
		"BODY anchors-owner: agent-a",
	} {
		if !strings.Contains(calls, want) {
			t.Errorf("a rejection should do %q:\n%s", want, calls)
		}
	}

	calls = run("approved")
	if strings.Contains(calls, "add-label=anchors:to-do") || strings.Contains(calls, "BODY anchors-owner: agent-a") {
		t.Errorf("an approval must leave the card for the merge:\n%s", calls)
	}
}

// pipelineBash runs a pipeline's script. The pipelines run on ubuntu-latest, and their
// scripts are Linux bash (GNU tools, `/tmp`, `date +%s`): on Windows there is nothing of
// theirs to prove, so the test that runs one skips there.
func pipelineBash(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the pipelines run on ubuntu-latest: their scripts are Linux bash")
	}
	return exec.Command("bash", args...)
}
