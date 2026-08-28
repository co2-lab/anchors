package main

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gate"
)

func nomes(gs []config.Gate) []string {
	out := make([]string, 0, len(gs))
	for _, g := range gs {
		out = append(out, g.Name)
	}
	return out
}

// `skip_on` é lista de EXCLUSÃO: quem não declara nada roda nas duas perspectivas. Se
// fosse lista de inclusão, todo gate escrito antes deste campo existir deixaria de
// rodar — o silêncio nunca pode desligar verificação.
func TestFiltrarGatesSkipOnEhPermissivoPorOmissao(t *testing.T) {
	gates := []config.Gate{{Name: "sem-declaracao"}}

	for _, p := range []string{config.PerspectiveChange, config.PerspectiveAll} {
		got := filtrarGates(gates, "", "", false, p, gate.Dispensa{})
		if len(got) != 1 {
			t.Errorf("perspectiva %q: gate sem `skip_on` deve rodar, veio %v", p, nomes(got))
		}
	}
}

// O caso de uso que motivou o campo: um gate que só responde bem sobre a foto completa
// (um detector de órfãos perguntado sobre um recorte acusaria tudo que o recorte não
// alcança) se declara fora do `--changed`.
func TestFiltrarGatesSkipOnChange(t *testing.T) {
	gates := []config.Gate{
		{Name: "so-no-all", SkipOn: []string{config.PerspectiveChange}},
		{Name: "sempre"},
	}

	noChange := filtrarGates(gates, "", "", false, config.PerspectiveChange, gate.Dispensa{})
	if len(noChange) != 1 || noChange[0].Name != "sempre" {
		t.Errorf("`skip_on: [change]` deve sair do --changed, veio %v", nomes(noChange))
	}

	noAll := filtrarGates(gates, "", "", false, config.PerspectiveAll, gate.Dispensa{})
	if len(noAll) != 2 {
		t.Errorf("`skip_on: [change]` deve continuar no --all, veio %v", nomes(noAll))
	}
}

// O inverso: gate caro demais para a varredura completa sai do `--all` sem sair do
// commit. É o que separa este eixo do `cost: slow`, que tira o gate dos dois.
func TestFiltrarGatesSkipOnAll(t *testing.T) {
	gates := []config.Gate{{Name: "so-no-recorte", SkipOn: []string{config.PerspectiveAll}}}

	if got := filtrarGates(gates, "", "", false, config.PerspectiveAll, gate.Dispensa{}); len(got) != 0 {
		t.Errorf("`skip_on: [all]` deve sair do --all, veio %v", nomes(got))
	}
	if got := filtrarGates(gates, "", "", false, config.PerspectiveChange, gate.Dispensa{}); len(got) != 1 {
		t.Errorf("`skip_on: [all]` deve continuar no --changed, veio %v", nomes(got))
	}
}

// As duas juntas desligam o gate — é declaração explícita, não acidente, e por isso o
// filtro a respeita em vez de tratar como contradição.
func TestFiltrarGatesSkipOnAmbasDesliga(t *testing.T) {
	gates := []config.Gate{{
		Name:   "desligado",
		SkipOn: []string{config.PerspectiveChange, config.PerspectiveAll},
	}}

	for _, p := range []string{config.PerspectiveChange, config.PerspectiveAll} {
		if got := filtrarGates(gates, "", "", false, p, gate.Dispensa{}); len(got) != 0 {
			t.Errorf("perspectiva %q: deveria estar desligado, veio %v", p, nomes(got))
		}
	}
}

// Os eixos são INDEPENDENTES: declarar perspectiva não pode alterar o efeito de
// fase/categoria/custo, senão a categorização deixaria de ser adotável aos poucos.
func TestFiltrarGatesEixosIndependentes(t *testing.T) {
	gates := []config.Gate{
		{Name: "lento", Cost: "slow", SkipOn: []string{config.PerspectiveChange}},
		{Name: "rapido"},
	}

	// no --all o lento entra; com --skip-slow, sai por CUSTO, não por perspectiva.
	if got := filtrarGates(gates, "", "", false, config.PerspectiveAll, gate.Dispensa{}); len(got) != 2 {
		t.Errorf("sem skip-slow os dois entram no --all, veio %v", nomes(got))
	}
	if got := filtrarGates(gates, "", "", true, config.PerspectiveAll, gate.Dispensa{}); len(got) != 1 {
		t.Errorf("skip-slow deve remover o lento, veio %v", nomes(got))
	}
}
