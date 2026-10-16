// @anchors
//   code: CMCNC
//   ref: CMMRC

// Package config carrega os defaults embutidos e a config do projeto.
package config

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

// Comments is how the files a glob matches write their comments (`dialect.comments`): the
// line markers, the block delimiters, and the string delimiters inside which a marker is
// text. The header (`@anchors`) and the chain's flags live in comments, and the evidence
// leaves comments out — so this is under the map, the evidence and every gate that reads
// code. The core knows none of it (DESIGN-core-without-syntax.md): a file no entry covers has
// no comment, and the gates that need one say so. Anchors' own formats are the exception —
// its specs and plans in Markdown, its features in Gherkin (ownComments).
type Comments struct {
	// Files is the glob of the files the entry speaks of (`**/*.{ts,tsx}`); empty is every
	// file.
	Files string `yaml:"files,omitempty"`
	// Line are the markers that open a comment running to the end of the line.
	Line []string `yaml:"line,omitempty"`
	// Block are the delimiters of a comment that may span lines, as `[open, close]`.
	Block []Delims `yaml:"block,omitempty"`
	// Strings are the delimiters of a string: inside one, a marker is text. A string ends at
	// its delimiter or at the end of its line, unless the delimiter is in Multiline.
	Strings []string `yaml:"strings,omitempty"`
	// Multiline are the string delimiters whose strings may span lines (a template literal,
	// a triple-quoted string).
	Multiline []string `yaml:"multiline,omitempty"`
	// NoLineAfter matches a line whose end takes no line comment — what follows is not code
	// (JSX text after a tag): a flag written there goes in the first block comment, right
	// after the code it speaks of.
	NoLineAfter string `yaml:"no_line_after,omitempty"`
}

// Delims are a block comment's opener and closer, written `[open, close]`.
type Delims struct {
	Open, Close string
}

// UnmarshalYAML reads `[open, close]`.
func (d *Delims) UnmarshalYAML(n *yaml.Node) error {
	var pair []string
	if err := n.Decode(&pair); err != nil || len(pair) != 2 || pair[0] == "" || pair[1] == "" {
		return fmt.Errorf("a block comment is `[open, close]`, got %v", n.Value)
	}
	d.Open, d.Close = pair[0], pair[1]
	return nil
}

// MarshalYAML writes `[open, close]`.
func (d Delims) MarshalYAML() (any, error) { return []string{d.Open, d.Close}, nil }

// ownComments are the comments of Anchors' own formats: its specs, plans and docs in
// Markdown, its features in Gherkin. They are Anchors' syntax, not a project's language.
var ownComments = []Comments{
	{Files: "**/*.md", Block: []Delims{{"<!--", "-->"}}},
	{Files: "**/*.feature", Line: []string{"#"}},
}

// CommentsFor is the comment syntax of a file — its path relative to the root, with `/` —:
// Anchors' own for its formats, else the first entry of `dialect.comments` whose glob matches
// it; nil when nothing says how the file comments.
func (c *Config) CommentsFor(rel string) *Comments {
	if c != nil && c.Dialect != nil {
		return c.Dialect.CommentsFor(rel)
	}
	return Dialect{}.CommentsFor(rel)
}

// CommentsFor is Config.CommentsFor on the dialect alone.
func (d Dialect) CommentsFor(rel string) *Comments {
	rel = strings.TrimPrefix(path.Clean(strings.ReplaceAll(rel, "\\", "/")), "./")
	// Anchors' own formats first: a spec is Markdown whatever an entry for every file says.
	for i := range ownComments {
		if ownComments[i].matches(rel) {
			return &ownComments[i]
		}
	}
	for i := range d.Comments {
		if d.Comments[i].matches(rel) {
			return &d.Comments[i]
		}
	}
	return nil
}

func (cs *Comments) matches(rel string) bool {
	if strings.TrimSpace(cs.Files) == "" {
		return true
	}
	if ok, _ := doublestar.Match(cs.Files, rel); ok {
		return true
	}
	// A glob written for the root's files also reads them at the root: `**/*.go` matches
	// `main.go`.
	ok, _ := doublestar.Match(strings.TrimPrefix(cs.Files, "**/"), rel)
	return ok
}

