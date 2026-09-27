package flow

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The message is read from the INTEGRATION BRANCH on the platform, raw — not from the
// agent's checkout, which can be days behind.
func TestNotificationsFromBranch_readsTheIntegrationBranch(t *testing.T) {
	t.Run("NTFCT-B03: Github mode reads the file raw from the integration branch", func(t *testing.T) {})
	t.Run("NTFCT-X01: Github mode does not read the agent's checkout", func(t *testing.T) {})
	var got []string
	orig := ghRaw
	defer func() { ghRaw = orig }()
	ghRaw = func(args ...string) ([]byte, error) {
		got = args
		return []byte("update the binary to 0.1.179\n"), nil
	}

	text, err := notificationsFromBranch("o/r", "develop")
	if err != nil || text != "update the binary to 0.1.179\n" {
		t.Fatalf("text=%q err=%v", text, err)
	}
	call := strings.Join(got, " ")
	for _, want := range []string{"api", "repos/o/r/contents/notifications.md?ref=develop", "application/vnd.github.raw"} {
		if !strings.Contains(call, want) {
			t.Errorf("the read should carry %q: gh %s", want, call)
		}
	}

	if _, err := notificationsFromBranch("o/r", "feat/x y"); err != nil {
		t.Fatal(err)
	}
	if call := strings.Join(got, " "); !strings.Contains(call, "notifications.md?ref=feat%2Fx+y") {
		t.Errorf("the branch must be escaped in the address: gh %s", call)
	}
	if _, err := notificationsFromBranch("o/r", ""); err != nil {
		t.Fatal(err)
	}
	if call := strings.Join(got, " "); strings.Contains(call, "?ref=") {
		t.Errorf("without a branch no ref is sent: gh %s", call)
	}
}

// Most of the time there is nothing to say: a missing file is silence, not an error. Any
// other failure is reported, so a broken read does not look like "no message".
func TestNotificationsFromBranch_missingIsSilentOtherErrorsAreNot(t *testing.T) {
	t.Run("NTFCT-B04: A file missing on the platform is silence", func(t *testing.T) {})
	t.Run("NTFCT-E01: A failed read is reported as one line", func(t *testing.T) {})
	orig := ghRaw
	defer func() { ghRaw = orig }()

	ghRaw = func(...string) ([]byte, error) {
		return []byte("gh: Not Found (HTTP 404)"), errors.New("exit status 1")
	}
	if text, err := notificationsFromBranch("o/r", "develop"); err != nil || text != "" {
		t.Errorf("a missing file must be silence: text=%q err=%v", text, err)
	}

	ghRaw = func(...string) ([]byte, error) {
		return []byte("gh: API rate limit exceeded (HTTP 403)"), errors.New("exit status 1")
	}
	if _, err := notificationsFromBranch("o/r", "develop"); err == nil {
		t.Error("a failed read must be reported, not taken for an empty file")
	}
}

// The file may explain itself in HTML comments and still count as empty: an explanation
// printed on every `next` would teach agents to skip the block.
func TestPrintNotifications(t *testing.T) {
	t.Run("NTFCT-B02: A file with only comments or nothing prints nothing", func(t *testing.T) {})
	t.Run("NTFCT-B05: The block names the file and its source and indents the message", func(t *testing.T) {})
	t.Run("NTFCT-I01: The explanatory comment is never printed", func(t *testing.T) {})
	explained := "<!-- Write here what every agent must read on `anchors next`. -->\n\n"

	if out := stdoutOf(t, func() { printNotifications(explained, "o/r@develop", nil) }); out != "" {
		t.Errorf("a file with only comments must print nothing, printed:\n%s", out)
	}
	if out := stdoutOf(t, func() { printNotifications("", "o/r@develop", nil) }); out != "" {
		t.Errorf("an empty file must print nothing, printed:\n%s", out)
	}

	out := stdoutOf(t, func() {
		printNotifications(explained+"Update anchors to 0.1.179.\nBugs go with `escalate --bug`.\n", "o/r@develop", nil)
	})
	for _, want := range []string{"NOTIFICATIONS", "notifications.md", "o/r@develop", "\n  Update anchors to 0.1.179.\n", "\n  Bugs go with `escalate --bug`.\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("the notice should show %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Write here") {
		t.Errorf("the explanatory comment must not be printed:\n%s", out)
	}

	out = stdoutOf(t, func() { printNotifications("Update anchors", "o/r@develop", errors.New("HTTP 403")) })
	if !strings.Contains(out, "could not read") || strings.Count(strings.TrimSpace(out), "\n") != 0 ||
		strings.Contains(out, "NOTIFICATIONS") {
		t.Errorf("a read failure is one line and no block, said:\n%s", out)
	}
}

// Local mode: the file at the project root is the message.
func TestNotificationsLocal(t *testing.T) {
	t.Run("NTFCT-B01: Local mode reads notifications.md at the project root", func(t *testing.T) {})
	root := t.TempDir()
	if notificationsLocal(root) != "" {
		t.Error("no file, no message")
	}
	if err := os.WriteFile(filepath.Join(root, NotificationsFile), []byte("hello agents"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := notificationsLocal(root); got != "hello agents" {
		t.Errorf("local message = %q", got)
	}
}
