// Package testlist reads the project's tests: which tests exist, in which file, on which
// line, and under which title.
//
// How a test is written belongs to the project's language and test library, not to
// Anchors. The project says it in one of two ways, and chooses which:
//
//   - a PATTERN: the regular expression of the call that opens a test, up to its opening
//     parenthesis (`\bt\.Run\(` in Go, `\b(?:it|test)\(` in Jest). The title is the string
//     literal right after it.
//   - a SCRIPT: a command the project provides, which Anchors runs at the project root and
//     which prints the tests on stdout under the contract below. The script can ask the
//     test library itself, so when the library changes, the answer changes with it.
//
// The script's contract, version 1, is a JSON object and nothing else:
//
//	{"version": 1, "tests": [{"file": "src/a.test.ts", "line": 12, "title": "…"}]}
//
// `file` is relative to the project root, `line` counts from one, `title` is the test's
// title as written. Any other field, a missing one, or another version is outside the
// contract and refused, naming what is wrong.
package testlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ContractVersion is the version of the script's output this package reads.
const ContractVersion = 1

// Test is one test of the project.
type Test struct {
	File  string `json:"file"`
	Line  int    `json:"line"`
	Title string `json:"title"`
}

// Source is how the project says its tests are written: a pattern or a script, never both.
type Source struct {
	Pattern string
	Script  string
}

// Declared reports whether the source says anything.
func (s Source) Declared() bool { return s.Pattern != "" || s.Script != "" }

// List reads the tests of the given test files (paths relative to root) through the
// source. A script lists the whole project on its own and the files are not passed to it.
func List(root string, files []string, src Source) ([]Test, error) {
	switch {
	case src.Pattern != "" && src.Script != "":
		return nil, fmt.Errorf("the tests are declared both by a pattern and by a script; declare one")
	case src.Script != "":
		return runScript(root, src.Script)
	case src.Pattern != "":
		return scanFiles(root, files, src.Pattern)
	}
	return nil, nil
}

// runScript runs the project's script at the root and reads its output under the contract.
func runScript(root, script string) ([]Test, error) {
	cmd := exec.Command("sh", "-c", script) //nolint:gosec // the command is declared by the project, as a gate's `run:`
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("the tests script `%s` failed (%v): %s", script, err, lastLine(stderr.String()))
	}
	tests, err := Parse(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf("the tests script `%s` answered outside the contract: %w", script, err)
	}
	return tests, nil
}

// Parse reads the script's output under the contract, strictly.
func Parse(b []byte) ([]Test, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var out struct {
		Version *int    `json:"version"`
		Tests   *[]Test `json:"tests"`
	}
	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("not the contract's JSON: %w", err)
	}
	if dec.More() {
		return nil, fmt.Errorf("more than one JSON value on stdout")
	}
	if out.Version == nil || *out.Version != ContractVersion {
		return nil, fmt.Errorf("`version` must be %d", ContractVersion)
	}
	if out.Tests == nil {
		return nil, fmt.Errorf("`tests` is missing")
	}
	for i, t := range *out.Tests {
		switch {
		case t.File == "":
			return nil, fmt.Errorf("tests[%d]: `file` is empty", i)
		case filepath.IsAbs(t.File) || strings.HasPrefix(path.Clean(filepath.ToSlash(t.File)), "../"):
			return nil, fmt.Errorf("tests[%d]: `file` %q is not relative to the project root", i, t.File)
		case t.Line < 1:
			return nil, fmt.Errorf("tests[%d]: `line` must count from 1", i)
		case strings.TrimSpace(t.Title) == "":
			return nil, fmt.Errorf("tests[%d]: `title` is empty", i)
		}
		(*out.Tests)[i].File = path.Clean(filepath.ToSlash(t.File))
	}
	return *out.Tests, nil
}

// scanFiles finds, in each file, every call the pattern opens followed by a string literal,
// and takes the literal as the test's title. A call whose title is not a literal (a
// variable, a template with placeholders built elsewhere) is not a test this reading can
// name, and is left out.
func scanFiles(root string, files []string, pattern string) ([]Test, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("the tests pattern does not compile: %w", err)
	}
	var out []Test
	sorted := append([]string(nil), files...)
	sort.Strings(sorted)
	for _, f := range sorted {
		b, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			continue // @resilient: a test file that cannot be read has no test this reading can name; the map notices a missing file
		}
		src := string(b)
		for _, loc := range re.FindAllStringIndex(src, -1) {
			title, ok := literalAt(src, loc[1])
			if !ok {
				continue
			}
			out = append(out, Test{File: filepath.ToSlash(f), Line: strings.Count(src[:loc[0]], "\n") + 1, Title: title})
		}
	}
	return out, nil
}

// literalAt reads the string literal that starts at s[i], after blanks: quoted by ', " or
// `, with backslash escapes inside the first two.
func literalAt(s string, i int) (string, bool) {
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	if i >= len(s) {
		return "", false
	}
	q := s[i]
	if q != '\'' && q != '"' && q != '`' {
		return "", false
	}
	var b strings.Builder
	for j := i + 1; j < len(s); j++ {
		c := s[j]
		switch {
		case c == q:
			return b.String(), true
		case c == '\\' && q != '`' && j+1 < len(s):
			j++
			b.WriteByte(s[j])
		case c == '\n' && q != '`':
			return "", false
		default:
			b.WriteByte(c)
		}
	}
	return "", false
}

// lastLine is the last non-empty line of a command's output: where it says why it stopped.
func lastLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}
