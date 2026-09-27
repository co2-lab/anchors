<!-- @anchors
  code: SGSTS
  updated_at: 2026-09-26
  layer: infra
-->
# SuggestionStore — a proposed fix, as a patch plus its reason, waiting for someone to decide

> **Code**: `SGSTS`

## Overview

The gates can accuse, but the reader then has to rebuild by hand the fix the tool already knew. A
suggestion carries that fix: a git patch that can be reviewed, applied with `git apply`, or refused at
no cost, plus the reason in prose. A patch rather than an edited file keeps the order right: Anchors
proposes, someone decides, and nothing changes under the feet of whoever is editing.

Suggestions live in the project's `suggestions/` folder, as reviewable project content. The state is the
folder a suggestion sits in: pending, approved or rejected, with no status field inside the file that
could disagree with it. A rejected suggestion is not garbage: it is the record that stops the same
proposal from coming back at the next scan as if it were new. A decision always carries its written
reason and who took it, and a decision taken by the AI under automatic judgment is marked as such, so
"nobody looked at this" never passes for a human approval.

A finding whose fix is not mechanical may still be recorded as a suggestion without a patch: it says so
explicitly instead of faking a fix.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the suggestion ID | a stable, file-name-safe identity derived from the gate and the target | an empty ID | this unit: opening is refused |
| the patch | a unified diff, or empty when the fix is not mechanical | — | the proposing gate or judgment |
| the decision | approved or rejected, with a non-blank reason | any other state, a blank reason | this unit: the decision is refused |

## Effects

| Effect | Description |
| --- | --- |
| `SGSTS-B01` | Opening a new suggestion writes it to `suggestions/pending/<ID>.md` with its target as title, the gate, origin, target and proposal date, and the reason under "Why". |
| `SGSTS-B02` | A suggestion with a patch carries it in a diff block, followed by the command that applies it. |
| `SGSTS-B03` | A suggestion without a patch has no diff block and says the fix is not mechanical and needs a human decision. |
| `SGSTS-B04` | Opening a suggestion whose ID is already pending creates nothing and answers the existing path. |
| `SGSTS-B05` | A suggestion already approved or rejected is never reopened: opening it again creates nothing, and nothing appears in pending. |
| `SGSTS-B06` | Deciding moves the suggestion from pending to approved or rejected, appending a Decision section with the state, the date, who decided and the reason. (`Decide`) |
| `SGSTS-B07` | A decision taken under automatic judgment is recorded as taken by the AI, distinct from a decision by a person. |
| `SGSTS-B08` | Listing a state gives the IDs of its suggestions, sorted, and nothing (without error) when the state's folder does not exist. |
| `SGSTS-B09` | The patch of a suggestion is extracted clean, without the diff fences, ready for `git apply`. (`PatchOf`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SGSTS-I01` | A suggestion is in exactly one state folder: once decided, it is no longer pending. | opens a suggestion, approves it, and checks approved lists it and pending does not |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SGSTS-E01` | A suggestion without an ID is opened. | Refused with "suggestion without ID"; nothing is written. | Without an ID there is no deduplication, and every scan would add the same proposal again. |
| `SGSTS-E02` | A decision names a state other than approved or rejected. | Refused with "invalid decision state"; the suggestion stays pending. | Moving to any other folder would invent a state the lifecycle does not have. |
| `SGSTS-E03` | A decision has a blank reason. | Refused with "a decision requires a written reason"; the suggestion stays pending. | A choice without its reason is indistinguishable from carelessness for whoever comes after. |
| `SGSTS-E04` | The suggestion to decide is not pending. | Refused with "pending suggestion not found". | Deciding what was never proposed, or was already decided, would fabricate a record. |
| `SGSTS-E05` | The patch is asked of a suggestion that has none. | Refused with "has no patch". | Handing an empty diff to `git apply` would pretend a fix exists. |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
