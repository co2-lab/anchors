package config

import "testing"

// O default é agnóstico: sem declaração, o Anchors não assume fluxo nenhum além do que o
// GitHub já dá a um repositório novo.
func TestBranchDefaultNaoAssumeFluxo(t *testing.T) {
	var w *Workflow
	if b := w.BranchDeIntegracao(); b != "main" {
		t.Errorf("sem config, o branch de integração é `main`, veio %q", b)
	}
	if p := w.BranchesProtegidos(); len(p) != 1 || p[0] != "main" {
		t.Errorf("sem config, só a main é protegida, veio %v", p)
	}
}

// O fluxo develop→staging→main é CONFIGURAÇÃO de um projeto, não regra do framework.
// Neste arranjo o trabalho chega em `develop` — e é lá que o card nasce, não na main.
func TestFluxoDeTresBranchesEhConfiguravel(t *testing.T) {
	w := &Workflow{
		IntegrationBranch: "develop",
		ProtectedBranches: []string{"develop", "staging", "main"},
	}

	if b := w.BranchDeIntegracao(); b != "develop" {
		t.Errorf("o trabalho chega em develop, veio %q", b)
	}
	if p := w.BranchesProtegidos(); len(p) != 3 {
		t.Errorf("os três branches são portas, veio %v", p)
	}
}

// Declarar só o branch de integração implica proteger também a main: um projeto que
// trabalha em `develop` certamente não quer push direto na produção.
func TestIntegracaoNaoMainProtegeMainTambem(t *testing.T) {
	w := &Workflow{IntegrationBranch: "develop"}

	p := w.BranchesProtegidos()
	tem := map[string]bool{}
	for _, b := range p {
		tem[b] = true
	}
	if !tem["develop"] || !tem["main"] {
		t.Errorf("esperava develop e main protegidos, veio %v", p)
	}
}
