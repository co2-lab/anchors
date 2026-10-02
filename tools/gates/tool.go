package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// findTool is the path of an external tool: on PATH, or in Go's install folder, where
// `go install` puts it and where a PATH without it still finds it.
func findTool(name string) string {
	if p, err := exec.LookPath(name); err == nil {
		return p
	}
	if home, err := os.UserHomeDir(); err == nil {
		p := filepath.Join(home, "go", "bin", name)
		if exe := p + exeSuffix(); fileExists(exe) {
			return exe
		}
	}
	return ""
}

// skipMissing says a gate did not run because its tool is not installed. It is a skip,
// never a pass: the line says so, and how to install the tool.
func skipMissing(name, hint string) int {
	fmt.Printf("%s is not installed — `%s`. Gate SKIPPED (not approved).\n", name, hint)
	return 0
}

// output runs a tool and returns its standard output alone. Its standard error is kept
// apart, so a warning a tool prints there (a download, the linker's `ld: warning`) never
// reaches what the gate parses.
func output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = nil
	return cmd.Output()
}

// passthrough runs a tool with its output on the gate's own, and returns its exit code.
func passthrough(name string, args ...string) int {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return exitCode(cmd.Run())
}

// exitCode is the exit code a command's error carries: 0 for none, the process's own for
// an exit, and 1 for a command that could not start.
func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	fmt.Fprintln(os.Stderr, err)
	return 1
}

func fileExists(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir()
}
