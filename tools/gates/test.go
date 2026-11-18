package main

import (
	"fmt"
	"os"
	"os/exec"
)

// runTest is the unit suite as `anchors test` runs it: the JUnit report proves the
// scenarios, and the lcov report measures the lines. Both are declared in anchors.yaml
// (`tests:`), and a tool that is missing only drops its own report — the suite still runs.
//
// The test run's standard error goes to the terminal and NEVER into the JUnit converter: a
// linker warning there (`ld: warning` on macOS) is not a test result, and piped in with
// the output it made the converter fail with no test failing.
func runTest(_ []string) int {
	if err := os.MkdirAll(".anchors", 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	junit, lcov := findTool("go-junit-report"), findTool("gcov2lcov")

	args := []string{"test", "-coverpkg=./...", "-coverprofile=.anchors/cover.out", "./..."}
	code := 0
	if junit == "" {
		code = passthrough("go", args...)
	} else {
		code = testToJUnit(junit, append([]string{"test", "-v"}, args[1:]...))
	}

	if lcov != "" {
		if c := passthrough(lcov, "-infile", ".anchors/cover.out", "-outfile", ".anchors/lcov.info"); c != 0 && code == 0 {
			code = c
		}
	}
	return code
}

// testToJUnit runs `go test` with its verbose output piped into the JUnit converter, and
// returns the failing exit code of either.
func testToJUnit(junit string, args []string) int {
	out, err := os.Create(".anchors/junit.xml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer out.Close()

	// An OS pipe between the two processes, the shape of `go test | go-junit-report`: each
	// child holds its own end, and the parent closes both once they started.
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	test := exec.Command("go", args...)
	test.Stdout, test.Stderr = w, os.Stderr
	conv := exec.Command(junit)
	conv.Stdin, conv.Stdout, conv.Stderr = r, out, os.Stderr
	if err := conv.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := test.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	_ = w.Close()
	_ = r.Close()
	testCode := exitCode(test.Wait())
	convCode := exitCode(conv.Wait())
	if testCode != 0 {
		return testCode
	}
	return convCode
}
