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
	var names []string
	for _, k := range containersFixture(t).InternalContainers() {
		names = append(names, k.Name)
	}
	if want := []string{"app", "api"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("InternalContainers = %v, want %v (the external db has no level 3)", names, want)
	}
}

func TestContainerOfLayer(t *testing.T) {
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
	c := containersFixture(t)
	got := c.OrphanLayers([]string{"screen", "lambda", "service", "infra"})
	if want := []string{"lambda", "infra"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("OrphanLayers = %v, want %v", got, want)
	}
	if got := (&Config{}).OrphanLayers([]string{"a"}); !reflect.DeepEqual(got, []string{"a"}) {
		t.Fatalf("with no containers every layer is an orphan, got %v", got)
	}
}
