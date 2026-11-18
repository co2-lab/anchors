// @anchors
//   code: TSINT
//   ref: TSTLS

// Package testlist reads the project's tests: which tests exist, in which file, on which
// line, and under which title.
//
// How a test is written belongs to the project's language and test library, not to
// Anchors. The project says it in one of two ways, and chooses which:
//
//   - a PATTERN: the regular expression of the call that opens a test, with its title in a
//     `title` group (`\bt\.Run\("(?P<title>[^"]*)"` in Go).
//   - a SCRIPT: a command the project provides, which Anchors runs at the project root and
//     which prints the tests on stdout under the contract below. The script can ask the
//     test library itself, so when the library changes, the answer changes with it.
//
// The script's contract, version 1, is a JSON object and nothing else:
//
//	{"version": 1, "tests": [{"file": "src/a.test.ts", "line": 12, "end": 30, "title": "…"}]}
//
// `file` is relative to the project root, `line` counts from one, `title` is the test's
// title as written. `end`, optional, is the test's last line: the script knows where a
// test ends, the reading of a pattern only estimates it. Any other field, a missing one,
// or another version is outside the contract and refused, naming what is wrong.
package testlist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/co2-lab/anchors/internal/shell"
)

// ContractVersion is the version of the script's output this package reads.
const ContractVersion = 1

// Test is one test of the project.
type Test struct {
	File string `json:"file"`
	Line int    `json:"line"`
	// End is the test's last line, when the source knows it; zero when it does not.
	End   int    `json:"end,omitempty"`
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
	return ListFrom(root, files, src, nil)
}

// ListFrom is List reading the files a pattern scans through `read` (a path relative to
// root) — the git index, under `--index` —, so a test's line is a line of the content the
// gates confront. Nil reads the tree. A script reads the files on its own.
func ListFrom(root string, files []string, src Source, read func(rel string) ([]byte, error)) ([]Test, error) {
	switch {
	case src.Pattern != "" && src.Script != "":
		return nil, fmt.Errorf("the tests are declared both by a pattern and by a script; declare one")
	case src.Script != "":
		return runScript(root, src.Script)
	case src.Pattern != "":
		return scanFiles(root, files, src.Pattern, read)
	}
	return nil, nil
}

// runScript runs the project's script at the root and reads its output under the contract.
func runScript(root, script string) ([]Test, error) {
	cmd, err := shell.Command(script) // the command is declared by the project, as a gate's `run:`
	if err != nil {
		return nil, fmt.Errorf("the tests script `%s` cannot run: %w", script, err)
	}
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
		case filepath.IsAbs(t.File) || path.IsAbs(filepath.ToSlash(t.File)) || strings.HasPrefix(path.Clean(filepath.ToSlash(t.File)), "../"):
			return nil, fmt.Errorf("tests[%d]: `file` %q is not relative to the project root", i, t.File)
		case t.Line < 1:
			return nil, fmt.Errorf("tests[%d]: `line` must count from 1", i)
		case t.End != 0 && t.End < t.Line:
			return nil, fmt.Errorf("tests[%d]: `end` is before `line`", i)
		case strings.TrimSpace(t.Title) == "":
			return nil, fmt.Errorf("tests[%d]: `title` is empty", i)
		}
		(*out.Tests)[i].File = path.Clean(filepath.ToSlash(t.File))
	}
	return *out.Tests, nil
}

// scanFiles finds, in each file, every match of the pattern, and takes as the test's title
// what its `title` group captured — the first group whose name starts with `title` and took
// part (`title`, `title2`… for the quotes a language writes a title in). How a title is
// quoted is the language's, so the pattern says it: the core reads no string literal. A match
// whose title group captured nothing — a variable, a template built elsewhere — is not a
// test this reading can name, and is left out; a pattern with no title group is an error.
func scanFiles(root string, files []string, pattern string, read func(rel string) ([]byte, error)) ([]Test, error) {
	if read == nil {
		read = func(rel string) ([]byte, error) { return os.ReadFile(filepath.Join(root, rel)) }
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("the tests pattern does not compile: %w", err)
	}
	var groups []int
	for i, name := range re.SubexpNames() {
		if strings.HasPrefix(name, "title") {
			groups = append(groups, i)
		}
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("the tests pattern has no `title` group: name where the title stands, e.g. `'(?P<title>[^']*)'`")
	}
	var out []Test
	sorted := append([]string(nil), files...)
	sort.Strings(sorted)
	for _, f := range sorted {
		b, err := read(f)
		if err != nil {
			continue // @resilient: a test file that cannot be read has no test this reading can name; the map notices a missing file
		}
		src := string(b)
		for _, m := range re.FindAllStringSubmatchIndex(src, -1) {
			title := ""
			for _, gi := range groups {
				if m[2*gi] >= 0 {
					title = src[m[2*gi]:m[2*gi+1]]
					break
				}
			}
			if title == "" {
				continue
			}
			out = append(out, Test{File: filepath.ToSlash(f), Line: strings.Count(src[:m[0]], "\n") + 1, Title: title})
		}
	}
	return out, nil
}

// lastLine is the last non-empty line of a command's output: where it says why it stopped.
func lastLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}
