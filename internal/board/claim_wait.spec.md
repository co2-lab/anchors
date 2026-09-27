<!-- @anchors
  code: CLWTC
  updated_at: 2026-09-26
  layer: apoio
-->
# ClaimWait — asking the claim pipeline for a card and waiting, bounded, for the answer

> **Code**: `CLWTC`

## Overview

In github mode an agent asks the claim pipeline for work instead of claiming a card itself (`BRCRB`).
The dispatch returns at once, and the card arrives seconds, or behind other agents' claims minutes,
later. `anchors next` used to stop right after dispatching and tell the agent to run it again, and each
re-run while the claim was still pending dispatched another claim: GitHub keeps one pending run per
concurrency group, so the new dispatch cancelled the previous one, or queued a second claim for the
same agent.

So the agent waits for the run it dispatched, and never dispatches while one of its runs is still
pending. `gh workflow run` returns no run id, so the agent's runs are told apart by their title, which
the claim pipeline writes, and its own run is the oldest one with an id above every run that existed
before; run ids only grow, so no clock is compared between the machine and GitHub. The wait ends as soon
as the card is the agent's, or when its run finishes without one, or when the time runs out. A run that
was cancelled (replaced in the group by another agent's dispatch) never ran, so it is asked again, a
bounded number of times.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the agent | the agent's identity, the same the pipeline writes in the run title | — | the caller (`anchors next`) |
| the wait bounds | a timeout and an interval, or zero for the defaults | — | this unit fills the defaults |
| the run list | the claim pipeline's recent runs as `gh run list` gives them | — | `gh` |

## Effects

| Effect | Description |
| --- | --- |
| `CLWTC-B01` | An agent's claim runs are the runs titled `anchors · claim for <agent>`; another agent's runs never count as its own. (`ClaimRunTitle`) |
| `CLWTC-B02` | While one of the agent's runs is still pending, no claim is dispatched: the wait follows the oldest pending run. |
| `CLWTC-B03` | The run the wait follows is the oldest run above every run of the agent that existed before; an older finished run is not its answer. |
| `CLWTC-B04` | The wait returns as soon as the agent owns a card under way, having dispatched once. |
| `CLWTC-B05` | When the followed run finishes, the board is looked at once more; with no card then, the wait ends with the run and without spending the timeout. |
| `CLWTC-B06` | A cancelled run is dispatched again, at most twice, and only before the deadline; the last cancelled run is reported. |
| `CLWTC-B07` | The wait is bounded (three minutes by default, looking every five seconds): when the time runs out the outcome says so and names the pending run, without a second dispatch. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CLWTC-E01` | The run list is not the expected JSON. | The wait fails with the error before dispatching anything. | Without the run list the agent cannot know whether a claim of its own is pending, and dispatching blind is what duplicated claims. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/board/board.go` | `Ask`, `Mine` | apoio — the dispatch and the look at the board (`BRCRB`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