// validate says what is wrong with an entry: a glob that is no glob, no marker at all, a
// marker or delimiter that is empty, an unknown NoLineAfter pattern.
func (cs Comments) validate(field string) error {
	if cs.Files != "" && !doublestar.ValidatePattern(cs.Files) {
		return fmt.Errorf("`%s.files` is not a valid glob: %s", field, cs.Files)
	}
	if len(cs.Line) == 0 && len(cs.Block) == 0 {
		return fmt.Errorf("`%s` declares no marker: a `line` or a `block`", field)
	}
	for _, m := range append(append([]string{}, cs.Line...), cs.Strings...) {
		if m == "" {
			return fmt.Errorf("`%s` has an empty marker", field)
		}
	}
	if cs.NoLineAfter != "" {
		if _, err := regexp.Compile(cs.NoLineAfter); err != nil {
			return fmt.Errorf("`%s.no_line_after` does not compile: %w", field, err)
		}
	}
	return nil
}

// LineMarker is the marker a flag or a header is written with in the file: the first line
// marker, or the first block's opener (closed by CloseFor). Empty for nil.
func (cs *Comments) LineMarker() string {
	switch {
	case cs == nil:
		return ""
	case len(cs.Line) > 0:
		return cs.Line[0]
	case len(cs.Block) > 0:
		return cs.Block[0].Open
	}
	return ""
}

// CloseFor is what closes a comment opened with the marker, with its leading space — ` -->`
// after `<!--`, ` *)` after `(*` —, or empty for a line marker, which the end of the line
// closes.
func (cs *Comments) CloseFor(marker string) string {
	if cs == nil {
		return ""
	}
	for _, b := range cs.Block {
		if b.Open == marker {
			return " " + b.Close
		}
	}
	return ""
}

// TakesLineComment says whether a comment may be written at the end of this line: not on a
// line NoLineAfter matches.
func (cs *Comments) TakesLineComment(line string) bool {
	if cs == nil || cs.NoLineAfter == "" {
		return true
	}
	re, err := regexp.Compile(cs.NoLineAfter)
	return err != nil || !re.MatchString(line)
}

// Blank is the content with every comment turned into spaces — its lines and offsets kept —,
// read with the strings it knows: a marker inside a string is text. Nil blanks nothing.
func (cs *Comments) Blank(src string) string {
	if cs == nil {
		return src
	}
	out := []byte(src)
	for _, sp := range cs.Spans(src) {
		blank(out, sp[0], sp[1])
	}
	return string(out)
}

// Spans are the byte ranges [from, to) of the content's comments, markers included, in
// order. Nil has none.
func (cs *Comments) Spans(src string) [][2]int {
	if cs == nil {
		return nil
	}
	var out [][2]int
	inString := ""
	for i := 0; i < len(src); {
		if inString != "" {
			switch {
			case src[i] == '\\' && len(inString) == 1:
				i += 2
			case strings.HasPrefix(src[i:], inString):
				i += len(inString)
				inString = ""
			case src[i] == '\n' && !cs.multiline(inString):
				inString = ""
				i++
			default:
				i++
			}
			continue
		}
		if s := cs.stringAt(src[i:]); s != "" {
			inString = s
			i += len(s)
			continue
		}
		if b, ok := cs.blockAt(src[i:]); ok {
			end := strings.Index(src[i+len(b.Open):], b.Close)
			to := len(src)
			if end >= 0 {
				to = i + len(b.Open) + end + len(b.Close)
			}
			out = append(out, [2]int{i, to})
			i = to
			continue
		}
		if m := cs.lineAt(src[i:]); m != "" {
			end := strings.IndexByte(src[i:], '\n')
			to := len(src)
			if end >= 0 {
				to = i + end
			}
			out = append(out, [2]int{i, to})
			i = to
			continue
		}
		i++
	}
	return out
}

func blank(b []byte, from, to int) {
	for j := from; j < to && j < len(b); j++ {
		if b[j] != '\n' {
			b[j] = ' '
		}
	}
}

// stringAt is the string delimiter that opens at the start of s, the longest first.
func (cs *Comments) stringAt(s string) string {
	best := ""
	for _, d := range cs.Strings {
		if strings.HasPrefix(s, d) && len(d) > len(best) {
			best = d
		}
	}
	return best
}

