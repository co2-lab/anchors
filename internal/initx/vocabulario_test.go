package initx

import (
	"testing"

	"github.com/co2-lab/anchors/internal/migra"
)

// O DE-PARA da migração tem de apontar para gates que EXISTEM.
//
// A tabela vive no passo `1→2` (`internal/migra/formato_2.go`) e converte o nome legado
// para o canônico. Um destino que não corresponde a gate nenhum converteria o projeto para
// um nome inexistente — e o gate sumiria do `check` sem nada acusar, que é pior que o nome
// velho.
func TestVocabularioAntigoApontaParaGateQueExiste(t *testing.T) {
	existentes := map[string]bool{}
	for _, g := range DefaultGates(todosOsArtefatos(), false) {
		existentes[g.Name] = true
	}
	for _, passo := range migraSteps() {
		for arquivo, chaves := range passo.RenameValues {
			for chave, depara := range chaves {
				if chave != "id" && chave != "gate" {
					continue // `check:` aponta para verificador interno, não para gate
				}
				for velho, novo := range depara {
					if !existentes[novo] {
						t.Errorf("%s/%s: %q → %q, e o gate %q não existe nos defaults",
							arquivo, chave, velho, novo, novo)
					}
				}
			}
		}
	}
}

func TestNenhumGateDefaultTemNomeEmPortugues(t *testing.T) {
	// Os nomes legados são exatamente a lista do que não pode mais aparecer, e ela agora
	// vive no passo de migração.
	proibidos := map[string]string{}
	for _, passo := range migraSteps() {
		for _, chaves := range passo.RenameValues {
			for _, depara := range chaves {
				for velho, novo := range depara {
					proibidos[velho] = novo
				}
			}
		}
	}
	for _, g := range DefaultGates(todosOsArtefatos(), false) {
		if novo, ehAntigo := proibidos[g.Name]; ehAntigo {
			t.Errorf("o gate default %q ainda usa o nome antigo — deveria ser %q", g.Name, novo)
		}
	}
}

// todosOsArtefatos habilita TODAS as escolhas de artefato.
//
// O de-para tem de ser confrontado contra o conjunto INTEIRO de gates default, não contra
// um subconjunto: um destino que só existe quando o projeto escolheu "plan" passaria
// despercebido num teste que não pede planos — e o gate viraria um nome inexistente
// justamente no projeto que o usa.
func todosOsArtefatos() map[string]bool {
	return map[string]bool{
		"spec": true, "feature": true, "test": true,
		"code": true, "plan": true, "guide": true, "doc": true,
	}
}

// migraSteps expõe os passos registrados, para as réguas acima os confrontarem.
func migraSteps() []migra.Step {
	return migra.AllSteps()
}
