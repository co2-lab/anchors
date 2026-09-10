package main

import (
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// As LETRAS do catálogo de seções e as CANÔNICAS do engine são duas metades da mesma
// decisão, e divergiram sem que nada acusasse: o catálogo declarava `invariants` com `I` e
// `errors` com `E`, e nenhuma das duas estava em DefaultRuleLetters.
//
// O efeito era o framework emitir uma spec que ele mesmo recusa — bastava escolher
// Invariantes no `anchors new` para nascer uma regra invisível para a rastreabilidade.
//
// Este teste é o que impede a volta: mexer numa lista sem a outra passa a falhar aqui, e
// não numa spec de usuário meses depois.
func TestCatalogoNaoUsaLetraForaDasCanonicas(t *testing.T) {
	canonicas := map[rune]bool{}
	for _, l := range config.DefaultRuleLetters {
		canonicas[l] = true
	}
	visto := map[string]bool{}
	for _, tpl := range templates {
		for _, s := range tpl.sections {
			if s.Realizes == "" || visto[s.Realizes] {
				continue
			}
			visto[s.Realizes] = true
			for _, l := range s.Realizes {
				if !canonicas[l] {
					t.Errorf("a seção %q realiza a letra %q, que não está nas canônicas (%s) — "+
						"uma spec gerada pelo `anchors new` nasceria com regra invisível para a "+
						"rastreabilidade", s.Key, s.Realizes, config.DefaultRuleLetters)
				}
			}
		}
	}
	if len(visto) == 0 {
		t.Fatal("nenhuma seção com Realizes encontrada — o teste não confrontou nada")
	}
}