func (cs *Comments) blockAt(s string) (Delims, bool) {
	var best Delims
	for _, b := range cs.Block {
		if strings.HasPrefix(s, b.Open) && len(b.Open) > len(best.Open) {
			best = b
		}
	}
	return best, best.Open != ""
}

func (cs *Comments) lineAt(s string) string {
	for _, m := range cs.Line {
		if strings.HasPrefix(s, m) {
			return m
		}
	}
	return ""
}

func (cs *Comments) multiline(d string) bool {
	for _, m := range cs.Multiline {
		if m == d {
			return true
		}
	}
	return false
}

// InString marks, by index, the lines that START inside a string spanning lines (a
// delimiter of Multiline): text, not code — a line of a template literal is no closer of a
// block. Nil marks none.
func (cs *Comments) InString(lines []string) []bool {
	out := make([]bool, len(lines))
	if cs == nil || len(cs.Multiline) == 0 {
		return out
	}
	src := strings.Join(lines, "\n")
	inString, inBlock, line := "", "", 0
	for i := 0; i < len(src); i++ {
		if src[i] == '\n' {
			line++
			if inString != "" && !cs.multiline(inString) {
				inString = ""
			}
			if line < len(out) {
				out[line] = inString != ""
			}
			continue
		}
		switch {
		case inBlock != "":
			if strings.HasPrefix(src[i:], inBlock) {
				i += len(inBlock) - 1
				inBlock = ""
			}
		case inString != "":
			switch {
			case src[i] == '\\' && len(inString) == 1:
				i++
			case strings.HasPrefix(src[i:], inString):
				i += len(inString) - 1
				inString = ""
			}
		default:
			if d := cs.stringAt(src[i:]); d != "" {
				inString = d
				i += len(d) - 1
			} else if b, ok := cs.blockAt(src[i:]); ok {
				inBlock = b.Close
				i += len(b.Open) - 1
			} else if cs.lineAt(src[i:]) != "" {
				for i+1 < len(src) && src[i+1] != '\n' {
					i++
				}
			}
		}
	}
	return out
}

// CommentOnly marks, by index, the lines of the content that hold a comment and nothing
// else: not blank, and blank once the comments are.
func (cs *Comments) CommentOnly(lines []string) []bool {
	out := make([]bool, len(lines))
	if cs == nil {
		return out
	}
	blanked := strings.Split(cs.Blank(strings.Join(lines, "\n")), "\n")
	for i := range lines {
		out[i] = strings.TrimSpace(lines[i]) != "" && i < len(blanked) && strings.TrimSpace(blanked[i]) == ""
	}
	return out
}

// Code is the content without its comments: each line as its code, trimmed on the right,
// and the lines that were only a comment dropped.
func (cs *Comments) Code(src string) string {
	lines := strings.Split(src, "\n")
	blanked := strings.Split(cs.Blank(src), "\n")
	var b strings.Builder
	for i, l := range blanked {
		if strings.TrimSpace(l) == "" && strings.TrimSpace(lines[i]) != "" {
			continue
		}
		b.WriteString(strings.TrimRight(l, " \t\r"))
		b.WriteString("\n")
	}
	return b.String()
}

// CutLine is one line without its comment, trimmed on the right.
func (cs *Comments) CutLine(line string) string {
	return strings.TrimRight(cs.Blank(line), " \t\r")
}

// OpensComment says whether a trimmed line starts with a marker or a block's opener.
func (cs *Comments) OpensComment(t string) bool {
	if cs == nil {
		return false
	}
	if cs.lineAt(t) != "" {
		return true
	}
	_, ok := cs.blockAt(t)
	return ok
}

