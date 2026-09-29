package config

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func containersFixture(t *testing.T) *Config {
	t.Helper()
	var c Config
	if err := yaml.Unmarshal([]byte(`
containers:
  - name: app
    description: the mobile app
    layers: [screen, hook]
    talks:
      - {to: api, protocol: HTTPS/JSON, why: queries}
  - name: api
    layers: [" Handler ", service]
  - name: db
    external: true
    layers: [table]
`), &c); err != nil {
		t.Fatal(err)
	}
	return &c
}

func TestContainers(t *testing.T) {
	t.Run("CNTNR-B01: The declared containers come back as written, and a missing config has none", func(t *testing.T) {})
	var nilCfg *Config
	if got := nilCfg.Containers(); got != nil {
		t.Fatalf("a nil config has no containers, got %v", got)
	}
	c := containersFixture(t)
	got := c.Containers()
	if len(got) != 3 || got[0].Name != "app" || !got[2].External {
		t.Fatalf("Containers = %+v", got)
	}
	if want := []Talk{{To: "api", Protocol: "HTTPS/JSON", Why: "queries"}}; !reflect.DeepEqual(got[0].Talks, want) {
		t.Fatalf("app talks = %+v, want %+v", got[0].Talks, want)
	}
}

func TestInternalContainers(t *testing.T) {
	t.Run("CNTNR-B02: The internal containers are the declared ones without the external, in declared order", func(t *testing.T) {})
	var names []string
	for _, k := range containersFixture(t).InternalContainers() {
		names = append(names, k.Name)
	}
	if want := []string{"app", "api"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("InternalContainers = %v, want %v (the external db has no level 3)", names, want)
	}
}

func TestContainerOfLayer(t *testing.T) {
	t.Run("CNTNR-B03: A layer is found in its container ignoring case and surrounding spaces", func(t *testing.T) {})
	t.Run("CNTNR-B04: A layer no container claims has no container, and that is an answer, not an error", func(t *testing.T) {})
	t.Run("CNTNR-B05: The layers of an external container are still claimed by it", func(t *testing.T) {})
	c := containersFixture(t)
	for layer, want := range map[string]string{
		"screen":  "app",
		"handler": "api", // declared as " Handler ": trimmed, case-insensitive
		"table":   "db",
		"lambda":  "",
	} {
		if got := c.ContainerOfLayer(layer); got != want {
			t.Errorf("ContainerOfLayer(%q) = %q, want %q", layer, got, want)
		}
	}
}

func TestOrphanLayers(t *testing.T) {
	t.Run("CNTNR-B06: The orphan layers are the given ones no container claims, in the given order", func(t *testing.T) {})
	t.Run("CNTNR-X01: With no container declared, no layer is placed by guessing: every layer is an orphan", func(t *testing.T) {})
	c := containersFixture(t)
	got := c.OrphanLayers([]string{"screen", "lambda", "service", "infra"})
	if want := []string{"lambda", "infra"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("OrphanLayers = %v, want %v", got, want)
	}
	if got := (&Config{}).OrphanLayers([]string{"a"}); !reflect.DeepEqual(got, []string{"a"}) {
		t.Fatalf("with no containers every layer is an orphan, got %v", got)
	}
}

func TestOrphanLayers_agreeWithContainerOfLayer(t *testing.T) {
	t.Run("CNTNR-I01: A layer is an orphan exactly when it has no container", func(t *testing.T) {})
	c := containersFixture(t)
	layers := []string{"screen", "HANDLER", "table", "lambda", "service", "infra"}
	orphan := map[string]bool{}
	for _, l := range c.OrphanLayers(layers) {
		orphan[l] = true
	}
	for _, l := range layers {
		if has := c.ContainerOfLayer(l) != ""; has == orphan[l] {
			t.Errorf("layer %q: in container=%v and orphan=%v disagree", l, has, orphan[l])
		}
	}
}

func TestOrphanLayers_onlyWhatRunsCode(t *testing.T) {
	t.Run("CNTNR-B07: A layer that runs no code is never an orphan", func(t *testing.T) {})
	c := &Config{Layers: map[string]Layer{
		"spec": {Kind: "spec"}, "e2e": {Kind: "test"}, "Landing-Component": {Kind: "code"}, "misc": {},
	}}
	got := c.OrphanLayers([]string{"spec", "E2E", "landing-component", "misc", "unknown"})
	if want := []string{"landing-component", "misc", "unknown"}; !reflect.DeepEqual(got, want) {
		t.Errorf("OrphanLayers = %v, want %v", got, want)
	}
	var none *Config
	if _, ok := none.layerDecl("x"); ok {
		t.Error("no configuration declares no layer")
	}
}
