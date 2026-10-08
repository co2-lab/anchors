<!-- @anchors
  code: SGCMS
  updated_at: 2026-10-08
  layer: comando
-->
# SuggestCommand — the proposed fixes are listed, shown, applied or rejected, and every decision keeps its record

> **Code**: `SGCMS`

## Overview

The gates and the AI judgment propose fixes; `suggest` is where a person or an agent
decides on them. A suggestion is a git patch plus the reason it was proposed, and its state
is the folder it lives in: pending, approved or rejected. Moving it is deciding. The command
proposes nothing itself: the proposal is automatic, the acceptance is not.

`list` shows the suggestions of one state (pending by default) and, for pending ones, points
at how to see a diff. `show` prints a suggestion's reason and diff. `apply` checks that the
patch still applies BEFORE touching the working tree: a patch proposed against a file that
has changed since fails with an explanation and changes nothing, instead of approving a fix
that never went in. Only after the patch is applied is the suggestion approved, with a
default reason when none was given. `reject` needs a reason and keeps the rejected record,
so the same proposal does not come back as news in the next sweep.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the suggestion id | an id present in the state being read | an id absent from it | this unit (`show`) and the suggestion package (`apply`, `reject`) refuse it |
| the state | pending, approved or rejected | any other folder name | the suggestion package lists what the folder holds |
| the working tree | a git repository where the patch applies | a tree outside git, or a file changed since the proposal | this unit: checks the patch first and explains the failure |
| the rejection reason | a non-empty text | an empty reason | the suggestion package refuses the decision |

## Effects

| Effect | Description |
| --- | --- |
| `SGCMS-B01` | `list` prints the ids of the given state (pending by default) with their count, says so when there is none, and only for pending suggestions points at `suggest show`. |
| `SGCMS-B02` | `show` prints the suggestion's reason and diff; an id absent from the state fails with "not found in <state>". |
| `SGCMS-B03` | `apply --dry-run` only checks that the patch applies and reports it; the file and the state are unchanged. |
| `SGCMS-B04` | `apply` applies the patch, then approves the suggestion; with no reason given, the record carries a default one. |
| `SGCMS-B05` | `reject` without a reason is refused and the suggestion stays pending; with a reason, the suggestion moves to rejected and its record keeps the reason and the automatic marker when `--auto` is given. |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `SGCMS-I01` | A suggestion is approved only if its patch went in: whenever `apply` fails, the file and the state are exactly as before. | applies a patch whose file changed since the proposal and checks both the file and the state |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `SGCMS-X01` | The command proposes no fix; it only lists, shows and decides the ones the gates and the judgment wrote. | The proposal is automatic and the acceptance is not; mixing them would let the tool approve itself. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `SGCMS-E01` | The patch no longer applies because the file changed since the proposal. | `apply` fails saying the patch no longer matches and asking for a new check; nothing is touched. | Applying half a patch, or approving one that did not apply, would leave the record lying about the tree. |
| `SGCMS-E02` | The project is not a git repository. | The failure names the missing repository and the suggestion patch it was trying to apply, not git's raw error. | A suggestion IS a patch; git's own message about a missing repository mentions neither. |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
