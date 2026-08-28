package initx

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
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

// A decisão da §7.13: o estado do trabalho é a COLUNA do Project, não uma label. Uma
// label de estado em paralelo duplicaria a verdade, e a dessincronia faria o board mentir
// sobre onde o trabalho está — num fluxo em que agentes escolhem o que pegar olhando o
// board, essa mentira vira trabalho duplicado.
func TestEstadoVivenoBoardNaoEmLabel(t *testing.T) {
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatalf("%s: %v", w.Arquivo, err)
		}
		texto := string(b)
		for _, proibido := range []string{`--add-label "doing"`, `--remove-label "doing"`} {
			if strings.Contains(texto, proibido) {
				t.Errorf("%s usa %s — o estado é a coluna do Project, não uma label", w.Arquivo, proibido)
			}
		}
	}
	// E o card criado tem de ENTRAR no board: fora dele, o claim nunca o encontra.
	b, _ := fs.ReadFile(workflowsFS, "workflows/anchors-identify.yml")
	if !strings.Contains(string(b), "gh project item-add") {
		t.Error("a issue criada precisa entrar no Project — fora do board ela é invisível ao claim")
	}
	// E o claim tem de escolher das colunas `READY TO ...`.
	b, _ = fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if !strings.Contains(string(b), `select(.status == "TO DO")`) {
		t.Error("o claim tem de tirar candidatos da coluna `TO DO`")
	}
}

// Toda coluna citada nos pipelines tem de existir na lista oficial. Um nome digitado
// diferente (`TODO` em vez de `TO DO`) não falha em lugar nenhum: o `item-edit` não acha
// a opção, o card não se move, e o trabalho some do fluxo em silêncio.
func TestPipelinesSoUsamColunasDeclaradas(t *testing.T) {
	valida := map[string]bool{}
	for _, c := range ColunasDoBoard {
		valida[c] = true
	}
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatalf("%s: %v", w.Arquivo, err)
		}
		for _, m := range regexp.MustCompile(`\.status == "([^"]+)"|select\(\.name == \$?d?\)|"(TO DO|IN PROGRESS|READY TO [A-Z]+|IN REVIEW|IN TEST|PRODUCTION)"`).FindAllStringSubmatch(string(b), -1) {
			for _, g := range m[1:] {
				if g != "" && !valida[g] {
					t.Errorf("%s usa a coluna %q, que não está em ColunasDoBoard", w.Arquivo, g)
				}
			}
		}
	}
}

// O Anchors escreve até `READY TO TEST` e para: as três últimas colunas são dos pipelines
// de entrega do projeto. Um pipeline do Anchors que movesse para lá passaria por cima de
// uma decisão que não é dele (quando um teste passou, quando um deploy aconteceu).
func TestAnchorsNaoEscreveAlemDeReadyToTest(t *testing.T) {
	naoNossas := map[string]bool{"IN TEST": true, "READY TO RELEASE": true, "PRODUCTION": true}
	for _, w := range WorkflowsDoFluxo {
		b, err := fs.ReadFile(workflowsFS, "workflows/"+w.Arquivo)
		if err != nil {
			t.Fatal(err)
		}
		for coluna := range naoNossas {
			if strings.Contains(string(b), `select(.name == "`+coluna+`")`) {
				t.Errorf("%s move para %q — além da alçada do Anchors", w.Arquivo, coluna)
			}
		}
	}
	if ColunaFinalDoAnchors != "READY TO TEST" {
		t.Errorf("a última coluna que o Anchors escreve mudou: %q", ColunaFinalDoAnchors)
	}
	if n := len(ColunasQueOAnchorsEscreve()); n != 5 {
		t.Errorf("o Anchors escreve em 5 colunas (TO DO..READY TO TEST), a lista diz %d", n)
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
		if !strings.Contains(texto, `--label "$LABEL"`) {
			t.Errorf("%s não filtra pela label do Anchors — pode tocar card de outro fluxo", w.Arquivo)
		}
	}

	// O claim é o caso mais exposto: ele lê o BOARD (que é de todos), não a lista de
	// issues do Anchors. A interseção com a label tem de ser explícita.
	b, err := fs.ReadFile(workflowsFS, "workflows/anchors-claim.yml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "sao-do-anchors") {
		t.Error("o claim seleciona do board compartilhado: precisa cruzar com a label antes de escolher")
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

	escritos, err := SemeiaWorkflows(dir)
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

	escritos, err := SemeiaWorkflows(dir)
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
