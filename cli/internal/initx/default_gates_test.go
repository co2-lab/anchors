package initx

import "testing"

func TestDefaultGates(t *testing.T) {
	// projeto com spec+feature+test → gates de teste presentes
	g := DefaultGates(map[string]bool{"spec": true, "feature": true, "test": true})
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
	g := DefaultGates(map[string]bool{"test": true})
	for _, x := range g {
		if x.Name == "scenario-coverage" {
			t.Error("scenario-coverage não deveria existir sem spec")
		}
	}
}

func TestDefaultGatesEmpty(t *testing.T) {
	if len(DefaultGates(map[string]bool{})) != 0 {
		t.Error("projeto sem artefatos não deveria semear gates")
	}
}
