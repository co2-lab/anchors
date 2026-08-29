package main

import (
	"encoding/json"
	"testing"
)

// O corpo da proteção precisa CHEGAR ao gh. Ele já foi montado e descartado com um
// `_ = body`, e o comando rodava com `--input -` sem nada no stdin: a API recebia corpo
// vazio e respondia 422 reclamando de campo obrigatório nulo. O `--fix` imprimia
// "não deu para proteger os branches" e a porta ficava aberta — que é o silêncio exato
// que a proteção existia para fechar.
func TestCorpoDeProtecaoTemOsCamposObrigatorios(t *testing.T) {
	var corpo map[string]any
	if err := json.Unmarshal([]byte(corpoDeProtecao()), &corpo); err != nil {
		t.Fatalf("o corpo não é JSON válido: %v", err)
	}
	// A API recusa (422) se qualquer um destes faltar, mesmo que o valor seja nulo.
	for _, campo := range []string{"required_status_checks", "enforce_admins",
		"required_pull_request_reviews", "restrictions"} {
		if _, ok := corpo[campo]; !ok {
			t.Errorf("falta `%s` — a API responde 422 sem ele", campo)
		}
	}
	// Zero aprovações: exigir revisor travaria um time de uma pessoa com agentes.
	rev, ok := corpo["required_pull_request_reviews"].(map[string]any)
	if !ok {
		t.Fatal("required_pull_request_reviews deveria ser objeto — é o que exige o PR")
	}
	if n, _ := rev["required_approving_review_count"].(float64); n != 0 {
		t.Errorf("aprovações exigidas = %v, esperado 0", n)
	}
}
