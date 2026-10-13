// @anchors
//   code: TFTTM
//   ref: TFSTM

package testsig

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

// The failures below are those MIF recorded on 2026-10-13, under a load of 13 to 77 on 10 CPUs.
const (
	hookTimeout = "Error: thrown: \"Exceeded timeout of 20000 ms for a hook.\n" +
		"Add a timeout value to this test to increase the timeout, if this is a long-running test.\"\n" +
		"    at Object.afterEach (node_modules/@testing-library/react-native/src/index.ts:15:5)"
	waitForRethrow = "Error: expect(received).toBe(expected) // Object.is equality\n" +
		"Expected: true\nReceived: false\n" +
		"    at Object.<anonymous> (apps/mobile/src/features/ir/hooks/useIRSnapshot.test.ts:3:18)"
)

func timeoutProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "apps", "mobile", "src", "features", "ir", "hooks")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := "it('busca', async () => {\n  render()\n  await waitFor(() => expect(result.current.isSuccess).toBe(true))\n  expect(a).toBe(b)\n})\n"
	if err := os.WriteFile(filepath.Join(dir, "useIRSnapshot.test.ts"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

var defaultPattern = regexp.MustCompile(DefaultTimeoutPattern)

func TestTimedOut_message(t *testing.T) {
	t.Run("TFSTM-B01: A failure past its timeout ran out of time", func(t *testing.T) {})
	c := CaseResult{Name: "IPSBI-V06#01", File: "apps/mobile/src/components/molecules/ImportStatusBanner.test.tsx", Failed: true, Failure: hookTimeout}
	if !TimedOut(c, defaultPattern, t.TempDir()) {
		t.Error("a hook past its timeout ran out of time")
	}
}

func TestTimedOut_waitLine(t *testing.T) {
	t.Run("TFSTM-B02: A wait that ran out is told by its line", func(t *testing.T) {})
	root := timeoutProject(t)
	c := CaseResult{Name: "UIRNU-B01#01", File: "apps/mobile/src/features/ir/hooks/useIRSnapshot.test.ts", Failed: true, Failure: waitForRethrow}
	if !TimedOut(c, defaultPattern, root) {
		t.Error("the rethrown assertion of a waitFor is a wait that ran out, told by its line")
	}
}

func TestTimedOut_assertion(t *testing.T) {
	t.Run("TFSTM-B03: A failed assertion is not a timeout", func(t *testing.T) {})
	root := timeoutProject(t)
	file := "apps/mobile/src/features/ir/hooks/useIRSnapshot.test.ts"
	outside := CaseResult{File: file, Failed: true, Failure: "Error: expect(a).toBe(b)\n    at Object.<anonymous> (" + file + ":4:3)"}
	other := CaseResult{File: "apps/mobile/src/other.test.ts", Failed: true, Failure: waitForRethrow}
	passed := CaseResult{File: file, Failure: hookTimeout}
	for name, c := range map[string]CaseResult{"outside a wait": outside, "another file's line": other, "passed": passed} {
		if TimedOut(c, defaultPattern, root) {
			t.Errorf("%s: not a timeout", name)
		}
	}
}

func TestTimedOut_unreadableFile(t *testing.T) {
	t.Run("TFSTM-B05: An unreadable test file counts for nothing", func(t *testing.T) {})
	c := CaseResult{File: "apps/mobile/src/features/ir/hooks/useIRSnapshot.test.ts", Failed: true, Failure: waitForRethrow}
	if TimedOut(c, defaultPattern, t.TempDir()) {
		t.Error("with no file to read, the failure alone decides")
	}
}

func TestRetried(t *testing.T) {
	t.Run("TFSTM-B04: The tests that passed on a retry are kept", func(t *testing.T) {})
	root := t.TempDir()
	if got := LoadRetried(root); len(got) != 0 {
		t.Errorf("no record, none: %v", got)
	}
	at := time.Now().Round(time.Second)
	if err := AddRetried(root, []Retried{{Name: "a", At: at, Load: 31, CPUs: 10}, {Name: "b", At: at}}); err != nil {
		t.Fatal(err)
	}
	if err := AddRetried(root, []Retried{{Name: "c", At: at}}); err != nil {
		t.Fatal(err)
	}
	got := LoadRetried(root)
	if len(got) != 3 || got[0].Name != "a" || got[0].Load != 31 || got[2].Name != "c" || !got[0].At.Equal(at) {
		t.Errorf("kept in order: %+v", got)
	}
	many := make([]Retried, retriedKeep)
	for i := range many {
		many[i] = Retried{Name: "n", At: at}
	}
	many[len(many)-1].Name = "last"
	if err := AddRetried(root, many); err != nil {
		t.Fatal(err)
	}
	if got := LoadRetried(root); len(got) != retriedKeep || got[0].Name != "n" || got[len(got)-1].Name != "last" {
		t.Errorf("only the latest %d kept: %d", retriedKeep, len(got))
	}
	if err := AddRetried(root, nil); err != nil {
		t.Errorf("nothing to add writes nothing: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(RetriedFile)), []byte("{broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := LoadRetried(root); len(got) != 0 {
		t.Errorf("an unreadable record reads none: %v", got)
	}
}
