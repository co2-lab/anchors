package initx

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

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
	for _, c := range EstadosDoTrabalho {
		valida[c] = true
	}
	// A label de ESCALAÇÃO é válida sem ser estado: ela marca QUEM destrava o card, e o
	// card continua na coluna onde o trabalho parou. Tratá-la como estado a faria sair
	// dessa coluna, e o board deixaria de mostrar onde o fluxo travou.
	valida[LabelPrecisaDoUsuario] = true
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
	if EstadoFinalDoAnchors != "anchors:ready-to-test" {
		t.Errorf("o último estado que o Anchors escreve mudou: %q", EstadoFinalDoAnchors)
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

	if faltam := FaltaWorkflow(dir); len(faltam) != len(WorkflowsDoFluxo) {
		t.Fatalf("projeto vazio: esperava %d faltando, veio %d", len(WorkflowsDoFluxo), len(faltam))
	}

	escritos, err := SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}})
	if err != nil {
		t.Fatalf("semear: %v", err)
	}
	if len(escritos) != len(WorkflowsDoFluxo) {
		t.Errorf("esperava %d escritos, veio %d", len(WorkflowsDoFluxo), len(escritos))
	}
	if faltam := FaltaWorkflow(dir); len(faltam) != 0 {
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

	escritos, err := SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}})
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
	for _, w := range FaltaWorkflow(dir) {
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
	if _, err := SemeiaWorkflows(dir, &config.Config{Workflow: &config.Workflow{}}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, ArquivoDoBoard))
	if err != nil {
		t.Fatalf("a página do board não foi semeada: %v", err)
	}
	if !strings.Contains(string(b), MarcadorDeTemplate) {
		t.Error("a página deveria trazer o marcador — sem ele o `--fix` nunca a atualiza")
	}
	// O HTML fica FORA de `.github/workflows/`: o GitHub executa tudo que está lá, e um
	// HTML naquele diretório vira um workflow inválido — erro de sintaxe permanente no
	// repositório de quem adotou.
	if strings.Contains(ArquivoDoBoard, DirWorkflows) {
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
	if !strings.Contains(texto, `index("`+LabelPrecisaDoUsuario+`") | not`) {
		t.Error("o claim precisa EXCLUIR o card escalado da lista de disponíveis — " +
			"senão a escalação vira só um rótulo")
	}
	// E precisa escalar em algum limite: contar sem agir deixaria o ciclo rodando.
	if !strings.Contains(texto, "--add-label \""+LabelPrecisaDoUsuario+"\"") {
		t.Error("o claim precisa APLICAR a label ao atingir o limite")
	}
}

// A label de escalação tem de ser CRIADA no repositório. `gh issue edit` com label
// inexistente não é erro fatal — ele falha em silêncio, e o card ficaria travado sem o
// sinalizador que diz por quê.
func TestLabelDeEscalacaoNaoEhEstado(t *testing.T) {
	for _, e := range EstadosDoTrabalho {
		if e == LabelPrecisaDoUsuario {
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
	i := strings.Index(texto, "dono:")
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
// Anchors (`anchors-owner:` e `anchors:sob-<n>`); quem sabe a SINTAXE é o `pr-body`.
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
	// A palavra em inglês NÃO pode estar hardcoded no pipeline: ela vive no `pr-body`,
	// que é quem conhece a plataforma.
	if strings.Contains(texto, "close[sd]?") || strings.Contains(texto, "Closes #") {
		t.Error("a sintaxe de fechamento não pode estar no pipeline: ela é da PLATAFORMA " +
			"e vive no `anchors pr-body`, que é quem a conhece")
	}
}
