package quality

import (
	"os"
	"strings"
	"testing"

	"github.com/co2-lab/anchors/internal/gate"
)

// A COUNT IS NOT AN ADDRESS.
//
// The earlier version printed only "%d targets awaiting judgment", and the reviewer knew
// there was work without knowing which, where, or what to ask. Measured in the reference
// app: 59 of 85 specs never received a verdict, with CI running on every PR.
func judgedFixture() []gate.Result {
	return []gate.Result{
		{Gate: "rule-fulfilled", Target: "a/One.spec.md", Verdict: gate.Judge},
		{Gate: "rule-fulfilled", Target: "b/Two.spec.md", Verdict: gate.Judge},
		{Gate: "test-proves", Target: "c/Three.test.ts", Verdict: gate.Judge},
	}
}

func brief(t *testing.T) string {
	t.Helper()
	return capturaSaida(t, func() {
		printJudgmentBrief(judgedFixture(),
			map[string]string{"rule-fulfilled": "SPEC.md"},
			map[string]string{"rule-fulfilled": "Does the excerpt DO what the rule describes?"})
	})
}

// The TARGET: without it the reviewer does not know where to look.
func TestJudgmentBrief_namesEveryTarget(t *testing.T) {
	out := brief(t)
	for _, target := range []string{"a/One.spec.md", "b/Two.spec.md", "c/Three.test.ts"} {
		if !strings.Contains(out, target) {
			t.Errorf("the brief does not name the target %q:\n%s", target, out)
		}
	}
}

// The QUESTION: it is the `ask` the gate declares, and without it the reviewer invents
// the criterion.
func TestJudgmentBrief_carriesTheQuestionAndTheGuide(t *testing.T) {
	out := brief(t)
	if !strings.Contains(out, "DO what the rule describes") {
		t.Errorf("the brief does not carry the gate's question:\n%s", out)
	}
	if !strings.Contains(out, "SPEC.md") {
		t.Errorf("the brief does not say where the ruler is:\n%s", out)
	}
}

// GROUPED BY GATE: a gate asks the same thing of several targets, and repeating the
// question per target would make the reviewer read it three times to answer once.
func TestJudgmentBrief_groupsByGate(t *testing.T) {
	out := brief(t)
	if n := strings.Count(out, "DO what the rule describes"); n != 1 {
		t.Errorf("the question appears %d times, want 1 (grouped by gate):\n%s", n, out)
	}
	if !strings.Contains(out, "rule-fulfilled") || !strings.Contains(out, "test-proves") {
		t.Errorf("the brief does not name both gates:\n%s", out)
	}
}

// WHO JUDGES: the pipeline hands over the list, the reviewer decides. Without saying so,
// the list reads as the machine's backlog — and nobody works through it.
func TestJudgmentBrief_saysThePipelineDoesNotJudge(t *testing.T) {
	out := brief(t)
	flat := strings.Join(strings.Fields(out), " ")
	if !strings.Contains(flat, "NAO julga") && !strings.Contains(flat, "NOT judge") {
		t.Errorf("the brief does not say the pipeline does not judge:\n%s", out)
	}
	if !strings.Contains(out, "REV-CK5") {
		t.Errorf("the brief does not point at the review checklist item:\n%s", out)
	}
}

// UP TO TEN, and the rest counted: sixty paths drown the question that comes before them.
func TestJudgmentBrief_truncatesALongList(t *testing.T) {
	var many []gate.Result
	for i := 0; i < 25; i++ {
		many = append(many, gate.Result{Gate: "rule-fulfilled",
			Target: "u/" + string(rune('a'+i)) + ".spec.md", Verdict: gate.Judge})
	}
	out := capturaSaida(t, func() {
		printJudgmentBrief(many, map[string]string{}, map[string]string{})
	})
	if !strings.Contains(out, "15") {
		t.Errorf("the brief does not count the ones left out:\n%s", out)
	}
	if strings.Count(out, ".spec.md") > 11 {
		t.Errorf("the brief listed more than ten targets:\n%s", out)
	}
}

// THE BRIEF IS A REPORT, NOT A RECORD — and its call has to live OUTSIDE the record block.
//
// It was born inside `if !noRecord`, and the effect was measured in the reference app:
// the pipeline runs `check --all --no-record` (on purpose — recording from there would
// open a card on every push), which is EXACTLY the mode in which the reviewer needs the
// list. The one place the brief mattered was the one place it did not print.
//
// The guard is on the ORDER of the source because that is where the defect lives: the
// brief called after `if !noRecord {` disappears in the pipeline again, and no output test
// would catch it without assembling a whole project in github mode.
func TestJudgmentBrief_staysOutsideTheRecordBlock(t *testing.T) {
	src, err := os.ReadFile("check.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(src)
	call := strings.Index(text, "printJudgmentBrief(profile.Judged")
	if call < 0 {
		t.Fatal("the brief's call disappeared from check's flow")
	}
	record := strings.Index(text, "if !noRecord {")
	if record < 0 {
		t.Fatal("the record block disappeared")
	}
	if call > record {
		t.Error("the brief went back INSIDE the record block — it disappears under " +
			"`--no-record`, which is the pipeline's mode and the only one where it matters")
	}
}
