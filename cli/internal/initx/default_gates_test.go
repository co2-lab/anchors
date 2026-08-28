package initx

import "testing"

func TestDefaultGates(t *testing.T) {
	// projeto com spec+feature+test → gates de teste presentes
	g := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true}, false)
	names := map[string]bool{}
	for _, x := range g {
		names[x.Name] = true
		if x.IsBlocking() {
			t.Errorf("gate %s deveria nascer informativo", x.Name)
		}
	}
	for _, want := range []string{"spec-completa", "feature-nao-vazia", "tests-green", "line-coverage", "scenario-coverage"} {
		if !names[want] {
			t.Errorf("gate padrão %q faltando", want)
		}
	}
}

func TestDefaultGatesNoScenarioWithoutSpec(t *testing.T) {
	// test sem spec → não semeia scenario-coverage (que cruza os dois)
	g := DefaultGates(map[string]bool{"test": true}, false)
	for _, x := range g {
		if x.Name == "scenario-coverage" {
			t.Error("scenario-coverage não deveria existir sem spec")
		}
	}
}

func TestDefaultGatesEmpty(t *testing.T) {
	if len(DefaultGates(map[string]bool{}, false)) != 0 {
		t.Error("projeto sem artefatos não deveria semear gates")
	}
}

// Num projeto NOVO os gates nascem bloqueantes, e a razão é que a premissa da maturação
// (QUALITY §7 — "impor o gate como bloqueante pararia o projeto") não se aplica: não há
// débito a acomodar. O gate não para nada; ele impede o PRIMEIRO desvio, que é quando
// corrigir custa menos.
func TestProjetoNovoNasceComGatesBloqueantes(t *testing.T) {
	novos := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true}, true)
	existentes := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true}, false)

	if len(novos) != len(existentes) {
		t.Fatalf("a lista de gates não muda com a idade do projeto: %d vs %d", len(novos), len(existentes))
	}

	var bloqueantes int
	for _, g := range novos {
		if g.Blocking != nil && *g.Blocking {
			bloqueantes++
		}
	}
	if bloqueantes == 0 {
		t.Fatal("projeto novo não ganhou nenhum gate bloqueante")
	}
	// E num projeto existente, nenhum: a maturação é o caminho de lá.
	for _, g := range existentes {
		if g.Blocking != nil && *g.Blocking {
			t.Errorf("projeto existente nasceu com %q bloqueante — pararia o trabalho no dia um", g.Name)
		}
	}
}

// Os gates que dependem de sinal INGERIDO (teste, cobertura, mutação) ficam informativos
// mesmo em projeto novo: sem relatório, eles barrariam por AUSÊNCIA DE DADO, e não por
// defeito — o commit falharia antes de a suíte sequer existir.
func TestGateSemSinalNaoBloqueiaNemEmProjetoNovo(t *testing.T) {
	g := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true, "code": true}, true)

	for _, gate := range g {
		if !dependemDeSinalIngerido[gate.Name] {
			continue
		}
		if gate.Blocking != nil && *gate.Blocking {
			t.Errorf("%q depende de sinal ingerido e nasceu bloqueante — barraria por falta de dado", gate.Name)
		}
	}
}
