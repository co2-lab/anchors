package initx

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/board"
	"github.com/co2-lab/anchors/internal/config"
)

// Todo pipeline declarado na lista precisa EXISTIR embutido. A lista é o que o doctor
// confere e o que o `--fix` semeia — um nome sem arquivo faria o doctor cobrar algo que
// o Anchors não sabe criar, e o `--fix` falharia no meio.
func TestTodoWorkflowDeclaradoTemTemplate(t *testing.T) {
	for _, w := range WorkflowsDoFluxo {
		if _, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo); err != nil {
			t.Errorf("%s está na lista e não tem template embutido: %v", w.Arquivo, err)
		}
		if w.Papel == "" {
			t.Errorf("%s sem Papel — a mensagem do doctor não teria o que dizer que fica sem acontecer", w.Arquivo)
		}
	}
}

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
		if !strings.Contains(texto, `--label "$LABEL"`) {
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

func TestFaltaWorkflowVeOQueNaoExiste(t *testing.T) {
	dir := t.TempDir()

	if faltam := MissingWorkflow(dir); len(faltam) != len(WorkflowsDoFluxo) {
		t.Fatalf("projeto vazio: esperava %d faltando, veio %d", len(WorkflowsDoFluxo), len(faltam))
	}

	escritos, _, err := SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}})
	if err != nil {
		t.Fatalf("semear: %v", err)
	}
	if len(escritos) != len(WorkflowsDoFluxo) {
		t.Errorf("esperava %d escritos, veio %d", len(WorkflowsDoFluxo), len(escritos))
	}
	if faltam := MissingWorkflow(dir); len(faltam) != 0 {
		t.Errorf("depois de semear nada deveria faltar: %v", faltam)
	}
	if quebrados := SemConcurrency(dir); len(quebrados) != 0 {
		t.Errorf("os templates semeados trazem concurrency: %v", quebrados)
	}
}

// Um pipeline que o time editou — outro ritmo de stale, uma permissão a mais, um passo
// próprio — é trabalho deliberado. Reescrevê-lo pelo padrão apagaria a customização sem
// avisar. Mesma régua do `install-hooks` com um pre-commit alheio.
func TestSemeiaNaoSobrescreveOQueOTimeEditou(t *testing.T) {
	dir := t.TempDir()
	wf := filepath.Join(dir, DirWorkflows)
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	meu := "# pipeline do time, editado à mão\nname: meu\n"
	alvo := filepath.Join(wf, WorkflowsDoFluxo[0].Arquivo)
	if err := os.WriteFile(alvo, []byte(meu), 0o644); err != nil {
		t.Fatal(err)
	}

	escritos, _, err := SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}})
	if err != nil {
		t.Fatalf("semear: %v", err)
	}

	b, _ := os.ReadFile(alvo)
	if string(b) != meu {
		t.Error("o pipeline editado pelo time foi sobrescrito")
	}
	for _, e := range escritos {
		if e == WorkflowsDoFluxo[0].Arquivo {
			t.Error("o arquivo existente não deveria constar como escrito")
		}
	}
}

