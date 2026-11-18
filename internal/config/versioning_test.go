// @anchors
//   code: VRSNG
//   ref: VRSNN

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadVersioning(t *testing.T, body string) (*Config, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "anchors.yaml")
	if err := os.WriteFile(p, []byte("version: 8\nversioning:\n"+body), 0o644); err != nil {
		t.Fatal(err)
	}
	return Load(p)
}

const appProduct = `  app:
    scheme: semver
    source:
      - file: app/manifest.json
        path: version
    tag: "app-v{version}"
    next: { strategy: semantic, initial: true, types: { perf: patch } }
    branches:
      "hotfix/*":
        next: { strategy: increment, part: patch }
      "re:hotfix/.*":
        next: none
`

func TestVersioning_loads(t *testing.T) {
	t.Run("VRSNN-B01: The versioning block loads its products, and a configuration without it versions nothing", func(t *testing.T) {
		c, err := loadVersioning(t, appProduct)
		if err != nil {
			t.Fatal(err)
		}
		app := c.Versioning["app"]
		if app == nil || app.Scheme != "semver" || len(app.Source) != 1 || app.Source[0].File != "app/manifest.json" ||
			app.Tag != "app-v{version}" || app.Next.Strategy != "semantic" || !app.Next.Initial || app.Next.Types["perf"] != "patch" || len(app.Branches) != 2 {
			t.Fatalf("got %+v", app)
		}
		p := filepath.Join(t.TempDir(), "anchors.yaml")
		os.WriteFile(p, []byte("version: 8\n"), 0o644)
		if c, err := Load(p); err != nil || len(c.Versioning) != 0 {
			t.Errorf("no block: %v %v", c.Versioning, err)
		}
	})
	t.Run("VRSNN-B02: Next is a strategy's name or a map, and manual when none is written", func(t *testing.T) {
		c, err := loadVersioning(t, `  a: { scheme: number, next: none }
  b: { scheme: semver, next: { strategy: increment, part: patch } }
  c: { scheme: semver, source: [{ file: v.txt, regex: '(.+)' }] }
`)
		if err != nil {
			t.Fatal(err)
		}
		for name, want := range map[string]string{"a": "none", "b": "increment/patch", "c": "manual"} {
			st := c.Versioning[name].Next.Resolved()
			got := st.Name
			if st.Part != "" {
				got += "/" + st.Part
			}
			if got != want {
				t.Errorf("%s: %s, want %s", name, got, want)
			}
		}
	})
	t.Run("VRSNN-B03: A branch takes the first rule it matches, in the order written, else the default", func(t *testing.T) {
		c, err := loadVersioning(t, appProduct)
		if err != nil {
			t.Fatal(err)
		}
		app := c.Versioning["app"]
		if st, pat := app.NextFor("hotfix/login"); st.Name != "increment" || pat != "hotfix/*" {
			t.Errorf("hotfix/login: %+v %q", st, pat)
		}
		if st, pat := app.NextFor("develop"); st.Name != "semantic" || pat != "" {
			t.Errorf("develop: %+v %q", st, pat)
		}
	})
	t.Run("VRSNN-B04: The tag is v{version} by default, and a tag is read back only with the product's text around it", func(t *testing.T) {
		plain, app := &Product{}, &Product{Tag: "app-v{version}"}
		if plain.TagFor("1.2.0") != "v1.2.0" || app.TagFor("1.2.0") != "app-v1.2.0" {
			t.Errorf("tags: %s %s", plain.TagFor("1.2.0"), app.TagFor("1.2.0"))
		}
		if v, ok := app.VersionOfTag("app-v1.2.0"); !ok || v != "1.2.0" {
			t.Errorf("app-v1.2.0: %q %v", v, ok)
		}
		for _, tag := range []string{"api-v1.2.0", "app-v"} {
			if v, ok := app.VersionOfTag(tag); ok {
				t.Errorf("%s read back as %q", tag, v)
			}
		}
	})
}

func files(m map[string]string) func(string) ([]byte, error) {
	return func(f string) ([]byte, error) {
		s, ok := m[f]
		if !ok {
			return nil, os.ErrNotExist
		}
		return []byte(s), nil
	}
}

