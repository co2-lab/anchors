// @anchors
//   code: TMFLT
//   ref: TFSTM

package testsig

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// DefaultTimeoutPattern tells, in a failure or in the line of the test where it failed, a
// test that ran out of time (DESIGN-retry-on-load.md): a test or a hook past its timeout
// (jest's "Exceeded timeout", vitest's "timed out", go test's "test timed out"), and a wait
// that ran out — jest's `waitFor` rethrows its last assertion and never says "timeout", so
// only the line `await waitFor(...)` tells it (measured at MIF).
const DefaultTimeoutPattern = `Exceeded timeout|[Tt]imed out|deadline exceeded|\bwaitFor\(|\bfind(All)?By[A-Z]`

// stackLineRE is a file and a line in a stack: `(apps/x/a.test.ts:66:18)` or `a_test.go:12:`.
var stackLineRE = regexp.MustCompile(`([\w./@-]+\.[A-Za-z]+):(\d+)`)

// TimedOut says whether a failed case ran out of time: the pattern matches its failure, or a
// line of its own test file the failure's stack points at, read under root.
func TimedOut(c CaseResult, pattern *regexp.Regexp, root string) bool {
	if !c.Failed {
		return false
	}
	if pattern.MatchString(c.Failure) {
		return true
	}
	if c.File == "" {
		return false
	}
	own := filepath.ToSlash(c.File)
	for _, m := range stackLineRE.FindAllStringSubmatch(c.Failure, -1) {
		if !strings.HasSuffix(own, filepath.ToSlash(m[1])) && !strings.HasSuffix(filepath.ToSlash(m[1]), own) {
			continue
		}
		n, _ := strconv.Atoi(m[2])
		if line := sourceLine(filepath.Join(root, filepath.FromSlash(own)), n); line != "" && pattern.MatchString(line) {
			return true
		}
	}
	return false
}

// sourceLine is the n-th line of a file, or empty.
func sourceLine(path string, n int) string {
	b, err := os.ReadFile(path)
	if err != nil || n < 1 {
		return ""
	}
	lines := strings.Split(string(b), "\n")
	if n > len(lines) {
		return ""
	}
	return lines[n-1]
}
