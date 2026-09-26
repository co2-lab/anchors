package flow

import (
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/initx"
)

// Each card drags what was born under it, and the lines come out in numeric order — so
// #101 does not sort before #45 as text would.
func TestPRBodyCmd_linksTheCardsAndWhatWasBornUnderThem(t *testing.T) {
	root := githubProject(t)
	calls := scriptedGH(t,
		ghRule{match: "issue list *--label " + initx.LabelSob("44") + " *", out: `[{"number":101},{"number":45}]`},
		ghRule{match: "issue list *--label " + initx.LabelSob("50") + " *", out: `[]`},
	)
	cmd := newPRBodyCmd()
	cmd.SetArgs([]string{"--root", root, "--cards", "#44, 50"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if out != "Refs #44\nRefs #45\nRefs #50\nRefs #101\n" {
		t.Errorf("unexpected body:\n%s", out)
	}
	if len(callsWith(calls(), "--repo acme/app")) != 2 {
		t.Errorf("the findings are looked up in the configured repo: %v", calls())
	}
}

// `--so-sob` prints only what was born under the given cards — what CI checks is present.
func TestPRBodyCmd_onlyUnderLeavesTheRootsOut(t *testing.T) {
	root := githubProject(t)
	scriptedGH(t, ghRule{match: "issue list *", out: `[{"number":45}]`})
	cmd := newPRBodyCmd()
	cmd.SetArgs([]string{"--root", root, "--cards", "44", "--so-sob"})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if out != "Refs #45\n" {
		t.Errorf("only the finding under #44 is printed, got:\n%s", out)
	}
}

// Without --cards, the cards are the agent's own on the board.
func TestPRBodyCmd_discoversTheAgentsCards(t *testing.T) {
	root := githubProject(t)
	scriptedGH(t,
		ghRule{match: "issue list *--json number,title,labels,comments*", out: "12\tmine\tanchors:in-progress"},
		ghRule{match: "issue list *", out: `[]`},
	)
	t.Setenv("ANCHORS_AGENT", "host/dev1")
	cmd := newPRBodyCmd()
	cmd.SetArgs([]string{"--root", root})
	var err error
	out := stdoutOf(t, func() { err = cmd.Execute() })
	if err != nil {
		t.Fatal(err)
	}
	if out != "Refs #12\n" {
		t.Errorf("the agent's card must be linked, got:\n%s", out)
	}
}

func TestPRBodyCmd_refusals(t *testing.T) {
	scriptedGH(t)
	cmd := newPRBodyCmd()
	cmd.SetArgs([]string{"--root", githubProject(t)})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "no card") {
		t.Errorf("without cards nothing can be linked, got %v", err)
	}

	cmd = newPRBodyCmd()
	cmd.SetArgs([]string{"--root", localProject(t), "--cards", "4"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "github mode") {
		t.Errorf("local mode must be refused, got %v", err)
	}
}

func TestCardsUnder_failureYieldsNothing(t *testing.T) {
	for _, r := range []ghRule{{match: "*", code: 1}, {match: "*", out: "not json"}} {
		scriptedGH(t, r)
		if got := cardsUnder(cfgGitHub(), "44"); got != nil {
			t.Errorf("a failed lookup yields nothing, got %v", got)
		}
	}
}
