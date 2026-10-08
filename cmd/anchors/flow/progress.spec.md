<!-- @anchors
  code: PLPRP
  updated_at: 2026-10-08
  layer: comando
-->
# PlanProgress — create a plan's progress file, the state that lives beside the decision and outside the map

> **Code**: `PLPRP`

## Overview

A plan is a decision: what will be done, in what order and why. Changing it must mean the decision
changed, and a gate charges a revision for that. While the phase checkboxes lived inside the plan,
marking one done was changing the plan: the gate charged a revision for finishing a phase, and the
judgment stamp of the plan-to-spec edge, which holds the revision of both ends, fell for the very spec
the phase delivered. Measured: one plan had two commits, its creation and a one-line checkbox change —
all its post-creation changes were progress.

So progress lives in a companion file beside the plan, with the plan's name and a fixed suffix, kept
out of the map so no gate confronts a file whose only job is to change. This unit creates that file:
one section per phase the plan declares in its headers, each with an item to mark. It never
overwrites an existing progress file, because that file holds recorded work.

`anchors new progress --for <plan>` exists for plans born before the companion did — measured: 17
plans, none with a progress file, all 17 with checkboxes inside the plan. The identity of the progress
comes from the plan's own header, so the pair stays locatable by code.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the plan | an existing plan file, relative to the project root or absolute | a missing path; no path at all | this unit: refuses without `--for`, fails reading a missing plan |
| the plan's header | an @anchors header declaring `code:` | a plan with no code | this unit: refuses, since the progress would be born without identity |
| the phase headers | level-two to level-four headers whose code is a project code followed by `-F` and two digits | codes of a length the project does not declare | the project configuration: code lengths are read from it |

## Effects

| Effect | Description |
| --- | --- |
| `PLPRP-B01` | The progress file sits beside the plan, with the plan's name and the `-progress.md` suffix in place of its extension. |
| `PLPRP-B02` | WriteInitialProgress creates the progress file with one `## <phase code> — <title>` section per phase declared in the plan's level-two to level-four headers, each with an unchecked item. |
| `PLPRP-B03` | The length of a phase's code follows the project's configured code lengths. |
| `PLPRP-B04` | A plan with no declared phase still gets its progress file, with a note saying the plan declares no phases yet and what to add when it does. |
| `PLPRP-B05` | An existing progress file is never overwritten; the creation fails and the file keeps its content. |
| `PLPRP-B06` | `new progress --for <plan>` (NewProgressCmd, exposed to the `new` domain rather than registered here) reads the plan relative to the project root, takes the code from the plan's header, creates the progress titled with that code, and says what it created. |
| `PLPRP-B07` | `new progress` without `--for` is refused, with an example of the flag. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `PLPRP-I01` | The suffix this unit writes is the one the scanner keeps out of the map. | a path built with this unit's suffix is recognised as a progress file by the scanner |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `PLPRP-X01` | The progress file's code is never taken from an argument, only from the plan's header. | An argument would open the door to a progress whose identity diverges from its plan's, and the pair would stop being locatable. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `PLPRP-E01` | The plan declares no `code:` in its header. | The command fails naming the plan, and no progress file is created. | A progress without identity cannot be paired with its plan. |
| `PLPRP-E02` | The plan cannot be read. | The command fails with "read the plan". | There is nothing to derive the phases and the code from. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
