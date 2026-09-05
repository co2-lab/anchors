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
var arquivos embed.FS

// Suportados são os idiomas que o Anchors aceita em `lang:`.
//
// A lista é FECHADA de propósito. Aceitar qualquer código produziria projetos declarando
// idiomas que não existem no catálogo, e o erro apareceria só quando alguém lesse a
// primeira mensagem — em inglês, sem entender por quê. Fechar a lista faz o `anchors init`
// e o `check` recusarem na hora, dizendo quais existem.
var Suportados = []string{"pt-BR", "en", "es"}

// Padrao é o idioma quando o projeto não declara `lang:`.
//
// Inglês, e não português: um projeto que não declarou nada é provavelmente um projeto
// novo de alguém que encontrou o Anchors — e o inglês é o que mais gente lê. Quem quer
// português declara.
const Padrao = "en"

var (
	mu       sync.RWMutex
	atual    = Padrao
	catalogo = map[string]map[string]string{}
)

// Suportado diz se o idioma está na lista fechada.
func Suportado(lang string) bool {
	for _, l := range Suportados {
		if l == lang {
			return true
		}
	}
	return false
}

// Definir troca o idioma corrente. Um idioma fora da lista é recusado com a lista junto —
// quem errou o código precisa saber quais existem, não só que errou.
func Definir(lang string) error {
	if lang == "" {
		lang = Padrao
	}
	if !Suportado(lang) {
		return fmt.Errorf("idioma %q não é suportado — os disponíveis são: %s",
			lang, strings.Join(Suportados, ", "))
	}
	mu.Lock()
	defer mu.Unlock()
	atual = lang
	return nil
}

// Atual devolve o idioma corrente.
func Atual() string {
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

	if s, ok := busca(lang, chave); ok {
		return format(s, args...)
	}
	if lang != Padrao {
		if s, ok := busca(Padrao, chave); ok {
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

func busca(lang, chave string) (string, bool) {
	mu.RLock()
	c, carregado := catalogo[lang]
	mu.RUnlock()
	if !carregado {
		c = carrega(lang)
	}
	s, ok := c[chave]
	return s, ok
}

// carrega lê o JSON do idioma UMA vez e o guarda.
//
// Um catálogo ausente vira mapa vazio, não erro: o binário tem de rodar mesmo que um
// arquivo de tradução falte, caindo para o fallback. Falhar aqui derrubaria o CLI inteiro
// por causa de uma tradução.
func carrega(lang string) map[string]string {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := catalogo[lang]; ok {
		return c
	}
	c := map[string]string{}
	if b, err := arquivos.ReadFile("locales/" + lang + ".json"); err == nil {
		_ = json.Unmarshal(b, &c)
	}
	catalogo[lang] = c
	return c
}

// Chaves devolve todas as chaves de um idioma. Serve ao gate que confronta os catálogos
// entre si — uma chave que existe em `en` e falta em `es` é tradução pendente, e sem isso
// ela sairia em inglês sem ninguém notar.
func Chaves(lang string) []string {
	c := carrega(lang)
	out := make([]string, 0, len(c))
	for k := range c {
		out = append(out, k)
	}
	return out
}
