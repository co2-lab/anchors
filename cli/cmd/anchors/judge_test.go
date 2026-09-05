package main

import "testing"

// O VALOR ANTIGO `dispensado` continua aceito, e vira `waived`.
//
// Ele já está em mapas commitados (`verdict: dispensado`) e em scripts. Recusá-lo faria
// o `check` reperguntar julgamentos que alguém já respondeu — com o carimbo ali, visível
// no arquivo, sem nada os ligando.
func TestVeredito_dispensadoAntigoAindaVale(t *testing.T) {
	if err := validaVeredito("dispensado", "a spec declara @TBD"); err != nil {
		t.Errorf("o valor antigo foi recusado: %v", err)
	}
	if err := validaVeredito("waived", "a spec declara @TBD"); err != nil {
		t.Errorf("o valor novo foi recusado: %v", err)
	}
}

// E o motivo continua obrigatório nos dois — a compatibilidade não pode afrouxar a régua.
func TestVeredito_dispensadoAntigoTambemExigeMotivo(t *testing.T) {
	if err := validaVeredito("dispensado", ""); err == nil {
		t.Error("o valor antigo sem motivo foi aceito — a compatibilidade afrouxou a régua")
	}
}
