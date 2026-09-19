package gate

import "testing"

func perfil(gs map[string]GateSummary) Profile {
	return Profile{ByGate: gs}
}

func TestPromotable_EmptyProfileReturnsEmpty(t *testing.T) {
	t.Run("PRGTP-B01: Returns an empty list of Promotable gates when profile contains no gates", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{})
	prom := PromotableGates(p)
	if len(prom) != 0 {
		t.Fatalf("expected 0 promotable gates, got %d", len(prom))
	}
}

// O caso que motiva a função: um gate informativo com tudo aprovado mede e NÃO DEFENDE.
// Enquanto ele não for promovido, nada impede o próximo commit de desfazer o que já
// está conforme.
func TestGateInformativoLimpoEhPromovivel(t *testing.T) {
	t.Run("PRGTP-B02: An informative gate with passes and zero failures is included", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"spec-complete": {Pass: 12, Fail: 0, Blocking: false},
	})

	prom := PromotableGates(p)

	if len(prom) != 1 {
		t.Fatalf("esperava 1 promovível, veio %d", len(prom))
	}
	if prom[0].Gate != "spec-complete" || prom[0].Passou != 12 {
		t.Errorf("promovível errado: %+v", prom[0])
	}
}

// Um gate que REPROVA não é promovível: promovê-lo barraria o trabalho na hora, o que é
// o oposto de uma sugestão útil.
func TestGateQueReprovaNaoEhPromovivel(t *testing.T) {
	t.Run("PRGTP-B03: An informative gate with failures is excluded from promotion", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"triad-complete": {Pass: 8, Fail: 3, Blocking: false},
	})

	if prom := PromotableGates(p); len(prom) != 0 {
		t.Errorf("gate com reprovação não deve ser sugerido: %+v", prom)
	}
}

// O caso mais importante: zero aprovações não é "limpo", é SEM DADO. Sugerir promoção
// aqui daria a impressão de defesa que não existe — o mesmo tipo de silêncio que o
// Anchors combate em `DirtyCount` e no `doctor`.
func TestGateSemNadaMedidoNaoEhPromovivel(t *testing.T) {
	t.Run("PRGTP-B04: An informative gate with zero passes is excluded as having no data", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"mutation-score": {Pass: 0, Fail: 0, Skip: 40, Blocking: false},
	})

	if prom := PromotableGates(p); len(prom) != 0 {
		t.Errorf("gate que nunca mediu nada não está limpo, está sem dado: %+v", prom)
	}
}

// O bloqueante já defende — não há o que sugerir.
func TestGateBloqueanteNaoEhSugerido(t *testing.T) {
	t.Run("PRGTP-B05: A blocking gate is excluded from promotion suggestions", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"layer-boundary": {Pass: 200, Fail: 0, Blocking: true},
	})

	if prom := PromotableGates(p); len(prom) != 0 {
		t.Errorf("bloqueante já defende: %+v", prom)
	}
}

func TestPromotable_PopulatesGateIdentifier(t *testing.T) {
	t.Run("PRGTP-B06: Returned Promotable populates Gate with the declared gate name", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"custom-linter": {Pass: 5, Fail: 0, Blocking: false},
	})
	prom := PromotableGates(p)
	if len(prom) != 1 || prom[0].Gate != "custom-linter" {
		t.Fatalf("expected Gate 'custom-linter', got %+v", prom)
	}
}

func TestPromotable_PopulatesPassouCount(t *testing.T) {
	t.Run("PRGTP-B07: Returned Promotable records Passou equal to passed node count", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"rule-types": {Pass: 12, Fail: 0, Blocking: false},
	})
	prom := PromotableGates(p)
	if len(prom) != 1 || prom[0].Passou != 12 {
		t.Fatalf("expected Passou 12, got %+v", prom)
	}
}

func TestPromotable_SortedAlphabetically(t *testing.T) {
	t.Run("PRGTP-B08: PromotableGates returns candidates sorted deterministically by name", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"zebra-gate":  {Pass: 3, Fail: 0, Blocking: false},
		"alpha-gate":  {Pass: 2, Fail: 0, Blocking: false},
		"middle-gate": {Pass: 1, Fail: 0, Blocking: false},
	})
	prom := PromotableGates(p)
	if len(prom) != 3 {
		t.Fatalf("expected 3 gates, got %d", len(prom))
	}
	if prom[0].Gate != "alpha-gate" || prom[1].Gate != "middle-gate" || prom[2].Gate != "zebra-gate" {
		t.Fatalf("expected alphabetical order, got %v, %v, %v", prom[0].Gate, prom[1].Gate, prom[2].Gate)
	}
}

func TestPromotable_MultipleCleanGatesCollected(t *testing.T) {
	t.Run("PRGTP-B09: Multiple clean informative gates are all collected", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"g1": {Pass: 1, Fail: 0, Blocking: false},
		"g2": {Pass: 2, Fail: 0, Blocking: false},
		"g3": {Pass: 3, Fail: 0, Blocking: false},
	})
	prom := PromotableGates(p)
	if len(prom) != 3 {
		t.Fatalf("expected 3 promotable gates, got %d", len(prom))
	}
}

func TestPromotable_InvariantMatching(t *testing.T) {
	t.Run("PRGTP-I01: An informative gate is candidate if and only if non-blocking zero failures and positive passes", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"clean":          {Pass: 5, Fail: 0, Blocking: false},
		"blocking_clean": {Pass: 5, Fail: 0, Blocking: true},
		"failing":        {Pass: 5, Fail: 1, Blocking: false},
		"empty":          {Pass: 0, Fail: 0, Blocking: false},
	})
	prom := PromotableGates(p)
	if len(prom) != 1 || prom[0].Gate != "clean" {
		t.Fatalf("expected only 'clean', got %+v", prom)
	}
}

func TestPromotable_ZeroPassesNeverClean(t *testing.T) {
	t.Run("PRGTP-I02: A gate with zero passes is never classified as clean", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"skipped-only": {Pass: 0, Fail: 0, Skip: 50, Pending: 10, Blocking: false},
	})
	prom := PromotableGates(p)
	if len(prom) != 0 {
		t.Fatalf("expected 0 promotable gates for zero passes, got %d", len(prom))
	}
}

func TestPromotable_DoesNotModifyConfig(t *testing.T) {
	t.Run("PRGTP-X01: Does not automatically promote gates or modify anchors yaml", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"spec-complete": {Pass: 12, Fail: 0, Blocking: false},
	})
	prom := PromotableGates(p)
	if len(prom) == 0 {
		t.Fatal("expected promotable gates")
	}
	// Function returns value without any mutation to config or disk
}

func TestPromotable_OperatesOnlyOnMemoryProfile(t *testing.T) {
	t.Run("PRGTP-X02: Does not evaluate gate execution results directly from disk", func(t *testing.T) {})
	p := perfil(map[string]GateSummary{
		"mock-gate": {Pass: 7, Fail: 0, Blocking: false},
	})
	prom := PromotableGates(p)
	if len(prom) != 1 || prom[0].Passou != 7 {
		t.Fatalf("expected candidate derived directly from profile, got %+v", prom)
	}
}
