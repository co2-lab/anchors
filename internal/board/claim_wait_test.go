package board

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

const waitAgent = "machine/session-1"

// fakeBoard stands in for `gh`: it answers the three calls `AskAndWait` makes — the
// dispatch, the run list and the issue list — from an in-memory board.
type fakeBoard struct {
	t          *testing.T
	runs       []ClaimRun // newest first, like `gh run list`
	nextID     int64
	dispatches int
	// onDispatch runs on every `gh workflow run`; by default it queues a run.
	onDispatch func(f *fakeBoard)
	// onPoll runs on every board read (`gh api graphql`, every look at the board).
	onPoll func(f *fakeBoard, poll int)
	polls  int
	card   bool // the claim handed the card: it shows up owned and in-progress
}

func (f *fakeBoard) queue(title, status, conclusion string) *ClaimRun {
	f.nextID++
	f.runs = append([]ClaimRun{{ID: f.nextID, Title: title, Status: status,
		Conclusion: conclusion}}, f.runs...)
	return &f.runs[0]
}

func (f *fakeBoard) find(id int64) *ClaimRun {
	for i := range f.runs {
		if f.runs[i].ID == id {
			return &f.runs[i]
		}
	}
	f.t.Fatalf("no run %d", id)
	return nil
}

func (f *fakeBoard) run(args ...string) ([]byte, error) {
	switch args[0] + " " + args[1] {
	case "workflow run":
		f.dispatches++
		if f.onDispatch != nil {
			f.onDispatch(f)
		} else {
			f.queue(ClaimRunTitle(waitAgent), "queued", "")
		}
		return nil, nil
	case "run list":
		return json.Marshal(f.runs)
	case "api graphql":
		f.polls++
		if f.onPoll != nil {
			f.onPoll(f, f.polls)
		}
		if !f.card {
			return nil, nil
		}
		return []byte(`{"number":4,"title":"[ABCDE] card","labels":[{"name":"anchors"},` +
			`{"name":"anchors:in-progress"}],"comments":[{"body":"anchors-owner: ` +
			waitAgent + `"}]}` + "\n"), nil
	}
	return nil, fmt.Errorf("unexpected gh call: %v", args)
}

func (f *fakeBoard) client() Client {
	return Client{Repo: "o/r", Labels: []string{"anchors"}, run: f.run}
}

// fakeClock advances only when the wait sleeps, so a timeout costs no real time.
func fakeClock() ClaimWait {
	now := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	return ClaimWait{
		Timeout: time.Minute, Interval: 5 * time.Second,
		Now:   func() time.Time { return now },
		Sleep: func(d time.Duration) { now = now.Add(d) },
	}
}

// `next` WAITS for its claim instead of telling the agent to run `next` again: the card
// that arrives a few looks later is returned by the same call, with ONE dispatch.
func TestAskAndWait_returnsTheCardOfTheRunItDispatched(t *testing.T) {
	f := &fakeBoard{t: t}
	f.onPoll = func(f *fakeBoard, poll int) {
		if poll == 3 {
			f.card = true
		}
	}
	out, err := f.client().AskAndWait(waitAgent, fakeClock())
	if err != nil {
		t.Fatal(err)
	}
	if out.Card == nil || out.Card.Number != 4 {
		t.Fatalf("the card the claim handed out was not returned: %+v", out)
	}
	if f.dispatches != 1 || !out.Dispatched {
		t.Errorf("dispatches = %d, want exactly 1", f.dispatches)
	}
}

// A CLAIM OF THIS AGENT STILL PENDING is waited on, never dispatched again — a second
// dispatch cancels the pending run or queues a duplicate claim (blue-eyes #679).
func TestAskAndWait_neverDispatchesWhileItsRunIsPending(t *testing.T) {
	f := &fakeBoard{t: t}
	pending := f.queue(ClaimRunTitle(waitAgent), "in_progress", "").ID
	f.onPoll = func(f *fakeBoard, poll int) {
		if poll == 2 {
			r := f.find(pending)
			r.Status, r.Conclusion = "completed", "success"
			f.card = true
		}
	}
	out, err := f.client().AskAndWait(waitAgent, fakeClock())
	if err != nil {
		t.Fatal(err)
	}
	if f.dispatches != 0 || out.Dispatched {
		t.Errorf("dispatched %d claim(s) while this agent's run was pending", f.dispatches)
	}
	if out.Card == nil {
		t.Errorf("the pending run's card was not returned: %+v", out)
	}
}

