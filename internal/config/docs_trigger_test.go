package config

import "testing"

func reguaDeDocs() *Config {
	return &Config{Docs: &Docs{Required: []DocArtifact{
		{Kind: KindSchema, Path: "docs/dados.md", Trigger: []string{"DTSTD"}},
		{Kind: KindOpenAPI, Path: "docs/api.yaml", Trigger: []string{"lambdas"}},
		{Kind: KindC4, Path: "docs/arquitetura.md"}, // sem trigger
	}}}
}

// A CAMADA ERRA SOZINHA, e é por isso que o trigger aceita unidade.
//
// Medido no projeto de referência: a camada `infra` tem NOVE unidades, e apenas uma toca
// esquema de dados — as outras oito são API Gateway, autenticação, mTLS, uma régua de
// processo. Cobrar o esquema das oito ensina o agente a ignorar o aviso, que é o pior
// resultado possível: o gate continua lá e ninguém o lê.
func TestRequiredFor_aUnidadeERAMaisPrecisaQueACamada(t *testing.T) {
	c := reguaDeDocs()

	// A unidade que TOCA esquema é cobrada.
	if d := c.RequiredFor("infra", "DTSTD"); len(d) != 1 || d[0].Kind != KindSchema {
		t.Errorf("DTSTD (o banco) devia dever o esquema, deu %v", d)
	}
	// A unidade vizinha, na MESMA camada, não é.
	if d := c.RequiredFor("infra", "GLCGL"); len(d) != 0 {
		t.Errorf("GLCGL (uma régua de processo) foi cobrado de %v — ele não toca tabela", d)
	}
}

// O trigger por CAMADA continua valendo: quem declarou `[lambdas]` quer a camada inteira.
func TestRequiredFor_oTriggerPorCamadaContinua(t *testing.T) {
	c := reguaDeDocs()
	for _, code := range []string{"", "QUALQ"} {
		if d := c.RequiredFor("lambdas", code); len(d) != 1 || d[0].Kind != KindOpenAPI {
			t.Errorf("code=%q: lambdas devia dever o OpenAPI, deu %v", code, d)
		}
	}
}

// Sem código, a resposta é a da CAMADA — é o que o `--layer` pergunta, em abstrato.
func TestRequiredFor_semCodigoRespondePelaCamada(t *testing.T) {
	c := reguaDeDocs()
	if d := c.RequiredFor("infra"); len(d) != 0 {
		t.Errorf("a camada `infra` não tem trigger por camada, e deu %v", d)
	}
}

// Doc SEM trigger não é cobrada de ninguém: é o caso do C4, que muda quando a ESTRUTURA
// muda — não quando uma unidade muda.
func TestRequiredFor_semTriggerNaoEhCobrada(t *testing.T) {
	c := reguaDeDocs()
	for _, par := range [][2]string{{"infra", "DTSTD"}, {"lambdas", "X"}, {"", ""}} {
		for _, d := range c.RequiredFor(par[0], par[1]) {
			if d.Kind == KindC4 {
				t.Errorf("o C4 foi cobrado de %v — ele não tem trigger", par)
			}
		}
	}
}
