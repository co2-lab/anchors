package gate

import (
	"os"
	"path/filepath"
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
	// as duas formas em que uma rota é registrada na prática
	src := `<Stack.Screen name="Perfil" component={P} />
type Params = {
  Ajustes: undefined
  Detalhe: { id: string }
}`
	if err := os.WriteFile(filepath.Join(nav, "Root.tsx"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANCHORS_TEST_ROOT", dir)
	return &config.Config{RouteRegistryGlobs: []string{"nav/**/*.tsx"}}
}

// TestRotaInexistenteEhAcusada guarda o defeito de um E2E real: uma spec declarou
// `> **Rota**: `MetadataEdit“, outra prometeu navegação para lá, e o gate BLOQUEANTE de
// rota deu ✓ nas duas — porque ele confere que a spec DECLAROU, não que a rota EXISTE.
// Duas specs descreviam caminho para uma tela inalcançável, tudo verde.
func TestRotaInexistenteEhAcusada(t *testing.T) {
	cfg := rootComRotas(t)
	root := os.Getenv("ANCHORS_TEST_ROOT")

	v, msg := checkRouteExists("> **Rota**: `NaoExiste`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg)
	if v != Fail {
		t.Fatalf("rota inexistente deve reprovar; veio %v (%s)", v, msg)
	}

	for _, rota := range []string{"Perfil", "Ajustes", "Detalhe"} {
		// `Perfil` vem de name="", os outros do tipo do stack — as duas formas contam.
		if v, msg := checkRouteExists("> **Rota**: `"+rota+"`", mapx.Node{Kind: mapx.KindSpec}, root, nil, cfg); v != Pass {
			t.Errorf("rota registrada %q deve passar; veio %v (%s)", rota, v, msg)
		}
	}
}

// Sem `route_registry` o gate NÃO pode afirmar que a rota existe — e também não deve
// acusar. O terceiro estado nomeia o que falta em vez de fingir aprovação.
func TestSemRegistroDeRotasFicaPendente(t *testing.T) {
	v, msg := checkRouteExists("> **Rota**: `Qualquer`", mapx.Node{Kind: mapx.KindSpec},
		t.TempDir(), nil, &config.Config{})
	if v != Pending {
		t.Errorf("sem `route_registry` o veredito é Pending; veio %v (%s)", v, msg)
	}
}

// Spec que não declara rota não é assunto deste gate (é do `route-declared`).
func TestSemRotaDeclaradaNaoEhAssunto(t *testing.T) {
	cfg := rootComRotas(t)
	if v, _ := checkRouteExists("# Uma spec qualquer", mapx.Node{Kind: mapx.KindSpec},
		os.Getenv("ANCHORS_TEST_ROOT"), nil, cfg); v != Skip {
		t.Errorf("sem rota declarada, Skip; veio %v", v)
	}
}
