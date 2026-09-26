package common

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/config"
)

// fakeGH puts a `gh` on PATH that records its arguments and prints output — no network,
// no real GitHub. It returns the file where the arguments land, one per line.
func fakeGH(t *testing.T, output string, exitCode int) string {
	t.Helper()
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args")
	outFile := filepath.Join(dir, "out")
	writeFile(t, outFile, output)
	script := "#!/bin/sh\nfor a in \"$@\"; do printf '%s\\n' \"$a\" >> '" + argsFile + "'; done\n" +
		"cat '" + outFile + "'\nexit " + string(rune('0'+exitCode)) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+"/bin"+string(os.PathListSeparator)+"/usr/bin")
	return argsFile
}

func workflowCfg() *config.Config {
	return &config.Config{Workflow: &config.Workflow{Repo: "acme/app", Labels: []string{"anchors"}}}
}

// The cards whose LAST owner is this agent, parsed from gh's tab-separated lines; the
// state label loses its `anchors:` prefix, and malformed lines are skipped.
func TestAgentCards_parsesTheCardsOfThisAgent(t *testing.T) {
	args := fakeGH(t, "12\tFix the parser\tanchors:doing\n40\tWrite docs\t\nbroken line\n\tno number\tx\n51\tShip it\tanchors:review\n", 0)
	t.Setenv("ANCHORS_AGENT", "host/worker-1")

	got := AgentCards(workflowCfg())
	want := []AgentCard{
		{Numero: "12", Titulo: "Fix the parser", Estado: "doing"},
		{Numero: "40", Titulo: "Write docs", Estado: ""},
		{Numero: "51", Titulo: "Ship it", Estado: "review"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
	b, err := os.ReadFile(args)
	if err != nil {
		t.Fatal(err)
	}
	called := string(b)
	for _, frag := range []string{"acme/app", "--label\nanchors\n", "anchors-owner: host/worker-1"} {
		if !strings.Contains(called, frag) {
			t.Errorf("gh was not called with %q:\n%s", frag, called)
		}
	}
}

// A card with no state label is kept in the LAST position too: its line ends in the tab
// of the empty field, and trimming the whole output ate it.
func TestAgentCards_lastCardWithoutStateIsKept(t *testing.T) {
	fakeGH(t, "12\tFix the parser\tanchors:doing\n40\tWrite docs\t\n", 0)
	t.Setenv("ANCHORS_AGENT", "host/worker-1")
	got := AgentCards(workflowCfg())
	want := []AgentCard{
		{Numero: "12", Titulo: "Fix the parser", Estado: "doing"},
		{Numero: "40", Titulo: "Write docs", Estado: ""},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// No agent name, no workflow, no gh, or gh failing: no cards, and no error to trip on.
func TestAgentCards_nothingWithoutTheIngredients(t *testing.T) {
	fakeGH(t, "12\tFix\tanchors:doing\n", 0)
	t.Setenv("ANCHORS_AGENT", "")
	if got := AgentCards(workflowCfg()); got != nil {
		t.Errorf("without ANCHORS_AGENT: %+v", got)
	}
	t.Setenv("ANCHORS_AGENT", "host/w")
	if got := AgentCards(&config.Config{}); got != nil {
		t.Errorf("without a workflow: %+v", got)
	}
	if got := AgentCards(nil); got != nil {
		t.Errorf("without a config: %+v", got)
	}

	fakeGH(t, "12\tFix\tanchors:doing\n", 1)
	if got := AgentCards(workflowCfg()); got != nil {
		t.Errorf("with gh failing: %+v", got)
	}

	t.Setenv("PATH", t.TempDir())
	if got := AgentCards(workflowCfg()); got != nil {
		t.Errorf("without gh installed: %+v", got)
	}
}

func TestFirstLineOfReason(t *testing.T) {
	if got := FirstLineOfReason("  the parser drops the header  \nmore detail\n"); got != "the parser drops the header" {
		t.Errorf("first line: %q", got)
	}
	long := strings.Repeat("x", 80)
	got := FirstLineOfReason(long)
	if len(got) != 70 || !strings.HasSuffix(got, "...") {
		t.Errorf("a long line must be cut to 70 with an ellipsis: %q (%d)", got, len(got))
	}
	if exact := strings.Repeat("y", 70); FirstLineOfReason(exact) != exact {
		t.Error("a line of exactly 70 must not be cut")
	}
}

func TestNumeroDaIssue(t *testing.T) {
	for url, want := range map[string]string{
		"https://github.com/acme/app/issues/433":   "433",
		"https://github.com/acme/app/issues/433\n": "433",
		"https://github.com/acme/app/issues/":      "",
		"https://github.com/acme/app/pull/12a":     "",
		"433":                                      "",
	} {
		if got := NumeroDaIssue(url); got != want {
			t.Errorf("NumeroDaIssue(%q) = %q, want %q", url, got, want)
		}
	}
}

// The card a PR declares is read from `Refs/Closes/Fixes/Resolves #N` at a line start.
func TestCardDoPR_readsTheLinkInTheBody(t *testing.T) {
	args := fakeGH(t, "Some summary.\n\nCloses #77\nRefs #80\n", 0)
	if got := CardDoPR("acme/app", " #15 "); got != "77" {
		t.Errorf("got %q, want 77", got)
	}
	b, _ := os.ReadFile(args)
	if !strings.HasPrefix(string(b), "pr\nview\n15\n--repo\nacme/app\n") {
		t.Errorf("gh called with the wrong arguments:\n%s", b)
	}

	fakeGH(t, "mentions closes #9 in the middle of a line\n", 0)
	if got := CardDoPR("acme/app", "15"); got != "" {
		t.Errorf("a link not at the line start must not count: %q", got)
	}
	fakeGH(t, "Closes #77\n", 1)
	if got := CardDoPR("acme/app", "15"); got != "" {
		t.Errorf("gh failing must give no card: %q", got)
	}
	if CardDoPR("", "15") != "" || CardDoPR("acme/app", " # ") != "" {
		t.Error("without repo or PR there is nothing to ask")
	}
}
