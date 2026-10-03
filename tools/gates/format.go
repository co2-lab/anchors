package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// runFormat confronts the Go files with gofmt's canonical form, over the whole tree, as the
// CI does.
//
// The CI checked it and the commit did not: a test file edited after the last local
// gofmt run went into a commit, and only the CI's format step caught it — on every platform
// at once, after the push. As a gate it runs in the pre-commit, where the fix is one
// command.
//
// It reads the files git tracks, as the CI's checkout has them: a copy of the tree another
// agent works in (`.claude/worktrees/…`) is no file of this repository.
func runFormat(_ []string) int {
	tracked, err := exec.Command("git", "ls-files", "--", "*.go").Output()
	if err != nil {
		fmt.Printf("git could not list the Go files: %v\n", err)
		return 1
	}
	files := strings.Fields(string(tracked))
	if len(files) == 0 {
		return reportFormat("")
	}
	// In batches: Windows caps a command line near 32k characters, and a repository's
	// paths pass it.
	var listed strings.Builder
	for len(files) > 0 {
		n := min(len(files), 200)
		out, err := exec.Command("gofmt", append([]string{"-l"}, files[:n]...)...).Output()
		if err != nil {
			fmt.Printf("gofmt could not run: %v\n", err)
			return 1
		}
		listed.Write(out)
		files = files[n:]
	}
	return reportFormat(listed.String())
}

// reportFormat says which files are out of gofmt's form, and fails when any is.
func reportFormat(listed string) int {
	var files []string
	for _, l := range strings.Split(strings.TrimSpace(listed), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			files = append(files, l)
		}
	}
	if len(files) == 0 {
		fmt.Println("✓ every Go file is in gofmt's form")
		return 0
	}
	fmt.Printf("%d file(s) out of gofmt's form — run `gofmt -w .`:\n  %s\n", len(files), strings.Join(files, "\n  "))
	return 1
}
