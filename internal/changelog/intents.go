// @anchors
//   code: INTNV
//   ref: INTNA

package changelog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// A change intent is what a change is to the people a changelog is for, written by whoever
// made it, in a file of its own: a front matter that names the product, how far the change
// moves its version and its type, and the text. The commits say what changed in the code;
// an intent says it in the product's words, and it is consumed by the release it ships in.

// Intent is one change intent.
type Intent struct {
	File    string // relative to the intents' directory
	Product string
	Bump    string // a part of the product's scheme, or none
	Type    string
	Text    string
}

// ParseIntent reads an intent: a `---` front matter with product, bump and type, then the
// text. A file without the front matter, or with a key it does not know, is an error naming
// the file: an intent half read would ship half a change.
func ParseIntent(file string, content []byte) (Intent, error) {
	in := Intent{File: file}
	s := strings.ReplaceAll(string(content), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return in, fmt.Errorf("%s: an intent opens with a --- front matter", file)
	}
	head, body, ok := strings.Cut(s[4:], "\n---")
	if !ok {
		return in, fmt.Errorf("%s: the front matter is not closed by ---", file)
	}
	var fm struct {
		Product string `yaml:"product"`
		Bump    string `yaml:"bump"`
		Type    string `yaml:"type"`
	}
	dec := yaml.NewDecoder(bytes.NewReader([]byte(head)))
	dec.KnownFields(true)
	if err := dec.Decode(&fm); err != nil {
		return in, fmt.Errorf("%s: the front matter: %v", file, err)
	}
	in.Product, in.Bump, in.Type = fm.Product, fm.Bump, fm.Type
	in.Text = strings.TrimSpace(strings.TrimPrefix(body, "\n"))
	return in, nil
}

// ReadIntents reads the intents of a directory, `*.md` sorted by name; a directory that does
// not exist holds none. Each file that does not parse is an error of its own, and the others
// are still read.
func ReadIntents(dir string) ([]Intent, []error) {
	paths, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	sort.Strings(paths)
	var out []Intent
	var errs []error
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		in, err := ParseIntent(filepath.Base(p), b)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		out = append(out, in)
	}
	return out, errs
}

// IntentRelease is what an intents template renders: the release's heading, as a release
// from the commits has it, and its intents by type, in the order the types are declared.
type IntentRelease struct {
	Version, Date, Candidate string
	Sections                 []IntentSection
}

// IntentSection is the intents of one type.
type IntentSection struct {
	Type    string
	Intents []Intent
}

// GroupIntents builds a release from intents, their types in the order given; an intent of
// another type goes after them, under its own.
func GroupIntents(version, date string, intents []Intent, types []string) IntentRelease {
	r := IntentRelease{Version: version, Date: date}
	idx := map[string]int{}
	add := func(t string) int {
		if i, ok := idx[t]; ok {
			return i
		}
		idx[t] = len(r.Sections)
		r.Sections = append(r.Sections, IntentSection{Type: t})
		return idx[t]
	}
	for _, t := range types {
		add(t)
	}
	for _, in := range intents {
		i := add(in.Type)
		r.Sections[i].Intents = append(r.Sections[i].Intents, in)
	}
	kept := r.Sections[:0]
	for _, s := range r.Sections {
		if len(s.Intents) > 0 {
			kept = append(kept, s)
		}
	}
	r.Sections = kept
	return r
}

// DefaultIntentsTemplate renders a release of intents: the heading the commits' template
// writes, and one section per type, its heading translated by `t` (`changelog.type.<type>`).
const DefaultIntentsTemplate = `## {{if .Version}}{{.Version}}{{if .Date}} — {{.Date}}{{end}}{{else}}{{if .Candidate}}{{.Candidate}} — {{end}}{{t "changelog.unreleased"}}{{end}}
{{- range .Sections}}

### {{t (printf "changelog.type.%s" .Type)}}
{{range .Intents}}
- {{.Text}}{{end}}{{end}}
`

// RenderIntents renders a release of intents with the template text.
func RenderIntents(tmpl string, r IntentRelease, t func(string) string) (string, error) {
	tp, err := template.New("intents").Funcs(template.FuncMap{"t": t}).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("changelog template: %w", err)
	}
	var b bytes.Buffer
	if err := tp.Execute(&b, r); err != nil {
		return "", fmt.Errorf("changelog template: %w", err)
	}
	return b.String(), nil
}
