// Package changelog builds a project's changelog from its commits.
//
// The commits are the source, written in Conventional Commits (the `commit-msg` hook keeps
// them that way): the TYPE says what a change is to whoever reads the changelog, and two
// footers say what the type cannot. `BREAKING CHANGE:` (or a `!` after the type) marks a
// change that breaks what was there. `Bug:` marks a `fix` of a defect that SHIPPED — reached
// a release or production — and says where it was seen; a `fix` without it corrects work
// that never reached anyone.
//
// A release is what the history holds between two tags. Each release lists its breaking
// changes, its features (`feat`), the bugs it fixed (`fix` with `Bug:`) and its other fixes
// (`fix` without it); the internal types (`refactor`, `test`, `chore`, …) are history, not
// changelog, and are left out.
package changelog

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

// Entry is one change of a release.
type Entry struct {
	Type     string
	Scope    string
	Subject  string
	Hash     string // short
	Bug      string // the `Bug:` footer: where the bug was seen
	Breaking bool
}

// Release is what the history holds between two tags, or after the last one.
type Release struct {
	Version  string // the tag; empty for what is not released yet
	Date     string // the tag's date, YYYY-MM-DD; empty when unreleased
	Breaking []Entry
	Features []Entry
	Bugs     []Entry
	Fixes    []Entry
}

// Empty says whether the release has nothing a changelog lists.
func (r Release) Empty() bool {
	return len(r.Breaking)+len(r.Features)+len(r.Bugs)+len(r.Fixes) == 0
}

// Commit is a commit as the history gives it: hash, subject, body.
type Commit struct {
	Hash, Subject, Body string
}

var headerRE = regexp.MustCompile(`^([a-z]+)(\(([^)]*)\))?(!)?: +(.+)$`)

// Classify puts the commits of a release in its sections, in the order given. A commit
// outside Conventional Commits, or of a type the changelog does not list, is left out,
// unless it declares a breaking change: that is never silent.
func Classify(version, date string, commits []Commit) Release {
	r := Release{Version: version, Date: date}
	for _, c := range commits {
		m := headerRE.FindStringSubmatch(strings.TrimSpace(c.Subject))
		if m == nil {
			continue
		}
		e := Entry{Type: m[1], Scope: m[3], Subject: m[5], Hash: c.Hash}
		footers := lastParagraph(c.Body)
		if b, ok := footers["BREAKING CHANGE"]; ok || m[4] == "!" {
			e.Breaking = true
			if ok && b != "" && m[4] != "!" {
				e.Subject = e.Subject + " — " + b
			}
		}
		e.Bug = footers["Bug"]
		switch {
		case e.Breaking:
			r.Breaking = append(r.Breaking, e)
		case e.Type == "feat":
			r.Features = append(r.Features, e)
		case e.Type == "fix" && e.Bug != "":
			r.Bugs = append(r.Bugs, e)
		case e.Type == "fix":
			r.Fixes = append(r.Fixes, e)
		}
	}
	return r
}

// lastParagraph reads the footers of a body: `Key: value` lines of its last paragraph,
// where git keeps trailers.
func lastParagraph(body string) map[string]string {
	paras := strings.Split(strings.TrimSpace(strings.ReplaceAll(body, "\r\n", "\n")), "\n\n")
	out := map[string]string{}
	if len(paras) == 0 {
		return out
	}
	for _, l := range strings.Split(paras[len(paras)-1], "\n") {
		k, v, ok := strings.Cut(strings.TrimSpace(l), ":")
		if !ok {
			continue
		}
		out[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return out
}

// Git reads releases from a repository.
type Git struct {
	Root string
}

func (g Git) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.Root
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
	}
	return out.String(), nil
}

