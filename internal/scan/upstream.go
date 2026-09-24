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
// Measured in blue-eyes: the map gave `anchors-claim.yml` no code, `anchors deliver`
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

// anchorsOpenerRE finds the token that opens the header: `@anchors`, and not
// `@anchors-shared-code` nor a mention in backticks.
var anchorsOpenerRE = regexp.MustCompile(`@anchors(?:\s|-->|$)`)

// AnchorsHeader returns the text of the file's `@anchors` header block, or nil when the file
// has none.
//
// The block opens on the first comment line carrying `@anchors` and ends with the comment:
// at `-->` for an HTML comment, at `*/` for a block comment, and at the first line that is
// not a comment for line comments (`//`, `#`, `*`, `--`). Reading a header key over the
// WHOLE file reads prose and code as declarations — a workflow's jq program with a line
// `parent: (...)` became a node's parent in the map.
func AnchorsHeader(content []byte) []byte {
	lines := bytes.Split(content, []byte("\n"))
	for i, l := range lines {
		t := strings.TrimSpace(string(l))
		if !isHeaderComment(t) || !anchorsOpenerRE.MatchString(t) {
			continue
		}
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
		return bytes.Join(lines[i:end], []byte("\n"))
	}
	return nil
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
