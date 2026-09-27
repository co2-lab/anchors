<!-- @anchors
  code: GHIGT
  updated_at: 2026-09-26
  layer: infra
-->
# GitHubIssues — the issue lifecycle on the repository's cards, when the project works on GitHub

> **Code**: `GHIGT`

## Overview

In a project whose workflow runs on GitHub, a gate's finding must become a card on the board, not a
local file: measured in a reference project, eleven gate findings went to local files and none reached
the board, which then showed the planned work and hid what the gates found. This unit is the card
backend of the issue lifecycle (`ISLFS`): an open issue is an open card carrying the workflow label and
the to-do label, a resolved issue is a closed card, and a finding that comes back reopens its card with
the new report as a comment.

Deduplication, which locally comes from looking for the file in the state folders, comes here from
looking for the card by the issue's key, written in the card body as a stable marker. GitHub's search is
textual and returns near misses, so the marker is confirmed exactly: a finding about `Foo.spec.md` must
never match the card of `FooBar.spec.md` and close the wrong card. The search covers closed cards too,
because a closed card is the memory that tells a new finding from one that came back.

Every label the backend applies is one `anchors init` creates, by the same constants: `gh` refuses the
whole command over a missing label, so a divergent label does not degrade, it erases the record.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the repository | `owner/name` from the project's `workflow:` | — | the configuration; `gh` refuses an unknown repository |
| the workflow label | the project's label, created by `anchors init` | — | `anchors init` |
| the `gh` answers | the JSON of `gh issue list` | anything else | this unit: an unreadable answer is an error |

## Effects

| Effect | Description |
| --- | --- |
| `GHIGT-B01` | A card is the issue's only when its body holds the key marker exactly, closed by its comment end; a card whose marker merely starts with the key is not it. |
| `GHIGT-B02` | The card title names the gate, what the finding is (per kind) and the target, without the key. |
| `GHIGT-B03` | The card is looked for among open and closed cards (up to 500) by the key, in the configured repository. |
| `GHIGT-B04` | A finding with no card creates one whose body is the issue's body plus the key marker, labelled with the workflow label and the to-do label, and with the needs-user label when the user owns it. |
| `GHIGT-B05` | An assumed debt's card carries only the workflow label, no flow state. |
| `GHIGT-B06` | A finding whose card is open changes nothing: only the search runs. |
| `GHIGT-B07` | A finding whose card is closed reopens it and comments the new report on it. |
| `GHIGT-B08` | Resolving closes the open card with a comment; a card already closed, or none, is left alone. |
| `GHIGT-B09` | The labels applied are the canonical ones `anchors init` creates, never the legacy Portuguese name. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `GHIGT-E01` | A `gh` call fails (search, create, reopen, comment or close), or the search answer is not readable. | The error is returned carrying `gh`'s output, and nothing is reported created or closed. | A finding reported as recorded that never reached the board is the silent loss this backend exists to end. |

## Dependencies

| Code | File | Method | Layer |
| --- | --- | --- | --- |
| DEP1 | `internal/initx/workflows.go` | `LabelToDo`, `LabelNeedsUser` | infra — the labels `anchors init` creates |
| DEP2 | `internal/issue/issue.go` | `Issue` | infra — the lifecycle this backend serves (`ISLFS`) |

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