// ANOTHER agent's pending run does not hold this one back: the title tells them apart.
func TestAskAndWait_anotherAgentsRunDoesNotCount(t *testing.T) {
	f := &fakeBoard{t: t}
	f.queue(ClaimRunTitle("other/session"), "queued", "")
	f.onPoll = func(f *fakeBoard, poll int) { f.card = true }
	if _, err := f.client().AskAndWait(waitAgent, fakeClock()); err != nil {
		t.Fatal(err)
	}
	if f.dispatches != 1 {
		t.Errorf("dispatches = %d: another agent's pending claim was taken for this one's",
			f.dispatches)
	}
}

// THE RUN FINISHED WITHOUT A CARD (empty board, frozen project): the wait ends there,
// with the run, instead of spending the whole timeout.
func TestAskAndWait_runDoneWithoutCardEndsTheWait(t *testing.T) {
	f := &fakeBoard{t: t}
	f.onDispatch = func(f *fakeBoard) {
		f.queue(ClaimRunTitle(waitAgent), "completed", "success")
	}
	w := fakeClock()
	start := w.Now()
	out, err := f.client().AskAndWait(waitAgent, w)
	if err != nil {
		t.Fatal(err)
	}
	if out.Card != nil || out.TimedOut || out.Run == nil || !out.Run.Done() {
		t.Fatalf("want a finished run and no card, got %+v", out)
	}
	if waited := w.Now().Sub(start); waited != 0 {
		t.Errorf("waited %s after the run had finished", waited)
	}
}

// AN OLDER FINISHED RUN of the same agent is not the answer to this request: taking it
// would end the wait with "no card" before the new claim even ran.
func TestAskAndWait_ignoresTheAgentsOlderRuns(t *testing.T) {
	f := &fakeBoard{t: t}
	f.queue(ClaimRunTitle(waitAgent), "completed", "success") // a previous `next`
	f.onPoll = func(f *fakeBoard, poll int) {
		if poll == 3 {
			f.card = true
		}
	}
	out, err := f.client().AskAndWait(waitAgent, fakeClock())
	if err != nil {
		t.Fatal(err)
	}
	if out.Card == nil {
		t.Errorf("the wait ended on the agent's previous run: %+v", out)
	}
}

// THE WAIT IS BOUNDED, and running out of time does not dispatch a second claim.
func TestAskAndWait_timesOutWithoutASecondDispatch(t *testing.T) {
	f := &fakeBoard{t: t}
	out, err := f.client().AskAndWait(waitAgent, fakeClock())
	if err != nil {
		t.Fatal(err)
	}
	if !out.TimedOut {
		t.Fatalf("a run that never finishes must end in a timeout: %+v", out)
	}
	if out.Run == nil {
		t.Error("the timeout should name the pending run, so the agent can follow it")
	}
	if f.dispatches != 1 {
		t.Errorf("dispatches = %d, want 1", f.dispatches)
	}
}

// A CANCELLED run never ran — another agent's dispatch replaced it in the concurrency
// group. Asking again is not a second claim, and it is bounded.
func TestAskAndWait_redispatchesACancelledRunBounded(t *testing.T) {
	f := &fakeBoard{t: t}
	f.onDispatch = func(f *fakeBoard) {
		f.queue(ClaimRunTitle(waitAgent), "completed", "cancelled")
	}
	out, err := f.client().AskAndWait(waitAgent, fakeClock())
	if err != nil {
		t.Fatal(err)
	}
	if want := 1 + maxClaimRedispatch; f.dispatches != want {
		t.Errorf("dispatches = %d, want %d (the first + the bounded re-dispatches)",
			f.dispatches, want)
	}
	if out.Run == nil || out.Run.Conclusion != "cancelled" {
		t.Errorf("the last cancelled run should be reported: %+v", out)
	}
}
