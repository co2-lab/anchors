package i18n

import (
	"reflect"
	"strings"
	"testing"
)

func TestSupportedLanguagesAndDefault(t *testing.T) {
	t.Run("INCTA-B01: The supported languages and the default", func(t *testing.T) {})
	if !reflect.DeepEqual(SupportedLangs, []string{"pt-BR", "en", "es"}) || Default != "en" {
		t.Fatalf("supported %v, default %q; want [pt-BR en es], en", SupportedLangs, Default)
	}
	if IsSupported("pt") || !IsSupported("pt-BR") {
		t.Error("only the exact supported codes are supported")
	}
}

// EVERY supported language has a catalog, and it LOADS; and the catalogs have the SAME KEYS.
//
// Without this, adding "fr" to the list without creating `fr.json` would pass — and a project
// declaring `lang: fr` would get everything in English by fallback, with nothing warning it.
// A key that exists in `en` and is missing in `es` comes out in English — better than vanishing,
// but a pending translation nobody sees. This test makes it visible, naming the key and language.
func TestEveryLanguageHasACatalogWithTheSameKeys(t *testing.T) {
	t.Run("INCTA-I01: Every language has a catalog with the same keys", func(t *testing.T) {})
	for _, lang := range SupportedLangs {
		if len(Keys(lang)) == 0 {
			t.Errorf("language %q is supported and has no catalog (locales/%s.json)", lang, lang)
		}
	}
	base := map[string]bool{}
	for _, k := range Keys(Default) {
		base[k] = true
	}
	if len(base) == 0 {
		t.Fatal("the default language's catalog is empty — the test would confront nothing")
	}

	for _, lang := range SupportedLangs {
		if lang == Default {
			continue
		}
		has := map[string]bool{}
		for _, k := range Keys(lang) {
			has[k] = true
		}
		for k := range base {
			if !has[k] {
				t.Errorf("%s: missing the key %q (exists in %s) — it would come out in %s unnoticed",
					lang, k, Default, Default)
			}
		}
		for k := range has {
			if !base[k] {
				t.Errorf("%s: has the key %q that does NOT exist in %s — leftover, or missing in the default",
					lang, k, Default)
			}
		}
	}
}

// The FALLBACK returns the KEY when not even the default has it — the step that prevents the
// worst outcome: an empty message.
func TestMissingKeyReturnsTheKeyItself(t *testing.T) {
	t.Run("INCTA-B05: A key missing everywhere resolves to itself", func(t *testing.T) {})
	t.Cleanup(func() { _ = Set(Default) })
	_ = Set("pt-BR")
	if got := T("nao.existe.esta.chave"); got != "nao.existe.esta.chave" {
		t.Errorf("a missing key returned %q — it should return the key itself", got)
	}
}

func TestMissingTranslationFallsBackToTheDefault(t *testing.T) {
	t.Run("INCTA-B04: A key missing in the current language falls back to English", func(t *testing.T) {})
	t.Cleanup(func() { _ = Set(Default) })
	es := load("es")
	partial := map[string]string{}
	for k, v := range es {
		if k != "freeze.done" {
			partial[k] = v
		}
	}
	mu.Lock()
	catalogo["es"] = partial
	mu.Unlock()
	t.Cleanup(func() { mu.Lock(); catalogo["es"] = es; mu.Unlock() })

	_ = Set("es")
	if got := T("freeze.done"); !strings.Contains(got, "FROZEN") {
		t.Errorf("a key missing in es = %q, want the English text", got)
	}
}

// The INTERPOLATION works the same in every language: a translation that loses the `%s` would
// produce a message without the datum — usually the part that matters.
func TestTranslatesInTheCurrentLanguage(t *testing.T) {
	t.Run("INCTA-B03: Messages come out in the current language with their arguments", func(t *testing.T) {})
	t.Cleanup(func() { _ = Set(Default) })
	for _, c := range []struct{ lang, holds string }{
		{"pt-BR", "CONGELADO"},
		{"en", "FROZEN"},
		{"es", "CONGELADO"},
	} {
		if err := Set(c.lang); err != nil {
			t.Fatal(err)
		}
		if got := T("freeze.done"); !strings.Contains(got, c.holds) {
			t.Errorf("%s: T(freeze.done) = %q, want it to hold %q", c.lang, got, c.holds)
		}
		if got := T("freeze.blocked.reason", "o plano 0002 quebrou"); !strings.Contains(got, "o plano 0002 quebrou") {
			t.Errorf("%s: the interpolation lost the argument: %q", c.lang, got)
		}
	}
}

// The LIST IS CLOSED, and the error says which exist: refusing without the options sends whoever
// erred to the documentation — and the most common mistake is the code itself (`pt` for `pt-BR`).
func TestUnsupportedLanguageIsRefusedWithTheOptions(t *testing.T) {
	t.Run("INCTA-E01: An unsupported language is refused with the options", func(t *testing.T) {})
	t.Cleanup(func() { _ = Set(Default) })
	_ = Set("en")
	err := Set("klingon")
	if err == nil {
		t.Fatal("a language that does not exist was accepted")
	}
	for _, l := range SupportedLangs {
		if !strings.Contains(err.Error(), l) {
			t.Errorf("the error does not list %q: %v", l, err)
		}
	}
	if Current() != "en" {
		t.Errorf("a refused language changed the current one to %q", Current())
	}
}

// Empty falls back to the default, and is not an error: a project with no `lang:` is the common case.
func TestEmptyFallsBackToTheDefault(t *testing.T) {
	t.Run("INCTA-B02: An empty language is the default", func(t *testing.T) {})
	t.Cleanup(func() { _ = Set(Default) })
	_ = Set("pt-BR")
	if err := Set(""); err != nil {
		t.Fatalf("an empty language gave an error: %v", err)
	}
	if Current() != Default {
		t.Errorf("empty became %q, want %q", Current(), Default)
	}
}

func TestAllTranslations_eachValueOnce(t *testing.T) {
	t.Run("INCTA-B06: A key's values across languages come back once each", func(t *testing.T) {})
	if got := AllTranslations("section.title.states"); !reflect.DeepEqual(got, []string{"Estados", "States"}) {
		t.Errorf("AllTranslations(section.title.states) = %v, want [Estados States]", got)
	}
}

func TestTIn_resolvesInAGivenLanguage(t *testing.T) {
	t.Run("INCTA-B07: A key resolves in a given language without changing the current one", func(t *testing.T) {})
	t.Cleanup(func() { _ = Set(Default) })
	_ = Set("en")
	if got := TIn("pt-BR", "section.title.rules"); got != "Regras" {
		t.Errorf("TIn(pt-BR, rules) = %q, want Regras", got)
	}
	if got := TIn("es", "nao.existe.esta.chave"); got != "" {
		t.Errorf("TIn of a missing key = %q, want empty", got)
	}
	if Current() != "en" {
		t.Errorf("TIn changed the current language to %q", Current())
	}
}

func TestSectionKeyFor_reverseLookup(t *testing.T) {
	t.Run("INCTA-B08: A written title gives back its key and language", func(t *testing.T) {})
	for title, want := range map[string][2]string{
		"  visão geral ":    {"section.title.overview", "pt-BR"},
		"Reglas":            {"section.title.rules", "es"},
		"Coisas do projeto": {"", ""},
		"":                  {"", ""},
	} {
		k, l := SectionKeyFor(title)
		if k != want[0] || l != want[1] {
			t.Errorf("SectionKeyFor(%q) = %q, %q; want %q, %q", title, k, l, want[0], want[1])
		}
	}
}
