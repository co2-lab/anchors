package main

import (
	"encoding/json"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
	"testing"
)

// O corpo da proteção precisa CHEGAR ao gh. Ele já foi montado e descartado com um
// `_ = body`, e o comando rodava com `--input -` sem nada no stdin: a API recebia corpo
// vazio e respondia 422 reclamando de campo obrigatório nulo. O `--fix` imprimia
// "não deu para proteger os branches" e a porta ficava aberta — que é o silêncio exato
// que a proteção existia para fechar.
func TestCorpoDeProtecaoTemOsCamposObrigatorios(t *testing.T) {
	var corpo map[string]any
	if err := json.Unmarshal([]byte(corpoDeProtecao(1)), &corpo); err != nil {
		t.Fatalf("o corpo não é JSON válido: %v", err)
	}
	// A API recusa (422) se qualquer um destes faltar, mesmo que o valor seja nulo.
	for _, campo := range []string{"required_status_checks", "enforce_admins",
		"required_pull_request_reviews", "restrictions"} {
		if _, ok := corpo[campo]; !ok {
			t.Errorf("falta `%s` — a API responde 422 sem ele", campo)
		}
	}
	// O objeto `required_pull_request_reviews` é o que EXIGE o PR — sem ele a proteção
	// existe e o push direto continua passando.
	rev, ok := corpo["required_pull_request_reviews"].(map[string]any)
	if !ok {
		t.Fatal("required_pull_request_reviews deveria ser objeto — é o que exige o PR")
	}
	// O NÚMERO é o que a chamada passa: era zero fixo, e passou a ser configurável com
	// padrão 1 (ver TestAprovacoesExigidasPadraoEhUma).
	if n, _ := rev["required_approving_review_count"].(float64); n != 1 {
		t.Errorf("aprovações exigidas = %v, esperado o que foi passado (1)", n)
	}
}

// O PADRÃO é UMA aprovação, e a distinção importa: era zero, e com zero o merge acontece
// sem que ninguém revise — a etapa que o card `ready-to-review` representa vira
// decoração. Medido: três PRs mesclados direto, com o board dizendo "esperando revisor"
// sobre trabalho que já estava na develop.
func TestAprovacoesExigidasPadraoEhUma(t *testing.T) {
	var nulo *config.Workflow
	if got := nulo.AprovacoesExigidas(); got != 1 {
		t.Errorf("sem config, o padrão deveria ser 1, veio %d", got)
	}
	if got := (&config.Workflow{}).AprovacoesExigidas(); got != 1 {
		t.Errorf("sem declarar, o padrão deveria ser 1, veio %d", got)
	}
	// ZERO declarado é deliberado e precisa valer: há projetos onde a revisão acontece
	// fora do GitHub. O ponteiro é o que distingue "não declarou" de "declarou zero" —
	// com int simples, o zero-value seria indistinguível da ausência.
	zero := 0
	if got := (&config.Workflow{RequiredApprovals: &zero}).AprovacoesExigidas(); got != 0 {
		t.Errorf("zero declarado deveria valer, veio %d", got)
	}
	// E o corpo enviado à API reflete o número.
	if !strings.Contains(corpoDeProtecao(2), `"required_approving_review_count":2`) {
		t.Error("o corpo deveria carregar o número de aprovações")
	}
}
