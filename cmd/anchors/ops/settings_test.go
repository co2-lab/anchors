package ops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/settings"
)

// withStdin runs fn with os.Stdin reading `input`.
func withStdin(t *testing.T, input string, fn func()) {
	t.Helper()
	p := writeFile(t, t.TempDir(), "stdin", input)
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	orig := os.Stdin
	os.Stdin = f
	defer func() { os.Stdin = orig }()
	fn()
}

func loadSettings(t *testing.T, root string) settings.Settings {
	t.Helper()
	s, err := settings.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// The date is stamped by whoever declares — the command does not read the clock.
func TestSettingsDecisionsRequireADate(t *testing.T) {
	root := t.TempDir()
	for _, sub := range []string{"role", "user-issues"} {
		err, _ := runCmd(t, newSettingsCmd(), sub, "--root", root, "dev")
		if err == nil || !strings.Contains(err.Error(), "--date") {
			t.Errorf("%s without --date: %v", sub, err)
		}
	}
	if _, err := os.Stat(settings.Path(root)); err == nil {
		t.Error("a decision without a date was recorded")
	}
}

// Declaring a role records it with the agent and the date, and drops the legacy
// user_issues field so the next read has one source; `show` then lists what it allows.
func TestSettingsRoleIsRecordedAndShown(t *testing.T) {
	root := t.TempDir()
	if err, _ := runCmd(t, newSettingsCmd(), "user-issues", "--root", root, "--date", "2026-09-01", "sim"); err != nil {
		t.Fatal(err)
	}
	err, out := runCmd(t, newSettingsCmd(), "role", "--root", root, "--date", "2026-09-26", "architect")
	if err != nil {
		t.Fatalf("role: %v", err)
	}
	s := loadSettings(t, root)
	if s.Role != settings.RoleArchitect || s.DecidedAt != "2026-09-26" || s.Agent == "" {
		t.Errorf("recorded %+v; want architect, 2026-09-26, and the agent", s)
	}
	if s.UserIssues != nil {
		t.Error("the legacy user_issues survived the role declaration: two sources for one answer")
	}
	if !strings.Contains(out, filepath.Join(".anchors", "settings.yaml")) {
		t.Errorf("the output does not say where it was recorded:\n%s", out)
	}

	err, out = runCmd(t, newSettingsCmd(), "show", "--root", root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "capabilities:") {
		t.Errorf("show does not list the capabilities:\n%s", out)
	}
	for _, c := range settings.RoleArchitect.Caps() {
		if !strings.Contains(out, string(c)) {
			t.Errorf("show omits the capability %q:\n%s", c, out)
		}
	}
}

func TestSettingsRejectsAnUnknownRoleOrAnswer(t *testing.T) {
	root := t.TempDir()
	if err, _ := runCmd(t, newSettingsCmd(), "role", "--root", root, "--date", "2026-09-26", "wizard"); err == nil ||
		!strings.Contains(err.Error(), `"wizard"`) {
		t.Errorf("an unknown role: %v", err)
	}
	if err, _ := runCmd(t, newSettingsCmd(), "user-issues", "--root", root, "--date", "2026-09-26", "maybe"); err == nil ||
		!strings.Contains(err.Error(), `"maybe"`) {
		t.Errorf("an unknown answer: %v", err)
	}
	if _, err := os.Stat(settings.Path(root)); err == nil {
		t.Error("an invalid declaration was recorded")
	}
}

// Without a role, `show` teaches how to declare one.
func TestSettingsShowWithoutRoleTeachesHowToDeclare(t *testing.T) {
	err, out := runCmd(t, newSettingsCmd(), "show", "--root", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "anchors settings role --date") || strings.Contains(out, "capabilities:") {
		t.Errorf("show without a role:\n%s", out)
	}
}

// Asked on the terminal, a reply it does not understand is asked AGAIN, never assumed.
func TestSettingsUserIssuesAsksUntilItUnderstands(t *testing.T) {
	root := t.TempDir()
	var err error
	var out string
	withStdin(t, "talvez\nnao\n", func() {
		err, out = runCmd(t, newSettingsCmd(), "user-issues", "--root", root, "--date", "2026-09-26")
	})
	if err != nil {
		t.Fatalf("user-issues: %v", err)
	}
	if !strings.Contains(out, `did not understand "talvez"`) {
		t.Errorf("the unclear reply was not pointed out:\n%s", out)
	}
	s := loadSettings(t, root)
	if s.UserIssues == nil || *s.UserIssues {
		t.Errorf("recorded %v, want false (the `nao` after the retry)", s.UserIssues)
	}
}

func TestAskUserIssuesGivesUpWithoutAnAnswer(t *testing.T) {
	var err error
	captureStdout(t, func() {
		withStdin(t, "a\nb\nc\n", func() { _, err = askUserIssues() })
	})
	if err == nil || !strings.Contains(err.Error(), "no recognized answer") {
		t.Errorf("three unclear replies: %v", err)
	}
	captureStdout(t, func() {
		withStdin(t, "", func() { _, err = askUserIssues() })
	})
	if err == nil {
		t.Error("a closed input must be an error, not a silent no")
	}
}

// Without the argument, the role is asked on the terminal.
func TestSettingsRoleAsksWhenNotGiven(t *testing.T) {
	root := t.TempDir()
	var err error
	withStdin(t, "qa\n", func() {
		err, _ = runCmd(t, newSettingsCmd(), "role", "--root", root, "--date", "2026-09-26")
	})
	if err != nil {
		t.Fatal(err)
	}
	if s := loadSettings(t, root); s.Role != settings.RoleQA {
		t.Errorf("role = %q, want qa", s.Role)
	}
}
