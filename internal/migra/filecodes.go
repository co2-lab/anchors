// @anchors
//   code: FLAPF
//   ref: MGFCD

package migra

import (
	"path"
	"regexp"
	"strings"

	"github.com/co2-lab/anchors/internal/code"
	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/scan"
)

// FORMAT 7 — every file has a code of its own, of five characters.
//
// A code always named a file: the spec carried its `code:`, and the code, the feature and
// the test only `ref:`d it — so a file could be named only through its unit, and a unit
// with more than one file could not name one of them. The dependency chain (`dep:` on the
// importer, `@used-by` on what it exports) needs every file to be addressable, so every
// governed file gets a `code:`, generated from its name and its type, beside the `ref:` it
// keeps. And the codes of four characters a project kept from before the default became
// five are widened, so every project that updates speaks the same length.
//
// The rewriting of the project is the command's (`anchors migrate`), which knows the disk
// and the map; this file holds the pure part: what a code becomes, what a file's code is,
// and where in a header the line goes.

// WidenedCode is the five-character code a four-character one becomes: the old code kept as
// the prefix, and one letter of the unit's name after it — `ARNA` of ArenaScreen becomes
// `ARNAS` —, so whoever knew the old code still reads it in the new one. The taken codes
// are never reused.
func WidenedCode(old, name string, taken map[string]bool) string {
	prev := code.Slots
	code.Slots = 5
	defer func() { code.Slots = prev }()
	return code.GenerateUniqueWithPrefix(name, old, taken)
}

// FileCodeName is what a file's code is generated from: its name without the artifact
// suffixes, and its type — the layer it belongs to, or its kind —, so the files of one unit
// (`ArenaScreen.tsx`, `ArenaScreen.feature`, `ArenaScreen.test.tsx`) get different codes.
func FileCodeName(f scan.File) string {
	b := path.Base(f.Path)
	for _, suf := range []string{".spec.md", ".feature", ".test", "_test", ".spec"} {
		if i := strings.Index(b, suf); i > 0 {
			b = b[:i]
		}
	}
	if i := strings.Index(b, "."); i > 0 {
		b = b[:i]
	}
	kind := f.Layer
	if kind == "" {
		kind = f.Kind
	}
	return b + " " + kind
}

// FileCode is a new five-character code for the file, unique among the taken ones.
func FileCode(f scan.File, taken map[string]bool) string {
	prev := code.Slots
	code.Slots = 5
	defer func() { code.Slots = prev }()
	return code.GenerateUnique(FileCodeName(f), taken)
}

// CanCarryCode says whether a file can be given a `code:` line: a text file whose type has
// a comment syntax the header is written in. A JSON, an image or a file of an unknown type
// would be broken by a comment.
func CanCarryCode(rel string, content []byte) bool {
	if strings.IndexByte(string(content), 0) >= 0 {
		return false
	}
	ext := strings.ToLower(path.Ext(rel))
	// Gherkin comments with `#`, and the scanner reads a feature's header as text: it is not
	// in the comment table, and a feature is a governed file like any other.
	if ext == ".feature" {
		return true
	}
	_, ok := config.CommentMarkers[ext]
	return ok
}

// WithHeaderCode is the content with `code: <c>` written in its `@anchors` header — right
// below the opener, so nothing someone wrote moves —, or with a header carrying it at the
// top when the file has none (after a shebang, which stays first).
func WithHeaderCode(content, rel, c string) string {
	marker := config.LineCommentFor(rel)
	field := "code: " + c
	if block := scan.AnchorsHeader([]byte(content)); block != nil {
		at := strings.Index(content, string(block))
		nl := strings.Index(string(block), "\n")
		if at < 0 || nl < 0 {
			return content
		}
		opener := at + nl
		line := marker + "   " + field
		if marker == "<!--" {
			line = "  " + field
		}
		return content[:opener] + "\n" + line + content[opener:]
	}
	header := marker + " @anchors\n" + marker + "   " + field + "\n\n"
	if marker == "<!--" {
		header = "<!-- @anchors\n  " + field + "\n-->\n\n"
	}
	if strings.HasPrefix(content, "#!") {
		if i := strings.Index(content, "\n"); i >= 0 {
			return content[:i+1] + header + content[i+1:]
		}
	}
	return header + content
}

// AsRefWithOwnCode is the content with the header's `code: <unit>` turned into `ref: <unit>`
// and a code of its own above it: the file that carried its spec's code — a feature, a code
// file written before the code was the file's — keeps naming its unit, and gets an address.
func AsRefWithOwnCode(content, unit, own string) string {
	block := scan.AnchorsHeader([]byte(content))
	if block == nil {
		return content
	}
	b := string(block)
	re := regexp.MustCompile(`(?m)^(.*?)code:(\s*)` + regexp.QuoteMeta(unit) + `\b`)
	m := re.FindStringSubmatchIndex(b)
	if m == nil {
		return content
	}
	prefix := b[m[2]:m[3]]
	nb := b[:m[0]] + prefix + "code: " + own + "\n" + prefix + "ref: " + unit + b[m[1]:]
	at := strings.Index(content, b)
	return content[:at] + nb + content[at+len(b):]
}
