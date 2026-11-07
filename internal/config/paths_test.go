// @anchors
//   code: PTHTA
//   ref: PTHSA

package config

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func loadAliased(t *testing.T, body string) (*Config, error) {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "@scope"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 8\n"+body), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

func TestPaths_aliasesExpandInEveryPathField(t *testing.T) {
	t.Run("PTHSA-B01: A path that opens with an alias reads the path the alias stands for, in every field that holds a path", func(t *testing.T) {})
	c, err := loadAliased(t, `aliases:
  "@guides": docs/engineering/guides/
  "@app": apps/mobile/src
layers:
  ui:
    pattern: "@app/screens/**/*.tsx"
    kind: code
    exclude: ["@app/screens/**/*.test.tsx"]
  scoped:
    pattern: "@scope/pkg/*.ts"
    kind: code
governs:
  - from: "@guides/SCREEN_GUIDE.md"
    governs: tela
derived:
  anchor: spec
  files:
    test: ["@app/{{name}}.test.tsx"]
gates:
  - name: line-coverage
    no_signal:
      "@app/generated/**": generated
changelog:
  template: "@guides"
`)
	if err != nil {
		t.Fatal(err)
	}
	if got := c.Layers["ui"].Pattern; got != "apps/mobile/src/screens/**/*.tsx" {
		t.Errorf("a layer's glob: %s", got)
	}
	if got := c.Layers["ui"].Exclude[0]; got != "apps/mobile/src/screens/**/*.test.tsx" {
		t.Errorf("a layer's exclusion: %s", got)
	}
	if got := c.Layers["scoped"].Pattern; got != "@scope/pkg/*.ts" {
		t.Errorf("a folder named with @ that no alias names is kept: %s", got)
	}
	if got := c.Governs[0].From; got != "docs/engineering/guides/SCREEN_GUIDE.md" {
		t.Errorf("a guide, the alias's trailing slash dropped: %s", got)
	}
	if got := c.Derived.Files["test"][0]; got != "apps/mobile/src/{{name}}.test.tsx" {
		t.Errorf("a derived template: %s", got)
	}
	var gate Gate
	for _, g := range c.Gates {
		if g.Name == "line-coverage" {
			gate = g
		}
	}
	if gate.NoSignal["apps/mobile/src/generated/**"] != "generated" || len(gate.NoSignal) != 1 {
		t.Errorf("a glob that is a map's key: %v", gate.NoSignal)
	}
	if got := c.Changelog.Template; got != "docs/engineering/guides" {
		t.Errorf("the alias alone is its path: %s", got)
	}
}

func TestPaths_aBadAliasRefusesTheConfiguration(t *testing.T) {
	t.Run("PTHSA-E01: An alias badly named, empty, a URL or opening with another, or a path citing an undeclared alias, refuses the configuration", func(t *testing.T) {})
	for name, tc := range map[string]struct{ body, want string }{
		"no @":      {"aliases:\n  guides: docs\n", "`aliases.guides`"},
		"empty":     {"aliases:\n  \"@g\": \"\"\n", "names no path"},
		"a URL":     {"aliases:\n  \"@g\": https://raw.githubusercontent.com/o/r/main\n", "is a URL"},
		"chained":   {"aliases:\n  \"@a\": docs\n  \"@b\": \"@a/guides\"\n", "opens with another alias"},
		"undeclare": {"layers:\n  ui:\n    pattern: \"@ap/**\"\n    kind: code\n", "`layers.ui.pattern` cites the alias `@ap`"},
	} {
		_, err := loadAliased(t, tc.body)
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: want %q, got %v", name, tc.want, err)
		}
	}
}

// pathNamed are the yaml names a field holding a path goes by. A field so named carries a
// `path:` tag — or is in notAPath, with the reason —, so a new field that holds a path
// cannot be born outside the aliases.
var pathNamed = regexp.MustCompile(`^(path|paths|pattern|from|guide|template|files|file_patterns|exclude|support|keep|workflows|unguarded|inputs|reports|report|junit|lcov|must_appear_in|route_registry|packs|pack_values|no_signal|coverage_floors)$`)

var notAPath = map[string]string{
	"TestsSource.pattern": "a regular expression over a test file's lines",
	"TestLevel.exclude":   "regular expressions over rule codes",
	"Requirement.from":    "a branch pattern",
}

func TestPaths_everyFieldNamedForAPathIsTagged(t *testing.T) {
	t.Run("PTHSA-B01: A path that opens with an alias reads the path the alias stands for, in every field that holds a path", func(t *testing.T) {})
	seen := map[reflect.Type]bool{}
	var walk func(t reflect.Type)
	var missing []string
	walk = func(ty reflect.Type) {
		for ty.Kind() == reflect.Pointer || ty.Kind() == reflect.Slice || ty.Kind() == reflect.Map {
			ty = ty.Elem()
		}
		if ty.Kind() != reflect.Struct || seen[ty] {
			return
		}
		seen[ty] = true
		for i := 0; i < ty.NumField(); i++ {
			sf := ty.Field(i)
			name := yamlName(sf)
			if !sf.IsExported() || name == "" {
				continue
			}
			key := ty.Name() + "." + name
			if _, tagged := sf.Tag.Lookup("path"); !tagged && pathNamed.MatchString(name) && holdsText(sf.Type) {
				if _, ok := notAPath[key]; !ok {
					missing = append(missing, key)
				}
			}
			walk(sf.Type)
		}
	}
	walk(reflect.TypeOf(Config{}))
	if len(missing) > 0 {
		t.Errorf("fields named for a path with no `path:` tag (tag them, or list them in notAPath with the reason): %s", strings.Join(missing, ", "))
	}
}

// holdsText says whether a field's type holds text: a string, or a list or map of them.
func holdsText(ty reflect.Type) bool {
	for ty.Kind() == reflect.Pointer || ty.Kind() == reflect.Slice || ty.Kind() == reflect.Map {
		if ty.Kind() == reflect.Map && ty.Key().Kind() == reflect.String {
			return true
		}
		ty = ty.Elem()
	}
	return ty.Kind() == reflect.String
}
