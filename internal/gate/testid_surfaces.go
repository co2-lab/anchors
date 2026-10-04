// @anchors
//   code: TSGTS
//   ref: TICTS

// Superfícies CONSUMIDORAS do testID: onde procurar quem se apoia no handle — o teste
// ligado pela unidade, o teste vizinho (compartilhado ou do pai) e os flows de ponta a
// ponta, que vivem fora do grafo.
//
// Este arquivo já foi `testid_honored.go`. O gate foi aposentado por `testid-coerente`
// — que também cobre a ponta que o antigo declarava fora de escopo (id consultado e
// inexistente), depois de a premissa "isso já falha ao rodar, ruidosamente" ter sido
// refutada por medição.
package gate

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/config"
)

// queriedID — o id aparece em alguma superfície consumidora? Para o template
// (`bdgt-item-*`), basta o PREFIXO ser consultado: o sufixo é dado de runtime, e o
// flow que casa `bdgt-item-.*` ou constrói `bdgt-item-${id}` está usando o contrato.
//
// A menção vale só numa FRONTEIRA de id: antes dela nada que continue um id, e (fora
// do curinga) nada depois. Era substring nua, e um teste que só consultava
// `btn-save` contava como consulta de `btn` — um handle que ninguém usa passava por
// consumido. A marca `:` não é caractere de id, então `:btn` casa pela mesma régua
// (a superfície pode ter sido escrita antes da convenção de marcação).
func queriedID(blob, id string) bool {
	nu := strings.TrimPrefix(id, ":")
	if strings.HasSuffix(nu, "-*") {
		return mentionsAtBoundary(blob, strings.TrimSuffix(nu, "*"), false)
	}
	return mentionsAtBoundary(blob, nu, true)
}

// mentionsAtBoundary procura `needle` em `blob` sem caractere de id colado antes e,
// quando `closed`, também depois.
func mentionsAtBoundary(blob, needle string, closed bool) bool {
	if needle == "" {
		return false
	}
	for from := 0; ; {
		i := strings.Index(blob[from:], needle)
		if i < 0 {
			return false
		}
		i += from
		end := i + len(needle)
		if (i == 0 || !isIDByte(blob[i-1])) && (!closed || end == len(blob) || !isIDByte(blob[end])) {
			return true
		}
		from = i + 1
	}
}

// isIDByte: os caracteres que um handle aceita (`[a-zA-Z0-9._-]`, ver handleRegex).
func isIDByte(c byte) bool {
	return c == '-' || c == '_' || c == '.' ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// readE2ESurface lê os flows de ponta a ponta declarados pelo projeto.
//
// A resolução é em DOIS passos, e confundi-los custa caro: `surfaces[e2e]` devolve a
// CHAVE da superfície (ex.: "e2e"), não um caminho — quem tem o caminho é
// `files[chave]`. Ler `surfaces` como se fosse path faz o gate procurar um diretório
// que não existe, achar zero consumidores e reportar como ÓRFÃO todo id usado apenas
// pelo E2E — acusando exatamente quem cumpre o contrato.
//
// Quando o projeto declara o regime mas não o arquivo (o app de referência hoje: `e2e: e2e` sem
// `files.e2e`), não há onde procurar. Devolver vazio aqui é honesto; o gate trata a
// falta de consumidor conhecido como Skip, não como reprovação.
func readE2ESurface(root string, cfg *config.Config) []string {
	if cfg == nil || cfg.Derived == nil {
		return nil
	}
	chave := cfg.Derived.Surfaces["e2e"]
	if chave == "" {
		return nil
	}
	// O PRIMEIRO padrão da camada: este gate resolve UM caminho de superfície, e uma
	// camada com vários padrões não muda o que ele pergunta.
	var padrao string
	if ps := cfg.Derived.PadroesDe()[chave]; len(ps) > 0 {
		padrao = ps[0]
	}
	if padrao == "" {
		// Sem template de arquivo para a superfície, procuramos por override — a
		// mesma precedência que o resto do framework usa.
		for _, ov := range cfg.Derived.Overrides {
			if ps := ov.PadroesDe()[chave]; len(ps) > 0 {
				padrao = ps[0]
				break
			}
		}
	}
	if padrao == "" {
		return nil
	}
	// A superfície é declarada como template de região (ex.: `apps/mobile/.maestro`);
	// aqui interessa a RAIZ dela, varrida por inteiro.
	dir := filepath.Join(root, firstStaticSegment(padrao))
	var out []string
	_ = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if b, e := os.ReadFile(p); e == nil {
			out = append(out, string(b))
		}
		return nil
	})
	return out
}

// readNeighborTests lê os testes da MESMA pasta e os da pasta-IRMÃ dentro do módulo.
//
// A pasta sozinha não basta, e o caso que mostra isso é o componente: `CategoryCard`
// vive em `trends/components/` e quem o exercita é `TrendsScreen.test.tsx`, em
// `trends/screens/` — a tela que o renderiza. Sem subir um nível, o gate acusa de
// órfão um handle que o teste consulta, e o achado aponta para o lugar errado.
//
// Sobe UM nível só (o módulo da feature), não a árvore inteira: ler tudo tornaria
// qualquer menção do repositório uma prova de consumo, e o gate deixaria de medir.
func readNeighborTests(root, specID string) []string {
	dir := filepath.Join(root, filepath.Dir(specID))
	dirs := []string{dir}
	// As irmãs dentro do módulo (`features/trends/{components,screens,hooks}`).
	if entradas, err := os.ReadDir(filepath.Dir(dir)); err == nil {
		for _, e := range entradas {
			if e.IsDir() {
				if d := filepath.Join(filepath.Dir(dir), e.Name()); d != dir {
					dirs = append(dirs, d)
				}
			}
		}
	}
	var out []string
	for _, d := range dirs {
		entradas, err := os.ReadDir(d)
		if err != nil {
			continue
		}
		for _, e := range entradas {
			if e.IsDir() || !strings.Contains(e.Name(), ".test.") {
				continue
			}
			if b, err := os.ReadFile(filepath.Join(d, e.Name())); err == nil {
				out = append(out, string(b))
			}
		}
	}
	return out
}

// firstStaticSegment corta o template no primeiro placeholder (`{{module}}`, `{unit}`)
// ou curinga de glob (`*`, `?`, `[`), devolvendo o DIRETÓRIO fixo antes dele — a raiz
// que dá para varrer.
//
// O glob não era cortado: `e2e/**/*.spec.ts` virava o caminho literal, o Walk não
// achava diretório nenhum, e a superfície inteira sumia em silêncio (todo id usado só
// pelos flows virava órfão). E o corte no MEIO de um segmento (`e2e/login-*.yaml`,
// `apps/x-{{m}}`) deixa um nome que não existe; a raiz é o diretório que o contém.
var placeholderRE = regexp.MustCompile(`\{\{?[a-zA-Z_]+\}?\}|[*?\[]`)

func firstStaticSegment(padrao string) string {
	if loc := placeholderRE.FindStringIndex(padrao); loc != nil {
		padrao = padrao[:loc[0]]
		if padrao != "" && !strings.HasSuffix(padrao, "/") {
			padrao = path.Dir(padrao)
		}
	}
	// `path`, not `filepath`: the pattern is written with `/` whatever the system, and
	// `filepath` turned `e2e/flows` into `e2e\flows` on Windows.
	return strings.TrimSuffix(path.Clean(padrao), "/")
}
