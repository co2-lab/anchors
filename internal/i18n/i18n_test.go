package i18n

import (
	"strings"
	"testing"
)

// TODO idioma suportado tem catálogo, e ele CARREGA.
//
// Sem este teste, acrescentar "fr" a `Suportados` sem criar o `fr.json` passaria — e o
// projeto que declarasse `lang: fr` receberia tudo em inglês pelo fallback, sem nada
// avisando que a tradução não existe.
func TestTodoSuportadoTemCatalogo(t *testing.T) {
	for _, lang := range SupportedLangs {
		if len(Keys(lang)) == 0 {
			t.Errorf("idioma %q está em Suportados e não tem catálogo (locales/%s.json)", lang, lang)
		}
	}
}

// OS CATÁLOGOS TÊM AS MESMAS CHAVES.
//
// Uma chave que existe em `en` e falta em `es` sai em inglês pelo fallback — o que é
// melhor que sumir, mas é tradução pendente que ninguém vê. Este teste é o que a torna
// visível: ele falha nomeando a chave e o idioma.
func TestCatalogosTemAsMesmasChaves(t *testing.T) {
	base := map[string]bool{}
	for _, k := range Keys(Default) {
		base[k] = true
	}
	if len(base) == 0 {
		t.Fatal("o catálogo do idioma padrão está vazio — o teste não confrontaria nada")
	}

	for _, lang := range SupportedLangs {
		if lang == Default {
			continue
		}
		tem := map[string]bool{}
		for _, k := range Keys(lang) {
			tem[k] = true
		}
		for k := range base {
			if !tem[k] {
				t.Errorf("%s: falta a chave %q (existe em %s) — ela sairia em %s sem ninguém notar",
					lang, k, Default, Default)
			}
		}
		for k := range tem {
			if !base[k] {
				t.Errorf("%s: tem a chave %q que NÃO existe em %s — ou é sobra, ou falta traduzir para o padrão",
					lang, k, Default)
			}
		}
	}
}

// O FALLBACK devolve a CHAVE quando nem o padrão a tem.
//
// É o terceiro degrau, e o que impede o pior desfecho: uma mensagem vazia. Um erro que
// não diz nada é pior que um erro em outro idioma — pelo menos a chave diz o que falta.
func TestChaveInexistenteDevolveAPropriaChave(t *testing.T) {
	t.Cleanup(func() { _ = Set(Default) })
	_ = Set("pt-BR")
	if got := T("nao.existe.esta.chave"); got != "nao.existe.esta.chave" {
		t.Errorf("chave inexistente devolveu %q — deveria devolver a própria chave", got)
	}
}

func TestTraduzNoIdiomaCorrente(t *testing.T) {
	t.Cleanup(func() { _ = Set(Default) })
	for _, c := range []struct{ lang, contem string }{
		{"pt-BR", "CONGELADO"},
		{"en", "FROZEN"},
		{"es", "CONGELADO"},
	} {
		if err := Set(c.lang); err != nil {
			t.Fatal(err)
		}
		if got := T("freeze.done"); !strings.Contains(got, c.contem) {
			t.Errorf("%s: T(freeze.done) = %q, esperava conter %q", c.lang, got, c.contem)
		}
	}
}

// A INTERPOLAÇÃO funciona igual em todos os idiomas.
//
// Uma tradução que perde o `%s` produziria uma mensagem sem o dado — e o dado costuma ser
// a parte que importa (o nome do arquivo, o motivo do congelamento).
func TestInterpolacaoValeEmTodosOsIdiomas(t *testing.T) {
	t.Cleanup(func() { _ = Set(Default) })
	for _, lang := range SupportedLangs {
		_ = Set(lang)
		got := T("freeze.blocked.reason", "o plano 0002 quebrou")
		if !strings.Contains(got, "o plano 0002 quebrou") {
			t.Errorf("%s: a interpolação perdeu o argumento: %q", lang, got)
		}
	}
}

// A LISTA É FECHADA, e o erro diz quais existem.
//
// Recusar sem dizer as opções obriga quem errou a ir procurar na documentação — e o
// erro mais comum é justamente o código do idioma (`pt` em vez de `pt-BR`).
func TestIdiomaForaDaListaEhRecusadoComAsOpcoes(t *testing.T) {
	t.Cleanup(func() { _ = Set(Default) })
	err := Set("klingon")
	if err == nil {
		t.Fatal("idioma inexistente foi aceito")
	}
	for _, l := range SupportedLangs {
		if !strings.Contains(err.Error(), l) {
			t.Errorf("o erro não lista %q: %v", l, err)
		}
	}
}

// Vazio cai no padrão, e não é erro: um projeto que não declarou `lang:` é o caso comum.
func TestVazioCaiNoPadrao(t *testing.T) {
	t.Cleanup(func() { _ = Set(Default) })
	if err := Set(""); err != nil {
		t.Fatalf("idioma vazio deu erro: %v", err)
	}
	if Current() != Default {
		t.Errorf("vazio virou %q, esperava %q", Current(), Default)
	}
}
