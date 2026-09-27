<!-- @anchors
  code: SCLTE
  updated_at: 2026-09-26
  layer: comando
-->
# Escalate — open the right card for a change the plan, the spec or the tool needs, and stop the work only when it must

> **Code**: `SCLTE`

## Overview

Whoever implements is whoever finds out the plan must change — by incoherence (the text contradicts
itself) or by a gap (the plan is coherent and did not cover something). `anchors escalate <reason>`
opens the card for it, and whoever found it chooses the exit by interpreting the impact:

- the default is an ordinary card: the change does not touch the project's direction, so it is born
  in to-do, enters the queue, and stops nobody;
- `--for-user` asserts the change impacts the direction: it becomes a decision for a person, and the
  card where it was found stops until the decision comes out;
- `--unsure` says the finder does not know whether it impacts: it stops the card like a decision, but
  it is marked as a framing question, because "decide between A and B" and "check whether this is
  yours" are different jobs and mixing them makes the second cost as much as the first;
- `--bug` says the pipeline or the tool is wrong: there is nothing to decide, but the fix lives where no
  agent in the queue edits, so it is not a decision either; with `--blocking` the card waits for the fix.

This is ergonomics that decide the outcome: while opening the right card is more work than fixing in
silence, agents fix in silence. So the command also finds the card the finding was born under — given,
derived from the pull request being reviewed, or the one card the agent holds — and links the new card
to it by label, because a sentence in a body cannot be queried. When the card must stop, it receives the
label that says which card it waits for, and the output prints the way back. Before creating, it warns
when the target file already has open cards, since two agents once delivered the same work on the same
day.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the reason | one or more words; the first line becomes the title | no reason at all | this unit: refuses |
| the exit | none, `--for-user`, `--unsure`, `--bug`, `--bug --blocking` | `--bug` with `--for-user` or `--unsure`; `--blocking` without `--bug` | this unit: refuses before any call |
| the origin card | `--card <n>`; or the pull request under review; or the agent's cards | two cards in the agent's hands | this unit: refuses to guess |
| the target | an optional path (`--about`) | — | the caller |
| the workflow mode | github mode, with at least one workflow label | local mode, where there is no shared queue | this unit refuses local mode; the configuration loader refuses github mode without labels |

## Effects

### Refusals

| Effect | Description |
| --- | --- |
| `SCLTE-B01` | Without a reason the command refuses to run. |
| `SCLTE-B02` | `--bug` with `--for-user` or `--unsure` is refused as "not a decision", and `--blocking` without `--bug` is refused; neither makes any call. |
| `SCLTE-B03` | Outside github mode the command refuses, saying there is no shared queue to take the card from. |

### The origin card

| Effect | Description |
| --- | --- |
| `SCLTE-B04` | Without `--card`, the pull request given with `--reviewing-pr` names the card: the first line of its body that starts with Refs, Closes, Fixes or Resolves followed by `#<n>`, in any case; a number cited mid-sentence does not count. Standard error says which card the pull request declares, or that it declares none. |
| `SCLTE-B05` | Without `--card` and without a card from a pull request, the single card the agent holds is the origin, and standard error says so; with two or more cards in hand the command refuses naming them and creates nothing. |
| `SCLTE-B06` | With no origin card at all the finding is still created, and standard error warns it is born WITHOUT provenance. |

### The new card

| Effect | Description |
| --- | --- |
| `SCLTE-B07` | The title is the exit's prefix followed by the reason's first line, cut to seventy characters ending in "..." when longer: `[plan]` for the default, `[decision]` for `--for-user`, `[framing]` for `--unsure`, `[bug]` for `--bug`. |
| `SCLTE-B08` | Every new card carries the first workflow label and the to-do label; a decision adds needs-user; a framing question adds needs-user and needs-framing; a bug adds the bug label. |
| `SCLTE-B09` | With an origin card, the new card carries `under-<card>`, created on demand; with `--reviewing-pr`, it also carries `from-pr-<pr>`, created on demand — the two coexist. |
| `SCLTE-B10` | The body says why the work stopped or not and how to go on: a decision says it is not the agent's and to remove needs-user after recording a revision; a framing question says the first question is the framing and how to return it to the queue; an ordinary card says it is born under the origin card, delivered in the same pull request, and to escalate for the user if it turns out to change direction; a bug body says there is nothing to decide, where the fix lives, and whether the card waits for it or goes on. |
| `SCLTE-B11` | When `--about` names a target that already has open cards, standard error warns before creating and lists up to three of them, counting the rest. |
| `SCLTE-B12` | The output names the kind and the address: "finding recorded", "decision opened" or "bug reported". |

### The origin card afterwards

| Effect | Description |
| --- | --- |
| `SCLTE-B13` | A decision or a framing question stops the origin card: it receives needs-user and `blocked-by-<new card>` (created on demand), a comment saying it stopped for an open decision, and the output prints the `anchors decided` command that resumes it. |
| `SCLTE-B14` | A blocking bug stops the origin card with `blocked-by-<new card>` alone, comments that a bug blocks it, and says it resumes when the bug closes. |
| `SCLTE-B15` | A non-blocking bug and an ordinary card leave the origin card going on, with a comment tracing the new card. |
| `SCLTE-B16` | The number of the new card is read from the address the platform answers only when its last segment is all digits; otherwise no blocked-by label is made. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SCLTE-I01` | Without an origin card, no card other than the new one is labelled or commented. | an escalation with no origin records no edit and no comment |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SCLTE-X01` | A bug never carries needs-user — neither the new card nor the origin card it blocks. | A bug waits for a fix, not for a person to choose; needs-user would list it under "waiting for you". |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SCLTE-E01` | The platform refuses to create the new card. | The command fails with "open the issue" and the platform's answer. | Nothing was recorded; saying it worked would lose the finding. |
| `SCLTE-E02` | The origin card cannot be labelled as stopped. | A warning says to label it by hand; the command succeeds and does not announce the card as stopped. | The new card exists; failing would hide it, and staying silent would let another agent take the card and redo the path. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/config/config.go` | `Load`, `AbsRoot`, `GitHubMode` | config — project configuration and workflow mode |
| DEP2 | `internal/initx/workflows.go` | `LabelNeedsUser`, `LabelNeedsFraming`, `LabelBug`, `LabelSob`, `LabelDePR`, `LabelBlockedBy` | the workflow labels |
| DEP3 | `cmd/anchors/flow/escalate_dup.go` | `openCardsAbout` | comando — the open cards about the target |
| DEP4 | `cmd/anchors/flow/pr_body.go` | `requestedCards` | comando — the agent's cards |
| DEP5 | `cmd/anchors/common` | `FirstLineOfReason`, `AliasDeFlag`, `ResolveAliases` | comando — the title line and the deprecated flag names |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
