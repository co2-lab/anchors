package board

import (
	"encoding/json"
	"fmt"
	"time"
)

// --- WAITING FOR THE CLAIM that `next` dispatched ---
//
// `Ask` only dispatches the claim pipeline; the card arrives seconds (or, behind other
// agents' claims, minutes) later. `anchors next` used to stop right there and tell the
// agent to run `anchors next` again — and every re-run while the claim was still pending
// dispatched ANOTHER claim. GitHub keeps one pending run per concurrency group, so each new
// dispatch cancelled the previous pending one or, once the first had started, queued a
// second claim for the same agent (blue-eyes #679).
//
// So `next` waits for the run it dispatched, and never asks twice while one of its runs is
// still pending. The wait is bounded: a claim that does not finish in time is reported
// with its run id, and a later `anchors next` picks the same run up instead of dispatching.

// ClaimRunTitle is the `run-name` the claim pipeline gives the run of each agent.
//
// `gh workflow run` returns no run id, and the run list does not show the inputs — the
// title is the only way to tell THIS agent's run from another agent's. The claim template
// (`internal/initx/workflows/anchors-claim.yml`) must write exactly this text.
func ClaimRunTitle(agent string) string {
	return "anchors · claim for " + agent
}

// ClaimRun is one run of the claim pipeline, as `gh run list` reports it.
type ClaimRun struct {
	ID         int64  `json:"databaseId"`
	Title      string `json:"displayTitle"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

// Done says whether the run finished — with any conclusion.
func (r ClaimRun) Done() bool { return r.Status == "completed" }

// ClaimWait bounds the wait. The zero value is completed with the defaults.
type ClaimWait struct {
	// Timeout is how long `AskAndWait` waits for the card. Default: DefaultClaimTimeout.
	Timeout time.Duration
	// Interval is the pause between two looks at the board. Default: 5s.
	Interval time.Duration
	// Now and Sleep replace the clock in tests.
	Now   func() time.Time
	Sleep func(time.Duration)
}

// DefaultClaimTimeout is how long `next` waits for its claim.
//
// A claim runs in well under a minute, but it is serialized: behind a few other agents'
// claims the run waits its turn. Three minutes covers that queue without turning a stuck
// pipeline into a session spent waiting.
const DefaultClaimTimeout = 3 * time.Minute

// maxClaimRedispatch bounds the re-dispatch of a run that was CANCELLED — which is what
// GitHub does to a pending run when another agent's dispatch replaces it in the group.
const maxClaimRedispatch = 2

// ClaimOutcome is what the wait found.
type ClaimOutcome struct {
	// Card is the card the claim handed to the agent, or nil.
	Card *Card
	// Run is this agent's claim run, when it could be identified.
	Run *ClaimRun
	// Dispatched says whether this call dispatched a claim (false: it waited on one that
	// was already pending).
	Dispatched bool
	// TimedOut says the wait ran out before the run finished or the card appeared.
	TimedOut bool
}

// claimRuns lists this agent's claim runs, newest first.
func (c Client) claimRuns(agent string) ([]ClaimRun, error) {
	out, err := c.gh("run", "list", "--workflow", ClaimWorkflow,
		"--event", "workflow_dispatch", "--limit", "50",
		"--json", "databaseId,displayTitle,status,conclusion")
	if err != nil {
		return nil, err
	}
	var all []ClaimRun
	if err := json.Unmarshal(out, &all); err != nil {
		return nil, fmt.Errorf("gh run list: response is not the expected JSON: %w", err)
	}
	title := ClaimRunTitle(agent)
	var mine []ClaimRun
	for _, r := range all {
		if r.Title == title {
			mine = append(mine, r)
		}
	}
	return mine, nil
}

// AskAndWait asks the pipeline for a card and waits, bounded, for the answer.
//
// It dispatches only when this agent has no claim run still pending; otherwise it waits
// on that one. Its own run is told apart by the title (`ClaimRunTitle`) and by an id
// greater than every run of this agent that existed before the dispatch — run ids only
// grow, so no clock comparison between this machine and GitHub is needed.
//
// It returns as soon as the card appears (`Mine`), or when the run finishes without one,
// or when the timeout runs out. A pipeline too old to carry the `run-name` has runs no
// title matches: the wait then falls back to looking for the card until the timeout.
func (c Client) AskAndWait(agent string, w ClaimWait) (ClaimOutcome, error) {
	if w.Timeout <= 0 {
		w.Timeout = DefaultClaimTimeout
	}
	if w.Interval <= 0 {
		w.Interval = 5 * time.Second
	}
	if w.Now == nil {
		w.Now = time.Now
	}
	if w.Sleep == nil {
		w.Sleep = time.Sleep
	}
	deadline := w.Now().Add(w.Timeout)

	before, err := c.claimRuns(agent)
	if err != nil {
		return ClaimOutcome{}, err
	}
	var outcome ClaimOutcome
	// `after` is the id above which a run is the one this call is waiting for.
	var after int64
	var pending *ClaimRun
	for i := range before {
		if before[i].ID > after {
			after = before[i].ID
		}
		if !before[i].Done() && (pending == nil || before[i].ID < pending.ID) {
			pending = &before[i]
		}
	}
	if pending != nil {
		// A claim of this agent is still queued or running: dispatching again is what
		// cancelled or duplicated claims. Wait on the oldest pending one.
		after = pending.ID - 1
	} else {
		if err := c.Ask(agent); err != nil {
			return ClaimOutcome{}, err
		}
		outcome.Dispatched = true
	}

	redispatched := 0
	for {
		card, err := c.Mine(agent)
		if err != nil {
			return outcome, err
		}
		if card != nil {
			outcome.Card = card
			return outcome, nil
		}

		runs, err := c.claimRuns(agent)
		if err != nil {
			return outcome, err
		}
		var run *ClaimRun
		for i := range runs {
			// The OLDEST run after the mark is the one this call started (or found
			// pending); a newer one belongs to a later `next`.
			if runs[i].ID > after && (run == nil || runs[i].ID < run.ID) {
				run = &runs[i]
			}
		}
		if run != nil {
			outcome.Run = run
			if run.Done() {
				// The card can land between the look at the board and the look at the
				// run: one more look before concluding there is none.
				card, err := c.Mine(agent)
				if err != nil {
					return outcome, err
				}
				if card != nil {
					outcome.Card = card
					return outcome, nil
				}
				if run.Conclusion == "cancelled" && redispatched < maxClaimRedispatch &&
					w.Now().Before(deadline) {
					// Replaced in the concurrency group by another agent's dispatch —
					// the request never ran, so asking again is not a second claim.
					if err := c.Ask(agent); err != nil {
						return outcome, err
					}
					redispatched++
					after = run.ID
					outcome.Run = nil
					continue
				}
				return outcome, nil
			}
		}

		if !w.Now().Before(deadline) {
			outcome.TimedOut = true
			return outcome, nil
		}
		w.Sleep(w.Interval)
	}
}
