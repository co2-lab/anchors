package main

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// O VALOR ANTIGO `dispensado` continua aceito, e vira `waived`.
//
// Ele já está em mapas commitados (`verdict: dispensado`) e em scripts. Recusá-lo faria
// o `check` reperguntar julgamentos que alguém já respondeu — com o carimbo ali, visível
// no arquivo, sem nada os ligando.
func TestVeredito_dispensadoAntigoAindaVale(t *testing.T) {
	if err := validateVerdict("dispensado", "a spec declara @TBD"); err != nil {
		t.Errorf("o valor antigo foi recusado: %v", err)
	}
	if err := validateVerdict("waived", "a spec declara @TBD"); err != nil {
		t.Errorf("o valor novo foi recusado: %v", err)
	}
}

// E o motivo continua obrigatório nos dois — a compatibilidade não pode afrouxar a régua.
func TestVeredito_dispensadoAntigoTambemExigeMotivo(t *testing.T) {
	if err := validateVerdict("dispensado", ""); err == nil {
		t.Error("o valor antigo sem motivo foi aceito — a compatibilidade afrouxou a régua")
	}
}

// O `Load` canoniza os nomes de gate, então o nome que a pessoa LÊ no anchors.yaml não é
// o que fica em `g.Name`. O `judge` comparava só contra o canônico e recusava o outro —
// medido no blue-eyes: `--gate mock-detect-cobre-o-dialeto` respondia "gate não existe"
// com o gate declarado três linhas acima no próprio arquivo.
func TestFindJudgmentGate_aceitaOsDoisNomes(t *testing.T) {
	cfg := &config.Config{Gates: []config.Gate{
		{Name: "mock-detect-covers-dialect", Measures: "judgment"},
	}}
	for _, nome := range []string{
		"mock-detect-covers-dialect",  // o canônico, que o mapa guarda
		"mock-detect-cobre-o-dialeto", // o que está ESCRITO no anchors.yaml do projeto
	} {
		if _, ok := findJudgmentGate(cfg, nome); !ok {
			t.Errorf("findJudgmentGate(%q) = false; o gate está declarado", nome)
		}
	}
	if _, ok := findJudgmentGate(cfg, "gate-que-nao-existe"); ok {
		t.Error("aceitou um gate inexistente — a canonização não pode virar vale-tudo")
	}
}
