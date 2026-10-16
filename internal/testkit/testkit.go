// Package testkit holds what the tests of every package share to run the same on Linux,
// macOS and Windows. It is imported only by tests.
package testkit

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
	"github.com/co2-lab/anchors/internal/shell"
)

// CaptureStdout returns what fn printed on standard output.
//
// The pipe is drained WHILE fn runs. Reading it only afterwards deadlocked on Windows,
// whose pipe buffer is a few kilobytes: a command printing a long prompt blocked on the
// write, and the test sat there until the ten-minute timeout.
func CaptureStdout(t *testing.T, fn func()) string {
	t.Helper()
	return capture(t, &os.Stdout, fn)
}

// CaptureStderr returns what fn printed on standard error.
func CaptureStderr(t *testing.T, fn func()) string {
	t.Helper()
	return capture(t, &os.Stderr, fn)
}

func capture(t *testing.T, std **os.File, fn func()) string {
	t.Helper()
	orig := *std
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()
	*std = w
	defer func() { *std = orig }()
	fn()
	w.Close()
	*std = orig
	return <-done
}

// FakeBin writes an executable named `name` into dir that runs the POSIX shell `script`,
// and returns its path.
//
// On Windows a script with no extension does not run, and `exec.LookPath` would find the
// real tool first. A `.cmd` wrapper was tried and lost arguments: `cmd.exe` reads `&`, `<`,
// `>` and `|` inside them — a GitHub URL with `?a=1&b=2` became two commands. So on Windows
// the fake is `name.exe`, a small launcher built once per test run, which hands its exact
// argv to the shell Git brings, running the script beside it.
func FakeBin(t *testing.T, dir, name, script string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	body := script
	if !strings.HasPrefix(body, "#!") {
		body = "#!/bin/sh\n" + body
	}
	if err := os.WriteFile(p, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		sh, err := shell.Path()
		if err != nil {
			t.Skipf("no POSIX shell to run the fake %s: %v", name, err)
		}
		exe, err := launcher()
		if err != nil {
			t.Fatalf("build the fake launcher: %v", err)
		}
		b, err := os.ReadFile(exe)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p+".exe", b, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv("ANCHORS_TESTKIT_SH", sh)
	}
	return p
}

// launcherSource runs the script that sits beside it, with the same name and no `.exe`,
// through the shell named in ANCHORS_TESTKIT_SH, passing every argument as it came.
const launcherSource = `package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	self, err := os.Executable()
	if err != nil {
		os.Exit(127)
	}
	script := strings.TrimSuffix(self, ".exe")
	cmd := exec.Command(os.Getenv("ANCHORS_TESTKIT_SH"), append([]string{script}, os.Args[1:]...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(127)
	}
}
`

var (
	launcherOnce sync.Once
	launcherPath string
	launcherErr  error
)

// launcher builds the fake launcher once per test binary.
func launcher() (string, error) {
	launcherOnce.Do(func() {
		dir, err := os.MkdirTemp("", "anchors-testkit-")
		if err != nil {
			launcherErr = err
			return
		}
		src := filepath.Join(dir, "main.go")
		if err := os.WriteFile(src, []byte(launcherSource), 0o644); err != nil {
			launcherErr = err
			return
		}
		out := filepath.Join(dir, "launcher.exe")
		cmd := exec.Command("go", "build", "-o", out, src)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GO111MODULE=off")
		if b, err := cmd.CombinedOutput(); err != nil {
			launcherErr = fmt.Errorf("%v: %s", err, b)
			return
		}
		launcherPath = out
	})
	return launcherPath, launcherErr
}

// OnPath puts dir first on PATH for the rest of the test.
func OnPath(t *testing.T, dir string) {
	t.Helper()
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// SkipWithoutPOSIXPermissions skips a test that needs a file the process cannot read or
// write, or a file mode kept across a rewrite: Windows does not enforce POSIX permissions
// (`chmod` only toggles read-only), so there is nothing there to observe.
func SkipWithoutPOSIXPermissions(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not enforce POSIX file permissions")
	}
}

// AllComments is a configuration that declares every family's comment syntax
// (config.AllComments): what a test reads a file's comments with, as a project that
// declared them all.
func AllComments() *config.Config {
	return &config.Config{Dialect: &config.Dialect{Comments: config.AllComments()}}
}

// Comments is the syntax a test reads a file at rel with (AllComments).
func Comments(rel string) *config.Comments { return AllComments().CommentsFor(rel) }

// WithComments gives a test's configuration every family's comment syntax, and AnyComment for
// any other file, when it declares none, and returns it.
func WithComments(cfg *config.Config) *config.Config {
	if cfg == nil {
		cfg = &config.Config{}
	}
	if cfg.Dialect == nil {
		cfg.Dialect = &config.Dialect{}
	}
	if len(cfg.Dialect.Comments) == 0 {
		// Every family's syntax, and last any file's (AnyComment): a test node is often
		// given no path.
		cfg.Dialect.Comments = append(config.AllComments(), *AnyComment())
	}
	FillFamily(cfg.Dialect)
	return cfg
}

