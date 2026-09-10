package initx

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// TODO nome antigo aponta para um gate que EXISTE de verdade.
//
// Mora no `initx` e não no `config` porque é aqui que os dois lados são visíveis: o
// de-para vive no `config`, a lista de gates vive aqui, e `config` não pode importar
// `initx` (seria ciclo).
//
// Um de-para com destino errado é PIOR que não ter de-para: o projeto carrega, o gate
// vira um nome que nenhum verificador conhece, e o `check` reporta "gate declarado sem
// nada para medir" — sem dizer que a causa foi a conversão do nome antigo.
func TestVocabularioAntigoApontaParaGateQueExiste(t *testing.T) {
	existe := map[string]bool{}
	for _, g := range DefaultGates(todosOsArtefatos(), false) {
		existe[g.Name] = true
	}
	if len(existe) == 0 {
		t.Fatal("nenhum gate default — o teste não confrontaria nada")
	}

	depara := config.LegacyNames()
	if len(depara) == 0 {
		t.Fatal("o de-para está vazio — ou a migração não aconteceu, ou o teste perdeu o alvo")
	}

	for antigo, novo := range depara {
		if !existe[novo] {
			t.Errorf("%q → %q, e nenhum gate default se chama %q — o de-para aponta para o vazio",
				antigo, novo, novo)
		}
	}
}

// E o INVERSO: nenhum gate default ainda tem nome em português.
//
// Sem isto, alguém acrescentaria um gate novo com nome em português e o vocabulário
// voltaria a se misturar — que é o que esta migração inteira existe para resolver.
func TestNenhumGateDefaultTemNomeEmPortugues(t *testing.T) {
	// Os nomes ANTIGOS são exatamente a lista do que não pode mais aparecer.
	proibidos := config.LegacyNames()
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
