// @anchors
//   code: CHVLV
//   ref: CHVLA

package gate

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/co2-lab/anchors/internal/changelog"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/gitmeta"
	"github.com/co2-lab/anchors/internal/i18n"
	"github.com/co2-lab/anchors/internal/mapx"
	"github.com/co2-lab/anchors/internal/version"
)

// changelog-valid: what a product's changelog outputs hold is confronted, whoever wrote it. A
// file is edited by hand, a template changes, and generating it once proves nothing of what
// it holds later. What does not depend on taste holds for every output — the markers read as
// versions, newest first, none twice, the unreleased block on top, a block for every tag and
// a tag for every block —; the rest is the output's `format:`, and an output with none is
// held to nothing more.

func checkChangelogValid(_ config.Gate, root string, _ *mapx.Graph, cfg *config.Config) (Verdict, string) {
	names := make([]string, 0, len(cfg.Versioning))
	for n, p := range cfg.Versioning {
		if p.Changelog != nil && len(p.Changelog.Outputs) > 0 {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	intentFails, hasIntents := confrontIntents(root, cfg)
	if len(names) == 0 && !hasIntents {
		return Skip, i18n.T("gate.changelog_valid.no_output")
	}
	fails := intentFails
	for _, name := range names {
		p := cfg.Versioning[name]
		rel, err := productTags(root, p)
		if err != nil {
			fails = append(fails, name+": "+err.Error())
			continue
		}
		if cur := sourcesVersion(root, p); cur != "" && rel.any[cur] == "" {
			rel.pending = cur
		}
		for _, o := range p.Changelog.Outputs {
			fails = append(fails, confrontOutput(root, p, o, rel)...)
		}
	}
	if len(fails) > 0 {
		return Fail, strings.Join(fails, "; ")
	}
	return Pass, i18n.T("gate.changelog_valid.pass", len(names))
}

// confrontIntents holds each change intent waiting for a release to what it may say, and
// says whether the project has intents to confront: declared, or a directory of them.
func confrontIntents(root string, cfg *config.Config) ([]string, bool) {
	var in *config.Intents
	if cfg.Changelog != nil {
		in = cfg.Changelog.Intents
	}
	dir := in.PathOrDefault()
	abs := filepath.Join(root, filepath.FromSlash(dir))
	if _, err := os.Stat(abs); err != nil {
		return nil, in != nil
	}
	intents, errs := changelog.ReadIntents(abs)
	var fails []string
	for _, err := range errs {
		fails = append(fails, dir+"/"+err.Error())
	}
	types := in.TypesOrDefault()
	latest := map[string]string{}
	for _, it := range intents {
		where := dir + "/" + it.File
		p, ok := cfg.Versioning[it.Product]
		if !ok {
			fails = append(fails, i18n.T("gate.changelog_valid.intent_product", where, it.Product))
			continue
		}
		s, _ := p.SchemeOf()
		if it.Bump != "none" && (it.Bump == "" || !s.ValidPart(it.Bump)) {
			fails = append(fails, i18n.T("gate.changelog_valid.intent_bump", where, it.Bump, strings.Join(append(s.Parts(), "none"), ", ")))
		}
		if !containsFile(types, it.Type) {
			fails = append(fails, i18n.T("gate.changelog_valid.intent_type", where, it.Type, strings.Join(types, ", ")))
		}
		if it.Text == "" {
			fails = append(fails, i18n.T("gate.changelog_valid.intent_empty", where))
		} else if n := utf8.RuneCountInString(it.Text); in != nil && in.MaxChars > 0 && n > in.MaxChars {
			fails = append(fails, i18n.T("gate.changelog_valid.too_long", where, n, in.MaxChars))
		}
		tag, seen := latest[it.Product]
		if !seen {
			if rel, err := productTags(root, p); err == nil {
				tag = rel.latest
			}
			latest[it.Product] = tag
		}
		if tag != "" && addedBefore(root, dir+"/"+it.File, tag) {
			fails = append(fails, i18n.T("gate.changelog_valid.intent_stale", where, tag))
		}
	}
	return fails, true
}

// addedBefore says whether the commit that added a file is in the release a tag records: an
// intent that release should have consumed.
func addedBefore(root, file, tag string) bool {
	out, err := exec.Command("git", "-C", root, "log", "--diff-filter=A", "--format=%H", "--", file).Output()
	lines := strings.Fields(string(out))
	if err != nil || len(lines) == 0 {
		return false
	}
	return exec.Command("git", "-C", root, "merge-base", "--is-ancestor", lines[len(lines)-1], tag).Run() == nil
}

// releases are a product's tags: those reachable from HEAD, every one the repository has,
// and the date of each.
type releases struct {
	scheme    version.Scheme
	reachable map[string]string // version → tag
	any       map[string]string // version → tag
	latest    string            // the tag of the highest reachable version
	date      map[string]string // version → YYYY-MM-DD
	pending   string            // the version the sources hold, when no tag records it yet
}

func productTags(root string, p *config.Product) (*releases, error) {
	s, _ := p.SchemeOf()
	r := &releases{scheme: s, reachable: map[string]string{}, any: map[string]string{}, date: map[string]string{}}
	list := func(args ...string) []string {
		out, _ := exec.Command("git", append([]string{"-C", root, "tag"}, args...)...).Output()
		return strings.Fields(string(out))
	}
	var best *version.V
	for _, t := range list("--merged", "HEAD") {
		raw, ok := p.VersionOfTag(t)
		if !ok {
			continue
		}
		v, err := s.Parse(raw)
		if err != nil {
			return nil, err
		}
		r.reachable[raw] = t
		if best == nil || version.Compare(v, *best) > 0 {
			best, r.latest = &v, t
		}
		out, _ := exec.Command("git", "-C", root, "log", "-1", "--format=%cs", t).Output()
		r.date[raw] = strings.TrimSpace(string(out))
	}
	for _, t := range list("--list") {
		if raw, ok := p.VersionOfTag(t); ok {
			r.any[raw] = t
		}
	}
	return r, nil
}

// fileBlock is a block of an output, with where it came from.
type fileBlock struct {
	changelog.WrittenBlock
	file string // the file it is in, relative to the root
}

// confrontOutput reads an output's blocks and confronts them; each failure is one string.
func confrontOutput(root string, p *config.Product, o config.ChangelogOutput, rel *releases) []string {
	var blocks []fileBlock
	var fails []string
	f := o.Format
	if f == nil {
		f = &config.ChangelogFormat{}
	}
	if o.ModeOrDefault() == "per_version" {
		pattern := f.File
		if pattern == "" {
			pattern = "{version}.md"
		}
		nameRE := regexp.MustCompile("^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), regexp.QuoteMeta("{version}"), `(.+)`) + "$")
		entries, _ := os.ReadDir(filepath.Join(root, filepath.FromSlash(o.Path)))
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			rel := o.Path + "/" + e.Name()
			b, _ := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
			if e.Name() == "unreleased.md" {
				blocks = append(blocks, fileBlock{changelog.WrittenBlock{Version: "unreleased", Text: string(b)}, rel})
				continue
			}
			m := nameRE.FindStringSubmatch(e.Name())
			if m == nil {
				fails = append(fails, i18n.T("gate.changelog_valid.file_name", rel, pattern))
				continue
			}
			blocks = append(blocks, fileBlock{changelog.WrittenBlock{Version: m[1], Text: string(b)}, rel})
		}
	} else {
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(o.Path)))
		if err != nil && len(rel.reachable) > 0 {
			return []string{o.Path + ": " + i18n.T("gate.changelog_valid.missing", len(rel.reachable))}
		}
		for _, wb := range changelog.Blocks(string(b)) {
			blocks = append(blocks, fileBlock{wb, o.Path})
		}
		for _, f := range orderOf(blocks, rel.scheme) {
			fails = append(fails, o.Path+": "+f)
		}
	}
	seen := map[string]bool{}
	for _, b := range blocks {
		seen[b.Version] = true
		fails = append(fails, confrontBlock(root, b, f, rel)...)
	}
	var missing []string
	for v := range rel.reachable {
		if !seen[v] {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		sortVersions(missing, rel.scheme)
		fails = append(fails, o.Path+": "+i18n.T("gate.changelog_valid.no_block", strings.Join(missing, ", ")))
	}
	if f.Required == "changed" {
		if msg := changedSinceRelease(root, p, o, rel); msg != "" {
			fails = append(fails, o.Path+": "+msg)
		}
	}
	return fails
}

// orderOf confronts the order of an incremental file's blocks: the unreleased one at most
// once and on top, the versions newest first and none twice.
func orderOf(blocks []fileBlock, s version.Scheme) []string {
	var fails []string
	var prev *version.V
	seen := map[string]bool{}
	for i, b := range blocks {
		if b.Version == "unreleased" {
			if i > 0 {
				fails = append(fails, i18n.T("gate.changelog_valid.unreleased_not_on_top"))
			}
			continue
		}
		if seen[b.Version] {
			fails = append(fails, i18n.T("gate.changelog_valid.twice", b.Version))
			continue
		}
		seen[b.Version] = true
		v, err := s.Parse(b.Version)
		if err != nil {
			continue // confrontBlock names it
		}
		if prev != nil && version.Compare(v, *prev) >= 0 {
			fails = append(fails, i18n.T("gate.changelog_valid.order", b.Version, prev.Raw))
		}
		prev = &v
	}
	return fails
}

var sectionRE = regexp.MustCompile(`(?m)^###\s+(.+?)\s*$`)

// confrontBlock holds one block to its version, its tag and the output's format.
func confrontBlock(root string, b fileBlock, f *config.ChangelogFormat, rel *releases) []string {
	var fails []string
	where := b.file
	if b.Version != "unreleased" {
		where += " " + b.Version
		if _, err := rel.scheme.Parse(b.Version); err != nil {
			return []string{where + ": " + err.Error()}
		}
		if _, ok := rel.any[b.Version]; !ok && b.Version != rel.pending {
			fails = append(fails, i18n.T("gate.changelog_valid.no_tag", where))
		}
	}
	text := strings.ReplaceAll(b.Text, "\r\n", "\n")
	if f.Heading != "" && b.Version != "unreleased" {
		re, _ := config.HeadingRE(f.Heading)
		first := strings.SplitN(strings.TrimLeft(text, "\n"), "\n", 2)[0]
		m := re.FindStringSubmatch(first)
		switch {
		case m == nil:
			fails = append(fails, i18n.T("gate.changelog_valid.heading", where, first, f.Heading))
		default:
			if i := re.SubexpIndex("version"); i > 0 && m[i] != b.Version {
				fails = append(fails, i18n.T("gate.changelog_valid.heading_version", where, m[i]))
			}
			if i := re.SubexpIndex("date"); i > 0 && rel.date[b.Version] != "" && m[i] != rel.date[b.Version] {
				fails = append(fails, i18n.T("gate.changelog_valid.date", where, m[i], rel.date[b.Version]))
			}
		}
	}
	sections := sectionRE.FindAllStringSubmatchIndex(text, -1)
	have := map[string]bool{}
	for i, s := range sections {
		title := text[s[2]:s[3]]
		have[title] = true
		end := len(text)
		if i+1 < len(sections) {
			end = sections[i+1][0]
		}
		if strings.TrimSpace(text[s[1]:end]) == "" {
			fails = append(fails, i18n.T("gate.changelog_valid.empty_section", where, title))
		}
	}
	for _, want := range f.RequiredSections {
		if !have[want] {
			fails = append(fails, i18n.T("gate.changelog_valid.no_section", where, want))
		}
	}
	if n := utf8.RuneCountInString(strings.TrimSpace(text)); f.MaxChars > 0 && n > f.MaxChars {
		fails = append(fails, i18n.T("gate.changelog_valid.too_long", where, n, f.MaxChars))
	}
	if _, tagged := rel.any[b.Version]; tagged {
		if before, ok := blockAtHead(root, b); ok && strings.TrimSpace(before) != strings.TrimSpace(text) {
			fails = append(fails, i18n.T("gate.changelog_valid.rewritten", where))
		}
	}
	return fails
}

// blockAtHead is the text the block had in the file as HEAD holds it; false when it had none.
func blockAtHead(root string, b fileBlock) (string, bool) {
	content, ok := gitmeta.AtHead(root, b.file)
	if !ok {
		return "", false
	}
	if strings.Contains(content, "<!-- anchors:changelog ") {
		for _, wb := range changelog.Blocks(content) {
			if wb.Version == b.Version {
				return strings.ReplaceAll(wb.Text, "\r\n", "\n"), true
			}
		}
		return "", false
	}
	return strings.ReplaceAll(content, "\r\n", "\n"), true
}

// changedSinceRelease asks an output to have changed since the latest release when the
// version the sources hold is not released yet: notes left as the previous release wrote
// them would ship with the next one.
func changedSinceRelease(root string, p *config.Product, o config.ChangelogOutput, rel *releases) string {
	if rel.latest == "" || rel.pending == "" {
		return ""
	}
	cur := rel.pending
	if exec.Command("git", "-C", root, "diff", "--quiet", rel.latest, "--", o.Path).Run() == nil {
		return i18n.T("gate.changelog_valid.unchanged", cur, rel.latest)
	}
	return ""
}

// sourcesVersion is the version a product's sources agree on in the tree; empty with no
// source, or when they cannot be read — version-forward confronts that.
func sourcesVersion(root string, p *config.Product) string {
	if len(p.Source) == 0 {
		return ""
	}
	svs, err := p.ReadSources(func(f string) ([]byte, error) { return os.ReadFile(filepath.Join(root, filepath.FromSlash(f))) })
	if err != nil {
		return ""
	}
	cur, _ := config.Agreed(svs)
	return cur
}

func sortVersions(vs []string, s version.Scheme) {
	sort.Slice(vs, func(i, j int) bool {
		a, _ := s.Parse(vs[i])
		b, _ := s.Parse(vs[j])
		return version.Compare(a, b) > 0
	})
}
