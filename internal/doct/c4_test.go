package doct

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

func projectWithContainers(t *testing.T) *Compiler {
	t.Helper()
	specs := map[string]string{
		"app/Tela.spec.md":   "---\ncode: TELAX\nlayer: screen\n---\n\n# Tela — a tela\n",
		"api/Rota.spec.md":   "---\ncode: ROTAX\nlayer: lambdas\n---\n\n# Rota — a rota\n",
		"api/Util.spec.md":   "---\ncode: UTILX\nlayer: shared\n---\n\n# Util — o util\n",
		"solto/Orfa.spec.md": "---\ncode: ORFAX\nlayer: orfa\n---\n\n# Orfa — sem contêiner\n",
	}
	root, g := projetoDeTeste(t, specs)
	c, err := New(root, g)
	if err != nil {
		t.Fatal(err)
	}
	c.Config = &config.Config{ContainersDecl: []config.Container{
		{Name: "app", Description: "a interface", Layers: []string{"screen"},
			Talks: []config.Talk{{To: "api", Protocol: "HTTPS/JSON"}}},
		{Name: "api", Description: "as rotas", Layers: []string{"lambdas", "shared"},
			Talks: []config.Talk{{To: "banco", Protocol: "SQL/TLS"}}},
		// The DATABASE is a container — a container is what runs OR STORES data, not
		// "a process we wrote". External because it is a third party's.
		{Name: "banco", Description: "o que persiste", External: true},
	}}
	return c
}

// LEVEL 3 IS PER CONTAINER, and the external one has none: we have no components inside it,
// and drawing them would claim a knowledge we do not have.
func TestInternalContainers_externalGetsNoLevel3(t *testing.T) {
	t.Run("C4CNC-B02: An external container gets no level 3", func(t *testing.T) {})
	t.Run("C4CNC-X01: The containers are only the declared ones", func(t *testing.T) {})
	c := projectWithContainers(t)

	all := c.fnContainers()
	var names []string
	for _, k := range all {
		names = append(names, k.Name)
	}
	if strings.Join(names, ",") != "app,api,banco" {
		t.Fatalf("containers = %v, want exactly the declared app, api, banco (the database IS a container)", names)
	}
	internal := c.fnInternalContainers()
	if len(internal) != 2 {
		t.Fatalf("internal = %d, want 2 — the external one has no level 3", len(internal))
	}
	for _, k := range internal {
		if k.Name == "banco" {
			t.Error("the database got a level 3 — we have no components inside it")
		}
	}
}

// A container's UNITS are those of the layers it declares — the bridge between the two
// vocabularies: a layer groups code, a component is a piece inside a container.
func TestContainers_unitsComeFromTheDeclaredLayers(t *testing.T) {
	c := projectWithContainers(t)
	by := map[string][]string{}
	for _, k := range c.fnContainers() {
		for _, u := range k.Units {
			by[k.Name] = append(by[k.Name], u.Code)
		}
	}
	if len(by["app"]) != 1 || by["app"][0] != "TELAX" {
		t.Errorf("app = %v, want [TELAX]", by["app"])
	}
	if len(by["api"]) != 2 {
		t.Errorf("api = %v, want ROTAX and UTILX (lambdas + shared)", by["api"])
	}
	if len(by["banco"]) != 0 {
		t.Errorf("banco = %v — it has no unit of ours", by["banco"])
	}
}

func TestContainers_layerMatchAndOrder(t *testing.T) {
	t.Run("C4CNC-B01: A container carries the specs of the layers it declares, in layer then code order", func(t *testing.T) {
		root, g := projetoDeTeste(t, map[string]string{
			"a/R.spec.md": "---\ncode: ROTAX\nlayer: lambdas\n---\n\n# R — r\n",
			"a/Y.spec.md": "---\ncode: UTILY\nlayer: shared\n---\n\n# Y — y\n",
			"a/X.spec.md": "---\ncode: UTILX\nlayer: shared\n---\n\n# X — x\n",
		})
		c, err := New(root, g)
		if err != nil {
			t.Fatal(err)
		}
		c.Config = &config.Config{ContainersDecl: []config.Container{
			{Name: "app", Layers: []string{" SHARED ", "lambdas"}},
		}}
		var got []string
		for _, u := range c.fnContainers()[0].Units {
			got = append(got, u.Code)
		}
		if strings.Join(got, ",") != "ROTAX,UTILX,UTILY" {
			t.Errorf("units = %v, want ROTAX, UTILX, UTILY", got)
		}
	})
}

// A layer OUTSIDE every container is named, not hidden.
func TestOrphanLayers_areNamed(t *testing.T) {
	t.Run("C4CNC-B03: A layer no container declares is named as orphan", func(t *testing.T) {})
	if orphans := projectWithContainers(t).fnOrphanLayers(); len(orphans) != 1 || orphans[0] != "orfa" {
		t.Errorf("orphans = %v, want [orfa]", orphans)
	}
}

func TestContainers_noConfig(t *testing.T) {
	t.Run("C4CNC-B04: No configuration means no containers and no orphans", func(t *testing.T) {
		c := projectWithContainers(t)
		c.Config = nil
		if k := c.fnContainers(); k != nil {
			t.Errorf("containers = %v, want none", k)
		}
		if k := c.fnInternalContainers(); k != nil {
			t.Errorf("internal containers = %v, want none", k)
		}
		if o := c.fnOrphanLayers(); o != nil {
			t.Errorf("orphans = %v, want none", o)
		}
	})
}

func TestContainers_everyLayerIsHeldOrOrphan(t *testing.T) {
	t.Run("C4CNC-I01: Every spec layer is held by a container or named orphan", func(t *testing.T) {
		c := projectWithContainers(t)
		seen := map[string]bool{}
		for _, k := range c.fnContainers() {
			for _, u := range k.Units {
				seen[u.Layer] = true
			}
		}
		for _, l := range c.fnOrphanLayers() {
			seen[l] = true
		}
		for _, l := range c.fnLayers() {
			if !seen[l] {
				t.Errorf("layer %q is neither in a container nor named orphan", l)
			}
		}
	})
}