// Um pipeline PRESENTE mas sem serialização é o pior caso: parece configurado e devolve
// a corrida em silêncio. Tem de ser um achado próprio, distinto de "está faltando".
func TestSemConcurrencyPegaPipelineQuePareceOK(t *testing.T) {
	dir := t.TempDir()
	wf := filepath.Join(dir, DirWorkflows)
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	semSerial := "name: claim\non:\n  workflow_dispatch:\njobs:\n  x:\n    runs-on: ubuntu-latest\n"
	if err := os.WriteFile(filepath.Join(wf, "anchors-claim.yml"), []byte(semSerial), 0o644); err != nil {
		t.Fatal(err)
	}

	quebrados := SemConcurrency(dir)

	achou := false
	for _, w := range quebrados {
		if w.Arquivo == "anchors-claim.yml" {
			achou = true
		}
	}
	if !achou {
		t.Error("claim sem `concurrency` atribuiria o mesmo card a dois agentes — tem de ser achado")
	}
	// E não pode ser contado como ausente: o arquivo está lá.
	for _, w := range MissingWorkflow(dir) {
		if w.Arquivo == "anchors-claim.yml" {
			t.Error("o arquivo existe — contá-lo como ausente reportaria o mesmo problema duas vezes")
		}
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

// A PÁGINA acompanha o pipeline que a publica: semear um sem o outro deixa o fluxo pela
// metade — o workflow roda e falha ao copiar um arquivo que não existe.
func TestSemeiaEscreveAPaginaDoBoard(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, BoardFile))
	if err != nil {
		t.Fatalf("a página do board não foi semeada: %v", err)
	}
	if !strings.Contains(string(b), MarcadorDeTemplate) {
		t.Error("a página deveria trazer o marcador — sem ele o `--fix` nunca a atualiza")
	}
	// O HTML fica FORA de `.github/workflows/`: o GitHub executa tudo que está lá, e um
	// HTML naquele diretório vira um workflow inválido — erro de sintaxe permanente no
	// repositório de quem adotou.
	if strings.Contains(BoardFile, DirWorkflows) {
		t.Errorf("a página não pode morar em %s: o GitHub tentaria executá-la", DirWorkflows)
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
				cmd := exec.Command("bash", "-n")
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
// Medido no blue-eyes: o board publicava 3 de 7 cards. Os 4 ausentes eram exatamente os
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

// O `doctor --fix` imprimia "os pipelines já existem e estão atualizados" enquanto
// REESCREVIA a página do board em silêncio. A mudança aparecia no `git status` de quem
// rodou o comando sem nada tê-la anunciado — e a única forma de saber de quem era era ler
// o diff inteiro.
//
// `semeiaBoard` passou a devolver O QUE FEZ. Estes testes cobram as três respostas, e a
// terceira é a que evita o ruído: idêntico NÃO é atualização.
func TestSemeiaBoard_dizOQueFez(t *testing.T) {
	cfg := &config.Config{Workflow: &config.Workflow{}}

	t.Run("não existia: criado", func(t *testing.T) {
		dir := t.TempDir()
		_, board, err := SemeiaWorkflows(dir, cfg)
		if err != nil {
			t.Fatalf("semear: %v", err)
		}
		if board != BoardCreated {
			t.Errorf("página nova deveria ser BoardCreated, e foi %v", board)
		}
	})

	t.Run("já idêntica: intocada, e não diz que atualizou", func(t *testing.T) {
		// Dizer "atualizei" quando nada mudou treina quem lê a ignorar o aviso — e aí ele
		// deixa de servir quando a mudança for real.
		dir := t.TempDir()
		if _, _, err := SemeiaWorkflows(dir, cfg); err != nil {
			t.Fatalf("primeira semeadura: %v", err)
		}
		_, board, err := SemeiaWorkflows(dir, cfg)
		if err != nil {
			t.Fatalf("segunda semeadura: %v", err)
		}
		if board != BoardUnchanged {
			t.Errorf("página idêntica deveria ser BoardUnchanged, e foi %v", board)
		}
	})

	t.Run("do Anchors e para trás: atualizada", func(t *testing.T) {
		dir := t.TempDir()
		if _, _, err := SemeiaWorkflows(dir, cfg); err != nil {
			t.Fatalf("primeira semeadura: %v", err)
		}
		// Uma página VELHA do Anchors: o marcador está lá, e o conteúdo difere.
		alvo := filepath.Join(dir, BoardFile)
		b, err := os.ReadFile(alvo)
		if err != nil {
			t.Fatalf("ler: %v", err)
		}
		velha := strings.Replace(string(b), "<meta charset=\"utf-8\">",
			"<meta charset=\"utf-8\">\n<!-- versão antiga -->", 1)
		if velha == string(b) {
			t.Fatal("o teste não conseguiu envelhecer a página: o alvo do replace mudou")
		}
		if err := os.WriteFile(alvo, []byte(velha), 0o644); err != nil {
			t.Fatalf("escrever: %v", err)
		}

		_, board, err := SemeiaWorkflows(dir, cfg)
		if err != nil {
			t.Fatalf("segunda semeadura: %v", err)
		}
		if board != BoardUpdated {
			t.Errorf("página do Anchors que ficou para trás deveria ser BoardUpdated, e foi %v", board)
		}
	})

	t.Run("editada pelo time: intocada, e o conteúdo fica", func(t *testing.T) {
		dir := t.TempDir()
		if _, _, err := SemeiaWorkflows(dir, cfg); err != nil {
			t.Fatalf("primeira semeadura: %v", err)
		}
		alvo := filepath.Join(dir, BoardFile)
		minha := "<!doctype html><p>a página é do time agora</p>"
		if err := os.WriteFile(alvo, []byte(minha), 0o644); err != nil {
			t.Fatalf("escrever: %v", err)
		}

		_, board, err := SemeiaWorkflows(dir, cfg)
		if err != nil {
			t.Fatalf("segunda semeadura: %v", err)
		}
		if board != BoardUnchanged {
			t.Errorf("página do time deveria ser BoardUnchanged, e foi %v", board)
		}
		b, _ := os.ReadFile(alvo)
		if string(b) != minha {
			t.Error("o Anchors sobrescreveu uma página que o time assumiu")
		}
	})
}

// O CLAIM SERIALIZA QUEM PEGA O CARD, e isso não basta.
//
// O card é solto ao fim da implementação (`anchors-owner: (liberado)`) e volta à fila para
// REVISÃO — mas a branch e o PR do primeiro agente continuam de pé. O claim seguinte o
// entrega como se fosse trabalho novo, e o agente que o recebe reimplementa do zero.
//
// Medido no blue-eyes: TRÊS agentes resolveram a issue #375 em paralelo, sete minutos entre
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
	i := strings.Index(texto, "gh issue list --state all")
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
// Um `## FNDTN-F01 — o CI` é uma SEÇÃO dentro do arquivo do plano: não tem issue, não tem
// trinca, e nenhum agente a pega porque não há o que pegar. Medido no projeto de referência:
// 43 fases, TODAS sem card.
//
// O plano, ao contrário, É trabalho — `plans/0014-alertas-incidentes.md` é um arquivo regido,
// com card próprio e trinca a cumprir; 18 deles existem, 2 já fechados por agentes.
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
// MEDIDO no blue-eyes: 59 de 95 cards abertos estavam FORA do fluxo -- 53 sem label de
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
  "api "*) exit 1 ;;
  "issue list")
    case "$labels" in
      *anchors:desbloqueia-*) fail unblock ;;
      *" anchors:to-do"*) echo "$FAKE_CARD" ;;
    esac ;;
  "issue view")
    case "$json" in
      labels) case "$jqx" in
          *blocked-by*) fail labels; [ -z "${FAKE_BLOCKED_BY:-}" ] || echo "$FAKE_BLOCKED_BY" ;;
          *) echo "anchors:to-do" ;;
        esac ;;
      state) fail state; echo "${FAKE_STATE:-CLOSED}" ;;
      title) fail title ;;
      comments) echo "" ;;
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
	if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(fakeGH), 0o755); err != nil {
		t.Fatal(err)
	}
	logFile := filepath.Join(dir, "gh.log")
	cmd := exec.Command("bash", "-c", script)
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
// URL. Measured in blue-eyes #679: a PR body reporting a mutation score of `4/4` matched
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
