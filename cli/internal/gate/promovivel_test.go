package gate

import "testing"

func perfil(gs map[string]GateSummary) Profile {
	return Profile{ByGate: gs}
}

// O caso que motiva a função: um gate informativo com tudo aprovado mede e NÃO DEFENDE.
// Enquanto ele não for promovido, nada impede o próximo commit de desfazer o que já
// está conforme.
func TestGateInformativoLimpoEhPromovivel(t *testing.T) {
	p := perfil(map[string]GateSummary{
		"spec-completa": {Pass: 12, Fail: 0, Blocking: false},
	})

	prom := GatesPromoviveis(p)

	if len(prom) != 1 {
		t.Fatalf("esperava 1 promovível, veio %d", len(prom))
	}
	if prom[0].Gate != "spec-completa" || prom[0].Passou != 12 {
		t.Errorf("promovível errado: %+v", prom[0])
	}
}

// Um gate que REPROVA não é promovível: promovê-lo barraria o trabalho na hora, o que é
// o oposto de uma sugestão útil.
func TestGateQueReprovaNaoEhPromovivel(t *testing.T) {
	p := perfil(map[string]GateSummary{
		"trinca-completa": {Pass: 8, Fail: 3, Blocking: false},
	})

	if prom := GatesPromoviveis(p); len(prom) != 0 {
		t.Errorf("gate com reprovação não deve ser sugerido: %+v", prom)
	}
}

// O caso mais importante: zero aprovações não é "limpo", é SEM DADO. Sugerir promoção
// aqui daria a impressão de defesa que não existe — o mesmo tipo de silêncio que o
// Anchors combate em `DirtyCount` e no `doctor`.
func TestGateSemNadaMedidoNaoEhPromovivel(t *testing.T) {
	p := perfil(map[string]GateSummary{
		"mutation-score": {Pass: 0, Fail: 0, Skip: 40, Blocking: false},
	})

	if prom := GatesPromoviveis(p); len(prom) != 0 {
		t.Errorf("gate que nunca mediu nada não está limpo, está sem dado: %+v", prom)
	}
}

// O bloqueante já defende — não há o que sugerir.
func TestGateBloqueanteNaoEhSugerido(t *testing.T) {
	p := perfil(map[string]GateSummary{
		"layer-boundary": {Pass: 200, Fail: 0, Blocking: true},
	})

	if prom := GatesPromoviveis(p); len(prom) != 0 {
		t.Errorf("bloqueante já defende: %+v", prom)
	}
}
