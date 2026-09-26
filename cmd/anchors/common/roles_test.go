package common

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/settings"
)

// The agent ID is <host>/<session>, and the session falls back from ANCHORS_SESSION to
// USER to "default" — so two sessions on one machine do not claim the same card.
func TestAgentID_sessionFallbacks(t *testing.T) {
	host, _ := os.Hostname()
	if host == "" {
		host = "local"
	}
	t.Setenv("ANCHORS_SESSION", "  worker-2 ")
	t.Setenv("USER", "alice")
	if got := AgentID(); got != host+"/worker-2" {
		t.Errorf("with ANCHORS_SESSION: %q", got)
	}
	t.Setenv("ANCHORS_SESSION", "")
	if got := AgentID(); got != host+"/alice" {
		t.Errorf("falling back to USER: %q", got)
	}
	t.Setenv("USER", "")
	if got := AgentID(); got != host+"/default" {
		t.Errorf("falling back to default: %q", got)
	}
}

func TestRoleList_listsEveryKnownRole(t *testing.T) {
	list := RoleList()
	for _, r := range settings.KnownRoles() {
		if !strings.Contains(list, string(r)) || !strings.Contains(list, r.Does()) {
			t.Errorf("role %q (or what it does) is missing from the list:\n%s", r, list)
		}
	}
}

// Whoever declares a role reads what it means — above all whether it acts on the
// escalated cards, which is the difference that most changes the work.
func TestPrintRole_saysWhetherItDecidesProduct(t *testing.T) {
	po := captureStdout(t, func() { PrintRole(settings.RolePO) })
	if !strings.Contains(po, settings.RolePO.Title()) || !strings.Contains(po, "`needs-user`") {
		t.Errorf("the product owner was not told it acts on escalated cards:\n%s", po)
	}
	dev := captureStdout(t, func() { PrintRole(settings.RoleDev) })
	if !strings.Contains(dev, "anchors escalate") {
		t.Errorf("the dev was not told to escalate instead of deciding:\n%s", dev)
	}
	sec := captureStdout(t, func() { PrintRole(settings.RoleReviewerSec) })
	if lens := settings.RoleReviewerSec.Lens(); lens == "" || !strings.Contains(sec, lens) {
		t.Errorf("the security reviewer's lens was not shown:\n%s", sec)
	}
}

// The terminal question accepts the abbreviation people type, and asks again (up to three
// times) for what it does not recognise instead of assuming a role.
func TestAskRole_retriesThenAccepts(t *testing.T) {
	withStdin(t, "wizard\npo\n")
	var got settings.Role
	var err error
	out := captureStdout(t, func() { got, err = AskRole() })
	if err != nil || got != settings.RolePO {
		t.Fatalf("AskRole = %q, %v; want product-owner", got, err)
	}
	if !strings.Contains(out, `"wizard"`) {
		t.Errorf("the unrecognised answer was not echoed back:\n%s", out)
	}
}

func TestAskRole_givesUpAfterThreeAttempts(t *testing.T) {
	withStdin(t, "a\nb\nc\npo\n")
	var err error
	captureStdout(t, func() { _, err = AskRole() })
	if err == nil || !strings.Contains(err.Error(), "anchors settings role") {
		t.Errorf("expected the give-up error pointing at `anchors settings role`, got %v", err)
	}
}

func TestAskRole_closedInputIsAnError(t *testing.T) {
	withStdin(t, "")
	var err error
	captureStdout(t, func() { _, err = AskRole() })
	if err == nil {
		t.Error("an input closed before any answer must be an error, not a role")
	}
}

// Nobody on the other side of a pipe: the question must not be asked.
func TestInteractiveTerminal_pipeIsNotInteractive(t *testing.T) {
	withStdin(t, "dev\n")
	if InteractiveTerminal() {
		t.Error("a pipe on stdin was taken as an interactive terminal")
	}
	nul, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer nul.Close()
	old := os.Stdin
	os.Stdin = nul
	defer func() { os.Stdin = old }()
	if InteractiveTerminal() {
		t.Error("/dev/null on stdin was taken as an interactive terminal")
	}
}

func TestDecidesProduct_followsTheDeclaredRole(t *testing.T) {
	root := t.TempDir()
	if DecidesProduct(root) {
		t.Error("with no settings, nothing is decided — the default is closed")
	}
	if err := settings.Save(root, settings.Settings{Role: settings.RolePO}); err != nil {
		t.Fatal(err)
	}
	if !DecidesProduct(root) {
		t.Error("a product owner decides the product")
	}
	if err := settings.Save(root, settings.Settings{Role: settings.RoleDev}); err != nil {
		t.Fatal(err)
	}
	if DecidesProduct(root) {
		t.Error("a dev does not decide the product")
	}
	writeFile(t, filepath.Join(root, ".anchors", "settings.yaml"), "role: [unclosed\n")
	if DecidesProduct(root) {
		t.Error("an unreadable settings file must not unlock the capability")
	}
}
