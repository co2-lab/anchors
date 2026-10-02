// @anchors
//   ref: UPOWP

package scan

import (
	"bytes"
	"regexp"
	"strings"
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

// anchorsOpenerRE finds the token that opens the header: `@anchors` as the FIRST word of the
// comment, and not `@anchors-shared-code`, a mention in backticks, nor prose that names
// the header — a package comment saying "the header @anchors (code:/ref:)" was read as the
// header, and a header fix wrote the unit's `ref:` into the middle of that prose.
var anchorsOpenerRE = regexp.MustCompile(`^(?://|#|/\*+|\*|--|<!--)\s*@anchors(?:\s|-->|$)`)

// AnchorsHeader returns the text of the file's `@anchors` header block, or nil when the file
// has none.
//
// The block opens on the first comment line carrying `@anchors` and ends with the comment:
// at `-->` for an HTML comment, at `*/` for a block comment, and at the first line that is
// not a comment for line comments (`//`, `#`, `*`, `--`). Reading a header key over the
// WHOLE file reads prose and code as declarations — a workflow's jq program with a line
// `parent: (...)` became a node's parent in the map.
//
// The header is at the TOP: before it, only blank lines, comments and a shebang. A block
// further down is text — an example of a header in a guide, a string in the code — and is
// not the file's header: `guide_header.go` carried one inside a string, and the map read
// `layer: screen` as its unit's layer. A file whose header must stand lower (a directive
// the language wants first) says so inside the block, with the reason:
// `@fixed-header: <why>`.
func AnchorsHeader(content []byte) []byte {
	lines := bytes.Split(content, []byte("\n"))
	codeBefore, inBlock := false, ""
	for i, l := range lines {
		t := strings.TrimSpace(string(l))
		if isHeaderComment(t) && anchorsOpenerRE.MatchString(t) {
			end := headerEnd(lines, i, t)
			block := bytes.Join(lines[i:end], []byte("\n"))
			if !codeBefore || fixedHeaderRE.Match(block) {
				return block
			}
			continue
		}
		switch {
		case inBlock != "":
			if strings.Contains(t, inBlock) {
				inBlock = ""
			}
		case strings.HasPrefix(t, "/*") && !strings.Contains(t, "*/"):
			inBlock = "*/"
		case strings.HasPrefix(t, "<!--") && !strings.Contains(t, "-->"):
			inBlock = "-->"
		case t == "" || isHeaderComment(t) || strings.HasPrefix(t, "#!"):
		default:
			codeBefore = true
		}
	}
	return nil
}

// fixedHeaderRE is the declaration of a header that stands below the top, with its reason.
var fixedHeaderRE = regexp.MustCompile(`@fixed-header[^\S\n]*:[^\S\n]*\S+`)

// HeaderOffTop says whether the file carries an `@anchors` block below the top that does not
// declare why (`@fixed-header: <why>`) — a header that is not read as one.
func HeaderOffTop(content []byte) bool {
	if AnchorsHeader(content) != nil {
		return false
	}
	for _, l := range bytes.Split(content, []byte("\n")) {
		t := strings.TrimSpace(string(l))
		if isHeaderComment(t) && anchorsOpenerRE.MatchString(t) {
			return true
		}
	}
	return false
}

// headerEnd is the index just past the header block opened on line i.
func headerEnd(lines [][]byte, i int, t string) int {
	end := i + 1
	switch {
	case strings.HasPrefix(t, "<!--"):
		if !strings.Contains(t[strings.Index(t, "@anchors"):], "-->") {
			end = closingLine(lines, i+1, "-->")
		}
	case strings.HasPrefix(t, "/*"):
		if !strings.Contains(t, "*/") {
			end = closingLine(lines, i+1, "*/")
		}
	default:
		for end < len(lines) && isHeaderComment(strings.TrimSpace(string(lines[end]))) {
			end++
		}
	}
	return end
}

// closingLine returns the index just past the line that carries `close`, or the end of the
// file when the comment never closes.
func closingLine(lines [][]byte, from int, close string) int {
	for j := from; j < len(lines); j++ {
		if bytes.Contains(lines[j], []byte(close)) {
			return j + 1
		}
	}
	return len(lines)
}

func isHeaderComment(t string) bool {
	for _, p := range []string{"//", "#", "/*", "*", "--", "<!--"} {
		if strings.HasPrefix(t, p) {
			return true
		}
	}
	return false
}
