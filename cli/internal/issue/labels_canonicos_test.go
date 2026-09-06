package issue

import (
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
