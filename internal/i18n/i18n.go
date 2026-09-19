// Package i18n é o catálogo de mensagens do Anchors.
//
// O produto é multi-idioma: o projeto declara `lang:` no `anchors.yaml`, e toda mensagem
// que uma PESSOA lê sai naquele idioma. O que NÃO se traduz é o vocabulário — nomes de
// gate, labels, flags — porque são identificadores, e traduzi-los faria o `anchors.yaml`
// de um projeto deixar de funcionar num time de outro idioma.
//
// A régua é essa: traduz-se o que se LÊ, não o que se ESCREVE na configuração.
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed locales/*.json
var files embed.FS

// SupportedLangs são os idiomas que o Anchors aceita em `lang:`.
//
// A lista é FECHADA de propósito. Aceitar qualquer código produziria projetos declarando
// idiomas que não existem no catálogo, e o erro apareceria só quando alguém lesse a
// primeira mensagem — em inglês, sem entender por quê. Fechar a lista faz o `anchors init`
// e o `check` recusarem na hora, dizendo quais existem.
var SupportedLangs = []string{"pt-BR", "en", "es"}

// Default é o idioma quando o projeto não declara `lang:`.
//
// Inglês, e não português: um projeto que não declarou nada é provavelmente um projeto
// novo de alguém que encontrou o Anchors — e o inglês é o que mais gente lê. Quem quer
// português declara.
const Default = "en"

var (
	mu       sync.RWMutex
	atual    = Default
	catalogo = map[string]map[string]string{}
)

// IsSupported diz se o idioma está na lista fechada.
func IsSupported(lang string) bool {
	for _, l := range SupportedLangs {
		if l == lang {
			return true
		}
	}
	return false
}

// Set troca o idioma corrente. Um idioma fora da lista é recusado com a lista junto —
// quem errou o código precisa saber quais existem, não só que errou.
func Set(lang string) error {
	if lang == "" {
		lang = Default
	}
	if !IsSupported(lang) {
		return fmt.Errorf("%s", T("lang.unsupported", lang, strings.Join(SupportedLangs, ", ")))
	}
	mu.Lock()
	defer mu.Unlock()
	atual = lang
	return nil
}

// Current devolve o idioma corrente.
func Current() string {
	mu.RLock()
	defer mu.RUnlock()
	return atual
}

// T traduz uma chave, interpolando os argumentos como `fmt.Sprintf`.
//
// O FALLBACK é em cascata, e cada degrau existe por um motivo:
//
//  1. o idioma corrente — o caso normal
//  2. o Padrao (en) — uma chave que ainda não foi traduzida sai em inglês, e não some
//  3. a própria chave — se nem em inglês existe, quem vê o texto sabe QUAL chave falta
//
// O terceiro degrau é o que impede o pior desfecho: uma mensagem vazia. Um erro que não
// diz nada é pior que um erro em outro idioma.
func T(chave string, args ...any) string {
	mu.RLock()
	lang := atual
	mu.RUnlock()

	if s, ok := lookup(lang, chave); ok {
		return format(s, args...)
	}
	if lang != Default {
		if s, ok := lookup(Default, chave); ok {
			return format(s, args...)
		}
	}
	return chave
}

func format(s string, args ...any) string {
	if len(args) == 0 {
		return s
	}
	return fmt.Sprintf(s, args...)
}

func lookup(lang, chave string) (string, bool) {
	mu.RLock()
	c, carregado := catalogo[lang]
	mu.RUnlock()
	if !carregado {
		c = load(lang)
	}
	s, ok := c[chave]
	return s, ok
}

// load lê o JSON do idioma UMA vez e o guarda.
//
// Um catálogo ausente vira mapa vazio, não erro: o binário tem de rodar mesmo que um
// arquivo de tradução falte, caindo para o fallback. Falhar aqui derrubaria o CLI inteiro
// por causa de uma tradução.
func load(lang string) map[string]string {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := catalogo[lang]; ok {
		return c
	}
	c := map[string]string{}
	if b, err := files.ReadFile("locales/" + lang + ".json"); err == nil {
		_ = json.Unmarshal(b, &c)
	}
	catalogo[lang] = c
	return c
}

// AllTranslations devolve o valor de uma chave em TODOS os idiomas suportados, sem
// duplicatas.
//
// Existe para o confronto com TEXTO JÁ ESCRITO EM DISCO, que é o caso em que o idioma
// atual não basta: um título de seção de spec foi escrito sob o idioma que o projeto
// tinha NAQUELE dia, e procurar só a tradução corrente acusaria ausência onde a seção
// existe — com outro nome. Vale também para o projeto que trocou de `lang:` depois de
// escrever metade do acervo.
//
// Não confundir com T(): aquela RESOLVE uma mensagem para o leitor; esta reconhece o que
// já está escrito.
func AllTranslations(chave string) []string {
	var out []string
	visto := map[string]bool{}
	for _, lang := range SupportedLangs {
		if s, ok := lookup(lang, chave); ok && s != "" && !visto[s] {
			visto[s] = true
			out = append(out, s)
		}
	}
	return out
}

// TIn resolve uma chave num idioma ESPECÍFICO, sem mexer no idioma atual.
//
// T() serve ao leitor (resolve no idioma corrente); esta serve ao confronto: "como este
// título se escreveria no `lang:` do projeto?". Devolve "" quando a chave não existe
// naquele idioma — quem chama decide o que fazer com a ausência, em vez de receber a
// chave crua de volta como a T() faz.
func TIn(lang, chave string) string {
	if s, ok := lookup(lang, chave); ok {
		return s
	}
	return ""
}

// SectionKeyFor faz o caminho INVERSO do catálogo: dado um título de seção já escrito,
// devolve a chave que ele realiza e o idioma em que está.
//
// Existe para o gate que cobra o idioma dos títulos. Sem ele, a alternativa seria uma
// lista de traduções escrita à mão dentro do gate — que divergiria do catálogo no dia em
// que alguém acrescentasse um idioma, e o gate passaria a acusar como "idioma errado" um
// título perfeitamente válido.
//
// Devolve ("", "") para título que não está no catálogo: é o caso do léxico PRÓPRIO do
// projeto (`## Fora de escopo`), que não é idioma errado — é seção que o framework não
// nomeia, e cobrá-la seria impor o vocabulário do engine ao projeto.
func SectionKeyFor(titulo string) (chave, lang string) {
	t := strings.ToLower(strings.TrimSpace(titulo))
	if t == "" {
		return "", ""
	}
	for _, l := range SupportedLangs {
		for k, v := range load(l) {
			if strings.HasPrefix(k, "section.title.") && strings.ToLower(v) == t {
				return k, l
			}
		}
	}
	return "", ""
}

// Keys devolve todas as chaves de um idioma. Serve ao gate que confronta os catálogos
// entre si — uma chave que existe em `en` e falta em `es` é tradução pendente, e sem isso
// ela sairia em inglês sem ninguém notar.
func Keys(lang string) []string {
	c := load(lang)
	out := make([]string, 0, len(c))
	for k := range c {
		out = append(out, k)
	}
	return out
}
