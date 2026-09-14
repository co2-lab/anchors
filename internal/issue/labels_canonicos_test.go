package issue

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// TODO LABEL QUE ESTE PACOTE APLICA TEM DE SER UM QUE O `anchors init` CRIA.
//
// O `gh` recusa o comando INTEIRO por um label inexistente:
//
//	could not add label: 'anchors:precisa-do-usuario' not found
//
// Então um literal divergente não degrada — ele apaga o registro. Medido no blue-eyes: o
// gate `open-questions-resolved` reprovou uma spec com duas perguntas em aberto (o
// comportamento certo) e a issue NÃO foi criada. O card que devia esperar a decisão do
// usuário nunca chegou ao board, e só um aviso no stderr disse algo.
//
// Um gate que acusa e não registra é pior do que não ter gate: o relatório da sessão diz
// que houve achado, e amanhã ninguém encontra.
//
// A causa foi a migração dos nomes de label do português para o inglês passando reto num
// literal — `LabelNeedsUser` já existia e este pacote não a usava. Este teste é o que
// impede a divergência de voltar por um terceiro literal.
func TestLabelsAplicadosSaoOsQueOInitCria(t *testing.T) {
	criados := map[string]bool{initx.LabelNeedsUser: true}
	for _, s := range initx.WorkStates {
		criados[s] = true
	}

	// os labels que `GitHub.Create` aplica além do label-raiz do projeto (esse vem da
	// configuração e é criado pelo `init` a partir dela)
	aplicados := []string{initx.LabelToDo, initx.LabelNeedsUser}

	for _, l := range aplicados {
		if !criados[l] {
			t.Errorf("o label %q é aplicado ao criar issue e o `init` não o cria — "+
				"o `gh` vai recusar o comando inteiro e o achado do gate não será registrado", l)
		}
	}
}

// O nome LEGADO não pode voltar a ser aplicado: ele não existe nos repositórios, e é
// exatamente o literal que apagava o registro.
func TestLabelLegadoNaoEhAplicado(t *testing.T) {
	if initx.LabelNeedsUser == initx.LabelNeedsUserLegacy {
		t.Fatal("o canônico e o legado ficaram iguais — a ponte de compatibilidade perdeu sentido")
	}
	if strings.Contains(initx.LabelNeedsUser, "precisa") {
		t.Errorf("LabelNeedsUser = %q — voltou ao nome em português", initx.LabelNeedsUser)
	}
}

// A PROSA QUE ENSINA também não pode citar o nome legado.
//
// A régua acima protege o que o código APLICA. Faltava o que a documentação ENSINA — e é
// por onde o defeito voltou: o `WORKFLOW.md` mandava rodar
//
//	gh issue list --label anchors:sob-44
//
// e o `BOOTSTRAP.md` dizia que o card ganha `anchors:precisa-do-usuario`. Nenhum dos dois
// existe num repositório inicializado pela versão atual. Quem seguisse o guia receberia uma
// lista VAZIA — e a lista vazia não parece erro, parece ausência de trabalho.
//
// Foi assim que três cards (#473, #483, #409) ficaram abertos com o trabalho já mergeado:
// o agente conferia se havia achados sob o card, a busca não devolvia nada, e o `Closes` do
// segundo card nunca entrou no corpo do PR.
//
// As constantes `*Legacy`/`*Antigo` continuam existindo e continuam certas — elas servem
// para LER cards antigos. O que este teste proíbe é ENSINAR o nome velho a quem cria
// trabalho novo.
func TestDocumentacaoNaoEnsinaLabelLegado(t *testing.T) {
	legados := map[string]string{
		initx.PrefixoLabelSobAntigo: initx.PrefixoLabelSob,
		initx.LabelNeedsUserLegacy:  initx.LabelNeedsUser,
	}

	raiz := repoRoot(t)
	docs, err := filepath.Glob(filepath.Join(raiz, "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) == 0 {
		t.Fatal("nenhum .md na raiz — o glob quebrou e o teste passaria vazio")
	}

	for _, doc := range docs {
		b, err := os.ReadFile(doc)
		if err != nil {
			t.Fatal(err)
		}
		texto := string(b)
		for velho, novo := range legados {
			// O bloco que EXPLICA a migração pode citar o nome velho — é do que ele
			// trata. O que não pode é a prosa ensinar o velho como se fosse o atual, e a
			// distinção é a menção ao novo na mesma linha.
			for _, linha := range strings.Split(texto, "\n") {
				if strings.Contains(linha, velho) && !strings.Contains(linha, novo) {
					t.Errorf("%s ensina o label legado %q (o atual é %q):\n\t%s",
						filepath.Base(doc), velho, novo, strings.TrimSpace(linha))
				}
			}
		}
	}
}

// repoRoot sobe até achar o `go.mod`: o teste roda com o diretório do PACOTE como corrente,
// e os documentos estão na raiz.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		pai := filepath.Dir(dir)
		if pai == dir {
			break
		}
		dir = pai
	}
	t.Fatal("não achei o go.mod subindo a partir do pacote")
	return ""
}