func TestVersioning_sources(t *testing.T) {
	disk := map[string]string{
		"app.json":  "{\n  \"name\": \"x\",\n  \"expo\": { \"version\": \"1.2.0\", \"build\": 7 }\n}\n",
		"list.yaml": "- version: '1.2.0'\n- version: 0.0.1\n",
		"build.txt": "NAME=x\nVERSION=1.3.0\nEND=1\n",
	}
	p := &Product{Source: []VersionSource{
		{File: "app.json", Path: "expo.version"}, {File: "list.yaml", Path: "0.version"}, {File: "build.txt", Regex: `VERSION=(\S+)`},
	}}
	t.Run("VRSNN-B05: A source is read at its path or in its group, and the sources must agree", func(t *testing.T) {
		svs, err := p.ReadSources(files(disk))
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, s := range svs {
			got = append(got, s.File+"="+s.Version)
		}
		if strings.Join(got, " ") != "app.json=1.2.0 list.yaml=1.2.0 build.txt=1.3.0" {
			t.Errorf("read %v", got)
		}
		if _, err := Agreed(svs); err == nil || !strings.Contains(err.Error(), "app.json 1.2.0") || !strings.Contains(err.Error(), "build.txt 1.3.0") {
			t.Errorf("disagreement: %v", err)
		}
		if v, err := Agreed(svs[:2]); err != nil || v != "1.2.0" {
			t.Errorf("agreement: %q %v", v, err)
		}
	})
	t.Run("VRSNN-B06: A source is rewritten in place, keeping every byte around the version", func(t *testing.T) {
		for file, want := range map[string]string{
			"app.json":  strings.Replace(disk["app.json"], "1.2.0", "1.3.0", 1),
			"list.yaml": strings.Replace(disk["list.yaml"], "1.2.0", "1.3.0", 1),
			"build.txt": strings.Replace(disk["build.txt"], "1.3.0", "1.4.0", 1),
		} {
			v := "1.3.0"
			if file == "build.txt" {
				v = "1.4.0"
			}
			got, err := p.Rewrite(file, []byte(disk[file]), v)
			if err != nil || string(got) != want {
				t.Errorf("%s: %q %v\nwant %q", file, got, err, want)
			}
		}
	})
	t.Run("VRSNN-E05: A version not found is an error naming the file", func(t *testing.T) {
		for _, c := range []struct {
			src  VersionSource
			disk map[string]string
		}{
			{VersionSource{File: "gone.json", Path: "version"}, map[string]string{}},
			{VersionSource{File: "app.json", Path: "expo.missing"}, disk},
			{VersionSource{File: "build.txt", Regex: `RELEASE=(\S+)`}, disk},
			{VersionSource{File: "bad.yaml", Path: "version"}, map[string]string{"bad.yaml": "a: [\n"}},
		} {
			_, err := (&Product{Source: []VersionSource{c.src}}).ReadSources(files(c.disk))
			if err == nil || !strings.Contains(err.Error(), c.src.File) {
				t.Errorf("%+v: %v", c.src, err)
			}
		}
	})
}

func TestVersioning_refused(t *testing.T) {
	refuses := func(t *testing.T, body string, says ...string) {
		t.Helper()
		_, err := loadVersioning(t, body)
		for _, s := range says {
			if err == nil || !strings.Contains(err.Error(), s) {
				t.Errorf("%s\n→ %v, want it to say %q", body, err, s)
			}
		}
	}
	t.Run("VRSNN-E01: A product whose scheme does not build refuses the configuration", func(t *testing.T) {
		refuses(t, "  app: { scheme: date, next: none }\n", "versioning.app", `"date"`)
		refuses(t, "  app:\n", "versioning.app")
	})
	t.Run("VRSNN-E02: A source that names no place refuses the configuration", func(t *testing.T) {
		refuses(t, "  app: { scheme: semver, source: [{ path: version }] }\n", "source[0]", "file")
		refuses(t, "  app: { scheme: semver, source: [{ file: a.json, path: v, regex: '(x)' }] }\n", "a.json", "one of them")
		refuses(t, "  app: { scheme: semver, source: [{ file: a.json }] }\n", "a.json", "one of them")
		refuses(t, "  app: { scheme: semver, source: [{ file: a.txt, regex: '(a)(b)' }] }\n", "a.txt", "2 groups")
	})
	t.Run("VRSNN-E03: A tag without its version, or with it twice, refuses the configuration", func(t *testing.T) {
		refuses(t, "  app: { scheme: number, tag: release, next: none }\n", `"release"`)
		refuses(t, "  app: { scheme: number, tag: 'v{version}-{version}', next: none }\n", "{version}-{version}")
	})
	t.Run("VRSNN-E04: A rule that cannot decide refuses the configuration", func(t *testing.T) {
		src := "source: [{ file: v.txt, regex: '(.+)' }]"
		for _, c := range []struct{ next, says string }{
			{"next: latest", `"latest"`},
			{"next: { strategy: increment, part: build }", `"build"`},
			{"next: { strategy: semantic, types: { feat: huge } }", `"huge"`},
			{"next: { strategy: semantic, part: patch }", "part is an option of increment"},
		} {
			refuses(t, "  app: { scheme: semver, "+src+", "+c.next+" }\n", "versioning.app", c.says)
		}
		refuses(t, "  app: { scheme: semver, "+src+", next: { strategy: semantic, bump: major } }\n", "next", "bump")
		refuses(t, "  app: { scheme: number, next: semantic }\n", "semantic needs scheme semver")
		refuses(t, "  app: { scheme: semver }\n", "manual reads the version from the sources")
		refuses(t, "  app:\n    scheme: number\n    next: none\n    branches:\n      review: { label: rc }\n", "label")
		refuses(t, "  app:\n    scheme: number\n    next: none\n    branches:\n      're:(': { next: none }\n", "branches", "re:(")
	})
}
