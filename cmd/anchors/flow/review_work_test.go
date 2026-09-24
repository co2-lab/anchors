package flow

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/board"
)

func stdoutOf(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = orig
	b, _ := io.ReadAll(r)
	return string(b)
}

// A review card whose body names no unit (a plan card, a finding under another card)
// printed `anchors work review --for ` with an empty target — seen in blue-eyes on #931.
// It now points to the PR that references the card.
func TestPrintReviewWork_cardWithNoUnitPointsToItsPR(t *testing.T) {
	out := stdoutOf(t, func() {
		printReviewWork(t.TempDir(), &board.Card{Number: 931, Title: "[plano] x", Body: "no unit here"}, "host/dev3")
	})
	if strings.Contains(out, "--for \n") || strings.Contains(out, "--for  ") {
		t.Errorf("the review command still has an empty target:\n%s", out)
	}
	if !strings.Contains(out, `#931 in:body`) {
		t.Errorf("the output does not point to the card's PR:\n%s", out)
	}
}

// A review ends with the verdict line on the PR, and `anchors next` is the only place the
// reviewer is told so (blue-eyes #998, #1001, #1009: reviewers ended with "Veredito: OK",
// `judge`, `decided`, and the card kept coming back). The output names the exact line for
// this agent, and no longer sends a reviewer to open a PR of their own.
func TestPrintReviewWork_endsWithTheVerdictLine(t *testing.T) {
	out := stdoutOf(t, func() {
		printReviewWork(t.TempDir(), &board.Card{Number: 869, Title: "[plano] x", Body: "no unit here"}, "host/dev3")
	})
	for _, want := range []string{"anchors-review: approved by host/dev3", "anchors-review: rejected by host/dev3", "releases you"} {
		if !strings.Contains(out, want) {
			t.Errorf("the review work should say %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "anchors pr-body") {
		t.Errorf("a reviewer is not told to open a PR of their own:\n%s", out)
	}
}
