package gate

import "testing"

func TestBacktickedSymbols(t *testing.T) {
	// o scan PRESERVA as crases do campo Método — elas são o sinal de "isto é um
	// SÍMBOLO" (verificável), separando-o da prosa descritiva. Formas com as bordas
	// já removidas seguem aceitas (mapas gravados antes dessa mudança).
	cases := map[string][]string{
		"`resolveVersion`":              {"resolveVersion"},              // forma canônica
		"`resolveVersion`, `applyEdit`": {"resolveVersion", "applyEdit"}, // 2 símbolos
		"resolveVersion":                nil,                             // SEM crase = prosa, não contrato
		"resolveVersion`, `applyEdit":   {"resolveVersion", "applyEdit"}, // bordas antigas
		"extractCaller`, `parseBody":    {"extractCaller", "parseBody"},  // bordas antigas
		"CRUD + consultas":              nil,                             // prosa → nada
		"CRUD de member/grupo":          nil,                             // prosa
		"":                              nil,                             // vazio
		"`RecurrenceRuleLike` (tipo)":   {"RecurrenceRuleLike"},          // símbolo + anotação
	}
	for in, want := range cases {
		got := backtickedSymbols(in)
		if len(got) != len(want) {
			t.Errorf("backtickedSymbols(%q) = %v, quer %v", in, got, want)
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("backtickedSymbols(%q)[%d] = %q, quer %q", in, i, got[i], want[i])
			}
		}
	}
}

func TestSymbolUsed(t *testing.T) {
	code := "import { resolveVersion } from './x'\nconst v = resolveVersion(h, m)\n"
	if !symbolUsed(code, "resolveVersion") {
		t.Error("resolveVersion deveria ser usado")
	}
	if symbolUsed(code, "applyEdit") {
		t.Error("applyEdit NÃO está no código")
	}
	// fronteira de palavra: 'resolve' não casa 'resolveVersion'
	if symbolUsed(code, "resolve") {
		t.Error("'resolve' não deveria casar 'resolveVersion' (fronteira)")
	}
}

func TestBacktickedSymbols_prosaNaoEhContrato(t *testing.T) {
	// células SEM crase são descrição livre — não viram promessa (o furo que
	// transformava "aportes"/"solicitações" em símbolo prometido).
	for _, prose := range []string{
		"solicitações", "aportes", "parcelas", "CRUD de IRDependent",
		"contagem de orçamentos", "leitura/atualização",
	} {
		if got := backtickedSymbols(prose); got != nil {
			t.Errorf("backtickedSymbols(%q) = %v, quer nil (prosa não é contrato)", prose, got)
		}
	}
	// com crase, continua sendo contrato
	if got := backtickedSymbols("`pingHeartbeat`"); len(got) != 1 || got[0] != "pingHeartbeat" {
		t.Errorf("símbolo crasado deveria ser extraído, got %v", got)
	}
}

func TestNearestSymbol(t *testing.T) {
	code := "import { pingHeartbeat } from './x'\nconst r = computeSeatsAmountCents(a)\n"
	if got := nearestSymbol(code, "ping"); got != "pingHeartbeat" {
		t.Errorf("nearestSymbol(ping) = %q, quer pingHeartbeat", got)
	}
	if got := nearestSymbol(code, "computeSeatsAmount"); got != "computeSeatsAmountCents" {
		t.Errorf("nearestSymbol(computeSeatsAmount) = %q, quer computeSeatsAmountCents", got)
	}
	// sem parentesco → sem palpite
	if got := nearestSymbol(code, "totalmenteOutro"); got != "" {
		t.Errorf("esperava sem sugestão, got %q", got)
	}
}
