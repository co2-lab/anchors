package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/i18n"
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
		return Skip, i18n.T("gate.route_exists.skip_not_spec")
	}
	rota := declaredRoute(content)
	if rota == "" {
		// Sem rota declarada não há o que confrontar. Cobrar a declaração é trabalho do
		// `route-declared`; duplicá-lo aqui produziria dois gates acusando o mesmo.
		return Skip, i18n.T("gate.route_exists.skip_no_route")
	}
	globs := cfg.RouteRegistry()
	if len(globs) == 0 {
		// Sem saber ONDE o projeto registra rotas, o gate não pode confrontar. Nomear o
		// que falta é melhor que passar em silêncio: um ✓ aqui afirmaria que a rota
		// existe, e o gate não olhou.
		return Pending, i18n.T("gate.route_exists.pending_no_route_registry", rota)
	}
	conhecidas, err := registeredRoutes(root, globs, cfg)
	if err != nil {
		return Pending, i18n.T("gate.route_exists.pending_read_err", err.Error())
	}
	if len(conhecidas) == 0 {
		return Pending, i18n.T("gate.route_exists.pending_no_routes_found")
	}
	if conhecidas[rota] {
		return Pass, ""
	}
	return Fail, i18n.T("gate.route_exists.fail_route_missing",
		rota, len(conhecidas), strings.Join(globs, ", "))
}

// declaredRouteRE casa a rota declarada pela spec, nas duas formas que aparecem:
//   - NOME de tela: `> **Rota**: `MetadataEdit“
//   - CAMINHO HTTP: `route: POST /manage-metadata` (interface de backend)
//
// A segunda existia e o gate não a via: uma spec de handler declarava `POST
// /manage-metadata` em duas linhas do cabeçalho, e os dois gates de rota respondiam `~`.
// Pior que o silêncio: o `~` ENSINA que não havia material a confrontar, quando havia uma
// promessa não cumprida — a Lambda não tinha rota, env var nem grant na infra, e a spec
// prometia a rota. O verbo HTTP é opcional e as crases também, porque as duas escritas
// aparecem no mesmo projeto.
var declaredRouteRE = regexp.MustCompile("(?mi)^>?\\s*\\*{0,2}(?:rota|route)\\*{0,2}\\s*:\\s*`?(?:(?:GET|POST|PUT|PATCH|DELETE)\\s+)?(/[a-z0-9][a-z0-9/_-]*|[A-Za-z][A-Za-z0-9_]*)`?")

func declaredRoute(content string) string {
	if m := declaredRouteRE.FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

// routeNameRE reconhece as formas em que uma rota é registrada (default de fallback):
//   - `name="Perfil"` — a prop do navegador (React Navigation, Expo Router e afins);
//   - `Perfil: undefined` / `Perfil: {` — a entrada no tipo do stack;
//   - `addResource('signup')` — a rota HTTP de um backend (API Gateway e afins).
//
// Projetos em outras stacks (Go, Python, Rails, etc.) declaram `derived.route_pattern`.
var routeNameRE = regexp.MustCompile(`name="([A-Za-z][A-Za-z0-9_]*)"|(?m)^\s{2,}([A-Za-z][A-Za-z0-9_]*)\s*:\s*(?:undefined|\{)|addResource\('([a-z0-9][a-z0-9/_-]*)'`)

// registeredRoutes lê os arquivos de registro de rota do projeto e devolve os nomes.
func registeredRoutes(root string, globs []string, cfg *config.Config) (map[string]bool, error) {
	out := map[string]bool{}
	fsys := os.DirFS(root)

	reRoute := routeNameRE
	if cfg != nil && cfg.Derived != nil && strings.TrimSpace(cfg.Derived.RoutePattern) != "" {
		re, err := regexp.Compile(cfg.Derived.RoutePattern)
		if err != nil || re.NumSubexp() < 1 {
			return nil, fmt.Errorf("invalid derived.route_pattern regex (must compile with at least 1 capture group)")
		}
		reRoute = re
	}

	for _, glob := range globs {
		arquivos, err := doublestar.Glob(fsys, glob, doublestar.WithFilesOnly())
		if err != nil {
			return nil, fmt.Errorf("%s", i18n.T("gate.route_exists.err_invalid_glob", glob, err))
		}
		for _, f := range arquivos {
			b, rerr := os.ReadFile(filepath.Join(root, f))
			if rerr != nil {
				// An unread registry file may be the one that registers the route: skipping
				// it accused a route that exists (RTEXR-E03).
				return nil, fmt.Errorf("%s", i18n.T("gate.route_exists.err_unreadable", f))
			}
			for _, m := range reRoute.FindAllStringSubmatch(string(b), -1) {
				for _, sub := range m[1:] {
					if sub != "" {
						nu := strings.TrimPrefix(sub, "/")
						out[sub] = true
						out[nu] = true
						out["/"+nu] = true
					}
				}
			}
		}
	}
	return out, nil
}