// FamilyComments is the family's example of `dialect.comments`, for the files of its stack
// and the configuration and scripts beside them: what `anchors init` and `anchors migrate`
// write and the doctor offers. It is never applied by itself. AllComments is every family's.
func FamilyComments(family string) []Comments {
	cLike := func(files string) Comments {
		return Comments{Files: files, Line: []string{"//"}, Block: []Delims{{"/*", "*/"}}, Strings: []string{`"`, "'"}}
	}
	hash := Comments{Files: "**/*.{yaml,yml,toml,sh,bash,zsh}", Line: []string{"#"}, Strings: []string{`"`, "'"}}
	switch strings.ToLower(family) {
	case "ts":
		return []Comments{
			{Files: "**/*.{ts,tsx,js,jsx,mjs,cjs,mts,cts}", Line: []string{"//"},
				// `{/* … */}` is a comment inside JSX.
				Block:   []Delims{{"{/*", "*/}"}, {"/*", "*/"}},
				Strings: []string{`"`, "'", "`"}, Multiline: []string{"`"},
				// After a JSX tag the line is text the app renders.
				NoLineAfter: `>\s*$`},
			hash,
		}
	case "go":
		goC := cLike("**/*.go")
		goC.Strings, goC.Multiline = []string{`"`, "'", "`"}, []string{"`"}
		return []Comments{goC, hash}
	case "python":
		return []Comments{
			{Files: "**/*.py", Line: []string{"#"}, Strings: []string{`"""`, "'''", `"`, "'"}, Multiline: []string{`"""`, "'''"}},
			hash,
		}
	case "java", "kotlin", "csharp", "scala", "swift":
		ext := map[string]string{"java": "java", "kotlin": "{kt,kts}", "csharp": "cs", "scala": "scala", "swift": "swift"}[strings.ToLower(family)]
		return []Comments{cLike("**/*." + ext), hash}
	case "rust":
		return []Comments{cLike("**/*.rs"), hash}
	case "php":
		php := cLike("**/*.php")
		php.Line = []string{"//", "#"}
		return []Comments{php, hash}
	case "ruby":
		return []Comments{{Files: "**/*.rb", Line: []string{"#"}, Strings: []string{`"`, "'"}}, hash}
	case "ocaml":
		return []Comments{
			{Files: "**/*.{ml,mli,mll}", Block: []Delims{{"(*", "*)"}}, Strings: []string{`"`}},
			{Files: "**/*.mly", Block: []Delims{{"(*", "*)"}, {"/*", "*/"}}, Strings: []string{`"`}},
			{Files: "**/{dune,dune-project,dune-file}", Line: []string{";"}, Strings: []string{`"`}},
		}
	}
	return nil
}

// CommentsCovering are the entries of the examples (AllComments) that cover the files, in the
// examples' order, and the extensions — or the names, for a file with none — that no example
// covers, sorted. Anchors' own formats are read by Anchors and need none. It is what
// `anchors init` and `anchors migrate` write for a project's files.
func CommentsCovering(files []string) ([]Comments, []string) {
	all := AllComments()
	used := make([]bool, len(all))
	uncovered := map[string]bool{}
	for _, f := range files {
		if (Dialect{}).CommentsFor(f) != nil {
			continue
		}
		f = strings.TrimPrefix(path.Clean(strings.ReplaceAll(f, "\\", "/")), "./")
		found := false
		for i := range all {
			if all[i].matches(f) {
				used[i], found = true, true
				break
			}
		}
		if !found {
			ext := path.Ext(f)
			if ext == "" {
				ext = path.Base(f)
			}
			uncovered[ext] = true
		}
	}
	var out []Comments
	for i, c := range all {
		if used[i] {
			out = append(out, c)
		}
	}
	exts := make([]string, 0, len(uncovered))
	for e := range uncovered {
		exts = append(exts, e)
	}
	sort.Strings(exts)
	return out, exts
}

// AllComments are the comment examples of every family, and of the formats beside them (SQL,
// Lua, Haskell, C, markup): `anchors migrate` writes into a project the ones whose files it
// has.
func AllComments() []Comments {
	out := []Comments{}
	seen := map[string]bool{}
	for _, fam := range []string{"ts", "go", "python", "java", "kotlin", "csharp", "scala", "swift", "rust", "php", "ruby", "ocaml"} {
		for _, c := range FamilyComments(fam) {
			if !seen[c.Files] {
				seen[c.Files] = true
				out = append(out, c)
			}
		}
	}
	return append(out,
		Comments{Files: "**/*.{c,h,cpp,hpp,cc}", Line: []string{"//"}, Block: []Delims{{"/*", "*/"}}, Strings: []string{`"`, "'"}},
		Comments{Files: "**/*.{sql,lua,hs}", Line: []string{"--"}, Strings: []string{"'", `"`}},
		Comments{Files: "**/*.{r,R}", Line: []string{"#"}, Strings: []string{`"`, "'"}},
		Comments{Files: "**/*.{html,xml,svg,vue}", Block: []Delims{{"<!--", "-->"}}},
	)
}
