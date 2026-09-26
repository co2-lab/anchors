package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/mapx"
)

func rootComRotas(t *testing.T) *config.Config {
	t.Helper()
	dir := t.TempDir()
	nav := filepath.Join(dir, "nav")
	if err := os.MkdirAll(nav, 0o755); err != nil {
		t.Fatal(err)
	}
	src := `<Stack.Screen name="Perfil" component={P} />
type Params = {
  Ajustes: undefined
  Detalhe: { id: string }
}
addResource('signup')
addResource('manage-metadata')
`
	if err := os.WriteFile(filepath.Join(nav, "Root.tsx"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANCHORS_TEST_ROOT", dir)
	return &config.Config{RouteRegistryGlobs: []string{"nav/**/*.tsx"}}
}

func TestRouteExists_B01_NaoSpecPula(t *testing.T) {
	t.Run("RTEXR-B01: Non-spec artifacts skip confrontation", func(t *testing.T) {})
	t.Run("RTEXR-I01: Route presence is validated only for specifications", func(t *testing.T) {})
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	for _, k := range []mapx.Kind{mapx.KindCode, mapx.KindFeature, mapx.KindTest} {
		v, _ := checkRouteExists("> **Rota**: `Perfil`", mapx.Node{Kind: k}, root, nil, cfg)
		if v != Skip {
			t.Fatalf("esperava Skip para kind %v, veio %v", k, v)
		}
	}
}

func TestRouteExists_B02_SemRotaDeclaradaPula(t *testing.T) {
	t.Run("RTEXR-B02: Specifications without declared routes skip confrontation", func(t *testing.T) {})
	t.Run("RTEXR-X01: The gate does not require every specification to declare a route", func(t *testing.T) {})
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	v, _ := checkRouteExists("# Uma spec sem rota declarada", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v != Skip {
		t.Fatalf("sem rota declarada, Skip; veio %v", v)
	}
}

func TestRouteExists_B03_SemRegistroFicaPendente(t *testing.T) {
	t.Run("RTEXR-B03: Specifications with declared routes return Pending when route registry is unconfigured", func(t *testing.T) {})
	t.Run("RTEXR-I02: The gate never approves route existence without inspecting registry files", func(t *testing.T) {})
	v, msg := checkRouteExists("> **Rota**: `Qualquer`", mapx.Node{Kind: mapx.KindSpec},
		t.TempDir(), nil, &config.Config{})
	if v != Pending {
		t.Fatalf("sem route_registry o veredito é Pending; veio %v (%s)", v, msg)
	}
}

func TestRouteExists_B04_GlobInvalidoFicaPendente(t *testing.T) {
	t.Run("RTEXR-B04: Invalid route registry glob patterns return Pending", func(t *testing.T) {})
	cfg := &config.Config{RouteRegistryGlobs: []string{"[a-"}}
	v, msg := checkRouteExists("> **Rota**: `Perfil`", mapx.Node{Kind: mapx.KindSpec},
		t.TempDir(), nil, cfg)
	if v != Pending {
		t.Fatalf("glob inválido deve retornar Pending; veio %v (%s)", v, msg)
	}
}

func TestRouteExists_B05_ZeroRotasEncontradasFicaPendente(t *testing.T) {
	t.Run("RTEXR-B05: Route registry files containing zero registered routes return Pending", func(t *testing.T) {})
	dir := t.TempDir()
	nav := filepath.Join(dir, "nav")
	if err := os.MkdirAll(nav, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nav, "Empty.tsx"), []byte("// sem rotas aqui\nconst x = 1;"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{RouteRegistryGlobs: []string{"nav/**/*.tsx"}}
	v, msg := checkRouteExists("> **Rota**: `Perfil`", mapx.Node{Kind: mapx.KindSpec}, dir, nil, cfg)
	if v != Pending {
		t.Fatalf("sem rotas no arquivo de registro deve retornar Pending; veio %v (%s)", v, msg)
	}
}

func TestRouteExists_B06_ScreenNamePropPassa(t *testing.T) {
	t.Run("RTEXR-B06: Declared screen route matching a component navigation prop passes", func(t *testing.T) {})
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	v, msg := checkRouteExists("> **Rota**: `Perfil`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("rota Perfil via name=\"Perfil\" deve passar; veio %v (%s)", v, msg)
	}
}

func TestRouteExists_B07_StackParamsTypePassa(t *testing.T) {
	t.Run("RTEXR-B07: Declared screen route matching a navigation stack parameter type entry passes", func(t *testing.T) {})
	t.Run("RTEXR-X02: Route parameter schemas and payload contracts are not evaluated by this gate", func(t *testing.T) {})
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	for _, r := range []string{"Ajustes", "Detalhe"} {
		v, msg := checkRouteExists("> **Rota**: `"+r+"`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
		if v != Pass {
			t.Fatalf("rota de tipo de stack %q deve passar; veio %v (%s)", r, v, msg)
		}
	}
}

func TestRouteExists_B08_AddResourcePassa(t *testing.T) {
	t.Run("RTEXR-B08: Declared backend route matching an HTTP resource registration passes", func(t *testing.T) {})
	t.Run("RTEXR-X03: Route access permissions and authentication middlewares are outside evaluation scope", func(t *testing.T) {})
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	v, msg := checkRouteExists("> **Rota**: `signup`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("rota HTTP via addResource deve passar; veio %v (%s)", v, msg)
	}
}

func TestRouteExists_B09_VerboHTTPRemovidoPassa(t *testing.T) {
	t.Run("RTEXR-B09: HTTP method verb prefixes are stripped when evaluating declared routes", func(t *testing.T) {})
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	v, msg := checkRouteExists("route: POST /manage-metadata", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v != Pass {
		t.Fatalf("rota com verbo POST deve ser reconhecida e passar; veio %v (%s)", v, msg)
	}
}

func TestRouteExists_B10_NormalizacaoDeBarraPassa(t *testing.T) {
	t.Run("RTEXR-B10: Route paths match regardless of leading slash differences between specification and code", func(t *testing.T) {})
	t.Run("RTEXR-I03: Route matching is slash-normalized between specification and code", func(t *testing.T) {})
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	// signup está registrado sem barra (addResource('signup'))
	// com barra na spec:
	v1, msg1 := checkRouteExists("> **Rota**: `/signup`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v1 != Pass {
		t.Fatalf("rota com barra /signup deve passar para addResource('signup'); veio %v (%s)", v1, msg1)
	}
	// sem barra na spec:
	v2, msg2 := checkRouteExists("> **Rota**: `signup`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v2 != Pass {
		t.Fatalf("rota sem barra signup deve passar para addResource('signup'); veio %v (%s)", v2, msg2)
	}
}

func TestRouteExists_B11_PadraoCustomizadoPassa(t *testing.T) {
	t.Run("RTEXR-B11: Custom route pattern regex matches custom route registration patterns", func(t *testing.T) {})
	dir := t.TempDir()
	routesFile := filepath.Join(dir, "routes.go")
	src := `package main
func registerRoutes(r chi.Router) {
    r.Get("/usuarios", ListUsers)
    r.Post("/usuarios", CreateUser)
}`
	if err := os.WriteFile(routesFile, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{
		RouteRegistryGlobs: []string{"routes.go"},
		Derived: &config.Derived{
			RoutePattern: `r\.(?:Get|Post)\("([^"]+)"`,
		},
	}
	v, msg := checkRouteExists("> **Rota**: `/usuarios`", mapx.Node{Kind: mapx.KindSpec}, dir, nil, cfg)
	if v != Pass {
		t.Fatalf("rota customizada deve passar; veio %v (%s)", v, msg)
	}
}

func TestRouteExists_B12_PadraoCustomizadoInvalidoPendente(t *testing.T) {
	t.Run("RTEXR-B12: Custom route pattern regex that is malformed or lacks capture groups returns Pending", func(t *testing.T) {})
	dir := t.TempDir()
	routesFile := filepath.Join(dir, "routes.go")
	if err := os.WriteFile(routesFile, []byte(`r.Get("/usuarios")`), 0o644); err != nil {
		t.Fatal(err)
	}
	// Malformed regex
	cfgInvalid := &config.Config{
		RouteRegistryGlobs: []string{"routes.go"},
		Derived: &config.Derived{
			RoutePattern: `[a-`,
		},
	}
	v1, _ := checkRouteExists("> **Rota**: `/usuarios`", mapx.Node{Kind: mapx.KindSpec}, dir, nil, cfgInvalid)
	if v1 != Pending {
		t.Fatalf("regex customizado malformado deve retornar Pending; veio %v", v1)
	}

	// Regex without capture groups
	cfgNoGroup := &config.Config{
		RouteRegistryGlobs: []string{"routes.go"},
		Derived: &config.Derived{
			RoutePattern: `r\.Get`,
		},
	}
	v2, msg2 := checkRouteExists("> **Rota**: `/usuarios`", mapx.Node{Kind: mapx.KindSpec}, dir, nil, cfgNoGroup)
	if v2 != Pending {
		t.Fatalf("regex customizado sem grupo de captura deve retornar Pending; veio %v", v2)
	}
	if !strings.Contains(msg2, "capture group") {
		t.Fatalf("esperava erro citando capture group, veio: %s", msg2)
	}
}

func TestRouteExists_B13_RotaInexistenteFalha(t *testing.T) {
	t.Run("RTEXR-B13: Declared routes missing from all registered route definitions fail", func(t *testing.T) {})
	t.Run("RTEXR-I04: Missing routes produce a blocking Fail verdict", func(t *testing.T) {})
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	v, msg := checkRouteExists("> **Rota**: `NaoExiste`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v != Fail {
		t.Fatalf("rota inexistente deve reprovar; veio %v (%s)", v, msg)
	}
	if !strings.Contains(msg, "NaoExiste") {
		t.Errorf("mensagem deve conter o nome da rota faltante; veio %s", msg)
	}
	if !strings.Contains(msg, "nav/**/*.tsx") {
		t.Errorf("mensagem deve conter o glob pesquisado; veio %s", msg)
	}
}

// An unread registry file may register the route: it never turns into an accusation.
func TestRouteExists_unreadableRegistryIsPending(t *testing.T) {
	t.Run("RTEXR-E03: An unreadable registry file leaves the route pending", func(t *testing.T) {})
	if os.Geteuid() == 0 {
		t.Skip("root reads a chmod 000 file")
	}
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")
	locked := filepath.Join(root, "nav", "Other.tsx")
	if err := os.WriteFile(locked, []byte(`<Stack.Screen name="Carteira" />`), 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(locked, 0o644)
	v, msg := checkRouteExists("> **Rota**: `Carteira`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v != Pending || !strings.Contains(msg, "Other.tsx") {
		t.Fatalf("an unreadable registry file must leave the route Pending naming it, got %v (%s)", v, msg)
	}
}
