package config

import "testing"

// O silêncio resolve para o canônico: um projeto que já rodava não pode mudar de
// comportamento por causa da chave nova.
func TestMutationFormat_DefaultEhMTE(t *testing.T) {
	c := &Config{Gates: []Gate{{Name: "mutation-score", Check: "mutation-score"}}}
	if got := c.MutationFormat(); got != FormatMTE {
		t.Errorf("MutationFormat() = %q, esperado %q", got, FormatMTE)
	}
	// Sem gate nenhum, idem.
	if got := (&Config{}).MutationFormat(); got != FormatMTE {
		t.Errorf("sem gate: MutationFormat() = %q, esperado %q", got, FormatMTE)
	}
}

func TestMutationFormat_Declarado(t *testing.T) {
	c := &Config{Gates: []Gate{
		{Name: "line-coverage", Check: "line-coverage"},
		{Name: "mutation-score", Check: "mutation-score", Format: "gremlins"},
	}}
	if got := c.MutationFormat(); got != FormatGremlins {
		t.Errorf("MutationFormat() = %q, esperado %q", got, FormatGremlins)
	}
}

// Caixa e espaço não podem decidir o formato — o yaml é escrito à mão.
func TestMutationFormat_Normaliza(t *testing.T) {
	c := &Config{Gates: []Gate{{Name: "mutation-score", Format: "  GREMLINS  "}}}
	if got := c.MutationFormat(); got != FormatGremlins {
		t.Errorf("MutationFormat() = %q, esperado %q", got, FormatGremlins)
	}
}

// O `format:` de OUTRO gate não vale para a mutação — a chave é lida do gate certo.
func TestMutationFormat_IgnoraOutroGate(t *testing.T) {
	c := &Config{Gates: []Gate{{Name: "line-coverage", Format: "gremlins"}}}
	if got := c.MutationFormat(); got != FormatMTE {
		t.Errorf("MutationFormat() = %q, esperado %q (o format é de outro gate)", got, FormatMTE)
	}
}