// FillFamily gives a test's dialect its family's example lexicon (config.FamilyLexicon) in
// the fields it leaves empty — as a project that wrote the example into its anchors.yaml.
// The core applies no family by itself; a test written for a family's reading declares it.
func FillFamily(d *config.Dialect) {
	if d == nil {
		return
	}
	ex := config.FamilyLexicon(d.Family)
	if ex == nil {
		return
	}
	fill := func(dst *string, src string) {
		if *dst == "" {
			*dst = src
		}
	}
	fill(&d.ExportedFunc, ex.ExportedFunc)
	fill(&d.ParamName, ex.ParamName)
	fill(&d.Loop, ex.Loop)
	fill(&d.Cursor, ex.Cursor)
	fill(&d.HTTPStatus, ex.HTTPStatus)
	fill(&d.HTTPStatusDynamic, ex.HTTPStatusDynamic)
	fill(&d.EnvRead, ex.EnvRead)
	fill(&d.Definition, ex.Definition)
	fill(&d.SetPromise, ex.SetPromise)
	fill(&d.SetSlice, ex.SetSlice)
	fill(&d.StatementEnd, ex.StatementEnd)
	fill(&d.SignatureEnd, ex.SignatureEnd)
	if len(d.HandlePatterns) == 0 {
		d.HandlePatterns = ex.HandlePatterns
	}
	if len(d.LogPatterns) == 0 {
		d.LogPatterns = ex.LogPatterns
	}
	if len(d.GuardPatterns) == 0 {
		d.GuardPatterns = ex.GuardPatterns
	}
	if len(d.Limits) == 0 {
		d.Limits = ex.Limits
	}
	switch {
	case d.Tests == nil:
		d.Tests = ex.Tests
	case ex.Tests != nil && d.Tests.Pattern == "" && d.Tests.Script == "":
		t := *ex.Tests
		if d.Tests.Assertion != "" {
			t.Assertion = d.Tests.Assertion
		}
		d.Tests = &t
	}
}

// CommentsYAML is the `comments:` block of every family's syntax, to put under `dialect:`
// in a test's anchors.yaml.
const CommentsYAML = `  comments:
    - files: '**/*.{ts,tsx,js,jsx,mjs,cjs,mts,cts}'
      line: ['//']
      block: [['{/*', '*/}'], ['/*', '*/']]
      strings: ['"', "'", '` + "`" + `']
      multiline: ['` + "`" + `']
      no_line_after: '>\s*$'
    - files: '**/*.go'
      line: ['//']
      block: [['/*', '*/']]
      strings: ['"', "'", '` + "`" + `']
      multiline: ['` + "`" + `']
    - files: '**/*.py'
      line: ['#']
      strings: ['"""', "'''", '"', "'"]
      multiline: ['"""', "'''"]
    - files: '**/*.{yaml,yml,toml,sh,bash,zsh,rb,r,R}'
      line: ['#']
      strings: ['"', "'"]
    - files: '**/*.{java,kt,kts,cs,scala,swift,rs,c,h,cpp,hpp,cc}'
      line: ['//']
      block: [['/*', '*/']]
      strings: ['"', "'"]
    - files: '**/*.php'
      line: ['//', '#']
      block: [['/*', '*/']]
      strings: ['"', "'"]
    - files: '**/*.{ml,mli,mll}'
      block: [['(*', '*)']]
      strings: ['"']
    - files: '**/*.mly'
      block: [['(*', '*)'], ['/*', '*/']]
      strings: ['"']
    - files: '**/{dune,dune-project,dune-file}'
      line: [';']
      strings: ['"']
    - files: '**/*.{sql,lua,hs}'
      line: ['--']
      strings: ["'", '"']
    - files: '**/*.{html,xml,svg,vue}'
      block: [['<!--', '-->']]
`

// AnyComment is a syntax with every marker the families use — `//`, `#`, `--`, `;`, and the
// blocks `/* */`, `<!-- -->`, `(* *)` — and no strings: for a test that reads a header or a
// flag in any language at once.
func AnyComment() *config.Comments {
	return &config.Comments{
		Line:  []string{"//", "#", "--", ";"},
		Block: []config.Delims{{Open: "/*", Close: "*/"}, {Open: "<!--", Close: "-->"}, {Open: "(*", Close: "*)"}},
	}
}

// AnyConfig is a configuration whose one comment syntax (AnyComment) covers every file: for
// a test that checks a file whose path it does not name.
func AnyConfig() *config.Config {
	return &config.Config{Dialect: &config.Dialect{Comments: []config.Comments{*AnyComment()}}}
}

// Chain is a config whose dialect declares the family's chain patterns and the comments of
// every family: what a test of a gate that reads imports runs on.
func Chain(family string) *config.Config {
	d := config.FamilyChain(family)
	if d == nil {
		d = &config.Dialect{}
	}
	d.Family = family
	return WithComments(&config.Config{Dialect: d})
}
