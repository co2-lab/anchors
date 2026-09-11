package config

import (
	"strings"
	"testing"
)

// `min_version` existe para o momento em que uma correção do Anchors precisa alcançar todo
// mundo antes que o trabalho continue. A conferência anterior comparava o binário com o
// `gerado_por` do mapa — campo DERIVADO, reescrito pelo próximo `map build` de qualquer
// agente. Medido: ele voltou a "dev" num projeto onde a release corrente era a v0.1.83.
//
// Aqui a declaração é de uma pessoa, no arquivo que o time versiona.

func TestAtendeMinVersion_aOrdemEhOrdinalNaoIgualdade(t *testing.T) {
	casos := []struct {
		minimo, rodando string
		atende          bool
	}{
		{"0.1.84", "0.1.84", true},  // exatamente o mínimo
		{"0.1.84", "0.1.85", true},  // mais novo atende
		{"0.1.84", "0.2.0", true},   // minor maior
		{"0.1.84", "1.0.0", true},   // major maior
		{"0.1.84", "0.1.83", false}, // patch menor
		{"0.1.84", "0.1.9", false},  // 9 < 84: é NÚMERO, não texto
		{"0.2.0", "0.1.99", false},  // minor manda sobre patch
		{"1.0.0", "0.99.99", false}, // major manda sobre tudo
	}
	for _, c := range casos {
		got, err := AtendeMinVersion(c.minimo, c.rodando)
		if err != nil {
			t.Errorf("min=%s rodando=%s: erro inesperado %v", c.minimo, c.rodando, err)
			continue
		}
		if got != c.atende {
			t.Errorf("min=%s rodando=%s: esperava atende=%v", c.minimo, c.rodando, c.atende)
		}
	}
}

// A comparação por TEXTO daria a resposta errada aqui, e este é o caso que a expõe:
// "0.1.9" > "0.1.84" em ordem lexicográfica, e é MENOR em versão.
func TestAtendeMinVersion_naoComparaComoTexto(t *testing.T) {
	if "0.1.9" <= "0.1.84" {
		t.Skip("a premissa do teste mudou: a ordem lexicográfica deixou de inverter aqui")
	}
	atende, err := AtendeMinVersion("0.1.84", "0.1.9")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if atende {
		t.Error("0.1.9 é MENOR que 0.1.84 — a comparação está sendo feita como texto")
	}
}

// Sem mínimo declarado, todo binário atende: a maioria dos projetos não declara, e exigir a
// declaração para poder trabalhar inverteria o padrão.
func TestAtendeMinVersion_semMinimoTodoMundoAtende(t *testing.T) {
	for _, rodando := range []string{"0.0.1", "dev", "", "qualquer-coisa"} {
		atende, err := AtendeMinVersion("", rodando)
		if err != nil {
			t.Errorf("rodando=%q: sem mínimo não deveria haver erro, veio %v", rodando, err)
		}
		if !atende {
			t.Errorf("rodando=%q: sem mínimo declarado, deveria atender", rodando)
		}
	}
}

// `dev` NÃO É ORDENÁVEL, e fingir que é seria pior que recusar a comparação.
//
// Um build local pode ser mais novo que qualquer release (a árvore de quem desenvolve o
// Anchors) ou mais velho que todas. Quem chama decide o que fazer com o erro; o que não se
// faz é responder "atende" ou "não atende" sem base.
func TestAtendeMinVersion_devNaoEhOrdenavel(t *testing.T) {
	for _, rodando := range []string{"dev", "0.1", "0.1.84-rc1", "v-nada", "1.2.3.4"} {
		_, err := AtendeMinVersion("0.1.84", rodando)
		if err == nil {
			t.Errorf("rodando=%q deveria recusar a comparação, e respondeu sem erro", rodando)
		}
	}
}

// O `v` da tag é aceito nos dois lados: a release se chama `v0.1.84` e o binário se
// reporta como `0.1.84`, e exigir que quem declara saiba qual dos dois usar é o tipo de
// detalhe que se erra uma vez e some.
func TestAtendeMinVersion_aceitaOPrefixoDaTag(t *testing.T) {
	for _, par := range [][2]string{
		{"v0.1.84", "0.1.84"},
		{"0.1.84", "v0.1.84"},
		{"v0.1.84", "v0.1.85"},
	} {
		atende, err := AtendeMinVersion(par[0], par[1])
		if err != nil {
			t.Errorf("min=%s rodando=%s: erro %v", par[0], par[1], err)
		}
		if !atende {
			t.Errorf("min=%s rodando=%s: deveria atender", par[0], par[1])
		}
	}
}

// O CAMPO só aceita `MAJOR.MINOR.PATCH`, e a recusa é na CARGA.
//
// Falhar só na comparação faria o efeito de um valor mal escrito ser o AVISO SUMIR — e o
// campo existe justamente para forçar uma atualização. Um `min_version: latest` produziria
// o silêncio que ele deveria quebrar.
func TestValidarMinVersion_recusaOQueNaoEhMmP(t *testing.T) {
	ruins := []string{
		"dev",        // um binário de desenvolvimento é um fato que se encontra, não uma decisão que se declara
		"latest",     // parece razoável e não se compara com nada
		"0.1",        // falta o patch
		"0.1.84-rc1", // pré-release comparada como final responderia "atende" a quem tem menos
		"1.2.3.4",
		"abc",
		"0.1.x",
	}
	for _, v := range ruins {
		c := &Config{MinVersion: v}
		if err := c.validarMinVersion(); err == nil {
			t.Errorf("min_version=%q deveria ser recusado na carga", v)
		}
	}
}

func TestValidarMinVersion_aceitaOFormatoEOVazio(t *testing.T) {
	for _, v := range []string{"", "0.1.84", "v0.1.84", "1.0.0", "10.20.30"} {
		c := &Config{MinVersion: v}
		if err := c.validarMinVersion(); err != nil {
			t.Errorf("min_version=%q deveria ser aceito, veio: %v", v, err)
		}
	}
}

// A mensagem de erro tem de dizer O FORMATO e A CONSEQUÊNCIA. Um "valor inválido" seco
// manda quem lê adivinhar, e o campo é raro o bastante para ninguém lembrar de cabeça.
func TestValidarMinVersion_aMensagemEnsina(t *testing.T) {
	c := &Config{MinVersion: "latest"}
	err := c.validarMinVersion()
	if err == nil {
		t.Fatal("esperava erro")
	}
	for _, quer := range []string{"MAJOR.MINOR.PATCH", "0.1.84", "silencia"} {
		if !strings.Contains(err.Error(), quer) {
			t.Errorf("a mensagem deveria conter %q; veio: %s", quer, err)
		}
	}
}