// Tags lists the tags reachable from `to`, oldest first by version order.
func (g Git) Tags(to string) ([]string, error) {
	out, err := g.run("tag", "--merged", to, "--sort=v:refname")
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, t := range strings.Split(strings.TrimSpace(out), "\n") {
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags, nil
}

// Commits lists the commits in (from, to], oldest first. An empty `from` means from the
// start of the history.
func (g Git) Commits(from, to string) ([]Commit, error) {
	rng := to
	if from != "" {
		rng = from + ".." + to
	}
	const sep, end = "\x1f", "\x1e"
	out, err := g.run("log", "--reverse", "--format=%h"+sep+"%s"+sep+"%b"+end, rng)
	if err != nil {
		return nil, err
	}
	var cs []Commit
	for _, rec := range strings.Split(out, end) {
		rec = strings.TrimLeft(rec, "\n")
		if rec == "" {
			continue
		}
		p := strings.SplitN(rec, sep, 3)
		if len(p) < 3 {
			continue
		}
		cs = append(cs, Commit{Hash: p[0], Subject: p[1], Body: p[2]})
	}
	return cs, nil
}

// Date is the date of a tag, YYYY-MM-DD.
func (g Git) Date(ref string) (string, error) {
	out, err := g.run("log", "-1", "--format=%cs", ref)
	return strings.TrimSpace(out), err
}

// Releases builds the releases whose tags are in `tags` (a sub-list of every tag, in
// order), each from the previous tag of the whole list. With `unreleased`, what came
// after the last tag up to `to` is added as a release with no version.
func (g Git) Releases(all, tags []string, to string, unreleased bool) ([]Release, error) {
	prev := map[string]string{}
	for i, t := range all {
		if i > 0 {
			prev[t] = all[i-1]
		}
	}
	var out []Release
	for _, t := range tags {
		cs, err := g.Commits(prev[t], t)
		if err != nil {
			return nil, err
		}
		date, err := g.Date(t)
		if err != nil {
			return nil, err
		}
		out = append(out, Classify(t, date, cs))
	}
	if unreleased {
		last := ""
		if len(all) > 0 {
			last = all[len(all)-1]
		}
		cs, err := g.Commits(last, to)
		if err != nil {
			return nil, err
		}
		out = append(out, Classify("", "", cs))
	}
	// Newest first: a changelog is read from the top.
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}

// DefaultTemplate renders one release. `t` translates a heading key; the project's own
// template (`changelog.template`) receives the same data and functions.
const DefaultTemplate = `## {{if .Version}}{{.Version}}{{if .Date}} — {{.Date}}{{end}}{{else}}{{t "changelog.unreleased"}}{{end}}
{{- if .Breaking}}

### {{t "changelog.breaking"}}
{{range .Breaking}}
- {{if .Scope}}**{{.Scope}}**: {{end}}{{.Subject}} ({{.Hash}}){{end}}{{end}}
{{- if .Features}}

### {{t "changelog.features"}}
{{range .Features}}
- {{if .Scope}}**{{.Scope}}**: {{end}}{{.Subject}} ({{.Hash}}){{end}}{{end}}
{{- if .Bugs}}

### {{t "changelog.bugs"}}
{{range .Bugs}}
- {{if .Scope}}**{{.Scope}}**: {{end}}{{.Subject}} ({{.Hash}}) — {{.Bug}}{{end}}{{end}}
{{- if .Fixes}}

### {{t "changelog.fixes"}}
{{range .Fixes}}
- {{if .Scope}}**{{.Scope}}**: {{end}}{{.Subject}} ({{.Hash}}){{end}}{{end}}
`

// Render renders a release with the template text, `t` translating the heading keys.
func Render(tmpl string, r Release, t func(string) string) (string, error) {
	tp, err := template.New("changelog").Funcs(template.FuncMap{"t": t}).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("changelog template: %w", err)
	}
	var b bytes.Buffer
	if err := tp.Execute(&b, r); err != nil {
		return "", fmt.Errorf("changelog template: %w", err)
	}
	return b.String(), nil
}

// Marker is the line a written release starts with: it is how a later write knows the
// release is already in the file. The text under it is the project's to edit.
func Marker(version string) string {
	if version == "" {
		version = "unreleased"
	}
	return "<!-- anchors:changelog " + version + " -->"
}

var markerRE = regexp.MustCompile(`(?m)^<!-- anchors:changelog (\S+) -->$`)

// Written lists the versions a changelog file already holds, by their markers; what is
// not released yet is `unreleased`.
func Written(existing string) map[string]bool {
	out := map[string]bool{}
	for _, m := range markerRE.FindAllStringSubmatch(existing, -1) {
		out[m[1]] = true
	}
	return out
}

// Block is a rendered release on its way into a file.
type Block struct {
	Version string // empty for what is not released yet
	Text    string
}

// Prepend writes the blocks at the top of an incremental changelog, newest first, and
// keeps everything the file already holds: a release already in it (by its marker) is
// not written again, and the text under a marker is never rewritten — it may have been
// edited by hand. The one exception is what is not released yet: that block is replaced,
// since it grows until the next tag. A file that starts with a `# ` title keeps it on
// top; an empty file gets `title`.
func Prepend(existing, title string, blocks []Block) string {
	existing = dropUnreleased(existing, blocks)
	have := Written(existing)
	var add []string
	for _, b := range blocks {
		key := b.Version
		if key == "" {
			key = "unreleased"
		}
		if have[key] {
			continue
		}
		add = append(add, Marker(b.Version)+"\n"+strings.TrimRight(b.Text, "\n")+"\n")
	}
	head, rest := "", existing
	if strings.HasPrefix(existing, "# ") {
		line, after, _ := strings.Cut(existing, "\n")
		head, rest = line+"\n\n", strings.TrimLeft(after, "\n")
	} else if strings.TrimSpace(existing) == "" {
		head, rest = "# "+title+"\n\n", ""
	}
	if len(add) == 0 {
		return head + rest
	}
	body := strings.Join(add, "\n")
	if rest != "" {
		body += "\n" + rest
	}
	return head + body
}

// dropUnreleased removes the unreleased block a previous write left, when a new one is
// on its way: from its marker to the next marker, or to the end of the file.
func dropUnreleased(existing string, blocks []Block) string {
	replacing := false
	for _, b := range blocks {
		replacing = replacing || b.Version == ""
	}
	start := strings.Index(existing, Marker(""))
	if !replacing || start < 0 {
		return existing
	}
	end := len(existing)
	if loc := markerRE.FindStringIndex(existing[start+len(Marker("")):]); loc != nil {
		end = start + len(Marker("")) + loc[0]
	}
	return existing[:start] + existing[end:]
}
