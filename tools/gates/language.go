package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// identifierTest is the test that confronts the language of the code's identifiers.
const identifierTest = "TestCodeLanguage_noProjectIdentifierInPortuguese"

// runLanguage confronts the LANGUAGE of the code's identifiers: everything in this
// repository is written in English, and the one exception is the translation catalog.
//
// It reads DECLARATIONS, through the test that holds the ruler. It refuses a run where the
// test did not run at all: the shell script it replaces named a test that had been
// renamed, `go test` answered "no tests to run" with exit 0, and the gate passed while
// measuring nothing.
func runLanguage(_ []string) int {
	out, err := exec.Command("go", "test", "./internal/gate/", "-run", "^"+identifierTest+"$", "-count=1", "-v").CombinedOutput()
	text := string(out)
	if !regexp.MustCompile(`(?m)^=== RUN\s+` + identifierTest + `\b`).MatchString(text) {
		fmt.Printf("%s did not run — the gate measured nothing. Is the test still named so?\n", identifierTest)
		return 1
	}
	if err != nil {
		n := 0
		for _, l := range strings.Split(text, "\n") {
			if strings.Contains(l, "identifier") && n < 20 {
				fmt.Println(strings.TrimSpace(l))
				n++
			}
		}
		return 1
	}
	fmt.Println("✓ no identifier in Portuguese")
	return 0
}
