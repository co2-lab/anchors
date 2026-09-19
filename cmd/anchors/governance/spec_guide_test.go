package governance

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func TestSpecGuideTemExemploCompleto(t *testing.T) {
	// A razão de este arquivo existir: o guide embutido (`anchors guide spec`) tem 94
	// linhas de regra e ZERO exemplos, e chega a instruir a descobrir o formato ERRANDO
	// ("rode `anchors check`; se reprovar, a mensagem do gate diz o formato esperado").
	// Ninguém escreve markdown conforme a partir de descrição — copia-se de um exemplo.
	g := renderSpecGuide(&config.Config{}, "LOGIX")
	for _, obrigatorio := range []string{
		"<!-- @anchors", // o cabeçalho, que é o mais fácil de errar
		"code: LOGIX",   // com o código no lugar certo
		"### LOGIX-S01", // o formato exato de regra catalogada
		// A seção vem do MESMO catálogo que o `anchors new spec` usa, resolvida pelo
		// `lang:` — aqui, sem `lang:` declarado, vale o default do framework (inglês).
		// Cravar o título em português fazia o teste exigir do guia o que a ferramenta
		// NÃO gera, que é a divergência que este arquivo existe para impedir.
		"## Open Decisions",
		"anchors new spec", // o comando que resolve, ANTES da doutrina
		"--list-sections",  // e o que lista as seções com critério de escolha
	} {
		if !strings.Contains(g, obrigatorio) {
			t.Errorf("o guide precisa conter %q — é o que falta no embutido", obrigatorio)
		}
	}
}

func TestSpecGuideUsaODialetoDoProjeto(t *testing.T) {
	// Um guide genérico não serve: o projeto declara as letras que USA, e um agente que lê
	// as canônicas escreve códigos que o gate `rule-types` reprova.
	cfg := &config.Config{
		RuleTypes:   []config.RuleType{{Letter: "I", Term: "Invariante"}},
		CodeLengths: []int{4},
	}
	g := renderSpecGuide(cfg, "ABCDX")
	if !strings.Contains(g, "Invariante") || !strings.Contains(g, "`I`") {
		t.Error("as letras declaradas em rule_types têm de aparecer no guide")
	}
	if !strings.Contains(g, "4 character") {
		t.Error("o comprimento declarado em code_lengths tem de aparecer")
	}
	// E o canônico NÃO deve ser oferecido quando o projeto declarou o seu.
	if strings.Contains(g, "canonical letters") {
		t.Error("projeto com rule_types não deve receber a lista canônica — ela contradiz a dele")
	}
}

func TestSpecGuideSemRuleTypesOferecOCanonico(t *testing.T) {
	g := renderSpecGuide(&config.Config{}, "ABCDX")
	if !strings.Contains(g, "canonical letters") {
		t.Error("sem rule_types, o guide tem de dizer quais letras valem")
	}
}
