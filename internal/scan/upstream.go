// @anchors
//   code: UPSCP
//   ref: UPOWP

package scan

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"

	"github.com/co2-lab/anchors/internal/config"
)

// UpstreamMarker identifies a file Anchors seeded and still owns — a pipeline that is
// still Anchors' template, not edited by the team. `anchors doctor --fix` keeps such a file
// up to date, and removing the marker is how a team takes the file over.
//
// Without the `#`: the marker lives both in a YAML (`# anchors:template`) and in an HTML
// (`<!-- anchors:template -->`), and tying it to one language's comment syntax would make
// the board page never match.
const UpstreamMarker = "anchors:template"

// UpstreamDir is where Anchors seeds the pipelines it owns.
const UpstreamDir = ".github/workflows"

// IsUpstreamOwned says whether a file is a VENDORED copy of an Anchors pipeline: it lives
// under `.github/workflows/` and still carries the template marker.
//
// Such a file belongs to the Anchors project, not to the one that received it. Its rules,
// its spec and its tests live upstream; the local copy is replaced whole by `doctor --fix`.
// Measured in the reference app: the map gave `anchors-claim.yml` no code, `anchors deliver`
// refused to record a change to it for lacking one, and `anchors-board.yml` got a code
// inferred from an example in its comments (`FNDTN`) — the project was being asked to own,
// and to write a spec for, files it does not own.
//
// The marker is the whole test, on purpose: a team that edits the file removes it, and from
// then on the file is theirs and governed like any other.
func IsUpstreamOwned(rel string, content []byte) bool {
	rel = strings.ReplaceAll(rel, `\`, "/")
	if !strings.HasPrefix(rel, UpstreamDir+"/") {
		return false
	}
	return bytes.Contains(content, []byte(UpstreamMarker))
}

// AnchorsHeader returns the text of the file's `@anchors` header block, or nil when the file
// has none — or when nothing says how the file writes a comment (`cs` nil): the header lives
// in a comment, and the core knows no language's (DESIGN-core-without-syntax.md).
//
// The block opens on the first line that is only a comment and whose comment STARTS with
// `@anchors` — not `@anchors-shared-code`, a mention in backticks, nor prose that names the
// header: a package comment saying "the header @anchors (code:/ref:)" was read as the
// header, and a header fix wrote the unit's `ref:` into the middle of that prose. It ends
// with the comment: at its block's closer, or at the first line that does not open with the
// same line marker. Reading a header key over the WHOLE file reads prose and code as
// declarations — a workflow's jq program with a line `parent: (...)` became a node's parent.
//
// The header is at the TOP: before it, only blank lines, comments and a shebang. A block
// further down is text — an example of a header in a guide, a string in the code — and is
// not the file's header: `guide_header.go` carried one inside a string, and the map read
// `layer: screen` as its unit's layer. A file whose header must stand lower (a directive
// the language wants first) says so inside the block, with the reason:
// `@fixed-header: <why>`.
func AnchorsHeader(content []byte, cs *config.Comments) []byte {
	// @no-guard: HeaderRange answers no header for a file with no comments
	lines := strings.Split(string(content), "\n")
	from, to, ok := HeaderRange(lines, cs)
	if !ok {
		return nil
	}
	return []byte(strings.Join(lines[from:to], "\n"))
}

// HeaderRange is where the header stands among the lines: [from, to), and whether there is
// one.
func HeaderRange(lines []string, cs *config.Comments) (int, int, bool) {
	if cs == nil {
		return 0, 0, false
	}
	only := cs.CommentOnly(lines)
	codeBefore := false
	for i, l := range lines {
		t := strings.TrimSpace(l)
		if only[i] && opensHeader(t) {
			end := headerEnd(lines, only, i, cs)
			if !codeBefore || fixedHeaderRE.MatchString(strings.Join(lines[i:end], "\n")) {
				return i, end, true
			}
			continue
		}
		if t != "" && !only[i] && !strings.HasPrefix(t, "#!") {
			codeBefore = true
		}
	}
	return 0, 0, false
}

// opensHeader says whether a comment line's text starts with `@anchors`, past the marker
// and any punctuation before it (`//`, `<!--`, ` * `), and ends the word there.
func opensHeader(t string) bool {
	rest := strings.TrimLeftFunc(t, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '@'
	})
	if !strings.HasPrefix(rest, "@anchors") {
		return false
	}
	after := rest[len("@anchors"):]
	return after == "" || after[0] == ' ' || after[0] == '\t' || !unicode.IsLetter(rune(after[0])) && after[0] != '-' && after[0] != '`' && after[0] != '_'
}

// fixedHeaderRE is the declaration of a header that stands below the top, with its reason.
var fixedHeaderRE = regexp.MustCompile(`@fixed-header[^\S\n]*:[^\S\n]*\S+`)

// HeaderOffTop says whether the file carries an `@anchors` block below the top that does not
// declare why (`@fixed-header: <why>`) — a header that is not read as one.
func HeaderOffTop(content []byte, cs *config.Comments) bool {
	if cs == nil || AnchorsHeader(content, cs) != nil {
		return false
	}
	lines := strings.Split(string(content), "\n")
	only := cs.CommentOnly(lines)
	for i, l := range lines {
		if only[i] && opensHeader(strings.TrimSpace(l)) {
			return true
		}
	}
	return false
}

// headerEnd is the index just past the header block opened on line i: past the line that
// closes its block comment, or past the run of lines opening with the same line marker.
func headerEnd(lines []string, only []bool, i int, cs *config.Comments) int {
	t := strings.TrimSpace(lines[i])
	for _, m := range cs.Line {
		if strings.HasPrefix(t, m) {
			end := i + 1
			for end < len(lines) && only[end] && strings.HasPrefix(strings.TrimSpace(lines[end]), m) {
				end++
			}
			return end
		}
	}
	// A block: the header ends where the comment holding it closes — the first line from here
	// whose comment's text ends before the line does.
	for j := i; j < len(lines); j++ {
		for _, b := range cs.Block {
			if strings.Contains(lines[j], b.Close) && (j > i || strings.Index(lines[j], b.Close) > strings.Index(lines[j], "@anchors")) {
				return j + 1
			}
		}
		if j > i && !only[j] {
			return j
		}
	}
	return len(lines)
}
