// Command gates runs this repository's external gates: the steps `anchors.yaml` declares
// with `run: "go run ./tools/gates <gate>"`.
//
// WHY GO, AND NOT SHELL.
//
// These were bash scripts, and they broke on what a script cannot see: the platform. BSD
// sed does not know `\s`, so a trailing-space trim did nothing on macOS; `2>&1` mixed a
// tool's download messages and the linker's `ld: warning` into the findings; and Windows
// has no `sh` unless git brought one. Each script also leaned on jq and python3 to read the
// tools' JSON. Written in the project's own language, a gate runs the same on macOS, Linux
// and Windows with nothing but the toolchain the project already needs, and it is vetted
// and tested like the rest of the code (reported from baas-proxy).
//
// The external tools stay what they were — gitleaks, osv-scanner, syft, semgrep,
// go-licenses, gremlins, go-junit-report, gcov2lcov — and a missing one skips its gate
// saying so, never passing it.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// gates are the subcommands, one per gate step.
var gates = map[string]func(args []string) int{
	"test":      runTest,
	"mutation":  runMutation,
	"secrets":   runSecrets,
	"deps-vuln": runDepsVuln,
	"licenses":  runLicenses,
	"sbom":      runSBOM,
	"semgrep":   runSemgrep,
	"language":  runLanguage,
	"format":    runFormat,
}

func main() {
	if len(os.Args) < 2 || gates[os.Args[1]] == nil {
		names := make([]string, 0, len(gates))
		for n := range gates {
			names = append(names, n)
		}
		sort.Strings(names)
		fmt.Fprintf(os.Stderr, "usage: go run ./tools/gates <%s> [args]\n", strings.Join(names, "|"))
		os.Exit(2)
	}
	// Every gate reads paths from the repository root, wherever it was called from.
	if out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output(); err == nil {
		if err := os.Chdir(strings.TrimSpace(string(out))); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	os.Exit(gates[os.Args[1]](os.Args[2:]))
}
