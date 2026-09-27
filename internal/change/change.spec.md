<!-- @anchors
  code: CHRCC
  updated_at: 2026-09-26
  layer: infra
-->
# ChangeRecord — the delivery record an agent leaves when it finishes a stage

> **Code**: `CHRCC`

## Overview

When an agent finishes a stage of work (a spec, the code, a feature, a test), it leaves a delivery
record: which unit it worked on, which files it touched, what it says it did, which decisions it took
that no ruler had decided, and what it knows is still unproven. The reviewer reads that record to know
the scope of the review, and confronts the declared intent against what is on disk; the divergence
between the two is itself a finding.

The record does not carry the change itself (the agent edits the repository directly, and the diff is
git's); it describes and points. Records live in the project's `changes/` folder, as project content
that is versioned and audited, not in the tool's ephemeral state. A record waiting for review sits
directly in `changes/`; once reviewed it moves to `changes/reviewed/`, so the folder encodes the state
and the history is kept.

The identity of a record is its stage plus its unit: a second delivery of the same stage over the same
unit replaces the first, because what matters is the current state of the delivery, not the attempts.

## Domain

| Input | Accepts | Outside the domain | Who guarantees |
| --- | --- | --- | --- |
| the stage | one of the artifact words of the work loop (spec, code, feature, test, plan) | an empty stage | the caller (`anchors deliver`); this unit writes whatever it receives |
| the unit | a path relative to the project root, with `/` or `\` separators | an absolute path | the caller; this unit only normalises the separators and the extension into the key |
| the record path given to review | a record that `Pending` returned | a path that no longer exists | this unit: the move fails and the error is returned |

## Effects

| Effect | Description |
| --- | --- |
| `CHRCC-B01` | The key of a record is its stage, two hyphens, and its unit with the extension dropped and every `/`, `\`, space and `.` turned into a hyphen (leading and trailing hyphens trimmed). |
| `CHRCC-B02` | A record waiting for review lives at `changes/<key>.md` under the project root. |
| `CHRCC-B03` | The rendered record opens with an @anchors header of layer `change` carrying the stage, the unit and the date; the agent line is written only when an agent is named. |
| `CHRCC-B04` | The rendered record states the delivery title, the intent (trimmed) under "What was done", and one line per touched file under "Files". |
| `CHRCC-B05` | The "Decisions the ruler did not decide" and "What is NOT proven" sections are always written: an empty one carries the sentence that says so ("none — …", "nothing — …") instead of being omitted. |
| `CHRCC-B06` | Saving writes the rendered record at its path, creating `changes/` when missing; saving a second delivery with the same stage and unit replaces the first. (`Save`) |
| `CHRCC-B07` | The pending list holds only the `.md` files directly inside `changes/`, sorted, and is empty (without error) when the folder does not exist. |
| `CHRCC-B08` | Marking a record reviewed moves it, under the same file name, to `changes/reviewed/`, creating that folder when missing. (`MarkReviewed`) |

## Invariants

| Rule | Always holds | How it is proven |
| --- | --- | --- |
| `CHRCC-I01` | A reviewed record is never lost: after the move it is no longer pending and it exists in the history. | saves two records, marks one reviewed, and checks it left `changes/`, exists in `changes/reviewed/`, and only the other stays pending |

## Constraints

| Rule | Boundary | Why |
| --- | --- | --- |
| `CHRCC-X01` | The pending list does not descend into subfolders, so the reviewed history (or any folder named like a record) is never listed as pending. | The state is the folder: a record in `changes/reviewed/` would otherwise be reviewed again. |

## Errors

| Code | Condition | Result | Why |
| --- | --- | --- | --- |
| `CHRCC-E01` | `changes/` exists but cannot be read as a folder (for example, it is a file). | Listing the pending records returns the error. | An empty list would say "nothing to review" about a folder that was never read. |
| `CHRCC-E02` | `changes/` cannot be created when a record is saved. | Saving returns the error and no path. | The agent must know its delivery was not recorded; a silent loss leaves the stage unreviewed. |
| `CHRCC-E03` | The record to mark reviewed no longer exists. | The move fails and the error is returned. | Reporting a review as recorded for a record that is not there would fake the history. |

## Dependencies

(none)

## Open Decisions

| Code | Question | Who decides | Becomes |
| --- | --- | --- | --- |

none
