package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

// route-exists: a rota que a spec DECLARA tem de existir no app.
//
// O `route-declared` confronta a spec contra si mesma — ela declara uma rota? Este
// confronta a spec contra o CÓDIGO: essa rota existe onde o app registra suas rotas?
//
// A diferença apareceu num E2E real. Uma spec de tela nova declarou `> **Rota**:
// `MetadataEdit“, e outra spec passou a prometer navegação para lá. `route-declared`
// (BLOQUEANTE) deu ✓ nas duas. A rota não existia em lugar nenhum do app: as duas specs
// descreviam um caminho para uma tela inalcançável, e o pipeline inteiro ficou verde.
//
// É a âncora que mente na variante mais difícil de ver: não falta nada, e tudo se
// referencia. A spec cita a rota, o gate confere que citou, e ninguém pergunta se ela
// existe.
//
// Medido em 96 specs de tela de um projeto real antes de ligar: 2 achados, ambos
// defeitos verdadeiros (a rota do E2E e uma spec cuja tela o app registra sob outro
// nome). Falso-positivo zero.
func checkRouteExists(content string, n mapx.Node, root string, g *mapx.Graph, cfg *config.Config) (Verdict, string) {
	if n.Kind != mapx.KindSpec {
		return Skip, "a rota é DECLARADA pela spec — é ela que promete o caminho"
	}
	rota := rotaDeclarada(content)
	if rota == "" {
		// Sem rota declarada não há o que confrontar. Cobrar a declaração é trabalho do
		// `route-declared`; duplicá-lo aqui produziria dois gates acusando o mesmo.
		return Skip, "a spec não declara rota"
	}
	globs := cfg.RouteRegistry()
	if len(globs) == 0 {
		// Sem saber ONDE o projeto registra rotas, o gate não pode confrontar. Nomear o
		// que falta é melhor que passar em silêncio: um ✓ aqui afirmaria que a rota
		// existe, e o gate não olhou.
		return Pending, "o projeto não declara `route_registry:` no anchors.yaml — sem " +
			"saber onde as rotas são registradas, não há como confrontar `" + rota + "`"
	}
	conhecidas, err := rotasRegistradas(root, globs)
	if err != nil {
		return Pending, "não foi possível ler o registro de rotas: " + err.Error()
	}
	if len(conhecidas) == 0 {
		return Pending, "nenhuma rota encontrada em `route_registry` — confira os globs " +
			"antes de concluir que a rota não existe"
	}
	if conhecidas[rota] {
		return Pass, ""
	}
	return Fail, fmt.Sprintf("a spec declara a rota `%s`, que NÃO existe no app "+
		"(%d rota(s) registrada(s) em %s). Uma rota declarada e inexistente não falha em "+
		"lugar nenhum: a spec fica bem-formada, quem a lê acredita que a tela é alcançável, "+
		"e outras specs passam a prometer navegação para lá. Registre a rota, ou corrija o "+
		"nome na spec para o que o app de fato usa",
		rota, len(conhecidas), strings.Join(globs, ", "))
}

// rotaDeclaradaRE casa a rota declarada pela spec, nas duas formas que aparecem:
//   - NOME de tela: `> **Rota**: `MetadataEdit“
//   - CAMINHO HTTP: `route: POST /manage-metadata` (interface de backend)
//
// A segunda existia e o gate não a via: uma spec de handler declarava `POST
// /manage-metadata` em duas linhas do cabeçalho, e os dois gates de rota respondiam `~`.
// Pior que o silêncio: o `~` ENSINA que não havia material a confrontar, quando havia uma
// promessa não cumprida — a Lambda não tinha rota, env var nem grant na infra, e a spec
// prometia a rota. O verbo HTTP é opcional e as crases também, porque as duas escritas
// aparecem no mesmo projeto.
var rotaDeclaradaRE = regexp.MustCompile("(?mi)^>?\\s*\\*{0,2}(?:rota|route)\\*{0,2}\\s*:\\s*`?(?:(?:GET|POST|PUT|PATCH|DELETE)\\s+)?(/[a-z0-9][a-z0-9/_-]*|[A-Za-z][A-Za-z0-9_]*)`?")

func rotaDeclarada(content string) string {
	if m := rotaDeclaradaRE.FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

// nomeDeRotaRE reconhece as formas em que uma rota é registrada:
//   - `name="Perfil"` — a prop do navegador (React Navigation, Expo Router e afins);
//   - `Perfil: undefined` / `Perfil: {` — a entrada no tipo do stack;
//   - `addResource('signup')` — a rota HTTP de um backend (API Gateway e afins).
//
// Todas contam porque projetos reais usam todas, e olhar só a primeira produz falso
// positivo: medido, 9 das 96 telas de um projeto declaravam rota que só aparecia no tipo,
// e 59 rotas HTTP viviam apenas na terceira forma.
var nomeDeRotaRE = regexp.MustCompile(`name="([A-Za-z][A-Za-z0-9_]*)"|(?m)^\s{2,}([A-Za-z][A-Za-z0-9_]*)\s*:\s*(?:undefined|\{)|addResource\('([a-z0-9][a-z0-9/_-]*)'`)

// rotasRegistradas lê os arquivos de registro de rota do projeto e devolve os nomes.
func rotasRegistradas(root string, globs []string) (map[string]bool, error) {
	out := map[string]bool{}
	fsys := os.DirFS(root)
	for _, glob := range globs {
		arquivos, err := doublestar.Glob(fsys, glob)
		if err != nil {
			return nil, fmt.Errorf("glob %q inválido: %w", glob, err)
		}
		for _, f := range arquivos {
			b, rerr := os.ReadFile(filepath.Join(root, f))
			if rerr != nil {
				continue
			}
			for _, m := range nomeDeRotaRE.FindAllStringSubmatch(string(b), -1) {
				if m[1] != "" {
					out[m[1]] = true
				}
				if m[2] != "" {
					out[m[2]] = true
				}
				// Rota HTTP: `addResource('signup')` registra `/signup`. Guardamos as duas
				// escritas porque a spec pode declarar com ou sem a barra inicial.
				if m[3] != "" {
					nu := strings.TrimPrefix(m[3], "/")
					out["/"+nu], out[nu] = true, true
				}
			}
		}
	}
	return out, nil
}
