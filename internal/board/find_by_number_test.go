package board

import (
	"strings"
	"testing"
)

// `deliver --card <n>` records on a card named directly: plan cards carry no `[CODE]` in
// the title and a change to generated files has no unit in the map (both measured in
// blue-eyes). The card must still be OPEN and an Anchors card — otherwise the record goes
// where nobody reads it, or onto someone else's issue.
func TestFindOpenByNumber(t *testing.T) {
	answer := ""
	c := Client{Repo: "o/r", Labels: []string{"anchors"}, run: func(args ...string) ([]byte, error) {
		if args[0] != "issue" || args[1] != "view" || args[2] != "801" {
			t.Fatalf("unexpected gh call: %v", args)
		}
		return []byte(answer), nil
	}}

	answer = `{"number":801,"title":"[plano] docs","body":"","state":"OPEN","labels":[{"name":"anchors"},{"name":"anchors:in-progress"}]}`
	card, err := c.FindOpenByNumber(801)
	if err != nil || card.Number != 801 {
		t.Fatalf("an open Anchors card was refused: %v %+v", err, card)
	}

	answer = `{"number":801,"title":"x","body":"","state":"CLOSED","labels":[{"name":"anchors"}]}`
	if _, err := c.FindOpenByNumber(801); err == nil || !strings.Contains(err.Error(), "closed") {
		t.Errorf("a closed card was accepted: %v", err)
	}

	answer = `{"number":801,"title":"x","body":"","state":"OPEN","labels":[{"name":"bug"}]}`
	if _, err := c.FindOpenByNumber(801); err == nil || !strings.Contains(err.Error(), "not an Anchors card") {
		t.Errorf("an issue without the Anchors label was accepted: %v", err)
	}
}
